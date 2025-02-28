// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package detection

import (
	"context"
	"fmt"
)

// ProjectInfo contains detected information about a project
type ProjectInfo struct {
	Type         string                 // Detected project type
	Name         string                 // Project name
	Version      string                 // Project version
	Path         string                 // Project path
	Framework    string                 // Framework used in the project
	Language     string                 // Primary language
	Dependencies map[string]string      // Project dependencies
	Scripts      map[string]string      // Build/run scripts
	Port         int                    // Default port
	HasDocker    bool                   // Whether the project has Docker setup
	Editor       string                 // Detected editor
	LLMProvider  string                 // AI/LLM provider if using an AI-powered IDE
	LLMModel     string                 // LLM model being used
	ImageTag     string                 // Docker image tag
	Confidence   float64                // Detection confidence (0-1)
	Metadata     map[string]interface{} // Additional metadata from detection
}

// Detector defines the interface for project detectors
type Detector interface {
	// Detect analyzes a directory to detect project information
	Detect(ctx context.Context, dir string) (*ProjectInfo, error)
	// Name returns the detector name
	Name() string
}

// BaseDetector provides common functionality for all detectors
type BaseDetector struct {
	name       string
	confidence float64
}

// Name returns the detector name
func (d *BaseDetector) Name() string {
	return d.name
}

// NewBaseDetector creates a new base detector with the given name and confidence
func NewBaseDetector(name string, confidence float64) *BaseDetector {
	return &BaseDetector{
		name:       name,
		confidence: confidence,
	}
}

// NewVSCodeDetector creates a new detector for VS Code
func NewVSCodeDetector() Detector {
	return &BaseDetector{
		name:       "vscode",
		confidence: 0.9,
	}
}

// NewJetBrainsDetector creates a new detector for JetBrains IDEs
func NewJetBrainsDetector() Detector {
	return &BaseDetector{
		name:       "jetbrains",
		confidence: 0.9,
	}
}

// NewLLMDetector creates a new detector for AI/LLM-powered IDEs
func NewLLMDetector() Detector {
	return &BaseDetector{
		name:       "llm",
		confidence: 0.9,
	}
}

// NewGoDetector creates a new detector for Go projects
func NewGoDetector() Detector {
	return &BaseDetector{
		name:       "go",
		confidence: 0.8,
	}
}

// NewNodeJSDetector creates a new detector for Node.js projects
func NewNodeJSDetector() Detector {
	return &BaseDetector{
		name:       "nodejs",
		confidence: 0.8,
	}
}

// NewPythonDetector creates a new detector for Python projects
func NewPythonDetector() Detector {
	return &BaseDetector{
		name:       "python",
		confidence: 0.8,
	}
}

// NewDockerDetector creates a new detector for Docker projects
func NewDockerDetector() Detector {
	return &BaseDetector{
		name:       "docker",
		confidence: 0.9,
	}
}

// Detect implements the Detector interface for the BaseDetector
// This is a base implementation that should be overridden by specific detectors
func (d *BaseDetector) Detect(ctx context.Context, dir string) (*ProjectInfo, error) {
	// Base implementation just returns a basic ProjectInfo with the detector name
	return &ProjectInfo{
		Type:       d.name,
		Path:       dir,
		Confidence: d.confidence,
		Metadata:   make(map[string]interface{}),
	}, nil
}

// EmitDeprecationWarning logs a deprecation warning for a detector
// This function should be called at the beginning of the Detect method in deprecated detectors
func EmitDeprecationWarning(detectorName string) {
	// Use a log level that's visible but not alarming
	fmt.Printf("DEPRECATED: %s is deprecated and will be removed in a future version. "+
		"Please migrate to the unified StackDetector. "+
		"See pkg/detection/MIGRATION_GUIDE.md for details.\n", detectorName)
}
