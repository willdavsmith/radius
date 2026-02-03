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

package generator

import (
	"regexp"
	"strings"

	"github.com/radius-project/radius/pkg/cli/discovery"
)

// ConnectionWirer wires connections between containers and resources.
type ConnectionWirer struct{}

// NewConnectionWirer creates a new ConnectionWirer.
func NewConnectionWirer() *ConnectionWirer {
	return &ConnectionWirer{}
}

// WireConnections creates connection definitions based on service dependencies.
func (w *ConnectionWirer) WireConnections(services []discovery.DetectedService, resources []discovery.ResourceDef) []discovery.ConnectionDef {
	connections := make([]discovery.ConnectionDef, 0)
	resourceMap := make(map[string]discovery.ResourceDef)
	for _, res := range resources {
		resourceMap[strings.ToLower(res.Name)] = res
	}

	for _, svc := range services {
		for _, depType := range svc.Dependencies {
			normalizedDep := strings.ToLower(depType)
			if _, found := resourceMap[normalizedDep]; found {
				conn := discovery.ConnectionDef{
					Container: svc.Name,
					Resource:  normalizedDep,
				}
				connections = append(connections, conn)
			}
		}
	}

	return connections
}

// InferConnections infers connections from dependency patterns.
// This can be extended to use environment variables once that field is available.
func (w *ConnectionWirer) InferConnections(services []discovery.DetectedService, resources []discovery.ResourceDef) []discovery.ConnectionDef {
	connections := make([]discovery.ConnectionDef, 0)
	resourceMap := make(map[string]discovery.ResourceDef)
	for _, res := range resources {
		resourceMap[strings.ToLower(res.Name)] = res
	}

	envMappings := getDefaultEnvMapping()

	for _, svc := range services {
		for _, dep := range svc.Dependencies {
			upperDep := strings.ToUpper(dep)
			for resType, patterns := range envMappings {
				for _, pattern := range patterns {
					if strings.Contains(upperDep, pattern) {
						if _, found := resourceMap[resType]; found {
							conn := discovery.ConnectionDef{
								Container: svc.Name,
								Resource:  resType,
							}
							if !connectionExists(connections, conn) {
								connections = append(connections, conn)
							}
						}
					}
				}
			}
		}
	}

	return connections
}

func connectionExists(connections []discovery.ConnectionDef, conn discovery.ConnectionDef) bool {
	for _, c := range connections {
		if c.Container == conn.Container && c.Resource == conn.Resource {
			return true
		}
	}
	return false
}

func getDefaultEnvMapping() map[string][]string {
	return map[string][]string{
		"postgresql": {"POSTGRES", "PGHOST", "PGUSER", "PGDATABASE", "PG_"},
		"redis":      {"REDIS", "CACHE_HOST", "REDIS_URL", "REDIS_HOST"},
		"mongodb":    {"MONGO", "MONGODB", "MONGO_URI", "MONGO_URL"},
		"rabbitmq":   {"RABBIT", "AMQP", "RABBITMQ"},
		"kafka":      {"KAFKA", "KAFKA_BROKER", "KAFKA_URL"},
		"mysql":      {"MYSQL", "MYSQL_HOST", "MYSQL_USER", "MYSQL_DATABASE"},
	}
}

func toUpperSnakeCase(s string) string {
	s = regexp.MustCompile(`([a-z])([A-Z])`).ReplaceAllString(s, "${1}_${2}")
	s = regexp.MustCompile(`[^a-zA-Z0-9]+`).ReplaceAllString(s, "_")
	return strings.ToUpper(s)
}
