// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/types"
)

func TestValidateNexlayerYAML(t *testing.T) {
	tests := []struct {
		name    string
		yaml    *types.NexlayerYAML
		wantErr bool
	}{
		{
			name: "valid yaml",
			yaml: &types.NexlayerYAML{
				Application: struct {
					Name         string              `yaml:"name" validate:"required,alphanum"`
					URL          string              `yaml:"url,omitempty" validate:"omitempty,url"`
					RegistryLogin *types.RegistryLogin `yaml:"registryLogin,omitempty" validate:"omitempty"`
					Pods         []types.Pod         `yaml:"pods" validate:"required,dive,min=1"`
				}{
					Name: "myapp",
					Pods: []types.Pod{
						{
							Name:  "frontend",
							Type:  "react",
							Image: "nginx:latest",
							Path:  "/",
							ServicePorts: []int{80},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid yaml - missing required fields",
			yaml: &types.NexlayerYAML{
				Application: struct {
					Name         string              `yaml:"name" validate:"required,alphanum"`
					URL          string              `yaml:"url,omitempty" validate:"omitempty,url"`
					RegistryLogin *types.RegistryLogin `yaml:"registryLogin,omitempty" validate:"omitempty"`
					Pods         []types.Pod         `yaml:"pods" validate:"required,dive,min=1"`
				}{
					Name: "",
					Pods: []types.Pod{},
				},
			},
			wantErr: true,
		},
		{
			name: "invalid yaml - invalid volume size",
			yaml: &types.NexlayerYAML{
				Application: struct {
					Name         string              `yaml:"name" validate:"required,alphanum"`
					URL          string              `yaml:"url,omitempty" validate:"omitempty,url"`
					RegistryLogin *types.RegistryLogin `yaml:"registryLogin,omitempty" validate:"omitempty"`
					Pods         []types.Pod         `yaml:"pods" validate:"required,dive,min=1"`
				}{
					Name: "myapp",
					Pods: []types.Pod{
						{
							Name:  "database",
							Type:  "postgres",
							Image: "postgres:latest",
							Volumes: []types.Volume{
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
			yaml: &types.NexlayerYAML{
				Application: struct {
					Name         string              `yaml:"name" validate:"required,alphanum"`
					URL          string              `yaml:"url,omitempty" validate:"omitempty,url"`
					RegistryLogin *types.RegistryLogin `yaml:"registryLogin,omitempty" validate:"omitempty"`
					Pods         []types.Pod         `yaml:"pods" validate:"required,dive,min=1"`
				}{
					Name: "myapp",
					RegistryLogin: &types.RegistryLogin{
						Registry:           "ghcr.io",
						Username:           "myuser",
						PersonalAccessToken: "token123",
					},
					Pods: []types.Pod{
						{
							Name:  "api",
							Type:  "fastapi",
							Image: "ghcr.io/myorg/api:latest",
							Path:  "/api",
							ServicePorts: []int{8080},
						},
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNexlayerYAML(tt.yaml)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateNexlayerYAML() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
