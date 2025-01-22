// Formatted with gofmt -s
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/nexlayer/nexlayer-cli/plugins/template-builder/types"
	openai "github.com/sashabaranov/go-openai"
)

// AI refinement interface
type AIRefiner interface {
	RefineTemplate(template *types.NexlayerTemplate, stack types.ProjectStack) error
}

// OpenAI refiner implementation
type OpenAIRefiner struct {
	client *openai.Client
}

// NewOpenAIRefiner creates a new OpenAI refiner
func NewOpenAIRefiner() *OpenAIRefiner {
	apiKey := os.Getenv("OPENAI_API_KEY")
	return &OpenAIRefiner{
		client: openai.NewClient(apiKey),
	}
}

// RefineTemplate uses OpenAI to refine the template
func (r *OpenAIRefiner) RefineTemplate(template *types.NexlayerTemplate, stack types.ProjectStack) error {
	// Convert template to JSON for API request
	templateJSON, err := json.Marshal(template)
	if err != nil {
		return fmt.Errorf("error marshaling template: %v", err)
	}

	// Create completion request
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4,
		Messages: []openai.ChatCompletionMessage{
			{
				Role: openai.ChatMessageRoleSystem,
				Content: fmt.Sprintf(`You are an expert in Kubernetes and cloud-native deployments.
				Your task is to analyze and improve a Nexlayer deployment template for a %s application using %s framework and %s database.
				Focus on:
				1. Resource optimization
				2. Security best practices
				3. Scalability
				4. Environment configuration
				5. Health checks and monitoring
				
				Respond with ONLY the improved template in JSON format.`,
					stack.Language, stack.Framework, stack.Database),
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("Please improve this template:\n%s", string(templateJSON)),
			},
		},
		MaxTokens:   2000,
		Temperature: 0.2,
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Get completion from OpenAI
	resp, err := r.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return fmt.Errorf("error getting completion: %v", err)
	}

	// Parse improved template
	var improved types.NexlayerTemplate
	if err := json.Unmarshal([]byte(resp.Choices[0].Message.Content), &improved); err != nil {
		return fmt.Errorf("error parsing improved template: %v", err)
	}

	// Update original template with improvements
	*template = improved
	return nil
}

// Claude refiner implementation
type ClaudeRefiner struct {
	apiKey string
}

// NewClaudeRefiner creates a new Claude refiner
func NewClaudeRefiner() *ClaudeRefiner {
	return &ClaudeRefiner{
		apiKey: os.Getenv("CLAUDE_API_KEY"),
	}
}

// RefineTemplate uses Claude to refine the template
func (r *ClaudeRefiner) RefineTemplate(template *types.NexlayerTemplate, stack types.ProjectStack) error {
	// TODO: Implement Claude API integration
	return fmt.Errorf("Claude integration not implemented yet")
}

func NewAIRefiner() AIRefiner {
	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		return NewOpenAIRefiner()
	}
	if key := os.Getenv("CLAUDE_API_KEY"); key != "" {
		return NewClaudeRefiner()
	}
	return nil
}
