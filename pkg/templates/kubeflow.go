package templates

import (
	initcmd "github.com/Nexlayer/nexlayer-cli/pkg/commands/init"
)

// CreateKubeflowConfig creates a Kubeflow pipeline template configuration
func CreateKubeflowConfig(projectName string) initcmd.Config {
	var config initcmd.Config
	config.Application.Template.Name = projectName
	config.Application.Template.DeploymentName = projectName

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

	return config
}
