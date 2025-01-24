package ai

import (
	"context"
	"fmt"
	"os"
	"strings"

	openai "github.com/sashabaranov/go-openai"
)

type Client struct {
	openaiClient *openai.Client
}

func NewClient() (*Client, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is not set")
	}

	return &Client{
		openaiClient: openai.NewClient(apiKey),
	}, nil
}

func (c *Client) HandleAIFlag(ctx context.Context, command string, args []string) error {
	suggestions, err := c.getSuggestions(ctx, command, args)
	if err != nil {
		return fmt.Errorf("failed to get AI suggestions: %w", err)
	}

	fmt.Printf("\n🤖 AI Deployment Assistant:\n%s\n", suggestions)
	return nil
}

func (c *Client) getSuggestions(ctx context.Context, command string, args []string) (string, error) {
	prompt := c.buildPrompt(command, args)

	resp, err := c.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `You are a Nexlayer deployment expert. Provide concise, practical suggestions focused on solving deployment issues and optimizing configurations.

Available API endpoints:
- POST /startUserDeployment/{applicationID}
- POST /saveCustomDomain/{applicationID}
- GET /getDeployments/{applicationID}
- GET /getDeploymentInfo/{namespace}/{applicationID}

Format your response in these sections:
1. Quick Fixes (if applicable)
2. YAML Configuration Tips
3. Monitoring Commands
4. Next Steps

Keep each section brief and focused on actionable items. Include specific commands and YAML examples where relevant.

Example YAML structure:
'''yaml
app:
  name: myapp
  env: production
  resources:
    cpu: 1
    memory: 512Mi
  scaling:
    min: 1
    max: 5
'''

Example monitoring commands:
'''bash
nexlayer deploy status --app myapp  # Check deployment status
nexlayer deploy logs --app myapp    # View deployment logs
'''`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			Temperature: 0.7,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to get suggestions: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no suggestions received")
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *Client) buildPrompt(command string, args []string) string {
	return fmt.Sprintf(`Analyze this Nexlayer deployment and provide practical suggestions:
Command: %s
Arguments: %s

Focus on:
1. Immediate fixes for common deployment issues
2. YAML configuration optimizations (resources, scaling, env vars)
3. Specific monitoring commands to debug issues
4. Clear next steps

Provide real examples that work with Nexlayer's API endpoints.`, command, strings.Join(args, " "))
}
