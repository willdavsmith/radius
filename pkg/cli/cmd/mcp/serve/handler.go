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
	"fmt"

	"github.com/radius-project/radius/pkg/cli/discovery/skills"
)

// MCPRequest represents an incoming MCP request.
type MCPRequest struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params,omitempty"`
}

// MCPResponse represents an outgoing MCP response.
type MCPResponse struct {
	JSONRPC string          `json:"jsonrpc"`
	ID      interface{}     `json:"id"`
	Result  json.RawMessage `json:"result,omitempty"`
	Error   *MCPError       `json:"error,omitempty"`
}

// MCPError represents an MCP error.
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ToolCallParams represents parameters for a tool call.
type ToolCallParams struct {
	Name      string          `json:"name"`
	Arguments json.RawMessage `json:"arguments"`
}

// ToolListResult represents the result of listing tools.
type ToolListResult struct {
	Tools []ToolInfo `json:"tools"`
}

// ToolInfo represents information about a tool.
type ToolInfo struct {
	Name        string          `json:"name"`
	Description string          `json:"description"`
	InputSchema json.RawMessage `json:"inputSchema"`
}

// ToolCallResult represents the result of a tool call.
type ToolCallResult struct {
	Content []ToolContent `json:"content"`
	IsError bool          `json:"isError,omitempty"`
}

// ToolContent represents content from a tool call.
type ToolContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

// InitializeResult represents the result of initialization.
type InitializeResult struct {
	ProtocolVersion string       `json:"protocolVersion"`
	Capabilities    Capabilities `json:"capabilities"`
	ServerInfo      ServerInfo   `json:"serverInfo"`
}

// Capabilities represents server capabilities.
type Capabilities struct {
	Tools *ToolsCapability `json:"tools,omitempty"`
}

// ToolsCapability represents tool capabilities.
type ToolsCapability struct {
	ListChanged bool `json:"listChanged,omitempty"`
}

// ServerInfo represents server information.
type ServerInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

// Handler processes MCP requests.
type Handler struct {
	registry *skills.Registry
}

// NewHandler creates a new MCP handler.
func NewHandler(registry *skills.Registry) *Handler {
	return &Handler{
		registry: registry,
	}
}

// Handle processes an MCP request and returns a response.
func (h *Handler) Handle(ctx context.Context, req *MCPRequest) *MCPResponse {
	switch req.Method {
	case "initialize":
		return h.handleInitialize(req)
	case "tools/list":
		return h.handleToolsList(req)
	case "tools/call":
		return h.handleToolsCall(ctx, req)
	case "notifications/initialized":
		return nil
	default:
		return h.errorResponse(req.ID, -32601, fmt.Sprintf("method not found: %s", req.Method))
	}
}

func (h *Handler) handleInitialize(req *MCPRequest) *MCPResponse {
	result := InitializeResult{
		ProtocolVersion: "2024-11-05",
		Capabilities: Capabilities{
			Tools: &ToolsCapability{
				ListChanged: false,
			},
		},
		ServerInfo: ServerInfo{
			Name:    "radius-mcp",
			Version: "1.0.0",
		},
	}

	resultJSON, _ := json.Marshal(result)
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  resultJSON,
	}
}

func (h *Handler) handleToolsList(req *MCPRequest) *MCPResponse {
	metadata := h.registry.ListMetadata()
	tools := make([]ToolInfo, 0, len(metadata))

	for _, m := range metadata {
		tools = append(tools, ToolInfo{
			Name:        m.Name,
			Description: m.Description,
			InputSchema: m.InputSchema,
		})
	}

	result := ToolListResult{Tools: tools}
	resultJSON, _ := json.Marshal(result)

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  resultJSON,
	}
}

func (h *Handler) handleToolsCall(ctx context.Context, req *MCPRequest) *MCPResponse {
	var params ToolCallParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		return h.errorResponse(req.ID, -32602, "invalid params")
	}

	output, err := h.registry.Execute(ctx, params.Name, params.Arguments)
	if err != nil {
		result := ToolCallResult{
			Content: []ToolContent{
				{Type: "text", Text: err.Error()},
			},
			IsError: true,
		}
		resultJSON, _ := json.Marshal(result)
		return &MCPResponse{
			JSONRPC: "2.0",
			ID:      req.ID,
			Result:  resultJSON,
		}
	}

	result := ToolCallResult{
		Content: []ToolContent{
			{Type: "text", Text: string(output)},
		},
	}
	resultJSON, _ := json.Marshal(result)

	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
		Result:  resultJSON,
	}
}

func (h *Handler) errorResponse(id interface{}, code int, message string) *MCPResponse {
	return &MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}
}
