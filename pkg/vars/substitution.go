package vars

import (
	"fmt"
	"regexp"
	"strings"
)

// Variable types
const (
	// PodReferenceVar represents pod reference variables like <pod-name>.pod
	PodReferenceVar = "pod-reference"

	// TemplateVar represents template variables like <% VAR_NAME %>
	TemplateVar = "template-var"

	// EnvVar represents environment variables
	EnvVar = "env-var"
)

// Common variables used in templates
const (
	URLVar      = "URL"
	RegistryVar = "REGISTRY"
)

// VariableContext provides the context for variable substitution
type VariableContext struct {
	// Map of pod names to internal DNS names
	Pods map[string]string

	// Map of variable names to values
	Variables map[string]string

	// Registry URL for image references
	Registry string

	// Application URL
	URL string
}

// NewVariableContext creates a new variable context
func NewVariableContext() *VariableContext {
	return &VariableContext{
		Pods:      make(map[string]string),
		Variables: make(map[string]string),
	}
}

// AddPod adds a pod to the context
func (c *VariableContext) AddPod(name string) {
	c.Pods[name] = fmt.Sprintf("%s.pod", name)
}

// SetVariable sets a variable value
func (c *VariableContext) SetVariable(name, value string) {
	c.Variables[name] = value
}

// SetRegistry sets the registry URL
func (c *VariableContext) SetRegistry(registry string) {
	c.Registry = registry
	c.Variables[RegistryVar] = registry
}

// SetURL sets the application URL
func (c *VariableContext) SetURL(url string) {
	c.URL = url
	c.Variables[URLVar] = url
}

// Patterns to match variable references
var (
	// Pattern for pod references: podname.pod
	podRefPattern = regexp.MustCompile(`\b([a-zA-Z0-9_-]+)\.pod\b`)

	// Pattern for template variables: <% VAR_NAME %>
	templateVarPattern = regexp.MustCompile(`<%\s*([A-Z_]+)\s*%>`)

	// Pattern for environment variables: ${VAR_NAME} or $VAR_NAME
	envVarPattern = regexp.MustCompile(`\${([A-Z_][A-Z0-9_]*)}|\$([A-Z_][A-Z0-9_]*)`)
)

// SubstituteVariables replaces Nexlayer variables in a string
func SubstituteVariables(input string, ctx *VariableContext) (string, error) {
	if ctx == nil {
		return input, fmt.Errorf("variable context is required")
	}

	if input == "" {
		return "", nil
	}

	result := input

	// Replace pod references (e.g., postgres.pod)
	result = podRefPattern.ReplaceAllStringFunc(result, func(match string) string {
		podName := strings.TrimSuffix(match, ".pod")
		if podRef, ok := ctx.Pods[podName]; ok {
			return podRef
		}
		return match // Keep original if pod not found
	})

	// Replace template variables (e.g., <% URL %>)
	result = templateVarPattern.ReplaceAllStringFunc(result, func(match string) string {
		// Extract variable name
		varName := templateVarPattern.FindStringSubmatch(match)[1]

		// Special case for URL
		if varName == URLVar && ctx.URL != "" {
			return ctx.URL
		}

		// Special case for REGISTRY
		if varName == RegistryVar && ctx.Registry != "" {
			return ctx.Registry
		}

		// Try to lookup in the variables map
		if value, ok := ctx.Variables[varName]; ok {
			return value
		}

		return match // Keep original if variable not found
	})

	// Replace environment variables (e.g., ${HOME} or $HOME)
	result = envVarPattern.ReplaceAllStringFunc(result, func(match string) string {
		var varName string
		if strings.HasPrefix(match, "${") {
			// Extract name from ${VAR_NAME}
			varName = match[2 : len(match)-1]
		} else {
			// Extract name from $VAR_NAME
			varName = match[1:]
		}

		// Try to lookup in the variables map
		if value, ok := ctx.Variables[varName]; ok {
			return value
		}

		return match // Keep original if variable not found
	})

	return result, nil
}

// ExtractVariables finds all variable references in a string
func ExtractVariables(input string) map[string]string {
	result := make(map[string]string)

	if input == "" {
		return result
	}

	// Extract pod references
	podMatches := podRefPattern.FindAllStringSubmatch(input, -1)
	for _, match := range podMatches {
		if len(match) >= 2 {
			podName := match[1]
			result[podName+".pod"] = PodReferenceVar
		}
	}

	// Extract template variables
	templateMatches := templateVarPattern.FindAllStringSubmatch(input, -1)
	for _, match := range templateMatches {
		if len(match) >= 2 {
			varName := match[1]
			result[varName] = TemplateVar
		}
	}

	// Extract environment variables
	envMatches := envVarPattern.FindAllStringSubmatch(input, -1)
	for _, match := range envMatches {
		var varName string
		if match[1] != "" {
			varName = match[1] // From ${VAR_NAME}
		} else if match[2] != "" {
			varName = match[2] // From $VAR_NAME
		}
		if varName != "" {
			result[varName] = EnvVar
		}
	}

	return result
}

// ValidateSubstitution checks if all variables in a string can be substituted
func ValidateSubstitution(input string, ctx *VariableContext) ([]string, error) {
	missingVars := make([]string, 0)

	if ctx == nil {
		return nil, fmt.Errorf("variable context is required")
	}

	if input == "" {
		return nil, nil
	}

	// Find all variables
	variables := ExtractVariables(input)

	// Check each variable
	for name, varType := range variables {
		switch varType {
		case PodReferenceVar:
			podName := strings.TrimSuffix(name, ".pod")
			if _, ok := ctx.Pods[podName]; !ok {
				missingVars = append(missingVars, name)
			}

		case TemplateVar:
			if name == URLVar {
				if ctx.URL == "" {
					missingVars = append(missingVars, name)
				}
			} else if name == RegistryVar {
				if ctx.Registry == "" {
					missingVars = append(missingVars, name)
				}
			} else if _, ok := ctx.Variables[name]; !ok {
				missingVars = append(missingVars, name)
			}

		case EnvVar:
			if _, ok := ctx.Variables[name]; !ok {
				missingVars = append(missingVars, name)
			}
		}
	}

	if len(missingVars) > 0 {
		return missingVars, fmt.Errorf("unresolvable variables: %s", strings.Join(missingVars, ", "))
	}

	return nil, nil
}

// SubstitutePod replaces all variables in pod configuration
func SubstitutePod(pod interface{}, ctx *VariableContext) error {
	// This is a placeholder for implementing full pod configuration substitution
	// It would recursively walk through the pod structure and substitute variables
	// in all string fields
	return nil
}
