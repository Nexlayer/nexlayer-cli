// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package detectors provides project-specific detection implementations.
// Note: Some detectors in this package are being gradually replaced by the unified StackDetector
// in pkg/detection/stack_detector.go, but remain available for backward compatibility.
package detectors

// NOTE: This detector's functionality is also available in the unified StackDetector,
// but this implementation is maintained for backward compatibility and specific use cases.
// Consider using StackDetector for new code that needs to detect Next.js projects.

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
)

// NextjsDetector detects Next.js framework usage in a project
type NextjsDetector struct {
	detection.BaseDetector
}

// NewNextjsDetector creates a new detector for Next.js
func NewNextjsDetector() *NextjsDetector {
	detector := &NextjsDetector{}
	detector.BaseDetector = *detection.NewBaseDetector("Next.js Detector", 0.9)
	return detector
}

// Detect implementation for NextjsDetector
func (d *NextjsDetector) Detect(ctx context.Context, dir string) (*detection.ProjectInfo, error) {
	// Emit deprecation warning
	detection.EmitDeprecationWarning("NextjsDetector")

	// Create a basic project info
	projectInfo := &detection.ProjectInfo{
		Type:       "unknown",
		Path:       dir,
		Confidence: 0.0,
		Metadata:   make(map[string]interface{}),
	}

	// Check package.json for Next.js dependency
	packageJSONPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		packageJSONContent, err := os.ReadFile(packageJSONPath)
		if err == nil {
			if strings.Contains(string(packageJSONContent), "\"next\"") {
				projectInfo.Type = "nextjs"
				projectInfo.Confidence = 0.9
				projectInfo.Language = "JavaScript"
				projectInfo.Framework = "Next.js"
				return projectInfo, nil
			}
		}
	}

	// Check for next.config.js
	nextConfigPath := filepath.Join(dir, "next.config.js")
	if _, err := os.Stat(nextConfigPath); err == nil {
		projectInfo.Type = "nextjs"
		projectInfo.Confidence = 0.9
		projectInfo.Language = "JavaScript"
		projectInfo.Framework = "Next.js"
		return projectInfo, nil
	}

	// Check for Next.js app directory structure
	appDirPath := filepath.Join(dir, "app")
	if _, err := os.Stat(appDirPath); err == nil {
		// Check for Next.js-specific files in app directory
		layoutPath := filepath.Join(appDirPath, "layout.js")
		layoutTsPath := filepath.Join(appDirPath, "layout.tsx")
		if _, err := os.Stat(layoutPath); err == nil {
			projectInfo.Type = "nextjs"
			projectInfo.Confidence = 0.9
			projectInfo.Language = "JavaScript"
			projectInfo.Framework = "Next.js"
			projectInfo.Metadata["app_dir"] = true
			return projectInfo, nil
		}
		if _, err := os.Stat(layoutTsPath); err == nil {
			projectInfo.Type = "nextjs"
			projectInfo.Confidence = 0.9
			projectInfo.Language = "TypeScript"
			projectInfo.Framework = "Next.js"
			projectInfo.Metadata["app_dir"] = true
			return projectInfo, nil
		}
	}

	// Check for Next.js pages directory structure
	pagesDirPath := filepath.Join(dir, "pages")
	if _, err := os.Stat(pagesDirPath); err == nil {
		// Check for Next.js-specific files in pages directory
		indexPath := filepath.Join(pagesDirPath, "index.js")
		indexTsPath := filepath.Join(pagesDirPath, "index.tsx")
		apiDirPath := filepath.Join(pagesDirPath, "api")

		// Check if any of the Next.js specific files exist
		indexJsInfo, _ := os.Stat(indexPath)
		indexTsInfo, _ := os.Stat(indexTsPath)
		apiDirInfo, _ := os.Stat(apiDirPath)

		if indexJsInfo != nil || indexTsInfo != nil || apiDirInfo != nil {
			projectInfo.Type = "nextjs"
			projectInfo.Confidence = 0.8
			projectInfo.Language = "JavaScript"
			projectInfo.Framework = "Next.js"
			projectInfo.Metadata["pages_dir"] = true
			return projectInfo, nil
		}
	}

	// Check for common Next.js imports in JavaScript/TypeScript files
	nextjsImportPattern := regexp.MustCompile(`import\s+.*?from\s+['"]next`)

	// Check JavaScript/TypeScript files
	jsExtensions := []string{".js", ".jsx", ".ts", ".tsx"}
	for _, ext := range jsExtensions {
		matches, _ := filepath.Glob(filepath.Join(dir, "**/*"+ext))
		for _, file := range matches {
			content, err := os.ReadFile(file)
			if err == nil {
				if nextjsImportPattern.MatchString(string(content)) {
					projectInfo.Type = "nextjs"
					projectInfo.Confidence = 0.7
					if strings.HasPrefix(ext, ".ts") {
						projectInfo.Language = "TypeScript"
					} else {
						projectInfo.Language = "JavaScript"
					}
					projectInfo.Framework = "Next.js"
					return projectInfo, nil
				}
			}
		}
	}

	return projectInfo, nil
}
