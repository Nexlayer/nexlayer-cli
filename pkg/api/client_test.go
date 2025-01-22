package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestStartDeployment(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "POST" {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/startUserDeployment/test-app" {
			t.Errorf("Expected path /startUserDeployment/test-app, got %s", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "text/x-yaml" {
			t.Errorf("Expected Content-Type text/x-yaml, got %s", r.Header.Get("Content-Type"))
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(StartDeploymentResponse{
			Message:   "Deployment started successfully",
			Namespace: "test-namespace",
			URL:      "https://test-url.com",
		})
	}))
	defer server.Close()

	// Create client
	client := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	// Test deployment
	resp, err := client.StartDeployment("test-app", []byte("test yaml"))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if resp.Message != "Deployment started successfully" {
		t.Errorf("Expected message 'Deployment started successfully', got '%s'", resp.Message)
	}
	if resp.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", resp.Namespace)
	}
	if resp.URL != "https://test-url.com" {
		t.Errorf("Expected URL 'https://test-url.com', got '%s'", resp.URL)
	}
}

func TestGetDeployments(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/getDeployments/test-app" {
			t.Errorf("Expected path /getDeployments/test-app, got %s", r.URL.Path)
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetDeploymentsResponse{
			Deployments: []DeploymentInfo{
				{
					Namespace:        "test-namespace",
					TemplateID:      "test-template",
					TemplateName:    "Test Template",
					DeploymentStatus: "running",
				},
			},
		})
	}))
	defer server.Close()

	// Create client
	client := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	// Test get deployments
	resp, err := client.GetDeployments("test-app")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if len(resp.Deployments) != 1 {
		t.Fatalf("Expected 1 deployment, got %d", len(resp.Deployments))
	}

	deployment := resp.Deployments[0]
	if deployment.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", deployment.Namespace)
	}
	if deployment.TemplateID != "test-template" {
		t.Errorf("Expected template ID 'test-template', got '%s'", deployment.TemplateID)
	}
}

func TestGetDeploymentInfo(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/getDeploymentInfo/test-namespace/test-app" {
			t.Errorf("Expected path /getDeploymentInfo/test-namespace/test-app, got %s", r.URL.Path)
		}

		// Return success response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(GetDeploymentInfoResponse{
			Deployment: DeploymentInfo{
				Namespace:        "test-namespace",
				TemplateID:      "test-template",
				TemplateName:    "Test Template",
				DeploymentStatus: "running",
			},
		})
	}))
	defer server.Close()

	// Create client
	client := &Client{
		httpClient: server.Client(),
		baseURL:    server.URL,
	}

	// Test get deployment info
	resp, err := client.GetDeploymentInfo("test-namespace", "test-app")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	deployment := resp.Deployment
	if deployment.Namespace != "test-namespace" {
		t.Errorf("Expected namespace 'test-namespace', got '%s'", deployment.Namespace)
	}
	if deployment.TemplateID != "test-template" {
		t.Errorf("Expected template ID 'test-template', got '%s'", deployment.TemplateID)
	}
}
