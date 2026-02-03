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
	"encoding/json"
	"path/filepath"
	"strconv"
	"strings"
)

// EntrypointDetector detects application entrypoints from various sources.
type EntrypointDetector struct{}

// NewEntrypointDetector creates a new entrypoint detector.
func NewEntrypointDetector() *EntrypointDetector {
	return &EntrypointDetector{}
}

// EntrypointInfo contains detected entrypoint information.
type EntrypointInfo struct {
	Type       string
	Entrypoint string
	Port       int
	Language   string
	Framework  string
	Confidence float64
}

// DetectFromDockerfile extracts entrypoint information from a Dockerfile.
func (d *EntrypointDetector) DetectFromDockerfile(content, filePath string) *EntrypointInfo {
	lines := strings.Split(content, "\n")
	var port int
	var cmd string
	var workdir string

	for _, line := range lines {
		line = strings.TrimSpace(line)
		upperLine := strings.ToUpper(line)

		// Parse EXPOSE
		if strings.HasPrefix(upperLine, "EXPOSE") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				portStr := strings.Split(parts[1], "/")[0]
				if p, err := strconv.Atoi(portStr); err == nil && p > 0 {
					port = p
				}
			}
		}

		// Parse CMD
		if strings.HasPrefix(upperLine, "CMD") {
			cmd = line[3:]
		}

		// Parse WORKDIR
		if strings.HasPrefix(upperLine, "WORKDIR") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				workdir = parts[1]
			}
		}
	}

	if port == 0 && cmd == "" {
		return nil
	}

	// Determine language from CMD
	language := d.detectLanguageFromCmd(cmd)
	_ = workdir // Reserved for future use

	return &EntrypointInfo{
		Type:       "http",
		Entrypoint: filePath,
		Port:       port,
		Language:   language,
		Confidence: 0.7,
	}
}

// detectLanguageFromCmd attempts to detect language from CMD instruction.
func (d *EntrypointDetector) detectLanguageFromCmd(cmd string) string {
	cmd = strings.ToLower(cmd)

	if strings.Contains(cmd, "node") || strings.Contains(cmd, "npm") {
		return "nodejs"
	}
	if strings.Contains(cmd, "python") || strings.Contains(cmd, "gunicorn") || strings.Contains(cmd, "uvicorn") {
		return "python"
	}
	if strings.Contains(cmd, "java") || strings.Contains(cmd, "jar") {
		return "java"
	}
	if strings.Contains(cmd, "dotnet") {
		return "dotnet"
	}
	if strings.Contains(cmd, "ruby") || strings.Contains(cmd, "rails") {
		return "ruby"
	}
	return ""
}

// DetectFromPackageJSON extracts entrypoint information from package.json.
func (d *EntrypointDetector) DetectFromPackageJSON(content, filePath string) *EntrypointInfo {
	var pkg struct {
		Main    string            `json:"main"`
		Scripts map[string]string `json:"scripts"`
	}

	if err := json.Unmarshal([]byte(content), &pkg); err != nil {
		return nil
	}

	// Check for start script
	startScript := pkg.Scripts["start"]
	if startScript == "" && pkg.Main == "" {
		return nil
	}

	entrypoint := pkg.Main
	if entrypoint == "" {
		entrypoint = "index.js"
	}

	// Detect framework from scripts
	framework := d.detectFrameworkFromPackageJSON(content)

	// Detect port from start script
	port := d.extractPortFromScript(startScript)
	if port == 0 {
		port = 3000 // Default Node.js port
	}

	return &EntrypointInfo{
		Type:       "http",
		Entrypoint: filepath.Join(filepath.Dir(filePath), entrypoint),
		Port:       port,
		Language:   "nodejs",
		Framework:  framework,
		Confidence: 0.8,
	}
}

// detectFrameworkFromPackageJSON detects the framework from package.json.
func (d *EntrypointDetector) detectFrameworkFromPackageJSON(content string) string {
	content = strings.ToLower(content)

	frameworks := map[string]string{
		"\"express\"":       "express",
		"\"fastify\"":       "fastify",
		"\"koa\"":           "koa",
		"\"hapi\"":          "hapi",
		"\"nest\":":         "nestjs",
		"\"next\"":          "nextjs",
		"\"nuxt\"":          "nuxtjs",
		"\"@remix-run\"":    "remix",
		"\"gatsby\"":        "gatsby",
		"\"react-scripts\"": "create-react-app",
	}

	for pattern, name := range frameworks {
		if strings.Contains(content, pattern) {
			return name
		}
	}
	return ""
}

// extractPortFromScript extracts port from a script command.
func (d *EntrypointDetector) extractPortFromScript(script string) int {
	// Look for PORT= or --port patterns
	patterns := []string{"PORT=", "--port ", "--port=", "-p "}

	for _, pattern := range patterns {
		idx := strings.Index(script, pattern)
		if idx == -1 {
			continue
		}

		remainder := script[idx+len(pattern):]
		var numStr strings.Builder
		for _, c := range remainder {
			if c >= '0' && c <= '9' {
				numStr.WriteRune(c)
			} else if numStr.Len() > 0 {
				break
			}
		}

		if port, err := strconv.Atoi(numStr.String()); err == nil && port > 0 && port < 65536 {
			return port
		}
	}

	return 0
}

// DetectFromPythonFile detects entrypoint from Python files.
func (d *EntrypointDetector) DetectFromPythonFile(content, filePath string) *EntrypointInfo {
	// Check for common Python web framework patterns
	patterns := []struct {
		Pattern   string
		Framework string
		Port      int
	}{
		{"from flask import", "flask", 5000},
		{"Flask(__name__)", "flask", 5000},
		{"from fastapi import", "fastapi", 8000},
		{"FastAPI()", "fastapi", 8000},
		{"from django", "django", 8000},
		{"application = get_wsgi_application", "django", 8000},
		{"from starlette", "starlette", 8000},
		{"from sanic import", "sanic", 8000},
	}

	for _, p := range patterns {
		if strings.Contains(content, p.Pattern) {
			return &EntrypointInfo{
				Type:       "http",
				Entrypoint: filePath,
				Port:       p.Port,
				Language:   "python",
				Framework:  p.Framework,
				Confidence: 0.7,
			}
		}
	}

	return nil
}

// DetectFromGoFile detects entrypoint from Go files.
func (d *EntrypointDetector) DetectFromGoFile(content, filePath string) *EntrypointInfo {
	// Must be main package
	if !strings.Contains(content, "package main") {
		return nil
	}

	patterns := []struct {
		Pattern   string
		Framework string
		Port      int
	}{
		{"gin.Default()", "gin", 8080},
		{"gin.New()", "gin", 8080},
		{"echo.New()", "echo", 8080},
		{"fiber.New()", "fiber", 3000},
		{"chi.NewRouter()", "chi", 8080},
		{"mux.NewRouter()", "gorilla/mux", 8080},
		{"http.ListenAndServe(", "net/http", 8080},
	}

	for _, p := range patterns {
		if strings.Contains(content, p.Pattern) {
			return &EntrypointInfo{
				Type:       "http",
				Entrypoint: filePath,
				Port:       p.Port,
				Language:   "go",
				Framework:  p.Framework,
				Confidence: 0.8,
			}
		}
	}

	return nil
}

// DetectFromJavaFile detects entrypoint from Java files.
func (d *EntrypointDetector) DetectFromJavaFile(content, filePath string) *EntrypointInfo {
	patterns := []struct {
		Pattern   string
		Framework string
		Port      int
	}{
		{"@SpringBootApplication", "spring-boot", 8080},
		{"SpringApplication.run", "spring-boot", 8080},
		{"public static void main", "java", 8080},
	}

	for _, p := range patterns {
		if strings.Contains(content, p.Pattern) {
			return &EntrypointInfo{
				Type:       "http",
				Entrypoint: filePath,
				Port:       p.Port,
				Language:   "java",
				Framework:  p.Framework,
				Confidence: 0.7,
			}
		}
	}

	return nil
}
