// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package detectors

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/types"
	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
)

// SupabaseDetector detects Supabase integration in projects
type SupabaseDetector struct {
	detection.BaseDetector
}

// NewSupabaseDetector creates a new detector for Supabase
func NewSupabaseDetector() *SupabaseDetector {
	detector := &SupabaseDetector{}
	detector.BaseDetector = *detection.NewBaseDetector("Supabase Detector", 0.8)
	return detector
}

// Priority returns the detector's priority level
func (d *SupabaseDetector) Priority() int {
	return 90
}

// Detect analyzes a directory to detect Supabase usage
func (d *SupabaseDetector) Detect(ctx context.Context, dir string) (*detection.ProjectInfo, error) {
	// Create a basic project info
	projectInfo := &detection.ProjectInfo{
		Type:       "unknown",
		Path:       dir,
		Confidence: 0.0,
		Metadata:   make(map[string]interface{}),
	}

	// Check package.json for Supabase dependencies
	packageJSONPath := filepath.Join(dir, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		packageJSONContent, err := os.ReadFile(packageJSONPath)
		if err == nil {
			contentStr := string(packageJSONContent)
			if strings.Contains(contentStr, "\"@supabase/supabase-js\"") ||
				strings.Contains(contentStr, "\"@supabase/auth-helpers\"") ||
				strings.Contains(contentStr, "\"@supabase/auth-ui-react\"") {
				projectInfo.Type = "supabase"
				projectInfo.Confidence = 0.9
				projectInfo.Language = "JavaScript"
				return projectInfo, nil
			}
		}
	}

	// Check for Supabase config files
	supabaseConfigPath := filepath.Join(dir, "supabase")
	if _, err := os.Stat(supabaseConfigPath); err == nil {
		// Check for Supabase configuration files
		configFilePath := filepath.Join(supabaseConfigPath, "config.toml")
		if _, err := os.Stat(configFilePath); err == nil {
			projectInfo.Type = "supabase"
			projectInfo.Confidence = 0.9
			projectInfo.Metadata["has_supabase_config"] = true
			return projectInfo, nil
		}
	}

	// Check for .env files containing Supabase URLs/keys
	envFiles := []string{".env", ".env.local", ".env.development", ".env.production"}
	supabaseEnvPattern := regexp.MustCompile(`(?i)(SUPABASE_URL|SUPABASE_KEY|NEXT_PUBLIC_SUPABASE)`)

	for _, envFile := range envFiles {
		envPath := filepath.Join(dir, envFile)
		if _, err := os.Stat(envPath); err == nil {
			content, err := os.ReadFile(envPath)
			if err == nil {
				if supabaseEnvPattern.MatchString(string(content)) {
					projectInfo.Type = "supabase"
					projectInfo.Confidence = 0.8
					projectInfo.Metadata["has_supabase_env"] = true
					return projectInfo, nil
				}
			}
		}
	}

	// Check for Supabase imports in JavaScript/TypeScript files
	supabaseImportPattern := regexp.MustCompile(`import\s+.*?from\s+['"]@supabase`)
	
	// Check JavaScript/TypeScript files
	jsExtensions := []string{".js", ".jsx", ".ts", ".tsx"}
	for _, ext := range jsExtensions {
		matches, _ := filepath.Glob(filepath.Join(dir, "**/*"+ext))
		for _, file := range matches {
			content, err := os.ReadFile(file)
			if err == nil {
				contentStr := string(content)
				
				// Check for Supabase imports
				if supabaseImportPattern.MatchString(contentStr) {
					projectInfo.Type = "supabase"
					projectInfo.Confidence = 0.8
					projectInfo.Language = strings.HasPrefix(ext, ".ts") ? "TypeScript" : "JavaScript"
					return projectInfo, nil
				}
				
				// Check for common Supabase client usage patterns
				if strings.Contains(contentStr, "createClient") && strings.Contains(contentStr, "supabase") {
					projectInfo.Type = "supabase"
					projectInfo.Confidence = 0.7
					projectInfo.Language = strings.HasPrefix(ext, ".ts") ? "TypeScript" : "JavaScript"
					return projectInfo, nil
				}
			}
		}
	}

	return projectInfo, nil
}

// processJSProject analyzes JavaScript/TypeScript projects for Supabase usage
func (d *SupabaseDetector) processJSProject(dir string, packageJSONPath string) (*types.ProjectInfo, bool, error) {
	projectInfo := &types.ProjectInfo{
		Type:         types.TypeUnknown,
		Dependencies: make(map[string]string),
		Scripts:      make(map[string]string),
	}

	// Read and parse package.json
	packageJSONBytes, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return projectInfo, false, err
	}

	var packageJSON map[string]interface{}
	if err := json.Unmarshal(packageJSONBytes, &packageJSON); err != nil {
		return projectInfo, false, err
	}

	// Extract project name and version
	if name, ok := packageJSON["name"].(string); ok {
		projectInfo.Name = name
	}
	if version, ok := packageJSON["version"].(string); ok {
		projectInfo.Version = version
	}

	// Check for dependencies
	dependencies := make(map[string]interface{})
	if deps, ok := packageJSON["dependencies"].(map[string]interface{}); ok {
		dependencies = deps
	}
	if devDeps, ok := packageJSON["devDependencies"].(map[string]interface{}); ok {
		for k, v := range devDeps {
			dependencies[k] = v
		}
	}

	// Extract all dependencies
	hasNextjs := false
	hasOpenAI := false
	hasLangChain := false
	hasSupabase := false

	for name, version := range dependencies {
		if vStr, ok := version.(string); ok {
			projectInfo.Dependencies[name] = vStr

			// Check for key dependencies
			if name == "@supabase/supabase-js" || name == "@supabase/auth-helpers-nextjs" {
				hasSupabase = true
			} else if name == "next" {
				hasNextjs = true
			} else if name == "openai" {
				hasOpenAI = true
			} else if name == "langchain" || name == "@langchain/core" {
				hasLangChain = true
			}
		}
	}

	// Check for scripts
	if scripts, ok := packageJSON["scripts"].(map[string]interface{}); ok {
		for name, cmd := range scripts {
			if cmdStr, ok := cmd.(string); ok {
				projectInfo.Scripts[name] = cmdStr
			}
		}
	}

	// Check for Docker configuration
	if _, err := os.Stat(filepath.Join(dir, "Dockerfile")); err == nil {
		projectInfo.HasDocker = true
	}
	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yml")); err == nil {
		projectInfo.HasDocker = true
	}
	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yaml")); err == nil {
		projectInfo.HasDocker = true
	}

	// Set default port based on framework
	if hasNextjs {
		projectInfo.Port = 3000
	} else {
		projectInfo.Port = 5173 // Default for Vite-based projects
	}

	// Check for .env files with Supabase config
	supabaseEnvDetected := d.checkSupabaseEnvFiles(dir)

	// Set Supabase as detected if we found dependencies or env vars
	hasSupabase = hasSupabase || supabaseEnvDetected

	// If we have Supabase and specific other technologies, set the appropriate type
	if hasSupabase && hasNextjs {
		if hasLangChain {
			projectInfo.Type = types.TypeNextjsSupabaseLangchain
			projectInfo.LLMProvider = "LangChain"
		} else if hasOpenAI {
			projectInfo.Type = types.TypeNextjsSupabaseOpenAI
			projectInfo.LLMProvider = "OpenAI"
		} else {
			projectInfo.Type = types.TypeNextjsSupabase
		}
	} else if hasSupabase {
		projectInfo.Type = types.TypeSupabase
	}

	return projectInfo, hasSupabase, nil
}

// processPythonProject analyzes Python projects for Supabase usage
func (d *SupabaseDetector) processPythonProject(dir string, requirementsPath string) (*types.ProjectInfo, bool, error) {
	projectInfo := &types.ProjectInfo{
		Type:         types.TypeUnknown,
		Dependencies: make(map[string]string),
		Scripts:      make(map[string]string),
		Name:         filepath.Base(dir), // Default to directory name for Python projects
	}

	// Read requirements.txt
	requirementsData, err := os.ReadFile(requirementsPath)
	if err != nil {
		return projectInfo, false, err
	}

	// Parse requirements
	lines := strings.Split(string(requirementsData), "\n")
	hasDjango := false
	hasFlask := false
	hasSupabase := false
	hasOpenAI := false
	hasLangChain := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		var packageName, version string
		if strings.Contains(line, "==") {
			parts := strings.Split(line, "==")
			packageName, version = parts[0], parts[1]
		} else if strings.Contains(line, ">=") {
			parts := strings.Split(line, ">=")
			packageName, version = parts[0], ">="+parts[1]
		} else {
			packageName, version = line, ""
		}

		// Store dependency
		packageName = strings.TrimSpace(packageName)
		projectInfo.Dependencies[packageName] = version

		// Check for key packages
		lcPackageName := strings.ToLower(packageName)
		if lcPackageName == "supabase" || lcPackageName == "postgrest-py" {
			hasSupabase = true
		} else if lcPackageName == "django" {
			hasDjango = true
		} else if lcPackageName == "flask" {
			hasFlask = true
		} else if lcPackageName == "openai" {
			hasOpenAI = true
		} else if lcPackageName == "langchain" {
			hasLangChain = true
		}
	}

	// Check for Docker configuration
	if _, err := os.Stat(filepath.Join(dir, "Dockerfile")); err == nil {
		projectInfo.HasDocker = true
	}
	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yml")); err == nil {
		projectInfo.HasDocker = true
	}
	if _, err := os.Stat(filepath.Join(dir, "docker-compose.yaml")); err == nil {
		projectInfo.HasDocker = true
	}

	// Check for .env files with Supabase config
	supabaseEnvDetected := d.checkSupabaseEnvFiles(dir)

	// Set Supabase as detected if we found dependencies or env vars
	hasSupabase = hasSupabase || supabaseEnvDetected

	// Set default port based on framework
	if hasDjango {
		projectInfo.Port = 8000
		if hasSupabase {
			if hasOpenAI {
				projectInfo.Type = types.TypeDjangoSupabaseOpenAI
				projectInfo.LLMProvider = "OpenAI"
			} else if hasLangChain {
				projectInfo.Type = types.TypeDjangoSupabaseLangchain
				projectInfo.LLMProvider = "LangChain"
			} else {
				projectInfo.Type = types.TypeDjangoSupabase
			}
		}
	} else if hasFlask {
		projectInfo.Port = 5000
		if hasSupabase {
			if hasOpenAI {
				projectInfo.Type = types.TypeFlaskSupabaseOpenAI
				projectInfo.LLMProvider = "OpenAI"
			} else if hasLangChain {
				projectInfo.Type = types.TypeFlaskSupabaseLangchain
				projectInfo.LLMProvider = "LangChain"
			} else {
				projectInfo.Type = types.TypeFlaskSupabase
			}
		}
	} else if hasSupabase {
		// Generic Python with Supabase
		projectInfo.Type = types.TypePythonSupabase
		projectInfo.Port = 8000
	}

	return projectInfo, hasSupabase, nil
}

// determineSupabaseProjectType sets the appropriate project type based on other detected frameworks
func (d *SupabaseDetector) determineSupabaseProjectType(dir string, projectInfo *types.ProjectInfo) {
	// Already handled in the process methods, but can be extended here
	// with additional file checks if needed
}

// detectPgVector checks for pgvector usage in the project
func (d *SupabaseDetector) detectPgVector(dir string, projectInfo *types.ProjectInfo) {
	// Look for pgvector in SQL migrations or setup files
	sqlMigrationsDir := filepath.Join(dir, "supabase", "migrations")
	if sqlDirExists, _ := exists(sqlMigrationsDir); sqlDirExists {
		files, err := os.ReadDir(sqlMigrationsDir)
		if err == nil {
			for _, file := range files {
				if strings.HasSuffix(file.Name(), ".sql") {
					filePath := filepath.Join(sqlMigrationsDir, file.Name())
					content, err := os.ReadFile(filePath)
					if err == nil && strings.Contains(string(content), "pgvector") {
						projectInfo.HasVectorDB = true
						break
					}
				}
			}
		}
	}

	// Check for pgvector references in TS/JS files
	tsFiles, _ := d.findFiles(dir, []string{".ts", ".js"}, []string{"node_modules", ".next", "dist"})
	for _, file := range tsFiles {
		content, err := os.ReadFile(file)
		if err == nil {
			if strings.Contains(string(content), "pgvector") ||
				strings.Contains(string(content), "vector(") ||
				strings.Contains(string(content), "createVectorStore") {
				projectInfo.HasVectorDB = true
				break
			}
		}
	}
}

// checkSupabaseEnvFiles checks environment files for Supabase configuration
func (d *SupabaseDetector) checkSupabaseEnvFiles(dir string) bool {
	envFiles := []string{
		".env",
		".env.local",
		".env.development",
		".env.production",
	}

	for _, envFile := range envFiles {
		envPath := filepath.Join(dir, envFile)
		if fileExists, _ := exists(envPath); fileExists {
			content, err := os.ReadFile(envPath)
			if err == nil {
				// Check for Supabase URL or API keys
				if strings.Contains(string(content), "SUPABASE_URL") ||
					strings.Contains(string(content), "SUPABASE_KEY") ||
					strings.Contains(string(content), "SUPABASE_SERVICE_KEY") {
					return true
				}
			}
		}
	}

	return false
}

// findFiles recursively finds files with specific extensions
func (d *SupabaseDetector) findFiles(root string, extensions []string, excludeDirs []string) ([]string, error) {
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip excluded directories
		if info.IsDir() {
			for _, excludeDir := range excludeDirs {
				if info.Name() == excludeDir {
					return filepath.SkipDir
				}
			}
			return nil
		}

		// Check file extensions
		for _, ext := range extensions {
			if strings.HasSuffix(path, ext) {
				files = append(files, path)
				break
			}
		}

		return nil
	})

	return files, err
}

// exists checks if a file or directory exists
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// TODO: Integrate this detector into the detector registry in pkg/detection/detectors.go
// by adding it to the detectors slice in the NewDetectorRegistry function.
