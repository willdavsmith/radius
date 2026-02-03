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

package resourcetypes

import (
	"strings"

	"github.com/radius-project/radius/pkg/cli/discovery"
)

// MatchResult contains the result of matching a dependency to a resource type.
type MatchResult struct {
	// Entry is the matched catalog entry.
	Entry *CatalogEntry

	// Confidence is the match confidence score (0.0-1.0).
	Confidence float64

	// MatchType is how the match was made: "exact", "alias", "fuzzy".
	MatchType string
}

// Matcher provides dependency-to-resource-type matching.
type Matcher struct {
	catalog *Catalog
}

// NewMatcher creates a new matcher.
func NewMatcher(catalog *Catalog) *Matcher {
	if catalog == nil {
		catalog = NewCatalog()
	}
	return &Matcher{catalog: catalog}
}

// Match attempts to match a dependency to a resource type.
func (m *Matcher) Match(dep discovery.DetectedDependency) *MatchResult {
	entry, found := m.catalog.Get(dep.Type)
	if found {
		return &MatchResult{
			Entry:      entry,
			Confidence: 1.0,
			MatchType:  "exact",
		}
	}

	normalizedType := m.normalizeType(dep.Type)
	entry, found = m.catalog.Get(normalizedType)
	if found {
		return &MatchResult{
			Entry:      entry,
			Confidence: 0.9,
			MatchType:  "alias",
		}
	}

	if result := m.fuzzyMatch(dep); result != nil {
		return result
	}

	return nil
}

// MatchAll matches multiple dependencies and returns results.
func (m *Matcher) MatchAll(deps []discovery.DetectedDependency) map[string]*MatchResult {
	results := make(map[string]*MatchResult)
	for _, dep := range deps {
		if result := m.Match(dep); result != nil {
			results[dep.Type] = result
		}
	}
	return results
}

// normalizeType normalizes a dependency type for matching.
func (m *Matcher) normalizeType(depType string) string {
	normalized := strings.ToLower(depType)
	normalized = strings.ReplaceAll(normalized, "-", "")
	normalized = strings.ReplaceAll(normalized, "_", "")
	return normalized
}

// fuzzyMatch attempts a fuzzy match against catalog entries.
func (m *Matcher) fuzzyMatch(dep discovery.DetectedDependency) *MatchResult {
	depType := strings.ToLower(dep.Type)

	for _, entry := range m.catalog.List() {
		entryType := strings.ToLower(entry.DependencyType)
		if strings.Contains(depType, entryType) || strings.Contains(entryType, depType) {
			return &MatchResult{
				Entry:      entry,
				Confidence: 0.7,
				MatchType:  "fuzzy",
			}
		}

		for _, alias := range entry.Aliases {
			aliasLower := strings.ToLower(alias)
			if strings.Contains(depType, aliasLower) || strings.Contains(aliasLower, depType) {
				return &MatchResult{
					Entry:      entry,
					Confidence: 0.6,
					MatchType:  "fuzzy",
				}
			}
		}
	}

	return nil
}

// GetSupportedTypes returns all supported dependency types.
func (m *Matcher) GetSupportedTypes() []string {
	return m.catalog.SupportedTypes()
}
