package detect

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/types"
)

// AnalyzeDirectory analyzes a directory to detect application type and configuration
func AnalyzeDirectory(dir string) (*types.AppConfig, error) {
	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", dir)
	}

	// Initialize default configuration
	config := &types.AppConfig{
		Name: filepath.Base(dir),
		Type: "generic",
		Container: &types.Container{
			Ports: []int{8080},
		},
		Resources: &types.Resources{
			CPU:    "100m",
			Memory: "128Mi",
		},
	}

	// Check for Dockerfile
	if _, err := os.Stat(filepath.Join(dir, "Dockerfile")); err == nil {
		config.Container.UseDockerfile = true
	}

	// Check for package.json (Node.js)
	if _, err := os.Stat(filepath.Join(dir, "package.json")); err == nil {
		config.Type = "nodejs"
		config.Container.Command = "npm start"
	}

	// Check for requirements.txt (Python)
	if _, err := os.Stat(filepath.Join(dir, "requirements.txt")); err == nil {
		config.Type = "python"
		config.Container.Command = "python app.py"
	}

	// Check for go.mod (Go)
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
		config.Type = "go"
		config.Container.Command = "./app"
	}

	// Check for .env file
	if envFile, err := os.Open(filepath.Join(dir, ".env")); err == nil {
		defer envFile.Close()
		if content, err := ioutil.ReadAll(envFile); err == nil {
			config.Env = []string{string(content)}
		}
	}

	return config, nil
}
