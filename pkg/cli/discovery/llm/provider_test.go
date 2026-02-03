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
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	require.NotNil(t, config)

	assert.Equal(t, "openai", config.Provider)
	assert.Equal(t, "gpt-4o", config.Model)
	assert.Equal(t, "OPENAI_API_KEY", config.APIKeyEnvVar)
	assert.Equal(t, 0.1, config.Temperature)
	assert.Equal(t, 4096, config.MaxTokens)
}

func TestNewProviderRegistry(t *testing.T) {
	registry := NewProviderRegistry()
	require.NotNil(t, registry)

	providers := registry.List()
	assert.Empty(t, providers)
}

func TestProviderRegistry_Register(t *testing.T) {
	registry := NewProviderRegistry()
	mock := NewMockProvider()

	err := registry.Register(mock)
	require.NoError(t, err)

	providers := registry.List()
	assert.Len(t, providers, 1)
	assert.Contains(t, providers, "mock")
}

func TestProviderRegistry_RegisterNil(t *testing.T) {
	registry := NewProviderRegistry()

	err := registry.Register(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestProviderRegistry_Get(t *testing.T) {
	registry := NewProviderRegistry()
	mock := NewMockProvider()
	err := registry.Register(mock)
	require.NoError(t, err)

	// Get existing provider
	provider, ok := registry.Get("mock")
	assert.True(t, ok)
	assert.Equal(t, "mock", provider.Name())

	// Get non-existing provider
	_, ok = registry.Get("nonexistent")
	assert.False(t, ok)
}

func TestMockProvider_Name(t *testing.T) {
	mock := NewMockProvider()
	assert.Equal(t, "mock", mock.Name())
}

func TestMockProvider_IsAvailable(t *testing.T) {
	mock := NewMockProvider()

	// Default is available
	assert.True(t, mock.IsAvailable(context.Background()))

	// Set to unavailable
	mock.SetAvailable(false)
	assert.False(t, mock.IsAvailable(context.Background()))

	// Set back to available
	mock.SetAvailable(true)
	assert.True(t, mock.IsAvailable(context.Background()))
}

func TestMockProvider_Analyze(t *testing.T) {
	mock := NewMockProvider()
	ctx := context.Background()

	// Default returns empty object
	result, err := mock.Analyze(ctx, "test prompt", nil)
	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(result))

	// With custom response for specific prompt
	expectedResponse := json.RawMessage(`{"key": "value"}`)
	mock.SetAnalyzeResponse("specific prompt", expectedResponse)

	result, err = mock.Analyze(ctx, "specific prompt", nil)
	require.NoError(t, err)
	assert.JSONEq(t, string(expectedResponse), string(result))

	// Different prompt still returns default
	result, err = mock.Analyze(ctx, "other prompt", nil)
	require.NoError(t, err)
	assert.JSONEq(t, `{}`, string(result))
}

func TestMockProvider_AnalyzeWithFunc(t *testing.T) {
	mock := NewMockProvider()
	ctx := context.Background()

	expectedResponse := json.RawMessage(`{"custom": "response"}`)
	mock.SetAnalyzeFunc(func(ctx context.Context, prompt string, schema json.RawMessage) (json.RawMessage, error) {
		return expectedResponse, nil
	})

	result, err := mock.Analyze(ctx, "any prompt", nil)
	require.NoError(t, err)
	assert.JSONEq(t, string(expectedResponse), string(result))
}

func TestMockProvider_Chat(t *testing.T) {
	mock := NewMockProvider()
	ctx := context.Background()

	// Default returns empty string
	result, err := mock.Chat(ctx, "test prompt")
	require.NoError(t, err)
	assert.Empty(t, result)

	// With custom response for specific prompt
	mock.SetChatResponse("hello", "world")

	result, err = mock.Chat(ctx, "hello")
	require.NoError(t, err)
	assert.Equal(t, "world", result)
}

func TestMockProvider_ChatWithFunc(t *testing.T) {
	mock := NewMockProvider()
	ctx := context.Background()

	mock.SetChatFunc(func(ctx context.Context, prompt string) (string, error) {
		return "custom: " + prompt, nil
	})

	result, err := mock.Chat(ctx, "input")
	require.NoError(t, err)
	assert.Equal(t, "custom: input", result)
}

func TestMockProvider_GetCalls(t *testing.T) {
	mock := NewMockProvider()
	ctx := context.Background()

	// Make some calls
	_, _ = mock.Analyze(ctx, "analyze1", nil)
	_, _ = mock.Chat(ctx, "chat1")
	_, _ = mock.Analyze(ctx, "analyze2", json.RawMessage(`{}`))

	calls := mock.GetCalls()
	require.Len(t, calls, 3)

	assert.Equal(t, "Analyze", calls[0].Method)
	assert.Equal(t, "analyze1", calls[0].Prompt)

	assert.Equal(t, "Chat", calls[1].Method)
	assert.Equal(t, "chat1", calls[1].Prompt)

	assert.Equal(t, "Analyze", calls[2].Method)
	assert.Equal(t, "analyze2", calls[2].Prompt)
}

func TestMockProvider_ClearCalls(t *testing.T) {
	mock := NewMockProvider()
	ctx := context.Background()

	_, _ = mock.Analyze(ctx, "test", nil)
	assert.Len(t, mock.GetCalls(), 1)

	mock.ClearCalls()
	assert.Empty(t, mock.GetCalls())
}

func TestMockProvider_Reset(t *testing.T) {
	mock := NewMockProvider()
	ctx := context.Background()

	// Configure mock
	mock.SetAvailable(false)
	mock.SetAnalyzeResponse("test", json.RawMessage(`{"x": 1}`))
	mock.SetChatResponse("test", "response")
	_, _ = mock.Analyze(ctx, "test", nil)

	// Reset should clear everything
	mock.Reset()

	assert.True(t, mock.IsAvailable(ctx))
	assert.Empty(t, mock.GetCalls())

	result, _ := mock.Analyze(ctx, "test", nil)
	assert.JSONEq(t, `{}`, string(result))
}

func TestMockProvider_WithDependencyResponse(t *testing.T) {
	mock := NewMockProvider()
	ctx := context.Background()

	deps := []map[string]any{
		{"type": "postgresql", "confidence": 0.95},
		{"type": "redis", "confidence": 0.90},
	}
	mock.WithDependencyResponse(deps)

	result, err := mock.Analyze(ctx, "any prompt", nil)
	require.NoError(t, err)

	var parsed map[string]any
	err = json.Unmarshal(result, &parsed)
	require.NoError(t, err)

	dependencies, ok := parsed["dependencies"].([]any)
	require.True(t, ok)
	assert.Len(t, dependencies, 2)
}

func TestMockProvider_WithServiceResponse(t *testing.T) {
	mock := NewMockProvider()
	ctx := context.Background()

	services := []map[string]any{
		{"name": "api", "type": "express", "port": 3000},
	}
	mock.WithServiceResponse(services)

	result, err := mock.Analyze(ctx, "any prompt", nil)
	require.NoError(t, err)

	var parsed map[string]any
	err = json.Unmarshal(result, &parsed)
	require.NoError(t, err)

	svcList, ok := parsed["services"].([]any)
	require.True(t, ok)
	assert.Len(t, svcList, 1)
}

func TestMockProvider_WithError(t *testing.T) {
	mock := NewMockProvider()
	ctx := context.Background()

	expectedError := errors.New("LLM unavailable")
	mock.WithError(expectedError)

	_, err := mock.Analyze(ctx, "test", nil)
	require.Error(t, err)
	assert.Equal(t, expectedError, err)

	_, err = mock.Chat(ctx, "test")
	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestOllamaProvider_Name(t *testing.T) {
	provider := NewOllamaProvider(nil)
	assert.Equal(t, "ollama", provider.Name())
}

func TestOllamaProvider_DefaultConfig(t *testing.T) {
	provider := NewOllamaProvider(nil)
	assert.Equal(t, "llama3.2", provider.GetModel())
}

func TestOllamaProvider_CustomConfig(t *testing.T) {
	config := &Config{
		Provider: "ollama",
		Model:    "codellama",
		Endpoint: "http://custom:11434",
	}
	provider := NewOllamaProvider(config)
	assert.Equal(t, "codellama", provider.GetModel())
}
