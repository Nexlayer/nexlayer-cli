// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
)

func TestValidateNexlayerYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    *schema.NexlayerYAML
		wantErr bool
	}{
		{
			name: "valid yaml",
			yaml: &schema.NexlayerYAML{
				Application: schema.Application{
					Name: "myapp",
					Pods: []schema.Pod{
						{
							Name:  "frontend",
							Image: "nginx:latest",
							Path:  "/",
							Ports: []schema.Port{
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
			yaml: &schema.NexlayerYAML{
				Application: schema.Application{
					Name: "",
					Pods: []schema.Pod{},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid yaml - invalid volume size",
			yaml: &schema.NexlayerYAML{
				Application: schema.Application{
					Name: "myapp",
					Pods: []schema.Pod{
						{
							Name:  "database",
							Image: "postgres:latest",
							Ports: []schema.Port{
								{
									ContainerPort: 5432,
									ServicePort:   5432,
									Name:          "postgres",
								},
							},
							Volumes: []schema.Volume{
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
			yaml: &schema.NexlayerYAML{
				Application: schema.Application{
					Name: "myapp",
					RegistryLogin: &schema.RegistryLogin{
						Registry:           "ghcr.io",
						Username:           "myuser",
						PersonalAccessToken: "token123",
					},
					Pods: []schema.Pod{
						{
							Name:  "api",
							Image: "ghcr.io/myorg/api:latest",
							Path:  "/api",
							Ports: []schema.Port{
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
			validator := NewValidator(false)
			errs := validator.ValidateYAML(tt.yaml)
			if (len(errs) > 0) != tt.wantErr {
				t.Errorf("ValidateNexlayerYAML() errors = %v, wantErr %v", errs, tt.wantErr)
			}
		})
	}
}
