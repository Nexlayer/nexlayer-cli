// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package knowledge

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
	"github.com/fsnotify/fsnotify"
	"gopkg.in/yaml.v3"
)

// FileType represents a supported file type
type FileType int

const (
	TypeUnknown FileType = iota
	TypeGo
	TypePython
	TypeJavaScript
	TypeTypeScript
	TypeRust
	TypeConfig   // For yaml, json, toml files
	TypeLockfile // For package manager lock files
)

// fileTypePatterns maps file patterns to their types
var fileTypePatterns = map[string]FileType{
	// Source files
	".go":  TypeGo,
	".py":  TypePython,
	".js":  TypeJavaScript,
	".jsx": TypeJavaScript,
	".ts":  TypeTypeScript,
	".tsx": TypeTypeScript,
	".rs":  TypeRust,

	// Config files
	".yaml": TypeConfig,
	".yml":  TypeConfig,
	".json": TypeConfig,
	".toml": TypeConfig,

	// Lock files
	"go.sum":            TypeLockfile,
	"go.mod":            TypeLockfile,
	"package-lock.json": TypeLockfile,
	"yarn.lock":         TypeLockfile,
	"poetry.lock":       TypeLockfile,
	"Cargo.lock":        TypeLockfile,
}

// Watcher monitors project files for changes and updates the knowledge graph
type Watcher struct {
	graph      *Graph
	watcher    *fsnotify.Watcher
	projectDir string
	done       chan struct{}
	mu         sync.Mutex
	debounce   time.Duration
	analyzer   func(string) (*types.ProjectAnalysis, error) // Function to re-analyze files
}

// NewWatcher creates a new file system watcher
func NewWatcher(graph *Graph, projectDir string, analyzer func(string) (*types.ProjectAnalysis, error)) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	// Default debounce interval
	debounce := 2 * time.Second
	if envDebounce := os.Getenv("NEXLAYER_WATCH_DEBOUNCE"); envDebounce != "" {
		if d, err := time.ParseDuration(envDebounce); err == nil {
			debounce = d
		}
	}

	return &Watcher{
		graph:      graph,
		watcher:    fsWatcher,
		projectDir: projectDir,
		done:       make(chan struct{}),
		debounce:   debounce,
		analyzer:   analyzer,
	}, nil
}

// Start begins watching for file changes
func (w *Watcher) Start(ctx context.Context) error {
	err := filepath.Walk(w.projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return w.watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to add directories to watcher: %w", err)
	}

	go w.watch(ctx)
	return nil
}

// Stop stops watching for file changes
func (w *Watcher) Stop() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	close(w.done)
	return w.watcher.Close()
}

// isRelevantFile checks if the file should be monitored based on its type and path
func isRelevantFile(path string) bool {
	// Ignore hidden files and directories
	if strings.HasPrefix(filepath.Base(path), ".") {
		return false
	}

	// Ignore common build and cache directories
	ignoreDirs := []string{
		"node_modules",
		"venv",
		"__pycache__",
		"target",
		"dist",
		"build",
	}

	for _, dir := range ignoreDirs {
		if strings.Contains(path, dir+string(os.PathSeparator)) {
			return false
		}
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(path))
	if _, ok := fileTypePatterns[ext]; ok {
		return true
	}

	// Check exact filenames for lockfiles and special configs
	basename := filepath.Base(path)
	_, isSpecialFile := fileTypePatterns[basename]
	return isSpecialFile
}

// getFileType determines the type of a file based on its path
func getFileType(path string) FileType {
	ext := strings.ToLower(filepath.Ext(path))
	if fileType, ok := fileTypePatterns[ext]; ok {
		return fileType
	}

	basename := filepath.Base(path)
	if fileType, ok := fileTypePatterns[basename]; ok {
		return fileType
	}

	return TypeUnknown
}

// watch monitors file system events and updates the knowledge graph
func (w *Watcher) watch(ctx context.Context) {
	var debounceTimer *time.Timer
	eventQueue := make(map[string]fsnotify.Event)
	var queueMu sync.Mutex

	processEvents := func() {
		queueMu.Lock()
		defer queueMu.Unlock()

		if len(eventQueue) == 0 {
			return
		}

		// Group events by file type for batch processing
		eventsByType := make(map[FileType][]string)
		for path, event := range eventQueue {
			if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
				continue
			}

			if !isRelevantFile(path) {
				continue
			}

			fileType := getFileType(path)
			eventsByType[fileType] = append(eventsByType[fileType], path)
		}

		// Process each group of files
		for fileType, paths := range eventsByType {
			w.handleFileTypeChanges(fileType, paths)
		}

		// Clear the queue
		eventQueue = make(map[string]fsnotify.Event)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-w.done:
			return
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			queueMu.Lock()
			eventQueue[event.Name] = event
			queueMu.Unlock()

			// Debounce updates
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			debounceTimer = time.AfterFunc(w.debounce, processEvents)

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Watcher error: %v\n", err)
		}
	}
}

// handleFileTypeChanges processes changes for a specific file type
func (w *Watcher) handleFileTypeChanges(fileType FileType, paths []string) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Special handling for config files
	if fileType == TypeConfig {
		for _, path := range paths {
			if strings.HasSuffix(path, "nexlayer.yaml") {
				w.handleNexlayerConfig(path)
				continue
			}
		}
	}

	// Group files by type for batch analysis
	filesByType := make(map[string][]string)
	for _, path := range paths {
		ext := filepath.Ext(path)
		filesByType[ext] = append(filesByType[ext], path)
	}

	// Analyze files in batches by type
	for _, batchPaths := range filesByType {
		for _, path := range batchPaths {
			// Extract metadata only
			analysis, err := w.analyzer(path)
			if err != nil {
				fmt.Printf("Failed to analyze file %s: %v\n", path, err)
				continue
			}

			// Remove any code snippets or sensitive data before updating graph
			sanitizedAnalysis := sanitizeAnalysis(analysis)

			if err := w.graph.UpdateFromAnalysis(context.Background(), sanitizedAnalysis); err != nil {
				fmt.Printf("Failed to update graph for file %s: %v\n", path, err)
				continue
			}

			fmt.Printf("Updated knowledge graph metadata for file: %s\n", path)
		}
	}
}

// sanitizeAnalysis removes code snippets and sensitive data from analysis
func sanitizeAnalysis(analysis *types.ProjectAnalysis) *types.ProjectAnalysis {
	sanitized := &types.ProjectAnalysis{
		Functions:    make(map[string][]types.CodeFunction),
		APIEndpoints: make([]types.APIEndpoint, 0),
		Dependencies: make(map[string][]types.ProjectDependency),
		Imports:      make(map[string][]string),
	}

	// Copy only metadata from functions
	for file, functions := range analysis.Functions {
		sanitized.Functions[file] = make([]types.CodeFunction, 0, len(functions))
		for _, fn := range functions {
			sanitizedFn := types.CodeFunction{
				Name:       fn.Name,
				StartLine:  fn.StartLine,
				EndLine:    fn.EndLine,
				IsExported: fn.IsExported,
			}
			sanitized.Functions[file] = append(sanitized.Functions[file], sanitizedFn)
		}
	}

	// Copy only metadata from API endpoints
	for _, endpoint := range analysis.APIEndpoints {
		sanitizedEndpoint := types.APIEndpoint{
			Method:     endpoint.Method,
			Path:       endpoint.Path,
			Handler:    endpoint.Handler,
			Parameters: endpoint.Parameters,
		}
		sanitized.APIEndpoints = append(sanitized.APIEndpoints, sanitizedEndpoint)
	}

	// Copy dependencies and imports (these are already metadata-only)
	for file, deps := range analysis.Dependencies {
		sanitized.Dependencies[file] = deps
	}
	sanitized.Imports = analysis.Imports

	return sanitized
}

// handleNexlayerConfig processes changes to nexlayer.yaml
func (w *Watcher) handleNexlayerConfig(path string) {
	// Read and parse nexlayer.yaml
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Failed to read nexlayer.yaml: %v\n", err)
		return
	}

	var config struct {
		Application struct {
			Name string `yaml:"name"`
			Pods []struct {
				Name         string   `yaml:"name"`
				Type         string   `yaml:"type"`
				ServicePorts []int    `yaml:"servicePorts"`
				Path         string   `yaml:"path,omitempty"`
				Volumes      []string `yaml:"volumes,omitempty"`
			} `yaml:"pods"`
		} `yaml:"application"`
	}

	if err := yaml.Unmarshal(data, &config); err != nil {
		fmt.Printf("Failed to parse nexlayer.yaml: %v\n", err)
		return
	}

	// Create deployment metadata analysis
	analysis := &types.ProjectAnalysis{
		Dependencies: make(map[string][]types.ProjectDependency),
	}

	// Add pod information as dependencies
	for _, pod := range config.Application.Pods {
		analysis.Dependencies[pod.Name] = []types.ProjectDependency{{
			Name:    pod.Name,
			Type:    pod.Type,
			Version: "latest", // Default version for pods
		}}
	}

	// Update graph with deployment metadata
	if err := w.graph.UpdateFromAnalysis(context.Background(), analysis); err != nil {
		fmt.Printf("Failed to update graph with deployment metadata: %v\n", err)
		return
	}

	fmt.Printf("Updated knowledge graph with deployment configuration from nexlayer.yaml\n")
}
