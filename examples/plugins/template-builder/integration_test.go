package templatebuilder

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func buildCLI(t *testing.T, dir string) string {
	execPath := filepath.Join(dir, "nexlayer")
	buildCmd := exec.Command("go", "build", "-o", execPath, ".")
	buildCmd.Dir = dir
	buildCmd.Env = append(os.Environ(),
		"GOOS=darwin",
		"GOARCH=amd64",
		"CGO_ENABLED=0",
	)

	buildOutput, err := buildCmd.CombinedOutput()
	require.NoError(t, err, "failed to build CLI: %s", buildOutput)

	// Set executable permissions
	err = os.Chmod(execPath, 0755)
	require.NoError(t, err, "failed to set executable permissions")

	return execPath
}

func TestCLICommands(t *testing.T) {
	t.Skip("Skipping CLI integration tests")
}

func TestTemplateRegistry(t *testing.T) {
	t.Skip("Skipping registry integration tests")
}
