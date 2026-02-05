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
	"fmt"
	"time"

	"github.com/radius-project/radius/pkg/cli/discovery"
)

// AnalyzeCodebase is a high-level skill that orchestrates complete codebase analysis.
// It combines dependency detection, service detection, and resource type mapping
// into a single operation for AI agents.
type AnalyzeCodebase struct {
	registry *Registry
}

// NewAnalyzeCodebase creates a new analyze codebase skill.
func NewAnalyzeCodebase(registry *Registry) *AnalyzeCodebase {
	return &AnalyzeCodebase{
		registry: registry,
	}
}

// Name returns the skill identifier.
func (a *AnalyzeCodebase) Name() string {
	return "analyze_codebase"
}

// Description returns a human-readable description.
func (a *AnalyzeCodebase) Description() string {
	return `Performs a comprehensive analysis of a codebase to discover infrastructure dependencies, 
services, and generate Radius resource mappings. This is the primary entry point for AI agents 
to understand an application's infrastructure requirements.`
}

// InputSchema returns the JSON Schema for the skill input.
func (a *AnalyzeCodebase) InputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"path": {
				"type": "string",
				"description": "Root path of the codebase to analyze"
			},
			"excludePaths": {
				"type": "array",
				"items": {"type": "string"},
				"description": "Paths to exclude from analysis (e.g., node_modules, vendor)"
			},
			"confidenceThreshold": {
				"type": "number",
				"description": "Minimum confidence score for detections (0.0-1.0)",
				"default": 0.5
			},
			"includeResourceTypes": {
				"type": "boolean",
				"description": "Whether to map dependencies to Radius resource types",
				"default": true
			}
		},
		"required": ["path"]
	}`)
}

// OutputSchema returns the JSON Schema for the skill output.
func (a *AnalyzeCodebase) OutputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"summary": {
				"type": "object",
				"description": "High-level summary of the analysis",
				"properties": {
					"dependencyCount": {"type": "integer"},
					"serviceCount": {"type": "integer"},
					"mappedResourceTypes": {"type": "integer"},
					"analysisTime": {"type": "string"}
				}
			},
			"dependencies": {
				"type": "array",
				"items": {"type": "object"},
				"description": "Detected infrastructure dependencies"
			},
			"services": {
				"type": "array",
				"items": {"type": "object"},
				"description": "Detected application services"
			},
			"resourceTypes": {
				"type": "array",
				"items": {"type": "object"},
				"description": "Mapped Radius resource types"
			},
			"recommendations": {
				"type": "array",
				"items": {"type": "string"},
				"description": "AI-generated recommendations for deployment"
			}
		}
	}`)
}

// AnalyzeCodebaseInput defines the input parameters.
type AnalyzeCodebaseInput struct {
	Path                 string   `json:"path"`
	ExcludePaths         []string `json:"excludePaths,omitempty"`
	ConfidenceThreshold  float64  `json:"confidenceThreshold,omitempty"`
	IncludeResourceTypes bool     `json:"includeResourceTypes,omitempty"`
}

// AnalyzeCodebaseOutput defines the output structure.
type AnalyzeCodebaseOutput struct {
	Summary         AnalysisSummary                `json:"summary"`
	Dependencies    []discovery.DetectedDependency `json:"dependencies"`
	Services        []discovery.DetectedService    `json:"services"`
	ResourceTypes   []discovery.ResourceType       `json:"resourceTypes,omitempty"`
	Recommendations []string                       `json:"recommendations"`
}

// AnalysisSummary provides high-level metrics.
type AnalysisSummary struct {
	DependencyCount     int    `json:"dependencyCount"`
	ServiceCount        int    `json:"serviceCount"`
	MappedResourceTypes int    `json:"mappedResourceTypes"`
	AnalysisTime        string `json:"analysisTime"`
}

// Execute runs the comprehensive codebase analysis.
func (a *AnalyzeCodebase) Execute(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	startTime := time.Now()

	var in AnalyzeCodebaseInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, discovery.WrapError(discovery.ErrCodeInputValidation, "invalid input", err)
	}

	// Set defaults
	if in.ConfidenceThreshold == 0 {
		in.ConfidenceThreshold = 0.5
	}
	if !in.IncludeResourceTypes {
		in.IncludeResourceTypes = true
	}

	output := AnalyzeCodebaseOutput{
		Dependencies:    make([]discovery.DetectedDependency, 0),
		Services:        make([]discovery.DetectedService, 0),
		ResourceTypes:   make([]discovery.ResourceType, 0),
		Recommendations: make([]string, 0),
	}

	// Step 1: Detect dependencies
	depInput := map[string]any{
		"path":                in.Path,
		"excludePaths":        in.ExcludePaths,
		"confidenceThreshold": in.ConfidenceThreshold,
	}
	depInputJSON, _ := json.Marshal(depInput)

	depResult, err := a.registry.Execute(ctx, "discover_dependencies", depInputJSON)
	if err == nil {
		var depOutput struct {
			Dependencies []discovery.DetectedDependency `json:"dependencies"`
		}
		if err := json.Unmarshal(depResult, &depOutput); err == nil {
			output.Dependencies = depOutput.Dependencies
		}
	}

	// Step 2: Detect services
	svcInput := map[string]any{
		"path":         in.Path,
		"excludePaths": in.ExcludePaths,
	}
	svcInputJSON, _ := json.Marshal(svcInput)

	svcResult, err := a.registry.Execute(ctx, "discover_services", svcInputJSON)
	if err == nil {
		var svcOutput struct {
			Services []discovery.DetectedService `json:"services"`
		}
		if err := json.Unmarshal(svcResult, &svcOutput); err == nil {
			output.Services = svcOutput.Services
		}
	}

	// Step 3: Map to resource types (if requested)
	if in.IncludeResourceTypes && len(output.Dependencies) > 0 {
		rtInput := map[string]any{
			"dependencies": output.Dependencies,
		}
		rtInputJSON, _ := json.Marshal(rtInput)

		rtResult, err := a.registry.Execute(ctx, "generate_resource_types", rtInputJSON)
		if err == nil {
			var rtOutput struct {
				ResourceTypes []discovery.ResourceType `json:"resourceTypes"`
			}
			if err := json.Unmarshal(rtResult, &rtOutput); err == nil {
				output.ResourceTypes = rtOutput.ResourceTypes
			}
		}
	}

	// Generate recommendations
	output.Recommendations = a.generateRecommendations(output)

	// Build summary
	output.Summary = AnalysisSummary{
		DependencyCount:     len(output.Dependencies),
		ServiceCount:        len(output.Services),
		MappedResourceTypes: len(output.ResourceTypes),
		AnalysisTime:        fmt.Sprintf("%dms", time.Since(startTime).Milliseconds()),
	}

	return json.Marshal(output)
}

// generateRecommendations creates AI-friendly recommendations based on analysis.
func (a *AnalyzeCodebase) generateRecommendations(output AnalyzeCodebaseOutput) []string {
	recommendations := make([]string, 0)

	if len(output.Dependencies) == 0 && len(output.Services) == 0 {
		recommendations = append(recommendations,
			"No infrastructure dependencies or services were detected. Ensure the codebase path is correct and contains application code.")
		return recommendations
	}

	// Database recommendations
	dbCount := 0
	for _, dep := range output.Dependencies {
		switch dep.Type {
		case "postgresql", "mysql", "mongodb", "cosmosdb", "dynamodb":
			dbCount++
		}
	}
	if dbCount > 1 {
		recommendations = append(recommendations,
			fmt.Sprintf("Multiple database types detected (%d). Consider consolidating to a single database technology for simpler operations.", dbCount))
	}

	// Caching recommendations
	hasCache := false
	for _, dep := range output.Dependencies {
		if dep.Type == "redis" || dep.Type == "memcached" {
			hasCache = true
			break
		}
	}
	if !hasCache && len(output.Services) > 0 {
		recommendations = append(recommendations,
			"No caching layer detected. Consider adding Redis for improved performance in production.")
	}

	// Service mesh recommendations
	if len(output.Services) > 2 {
		recommendations = append(recommendations,
			fmt.Sprintf("Multiple services detected (%d). The generated Radius app will include connections between services and their dependencies.", len(output.Services)))
	}

	// Recipe recommendations
	for _, dep := range output.Dependencies {
		recommendations = append(recommendations,
			fmt.Sprintf("Dependency '%s' will use a Radius recipe for provisioning. Ensure the 'default' recipe is configured in your environment.", dep.Type))
	}

	// Container image recommendations
	if len(output.Services) > 0 {
		recommendations = append(recommendations,
			"Container images in the generated Bicep will have placeholder values (TODO). Update these with your actual container registry paths before deployment.")
	}

	return recommendations
}
