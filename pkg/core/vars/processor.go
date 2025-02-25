// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package vars

import (
	"fmt"
	"regexp"
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
		templateVarRegex: regexp.MustCompile(`<%\s*([A-Z_]+)(?:\s*:\s*([^%]*))?%>`),
		podRefRegex:      regexp.MustCompile(`([a-z][a-z0-9-]*).pod`),
	}

	// Register default variables
	p.RegisterVariable("REGISTRY", func(ctx VariableContext, _ ...string) (string, error) {
		registry := ctx.GetRegistry()
		if registry == "" {
			return "", fmt.Errorf("registry not set in context")
		}
		return registry, nil
	})

	p.RegisterVariable("URL", func(ctx VariableContext, _ ...string) (string, error) {
		url := ctx.GetURL()
		if url == "" {
			return "", fmt.Errorf("URL not set in context")
		}
		return url, nil
	})

	return p
}

// Process replaces template variables in the given string with their values
func (p *DefaultVariableProcessor) Process(input string, context VariableContext) (string, error) {
	if input == "" {
		return "", nil
	}

	// Process pod references (e.g., postgres.pod)
	result := p.podRefRegex.ReplaceAllStringFunc(input, func(match string) string {
		// Extract pod name
		submatches := p.podRefRegex.FindStringSubmatch(match)
		if len(submatches) < 2 {
			return match
		}

		podName := submatches[1]
		return context.GetPodName(podName)
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

	// Extract template variables
	var variables []string
	matches := p.templateVarRegex.FindAllStringSubmatch(input, -1)
	for _, match := range matches {
		if len(match) >= 2 {
			variables = append(variables, match[1])
		}
	}

	// Extract pod references
	podMatches := p.podRefRegex.FindAllStringSubmatch(input, -1)
	for _, match := range podMatches {
		if len(match) >= 2 {
			variables = append(variables, fmt.Sprintf("POD.%s", match[1]))
		}
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

// DefaultVariableContext is the default implementation of VariableContext
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

// WithRegistry returns a new context with the registry set
func (c *DefaultVariableContext) WithRegistry(registry string) *DefaultVariableContext {
	clone := c.clone()
	clone.registry = registry
	return clone
}

// WithURL returns a new context with the URL set
func (c *DefaultVariableContext) WithURL(url string) *DefaultVariableContext {
	clone := c.clone()
	clone.url = url
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

// clone creates a deep copy of the context
func (c *DefaultVariableContext) clone() *DefaultVariableContext {
	clone := &DefaultVariableContext{
		registry: c.registry,
		url:      c.url,
		pods:     make(map[string]string, len(c.pods)),
		values:   make(map[string]string, len(c.values)),
	}

	for k, v := range c.pods {
		clone.pods[k] = v
	}

	for k, v := range c.values {
		clone.values[k] = v
	}

	return clone
}
