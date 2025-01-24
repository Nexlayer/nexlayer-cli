package models

// Config holds the AI configuration
type Config struct {
	APIKey string
	Model  string
}

// DefaultConfig returns the default AI configuration
func DefaultConfig() *Config {
	return &Config{
		Model: "gpt-4",
	}
}

// StackAnalysis represents the analysis of a project stack
type StackAnalysis struct {
	ContainerImage string   `json:"container_image,omitempty"`
	Dependencies  []string `json:"dependencies,omitempty"`
	Ports         []int    `json:"ports,omitempty"`
	Resources     *ResourceRequests `json:"resources,omitempty"`
	EnvVars       []string `json:"env_vars,omitempty"`
	Suggestions   []string `json:"suggestions,omitempty"`
}

// ResourceRequests represents the resource requirements for a service
type ResourceRequests struct {
	CPU    string `json:"cpu,omitempty"`
	Memory string `json:"memory,omitempty"`
}

// PortConfig represents a port configuration
type PortConfig struct {
	Name        string `json:"name,omitempty"`
	Port        int    `json:"port"`
	TargetPort  int    `json:"target_port,omitempty"`
	Protocol    string `json:"protocol,omitempty"`
	Public      bool   `json:"public,omitempty"`
	Healthcheck bool   `json:"healthcheck,omitempty"`
}
