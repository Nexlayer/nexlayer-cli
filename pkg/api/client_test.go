package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/Nexlayer/nexlayer-cli/pkg/api/types"
)

func TestClient_CreateApplication(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/apps", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Send response
		app := types.App{
			ID:        "test-app",
			Name:      "Test App",
			Status:    "created",
			CreatedAt: time.Now(),
		}
		json.NewEncoder(w).Encode(app)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL)
	client.SetToken("test-token")

	// Test
	req := &types.CreateAppRequest{Name: "Test App"}
	app, err := client.CreateApplication(context.Background(), req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, app)
	assert.Equal(t, "test-app", app.ID)
	assert.Equal(t, "Test App", app.Name)
	assert.Equal(t, "created", app.Status)
}

func TestClient_ListApplications(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/apps", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Send response
		apps := []types.App{
			{
				ID:        "app-1",
				Name:      "App 1",
				Status:    "running",
				CreatedAt: time.Now(),
			},
			{
				ID:        "app-2",
				Name:      "App 2",
				Status:    "stopped",
				CreatedAt: time.Now(),
			},
		}
		json.NewEncoder(w).Encode(apps)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL)
	client.SetToken("test-token")

	// Test
	apps, err := client.ListApplications(context.Background())

	// Assert
	assert.NoError(t, err)
	assert.Len(t, apps, 2)
	assert.Equal(t, "app-1", apps[0].ID)
	assert.Equal(t, "app-2", apps[1].ID)
}

func TestClient_StartUserDeployment(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/apps/test-app/deployments", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Send response
		deployment := types.Deployment{
			ID:            "deploy-1",
			ApplicationID: "test-app",
			Status:        "pending",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		json.NewEncoder(w).Encode(deployment)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL)
	client.SetToken("test-token")

	// Test
	req := &types.DeployRequest{YAML: "test: yaml"}
	deployment, err := client.StartUserDeployment(context.Background(), "test-app", req)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, deployment)
	assert.Equal(t, "deploy-1", deployment.ID)
	assert.Equal(t, "test-app", deployment.ApplicationID)
	assert.Equal(t, "pending", deployment.Status)
}

func TestClient_SaveCustomDomain(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/apps/test-app/domains", r.URL.Path)
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Send response
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL)
	client.SetToken("test-token")

	// Test
	err := client.SaveCustomDomain(context.Background(), "test-app", "test.example.com")

	// Assert
	assert.NoError(t, err)
}

func TestClient_GetDeployments(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/apps/test-app/deployments", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Send response
		deployments := []types.Deployment{
			{
				ID:            "deploy-1",
				ApplicationID: "test-app",
				Status:        "running",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			{
				ID:            "deploy-2",
				ApplicationID: "test-app",
				Status:        "failed",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}
		json.NewEncoder(w).Encode(deployments)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL)
	client.SetToken("test-token")

	// Test
	deployments, err := client.GetDeployments(context.Background(), "test-app")

	// Assert
	assert.NoError(t, err)
	assert.Len(t, deployments, 2)
	assert.Equal(t, "deploy-1", deployments[0].ID)
	assert.Equal(t, "deploy-2", deployments[1].ID)
}

func TestClient_GetDeploymentInfo(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/deployments/test-ns/test-app", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Send response
		info := types.DeploymentInfo{
			ID:            "deploy-1",
			ApplicationID: "test-app",
			Status:        "running",
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
			Namespace:     "test-ns",
			Config:        "test: yaml",
		}
		json.NewEncoder(w).Encode(info)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL)
	client.SetToken("test-token")

	// Test
	info, err := client.GetDeploymentInfo(context.Background(), "test-ns", "test-app")

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, info)
	assert.Equal(t, "deploy-1", info.ID)
	assert.Equal(t, "test-app", info.ApplicationID)
	assert.Equal(t, "running", info.Status)
	assert.Equal(t, "test-ns", info.Namespace)
}

func TestClient_GetAppInfo(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/apps/test-app", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Return response
		app := types.App{
			ID:           "test-app",
			Name:         "Test App",
			CreatedAt:    time.Now(),
			LastDeployAt: time.Now(),
			Status:       "active",
		}
		json.NewEncoder(w).Encode(app)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL)
	client.SetToken("test-token")

	// Test GetAppInfo
	app, err := client.GetAppInfo(context.Background(), "test-app")
	assert.NoError(t, err)
	assert.Equal(t, "test-app", app.ID)
	assert.Equal(t, "Test App", app.Name)
	assert.Equal(t, "active", app.Status)
}

func TestClient_GetDomains(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/apps/test-app/domains", r.URL.Path)
		assert.Equal(t, "GET", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Return response
		domains := []types.Domain{
			{Domain: "example.com", Status: "active"},
			{Domain: "test.com", Status: "pending"},
		}
		json.NewEncoder(w).Encode(domains)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL)
	client.SetToken("test-token")

	// Test GetDomains
	domains, err := client.GetDomains("test-app")
	assert.NoError(t, err)
	assert.Len(t, domains, 2)
	assert.Equal(t, "example.com", domains[0].Domain)
	assert.Equal(t, "active", domains[0].Status)
}

func TestClient_RemoveDomain(t *testing.T) {
	// Setup test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request
		assert.Equal(t, "/apps/test-app/domains/example.com", r.URL.Path)
		assert.Equal(t, "DELETE", r.Method)
		assert.Equal(t, "Bearer test-token", r.Header.Get("Authorization"))

		// Return success
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	// Create client
	client := NewClient(server.URL)
	client.SetToken("test-token")

	// Test RemoveDomain
	err := client.RemoveDomain("test-app", "example.com")
	assert.NoError(t, err)
}

func TestClient_ErrorCases(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
		response   string
		wantErr    string
	}{
		{
			name:       "Unauthorized",
			statusCode: http.StatusUnauthorized,
			response:   `{"message": "unauthorized", "code": "auth_error"}`,
			wantErr:    "failed to get app info: request failed: unauthorized",
		},
		{
			name:       "Not Found",
			statusCode: http.StatusNotFound,
			response:   `{"message": "application not found", "code": "not_found"}`,
			wantErr:    "failed to get app info: request failed: application not found",
		},
		{
			name:       "Internal Server Error",
			statusCode: http.StatusInternalServerError,
			response:   `{"message": "internal server error", "code": "server_error"}`,
			wantErr:    "failed to get app info: request failed: internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.response))
			}))
			defer server.Close()

			// Create client
			client := NewClient(server.URL)
			client.SetToken("test-token")

			// Test GetAppInfo
			_, err := client.GetAppInfo(context.Background(), "test-app")
			assert.Error(t, err)
			assert.Equal(t, tt.wantErr, err.Error())
		})
	}
}

func TestClient_InvalidURL(t *testing.T) {
	client := NewClient("invalid-url")
	client.SetToken("test-token")

	_, err := client.GetAppInfo(context.Background(), "test-app")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unsupported protocol scheme")
}
