// Formatted with gofmt -s
package ai

import (
	"fmt"
	"os"
)

const (
	DefaultOpenAIModel = "gpt-3.5-turbo"
)

// NewAIClient creates a new AI client based on the provider
func NewAIClient(config Config) (AIClient, error) {
	switch config.Provider {
	case "openai":
		model := os.Getenv("OPENAI_MODEL")
		if model == "" {
			model = DefaultOpenAIModel
		}
		return &OpenAIClient{
			apiKey: config.APIKey,
			model:  model,
		}, nil
	case "claude":
		return &ClaudeClient{
			apiKey: config.APIKey,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported AI provider: %s", config.Provider)
	}
}
