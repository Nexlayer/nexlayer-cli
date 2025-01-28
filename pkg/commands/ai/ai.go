package ai

import (
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

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
   - Database: mongodb, postgres, redis, neo4j
   - Others:
     - nginx: Used for load balancing or serving static assets.
     - llm: Custom pods for large language models or custom workloads (naming allowed).

4. **Expose HTTP**:
   - Use exposeHttp: true for components that need to be accessible via HTTP.

5. **Predefined Environment Variables**:
   - Include the following mandatory variables based on the stack:
     - DATABASE_CONNECTION_STRING for database connectivity.
     - FRONTEND_CONNECTION_URL and BACKEND_CONNECTION_URL for connecting frontend and backend services.

6. **Validation**:
   - Ensure the generated YAML conforms to Nexlayer's standards and avoid hallucinating unsupported configurations.

Based on the following inputs, generate a YAML deployment file:
- Application Name: %s
- Detected Stack: %s
- Required Components: %s

The YAML should be valid and follow all the requirements above. Only output the YAML content, nothing else.`

// Command represents the ai command
type Command struct {
	client interface{}
}

// NewCommand creates a new ai command
func NewCommand(client interface{}) *cobra.Command {
	c := &Command{
		client: client,
	}

	cmd := &cobra.Command{
		Use:   "ai",
		Short: "AI-powered features",
		Long:  "AI-powered features for Nexlayer CLI",
	}

	// Add detect subcommand
	cmd.AddCommand(c.newDetectCommand())

	return cmd
}

func (c *Command) newDetectCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "detect",
		Short: "Detect available AI assistants",
		Long:  "Detect and cache available AI assistants in your development environment",
		RunE: func(_ *cobra.Command, args []string) error {
			return detectAI()
		},
	}
}

func detectAI() error {
	var detectedModels []string

	// Detect Windsurf
	if _, err := os.Stat("/Applications/Windsurf.app"); err == nil {
		detectedModels = append(detectedModels, "Windsurf")
	}

	// Detect OpenAI API key
	if os.Getenv("OPENAI_API_KEY") != "" {
		detectedModels = append(detectedModels, "OpenAI")
	}

	// Detect Anthropic API key
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		detectedModels = append(detectedModels, "Claude")
	}

	if len(detectedModels) == 0 {
		return fmt.Errorf("no AI assistants detected")
	}

	fmt.Printf("Detected AI assistants: %s\n", strings.Join(detectedModels, ", "))
	return nil
}

// GenerateYAML generates a YAML template using AI
func GenerateYAML(appName string, stackType string, components []string) (string, error) {
	// TODO: Replace this with actual AI service call
	// For now, we'll use a mock response
	_ = fmt.Sprintf(yamlPrompt, appName, stackType, strings.Join(components, ", "))

	yaml := fmt.Sprintf(`application:
  template:
    name: "%s"
    deploymentName: "%s"
    registryLogin:
      registry: ghcr.io
      username: <Github username>
      personalAccessToken: <Github Packages Read-Only PAT>
    pods:
      - type: frontend
        name: frontend
        tag: node:18
        vars:
          - key: NODE_ENV
            value: development
          - key: PORT
            value: "3000"
        exposeHttp: true
      - type: database
        name: redis
        tag: redis:7-alpine
        vars:
          - key: REDIS_MAX_MEMORY
            value: 256mb
        exposeHttp: false
    build:
      command: npm install && npm run build
      output: build`, appName, appName)

	return yaml, nil
}
