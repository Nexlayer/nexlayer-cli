// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package watch

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/cobra"
)

// NewWatchCommand creates a new watch command for automatic redeployment
func NewWatchCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "watch",
		Short: "Watch for file changes and auto-redeploy",
		Long:  "Watch your project directory for file changes and automatically redeploy your application.",
		RunE: func(cmd *cobra.Command, args []string) error {
			return watchProject()
		},
	}
}

func watchProject() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("failed to create file watcher: %w", err)
	}
	defer watcher.Close()

	// Channel to keep the process running
	done := make(chan bool)
	
	// Channel for debouncing events
	events := make(chan fsnotify.Event)
	
	// Start event debouncer
	go func() {
		var timer *time.Timer
		for range events {
			if timer != nil {
				timer.Stop()
			}
			timer = time.AfterFunc(2*time.Second, func() {
				log.Printf("Changes detected, redeploying...")
				cmd := exec.Command("nexlayer", "deploy")
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					log.Printf("Error: %v", err)
				}
			})
		}
	}()

	// Start watching for events
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Skip common ignored directories and files
				path := event.Name
				if strings.Contains(path, "node_modules") ||
					strings.Contains(path, ".git") ||
					strings.Contains(path, "vendor") ||
					strings.Contains(path, ".DS_Store") {
					continue
				}

				// Only trigger on write or create events
				if event.Op&(fsnotify.Write|fsnotify.Create) != 0 {
					events <- event
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	// Add directories to watch
	log.Println("Setting up file watchers...")
	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			// Skip common ignored directories
			if strings.Contains(path, "node_modules") ||
				strings.Contains(path, ".git") ||
				strings.Contains(path, "vendor") {
				return filepath.SkipDir
			}
			return watcher.Add(path)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("failed to set up watchers: %w", err)
	}

	log.Println("Watching for changes. Press Ctrl+C to exit.")
	<-done
	return nil
}
