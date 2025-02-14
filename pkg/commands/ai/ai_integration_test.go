package ai

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/Nexlayer/nexlayer-cli/pkg/analysis"
	"github.com/stretchr/testify/assert"
)

func TestGenerateTemplateWithAnalysis(t *testing.T) {
	// Create a temporary test directory
	tmpDir, err := os.MkdirTemp("", "nexlayer-test-*")
	assert.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create test files
	files := map[string]string{
		"go.mod": `module test-app

go 1.22.0

require (
	github.com/lib/pq v1.10.9
)`,
		"main.go": `package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/data", handleData)
	http.ListenAndServe(":8080", nil)
}

func handleData(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}`,
		"database.go": `package main

import (
	"database/sql"
	_ "github.com/lib/pq"
)

var db *sql.DB

func initDB() error {
	var err error
	db, err = sql.Open("postgres", "postgres://user:pass@localhost:5432/app?sslmode=disable")
	return err
}`,
	}

	for name, content := range files {
		path := filepath.Join(tmpDir, name)
		err := os.WriteFile(path, []byte(content), 0o644)
		assert.NoError(t, err)
		if filepath.Base(path) == "go.mod" {
			err = os.Chmod(path, 0o644)
			assert.NoError(t, err)
		}
	}

	// Test template generation
	req := TemplateRequest{
		ProjectName: "test-app",
		ProjectDir:  tmpDir,
	}

	yamlOut, err := GenerateTemplate(context.Background(), req)
	assert.NoError(t, err)
	assert.NotEmpty(t, yamlOut)

	// Verify template contains expected elements
	assert.Contains(t, yamlOut, "test-app")     // Project name
	assert.Contains(t, yamlOut, ":8080")        // Detected port
	assert.Contains(t, yamlOut, "DB_HOST")      // Database env vars
	assert.Contains(t, yamlOut, "postgres")     // Detected database type
	assert.Contains(t, yamlOut, "/api/v1/data") // Detected API endpoint
}

func TestEnhancePromptWithAnalysis(t *testing.T) {
	analysis := &analysis.ProjectAnalysis{
		Frameworks: []string{"gin", "gorm"},
		APIEndpoints: []analysis.APIEndpoint{
			{
				Method:  "GET",
				Path:    "/api/v1/users",
				Handler: "handleUsers",
			},
		},
		DatabaseTypes: []string{"postgresql"},
	}

	basePrompt := "Generate a template"
	enhanced := enhancePromptWithAnalysis(basePrompt, analysis)

	assert.Contains(t, enhanced, "Generate a template")
	assert.Contains(t, enhanced, "Frameworks:")
	assert.Contains(t, enhanced, "gin")
	assert.Contains(t, enhanced, "gorm")
	assert.Contains(t, enhanced, "API Endpoints:")
	assert.Contains(t, enhanced, "GET /api/v1/users")
	assert.Contains(t, enhanced, "Databases:")
	assert.Contains(t, enhanced, "postgresql")
}

func TestGenerateDatabaseEnvVars(t *testing.T) {
	tests := []struct {
		name     string
		dbTypes  []string
		wantVars int
		wantKeys []string
	}{
		{
			name:     "postgresql",
			dbTypes:  []string{"postgresql"},
			wantVars: 5,
			wantKeys: []string{"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD"},
		},
		{
			name:     "mysql",
			dbTypes:  []string{"mysql"},
			wantVars: 5,
			wantKeys: []string{"DB_HOST", "DB_PORT", "DB_NAME", "DB_USER", "DB_PASSWORD"},
		},
		{
			name:     "mongodb",
			dbTypes:  []string{"mongodb"},
			wantVars: 1,
			wantKeys: []string{"MONGODB_URI"},
		},
		{
			name:     "multiple databases",
			dbTypes:  []string{"postgresql", "mongodb"},
			wantVars: 6,
			wantKeys: []string{"DB_HOST", "MONGODB_URI"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vars := generateDatabaseEnvVars(tt.dbTypes)
			assert.Len(t, vars, tt.wantVars)
			for _, key := range tt.wantKeys {
				found := false
				for _, v := range vars {
					if v.Key == key {
						found = true
						break
					}
				}
				assert.True(t, found, "Expected to find env var with key %s", key)
			}
		})
	}
}

func TestExtractPortFromEndpoint(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		wantPort int
	}{
		{
			name:     "standard http port",
			path:     "http://localhost:8080/api",
			wantPort: 8080,
		},
		{
			name:     "https port",
			path:     "https://example.com:443",
			wantPort: 443,
		},
		{
			name:     "no port",
			path:     "http://localhost/api",
			wantPort: 0,
		},
		{
			name:     "invalid port",
			path:     "http://localhost:invalid/api",
			wantPort: 0,
		},
		{
			name:     "empty path",
			path:     "",
			wantPort: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			port := extractPortFromEndpoint(tt.path)
			assert.Equal(t, tt.wantPort, port)
		})
	}
}
