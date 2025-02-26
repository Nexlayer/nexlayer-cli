package schema

import (
	"fmt"
	"strings"
)

// Generator handles generation of Nexlayer YAML configurations
type Generator struct {
	defaultPort int
}

// NewGenerator creates a new schema generator
func NewGenerator() *Generator {
	return &Generator{
		defaultPort: 3000,
	}
}

// GenerateFromProjectInfo generates a Nexlayer YAML configuration from project info
func (g *Generator) GenerateFromProjectInfo(name string, projectType string, port int) (*NexlayerYAML, error) {
	if name == "" {
		name = "my-app"
	}

	if port == 0 {
		port = g.defaultPort
	}

	// Create base configuration
	config := &NexlayerYAML{
		Application: Application{
			Name: name,
			Pods: make([]Pod, 0),
		},
	}

	// Add main pod based on project type
	mainPod := g.generateMainPod(projectType, port)
	config.Application.Pods = append(config.Application.Pods, mainPod)

	return config, nil
}

// generateMainPod creates the main application pod
func (g *Generator) generateMainPod(projectType string, port int) Pod {
	pod := Pod{
		Name:  "main",
		Type:  projectType,
		Image: fmt.Sprintf("<%% REGISTRY %%>/%s:latest", strings.ToLower(projectType)),
		ServicePorts: []interface{}{
			port,
		},
	}

	// Add type-specific configuration
	switch projectType {
	case "nextjs", "react":
		pod.Name = "web"
		pod.Path = "/"
		pod.Vars = append(pod.Vars, EnvVar{
			Key:   "NODE_ENV",
			Value: "production",
		})

	case "node", "python", "go":
		pod.Name = "api"
		pod.Path = "/api"
		pod.Vars = append(pod.Vars, EnvVar{
			Key:   "PORT",
			Value: fmt.Sprintf("%d", port),
		})
	}

	return pod
}

// AddPod adds a new pod to the configuration
func (g *Generator) AddPod(config *NexlayerYAML, podType string, port int) error {
	if port == 0 {
		port = g.getDefaultPortForType(podType)
	}

	pod := Pod{
		Name:  g.generatePodName(podType, len(config.Application.Pods)),
		Type:  podType,
		Image: g.getDefaultImageForType(podType),
		ServicePorts: []interface{}{
			port,
		},
	}

	// Add type-specific configuration
	switch podType {
	case "postgres", "mysql", "mongodb":
		pod.Volumes = []Volume{
			{
				Name: fmt.Sprintf("%s-data", pod.Name),
				Path: g.getDefaultPathForType(podType),
				Size: "1Gi",
			},
		}
		pod.Vars = g.getDefaultEnvVarsForType(podType)

	case "redis":
		pod.Vars = []EnvVar{
			{Key: "REDIS_PASSWORD", Value: "<% REDIS_PASSWORD %>"},
		}

	case "minio":
		pod.Volumes = []Volume{
			{
				Name: "minio-data",
				Path: "/data",
				Size: "5Gi",
			},
		}
		pod.Vars = []EnvVar{
			{Key: "MINIO_ROOT_USER", Value: "<% MINIO_ROOT_USER %>"},
			{Key: "MINIO_ROOT_PASSWORD", Value: "<% MINIO_ROOT_PASSWORD %>"},
		}
	}

	config.Application.Pods = append(config.Application.Pods, pod)
	return nil
}

// Helper functions

func (g *Generator) generatePodName(podType string, index int) string {
	base := strings.ToLower(podType)
	if index == 0 {
		return base
	}
	return fmt.Sprintf("%s-%d", base, index+1)
}

func (g *Generator) getDefaultPortForType(podType string) int {
	switch podType {
	case "postgres":
		return 5432
	case "mysql":
		return 3306
	case "mongodb":
		return 27017
	case "redis":
		return 6379
	case "minio":
		return 9000
	default:
		return g.defaultPort
	}
}

func (g *Generator) getDefaultPortNameForType(podType string) string {
	switch podType {
	case "postgres", "mysql", "mongodb", "redis":
		return "db"
	case "minio":
		return "api"
	default:
		return "http"
	}
}

func (g *Generator) getDefaultImageForType(podType string) string {
	switch podType {
	case "postgres":
		return "postgres:latest"
	case "mysql":
		return "mysql:8"
	case "mongodb":
		return "mongo:latest"
	case "redis":
		return "redis:alpine"
	case "minio":
		return "minio/minio:latest"
	default:
		return fmt.Sprintf("<%% REGISTRY %%>/%s:latest", podType)
	}
}

func (g *Generator) getDefaultPathForType(podType string) string {
	switch podType {
	case "postgres":
		return "/var/lib/postgresql/data"
	case "mysql":
		return "/var/lib/mysql"
	case "mongodb":
		return "/data/db"
	default:
		return "/data"
	}
}

func (g *Generator) getDefaultEnvVarsForType(podType string) []EnvVar {
	switch podType {
	case "postgres":
		return []EnvVar{
			{Key: "POSTGRES_USER", Value: "<% POSTGRES_USER %>"},
			{Key: "POSTGRES_PASSWORD", Value: "<% POSTGRES_PASSWORD %>"},
			{Key: "POSTGRES_DB", Value: "<% POSTGRES_DB %>"},
		}
	case "mysql":
		return []EnvVar{
			{Key: "MYSQL_ROOT_PASSWORD", Value: "<% MYSQL_ROOT_PASSWORD %>"},
			{Key: "MYSQL_DATABASE", Value: "<% MYSQL_DATABASE %>"},
			{Key: "MYSQL_USER", Value: "<% MYSQL_USER %>"},
			{Key: "MYSQL_PASSWORD", Value: "<% MYSQL_PASSWORD %>"},
		}
	case "mongodb":
		return []EnvVar{
			{Key: "MONGO_INITDB_ROOT_USERNAME", Value: "<% MONGO_ROOT_USERNAME %>"},
			{Key: "MONGO_INITDB_ROOT_PASSWORD", Value: "<% MONGO_ROOT_PASSWORD %>"},
		}
	default:
		return nil
	}
}

// AddAIConfigurations adds AI-specific configurations to pods
func (g *Generator) AddAIConfigurations(config *NexlayerYAML, llmProvider string) {
	for i := range config.Application.Pods {
		if config.Application.Pods[i].Annotations == nil {
			config.Application.Pods[i].Annotations = make(map[string]string)
		}
		config.Application.Pods[i].Annotations["ai.nexlayer.io/provider"] = llmProvider
		config.Application.Pods[i].Annotations["ai.nexlayer.io/enabled"] = "true"
	}
}

// AddEnvironmentVars adds environment variables to a pod
func (g *Generator) AddEnvironmentVars(pod *Pod, vars map[string]string) {
	for k, v := range vars {
		pod.Vars = append(pod.Vars, EnvVar{
			Key:   k,
			Value: v,
		})
	}
}

// AddVolume adds a volume to a pod
func (g *Generator) AddVolume(pod *Pod, name, path, size string) {
	pod.Volumes = append(pod.Volumes, Volume{
		Name: name,
		Path: path,
		Size: size,
	})
}

// AddSecret adds a secret to a pod
func (g *Generator) AddSecret(pod *Pod, name, data, path, fileName string) {
	pod.Secrets = append(pod.Secrets, Secret{
		Name:     name,
		Data:     data,
		Path:     path,
		FileName: fileName,
	})
}

// AddServicePort adds a service port to a pod
func (g *Generator) AddServicePort(pod *Pod, name string, port, targetPort int) {
	// If simple port (same port and target port with default name), use integer format
	if name == fmt.Sprintf("port-%d", port) && port == targetPort {
		pod.ServicePorts = append(pod.ServicePorts, port)
	} else {
		// Otherwise use the structured format
		pod.ServicePorts = append(pod.ServicePorts, map[string]interface{}{
			"name":       name,
			"port":       port,
			"targetPort": targetPort,
		})
	}
}

// SetRegistryLogin sets the registry login configuration
func (g *Generator) SetRegistryLogin(config *NexlayerYAML, registry, username, token string) {
	config.Application.RegistryLogin = &RegistryLogin{
		Registry:            registry,
		Username:            username,
		PersonalAccessToken: token,
	}
}

// SetCustomDomain sets the custom domain for the application
func (g *Generator) SetCustomDomain(config *NexlayerYAML, domain string) {
	config.Application.URL = domain
}
