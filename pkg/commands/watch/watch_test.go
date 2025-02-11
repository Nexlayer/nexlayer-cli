package watch

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewWatchCommand(t *testing.T) {
	cmd := NewWatchCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "watch", cmd.Use)
	assert.Equal(t, "Watch for file changes and auto-redeploy", cmd.Short)
	assert.Equal(t, "Watch your project directory for file changes and automatically redeploy your application.", cmd.Long)
}

func TestWatchProject(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "nexlayer-watch-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Change to the temporary directory
	originalDir, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Create test directories and files
	dirs := []string{
		"src",
		"node_modules",
		".git",
		"vendor",
	}
	for _, dir := range dirs {
		err := os.MkdirAll(dir, 0755)
		assert.NoError(t, err)
	}

	files := []string{
		"src/main.go",
		"src/app.js",
		"node_modules/test.js",
		".git/config",
		"vendor/module.go",
		".DS_Store",
	}
	for _, file := range files {
		err := os.WriteFile(file, []byte("test"), 0644)
		assert.NoError(t, err)
	}

	// Start watching in a goroutine
	done := make(chan bool)
	go func() {
		err := watchProject()
		assert.NoError(t, err)
		done <- true
	}()

	// Give the watcher time to set up
	time.Sleep(100 * time.Millisecond)

	// Test file changes
	testCases := []struct {
		name     string
		file     string
		content  string
		expected bool // true if should trigger redeploy
	}{
		{
			name:     "source file change",
			file:     "src/main.go",
			content:  "package main\n\nfunc main() {}\n",
			expected: true,
		},
		{
			name:     "node_modules change",
			file:     "node_modules/test.js",
			content:  "console.log('test')",
			expected: false,
		},
		{
			name:     "git file change",
			file:     ".git/config",
			content:  "[core]\n\tbare = false\n",
			expected: false,
		},
		{
			name:     "vendor file change",
			file:     "vendor/module.go",
			content:  "package vendor\n",
			expected: false,
		},
		{
			name:     "DS_Store change",
			file:     ".DS_Store",
			content:  "test",
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Write to the file
			err := os.WriteFile(filepath.Join(tmpDir, tc.file), []byte(tc.content), 0644)
			assert.NoError(t, err)

			// Give the watcher time to process
			time.Sleep(100 * time.Millisecond)
		})
	}

	// Clean up
	close(done)
}
