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
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"sync"
)

var skillNamePattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

// Registry manages skill registration and lookup.
type Registry struct {
	mu     sync.RWMutex
	skills map[string]Skill
}

// NewRegistry creates a new empty skill registry.
func NewRegistry() *Registry {
	return &Registry{
		skills: make(map[string]Skill),
	}
}

// Register adds a skill to the registry.
func (r *Registry) Register(skill Skill) error {
	if skill == nil {
		return fmt.Errorf("skill cannot be nil")
	}

	name := skill.Name()
	if !skillNamePattern.MatchString(name) {
		return fmt.Errorf("invalid skill name %q: must match pattern %s", name, skillNamePattern.String())
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.skills[name]; exists {
		return fmt.Errorf("skill %q is already registered", name)
	}

	r.skills[name] = skill
	return nil
}

// Get returns a skill by name, or nil if not found.
func (r *Registry) Get(name string) Skill {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.skills[name]
}

// List returns all registered skills sorted by name.
func (r *Registry) List() []Skill {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.skills))
	for name := range r.skills {
		names = append(names, name)
	}
	sort.Strings(names)

	skills := make([]Skill, 0, len(r.skills))
	for _, name := range names {
		skills = append(skills, r.skills[name])
	}
	return skills
}

// ListMetadata returns metadata for all registered skills.
func (r *Registry) ListMetadata() []SkillMetadata {
	skills := r.List()
	metadata := make([]SkillMetadata, 0, len(skills))

	for _, skill := range skills {
		metadata = append(metadata, SkillMetadata{
			Name:         skill.Name(),
			Description:  skill.Description(),
			InputSchema:  skill.InputSchema(),
			OutputSchema: skill.OutputSchema(),
		})
	}

	return metadata
}

// Execute runs a skill by name with the given input.
func (r *Registry) Execute(ctx context.Context, name string, input json.RawMessage) (json.RawMessage, error) {
	skill := r.Get(name)
	if skill == nil {
		return nil, fmt.Errorf("skill %q not found", name)
	}

	return skill.Execute(ctx, input)
}

// Unregister removes a skill from the registry.
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.skills[name]; !exists {
		return fmt.Errorf("skill %q not found", name)
	}

	delete(r.skills, name)
	return nil
}

// Clear removes all skills from the registry.
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.skills = make(map[string]Skill)
}

// Count returns the number of registered skills.
func (r *Registry) Count() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.skills)
}

// defaultRegistry is the global skill registry.
var defaultRegistry = NewRegistry()

// DefaultRegistry returns the global skill registry.
func DefaultRegistry() *Registry {
	return defaultRegistry
}

// RegisterSkill registers a skill with the default registry.
func RegisterSkill(skill Skill) error {
	return defaultRegistry.Register(skill)
}

// GetSkill returns a skill from the default registry.
func GetSkill(name string) Skill {
	return defaultRegistry.Get(name)
}

// ListSkills returns all skills from the default registry.
func ListSkills() []Skill {
	return defaultRegistry.List()
}

// ExecuteSkill executes a skill from the default registry.
func ExecuteSkill(ctx context.Context, name string, input json.RawMessage) (json.RawMessage, error) {
	return defaultRegistry.Execute(ctx, name, input)
}
