package types

// Template represents the application template configuration
type Template struct {
	Name           string       `yaml:"name"`
	DeploymentName string       `yaml:"deploymentName"`
	RegistryLogin  RegistryAuth `yaml:"registryLogin"`
	Pods           []PodConfig  `yaml:"pods"`
	Build          struct {
		Command string `yaml:"command"`
		Output  string `yaml:"output"`
	} `yaml:"build"`
}

// Application represents the application configuration
type Application struct {
	Template Template `yaml:"template"`
}

// Config represents the application configuration
type Config struct {
	Application Application `yaml:"application"`
}

// RegistryAuth represents registry authentication configuration
type RegistryAuth struct {
	Registry            string `yaml:"registry"`
	Username            string `yaml:"username"`
	PersonalAccessToken string `yaml:"personalAccessToken"`
}

// PodConfig represents a pod configuration
type PodConfig struct {
	Type       string    `yaml:"type"`
	Name       string    `yaml:"name"`
	Tag        string    `yaml:"tag"`
	Vars       []VarPair `yaml:"vars"`
	ExposeHttp bool      `yaml:"exposeHttp"`
}

// VarPair represents a key-value pair for environment variables
type VarPair struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
