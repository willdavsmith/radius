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
	"context"
	"encoding/json"
	"testing"

	"github.com/radius-project/radius/pkg/cli/discovery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewGenerator(t *testing.T) {
	g := NewGenerator(nil)
	require.NotNil(t, g)
	assert.NotNil(t, g.catalog)
}

func TestGenerator_Name(t *testing.T) {
	g := NewGenerator(nil)
	assert.Equal(t, "generate_resource_types", g.Name())
}

func TestGenerator_Description(t *testing.T) {
	g := NewGenerator(nil)
	assert.Contains(t, g.Description(), "Resource Type")
}

func TestGenerator_InputSchema(t *testing.T) {
	g := NewGenerator(nil)
	schema := g.InputSchema()

	var s map[string]interface{}
	err := json.Unmarshal(schema, &s)
	require.NoError(t, err)
	assert.Equal(t, "object", s["type"])
}

func TestGenerator_OutputSchema(t *testing.T) {
	g := NewGenerator(nil)
	schema := g.OutputSchema()

	var s map[string]interface{}
	err := json.Unmarshal(schema, &s)
	require.NoError(t, err)
	assert.Equal(t, "object", s["type"])
}

func TestGenerator_Execute_KnownDependencies(t *testing.T) {
	g := NewGenerator(nil)

	input := GeneratorInput{
		Dependencies: []discovery.DetectedDependency{
			{Type: "postgresql", Technology: "PostgreSQL", Confidence: 0.9},
			{Type: "redis", Technology: "Redis", Confidence: 0.8},
		},
	}
	inputJSON, err := json.Marshal(input)
	require.NoError(t, err)

	outputJSON, err := g.Execute(context.Background(), inputJSON)
	require.NoError(t, err)

	var output GeneratorOutput
	err = json.Unmarshal(outputJSON, &output)
	require.NoError(t, err)

	assert.Len(t, output.ResourceTypes, 2)
	assert.Empty(t, output.UnmappedTypes)

	typeNames := make([]string, 0, len(output.ResourceTypes))
	for _, rt := range output.ResourceTypes {
		typeNames = append(typeNames, rt.Type)
	}
	assert.Contains(t, typeNames, "Applications.Datastores/sqlDatabases")
	assert.Contains(t, typeNames, "Applications.Datastores/redisCaches")
}

func TestGenerator_Execute_UnknownDependency(t *testing.T) {
	g := NewGenerator(nil)

	input := GeneratorInput{
		Dependencies: []discovery.DetectedDependency{
			{Type: "unknown-service", Technology: "Unknown", Confidence: 0.5},
		},
	}
	inputJSON, err := json.Marshal(input)
	require.NoError(t, err)

	outputJSON, err := g.Execute(context.Background(), inputJSON)
	require.NoError(t, err)

	var output GeneratorOutput
	err = json.Unmarshal(outputJSON, &output)
	require.NoError(t, err)

	assert.Len(t, output.ResourceTypes, 1)
	assert.Equal(t, "Applications.Core/extenders", output.ResourceTypes[0].Type)
	assert.NotEmpty(t, output.Warnings)
}

func TestGenerator_Execute_DuplicateDependencies(t *testing.T) {
	g := NewGenerator(nil)

	input := GeneratorInput{
		Dependencies: []discovery.DetectedDependency{
			{Type: "postgresql", Technology: "PostgreSQL", Confidence: 0.9},
			{Type: "postgresql", Technology: "PostgreSQL", Confidence: 0.8},
		},
	}
	inputJSON, err := json.Marshal(input)
	require.NoError(t, err)

	outputJSON, err := g.Execute(context.Background(), inputJSON)
	require.NoError(t, err)

	var output GeneratorOutput
	err = json.Unmarshal(outputJSON, &output)
	require.NoError(t, err)

	assert.Len(t, output.ResourceTypes, 1)
}

func TestGenerator_Execute_InvalidInput(t *testing.T) {
	g := NewGenerator(nil)

	_, err := g.Execute(context.Background(), json.RawMessage(`invalid json`))
	require.Error(t, err)
}

func TestGenerator_Execute_EmptyDependencies(t *testing.T) {
	g := NewGenerator(nil)

	input := GeneratorInput{
		Dependencies: []discovery.DetectedDependency{},
	}
	inputJSON, err := json.Marshal(input)
	require.NoError(t, err)

	outputJSON, err := g.Execute(context.Background(), inputJSON)
	require.NoError(t, err)

	var output GeneratorOutput
	err = json.Unmarshal(outputJSON, &output)
	require.NoError(t, err)

	assert.Empty(t, output.ResourceTypes)
	assert.Empty(t, output.UnmappedTypes)
}
