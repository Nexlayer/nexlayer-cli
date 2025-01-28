package wizard

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api"
	"github.com/Nexlayer/nexlayer-cli/pkg/ui"
)

type Template struct {
	Name           string `yaml:"name"`
	DeploymentName string `yaml:"deploymentName"`
}

type Pod struct {
	Type       string `yaml:"type"`
	ExposeHttp bool   `yaml:"exposeHttp"`
	Name       string `yaml:"name"`
	Tag        string `yaml:"tag,omitempty"`
	PrivateTag bool   `yaml:"privateTag,omitempty"`
	Vars       []Var  `yaml:"vars,omitempty"`
}

type Var struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type Application struct {
	Template Template `yaml:"template"`
	Pods     []Pod    `yaml:"pods,omitempty"`
}

type DeploymentConfig struct {
	Application Application `yaml:"application"`
}

var validTemplates = map[string]bool{
	"langchain-nextjs":  true,
	"langchain-fastapi": true,
	"mern":             true,
	"pern":             true,
	"mean":             true,
}

var templateToPods = map[string][]Pod{
	"langchain-nextjs": {
		{
			Type:       "nextjs",
			ExposeHttp: true,
			Name:      "app",
			Vars: []Var{
				{Key: "OPENAI_API_KEY", Value: "your-key"},
				{Key: "LANGCHAIN_TRACING_V2", Value: "true"},
			},
		},
	},
	"langchain-fastapi": {
		{
			Type:       "fastapi",
			ExposeHttp: true,
			Name:      "backend",
			Vars: []Var{
				{Key: "OPENAI_API_KEY", Value: "your-key"},
				{Key: "PINECONE_API_KEY", Value: "your-key"},
				{Key: "PINECONE_ENVIRONMENT", Value: "gcp-starter"},
			},
		},
	},
	"mern": {
		{
			Type:       "database",
			ExposeHttp: false,
			Name:      "mongodb",
			Vars: []Var{
				{Key: "MONGO_INITDB_DATABASE", Value: "myapp"},
			},
		},
		{
			Type:       "express",
			ExposeHttp: false,
			Name:      "backend",
			Vars: []Var{
				{Key: "MONGODB_URL", Value: "DATABASE_CONNECTION_STRING"},
			},
		},
		{
			Type:       "nginx",
			ExposeHttp: true,
			Name:      "frontend",
			Vars: []Var{
				{Key: "EXPRESS_URL", Value: "BACKEND_CONNECTION_URL"},
			},
		},
	},
}

func NewCommand(client *api.Client) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wizard",
		Short: "Interactive deployment wizard",
		Long: `Create a new deployment using an interactive wizard.
		
The wizard will guide you through:
1. Choosing a project name
2. Selecting a template
3. Creating a deployment configuration file`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runWizard(cmd)
		},
	}

	return cmd
}

func runWizard(cmd *cobra.Command) error {
	cmd.Println(ui.RenderTitleWithBorder("Deployment Wizard"))

	// Get project name
	var projectName string
	cmd.Print("Enter project name: ")
	fmt.Scanln(&projectName)

	if projectName == "" {
		return fmt.Errorf("project name is required")
	}

	// Get template name
	var templateName string
	cmd.Print("Enter template name (e.g., langchain-nextjs, langchain-fastapi): ")
	fmt.Scanln(&templateName)

	if !validTemplates[templateName] {
		return fmt.Errorf("invalid template: %s", templateName)
	}

	// Create config
	config := DeploymentConfig{
		Application: Application{
			Template: Template{
				Name:           templateName,
				DeploymentName: projectName,
			},
			Pods: templateToPods[templateName],
		},
	}

	// Save config to file
	configPath := "nexlayer.yaml"
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to create config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	cmd.Printf("\nConfiguration saved to %s\n", configPath)
	cmd.Println("\nTo deploy your application, run:")
	cmd.Printf("  nexlayer deploy\n")

	return nil
}
