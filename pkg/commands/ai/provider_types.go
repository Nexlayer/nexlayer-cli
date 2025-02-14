// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style license
// that can be found in the LICENSE file.

package ai

import (
	"context"
)

// Capability represents an AI provider's capabilities.
type Capability int

const (
	CapCodeGeneration Capability = 1 << iota
	CapCodeCompletion
	CapDeploymentAssistance
	CapErrorDiagnosis
)

// AIProvider represents an AI code assistant provider.
type AIProvider struct {
	Name         string
	Description  string
	EnvVarKey    string
	Capabilities Capability
}

// GenerateText generates text using the AI provider.
func (p *AIProvider) GenerateText(ctx context.Context, prompt string) (string, error) {
	// TODO: Implement provider-specific text generation.
	// For now, return a basic template matching the new Nexlayer schema.
	return `application:
  name: "generated-app"
  pods:
    - name: "app"
      image: "nginx:latest"
      servicePorts: [80]`, nil
}
