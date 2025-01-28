package ai

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// AIProvider represents an AI code assistant provider
type AIProvider struct {
	Name        string
	EnvVarKey   string
	Description string
	Endpoint    string
}

// Supported AI providers
var (
	GitHubCopilot = AIProvider{
		Name:        "GitHub Copilot",
		EnvVarKey:   "GITHUB_COPILOT_TOKEN",
		Description: "GitHub's AI pair programmer",
		Endpoint:    "https://api.github.com/copilot",
	}
	CursorAI = AIProvider{
		Name:        "Cursor AI",
		EnvVarKey:   "CURSOR_AI_KEY",
		Description: "AI-powered code editor",
		Endpoint:    "https://api.cursor.sh",
	}
	WindsurfEditor = AIProvider{
		Name:        "Windsurf Editor by Codeium",
		EnvVarKey:   "WINDSURF_API_KEY",
		Description: "World's first agentic IDE",
		Endpoint:    "https://api.codeium.com/windsurf",
	}
	Tabnine = AIProvider{
		Name:        "Tabnine",
		EnvVarKey:   "TABNINE_API_KEY",
		Description: "AI code completion assistant",
		Endpoint:    "https://api.tabnine.com",
	}
	JetBrainsAI = AIProvider{
		Name:        "JetBrains AI",
		EnvVarKey:   "JETBRAINS_AI_KEY",
		Description: "AI-powered development in JetBrains IDEs",
		Endpoint:    "https://api.jetbrains.com/ai",
	}
	IntelliCode = AIProvider{
		Name:        "Microsoft IntelliCode",
		EnvVarKey:   "INTELLICODE_KEY",
		Description: "AI-assisted development in Visual Studio",
		Endpoint:    "https://api.intellicode.microsoft.com",
	}
	CodeWhisperer = AIProvider{
		Name:        "Amazon CodeWhisperer",
		EnvVarKey:   "CODEWHISPERER_KEY",
		Description: "Amazon's AI code companion",
		Endpoint:    "https://api.aws.amazon.com/codewhisperer",
	}
	ClaudeSonnet = AIProvider{
		Name:        "Claude Sonnet 3.5",
		EnvVarKey:   "CLAUDE_API_KEY",
		Description: "Anthropic's advanced code assistant",
		Endpoint:    "https://api.anthropic.com/v1",
	}
	ChatGPTCode = AIProvider{
		Name:        "ChatGPT Code Interpreter",
		EnvVarKey:   "OPENAI_API_KEY",
		Description: "OpenAI's code interpreter",
		Endpoint:    "https://api.openai.com/v1",
	}
)

// All supported AI providers
var AllProviders = []AIProvider{
	GitHubCopilot,
	CursorAI,
	WindsurfEditor,
	Tabnine,
	JetBrainsAI,
	IntelliCode,
	CodeWhisperer,
	ClaudeSonnet,
	ChatGPTCode,
}

const yamlPrompt = `You are an AI assistant integrated into the Nexlayer CLI. Your job is to generate YAML deployment templates that strictly follow the Nexlayer template standards and specifications.

### Nexlayer Template Requirements:
1. **General Structure**:
   - The YAML must include application.template.name, application.template.deploymentName, and application.template.registryLogin fields.
   - The pods array must define the application's components, such as frontend, backend, database, or others (e.g., nginx, llm).

2. **Pod Configuration**:
   - Each pod must include:
     - type: The type of component (frontend, backend, database, nginx, llm, etc.).
     - name: A specific, descriptive name for the pod (e.g., backend-api, frontend-react, custom-llm-service).
     - tag: A valid Docker image tag (e.g., node:18, redis:7, nginx:1.23).
     - vars: Environment variables for the pod in the format [{ "key": "<VAR_NAME>", "value": "<VALUE>" }].

3. **Supported Pod Types**:
   - Frontend: react, angular, vue
   - Backend: express, django, fastapi
   - Database: mongodb, postgres, redis, neo4j, pinecone
   - Others:
     - nginx: Used for load balancing or serving static assets.
     - llm: Custom pods for large language models or custom workloads (naming allowed).

4. **Expose HTTP**:
   - Use exposeHttp: true for components that need to be accessible via HTTP.

5. **Predefined Environment Variables**:
   - Include the following mandatory variables based on the stack:
     - DATABASE_CONNECTION_STRING for database connectivity.
     - FRONTEND_CONNECTION_URL and BACKEND_CONNECTION_URL for connecting frontend and backend services.
     - PINECONE_API_KEY, PINECONE_ENVIRONMENT, PINECONE_INDEX for Pinecone vector database.

6. **Validation**:
   - Ensure the generated YAML conforms to Nexlayer's standards and avoid hallucinating unsupported configurations.

Based on the following inputs, generate a YAML deployment file:
- Application Name: %s
- Detected Stack: %s
- Required Components: %s

The YAML should be valid and follow all the requirements above. Only output the YAML content, nothing else.`

// GenerateYAML generates a YAML template using AI
func GenerateYAML(appName string, stackType string, components []string) (string, error) {
	// Get preferred AI provider from environment
	provider := getPreferredProvider()
	if provider == nil {
		return "", fmt.Errorf("no AI provider configured. Set one of the following environment variables: %s", getSupportedEnvVars())
	}

	// Format prompt
	prompt := fmt.Sprintf(yamlPrompt, appName, stackType, strings.Join(components, ", "))

	// TODO: Make actual API call to the provider
	// For now, return a mock response
	return mockGenerateYAML(appName, stackType, components), nil
}

// getPreferredProvider returns the first configured AI provider
func getPreferredProvider() *AIProvider {
	for _, provider := range AllProviders {
		if os.Getenv(provider.EnvVarKey) != "" {
			return &provider
		}
	}
	return nil
}

// getSupportedEnvVars returns a list of supported environment variables
func getSupportedEnvVars() string {
	var vars []string
	for _, provider := range AllProviders {
		vars = append(vars, provider.EnvVarKey)
	}
	return strings.Join(vars, ", ")
}

// mockGenerateYAML generates a mock YAML response
func mockGenerateYAML(appName string, stackType string, components []string) string {
	// ... (existing mock implementation)
	return "" // TODO: Implement mock response
}

// Command represents the ai command
func NewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI-powered features",
		Long:  "AI-powered features for generating and optimizing deployment templates",
	}

	cmd.AddCommand(newGenerateCommand())
	cmd.AddCommand(newDetectCommand())

	return cmd
}

func newGenerateCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "generate [app-name]",
		Short: "Generate deployment template using AI",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			appName := args[0]
			stackType, components := detectStack(".")
			yaml, err := GenerateYAML(appName, stackType, components)
			if err != nil {
				return err
			}
			fmt.Println(yaml)
			return nil
		},
	}
}

func newDetectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "detect",
		Short: "Detect available AI providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("Available AI Providers:")
			for _, provider := range AllProviders {
				configured := ""
				if os.Getenv(provider.EnvVarKey) != "" {
					configured = " (configured)"
				}
				fmt.Printf("- %s: %s%s\n", provider.Name, provider.Description, configured)
			}
			return nil
		},
	}
}
