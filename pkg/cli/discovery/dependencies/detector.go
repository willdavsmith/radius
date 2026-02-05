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

// Package dependencies provides the dependency detection skill.
package dependencies

import (
	"context"
	"encoding/json"

	"github.com/radius-project/radius/pkg/cli/discovery"
	"github.com/radius-project/radius/pkg/cli/discovery/llm"
)

// Detector is the skill that detects infrastructure dependencies from a codebase.
type Detector struct {
	llmProvider llm.Provider
	walker      *Walker
}

// NewDetector creates a new dependency detector.
func NewDetector(llmProvider llm.Provider) *Detector {
	return &Detector{
		llmProvider: llmProvider,
		walker:      NewWalker(nil),
	}
}

// Name returns the skill identifier.
func (d *Detector) Name() string {
	return "discover_dependencies"
}

// Description returns a human-readable description.
func (d *Detector) Description() string {
	return "Analyzes a codebase to detect infrastructure dependencies such as databases, message queues, and external services"
}

// InputSchema returns the JSON Schema for the skill input.
func (d *Detector) InputSchema() json.RawMessage {
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
				"description": "Paths to exclude from analysis"
			},
			"confidenceThreshold": {
				"type": "number",
				"description": "Minimum confidence score for detections (0.0-1.0)",
				"default": 0.5
			}
		},
		"required": ["path"]
	}`)
}

// OutputSchema returns the JSON Schema for the skill output.
func (d *Detector) OutputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"dependencies": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"type": {"type": "string"},
						"technology": {"type": "string"},
						"confidence": {"type": "number"},
						"evidence": {"type": "array", "items": {"type": "object"}}
					}
				}
			}
		}
	}`)
}

// DetectorInput is the input for the dependency detector.
type DetectorInput struct {
	Path                string   `json:"path"`
	ExcludePaths        []string `json:"excludePaths,omitempty"`
	ConfidenceThreshold float64  `json:"confidenceThreshold,omitempty"`
}

// DetectorOutput is the output of the dependency detector.
type DetectorOutput struct {
	Dependencies []discovery.DetectedDependency `json:"dependencies"`
}

// Execute runs the dependency detection.
func (d *Detector) Execute(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	var in DetectorInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, discovery.WrapError(discovery.ErrCodeInputValidation, "invalid input", err)
	}

	if in.ConfidenceThreshold == 0 {
		in.ConfidenceThreshold = 0.5
	}

	// Update walker exclusions
	d.walker = NewWalker(in.ExcludePaths)

	// Detect dependencies
	dependencies, err := d.detectDependencies(ctx, in.Path, in.ConfidenceThreshold)
	if err != nil {
		return nil, err
	}

	output := DetectorOutput{
		Dependencies: dependencies,
	}
	return json.Marshal(output)
}

// detectDependencies analyzes the codebase for infrastructure dependencies.
func (d *Detector) detectDependencies(ctx context.Context, path string, threshold float64) ([]discovery.DetectedDependency, error) {
	// Collect relevant files
	files, err := d.walker.WalkDependencyFiles(path)
	if err != nil {
		return nil, discovery.WrapError(discovery.ErrCodeFileNotFound, "failed to walk directory", err)
	}

	// Analyze each file
	var allDeps []discovery.DetectedDependency
	for _, file := range files {
		deps, err := d.analyzeFile(ctx, file)
		if err != nil {
			// Log warning but continue with other files
			continue
		}
		allDeps = append(allDeps, deps...)
	}

	// Filter by confidence threshold and deduplicate
	filtered := d.filterAndDeduplicate(allDeps, threshold)
	return filtered, nil
}

// analyzeFile uses the LLM to analyze a single file for dependencies.
func (d *Detector) analyzeFile(ctx context.Context, file FileInfo) ([]discovery.DetectedDependency, error) {
	if d.llmProvider == nil {
		// Return pattern-based detection if no LLM is available
		return d.patternBasedDetection(file), nil
	}

	prompt := d.buildPrompt(file)
	schema := d.OutputSchema()

	result, err := d.llmProvider.Analyze(ctx, prompt, schema)
	if err != nil {
		// Fall back to pattern-based detection
		return d.patternBasedDetection(file), nil
	}

	var output DetectorOutput
	if err := json.Unmarshal(result, &output); err != nil {
		return d.patternBasedDetection(file), nil
	}

	// Add file source to evidence
	for i := range output.Dependencies {
		for j := range output.Dependencies[i].Evidence {
			if output.Dependencies[i].Evidence[j].Source == "" {
				output.Dependencies[i].Evidence[j].Source = file.RelativePath
			}
		}
	}

	return output.Dependencies, nil
}

// buildPrompt creates the LLM prompt for dependency detection.
func (d *Detector) buildPrompt(file FileInfo) string {
	return "Analyze this file for infrastructure dependencies.\n\n" +
		"FILE PATH: " + file.RelativePath + "\n\n" +
		"CONTENT:\n" + file.Content
}

// patternBasedDetection performs simple pattern matching for common dependencies.
func (d *Detector) patternBasedDetection(file FileInfo) []discovery.DetectedDependency {
	// Match common patterns in the file content
	patterns := map[string]struct {
		Type       string
		Technology string
	}{
		"pg":            {"postgresql", "PostgreSQL"},
		"postgres":      {"postgresql", "PostgreSQL"},
		"mysql":         {"mysql", "MySQL"},
		"mongodb":       {"mongodb", "MongoDB"},
		"redis":         {"redis", "Redis"},
		"rabbitmq":      {"rabbitmq", "RabbitMQ"},
		"kafka":         {"kafka", "Apache Kafka"},
		"elasticsearch": {"elasticsearch", "Elasticsearch"},
		"dynamodb":      {"dynamodb", "Amazon DynamoDB"},
		"cosmosdb":      {"cosmosdb", "Azure CosmosDB"},
	}

	var deps []discovery.DetectedDependency
	for pattern, dep := range patterns {
		if containsPattern(file.Content, pattern) {
			deps = append(deps, discovery.DetectedDependency{
				Type:       dep.Type,
				Technology: dep.Technology,
				Confidence: 0.7,
				Evidence: []discovery.Evidence{
					{
						Source: file.RelativePath,
						Type:   "pattern",
						Value:  pattern,
					},
				},
			})
		}
	}

	return deps
}

// filterAndDeduplicate removes low-confidence and duplicate detections.
func (d *Detector) filterAndDeduplicate(deps []discovery.DetectedDependency, threshold float64) []discovery.DetectedDependency {
	seen := make(map[string]*discovery.DetectedDependency)

	for _, dep := range deps {
		if dep.Confidence < threshold {
			continue
		}

		if existing, ok := seen[dep.Type]; ok {
			// Keep the one with higher confidence
			if dep.Confidence > existing.Confidence {
				seen[dep.Type] = &dep
			}
			// Merge evidence
			existing.Evidence = append(existing.Evidence, dep.Evidence...)
		} else {
			depCopy := dep
			seen[dep.Type] = &depCopy
		}
	}

	result := make([]discovery.DetectedDependency, 0, len(seen))
	for _, dep := range seen {
		result = append(result, *dep)
	}
	return result
}

// containsPattern checks if content contains a pattern (case-insensitive).
func containsPattern(content, pattern string) bool {
	return len(content) > 0 && len(pattern) > 0 &&
		(contains(content, pattern) || contains(content, `"`+pattern+`"`))
}

func contains(s, substr string) bool {
	for i := 0; i+len(substr) <= len(s); i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			c1 := s[i+j]
			c2 := substr[j]
			// Case-insensitive comparison
			if c1 != c2 && c1 != c2+32 && c1 != c2-32 {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
