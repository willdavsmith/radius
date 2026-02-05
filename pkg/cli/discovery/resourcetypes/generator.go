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

package resourcetypes

import (
	"context"
	"encoding/json"

	"github.com/radius-project/radius/pkg/cli/discovery"
	"github.com/radius-project/radius/pkg/cli/discovery/llm"
)

// Generator is a skill that generates Resource Types from detected dependencies.
type Generator struct {
	catalog     *Catalog
	llmProvider llm.Provider
}

// GeneratorInput is the input for the resource type generator.
type GeneratorInput struct {
	// Dependencies are the detected dependencies.
	Dependencies []discovery.DetectedDependency `json:"dependencies"`
}

// GeneratorOutput is the output from the resource type generator.
type GeneratorOutput struct {
	// ResourceTypes are the generated resource types.
	ResourceTypes []discovery.ResourceType `json:"resourceTypes"`

	// UnmappedTypes lists dependency types that couldn't be mapped.
	UnmappedTypes []string `json:"unmappedTypes,omitempty"`

	// Warnings lists any warnings encountered.
	Warnings []string `json:"warnings,omitempty"`
}

// NewGenerator creates a new resource type generator.
func NewGenerator(llmProvider llm.Provider) *Generator {
	return &Generator{
		catalog:     NewCatalog(),
		llmProvider: llmProvider,
	}
}

// Name returns the skill name.
func (g *Generator) Name() string {
	return "generate_resource_types"
}

// Description returns the skill description.
func (g *Generator) Description() string {
	return "Maps detected infrastructure dependencies to Radius Resource Type definitions."
}

// InputSchema returns the JSON Schema for the input.
func (g *Generator) InputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"dependencies": {
				"type": "array",
				"description": "List of detected dependencies to map",
				"items": {
					"type": "object",
					"properties": {
						"type": {"type": "string"},
						"technology": {"type": "string"},
						"confidence": {"type": "number"}
					}
				}
			}
		},
		"required": ["dependencies"]
	}`)
}

// OutputSchema returns the JSON Schema for the output.
func (g *Generator) OutputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"resourceTypes": {
				"type": "array",
				"items": {"type": "object"},
				"description": "Generated Resource Type definitions"
			},
			"unmappedTypes": {
				"type": "array",
				"items": {"type": "string"},
				"description": "Dependency types that couldn't be mapped"
			},
			"warnings": {
				"type": "array",
				"items": {"type": "string"},
				"description": "Warnings encountered during generation"
			}
		}
	}`)
}

// Execute runs the resource type generation.
func (g *Generator) Execute(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	var in GeneratorInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}

	output := GeneratorOutput{
		ResourceTypes: make([]discovery.ResourceType, 0),
		UnmappedTypes: make([]string, 0),
		Warnings:      make([]string, 0),
	}

	seen := make(map[string]bool)
	for _, dep := range in.Dependencies {
		if seen[dep.Type] {
			continue
		}
		seen[dep.Type] = true

		entry, found := g.catalog.Get(dep.Type)
		if found {
			output.ResourceTypes = append(output.ResourceTypes, entry.ResourceType)
		} else {
			fallbackType := g.generateFallback(dep)
			if fallbackType != nil {
				output.ResourceTypes = append(output.ResourceTypes, *fallbackType)
				output.Warnings = append(output.Warnings,
					"Generated fallback Resource Type for "+dep.Type+". Consider defining a proper type.")
			} else {
				output.UnmappedTypes = append(output.UnmappedTypes, dep.Type)
			}
		}
	}

	return json.Marshal(output)
}

// generateFallback creates a minimal Resource Type for an unknown dependency.
func (g *Generator) generateFallback(dep discovery.DetectedDependency) *discovery.ResourceType {
	return &discovery.ResourceType{
		Name:       dep.Technology,
		Type:       "Applications.Core/extenders",
		APIVersion: "2023-10-01-preview",
		Maturity:   "alpha",
		Schema: &discovery.ResourceSchema{
			Properties: map[string]discovery.PropertyDef{
				"resourceType": {
					Type:        "string",
					Description: "The underlying resource type",
					Default:     dep.Type,
				},
			},
		},
		Outputs: []discovery.OutputDefinition{
			{Name: "connectionString", Type: "string", Description: "Connection string for the resource", Sensitive: true},
		},
	}
}
