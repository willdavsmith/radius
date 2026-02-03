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

// Package generator provides Bicep application generation from discovery results.
package generator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"text/template"

	"github.com/radius-project/radius/pkg/cli/discovery"
	"github.com/radius-project/radius/pkg/cli/discovery/generator/templates"
	"github.com/radius-project/radius/pkg/cli/discovery/resourcetypes"
)

// BicepGenerator generates Bicep application definitions from discovery results.
type BicepGenerator struct {
	catalog *resourcetypes.Catalog
}

// BicepGeneratorInput is the input for the Bicep generator.
type BicepGeneratorInput struct {
	AppName      string                         `json:"appName"`
	Dependencies []discovery.DetectedDependency `json:"dependencies"`
	Services     []discovery.DetectedService    `json:"services"`
}

// BicepGeneratorOutput is the output from the Bicep generator.
type BicepGeneratorOutput struct {
	BicepContent  string                  `json:"bicepContent"`
	AppDefinition discovery.AppDefinition `json:"appDefinition"`
	Warnings      []string                `json:"warnings,omitempty"`
}

// NewBicepGenerator creates a new Bicep generator.
func NewBicepGenerator() *BicepGenerator {
	return &BicepGenerator{
		catalog: resourcetypes.NewCatalog(),
	}
}

// Name returns the skill name.
func (g *BicepGenerator) Name() string {
	return "generate_app_definition"
}

// Description returns the skill description.
func (g *BicepGenerator) Description() string {
	return "Generates a Radius application definition (app.bicep) from discovery results."
}

// InputSchema returns the JSON Schema for the input.
func (g *BicepGenerator) InputSchema() json.RawMessage {
	return json.RawMessage(`{"type":"object","properties":{"appName":{"type":"string"},"dependencies":{"type":"array"},"services":{"type":"array"}},"required":["appName","dependencies","services"]}`)
}

// OutputSchema returns the JSON Schema for the output.
func (g *BicepGenerator) OutputSchema() json.RawMessage {
	return json.RawMessage(`{"type":"object","properties":{"bicepContent":{"type":"string"},"appDefinition":{"type":"object"},"warnings":{"type":"array"}}}`)
}

// Execute runs the Bicep generation.
func (g *BicepGenerator) Execute(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	var in BicepGeneratorInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}
	output, err := g.Generate(in)
	if err != nil {
		return nil, err
	}
	return json.Marshal(output)
}

// Generate creates the Bicep content from input.
func (g *BicepGenerator) Generate(input BicepGeneratorInput) (*BicepGeneratorOutput, error) {
	output := &BicepGeneratorOutput{
		Warnings: make([]string, 0),
	}

	appDef := discovery.AppDefinition{
		Name:        input.AppName,
		Containers:  make([]discovery.ContainerDef, 0),
		Resources:   make([]discovery.ResourceDef, 0),
		Connections: make([]discovery.ConnectionDef, 0),
	}

	for _, dep := range input.Dependencies {
		resourceDef := g.createResourceDef(dep)
		appDef.Resources = append(appDef.Resources, resourceDef)
	}

	for _, svc := range input.Services {
		containerDef := g.createContainerDef(svc)
		appDef.Containers = append(appDef.Containers, containerDef)
		for _, depType := range svc.Dependencies {
			conn := discovery.ConnectionDef{
				Container: svc.Name,
				Resource:  g.normalizeResourceName(depType),
			}
			appDef.Connections = append(appDef.Connections, conn)
		}
	}

	bicepContent, err := g.renderBicep(input.AppName, appDef)
	if err != nil {
		return nil, fmt.Errorf("failed to render Bicep: %w", err)
	}

	appDef.BicepContent = bicepContent
	output.BicepContent = bicepContent
	output.AppDefinition = appDef

	return output, nil
}

func (g *BicepGenerator) createResourceDef(dep discovery.DetectedDependency) discovery.ResourceDef {
	entry, found := g.catalog.Get(dep.Type)
	resourceType := "Applications.Core/extenders"
	if found {
		resourceType = entry.ResourceType.Type
	}
	return discovery.ResourceDef{
		Name:   g.normalizeResourceName(dep.Type),
		Type:   resourceType,
		Recipe: "default",
	}
}

func (g *BicepGenerator) createContainerDef(svc discovery.DetectedService) discovery.ContainerDef {
	container := discovery.ContainerDef{
		Name:  svc.Name,
		Image: fmt.Sprintf("TODO: %s", svc.Name),
		Env:   make(map[string]discovery.EnvValue),
	}
	if svc.Port > 0 {
		container.Ports = []discovery.PortDef{
			{ContainerPort: svc.Port, Protocol: "TCP"},
		}
	}
	return container
}

func (g *BicepGenerator) normalizeResourceName(name string) string {
	name = strings.ToLower(name)
	name = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(name, "_")
	name = strings.Trim(name, "_")
	return name
}

func (g *BicepGenerator) renderBicep(appName string, appDef discovery.AppDefinition) (string, error) {
	tmpl, err := template.New("app").Parse(templates.AppTemplate)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	data := struct {
		Name       string
		Resources  []resourceTemplateData
		Containers []containerTemplateData
	}{
		Name:       appName,
		Resources:  make([]resourceTemplateData, 0),
		Containers: make([]containerTemplateData, 0),
	}

	for _, res := range appDef.Resources {
		data.Resources = append(data.Resources, resourceTemplateData{
			Name:         res.Name,
			ResourceName: g.normalizeResourceName(res.Name),
			RadiusType:   res.Type,
			Recipe:       res.Recipe,
			BicepContent: g.renderResourceBicep(res),
		})
	}

	for _, cont := range appDef.Containers {
		data.Containers = append(data.Containers, containerTemplateData{
			Name:         cont.Name,
			ResourceName: g.normalizeResourceName(cont.Name),
			Image:        cont.Image,
			BicepContent: g.renderContainerBicep(cont, appDef.Connections),
		})
	}

	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return buf.String(), nil
}

type resourceTemplateData struct {
	Name         string
	ResourceName string
	RadiusType   string
	Recipe       string
	BicepContent string
}

type containerTemplateData struct {
	Name         string
	ResourceName string
	Image        string
	BicepContent string
}

func (g *BicepGenerator) renderResourceBicep(res discovery.ResourceDef) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("resource %s '%s@2023-10-01-preview' = {\n", g.normalizeResourceName(res.Name), res.Type))
	buf.WriteString(fmt.Sprintf("  name: '%s'\n", res.Name))
	buf.WriteString("  properties: {\n")
	buf.WriteString("    application: app.id\n")
	buf.WriteString("    environment: environment\n")
	if res.Recipe != "" {
		buf.WriteString("    recipe: {\n")
		buf.WriteString(fmt.Sprintf("      name: '%s'\n", res.Recipe))
		buf.WriteString("    }\n")
	}
	buf.WriteString("  }\n")
	buf.WriteString("}\n")
	return buf.String()
}

func (g *BicepGenerator) renderContainerBicep(cont discovery.ContainerDef, connections []discovery.ConnectionDef) string {
	var buf bytes.Buffer
	resourceName := g.normalizeResourceName(cont.Name)

	buf.WriteString(fmt.Sprintf("resource %s 'Applications.Core/containers@2023-10-01-preview' = {\n", resourceName))
	buf.WriteString(fmt.Sprintf("  name: '%s'\n", cont.Name))
	buf.WriteString("  properties: {\n")
	buf.WriteString("    application: app.id\n")
	buf.WriteString("    container: {\n")
	buf.WriteString(fmt.Sprintf("      image: '%s'\n", cont.Image))

	if len(cont.Ports) > 0 {
		buf.WriteString("      ports: {\n")
		for _, port := range cont.Ports {
			buf.WriteString("        http: {\n")
			buf.WriteString(fmt.Sprintf("          containerPort: %d\n", port.ContainerPort))
			buf.WriteString("        }\n")
		}
		buf.WriteString("      }\n")
	}

	buf.WriteString("    }\n")

	containerConnections := make([]discovery.ConnectionDef, 0)
	for _, conn := range connections {
		if conn.Container == cont.Name {
			containerConnections = append(containerConnections, conn)
		}
	}

	if len(containerConnections) > 0 {
		buf.WriteString("    connections: {\n")
		for _, conn := range containerConnections {
			buf.WriteString(fmt.Sprintf("      %s: {\n", g.normalizeResourceName(conn.Resource)))
			buf.WriteString(fmt.Sprintf("        source: %s.id\n", g.normalizeResourceName(conn.Resource)))
			buf.WriteString("      }\n")
		}
		buf.WriteString("    }\n")
	}

	buf.WriteString("  }\n")
	buf.WriteString("}\n")

	return buf.String()
}
