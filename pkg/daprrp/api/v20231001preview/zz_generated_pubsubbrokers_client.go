//go:build go1.18
// +build go1.18

// Licensed under the Apache License, Version 2.0 . See LICENSE in the repository root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator. DO NOT EDIT.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20231001preview

import (
	"context"
	"errors"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/arm"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"net/http"
	"net/url"
	"strings"
)

// PubSubBrokersClient contains the methods for the PubSubBrokers group.
// Don't use this type directly, use NewPubSubBrokersClient() instead.
type PubSubBrokersClient struct {
	internal *arm.Client
	rootScope string
}

// NewPubSubBrokersClient creates a new instance of PubSubBrokersClient with the specified values.
//   - rootScope - The scope in which the resource is present. UCP Scope is /planes/{planeType}/{planeName}/resourceGroup/{resourcegroupID}
//     and Azure resource scope is
//     /subscriptions/{subscriptionID}/resourceGroup/{resourcegroupID}
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewPubSubBrokersClient(rootScope string, credential azcore.TokenCredential, options *arm.ClientOptions) (*PubSubBrokersClient, error) {
	cl, err := arm.NewClient(moduleName+".PubSubBrokersClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &PubSubBrokersClient{
		rootScope: rootScope,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a DaprPubSubBrokerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - pubSubBrokerName - PubSubBroker name
//   - resource - Resource create parameters.
//   - options - PubSubBrokersClientBeginCreateOrUpdateOptions contains the optional parameters for the PubSubBrokersClient.BeginCreateOrUpdate
//     method.
func (client *PubSubBrokersClient) BeginCreateOrUpdate(ctx context.Context, pubSubBrokerName string, resource DaprPubSubBrokerResource, options *PubSubBrokersClientBeginCreateOrUpdateOptions) (*runtime.Poller[PubSubBrokersClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, pubSubBrokerName, resource, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PubSubBrokersClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[PubSubBrokersClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Create a DaprPubSubBrokerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
func (client *PubSubBrokersClient) createOrUpdate(ctx context.Context, pubSubBrokerName string, resource DaprPubSubBrokerResource, options *PubSubBrokersClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, pubSubBrokerName, resource, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *PubSubBrokersClient) createOrUpdateCreateRequest(ctx context.Context, pubSubBrokerName string, resource DaprPubSubBrokerResource, options *PubSubBrokersClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Dapr/pubSubBrokers/{pubSubBrokerName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if pubSubBrokerName == "" {
		return nil, errors.New("parameter pubSubBrokerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{pubSubBrokerName}", url.PathEscape(pubSubBrokerName))
	req, err := runtime.NewRequest(ctx, http.MethodPut, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, resource); err != nil {
	return nil, err
}
	return req, nil
}

// BeginDelete - Delete a DaprPubSubBrokerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - pubSubBrokerName - PubSubBroker name
//   - options - PubSubBrokersClientBeginDeleteOptions contains the optional parameters for the PubSubBrokersClient.BeginDelete
//     method.
func (client *PubSubBrokersClient) BeginDelete(ctx context.Context, pubSubBrokerName string, options *PubSubBrokersClientBeginDeleteOptions) (*runtime.Poller[PubSubBrokersClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, pubSubBrokerName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PubSubBrokersClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[PubSubBrokersClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Delete a DaprPubSubBrokerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
func (client *PubSubBrokersClient) deleteOperation(ctx context.Context, pubSubBrokerName string, options *PubSubBrokersClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, pubSubBrokerName, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusAccepted, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// deleteCreateRequest creates the Delete request.
func (client *PubSubBrokersClient) deleteCreateRequest(ctx context.Context, pubSubBrokerName string, options *PubSubBrokersClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Dapr/pubSubBrokers/{pubSubBrokerName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if pubSubBrokerName == "" {
		return nil, errors.New("parameter pubSubBrokerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{pubSubBrokerName}", url.PathEscape(pubSubBrokerName))
	req, err := runtime.NewRequest(ctx, http.MethodDelete, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// Get - Get a DaprPubSubBrokerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - pubSubBrokerName - PubSubBroker name
//   - options - PubSubBrokersClientGetOptions contains the optional parameters for the PubSubBrokersClient.Get method.
func (client *PubSubBrokersClient) Get(ctx context.Context, pubSubBrokerName string, options *PubSubBrokersClientGetOptions) (PubSubBrokersClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, pubSubBrokerName, options)
	if err != nil {
		return PubSubBrokersClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return PubSubBrokersClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return PubSubBrokersClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *PubSubBrokersClient) getCreateRequest(ctx context.Context, pubSubBrokerName string, options *PubSubBrokersClientGetOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Dapr/pubSubBrokers/{pubSubBrokerName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if pubSubBrokerName == "" {
		return nil, errors.New("parameter pubSubBrokerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{pubSubBrokerName}", url.PathEscape(pubSubBrokerName))
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// getHandleResponse handles the Get response.
func (client *PubSubBrokersClient) getHandleResponse(resp *http.Response) (PubSubBrokersClientGetResponse, error) {
	result := PubSubBrokersClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DaprPubSubBrokerResource); err != nil {
		return PubSubBrokersClientGetResponse{}, err
	}
	return result, nil
}

// NewListByScopePager - List DaprPubSubBrokerResource resources by Scope
//
// Generated from API version 2023-10-01-preview
//   - options - PubSubBrokersClientListByScopeOptions contains the optional parameters for the PubSubBrokersClient.NewListByScopePager
//     method.
func (client *PubSubBrokersClient) NewListByScopePager(options *PubSubBrokersClientListByScopeOptions) (*runtime.Pager[PubSubBrokersClientListByScopeResponse]) {
	return runtime.NewPager(runtime.PagingHandler[PubSubBrokersClientListByScopeResponse]{
		More: func(page PubSubBrokersClientListByScopeResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *PubSubBrokersClientListByScopeResponse) (PubSubBrokersClientListByScopeResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByScopeCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return PubSubBrokersClientListByScopeResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return PubSubBrokersClientListByScopeResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return PubSubBrokersClientListByScopeResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByScopeHandleResponse(resp)
		},
	})
}

// listByScopeCreateRequest creates the ListByScope request.
func (client *PubSubBrokersClient) listByScopeCreateRequest(ctx context.Context, options *PubSubBrokersClientListByScopeOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Dapr/pubSubBrokers"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	req, err := runtime.NewRequest(ctx, http.MethodGet, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	return req, nil
}

// listByScopeHandleResponse handles the ListByScope response.
func (client *PubSubBrokersClient) listByScopeHandleResponse(resp *http.Response) (PubSubBrokersClientListByScopeResponse, error) {
	result := PubSubBrokersClientListByScopeResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.DaprPubSubBrokerResourceListResult); err != nil {
		return PubSubBrokersClientListByScopeResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update a DaprPubSubBrokerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - pubSubBrokerName - PubSubBroker name
//   - properties - The resource properties to be updated.
//   - options - PubSubBrokersClientBeginUpdateOptions contains the optional parameters for the PubSubBrokersClient.BeginUpdate
//     method.
func (client *PubSubBrokersClient) BeginUpdate(ctx context.Context, pubSubBrokerName string, properties DaprPubSubBrokerResourceUpdate, options *PubSubBrokersClientBeginUpdateOptions) (*runtime.Poller[PubSubBrokersClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, pubSubBrokerName, properties, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[PubSubBrokersClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[PubSubBrokersClientUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Update - Update a DaprPubSubBrokerResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
func (client *PubSubBrokersClient) update(ctx context.Context, pubSubBrokerName string, properties DaprPubSubBrokerResourceUpdate, options *PubSubBrokersClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, pubSubBrokerName, properties, options)
	if err != nil {
		return nil, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return nil, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusAccepted) {
		err = runtime.NewResponseError(httpResp)
		return nil, err
	}
	return httpResp, nil
}

// updateCreateRequest creates the Update request.
func (client *PubSubBrokersClient) updateCreateRequest(ctx context.Context, pubSubBrokerName string, properties DaprPubSubBrokerResourceUpdate, options *PubSubBrokersClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Dapr/pubSubBrokers/{pubSubBrokerName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if pubSubBrokerName == "" {
		return nil, errors.New("parameter pubSubBrokerName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{pubSubBrokerName}", url.PathEscape(pubSubBrokerName))
	req, err := runtime.NewRequest(ctx, http.MethodPatch, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, properties); err != nil {
	return nil, err
}
	return req, nil
}

