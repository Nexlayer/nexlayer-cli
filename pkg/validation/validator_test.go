// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package validation

import (
	"testing"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/api/schema"
	"github.com/stretchr/testify/assert"
)

func TestValidateYAML(t *testing.T) {
	tests := []struct {
		name          string
		yaml          *schema.NexlayerYAML
		expectedError bool
		errorCount    int
	}{
		{
			name: "Valid configuration",
			yaml: &schema.NexlayerYAML{
				Application: schema.Application{
					Name: "test-app",
					Pods: []schema.Pod{
						{
							Name:  "web",
							Image: "nginx:latest",
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
			expectedError: false,
		},
		{
			name: "Missing application name",
			yaml: &schema.NexlayerYAML{
				Application: schema.Application{
					Pods: []schema.Pod{
						{
							Name:  "web",
							Image: "nginx:latest",
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
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Invalid pod name",
			yaml: &schema.NexlayerYAML{
				Application: schema.Application{
					Name: "test-app",
					Pods: []schema.Pod{
						{
							Name:  "Web_Server",
							Image: "nginx:latest",
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
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Invalid image format",
			yaml: &schema.NexlayerYAML{
				Application: schema.Application{
					Name: "test-app",
					Pods: []schema.Pod{
						{
							Name:  "web",
							Image: "nginx::latest",
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
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Invalid volume size",
			yaml: &schema.NexlayerYAML{
				Application: schema.Application{
					Name: "test-app",
					Pods: []schema.Pod{
						{
							Name:  "web",
							Image: "nginx:latest",
							Ports: []schema.Port{
								{
									ContainerPort: 80,
									ServicePort:   80,
									Name:          "web",
								},
							},
							Volumes: []schema.Volume{
								{
									Name:      "data",
									Size:      "1G",
									MountPath: "/data",
								},
							},
						},
					},
				},
			},
			expectedError: true,
			errorCount:    1,
		},
		{
			name: "Missing registry credentials",
			yaml: &schema.NexlayerYAML{
				Application: schema.Application{
					Name: "test-app",
					RegistryLogin: &schema.RegistryLogin{
						Registry: "registry.example.com",
					},
					Pods: []schema.Pod{
						{
							Name:  "web",
							Image: "nginx:latest",
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
			expectedError: true,
			errorCount:    2, // Missing username and token
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator := NewValidator(true)
			errors := validator.ValidateYAML(tt.yaml)

			if tt.expectedError {
				assert.Len(t, errors, tt.errorCount)
			} else {
				assert.Empty(t, errors)
			}
		})
	}
}

func TestValidationHelpers(t *testing.T) {
	// Test isValidName
	validNames := []string{"web", "api-server", "db-1"}
	invalidNames := []string{"Web", "api_server", "-web", "web-", "a@b"}

	for _, name := range validNames {
		assert.True(t, isValidName(name), "Expected %s to be valid", name)
	}

	for _, name := range invalidNames {
		assert.False(t, isValidName(name), "Expected %s to be invalid", name)
	}

	// Test isValidImageName
	validImages := []string{
		"nginx",
		"nginx:latest",
		"docker.io/library/nginx:latest",
		"registry.example.com/org/app:v1.2.3",
	}
	invalidImages := []string{
		"",
		"nginx:",
		":latest",
		"nginx::latest",
		"/nginx",
		"nginx/",
		"a/b/c/d",
	}

	for _, image := range validImages {
		assert.True(t, isValidImageName(image), "Expected %s to be valid", image)
	}

	for _, image := range invalidImages {
		assert.False(t, isValidImageName(image), "Expected %s to be invalid", image)
	}

	// Test isValidVolumeSize
	validSizes := []string{"1Ki", "500Mi", "10Gi", "1Ti"}
	invalidSizes := []string{"1K", "500M", "10G", "1T", "1.5Gi", "-1Gi"}

	for _, size := range validSizes {
		assert.True(t, isValidVolumeSize(size), "Expected %s to be valid", size)
	}

	for _, size := range invalidSizes {
		assert.False(t, isValidVolumeSize(size), "Expected %s to be invalid", size)
	}
}
