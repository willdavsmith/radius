//go:build go1.16
// +build go1.16

// Copyright (c) Microsoft Corporation. All rights reserved.
// Licensed under the MIT License. See License.txt in the project root for license information.
// Code generated by Microsoft (R) AutoRest Code Generator.
// Changes may cause incorrect behavior and will be lost if the code is regenerated.

package v20220315privatepreview

import (
	"encoding/json"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"reflect"
	"time"
)

// ApplicationProperties - Application properties
type ApplicationProperties struct {
	// REQUIRED; The resource id of the environment linked to application.
	Environment *string `json:"environment,omitempty"`

	// READ-ONLY; Provisioning state of the application at the time the operation was called.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// ApplicationResource - Radius Application.
type ApplicationResource struct {
	TrackedResource
	// REQUIRED; Application properties
	Properties *ApplicationProperties `json:"properties,omitempty"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type ApplicationResource.
func (a ApplicationResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	a.TrackedResource.marshalInternal(objectMap)
	populate(objectMap, "properties", a.Properties)
	populate(objectMap, "systemData", a.SystemData)
	return json.Marshal(objectMap)
}

// ApplicationResourceList - The list of applications.
type ApplicationResourceList struct {
	// The link used to get the next page of applications list.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of applications.
	Value []*ApplicationResource `json:"value,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type ApplicationResourceList.
func (a ApplicationResourceList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", a.NextLink)
	populate(objectMap, "value", a.Value)
	return json.Marshal(objectMap)
}

// ApplicationsCreateOrUpdateOptions contains the optional parameters for the Applications.CreateOrUpdate method.
type ApplicationsCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// ApplicationsDeleteOptions contains the optional parameters for the Applications.Delete method.
type ApplicationsDeleteOptions struct {
	// placeholder for future optional parameters
}

// ApplicationsGetOptions contains the optional parameters for the Applications.Get method.
type ApplicationsGetOptions struct {
	// placeholder for future optional parameters
}

// ApplicationsListBySubscriptionOptions contains the optional parameters for the Applications.ListBySubscription method.
type ApplicationsListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// ApplicationsListOptions contains the optional parameters for the Applications.List method.
type ApplicationsListOptions struct {
	// placeholder for future optional parameters
}

// ApplicationsUpdateOptions contains the optional parameters for the Applications.Update method.
type ApplicationsUpdateOptions struct {
	// placeholder for future optional parameters
}

// BasicResourceProperties - Basic properties of a Radius resource.
type BasicResourceProperties struct {
	// Status of the resource
	Status *ResourceStatus `json:"status,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type BasicResourceProperties.
func (b BasicResourceProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	b.marshalInternal(objectMap)
	return json.Marshal(objectMap)
}

func (b BasicResourceProperties) marshalInternal(objectMap map[string]interface{}) {
	populate(objectMap, "status", b.Status)
}

// EnvironmentCompute - Compute resource used by application environment resource.
type EnvironmentCompute struct {
	// REQUIRED; Type of compute resource.
	Kind *EnvironmentComputeKind `json:"kind,omitempty"`

	// REQUIRED; The resource id of the compute resource for application environment.
	ResourceID *string `json:"resourceId,omitempty"`
}

// EnvironmentProperties - Application environment properties
type EnvironmentProperties struct {
	// REQUIRED; The compute resource used by application environment.
	Compute *EnvironmentCompute `json:"compute,omitempty"`

	// READ-ONLY; Provisioning state of the environment at the time the operation was called.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// EnvironmentResource - Application environment.
type EnvironmentResource struct {
	TrackedResource
	// REQUIRED; Application environment properties
	Properties *EnvironmentProperties `json:"properties,omitempty"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type EnvironmentResource.
func (e EnvironmentResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	e.TrackedResource.marshalInternal(objectMap)
	populate(objectMap, "properties", e.Properties)
	populate(objectMap, "systemData", e.SystemData)
	return json.Marshal(objectMap)
}

// EnvironmentResourceList - The list of environments.
type EnvironmentResourceList struct {
	// The link used to get the next page of environments list.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of environments.
	Value []*EnvironmentResource `json:"value,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type EnvironmentResourceList.
func (e EnvironmentResourceList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", e.NextLink)
	populate(objectMap, "value", e.Value)
	return json.Marshal(objectMap)
}

// EnvironmentsCreateOrUpdateOptions contains the optional parameters for the Environments.CreateOrUpdate method.
type EnvironmentsCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// EnvironmentsDeleteOptions contains the optional parameters for the Environments.Delete method.
type EnvironmentsDeleteOptions struct {
	// placeholder for future optional parameters
}

// EnvironmentsGetOptions contains the optional parameters for the Environments.Get method.
type EnvironmentsGetOptions struct {
	// placeholder for future optional parameters
}

// EnvironmentsListBySubscriptionOptions contains the optional parameters for the Environments.ListBySubscription method.
type EnvironmentsListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// EnvironmentsListOptions contains the optional parameters for the Environments.List method.
type EnvironmentsListOptions struct {
	// placeholder for future optional parameters
}

// EnvironmentsUpdateOptions contains the optional parameters for the Environments.Update method.
type EnvironmentsUpdateOptions struct {
	// placeholder for future optional parameters
}

// ErrorAdditionalInfo - The resource management error additional info.
type ErrorAdditionalInfo struct {
	// READ-ONLY; The additional info.
	Info map[string]interface{} `json:"info,omitempty" azure:"ro"`

	// READ-ONLY; The additional info type.
	Type *string `json:"type,omitempty" azure:"ro"`
}

// ErrorDetail - The error detail.
type ErrorDetail struct {
	// READ-ONLY; The error additional info.
	AdditionalInfo []*ErrorAdditionalInfo `json:"additionalInfo,omitempty" azure:"ro"`

	// READ-ONLY; The error code.
	Code *string `json:"code,omitempty" azure:"ro"`

	// READ-ONLY; The error details.
	Details []*ErrorDetail `json:"details,omitempty" azure:"ro"`

	// READ-ONLY; The error message.
	Message *string `json:"message,omitempty" azure:"ro"`

	// READ-ONLY; The error target.
	Target *string `json:"target,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type ErrorDetail.
func (e ErrorDetail) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "additionalInfo", e.AdditionalInfo)
	populate(objectMap, "code", e.Code)
	populate(objectMap, "details", e.Details)
	populate(objectMap, "message", e.Message)
	populate(objectMap, "target", e.Target)
	return json.Marshal(objectMap)
}

// ErrorResponse - Common error response for all Azure Resource Manager APIs to return error details for failed operations. (This also follows the OData
// error response format.).
// Implements the error and azcore.HTTPResponse interfaces.
type ErrorResponse struct {
	raw string
	// The error object.
	InnerError *ErrorDetail `json:"error,omitempty"`
}

// Error implements the error interface for type ErrorResponse.
// The contents of the error text are not contractual and subject to change.
func (e ErrorResponse) Error() string {
	return e.raw
}

// GatewayProperties - Gateway properties
type GatewayProperties struct {
	BasicResourceProperties
	// REQUIRED; The resource id of the application linked to Gateway resource.
	Application *string `json:"application,omitempty"`

	// Declare hostname information for the Gateway. Leaving the hostname empty auto-assigns one: mygateway.myapp.PUBLICHOSTNAMEORIP.nip.io.
	Hostname *GatewayPropertiesHostname `json:"hostname,omitempty"`

	// Sets Gateway to not be exposed externally (no public IP address associated). Defaults to false (exposed to internet).
	Internal *bool `json:"internal,omitempty"`

	// Routes attached to this Gateway
	Routes []*GatewayRoute `json:"routes,omitempty"`

	// READ-ONLY; Provisioning state of the Gateway at the time the operation was called.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type GatewayProperties.
func (g GatewayProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	g.BasicResourceProperties.marshalInternal(objectMap)
	populate(objectMap, "application", g.Application)
	populate(objectMap, "hostname", g.Hostname)
	populate(objectMap, "internal", g.Internal)
	populate(objectMap, "provisioningState", g.ProvisioningState)
	populate(objectMap, "routes", g.Routes)
	return json.Marshal(objectMap)
}

// GatewayPropertiesHostname - Declare hostname information for the Gateway. Leaving the hostname empty auto-assigns one: mygateway.myapp.PUBLICHOSTNAMEORIP.nip.io.
type GatewayPropertiesHostname struct {
	// Specify a fully-qualified domain name: myapp.mydomain.com. Mutually exclusive with 'prefix' and will take priority if both are defined.
	FullyQualifiedHostname *string `json:"fullyQualifiedHostname,omitempty"`

	// Specify a prefix for the hostname: myhostname.myapp.PUBLICHOSTNAMEORIP.nip.io. Mutually exclusive with 'fullyQualifiedHostname' and will be overridden
// if both are defined.
	Prefix *string `json:"prefix,omitempty"`
}

// GatewayResource - Gateway Resource that specifies how traffic is exposed to the application.
type GatewayResource struct {
	TrackedResource
	// REQUIRED; Gateway properties
	Properties *GatewayProperties `json:"properties,omitempty"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type GatewayResource.
func (g GatewayResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	g.TrackedResource.marshalInternal(objectMap)
	populate(objectMap, "properties", g.Properties)
	populate(objectMap, "systemData", g.SystemData)
	return json.Marshal(objectMap)
}

// GatewayResourceList - The list of Gateways.
type GatewayResourceList struct {
	// The link used to get the next page of Gateways list.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of Gateways.
	Value []*GatewayResource `json:"value,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type GatewayResourceList.
func (g GatewayResourceList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", g.NextLink)
	populate(objectMap, "value", g.Value)
	return json.Marshal(objectMap)
}

type GatewayRoute struct {
	// The HttpRoute to route to. Ex - myserviceroute.id.
	Destination *string `json:"destination,omitempty"`

	// The path to match the incoming request path on. Ex - /myservice.
	Path *string `json:"path,omitempty"`

	// Optionally update the prefix when sending the request to the service. Ex - replacePrefix: '/' and path: '/myservice' will transform '/myservice/myroute'
// to '/myroute'
	ReplacePrefix *string `json:"replacePrefix,omitempty"`
}

// GatewaysCreateOrUpdateOptions contains the optional parameters for the Gateways.CreateOrUpdate method.
type GatewaysCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// GatewaysDeleteOptions contains the optional parameters for the Gateways.Delete method.
type GatewaysDeleteOptions struct {
	// placeholder for future optional parameters
}

// GatewaysGetOptions contains the optional parameters for the Gateways.Get method.
type GatewaysGetOptions struct {
	// placeholder for future optional parameters
}

// GatewaysListBySubscriptionOptions contains the optional parameters for the Gateways.ListBySubscription method.
type GatewaysListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// GatewaysListOptions contains the optional parameters for the Gateways.List method.
type GatewaysListOptions struct {
	// placeholder for future optional parameters
}

// GatewaysUpdateOptions contains the optional parameters for the Gateways.Update method.
type GatewaysUpdateOptions struct {
	// placeholder for future optional parameters
}

// HTTPRouteProperties - HTTP Route properties
type HTTPRouteProperties struct {
	BasicResourceProperties
	// REQUIRED; The resource id of the application linked to HTTP Route resource.
	Application *string `json:"application,omitempty"`

	// The internal hostname accepting traffic for the HTTP Route. Readonly.
	Hostname *string `json:"hostname,omitempty"`

	// The port number for the HTTP Route. Defaults to 80. Readonly.
	Port *int32 `json:"port,omitempty"`

	// The scheme used for traffic. Readonly.
	Scheme *string `json:"scheme,omitempty"`

	// A stable URL that that can be used to route traffic to a resource. Readonly.
	URL *string `json:"url,omitempty"`

	// READ-ONLY; Provisioning state of the HTTP Route at the time the operation was called.
	ProvisioningState *ProvisioningState `json:"provisioningState,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type HTTPRouteProperties.
func (h HTTPRouteProperties) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	h.BasicResourceProperties.marshalInternal(objectMap)
	populate(objectMap, "application", h.Application)
	populate(objectMap, "hostname", h.Hostname)
	populate(objectMap, "port", h.Port)
	populate(objectMap, "provisioningState", h.ProvisioningState)
	populate(objectMap, "scheme", h.Scheme)
	populate(objectMap, "url", h.URL)
	return json.Marshal(objectMap)
}

// HTTPRouteResource - Radius HTTP Route Resource.
type HTTPRouteResource struct {
	TrackedResource
	// REQUIRED; HTTP Route properties
	Properties *HTTPRouteProperties `json:"properties,omitempty"`

	// READ-ONLY; Metadata pertaining to creation and last modification of the resource.
	SystemData *SystemData `json:"systemData,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type HTTPRouteResource.
func (h HTTPRouteResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	h.TrackedResource.marshalInternal(objectMap)
	populate(objectMap, "properties", h.Properties)
	populate(objectMap, "systemData", h.SystemData)
	return json.Marshal(objectMap)
}

// HTTPRouteResourceList - The list of HTTP Routes.
type HTTPRouteResourceList struct {
	// The link used to get the next page of HTTP Routes list.
	NextLink *string `json:"nextLink,omitempty"`

	// The list of HTTP Route.
	Value []*HTTPRouteResource `json:"value,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type HTTPRouteResourceList.
func (h HTTPRouteResourceList) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "nextLink", h.NextLink)
	populate(objectMap, "value", h.Value)
	return json.Marshal(objectMap)
}

// HTTPRoutesCreateOrUpdateOptions contains the optional parameters for the HTTPRoutes.CreateOrUpdate method.
type HTTPRoutesCreateOrUpdateOptions struct {
	// placeholder for future optional parameters
}

// HTTPRoutesDeleteOptions contains the optional parameters for the HTTPRoutes.Delete method.
type HTTPRoutesDeleteOptions struct {
	// placeholder for future optional parameters
}

// HTTPRoutesGetOptions contains the optional parameters for the HTTPRoutes.Get method.
type HTTPRoutesGetOptions struct {
	// placeholder for future optional parameters
}

// HTTPRoutesListBySubscriptionOptions contains the optional parameters for the HTTPRoutes.ListBySubscription method.
type HTTPRoutesListBySubscriptionOptions struct {
	// placeholder for future optional parameters
}

// HTTPRoutesListOptions contains the optional parameters for the HTTPRoutes.List method.
type HTTPRoutesListOptions struct {
	// placeholder for future optional parameters
}

// HTTPRoutesUpdateOptions contains the optional parameters for the HTTPRoutes.Update method.
type HTTPRoutesUpdateOptions struct {
	// placeholder for future optional parameters
}

// Resource - Common fields that are returned in the response for all Azure Resource Manager resources
type Resource struct {
	// READ-ONLY; Fully qualified resource ID for the resource. Ex - /subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/{resourceProviderNamespace}/{resourceType}/{resourceName}
	ID *string `json:"id,omitempty" azure:"ro"`

	// READ-ONLY; The name of the resource
	Name *string `json:"name,omitempty" azure:"ro"`

	// READ-ONLY; The type of the resource. E.g. "Microsoft.Compute/virtualMachines" or "Microsoft.Storage/storageAccounts"
	Type *string `json:"type,omitempty" azure:"ro"`
}

// MarshalJSON implements the json.Marshaller interface for type Resource.
func (r Resource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	r.marshalInternal(objectMap)
	return json.Marshal(objectMap)
}

func (r Resource) marshalInternal(objectMap map[string]interface{}) {
	populate(objectMap, "id", r.ID)
	populate(objectMap, "name", r.Name)
	populate(objectMap, "type", r.Type)
}

// ResourceStatus - Status of a resource.
type ResourceStatus struct {
	OutputResources []map[string]interface{} `json:"outputResources,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type ResourceStatus.
func (r ResourceStatus) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "outputResources", r.OutputResources)
	return json.Marshal(objectMap)
}

// SystemData - Metadata pertaining to creation and last modification of the resource.
type SystemData struct {
	// The timestamp of resource creation (UTC).
	CreatedAt *time.Time `json:"createdAt,omitempty"`

	// The identity that created the resource.
	CreatedBy *string `json:"createdBy,omitempty"`

	// The type of identity that created the resource.
	CreatedByType *CreatedByType `json:"createdByType,omitempty"`

	// The timestamp of resource last modification (UTC)
	LastModifiedAt *time.Time `json:"lastModifiedAt,omitempty"`

	// The identity that last modified the resource.
	LastModifiedBy *string `json:"lastModifiedBy,omitempty"`

	// The type of identity that last modified the resource.
	LastModifiedByType *CreatedByType `json:"lastModifiedByType,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type SystemData.
func (s SystemData) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	populate(objectMap, "createdAt", (*timeRFC3339)(s.CreatedAt))
	populate(objectMap, "createdBy", s.CreatedBy)
	populate(objectMap, "createdByType", s.CreatedByType)
	populate(objectMap, "lastModifiedAt", (*timeRFC3339)(s.LastModifiedAt))
	populate(objectMap, "lastModifiedBy", s.LastModifiedBy)
	populate(objectMap, "lastModifiedByType", s.LastModifiedByType)
	return json.Marshal(objectMap)
}

// UnmarshalJSON implements the json.Unmarshaller interface for type SystemData.
func (s *SystemData) UnmarshalJSON(data []byte) error {
	var rawMsg map[string]json.RawMessage
	if err := json.Unmarshal(data, &rawMsg); err != nil {
		return err
	}
	for key, val := range rawMsg {
		var err error
		switch key {
		case "createdAt":
				var aux timeRFC3339
				err = unpopulate(val, &aux)
				s.CreatedAt = (*time.Time)(&aux)
				delete(rawMsg, key)
		case "createdBy":
				err = unpopulate(val, &s.CreatedBy)
				delete(rawMsg, key)
		case "createdByType":
				err = unpopulate(val, &s.CreatedByType)
				delete(rawMsg, key)
		case "lastModifiedAt":
				var aux timeRFC3339
				err = unpopulate(val, &aux)
				s.LastModifiedAt = (*time.Time)(&aux)
				delete(rawMsg, key)
		case "lastModifiedBy":
				err = unpopulate(val, &s.LastModifiedBy)
				delete(rawMsg, key)
		case "lastModifiedByType":
				err = unpopulate(val, &s.LastModifiedByType)
				delete(rawMsg, key)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// TrackedResource - The resource model definition for an Azure Resource Manager tracked top level resource which has 'tags' and a 'location'
type TrackedResource struct {
	Resource
	// REQUIRED; The geo-location where the resource lives
	Location *string `json:"location,omitempty"`

	// Resource tags.
	Tags map[string]*string `json:"tags,omitempty"`
}

// MarshalJSON implements the json.Marshaller interface for type TrackedResource.
func (t TrackedResource) MarshalJSON() ([]byte, error) {
	objectMap := make(map[string]interface{})
	t.marshalInternal(objectMap)
	return json.Marshal(objectMap)
}

func (t TrackedResource) marshalInternal(objectMap map[string]interface{}) {
	t.Resource.marshalInternal(objectMap)
	populate(objectMap, "location", t.Location)
	populate(objectMap, "tags", t.Tags)
}

func populate(m map[string]interface{}, k string, v interface{}) {
	if v == nil {
		return
	} else if azcore.IsNullValue(v) {
		m[k] = nil
	} else if !reflect.ValueOf(v).IsNil() {
		m[k] = v
	}
}

func unpopulate(data json.RawMessage, v interface{}) error {
	if data == nil {
		return nil
	}
	return json.Unmarshal(data, v)
}

