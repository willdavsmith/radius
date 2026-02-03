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
	"sync"
)

// MockProvider is a mock LLM provider for testing.
type MockProvider struct {
	mu            sync.Mutex
	responses     map[string]json.RawMessage
	chatResponses map[string]string
	analyzeFunc   func(ctx context.Context, prompt string, schema json.RawMessage) (json.RawMessage, error)
	chatFunc      func(ctx context.Context, prompt string) (string, error)
	available     bool
	calls         []MockCall
}

// MockCall records a call to the mock provider.
type MockCall struct {
	Method string
	Prompt string
	Schema json.RawMessage
}

// NewMockProvider creates a new mock provider.
func NewMockProvider() *MockProvider {
	return &MockProvider{
		responses:     make(map[string]json.RawMessage),
		chatResponses: make(map[string]string),
		available:     true,
		calls:         make([]MockCall, 0),
	}
}

// Name returns the provider identifier.
func (m *MockProvider) Name() string {
	return "mock"
}

// Analyze returns a pre-configured response or calls the custom function.
func (m *MockProvider) Analyze(ctx context.Context, prompt string, schema json.RawMessage) (json.RawMessage, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.calls = append(m.calls, MockCall{
		Method: "Analyze",
		Prompt: prompt,
		Schema: schema,
	})

	if m.analyzeFunc != nil {
		return m.analyzeFunc(ctx, prompt, schema)
	}

	if response, ok := m.responses[prompt]; ok {
		return response, nil
	}

	return json.RawMessage(`{}`), nil
}

// Chat returns a pre-configured response or calls the custom function.
func (m *MockProvider) Chat(ctx context.Context, prompt string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.calls = append(m.calls, MockCall{
		Method: "Chat",
		Prompt: prompt,
	})

	if m.chatFunc != nil {
		return m.chatFunc(ctx, prompt)
	}

	if response, ok := m.chatResponses[prompt]; ok {
		return response, nil
	}

	return "", nil
}

// IsAvailable returns whether the mock provider is available.
func (m *MockProvider) IsAvailable(_ context.Context) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.available
}

// SetAnalyzeResponse sets a response for a specific prompt.
func (m *MockProvider) SetAnalyzeResponse(prompt string, response json.RawMessage) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses[prompt] = response
}

// SetChatResponse sets a chat response for a specific prompt.
func (m *MockProvider) SetChatResponse(prompt, response string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.chatResponses[prompt] = response
}

// SetAnalyzeFunc sets a custom function for Analyze calls.
func (m *MockProvider) SetAnalyzeFunc(fn func(ctx context.Context, prompt string, schema json.RawMessage) (json.RawMessage, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.analyzeFunc = fn
}

// SetChatFunc sets a custom function for Chat calls.
func (m *MockProvider) SetChatFunc(fn func(ctx context.Context, prompt string) (string, error)) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.chatFunc = fn
}

// SetAvailable sets whether the mock provider is available.
func (m *MockProvider) SetAvailable(available bool) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.available = available
}

// GetCalls returns all recorded calls.
func (m *MockProvider) GetCalls() []MockCall {
	m.mu.Lock()
	defer m.mu.Unlock()
	calls := make([]MockCall, len(m.calls))
	copy(calls, m.calls)
	return calls
}

// ClearCalls clears all recorded calls.
func (m *MockProvider) ClearCalls() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.calls = make([]MockCall, 0)
}

// Reset resets the mock provider to its initial state.
func (m *MockProvider) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.responses = make(map[string]json.RawMessage)
	m.chatResponses = make(map[string]string)
	m.analyzeFunc = nil
	m.chatFunc = nil
	m.available = true
	m.calls = make([]MockCall, 0)
}

// WithDependencyResponse configures a standard response for dependency detection.
func (m *MockProvider) WithDependencyResponse(dependencies []map[string]any) *MockProvider {
	data, err := json.Marshal(map[string]any{
		"dependencies": dependencies,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal dependency response: %v", err))
	}

	m.SetAnalyzeFunc(func(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
		return data, nil
	})
	return m
}

// WithServiceResponse configures a standard response for service detection.
func (m *MockProvider) WithServiceResponse(services []map[string]any) *MockProvider {
	data, err := json.Marshal(map[string]any{
		"services": services,
	})
	if err != nil {
		panic(fmt.Sprintf("failed to marshal service response: %v", err))
	}

	m.SetAnalyzeFunc(func(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
		return data, nil
	})
	return m
}

// WithError configures the mock to return an error.
func (m *MockProvider) WithError(err error) *MockProvider {
	m.SetAnalyzeFunc(func(_ context.Context, _ string, _ json.RawMessage) (json.RawMessage, error) {
		return nil, err
	})
	m.SetChatFunc(func(_ context.Context, _ string) (string, error) {
		return "", err
	})
	return m
}
