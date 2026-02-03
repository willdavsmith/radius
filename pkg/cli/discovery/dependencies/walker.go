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

package dependencies

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
}

// Dependency manifest files to analyze.
var dependencyFiles = []string{
	"package.json",
	"package-lock.json",
	"requirements.txt",
	"Pipfile",
	"pyproject.toml",
	"go.mod",
	"go.sum",
	"pom.xml",
	"build.gradle",
	"Gemfile",
	"composer.json",
	"Cargo.toml",
	"*.csproj",
	"packages.config",
}

// Configuration files that may contain dependency info.
var configFiles = []string{
	".env",
	".env.*",
	"docker-compose.yml",
	"docker-compose.yaml",
	"*.dockerfile",
	"Dockerfile",
	"*.tf",
	"*.bicep",
	"*.yaml",
	"*.yml",
	"*.json",
}

// FileInfo contains information about a file for analysis.
type FileInfo struct {
	Path         string
	RelativePath string
	Content      string
	Type         string // "dependency", "config", "source"
}

// Walker walks a directory and collects files for analysis.
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
		maxFileSize:  1024 * 1024, // 1MB max file size
	}
}

// WalkDependencyFiles walks the directory and returns dependency-related files.
func (w *Walker) WalkDependencyFiles(root string) ([]FileInfo, error) {
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

		// Check if file matches dependency patterns
		if w.matchesDependencyPattern(d.Name()) {
			info, err := w.readFileInfo(path, root, "dependency")
			if err == nil {
				files = append(files, info)
			}
		}

		return nil
	})

	return files, err
}

// WalkConfigFiles walks the directory and returns configuration files.
func (w *Walker) WalkConfigFiles(root string) ([]FileInfo, error) {
	var files []FileInfo

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}

		if d.IsDir() {
			if w.isExcluded(path, root) {
				return filepath.SkipDir
			}
			return nil
		}

		if w.matchesConfigPattern(d.Name()) {
			info, err := w.readFileInfo(path, root, "config")
			if err == nil {
				files = append(files, info)
			}
		}

		return nil
	})

	return files, err
}

// WalkAllRelevant walks and returns all relevant files for analysis.
func (w *Walker) WalkAllRelevant(root string) ([]FileInfo, error) {
	deps, err := w.WalkDependencyFiles(root)
	if err != nil {
		return nil, err
	}

	configs, err := w.WalkConfigFiles(root)
	if err != nil {
		return nil, err
	}

	return append(deps, configs...), nil
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

// matchesDependencyPattern checks if a filename matches dependency file patterns.
func (w *Walker) matchesDependencyPattern(name string) bool {
	for _, pattern := range dependencyFiles {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
		// Exact match for non-wildcard patterns
		if !strings.Contains(pattern, "*") && name == pattern {
			return true
		}
	}
	return false
}

// matchesConfigPattern checks if a filename matches config file patterns.
func (w *Walker) matchesConfigPattern(name string) bool {
	for _, pattern := range configFiles {
		if matched, _ := filepath.Match(pattern, name); matched {
			return true
		}
	}
	return false
}

// readFileInfo reads a file and returns its info.
func (w *Walker) readFileInfo(path, root, fileType string) (FileInfo, error) {
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
		Type:         fileType,
	}, nil
}
