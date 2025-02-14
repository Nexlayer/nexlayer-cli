// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package template

// DefaultEnvVars provides common environment variables for different pod types
var DefaultEnvVars = map[PodType][]EnvVar{
	Frontend: {
		{Key: "REACT_APP_API_URL", Value: "http://CANDIDATE_DEPENDENCY_URL_0"},
		{Key: "NODE_ENV", Value: "production"},
	},
	Backend: {
		{Key: "PORT", Value: "8080"},
		{Key: "DATABASE_URL", Value: "postgresql://CANDIDATE_DEPENDENCY_URL_1:5432/db"},
	},
	Postgres: {
		{Key: "POSTGRES_USER", Value: "postgres"},
		{Key: "POSTGRES_DB", Value: "postgres"},
	},
	MongoDB: {
		{Key: "MONGO_INITDB_ROOT_USERNAME", Value: "root"},
	},
}

// DefaultPorts maps pod types to their default port configurations
var DefaultPorts = map[PodType][]Port{
	// Frontend ports
	React: {
		{ContainerPort: 3000, ServicePort: 80, Name: "web"},
	},
	NextJS: {
		{ContainerPort: 3000, ServicePort: 80, Name: "web"},
	},
	Vue: {
		{ContainerPort: 8080, ServicePort: 80, Name: "web"},
	},

	// Backend ports
	FastAPI: {
		{ContainerPort: 8000, ServicePort: 8000, Name: "api"},
	},
	Express: {
		{ContainerPort: 3000, ServicePort: 3000, Name: "api"},
	},
	Django: {
		{ContainerPort: 8000, ServicePort: 8000, Name: "api"},
	},

	// Database ports
	Postgres: {
		{ContainerPort: 5432, ServicePort: 5432, Name: "db"},
	},
	MongoDB: {
		{ContainerPort: 27017, ServicePort: 27017, Name: "db"},
	},
	Redis: {
		{ContainerPort: 6379, ServicePort: 6379, Name: "cache"},
	},

	// AI/ML ports
	Ollama: {
		{ContainerPort: 11434, ServicePort: 11434, Name: "llm"},
	},
	Jupyter: {
		{ContainerPort: 8888, ServicePort: 8888, Name: "notebook"},
	},
}

// DefaultVolumes defines standard volume configurations for stateful pods
var DefaultVolumes = map[PodType][]struct {
	Name      string
	Size      string
	MountPath string
}{
	Postgres: {
		{
			Name:      "data",
			Size:      "1Gi",
			MountPath: "/var/lib/postgresql/data",
		},
	},
	MongoDB: {
		{
			Name:      "data",
			Size:      "1Gi",
			MountPath: "/data/db",
		},
	},
	Redis: {
		{
			Name:      "data",
			Size:      "1Gi",
			MountPath: "/data",
		},
	},
}
