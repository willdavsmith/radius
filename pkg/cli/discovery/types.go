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

package discovery

import (
	"time"
)

// DiscoveryResult is the output of the discovery phase containing all detected information.
type DiscoveryResult struct {
	// Version is the schema version (e.g., "1.0").
	Version string `json:"version" yaml:"version"`

	// AnalyzedAt is when the analysis was performed.
	AnalyzedAt time.Time `json:"analyzedAt" yaml:"analyzed_at"`

	// CodebasePath is the absolute path to the analyzed codebase.
	CodebasePath string `json:"codebasePath" yaml:"codebase_path"`

	// Dependencies lists detected infrastructure dependencies.
	Dependencies []DetectedDependency `json:"dependencies" yaml:"dependencies"`

	// Services lists detected deployable services.
	Services []DetectedService `json:"services" yaml:"services"`

	// Practices contains detected or configured team practices.
	Practices *TeamPractices `json:"practices,omitempty" yaml:"practices,omitempty"`

	// Warnings lists partial failures or low-confidence detections.
	Warnings []Warning `json:"warnings,omitempty" yaml:"warnings,omitempty"`

	// State indicates the current state of the discovery process.
	State DiscoveryState `json:"state" yaml:"state"`
}

// DiscoveryState represents the state of a discovery operation.
type DiscoveryState string

const (
	// StatePending indicates discovery has not started.
	StatePending DiscoveryState = "pending"

	// StateAnalyzing indicates discovery is in progress.
	StateAnalyzing DiscoveryState = "analyzing"

	// StateComplete indicates discovery completed successfully.
	StateComplete DiscoveryState = "complete"

	// StateFailed indicates discovery failed.
	StateFailed DiscoveryState = "failed"
)

// DetectedDependency is an infrastructure dependency detected from codebase analysis.
type DetectedDependency struct {
	// Type is the normalized dependency type (e.g., "postgresql", "redis").
	Type string `json:"type" yaml:"type"`

	// Technology is the specific technology/service (e.g., "PostgreSQL", "Azure Redis").
	Technology string `json:"technology" yaml:"technology"`

	// Confidence is the detection confidence score (0.0-1.0).
	Confidence float64 `json:"confidence" yaml:"confidence"`

	// Evidence lists sources supporting this detection.
	Evidence []Evidence `json:"evidence" yaml:"evidence"`

	// ConnectionInfo contains detected connection configuration.
	ConnectionInfo *ConnectionInfo `json:"connectionInfo,omitempty" yaml:"connection_info,omitempty"`
}

// ConfidenceLevel returns the confidence tier based on the score.
func (d *DetectedDependency) ConfidenceLevel() string {
	switch {
	case d.Confidence >= 0.80:
		return "high"
	case d.Confidence >= 0.50:
		return "medium"
	default:
		return "low"
	}
}

// Evidence is supporting evidence for a detected dependency.
type Evidence struct {
	// Source is the file path relative to codebase root.
	Source string `json:"source" yaml:"source"`

	// Type is the evidence type: "package", "import", "config", "connection_string".
	Type string `json:"type" yaml:"type"`

	// Value is the detected value (e.g., "pg@8.11.0").
	Value string `json:"value" yaml:"value"`

	// Line is the line number in the source file.
	Line int `json:"line,omitempty" yaml:"line,omitempty"`
}

// ConnectionInfo is detected connection configuration for a dependency.
type ConnectionInfo struct {
	// EnvVar is the environment variable name (e.g., "DATABASE_URL").
	EnvVar string `json:"envVar,omitempty" yaml:"env_var,omitempty"`

	// DefaultPort is the default port for the service.
	DefaultPort int `json:"defaultPort,omitempty" yaml:"default_port,omitempty"`

	// HostPattern is the detected hostname pattern.
	HostPattern string `json:"hostPattern,omitempty" yaml:"host_pattern,omitempty"`
}

// DetectedService is a deployable service/entrypoint detected in the codebase.
type DetectedService struct {
	// Name is the service name (derived from directory or config).
	Name string `json:"name" yaml:"name"`

	// Type is the service type (e.g., "http", "grpc", "worker").
	Type string `json:"type" yaml:"type"`

	// Entrypoint is the main file or command.
	Entrypoint string `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`

	// Port is the detected or configured port.
	Port int `json:"port,omitempty" yaml:"port,omitempty"`

	// Framework is the framework (e.g., "express", "flask", "spring").
	Framework string `json:"framework,omitempty" yaml:"framework,omitempty"`

	// Language is the programming language (e.g., "nodejs", "python", "go").
	Language string `json:"language,omitempty" yaml:"language,omitempty"`

	// Confidence is the detection confidence score (0.0-1.0).
	Confidence float64 `json:"confidence,omitempty" yaml:"confidence,omitempty"`

	// Dockerfile is the path to Dockerfile if detected.
	Dockerfile string `json:"dockerfile,omitempty" yaml:"dockerfile,omitempty"`

	// BuildContext is the build context directory.
	BuildContext string `json:"buildContext,omitempty" yaml:"build_context,omitempty"`

	// Dependencies lists dependencies specific to this service.
	Dependencies []string `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
}

// TeamPractices contains infrastructure conventions detected or configured for the team.
type TeamPractices struct {
	// Naming is the naming pattern for resources.
	Naming *NamingConvention `json:"naming,omitempty" yaml:"naming,omitempty"`

	// Tags is the tag policy.
	Tags *TagPolicy `json:"tags,omitempty" yaml:"tags,omitempty"`

	// Sizing contains environment-specific sizing policies.
	Sizing map[string]*SizingPolicy `json:"sizing,omitempty" yaml:"sizing,omitempty"`

	// Sources lists where practices were detected from.
	Sources []PracticeSource `json:"sources,omitempty" yaml:"sources,omitempty"`
}

// NamingConvention is the naming pattern for infrastructure resources.
type NamingConvention struct {
	// Pattern is the pattern with placeholders (e.g., "{env}-{service}-{resource}").
	Pattern string `json:"pattern" yaml:"pattern"`

	// MaxLength is the maximum resource name length.
	MaxLength int `json:"maxLength,omitempty" yaml:"max_length,omitempty"`

	// EnvironmentAbbreviations maps environment names to short forms.
	EnvironmentAbbreviations map[string]string `json:"environmentAbbreviations,omitempty" yaml:"environment_abbreviations,omitempty"`
}

// TagPolicy defines required and recommended tags for resources.
type TagPolicy struct {
	// Required lists tags that must be present.
	Required []string `json:"required" yaml:"required"`

	// Recommended lists tags that should be present.
	Recommended []string `json:"recommended,omitempty" yaml:"recommended,omitempty"`

	// Defaults provides default values for tags.
	Defaults map[string]string `json:"defaults,omitempty" yaml:"defaults,omitempty"`
}

// SizingPolicy defines environment-specific sizing and configuration.
type SizingPolicy struct {
	// DefaultTier is the default pricing tier.
	DefaultTier string `json:"defaultTier" yaml:"default_tier"`

	// HAEnabled indicates if high availability is enabled.
	HAEnabled bool `json:"haEnabled" yaml:"ha_enabled"`

	// BackupEnabled indicates if backup is enabled.
	BackupEnabled bool `json:"backupEnabled,omitempty" yaml:"backup_enabled,omitempty"`

	// GeoRedundant indicates if geo-redundancy is enabled.
	GeoRedundant bool `json:"geoRedundant,omitempty" yaml:"geo_redundant,omitempty"`
}

// PracticeSource is a location where a practice was detected.
type PracticeSource struct {
	// Type is the source type: "iac", "config", "wiki", "manual".
	Type string `json:"type" yaml:"type"`

	// Location is the file path or URL.
	Location string `json:"location" yaml:"location"`

	// DetectedAt is when this source was analyzed.
	DetectedAt time.Time `json:"detectedAt" yaml:"detected_at"`
}

// Warning represents a partial failure or low-confidence detection.
type Warning struct {
	// Code is a machine-readable warning code.
	Code string `json:"code" yaml:"code"`

	// Message is a human-readable warning message.
	Message string `json:"message" yaml:"message"`

	// Source is the source file or component that generated the warning.
	Source string `json:"source,omitempty" yaml:"source,omitempty"`
}

// ResourceType is an interface defining what an application consumes from infrastructure.
type ResourceType struct {
	// Name is the human-friendly name (e.g., "PostgreSQL Database").
	Name string `json:"name" yaml:"name"`

	// Type is the Radius type (e.g., "Applications.Datastores/postgreSql").
	Type string `json:"type" yaml:"type"`

	// APIVersion is the API version.
	APIVersion string `json:"apiVersion" yaml:"api_version"`

	// Schema is the input/output schema.
	Schema *ResourceSchema `json:"schema,omitempty" yaml:"schema,omitempty"`

	// Outputs lists outputs provided to containers.
	Outputs []OutputDefinition `json:"outputs" yaml:"outputs"`

	// Maturity is "alpha", "beta", or "stable".
	Maturity string `json:"maturity" yaml:"maturity"`
}

// ResourceSchema defines schema for a Resource Type.
type ResourceSchema struct {
	// Properties maps property names to definitions.
	Properties map[string]PropertyDef `json:"properties" yaml:"properties"`

	// Required lists required property names.
	Required []string `json:"required,omitempty" yaml:"required,omitempty"`
}

// PropertyDef defines a single property.
type PropertyDef struct {
	// Type is the JSON type.
	Type string `json:"type" yaml:"type"`

	// Description is the property description.
	Description string `json:"description" yaml:"description"`

	// Default is the default value.
	Default any `json:"default,omitempty" yaml:"default,omitempty"`

	// Enum lists allowed values.
	Enum []any `json:"enum,omitempty" yaml:"enum,omitempty"`

	// Sensitive indicates if the value is sensitive.
	Sensitive bool `json:"sensitive,omitempty" yaml:"sensitive,omitempty"`
}

// OutputDefinition is an output provided by a Resource Type to containers.
type OutputDefinition struct {
	// Name is the output name (e.g., "connectionString").
	Name string `json:"name" yaml:"name"`

	// Type is the output type.
	Type string `json:"type" yaml:"type"`

	// Description describes the output.
	Description string `json:"description" yaml:"description"`

	// Sensitive indicates if the output is sensitive.
	Sensitive bool `json:"sensitive,omitempty" yaml:"sensitive,omitempty"`
}

// Recipe is an IaC implementation that provisions infrastructure.
type Recipe struct {
	// Name is the recipe name.
	Name string `json:"name" yaml:"name"`

	// Description is a human-readable description.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// IaCType is "bicep" or "terraform".
	IaCType string `json:"iacType" yaml:"iac_type"`

	// TemplatePath is the registry path or URL.
	TemplatePath string `json:"templatePath" yaml:"template_path"`

	// Parameters defines recipe parameters.
	Parameters map[string]ParameterDef `json:"parameters,omitempty" yaml:"parameters,omitempty"`

	// Source is where the recipe was found.
	Source *RecipeSource `json:"source" yaml:"source"`

	// Environments lists suitable environments.
	Environments []string `json:"environments,omitempty" yaml:"environments,omitempty"`

	// ResourceType is the target Resource Type.
	ResourceType string `json:"resourceType" yaml:"resource_type"`

	// Verified indicates if from a verified source (e.g., AVM).
	Verified bool `json:"verified,omitempty" yaml:"verified,omitempty"`
}

// ParameterDef defines a recipe parameter.
type ParameterDef struct {
	// Type is the parameter type.
	Type string `json:"type" yaml:"type"`

	// Description describes the parameter.
	Description string `json:"description,omitempty" yaml:"description,omitempty"`

	// Default is the default value.
	Default any `json:"default,omitempty" yaml:"default,omitempty"`

	// Required indicates if the parameter is required.
	Required bool `json:"required,omitempty" yaml:"required,omitempty"`
}

// RecipeSource is a location for discovering Recipes.
type RecipeSource struct {
	// Type is "avm", "oci", "terraform-registry", or "filesystem".
	Type string `json:"type" yaml:"type"`

	// Location is the registry URL or path.
	Location string `json:"location" yaml:"location"`

	// Priority is the search priority (lower = higher priority).
	Priority int `json:"priority" yaml:"priority"`

	// AuthConfig is the authentication configuration.
	AuthConfig *AuthConfig `json:"authConfig,omitempty" yaml:"auth_config,omitempty"`

	// RefreshFrequency is how often to refresh the cache.
	RefreshFrequency string `json:"refreshFrequency,omitempty" yaml:"refresh_frequency,omitempty"`
}

// AuthConfig is authentication configuration for recipe sources.
type AuthConfig struct {
	// Type is "none", "token", "oidc", or "env".
	Type string `json:"type" yaml:"type"`

	// TokenEnvVar is the environment variable containing the token.
	TokenEnvVar string `json:"tokenEnvVar,omitempty" yaml:"token_env_var,omitempty"`

	// CredentialHelper is the credential helper command.
	CredentialHelper string `json:"credentialHelper,omitempty" yaml:"credential_helper,omitempty"`
}

// RecipeSelection is a user's selection of a recipe for a dependency.
type RecipeSelection struct {
	// Dependency is the dependency type.
	Dependency string `json:"dependency" yaml:"dependency"`

	// Recipe is the selected recipe.
	Recipe *Recipe `json:"recipe" yaml:"recipe"`

	// Parameters are parameter overrides.
	Parameters map[string]any `json:"parameters,omitempty" yaml:"parameters,omitempty"`

	// Environment is the target environment.
	Environment string `json:"environment,omitempty" yaml:"environment,omitempty"`
}

// AppDefinition is a generated Radius application definition.
type AppDefinition struct {
	// Name is the application name.
	Name string `json:"name" yaml:"name"`

	// Containers lists container definitions.
	Containers []ContainerDef `json:"containers" yaml:"containers"`

	// Resources lists resource definitions.
	Resources []ResourceDef `json:"resources" yaml:"resources"`

	// Connections lists container-to-resource connections.
	Connections []ConnectionDef `json:"connections" yaml:"connections"`

	// BicepContent is the generated Bicep content.
	BicepContent string `json:"bicepContent" yaml:"bicep_content"`

	// Valid indicates if Bicep validated successfully.
	Valid bool `json:"valid" yaml:"valid"`

	// Errors lists validation errors.
	Errors []ValidationError `json:"errors,omitempty" yaml:"errors,omitempty"`
}

// ContainerDef is a container definition in the application.
type ContainerDef struct {
	// Name is the container name.
	Name string `json:"name" yaml:"name"`

	// Image is the container image (or TODO placeholder).
	Image string `json:"image,omitempty" yaml:"image,omitempty"`

	// Ports lists exposed ports.
	Ports []PortDef `json:"ports,omitempty" yaml:"ports,omitempty"`

	// Env maps environment variable names to values.
	Env map[string]EnvValue `json:"env,omitempty" yaml:"env,omitempty"`
}

// PortDef defines a container port.
type PortDef struct {
	// ContainerPort is the port inside the container.
	ContainerPort int `json:"containerPort" yaml:"container_port"`

	// Protocol is the protocol (TCP, UDP).
	Protocol string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}

// EnvValue is an environment variable value.
type EnvValue struct {
	// Value is a literal value.
	Value string `json:"value,omitempty" yaml:"value,omitempty"`

	// ValueFrom references a value from another source.
	ValueFrom *EnvValueFrom `json:"valueFrom,omitempty" yaml:"value_from,omitempty"`
}

// EnvValueFrom references an environment variable value from another source.
type EnvValueFrom struct {
	// Resource is the resource name.
	Resource string `json:"resource,omitempty" yaml:"resource,omitempty"`

	// Output is the output name from the resource.
	Output string `json:"output,omitempty" yaml:"output,omitempty"`
}

// ResourceDef is a resource definition in the application.
type ResourceDef struct {
	// Name is the resource name.
	Name string `json:"name" yaml:"name"`

	// Type is the Resource Type.
	Type string `json:"type" yaml:"type"`

	// Recipe is the recipe reference.
	Recipe string `json:"recipe,omitempty" yaml:"recipe,omitempty"`

	// Properties are resource properties.
	Properties map[string]any `json:"properties,omitempty" yaml:"properties,omitempty"`
}

// ConnectionDef is a connection between container and resource.
type ConnectionDef struct {
	// Container is the container name.
	Container string `json:"container" yaml:"container"`

	// Resource is the resource name.
	Resource string `json:"resource" yaml:"resource"`

	// EnvMapping maps output names to environment variables.
	EnvMapping map[string]string `json:"envMapping,omitempty" yaml:"env_mapping,omitempty"`
}

// ValidationError is a Bicep validation error.
type ValidationError struct {
	// Code is a machine-readable error code.
	Code string `json:"code" yaml:"code"`

	// Message is a human-readable error message.
	Message string `json:"message" yaml:"message"`

	// Line is the line number in the Bicep file.
	Line int `json:"line,omitempty" yaml:"line,omitempty"`

	// Column is the column number.
	Column int `json:"column,omitempty" yaml:"column,omitempty"`
}
