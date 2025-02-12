package initcmd

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewCommand ensures the command initializes correctly
func TestNewCommand(t *testing.T) {
	cmd := NewCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "init [project-name]", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.Contains(t, cmd.Long, "Initialize")
}

// isHiddenFile returns true if the file name starts with a dot
func isHiddenFile(name string) bool {
	return len(name) > 0 && name[0] == '.'
}

// TestIsHiddenFile ensures hidden files are correctly detected
func TestIsHiddenFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		expected bool
	}{
		{"Hidden file", ".gitignore", true},
		{"Regular file", "main.go", false},
		{"Hidden directory", ".git", true},
		{"Visible directory", "src", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := isHiddenFile(tt.filename)
			assert.Equal(t, tt.expected, actual)
		})
	}
}

// mockAIProvider implements the AIProvider interface for testing
type mockAIProvider struct{}

func (m *mockAIProvider) GenerateTemplate(projectName string) (string, error) {
	return `application:
  name: test-app
  pods:
    - name: web
      type: react
      image: node:18-alpine
      ports:
        - 3000`, nil
}

// TestInitCommand_Execute tests the full init process
func TestInitCommand_Execute(t *testing.T) {
	tmpDir := t.TempDir() // Auto-cleans temp directory after test

	// Change to temp directory
	originalDir, err := os.Getwd()
	assert.NoError(t, err, "Failed to get current directory")
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Create a test file to simulate an existing project
	testFile := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(testFile, []byte("package main\n\nfunc main() {}\n"), 0644)
	assert.NoError(t, err, "Failed to create test file")

	// Set up mock AI provider
	SetAIProvider(&mockAIProvider{})

	// Capture command output
	var stdout, stderr bytes.Buffer
	cmd := NewCommand()
	cmd.SetArgs([]string{"test-app"})
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	// Execute command
	err = cmd.Execute()
	assert.NoError(t, err, "Command execution failed")

	// Verify nexlayer.yaml was created
	assert.FileExists(t, "nexlayer.yaml", "nexlayer.yaml should be created")

	// Verify contents
	content, err := os.ReadFile("nexlayer.yaml")
	assert.NoError(t, err, "Failed to read nexlayer.yaml")
	assert.Contains(t, string(content), "test-app")
	assert.Contains(t, string(content), "web")
	assert.Contains(t, string(content), "react")

	// Verify no errors in stderr
	assert.Empty(t, stderr.String(), "Expected no errors in stderr")

	// Reset AI provider
	SetAIProvider(nil)
}

// TestInitCommand_Fallback tests if init falls back correctly when AI provider fails
func TestInitCommand_Fallback(t *testing.T) {
	tmpDir := t.TempDir()
	originalDir, err := os.Getwd()
	assert.NoError(t, err, "Failed to get current directory")
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Set AI provider to nil to force fallback
	SetAIProvider(nil)

	// Capture output
	var stdout, stderr bytes.Buffer
	cmd := NewCommand()
	cmd.SetArgs([]string{"fallback-project"})
	cmd.SetOut(&stdout)
	cmd.SetErr(&stderr)

	// Execute command
	err = cmd.Execute()
	assert.NoError(t, err, "Command execution failed")

	// Verify nexlayer.yaml was created
	assert.FileExists(t, "nexlayer.yaml", "Expected nexlayer.yaml to be created")

	// Verify fallback content
	content, err := os.ReadFile("nexlayer.yaml")
	assert.NoError(t, err, "Failed to read nexlayer.yaml")
	assert.Contains(t, string(content), "fallback-project")
	assert.Contains(t, string(content), "nginx:latest", "Expected fallback image nginx:latest")

	// Reset AI provider
	SetAIProvider(nil)
}

// TestInitCommand_InvalidArgs ensures invalid arguments are handled properly
func TestInitCommand_InvalidArgs(t *testing.T) {
	cmd := NewCommand()
	cmd.SetArgs([]string{"too", "many", "args"})

	// Capture output
	var stderr bytes.Buffer
	cmd.SetErr(&stderr)

	err := cmd.Execute()
	assert.Error(t, err, "Expected error for too many arguments")
	assert.Contains(t, stderr.String(), "invalid argument", "Expected error message for invalid arguments")
}
