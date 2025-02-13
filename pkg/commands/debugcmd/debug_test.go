package debugcmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDebugCommand(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "debug-test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	// Create test nexlayer.yaml
	yamlPath := filepath.Join(tmpDir, "nexlayer.yaml")
	err = os.WriteFile(yamlPath, []byte(`application:
  name: test-app
  pods:
    - name: app
      image: nginx:latest
      servicePorts:
        - containerPort: 80
          servicePort: 80
          name: web
`), 0644)
	require.NoError(t, err)

	// Create command
	cmd := NewCommand()
	require.NotNil(t, cmd)

	// Verify command properties
	assert.Equal(t, "debug", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)

	// Test flags
	flags := cmd.Flags()
	assert.True(t, flags.HasFlags())

	fileFlag := flags.Lookup("file")
	require.NotNil(t, fileFlag)
	assert.Equal(t, "nexlayer.yaml", fileFlag.DefValue)

	fullFlag := flags.Lookup("full")
	require.NotNil(t, fullFlag)
	assert.Equal(t, "false", fullFlag.DefValue)

	jsonFlag := flags.Lookup("json")
	require.NotNil(t, jsonFlag)
	assert.Equal(t, "false", jsonFlag.DefValue)

	// Test configuration check
	result := checkConfiguration(yamlPath)
	assert.Equal(t, "success", result.Status)
	assert.Contains(t, result.Message, "valid")

	// Test registry access check
	result = checkRegistryAccess("nginx:latest")
	assert.Contains(t, []string{"success", "warning", "error"}, result.Status)

	// Test invalid configuration
	invalidPath := filepath.Join(tmpDir, "invalid.yaml")
	result = checkConfiguration(invalidPath)
	assert.Equal(t, "error", result.Status)
	assert.NotEmpty(t, result.Fixes)

	// Test deployment status check
	result = checkDeploymentStatus("test-app")
	assert.Equal(t, "success", result.Status)
	assert.Contains(t, result.Message, "test-app")

	// Test log analysis
	result = analyzeLogs()
	assert.Contains(t, []string{"success", "warning"}, result.Status)
}
