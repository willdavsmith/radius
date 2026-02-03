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
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/radius-project/radius/pkg/cli/discovery/dependencies"
	"github.com/radius-project/radius/pkg/cli/discovery/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDependencyDetector_NodeJSExpress(t *testing.T) {
	testdataPath := filepath.Join("testdata", "nodejs-express")
	if _, err := os.Stat(testdataPath); os.IsNotExist(err) {
		t.Skip("testdata/nodejs-express not found")
	}

	absPath, err := filepath.Abs(testdataPath)
	require.NoError(t, err)

	detector := dependencies.NewDetector(nil)
	input := dependencies.DetectorInput{
		Path:                absPath,
		ConfidenceThreshold: 0.5,
	}
	inputJSON, err := json.Marshal(input)
	require.NoError(t, err)

	resultJSON, err := detector.Execute(context.Background(), inputJSON)
	require.NoError(t, err)

	var result dependencies.DetectorOutput
	err = json.Unmarshal(resultJSON, &result)
	require.NoError(t, err)

	depTypes := make([]string, 0, len(result.Dependencies))
	for _, dep := range result.Dependencies {
		depTypes = append(depTypes, dep.Type)
	}

	assert.Contains(t, depTypes, "postgresql", "should detect PostgreSQL")
	assert.Contains(t, depTypes, "redis", "should detect Redis")
}

func TestDependencyDetector_PythonFlask(t *testing.T) {
	testdataPath := filepath.Join("testdata", "python-flask")
	if _, err := os.Stat(testdataPath); os.IsNotExist(err) {
		t.Skip("testdata/python-flask not found")
	}

	absPath, err := filepath.Abs(testdataPath)
	require.NoError(t, err)

	detector := dependencies.NewDetector(nil)
	input := dependencies.DetectorInput{
		Path:                absPath,
		ConfidenceThreshold: 0.5,
	}
	inputJSON, err := json.Marshal(input)
	require.NoError(t, err)

	resultJSON, err := detector.Execute(context.Background(), inputJSON)
	require.NoError(t, err)

	var result dependencies.DetectorOutput
	err = json.Unmarshal(resultJSON, &result)
	require.NoError(t, err)

	depTypes := make([]string, 0, len(result.Dependencies))
	for _, dep := range result.Dependencies {
		depTypes = append(depTypes, dep.Type)
	}

	assert.Contains(t, depTypes, "postgresql", "should detect PostgreSQL")
	assert.Contains(t, depTypes, "redis", "should detect Redis")
}

func TestDependencyDetector_GoGin(t *testing.T) {
	testdataPath := filepath.Join("testdata", "go-gin")
	if _, err := os.Stat(testdataPath); os.IsNotExist(err) {
		t.Skip("testdata/go-gin not found")
	}

	absPath, err := filepath.Abs(testdataPath)
	require.NoError(t, err)

	detector := dependencies.NewDetector(nil)
	input := dependencies.DetectorInput{
		Path:                absPath,
		ConfidenceThreshold: 0.5,
	}
	inputJSON, err := json.Marshal(input)
	require.NoError(t, err)

	resultJSON, err := detector.Execute(context.Background(), inputJSON)
	require.NoError(t, err)

	var result dependencies.DetectorOutput
	err = json.Unmarshal(resultJSON, &result)
	require.NoError(t, err)

	depTypes := make([]string, 0, len(result.Dependencies))
	for _, dep := range result.Dependencies {
		depTypes = append(depTypes, dep.Type)
	}

	assert.Contains(t, depTypes, "postgresql", "should detect PostgreSQL")
	assert.Contains(t, depTypes, "redis", "should detect Redis")
}

func TestServiceDetector_NodeJSExpress(t *testing.T) {
	testdataPath := filepath.Join("testdata", "nodejs-express")
	if _, err := os.Stat(testdataPath); os.IsNotExist(err) {
		t.Skip("testdata/nodejs-express not found")
	}

	absPath, err := filepath.Abs(testdataPath)
	require.NoError(t, err)

	detector := services.NewDetector(nil)
	input := services.DetectorInput{
		Path: absPath,
	}
	inputJSON, err := json.Marshal(input)
	require.NoError(t, err)

	resultJSON, err := detector.Execute(context.Background(), inputJSON)
	require.NoError(t, err)

	var result services.DetectorOutput
	err = json.Unmarshal(resultJSON, &result)
	require.NoError(t, err)

	assert.NotEmpty(t, result.Services, "should detect at least one service")

	found := false
	for _, svc := range result.Services {
		if svc.Framework == "express" || svc.Type == "http" {
			found = true
			break
		}
	}
	assert.True(t, found, "should detect Express service")
}

func TestServiceDetector_PythonFlask(t *testing.T) {
	testdataPath := filepath.Join("testdata", "python-flask")
	if _, err := os.Stat(testdataPath); os.IsNotExist(err) {
		t.Skip("testdata/python-flask not found")
	}

	absPath, err := filepath.Abs(testdataPath)
	require.NoError(t, err)

	detector := services.NewDetector(nil)
	input := services.DetectorInput{
		Path: absPath,
	}
	inputJSON, err := json.Marshal(input)
	require.NoError(t, err)

	resultJSON, err := detector.Execute(context.Background(), inputJSON)
	require.NoError(t, err)

	var result services.DetectorOutput
	err = json.Unmarshal(resultJSON, &result)
	require.NoError(t, err)

	assert.NotEmpty(t, result.Services, "should detect at least one service")

	found := false
	for _, svc := range result.Services {
		if svc.Framework == "flask" || svc.Type == "http" {
			found = true
			break
		}
	}
	assert.True(t, found, "should detect Flask service")
}

func TestServiceDetector_GoGin(t *testing.T) {
	testdataPath := filepath.Join("testdata", "go-gin")
	if _, err := os.Stat(testdataPath); os.IsNotExist(err) {
		t.Skip("testdata/go-gin not found")
	}

	absPath, err := filepath.Abs(testdataPath)
	require.NoError(t, err)

	detector := services.NewDetector(nil)
	input := services.DetectorInput{
		Path: absPath,
	}
	inputJSON, err := json.Marshal(input)
	require.NoError(t, err)

	resultJSON, err := detector.Execute(context.Background(), inputJSON)
	require.NoError(t, err)

	var result services.DetectorOutput
	err = json.Unmarshal(resultJSON, &result)
	require.NoError(t, err)

	assert.NotEmpty(t, result.Services, "should detect at least one service")

	found := false
	for _, svc := range result.Services {
		if svc.Framework == "gin" || svc.Type == "http" {
			found = true
			break
		}
	}
	assert.True(t, found, "should detect Gin service")
}
