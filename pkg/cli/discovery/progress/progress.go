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

// Package progress provides progress indication for long-running discovery operations.
package progress

import (
	"fmt"
	"io"
	"sync"
)

// Tracker tracks progress of multi-step operations.
type Tracker struct {
	writer      io.Writer
	totalSteps  int
	currentStep int
	quiet       bool
	mu          sync.Mutex
}

// NewTracker creates a new progress tracker.
func NewTracker(writer io.Writer, totalSteps int) *Tracker {
	return &Tracker{
		writer:     writer,
		totalSteps: totalSteps,
	}
}

// SetQuiet enables or disables quiet mode.
func (t *Tracker) SetQuiet(quiet bool) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.quiet = quiet
}

// Step advances progress and displays a message.
func (t *Tracker) Step(message string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	t.currentStep++
	if !t.quiet {
		fmt.Fprintf(t.writer, "[%d/%d] %s\n", t.currentStep, t.totalSteps, message)
	}
}

// Complete marks the operation as complete.
func (t *Tracker) Complete(message string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.quiet {
		fmt.Fprintf(t.writer, "✅ %s\n", message)
	}
}

// Error displays an error message.
func (t *Tracker) Error(message string) {
	t.mu.Lock()
	defer t.mu.Unlock()

	fmt.Fprintf(t.writer, "❌ %s\n", message)
}
