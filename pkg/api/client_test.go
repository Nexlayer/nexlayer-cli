package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestDeploy(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and headers
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected application/json content type, got %s", r.Header.Get("Content-Type"))
		}
		if r.URL.Path != "/api/v1/deploy" {
			t.Errorf("Expected /api/v1/deploy path, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Bearer test-token auth header, got %s", r.Header.Get("Authorization"))
		}

		// Verify request body
		var config ServiceConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if config.AppName != "test-app" {
			t.Errorf("Expected app name test-app, got %s", config.AppName)
		}
		if config.ServiceName != "test-service" {
			t.Errorf("Expected service name test-service, got %s", config.ServiceName)
		}

		// Send response
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	// Create client and test deploy
	client := NewClient(server.URL)
	err := client.Deploy("test-token", "test-app", "test-service", []string{"DB_URL=test"})
	if err != nil {
		t.Fatalf("Deploy failed: %v", err)
	}
}

func TestConfigure(t *testing.T) {
	// Set auth token for test
	os.Setenv("NEXLAYER_AUTH_TOKEN", "test-token")
	defer os.Unsetenv("NEXLAYER_AUTH_TOKEN")

	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and headers
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected application/json content type, got %s", r.Header.Get("Content-Type"))
		}
		if r.URL.Path != "/api/v1/services/configure" {
			t.Errorf("Expected /api/v1/services/configure path, got %s", r.URL.Path)
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("Expected Bearer test-token auth header, got %s", r.Header.Get("Authorization"))
		}

		// Verify request body
		var config ServiceConfig
		if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
			t.Errorf("Failed to decode request body: %v", err)
		}
		if config.AppName != "test-app" {
			t.Errorf("Expected app name test-app, got %s", config.AppName)
		}
		if config.ServiceName != "test-service" {
			t.Errorf("Expected service name test-service, got %s", config.ServiceName)
		}

		// Send response
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	// Create client and test configure
	client := NewClient(server.URL)
	err := client.Configure("test-app", "test-service", []string{"DB_URL=test"})
	if err != nil {
		t.Fatalf("Configure failed: %v", err)
	}
}

func TestStartUserDeployment(t *testing.T) {
	// Mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "text/x-yaml" {
			t.Errorf("Expected text/x-yaml content type, got %s", r.Header.Get("Content-Type"))
		}
		if r.URL.Path != "/startUserDeployment/test-app" {
			t.Errorf("Expected /startUserDeployment/test-app path, got %s", r.URL.Path)
		}

		// Send response
		resp := StartUserDeploymentResponse{
			Message:   "Deployment started successfully",
			Namespace: "fantastic-fox",
			URL:      "https://fantastic-fox-my-mern-app.alpha.nexlayer.ai",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL)

	// Test deployment
	yamlContent := []byte("name: test-app\nversion: 1.0.0")
	resp, err := client.StartUserDeployment("test-app", yamlContent)
	if err != nil {
		t.Fatalf("StartUserDeployment failed: %v", err)
	}

	// Verify response
	if resp.Message != "Deployment started successfully" {
		t.Errorf("Expected 'Deployment started successfully', got %s", resp.Message)
	}
	if resp.Namespace != "fantastic-fox" {
		t.Errorf("Expected 'fantastic-fox', got %s", resp.Namespace)
	}
}

func TestSaveCustomDomain(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected application/json content type, got %s", r.Header.Get("Content-Type"))
		}
		if r.URL.Path != "/saveCustomDomain/test-app" {
			t.Errorf("Expected /saveCustomDomain/test-app path, got %s", r.URL.Path)
		}

		resp := SaveCustomDomainResponse{
			Message: "Custom domain saved successfully",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.SaveCustomDomain("test-app", "example.com")
	if err != nil {
		t.Fatalf("SaveCustomDomain failed: %v", err)
	}

	if resp.Message != "Custom domain saved successfully" {
		t.Errorf("Expected 'Custom domain saved successfully', got %s", resp.Message)
	}
}

func TestGetDeployments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/getDeployments/test-app" {
			t.Errorf("Expected /getDeployments/test-app path, got %s", r.URL.Path)
		}

		resp := GetDeploymentsResponse{
			Deployments: []struct {
				Namespace        string `json:"namespace"`
				TemplateID       string `json:"templateID"`
				TemplateName     string `json:"templateName"`
				DeploymentStatus string `json:"deploymentStatus"`
			}{
				{
					Namespace:        "ecstatic-frog",
					TemplateID:       "0001",
					TemplateName:     "MERN Stack",
					DeploymentStatus: "Running",
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.GetDeployments("test-app")
	if err != nil {
		t.Fatalf("GetDeployments failed: %v", err)
	}

	if len(resp.Deployments) != 1 {
		t.Errorf("Expected 1 deployment, got %d", len(resp.Deployments))
	}
	if resp.Deployments[0].Namespace != "ecstatic-frog" {
		t.Errorf("Expected namespace ecstatic-frog, got %s", resp.Deployments[0].Namespace)
	}
}

func TestGetDeploymentInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/getDeploymentInfo/ecstatic-frog/test-app" {
			t.Errorf("Expected /getDeploymentInfo/ecstatic-frog/test-app path, got %s", r.URL.Path)
		}

		resp := GetDeploymentInfoResponse{
			Deployment: struct {
				Namespace        string `json:"namespace"`
				TemplateID       string `json:"templateID"`
				TemplateName     string `json:"templateName"`
				DeploymentStatus string `json:"deploymentStatus"`
			}{
				Namespace:        "ecstatic-frog",
				TemplateID:       "0001",
				TemplateName:     "K-d chat",
				DeploymentStatus: "running",
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	resp, err := client.GetDeploymentInfo("ecstatic-frog", "test-app")
	if err != nil {
		t.Fatalf("GetDeploymentInfo failed: %v", err)
	}

	deployment := resp.Deployment
	if deployment.Namespace != "ecstatic-frog" {
		t.Errorf("Expected namespace 'ecstatic-frog', got %s", deployment.Namespace)
	}
	if deployment.TemplateID != "0001" {
		t.Errorf("Expected templateID '0001', got %s", deployment.TemplateID)
	}
}
