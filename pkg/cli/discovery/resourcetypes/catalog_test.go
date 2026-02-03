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
	"testing"

	"github.com/radius-project/radius/pkg/cli/discovery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCatalog(t *testing.T) {
	c := NewCatalog()
	require.NotNil(t, c)
	assert.NotEmpty(t, c.entries, "catalog should have builtin entries")
}

func TestCatalog_Get_ExactMatch(t *testing.T) {
	c := NewCatalog()

	entry, found := c.Get("postgresql")
	require.True(t, found)
	assert.Equal(t, "postgresql", entry.DependencyType)
}

func TestCatalog_Get_AliasMatch(t *testing.T) {
	c := NewCatalog()

	entry, found := c.Get("pg")
	require.True(t, found)
	assert.Equal(t, "postgresql", entry.DependencyType)

	entry, found = c.Get("postgres")
	require.True(t, found)
	assert.Equal(t, "postgresql", entry.DependencyType)
}

func TestCatalog_Get_NotFound(t *testing.T) {
	c := NewCatalog()

	_, found := c.Get("unknown-dependency")
	assert.False(t, found)
}

func TestCatalog_List(t *testing.T) {
	c := NewCatalog()

	entries := c.List()
	assert.NotEmpty(t, entries)

	types := make([]string, 0, len(entries))
	for _, e := range entries {
		types = append(types, e.DependencyType)
	}

	assert.Contains(t, types, "postgresql")
	assert.Contains(t, types, "redis")
	assert.Contains(t, types, "mongodb")
}

func TestCatalog_Register(t *testing.T) {
	c := NewCatalog()

	customEntry := &CatalogEntry{
		DependencyType: "custom-db",
		Aliases:        []string{"custom"},
		ResourceType: discovery.ResourceType{
			Name: "Custom Database",
			Type: "Applications.Custom/databases",
		},
	}

	c.Register(customEntry)

	entry, found := c.Get("custom-db")
	require.True(t, found)
	assert.Equal(t, "Custom Database", entry.ResourceType.Name)

	entry, found = c.Get("custom")
	require.True(t, found)
	assert.Equal(t, "custom-db", entry.DependencyType)
}

func TestCatalog_Register_Nil(t *testing.T) {
	c := NewCatalog()
	initialCount := len(c.List())

	c.Register(nil)

	assert.Equal(t, initialCount, len(c.List()))
}

func TestCatalog_SupportedTypes(t *testing.T) {
	c := NewCatalog()

	types := c.SupportedTypes()
	assert.NotEmpty(t, types)
	assert.Contains(t, types, "postgresql")
	assert.Contains(t, types, "redis")
}

func TestMatcher_Match_Exact(t *testing.T) {
	m := NewMatcher(nil)

	dep := discovery.DetectedDependency{
		Type:       "postgresql",
		Technology: "PostgreSQL",
		Confidence: 0.9,
	}

	result := m.Match(dep)
	require.NotNil(t, result)
	assert.Equal(t, "exact", result.MatchType)
	assert.Equal(t, 1.0, result.Confidence)
}

func TestMatcher_Match_Alias(t *testing.T) {
	m := NewMatcher(nil)

	dep := discovery.DetectedDependency{
		Type:       "pg",
		Technology: "PostgreSQL",
		Confidence: 0.8,
	}

	result := m.Match(dep)
	require.NotNil(t, result)
	assert.Equal(t, "postgresql", result.Entry.DependencyType)
}

func TestMatcher_Match_NotFound(t *testing.T) {
	m := NewMatcher(nil)

	dep := discovery.DetectedDependency{
		Type:       "completely-unknown-service",
		Technology: "Unknown",
		Confidence: 0.5,
	}

	result := m.Match(dep)
	assert.Nil(t, result)
}

func TestMatcher_MatchAll(t *testing.T) {
	m := NewMatcher(nil)

	deps := []discovery.DetectedDependency{
		{Type: "postgresql", Technology: "PostgreSQL"},
		{Type: "redis", Technology: "Redis"},
		{Type: "unknown", Technology: "Unknown"},
	}

	results := m.MatchAll(deps)
	assert.Len(t, results, 2)
	assert.Contains(t, results, "postgresql")
	assert.Contains(t, results, "redis")
	assert.NotContains(t, results, "unknown")
}

func TestFallbackGenerator_Generate(t *testing.T) {
	fg := NewFallbackGenerator()

	dep := discovery.DetectedDependency{
		Type:       "custom-database",
		Technology: "CustomDB",
		Confidence: 0.7,
	}

	rt := fg.Generate(dep)
	assert.Equal(t, "CustomDB", rt.Name)
	assert.Equal(t, "Applications.Core/extenders", rt.Type)
	assert.Equal(t, "alpha", rt.Maturity)
	assert.NotNil(t, rt.Schema)
	assert.NotEmpty(t, rt.Outputs)
}

func TestFallbackGenerator_Generate_WithConnectionInfo(t *testing.T) {
	fg := NewFallbackGenerator()

	dep := discovery.DetectedDependency{
		Type:       "custom-cache",
		Technology: "CustomCache",
		ConnectionInfo: &discovery.ConnectionInfo{
			EnvVar:      "CACHE_URL",
			DefaultPort: 6380,
		},
	}

	rt := fg.Generate(dep)
	assert.NotNil(t, rt.Schema)
	assert.Contains(t, rt.Schema.Properties, "connectionEnvVar")
	assert.Contains(t, rt.Schema.Properties, "port")
}
