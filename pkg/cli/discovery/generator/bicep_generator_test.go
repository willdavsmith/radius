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

package generator
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

package generator

import (
	"strings"
	"testing"

	"github.com/radius-project/radius/pkg/cli/discovery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBicepGenerator_Name(t *testing.T) {
	gen := NewBicepGenerator()
	assert.Equal(t, "generate_app_definition", gen.Name())
}

func TestBicepGenerator_Description(t *testing.T) {
	gen := NewBicepGenerator()
	assert.Contains(t, gen.Description(), "Radius application definition")
}

func TestBicepGenerator_Generate_BasicApp(t *testing.T) {
	gen := NewBicepGenerator()

	input := BicepGeneratorInput{
		AppName:      "testapp",
		Dependencies: []discovery.DetectedDependency{},
		Services:     []discovery.DetectedService{},
	}

	output, err := gen.Generate(input)
	require.NoError(t, err)
	require.NotNil(t, output)

	assert.Contains(t, output.BicepContent, "extension radius")
	assert.Contains(t, output.BicepContent, "param environment string")
	assert.Contains(t, output.BicepContent, "testapp")
	assert.Equal(t, "testapp", output.AppDefinition.Name)
}

func TestBicepGenerator_Generate_WithDependencies(t *testing.T) {
	gen := NewBicepGenerator()

	input := BicepGeneratorInput{
		AppName: "myapp",
		Dependencies: []discovery.DetectedDependency{
			{Type: "postgresql", Confidence: 0.9},
			{Type: "redis", Confidence: 0.8},
		},
		Services: []discovery.DetectedService{},
	}

	output, err := gen.Generate(input)
	require.NoError(t, err)
	require.NotNil(t, output)

	assert.Len(t, output.AppDefinition.Resources, 2)
	assert.Equal(t, "postgresql", output.AppDefinition.Resources[0].Name)
	assert.Equal(t, "redis", output.AppDefinition.Resources[1].Name)
	assert.Contains(t, output.BicepContent, "Applications.Datastores/sqlDatabases")
	assert.Contains(t, output.BicepContent, "Applications.Datastores/redisCaches")
}

func TestBicepGenerator_Generate_WithServices(t *testing.T) {
	gen := NewBicepGenerator()

	input := BicepGeneratorInput{
		AppName:      "myapp",
		Dependencies: []discovery.DetectedDependency{},
		Services: []discovery.DetectedService{
			{Name: "api", Type: "http", Port: 8080, Framework: "express"},
			{Name: "worker", Type: "worker"},
		},
	}

	output, err := gen.Generate(input)
	require.NoError(t, err)
	require.NotNil(t, output)

	assert.Len(t, output.AppDefinition.Containers, 2)
	assert.Equal(t, "api", output.AppDefinition.Containers[0].Name)
	assert.Equal(t, "worker", output.AppDefinition.Containers[1].Name)
	assert.Contains(t, output.BicepContent, "Applications.Core/containers")
	assert.Contains(t, output.BicepContent, "containerPort: 8080")
}

func TestBicepGenerator_Generate_WithConnections(t *testing.T) {
	gen := NewBicepGenerator()

	input := BicepGeneratorInput{
		AppName: "myapp",
		Dependencies: []discovery.DetectedDependency{
			{Type: "postgresql", Confidence: 0.9},
		},
		Services: []discovery.DetectedService{
			{Name: "api", Type: "http", Port: 8080, Dependencies: []string{"postgresql"}},
		},
	}

	output, err := gen.Generate(input)
	require.NoError(t, err)
	require.NotNil(t, output)

	assert.Len(t, output.AppDefinition.Connections, 1)
	assert.Equal(t, "api", output.AppDefinition.Connections[0].Container)
	assert.Equal(t, "postgresql", output.AppDefinition.Connections[0].Resource)
	assert.Contains(t, output.BicepContent, "connections:")
	assert.Contains(t, output.BicepContent, "source: postgresql.id")
}

func TestBicepGenerator_NormalizeResourceName(t *testing.T) {
	gen := NewBicepGenerator()

	tests := []struct {
		input    string
		expected string
	}{
		{"PostgreSQL", "postgresql"},
		{"Redis-Cache", "redis_cache"},
		{"my.resource.name", "my_resource_name"},
		{"  spaces  ", "spaces"},
		{"CamelCase123", "camelcase123"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := gen.normalizeResourceName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestBicepGenerator_InputOutputSchemas(t *testing.T) {
	gen := NewBicepGenerator()

	inputSchema := gen.InputSchema()
	assert.NotNil(t, inputSchema)
	assert.True(t, strings.Contains(string(inputSchema), "appName"))

	outputSchema := gen.OutputSchema()
	assert.NotNil(t, outputSchema)
	assert.True(t, strings.Contains(string(outputSchema), "bicepContent"))
}

func TestConnectionWirer_WireConnections(t *testing.T) {
	wirer := NewConnectionWirer()

	services := []discovery.DetectedService{
		{Name: "api", Dependencies: []string{"postgresql", "redis"}},
		{Name: "worker", Dependencies: []string{"redis"}},
	}

	resources := []discovery.ResourceDef{
		{Name: "postgresql", Type: "Applications.Datastores/sqlDatabases"},
		{Name: "redis", Type: "Applications.Datastores/redisCaches"},
	}

	connections := wirer.WireConnections(services, resources)

	assert.Len(t, connections, 3)

	apiPostgres := false
	apiRedis := false
	workerRedis := false
	for _, conn := range connections {
		if conn.Container == "api" && conn.Resource == "postgresql" {
			apiPostgres = true
		}
		if conn.Container == "api" && conn.Resource == "redis" {
			apiRedis = true
		}
		if conn.Container == "worker" && conn.Resource == "redis" {
			workerRedis = true
		}
	}
	assert.True(t, apiPostgres, "expected api->postgresql connection")
	assert.True(t, apiRedis, "expected api->redis connection")
	assert.True(t, workerRedis, "expected worker->redis connection")
}

func TestConnectionWirer_WireConnections_NoMatch(t *testing.T) {
	wirer := NewConnectionWirer()

	services := []discovery.DetectedService{
		{Name: "api", Dependencies: []string{"mongodb"}},
	}

	resources := []discovery.ResourceDef{
		{Name: "postgresql", Type: "Applications.Datastores/sqlDatabases"},
	}

	connections := wirer.WireConnections(services, resources)
	assert.Len(t, connections, 0)
}
