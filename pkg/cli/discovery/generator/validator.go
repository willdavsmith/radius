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
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

// Validator validates generated Bicep content.
type Validator struct{}

// ValidatorInput is the input for the Validator skill.
type ValidatorInput struct {
	BicepContent string `json:"bicepContent"`
}

// ValidatorOutput is the output from the Validator skill.
type ValidatorOutput struct {
	Valid    bool         `json:"valid"`
	Errors   []ValidError `json:"errors,omitempty"`
	Warnings []string     `json:"warnings,omitempty"`
}

// ValidError represents a validation error.
type ValidError struct {
	Line    int    `json:"line"`
	Message string `json:"message"`
}

// NewValidator creates a new Validator.
func NewValidator() *Validator {
	return &Validator{}
}

// Name returns the skill name.
func (v *Validator) Name() string {
	return "validate_app_definition"
}

// Description returns the skill description.
func (v *Validator) Description() string {
	return "Validates a generated Radius application definition (app.bicep) for correctness."
}

// InputSchema returns the JSON Schema for the input.
func (v *Validator) InputSchema() json.RawMessage {
	return json.RawMessage(`{"type":"object","properties":{"bicepContent":{"type":"string"}},"required":["bicepContent"]}`)
}

// OutputSchema returns the JSON Schema for the output.
func (v *Validator) OutputSchema() json.RawMessage {
	return json.RawMessage(`{"type":"object","properties":{"valid":{"type":"boolean"},"errors":{"type":"array"},"warnings":{"type":"array"}}}`)
}

// Execute runs the validation.
func (v *Validator) Execute(ctx context.Context, input json.RawMessage) (json.RawMessage, error) {
	var in ValidatorInput
	if err := json.Unmarshal(input, &in); err != nil {
		return nil, err
	}
	output := v.Validate(in.BicepContent)
	return json.Marshal(output)
}

// Validate validates the Bicep content.
func (v *Validator) Validate(bicepContent string) *ValidatorOutput {
	output := &ValidatorOutput{
		Valid:    true,
		Errors:   make([]ValidError, 0),
		Warnings: make([]string, 0),
	}

	if strings.TrimSpace(bicepContent) == "" {
		output.Valid = false
		output.Errors = append(output.Errors, ValidError{
			Line:    0,
			Message: "Bicep content is empty",
		})
		return output
	}

	syntaxErrors := v.validateSyntax(bicepContent)
	if len(syntaxErrors) > 0 {
		output.Valid = false
		output.Errors = append(output.Errors, syntaxErrors...)
	}

	semanticErrors, warnings := v.validateSemantics(bicepContent)
	if len(semanticErrors) > 0 {
		output.Valid = false
		output.Errors = append(output.Errors, semanticErrors...)
	}
	output.Warnings = append(output.Warnings, warnings...)

	return output
}

func (v *Validator) validateSyntax(content string) []ValidError {
	errors := make([]ValidError, 0)
	lines := strings.Split(content, "\n")

	braceCount := 0
	for i, line := range lines {
		braceCount += strings.Count(line, "{") - strings.Count(line, "}")
		if braceCount < 0 {
			errors = append(errors, ValidError{
				Line:    i + 1,
				Message: "Unmatched closing brace",
			})
		}
	}
	if braceCount != 0 {
		errors = append(errors, ValidError{
			Line:    len(lines),
			Message: fmt.Sprintf("Unbalanced braces: %d unclosed", braceCount),
		})
	}

	return errors
}

func (v *Validator) validateSemantics(content string) ([]ValidError, []string) {
	errors := make([]ValidError, 0)
	warnings := make([]string, 0)

	if !strings.Contains(content, "extension radius") {
		errors = append(errors, ValidError{
			Line:    1,
			Message: "Missing 'extension radius' declaration",
		})
	}

	if !strings.Contains(content, "param environment string") {
		warnings = append(warnings, "Missing 'param environment string' parameter declaration")
	}

	if strings.Contains(content, "TODO:") {
		warnings = append(warnings, "Content contains TODO placeholders that need to be replaced")
	}

	return errors, warnings
}
