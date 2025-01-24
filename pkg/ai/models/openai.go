package models

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/sashabaranov/go-openai"
)

const nexlayerSystemPrompt = `You are Nexlayer Deploy Assistant, an AI specifically designed to help developers deploy and manage applications on Nexlayer's cloud platform. You are NOT a general coding assistant.

Focus areas:
1. Deployment Configuration
2. Resource Optimization
3. Security Best Practices
4. Scaling Strategies
5. Monitoring Setup

Your goal is to help developers rapidly deploy and scale their existing applications on Nexlayer, not to help develop the applications themselves.`

// OpenAIProvider implements the Provider interface using OpenAI
type OpenAIProvider struct {
	client *openai.Client
}

// NewOpenAIProvider creates a new OpenAI provider
func NewOpenAIProvider(apiKey string) (*OpenAIProvider, error) {
	client := openai.NewClient(apiKey)
	return &OpenAIProvider{
		client: client,
	}, nil
}

// GetSuggestions gets AI suggestions for the given query
func (p *OpenAIProvider) GetSuggestions(ctx context.Context, query string) ([]string, error) {
	resp, err := p.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: nexlayerSystemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: query,
				},
			},
			Temperature: 0.7,
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAI response: %v", err)
	}

	// Parse response into suggestions
	suggestions := make([]string, 0)
	if len(resp.Choices) > 0 {
		// Split response into bullet points
		content := resp.Choices[0].Message.Content
		suggestions = append(suggestions, content)
	}

	return suggestions, nil
}

// AnalyzeStack analyzes a project stack and returns deployment suggestions
func (p *OpenAIProvider) AnalyzeStack(ctx context.Context, projectPath string) (*StackAnalysis, error) {
	// Read project files
	files, err := readProjectFiles(projectPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read project files: %v", err)
	}

	// Create analysis request
	resp, err := p.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: nexlayerSystemPrompt,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Analyze project at path: %s\n\nFiles:\n%s", projectPath, files),
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to analyze stack: %v", err)
	}

	// Create basic analysis
	analysis := &StackAnalysis{
		ContainerImage: "default/image:latest",
		Dependencies:   []string{"node", "npm"},
		Ports:          []int{3000},
		EnvVars:        []string{"NODE_ENV=production"},
		Suggestions:    []string{resp.Choices[0].Message.Content},
	}

	return analysis, nil
}

// readProjectFiles reads important configuration files from the project
func readProjectFiles(projectPath string) (string, error) {
	var files string

	// List of important file patterns to look for
	patterns := []string{
		"Dockerfile",
		"docker-compose.yml",
		"package.json",
		"go.mod",
		"requirements.txt",
		"*.yaml",
		"*.yml",
	}

	// Walk through project directory
	err := filepath.Walk(projectPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check if file matches any pattern
		for _, pattern := range patterns {
			matched, err := filepath.Match(pattern, filepath.Base(path))
			if err != nil {
				return err
			}
			if matched {
				// Read file content
				content, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				// Add file content to analysis
				files += fmt.Sprintf("\n--- %s ---\n%s\n", path, string(content))
				break
			}
		}
		return nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to read project files: %v", err)
	}

	return files, nil
}
