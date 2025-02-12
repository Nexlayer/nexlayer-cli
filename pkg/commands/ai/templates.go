package ai

// Template represents a Nexlayer YAML template
type Template struct {
	Application Application `yaml:"application"`
}

// Application represents the top-level application configuration
type Application struct {
	Name         string        `yaml:"name"`
	URL          string        `yaml:"url,omitempty"`
	RegistryLogin *RegistryLogin `yaml:"registryLogin,omitempty"`
	Pods         []Pod         `yaml:"pods"`
}

// RegistryLogin represents private registry authentication
type RegistryLogin struct {
	Registry           string `yaml:"registry"`
	Username           string `yaml:"username"`
	PersonalAccessToken string `yaml:"personalAccessToken"`
}

// Pod represents a single pod configuration
type Pod struct {
	Name         string        `yaml:"name"`
	Type         string        `yaml:"type"`
	Path         string        `yaml:"path,omitempty"`
	Image        string        `yaml:"image"`
	Volumes      []Volume      `yaml:"volumes,omitempty"`
	Secrets      []Secret      `yaml:"secrets,omitempty"`
	EnvVars      map[string]string `yaml:"env_vars,omitempty"`
	ServicePorts []ServicePort `yaml:"servicePorts"`
}

// Volume represents persistent storage configuration
type Volume struct {
	Name      string `yaml:"name"`
	Size      string `yaml:"size"`
	MountPath string `yaml:"mountPath"`
}

// Secret represents secret configuration
type Secret struct {
	Name      string `yaml:"name"`
	Data      string `yaml:"data"`
	MountPath string `yaml:"mountPath"`
	FileName  string `yaml:"fileName"`
}

// ServicePort represents port configuration
type ServicePort struct {
	ContainerPort int    `yaml:"containerPort"`
	ServicePort   int    `yaml:"servicePort"`
	Name          string `yaml:"name"`
}
