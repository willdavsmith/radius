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

package discovery

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"text/template"
	"time"

	"gopkg.in/yaml.v3"
)

// OutputWriter writes discovery results to files.
type OutputWriter struct {
	outputDir string
}

// NewOutputWriter creates a new output writer.
func NewOutputWriter(outputDir string) *OutputWriter {
	return &OutputWriter{
		outputDir: outputDir,
	}
}

// WriteDiscoveryResult is a convenience function to write a discovery result to a file.
func WriteDiscoveryResult(path string, result DiscoveryResult) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	w := NewOutputWriter(dir)
	content, err := w.generateMarkdown(&result)
	if err != nil {
		return err
	}

	return os.WriteFile(path, []byte(content), 0644)
}

// WriteDiscoveryResult writes the discovery result to discovery.md.
func (w *OutputWriter) WriteDiscoveryResult(result *DiscoveryResult) error {
	if err := os.MkdirAll(w.outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate Markdown content
	content, err := w.generateMarkdown(result)
	if err != nil {
		return fmt.Errorf("failed to generate markdown: %w", err)
	}

	// Write to file
	outputPath := filepath.Join(w.outputDir, "discovery.md")
	if err := os.WriteFile(outputPath, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write discovery.md: %w", err)
	}

	return nil
}

// generateMarkdown generates the Markdown content with YAML frontmatter.
func (w *OutputWriter) generateMarkdown(result *DiscoveryResult) (string, error) {
	// Sort dependencies and services for deterministic output
	sortedDeps := make([]DetectedDependency, len(result.Dependencies))
	copy(sortedDeps, result.Dependencies)
	sort.Slice(sortedDeps, func(i, j int) bool {
		return sortedDeps[i].Type < sortedDeps[j].Type
	})

	sortedServices := make([]DetectedService, len(result.Services))
	copy(sortedServices, result.Services)
	sort.Slice(sortedServices, func(i, j int) bool {
		return sortedServices[i].Name < sortedServices[j].Name
	})

	// Create frontmatter data
	frontmatter := map[string]any{
		"version":       result.Version,
		"analyzed_at":   result.AnalyzedAt.Format(time.RFC3339),
		"codebase_path": result.CodebasePath,
		"dependencies":  w.dependenciesToFrontmatter(sortedDeps),
		"services":      w.servicesToFrontmatter(sortedServices),
	}

	if result.Practices != nil {
		frontmatter["practices"] = result.Practices
	}

	// Generate YAML frontmatter
	var yamlBuf bytes.Buffer
	yamlEncoder := yaml.NewEncoder(&yamlBuf)
	yamlEncoder.SetIndent(2)
	if err := yamlEncoder.Encode(frontmatter); err != nil {
		return "", fmt.Errorf("failed to encode frontmatter: %w", err)
	}

	// Generate Markdown body
	tmpl, err := template.New("discovery").Parse(discoveryTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var bodyBuf bytes.Buffer
	data := templateData{
		Result:       result,
		Dependencies: sortedDeps,
		Services:     sortedServices,
	}
	if err := tmpl.Execute(&bodyBuf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	// Combine frontmatter and body
	return fmt.Sprintf("---\n%s---\n\n%s", yamlBuf.String(), bodyBuf.String()), nil
}

type templateData struct {
	Result       *DiscoveryResult
	Dependencies []DetectedDependency
	Services     []DetectedService
}

func (w *OutputWriter) dependenciesToFrontmatter(deps []DetectedDependency) []map[string]any {
	result := make([]map[string]any, 0, len(deps))
	for _, dep := range deps {
		evidence := make([]string, 0, len(dep.Evidence))
		for _, e := range dep.Evidence {
			evidence = append(evidence, fmt.Sprintf("%s:%s", e.Source, e.Value))
		}

		result = append(result, map[string]any{
			"type":       dep.Type,
			"technology": dep.Technology,
			"confidence": dep.Confidence,
			"evidence":   evidence,
		})
	}
	return result
}

func (w *OutputWriter) servicesToFrontmatter(services []DetectedService) []map[string]any {
	result := make([]map[string]any, 0, len(services))
	for _, svc := range services {
		m := map[string]any{
			"name": svc.Name,
			"type": svc.Type,
		}
		if svc.Port > 0 {
			m["port"] = svc.Port
		}
		if svc.Entrypoint != "" {
			m["entrypoint"] = svc.Entrypoint
		}
		result = append(result, m)
	}
	return result
}

const discoveryTemplate = `# Discovery Report

Generated: {{.Result.AnalyzedAt.Format "2006-01-02 15:04:05"}}
Codebase: {{.Result.CodebasePath}}

## Dependencies Detected

{{if .Dependencies -}}
| Type | Technology | Confidence | Evidence |
|------|------------|------------|----------|
{{range .Dependencies -}}
| {{.Type}} | {{.Technology}} | {{printf "%.0f" (mul .Confidence 100)}}% | {{range $i, $e := .Evidence}}{{if $i}}, {{end}}{{$e.Source}}{{end}} |
{{end}}
{{else -}}
No dependencies detected.
{{end}}

## Services Detected

{{if .Services -}}
| Name | Type | Port | Entrypoint |
|------|------|------|------------|
{{range .Services -}}
| {{.Name}} | {{.Type}} | {{if .Port}}{{.Port}}{{else}}-{{end}} | {{if .Entrypoint}}{{.Entrypoint}}{{else}}-{{end}} |
{{end}}
{{else -}}
No services detected.
{{end}}

{{if .Result.Warnings -}}
## Warnings

{{range .Result.Warnings -}}
- **{{.Code}}**: {{.Message}}{{if .Source}} ({{.Source}}){{end}}
{{end}}
{{end}}

## Next Steps

1. Review the detected dependencies and services above
2. Run ` + "`rad app generate`" + ` to create app.bicep
3. Customize the generated application definition as needed
4. Run ` + "`rad deploy`" + ` to deploy your application
`

func init() {
	// Register template function
	template.Must(template.New("").Funcs(template.FuncMap{
		"mul": func(a, b float64) float64 { return a * b },
	}).Parse(""))
}
