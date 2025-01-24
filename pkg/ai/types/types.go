package types

// AIProvider represents the type of AI model provider
type AIProvider string

const (
	OpenAI    AIProvider = "openai"
	Anthropic AIProvider = "claude"
)

// AIConfig holds the configuration for AI services
type AIConfig struct {
	Provider     AIProvider
	APIKey       string
	Temperature  float64
	MaxTokens    int
	DocsPath     string
	TemplatesDir string
}

// NexlayerTemplate represents a Nexlayer application template
type NexlayerTemplate struct {
	Name     string                 `json:"name" yaml:"name"`
	Version  string                 `json:"version" yaml:"version"`
	Stack    ProjectStack           `json:"stack" yaml:"stack"`
	Services []Service             `json:"services,omitempty" yaml:"services,omitempty"`
	Config   map[string]string     `json:"config,omitempty" yaml:"config,omitempty"`
	Resources []Resource           `json:"resources,omitempty" yaml:"resources,omitempty"`
}

// ProjectStack represents the detected project stack
type ProjectStack struct {
	Language    string   `json:"language" yaml:"language"`
	Framework   string   `json:"framework,omitempty" yaml:"framework,omitempty"`
	Database    string   `json:"database,omitempty" yaml:"database,omitempty"`
	Dependencies []string `json:"dependencies,omitempty" yaml:"dependencies,omitempty"`
}

// Service represents a service in the application
type Service struct {
	Name        string            `json:"name" yaml:"name"`
	Image       string            `json:"image" yaml:"image"`
	Ports       []PortConfig      `json:"ports,omitempty" yaml:"ports,omitempty"`
	Environment map[string]string `json:"environment,omitempty" yaml:"environment,omitempty"`
	Resources   ResourceRequests  `json:"resources,omitempty" yaml:"resources,omitempty"`
	Healthcheck *HealthcheckConfig `json:"healthcheck,omitempty" yaml:"healthcheck,omitempty"`
}

// PortConfig represents a port configuration
type PortConfig struct {
	Name       string `json:"name" yaml:"name"`
	Port       int    `json:"port" yaml:"port"`
	TargetPort int    `json:"targetPort" yaml:"targetPort"`
	Protocol   string `json:"protocol" yaml:"protocol"`
	Host       bool   `json:"host" yaml:"host"`
	Public     bool   `json:"public" yaml:"public"`
}

// ResourceRequests represents container resource requirements
type ResourceRequests struct {
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
	GPU    string `json:"gpu,omitempty" yaml:"gpu,omitempty"`
}

// HealthcheckConfig represents a service healthcheck configuration
type HealthcheckConfig struct {
	Path     string `json:"path" yaml:"path"`
	Port     int    `json:"port" yaml:"port"`
	Protocol string `json:"protocol" yaml:"protocol"`
}

// Resource represents a cloud resource
type Resource struct {
	Type     string            `json:"type" yaml:"type"`
	Name     string            `json:"name" yaml:"name"`
	Provider string            `json:"provider" yaml:"provider"`
	Config   map[string]string `json:"config" yaml:"config"`
}

// SecurityIssue represents a detected security issue
type SecurityIssue struct {
	Level       string `json:"level" yaml:"level"`
	Component   string `json:"component" yaml:"component"`
	Description string `json:"description" yaml:"description"`
	Mitigation  string `json:"mitigation" yaml:"mitigation"`
}

// CostEstimate represents estimated resource costs
type CostEstimate struct {
	Monthly     float64        `json:"monthly" yaml:"monthly"`
	Hourly      float64        `json:"hourly" yaml:"hourly"`
	Currency    string         `json:"currency" yaml:"currency"`
	Components  []ResourceCost `json:"components" yaml:"components"`
}

// ResourceCost represents the cost for a specific resource
type ResourceCost struct {
	Name     string  `json:"name" yaml:"name"`
	Type     string  `json:"type" yaml:"type"`
	Monthly  float64 `json:"monthly" yaml:"monthly"`
	Hourly   float64 `json:"hourly" yaml:"hourly"`
}

// StackAnalysis represents AI analysis of a project stack
type StackAnalysis struct {
	ContainerImage string            `json:"container_image"`
	EnvVars       []EnvVar          `json:"env_vars"`
	Ports         []Port            `json:"ports"`
	Resources     ResourceRequests  `json:"resources"`
	Dependencies  []string          `json:"dependencies"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
}

// Port represents a port configuration
type Port struct {
	Number      int    `json:"number"`
	Protocol    string `json:"protocol"`
	Purpose     string `json:"purpose"`
	Description string `json:"description"`
}
