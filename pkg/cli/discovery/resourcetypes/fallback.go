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
	"strings"

	"github.com/radius-project/radius/pkg/cli/discovery"
)

// FallbackGenerator generates fallback Resource Types for unknown dependencies.
type FallbackGenerator struct {
	defaultAPIVersion string
}

// NewFallbackGenerator creates a new fallback generator.
func NewFallbackGenerator() *FallbackGenerator {
	return &FallbackGenerator{
		defaultAPIVersion: "2023-10-01-preview",
	}
}

// Generate creates a fallback Resource Type for an unknown dependency.
func (f *FallbackGenerator) Generate(dep discovery.DetectedDependency) discovery.ResourceType {
	name := f.generateName(dep)
	radiusType := f.determineRadiusType(dep)

	return discovery.ResourceType{
		Name:       name,
		Type:       radiusType,
		APIVersion: f.defaultAPIVersion,
		Maturity:   "alpha",
		Schema:     f.generateSchema(dep),
		Outputs:    f.generateOutputs(dep),
	}
}

// generateName creates a human-readable name from the dependency.
func (f *FallbackGenerator) generateName(dep discovery.DetectedDependency) string {
	if dep.Technology != "" {
		return dep.Technology
	}
	name := strings.ReplaceAll(dep.Type, "-", " ")
	name = strings.ReplaceAll(name, "_", " ")
	return strings.Title(name)
}

// determineRadiusType determines the appropriate Radius type based on the dependency.
func (f *FallbackGenerator) determineRadiusType(dep discovery.DetectedDependency) string {
	depType := strings.ToLower(dep.Type)

	if containsAny(depType, []string{"sql", "db", "database", "postgres", "mysql", "mongo", "redis", "cache", "storage"}) {
		return "Applications.Core/extenders"
	}

	if containsAny(depType, []string{"queue", "message", "kafka", "rabbit", "amqp", "pubsub"}) {
		return "Applications.Core/extenders"
	}

	if containsAny(depType, []string{"dapr"}) {
		return "Applications.Dapr/secretStores"
	}

	return "Applications.Core/extenders"
}

// generateSchema creates a minimal schema for the fallback type.
func (f *FallbackGenerator) generateSchema(dep discovery.DetectedDependency) *discovery.ResourceSchema {
	properties := map[string]discovery.PropertyDef{
		"resourceType": {
			Type:        "string",
			Description: "The underlying resource type identifier",
			Default:     dep.Type,
		},
	}

	if dep.ConnectionInfo != nil {
		if dep.ConnectionInfo.EnvVar != "" {
			properties["connectionEnvVar"] = discovery.PropertyDef{
				Type:        "string",
				Description: "Environment variable for connection string",
				Default:     dep.ConnectionInfo.EnvVar,
			}
		}
		if dep.ConnectionInfo.DefaultPort != 0 {
			properties["port"] = discovery.PropertyDef{
				Type:        "integer",
				Description: "Default port for the service",
				Default:     dep.ConnectionInfo.DefaultPort,
			}
		}
	}

	return &discovery.ResourceSchema{
		Properties: properties,
	}
}

// generateOutputs creates outputs based on the dependency type.
func (f *FallbackGenerator) generateOutputs(dep discovery.DetectedDependency) []discovery.OutputDefinition {
	outputs := []discovery.OutputDefinition{
		{
			Name:        "connectionString",
			Type:        "string",
			Description: "Connection string for the resource",
			Sensitive:   true,
		},
	}

	if dep.ConnectionInfo != nil && dep.ConnectionInfo.DefaultPort != 0 {
		outputs = append(outputs, discovery.OutputDefinition{
			Name:        "port",
			Type:        "integer",
			Description: "Port for the resource",
		})
	}

	return outputs
}

// containsAny checks if the string contains any of the substrings.
func containsAny(s string, substrings []string) bool {
	for _, sub := range substrings {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
