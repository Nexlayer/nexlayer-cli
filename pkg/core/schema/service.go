// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package schema

import (
	"context"
	"fmt"
	"sync"
)

// contextKey is a private type for context keys
type contextKey int

// processorKey is the key for the template processor in the context
const processorKey contextKey = iota

var (
	defaultProcessor     VariableProcessor
	defaultProcessorOnce sync.Once
)

// GetProcessor returns the template processor from the context or the default processor
func GetProcessor(ctx context.Context) VariableProcessor {
	if p, ok := ctx.Value(processorKey).(VariableProcessor); ok {
		return p
	}
	return getDefaultProcessor()
}

// WithProcessor returns a new context with the template processor
func WithProcessor(ctx context.Context, processor VariableProcessor) context.Context {
	return context.WithValue(ctx, processorKey, processor)
}

// getDefaultProcessor returns the default template processor
func getDefaultProcessor() VariableProcessor {
	defaultProcessorOnce.Do(func() {
		defaultProcessor = NewVariableProcessor()
	})
	return defaultProcessor
}

// ProcessString is a convenience function to process a template string
func ProcessString(ctx context.Context, input string, varCtx VariableContext) (string, error) {
	return GetProcessor(ctx).Process(input, varCtx)
}

// ProcessMap is a convenience function to process a map of template strings
func ProcessMap(ctx context.Context, input map[string]string, varCtx VariableContext) (map[string]string, error) {
	return GetProcessor(ctx).ProcessMap(input, varCtx)
}

// Extract is a convenience function to extract template variables from a string
func Extract(ctx context.Context, input string) []string {
	return GetProcessor(ctx).Extract(input)
}

// RegisterVariable is a convenience function to register a custom variable processor
func RegisterVariable(ctx context.Context, name string, processor VariableFunc) {
	GetProcessor(ctx).RegisterVariable(name, processor)
}

// NewContextFromConfig is a convenience function to create a variable context from a configuration
func NewContextFromConfig(config *NexlayerYAML) VariableContext {
	return NewVariableContextFromConfig(config)
}

// NewContext is a convenience function to create a new variable context
func NewContext() *DefaultVariableContext {
	return NewVariableContext()
}

// Service provides high-level schema management functionality
type Service struct {
	processor VariableProcessor
}

// NewService creates a new schema service
func NewService() *Service {
	return &Service{
		processor: NewVariableProcessor(),
	}
}

// ProcessConfig processes a Nexlayer YAML configuration
func (s *Service) ProcessConfig(config *NexlayerYAML) error {
	ctx := NewContextFromConfig(config)

	// Process all string fields in the configuration
	for i, pod := range config.Application.Pods {
		// Process image
		processedImage, err := s.processor.Process(pod.Image, ctx)
		if err != nil {
			return fmt.Errorf("failed to process pod %s image: %w", pod.Name, err)
		}
		config.Application.Pods[i].Image = processedImage

		// Process environment variables
		if len(pod.Vars) > 0 {
			for j, env := range pod.Vars {
				processedValue, err := s.processor.Process(env.Value, ctx)
				if err != nil {
					return fmt.Errorf("failed to process pod %s env var %s: %w", pod.Name, env.Key, err)
				}
				config.Application.Pods[i].Vars[j].Value = processedValue
			}
		}

		// Process annotations if present
		if len(pod.Annotations) > 0 {
			processedAnnotations, err := s.processor.ProcessMap(pod.Annotations, ctx)
			if err != nil {
				return fmt.Errorf("failed to process pod %s annotations: %w", pod.Name, err)
			}
			config.Application.Pods[i].Annotations = processedAnnotations
		}
	}

	return nil
}
