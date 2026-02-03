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
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Default paths to exclude from analysis.
var defaultExcludePaths = []string{
	"node_modules",
	"vendor",
	".git",
	"__pycache__",
	".venv",
	"venv",
	"dist",
	"build",
	"target",
	".next",
	"coverage",
	".nyc_output",
	".idea",
	".vscode",
	"bin",
	"obj",
	"test",
	"tests",
	"spec",
	"__tests__",
}

// Source file extensions to analyze.
var sourceExtensions = []string{
	".js", ".ts", ".jsx", ".tsx",
	".py",
	".go",
	".java",
	".rb",
	".cs",
	".rs",
	".php",
}

// Special files to analyze (like Dockerfile).
var specialFiles = []string{
	"Dockerfile",
	"docker-compose.yml",
	"docker-compose.yaml",
	"main.go",
	"main.py",
	"app.py",
	"index.js",
	"server.js",
	"app.js",
	"manage.py",
	"config.ru",
	"Program.cs",
	"Main.java",
}

// FileInfo contains information about a file for analysis.
type FileInfo struct {
	Path         string
	RelativePath string
	Content      string
}

// Walker walks a directory and collects source files for analysis.
type Walker struct {
	excludePaths []string
	maxFileSize  int64
}

// NewWalker creates a new Walker with the given exclusion paths.
func NewWalker(excludePaths []string) *Walker {
	paths := defaultExcludePaths
	if len(excludePaths) > 0 {
		paths = append(paths, excludePaths...)
	}
	return &Walker{
		excludePaths: paths,
		maxFileSize:  512 * 1024, // 512KB max file size for source files
	}
}

// WalkSourceFiles walks the directory and returns source files.
func (w *Walker) WalkSourceFiles(root string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip errors and continue
		}

		// Check if directory should be excluded
		if d.IsDir() {
			if w.isExcluded(path, root) {
				return filepath.SkipDir
			}
			return nil
		}

		// Check if file is a source file or special file
		if w.isSourceFile(d.Name()) || w.isSpecialFile(d.Name()) {
			info, err := w.readFileInfo(path, root)
			if err == nil {
				files = append(files, info)
			}
		}

		return nil
	})

	return files, err
}

// isExcluded checks if a path should be excluded.
func (w *Walker) isExcluded(path, root string) bool {
	relPath, err := filepath.Rel(root, path)
	if err != nil {
		return false
	}

	for _, excluded := range w.excludePaths {
		if relPath == excluded || strings.HasPrefix(relPath, excluded+string(filepath.Separator)) {
			return true
		}

		// Check if any path component matches
		parts := strings.Split(relPath, string(filepath.Separator))
		for _, part := range parts {
			if part == excluded {
				return true
			}
		}
	}

	return false
}

// isSourceFile checks if a file is a source file by extension.
func (w *Walker) isSourceFile(name string) bool {
	ext := strings.ToLower(filepath.Ext(name))
	for _, srcExt := range sourceExtensions {
		if ext == srcExt {
			return true
		}
	}
	return false
}

// isSpecialFile checks if a file is a special file like Dockerfile.
func (w *Walker) isSpecialFile(name string) bool {
	lowerName := strings.ToLower(name)
	for _, special := range specialFiles {
		if lowerName == strings.ToLower(special) {
			return true
		}
	}
	return false
}

// readFileInfo reads a file and returns its info.
func (w *Walker) readFileInfo(path, root string) (FileInfo, error) {
	info, err := os.Stat(path)
	if err != nil {
		return FileInfo{}, err
	}

	// Skip files that are too large
	if info.Size() > w.maxFileSize {
		return FileInfo{}, fs.ErrInvalid
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return FileInfo{}, err
	}

	relPath, err := filepath.Rel(root, path)
	if err != nil {
		relPath = path
	}

	return FileInfo{
		Path:         path,
		RelativePath: relPath,
		Content:      string(content),
	}, nil
}
