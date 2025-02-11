package initcmd

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	cmd := NewCommand()
	assert.NotNil(t, cmd)
	assert.Equal(t, "init [project-name]", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.Contains(t, cmd.Long, "Initialize")
}

func TestIsHiddenFile(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		want     bool
	}{
		{
			name:     "hidden file",
			filename: ".gitignore",
			want:     true,
		},
		{
			name:     "regular file",
			filename: "main.go",
			want:     false,
		},
		{
			name:     "hidden directory",
			filename: ".git",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isHiddenFile(tt.filename)
			assert.Equal(t, tt.want, got)
		})
	}
}

// mockAIProvider implements the AI provider interface for testing
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

func TestInitCommand_Execute(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "nexlayer-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Change to the temporary directory
	originalDir, err := os.Getwd()
	assert.NoError(t, err)
	defer os.Chdir(originalDir)
	os.Chdir(tmpDir)

	// Create a test file
	testFile := filepath.Join(tmpDir, "main.go")
	err = os.WriteFile(testFile, []byte("package main\n\nfunc main() {}\n"), 0644)
	assert.NoError(t, err)

	// Set up mock AI provider
	SetAIProvider(&mockAIProvider{})

	// Test the init command
	cmd := NewCommand()
	cmd.SetArgs([]string{"test-app"})

	// Capture output to prevent terminal clutter during tests
	cmd.SetOut(io.Discard)
	cmd.SetErr(io.Discard)

	err = cmd.Execute()
	assert.NoError(t, err)

	// Verify nexlayer.yaml was created and contains expected content
	content, err := os.ReadFile("nexlayer.yaml")
	assert.NoError(t, err)
	assert.Contains(t, string(content), "test-app")
	assert.Contains(t, string(content), "web")
	assert.Contains(t, string(content), "react")

	// Reset AI provider
	SetAIProvider(nil)
}
