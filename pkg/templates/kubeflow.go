package templates

import (
	"github.com/Nexlayer/nexlayer-cli/pkg/types"
)

// CreateKubeflowConfig creates a Kubeflow pipeline template configuration
func CreateKubeflowConfig(projectName string) types.Config {
	config := types.Config{
		Application: types.Application{
			Template: types.Template{
				Name:           projectName,
				DeploymentName: projectName,
				RegistryLogin: types.RegistryAuth{
					Registry:            "ghcr.io",
					Username:           "<Github username>",
					PersonalAccessToken: "<Github Packages Read-Only PAT>",
				},
				Pods: []types.PodConfig{
					{
						Type: "llm",
						Name: "ml-pipeline",
						Tag:  "python:3.11-slim",
						Vars: []types.VarPair{
							{Key: "PIPELINE_ROOT", Value: "/tmp/pipeline"},
							{Key: "DATA_PATH", Value: "/tmp/data"},
							{Key: "MODEL_PATH", Value: "/tmp/model"},
							{Key: "KUBEFLOW_URL", Value: "http://localhost:8080"},
						},
						ExposeHttp: true,
					},
				},
				Build: struct {
					Command string `yaml:"command"`
					Output  string `yaml:"output"`
				}{
					Command: "pip install -r requirements.txt",
					Output:  "dist",
				},
			},
		},
	}

	return config
}
