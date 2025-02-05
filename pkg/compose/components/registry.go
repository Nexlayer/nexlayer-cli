package components

import "fmt"

const (
	NexlayerRegistry = "us-east1-docker.pkg.dev/nexlayer/components"
)

var (
	// ComponentRegistry maps component types to their configurations
	ComponentRegistry = map[string]ComponentConfig{
		"langfuse-ui": {
			Image: fmt.Sprintf("%s/langfuse:3", NexlayerRegistry),
			Ports: []Port{
				{Container: 3000, Host: 3000, Protocol: "tcp", Name: "http"},
			},
			Environment: []EnvVar{
				{Key: "DATABASE_URL", Value: "postgresql://postgres:postgres@postgres-db:5432/postgres", Required: true},
				{Key: "NEXTAUTH_URL", Value: "http://localhost:3000", Required: true},
				{Key: "NEXTAUTH_SECRET", Value: "mysecret", Required: true},
			},
		},
		"langfuse-worker": {
			Image: fmt.Sprintf("%s/langfuse-worker:3", NexlayerRegistry),
			Ports: []Port{
				{Container: 3030, Host: 3030, Protocol: "tcp", Name: "http"},
			},
			Environment: []EnvVar{
				{Key: "DATABASE_URL", Value: "postgresql://postgres:postgres@postgres-db:5432/postgres", Required: true},
				{Key: "SALT", Value: "mysalt", Required: true},
				{Key: "ENCRYPTION_KEY", Value: "0000000000000000000000000000000000000000000000000000000000000000", Required: true},
				{Key: "TELEMETRY_ENABLED", Value: "true", Required: false},
				{Key: "LANGFUSE_ENABLE_EXPERIMENTAL_FEATURES", Value: "true", Required: false},
			},
		},
		"postgres": {
			Image: fmt.Sprintf("%s/postgres:latest", NexlayerRegistry),
			Ports: []Port{
				{Container: 5432, Host: 5432, Protocol: "tcp", Name: "postgres"},
			},
			Environment: []EnvVar{
				{Key: "POSTGRES_USER", Value: "postgres", Required: true},
				{Key: "POSTGRES_PASSWORD", Value: "postgres", Required: true},
				{Key: "POSTGRES_DB", Value: "postgres", Required: true},
			},
			Healthcheck: &Healthcheck{
				Test:     []string{"CMD-SHELL", "pg_isready -U postgres"},
				Interval: "5s",
				Timeout:  "5s",
				Retries:  5,
			},
		},
		"redis": {
			Image: fmt.Sprintf("%s/redis:7", NexlayerRegistry),
			Ports: []Port{
				{Container: 6379, Host: 6379, Protocol: "tcp", Name: "redis"},
			},
			Command: []string{"redis-server", "--requirepass", "${REDIS_PASSWORD:-redis}"},
			Environment: []EnvVar{
				{Key: "REDIS_PASSWORD", Value: "redis", Required: false},
			},
			Healthcheck: &Healthcheck{
				Test:     []string{"CMD", "redis-cli", "ping"},
				Interval: "5s",
				Timeout:  "5s",
				Retries:  5,
			},
		},
		"clickhouse": {
			Image: fmt.Sprintf("%s/clickhouse:latest", NexlayerRegistry),
			Ports: []Port{
				{Container: 8123, Host: 8123, Protocol: "tcp", Name: "http"},
				{Container: 9000, Host: 9000, Protocol: "tcp", Name: "native"},
			},
			Environment: []EnvVar{
				{Key: "CLICKHOUSE_DB", Value: "default", Required: true},
				{Key: "CLICKHOUSE_USER", Value: "default", Required: true},
				{Key: "CLICKHOUSE_PASSWORD", Value: "default", Required: true},
			},
			Healthcheck: &Healthcheck{
				Test:     []string{"CMD", "wget", "--no-verbose", "--tries=1", "--spider", "http://localhost:8123/ping"},
				Interval: "5s",
				Timeout:  "5s",
				Retries:  5,
			},
		},
		"minio": {
			Image: fmt.Sprintf("%s/minio:latest", NexlayerRegistry),
			Ports: []Port{
				{Container: 9000, Host: 9090, Protocol: "tcp", Name: "api"},
				{Container: 9001, Host: 9091, Protocol: "tcp", Name: "console"},
			},
			Command: []string{"sh", "-c", "mkdir -p /data && minio server --address ':9000' --console-address ':9001' /data"},
			Environment: []EnvVar{
				{Key: "MINIO_ROOT_USER", Value: "minio", Required: true},
				{Key: "MINIO_ROOT_PASSWORD", Value: "miniosecret", Required: true},
			},
			Volumes: []Volume{
				{Source: "minio-data", Target: "/data", Type: "volume", Persistent: true},
			},
			Healthcheck: &Healthcheck{
				Test:     []string{"CMD", "curl", "-f", "http://localhost:9000/minio/health/live"},
				Interval: "5s",
				Timeout:  "5s",
				Retries:  5,
			},
		},
		"mongodb": {
			Image: "docker.io/library/mongo:latest",
			Ports: []Port{
				{Container: 27017, Host: 27017, Protocol: "tcp", Name: "mongodb"},
			},
			Environment: []EnvVar{
				{Key: "MONGO_INITDB_ROOT_USERNAME", Value: "mongo", Required: true},
				{Key: "MONGO_INITDB_ROOT_PASSWORD", Value: "mongo", Required: true},
			},
			Volumes: []Volume{
				{Source: "mongodb-data", Target: "/data/db", Type: "volume", Persistent: true},
			},
			Healthcheck: &Healthcheck{
				Test:     []string{"CMD", "mongosh", "--eval", "db.adminCommand('ping')"},
				Interval: "5s",
				Timeout:  "5s",
				Retries:  5,
			},
		},
	}

	// Add more component types as needed...
)

// GetComponentConfig returns the configuration for a given component type
func GetComponentConfig(componentType string) (ComponentConfig, error) {
	config, exists := ComponentRegistry[componentType]
	if !exists {
		return ComponentConfig{}, fmt.Errorf("unknown component type: %s", componentType)
	}
	return config, nil
}

// DetectComponentType analyzes the pod configuration and returns the detected component type
func DetectComponentType(pod interface{}) string {
	// TODO: Implement component detection logic based on:
	// 1. Explicit type in configuration
	// 2. Image name patterns
	// 3. Environment variables
	// 4. Port configurations
	// 5. Volume mounts
	return ""
}
