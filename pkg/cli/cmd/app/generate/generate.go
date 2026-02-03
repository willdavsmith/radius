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

// Package generate provides the rad app generate command.
package generate

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/radius-project/radius/pkg/cli/discovery"
	"github.com/radius-project/radius/pkg/cli/discovery/generator"
	"github.com/radius-project/radius/pkg/cli/framework"
	"github.com/radius-project/radius/pkg/cli/output"
	"github.com/radius-project/radius/pkg/cli/prompt"
)

const (
	defaultDiscoveryFile = ".discovery/result.json"
	defaultOutputFile    = "app.bicep"
)

// NewCommand creates a new cobra command for rad app generate.
func NewCommand(factory framework.Factory) (*cobra.Command, framework.Runner) {
	runner := NewRunner(factory)
	cmd := &cobra.Command{
		Use:   "generate",
		Short: "Generate a Radius application definition from discovery results",
		Long: `Generate a Radius application definition (app.bicep) from discovery results.

This command reads the discovery output from 'rad app discover' and generates a Bicep
file that defines the Radius application with containers and resources.

Examples:
  # Generate app.bicep from default discovery results
  rad app generate

  # Generate with a custom app name
  rad app generate --name myapp

  # Generate from custom discovery file
  rad app generate --path .discovery/custom-result.json

  # Generate to a custom output file
  rad app generate --output infra/app.bicep

  # Skip all prompts (use defaults)
  rad app generate --accept-defaults
`,
		Example: `# Generate app.bicep from discovery results
rad app generate

# Generate with custom name and output
rad app generate --name myapp --output infra/app.bicep`,
		RunE: framework.RunCommand(runner),
	}

	cmd.Flags().StringP("path", "p", "", "Path to discovery result file (default: .discovery/result.json)")
	cmd.Flags().StringP("output", "o", "", "Output file path (default: app.bicep)")
	cmd.Flags().String("name", "", "Application name (default: inferred from codebase path)")
	cmd.Flags().Bool("accept-defaults", false, "Accept all defaults without prompts")
	cmd.Flags().Bool("force", false, "Overwrite existing files without prompting")

	return cmd, runner
}

// Runner is the runner for the generate command.
type Runner struct {
	Output   output.Interface
	Prompter prompt.Interface

	inputPath      string
	outputPath     string
	appName        string
	acceptDefaults bool
	force          bool
}

// NewRunner creates a new Runner.
func NewRunner(factory framework.Factory) *Runner {
	return &Runner{
		Output:   factory.GetOutput(),
		Prompter: factory.GetPrompter(),
	}
}

// Validate validates the command arguments.
func (r *Runner) Validate(cmd *cobra.Command, args []string) error {
	r.inputPath, _ = cmd.Flags().GetString("path")
	r.outputPath, _ = cmd.Flags().GetString("output")
	r.appName, _ = cmd.Flags().GetString("name")
	r.acceptDefaults, _ = cmd.Flags().GetBool("accept-defaults")
	r.force, _ = cmd.Flags().GetBool("force")

	if r.inputPath == "" {
		r.inputPath = defaultDiscoveryFile
	}

	if r.outputPath == "" {
		r.outputPath = defaultOutputFile
	}

	if _, err := os.Stat(r.inputPath); os.IsNotExist(err) {
		return fmt.Errorf("discovery result file not found: %s\nRun 'rad app discover' first to analyze your codebase", r.inputPath)
	}

	return nil
}

// Run executes the generate command.
func (r *Runner) Run(ctx context.Context) error {
	r.Output.LogInfo("Generating Radius application definition...")

	result, err := r.loadDiscoveryResult(r.inputPath)
	if err != nil {
		return fmt.Errorf("failed to load discovery result: %w", err)
	}

	if r.appName == "" {
		if result.CodebasePath != "" {
			r.appName = filepath.Base(result.CodebasePath)
		} else {
			r.appName = "app"
		}
	}

	if _, err := os.Stat(r.outputPath); err == nil {
		if !r.force && !r.acceptDefaults {
			overwrite, err := r.promptExistingFile(r.outputPath)
			if err != nil {
				return err
			}
			if !overwrite {
				r.Output.LogInfo("Generation cancelled.")
				return nil
			}
		}
	}

	gen := generator.NewBicepGenerator()
	genInput := generator.BicepGeneratorInput{
		AppName:      r.appName,
		Dependencies: result.Dependencies,
		Services:     result.Services,
	}

	genOutput, err := gen.Generate(genInput)
	if err != nil {
		return fmt.Errorf("failed to generate Bicep: %w", err)
	}

	validator := generator.NewValidator()
	valResult := validator.Validate(genOutput.BicepContent)

	if !valResult.Valid {
		for _, e := range valResult.Errors {
			r.Output.LogInfo("Validation error at line %d: %s", e.Line, e.Message)
		}
		return fmt.Errorf("generated Bicep failed validation")
	}

	for _, w := range valResult.Warnings {
		r.Output.LogInfo("Warning: %s", w)
	}

	outDir := filepath.Dir(r.outputPath)
	if outDir != "" && outDir != "." {
		if err := os.MkdirAll(outDir, 0755); err != nil {
			return fmt.Errorf("failed to create output directory: %w", err)
		}
	}

	if err := os.WriteFile(r.outputPath, []byte(genOutput.BicepContent), 0644); err != nil {
		return fmt.Errorf("failed to write output file: %w", err)
	}

	r.Output.LogInfo("")
	r.Output.LogInfo("✅ Generated %s", r.outputPath)
	r.Output.LogInfo("")
	r.Output.LogInfo("Application: %s", r.appName)
	r.Output.LogInfo("Containers:  %d", len(genOutput.AppDefinition.Containers))
	r.Output.LogInfo("Resources:   %d", len(genOutput.AppDefinition.Resources))
	r.Output.LogInfo("Connections: %d", len(genOutput.AppDefinition.Connections))
	r.Output.LogInfo("")
	r.Output.LogInfo("Next steps:")
	r.Output.LogInfo("  1. Review and customize %s", r.outputPath)
	r.Output.LogInfo("  2. Update container images with actual image references")
	r.Output.LogInfo("  3. Deploy with: rad deploy %s", r.outputPath)

	return nil
}

func (r *Runner) loadDiscoveryResult(path string) (*discovery.DiscoveryResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var result discovery.DiscoveryResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, fmt.Errorf("invalid discovery result format: %w", err)
	}

	return &result, nil
}

func (r *Runner) promptExistingFile(path string) (bool, error) {
	message := fmt.Sprintf("File '%s' already exists. Overwrite?", path)
	options := []string{"Yes", "No"}

	response, err := r.Prompter.GetListInput(options, message)
	if err != nil {
		return false, err
	}

	return response == "Yes", nil
}

var _ framework.Runner = (*Runner)(nil)
