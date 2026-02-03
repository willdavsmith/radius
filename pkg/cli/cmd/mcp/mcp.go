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

// Package mcp provides the rad mcp command group for MCP (Model Context Protocol) server functionality.
package mcp

import (
	"github.com/spf13/cobra"

	"github.com/radius-project/radius/pkg/cli/framework"
)

// NewCommand creates a new cobra command for the mcp command group.
func NewCommand(factory framework.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "mcp",
		Short: "Manage MCP (Model Context Protocol) server for AI agent integration",
		Long: `The mcp command group provides functionality for running an MCP server
that exposes Radius discovery and generation skills to AI coding agents.

MCP (Model Context Protocol) enables AI assistants like GitHub Copilot and
other compatible agents to interact with Radius for application discovery
and deployment automation.

Available Commands:
  serve    Start the MCP server

Examples:
  # Start MCP server with stdio transport (for VS Code integration)
  rad mcp serve

  # Start MCP server with HTTP transport
  rad mcp serve --transport http --port 8080
`,
	}

	return cmd
}
