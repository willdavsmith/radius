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

package serve

import (
	"github.com/radius-project/radius/pkg/cli/discovery/llm"
	"github.com/radius-project/radius/pkg/cli/discovery/skills"
)

// RegisterAllTools registers all available tools with the skill registry.
// This function ensures CLI commands use the same skill implementations as the MCP server.
func RegisterAllTools(registry *skills.Registry) error {
	return skills.RegisterBuiltinSkills(registry, llm.NewMockProvider())
}

// GetToolNames returns a list of all registered tool names.
func GetToolNames(registry *skills.Registry) []string {
	metadata := registry.ListMetadata()
	names := make([]string, 0, len(metadata))
	for _, m := range metadata {
		names = append(names, m.Name)
	}
	return names
}
