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
	"bytes"
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/radius-project/radius/pkg/cli/cmd/mcp/serve"
	"github.com/radius-project/radius/pkg/cli/discovery/skills"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPServer_Initialize(t *testing.T) {
	registry := skills.NewRegistry()
	err := serve.RegisterAllTools(registry)
	require.NoError(t, err)

	handler := serve.NewHandler(registry)

	req := &serve.MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
	}

	resp := handler.Handle(context.Background(), req)
	require.NotNil(t, resp)
	assert.Nil(t, resp.Error)

	var result serve.InitializeResult
	err = json.Unmarshal(resp.Result, &result)
	require.NoError(t, err)

	assert.Equal(t, "2024-11-05", result.ProtocolVersion)
	assert.Equal(t, "radius-mcp", result.ServerInfo.Name)
}

func TestMCPServer_ListTools(t *testing.T) {
	registry := skills.NewRegistry()
	err := serve.RegisterAllTools(registry)
	require.NoError(t, err)

	handler := serve.NewHandler(registry)

	req := &serve.MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "tools/list",
	}

	resp := handler.Handle(context.Background(), req)
	require.NotNil(t, resp)
	assert.Nil(t, resp.Error)

	var result serve.ToolListResult
	err = json.Unmarshal(resp.Result, &result)
	require.NoError(t, err)

	toolNames := make(map[string]bool)
	for _, tool := range result.Tools {
		toolNames[tool.Name] = true
	}

	expectedTools := []string{
		"discover_dependencies",
		"discover_services",
		"generate_resource_types",
		"generate_app_definition",
		"validate_app_definition",
	}

	for _, expected := range expectedTools {
		assert.True(t, toolNames[expected], "expected tool %s to be registered", expected)
	}
}

func TestMCPServer_InvokeGenerateAppDefinition(t *testing.T) {
	registry := skills.NewRegistry()
	err := serve.RegisterAllTools(registry)
	require.NoError(t, err)

	handler := serve.NewHandler(registry)

	params := serve.ToolCallParams{
		Name: "generate_app_definition",
		Arguments: json.RawMessage(`{
			"appName": "test-app",
			"dependencies": [{"type": "redis", "confidence": 0.9}],
			"services": [{"name": "api", "type": "http", "port": 8080, "dependencies": ["redis"]}]
		}`),
	}
	paramsJSON, _ := json.Marshal(params)

	req := &serve.MCPRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/call",
		Params:  paramsJSON,
	}

	resp := handler.Handle(context.Background(), req)
	require.NotNil(t, resp)
	assert.Nil(t, resp.Error)

	var result serve.ToolCallResult
	err = json.Unmarshal(resp.Result, &result)
	require.NoError(t, err)

	assert.False(t, result.IsError)
	assert.NotEmpty(t, result.Content)
	assert.Equal(t, "text", result.Content[0].Type)

	content := result.Content[0].Text
	assert.Contains(t, content, "bicepContent")
	assert.Contains(t, content, "test-app")
}

func TestMCPServer_InvokeValidator(t *testing.T) {
	registry := skills.NewRegistry()
	err := serve.RegisterAllTools(registry)
	require.NoError(t, err)

	handler := serve.NewHandler(registry)

	validBicep := "extension radius\nparam environment string\nresource app 'Applications.Core/applications@2023-10-01-preview' = {\n  name: 'test'\n  properties: { environment: environment }\n}"

	args := map[string]string{"bicepContent": validBicep}
	argsJSON, _ := json.Marshal(args)

	params := serve.ToolCallParams{
		Name:      "validate_app_definition",
		Arguments: argsJSON,
	}
	paramsJSON, _ := json.Marshal(params)

	req := &serve.MCPRequest{
		JSONRPC: "2.0",
		ID:      3,
		Method:  "tools/call",
		Params:  paramsJSON,
	}

	resp := handler.Handle(context.Background(), req)
	require.NotNil(t, resp)
	assert.Nil(t, resp.Error)

	var result serve.ToolCallResult
	err = json.Unmarshal(resp.Result, &result)
	require.NoError(t, err)

	assert.False(t, result.IsError)
	assert.Contains(t, result.Content[0].Text, "valid")
}

func TestMCPServer_StdioTransport(t *testing.T) {
	registry := skills.NewRegistry()
	err := serve.RegisterAllTools(registry)
	require.NoError(t, err)

	handler := serve.NewHandler(registry)

	initReq := `{"jsonrpc":"2.0","id":1,"method":"initialize"}` + "\n"
	input := bytes.NewBufferString(initReq)
	output := &bytes.Buffer{}

	server := serve.NewStdioServerWithIO(handler, input, output)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	go func() {
		_ = server.Serve(ctx)
	}()

	<-ctx.Done()

	outputStr := output.String()
	assert.Contains(t, outputStr, "radius-mcp")
	assert.Contains(t, outputStr, "2024-11-05")
}

func TestMCPServer_EndToEndWorkflow(t *testing.T) {
	registry := skills.NewRegistry()
	err := serve.RegisterAllTools(registry)
	require.NoError(t, err)

	handler := serve.NewHandler(registry)
	ctx := context.Background()

	initReq := &serve.MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
	}
	resp := handler.Handle(ctx, initReq)
	require.NotNil(t, resp)
	assert.Nil(t, resp.Error)

	listReq := &serve.MCPRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
	}
	resp = handler.Handle(ctx, listReq)
	require.NotNil(t, resp)
	assert.Nil(t, resp.Error)

	var listResult serve.ToolListResult
	json.Unmarshal(resp.Result, &listResult)
	assert.GreaterOrEqual(t, len(listResult.Tools), 5)

	generateParams := serve.ToolCallParams{
		Name: "generate_app_definition",
		Arguments: json.RawMessage(`{
			"appName": "myapp",
			"dependencies": [
				{"type": "postgresql", "confidence": 0.9},
				{"type": "redis", "confidence": 0.85}
			],
			"services": [
				{"name": "web", "type": "http", "port": 3000, "dependencies": ["postgresql", "redis"]},
				{"name": "worker", "type": "worker", "dependencies": ["redis"]}
			]
		}`),
	}
	paramsJSON, _ := json.Marshal(generateParams)

	generateReq := &serve.MCPRequest{
		JSONRPC: "2.0",
		ID:      3,
		Method:  "tools/call",
		Params:  paramsJSON,
	}
	resp = handler.Handle(ctx, generateReq)
	require.NotNil(t, resp)
	assert.Nil(t, resp.Error)

	var generateResult serve.ToolCallResult
	json.Unmarshal(resp.Result, &generateResult)
	assert.False(t, generateResult.IsError)
	assert.Contains(t, generateResult.Content[0].Text, "bicepContent")
	assert.Contains(t, generateResult.Content[0].Text, "myapp")
}
