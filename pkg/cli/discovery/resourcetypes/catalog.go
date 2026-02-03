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

// Package resourcetypes provides resource type catalog and generation for discovered dependencies.
package resourcetypes

import (
	"sync"

	"github.com/radius-project/radius/pkg/cli/discovery"
)

// Catalog provides access to known Resource Type definitions.
type Catalog struct {
	entries map[string]*CatalogEntry
	mu      sync.RWMutex
}

// CatalogEntry is a resource type definition in the catalog.
type CatalogEntry struct {
	// DependencyType is the detected dependency type (e.g., "postgresql", "redis").
	DependencyType string

	// ResourceType is the Radius resource type definition.
	ResourceType discovery.ResourceType

	// Aliases are alternative names for this dependency type.
	Aliases []string

	// DefaultRecipe is the default recipe name for this resource type.
	DefaultRecipe string
}

// NewCatalog creates a new resource type catalog.
func NewCatalog() *Catalog {
	c := &Catalog{
		entries: make(map[string]*CatalogEntry),
	}
	c.loadBuiltinEntries()
	return c
}

// Get retrieves a catalog entry by dependency type.
func (c *Catalog) Get(dependencyType string) (*CatalogEntry, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.entries[dependencyType]
	if ok {
		return entry, true
	}

	for _, e := range c.entries {
		for _, alias := range e.Aliases {
			if alias == dependencyType {
				return e, true
			}
		}
	}

	return nil, false
}

// List returns all catalog entries.
func (c *Catalog) List() []*CatalogEntry {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]*CatalogEntry, 0, len(c.entries))
	for _, entry := range c.entries {
		result = append(result, entry)
	}
	return result
}

// Register adds a new entry to the catalog.
func (c *Catalog) Register(entry *CatalogEntry) {
	if entry == nil {
		return
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[entry.DependencyType] = entry
}

// SupportedTypes returns the list of supported dependency types.
func (c *Catalog) SupportedTypes() []string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	types := make([]string, 0, len(c.entries))
	for t := range c.entries {
		types = append(types, t)
	}
	return types
}

// loadBuiltinEntries loads the built-in catalog entries.
func (c *Catalog) loadBuiltinEntries() {
	for _, entry := range builtinCatalogEntries {
		c.entries[entry.DependencyType] = entry
	}
}
