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

package discovery

import (
	"errors"
	"fmt"
)

// Error codes for discovery operations.
const (
	// ErrCodeInvalidSkillName indicates an invalid skill name format.
	ErrCodeInvalidSkillName = "INVALID_SKILL_NAME"

	// ErrCodeSkillNotFound indicates a skill was not found.
	ErrCodeSkillNotFound = "SKILL_NOT_FOUND"

	// ErrCodeSkillAlreadyRegistered indicates a skill is already registered.
	ErrCodeSkillAlreadyRegistered = "SKILL_ALREADY_REGISTERED"

	// ErrCodeInvalidConfidence indicates an invalid confidence score.
	ErrCodeInvalidConfidence = "INVALID_CONFIDENCE"

	// ErrCodeInvalidServiceName indicates an invalid service name.
	ErrCodeInvalidServiceName = "INVALID_SERVICE_NAME"

	// ErrCodeInvalidPort indicates an invalid port number.
	ErrCodeInvalidPort = "INVALID_PORT"

	// ErrCodeInvalidResourceType indicates an invalid Resource Type format.
	ErrCodeInvalidResourceType = "INVALID_RESOURCE_TYPE"

	// ErrCodeInvalidIaCType indicates an invalid IaC type.
	ErrCodeInvalidIaCType = "INVALID_IAC_TYPE"

	// ErrCodeInvalidPriority indicates an invalid priority value.
	ErrCodeInvalidPriority = "INVALID_PRIORITY"

	// ErrCodeLLMUnavailable indicates the LLM provider is unavailable.
	ErrCodeLLMUnavailable = "LLM_UNAVAILABLE"

	// ErrCodeLLMError indicates an error from the LLM provider.
	ErrCodeLLMError = "LLM_ERROR"

	// ErrCodeRecipeSourceUnavailable indicates a recipe source is unavailable.
	ErrCodeRecipeSourceUnavailable = "RECIPE_SOURCE_UNAVAILABLE"

	// ErrCodeRecipeNotFound indicates no recipe was found.
	ErrCodeRecipeNotFound = "RECIPE_NOT_FOUND"

	// ErrCodeBicepValidation indicates a Bicep validation error.
	ErrCodeBicepValidation = "BICEP_VALIDATION"

	// ErrCodeFileNotFound indicates a file was not found.
	ErrCodeFileNotFound = "FILE_NOT_FOUND"

	// ErrCodePermissionDenied indicates permission was denied.
	ErrCodePermissionDenied = "PERMISSION_DENIED"

	// ErrCodeTimeout indicates an operation timed out.
	ErrCodeTimeout = "TIMEOUT"

	// ErrCodePartialFailure indicates some operations failed but others succeeded.
	ErrCodePartialFailure = "PARTIAL_FAILURE"

	// ErrCodeInputValidation indicates input validation failed.
	ErrCodeInputValidation = "INPUT_VALIDATION"
)

// DiscoveryError is the base error type for discovery operations.
type DiscoveryError struct {
	// Code is a machine-readable error code.
	Code string

	// Message is a human-readable error message.
	Message string

	// Cause is the underlying error.
	Cause error

	// Details contains additional error context.
	Details map[string]any
}

// Error implements the error interface.
func (e *DiscoveryError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s: %v", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// Unwrap returns the underlying error.
func (e *DiscoveryError) Unwrap() error {
	return e.Cause
}

// WithDetail adds a detail to the error.
func (e *DiscoveryError) WithDetail(key string, value any) *DiscoveryError {
	if e.Details == nil {
		e.Details = make(map[string]any)
	}
	e.Details[key] = value
	return e
}

// NewError creates a new DiscoveryError.
func NewError(code, message string) *DiscoveryError {
	return &DiscoveryError{
		Code:    code,
		Message: message,
	}
}

// WrapError wraps an existing error with a DiscoveryError.
func WrapError(code, message string, cause error) *DiscoveryError {
	return &DiscoveryError{
		Code:    code,
		Message: message,
		Cause:   cause,
	}
}

// Common sentinel errors.
var (
	// ErrLLMUnavailable indicates the LLM provider is not available.
	ErrLLMUnavailable = NewError(ErrCodeLLMUnavailable, "LLM provider is not available")

	// ErrSkillNotFound indicates a skill was not found in the registry.
	ErrSkillNotFound = NewError(ErrCodeSkillNotFound, "skill not found")

	// ErrRecipeNotFound indicates no matching recipe was found.
	ErrRecipeNotFound = NewError(ErrCodeRecipeNotFound, "no matching recipe found")
)

// IsDiscoveryError checks if an error is a DiscoveryError with the given code.
func IsDiscoveryError(err error, code string) bool {
	var de *DiscoveryError
	if errors.As(err, &de) {
		return de.Code == code
	}
	return false
}

// IsLLMError checks if an error is an LLM-related error.
func IsLLMError(err error) bool {
	return IsDiscoveryError(err, ErrCodeLLMUnavailable) || IsDiscoveryError(err, ErrCodeLLMError)
}

// IsPartialFailure checks if an error indicates partial failure.
func IsPartialFailure(err error) bool {
	return IsDiscoveryError(err, ErrCodePartialFailure)
}

// PartialFailureError represents an operation that partially succeeded.
type PartialFailureError struct {
	*DiscoveryError

	// Warnings lists non-fatal warnings.
	Warnings []Warning

	// FailedOperations lists operations that failed.
	FailedOperations []string

	// SuccessfulOperations lists operations that succeeded.
	SuccessfulOperations []string
}

// NewPartialFailureError creates a new PartialFailureError.
func NewPartialFailureError(message string, warnings []Warning) *PartialFailureError {
	return &PartialFailureError{
		DiscoveryError: &DiscoveryError{
			Code:    ErrCodePartialFailure,
			Message: message,
		},
		Warnings: warnings,
	}
}

// AddFailed adds a failed operation.
func (e *PartialFailureError) AddFailed(op string) *PartialFailureError {
	e.FailedOperations = append(e.FailedOperations, op)
	return e
}

// AddSuccessful adds a successful operation.
func (e *PartialFailureError) AddSuccessful(op string) *PartialFailureError {
	e.SuccessfulOperations = append(e.SuccessfulOperations, op)
	return e
}

// ValidationErrors aggregates multiple validation errors.
type ValidationErrors struct {
	*DiscoveryError
	Errors []ValidationError
}

// NewValidationErrors creates a new ValidationErrors.
func NewValidationErrors(errors []ValidationError) *ValidationErrors {
	return &ValidationErrors{
		DiscoveryError: &DiscoveryError{
			Code:    ErrCodeInputValidation,
			Message: fmt.Sprintf("%d validation error(s)", len(errors)),
		},
		Errors: errors,
	}
}

// Add adds a validation error.
func (e *ValidationErrors) Add(code, message string) {
	e.Errors = append(e.Errors, ValidationError{
		Code:    code,
		Message: message,
	})
	e.Message = fmt.Sprintf("%d validation error(s)", len(e.Errors))
}

// HasErrors returns true if there are validation errors.
func (e *ValidationErrors) HasErrors() bool {
	return len(e.Errors) > 0
}

// BicepValidationError represents a Bicep validation failure.
type BicepValidationError struct {
	*DiscoveryError

	// File is the Bicep file path.
	File string

	// ValidationErrors lists specific validation errors.
	ValidationErrors []ValidationError
}

// NewBicepValidationError creates a new BicepValidationError.
func NewBicepValidationError(file string, errors []ValidationError) *BicepValidationError {
	return &BicepValidationError{
		DiscoveryError: &DiscoveryError{
			Code:    ErrCodeBicepValidation,
			Message: fmt.Sprintf("Bicep validation failed for %s", file),
		},
		File:             file,
		ValidationErrors: errors,
	}
}
