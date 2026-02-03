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

package config

// DiscoveryConfig contains configuration for the automatic application discovery feature.
type DiscoveryConfig struct {
	// OutputPath is the directory where discovery outputs are written (default: "./radius").
	OutputPath string `json:"outputPath,omitempty" mapstructure:"outputPath" yaml:"outputPath,omitempty"`

	// ConfidenceThreshold is the minimum confidence score for detections (default: 0.5).
	ConfidenceThreshold float64 `json:"confidenceThreshold,omitempty" mapstructure:"confidenceThreshold" yaml:"confidenceThreshold,omitempty"`

	// ExcludePaths lists paths to exclude from analysis.
	ExcludePaths []string `json:"excludePaths,omitempty" mapstructure:"excludePaths" yaml:"excludePaths,omitempty"`

	// MaxFiles limits the number of files to analyze (default: 1000).
	MaxFiles int `json:"maxFiles,omitempty" mapstructure:"maxFiles" yaml:"maxFiles,omitempty"`

	// LLM contains LLM provider configuration.
	LLM *LLMConfig `json:"llm,omitempty" mapstructure:"llm" yaml:"llm,omitempty"`

	// RecipeSources lists configured recipe sources.
	RecipeSources []RecipeSourceConfig `json:"recipeSources,omitempty" mapstructure:"recipeSources" yaml:"recipeSources,omitempty"`
}

// LLMConfig contains configuration for the LLM provider.
type LLMConfig struct {
	// Provider is the LLM provider type ("openai", "ollama", "anthropic").
	Provider string `json:"provider" mapstructure:"provider" yaml:"provider"`

	// Model is the model name (e.g., "gpt-4o", "llama3.2").
	Model string `json:"model,omitempty" mapstructure:"model" yaml:"model,omitempty"`

	// Endpoint is the custom endpoint URL for self-hosted providers.
	Endpoint string `json:"endpoint,omitempty" mapstructure:"endpoint" yaml:"endpoint,omitempty"`

	// APIKeyEnvVar is the environment variable containing the API key.
	APIKeyEnvVar string `json:"apiKeyEnvVar,omitempty" mapstructure:"apiKeyEnvVar" yaml:"apiKeyEnvVar,omitempty"`

	// Temperature controls randomness (0.0-1.0). Lower is more deterministic.
	Temperature float64 `json:"temperature,omitempty" mapstructure:"temperature" yaml:"temperature,omitempty"`

	// MaxTokens limits the response length.
	MaxTokens int `json:"maxTokens,omitempty" mapstructure:"maxTokens" yaml:"maxTokens,omitempty"`
}

// RecipeSourceConfig contains configuration for a recipe source.
type RecipeSourceConfig struct {
	// Name is a friendly name for the source.
	Name string `json:"name" mapstructure:"name" yaml:"name"`

	// Type is the source type ("avm", "oci", "terraform-registry", "filesystem").
	Type string `json:"type" mapstructure:"type" yaml:"type"`

	// Location is the registry URL or path.
	Location string `json:"location" mapstructure:"location" yaml:"location"`

	// Priority is the search priority (lower = higher priority).
	Priority int `json:"priority,omitempty" mapstructure:"priority" yaml:"priority,omitempty"`

	// AuthConfig contains authentication settings.
	AuthConfig *RecipeSourceAuthConfig `json:"authConfig,omitempty" mapstructure:"authConfig" yaml:"authConfig,omitempty"`
}

// RecipeSourceAuthConfig contains authentication configuration for a recipe source.
type RecipeSourceAuthConfig struct {
	// Type is the auth type ("none", "token", "oidc", "env").
	Type string `json:"type" mapstructure:"type" yaml:"type"`

	// TokenEnvVar is the environment variable containing the token.
	TokenEnvVar string `json:"tokenEnvVar,omitempty" mapstructure:"tokenEnvVar" yaml:"tokenEnvVar,omitempty"`

	// CredentialHelper is the credential helper command.
	CredentialHelper string `json:"credentialHelper,omitempty" mapstructure:"credentialHelper" yaml:"credentialHelper,omitempty"`
}

// DefaultDiscoveryConfig returns the default discovery configuration.
func DefaultDiscoveryConfig() *DiscoveryConfig {
	return &DiscoveryConfig{
		OutputPath:          "./radius",
		ConfidenceThreshold: 0.5,
		ExcludePaths: []string{
			"node_modules",
			"vendor",
			".git",
			"__pycache__",
			".venv",
			"venv",
			"dist",
			"build",
			"target",
			".next",
			"coverage",
			".nyc_output",
		},
		MaxFiles: 1000,
		LLM: &LLMConfig{
			Provider:     "openai",
			Model:        "gpt-4o",
			APIKeyEnvVar: "OPENAI_API_KEY",
			Temperature:  0.1,
			MaxTokens:    4096,
		},
		RecipeSources: []RecipeSourceConfig{
			{
				Name:     "avm",
				Type:     "avm",
				Location: "https://github.com/Azure/bicep-registry-modules",
				Priority: 1,
			},
		},
	}
}

// Merge merges another configuration into this one, preferring non-zero values from other.
func (c *DiscoveryConfig) Merge(other *DiscoveryConfig) {
	if other == nil {
		return
	}

	if other.OutputPath != "" {
		c.OutputPath = other.OutputPath
	}
	if other.ConfidenceThreshold > 0 {
		c.ConfidenceThreshold = other.ConfidenceThreshold
	}
	if len(other.ExcludePaths) > 0 {
		c.ExcludePaths = other.ExcludePaths
	}
	if other.MaxFiles > 0 {
		c.MaxFiles = other.MaxFiles
	}
	if other.LLM != nil {
		if c.LLM == nil {
			c.LLM = &LLMConfig{}
		}
		if other.LLM.Provider != "" {
			c.LLM.Provider = other.LLM.Provider
		}
		if other.LLM.Model != "" {
			c.LLM.Model = other.LLM.Model
		}
		if other.LLM.Endpoint != "" {
			c.LLM.Endpoint = other.LLM.Endpoint
		}
		if other.LLM.APIKeyEnvVar != "" {
			c.LLM.APIKeyEnvVar = other.LLM.APIKeyEnvVar
		}
		if other.LLM.Temperature > 0 {
			c.LLM.Temperature = other.LLM.Temperature
		}
		if other.LLM.MaxTokens > 0 {
			c.LLM.MaxTokens = other.LLM.MaxTokens
		}
	}
	if len(other.RecipeSources) > 0 {
		c.RecipeSources = append(c.RecipeSources, other.RecipeSources...)
	}
}
