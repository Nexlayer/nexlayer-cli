package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

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
					TemplateName:     "K-d chat",
					DeploymentStatus: "running",
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
		t.Fatalf("Expected 1 deployment, got %d", len(resp.Deployments))
	}

	deployment := resp.Deployments[0]
	if deployment.Namespace != "ecstatic-frog" {
		t.Errorf("Expected namespace 'ecstatic-frog', got %s", deployment.Namespace)
	}
	if deployment.TemplateID != "0001" {
		t.Errorf("Expected templateID '0001', got %s", deployment.TemplateID)
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
