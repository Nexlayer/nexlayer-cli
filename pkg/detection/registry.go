// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package detection

import (
	"context"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
)

// Registry manages and orchestrates project detectors
type Registry struct {
	mu        sync.RWMutex
	detectors []Detector
}

// NewRegistry creates a new detector registry with default detectors
func NewRegistry() *Registry {
	r := &Registry{
		detectors: make([]Detector, 0),
	}

	// Register default detectors in priority order
	r.Register(
		&LLMDetector{},  // AI/LLM detection
		&MERNDetector{}, // Full-stack detectors
		&PERNDetector{},
		&MEANDetector{},
		&NextjsDetector{}, // Framework detectors
		&ReactDetector{},
		&NodeDetector{},
		&PythonDetector{},
		&GoDetector{},
		&DockerDetector{}, // Infrastructure detectors
	)

	return r
}

// Register adds one or more detectors to the registry
func (r *Registry) Register(detectors ...Detector) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.detectors = append(r.detectors, detectors...)
	r.sortDetectors()
}

// sortDetectors sorts detectors by priority (highest first)
func (r *Registry) sortDetectors() {
	sort.Slice(r.detectors, func(i, j int) bool {
		return r.detectors[i].Priority() > r.detectors[j].Priority()
	})
}

// DetectProject attempts to detect project information using all registered detectors
func (r *Registry) DetectProject(ctx context.Context, dir string) (*types.ProjectInfo, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Try each detector in priority order
	for _, detector := range r.detectors {
		select {
		case <-ctx.Done():
			return nil, NewDetectionError(ErrorTypeInternal, "detection cancelled", ctx.Err())
		default:
			if info, err := detector.Detect(dir); err == nil && info != nil {
				return info, nil
			}
		}
	}

	// If no detector succeeded, return unknown type
	return &types.ProjectInfo{
		Type: types.TypeUnknown,
		Name: getProjectName(dir),
	}, nil
}

// getProjectName extracts a clean project name from a directory path
func getProjectName(dir string) string {
	name := filepath.Base(dir)
	return sanitizeName(name)
}

// sanitizeName ensures the name follows Nexlayer naming conventions
func sanitizeName(name string) string {
	// Convert to lowercase
	name = strings.ToLower(name)

	// Replace invalid characters with hyphens
	name = strings.Map(func(r rune) rune {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			return r
		}
		return '-'
	}, name)

	// Ensure starts with a letter
	if len(name) > 0 && (name[0] < 'a' || name[0] > 'z') {
		name = "app-" + name
	}

	// If empty after sanitization, use default
	if name == "" {
		name = "app"
	}

	return name
}

// GetDetectors returns a copy of the registered detectors
func (r *Registry) GetDetectors() []Detector {
	r.mu.RLock()
	defer r.mu.RUnlock()

	detectors := make([]Detector, len(r.detectors))
	copy(detectors, r.detectors)
	return detectors
}

// Clear removes all registered detectors
func (r *Registry) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.detectors = make([]Detector, 0)
}
