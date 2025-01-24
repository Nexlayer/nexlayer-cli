package help

// Command help texts with examples
const (
	// Root command
	RootLongDesc = `Nexlayer CLI is a command line interface for managing your cloud-native applications.
It provides commands for deploying, managing, and monitoring your applications in a modern cloud environment.

Key Features:
• AI-powered deployment optimization (use --ai flag)
• Interactive deployment wizard with AI assistance
• Custom domain management
• Service scaling and monitoring
• Container registry integration
• Plugin system for extensibility

Environment Variables:
  NEXLAYER_CONFIG    Path to config file (default: $HOME/.nexlayer.yaml)
  NEXLAYER_TOKEN     Authentication token (optional)
  NEXLAYER_DEBUG     Enable debug logging (1=true, 0=false)
  NEXLAYER_AI_KEY    AI provider API key for AI features`

	RootExample = `  # Initialize a new project
  nexlayer init my-app

  # Deploy using the AI-powered wizard (recommended)
  nexlayer wizard

  # Deploy with AI optimization
  nexlayer deploy -f stack.yaml --ai

  # Add a custom domain
  nexlayer domain add my-app example.com`

	// Deploy command
	DeployLongDesc = `Deploy an application to Nexlayer cloud platform.
You can deploy using a YAML configuration file or the interactive wizard.

The deployment process includes:
1. Validating your configuration
2. Building and pushing container images
3. Creating necessary cloud resources
4. Setting up networking and routing
5. Starting your application containers

Use the --ai flag to get AI-powered suggestions for:
• Resource optimization
• Security improvements
• Cost reduction
• Performance tuning`

	DeployExample = `  # Deploy using a configuration file
  nexlayer deploy -f stack.yaml

  # Deploy with AI optimization
  nexlayer deploy -f stack.yaml --ai

  # Deploy with environment variables
  nexlayer deploy -f stack.yaml --env-file .env

  # Deploy with a specific namespace
  nexlayer deploy -f stack.yaml --namespace prod`

	// Service command
	ServiceLongDesc = `Manage services running on Nexlayer platform.
Services represent the running instances of your application components.
You can scale services, view logs, and manage configurations.`

	ServiceExample = `  # List all services
  nexlayer service list

  # Scale a service
  nexlayer service scale frontend --replicas 3

  # View service logs
  nexlayer service logs backend --tail 100

  # Restart a service
  nexlayer service restart frontend`

	// Domain command
	DomainLongDesc = `Manage custom domains for your applications.
Nexlayer automatically provisions SSL certificates and sets up DNS records.
You can add multiple domains to a single application.`

	DomainExample = `  # Add a domain
  nexlayer domain add my-app example.com

  # List domains
  nexlayer domain list my-app

  # Remove a domain
  nexlayer domain remove my-app example.com

  # Check domain status
  nexlayer domain status my-app example.com`

	// Registry command
	RegistryLongDesc = `Manage container registry operations.
Nexlayer provides a secure container registry for storing your application images.
You can push, pull, and manage container images.`

	RegistryExample = `  # Login to registry
  nexlayer registry login

  # Push an image
  nexlayer registry push my-app:latest

  # List images
  nexlayer registry list

  # Remove an image
  nexlayer registry remove my-app:old`
)
