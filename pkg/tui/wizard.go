package tui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"gopkg.in/yaml.v2"
)

// DeploymentWizard handles the interactive deployment creation process
type DeploymentWizard struct {
	model Model
}

// NewDeploymentWizard creates a new deployment wizard
func NewDeploymentWizard() *DeploymentWizard {
	return &DeploymentWizard{
		model: NewModel(),
	}
}

// Run starts the deployment wizard
func (w *DeploymentWizard) Run() error {
	p := tea.NewProgram(w.model)
	model, err := p.Run()
	if err != nil {
		return fmt.Errorf("error running wizard: %w", err)
	}

	finalModel := model.(Model)
	if finalModel.err != nil {
		return finalModel.err
	}

	// Generate YAML configuration
	config := w.generateConfig(finalModel.config)

	// Save configuration
	if err := w.saveConfig(config); err != nil {
		return fmt.Errorf("error saving configuration: %w", err)
	}

	return nil
}

type AppConfig struct {
	Application struct {
		Template struct {
			Name           string `yaml:"name"`
			DeploymentName string `yaml:"deploymentName"`
			RegistryLogin  struct {
				Registry            string `yaml:"registry"`
				Username            string `yaml:"username"`
				PersonalAccessToken string `yaml:"personalAccessToken"`
			} `yaml:"registryLogin"`
			Pods []struct {
				Type       string `yaml:"type"`
				ExposeHttp bool   `yaml:"exposeHttp"`
				Name       string `yaml:"name"`
				Tag        string `yaml:"tag"`
				PrivateTag bool   `yaml:"privateTag"`
				Vars       []struct {
					Key   string `yaml:"key"`
					Value string `yaml:"value"`
				} `yaml:"vars"`
			} `yaml:"pods"`
		} `yaml:"template"`
	} `yaml:"application"`
}

func (w *DeploymentWizard) generateConfig(config DeploymentConfig) *AppConfig {
	appConfig := &AppConfig{}

	// Generate template name based on stack
	stack := strings.ToLower(fmt.Sprintf("%s-%s-%s",
		config.DatabaseType,
		config.BackendType,
		config.FrontendType,
	))

	appConfig.Application.Template.Name = stack
	appConfig.Application.Template.DeploymentName = config.DeploymentName

	// Registry login
	appConfig.Application.Template.RegistryLogin.Registry = "ghcr.io"
	appConfig.Application.Template.RegistryLogin.Username = config.GithubUsername
	appConfig.Application.Template.RegistryLogin.PersonalAccessToken = config.GithubToken

	// Generate image tags based on app name
	appName := strings.ToLower(config.AppName)
	username := strings.ToLower(config.GithubUsername)

	// Pods configuration
	appConfig.Application.Template.Pods = []struct {
		Type       string `yaml:"type"`
		ExposeHttp bool   `yaml:"exposeHttp"`
		Name       string `yaml:"name"`
		Tag        string `yaml:"tag"`
		PrivateTag bool   `yaml:"privateTag"`
		Vars       []struct {
			Key   string `yaml:"key"`
			Value string `yaml:"value"`
		} `yaml:"vars"`
	}{
		{
			Type:       "database",
			ExposeHttp: false,
			Name:       strings.ToLower(config.DatabaseType),
			Tag:        fmt.Sprintf("ghcr.io/%s/%s-%s:v0.0.1", username, appName, strings.ToLower(config.DatabaseType)),
			PrivateTag: true,
			Vars: []struct {
				Key   string `yaml:"key"`
				Value string `yaml:"value"`
			}{
				{Key: "DB_USERNAME", Value: "admin"},
				{Key: "DB_PASSWORD", Value: "passw0rd"},
				{Key: "DB_NAME", Value: strings.ToLower(appName)},
			},
		},
		{
			Type:       "backend",
			ExposeHttp: false,
			Name:       strings.ToLower(config.BackendType),
			Tag:        fmt.Sprintf("ghcr.io/%s/%s-%s:v0.0.1", username, appName, strings.ToLower(config.BackendType)),
			PrivateTag: true,
			Vars: []struct {
				Key   string `yaml:"key"`
				Value string `yaml:"value"`
			}{
				{Key: "DATABASE_URL", Value: "DATABASE_CONNECTION_STRING"},
			},
		},
		{
			Type:       "frontend",
			ExposeHttp: true,
			Name:       strings.ToLower(config.FrontendType),
			Tag:        fmt.Sprintf("ghcr.io/%s/%s-%s:v0.0.1", username, appName, strings.ToLower(config.FrontendType)),
			PrivateTag: true,
			Vars: []struct {
				Key   string `yaml:"key"`
				Value string `yaml:"value"`
			}{
				{Key: "BACKEND_URL", Value: "BACKEND_CONNECTION_URL"},
			},
		},
	}

	return appConfig
}

func (w *DeploymentWizard) saveConfig(config *AppConfig) error {
	// Create deployment directory if it doesn't exist
	deployDir := "deployment"
	if err := os.MkdirAll(deployDir, 0755); err != nil {
		return err
	}

	// Marshal configuration to YAML
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Save to deployment.yaml
	configPath := filepath.Join(deployDir, "deployment.yaml")
	if err := os.WriteFile(configPath, yamlData, 0644); err != nil {
		return err
	}

	fmt.Printf("\nConfiguration saved to %s\n", configPath)
	fmt.Println("\nTo deploy your application, run:")
	fmt.Printf("nexlayer deploy -f %s\n\n", configPath)

	return nil
}

// GetDeploymentConfig returns the deployment configuration
func (w *DeploymentWizard) GetDeploymentConfig() (map[string]interface{}, error) {
	// TODO: Implement configuration gathering from model
	return nil, nil
}

// SaveDeploymentConfig saves the deployment configuration to disk
func (w *DeploymentWizard) SaveDeploymentConfig(config map[string]interface{}) error {
	// TODO: Implement configuration saving
	return nil
}
