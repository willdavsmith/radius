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

// RedisCachesClient contains the methods for the RedisCaches group.
// Don't use this type directly, use NewRedisCachesClient() instead.
type RedisCachesClient struct {
	internal *arm.Client
	rootScope string
}

// NewRedisCachesClient creates a new instance of RedisCachesClient with the specified values.
//   - rootScope - The scope in which the resource is present. UCP Scope is /planes/{planeType}/{planeName}/resourceGroup/{resourcegroupID}
//     and Azure resource scope is
//     /subscriptions/{subscriptionID}/resourceGroup/{resourcegroupID}
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewRedisCachesClient(rootScope string, credential azcore.TokenCredential, options *arm.ClientOptions) (*RedisCachesClient, error) {
	cl, err := arm.NewClient(moduleName+".RedisCachesClient", moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &RedisCachesClient{
		rootScope: rootScope,
	internal: cl,
	}
	return client, nil
}

// BeginCreateOrUpdate - Create a RedisCacheResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - redisCacheName - The name of the RedisCache portable resource resource
//   - resource - Resource create parameters.
//   - options - RedisCachesClientBeginCreateOrUpdateOptions contains the optional parameters for the RedisCachesClient.BeginCreateOrUpdate
//     method.
func (client *RedisCachesClient) BeginCreateOrUpdate(ctx context.Context, redisCacheName string, resource RedisCacheResource, options *RedisCachesClientBeginCreateOrUpdateOptions) (*runtime.Poller[RedisCachesClientCreateOrUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.createOrUpdate(ctx, redisCacheName, resource, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[RedisCachesClientCreateOrUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaAzureAsyncOp,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[RedisCachesClientCreateOrUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// CreateOrUpdate - Create a RedisCacheResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
func (client *RedisCachesClient) createOrUpdate(ctx context.Context, redisCacheName string, resource RedisCacheResource, options *RedisCachesClientBeginCreateOrUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.createOrUpdateCreateRequest(ctx, redisCacheName, resource, options)
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
func (client *RedisCachesClient) createOrUpdateCreateRequest(ctx context.Context, redisCacheName string, resource RedisCacheResource, options *RedisCachesClientBeginCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Datastores/redisCaches/{redisCacheName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if redisCacheName == "" {
		return nil, errors.New("parameter redisCacheName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{redisCacheName}", url.PathEscape(redisCacheName))
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

// BeginDelete - Delete a RedisCacheResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - redisCacheName - The name of the RedisCache portable resource resource
//   - options - RedisCachesClientBeginDeleteOptions contains the optional parameters for the RedisCachesClient.BeginDelete method.
func (client *RedisCachesClient) BeginDelete(ctx context.Context, redisCacheName string, options *RedisCachesClientBeginDeleteOptions) (*runtime.Poller[RedisCachesClientDeleteResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.deleteOperation(ctx, redisCacheName, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[RedisCachesClientDeleteResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[RedisCachesClientDeleteResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Delete - Delete a RedisCacheResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
func (client *RedisCachesClient) deleteOperation(ctx context.Context, redisCacheName string, options *RedisCachesClientBeginDeleteOptions) (*http.Response, error) {
	var err error
	req, err := client.deleteCreateRequest(ctx, redisCacheName, options)
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
func (client *RedisCachesClient) deleteCreateRequest(ctx context.Context, redisCacheName string, options *RedisCachesClientBeginDeleteOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Datastores/redisCaches/{redisCacheName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if redisCacheName == "" {
		return nil, errors.New("parameter redisCacheName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{redisCacheName}", url.PathEscape(redisCacheName))
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

// Get - Get a RedisCacheResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - redisCacheName - The name of the RedisCache portable resource resource
//   - options - RedisCachesClientGetOptions contains the optional parameters for the RedisCachesClient.Get method.
func (client *RedisCachesClient) Get(ctx context.Context, redisCacheName string, options *RedisCachesClientGetOptions) (RedisCachesClientGetResponse, error) {
	var err error
	req, err := client.getCreateRequest(ctx, redisCacheName, options)
	if err != nil {
		return RedisCachesClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return RedisCachesClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return RedisCachesClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *RedisCachesClient) getCreateRequest(ctx context.Context, redisCacheName string, options *RedisCachesClientGetOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Datastores/redisCaches/{redisCacheName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if redisCacheName == "" {
		return nil, errors.New("parameter redisCacheName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{redisCacheName}", url.PathEscape(redisCacheName))
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
func (client *RedisCachesClient) getHandleResponse(resp *http.Response) (RedisCachesClientGetResponse, error) {
	result := RedisCachesClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RedisCacheResource); err != nil {
		return RedisCachesClientGetResponse{}, err
	}
	return result, nil
}

// NewListByScopePager - List RedisCacheResource resources by Scope
//
// Generated from API version 2023-10-01-preview
//   - options - RedisCachesClientListByScopeOptions contains the optional parameters for the RedisCachesClient.NewListByScopePager
//     method.
func (client *RedisCachesClient) NewListByScopePager(options *RedisCachesClientListByScopeOptions) (*runtime.Pager[RedisCachesClientListByScopeResponse]) {
	return runtime.NewPager(runtime.PagingHandler[RedisCachesClientListByScopeResponse]{
		More: func(page RedisCachesClientListByScopeResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *RedisCachesClientListByScopeResponse) (RedisCachesClientListByScopeResponse, error) {
			var req *policy.Request
			var err error
			if page == nil {
				req, err = client.listByScopeCreateRequest(ctx, options)
			} else {
				req, err = runtime.NewRequest(ctx, http.MethodGet, *page.NextLink)
			}
			if err != nil {
				return RedisCachesClientListByScopeResponse{}, err
			}
			resp, err := client.internal.Pipeline().Do(req)
			if err != nil {
				return RedisCachesClientListByScopeResponse{}, err
			}
			if !runtime.HasStatusCode(resp, http.StatusOK) {
				return RedisCachesClientListByScopeResponse{}, runtime.NewResponseError(resp)
			}
			return client.listByScopeHandleResponse(resp)
		},
	})
}

// listByScopeCreateRequest creates the ListByScope request.
func (client *RedisCachesClient) listByScopeCreateRequest(ctx context.Context, options *RedisCachesClientListByScopeOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Datastores/redisCaches"
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
func (client *RedisCachesClient) listByScopeHandleResponse(resp *http.Response) (RedisCachesClientListByScopeResponse, error) {
	result := RedisCachesClientListByScopeResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RedisCacheResourceListResult); err != nil {
		return RedisCachesClientListByScopeResponse{}, err
	}
	return result, nil
}

// ListSecrets - Lists secrets values for the specified RedisCache resource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - redisCacheName - The name of the RedisCache portable resource resource
//   - body - The content of the action request
//   - options - RedisCachesClientListSecretsOptions contains the optional parameters for the RedisCachesClient.ListSecrets method.
func (client *RedisCachesClient) ListSecrets(ctx context.Context, redisCacheName string, body map[string]any, options *RedisCachesClientListSecretsOptions) (RedisCachesClientListSecretsResponse, error) {
	var err error
	req, err := client.listSecretsCreateRequest(ctx, redisCacheName, body, options)
	if err != nil {
		return RedisCachesClientListSecretsResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return RedisCachesClientListSecretsResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return RedisCachesClientListSecretsResponse{}, err
	}
	resp, err := client.listSecretsHandleResponse(httpResp)
	return resp, err
}

// listSecretsCreateRequest creates the ListSecrets request.
func (client *RedisCachesClient) listSecretsCreateRequest(ctx context.Context, redisCacheName string, body map[string]any, options *RedisCachesClientListSecretsOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Datastores/redisCaches/{redisCacheName}/listSecrets"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if redisCacheName == "" {
		return nil, errors.New("parameter redisCacheName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{redisCacheName}", url.PathEscape(redisCacheName))
	req, err := runtime.NewRequest(ctx, http.MethodPost, runtime.JoinPaths(client.internal.Endpoint(), urlPath))
	if err != nil {
		return nil, err
	}
	reqQP := req.Raw().URL.Query()
	reqQP.Set("api-version", "2023-10-01-preview")
	req.Raw().URL.RawQuery = reqQP.Encode()
	req.Raw().Header["Accept"] = []string{"application/json"}
	if err := runtime.MarshalAsJSON(req, body); err != nil {
	return nil, err
}
	return req, nil
}

// listSecretsHandleResponse handles the ListSecrets response.
func (client *RedisCachesClient) listSecretsHandleResponse(resp *http.Response) (RedisCachesClientListSecretsResponse, error) {
	result := RedisCachesClientListSecretsResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.RedisCacheListSecretsResult); err != nil {
		return RedisCachesClientListSecretsResponse{}, err
	}
	return result, nil
}

// BeginUpdate - Update a RedisCacheResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - redisCacheName - The name of the RedisCache portable resource resource
//   - properties - The resource properties to be updated.
//   - options - RedisCachesClientBeginUpdateOptions contains the optional parameters for the RedisCachesClient.BeginUpdate method.
func (client *RedisCachesClient) BeginUpdate(ctx context.Context, redisCacheName string, properties RedisCacheResourceUpdate, options *RedisCachesClientBeginUpdateOptions) (*runtime.Poller[RedisCachesClientUpdateResponse], error) {
	if options == nil || options.ResumeToken == "" {
		resp, err := client.update(ctx, redisCacheName, properties, options)
		if err != nil {
			return nil, err
		}
		poller, err := runtime.NewPoller(resp, client.internal.Pipeline(), &runtime.NewPollerOptions[RedisCachesClientUpdateResponse]{
			FinalStateVia: runtime.FinalStateViaLocation,
		})
		return poller, err
	} else {
		return runtime.NewPollerFromResumeToken[RedisCachesClientUpdateResponse](options.ResumeToken, client.internal.Pipeline(), nil)
	}
}

// Update - Update a RedisCacheResource
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
func (client *RedisCachesClient) update(ctx context.Context, redisCacheName string, properties RedisCacheResourceUpdate, options *RedisCachesClientBeginUpdateOptions) (*http.Response, error) {
	var err error
	req, err := client.updateCreateRequest(ctx, redisCacheName, properties, options)
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
func (client *RedisCachesClient) updateCreateRequest(ctx context.Context, redisCacheName string, properties RedisCacheResourceUpdate, options *RedisCachesClientBeginUpdateOptions) (*policy.Request, error) {
	urlPath := "/{rootScope}/providers/Applications.Datastores/redisCaches/{redisCacheName}"
	urlPath = strings.ReplaceAll(urlPath, "{rootScope}", client.rootScope)
	if redisCacheName == "" {
		return nil, errors.New("parameter redisCacheName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{redisCacheName}", url.PathEscape(redisCacheName))
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

