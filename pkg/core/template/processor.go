// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package template

import (
	"fmt"
	"regexp"
	"strings"
)

// VariableProcessor defines the interface for processing template variables
type VariableProcessor interface {
	// Process replaces template variables in the given string with their values
	Process(input string, context VariableContext) (string, error)

	// ProcessMap replaces template variables in all values of the given map
	ProcessMap(input map[string]string, context VariableContext) (map[string]string, error)

	// Extract extracts template variables from the given string
	Extract(input string) []string

	// RegisterVariable registers a custom variable processor
	RegisterVariable(name string, processor VariableFunc)
}

// VariableContext provides context for variable processing
type VariableContext interface {
	// GetRegistry returns the registry URL
	GetRegistry() string

	// GetURL returns the application URL
	GetURL() string

	// GetPodName returns the internal DNS name for a pod
	GetPodName(name string) string

	// GetValue returns a custom variable value
	GetValue(name string) (string, bool)
}

// VariableFunc is a function that processes a template variable
type VariableFunc func(context VariableContext, args ...string) (string, error)

// DefaultVariableProcessor is the default implementation of VariableProcessor
type DefaultVariableProcessor struct {
	variables map[string]VariableFunc
	// Regular expressions for matching template variables
	templateVarRegex *regexp.Regexp
	podRefRegex      *regexp.Regexp
}

// NewVariableProcessor creates a new DefaultVariableProcessor
func NewVariableProcessor() *DefaultVariableProcessor {
	p := &DefaultVariableProcessor{
		variables:        make(map[string]VariableFunc),
		templateVarRegex: regexp.MustCompile(`<%\s*([^%>]+)\s*%>`),
		podRefRegex:      regexp.MustCompile(`(\w+)\.pod`),
	}

	// Register default variables
	p.RegisterVariable("REGISTRY", func(ctx VariableContext, _ ...string) (string, error) {
		if value, ok := ctx.GetValue("REGISTRY"); ok {
			return value, nil
		}
		return "", fmt.Errorf("registry not found in context")
	})

	p.RegisterVariable("URL", func(ctx VariableContext, _ ...string) (string, error) {
		if value, ok := ctx.GetValue("URL"); ok {
			return value, nil
		}
		return "", fmt.Errorf("url not found in context")
	})

	return p
}

// Process replaces template variables in the given string with their values
func (p *DefaultVariableProcessor) Process(input string, context VariableContext) (string, error) {
	if input == "" {
		return "", nil
	}

	// Process template variables (e.g., <% REGISTRY %> or <% REPEAT: hello, 2 %>)
	result := p.templateVarRegex.ReplaceAllStringFunc(input, func(match string) string {
		// Extract variable name and arguments
		submatches := p.templateVarRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		parts := strings.Split(submatches[1], ":")
		name := strings.TrimSpace(parts[0])
		var args []string
		if len(parts) > 1 {
			argStr := strings.TrimSpace(parts[1])
			args = strings.Split(argStr, ",")
			for i, arg := range args {
				args[i] = strings.TrimSpace(arg)
			}
		}

		// First check if it's a custom value in the context
		if value, ok := context.GetValue(name); ok {
			return value
		}

		// Then check if it's a registered processor
		if processor, ok := p.variables[name]; ok {
			if value, err := processor(context, args...); err == nil {
				return value
			}
		}

		return match
	})

	// Process pod references (e.g., postgres.pod)
	result = p.podRefRegex.ReplaceAllStringFunc(result, func(match string) string {
		submatches := p.podRefRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		podName := submatches[1]
		dnsName := context.GetPodName(podName)
		if dnsName == fmt.Sprintf("%s.pod", podName) && !strings.HasSuffix(match, ".pod") {
			return match
		}
		return dnsName
	})

	return result, nil
}

// ProcessMap replaces template variables in all values of the given map
func (p *DefaultVariableProcessor) ProcessMap(input map[string]string, context VariableContext) (map[string]string, error) {
	if input == nil {
		return nil, nil
	}

	result := make(map[string]string, len(input))
	for key, value := range input {
		processed, err := p.Process(value, context)
		if err != nil {
			return nil, fmt.Errorf("error processing value for key %q: %w", key, err)
		}
		result[key] = processed
	}

	return result, nil
}

// Extract extracts template variables from the given string
func (p *DefaultVariableProcessor) Extract(input string) []string {
	if input == "" {
		return nil
	}

	var variables []string
	seen := make(map[string]bool)

	// Extract pod references first
	podMatches := p.podRefRegex.FindAllStringSubmatch(input, -1)
	for _, match := range podMatches {
		if len(match) >= 2 && strings.HasSuffix(match[0], ".pod") {
			name := fmt.Sprintf("POD.%s", match[1])
			if !seen[name] {
				variables = append(variables, name)
				seen[name] = true
			}
		}
	}

	// Then extract template variables
	matches := p.templateVarRegex.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			parts := strings.Split(match[1], ":")
			name := strings.TrimSpace(parts[0])
			// Only include variables that are explicitly used in the template
			if !seen[name] && strings.Contains(input, fmt.Sprintf("<%% %s %%>", strings.TrimSpace(name))) {
				variables = append(variables, name)
				seen[name] = true
			}
		}
	}

	// Return nil if no variables were found
	if len(variables) == 0 {
		return nil
	}

	return variables
}

// RegisterVariable registers a custom variable processor
func (p *DefaultVariableProcessor) RegisterVariable(name string, processor VariableFunc) {
	p.variables[name] = processor
}

// ExtractPodReferences extracts pod references from the given string
func (p *DefaultVariableProcessor) ExtractPodReferences(input string) []string {
	if input == "" {
		return nil
	}

	var pods []string
	matches := p.podRefRegex.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			pods = append(pods, match[1])
		}
	}

	return pods
}

// DefaultVariableContext represents a context for variable processing
type DefaultVariableContext struct {
	registry string
	url      string
	pods     map[string]string
	values   map[string]string
}

// NewVariableContext creates a new DefaultVariableContext
func NewVariableContext() *DefaultVariableContext {
	return &DefaultVariableContext{
		pods:   make(map[string]string),
		values: make(map[string]string),
	}
}

// GetRegistry returns the registry URL
func (c *DefaultVariableContext) GetRegistry() string {
	return c.registry
}

// SetRegistry sets the registry URL
func (c *DefaultVariableContext) SetRegistry(registry string) {
	c.registry = registry
}

// GetURL returns the application URL
func (c *DefaultVariableContext) GetURL() string {
	return c.url
}

// SetURL sets the application URL
func (c *DefaultVariableContext) SetURL(url string) {
	c.url = url
}

// GetPodName returns the internal DNS name for a pod
func (c *DefaultVariableContext) GetPodName(name string) string {
	if pod, ok := c.pods[name]; ok {
		return pod
	}
	return fmt.Sprintf("%s.pod", name)
}

// AddPod adds a pod to the context
func (c *DefaultVariableContext) AddPod(name, dnsName string) {
	c.pods[name] = dnsName
}

// GetValue returns a custom variable value
func (c *DefaultVariableContext) GetValue(name string) (string, bool) {
	value, ok := c.values[name]
	return value, ok
}

// SetValue sets a custom variable value
func (c *DefaultVariableContext) SetValue(name, value string) {
	c.values[name] = value
}

// clone creates a deep copy of the context
func (c *DefaultVariableContext) clone() *DefaultVariableContext {
	newContext := &DefaultVariableContext{
		registry: c.registry,
		url:      c.url,
		pods:     make(map[string]string),
		values:   make(map[string]string),
	}
	for k, v := range c.pods {
		newContext.pods[k] = v
	}
	for k, v := range c.values {
		newContext.values[k] = v
	}
	return newContext
}

// WithRegistry returns a new context with the registry set
func (c *DefaultVariableContext) WithRegistry(registry string) *DefaultVariableContext {
	clone := c.clone()
	clone.registry = registry
	clone.values["REGISTRY"] = registry
	return clone
}

// WithURL returns a new context with the URL set
func (c *DefaultVariableContext) WithURL(url string) *DefaultVariableContext {
	clone := c.clone()
	clone.url = url
	clone.values["URL"] = url
	return clone
}

// WithPod returns a new context with a pod added
func (c *DefaultVariableContext) WithPod(name, dnsName string) *DefaultVariableContext {
	clone := c.clone()
	clone.pods[name] = dnsName
	return clone
}

// WithValue returns a new context with a custom value set
func (c *DefaultVariableContext) WithValue(name, value string) *DefaultVariableContext {
	clone := c.clone()
	clone.values[name] = value
	return clone
}

// NewVariableContextFromConfig creates a variable context from a NexlayerYAML configuration
func NewVariableContextFromConfig(config *NexlayerYAML) *DefaultVariableContext {
	ctx := NewVariableContext()

	if config == nil {
		return ctx
	}

	// Set application URL
	if config.Application.URL != "" {
		ctx.SetURL(config.Application.URL)
	}

	// Set registry from registry login
	if config.Application.RegistryLogin != nil && config.Application.RegistryLogin.Registry != "" {
		ctx.SetRegistry(config.Application.RegistryLogin.Registry)
	}

	// Add pods
	for _, pod := range config.Application.Pods {
		ctx.AddPod(pod.Name, fmt.Sprintf("%s.pod", pod.Name))
	}

	return ctx
}
