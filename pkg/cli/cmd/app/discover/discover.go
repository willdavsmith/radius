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

package discover

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/radius-project/radius/pkg/cli/cmd/commonflags"
	"github.com/radius-project/radius/pkg/cli/discovery"
	"github.com/radius-project/radius/pkg/cli/discovery/dependencies"
	"github.com/radius-project/radius/pkg/cli/discovery/llm"
	"github.com/radius-project/radius/pkg/cli/discovery/progress"
	"github.com/radius-project/radius/pkg/cli/discovery/services"
	"github.com/radius-project/radius/pkg/cli/framework"
	"github.com/radius-project/radius/pkg/cli/output"
	"github.com/spf13/cobra"
)

const (
	flagPath        = "path"
	flagOutputFile  = "output-file"
	flagConfidence  = "confidence"
	flagExclude     = "exclude"
	flagNoLLM       = "no-llm"
	flagLLMProvider = "llm-provider"
	flagLLMModel    = "llm-model"
)

// NewCommand creates an instance of the `rad app discover` command and runner.
func NewCommand(factory framework.Factory) (*cobra.Command, framework.Runner) {
	runner := NewRunner(factory)

	cmd := &cobra.Command{
		Use:   "discover [path]",
		Short: "Discover application dependencies and services",
		Long: `Analyzes an application codebase to detect infrastructure dependencies and services.

The discover command scans your codebase for:
- Infrastructure dependencies (databases, caches, message queues)
- Application services and entrypoints
- Connection patterns and environment variables

Results are written to discovery.md in the codebase directory.`,
		Args: cobra.MaximumNArgs(1),
		Example: `
# Discover from current directory
rad app discover

# Discover from a specific path
rad app discover ./my-app

# Discover with custom output file
rad app discover --output-file ./analysis.md

# Discover with higher confidence threshold
rad app discover --confidence 0.8

# Discover without LLM (pattern-based only)
rad app discover --no-llm
`,
		RunE: framework.RunCommand(runner),
	}

	cmd.Flags().StringP(flagPath, "p", ".", "Path to codebase directory")
	cmd.Flags().String(flagOutputFile, "", "Output file path (default: discovery.md)")
	cmd.Flags().Float64(flagConfidence, 0.5, "Minimum confidence threshold (0.0-1.0)")
	cmd.Flags().StringSlice(flagExclude, nil, "Paths to exclude from analysis")
	cmd.Flags().Bool(flagNoLLM, false, "Disable LLM-based analysis")
	cmd.Flags().String(flagLLMProvider, "", "LLM provider (openai, ollama)")
	cmd.Flags().String(flagLLMModel, "", "LLM model name")
	commonflags.AddOutputFlag(cmd)

	return cmd, runner
}

// Runner is the Runner implementation for the `rad app discover` command.
type Runner struct {
	Output output.Interface

	CodebasePath        string
	OutputPath          string
	ConfidenceThreshold float64
	ExcludePaths        []string
	NoLLM               bool
	LLMProvider         string
	LLMModel            string
	Format              string
}

// NewRunner creates an instance of the runner for the `rad app discover` command.
func NewRunner(factory framework.Factory) *Runner {
	return &Runner{
		Output: factory.GetOutput(),
	}
}

// Validate runs validation for the `rad app discover` command.
func (r *Runner) Validate(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		r.CodebasePath = args[0]
	} else {
		r.CodebasePath, _ = cmd.Flags().GetString(flagPath)
	}

	absPath, err := filepath.Abs(r.CodebasePath)
	if err != nil {
		return fmt.Errorf("invalid path: %w", err)
	}
	r.CodebasePath = absPath

	info, err := os.Stat(r.CodebasePath)
	if err != nil {
		return fmt.Errorf("path does not exist: %s", r.CodebasePath)
	}
	if !info.IsDir() {
		return fmt.Errorf("path is not a directory: %s", r.CodebasePath)
	}

	r.OutputPath, _ = cmd.Flags().GetString(flagOutputFile)
	if r.OutputPath == "" {
		r.OutputPath = filepath.Join(r.CodebasePath, "discovery.md")
	}

	r.ConfidenceThreshold, _ = cmd.Flags().GetFloat64(flagConfidence)
	if r.ConfidenceThreshold < 0 || r.ConfidenceThreshold > 1 {
		return fmt.Errorf("confidence must be between 0.0 and 1.0")
	}

	r.ExcludePaths, _ = cmd.Flags().GetStringSlice(flagExclude)
	r.NoLLM, _ = cmd.Flags().GetBool(flagNoLLM)
	r.LLMProvider, _ = cmd.Flags().GetString(flagLLMProvider)
	r.LLMModel, _ = cmd.Flags().GetString(flagLLMModel)
	r.Format, _ = cmd.Flags().GetString("output")
	if r.Format == "" {
		r.Format = "table"
	}

	return nil
}

// Run runs the `rad app discover` command.
func (r *Runner) Run(ctx context.Context) error {
	tracker := progress.NewTracker(os.Stdout, 4)
	var warnings []discovery.Warning

	tracker.Step("Initializing discovery analysis")
	r.Output.LogInfo("Analyzing codebase at %s", r.CodebasePath)

	var llmProvider llm.Provider
	if !r.NoLLM {
		var err error
		llmProvider, err = r.createLLMProvider()
		if err != nil {
			r.Output.LogInfo("LLM not available, using pattern-based detection: %s", err)
			warnings = append(warnings, discovery.Warning{
				Code:    "LLM_UNAVAILABLE",
				Message: fmt.Sprintf("LLM provider not available: %s. Using pattern-based detection.", err),
			})
		}
	}

	depDetector := dependencies.NewDetector(llmProvider)
	svcDetector := services.NewDetector(llmProvider)

	tracker.Step("Detecting infrastructure dependencies")
	depInput := dependencies.DetectorInput{
		Path:                r.CodebasePath,
		ExcludePaths:        r.ExcludePaths,
		ConfidenceThreshold: r.ConfidenceThreshold,
	}
	depInputJSON, _ := json.Marshal(depInput)

	var depResult dependencies.DetectorOutput
	depResultJSON, err := depDetector.Execute(ctx, depInputJSON)
	if err != nil {
		r.Output.LogInfo("Warning: Dependency detection encountered issues: %s", err)
		warnings = append(warnings, discovery.Warning{
			Code:    "DEPENDENCY_DETECTION_PARTIAL",
			Message: fmt.Sprintf("Dependency detection had issues: %s", err),
			Source:  "dependencies",
		})
	} else {
		if err := json.Unmarshal(depResultJSON, &depResult); err != nil {
			warnings = append(warnings, discovery.Warning{
				Code:    "DEPENDENCY_PARSE_ERROR",
				Message: fmt.Sprintf("Failed to parse dependency results: %s", err),
				Source:  "dependencies",
			})
		}
	}

	tracker.Step("Detecting services and entrypoints")
	svcInput := services.DetectorInput{
		Path:         r.CodebasePath,
		ExcludePaths: r.ExcludePaths,
	}
	svcInputJSON, _ := json.Marshal(svcInput)

	var svcResult services.DetectorOutput
	svcResultJSON, err := svcDetector.Execute(ctx, svcInputJSON)
	if err != nil {
		r.Output.LogInfo("Warning: Service detection encountered issues: %s", err)
		warnings = append(warnings, discovery.Warning{
			Code:    "SERVICE_DETECTION_PARTIAL",
			Message: fmt.Sprintf("Service detection had issues: %s", err),
			Source:  "services",
		})
	} else {
		if err := json.Unmarshal(svcResultJSON, &svcResult); err != nil {
			warnings = append(warnings, discovery.Warning{
				Code:    "SERVICE_PARSE_ERROR",
				Message: fmt.Sprintf("Failed to parse service results: %s", err),
				Source:  "services",
			})
		}
	}

	tracker.Step("Generating discovery report")

	state := discovery.StateComplete
	if len(depResult.Dependencies) == 0 && len(svcResult.Services) == 0 {
		if len(warnings) > 0 {
			state = discovery.StateFailed
		}
	}

	result := discovery.DiscoveryResult{
		Version:      "1.0",
		AnalyzedAt:   time.Now(),
		CodebasePath: r.CodebasePath,
		Dependencies: depResult.Dependencies,
		Services:     svcResult.Services,
		Warnings:     warnings,
		State:        state,
	}

	if err := discovery.WriteDiscoveryResult(r.OutputPath, result); err != nil {
		return fmt.Errorf("failed to write discovery output: %w", err)
	}

	// Also write JSON format for rad app generate to consume
	jsonPath := filepath.Join(filepath.Dir(r.OutputPath), ".discovery", "result.json")
	if err := os.MkdirAll(filepath.Dir(jsonPath), 0755); err != nil {
		return fmt.Errorf("failed to create .discovery directory: %w", err)
	}
	jsonData, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal discovery result: %w", err)
	}
	if err := os.WriteFile(jsonPath, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write discovery JSON: %w", err)
	}

	tracker.Complete("Discovery complete")

	r.Output.LogInfo("")
	r.Output.LogInfo("Discovery complete!")
	r.Output.LogInfo("  Dependencies found: %d", len(result.Dependencies))
	r.Output.LogInfo("  Services found: %d", len(result.Services))
	if len(warnings) > 0 {
		r.Output.LogInfo("  Warnings: %d", len(warnings))
		for _, w := range warnings {
			r.Output.LogInfo("    - [%s] %s", w.Code, w.Message)
		}
	}
	r.Output.LogInfo("  Output: %s", r.OutputPath)

	return nil
}

func (r *Runner) createLLMProvider() (llm.Provider, error) {
	registry := llm.NewProviderRegistry()

	if r.LLMProvider == "" || r.LLMProvider == "openai" {
		apiKey := os.Getenv("OPENAI_API_KEY")
		if apiKey != "" {
			model := r.LLMModel
			if model == "" {
				model = "gpt-4o-mini"
			}
			openaiConfig := &llm.Config{
				Provider:     "openai",
				Model:        model,
				APIKeyEnvVar: "OPENAI_API_KEY",
			}
			openaiProvider, err := llm.NewOpenAIProvider(openaiConfig)
			if err == nil {
				registry.Register(openaiProvider)
			}
		}
	}

	if r.LLMProvider == "" || r.LLMProvider == "ollama" {
		model := r.LLMModel
		if model == "" {
			model = "llama3.2"
		}
		ollamaConfig := &llm.Config{
			Provider: "ollama",
			Model:    model,
			Endpoint: "http://localhost:11434",
		}
		ollamaProvider := llm.NewOllamaProvider(ollamaConfig)
		registry.Register(ollamaProvider)
	}

	provider := registry.GetPreferred(r.LLMProvider)
	if provider == nil {
		return nil, fmt.Errorf("no LLM provider available")
	}

	return provider, nil
}
