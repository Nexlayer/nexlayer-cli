package models

import (
	"context"
)

// Provider defines the interface for AI providers
type Provider interface {
	// GetSuggestions gets AI suggestions for the given query
	GetSuggestions(ctx context.Context, query string) ([]string, error)

	// AnalyzeStack analyzes a project stack and returns deployment suggestions
	AnalyzeStack(ctx context.Context, projectPath string) (*StackAnalysis, error)
}

// MockProvider implements Provider interface for testing
type MockProvider struct{}

func NewMockProvider() *MockProvider {
	return &MockProvider{}
}

func (p *MockProvider) GetSuggestions(ctx context.Context, query string) ([]string, error) {
	return []string{"Mock suggestion 1", "Mock suggestion 2"}, nil
}

func (p *MockProvider) AnalyzeStack(ctx context.Context, projectPath string) (*StackAnalysis, error) {
	return &StackAnalysis{
		ContainerImage: "mock/image:latest",
		Dependencies:  []string{"mock-dep-1", "mock-dep-2"},
		Ports:        []int{8080},
		EnvVars:      []string{"MOCK_ENV=true"},
		Suggestions:  []string{"Mock suggestion"},
	}, nil
}
