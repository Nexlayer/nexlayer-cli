// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/Nexlayer/nexlayer-cli/pkg/template"
)

func TestValidateNexlayerYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    *template.NexlayerYAML
		wantErr bool
	}{
		{
			name: "valid yaml",
			yaml: &template.NexlayerYAML{
				Application: template.Application{
					Name: "myapp",
					Pods: []template.Pod{
						{
							Name:  "frontend",
							Image: "nginx:latest",
							Path:  "/",
							Ports: []template.Port{
								{
									ContainerPort: 80,
									ServicePort:   80,
									Name:          "web",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid yaml - missing required fields",
			yaml: &template.NexlayerYAML{
				Application: template.Application{
					Name: "",
					Pods: []template.Pod{},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid yaml - invalid volume size",
			yaml: &template.NexlayerYAML{
				Application: template.Application{
					Name: "myapp",
					Pods: []template.Pod{
						{
							Name:  "database",
							Image: "postgres:latest",
							Ports: []template.Port{
								{
									ContainerPort: 5432,
									ServicePort:   5432,
									Name:          "postgres",
								},
							},
							Volumes: []template.Volume{
								{
									Name:      "data",
									Size:      "invalid",
									MountPath: "/var/lib/postgresql/data",
								},
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "valid yaml with registry login",
			yaml: &template.NexlayerYAML{
				Application: template.Application{
					Name: "myapp",
					RegistryLogin: &template.RegistryLogin{
						Registry:           "ghcr.io",
						Username:           "myuser",
						PersonalAccessToken: "token123",
					},
					Pods: []template.Pod{
						{
							Name:  "api",
							Image: "ghcr.io/myorg/api:latest",
							Path:  "/api",
							Ports: []template.Port{
								{
									ContainerPort: 8080,
									ServicePort:   8080,
									Name:          "api",
								},
							},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTemplate(tt.yaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNexlayerYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
