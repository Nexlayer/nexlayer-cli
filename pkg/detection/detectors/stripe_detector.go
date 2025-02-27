// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package detectors provides project-specific detection implementations.

// DEPRECATED: This specialized detector is deprecated and will be removed in a future version.
// The functionality has been replaced by the unified StackDetector in pkg/detection/stack_detector.go,
// which uses pattern-based detection to identify technology components including Stripe integrations.
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

// StripeDetector detects the presence of Stripe integration in a project
type StripeDetector struct {
	detection.BaseDetector
}

// NewStripeDetector creates a new detector for Stripe integration
func NewStripeDetector() *StripeDetector {
	detector := &StripeDetector{}
	detector.BaseDetector = *detection.NewBaseDetector("Stripe Detector", 0.8)
	return detector
}

// detectStripeJS detects Stripe integration in JavaScript/TypeScript projects
func (d *StripeDetector) detectStripeJS(projectPath string) (bool, float64) {
	// Check for Stripe dependency in package.json
	packageJSONPath := filepath.Join(projectPath, "package.json")
	if _, err := os.Stat(packageJSONPath); err == nil {
		content, err := os.ReadFile(packageJSONPath)
		if err == nil {
			if strings.Contains(string(content), "\"stripe\"") ||
				strings.Contains(string(content), "\"@stripe/stripe-js\"") ||
				strings.Contains(string(content), "\"@stripe/react-stripe-js\"") {
				return true, 0.9
			}
		}
	}

	// Check for import statements in JS/TS files
	stripejsImportPattern := regexp.MustCompile(`import\s+.*?from\s+['"]stripe['"]`)
	stripejsClientImportPattern := regexp.MustCompile(`import\s+.*?from\s+['"]@stripe/stripe-js['"]`)
	stripejsReactImportPattern := regexp.MustCompile(`import\s+.*?from\s+['"]@stripe/react-stripe-js['"]`)

	// Check for Stripe Elements usage
	stripeElementsPattern := regexp.MustCompile(`(Elements|CardElement|PaymentElement|StripeProvider)`)

	// Check for Stripe API key patterns
	stripeKeyPattern := regexp.MustCompile(`(pk_test_|sk_test_|pk_live_|sk_live_)[a-zA-Z0-9]{24,}`)

	// Check JavaScript/TypeScript files
	jsExtensions := []string{".js", ".jsx", ".ts", ".tsx"}
	for _, ext := range jsExtensions {
		matches, _ := filepath.Glob(filepath.Join(projectPath, "**/*"+ext))
		for _, file := range matches {
			content, err := os.ReadFile(file)
			if err == nil {
				contentStr := string(content)
				// Check for imports
				if stripejsImportPattern.MatchString(contentStr) ||
					stripejsClientImportPattern.MatchString(contentStr) ||
					stripejsReactImportPattern.MatchString(contentStr) {
					return true, 0.9
				}

				// Check for Stripe Elements
				if stripeElementsPattern.MatchString(contentStr) {
					return true, 0.8
				}

				// Check for API keys
				if stripeKeyPattern.MatchString(contentStr) {
					return true, 0.9
				}

				// Check for common Stripe usage
				if strings.Contains(contentStr, "new Stripe(") ||
					strings.Contains(contentStr, "stripe.customers") ||
					strings.Contains(contentStr, "stripe.checkout") ||
					strings.Contains(contentStr, "stripe.paymentIntents") {
					return true, 0.9
				}
			}
		}
	}

	// Check for environment variables in .env or similar files
	envMatches, _ := filepath.Glob(filepath.Join(projectPath, ".env*"))
	for _, file := range envMatches {
		content, err := os.ReadFile(file)
		if err == nil {
			contentStr := string(content)
			if strings.Contains(contentStr, "STRIPE_") ||
				strings.Contains(contentStr, "NEXT_PUBLIC_STRIPE_") ||
				stripeKeyPattern.MatchString(contentStr) {
				return true, 0.8
			}
		}
	}

	return false, 0.0
}

// detectStripePython detects Stripe integration in Python projects
func (d *StripeDetector) detectStripePython(projectPath string) (bool, float64) {
	// Check requirements.txt for stripe
	requirementsPath := filepath.Join(projectPath, "requirements.txt")
	if _, err := os.Stat(requirementsPath); err == nil {
		content, err := os.ReadFile(requirementsPath)
		if err == nil {
			if strings.Contains(string(content), "stripe") {
				return true, 0.9
			}
		}
	}

	// Check Poetry pyproject.toml
	pyprojectPath := filepath.Join(projectPath, "pyproject.toml")
	if _, err := os.Stat(pyprojectPath); err == nil {
		content, err := os.ReadFile(pyprojectPath)
		if err == nil {
			if strings.Contains(string(content), "stripe") {
				return true, 0.9
			}
		}
	}

	// Check for import statements in Python files
	importPattern := regexp.MustCompile(`(?m)^import\s+stripe`)
	fromImportPattern := regexp.MustCompile(`(?m)^from\s+stripe\s+import`)

	// Check for Stripe API key patterns
	stripeKeyPattern := regexp.MustCompile(`(pk_test_|sk_test_|pk_live_|sk_live_)[a-zA-Z0-9]{24,}`)

	// Check Python files
	pyMatches, _ := filepath.Glob(filepath.Join(projectPath, "**/*.py"))
	for _, file := range pyMatches {
		content, err := os.ReadFile(file)
		if err == nil {
			contentStr := string(content)
			// Check for imports
			if importPattern.MatchString(contentStr) || fromImportPattern.MatchString(contentStr) {
				return true, 0.9
			}

			// Check for API keys
			if stripeKeyPattern.MatchString(contentStr) {
				return true, 0.9
			}

			// Check for common Stripe usage
			if strings.Contains(contentStr, "stripe.Customer") ||
				strings.Contains(contentStr, "stripe.Charge") ||
				strings.Contains(contentStr, "stripe.PaymentIntent") ||
				strings.Contains(contentStr, "stripe.Subscription") {
				return true, 0.9
			}
		}
	}

	// Check for environment variables
	envMatches, _ := filepath.Glob(filepath.Join(projectPath, ".env*"))
	for _, file := range envMatches {
		content, err := os.ReadFile(file)
		if err == nil {
			contentStr := string(content)
			if strings.Contains(contentStr, "STRIPE_") || stripeKeyPattern.MatchString(contentStr) {
				return true, 0.8
			}
		}
	}

	return false, 0.0
}

// Detect implementation for StripeDetector
func (d *StripeDetector) Detect(ctx context.Context, dir string) (*detection.ProjectInfo, error) {
	// Emit deprecation warning
	detection.EmitDeprecationWarning("StripeDetector")

	// Create a basic project info
	projectInfo := &detection.ProjectInfo{
		Type:       "unknown",
		Path:       dir,
		Confidence: 0.0,
		Metadata:   make(map[string]interface{}),
	}

	// Try to detect Stripe in JavaScript/TypeScript
	jsFound, jsConf := d.detectStripeJS(dir)
	if jsFound {
		projectInfo.Type = "stripe"
		projectInfo.Confidence = jsConf
		return projectInfo, nil
	}

	// Try to detect Stripe in Python
	pyFound, pyConf := d.detectStripePython(dir)
	if pyFound {
		projectInfo.Type = "stripe"
		projectInfo.Confidence = pyConf
		projectInfo.Language = "Python"
		return projectInfo, nil
	}

	return projectInfo, nil
}
