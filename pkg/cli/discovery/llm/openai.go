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

package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

// OpenAIProvider implements the LLM Provider interface using OpenAI.
type OpenAIProvider struct {
	client *openai.Client
	config *Config
}

// NewOpenAIProvider creates a new OpenAI provider with the given configuration.
func NewOpenAIProvider(config *Config) (*OpenAIProvider, error) {
	if config == nil {
		config = DefaultConfig()
	}

	apiKey := os.Getenv(config.APIKeyEnvVar)
	if apiKey == "" {
		return nil, fmt.Errorf("API key not found in environment variable %s", config.APIKeyEnvVar)
	}

	clientConfig := openai.DefaultConfig(apiKey)
	if config.Endpoint != "" {
		clientConfig.BaseURL = config.Endpoint
	}

	client := openai.NewClientWithConfig(clientConfig)

	return &OpenAIProvider{
		client: client,
		config: config,
	}, nil
}

// Name returns the provider identifier.
func (p *OpenAIProvider) Name() string {
	return "openai"
}

// Analyze sends a prompt to OpenAI and returns a structured JSON response.
func (p *OpenAIProvider) Analyze(ctx context.Context, prompt string, schema json.RawMessage) (json.RawMessage, error) {
	model := p.config.Model
	if model == "" {
		model = "gpt-4o"
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: "You are a code analysis assistant. Analyze the provided code and return structured JSON output according to the specified schema.",
		},
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	// Set temperature if configured
	if p.config.Temperature > 0 {
		temperature := float32(p.config.Temperature)
		req.Temperature = temperature
	}

	// Set max tokens if configured
	if p.config.MaxTokens > 0 {
		req.MaxTokens = p.config.MaxTokens
	}

	// Request JSON response format if schema is provided
	if len(schema) > 0 {
		req.ResponseFormat = &openai.ChatCompletionResponseFormat{
			Type: openai.ChatCompletionResponseFormatTypeJSONObject,
		}
	}

	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from OpenAI")
	}

	content := resp.Choices[0].Message.Content

	// Validate that the response is valid JSON
	var result json.RawMessage
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON response from OpenAI: %w", err)
	}

	return result, nil
}

// Chat sends a prompt to OpenAI and returns a text response.
func (p *OpenAIProvider) Chat(ctx context.Context, prompt string) (string, error) {
	model := p.config.Model
	if model == "" {
		model = "gpt-4o"
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: prompt,
		},
	}

	req := openai.ChatCompletionRequest{
		Model:    model,
		Messages: messages,
	}

	if p.config.Temperature > 0 {
		req.Temperature = float32(p.config.Temperature)
	}

	if p.config.MaxTokens > 0 {
		req.MaxTokens = p.config.MaxTokens
	}

	resp, err := p.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no response from OpenAI")
	}

	return resp.Choices[0].Message.Content, nil
}

// IsAvailable checks if the OpenAI provider is configured and accessible.
func (p *OpenAIProvider) IsAvailable(ctx context.Context) bool {
	// Check if API key is set
	apiKey := os.Getenv(p.config.APIKeyEnvVar)
	if apiKey == "" {
		return false
	}

	// Try a minimal API call to verify connectivity
	_, err := p.client.ListModels(ctx)
	return err == nil
}

// GetModel returns the configured model name.
func (p *OpenAIProvider) GetModel() string {
	if p.config.Model != "" {
		return p.config.Model
	}
	return "gpt-4o"
}
