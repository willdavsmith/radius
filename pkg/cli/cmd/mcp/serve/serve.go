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

// Package serve provides the rad mcp serve command for starting the MCP server.
package serve

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/radius-project/radius/pkg/cli/discovery/llm"
	"github.com/radius-project/radius/pkg/cli/discovery/skills"
	"github.com/radius-project/radius/pkg/cli/framework"
	"github.com/radius-project/radius/pkg/cli/output"
)

const (
	defaultPort    = 8080
	transportStdio = "stdio"
	transportHTTP  = "http"
)

// NewCommand creates a new cobra command for rad mcp serve.
func NewCommand(factory framework.Factory) (*cobra.Command, framework.Runner) {
	runner := NewRunner(factory)
	cmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the MCP server for AI agent integration",
		Long: `Start the MCP (Model Context Protocol) server to expose Radius discovery
and generation skills to AI coding agents.

The MCP server supports two transport modes:
  - stdio: Communicates via stdin/stdout (default, for VS Code integration)
  - http: Communicates via HTTP/REST API (for remote agents)

Examples:
  # Start with stdio transport (for VS Code)
  rad mcp serve

  # Start with HTTP transport on default port (8080)
  rad mcp serve --transport http

  # Start with HTTP transport on custom port
  rad mcp serve --transport http --port 9000
`,
		Example: `# Start MCP server with stdio transport
rad mcp serve

# Start MCP server with HTTP transport
rad mcp serve --transport http --port 8080`,
		RunE: framework.RunCommand(runner),
	}

	cmd.Flags().String("transport", transportStdio, "Transport mode: stdio or http")
	cmd.Flags().Int("port", defaultPort, "Port for HTTP transport")
	cmd.Flags().StringSlice("allowed-origins", []string{"*"}, "Allowed CORS origins for HTTP transport")

	return cmd, runner
}

// Runner is the runner for the serve command.
type Runner struct {
	Output output.Interface

	transport      string
	port           int
	allowedOrigins []string
}

// NewRunner creates a new Runner.
func NewRunner(factory framework.Factory) *Runner {
	return &Runner{
		Output: factory.GetOutput(),
	}
}

// Validate validates the command arguments.
func (r *Runner) Validate(cmd *cobra.Command, args []string) error {
	r.transport, _ = cmd.Flags().GetString("transport")
	r.port, _ = cmd.Flags().GetInt("port")
	r.allowedOrigins, _ = cmd.Flags().GetStringSlice("allowed-origins")

	if r.transport != transportStdio && r.transport != transportHTTP {
		return fmt.Errorf("invalid transport: %s (must be 'stdio' or 'http')", r.transport)
	}

	if r.transport == transportHTTP && (r.port < 1 || r.port > 65535) {
		return fmt.Errorf("invalid port: %d (must be between 1 and 65535)", r.port)
	}

	return nil
}

// Run executes the serve command.
func (r *Runner) Run(ctx context.Context) error {
	registry := skills.NewRegistry()

	if err := skills.RegisterBuiltinSkills(registry, llm.NewMockProvider()); err != nil {
		return fmt.Errorf("failed to register skills: %w", err)
	}

	handler := NewHandler(registry)

	if r.transport == transportStdio {
		r.Output.LogInfo("Starting MCP server with stdio transport...")
		return r.runStdio(ctx, handler)
	}

	r.Output.LogInfo("Starting MCP server on port %d...", r.port)
	return r.runHTTP(ctx, handler)
}

func (r *Runner) runStdio(ctx context.Context, handler *Handler) error {
	server := NewStdioServer(handler)
	return server.Serve(ctx)
}

func (r *Runner) runHTTP(ctx context.Context, handler *Handler) error {
	server := NewHTTPServer(handler, r.port, r.allowedOrigins)
	return server.Serve(ctx)
}

var _ framework.Runner = (*Runner)(nil)
