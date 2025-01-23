package detector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDetectStack(t *testing.T) {
	tests := []struct {
		name     string
		files    map[string]string
		expected struct {
			Language  string
			Framework string
			Database  string
		}
	}{
		{
			name: "Node.js with React and MongoDB",
			files: map[string]string{
				"package.json": `{
					"dependencies": {
						"react": "^17.0.2",
						"mongodb": "^4.1.0"
					}
				}`,
			},
			expected: struct {
				Language  string
				Framework string
				Database  string
			}{
				Language:  "nodejs",
				Framework: "react",
				Database:  "mongodb",
			},
		},
		{
			name: "Python with Django",
			files: map[string]string{
				"requirements.txt": `
					Django==3.2.0
					psycopg2==2.9.1
				`,
			},
			expected: struct {
				Language  string
				Framework string
				Database  string
			}{
				Language:  "python",
				Framework: "django",
				Database:  "postgres",
			},
		},
		{
			name: "Go with Gin",
			files: map[string]string{
				"go.mod": `
					module myapp
					require github.com/gin-gonic/gin v1.7.0
				`,
			},
			expected: struct {
				Language  string
				Framework string
				Database  string
			}{
				Language:  "go",
				Framework: "gin",
				Database:  "",
			},
		},
		{
			name: "Java with Spring Boot",
			files: map[string]string{
				"pom.xml": `
					<project>
						<dependencies>
							<dependency>
								<groupId>org.springframework.boot</groupId>
								<artifactId>spring-boot-starter-web</artifactId>
							</dependency>
						</dependencies>
					</project>
				`,
			},
			expected: struct {
				Language  string
				Framework string
				Database  string
			}{
				Language:  "java",
				Framework: "spring",
				Database:  "",
			},
		},
		{
			name:  "Empty directory",
			files: map[string]string{},
			expected: struct {
				Language  string
				Framework string
				Database  string
			}{
				Language:  "",
				Framework: "",
				Database:  "",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary directory
			tmpDir, err := os.MkdirTemp("", "test-*")
			require.NoError(t, err)
			defer os.RemoveAll(tmpDir)

			// Create test files
			for path, content := range tt.files {
				fullPath := filepath.Join(tmpDir, path)
				err := os.MkdirAll(filepath.Dir(fullPath), 0755)
				require.NoError(t, err)
				err = os.WriteFile(fullPath, []byte(content), 0644)
				require.NoError(t, err)
			}

			// Run detector
			stack, err := DetectStack(tmpDir)
			require.NoError(t, err)

			// Verify results
			assert.Equal(t, tt.expected.Language, stack.Language)
			assert.Equal(t, tt.expected.Framework, stack.Framework)
			assert.Equal(t, tt.expected.Database, stack.Database)
		})
	}
}

func TestDetectStackErrors(t *testing.T) {
	tests := []struct {
		name          string
		setup         func(t *testing.T) string
		expectedError string
	}{
		{
			name: "Invalid package.json",
			setup: func(t *testing.T) string {
				tmpDir, err := os.MkdirTemp("", "test-*")
				require.NoError(t, err)
				err = os.WriteFile(filepath.Join(tmpDir, "package.json"), []byte("invalid json"), 0644)
				require.NoError(t, err)
				return tmpDir
			},
			expectedError: "invalid character 'i' looking for beginning of value",
		},
		{
			name: "Unreadable directory",
			setup: func(t *testing.T) string {
				tmpDir, err := os.MkdirTemp("", "test-*")
				require.NoError(t, err)
				err = os.Chmod(tmpDir, 0000)
				require.NoError(t, err)
				return tmpDir
			},
			expectedError: "permission denied",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := tt.setup(t)
			defer os.RemoveAll(tmpDir)

			_, err := DetectStack(tmpDir)
			if tt.expectedError != "" {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
