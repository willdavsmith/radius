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

package discovery

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/radius-project/radius/pkg/cli/discovery"
	"github.com/radius-project/radius/pkg/cli/discovery/generator"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBicepGenerator_NodeJSApp(t *testing.T) {
	tempDir := t.TempDir()

	result := &discovery.DiscoveryResult{
		Version:      "1.0",
		AnalyzedAt:   time.Now(),
		CodebasePath: tempDir,
		Dependencies: []discovery.DetectedDependency{
			{Type: "postgresql", Confidence: 0.9},
			{Type: "redis", Confidence: 0.8},
		},
		Services: []discovery.DetectedService{
			{
				Name:         "api",
				Type:         "http",
				Port:         3000,
				Framework:    "express",
				Language:     "nodejs",
				Dependencies: []string{"postgresql", "redis"},
			},
		},
		State: discovery.StateComplete,
	}

	gen := generator.NewBicepGenerator()
	input := generator.BicepGeneratorInput{
		AppName:      "nodejs-app",
		Dependencies: result.Dependencies,
		Services:     result.Services,
	}

	output, err := gen.Generate(input)
	require.NoError(t, err)
	require.NotNil(t, output)

	assert.Contains(t, output.BicepContent, "extension radius")
	assert.Contains(t, output.BicepContent, "nodejs-app")
	assert.Contains(t, output.BicepContent, "Applications.Core/applications")
	assert.Contains(t, output.BicepContent, "Applications.Core/containers")
	assert.Contains(t, output.BicepContent, "Applications.Datastores/sqlDatabases")
	assert.Contains(t, output.BicepContent, "Applications.Datastores/redisCaches")
	assert.Contains(t, output.BicepContent, "containerPort: 3000")
	assert.Contains(t, output.BicepContent, "connections:")

	assert.Equal(t, "nodejs-app", output.AppDefinition.Name)
	assert.Len(t, output.AppDefinition.Resources, 2)
	assert.Len(t, output.AppDefinition.Containers, 1)
	assert.Len(t, output.AppDefinition.Connections, 2)
}

func TestBicepGenerator_PythonApp(t *testing.T) {
	tempDir := t.TempDir()

	result := &discovery.DiscoveryResult{
		Version:      "1.0",
		AnalyzedAt:   time.Now(),
		CodebasePath: tempDir,
		Dependencies: []discovery.DetectedDependency{
			{Type: "mongodb", Confidence: 0.85},
		},
		Services: []discovery.DetectedService{
			{
				Name:         "flask-api",
				Type:         "http",
				Port:         5000,
				Framework:    "flask",
				Language:     "python",
				Dependencies: []string{"mongodb"},
			},
		},
		State: discovery.StateComplete,
	}

	gen := generator.NewBicepGenerator()
	input := generator.BicepGeneratorInput{
		AppName:      "python-app",
		Dependencies: result.Dependencies,
		Services:     result.Services,
	}

	output, err := gen.Generate(input)
	require.NoError(t, err)
	require.NotNil(t, output)

	assert.Contains(t, output.BicepContent, "python-app")
	assert.Contains(t, output.BicepContent, "Applications.Datastores/mongoDatabases")
	assert.Contains(t, output.BicepContent, "containerPort: 5000")
	assert.Len(t, output.AppDefinition.Connections, 1)
}

func TestBicepGenerator_MultiServiceApp(t *testing.T) {
	tempDir := t.TempDir()

	result := &discovery.DiscoveryResult{
		Version:      "1.0",
		AnalyzedAt:   time.Now(),
		CodebasePath: tempDir,
		Dependencies: []discovery.DetectedDependency{
			{Type: "postgresql", Confidence: 0.9},
			{Type: "rabbitmq", Confidence: 0.85},
		},
		Services: []discovery.DetectedService{
			{
				Name:         "web",
				Type:         "http",
				Port:         8080,
				Dependencies: []string{"postgresql"},
			},
			{
				Name:         "worker",
				Type:         "worker",
				Dependencies: []string{"postgresql", "rabbitmq"},
			},
			{
				Name:         "scheduler",
				Type:         "worker",
				Dependencies: []string{"rabbitmq"},
			},
		},
		State: discovery.StateComplete,
	}

	gen := generator.NewBicepGenerator()
	input := generator.BicepGeneratorInput{
		AppName:      "multi-service-app",
		Dependencies: result.Dependencies,
		Services:     result.Services,
	}

	output, err := gen.Generate(input)
	require.NoError(t, err)
	require.NotNil(t, output)

	assert.Len(t, output.AppDefinition.Containers, 3)
	assert.Len(t, output.AppDefinition.Resources, 2)
	assert.Len(t, output.AppDefinition.Connections, 4)

	webConnections := 0
	workerConnections := 0
	schedulerConnections := 0
	for _, conn := range output.AppDefinition.Connections {
		switch conn.Container {
		case "web":
			webConnections++
		case "worker":
			workerConnections++
		case "scheduler":
			schedulerConnections++
		}
	}
	assert.Equal(t, 1, webConnections)
	assert.Equal(t, 2, workerConnections)
	assert.Equal(t, 1, schedulerConnections)
}

func TestBicepGenerator_WriteAndValidate(t *testing.T) {
	tempDir := t.TempDir()

	discoveryDir := filepath.Join(tempDir, ".discovery")
	err := os.MkdirAll(discoveryDir, 0755)
	require.NoError(t, err)

	result := &discovery.DiscoveryResult{
		Version:      "1.0",
		AnalyzedAt:   time.Now(),
		CodebasePath: tempDir,
		Dependencies: []discovery.DetectedDependency{
			{Type: "redis", Confidence: 0.9},
		},
		Services: []discovery.DetectedService{
			{Name: "api", Type: "http", Port: 8080, Dependencies: []string{"redis"}},
		},
		State: discovery.StateComplete,
	}

	resultData, err := json.MarshalIndent(result, "", "  ")
	require.NoError(t, err)
	resultPath := filepath.Join(discoveryDir, "result.json")
	err = os.WriteFile(resultPath, resultData, 0644)
	require.NoError(t, err)

	gen := generator.NewBicepGenerator()
	input := generator.BicepGeneratorInput{
		AppName:      "test-app",
		Dependencies: result.Dependencies,
		Services:     result.Services,
	}

	output, err := gen.Generate(input)
	require.NoError(t, err)

	bicepPath := filepath.Join(tempDir, "app.bicep")
	err = os.WriteFile(bicepPath, []byte(output.BicepContent), 0644)
	require.NoError(t, err)

	content, err := os.ReadFile(bicepPath)
	require.NoError(t, err)

	validator := generator.NewValidator()
	valResult := validator.Validate(string(content))

	assert.True(t, valResult.Valid, "Generated Bicep should be valid")
	assert.Empty(t, valResult.Errors, "Generated Bicep should have no errors")

	lines := strings.Split(string(content), "\n")
	assert.Greater(t, len(lines), 10, "Generated Bicep should have multiple lines")
}

func TestBicepGenerator_UnknownDependency(t *testing.T) {
	gen := generator.NewBicepGenerator()

	input := generator.BicepGeneratorInput{
		AppName: "test-app",
		Dependencies: []discovery.DetectedDependency{
			{Type: "unknown-database", Confidence: 0.7},
		},
		Services: []discovery.DetectedService{},
	}

	output, err := gen.Generate(input)
	require.NoError(t, err)
	require.NotNil(t, output)

	assert.Len(t, output.AppDefinition.Resources, 1)
	assert.Equal(t, "Applications.Core/extenders", output.AppDefinition.Resources[0].Type)
}

func TestBicepGenerator_EmptyApp(t *testing.T) {
	gen := generator.NewBicepGenerator()

	input := generator.BicepGeneratorInput{
		AppName:      "empty-app",
		Dependencies: []discovery.DetectedDependency{},
		Services:     []discovery.DetectedService{},
	}

	output, err := gen.Generate(input)
	require.NoError(t, err)
	require.NotNil(t, output)

	assert.Contains(t, output.BicepContent, "extension radius")
	assert.Contains(t, output.BicepContent, "empty-app")
	assert.Len(t, output.AppDefinition.Resources, 0)
	assert.Len(t, output.AppDefinition.Containers, 0)
	assert.Len(t, output.AppDefinition.Connections, 0)

	validator := generator.NewValidator()
	valResult := validator.Validate(output.BicepContent)
	assert.True(t, valResult.Valid)
}
