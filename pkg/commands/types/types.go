package types

type Application struct {
	Template Template
}

type Template struct {
	Name           string
	DeploymentName string
	RegistryLogin  RegistryAuth
	Pods          []PodConfig
	Build         BuildConfig
}

type RegistryAuth struct {
	Registry            string
	Username           string
	PersonalAccessToken string
}

type PodConfig struct {
	Type       string
	Name       string
	Tag        string
	ExposeHttp bool
	Vars       []VarPair
}

type VarPair struct {
	Key   string
	Value string
}

type BuildConfig struct {
	Command string
	Output  string
}

type Config struct {
	Application Application
}
