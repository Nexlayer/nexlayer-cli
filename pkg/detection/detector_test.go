package detection

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectDetection(t *testing.T) {
	// Create temp test directory
	testDir, err := os.MkdirTemp("", "project-detection-test")
	require.NoError(t, err)
	defer os.RemoveAll(testDir)

	tests := []struct {
		name     string
		files    map[string]string
		expected ProjectType
	}{
		{
			name: "nextjs_project",
			files: map[string]string{
				"package.json": `{
					"name": "test-next",
					"dependencies": {
						"next": "12.0.0",
						"react": "17.0.2"
					}
				}`,
				"Dockerfile": "FROM node:16\nEXPOSE 3000",
			},
			expected: TypeNextjs,
		},
		{
			name: "react_project",
			files: map[string]string{
				"package.json": `{
					"name": "test-react",
					"dependencies": {
						"react": "17.0.2"
					}
				}`,
			},
			expected: TypeReact,
		},
		{
			name: "node_project",
			files: map[string]string{
				"package.json": `{
					"name": "test-node",
					"dependencies": {
						"express": "4.17.1"
					}
				}`,
			},
			expected: TypeNode,
		},
		{
			name: "python_project",
			files: map[string]string{
				"requirements.txt": "flask==2.0.0\n",
				"app.py":          "from flask import Flask",
			},
			expected: TypePython,
		},
		{
			name: "go_project",
			files: map[string]string{
				"go.mod": "module example.com/test\n\ngo 1.17\n",
				"main.go": "package main",
			},
			expected: TypeGo,
		},
		{
			name: "docker_project",
			files: map[string]string{
				"Dockerfile": "FROM nginx\nEXPOSE 80",
			},
			expected: TypeDockerRaw,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create project directory
			projectDir := filepath.Join(testDir, tt.name)
			err := os.MkdirAll(projectDir, 0755)
			require.NoError(t, err)

			// Create test files
			for name, content := range tt.files {
				err := os.WriteFile(filepath.Join(projectDir, name), []byte(content), 0644)
				require.NoError(t, err)
			}

			// Run detection
			registry := NewDetectorRegistry()
			info, err := registry.DetectProject(projectDir)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, info.Type)
		})
	}
}

func TestTemplateGeneration(t *testing.T) {
	tests := []struct {
		name     string
		info     *ProjectInfo
		contains []string
	}{
		{
			name: "nextjs_template",
			info: &ProjectInfo{
				Type:      TypeNextjs,
				Name:      "test-next",
				Port:      3000,
				HasDocker: true,
			},
			contains: []string{
				"name: test-next",
				"servicePorts:",
				"- 3000",
				"path: /",
			},
		},
		{
			name: "python_template",
			info: &ProjectInfo{
				Type:      TypePython,
				Name:      "test-python",
				Port:      8000,
				HasDocker: true,
			},
			contains: []string{
				"name: test-python",
				"servicePorts:",
				"- 8000",
				"# Python web application",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			yaml, err := GenerateYAML(tt.info)
			require.NoError(t, err)
			for _, s := range tt.contains {
				assert.Contains(t, yaml, s)
			}
		})
	}
}
