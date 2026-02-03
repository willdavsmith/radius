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

// Package services provides the service detection skill.
package services

import (
	"context"
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/radius-project/radius/pkg/cli/discovery"
	"github.com/radius-project/radius/pkg/cli/discovery/llm"
)

// Detector detects services and entrypoints in a codebase.
type Detector struct {
	llmProvider llm.Provider
	walker      *Walker
}

// NewDetector creates a new service detector.
func NewDetector(llmProvider llm.Provider) *Detector {
	return &Detector{
		llmProvider: llmProvider,
		walker:      NewWalker(nil),
	}
}

// Name returns the skill identifier.
func (d *Detector) Name() string {
	return "discover_services"
}

// Description returns a human-readable description.
func (d *Detector) Description() string {
	return "Detects services and application entrypoints in a codebase"
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
			"services": {
				"type": "array",
				"items": {
					"type": "object",
					"properties": {
						"name": {"type": "string"},
						"type": {"type": "string"},
						"port": {"type": "integer"},
						"entrypoint": {"type": "string"},
						"framework": {"type": "string"},
						"language": {"type": "string"},
						"confidence": {"type": "number"}
					}
				}
			}
		}
	}`)
}

// DetectorInput is the input for the service detector.
type DetectorInput struct {
	Path         string   `json:"path"`
	ExcludePaths []string `json:"excludePaths,omitempty"`
}

// DetectorOutput is the output of the service detector.
type DetectorOutput struct {
	Services []discovery.DetectedService `json:"services"`
}

// Execute runs the service detection.
func (d *Detector) Execute(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	var in DetectorInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, discovery.WrapError(discovery.ErrCodeInputValidation, "invalid input", err)
	}

	d.walker = NewWalker(in.ExcludePaths)

	services, err := d.detectServices(ctx, in.Path)
	if err != nil {
		return nil, err
	}

	output := DetectorOutput{
		Services: services,
	}
	return json.Marshal(output)
}

// detectServices analyzes the codebase to find service entrypoints.
func (d *Detector) detectServices(ctx context.Context, path string) ([]discovery.DetectedService, error) {
	files, err := d.walker.WalkSourceFiles(path)
	if err != nil {
		return nil, discovery.WrapError(discovery.ErrCodeFileNotFound, "failed to walk directory", err)
	}

	var services []discovery.DetectedService
	for _, file := range files {
		service := d.detectEntrypoint(ctx, file)
		if service != nil {
			services = append(services, *service)
		}
	}

	return d.deduplicateServices(services), nil
}

// detectEntrypoint checks if a file is an application entrypoint.
func (d *Detector) detectEntrypoint(ctx context.Context, file FileInfo) *discovery.DetectedService {
	if d.llmProvider != nil {
		service := d.llmAnalyzeEntrypoint(ctx, file)
		if service != nil {
			return service
		}
	}
	return d.patternBasedEntrypoint(file)
}

// llmAnalyzeEntrypoint uses the LLM to detect entrypoints.
func (d *Detector) llmAnalyzeEntrypoint(ctx context.Context, file FileInfo) *discovery.DetectedService {
	prompt := "Analyze this file and determine if it's an application entrypoint.\n\n" +
		"FILE PATH: " + file.RelativePath + "\n\n" +
		"CONTENT:\n" + file.Content

	result, err := d.llmProvider.Analyze(ctx, prompt, d.OutputSchema())
	if err != nil {
		return nil
	}

	var output DetectorOutput
	if err := json.Unmarshal(result, &output); err != nil {
		return nil
	}

	if len(output.Services) > 0 {
		return &output.Services[0]
	}
	return nil
}

// patternBasedEntrypoint detects entrypoints using patterns.
func (d *Detector) patternBasedEntrypoint(file FileInfo) *discovery.DetectedService {
	type entryPattern struct {
		Language    string
		Framework   string
		Type        string
		FileMatch   string
		Content     []string
		DefaultPort int
	}

	patterns := []entryPattern{
		{"nodejs", "express", "http", "*.js", []string{"express()", "app.listen("}, 3000},
		{"nodejs", "fastify", "http", "*.js", []string{"fastify(", "server.listen("}, 3000},
		{"python", "flask", "http", "*.py", []string{"from flask import", "Flask("}, 5000},
		{"python", "fastapi", "http", "*.py", []string{"from fastapi import", "FastAPI("}, 8000},
		{"go", "gin", "http", "*.go", []string{"gin.Default()", "gin.New()"}, 8080},
		{"go", "net/http", "http", "*.go", []string{"http.ListenAndServe("}, 8080},
		{"java", "spring-boot", "http", "*.java", []string{"@SpringBootApplication"}, 8080},
		{"ruby", "rails", "http", "config.ru", []string{"Rails.application"}, 3000},
	}

	for _, p := range patterns {
		if matched, _ := filepath.Match(p.FileMatch, filepath.Base(file.RelativePath)); !matched {
			continue
		}
		for _, content := range p.Content {
			if strings.Contains(file.Content, content) {
				port := d.extractPort(file.Content, p.DefaultPort)
				name := d.generateServiceName(file.RelativePath)
				return &discovery.DetectedService{
					Name:       name,
					Type:       p.Type,
					Port:       port,
					Entrypoint: file.RelativePath,
					Framework:  p.Framework,
					Language:   p.Language,
					Confidence: 0.7,
				}
			}
		}
	}

	if strings.ToLower(filepath.Base(file.RelativePath)) == "dockerfile" {
		return d.detectFromDockerfile(file)
	}
	return nil
}

// detectFromDockerfile extracts service info from a Dockerfile.
func (d *Detector) detectFromDockerfile(file FileInfo) *discovery.DetectedService {
	var port int
	for _, line := range strings.Split(file.Content, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(strings.ToUpper(line), "EXPOSE") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				portStr := strings.Split(parts[1], "/")[0]
				if p, err := strconv.Atoi(portStr); err == nil {
					port = p
				}
			}
		}
	}
	if port == 0 {
		return nil
	}

	dir := filepath.Dir(file.RelativePath)
	name := d.generateServiceName(dir)
	if name == "" || name == "." {
		name = "app"
	}
	return &discovery.DetectedService{
		Name:       name,
		Type:       "http",
		Port:       port,
		Entrypoint: file.RelativePath,
		Confidence: 0.6,
	}
}

// extractPort tries to extract a port number from code content.
func (d *Detector) extractPort(content string, defaultPort int) int {
	portPatterns := []string{"PORT", "port", "listen(", "ListenAndServe("}
	for _, pattern := range portPatterns {
		idx := strings.Index(content, pattern)
		if idx == -1 {
			continue
		}
		remainder := content[idx:]
		for i := 0; i < len(remainder) && i < 50; i++ {
			if remainder[i] >= '0' && remainder[i] <= '9' {
				var numStr strings.Builder
				for j := i; j < len(remainder) && remainder[j] >= '0' && remainder[j] <= '9'; j++ {
					numStr.WriteByte(remainder[j])
				}
				if port, err := strconv.Atoi(numStr.String()); err == nil && port > 0 && port < 65536 {
					return port
				}
			}
		}
	}
	return defaultPort
}

// generateServiceName creates a service name from a file path.
func (d *Detector) generateServiceName(path string) string {
	dir := filepath.Dir(path)
	if dir == "." || dir == "/" {
		base := filepath.Base(path)
		ext := filepath.Ext(base)
		return strings.TrimSuffix(base, ext)
	}
	return filepath.Base(dir)
}

// deduplicateServices removes duplicate services.
func (d *Detector) deduplicateServices(services []discovery.DetectedService) []discovery.DetectedService {
	seen := make(map[string]discovery.DetectedService)
	for _, svc := range services {
		key := svc.Name + "-" + svc.Entrypoint
		if existing, ok := seen[key]; ok {
			if svc.Confidence > existing.Confidence {
				seen[key] = svc
			}
		} else {
			seen[key] = svc
		}
	}
	result := make([]discovery.DetectedService, 0, len(seen))
	for _, svc := range seen {
		result = append(result, svc)
	}
	return result
}
