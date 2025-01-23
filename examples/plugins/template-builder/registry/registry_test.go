package registry

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	
	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTemplateRegistry(t *testing.T) {
	// Create a mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/templates":
			templates := []types.NexlayerTemplate{
				{
					Name:    "test-template",
					Version: "1.0.0",
				},
			}
			json.NewEncoder(w).Encode(templates)
		case "/templates/test-template":
			template := types.NexlayerTemplate{
				Name:    "test-template",
				Version: "1.0.0",
			}
			json.NewEncoder(w).Encode(template)
		case "/templates/test-template/versions":
			versions := []string{"1.0.0", "1.1.0", "2.0.0"}
			json.NewEncoder(w).Encode(versions)
		default:
			http.NotFound(w, r)
		}
	}))
	defer server.Close()
	
	// Create a registry client
	client := NewClient(server.URL)
	
	t.Run("List Templates", func(t *testing.T) {
		templates, err := client.ListTemplates(context.Background())
		require.NoError(t, err)
		assert.Len(t, templates, 1)
		assert.Equal(t, "test-template", templates[0].Name)
	})
	
	t.Run("Get Template", func(t *testing.T) {
		template, err := client.GetTemplate(context.Background(), "test-template")
		require.NoError(t, err)
		assert.Equal(t, "test-template", template.Name)
		assert.Equal(t, "1.0.0", template.Version)
	})
	
	t.Run("Get Template Versions", func(t *testing.T) {
		versions, err := client.GetTemplateVersions(context.Background(), "test-template")
		require.NoError(t, err)
		assert.Len(t, versions, 3)
		assert.Contains(t, versions, "1.0.0")
		assert.Contains(t, versions, "1.1.0")
		assert.Contains(t, versions, "2.0.0")
	})
	
	t.Run("Template Not Found", func(t *testing.T) {
		_, err := client.GetTemplate(context.Background(), "nonexistent")
		assert.Error(t, err)
	})
}

func TestTemplateRegistryErrors(t *testing.T) {
	// Create a server that always returns errors
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()
	
	client := NewClient(server.URL)
	
	t.Run("List Templates Error", func(t *testing.T) {
		_, err := client.ListTemplates(context.Background())
		assert.Error(t, err)
	})
	
	t.Run("Get Template Error", func(t *testing.T) {
		_, err := client.GetTemplate(context.Background(), "test")
		assert.Error(t, err)
	})
	
	t.Run("Get Versions Error", func(t *testing.T) {
		_, err := client.GetTemplateVersions(context.Background(), "test")
		assert.Error(t, err)
	})
}

func TestTemplateRegistryValidation(t *testing.T) {
	client := NewClient("http://localhost")
	
	t.Run("Invalid Template Name", func(t *testing.T) {
		_, err := client.GetTemplate(context.Background(), "")
		assert.Error(t, err)
	})
	
	t.Run("Invalid Context", func(t *testing.T) {
		_, err := client.ListTemplates(nil)
		assert.Error(t, err)
	})
}

func TestTemplateRegistryTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate a timeout by not responding
		select {}
	}))
	defer server.Close()
	
	client := NewClient(server.URL)
	ctx, cancel := context.WithTimeout(context.Background(), 100)
	defer cancel()
	
	_, err := client.ListTemplates(ctx)
	assert.Error(t, err)
}
