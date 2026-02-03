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

package serve

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/radius-project/radius/pkg/cli/discovery/skills"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_Initialize(t *testing.T) {
	registry := skills.NewRegistry()
	handler := NewHandler(registry)

	req := &MCPRequest{
		JSONRPC: "2.0",
		ID:      1,
		Method:  "initialize",
	}

	resp := handler.Handle(context.Background(), req)
	require.NotNil(t, resp)
	assert.Nil(t, resp.Error)
	assert.Equal(t, "2.0", resp.JSONRPC)
	assert.Equal(t, 1, resp.ID)

	var result InitializeResult
	err := json.Unmarshal(resp.Result, &result)
	require.NoError(t, err)

	assert.Equal(t, "2024-11-05", result.ProtocolVersion)
	assert.Equal(t, "radius-mcp", result.ServerInfo.Name)
	assert.NotNil(t, result.Capabilities.Tools)
}

func TestHandler_ToolsList(t *testing.T) {
	registry := skills.NewRegistry()
	err := RegisterAllTools(registry)
	require.NoError(t, err)

	handler := NewHandler(registry)

	req := &MCPRequest{
		JSONRPC: "2.0",
		ID:      2,
		Method:  "tools/list",
	}

	resp := handler.Handle(context.Background(), req)
	require.NotNil(t, resp)
	assert.Nil(t, resp.Error)

	var result ToolListResult
	err = json.Unmarshal(resp.Result, &result)
	require.NoError(t, err)

	assert.NotEmpty(t, result.Tools)

	foundDetector := false
	for _, tool := range result.Tools {
		if tool.Name == "discover_dependencies" {
			foundDetector = true
			assert.NotEmpty(t, tool.Description)
			assert.NotNil(t, tool.InputSchema)
		}
	}
	assert.True(t, foundDetector, "expected discover_dependencies tool to be registered")
}

func TestHandler_ToolsCall(t *testing.T) {
	registry := skills.NewRegistry()
	err := RegisterAllTools(registry)
	require.NoError(t, err)

	handler := NewHandler(registry)

	params := ToolCallParams{
		Name:      "generate_app_definition",
		Arguments: json.RawMessage(`{"appName":"test","dependencies":[],"services":[]}`),
	}
	paramsJSON, _ := json.Marshal(params)

	req := &MCPRequest{
		JSONRPC: "2.0",
		ID:      3,
		Method:  "tools/call",
		Params:  paramsJSON,
	}

	resp := handler.Handle(context.Background(), req)
	require.NotNil(t, resp)
	assert.Nil(t, resp.Error)

	var result ToolCallResult
	err = json.Unmarshal(resp.Result, &result)
	require.NoError(t, err)

	assert.False(t, result.IsError)
	assert.NotEmpty(t, result.Content)
	assert.Equal(t, "text", result.Content[0].Type)
}

func TestHandler_ToolsCall_NotFound(t *testing.T) {
	registry := skills.NewRegistry()
	handler := NewHandler(registry)

	params := ToolCallParams{
		Name:      "nonexistent_tool",
		Arguments: json.RawMessage(`{}`),
	}
	paramsJSON, _ := json.Marshal(params)

	req := &MCPRequest{
		JSONRPC: "2.0",
		ID:      4,
		Method:  "tools/call",
		Params:  paramsJSON,
	}

	resp := handler.Handle(context.Background(), req)
	require.NotNil(t, resp)
	assert.Nil(t, resp.Error)

	var result ToolCallResult
	err := json.Unmarshal(resp.Result, &result)
	require.NoError(t, err)

	assert.True(t, result.IsError)
}

func TestHandler_MethodNotFound(t *testing.T) {
	registry := skills.NewRegistry()
	handler := NewHandler(registry)

	req := &MCPRequest{
		JSONRPC: "2.0",
		ID:      5,
		Method:  "unknown/method",
	}

	resp := handler.Handle(context.Background(), req)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, -32601, resp.Error.Code)
}

func TestHandler_NotificationsIgnored(t *testing.T) {
	registry := skills.NewRegistry()
	handler := NewHandler(registry)

	req := &MCPRequest{
		JSONRPC: "2.0",
		Method:  "notifications/initialized",
	}

	resp := handler.Handle(context.Background(), req)
	assert.Nil(t, resp)
}

func TestHandler_InvalidParams(t *testing.T) {
	registry := skills.NewRegistry()
	handler := NewHandler(registry)

	req := &MCPRequest{
		JSONRPC: "2.0",
		ID:      6,
		Method:  "tools/call",
		Params:  json.RawMessage(`invalid json`),
	}

	resp := handler.Handle(context.Background(), req)
	require.NotNil(t, resp)
	assert.NotNil(t, resp.Error)
	assert.Equal(t, -32602, resp.Error.Code)
}
