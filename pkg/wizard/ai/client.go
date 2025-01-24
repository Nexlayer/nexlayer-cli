package ai

import (
	"context"
	"fmt"
	"os"

	"github.com/sashabaranov/go-openai"
)

// Client represents an AI client for the wizard
type Client struct {
	openAIClient *openai.Client
}

// NewClient creates a new AI client using OpenAI
func NewClient() (*Client, error) {
	openAIKey := os.Getenv("OPENAI_API_KEY")
	if openAIKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	client := openai.NewClient(openAIKey)
	return &Client{
		openAIClient: client,
	}, nil
}

// AnalyzeStack analyzes the current directory and returns deployment suggestions
func (c *Client) AnalyzeStack(ctx context.Context, dir string) ([]string, error) {
	// Create a prompt for OpenAI
	prompt := fmt.Sprintf(`Analyze the following project directory for deployment optimization:
Directory: %s

Please provide:
1. Resource optimization suggestions
2. Security best practices
3. Scaling recommendations
4. Performance tips

Format the response as a list of suggestions.`, dir)

	resp, err := c.openAIClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get AI suggestions: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no suggestions received from AI")
	}

	// Parse the response into a list of suggestions
	suggestions := []string{
		" AI Suggestions:",
		resp.Choices[0].Message.Content,
	}

	return suggestions, nil
}

// DetectProjectType attempts to determine the type of project in the directory
func (c *Client) DetectProjectType(dir string) (string, error) {
	// Implementation for project type detection
	return "", nil
}
