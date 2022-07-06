// ------------------------------------------------------------
// Copyright (c) Microsoft Corporation.
// Licensed under the MIT License.
// ------------------------------------------------------------

package deployment

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-openapi/jsonpointer"
	"github.com/project-radius/radius/pkg/armrpc/api/conv"
	"github.com/project-radius/radius/pkg/connectorrp/model"
	"github.com/project-radius/radius/pkg/connectorrp/renderers"
	"github.com/project-radius/radius/pkg/radlogger"
	"github.com/project-radius/radius/pkg/radrp/outputresource"
	"github.com/project-radius/radius/pkg/resourcemodel"
	"github.com/project-radius/radius/pkg/rp"
	"github.com/project-radius/radius/pkg/ucp/resources"
	"github.com/project-radius/radius/pkg/ucp/store"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

//go:generate mockgen -destination=./mock_deploymentprocessor.go -package=deployment -self_package github.com/project-radius/radius/pkg/connectorrp/frontend/deployment github.com/project-radius/radius/pkg/connectorrp/frontend/deployment DeploymentProcessor

type DeploymentProcessor interface {
	Render(ctx context.Context, id resources.ID, resource conv.DataModelInterface) (renderers.RendererOutput, error)
	Deploy(ctx context.Context, id resources.ID, rendererOutput renderers.RendererOutput) (DeploymentOutput, error)
	Delete(ctx context.Context, id resources.ID, outputResources []outputresource.OutputResource) error
	FetchSecrets(ctx context.Context, resource ResourceData) (map[string]interface{}, error)
}

func NewDeploymentProcessor(appmodel model.ApplicationModel, storageClient store.StorageClient, secretClient renderers.SecretValueClient, k8s client.Client) DeploymentProcessor {
	return &deploymentProcessor{appmodel: appmodel, store: storageClient, secretClient: secretClient, k8s: k8s}
}

var _ DeploymentProcessor = (*deploymentProcessor)(nil)

type deploymentProcessor struct {
	appmodel     model.ApplicationModel
	store        store.StorageClient
	secretClient renderers.SecretValueClient
	k8s          client.Client
}

type DeploymentOutput struct {
	Resources      []outputresource.OutputResource
	ComputedValues map[string]interface{}
	SecretValues   map[string]rp.SecretValueReference
}

type ResourceData struct {
	ID              resources.ID
	Resource        conv.DataModelInterface
	OutputResources []outputresource.OutputResource
	ComputedValues  map[string]interface{}
	SecretValues    map[string]rp.SecretValueReference
}

func (dp *deploymentProcessor) Render(ctx context.Context, id resources.ID, resource conv.DataModelInterface) (renderers.RendererOutput, error) {
	logger := radlogger.GetLogger(ctx).WithValues(radlogger.LogFieldResourceID, id.String())
	logger.Info("Rendering resource")

	renderer, err := dp.getResourceRenderer(id)
	if err != nil {
		return renderers.RendererOutput{}, err
	}

	rendererOutput, err := renderer.Render(ctx, resource)
	if err != nil {
		return renderers.RendererOutput{}, err
	}

	// Check if the output resources have the corresponding provider supported in Radius
	for _, or := range rendererOutput.Resources {
		if or.ResourceType.Provider == "" {
			err = fmt.Errorf("output resource %q does not have a provider specified", or.LocalID)
			return renderers.RendererOutput{}, err
		}
		if !dp.appmodel.IsProviderSupported(or.ResourceType.Provider) {
			err := fmt.Errorf("provider %s is not configured. Cannot support resource type %s", or.ResourceType.Provider, or.ResourceType.Type)
			return renderers.RendererOutput{}, err
		}
	}

	return rendererOutput, nil
}

func (dp *deploymentProcessor) getResourceRenderer(id resources.ID) (renderers.Renderer, error) {
	radiusResourceModel, err := dp.appmodel.LookupRadiusResourceModel(id.Type()) // Lookup using resource type
	if err != nil {
		return nil, err
	}

	return radiusResourceModel.Renderer, nil
}

// Deploys rendered output resources in order of dependencies
// returns updated outputresource properties and computed values
func (dp *deploymentProcessor) Deploy(ctx context.Context, id resources.ID, rendererOutput renderers.RendererOutput) (DeploymentOutput, error) {
	logger := radlogger.GetLogger(ctx).WithValues(radlogger.LogFieldResourceID, id.String())
	// Deploy
	logger.Info("Deploying radius resource")

	// Order output resources in deployment dependency order
	orderedOutputResources, err := outputresource.OrderOutputResources(rendererOutput.Resources)
	if err != nil {
		return DeploymentOutput{}, err
	}

	updatedOutputResources := []outputresource.OutputResource{}
	var computedValues map[string]interface{}
	for _, outputResource := range orderedOutputResources {
		deployedComputedValues, err := dp.deployOutputResource(ctx, id, &outputResource, rendererOutput)
		if err != nil {
			return DeploymentOutput{}, err
		}

		updatedOutputResources = append(updatedOutputResources, outputResource)
		computedValues = deployedComputedValues
	}

	// Update static values for connections
	for k, computedValue := range rendererOutput.ComputedValues {
		if computedValue.Value != nil {
			computedValues[k] = computedValue.Value
		}
	}

	return DeploymentOutput{
		Resources:      updatedOutputResources,
		ComputedValues: computedValues,
		SecretValues:   rendererOutput.SecretValues,
	}, nil
}

func (dp *deploymentProcessor) deployOutputResource(ctx context.Context, id resources.ID, outputResource *outputresource.OutputResource, rendererOutput renderers.RendererOutput) (computedValues map[string]interface{}, err error) {
	logger := radlogger.GetLogger(ctx)
	logger.Info(fmt.Sprintf("Deploying output resource: LocalID: %s, resource type: %q\n", outputResource.LocalID, outputResource.ResourceType))

	outputResourceModel, err := dp.appmodel.LookupOutputResourceModel(outputResource.ResourceType)
	if err != nil {
		return nil, err
	}

	resourceIdentity, properties, err := outputResourceModel.ResourceHandler.Put(ctx, outputResource)
	if err != nil {
		return nil, err
	}
	if (resourceIdentity != resourcemodel.ResourceIdentity{}) {
		outputResource.Identity = resourceIdentity
	}

	if outputResource.Identity.ResourceType == nil {
		err = fmt.Errorf("output resource %q does not have an identity. This is a bug in the handler or renderer", outputResource.LocalID)
		return nil, err
	}

	// Values consumed by other Radius resource types through connections
	computedValues = map[string]interface{}{}
	// Copy deployed output resource property values into corresponding expected computed values
	for k, v := range rendererOutput.ComputedValues {
		// A computed value might be a reference to a 'property' returned in preserved properties
		if outputResource.LocalID == v.LocalID && v.PropertyReference != "" {
			computedValues[k] = properties[v.PropertyReference]
			continue
		}

		// A computed value might be a 'pointer' into the deployed resource
		if outputResource.LocalID == v.LocalID && v.JSONPointer != "" {
			pointer, err := jsonpointer.New(v.JSONPointer)
			if err != nil {
				err = fmt.Errorf("failed to process JSON Pointer %q for resource: %w", v.JSONPointer, err)
				return nil, err
			}

			value, _, err := pointer.Get(outputResource.Resource)
			if err != nil {
				err = fmt.Errorf("failed to process JSON Pointer %q for resource: %w", v.JSONPointer, err)
				return nil, err
			}
			computedValues[k] = value
		}
	}

	return computedValues, nil
}

func (dp *deploymentProcessor) Delete(ctx context.Context, resourceID resources.ID, outputResources []outputresource.OutputResource) error {
	logger := radlogger.GetLogger(ctx).WithValues(radlogger.LogFieldResourceID, resourceID)

	orderedOutputResources, err := outputresource.OrderOutputResources(outputResources)
	if err != nil {
		return err
	}

	// Loop over each output resource and delete in reverse dependency order
	for i := len(orderedOutputResources) - 1; i >= 0; i-- {
		outputResource := orderedOutputResources[i]
		outputResourceModel, err := dp.appmodel.LookupOutputResourceModel(outputResource.ResourceType)
		if err != nil {
			return err
		}

		logger.Info(fmt.Sprintf("Deleting output resource: %v, LocalID: %s, resource type: %s\n", outputResource.Identity, outputResource.LocalID, outputResource.ResourceType.Type))
		err = outputResourceModel.ResourceHandler.Delete(ctx, &outputResource)
		if err != nil {
			return err
		}
	}

	return nil
}

func (dp *deploymentProcessor) FetchSecrets(ctx context.Context, resourceData ResourceData) (map[string]interface{}, error) {
	secretValues := map[string]interface{}{}
	for k, secretReference := range resourceData.SecretValues {
		secret, err := dp.fetchSecret(ctx, resourceData.OutputResources, secretReference)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch secret %s for resource %s: %w", k, resourceData.ID.String(), err)
		}

		if (secretReference.Transformer != resourcemodel.ResourceType{}) {
			outputResourceModel, err := dp.appmodel.LookupOutputResourceModel(secretReference.Transformer)
			if err != nil {
				return nil, err
			} else if outputResourceModel.SecretValueTransformer == nil {
				return nil, fmt.Errorf("could not find a secret transformer for %q", secretReference.Transformer)
			}

			secret, err = outputResourceModel.SecretValueTransformer.Transform(ctx, resourceData.Resource, secret)
			if err != nil {
				return nil, fmt.Errorf("failed to transform secret %s for resource %s: %w", k, resourceData.ID.String(), err)
			}
		}

		secretValues[k] = secret
	}

	return secretValues, nil
}

func (dp *deploymentProcessor) fetchSecret(ctx context.Context, outputResources []outputresource.OutputResource, reference rp.SecretValueReference) (interface{}, error) {
	if reference.Value != "" {
		// The secret reference contains the value itself
		return reference.Value, nil
	}

	// Reference to operations to fetch secrets is currently only supported for Azure resources
	if dp.secretClient == nil {
		return nil, errors.New("no Azure credentials provided to fetch secret")
	}

	// Find the output resource that maps to the secret value reference
	for _, outputResource := range outputResources {
		if outputResource.LocalID == reference.LocalID {
			return dp.secretClient.FetchSecret(ctx, outputResource.Identity, reference.Action, reference.ValueSelector)
		}
	}

	return nil, fmt.Errorf("cannot find an output resource matching LocalID %s", reference.LocalID)
}