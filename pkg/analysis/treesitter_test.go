package analysis

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewParser(t *testing.T) {
	parser := NewParser()
	assert.NotNil(t, parser)
	assert.NotNil(t, parser.parsers)
	assert.NotNil(t, parser.queries)
	assert.NotNil(t, parser.language)
}

func TestIsSupportedFile(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected bool
	}{
		{
			name:     "go file",
			path:     "main.go",
			expected: true,
		},
		{
			name:     "javascript file",
			path:     "app.js",
			expected: true,
		},
		{
			name:     "typescript file",
			path:     "service.ts",
			expected: true,
		},
		{
			name:     "python file",
			path:     "script.py",
			expected: true,
		},
		{
			name:     "unsupported file",
			path:     "readme.md",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isSupportedFile(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected Language
	}{
		{
			name:     "go file",
			path:     "main.go",
			expected: Go,
		},
		{
			name:     "javascript file",
			path:     "app.js",
			expected: JavaScript,
		},
		{
			name:     "typescript file",
			path:     "service.ts",
			expected: JavaScript,
		},
		{
			name:     "python file",
			path:     "script.py",
			expected: Python,
		},
		{
			name:     "unsupported file",
			path:     "readme.md",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := detectLanguage(tt.path)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAnalyzeProject(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := os.MkdirTemp("", "treesitter-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := map[string]string{
		"main.go": `package main

import "fmt"

func main() string {
	fmt.Println("Hello, World!")
	return "success"
}`,
		"app.js": `
const express = require('express');
const app = express();

app.get('/', (req, res) => {
	res.send('Hello, World!');
});

app.listen(3000);`,
		"script.py": `
def greet(name):
	print(f"Hello, {name}!")
	return name

if __name__ == "__main__":
	greet("World")`,
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		err := os.WriteFile(path, []byte(content), 0o644)
		assert.NoError(t, err)
	}

	// Run analysis
	parser := NewParser()
	analysis, err := parser.AnalyzeProject(context.Background(), tmpDir)
	assert.NoError(t, err)
	assert.NotNil(t, analysis)

	// Validate Go analysis
	assert.Contains(t, analysis.Functions, "main")
	assert.Contains(t, analysis.Imports, "fmt")

	// Validate JavaScript analysis
	assert.Contains(t, analysis.Dependencies, "express")

	// Validate Python analysis
	assert.Contains(t, analysis.Functions, "greet")
}
