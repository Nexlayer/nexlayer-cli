package templatebuilder

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/v2/types"
	openai "github.com/sashabaranov/go-openai"
)

// AIRefiner is the interface for template refinement
type AIRefiner interface {
	RefineTemplate(stack types.ProjectStack, template *types.NexlayerTemplate) (*types.NexlayerTemplate, error)
}

// OpenAIRefiner implements AIRefiner using OpenAI
type OpenAIRefiner struct {
	client *openai.Client
}

// NewOpenAIRefiner creates a new OpenAI refiner
func NewOpenAIRefiner() (*OpenAIRefiner, error) {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required")
	}

	return &OpenAIRefiner{
		client: openai.NewClient(apiKey),
	}, nil
}

// RefineTemplate uses OpenAI to refine the template
func (r *OpenAIRefiner) RefineTemplate(stack types.ProjectStack, template *types.NexlayerTemplate) (*types.NexlayerTemplate, error) {
	// Convert template to JSON for API request
	templateJSON, err := json.Marshal(template)
	if err != nil {
		return nil, fmt.Errorf("error marshaling template: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := r.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: fmt.Sprintf(`You are an expert in Nexlayer deployments. 
Analyze and improve this template for a %s application using %s framework and %s database.
Focus on:
1. Resource optimization
2. Security best practices
3. Scalability
4. Monitoring and health checks
5. Environment configuration

Respond with the improved template in JSON format.`,
						stack.Language, stack.Framework, stack.Database),
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: fmt.Sprintf("Please improve this template:\n%s", string(templateJSON)),
				},
			},
			MaxTokens:   2000,
			Temperature: 0.2,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %v", err)
	}

	// Parse improved template
	var improved types.NexlayerTemplate
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &improved); err != nil {
		return nil, fmt.Errorf("error parsing improved template: %v", err)
	}

	return &improved, nil
}

// ClaudeRefiner implements AIRefiner using Anthropic's Claude
type ClaudeRefiner struct {
	apiKey string
}

// NewClaudeRefiner creates a new Claude refiner
func NewClaudeRefiner() (*ClaudeRefiner, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
	}
	return &ClaudeRefiner{apiKey: apiKey}, nil
}

// RefineTemplate uses Claude to refine the template
func (r *ClaudeRefiner) RefineTemplate(stack types.ProjectStack, template *types.NexlayerTemplate) (*types.NexlayerTemplate, error) {
	// TODO: Implement Claude API integration
	return nil, fmt.Errorf("Claude integration not implemented yet")
}

// NewAIRefiner creates a new AI refiner based on environment configuration
func NewAIRefiner() (AIRefiner, error) {
	// Try OpenAI first
	if refiner, err := NewOpenAIRefiner(); err == nil {
		return refiner, nil
	}

	// Fall back to Claude
	if refiner, err := NewClaudeRefiner(); err == nil {
		return refiner, nil
	}

	return nil, fmt.Errorf("no AI refiner available - set OPENAI_API_KEY or ANTHROPIC_API_KEY")
}
