package templates

import (
	"github.com/Nexlayer/nexlayer-cli/pkg/commands/initcmd"
)

// CreateKubeflowConfig creates a Kubeflow pipeline template configuration
func CreateKubeflowConfig(projectName string) initcmd.Config {
	var config initcmd.Config
	config.Application.Template.Name = projectName
	config.Application.Template.DeploymentName = projectName

	// Set registry login
	config.Application.Template.RegistryLogin = initcmd.RegistryAuth{
		Registry: "ghcr.io",
		Username: "<Github username>",
		PersonalAccessToken: "<Github Packages Read-Only PAT>",
	}

	// Add Kubeflow pipeline pod
	config.Application.Template.Pods = []initcmd.PodConfig{
		{
			Type: "llm",
			Name: "ml-pipeline",
			Tag:  "python:3.11-slim",
			Vars: []initcmd.VarPair{
				{Key: "PIPELINE_ROOT", Value: "/tmp/pipeline"},
				{Key: "DATA_PATH", Value: "/tmp/data"},
				{Key: "MODEL_PATH", Value: "/tmp/model"},
				{Key: "KUBEFLOW_URL", Value: "http://localhost:8080"},
			},
			ExposeHttp: true,
		},
	}

	// Set build configuration
	config.Application.Template.Build.Command = "pip install -r requirements.txt"
	config.Application.Template.Build.Output = "dist"

	return config
}
