package types

// NexlayerTemplate represents a complete deployment template
type NexlayerTemplate struct {
	Application struct {
		Template struct {
			Name           string    `yaml:"name"`
			TemplateID     string    `yaml:"templateID"`
			DeploymentName string    `yaml:"deploymentName"`
			RegistryLogin  *Registry `yaml:"registryLogin,omitempty"`
			Pods          []Pod     `yaml:"pods"`
		} `yaml:"template"`
	} `yaml:"application"`
}

// Registry represents container registry login information
type Registry struct {
	Registry           string `yaml:"registry"`
	Username          string `yaml:"username"`
	PersonalAccessToken string `yaml:"personalAccessToken"`
}

// Pod represents a container in the deployment
type Pod struct {
	Type       string   `yaml:"type"`
	Name       string   `yaml:"name"`
	Image      string   `yaml:"image"`
	Tag        string   `yaml:"tag"`
	PrivateTag bool     `yaml:"privateTag"`
	ExposeHttp bool     `yaml:"exposeHttp"`
	Vars       []EnvVar `yaml:"vars,omitempty"`
}

// EnvVar represents an environment variable
type EnvVar struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// ProjectStack represents the detected project stack
type ProjectStack struct {
	Language      string
	Framework     string
	Database      string
	Frontend      string
	HasDocker     bool
	HasKubernetes bool
	Dependencies  map[string]string
}

// HasFrontend returns true if the stack has a frontend framework
func (s *ProjectStack) HasFrontend() bool {
	return s.Frontend != ""
}

// HasBackend returns true if the stack has a backend framework
func (s *ProjectStack) HasBackend() bool {
	return s.Framework != ""
}

// HasDatabase returns true if the stack has a database
func (s *ProjectStack) HasDatabase() bool {
	return s.Database != ""
}

// GetDatabaseType returns the standardized database type
func (s *ProjectStack) GetDatabaseType() string {
	switch s.Database {
	case "postgresql", "postgres":
		return "postgres"
	case "mongodb", "mongo":
		return "mongodb"
	case "mysql":
		return "mysql"
	default:
		return "postgres"
	}
}
