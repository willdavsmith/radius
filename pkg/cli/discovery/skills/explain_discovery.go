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
	"strings"

	"github.com/radius-project/radius/pkg/cli/discovery"
)

// ExplainDiscovery generates natural language explanations of discovery results
// for AI agents to communicate to users.
type ExplainDiscovery struct{}

// NewExplainDiscovery creates a new explain discovery skill.
func NewExplainDiscovery() *ExplainDiscovery {
	return &ExplainDiscovery{}
}

// Name returns the skill identifier.
func (e *ExplainDiscovery) Name() string {
	return "explain_discovery"
}

// Description returns a human-readable description.
func (e *ExplainDiscovery) Description() string {
	return `Generates natural language explanations of discovery results.
Produces user-friendly summaries, deployment guidance, and next steps
that AI agents can use to communicate findings to users.`
}

// InputSchema returns the JSON Schema for the skill input.
func (e *ExplainDiscovery) InputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
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
			"appName": {
				"type": "string",
				"description": "Name of the application"
			},
			"format": {
				"type": "string",
				"description": "Output format: summary, detailed, or markdown",
				"default": "detailed"
			}
		},
		"required": ["dependencies", "services"]
	}`)
}

// OutputSchema returns the JSON Schema for the skill output.
func (e *ExplainDiscovery) OutputSchema() json.RawMessage {
	return json.RawMessage(`{
		"type": "object",
		"properties": {
			"explanation": {
				"type": "string",
				"description": "Natural language explanation of the discovery results"
			},
			"sections": {
				"type": "object",
				"properties": {
					"overview": {"type": "string"},
					"dependencies": {"type": "string"},
					"services": {"type": "string"},
					"architecture": {"type": "string"},
					"nextSteps": {"type": "string"}
				}
			}
		}
	}`)
}

// ExplainDiscoveryInput defines the input parameters.
type ExplainDiscoveryInput struct {
	Dependencies []discovery.DetectedDependency `json:"dependencies"`
	Services     []discovery.DetectedService    `json:"services"`
	AppName      string                         `json:"appName,omitempty"`
	Format       string                         `json:"format,omitempty"`
}

// ExplainDiscoveryOutput defines the output structure.
type ExplainDiscoveryOutput struct {
	Explanation string           `json:"explanation"`
	Sections    ExplanationParts `json:"sections"`
}

// ExplanationParts contains individual sections of the explanation.
type ExplanationParts struct {
	Overview     string `json:"overview"`
	Dependencies string `json:"dependencies"`
	Services     string `json:"services"`
	Architecture string `json:"architecture"`
	NextSteps    string `json:"nextSteps"`
}

// Execute runs the explanation generation.
func (e *ExplainDiscovery) Execute(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	var in ExplainDiscoveryInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, discovery.WrapError(discovery.ErrCodeInputValidation, "invalid input", err)
	}

	if in.Format == "" {
		in.Format = "detailed"
	}
	if in.AppName == "" {
		in.AppName = "your application"
	}

	sections := e.generateSections(in)
	explanation := e.formatExplanation(sections, in.Format)

	output := ExplainDiscoveryOutput{
		Explanation: explanation,
		Sections:    sections,
	}

	return json.Marshal(output)
}

// generateSections creates the individual explanation sections.
func (e *ExplainDiscovery) generateSections(in ExplainDiscoveryInput) ExplanationParts {
	return ExplanationParts{
		Overview:     e.generateOverview(in),
		Dependencies: e.generateDependenciesSection(in.Dependencies),
		Services:     e.generateServicesSection(in.Services),
		Architecture: e.generateArchitectureSection(in),
		NextSteps:    e.generateNextStepsSection(in),
	}
}

func (e *ExplainDiscovery) generateOverview(in ExplainDiscoveryInput) string {
	var sb strings.Builder

	depCount := len(in.Dependencies)
	svcCount := len(in.Services)

	if depCount == 0 && svcCount == 0 {
		return fmt.Sprintf("No infrastructure dependencies or services were detected in %s. "+
			"This could mean the codebase uses inline configuration, or the detection patterns "+
			"didn't match your technology stack.", in.AppName)
	}

	sb.WriteString(fmt.Sprintf("Analysis of %s found ", in.AppName))

	if depCount > 0 {
		sb.WriteString(fmt.Sprintf("%d infrastructure %s", depCount, e.pluralize("dependency", "dependencies", depCount)))
	}

	if depCount > 0 && svcCount > 0 {
		sb.WriteString(" and ")
	}

	if svcCount > 0 {
		sb.WriteString(fmt.Sprintf("%d %s", svcCount, e.pluralize("service", "services", svcCount)))
	}

	sb.WriteString(". ")

	// Add high-level characterization
	if e.hasDatabase(in.Dependencies) && svcCount > 0 {
		sb.WriteString("This appears to be a data-driven application ")
		if e.hasCache(in.Dependencies) {
			sb.WriteString("with caching for improved performance")
		} else {
			sb.WriteString("that would benefit from adding a caching layer")
		}
		sb.WriteString(".")
	}

	return sb.String()
}

func (e *ExplainDiscovery) generateDependenciesSection(deps []discovery.DetectedDependency) string {
	if len(deps) == 0 {
		return "No external infrastructure dependencies were detected."
	}

	var sb strings.Builder
	sb.WriteString("The following infrastructure dependencies were detected:\n\n")

	// Group by category
	databases := make([]discovery.DetectedDependency, 0)
	caches := make([]discovery.DetectedDependency, 0)
	messaging := make([]discovery.DetectedDependency, 0)
	other := make([]discovery.DetectedDependency, 0)

	for _, dep := range deps {
		switch dep.Type {
		case "postgresql", "mysql", "mongodb", "cosmosdb", "dynamodb", "sqlite":
			databases = append(databases, dep)
		case "redis", "memcached":
			caches = append(caches, dep)
		case "rabbitmq", "kafka", "sqs", "servicebus":
			messaging = append(messaging, dep)
		default:
			other = append(other, dep)
		}
	}

	if len(databases) > 0 {
		sb.WriteString("**Databases:**\n")
		for _, db := range databases {
			sb.WriteString(fmt.Sprintf("- %s (%s confidence: %.0f%%)\n",
				db.Technology, db.Type, db.Confidence*100))
		}
		sb.WriteString("\n")
	}

	if len(caches) > 0 {
		sb.WriteString("**Caching:**\n")
		for _, cache := range caches {
			sb.WriteString(fmt.Sprintf("- %s (%s confidence: %.0f%%)\n",
				cache.Technology, cache.Type, cache.Confidence*100))
		}
		sb.WriteString("\n")
	}

	if len(messaging) > 0 {
		sb.WriteString("**Messaging:**\n")
		for _, msg := range messaging {
			sb.WriteString(fmt.Sprintf("- %s (%s confidence: %.0f%%)\n",
				msg.Technology, msg.Type, msg.Confidence*100))
		}
		sb.WriteString("\n")
	}

	if len(other) > 0 {
		sb.WriteString("**Other Services:**\n")
		for _, o := range other {
			sb.WriteString(fmt.Sprintf("- %s (%s confidence: %.0f%%)\n",
				o.Technology, o.Type, o.Confidence*100))
		}
	}

	return sb.String()
}

func (e *ExplainDiscovery) generateServicesSection(services []discovery.DetectedService) string {
	if len(services) == 0 {
		return "No application services or entrypoints were detected."
	}

	var sb strings.Builder
	sb.WriteString("The following services were detected:\n\n")

	for _, svc := range services {
		sb.WriteString(fmt.Sprintf("- **%s**", svc.Name))
		if svc.Framework != "" {
			sb.WriteString(fmt.Sprintf(" (%s", svc.Framework))
			if svc.Language != "" {
				sb.WriteString(fmt.Sprintf("/%s", svc.Language))
			}
			sb.WriteString(")")
		}
		if svc.Port > 0 {
			sb.WriteString(fmt.Sprintf(" on port %d", svc.Port))
		}
		if svc.Entrypoint != "" {
			sb.WriteString(fmt.Sprintf(" - entrypoint: `%s`", svc.Entrypoint))
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

func (e *ExplainDiscovery) generateArchitectureSection(in ExplainDiscoveryInput) string {
	if len(in.Dependencies) == 0 && len(in.Services) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("Based on the detected components, the application architecture includes:\n\n")

	// Describe the architecture pattern
	svcCount := len(in.Services)
	if svcCount == 1 {
		sb.WriteString("- **Monolithic service** with external dependencies\n")
	} else if svcCount > 1 {
		sb.WriteString(fmt.Sprintf("- **%d services** suggesting a distributed/microservices architecture\n", svcCount))
	}

	if e.hasDatabase(in.Dependencies) {
		sb.WriteString("- **Persistent data layer** for storing application state\n")
	}

	if e.hasCache(in.Dependencies) {
		sb.WriteString("- **Caching layer** for performance optimization\n")
	}

	if e.hasMessaging(in.Dependencies) {
		sb.WriteString("- **Message-based communication** for async processing or service decoupling\n")
	}

	return sb.String()
}

func (e *ExplainDiscovery) generateNextStepsSection(in ExplainDiscoveryInput) string {
	var sb strings.Builder
	sb.WriteString("Recommended next steps:\n\n")

	sb.WriteString("1. **Generate Application Definition**: Run `rad app generate` to create a Bicep file\n")
	sb.WriteString("2. **Review Generated Code**: Check the app.bicep file and update container image references\n")
	sb.WriteString("3. **Configure Recipes**: Ensure your Radius environment has recipes configured for the detected dependencies\n")
	sb.WriteString("4. **Deploy**: Run `rad deploy app.bicep` to deploy your application\n")

	if len(in.Dependencies) > 0 {
		sb.WriteString("\n**For production deployment:**\n")
		sb.WriteString("- Configure environment-specific recipes (Azure, AWS) instead of local-dev recipes\n")
		sb.WriteString("- Set up proper secrets management for connection strings\n")
		sb.WriteString("- Consider adding monitoring and observability\n")
	}

	return sb.String()
}

func (e *ExplainDiscovery) formatExplanation(sections ExplanationParts, format string) string {
	var sb strings.Builder

	switch format {
	case "summary":
		sb.WriteString(sections.Overview)
	case "markdown":
		sb.WriteString("# Discovery Results\n\n")
		sb.WriteString("## Overview\n\n")
		sb.WriteString(sections.Overview)
		sb.WriteString("\n\n## Infrastructure Dependencies\n\n")
		sb.WriteString(sections.Dependencies)
		sb.WriteString("\n\n## Services\n\n")
		sb.WriteString(sections.Services)
		if sections.Architecture != "" {
			sb.WriteString("\n\n## Architecture\n\n")
			sb.WriteString(sections.Architecture)
		}
		sb.WriteString("\n\n## Next Steps\n\n")
		sb.WriteString(sections.NextSteps)
	default: // detailed
		sb.WriteString(sections.Overview)
		sb.WriteString("\n\n")
		sb.WriteString(sections.Dependencies)
		sb.WriteString("\n")
		sb.WriteString(sections.Services)
		if sections.Architecture != "" {
			sb.WriteString("\n")
			sb.WriteString(sections.Architecture)
		}
		sb.WriteString("\n")
		sb.WriteString(sections.NextSteps)
	}

	return sb.String()
}

func (e *ExplainDiscovery) pluralize(singular, plural string, count int) string {
	if count == 1 {
		return singular
	}
	return plural
}

func (e *ExplainDiscovery) hasDatabase(deps []discovery.DetectedDependency) bool {
	for _, dep := range deps {
		switch dep.Type {
		case "postgresql", "mysql", "mongodb", "cosmosdb", "dynamodb", "sqlite":
			return true
		}
	}
	return false
}

func (e *ExplainDiscovery) hasCache(deps []discovery.DetectedDependency) bool {
	for _, dep := range deps {
		if dep.Type == "redis" || dep.Type == "memcached" {
			return true
		}
	}
	return false
}

func (e *ExplainDiscovery) hasMessaging(deps []discovery.DetectedDependency) bool {
	for _, dep := range deps {
		switch dep.Type {
		case "rabbitmq", "kafka", "sqs", "servicebus":
			return true
		}
	}
	return false
}
