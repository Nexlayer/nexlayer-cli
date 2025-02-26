// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package template

import (
	"fmt"

	"github.com/Nexlayer/nexlayer-cli/pkg/schema"
)

// NexlayerYAML represents the structure of a Nexlayer configuration file
// This is a serialization-friendly version of schema.NexlayerYAML
type NexlayerYAML struct {
	Application ApplicationYAML `yaml:"application"`
}

// ApplicationYAML represents the application section of a Nexlayer configuration
// This is a serialization-friendly version of schema.Application
type ApplicationYAML struct {
	Name          string             `yaml:"name"`
	URL           string             `yaml:"url"`
	RegistryLogin *RegistryLoginYAML `yaml:"registry_login,omitempty"`
	Pods          []PodYAML          `yaml:"pods"`
}

// RegistryLoginYAML represents registry login information
// This is a serialization-friendly version of schema.RegistryLogin
type RegistryLoginYAML struct {
	Registry string `yaml:"registry"`
	Username string `yaml:"username"`
	Password string `yaml:"password,omitempty"`
}

// PodYAML represents a pod configuration
// This is a serialization-friendly version of schema.Pod
type PodYAML struct {
	Name         string            `yaml:"name"`
	Type         string            `yaml:"type,omitempty"`
	Path         string            `yaml:"path,omitempty"`
	Image        string            `yaml:"image"`
	Command      string            `yaml:"command,omitempty"`
	Entrypoint   string            `yaml:"entrypoint,omitempty"`
	ServicePorts []ServicePort     `yaml:"servicePorts,omitempty"`
	Vars         []EnvVar          `yaml:"vars,omitempty"`
	Volumes      []Volume          `yaml:"volumes,omitempty"`
	Secrets      []Secret          `yaml:"secrets,omitempty"`
	Annotations  map[string]string `yaml:"annotations,omitempty"`
}

// Application represents the application configuration
type Application struct {
	Name          string         `yaml:"name" validate:"required,name"`
	URL           string         `yaml:"url,omitempty" validate:"omitempty,url"`
	RegistryLogin *RegistryLogin `yaml:"registryLogin,omitempty" validate:"omitempty"`
	Pods          []Pod          `yaml:"pods" validate:"required,min=1,dive"`
}

// RegistryLogin represents container registry authentication
type RegistryLogin struct {
	Registry            string `yaml:"registry" validate:"required"`
	Username            string `yaml:"username" validate:"required"`
	PersonalAccessToken string `yaml:"personalAccessToken" validate:"required"`
}

// Pod represents a container in the deployment
type Pod struct {
	Name         string            `yaml:"name" validate:"required,name"`
	Type         string            `yaml:"type" validate:"required,oneof=frontend backend nextjs react node python go raw"`
	Path         string            `yaml:"path,omitempty" validate:"omitempty,startswith=/"`
	Image        string            `yaml:"image" validate:"required,image"`
	Command      string            `yaml:"command,omitempty"`
	Entrypoint   string            `yaml:"entrypoint,omitempty"`
	ServicePorts []ServicePort     `yaml:"servicePorts" validate:"required,min=1,dive"`
	Vars         []EnvVar          `yaml:"vars,omitempty" validate:"omitempty,dive"`
	Volumes      []Volume          `yaml:"volumes,omitempty" validate:"omitempty,dive"`
	Secrets      []Secret          `yaml:"secrets,omitempty" validate:"omitempty,dive"`
	Annotations  map[string]string `yaml:"annotations,omitempty"`
}

// ServicePort represents a service port configuration
type ServicePort struct {
	Name       string `yaml:"name,omitempty" validate:"omitempty"`
	Port       int    `yaml:"port" validate:"required,min=1,max=65535"`
	TargetPort int    `yaml:"targetPort,omitempty" validate:"omitempty,min=1,max=65535"`
	Protocol   string `yaml:"protocol,omitempty" validate:"omitempty,oneof=TCP UDP"`
}

// UnmarshalYAML implements custom unmarshaling for ServicePort to support both formats
func (sp *ServicePort) UnmarshalYAML(unmarshal func(interface{}) error) error {
	// Try simple format (just port number)
	var port int
	if err := unmarshal(&port); err == nil {
		sp.Port = port
		sp.TargetPort = port
		sp.Name = fmt.Sprintf("port-%d", port)
		sp.Protocol = ProtocolTCP
		return nil
	}

	// Try full format
	type fullServicePort ServicePort
	var full fullServicePort
	if err := unmarshal(&full); err != nil {
		return err
	}

	sp.Name = full.Name
	sp.Port = full.Port
	sp.TargetPort = full.TargetPort
	if sp.TargetPort == 0 {
		sp.TargetPort = sp.Port
	}
	sp.Protocol = full.Protocol
	if sp.Protocol == "" {
		sp.Protocol = ProtocolTCP
	}
	if sp.Name == "" {
		sp.Name = fmt.Sprintf("port-%d", sp.Port)
	}

	return nil
}

// MarshalYAML implements custom marshaling for ServicePort to support the simple format
func (sp ServicePort) MarshalYAML() (interface{}, error) {
	// If it's a simple port with default values, return just the port number
	if sp.Name == fmt.Sprintf("port-%d", sp.Port) &&
		sp.TargetPort == sp.Port &&
		(sp.Protocol == ProtocolTCP || sp.Protocol == "") {
		return sp.Port, nil
	}

	// Otherwise return the full structure
	return struct {
		Name       string `yaml:"name,omitempty"`
		Port       int    `yaml:"port"`
		TargetPort int    `yaml:"targetPort,omitempty"`
		Protocol   string `yaml:"protocol,omitempty"`
	}{
		Name:       sp.Name,
		Port:       sp.Port,
		TargetPort: sp.TargetPort,
		Protocol:   sp.Protocol,
	}, nil
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `yaml:"key" validate:"required,envvar"`
	Value string `yaml:"value" validate:"required"`
}

// Volume represents a persistent storage volume
type Volume struct {
	Name     string `yaml:"name" validate:"required,name"`
	Path     string `yaml:"path" validate:"required,startswith=/"`
	Size     string `yaml:"size,omitempty" validate:"omitempty,volumesize"`
	Type     string `yaml:"type,omitempty" validate:"omitempty,oneof=persistent ephemeral"`
	ReadOnly bool   `yaml:"readOnly,omitempty"`
}

// Secret represents a secret configuration
type Secret struct {
	Name     string `yaml:"name" validate:"required,name"`
	Data     string `yaml:"data" validate:"required"`
	Path     string `yaml:"path" validate:"required,startswith=/"`
	FileName string `yaml:"fileName" validate:"required"`
}

// Conversion functions between schema types and template types

// ToSchemaType converts the template NexlayerYAML to the schema NexlayerYAML
func (t *NexlayerYAML) ToSchemaType() *schema.NexlayerYAML {
	result := &schema.NexlayerYAML{
		Application: schema.Application{
			Name: t.Application.Name,
			URL:  t.Application.URL,
		},
	}

	// Convert registry login
	if t.Application.RegistryLogin != nil {
		result.Application.RegistryLogin = &schema.RegistryLogin{
			Registry:            t.Application.RegistryLogin.Registry,
			Username:            t.Application.RegistryLogin.Username,
			PersonalAccessToken: t.Application.RegistryLogin.Password, // Map Password to PersonalAccessToken
		}
	}

	// Convert pods
	for _, pod := range t.Application.Pods {
		schemaPod := schema.Pod{
			Name:        pod.Name,
			Type:        pod.Type,
			Path:        pod.Path,
			Image:       pod.Image,
			Command:     pod.Command,
			Entrypoint:  pod.Entrypoint,
			Annotations: pod.Annotations,
		}

		// Convert service ports
		for _, port := range pod.ServicePorts {
			// For simple integer ports, convert to the simple format
			if port.Name == fmt.Sprintf("port-%d", port.Port) &&
				port.TargetPort == port.Port &&
				(port.Protocol == ProtocolTCP || port.Protocol == "") {
				schemaPod.ServicePorts = append(schemaPod.ServicePorts, port.Port)
			} else {
				// For complex ports, use the full structure
				schemaPod.ServicePorts = append(schemaPod.ServicePorts, map[string]interface{}{
					"name":       port.Name,
					"port":       port.Port,
					"targetPort": port.TargetPort,
					"protocol":   port.Protocol,
				})
			}
		}

		// Convert environment variables
		for _, envVar := range pod.Vars {
			schemaPod.Vars = append(schemaPod.Vars, schema.EnvVar{
				Key:   envVar.Key,
				Value: envVar.Value,
			})
		}

		// Convert volumes
		for _, volume := range pod.Volumes {
			schemaPod.Volumes = append(schemaPod.Volumes, schema.Volume{
				Name:     volume.Name,
				Path:     volume.Path,
				Size:     volume.Size,
				Type:     volume.Type,
				ReadOnly: volume.ReadOnly,
			})
		}

		// Convert secrets
		for _, secret := range pod.Secrets {
			schemaPod.Secrets = append(schemaPod.Secrets, schema.Secret{
				Name:     secret.Name,
				Data:     secret.Data,
				Path:     secret.Path,
				FileName: secret.FileName,
			})
		}

		result.Application.Pods = append(result.Application.Pods, schemaPod)
	}

	return result
}

// FromSchemaType converts a schema.NexlayerYAML to template.NexlayerYAML
func FromSchemaType(s *schema.NexlayerYAML) *NexlayerYAML {
	result := &NexlayerYAML{
		Application: ApplicationYAML{
			Name: s.Application.Name,
			URL:  s.Application.URL,
		},
	}

	// Convert registry login
	if s.Application.RegistryLogin != nil {
		result.Application.RegistryLogin = &RegistryLoginYAML{
			Registry: s.Application.RegistryLogin.Registry,
			Username: s.Application.RegistryLogin.Username,
			Password: s.Application.RegistryLogin.PersonalAccessToken, // Map PersonalAccessToken to Password
		}
	}

	// Convert pods
	for _, pod := range s.Application.Pods {
		templatePod := PodYAML{
			Name:        pod.Name,
			Type:        pod.Type,
			Path:        pod.Path,
			Image:       pod.Image,
			Command:     pod.Command,
			Entrypoint:  pod.Entrypoint,
			Annotations: pod.Annotations,
		}

		// Convert service ports
		for _, port := range pod.ServicePorts {
			switch p := port.(type) {
			case int:
				// Simple integer port
				templatePod.ServicePorts = append(templatePod.ServicePorts, ServicePort{
					Name:       fmt.Sprintf("port-%d", p),
					Port:       p,
					TargetPort: p,
					Protocol:   ProtocolTCP,
				})
			case float64:
				// JSON numbers might be unmarshaled as float64
				intPort := int(p)
				templatePod.ServicePorts = append(templatePod.ServicePorts, ServicePort{
					Name:       fmt.Sprintf("port-%d", intPort),
					Port:       intPort,
					TargetPort: intPort,
					Protocol:   ProtocolTCP,
				})
			case map[string]interface{}:
				// Object format
				sp := ServicePort{
					Protocol: ProtocolTCP, // Default protocol
				}

				if name, ok := p["name"].(string); ok {
					sp.Name = name
				}

				if portVal, ok := p["port"].(float64); ok {
					sp.Port = int(portVal)
					sp.TargetPort = int(portVal) // Default target port to port
				} else if portVal, ok := p["port"].(int); ok {
					sp.Port = portVal
					sp.TargetPort = portVal
				}

				if targetPort, ok := p["targetPort"].(float64); ok {
					sp.TargetPort = int(targetPort)
				} else if targetPort, ok := p["targetPort"].(int); ok {
					sp.TargetPort = targetPort
				}

				if protocol, ok := p["protocol"].(string); ok {
					sp.Protocol = protocol
				}

				// Generate name if not provided
				if sp.Name == "" {
					sp.Name = fmt.Sprintf("port-%d", sp.Port)
				}

				templatePod.ServicePorts = append(templatePod.ServicePorts, sp)
			}
		}

		// Convert environment variables
		for _, envVar := range pod.Vars {
			templatePod.Vars = append(templatePod.Vars, EnvVar{
				Key:   envVar.Key,
				Value: envVar.Value,
			})
		}

		// Convert volumes
		for _, volume := range pod.Volumes {
			templatePod.Volumes = append(templatePod.Volumes, Volume{
				Name:     volume.Name,
				Path:     volume.Path,
				Size:     volume.Size,
				Type:     volume.Type,
				ReadOnly: volume.ReadOnly,
			})
		}

		// Convert secrets
		for _, secret := range pod.Secrets {
			templatePod.Secrets = append(templatePod.Secrets, Secret{
				Name:     secret.Name,
				Data:     secret.Data,
				Path:     secret.Path,
				FileName: secret.FileName,
			})
		}

		result.Application.Pods = append(result.Application.Pods, templatePod)
	}

	return result
}

// ConvertPodToPodYAML converts a schema.Pod to a template.PodYAML
func ConvertPodToPodYAML(pod schema.Pod) PodYAML {
	result := PodYAML{
		Name:         pod.Name,
		Type:         pod.Type,
		Path:         pod.Path,
		Image:        pod.Image,
		Command:      pod.Command,
		Entrypoint:   pod.Entrypoint,
		ServicePorts: make([]ServicePort, 0, len(pod.ServicePorts)),
		Vars:         ConvertEnvVars(pod.Vars),
		Volumes:      ConvertVolumes(pod.Volumes),
		Secrets:      ConvertSecrets(pod.Secrets),
		Annotations:  pod.Annotations,
	}

	// Convert service ports based on type
	for _, port := range pod.ServicePorts {
		switch p := port.(type) {
		case int:
			// Simple integer port
			result.ServicePorts = append(result.ServicePorts, ServicePort{
				Name:       fmt.Sprintf("port-%d", p),
				Port:       p,
				TargetPort: p,
				Protocol:   ProtocolTCP,
			})
		case map[string]interface{}:
			// Object format
			sp := ServicePort{
				Protocol: ProtocolTCP, // Default protocol
			}

			if name, ok := p["name"].(string); ok {
				sp.Name = name
			}

			if portVal, ok := p["port"].(float64); ok {
				sp.Port = int(portVal)
				sp.TargetPort = int(portVal) // Default target port to port
			}

			if targetPort, ok := p["targetPort"].(float64); ok {
				sp.TargetPort = int(targetPort)
			}

			if protocol, ok := p["protocol"].(string); ok {
				sp.Protocol = protocol
			}

			// Generate name if not provided
			if sp.Name == "" {
				sp.Name = fmt.Sprintf("port-%d", sp.Port)
			}

			result.ServicePorts = append(result.ServicePorts, sp)
		}
	}

	return result
}

// ConvertServicePorts converts schema service ports to template service ports
func ConvertServicePorts(ports []interface{}) []ServicePort {
	result := make([]ServicePort, 0, len(ports))

	for _, port := range ports {
		switch p := port.(type) {
		case int:
			// Simple integer port
			result = append(result, ServicePort{
				Name:       fmt.Sprintf("port-%d", p),
				Port:       p,
				TargetPort: p,
				Protocol:   ProtocolTCP,
			})
		case map[string]interface{}:
			// Object format
			sp := ServicePort{}

			if name, ok := p["name"].(string); ok {
				sp.Name = name
			}

			if portVal, ok := p["port"].(int); ok {
				sp.Port = portVal
				sp.TargetPort = portVal // Default target port to port
			}

			if targetPort, ok := p["targetPort"].(int); ok {
				sp.TargetPort = targetPort
			}

			if protocol, ok := p["protocol"].(string); ok {
				sp.Protocol = protocol
			} else {
				sp.Protocol = ProtocolTCP
			}

			// Generate name if not provided
			if sp.Name == "" {
				sp.Name = fmt.Sprintf("port-%d", sp.Port)
			}

			result = append(result, sp)
		case ServicePort:
			// Already a ServicePort
			result = append(result, p)
		}
	}

	return result
}

// ConvertEnvVars converts schema.EnvVar slice to template.EnvVar slice
func ConvertEnvVars(vars []schema.EnvVar) []EnvVar {
	var result []EnvVar
	for _, v := range vars {
		result = append(result, EnvVar{
			Key:   v.Key,
			Value: v.Value,
		})
	}
	return result
}

// ConvertVolumes converts schema.Volume slice to template.Volume slice
func ConvertVolumes(volumes []schema.Volume) []Volume {
	var result []Volume
	for _, v := range volumes {
		result = append(result, Volume{
			Name:     v.Name,
			Path:     v.Path,
			Size:     v.Size,
			Type:     v.Type,
			ReadOnly: v.ReadOnly,
		})
	}
	return result
}

// ConvertSecrets converts schema.Secret slice to template.Secret slice
func ConvertSecrets(secrets []schema.Secret) []Secret {
	var result []Secret
	for _, s := range secrets {
		result = append(result, Secret{
			Name:     s.Name,
			Data:     s.Data,
			Path:     s.Path,
			FileName: s.FileName,
		})
	}
	return result
}
