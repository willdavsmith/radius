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

package dependencies

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewDetector(t *testing.T) {
	d := NewDetector(nil)
	require.NotNil(t, d)
	assert.Nil(t, d.llmProvider)
	assert.NotNil(t, d.walker)
}

func TestDetector_Name(t *testing.T) {
	d := NewDetector(nil)
	assert.Equal(t, "discover_dependencies", d.Name())
}

func TestDetector_Description(t *testing.T) {
	d := NewDetector(nil)
	assert.Contains(t, d.Description(), "infrastructure dependencies")
}

func TestDetector_InputSchema(t *testing.T) {
	d := NewDetector(nil)
	schema := d.InputSchema()

	var s map[string]interface{}
	err := json.Unmarshal(schema, &s)
	require.NoError(t, err)
	assert.Equal(t, "object", s["type"])
}

func TestDetector_OutputSchema(t *testing.T) {
	d := NewDetector(nil)
	schema := d.OutputSchema()

	var s map[string]interface{}
	err := json.Unmarshal(schema, &s)
	require.NoError(t, err)
	assert.Equal(t, "object", s["type"])
}

func TestDetector_PatternBasedDetection_PostgreSQL(t *testing.T) {
	d := NewDetector(nil)
	file := FileInfo{
		Path:         "/test/file.js",
		RelativePath: "file.js",
		Content:      `const { Pool } = require("pg");`,
		Type:         "dependency",
	}

	deps := d.patternBasedDetection(file)
	require.NotEmpty(t, deps)

	found := false
	for _, dep := range deps {
		if dep.Type == "postgresql" {
			found = true
			break
		}
	}
	assert.True(t, found, "expected to find postgresql")
}

func TestDetector_PatternBasedDetection_Redis(t *testing.T) {
	d := NewDetector(nil)
	file := FileInfo{
		Path:         "/test/file.js",
		RelativePath: "file.js",
		Content:      `import redis from "redis";`,
		Type:         "dependency",
	}

	deps := d.patternBasedDetection(file)
	require.NotEmpty(t, deps)

	found := false
	for _, dep := range deps {
		if dep.Type == "redis" {
			found = true
			break
		}
	}
	assert.True(t, found, "expected to find redis")
}

func TestDetector_PatternBasedDetection_NoMatch(t *testing.T) {
	d := NewDetector(nil)
	file := FileInfo{
		Path:         "/test/file.js",
		RelativePath: "file.js",
		Content:      `console.log("hello world");`,
		Type:         "dependency",
	}

	deps := d.patternBasedDetection(file)
	assert.Empty(t, deps)
}

func TestDetector_Execute_InvalidInput(t *testing.T) {
	d := NewDetector(nil)
	ctx := context.Background()

	_, err := d.Execute(ctx, json.RawMessage(`invalid json`))
	require.Error(t, err)
}

func TestContainsPattern(t *testing.T) {
	assert.True(t, containsPattern("hello pg world", "pg"))
	assert.True(t, containsPattern(`"pg"`, "pg"))
	assert.True(t, containsPattern("POSTGRES connection", "postgres"))
	assert.False(t, containsPattern("hello world", "redis"))
	assert.False(t, containsPattern("", "pg"))
	assert.False(t, containsPattern("hello", ""))
}

func TestNewWalker(t *testing.T) {
	w := NewWalker(nil)
	require.NotNil(t, w)
	assert.NotEmpty(t, w.excludePaths)

	w2 := NewWalker([]string{"custom", "paths"})
	assert.Contains(t, w2.excludePaths, "custom")
	assert.Contains(t, w2.excludePaths, "paths")
}

func TestWalker_IsExcluded(t *testing.T) {
	w := NewWalker(nil)

	assert.True(t, w.isExcluded("/project/node_modules", "/project"))
	assert.True(t, w.isExcluded("/project/vendor", "/project"))
	assert.True(t, w.isExcluded("/project/.git", "/project"))
	assert.False(t, w.isExcluded("/project/src", "/project"))
}

func TestWalker_MatchesDependencyPattern(t *testing.T) {
	w := NewWalker(nil)

	assert.True(t, w.matchesDependencyPattern("package.json"))
	assert.True(t, w.matchesDependencyPattern("requirements.txt"))
	assert.True(t, w.matchesDependencyPattern("go.mod"))
	assert.False(t, w.matchesDependencyPattern("random.js"))
	assert.True(t, w.matchesDependencyPattern("Gemfile"))
	assert.True(t, w.matchesDependencyPattern("pom.xml"))
}
