package types

// NexlayerTemplate represents a complete infrastructure template
type NexlayerTemplate struct {
	Name      string            `json:"name" yaml:"name"`
	Version   string            `json:"version" yaml:"version"`
	Stack     ProjectStack      `json:"stack" yaml:"stack"`
	Services  []Service         `json:"services" yaml:"services"`
	Resources map[string]Resource `json:"resources,omitempty" yaml:"resources,omitempty"`
	Config    map[string]string `json:"config,omitempty" yaml:"config,omitempty"`
	Variables map[string]string `json:"variables,omitempty" yaml:"variables,omitempty"`
	Secrets   map[string]string `json:"secrets,omitempty" yaml:"secrets,omitempty"`
}

// ProjectStack represents the technology stack of a project
type ProjectStack struct {
	Language  string `json:"language" yaml:"language"`
	Framework string `json:"framework" yaml:"framework"`
	Database  string `json:"database" yaml:"database"`
}

// Service represents a deployable service
type Service struct {
	Name        string            `json:"name" yaml:"name"`
	Image       string            `json:"image" yaml:"image"`
	Command     []string          `json:"command,omitempty" yaml:"command,omitempty"`
	Environment map[string]string `json:"environment,omitempty" yaml:"environment,omitempty"`
	Ports       []PortConfig      `json:"ports,omitempty" yaml:"ports,omitempty"`
	Resources   ResourceRequests  `json:"resources" yaml:"resources"`
	Healthcheck *HealthcheckConfig `json:"healthcheck,omitempty" yaml:"healthcheck,omitempty"`
}

// PortConfig defines port mapping configuration
type PortConfig struct {
	Name        string `json:"name" yaml:"name"`
	Port        int    `json:"port" yaml:"port"`
	TargetPort  int    `json:"targetPort" yaml:"targetPort"`
	Protocol    string `json:"protocol" yaml:"protocol"`
	Host        bool   `json:"host" yaml:"host"`
	Public      bool   `json:"public" yaml:"public"`
	Healthcheck bool   `json:"healthcheck" yaml:"healthcheck"`
}

// ResourceRequests defines resource requirements
type ResourceRequests struct {
	CPU    string `json:"cpu" yaml:"cpu"`
	Memory string `json:"memory" yaml:"memory"`
}

// HealthcheckConfig defines health check configuration
type HealthcheckConfig struct {
	Path     string `json:"path" yaml:"path"`
	Port     int    `json:"port" yaml:"port"`
	Protocol string `json:"protocol" yaml:"protocol"`
}

// Resource defines a deployable resource
type Resource struct {
	Type       string            `json:"type" yaml:"type"`
	Version    string            `json:"version" yaml:"version"`
	Config     map[string]string `json:"config,omitempty" yaml:"config,omitempty"`
	Storage    []Storage         `json:"storage,omitempty" yaml:"storage,omitempty"`
	Network    Network           `json:"network,omitempty" yaml:"network,omitempty"`
}

// TemplateInfo represents template metadata for registry listings
type TemplateInfo struct {
	Name        string   `json:"name" yaml:"name"`
	Version     string   `json:"version" yaml:"version"`
	Description string   `json:"description" yaml:"description"`
	Tags        []string `json:"tags" yaml:"tags"`
	Author      string   `json:"author" yaml:"author"`
	Downloads   int      `json:"downloads" yaml:"downloads"`
	CreatedAt   string   `json:"created_at" yaml:"created_at"`
	UpdatedAt   string   `json:"updated_at" yaml:"updated_at"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Name  string `json:"name" yaml:"name"`
	Value string `json:"value" yaml:"value"`
}

// Port represents a network port configuration
type Port struct {
	Number   int    `json:"number" yaml:"number"`
	Protocol string `json:"protocol" yaml:"protocol"`
}

// Storage represents storage configuration
type Storage struct {
	Size string `json:"size" yaml:"size"`
	Type string `json:"type" yaml:"type"`
}

// Network represents network configuration
type Network struct {
	Ingress  string `json:"ingress" yaml:"ingress"`
	Egress   string `json:"egress" yaml:"egress"`
	Requests string `json:"requests" yaml:"requests"`
}

// Resources represents resource configuration
type Resources struct {
	CPU     string    `json:"cpu" yaml:"cpu"`
	Memory  string    `json:"memory" yaml:"memory"`
	Storage []Storage `json:"storage,omitempty" yaml:"storage,omitempty"`
	Network Network   `json:"network,omitempty" yaml:"network,omitempty"`
}

// SecurityIssue represents a security issue found in a template
type SecurityIssue struct {
	Type        string            `json:"type" yaml:"type"`
	Severity    string            `json:"severity" yaml:"severity"`
	Description string            `json:"description" yaml:"description"`
	Context     map[string]string `json:"context,omitempty" yaml:"context,omitempty"`
}

// ResourceCost represents the cost of a specific resource
type ResourceCost struct {
	Type         string  `json:"type" yaml:"type"`
	MonthlyCost  float64 `json:"monthly_cost" yaml:"monthly_cost"`
	Description  string  `json:"description" yaml:"description"`
}

// CostEstimate represents a complete cost estimation
type CostEstimate struct {
	TotalCost     float64        `json:"total_cost" yaml:"total_cost"`
	ResourceCosts []ResourceCost `json:"resource_costs" yaml:"resource_costs"`
	Currency      string         `json:"currency" yaml:"currency"`
}
