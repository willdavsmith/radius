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

// Package llm provides abstraction for LLM providers used in codebase analysis.
package llm

import (
	"context"
	"encoding/json"
	"fmt"
)

// Provider defines the interface for LLM providers.
type Provider interface {
	// Name returns the provider identifier.
	Name() string

	// Analyze sends a prompt and returns a structured JSON response.
	Analyze(ctx context.Context, prompt string, schema json.RawMessage) (json.RawMessage, error)

	// Chat sends a prompt and returns a text response.
	Chat(ctx context.Context, prompt string) (string, error)

	// IsAvailable checks if the provider is accessible.
	IsAvailable(ctx context.Context) bool
}

// Config holds configuration for an LLM provider.
type Config struct {
	Provider     string  `json:"provider" yaml:"provider"`
	Model        string  `json:"model,omitempty" yaml:"model,omitempty"`
	Endpoint     string  `json:"endpoint,omitempty" yaml:"endpoint,omitempty"`
	APIKeyEnvVar string  `json:"apiKeyEnvVar,omitempty" yaml:"apiKeyEnvVar,omitempty"`
	Temperature  float64 `json:"temperature,omitempty" yaml:"temperature,omitempty"`
	MaxTokens    int     `json:"maxTokens,omitempty" yaml:"maxTokens,omitempty"`
}

// DefaultConfig returns default LLM configuration.
func DefaultConfig() *Config {
	return &Config{
		Provider:     "openai",
		Model:        "gpt-4o",
		APIKeyEnvVar: "OPENAI_API_KEY",
		Temperature:  0.1,
		MaxTokens:    4096,
	}
}

// ProviderRegistry manages LLM provider instances.
type ProviderRegistry struct {
	providers map[string]Provider
}

// NewProviderRegistry creates a new provider registry.
func NewProviderRegistry() *ProviderRegistry {
	return &ProviderRegistry{
		providers: make(map[string]Provider),
	}
}

// Register adds a provider to the registry.
func (r *ProviderRegistry) Register(provider Provider) error {
	if provider == nil {
		return fmt.Errorf("provider cannot be nil")
	}
	r.providers[provider.Name()] = provider
	return nil
}

// Get returns a provider by name.
func (r *ProviderRegistry) Get(name string) (Provider, bool) {
	provider, ok := r.providers[name]
	return provider, ok
}

// List returns all registered provider names.
func (r *ProviderRegistry) List() []string {
	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

// GetPreferred returns the first available provider, optionally filtering by name.
func (r *ProviderRegistry) GetPreferred(preferred string) Provider {
	// If a specific provider is requested, return it if available
	if preferred != "" {
		if p, ok := r.providers[preferred]; ok {
			return p
		}
		return nil
	}

	// Return the first registered provider
	for _, p := range r.providers {
		return p
	}
	return nil
}

// AnalysisRequest represents a request to analyze content.
type AnalysisRequest struct {
	Content string          `json:"content"`
	Context string          `json:"context,omitempty"`
	Task    string          `json:"task"`
	Schema  json.RawMessage `json:"schema,omitempty"`
}

// AnalysisResult represents the result of an analysis.
type AnalysisResult struct {
	Data       json.RawMessage `json:"data"`
	Confidence float64         `json:"confidence"`
	Tokens     TokenUsage      `json:"tokens"`
}

// TokenUsage tracks token consumption.
type TokenUsage struct {
	PromptTokens     int `json:"promptTokens"`
	CompletionTokens int `json:"completionTokens"`
	TotalTokens      int `json:"totalTokens"`
}
