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

package services

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
	assert.Equal(t, "discover_services", d.Name())
}

func TestDetector_Description(t *testing.T) {
	d := NewDetector(nil)
	assert.Contains(t, d.Description(), "services")
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

func TestDetector_PatternBasedEntrypoint_Express(t *testing.T) {
	d := NewDetector(nil)
	file := FileInfo{
		RelativePath: "index.js",
		Content:      `const app = express(); app.listen(3000);`,
	}
	result := d.patternBasedEntrypoint(file)
	require.NotNil(t, result)
	assert.Equal(t, "nodejs", result.Language)
	assert.Equal(t, "express", result.Framework)
}

func TestDetector_PatternBasedEntrypoint_Flask(t *testing.T) {
	d := NewDetector(nil)
	file := FileInfo{
		RelativePath: "app.py",
		Content:      `from flask import Flask`,
	}
	result := d.patternBasedEntrypoint(file)
	require.NotNil(t, result)
	assert.Equal(t, "python", result.Language)
	assert.Equal(t, "flask", result.Framework)
}

func TestDetector_PatternBasedEntrypoint_NoMatch(t *testing.T) {
	d := NewDetector(nil)
	file := FileInfo{
		RelativePath: "util.txt",
		Content:      `some text`,
	}
	result := d.patternBasedEntrypoint(file)
	assert.Nil(t, result)
}

func TestDetector_DetectFromDockerfile(t *testing.T) {
	d := NewDetector(nil)
	file := FileInfo{
		RelativePath: "Dockerfile",
		Content:      "FROM node:18\nEXPOSE 3000\nCMD [\"node\", \"app.js\"]",
	}
	result := d.detectFromDockerfile(file)
	require.NotNil(t, result)
	assert.Equal(t, 3000, result.Port)
}

func TestDetector_ExtractPort(t *testing.T) {
	d := NewDetector(nil)

	assert.Equal(t, 3000, d.extractPort("PORT=3000", 8080))
	assert.Equal(t, 5000, d.extractPort("app.listen(5000)", 8080))
	assert.Equal(t, 8080, d.extractPort("no port", 8080))
}

func TestDetector_GenerateServiceName(t *testing.T) {
	d := NewDetector(nil)

	assert.Equal(t, "api", d.generateServiceName("services/api/main.go"))
	assert.Equal(t, "main", d.generateServiceName("main.go"))
}

func TestDetector_Execute_InvalidInput(t *testing.T) {
	d := NewDetector(nil)
	ctx := context.Background()
	_, err := d.Execute(ctx, json.RawMessage(`invalid`))
	require.Error(t, err)
}

func TestWalker_IsSourceFile(t *testing.T) {
	w := NewWalker(nil)

	assert.True(t, w.isSourceFile("app.js"))
	assert.True(t, w.isSourceFile("app.py"))
	assert.True(t, w.isSourceFile("main.go"))
	assert.False(t, w.isSourceFile("readme.txt"))
}

func TestWalker_IsSpecialFile(t *testing.T) {
	w := NewWalker(nil)

	assert.True(t, w.isSpecialFile("Dockerfile"))
	assert.True(t, w.isSpecialFile("main.go"))
	assert.False(t, w.isSpecialFile("random.txt"))
}

func TestEntrypointDetector_DetectFromPackageJSON(t *testing.T) {
	d := NewEntrypointDetector()
	content := `{"main": "index.js", "scripts": {"start": "node index.js"}, "dependencies": {"express": "^4.18.0"}}`
	result := d.DetectFromPackageJSON(content, "package.json")
	require.NotNil(t, result)
	assert.Equal(t, 3000, result.Port)
	assert.Equal(t, "express", result.Framework)
}

func TestEntrypointDetector_DetectFromDockerfile(t *testing.T) {
	d := NewEntrypointDetector()
	content := "FROM node:18\nEXPOSE 3000\nCMD [\"node\", \"server.js\"]"
	result := d.DetectFromDockerfile(content, "Dockerfile")
	require.NotNil(t, result)
	assert.Equal(t, 3000, result.Port)
	assert.Equal(t, "nodejs", result.Language)
}
