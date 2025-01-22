package vars

// API variables
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

	// EnvVars are environment variables
	EnvVars []string

	// URL is the API URL override
	URL string
)

// CI variables
var (
	// BuildContext is the directory containing the Dockerfile
	BuildContext string

	// ImageTag is the tag to apply to the built Docker image
	ImageTag string
)

// Graph variables
var (
	// Depth is the maximum depth to traverse when visualizing dependencies
	Depth int

	// OutputFormat is the format to use when visualizing dependencies
	OutputFormat string

	// OutputFile is the file to write visualization output to
	OutputFile string

	// EnvPairs are key-value pairs for environment variables
	EnvPairs []string
)
