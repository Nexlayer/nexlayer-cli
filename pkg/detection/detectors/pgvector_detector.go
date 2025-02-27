// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package detectors provides project-specific detection implementations.
package detectors

// DEPRECATED: This specialized detector is deprecated and will be removed in a future version.
// The functionality has been replaced by the unified StackDetector in pkg/detection/stack_detector.go,
// which uses pattern-based detection to identify technology components including pgvector.
// Please migrate any direct references to this detector to use StackDetector instead.
// Planned removal: v1.x.0 (next major/minor release)

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
)

// PgvectorDetector detects the presence of pgvector in a project
type PgvectorDetector struct {
	detection.BaseDetector
}

// NewPgvectorDetector creates a new detector for pgvector
func NewPgvectorDetector() *PgvectorDetector {
	detector := &PgvectorDetector{}
	detector.BaseDetector = *detection.NewBaseDetector("Pgvector Detector", 0.8)
	return detector
}

// detectPgvectorJS detects pgvector usage in JavaScript/TypeScript projects
func (d *PgvectorDetector) detectPgvectorJS(projectPath string) (bool, float64) {
	// Check package.json for pgvector related dependencies
	packageJSONPath := filepath.Join(projectPath, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		packageJSONContent, err := os.ReadFile(packageJSONPath)
		if err == nil {
			contentStr := string(packageJSONContent)
			if strings.Contains(contentStr, "\"pgvector\"") ||
				strings.Contains(contentStr, "\"@pgvector/pgvector\"") ||
				strings.Contains(contentStr, "\"pg-vector\"") {
				return true, 0.9
			}

			// Check for Prisma with pgvector
			if strings.Contains(contentStr, "\"prisma\"") && strings.Contains(contentStr, "\"@prisma/client\"") {
				// Let's also look for Prisma schema
				prismaSchemaPath := filepath.Join(projectPath, "prisma/schema.prisma")
				if _, err := os.Stat(prismaSchemaPath); err == nil {
					prismaContent, err := os.ReadFile(prismaSchemaPath)
					if err == nil && strings.Contains(string(prismaContent), "vector") {
						return true, 0.8
					}
				}
			}
		}
	}

	// Check for pgvector SQL statements in migration files
	migrationPattern := regexp.MustCompile(`(?i)create\s+extension\s+(?:if\s+not\s+exists\s+)?vector`)
	vectorColumnPattern := regexp.MustCompile(`(?i)column\s+\w+\s+vector\s*\(`)

	// Check SQL files for pgvector extension or vector columns
	sqlMatches, _ := filepath.Glob(filepath.Join(projectPath, "**/*.sql"))
	for _, file := range sqlMatches {
		content, err := os.ReadFile(file)
		if err == nil {
			contentStr := string(content)
			if migrationPattern.MatchString(contentStr) || vectorColumnPattern.MatchString(contentStr) {
				return true, 0.9
			}
		}
	}

	// Check JavaScript/TypeScript files for pgvector usage
	pgvectorImportPattern := regexp.MustCompile(`import\s+.*?from\s+['"]pgvector`)

	jsExtensions := []string{".js", ".jsx", ".ts", ".tsx"}
	for _, ext := range jsExtensions {
		matches, _ := filepath.Glob(filepath.Join(projectPath, "**/*"+ext))
		for _, file := range matches {
			content, err := os.ReadFile(file)
			if err == nil {
				contentStr := string(content)
				// Check for pgvector imports
				if pgvectorImportPattern.MatchString(contentStr) {
					return true, 0.9
				}

				// Check for common pgvector method calls
				if strings.Contains(contentStr, "createVector") ||
					strings.Contains(contentStr, "knexCosmetic.raw('create extension vector')") ||
					strings.Contains(contentStr, "cosineDistance") {
					return true, 0.8
				}
			}
		}
	}

	// Check for Supabase configuration with pgvector
	supabaseConfigPath := filepath.Join(projectPath, "supabase/migrations")
	if _, err := os.Stat(supabaseConfigPath); err == nil {
		migrationsMatches, _ := filepath.Glob(filepath.Join(supabaseConfigPath, "**/*.sql"))
		for _, file := range migrationsMatches {
			content, err := os.ReadFile(file)
			if err == nil {
				contentStr := string(content)
				if migrationPattern.MatchString(contentStr) || vectorColumnPattern.MatchString(contentStr) {
					return true, 0.9
				}
			}
		}
	}

	return false, 0.0
}

// detectPgvectorPython detects pgvector usage in Python projects
func (d *PgvectorDetector) detectPgvectorPython(projectPath string) (bool, float64) {
	// Check requirements.txt for pgvector
	requirementsPath := filepath.Join(projectPath, "requirements.txt")
	if _, err := os.Stat(requirementsPath); err == nil {
		content, err := os.ReadFile(requirementsPath)
		if err == nil {
			contentStr := string(content)
			if strings.Contains(contentStr, "pgvector") ||
				strings.Contains(contentStr, "sqlalchemy-pgvector") {
				return true, 0.9
			}
		}
	}

	// Check Poetry pyproject.toml
	pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
	if _, err := os.Stat(pyprojectPath); err == nil {
		content, err := os.ReadFile(pyprojectPath)
		if err == nil {
			contentStr := string(content)
			if strings.Contains(contentStr, "pgvector") ||
				strings.Contains(contentStr, "sqlalchemy-pgvector") {
				return true, 0.9
			}
		}
	}

	// Check for import statements in Python files
	importPattern := regexp.MustCompile(`(?m)^import\s+pgvector`)
	fromImportPattern := regexp.MustCompile(`(?m)^from\s+pgvector\s+import`)
	sqlalchemyPattern := regexp.MustCompile(`(?m)^from\s+sqlalchemy_pgvector\s+import`)

	// Check Python files
	pyMatches, _ := filepath.Glob(filepath.Join(projectPath, "**/*.py"))
	for _, file := range pyMatches {
		content, err := os.ReadFile(file)
		if err == nil {
			contentStr := string(content)
			if importPattern.MatchString(contentStr) ||
				fromImportPattern.MatchString(contentStr) ||
				sqlalchemyPattern.MatchString(contentStr) {
				return true, 0.9
			}

			// Check for pgvector related code
			if strings.Contains(contentStr, "Vector(") ||
				strings.Contains(contentStr, "cosine_distance") ||
				strings.Contains(contentStr, "CREATE EXTENSION vector") {
				return true, 0.8
			}
		}
	}

	return false, 0.0
}

// Detect implementation for PgvectorDetector
func (d *PgvectorDetector) Detect(ctx context.Context, dir string) (*detection.ProjectInfo, error) {
	// Emit deprecation warning
	detection.EmitDeprecationWarning("PgvectorDetector")

	// Create a basic project info
	projectInfo := &detection.ProjectInfo{
		Type:       "unknown",
		Path:       dir,
		Confidence: 0.0,
		Metadata:   make(map[string]interface{}),
	}

	// Try to detect pgvector in JavaScript/TypeScript
	jsFound, jsConf := d.detectPgvectorJS(dir)
	if jsFound {
		projectInfo.Type = "pgvector"
		projectInfo.Confidence = jsConf
		projectInfo.Language = "JavaScript"
		return projectInfo, nil
	}

	// Try to detect pgvector in Python
	pyFound, pyConf := d.detectPgvectorPython(dir)
	if pyFound {
		projectInfo.Type = "pgvector"
		projectInfo.Confidence = pyConf
		projectInfo.Language = "Python"
		return projectInfo, nil
	}

	return projectInfo, nil
}
