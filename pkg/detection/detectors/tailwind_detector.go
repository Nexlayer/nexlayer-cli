// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package detectors provides project-specific detection implementations.

// DEPRECATED: This specialized detector is deprecated and will be removed in a future version.
// The functionality has been replaced by the unified StackDetector in pkg/detection/stack_detector.go,
// which uses pattern-based detection to identify technology components including Tailwind CSS.
// Please migrate any direct references to this detector to use StackDetector instead.
// Planned removal: v1.x.0 (next major/minor release)

package detectors

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
)

// TailwindDetector detects the presence of Tailwind CSS in a project
type TailwindDetector struct {
	detection.BaseDetector
}

// NewTailwindDetector creates a new detector for Tailwind CSS
func NewTailwindDetector() *TailwindDetector {
	detector := &TailwindDetector{}
	detector.BaseDetector = *detection.NewBaseDetector("Tailwind CSS Detector", 0.8)
	return detector
}

// Detect implementation for TailwindDetector
func (d *TailwindDetector) Detect(ctx context.Context, dir string) (*detection.ProjectInfo, error) {
	// Emit deprecation warning
	detection.EmitDeprecationWarning("TailwindDetector")

	// Create a basic project info
	projectInfo := &detection.ProjectInfo{
		Type:       "unknown",
		Path:       dir,
		Confidence: 0.0,
		Metadata:   make(map[string]interface{}),
	}

	// Check package.json for Tailwind dependencies
	packageJSONPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		packageJSONContent, err := os.ReadFile(packageJSONPath)
		if err == nil {
			contentStr := string(packageJSONContent)
			if strings.Contains(contentStr, "\"tailwindcss\"") {
				projectInfo.Type = "tailwind"
				projectInfo.Confidence = 0.9
				projectInfo.Language = "JavaScript"
				return projectInfo, nil
			}
		}
	}

	// Check for tailwind.config.js
	tailwindConfigPath := filepath.Join(dir, "tailwind.config.js")
	tailwindConfigCJSPath := filepath.Join(dir, "tailwind.config.cjs")
	tailwindConfigTSPath := filepath.Join(dir, "tailwind.config.ts")

	if _, err := os.Stat(tailwindConfigPath); err == nil {
		projectInfo.Type = "tailwind"
		projectInfo.Confidence = 0.9
		projectInfo.Language = "JavaScript"
		return projectInfo, nil
	}

	if _, err := os.Stat(tailwindConfigCJSPath); err == nil {
		projectInfo.Type = "tailwind"
		projectInfo.Confidence = 0.9
		projectInfo.Language = "JavaScript"
		return projectInfo, nil
	}

	if _, err := os.Stat(tailwindConfigTSPath); err == nil {
		projectInfo.Type = "tailwind"
		projectInfo.Confidence = 0.9
		projectInfo.Language = "TypeScript"
		return projectInfo, nil
	}

	// Check PostCSS config for Tailwind
	postcssConfigPath := filepath.Join(dir, "postcss.config.js")
	if _, err := os.Stat(postcssConfigPath); err == nil {
		content, err := os.ReadFile(postcssConfigPath)
		if err == nil && strings.Contains(string(content), "tailwindcss") {
			projectInfo.Type = "tailwind"
			projectInfo.Confidence = 0.9
			projectInfo.Language = "JavaScript"
			return projectInfo, nil
		}
	}

	// Check for Tailwind directives in CSS files
	tailwindDirectivePattern := regexp.MustCompile(`@tailwind\s+(?:base|components|utilities)`)

	cssMatches, _ := filepath.Glob(filepath.Join(dir, "**/*.css"))
	for _, file := range cssMatches {
		content, err := os.ReadFile(file)
		if err == nil {
			if tailwindDirectivePattern.MatchString(string(content)) {
				projectInfo.Type = "tailwind"
				projectInfo.Confidence = 0.9
				return projectInfo, nil
			}
		}
	}

	// Check for Tailwind class usage in HTML/JSX/TSX files
	tailwindClassPattern := regexp.MustCompile(`class(?:Name)?=["'](?:[^"']*\s)?(?:bg-|text-|flex|grid|p-|m-|rounded|shadow|hover:|focus:|sm:|md:|lg:|xl:)`)

	// Extensions to check for Tailwind class usage
	extensions := []string{".html", ".js", ".jsx", ".ts", ".tsx", ".vue", ".svelte"}

	for _, ext := range extensions {
		matches, _ := filepath.Glob(filepath.Join(dir, "**/*"+ext))
		for _, file := range matches {
			content, err := os.ReadFile(file)
			if err == nil {
				if tailwindClassPattern.MatchString(string(content)) {
					projectInfo.Type = "tailwind"
					// Lower confidence since these could be from other libraries too
					projectInfo.Confidence = 0.7
					return projectInfo, nil
				}
			}
		}
	}

	return projectInfo, nil
}
