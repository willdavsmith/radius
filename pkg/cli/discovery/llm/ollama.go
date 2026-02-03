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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultOllamaEndpoint = "http://localhost:11434"
	defaultOllamaModel    = "llama3.2"
)

// OllamaProvider implements the LLM Provider interface using Ollama.
type OllamaProvider struct {
	config     *Config
	httpClient *http.Client
}

// NewOllamaProvider creates a new Ollama provider with the given configuration.
func NewOllamaProvider(config *Config) *OllamaProvider {
	if config == nil {
		config = &Config{
			Provider: "ollama",
			Model:    defaultOllamaModel,
			Endpoint: defaultOllamaEndpoint,
		}
	}

	if config.Endpoint == "" {
		config.Endpoint = defaultOllamaEndpoint
	}

	if config.Model == "" {
		config.Model = defaultOllamaModel
	}

	return &OllamaProvider{
		config: config,
		httpClient: &http.Client{
			Timeout: 5 * time.Minute, // LLM calls can be slow
		},
	}
}

// Name returns the provider identifier.
func (p *OllamaProvider) Name() string {
	return "ollama"
}

// ollamaRequest represents a request to the Ollama API.
type ollamaRequest struct {
	Model   string `json:"model"`
	Prompt  string `json:"prompt"`
	Stream  bool   `json:"stream"`
	Options struct {
		Temperature float64 `json:"temperature,omitempty"`
	} `json:"options,omitempty"`
	Format string `json:"format,omitempty"`
}

// ollamaResponse represents a response from the Ollama API.
type ollamaResponse struct {
	Model    string `json:"model"`
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Analyze sends a prompt to Ollama and returns a structured JSON response.
func (p *OllamaProvider) Analyze(ctx context.Context, prompt string, schema json.RawMessage) (json.RawMessage, error) {
	systemPrompt := "You are a code analysis assistant. Analyze the provided code and return structured JSON output according to the specified schema. Only output valid JSON, no other text."

	fullPrompt := fmt.Sprintf("%s\n\n%s", systemPrompt, prompt)

	req := ollamaRequest{
		Model:  p.config.Model,
		Prompt: fullPrompt,
		Stream: false,
		Format: "json",
	}

	if p.config.Temperature > 0 {
		req.Options.Temperature = p.config.Temperature
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/generate", p.config.Endpoint)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("Ollama API error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return nil, fmt.Errorf("failed to parse Ollama response: %w", err)
	}

	// Validate that the response is valid JSON
	var result json.RawMessage
	if err := json.Unmarshal([]byte(ollamaResp.Response), &result); err != nil {
		return nil, fmt.Errorf("invalid JSON response from Ollama: %w", err)
	}

	return result, nil
}

// Chat sends a prompt to Ollama and returns a text response.
func (p *OllamaProvider) Chat(ctx context.Context, prompt string) (string, error) {
	req := ollamaRequest{
		Model:  p.config.Model,
		Prompt: prompt,
		Stream: false,
	}

	if p.config.Temperature > 0 {
		req.Options.Temperature = p.config.Temperature
	}

	body, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/api/generate", p.config.Endpoint)
	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := p.httpClient.Do(httpReq)
	if err != nil {
		return "", fmt.Errorf("Ollama API error: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("Ollama API returned status %d: %s", resp.StatusCode, string(respBody))
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp ollamaResponse
	if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to parse Ollama response: %w", err)
	}

	return ollamaResp.Response, nil
}

// IsAvailable checks if Ollama is running and accessible.
func (p *OllamaProvider) IsAvailable(ctx context.Context) bool {
	endpoint := fmt.Sprintf("%s/api/tags", p.config.Endpoint)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return false
	}

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// GetModel returns the configured model name.
func (p *OllamaProvider) GetModel() string {
	return p.config.Model
}
