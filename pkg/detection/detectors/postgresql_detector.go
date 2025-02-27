package detectors

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
)

// PostgreSQLDetector detects PostgreSQL database usage in a project
type PostgreSQLDetector struct {
	detection.BaseDetector
}

// NewPostgreSQLDetector creates a new detector for PostgreSQL database integration
func NewPostgreSQLDetector() *PostgreSQLDetector {
	detector := &PostgreSQLDetector{}
	detector.BaseDetector = *detection.NewBaseDetector("PostgreSQL Detector", 0.8)
	return detector
}

// detectPostgresJS detects PostgreSQL usage in JavaScript/TypeScript projects
func (d *PostgreSQLDetector) detectPostgresJS(projectPath string) (bool, float64) {
	// Define regex patterns
	postgresConnectionPattern := regexp.MustCompile(`(?i)postgres(?:ql)?:\/\/[^\/\s]+\/\w+`)
	postgresImportPattern := regexp.MustCompile(`(?i)(?:import|require).*(?:pg|postgres)`)

	// Check for PostgreSQL dependency in package.json
	packageJSONPath := filepath.Join(projectPath, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		content, err := os.ReadFile(packageJSONPath)
		if err == nil {
			contentStr := string(content)
			// Direct PostgreSQL client packages
			if strings.Contains(contentStr, "\"pg\"") ||
				strings.Contains(contentStr, "\"postgres\"") ||
				strings.Contains(contentStr, "\"postgresql\"") ||
				strings.Contains(contentStr, "\"node-postgres\"") {
				return true, 0.9
			}

			// ORMs with PostgreSQL support
			if strings.Contains(contentStr, "\"sequelize\"") ||
				strings.Contains(contentStr, "\"typeorm\"") ||
				strings.Contains(contentStr, "\"prisma\"") ||
				strings.Contains(contentStr, "\"knex\"") {
				return true, 0.7 // Lower confidence since these can be used with other DBs
			}
		}
	}

	// Check for PostgreSQL connection strings in JavaScript/TypeScript files
	jsFiles, err := filepath.Glob(filepath.Join(projectPath, "**/*.{js,jsx,ts,tsx}"))
	if err == nil {
		for _, file := range jsFiles {
			content, err := os.ReadFile(file)
			if err == nil {
				// PostgreSQL connection strings
				if postgresConnectionPattern.MatchString(string(content)) {
					return true, 0.9
				}

				// PostgreSQL imports
				if postgresImportPattern.MatchString(string(content)) {
					return true, 0.8
				}

				// Other PostgreSQL indicators
				if strings.Contains(string(content), "createPool") && strings.Contains(string(content), "postgres") {
					return true, 0.9
				}
			}
		}
	}

	// Check for PostgreSQL connection strings in .env files
	envFiles, err := filepath.Glob(filepath.Join(projectPath, ".env*"))
	if err == nil {
		for _, file := range envFiles {
			content, err := os.ReadFile(file)
			if err == nil {
				if strings.Contains(string(content), "POSTGRES") ||
					strings.Contains(string(content), "PG_") ||
					postgresConnectionPattern.MatchString(string(content)) {
					return true, 0.8
				}
			}
		}
	}

	// Check for PostgreSQL configuration in JSON/YAML files
	configFiles, err := filepath.Glob(filepath.Join(projectPath, "**/*.{json,yaml,yml}"))
	if err == nil {
		for _, file := range configFiles {
			content, err := os.ReadFile(file)
			if err == nil {
				if strings.Contains(string(content), "dialect") && strings.Contains(string(content), "postgres") {
					return true, 0.9
				}
				if strings.Contains(string(content), "provider") && strings.Contains(string(content), "postgresql") {
					return true, 0.9
				}
			}
		}
	}

	return false, 0.0
}

// detectPostgresPython detects PostgreSQL usage in Python projects
func (d *PostgreSQLDetector) detectPostgresPython(projectPath string) (bool, float64) {
	// Define regex patterns
	psycopgImportPattern := regexp.MustCompile(`(?m)^import\s+psycopg2`)
	fromPsycopgImportPattern := regexp.MustCompile(`(?m)^from\s+psycopg2\s+import`)
	postgresConnectionPattern := regexp.MustCompile(`(?i)postgres(?:ql)?:\/\/[^\/\s]+\/\w+`)

	// Check for PostgreSQL dependency in requirements.txt
	requirementsPath := filepath.Join(projectPath, "requirements.txt")
	if _, err := os.Stat(requirementsPath); err == nil {
		content, err := os.ReadFile(requirementsPath)
		if err == nil {
			if strings.Contains(string(content), "psycopg2") ||
				strings.Contains(string(content), "psycopg2-binary") ||
				strings.Contains(string(content), "pg8000") ||
				strings.Contains(string(content), "sqlalchemy") {
				return true, 0.8
			}
		}
	}

	// Check for PostgreSQL dependency in pyproject.toml
	pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
	if _, err := os.Stat(pyprojectPath); err == nil {
		content, err := os.ReadFile(pyprojectPath)
		if err == nil {
			if strings.Contains(string(content), "psycopg2") ||
				strings.Contains(string(content), "pg8000") {
				return true, 0.9
			}
		}
	}

	// Check for PostgreSQL import statements in Python files
	pyFiles, err := filepath.Glob(filepath.Join(projectPath, "**/*.py"))
	if err == nil {
		for _, file := range pyFiles {
			content, err := os.ReadFile(file)
			if err == nil {
				// PostgreSQL connection strings
				if postgresConnectionPattern.MatchString(string(content)) {
					return true, 0.9
				}

				// PostgreSQL adapter imports
				if psycopgImportPattern.MatchString(string(content)) || fromPsycopgImportPattern.MatchString(string(content)) {
					return true, 0.9
				}

				// SQLAlchemy with PostgreSQL dialect
				if strings.Contains(string(content), "sqlalchemy") && strings.Contains(string(content), "postgresql") {
					return true, 0.9
				}

				// Django with PostgreSQL backend
				if strings.Contains(string(content), "ENGINE") && strings.Contains(string(content), "django.db.backends.postgresql") {
					return true, 0.9
				}
			}
		}
	}

	// Check for PostgreSQL connection strings in .env files
	envFiles, err := filepath.Glob(filepath.Join(projectPath, ".env*"))
	if err == nil {
		for _, file := range envFiles {
			content, err := os.ReadFile(file)
			if err == nil {
				if strings.Contains(string(content), "POSTGRES") ||
					strings.Contains(string(content), "PG_") ||
					postgresConnectionPattern.MatchString(string(content)) {
					return true, 0.8
				}
			}
		}
	}

	return false, 0.0
}

// detectPostgresGo detects PostgreSQL usage in Go projects
func (d *PostgreSQLDetector) detectPostgresGo(projectPath string) (bool, float64) {
	// Define regex patterns
	pgImportPattern1 := regexp.MustCompile(`(?m)import\s+\(.*?["']github\.com/lib/pq["'].*?\)`)
	pgImportPattern2 := regexp.MustCompile(`(?m)import\s+["']github\.com/lib/pq["']`)
	pgImportPattern3 := regexp.MustCompile(`(?m)import\s+\(.*?["']github\.com/jackc/pgx["'].*?\)`)
	pgImportPattern4 := regexp.MustCompile(`(?m)import\s+["']github\.com/jackc/pgx["']`)
	postgresConnectionPattern := regexp.MustCompile(`(?i)postgres(?:ql)?:\/\/[^\/\s]+\/\w+`)

	// Check for PostgreSQL dependency in go.mod
	goModPath := filepath.Join(projectPath, "go.mod")
	if _, err := os.Stat(goModPath); err == nil {
		content, err := os.ReadFile(goModPath)
		if err == nil {
			if strings.Contains(string(content), "github.com/lib/pq") ||
				strings.Contains(string(content), "github.com/jackc/pgx") ||
				strings.Contains(string(content), "github.com/go-pg/pg") {
				return true, 0.9
			}
		}
	}

	// Check for PostgreSQL import statements in Go files
	goFiles, err := filepath.Glob(filepath.Join(projectPath, "**/*.go"))
	if err == nil {
		for _, file := range goFiles {
			content, err := os.ReadFile(file)
			if err == nil {
				// PostgreSQL driver imports
				if pgImportPattern1.MatchString(string(content)) ||
					pgImportPattern2.MatchString(string(content)) ||
					pgImportPattern3.MatchString(string(content)) ||
					pgImportPattern4.MatchString(string(content)) {
					return true, 0.9
				}

				// PostgreSQL connection strings
				if postgresConnectionPattern.MatchString(string(content)) {
					return true, 0.9
				}

				// Other PostgreSQL indicators
				if strings.Contains(string(content), "postgres") && strings.Contains(string(content), "sql.Open") {
					return true, 0.9
				}
			}
		}
	}

	// Check for PostgreSQL connection strings in .env files
	envFiles, err := filepath.Glob(filepath.Join(projectPath, ".env*"))
	if err == nil {
		for _, file := range envFiles {
			content, err := os.ReadFile(file)
			if err == nil {
				if strings.Contains(string(content), "POSTGRES") ||
					strings.Contains(string(content), "PG_") ||
					postgresConnectionPattern.MatchString(string(content)) {
					return true, 0.8
				}
			}
		}
	}

	return false, 0.0
}

// Detect checks for PostgreSQL usage in the project
func (d *PostgreSQLDetector) Detect(ctx context.Context, dir string) (*detection.ProjectInfo, error) {
	info := &detection.ProjectInfo{
		Name:         filepath.Base(dir),
		Type:         "unknown",
		Dependencies: make(map[string]string),
		Metadata:     make(map[string]interface{}),
	}

	// Check for PostgreSQL in JavaScript/TypeScript projects
	jsDetected, jsConf := d.detectPostgresJS(dir)

	// Check for PostgreSQL in Python projects
	pyDetected, pyConf := d.detectPostgresPython(dir)

	// Check for PostgreSQL in Go projects
	goDetected, goConf := d.detectPostgresGo(dir)

	// Calculate overall confidence
	detected := jsDetected || pyDetected || goDetected
	confidence := 0.0
	count := 0

	if jsDetected {
		confidence += jsConf
		count++
	}

	if pyDetected {
		confidence += pyConf
		count++
	}

	if goDetected {
		confidence += goConf
		count++
	}

	if count > 0 {
		confidence /= float64(count)
	}

	if detected {
		info.Type = "postgresql"
		info.Confidence = confidence
		info.Dependencies["database"] = "postgresql"
	}

	return info, nil
}
