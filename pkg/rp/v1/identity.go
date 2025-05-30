/*
Copyright 2023 The Radius Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import "errors"

// IdentitySettingKind represents the kind of identity setting.
type IdentitySettingKind string

const (
	// IdentityNone represents unknown identity.
	IdentityNone IdentitySettingKind = "None"
	// AzureIdentityWorkload represents Azure Workload identity.
	AzureIdentityWorkload IdentitySettingKind = "azure.com.workload"
	// UserAssigned represents Azure User Assigned Managed Identity.
	UserAssigned IdentitySettingKind = "userAssigned"
	// SystemAssigned represents Azure System Assigned Managed Identity.
	SystemAssigned IdentitySettingKind = "systemAssigned"
	// SystemAssignedUserAssigned represents Azure System Assigned and User Assigned Managed Identity.
	SystemAssignedUserAssigned IdentitySettingKind = "systemAssignedUserAssigned"
)

// IdentitySettings represents the identity info to access azure resource, such as Key vault.
type IdentitySettings struct {
	// Kind represents the type of authentication.
	Kind IdentitySettingKind `json:"kind"`
	// OIDCIssuer represents the name of OIDC issuer.
	OIDCIssuer string `json:"oidcIssuer,omitempty"`
	// Resource represents the resource id of managed identity.
	Resource string `json:"resource,omitempty"`
	// ManagedIdentity represents the list of user assigned managed identities.
	ManagedIdentity []string `json:"managedIdentity,omitempty"`
}

// Validate checks if the IdentitySettings struct is nil and if the Kind is AzureIdentityWorkload, checks if the OIDCIssuer
// is empty and if the Resource is not empty. It returns an error if any of these conditions are not met.
func (is *IdentitySettings) Validate() error {
	if is == nil {
		return nil
	}

	if is.Kind == AzureIdentityWorkload {
		if is.OIDCIssuer == "" {
			return errors.New(".properties.oidcIssuer is required for workload identity")
		}
		if is.Resource != "" {
			return errors.New(".properties.resource is read-only property")
		}
	}

	if is.Kind == UserAssigned || is.Kind == SystemAssignedUserAssigned {
		if is.ManagedIdentity == nil {
			return errors.New(".properties.managedIdentity is required for user assigned identity")
		}
	}

	if is.Kind == SystemAssigned {
		if is.ManagedIdentity != nil {
			return errors.New("no managed identities can be set for system assigned identity")
		}
	}

	return nil
}
