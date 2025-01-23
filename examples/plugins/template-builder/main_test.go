package templatebuilder

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuildTemplate(t *testing.T) {
	// Create a temporary test project
	tmpDir := t.TempDir()

	// Create test files
	files := map[string]string{
		"package.json": `{
			"name": "test-app",
			"version": "1.0.0",
			"dependencies": {
				"express": "^4.17.1",
				"react": "^17.0.2"
			}
		}`,
		"Dockerfile": `FROM node:16
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
CMD ["npm", "start"]`,
	}

	for name, content := range files {
		err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644)
		require.NoError(t, err)
	}

	// Test template generation
	template, err := BuildTemplate(tmpDir)
	require.NoError(t, err)
	assert.NotNil(t, template)

	// Verify template contents
	assert.Equal(t, filepath.Base(tmpDir), template.Name)
	assert.Equal(t, "nodejs", template.Stack.Language)
	assert.Equal(t, "react", template.Stack.Framework)
	assert.NotEmpty(t, template.Services)
}

func TestSaveTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	template := &types.NexlayerTemplate{
		Name:    "test-app",
		Version: "1.0.0",
		Stack: types.ProjectStack{
			Language:  "javascript",
			Framework: "react",
		},
	}

	tests := []struct {
		name     string
		format   string
		wantErr  bool
		validate func(t *testing.T, path string)
	}{
		{
			name:   "YAML format",
			format: "yaml",
			validate: func(t *testing.T, path string) {
				content, err := os.ReadFile(path)
				require.NoError(t, err)
				assert.Contains(t, string(content), "name: test-app")
			},
		},
		{
			name:   "JSON format",
			format: "json",
			validate: func(t *testing.T, path string) {
				content, err := os.ReadFile(path)
				require.NoError(t, err)
				assert.Contains(t, string(content), `"name": "test-app"`)
			},
		},
		{
			name:    "Invalid format",
			format:  "invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := filepath.Join(tmpDir, "template."+tt.format)
			err := SaveTemplate(template, path)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			tt.validate(t, path)
		})
	}
}

func TestTemplateVersioning(t *testing.T) {
	tests := []struct {
		name    string
		version string
		want    string
	}{
		{
			name:    "Major version",
			version: "1.0.0",
			want:    "1.0.0-next",
		},
		{
			name:    "With pre-release",
			version: "1.0.0-alpha",
			want:    "1.0.0-alpha-next",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := incrementVersion(tt.version)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestTemplateDiff(t *testing.T) {
	t1 := &types.NexlayerTemplate{
		Name:    "app",
		Version: "1.0.0",
		Stack: types.ProjectStack{
			Language:  "javascript",
			Framework: "react",
		},
	}

	t2 := &types.NexlayerTemplate{
		Name:    "app",
		Version: "1.1.0",
		Stack: types.ProjectStack{
			Language:  "javascript",
			Framework: "react",
			Database: "mongodb",
		},
	}

	diff := compareTemplates(t1, t2)
	assert.NotEmpty(t, diff)
}
