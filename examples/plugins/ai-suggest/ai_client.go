package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	openai "github.com/sashabaranov/go-openai"
	"gopkg.in/yaml.v3"
)

// AIClient interface for interacting with AI models
type AIClient interface {
	GetSuggestions(ctx context.Context, category string, appInfo map[string]interface{}) ([]string, error)
	StreamSuggestions(ctx context.Context, category string, appInfo map[string]interface{}) (<-chan string, <-chan error)
}

// OpenAIClient implements AIClient using OpenAI's GPT-4
type OpenAIClient struct {
	client *openai.Client
	docs   *DocSearch
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(docsPath, templatesPath string) (*OpenAIClient, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	apiKey := config.OpenAIKey
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	if apiKey == "" {
		return nil, fmt.Errorf("missing OpenAI API key. Please set it in the configuration file or as an environment variable")
	}

	docs, err := NewDocSearch(docsPath, templatesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load documentation: %w", err)
	}

	return &OpenAIClient{
		client: openai.NewClient(apiKey),
		docs:   docs,
	}, nil
}

// ClaudeClient implements AIClient using Anthropic's Claude
type ClaudeClient struct {
	apiKey string
	docs   *DocsContent
}

// NewClaudeClient creates a new Claude client
func NewClaudeClient(docsPath, templatesPath string) (*ClaudeClient, error) {
	config, err := loadConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	apiKey := config.AnthropicKey
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
	}

	if apiKey == "" {
		return nil, fmt.Errorf("missing Anthropic API key. Please set it in the configuration file or as an environment variable")
	}

	docs, err := LoadDocumentation(docsPath, templatesPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load documentation: %w", err)
	}

	return &ClaudeClient{
		apiKey: apiKey,
		docs:   docs,
	}, nil
}

// NewAIClient creates the appropriate AI client based on the model
func NewAIClient(model, docsPath, templatesPath string) (AIClient, error) {
	switch strings.ToLower(model) {
	case "openai":
		return NewOpenAIClient(docsPath, templatesPath)
	case "claude":
		return NewClaudeClient(docsPath, templatesPath)
	default:
		return nil, fmt.Errorf("unsupported AI model: %s. Supported models are 'openai' and 'claude'", model)
	}
}

// Config represents the configuration file structure
type Config struct {
	OpenAIKey    string `yaml:"openai_api_key"`
	AnthropicKey string `yaml:"anthropic_api_key"`
}

// loadConfig loads the configuration from a file
func loadConfig() (*Config, error) {
	file, err := os.Open("config.yaml")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	config := &Config{}
	if err := decoder.Decode(config); err != nil {
		return nil, err
	}

	return config, nil
}

// GetSuggestions gets AI-powered suggestions from OpenAI
func (c *OpenAIClient) GetSuggestions(ctx context.Context, category string, appInfo map[string]interface{}) ([]string, error) {
	prompt := c.generatePrompt(category, appInfo)

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: openai.GPT4,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `You are an AI assistant helping developers with their Nexlayer applications.
Format your responses in Markdown with the following structure:
1. Start each suggestion with a clear title on its own line
2. Follow with a detailed explanation
3. If providing code, use proper Markdown code blocks with language tags
4. Use bullet points and headings for better readability
5. Keep each suggestion focused and concise`,
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get suggestions from OpenAI: %w. Please check your network connection and API key", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no suggestions received from OpenAI. Please try again later or adjust your query")
	}

	// Split response into separate suggestions
	content := resp.Choices[0].Message.Content
	suggestions := strings.Split(content, "\n\n---\n\n")

	// Clean up each suggestion
	for i, s := range suggestions {
		suggestions[i] = strings.TrimSpace(s)
	}

	return suggestions, nil
}

// StreamSuggestions streams AI-powered suggestions from OpenAI
func (c *OpenAIClient) StreamSuggestions(ctx context.Context, category string, appInfo map[string]interface{}) (<-chan string, <-chan error) {
	prompt := c.generatePrompt(category, appInfo)

	suggestions := make(chan string)
	errors := make(chan error)

	go func() {
		defer close(suggestions)
		defer close(errors)

		stream, err := c.client.CreateChatCompletionStream(
			ctx,
			openai.ChatCompletionRequest{
				Model: openai.GPT4,
				Messages: []openai.ChatCompletionMessage{
					{
						Role: openai.ChatMessageRoleSystem,
						Content: `You are an AI assistant helping developers with their Nexlayer applications.
Format your responses in Markdown with the following structure:
1. Start each suggestion with a clear title on its own line
2. Follow with a detailed explanation
3. If providing code, use proper Markdown code blocks with language tags
4. Use bullet points and headings for better readability
5. Keep each suggestion focused and concise`,
					},
					{
						Role:    openai.ChatMessageRoleUser,
						Content: prompt,
					},
				},
			},
		)
		if err != nil {
			errors <- fmt.Errorf("failed to start streaming suggestions from OpenAI: %w", err)
			return
		}

		// Correctly handle the stream
		for {
			response, err := stream.Recv()
			if err != nil {
				if err == io.EOF {
					break
				}
				errors <- fmt.Errorf("error receiving from stream: %w", err)
				return
			}

			// Process the response
			for _, choice := range response.Choices {
				suggestions <- choice.Delta.Content
			}
		}
	}()
	return suggestions, errors
}

// generatePrompt creates a prompt for the AI model
func (c *OpenAIClient) generatePrompt(category string, appInfo map[string]interface{}) string {
	var prompt strings.Builder

	// Base prompt
	prompt.WriteString("You are an expert in Nexlayer deployments and cloud-native applications.\n\n")
	prompt.WriteString(fmt.Sprintf("Analyze this Nexlayer application and provide %s optimization suggestions based on Nexlayer's best practices:\n\n", category))

	// Add app info
	prompt.WriteString("Application Configuration:\n")
	for k, v := range appInfo {
		prompt.WriteString(fmt.Sprintf("%s: %v\n", k, v))
	}

	// Find relevant documentation based on category and app configuration
	searchQueries := c.generateSearchQueries(category, appInfo)
	var relevantDocs []string

	for _, query := range searchQueries {
		// Get relevant documentation
		docResults := c.docs.Search(query)
		var docContent []string
		for i, path := range docResults {
			if i >= 3 { // Limit to top 3 results
				break
			}
			content := c.docs.GetContent(path)
			if content != "" {
				docContent = append(docContent, fmt.Sprintf("From %s:\n%s", path, content))
			}
		}
		relevantDocs = append(relevantDocs, docContent...)
	}

	if len(relevantDocs) > 0 {
		prompt.WriteString("\nNexlayer Best Practices:\n")
		for _, doc := range relevantDocs {
			prompt.WriteString(doc + "\n")
		}
	}

	// Add category-specific instructions
	prompt.WriteString("\nProvide suggestions that follow Nexlayer's approach for:\n")
	switch category {
	case "deployment":
		prompt.WriteString(`
1. YAML configuration using startUserDeployment API
2. Resource allocation and scaling
3. Service dependencies and communication
4. Health checks and monitoring
5. Deployment strategies specific to Nexlayer
6. Environment configuration`)

	case "domain":
		prompt.WriteString(`
1. Custom domain setup using saveCustomDomain API
2. DNS configuration and verification
3. SSL/TLS management
4. Domain routing strategies
5. Subdomain organization
6. URL structure optimization`)

	case "status":
		prompt.WriteString(`
1. Deployment status monitoring via getDeployments API
2. Health check implementation
3. Log aggregation and analysis
4. Performance metrics collection
5. Alert configuration
6. Troubleshooting approaches`)

	case "template":
		prompt.WriteString(`
1. Template selection and customization
2. Resource optimization
3. Service configuration
4. Environment variables
5. Volume management
6. Template migration strategies`)
	}

	prompt.WriteString("\n\nFor each suggestion:\n")
	prompt.WriteString("1. Explain WHY it follows Nexlayer's best practices\n")
	prompt.WriteString("2. Show HOW to implement it using Nexlayer commands\n")
	prompt.WriteString("3. Include relevant configuration examples\n")
	prompt.WriteString("4. Reference specific Nexlayer features or APIs\n")

	return prompt.String()
}

func (c *OpenAIClient) generateSearchQueries(category string, appInfo map[string]interface{}) []string {
	queries := []string{
		// Base category query
		fmt.Sprintf("%s best practices", category),
	}

	// Add specific queries based on category and app info
	switch category {
	case "deployment":
		if services, ok := appInfo["services"].([]string); ok {
			for _, service := range services {
				queries = append(queries, fmt.Sprintf("deploy %s service", service))
			}
		}
		queries = append(queries,
			"deployment configuration yaml",
			"zero downtime deployment",
			"resource allocation",
		)

	case "domain":
		queries = append(queries,
			"custom domain configuration",
			"domain verification",
			"ssl setup",
			"subdomain management",
		)

	case "status":
		queries = append(queries,
			"deployment status monitoring",
			"health checks",
			"logging setup",
			"performance monitoring",
		)

	case "template":
		if templateName, ok := appInfo["templateName"].(string); ok {
			queries = append(queries, fmt.Sprintf("%s template configuration", templateName))
		}
		queries = append(queries,
			"template customization",
			"template resources",
			"environment configuration",
		)
	}

	return queries
}

// GetSuggestions gets AI-powered suggestions from Claude
func (c *ClaudeClient) GetSuggestions(ctx context.Context, category string, appInfo map[string]interface{}) ([]string, error) {
	// TODO: Implement Claude API integration
	// This will be similar to OpenAI but using Claude's API
	return nil, fmt.Errorf("Claude integration coming soon")
}

// StreamSuggestions streams AI-powered suggestions from Claude
func (c *ClaudeClient) StreamSuggestions(ctx context.Context, category string, appInfo map[string]interface{}) (<-chan string, <-chan error) {
	results := make(chan string)
	errors := make(chan error)
	go func() {
		defer close(results)
		defer close(errors)
		// Simulate streaming suggestions
		for i := 0; i < 5; i++ {
			results <- fmt.Sprintf("Suggestion %d", i)
			time.Sleep(1 * time.Second)
		}
	}()
	return results, errors
}
