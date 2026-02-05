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

package skills

import (
	"github.com/radius-project/radius/pkg/cli/discovery/dependencies"
	"github.com/radius-project/radius/pkg/cli/discovery/generator"
	"github.com/radius-project/radius/pkg/cli/discovery/llm"
	"github.com/radius-project/radius/pkg/cli/discovery/resourcetypes"
	"github.com/radius-project/radius/pkg/cli/discovery/services"
)

// RegisterBuiltinSkills registers all built-in discovery skills with the given registry.
func RegisterBuiltinSkills(registry *Registry, llmProvider llm.Provider) error {
	// Core detection skills
	depDetector := dependencies.NewDetector(llmProvider)
	if err := registry.Register(depDetector); err != nil {
		return err
	}

	svcDetector := services.NewDetector(llmProvider)
	if err := registry.Register(svcDetector); err != nil {
		return err
	}

	// Resource type mapping
	rtGenerator := resourcetypes.NewGenerator(llmProvider)
	if err := registry.Register(rtGenerator); err != nil {
		return err
	}

	// Bicep generation and validation
	bicepGenerator := generator.NewBicepGenerator()
	if err := registry.Register(bicepGenerator); err != nil {
		return err
	}

	validator := generator.NewValidator()
	if err := registry.Register(validator); err != nil {
		return err
	}

	// High-level AI agent skills
	analyzeCodebase := NewAnalyzeCodebase(registry)
	if err := registry.Register(analyzeCodebase); err != nil {
		return err
	}

	suggestRecipes := NewSuggestRecipes()
	if err := registry.Register(suggestRecipes); err != nil {
		return err
	}

	explainDiscovery := NewExplainDiscovery()
	if err := registry.Register(explainDiscovery); err != nil {
		return err
	}

	generateDeployment := NewGenerateDeployment(registry)
	if err := registry.Register(generateDeployment); err != nil {
		return err
	}

	return nil
}

// RegisterBuiltinSkillsDefault registers all built-in skills with the default registry.
func RegisterBuiltinSkillsDefault(llmProvider llm.Provider) error {
	return RegisterBuiltinSkills(defaultRegistry, llmProvider)
}
