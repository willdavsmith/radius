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

package skills

import (
	"context"
	"encoding/json"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type mockSkill struct {
	name         string
	description  string
	inputSchema  json.RawMessage
	outputSchema json.RawMessage
	executeFunc  func(ctx context.Context, input json.RawMessage) (json.RawMessage, error)
}

func (m *mockSkill) Name() string                  { return m.name }
func (m *mockSkill) Description() string           { return m.description }
func (m *mockSkill) InputSchema() json.RawMessage  { return m.inputSchema }
func (m *mockSkill) OutputSchema() json.RawMessage { return m.outputSchema }

func (m *mockSkill) Execute(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	if m.executeFunc != nil {
		return m.executeFunc(ctx, input)
	}
	return json.RawMessage(`{"result": "ok"}`), nil
}

func newMockSkill(name string) *mockSkill {
	return &mockSkill{
		name:         name,
		description:  "A mock skill for testing",
		inputSchema:  json.RawMessage(`{"type": "object"}`),
		outputSchema: json.RawMessage(`{"type": "object"}`),
	}
}

func TestNewRegistry(t *testing.T) {
	registry := NewRegistry()
	require.NotNil(t, registry)
	assert.Equal(t, 0, registry.Count())
}

func TestRegistry_Register(t *testing.T) {
	tests := []struct {
		name        string
		skillName   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid skill name",
			skillName:   "discover_dependencies",
			expectError: false,
		},
		{
			name:        "valid skill with numbers",
			skillName:   "skill1",
			expectError: false,
		},
		{
			name:        "valid skill with underscores",
			skillName:   "my_skill_name",
			expectError: false,
		},
		{
			name:        "invalid skill starts with number",
			skillName:   "1skill",
			expectError: true,
			errorMsg:    "invalid skill name",
		},
		{
			name:        "invalid skill with hyphen",
			skillName:   "my-skill",
			expectError: true,
			errorMsg:    "invalid skill name",
		},
		{
			name:        "invalid skill with uppercase",
			skillName:   "MySkill",
			expectError: true,
			errorMsg:    "invalid skill name",
		},
		{
			name:        "empty skill name",
			skillName:   "",
			expectError: true,
			errorMsg:    "invalid skill name",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			registry := NewRegistry()
			skill := newMockSkill(tt.skillName)

			err := registry.Register(skill)

			if tt.expectError {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errorMsg)
			} else {
				require.NoError(t, err)
				assert.Equal(t, 1, registry.Count())
			}
		})
	}
}

func TestRegistry_RegisterNil(t *testing.T) {
	registry := NewRegistry()
	err := registry.Register(nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot be nil")
}

func TestRegistry_RegisterDuplicate(t *testing.T) {
	registry := NewRegistry()
	skill := newMockSkill("test_skill")

	err := registry.Register(skill)
	require.NoError(t, err)

	err = registry.Register(skill)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already registered")
}

func TestRegistry_Get(t *testing.T) {
	registry := NewRegistry()
	skill := newMockSkill("test_skill")
	err := registry.Register(skill)
	require.NoError(t, err)

	result := registry.Get("test_skill")
	require.NotNil(t, result)
	assert.Equal(t, "test_skill", result.Name())

	result = registry.Get("nonexistent")
	assert.Nil(t, result)
}

func TestRegistry_List(t *testing.T) {
	registry := NewRegistry()

	skills := []string{"zebra_skill", "alpha_skill", "middle_skill"}
	for _, name := range skills {
		err := registry.Register(newMockSkill(name))
		require.NoError(t, err)
	}

	list := registry.List()
	require.Len(t, list, 3)
	assert.Equal(t, "alpha_skill", list[0].Name())
	assert.Equal(t, "middle_skill", list[1].Name())
	assert.Equal(t, "zebra_skill", list[2].Name())
}

func TestRegistry_ListMetadata(t *testing.T) {
	registry := NewRegistry()
	skill := &mockSkill{
		name:         "test_skill",
		description:  "Test description",
		inputSchema:  json.RawMessage(`{"type": "object"}`),
		outputSchema: json.RawMessage(`{"type": "object"}`),
	}
	err := registry.Register(skill)
	require.NoError(t, err)

	metadata := registry.ListMetadata()
	require.Len(t, metadata, 1)
	assert.Equal(t, "test_skill", metadata[0].Name)
	assert.Equal(t, "Test description", metadata[0].Description)
}

func TestRegistry_Execute(t *testing.T) {
	registry := NewRegistry()
	expectedOutput := json.RawMessage(`{"status": "success"}`)

	skill := &mockSkill{
		name: "test_skill",
		executeFunc: func(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
			return expectedOutput, nil
		},
	}
	err := registry.Register(skill)
	require.NoError(t, err)

	result, err := registry.Execute(context.Background(), "test_skill", nil)
	require.NoError(t, err)
	assert.JSONEq(t, string(expectedOutput), string(result))

	_, err = registry.Execute(context.Background(), "nonexistent", nil)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRegistry_ExecuteWithError(t *testing.T) {
	registry := NewRegistry()
	expectedError := errors.New("skill execution failed")

	skill := &mockSkill{
		name: "failing_skill",
		executeFunc: func(_ context.Context, _ json.RawMessage) (json.RawMessage, error) {
			return nil, expectedError
		},
	}
	err := registry.Register(skill)
	require.NoError(t, err)

	_, err = registry.Execute(context.Background(), "failing_skill", nil)
	require.Error(t, err)
	assert.Equal(t, expectedError, err)
}

func TestRegistry_Unregister(t *testing.T) {
	registry := NewRegistry()
	skill := newMockSkill("test_skill")
	err := registry.Register(skill)
	require.NoError(t, err)
	assert.Equal(t, 1, registry.Count())

	err = registry.Unregister("test_skill")
	require.NoError(t, err)
	assert.Equal(t, 0, registry.Count())

	err = registry.Unregister("nonexistent")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestRegistry_Clear(t *testing.T) {
	registry := NewRegistry()

	for i := 0; i < 5; i++ {
		skill := newMockSkill("skill" + string(rune('a'+i)))
		err := registry.Register(skill)
		require.NoError(t, err)
	}
	assert.Equal(t, 5, registry.Count())

	registry.Clear()
	assert.Equal(t, 0, registry.Count())
}

func TestDefaultExecutionConfig(t *testing.T) {
	config := DefaultExecutionConfig()
	require.NotNil(t, config)

	assert.Equal(t, 0.5, config.ConfidenceThreshold)
	assert.Equal(t, 1000, config.MaxFiles)
	assert.Equal(t, "openai", config.LLMProvider)
	assert.Contains(t, config.ExcludePaths, "node_modules")
	assert.Contains(t, config.ExcludePaths, ".git")
}
