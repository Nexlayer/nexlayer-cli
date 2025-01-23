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

		// Add template guidance
		prompt.WriteString(`
Template Guidance:
You are generating a Nexlayer infrastructure template. Please follow these guidelines:

1. Use the standard Nexlayer YAML structure as shown in the documentation
2. Include all necessary sections: infrastructure, application, environment, and networking
3. Utilize Nexlayer-provided environment variables instead of creating custom ones
4. Follow Nexlayer's best practices for resource naming and configuration
5. Ensure the template includes proper health checks and scaling configurations
6. Add appropriate comments to explain configuration choices
7. Include all required dependencies and their versions
8. Configure proper resource limits and requests
9. Set up appropriate logging and monitoring
10. Include security best practices

The template should be production-ready and follow cloud-native best practices.

Available Nexlayer environment variables:
- Core: NEXLAYER_APP_NAME, NEXLAYER_APP_VERSION, NEXLAYER_ENVIRONMENT, NEXLAYER_NAMESPACE, NEXLAYER_CLUSTER, NEXLAYER_REGION
- Networking: NEXLAYER_SERVICE_HOST, NEXLAYER_SERVICE_PORT, NEXLAYER_INGRESS_HOST, NEXLAYER_INGRESS_PORT, NEXLAYER_INTERNAL_URL, NEXLAYER_EXTERNAL_URL
- Database: NEXLAYER_DB_HOST, NEXLAYER_DB_PORT, NEXLAYER_DB_NAME, NEXLAYER_DB_USER, NEXLAYER_DB_PASSWORD, NEXLAYER_DB_URL
- Cache: NEXLAYER_CACHE_HOST, NEXLAYER_CACHE_PORT, NEXLAYER_CACHE_PASSWORD, NEXLAYER_CACHE_URL
- Storage: NEXLAYER_STORAGE_BUCKET, NEXLAYER_STORAGE_REGION, NEXLAYER_STORAGE_ACCESS_KEY, NEXLAYER_STORAGE_SECRET_KEY
- Security: NEXLAYER_API_KEY, NEXLAYER_JWT_SECRET, NEXLAYER_ENCRYPTION_KEY
- Monitoring: NEXLAYER_METRICS_PORT, NEXLAYER_TRACING_ENDPOINT, NEXLAYER_LOGGING_LEVEL`)
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

const (
	defaultModel = "gpt-4"
	maxTokens    = 2000
)

const TemplateGuidance = `You are generating a Nexlayer infrastructure template. Please follow these guidelines:

1. Use the standard Nexlayer pod-based YAML structure
2. Include required template fields: name, deploymentName, and optional registryLogin
3. Define pods with correct type, name, tag, and exposeHttp settings
4. Use Nexlayer-provided environment variables for service connections
5. Configure appropriate environment variables for each pod type
6. Follow pod naming conventions for databases, frontends, and backends
7. Set proper exposure settings for public-facing services
8. Include registry configuration for private images
9. Add descriptive comments for maintainability
10. Follow security best practices for sensitive variables

Template Structure Example:
application:
  template:
    name: "my-stack-name"
    deploymentName: "My Application"
    registryLogin:
      registry: ghcr.io
      username: <username>
      personalAccessToken: <token>
    pods:
    - type: database
      exposeHttp: false
      name: mongoDB
      tag: mongo:latest
      privateTag: false
      vars:
      - key: MONGO_INITDB_ROOT_USERNAME
        value: mongo

Supported Pod Types:
- Database: postgres, mysql, neo4j, redis, mongodb
- Frontend: react, angular, vue
- Backend: django, fastapi, express
- Others: nginx, llm

Available Nexlayer environment variables:
- Core: PROXY_URL, PROXY_DOMAIN
- Database: DATABASE_HOST, NEO4J_URI, DATABASE_CONNECTION_STRING
- Service URLs: FRONTEND_CONNECTION_URL, BACKEND_CONNECTION_URL, LLM_CONNECTION_URL
- Service Domains: FRONTEND_CONNECTION_DOMAIN, BACKEND_CONNECTION_DOMAIN, LLM_CONNECTION_DOMAIN`

type Client struct {
	openaiClient *openai.Client
	model       string
}

func NewClient(apiKey string) *Client {
	client := openai.NewClient(apiKey)
	return &Client{
		openaiClient: client,
		model:       defaultModel,
	}
}

func (c *Client) SetModel(model string) {
	c.model = model
}

func (c *Client) SuggestTemplate(ctx context.Context, requirements string) (string, error) {
	prompt := fmt.Sprintf("%s\n\nRequirements:\n%s\n\nPlease generate a Nexlayer template that meets these requirements:", TemplateGuidance, requirements)

	resp, err := c.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a Nexlayer infrastructure expert. Generate infrastructure templates that follow Nexlayer's pod-based architecture and best practices.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: maxTokens,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to get template suggestion: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no suggestions received")
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *Client) RefineTemplate(ctx context.Context, template, feedback string) (string, error) {
	prompt := fmt.Sprintf("%s\n\nCurrent template:\n%s\n\nFeedback:\n%s\n\nPlease refine this template based on the feedback:", TemplateGuidance, template, feedback)

	resp, err := c.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a Nexlayer infrastructure expert. Refine infrastructure templates to follow Nexlayer's pod-based architecture and best practices.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: maxTokens,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to get template refinement: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no refinements received")
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *Client) AnalyzeTemplate(ctx context.Context, template string) (string, error) {
	prompt := fmt.Sprintf("%s\n\nTemplate to analyze:\n%s\n\nPlease analyze this template and provide suggestions for improvement:", TemplateGuidance, template)

	resp, err := c.openaiClient.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleSystem,
					Content: "You are a Nexlayer infrastructure expert. Analyze infrastructure templates and provide suggestions for improvement based on Nexlayer's pod-based architecture and best practices.",
				},
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens: maxTokens,
		},
	)

	if err != nil {
		return "", fmt.Errorf("failed to analyze template: %w", err)
	}

	if len(resp.Choices) == 0 {
		return "", fmt.Errorf("no analysis received")
	}

	return resp.Choices[0].Message.Content, nil
}
