package vars

// Global variables used across different commands
var (
	// APIURL is the base URL for the Nexlayer API
	APIURL = "https://app.staging.nexlayer.io"

	// AppID is the ID of the application to operate on
	AppID string

	// AppName is the name of the application
	AppName string

	// Namespace is the deployment namespace
	Namespace string

	// Domain is the custom domain to set
	Domain string

	// ConfigFile is the path to the YAML configuration file
	ConfigFile string

	// ServiceName is the name of the service
	ServiceName string

	// Service is the service identifier
	Service string

	// EnvVars are environment variables
	EnvVars []string

	// EnvPairs are key-value pairs for environment variables
	EnvPairs []string

	// URL is the API URL override
	URL string

	// Registry configuration
	RegistryType     string // Container registry type (ghcr or dockerhub)
	Registry         string // Container registry URL
	RegistryUsername string // Registry username
	
	// Build configuration
	BuildContext string // Docker build context path
	ImageTag     string // Docker image tag

	// Graph configuration
	Depth        int    // Maximum depth to traverse when visualizing dependencies
	OutputFormat string // Format to use when visualizing dependencies
	OutputFile   string // File to write visualization output to
)
