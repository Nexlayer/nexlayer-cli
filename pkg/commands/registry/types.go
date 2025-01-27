package registry

// ImageConfig represents configuration for a Docker image
type ImageConfig struct {
	ServiceName string
	Path        string
	Tags        []string
	Namespace   string
}

// BuildConfig represents configuration for building multiple images
type BuildConfig struct {
	Images    []ImageConfig
	Namespace string
	Tags      []string
}

// RegistryConfig represents configuration for container registry
type RegistryConfig struct {
	Username string
	Token    string
	Registry string // e.g., "ghcr.io"
}
