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
	"os"
	"path/filepath"

	"github.com/radius-project/radius/pkg/cli/discovery"
)

// GenerateDeployment is a high-level skill that generates complete deployment artifacts.
// It combines discovery and generation into a single operation that produces ready-to-deploy files.
type GenerateDeployment struct {
	registry *Registry
}

// NewGenerateDeployment creates a new generate deployment skill.
func NewGenerateDeployment(registry *Registry) *GenerateDeployment {
	return &GenerateDeployment{
		registry: registry,
	}
}

// Name returns the skill identifier.
func (g *GenerateDeployment) Name() string {
	return "generate_deployment"
}

// Description returns a human-readable description.
func (g *GenerateDeployment) Description() string {
	return `Generates complete Radius deployment artifacts from a codebase.
Performs end-to-end analysis and generates app.bicep and supporting files
ready for deployment with 'rad deploy'.`
}

// InputSchema returns the JSON Schema for the skill input.
func (g *GenerateDeployment) InputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"codebasePath": {
				"type": "string",
				"description": "Root path of the codebase to analyze"
			},
			"appName": {
				"type": "string",
				"description": "Name for the Radius application"
			},
			"outputPath": {
				"type": "string",
				"description": "Directory to write generated files (default: codebase root)"
			},
			"containerImages": {
				"type": "object",
				"description": "Map of service name to container image URI",
				"additionalProperties": {"type": "string"}
			},
			"dryRun": {
				"type": "boolean",
				"description": "If true, returns generated content without writing files",
				"default": false
			}
		},
		"required": ["codebasePath"]
	}`)
}

// OutputSchema returns the JSON Schema for the skill output.
func (g *GenerateDeployment) OutputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"success": {"type": "boolean"},
			"appName": {"type": "string"},
			"generatedFiles": {
				"type": "array",
				"items": {"type": "string"}
			},
			"bicepContent": {"type": "string"},
			"summary": {
				"type": "object",
				"properties": {
					"containers": {"type": "integer"},
					"resources": {"type": "integer"},
					"connections": {"type": "integer"}
				}
			},
			"deployCommand": {"type": "string"},
			"warnings": {"type": "array", "items": {"type": "string"}}
		}
	}`)
}

// GenerateDeploymentInput defines the input parameters.
type GenerateDeploymentInput struct {
	CodebasePath    string            `json:"codebasePath"`
	AppName         string            `json:"appName,omitempty"`
	OutputPath      string            `json:"outputPath,omitempty"`
	ContainerImages map[string]string `json:"containerImages,omitempty"`
	DryRun          bool              `json:"dryRun,omitempty"`
}

// GenerateDeploymentOutput defines the output structure.
type GenerateDeploymentOutput struct {
	Success        bool              `json:"success"`
	AppName        string            `json:"appName"`
	GeneratedFiles []string          `json:"generatedFiles"`
	BicepContent   string            `json:"bicepContent"`
	Summary        DeploymentSummary `json:"summary"`
	DeployCommand  string            `json:"deployCommand"`
	Warnings       []string          `json:"warnings"`
}

// DeploymentSummary provides counts of generated resources.
type DeploymentSummary struct {
	Containers  int `json:"containers"`
	Resources   int `json:"resources"`
	Connections int `json:"connections"`
}

// Execute runs the deployment generation.
func (g *GenerateDeployment) Execute(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	var in GenerateDeploymentInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, discovery.WrapError(discovery.ErrCodeInputValidation, "invalid input", err)
	}

	output := GenerateDeploymentOutput{
		Success:        false,
		GeneratedFiles: make([]string, 0),
		Warnings:       make([]string, 0),
	}

	// Set defaults
	if in.AppName == "" {
		in.AppName = filepath.Base(in.CodebasePath)
	}
	if in.OutputPath == "" {
		in.OutputPath = in.CodebasePath
	}

	output.AppName = in.AppName

	// Step 1: Run comprehensive analysis
	analyzeInput := map[string]any{
		"path":                 in.CodebasePath,
		"includeResourceTypes": true,
	}
	analyzeInputJSON, _ := json.Marshal(analyzeInput)

	analyzeResult, err := g.registry.Execute(ctx, "analyze_codebase", analyzeInputJSON)
	if err != nil {
		return nil, discovery.WrapError(discovery.ErrCodeInternal, "analysis failed", err)
	}

	var analysisOutput AnalyzeCodebaseOutput
	if err := json.Unmarshal(analyzeResult, &analysisOutput); err != nil {
		return nil, discovery.WrapError(discovery.ErrCodeInternal, "failed to parse analysis results", err)
	}

	// Step 2: Generate Bicep
	genInput := map[string]any{
		"appName":      in.AppName,
		"dependencies": analysisOutput.Dependencies,
		"services":     analysisOutput.Services,
	}
	genInputJSON, _ := json.Marshal(genInput)

	genResult, err := g.registry.Execute(ctx, "generate_app_definition", genInputJSON)
	if err != nil {
		return nil, discovery.WrapError(discovery.ErrCodeInternal, "generation failed", err)
	}

	var genOutput struct {
		BicepContent  string                  `json:"bicepContent"`
		AppDefinition discovery.AppDefinition `json:"appDefinition"`
		Warnings      []string                `json:"warnings"`
	}
	if err := json.Unmarshal(genResult, &genOutput); err != nil {
		return nil, discovery.WrapError(discovery.ErrCodeInternal, "failed to parse generation results", err)
	}

	// Apply container image overrides
	bicepContent := genOutput.BicepContent
	if len(in.ContainerImages) > 0 {
		for serviceName, image := range in.ContainerImages {
			placeholder := "TODO: " + serviceName
			if contains(bicepContent, placeholder) {
				bicepContent = replaceAll(bicepContent, placeholder, image)
			}
		}
	}

	// Add location: 'global' to all resources (fix for Radius deployment)
	bicepContent = g.addLocationToResources(bicepContent)

	output.BicepContent = bicepContent
	output.Summary = DeploymentSummary{
		Containers:  len(genOutput.AppDefinition.Containers),
		Resources:   len(genOutput.AppDefinition.Resources),
		Connections: len(genOutput.AppDefinition.Connections),
	}
	output.Warnings = append(output.Warnings, genOutput.Warnings...)

	// Check for remaining TODOs
	if contains(bicepContent, "TODO:") {
		output.Warnings = append(output.Warnings,
			"Generated Bicep contains TODO placeholders. Update container image references before deployment.")
	}

	// Write files if not dry run
	if !in.DryRun {
		bicepPath := filepath.Join(in.OutputPath, "app.bicep")
		if err := os.WriteFile(bicepPath, []byte(bicepContent), 0644); err != nil {
			return nil, discovery.WrapError(discovery.ErrCodeFileNotFound, "failed to write app.bicep", err)
		}
		output.GeneratedFiles = append(output.GeneratedFiles, bicepPath)

		// Copy bicepconfig.json if it doesn't exist
		bicepConfigPath := filepath.Join(in.OutputPath, "bicepconfig.json")
		if _, err := os.Stat(bicepConfigPath); os.IsNotExist(err) {
			bicepConfig := `{
	"experimentalFeaturesEnabled": {
		"extensibility": true
	},
	"extensions": {
		"radius": "br:biceptypes.azurecr.io/radius:latest",
		"aws": "br:biceptypes.azurecr.io/aws:latest"
	}
}`
			if err := os.WriteFile(bicepConfigPath, []byte(bicepConfig), 0644); err == nil {
				output.GeneratedFiles = append(output.GeneratedFiles, bicepConfigPath)
			}
		}
	}

	output.Success = true
	output.DeployCommand = "rad deploy app.bicep -a " + in.AppName

	return json.Marshal(output)
}

// addLocationToResources adds location: 'global' to resource definitions that are missing it.
func (g *GenerateDeployment) addLocationToResources(content string) string {
	// This is a simple approach - in production you'd use proper Bicep parsing
	// Add location after "name:" lines for resources that don't have it
	lines := splitLines(content)
	result := make([]string, 0, len(lines))

	inResource := false
	hasLocation := false
	resourceIndent := ""

	for i, line := range lines {
		result = append(result, line)

		// Detect resource block start
		if contains(line, "resource ") && contains(line, " = {") {
			inResource = true
			hasLocation = false
			// Determine indentation
			for _, c := range line {
				if c == ' ' || c == '\t' {
					resourceIndent += string(c)
				} else {
					break
				}
			}
		}

		// Check if location already exists
		if inResource && contains(line, "location:") {
			hasLocation = true
		}

		// After name: line, add location if needed
		if inResource && contains(line, "name:") && !hasLocation {
			// Look ahead to see if location is on next line
			if i+1 < len(lines) && !contains(lines[i+1], "location:") {
				indent := resourceIndent + "  "
				result = append(result, indent+"location: 'global'")
				hasLocation = true
			}
		}

		// Detect resource block end
		if inResource && line == resourceIndent+"}" {
			inResource = false
			resourceIndent = ""
		}
	}

	return joinLines(result)
}

func splitLines(s string) []string {
	lines := make([]string, 0)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func joinLines(lines []string) string {
	result := ""
	for i, line := range lines {
		if i > 0 {
			result += "\n"
		}
		result += line
	}
	return result
}

func replaceAll(s, old, new string) string {
	result := ""
	for {
		idx := indexOf(s, old)
		if idx == -1 {
			result += s
			break
		}
		result += s[:idx] + new
		s = s[idx+len(old):]
	}
	return result
}

func indexOf(s, substr string) int {
	for i := 0; i+len(substr) <= len(s); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}

func contains(s, substr string) bool {
	return indexOf(s, substr) != -1
}
