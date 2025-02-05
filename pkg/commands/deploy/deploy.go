// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package deploy

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/sysinfo"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// findDeploymentFile looks for a deployment file in the current directory
func findDeploymentFile() (string, error) {
	// List of possible deployment file names
	possibleFiles := []string{
		"deployment.yaml",
		"deployment.yml",
		"nexlayer.yaml",
		"nexlayer.yml",
	}

	for _, file := range possibleFiles {
		if _, err := os.Stat(file); err == nil {
			return file, nil
		}
	}

	return "", fmt.Errorf("no deployment file found in current directory. Expected one of: %v", possibleFiles)
}

// NewCommand creates a new deploy command
func NewCommand(client api.APIClient) *cobra.Command {
	var yamlFile string

	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy an application",
		Long: `Deploy an application using a YAML configuration file.

If no file is specified with -f flag, it will look for one of these files in the current directory:
- deployment.yaml
- deployment.yml
- nexlayer.yaml
- nexlayer.yml

Example:
  nexlayer deploy
  nexlayer deploy -f custom-deploy.yaml`,
		Args: cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			// If no file specified, try to find one
			if yamlFile == "" {
				file, err := findDeploymentFile()
				if err != nil {
					return err
				}
				yamlFile = file
				cmd.Printf("Using deployment file: %s\n", yamlFile)
			}

			return runDeploy(cmd, yamlFile)
		},
	}

	cmd.Flags().StringVarP(&yamlFile, "file", "f", "", "Path to YAML/JSON configuration file (optional)")

	return cmd
}

func runDeploy(cmd *cobra.Command, yamlFile string) error {
	cmd.Println(ui.RenderTitleWithBorder("Deploying Application"))

	// Create API client for sending feedback
	apiClient := api.NewClient("")

	// Read the configuration file
	fileContent, err := ioutil.ReadFile(yamlFile)
	if err != nil {
		return fmt.Errorf("failed to read configuration file: %w", err)
	}

	// Parse YAML to get deployment name
	var config struct {
		Application struct {
			Template struct {
				Name           string `yaml:"name"`
				DeploymentName string `yaml:"deploymentName"`
			} `yaml:"template"`
		} `yaml:"application"`
	}

	if err := yaml.Unmarshal(fileContent, &config); err != nil {
		cmd.Printf("Warning: Could not parse deployment name from config: %v\n", err)
	}

	// Create HTTP request
	url := "https://app.staging.nexlayer.io/startUserDeployment"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(fileContent))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set content type based on file extension
	contentType := "application/json"
	if filepath.Ext(yamlFile) == ".yaml" || filepath.Ext(yamlFile) == ".yml" {
		contentType = "text/x-yaml"
	}
	req.Header.Set("Content-Type", contentType)

	// Create HTTP client that skips SSL verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send deployment request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("deployment failed: %s", string(body))
	}

	cmd.Printf("Deployment started successfully!\n")
	cmd.Printf("Response: %s\n", string(body))

	// Send automatic feedback
	go func() {
		// Get system information
		sysInfo := sysinfo.GetSystemInfo()

		// Extract deployment URL from response
		var deployResp struct {
			URL string `json:"url"`
		}
		if err := json.Unmarshal(body, &deployResp); err == nil {
			sysInfo.DeploymentURL = deployResp.URL
		}

		// Format feedback message
		deploymentName := config.Application.Template.DeploymentName
		if deploymentName == "" {
			deploymentName = filepath.Base(yamlFile)
		}
		feedbackMsg := sysInfo.FormatFeedback(deploymentName)

		// Send feedback
		if err := apiClient.SendFeedback(context.Background(), feedbackMsg); err != nil {
			cmd.Printf("Warning: Failed to send automatic feedback: %v\n", err)
		}
	}()

	return nil
}
