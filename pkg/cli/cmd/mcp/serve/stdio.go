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
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sync"
)

// StdioServer implements MCP over stdin/stdout.
type StdioServer struct {
	handler *Handler
	reader  io.Reader
	writer  io.Writer
	mu      sync.Mutex
}

// NewStdioServer creates a new stdio server.
func NewStdioServer(handler *Handler) *StdioServer {
	return &StdioServer{
		handler: handler,
		reader:  os.Stdin,
		writer:  os.Stdout,
	}
}

// NewStdioServerWithIO creates a new stdio server with custom IO.
func NewStdioServerWithIO(handler *Handler, reader io.Reader, writer io.Writer) *StdioServer {
	return &StdioServer{
		handler: handler,
		reader:  reader,
		writer:  writer,
	}
}

// Serve starts the stdio server and processes requests until context is cancelled.
func (s *StdioServer) Serve(ctx context.Context) error {
	scanner := bufio.NewScanner(s.reader)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			if !scanner.Scan() {
				if err := scanner.Err(); err != nil {
					return fmt.Errorf("error reading input: %w", err)
				}
				return nil
			}

			line := scanner.Text()
			if line == "" {
				continue
			}

			var req MCPRequest
			if err := json.Unmarshal([]byte(line), &req); err != nil {
				s.writeError(nil, -32700, "parse error")
				continue
			}

			resp := s.handler.Handle(ctx, &req)
			if resp != nil {
				s.writeResponse(resp)
			}
		}
	}
}

func (s *StdioServer) writeResponse(resp *MCPResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()

	data, err := json.Marshal(resp)
	if err != nil {
		return
	}

	fmt.Fprintln(s.writer, string(data))
}

func (s *StdioServer) writeError(id interface{}, code int, message string) {
	resp := &MCPResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &MCPError{
			Code:    code,
			Message: message,
		},
	}
	s.writeResponse(resp)
}
