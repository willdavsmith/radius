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

// Package skills defines the skill interface and registry for the discovery feature.
package skills

import (
	"context"
	"encoding/json"
)

// Skill defines a composable capability that can be invoked via CLI or MCP.
type Skill interface {
	// Name returns the unique skill identifier.
	Name() string

	// Description returns a human-readable description.
	Description() string

	// InputSchema returns the JSON Schema defining valid input.
	InputSchema() json.RawMessage

	// OutputSchema returns the JSON Schema defining output.
	OutputSchema() json.RawMessage

	// Execute runs the skill with the given input.
	Execute(ctx context.Context, input json.RawMessage) (json.RawMessage, error)
}

// SkillMetadata contains metadata about a skill.
type SkillMetadata struct {
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	Category     SkillCategory   `json:"category"`
	InputSchema  json.RawMessage `json:"inputSchema"`
	OutputSchema json.RawMessage `json:"outputSchema"`
}

// SkillCategory groups related skills.
type SkillCategory string

const (
	// CategoryDiscovery includes skills that analyze codebases.
	CategoryDiscovery SkillCategory = "discovery"

	// CategoryGeneration includes skills that generate artifacts.
	CategoryGeneration SkillCategory = "generation"

	// CategoryValidation includes skills that validate artifacts.
	CategoryValidation SkillCategory = "validation"
)

// ExecutionContext provides context for skill execution.
type ExecutionContext struct {
	WorkingDir       string
	Config           *ExecutionConfig
	ProgressCallback func(message string, percentComplete float64)
}

// ExecutionConfig contains configuration for skill execution.
type ExecutionConfig struct {
	ConfidenceThreshold float64  `json:"confidenceThreshold"`
	ExcludePaths        []string `json:"excludePaths"`
	MaxFiles            int      `json:"maxFiles"`
	LLMProvider         string   `json:"llmProvider"`
}

// DefaultExecutionConfig returns the default configuration.
func DefaultExecutionConfig() *ExecutionConfig {
	return &ExecutionConfig{
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
		},
		MaxFiles:    1000,
		LLMProvider: "openai",
	}
}
