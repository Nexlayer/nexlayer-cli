package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/types"
)

func TestStartDeployment(t *testing.T) {
	// Create a temporary YAML file
	tmpDir := t.TempDir()
	yamlPath := filepath.Join(tmpDir, "test.yaml")
	err := os.WriteFile(yamlPath, []byte("template: python"), 0644)
	require.NoError(t, err)

	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/startUserDeployment/test-app", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		resp := types.StartDeploymentResponse{
			Namespace: "test-ns",
			URL:       "https://test-ns.nexlayer.com",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	// Create client and make request
	client := NewClient(server.URL)
	resp, err := client.StartDeployment(context.Background(), "test-app", yamlPath)
	require.NoError(t, err)
	assert.Equal(t, "test-ns", resp.Namespace)
	assert.Equal(t, "https://test-ns.nexlayer.com", resp.URL)
}

func TestSaveCustomDomain(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/saveCustomDomain/test-app", r.URL.Path)
		assert.Equal(t, "POST", r.Method)

		resp := types.SaveCustomDomainResponse{
			Success: true,
			Message: "Domain saved successfully",
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	err := client.SaveCustomDomain(context.Background(), "test-app", "example.com")
	require.NoError(t, err)
}

func TestGetDeployments(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/getDeployments/test-app", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		deployments := []types.Deployment{
			{
				Namespace:        "test-ns-1",
				TemplateName:     "python",
				TemplateID:       "123",
				DeploymentStatus: "running",
			},
			{
				Namespace:        "test-ns-2",
				TemplateName:     "node",
				TemplateID:       "456",
				DeploymentStatus: "pending",
			},
		}
		json.NewEncoder(w).Encode(deployments)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	deployments, err := client.GetDeployments(context.Background(), "test-app")
	require.NoError(t, err)
	assert.Len(t, deployments, 2)
	assert.Equal(t, "test-ns-1", deployments[0].Namespace)
	assert.Equal(t, "test-ns-2", deployments[1].Namespace)
}

func TestGetDeploymentInfo(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/getDeploymentInfo/test-ns/test-app", r.URL.Path)
		assert.Equal(t, "GET", r.Method)

		deployment := types.DeploymentInfo{
			Namespace:        "test-ns",
			TemplateName:     "python",
			TemplateID:       "123",
			DeploymentStatus: "running",
		}
		json.NewEncoder(w).Encode(deployment)
	}))
	defer server.Close()

	client := NewClient(server.URL)
	info, err := client.GetDeploymentInfo(context.Background(), "test-ns", "test-app")
	require.NoError(t, err)
	assert.Equal(t, "test-ns", info.Namespace)
	assert.Equal(t, "python", info.TemplateName)
	assert.Equal(t, "123", info.TemplateID)
	assert.Equal(t, "running", info.DeploymentStatus)
}

func TestClientErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Internal server error"}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)

	t.Run("StartDeployment error", func(t *testing.T) {
		tmpDir := t.TempDir()
		yamlPath := filepath.Join(tmpDir, "test.yaml")
		err := os.WriteFile(yamlPath, []byte("template: python"), 0644)
		require.NoError(t, err)

		_, err = client.StartDeployment(context.Background(), "test-app", yamlPath)
		assert.Error(t, err)
	})

	t.Run("SaveCustomDomain error", func(t *testing.T) {
		err := client.SaveCustomDomain(context.Background(), "test-app", "example.com")
		assert.Error(t, err)
	})

	t.Run("GetDeployments error", func(t *testing.T) {
		_, err := client.GetDeployments(context.Background(), "test-app")
		assert.Error(t, err)
	})

	t.Run("GetDeploymentInfo error", func(t *testing.T) {
		_, err := client.GetDeploymentInfo(context.Background(), "test-ns", "test-app")
		assert.Error(t, err)
	})
}
