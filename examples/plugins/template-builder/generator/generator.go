package generator

import (
	"fmt"
	"strings"

	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/types"
)

// GenerateTemplate creates a template based on project name and detected stack
func GenerateTemplate(projectName string, stack *types.ProjectStack) *types.NexlayerTemplate {
	tmpl := &types.NexlayerTemplate{}
	t := &tmpl.Application.Template

	// Basic template info
	t.Name = projectName
	t.TemplateID = fmt.Sprintf("%s-NrqurwWL", strings.ToLower(projectName))
	t.DeploymentName = projectName

	// Registry login info
	t.RegistryLogin = &types.Registry{
		Registry:           "ghcr.io",
		Username:          "nexlayer",
		PersonalAccessToken: "${GITHUB_TOKEN}",
	}

	// Add pods based on stack
	if stack.HasFrontend() {
		t.Pods = append(t.Pods, generateFrontendPod(projectName, stack))
	}
	if stack.HasBackend() {
		t.Pods = append(t.Pods, generateBackendPod(projectName, stack))
	}
	if stack.HasDatabase() {
		t.Pods = append(t.Pods, generateDatabasePod(projectName, stack))
	}

	return tmpl
}

func generateFrontendPod(projectName string, stack *types.ProjectStack) types.Pod {
	vars := []types.EnvVar{
		{
			Key:   "API_URL",
			Value: "BACKEND_CONNECTION_URL",
		},
	}

	// Add NODE_ENV for Node.js projects
	if stack.Language == "nodejs" {
		vars = append(vars, types.EnvVar{
			Key:   "NODE_ENV",
			Value: "development",
		})
	}

	return types.Pod{
		Type:       "web",
		Name:       "frontend",
		Image:      fmt.Sprintf("ghcr.io/nexlayer/%s-frontend", strings.ToLower(projectName)),
		Tag:        "latest",
		PrivateTag: true,
		ExposeHttp: true,
		Vars:       vars,
	}
}

func generateBackendPod(projectName string, stack *types.ProjectStack) types.Pod {
	vars := []types.EnvVar{
		{
			Key:   "DATABASE_URL",
			Value: "DATABASE_CONNECTION_STRING",
		},
	}

	// Add language-specific env vars
	switch stack.Language {
	case "nodejs":
		vars = append(vars, types.EnvVar{
			Key:   "NODE_ENV",
			Value: "development",
		})
	case "python":
		vars = append(vars, types.EnvVar{
			Key:   "PYTHONPATH",
			Value: "/app",
		})
	case "go":
		vars = append(vars, types.EnvVar{
			Key:   "GO_ENV",
			Value: "development",
		})
	}

	return types.Pod{
		Type:       "api",
		Name:       "backend",
		Image:      fmt.Sprintf("ghcr.io/nexlayer/%s-backend", strings.ToLower(projectName)),
		Tag:        "latest",
		PrivateTag: true,
		ExposeHttp: true,
		Vars:       vars,
	}
}

func generateDatabasePod(projectName string, stack *types.ProjectStack) types.Pod {
	dbType := stack.GetDatabaseType()
	vars := []types.EnvVar{
		{
			Key:   "POSTGRES_USER",
			Value: "postgres",
		},
		{
			Key:   "POSTGRES_PASSWORD",
			Value: "postgres",
		},
		{
			Key:   "POSTGRES_DB",
			Value: strings.ToLower(projectName),
		},
	}

	// Add database-specific env vars
	switch dbType {
	case "mongodb":
		vars = []types.EnvVar{
			{
				Key:   "MONGO_INITDB_ROOT_USERNAME",
				Value: "mongo",
			},
			{
				Key:   "MONGO_INITDB_ROOT_PASSWORD",
				Value: "mongo",
			},
			{
				Key:   "MONGO_INITDB_DATABASE",
				Value: strings.ToLower(projectName),
			},
		}
	}

	return types.Pod{
		Type:       "database",
		Name:       "database",
		Image:      fmt.Sprintf("ghcr.io/nexlayer/%s-%s", strings.ToLower(projectName), dbType),
		Tag:        "latest",
		PrivateTag: true,
		ExposeHttp: false,
		Vars:       vars,
	}
}
