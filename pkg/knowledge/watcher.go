// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package knowledge

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher monitors project files for changes and updates the knowledge graph
type Watcher struct {
	graph      *Graph
	watcher    *fsnotify.Watcher
	projectDir string
	done       chan struct{}
	mu         sync.Mutex
}

// NewWatcher creates a new file system watcher
func NewWatcher(graph *Graph, projectDir string) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &Watcher{
		graph:      graph,
		watcher:    fsWatcher,
		projectDir: projectDir,
		done:       make(chan struct{}),
	}, nil
}

// Start begins watching for file changes
func (w *Watcher) Start(ctx context.Context) error {
	// Add project directory to watcher
	if err := filepath.Walk(w.projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return w.watcher.Add(path)
		}
		return nil
	}); err != nil {
		return fmt.Errorf("failed to add directories to watcher: %w", err)
	}

	// Start watching for changes
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

// watch monitors file system events and updates the knowledge graph
func (w *Watcher) watch(ctx context.Context) {
	// Debounce events to prevent rapid updates
	var debounceTimer *time.Timer
	debounceInterval := 2 * time.Second

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

			// Reset or start debounce timer
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			debounceTimer = time.AfterFunc(debounceInterval, func() {
				w.handleFileChange(event)
			})

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			fmt.Printf("Watcher error: %v\n", err)
		}
	}
}

// handleFileChange processes a file system event and updates the knowledge graph
func (w *Watcher) handleFileChange(event fsnotify.Event) {
	w.mu.Lock()
	defer w.mu.Unlock()

	// Only process write and create events
	if event.Op&(fsnotify.Write|fsnotify.Create) == 0 {
		return
	}

	// Only process Go files
	if filepath.Ext(event.Name) != ".go" {
		return
	}

	// Update knowledge graph (implementation depends on your needs)
	// This is a placeholder for the actual update logic
	fmt.Printf("File changed: %s\n", event.Name)
}
