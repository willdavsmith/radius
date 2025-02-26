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

// AwsCredentialsClient contains the methods for the AwsCredentials group.
// Don't use this type directly, use NewAwsCredentialsClient() instead.
type AwsCredentialsClient struct {
	internal *arm.Client
}

// NewAwsCredentialsClient creates a new instance of AwsCredentialsClient with the specified values.
//   - credential - used to authorize requests. Usually a credential from azidentity.
//   - options - pass nil to accept the default values.
func NewAwsCredentialsClient(credential azcore.TokenCredential, options *arm.ClientOptions) (*AwsCredentialsClient, error) {
	cl, err := arm.NewClient(moduleName, moduleVersion, credential, options)
	if err != nil {
		return nil, err
	}
	client := &AwsCredentialsClient{
	internal: cl,
	}
	return client, nil
}

// CreateOrUpdate - Create or update an AWS credential
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - planeName - The name of AWS plane
//   - credentialName - The AWS credential name.
//   - resource - Resource create parameters.
//   - options - AwsCredentialsClientCreateOrUpdateOptions contains the optional parameters for the AwsCredentialsClient.CreateOrUpdate
//     method.
func (client *AwsCredentialsClient) CreateOrUpdate(ctx context.Context, planeName string, credentialName string, resource AwsCredentialResource, options *AwsCredentialsClientCreateOrUpdateOptions) (AwsCredentialsClientCreateOrUpdateResponse, error) {
	var err error
	const operationName = "AwsCredentialsClient.CreateOrUpdate"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.createOrUpdateCreateRequest(ctx, planeName, credentialName, resource, options)
	if err != nil {
		return AwsCredentialsClientCreateOrUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AwsCredentialsClientCreateOrUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusCreated) {
		err = runtime.NewResponseError(httpResp)
		return AwsCredentialsClientCreateOrUpdateResponse{}, err
	}
	resp, err := client.createOrUpdateHandleResponse(httpResp)
	return resp, err
}

// createOrUpdateCreateRequest creates the CreateOrUpdate request.
func (client *AwsCredentialsClient) createOrUpdateCreateRequest(ctx context.Context, planeName string, credentialName string, resource AwsCredentialResource, _ *AwsCredentialsClientCreateOrUpdateOptions) (*policy.Request, error) {
	urlPath := "/planes/aws/{planeName}/providers/System.AWS/credentials/{credentialName}"
	urlPath = strings.ReplaceAll(urlPath, "{planeName}", planeName)
	if credentialName == "" {
		return nil, errors.New("parameter credentialName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{credentialName}", url.PathEscape(credentialName))
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
;	return req, nil
}

// createOrUpdateHandleResponse handles the CreateOrUpdate response.
func (client *AwsCredentialsClient) createOrUpdateHandleResponse(resp *http.Response) (AwsCredentialsClientCreateOrUpdateResponse, error) {
	result := AwsCredentialsClientCreateOrUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AwsCredentialResource); err != nil {
		return AwsCredentialsClientCreateOrUpdateResponse{}, err
	}
	return result, nil
}

// Delete - Delete an AWS credential
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - planeName - The name of AWS plane
//   - credentialName - The AWS credential name.
//   - options - AwsCredentialsClientDeleteOptions contains the optional parameters for the AwsCredentialsClient.Delete method.
func (client *AwsCredentialsClient) Delete(ctx context.Context, planeName string, credentialName string, options *AwsCredentialsClientDeleteOptions) (AwsCredentialsClientDeleteResponse, error) {
	var err error
	const operationName = "AwsCredentialsClient.Delete"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.deleteCreateRequest(ctx, planeName, credentialName, options)
	if err != nil {
		return AwsCredentialsClientDeleteResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AwsCredentialsClientDeleteResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK, http.StatusNoContent) {
		err = runtime.NewResponseError(httpResp)
		return AwsCredentialsClientDeleteResponse{}, err
	}
	return AwsCredentialsClientDeleteResponse{}, nil
}

// deleteCreateRequest creates the Delete request.
func (client *AwsCredentialsClient) deleteCreateRequest(ctx context.Context, planeName string, credentialName string, _ *AwsCredentialsClientDeleteOptions) (*policy.Request, error) {
	urlPath := "/planes/aws/{planeName}/providers/System.AWS/credentials/{credentialName}"
	urlPath = strings.ReplaceAll(urlPath, "{planeName}", planeName)
	if credentialName == "" {
		return nil, errors.New("parameter credentialName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{credentialName}", url.PathEscape(credentialName))
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

// Get - Get an AWS credential
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - planeName - The name of AWS plane
//   - credentialName - The AWS credential name.
//   - options - AwsCredentialsClientGetOptions contains the optional parameters for the AwsCredentialsClient.Get method.
func (client *AwsCredentialsClient) Get(ctx context.Context, planeName string, credentialName string, options *AwsCredentialsClientGetOptions) (AwsCredentialsClientGetResponse, error) {
	var err error
	const operationName = "AwsCredentialsClient.Get"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.getCreateRequest(ctx, planeName, credentialName, options)
	if err != nil {
		return AwsCredentialsClientGetResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AwsCredentialsClientGetResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AwsCredentialsClientGetResponse{}, err
	}
	resp, err := client.getHandleResponse(httpResp)
	return resp, err
}

// getCreateRequest creates the Get request.
func (client *AwsCredentialsClient) getCreateRequest(ctx context.Context, planeName string, credentialName string, _ *AwsCredentialsClientGetOptions) (*policy.Request, error) {
	urlPath := "/planes/aws/{planeName}/providers/System.AWS/credentials/{credentialName}"
	urlPath = strings.ReplaceAll(urlPath, "{planeName}", planeName)
	if credentialName == "" {
		return nil, errors.New("parameter credentialName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{credentialName}", url.PathEscape(credentialName))
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
func (client *AwsCredentialsClient) getHandleResponse(resp *http.Response) (AwsCredentialsClientGetResponse, error) {
	result := AwsCredentialsClientGetResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AwsCredentialResource); err != nil {
		return AwsCredentialsClientGetResponse{}, err
	}
	return result, nil
}

// NewListPager - List AWS credentials
//
// Generated from API version 2023-10-01-preview
//   - planeName - The name of AWS plane
//   - options - AwsCredentialsClientListOptions contains the optional parameters for the AwsCredentialsClient.NewListPager method.
func (client *AwsCredentialsClient) NewListPager(planeName string, options *AwsCredentialsClientListOptions) (*runtime.Pager[AwsCredentialsClientListResponse]) {
	return runtime.NewPager(runtime.PagingHandler[AwsCredentialsClientListResponse]{
		More: func(page AwsCredentialsClientListResponse) bool {
			return page.NextLink != nil && len(*page.NextLink) > 0
		},
		Fetcher: func(ctx context.Context, page *AwsCredentialsClientListResponse) (AwsCredentialsClientListResponse, error) {
		ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, "AwsCredentialsClient.NewListPager")
			nextLink := ""
			if page != nil {
				nextLink = *page.NextLink
			}
			resp, err := runtime.FetcherForNextLink(ctx, client.internal.Pipeline(), nextLink, func(ctx context.Context) (*policy.Request, error) {
				return client.listCreateRequest(ctx, planeName, options)
			}, nil)
			if err != nil {
				return AwsCredentialsClientListResponse{}, err
			}
			return client.listHandleResponse(resp)
			},
		Tracer: client.internal.Tracer(),
	})
}

// listCreateRequest creates the List request.
func (client *AwsCredentialsClient) listCreateRequest(ctx context.Context, planeName string, _ *AwsCredentialsClientListOptions) (*policy.Request, error) {
	urlPath := "/planes/aws/{planeName}/providers/System.AWS/credentials"
	urlPath = strings.ReplaceAll(urlPath, "{planeName}", planeName)
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

// listHandleResponse handles the List response.
func (client *AwsCredentialsClient) listHandleResponse(resp *http.Response) (AwsCredentialsClientListResponse, error) {
	result := AwsCredentialsClientListResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AwsCredentialResourceListResult); err != nil {
		return AwsCredentialsClientListResponse{}, err
	}
	return result, nil
}

// Update - Update an AWS credential
// If the operation fails it returns an *azcore.ResponseError type.
//
// Generated from API version 2023-10-01-preview
//   - planeName - The name of AWS plane
//   - credentialName - The AWS credential name.
//   - properties - The resource properties to be updated.
//   - options - AwsCredentialsClientUpdateOptions contains the optional parameters for the AwsCredentialsClient.Update method.
func (client *AwsCredentialsClient) Update(ctx context.Context, planeName string, credentialName string, properties AwsCredentialResourceTagsUpdate, options *AwsCredentialsClientUpdateOptions) (AwsCredentialsClientUpdateResponse, error) {
	var err error
	const operationName = "AwsCredentialsClient.Update"
	ctx = context.WithValue(ctx, runtime.CtxAPINameKey{}, operationName)
	ctx, endSpan := runtime.StartSpan(ctx, operationName, client.internal.Tracer(), nil)
	defer func() { endSpan(err) }()
	req, err := client.updateCreateRequest(ctx, planeName, credentialName, properties, options)
	if err != nil {
		return AwsCredentialsClientUpdateResponse{}, err
	}
	httpResp, err := client.internal.Pipeline().Do(req)
	if err != nil {
		return AwsCredentialsClientUpdateResponse{}, err
	}
	if !runtime.HasStatusCode(httpResp, http.StatusOK) {
		err = runtime.NewResponseError(httpResp)
		return AwsCredentialsClientUpdateResponse{}, err
	}
	resp, err := client.updateHandleResponse(httpResp)
	return resp, err
}

// updateCreateRequest creates the Update request.
func (client *AwsCredentialsClient) updateCreateRequest(ctx context.Context, planeName string, credentialName string, properties AwsCredentialResourceTagsUpdate, _ *AwsCredentialsClientUpdateOptions) (*policy.Request, error) {
	urlPath := "/planes/aws/{planeName}/providers/System.AWS/credentials/{credentialName}"
	urlPath = strings.ReplaceAll(urlPath, "{planeName}", planeName)
	if credentialName == "" {
		return nil, errors.New("parameter credentialName cannot be empty")
	}
	urlPath = strings.ReplaceAll(urlPath, "{credentialName}", url.PathEscape(credentialName))
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
;	return req, nil
}

// updateHandleResponse handles the Update response.
func (client *AwsCredentialsClient) updateHandleResponse(resp *http.Response) (AwsCredentialsClientUpdateResponse, error) {
	result := AwsCredentialsClientUpdateResponse{}
	if err := runtime.UnmarshalAsJSON(resp, &result.AwsCredentialResource); err != nil {
		return AwsCredentialsClientUpdateResponse{}, err
	}
	return result, nil
}

