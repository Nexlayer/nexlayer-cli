package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func setupTestConfig(t *testing.T, token string) (string, func()) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "nexlayer-test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create config directory
	configDir := filepath.Join(tmpDir, ".nexlayer")
	if err := os.MkdirAll(configDir, 0700); err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	// Create config file
	configPath := filepath.Join(configDir, "config")
	config := Config{Token: token}
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	// Return cleanup function
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}

	return tmpDir, cleanup
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
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check path
		if r.URL.Path != "/startUserDeployment/test-app" {
			t.Errorf("Expected path /startUserDeployment/test-app, got %s", r.URL.Path)
		}

		// Check auth header
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization: Bearer test-token, got %s", r.Header.Get("Authorization"))
		}

		// Return response
		resp := DeploymentResponse{
			Message:   "Deployment started",
			URL:      "https://test-app.nexlayer.io",
			Namespace: "test-namespace",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Setup test config
	tmpDir, cleanup := setupTestConfig(t, "test-token")
	defer cleanup()

	// Set HOME environment variable
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Create test YAML file
	yamlContent := []byte("name: test-app\nversion: 1.0.0")
	yamlPath := filepath.Join(tmpDir, "test.yaml")
	if err := os.WriteFile(yamlPath, yamlContent, 0600); err != nil {
		t.Fatalf("Failed to write test YAML: %v", err)
	}

	// Create client
	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test deployment
	resp, err := client.StartUserDeployment("test-app", yamlPath)
	if err != nil {
		t.Fatalf("StartUserDeployment() error = %v", err)
	}

	// Check response
	if resp.Message != "Deployment started" {
		t.Errorf("Expected message 'Deployment started', got %q", resp.Message)
	}
	if resp.URL != "https://test-app.nexlayer.io" {
		t.Errorf("Expected URL 'https://test-app.nexlayer.io', got %q", resp.URL)
	}
	if resp.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got %q", resp.Namespace)
	}
}

func TestSaveCustomDomain(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// Check path
		if r.URL.Path != "/saveCustomDomain/test-app" {
			t.Errorf("Expected path /saveCustomDomain/test-app, got %s", r.URL.Path)
		}

		// Check auth header
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization: Bearer test-token, got %s", r.Header.Get("Authorization"))
		}

		// Return response
		resp := CustomDomainResponse{
			Message: "Domain configured",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Setup test config
	tmpDir, cleanup := setupTestConfig(t, "test-token")
	defer cleanup()

	// Set HOME environment variable
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Create client
	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test domain configuration
	resp, err := client.SaveCustomDomain("test-app", "test.example.com")
	if err != nil {
		t.Fatalf("SaveCustomDomain() error = %v", err)
	}

	// Check response
	if resp.Message != "Domain configured" {
		t.Errorf("Expected message 'Domain configured', got %q", resp.Message)
	}
}

func TestGetDeployments(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check method
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// Check path
		if r.URL.Path != "/getDeployments/test-app" {
			t.Errorf("Expected path /getDeployments/test-app, got %s", r.URL.Path)
		}

		// Check auth header
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Authorization: Bearer test-token, got %s", r.Header.Get("Authorization"))
		}

		// Return response
		resp := DeploymentsResponse{
			Deployments: []DeploymentInfo{
				{
					Namespace:        "test-namespace",
					ApplicationID:    "test-app",
					DeploymentStatus: "running",
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Setup test config
	tmpDir, cleanup := setupTestConfig(t, "test-token")
	defer cleanup()

	// Set HOME environment variable
	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	defer os.Setenv("HOME", oldHome)

	// Create client
	client, err := NewClient(server.URL)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	// Test getting deployments
	resp, err := client.GetDeployments("test-app")
	if err != nil {
		t.Fatalf("GetDeployments() error = %v", err)
	}

	// Check response
	if len(resp.Deployments) != 1 {
		t.Errorf("Expected 1 deployment, got %d", len(resp.Deployments))
	}
	d := resp.Deployments[0]
	if d.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got %q", d.Namespace)
	}
	if d.ApplicationID != "test-app" {
		t.Errorf("Expected application ID 'test-app', got %q", d.ApplicationID)
	}
	if d.DeploymentStatus != "running" {
		t.Errorf("Expected status 'running', got %q", d.DeploymentStatus)
	}
}

func contains(s, substr string) bool {
	return s != "" && substr != "" && s != substr && s[len(s)-1] != '/' && s[0] != '/' && s[len(s)-1] != '\\' && s[0] != '\\' && s[len(s)-1] != '.' && s[0] != '.' && s[len(s)-1] != '_' && s[0] != '_' && s[len(s)-1] != '-' && s[0] != '-' && s[len(s)-1] != ' ' && s[0] != ' ' && s[len(s)-1] != '\t' && s[0] != '\t' && s[len(s)-1] != '\n' && s[0] != '\n' && s[len(s)-1] != '\r' && s[0] != '\r'
}
