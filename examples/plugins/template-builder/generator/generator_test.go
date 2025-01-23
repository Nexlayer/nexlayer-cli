package generator

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateTemplate(t *testing.T) {
	tests := []struct {
		name           string
		stack          *types.ProjectStack
		files          map[string]string
		expectedConfig map[string]string
		expectedPorts  []types.PortConfig
	}{
		{
			name: "Node.js API with MongoDB",
			stack: &types.ProjectStack{
				Language:  "javascript",
				Framework: "express",
				Database:  "mongodb",
			},
			files: map[string]string{
				"package.json": `{
					"name": "test-api",
					"dependencies": {
						"express": "^4.17.1",
						"mongodb": "^4.0.0"
					}
				}`,
				".env": `
					PORT=3000
					MONGODB_URI=mongodb://localhost:27017
				`,
			},
			expectedConfig: map[string]string{
				"PORT":        "3000",
				"MONGODB_URI": "mongodb://localhost:27017",
			},
			expectedPorts: []types.PortConfig{
				{
					Name:        "http",
					Port:        3000,
					TargetPort:  3000,
					Protocol:    "TCP",
					Host:        false,
					Public:      true,
					Healthcheck: true,
				},
				{
					Name:       "mongodb",
					Port:       27017,
					TargetPort: 27017,
					Protocol:   "TCP",
					Host:       false,
					Public:     false,
				},
			},
		},
		{
			name: "Python Web App",
			stack: &types.ProjectStack{
				Language:  "python",
				Framework: "flask",
				Database:  "postgres",
			},
			files: map[string]string{
				"requirements.txt": `
					flask==2.0.1
					psycopg2==2.9.1
				`,
				".env.example": `
					FLASK_APP=app.py
					FLASK_ENV=development
					DATABASE_URL=postgresql://localhost:5432/db
				`,
			},
			expectedConfig: map[string]string{
				"FLASK_APP":    "app.py",
				"FLASK_ENV":    "development",
				"DATABASE_URL": "postgresql://localhost:5432/db",
			},
			expectedPorts: []types.PortConfig{
				{
					Name:        "http",
					Port:        5000,
					TargetPort:  5000,
					Protocol:    "TCP",
					Host:        false,
					Public:      true,
					Healthcheck: true,
				},
				{
					Name:       "postgres",
					Port:       5432,
					TargetPort: 5432,
					Protocol:   "TCP",
					Host:       false,
					Public:     false,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary test directory
			tmpDir := t.TempDir()

			// Create test files
			for name, content := range tt.files {
				err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644)
				require.NoError(t, err, "Failed to create test file")
			}

			// Generate template
			template, err := GenerateTemplate(tmpDir, tt.stack)
			require.NoError(t, err)

			// Verify template
			assert.NotNil(t, template)
			assert.Equal(t, filepath.Base(tmpDir), template.Name)
			assert.Equal(t, tt.stack.Language, template.Stack.Language)
			assert.Equal(t, tt.stack.Framework, template.Stack.Framework)
			assert.Equal(t, tt.stack.Database, template.Stack.Database)

			// Verify environment variables
			if len(tt.expectedConfig) > 0 {
				assert.Equal(t, tt.expectedConfig, template.Config)
			}

			// Verify ports
			if len(tt.expectedPorts) > 0 {
				require.Len(t, template.Services, 1)
				assert.ElementsMatch(t, tt.expectedPorts, template.Services[0].Ports)
			}
		})
	}
}

func TestGenerateTemplateErrors(t *testing.T) {
	tests := []struct {
		name          string
		stack         *types.ProjectStack
		setupFunc     func(dir string) error
		expectedError string
	}{
		{
			name: "Invalid project directory",
			stack: &types.ProjectStack{
				Language: "javascript",
			},
			setupFunc: func(dir string) error {
				return os.Chmod(dir, 0000)
			},
			expectedError: "permission denied",
		},
		{
			name:          "Nil stack",
			stack:         nil,
			expectedError: "stack cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			defer os.Chmod(tmpDir, 0755) // Restore permissions

			if tt.setupFunc != nil {
				err := tt.setupFunc(tmpDir)
				require.NoError(t, err, "Failed to setup test")
			}

			_, err := GenerateTemplate(tmpDir, tt.stack)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), tt.expectedError)
		})
	}
}

func TestGenerateServiceConfig(t *testing.T) {
	tests := []struct {
		name            string
		stack           *types.ProjectStack
		expectedService *types.Service
	}{
		{
			name: "Node.js API service",
			stack: &types.ProjectStack{
				Language:  "javascript",
				Framework: "express",
			},
			expectedService: &types.Service{
				Name:  "api",
				Image: "node:16-alpine",
				Ports: []types.PortConfig{
					{
						Name:        "http",
						Port:        3000,
						TargetPort:  3000,
						Protocol:    "TCP",
						Host:        false,
						Public:      true,
						Healthcheck: true,
					},
				},
				Resources: types.ResourceRequests{
					CPU:    "100m",
					Memory: "128Mi",
				},
				Healthcheck: &types.HealthcheckConfig{
					Path:     "/health",
					Port:     3000,
					Protocol: "HTTP",
				},
			},
		},
		{
			name: "Python web service",
			stack: &types.ProjectStack{
				Language:  "python",
				Framework: "flask",
			},
			expectedService: &types.Service{
				Name:  "web",
				Image: "python:3.9-slim",
				Ports: []types.PortConfig{
					{
						Name:        "http",
						Port:        5000,
						TargetPort:  5000,
						Protocol:    "TCP",
						Host:        false,
						Public:      true,
						Healthcheck: true,
					},
				},
				Resources: types.ResourceRequests{
					CPU:    "100m",
					Memory: "256Mi",
				},
				Healthcheck: &types.HealthcheckConfig{
					Path:     "/health",
					Port:     5000,
					Protocol: "HTTP",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			service := generateServiceConfig(tt.stack)
			assert.Equal(t, tt.expectedService, service)
		})
	}
}
