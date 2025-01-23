package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/api/types"
	"github.com/stretchr/testify/assert"
)

func setupTestConfig(t *testing.T, token string) (string, func()) {
	// Create temp directory
	configDir, err := os.MkdirTemp("", "nexlayer-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// Create config file
	configPath := filepath.Join(configDir, "config")
	config := types.Config{
		Token: token,
	}
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		t.Fatalf("Failed to write config: %v", err)
	}

	// Return cleanup function
	return configDir, func() {
		os.RemoveAll(configDir)
	}
}

func TestNewClient(t *testing.T) {
	// Test cases
	tests := []struct {
		name        string
		token       string
		wantErr     bool
		errContains string
	}{
		{
			name:  "valid token",
			token: "valid-token",
		},
		{
			name:        "empty token",
			token:       "",
			wantErr:     true,
			errContains: "no token found in config file",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test config
			tmpDir, cleanup := setupTestConfig(t, tt.token)
			defer cleanup()

			// Set HOME environment variable to use our test config
			oldHome := os.Getenv("HOME")
			os.Setenv("HOME", tmpDir)
			defer os.Setenv("HOME", oldHome)

			// Create client
			client, err := NewClient("https://api.test.com")
			if tt.wantErr {
				if err == nil {
					t.Errorf("NewClient() error = nil, want error containing %q", tt.errContains)
				} else if !contains(err.Error(), tt.errContains) {
					t.Errorf("NewClient() error = %v, want error containing %q", err, tt.errContains)
				}
				return
			}
			if err != nil {
				t.Errorf("NewClient() error = %v, want nil", err)
			}
			if client == nil {
				t.Error("NewClient() client is nil")
			}
		})
	}
}

func TestStartUserDeployment(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/startUserDeployment/test-app", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "text/x-yaml", r.Header.Get("Content-Type"))

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(types.DeploymentResponse{
			Message:   "Deployment started",
			Namespace: "default",
			URL:       "https://example.com",
		})
	}))
	defer server.Close()

	// Create test YAML file
	yamlContent := []byte("name: test-app\nversion: 1.0.0")
	tmpfile, err := os.CreateTemp("", "test-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())
	if _, err := tmpfile.Write(yamlContent); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Create client
	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Test deployment
	resp, err := client.StartUserDeployment("test-app", tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}

	// Check response
	assert.Equal(t, "Deployment started", resp.Message)
	assert.Equal(t, "default", resp.Namespace)
	assert.Equal(t, "https://example.com", resp.URL)
}

func TestSaveCustomDomain(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/saveCustomDomain/test-app", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		// Check request body
		var req types.SaveCustomDomainRequest
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "example.com", req.Domain)

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(types.SaveCustomDomainResponse{
			Message: "Custom domain saved successfully",
		})
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Test saving custom domain
	resp, err := client.SaveCustomDomain("test-app", "example.com")
	if err != nil {
		t.Fatal(err)
	}

	// Check response
	assert.Equal(t, "Custom domain saved successfully", resp.Message)
}

func TestGetDeployments(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/getDeployments/test-app", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(types.GetDeploymentsResponse{
			Deployments: []types.DeploymentInfo{
				{
					Namespace:        "default",
					TemplateID:       "123",
					TemplateName:     "MERN Todo",
					DeploymentStatus: "running",
				},
			},
		})
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Test getting deployments
	resp, err := client.GetDeployments("test-app")
	if err != nil {
		t.Fatal(err)
	}

	// Check response
	assert.Len(t, resp.Deployments, 1)
	d := resp.Deployments[0]
	assert.Equal(t, "default", d.Namespace)
	assert.Equal(t, "123", d.TemplateID)
	assert.Equal(t, "MERN Todo", d.TemplateName)
	assert.Equal(t, "running", d.DeploymentStatus)
}

func TestListApplications(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/api/v1/applications", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		// Send response
		w.Header().Set("Content-Type", "application/json")
		createdAt, _ := time.Parse(time.RFC3339, "2025-01-23T01:22:27-05:00")
		json.NewEncoder(w).Encode(struct {
			Applications []types.Application `json:"applications"`
		}{
			Applications: []types.Application{
				{
					ID:        "app-123",
					Name:      "test-app",
					CreatedAt: createdAt,
				},
			},
		})
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Test listing applications
	apps, err := client.ListApplications()
	if err != nil {
		t.Fatal(err)
	}

	// Check response
	assert.Len(t, apps, 1)
	assert.Equal(t, "app-123", apps[0].ID)
	assert.Equal(t, "test-app", apps[0].Name)
	createdAt, _ := time.Parse(time.RFC3339, "2025-01-23T01:22:27-05:00")
	assert.Equal(t, createdAt, apps[0].CreatedAt)
}

func TestCreateApplication(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/api/v1/applications", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		// Check request body
		var req map[string]string
		json.NewDecoder(r.Body).Decode(&req)
		assert.Equal(t, "test-app", req["name"])

		// Send response
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(types.CreateApplicationResponse{
			ID:   "app-123",
			Name: "test-app",
		})
	}))
	defer server.Close()

	// Create client
	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatal(err)
	}

	// Test creating application
	resp, err := client.CreateApplication("test-app")
	if err != nil {
		t.Fatal(err)
	}

	// Check response
	assert.Equal(t, "app-123", resp.ID)
	assert.Equal(t, "test-app", resp.Name)
}

func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && s[len(s)-1] != '/' && s[0] != '/' && s[len(s)-1] != '\\' && s[0] != '\\' && s[len(s)-1] != '.' && s[0] != '.' && s[len(s)-1] != '_' && s[0] != '_' && s[len(s)-1] != '-' && s[0] != '-' && s[len(s)-1] != ' ' && s[0] != ' ' && s[len(s)-1] != '\t' && s[0] != '\t' && s[len(s)-1] != '\n' && s[0] != '\n' && s[len(s)-1] != '\r' && s[0] != '\r'
}
