package watchcmd

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWatchCommand(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "watch-test")
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
        - 80
`), 0644)
	require.NoError(t, err)

	// Create command
	cmd := NewCommand()
	require.NotNil(t, cmd)

	// Verify command properties
	assert.Equal(t, "watch", cmd.Use)
	assert.NotEmpty(t, cmd.Short)
	assert.NotEmpty(t, cmd.Long)

	// Test flags
	flags := cmd.Flags()
	assert.True(t, flags.HasFlags())

	fileFlag := flags.Lookup("file")
	require.NotNil(t, fileFlag)
	assert.Equal(t, "nexlayer.yaml", fileFlag.DefValue)

	noSyncFlag := flags.Lookup("no-sync")
	require.NotNil(t, noSyncFlag)
	assert.Equal(t, "false", noSyncFlag.DefValue)

	noDeployFlag := flags.Lookup("no-deploy")
	require.NotNil(t, noDeployFlag)
	assert.Equal(t, "false", noDeployFlag.DefValue)

	// Test file watching (in a goroutine to avoid blocking)
	done := make(chan bool)
	go func() {
		// Modify the file after a short delay
		time.Sleep(100 * time.Millisecond)
		err := os.WriteFile(yamlPath, []byte(`application:
  name: test-app-modified
  pods:
    - name: app
      image: nginx:latest
      servicePorts:
        - 80
`), 0644)
		assert.NoError(t, err)
		done <- true
	}()

	// Set flags for testing
	flags.Set("file", yamlPath)
	flags.Set("no-sync", "true")  // Disable sync for testing
	flags.Set("no-deploy", "true") // Disable deploy for testing

	// Run command with timeout
	errChan := make(chan error)
	go func() {
		errChan <- cmd.RunE(cmd, []string{})
	}()

	// Wait for either completion or timeout
	select {
	case err := <-errChan:
		assert.NoError(t, err)
	case <-done:
		// Test passed
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
}
