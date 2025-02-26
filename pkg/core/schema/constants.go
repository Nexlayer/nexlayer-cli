// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package schema

// Pod types
const (
	// Frontend pod types
	PodTypeFrontend = "frontend"
	PodTypeReact    = "react"
	PodTypeNextJS   = "nextjs"
	PodTypeVue      = "vue"

	// Backend pod types
	PodTypeBackend = "backend"
	PodTypeExpress = "express"
	PodTypeDjango  = "django"
	PodTypeFastAPI = "fastapi"
	PodTypeNode    = "node"
	PodTypePython  = "python"
	PodTypeGolang  = "golang"
	PodTypeJava    = "java"

	// Database pod types
	PodTypeDatabase   = "database"
	PodTypeMongoDB    = "mongodb"
	PodTypePostgres   = "postgres"
	PodTypeRedis      = "redis"
	PodTypeMySQL      = "mysql"
	PodTypeClickhouse = "clickhouse"

	// Message Queue types
	PodTypeRabbitMQ = "rabbitmq"
	PodTypeKafka    = "kafka"

	// Storage types
	PodTypeMinio   = "minio"
	PodTypeElastic = "elasticsearch"

	// Web Server types
	PodTypeNginx   = "nginx"
	PodTypeTraefik = "traefik"

	// AI/ML pod types
	PodTypeLLM      = "llm"
	PodTypeOllama   = "ollama"
	PodTypeHFModel  = "huggingface"
	PodTypeVertexAI = "vertexai"
	PodTypeJupyter  = "jupyter"
)

// Protocol types
const (
	ProtocolTCP = "TCP"
	ProtocolUDP = "UDP"
)

// Volume types
const (
	VolumeTypePersistent = "persistent"
	VolumeTypeEphemeral  = "ephemeral"
)

// Registry and image defaults
const (
	DefaultRegistry = "ghcr.io/nexlayer"
	DefaultTag      = "latest"
	// Template placeholders
	RegistryPlaceholder = "<% REGISTRY %>"
	URLPlaceholder      = "<% URL %>"
)

// Default ports for different pod types
var DefaultPorts = map[string]int{
	PodTypeReact:    3000,
	PodTypeNextJS:   3000,
	PodTypeVue:      3000,
	PodTypeExpress:  3000,
	PodTypeDjango:   8000,
	PodTypeFastAPI:  8000,
	PodTypeNode:     3000,
	PodTypePython:   8000,
	PodTypeGolang:   8080,
	PodTypeJava:     8080,
	PodTypePostgres: 5432,
	PodTypeMongoDB:  27017,
	PodTypeRedis:    6379,
	PodTypeMySQL:    3306,
	PodTypeOllama:   11434,
	PodTypeJupyter:  8888,
}

// Default environment variables for different pod types
var DefaultEnvVars = map[string][]EnvVar{
	PodTypeReact: {
		{Key: "NODE_ENV", Value: "production"},
	},
	PodTypeNextJS: {
		{Key: "NODE_ENV", Value: "production"},
	},
	PodTypeVue: {
		{Key: "NODE_ENV", Value: "production"},
	},
	PodTypeExpress: {
		{Key: "NODE_ENV", Value: "production"},
		{Key: "PORT", Value: "3000"},
	},
	PodTypeDjango: {
		{Key: "DJANGO_SETTINGS_MODULE", Value: "config.settings.production"},
		{Key: "DJANGO_SECRET_KEY", Value: "<% DJANGO_SECRET_KEY %>"},
	},
	PodTypeFastAPI: {
		{Key: "PORT", Value: "8000"},
	},
	PodTypePostgres: {
		{Key: "POSTGRES_USER", Value: "postgres"},
		{Key: "POSTGRES_PASSWORD", Value: "<% POSTGRES_PASSWORD %>"},
		{Key: "POSTGRES_DB", Value: "app"},
	},
	PodTypeMongoDB: {
		{Key: "MONGO_INITDB_ROOT_USERNAME", Value: "root"},
		{Key: "MONGO_INITDB_ROOT_PASSWORD", Value: "<% MONGO_ROOT_PASSWORD %>"},
	},
	PodTypeRedis: {
		{Key: "REDIS_PASSWORD", Value: "<% REDIS_PASSWORD %>"},
	},
}
