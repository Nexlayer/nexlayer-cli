// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package compose provides Docker Compose to Nexlayer YAML conversion functionality.
package compose

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/ai"
	"github.com/Nexlayer/nexlayer-cli/pkg/core/schema"
	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
	"github.com/Nexlayer/nexlayer-cli/pkg/knowledge"
	"github.com/Nexlayer/nexlayer-cli/pkg/vars"
)

// DockerComposeService represents a service in docker-compose.yml
type DockerComposeService struct {
	Image         string                 `yaml:"image"`
	Build         interface{}            `yaml:"build,omitempty"`
	Command       interface{}            `yaml:"command,omitempty"`
	Entrypoint    interface{}            `yaml:"entrypoint,omitempty"`
	Environment   interface{}            `yaml:"environment,omitempty"`
	EnvFile       interface{}            `yaml:"env_file,omitempty"`
	Ports         interface{}            `yaml:"ports,omitempty"`
	Volumes       interface{}            `yaml:"volumes,omitempty"`
	DependsOn     interface{}            `yaml:"depends_on,omitempty"`
	Networks      interface{}            `yaml:"networks,omitempty"`
	Restart       string                 `yaml:"restart,omitempty"`
	Links         []string               `yaml:"links,omitempty"`
	ExtraHosts    []string               `yaml:"extra_hosts,omitempty"`
	ExtraSettings map[string]interface{} `yaml:",inline,omitempty"`
	Secrets       []interface{}          `yaml:"secrets,omitempty"`
}

// DockerComposeConfig represents the structure of a docker-compose.yml file
type DockerComposeConfig struct {
	Version    string                          `yaml:"version,omitempty"`
	Services   map[string]DockerComposeService `yaml:"services"`
	Volumes    map[string]interface{}          `yaml:"volumes,omitempty"`
	Networks   map[string]interface{}          `yaml:"networks,omitempty"`
	Secrets    map[string]interface{}          `yaml:"secrets,omitempty"`
	ConfigPath string
}

// ConvertOptions provides configuration options for the conversion process
type ConvertOptions struct {
	ProjectDir      string
	ApplicationName string
	ForceConversion bool
	ComposeFileName string
	ApplicationURL  string
	RegistryURL     string
	UseAI           bool
}

// DefaultPorts maps common images to their default ports for intelligent port assignment
var DefaultPorts = map[string]int{
	"postgres":   5432,
	"mysql":      3306,
	"redis":      6379,
	"nginx":      80,
	"apache":     80,
	"node":       3000,
	"mongo":      27017,
	"clickhouse": 8123,
	"minio":      9000,
}

// DefaultVolumeSizes maps service types to default volume sizes
var DefaultVolumeSizes = map[string]string{
	"postgres":   "10Gi",
	"mysql":      "10Gi",
	"mongo":      "10Gi",
	"clickhouse": "10Gi",
	"minio":      "10Gi",
	"default":    "1Gi",
}

// ParsePortMapping parses a Docker port mapping string like "8080:80/tcp"
func ParsePortMapping(portStr, serviceName string) (int, int, string, error) {
	protocol := "TCP"
	portParts := strings.Split(portStr, "/")
	if len(portParts) > 1 {
		protocol = strings.ToUpper(portParts[1])
	}

	ports := strings.Split(portParts[0], ":")
	if len(ports) == 1 {
		port, err := strconv.Atoi(ports[0])
		if err != nil {
			log.Printf("Warning: Invalid port '%s' for service '%s'", ports[0], serviceName)
			return 0, 0, "", fmt.Errorf("invalid port: %s", ports[0])
		}
		return port, port, protocol, nil
	} else if len(ports) == 2 {
		externalPort, err := strconv.Atoi(ports[0])
		if err != nil {
			log.Printf("Warning: Invalid external port '%s' for service '%s'", ports[0], serviceName)
			return 0, 0, "", fmt.Errorf("invalid external port: %s", ports[0])
		}
		internalPort, err := strconv.Atoi(ports[1])
		if err != nil {
			log.Printf("Warning: Invalid internal port '%s' for service '%s'", ports[1], serviceName)
			return 0, 0, "", fmt.Errorf("invalid internal port: %s", ports[1])
		}
		return externalPort, internalPort, protocol, nil
	}
	log.Printf("Warning: Invalid port mapping '%s' for service '%s'", portStr, serviceName)
	return 0, 0, "", fmt.Errorf("invalid port mapping: %s", portStr)
}

// ParseVolumeMapping parses a Docker volume mapping string like "/host/path:/container/path:ro"
func ParseVolumeMapping(volumeStr, serviceName string) (string, string, bool, error) {
	readOnly := false
	parts := strings.Split(volumeStr, ":")
	if len(parts) == 1 {
		return parts[0], "/" + parts[0], readOnly, nil
	} else if len(parts) >= 2 {
		hostPath := parts[0]
		containerPath := parts[1]
		if len(parts) > 2 && parts[2] == "ro" {
			readOnly = true
		}
		return hostPath, containerPath, readOnly, nil
	}
	log.Printf("Warning: Invalid volume mapping '%s' for service '%s'", volumeStr, serviceName)
	return "", "", false, fmt.Errorf("invalid volume mapping: %s", volumeStr)
}

// Convert converts a Docker Compose configuration to Nexlayer YAML with enhanced validation
func Convert(composeFilePath string, opts ConvertOptions) (*schema.NexlayerYAML, error) {
	// Use AI enhancement by default unless explicitly disabled
	if !opts.UseAI {
		return convertBasic(composeFilePath, opts)
	}

	// Create a context with timeout for AI operations
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Perform the basic conversion first
	config, err := convertBasic(composeFilePath, opts)
	if err != nil {
		return nil, err
	}

	// Get the project directory containing the compose file
	projectDir := filepath.Dir(composeFilePath)

	// Initialize AI enhancement components if enabled
	var enhancer *ai.Enhancer
	detectionManager := initializeDetectionManager()
	llmEnricher := initializeLLMEnricher()

	if llmEnricher != nil {
		enhancer = ai.NewEnhancer(llmEnricher, detectionManager)

		// Perform AI enhancement in the background
		resultCh, errCh := enhancer.EnhanceAsync(ctx, config, projectDir)

		// Wait for enhancement to complete or timeout
		select {
		case result := <-resultCh:
			if result != nil {
				// Apply AI suggestions as comments in the YAML
				applyAISuggestions(config, result)

				// Display enhancement suggestions to the user
				printEnhancementSuggestions(result)

				// Check for critical issues and warn user
				if hasCriticalIssues(result) {
					log.Printf("⚠️ Warning: The generated configuration has potential issues. Review the suggestions above.")
				} else {
					log.Printf("✅ AI analysis complete: Configuration looks good!")
				}
			}
		case err := <-errCh:
			log.Printf("⚠️ AI enhancement failed: %v", err)
		case <-ctx.Done():
			log.Printf("⚠️ AI enhancement timed out, proceeding with basic configuration")
		}
	}

	return config, nil
}

// convertBasic performs the basic Docker Compose to Nexlayer YAML conversion
// This is the original conversion logic from before AI enhancement
func convertBasic(composeFilePath string, opts ConvertOptions) (*schema.NexlayerYAML, error) {
	content, err := os.ReadFile(composeFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Docker Compose file: %w", err)
	}

	var composeConfig DockerComposeConfig
	if err := yaml.Unmarshal(content, &composeConfig); err != nil {
		return nil, fmt.Errorf("failed to parse Docker Compose file: %w", err)
	}

	composeConfig.ConfigPath = composeFilePath

	// Setup variable context for substitution
	varCtx := vars.NewVariableContext()

	// Detect project URL if available (for template variables)
	if opts.ApplicationURL != "" {
		varCtx.SetURL(opts.ApplicationURL)
	}

	// Detect registry if available
	if opts.RegistryURL != "" {
		varCtx.SetRegistry(opts.RegistryURL)
	}

	nexlayerConfig := &schema.NexlayerYAML{
		Application: schema.Application{
			Name: opts.ApplicationName,
			Pods: make([]schema.Pod, 0, len(composeConfig.Services)),
		},
	}

	// Create a detector registry to help with project type detection
	registry := detection.NewDetectorRegistry()

	// Detect if we're running in an AI-powered IDE
	projectInfo, err := registry.DetectProject(".")
	if err == nil && projectInfo != nil && projectInfo.LLMProvider != "" {
		// Add AI-specific annotations to the application
		nexlayerConfig.Application.Annotations = map[string]string{
			"nexlayer.ai/llm-provider": projectInfo.LLMProvider,
			"nexlayer.ai/llm-model":    projectInfo.LLMModel,
		}
	}

	for serviceName, service := range composeConfig.Services {
		pod, err := convertServiceToPod(serviceName, service, composeConfig)
		if err != nil {
			log.Printf("Error converting service '%s': %v", serviceName, err)
			if !opts.ForceConversion {
				return nil, fmt.Errorf("failed to convert service '%s': %w", serviceName, err)
			}
			continue
		}
		nexlayerConfig.Application.Pods = append(nexlayerConfig.Application.Pods, *pod)

		// Add this pod to the variable context
		varCtx.AddPod(pod.Name)
	}

	// Process variable substitutions
	for i, pod := range nexlayerConfig.Application.Pods {
		// Process environment variables
		for j, envVar := range pod.Vars {
			if vars.ExtractVariables(envVar.Value) != nil {
				processed, err := vars.SubstituteVariables(envVar.Value, varCtx)
				if err == nil && processed != envVar.Value {
					nexlayerConfig.Application.Pods[i].Vars[j].Value = processed
				}
			}
		}

		// Process pod image for template variables
		if strings.Contains(pod.Image, "<%") {
			processed, err := vars.SubstituteVariables(pod.Image, varCtx)
			if err == nil && processed != pod.Image {
				nexlayerConfig.Application.Pods[i].Image = processed
			}
		}

		// Process volume paths
		for k, vol := range pod.Volumes {
			if strings.Contains(vol.Path, "<%") || strings.Contains(vol.Path, "${") {
				processed, err := vars.SubstituteVariables(vol.Path, varCtx)
				if err == nil && processed != vol.Path {
					nexlayerConfig.Application.Pods[i].Volumes[k].Path = processed
				}
			}
		}
	}

	// Process traditional pod references (maintaining backward compatibility)
	nexlayerConfig = addPodReferences(nexlayerConfig, composeConfig)
	nexlayerConfig = reorderPods(nexlayerConfig)

	// Validate the configuration
	if err := validateNexlayerConfig(nexlayerConfig); err != nil {
		if !opts.ForceConversion {
			return nil, fmt.Errorf("generated Nexlayer YAML is invalid: %w", err)
		}
		log.Printf("Warning: Generated Nexlayer YAML has validation errors: %v", err)
	}

	return nexlayerConfig, nil
}

// initializeDetectionManager creates a new detection manager with default tasks
func initializeDetectionManager() *detection.DetectionManager {
	// Create detector registry
	registry := detection.NewDetectorRegistry()

	// Create detection manager
	manager := detection.NewDetectionManager(registry)

	// Register default detection tasks
	manager.RegisterDefaultTasks()

	return manager
}

// initializeLLMEnricher creates a new LLM enricher for AI-powered analysis
func initializeLLMEnricher() *knowledge.LLMEnricher {
	// Check if LLM environment variables are set
	llmEnabled := os.Getenv("NEXLAYER_LLM_ENABLED")
	if llmEnabled == "false" {
		return nil
	}

	// Create knowledge graph
	graph := knowledge.NewGraph()

	// Get metadata directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Warning: Could not determine user home directory: %v", err)
		return nil
	}

	metadataDir := filepath.Join(homeDir, ".nexlayer", "metadata")

	// Create LLM enricher
	enricher := knowledge.NewLLMEnricher(graph, metadataDir)

	// Load metadata (non-blocking)
	go func() {
		if err := enricher.LoadMetadata(); err != nil {
			log.Printf("Warning: Could not load LLM metadata: %v", err)
		}
	}()

	return enricher
}

// applyAISuggestions applies AI suggestions to the Nexlayer configuration
func applyAISuggestions(config *schema.NexlayerYAML, result *ai.EnhancementResult) {
	// Add comments to the configuration
	if config.Comments == nil {
		config.Comments = make(map[string]string)
	}

	// Add field-specific comments
	for field, comment := range result.Comments {
		config.Comments[field] = comment
	}

	// Add general suggestions as application-level comments
	if len(result.Suggestions) > 0 {
		for i, suggestion := range result.Suggestions {
			config.Comments[fmt.Sprintf("suggestion.%d", i+1)] = suggestion
		}
	}

	// Auto-fix minor issues if possible
	for _, issue := range result.Issues {
		if issue.Type == "warning" {
			// Handle specific warnings that can be auto-fixed
			if strings.Contains(issue.Field, "servicePorts") && strings.Contains(issue.Message, "no name") {
				// Try to find the pod and port to fix
				parts := strings.Split(issue.Field, ".")
				if len(parts) >= 2 {
					podName := parts[1]
					for i, pod := range config.Application.Pods {
						if pod.Name == podName {
							for j, port := range pod.ServicePorts {
								if port.Name == "" {
									// Add a default name
									config.Application.Pods[i].ServicePorts[j].Name = fmt.Sprintf("%s-port-%d", podName, port.Port)
									break
								}
							}
							break
						}
					}
				}
			}
		}
	}
}

// printEnhancementSuggestions prints AI suggestions to the user
func printEnhancementSuggestions(result *ai.EnhancementResult) {
	if len(result.Issues) == 0 && len(result.Suggestions) == 0 {
		return
	}

	fmt.Println("\n🧠 AI Analysis Results:")

	// Print issues
	if len(result.Issues) > 0 {
		fmt.Println("\n⚠️ Potential Issues:")
		for _, issue := range result.Issues {
			var prefix string
			switch issue.Type {
			case "error":
				prefix = "❌ ERROR"
			case "warning":
				prefix = "⚠️ WARNING"
			default:
				prefix = "ℹ️ NOTE"
			}

			fmt.Printf("%s: %s\n", prefix, issue.Message)

			if len(issue.Suggestions) > 0 {
				fmt.Println("  Suggestions:")
				for _, suggestion := range issue.Suggestions {
					fmt.Printf("  - %s\n", suggestion)
				}
			}
			fmt.Println()
		}
	}

	// Print general suggestions
	if len(result.Suggestions) > 0 {
		fmt.Println("\n💡 Improvement Suggestions:")
		for _, suggestion := range result.Suggestions {
			fmt.Printf("- %s\n", suggestion)
		}
		fmt.Println()
	}
}

// hasCriticalIssues checks if the enhancement result contains critical issues
func hasCriticalIssues(result *ai.EnhancementResult) bool {
	for _, issue := range result.Issues {
		if issue.Type == "error" {
			return true
		}
	}
	return false
}

// convertServiceToPod converts a Docker Compose service to a Nexlayer pod
func convertServiceToPod(serviceName string, service DockerComposeService, composeConfig DockerComposeConfig) (*schema.Pod, error) {
	pod := &schema.Pod{
		Name:  serviceName,
		Type:  "docker",
		Image: service.Image,
	}

	// Set path for web services (case-insensitive)
	serviceNameLower := strings.ToLower(serviceName)
	if strings.Contains(serviceNameLower, "web") ||
		strings.Contains(serviceNameLower, "frontend") ||
		strings.Contains(serviceNameLower, "ui") {
		pod.Path = "/"
	}

	// Handle command
	if service.Command != nil {
		switch cmd := service.Command.(type) {
		case string:
			pod.Command = cmd
		case []interface{}:
			cmdParts := make([]string, 0, len(cmd))
			for _, part := range cmd {
				if strPart, ok := part.(string); ok {
					cmdParts = append(cmdParts, strPart)
				}
			}
			pod.Command = strings.Join(cmdParts, " ")
		}
	}

	// Handle entrypoint
	if service.Entrypoint != nil {
		switch entry := service.Entrypoint.(type) {
		case string:
			pod.Entrypoint = entry
		case []interface{}:
			entryParts := make([]string, 0, len(entry))
			for _, part := range entry {
				if strPart, ok := part.(string); ok {
					entryParts = append(entryParts, strPart)
				}
			}
			pod.Entrypoint = strings.Join(entryParts, " ")
		}
	}

	// Handle ports with intelligent defaults
	pod.ServicePorts = make([]schema.ServicePort, 0)
	if service.Ports != nil {
		switch ports := service.Ports.(type) {
		case []interface{}:
			for i, portDef := range ports {
				if portStr, ok := portDef.(string); ok {
					externalPort, internalPort, protocol, err := ParsePortMapping(portStr, serviceName)
					if err != nil {
						continue
					}
					pod.ServicePorts = append(pod.ServicePorts, schema.ServicePort{
						Name:       fmt.Sprintf("%s-port-%d", serviceName, i+1),
						Port:       externalPort,
						TargetPort: internalPort,
						Protocol:   protocol,
					})
				}
			}
		}
	}
	if len(pod.ServicePorts) == 0 {
		defaultPort := 80
		for img, port := range DefaultPorts {
			if strings.Contains(strings.ToLower(service.Image), img) {
				defaultPort = port
				break
			}
		}
		pod.ServicePorts = append(pod.ServicePorts, schema.ServicePort{
			Name:       fmt.Sprintf("%s-port-1", serviceName),
			Port:       defaultPort,
			TargetPort: defaultPort,
			Protocol:   "TCP",
		})
		log.Printf("Warning: No ports specified for service '%s', using default port %d", serviceName, defaultPort)
	}

	// Handle volumes with intelligent sizing
	pod.Volumes = make([]schema.Volume, 0)
	if service.Volumes != nil {
		switch volumes := service.Volumes.(type) {
		case []interface{}:
			for i, volumeDef := range volumes {
				if volumeStr, ok := volumeDef.(string); ok {
					volumeName, containerPath, readOnly, err := ParseVolumeMapping(volumeStr, serviceName)
					if err != nil {
						continue
					}

					// Determine appropriate size based on service type
					size := DefaultVolumeSizes["default"]
					for key, defaultSize := range DefaultVolumeSizes {
						if key != "default" && strings.Contains(strings.ToLower(service.Image), key) {
							size = defaultSize
							break
						}
					}

					// Use volume name directly if it's a named volume in compose file
					if _, ok := composeConfig.Volumes[volumeName]; ok {
						// Keep the volume name but make it more readable
						volumeName = strings.ReplaceAll(volumeName, "_", "-")
					} else {
						// For host paths, use a standard name based on the service
						volumeName = fmt.Sprintf("%s-data", serviceName)
						if i > 0 {
							volumeName = fmt.Sprintf("%s-%d", volumeName, i+1)
						}
					}

					pod.Volumes = append(pod.Volumes, schema.Volume{
						Name:     volumeName,
						Path:     containerPath,
						ReadOnly: readOnly,
						Size:     size,
					})
				}
			}
		}
	}

	// Handle environment variables
	pod.Vars = make([]schema.EnvVar, 0)
	if service.Environment != nil {
		switch env := service.Environment.(type) {
		case map[string]interface{}:
			for key, value := range env {
				if value != nil {
					pod.Vars = append(pod.Vars, schema.EnvVar{
						Key:   key,
						Value: fmt.Sprintf("%v", value),
					})
				}
			}
		case []interface{}:
			for _, item := range env {
				if strItem, ok := item.(string); ok {
					parts := strings.SplitN(strItem, "=", 2)
					if len(parts) == 2 {
						pod.Vars = append(pod.Vars, schema.EnvVar{
							Key:   parts[0],
							Value: parts[1],
						})
					}
				}
			}
		}
	}

	// Handle env_file
	if service.EnvFile != nil {
		envFiles := parseEnvFiles(service.EnvFile)
		for _, envFile := range envFiles {
			pod.Vars = append(pod.Vars, parseEnvFile(envFile)...)
		}
	}

	// Handle secrets
	if service.Secrets != nil {
		pod.Secrets = make([]schema.Secret, 0)
		for _, secretDef := range service.Secrets {
			var secretName string
			switch secret := secretDef.(type) {
			case string:
				secretName = secret
			case map[string]interface{}:
				if name, ok := secret["source"].(string); ok {
					secretName = name
				}
			}
			if secretName != "" {
				if secret, err := createSecret(secretName); err == nil {
					pod.Secrets = append(pod.Secrets, secret)
				}
			}
		}
	}

	return pod, nil
}

// parseEnvFiles extracts env file paths from various formats
func parseEnvFiles(envFilesDef interface{}) []string {
	envFiles := make([]string, 0)

	switch ef := envFilesDef.(type) {
	case string:
		envFiles = append(envFiles, ef)
	case []interface{}:
		for _, file := range ef {
			if strFile, ok := file.(string); ok {
				envFiles = append(envFiles, strFile)
			}
		}
	}

	return envFiles
}

// parseEnvFile reads and parses a .env file into environment variables
func parseEnvFile(filePath string) []schema.EnvVar {
	vars := make([]schema.EnvVar, 0)
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Warning: Failed to read env file '%s': %v", filePath, err)
		return vars
	}

	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// Remove surrounding quotes if present
			value = strings.Trim(value, `"'`)

			vars = append(vars, schema.EnvVar{
				Key:   key,
				Value: value,
			})
		}
	}

	return vars
}

// createSecret creates a Nexlayer secret from a Docker Compose secret
func createSecret(secretName string) (schema.Secret, error) {
	if secretName == "" {
		return schema.Secret{}, fmt.Errorf("secret name cannot be empty")
	}

	return schema.Secret{
		Name:     secretName,
		Path:     "/var/secrets/" + secretName,
		FileName: secretName + ".txt",
	}, nil
}

// addPodReferences modifies environment variables to use pod references - legacy method for compatibility
func addPodReferences(config *schema.NexlayerYAML, composeConfig DockerComposeConfig) *schema.NexlayerYAML {
	// Build service name to pod name map (legacy approach for compatibility)
	serviceMap := make(map[string]string, len(config.Application.Pods))
	for _, pod := range config.Application.Pods {
		serviceMap[pod.Name] = pod.Name + ".pod"
	}

	// Process each pod's environment variables for service references
	for i, pod := range config.Application.Pods {
		for j, envVar := range pod.Vars {
			// Look for service references in environment variables
			for serviceName, podRef := range serviceMap {
				// Replace serviceName:port with serviceName.pod:port using regex
				re := regexp.MustCompile(fmt.Sprintf(`\b%s:(\d+)\b`, regexp.QuoteMeta(serviceName)))
				envVar.Value = re.ReplaceAllString(envVar.Value, fmt.Sprintf("%s:$1", podRef))

				// Replace URLs like http://service:port with http://service.pod:port
				urlRe := regexp.MustCompile(fmt.Sprintf(`(://)%s:(\d+)`, regexp.QuoteMeta(serviceName)))
				envVar.Value = urlRe.ReplaceAllString(envVar.Value, fmt.Sprintf("$1%s:$2", podRef))

				// Replace standalone serviceName with pod reference
				if envVar.Value == serviceName {
					envVar.Value = podRef
				}
			}

			// Replace localhost or 127.0.0.1 with <% URL %>
			envVar.Value = strings.ReplaceAll(envVar.Value, "localhost", "<% URL %>")
			envVar.Value = strings.ReplaceAll(envVar.Value, "127.0.0.1", "<% URL %>")

			config.Application.Pods[i].Vars[j] = envVar
		}
	}

	return config
}

// validateNexlayerConfig performs validation on the generated Nexlayer YAML
func validateNexlayerConfig(config *schema.NexlayerYAML) error {
	if config.Application.Name == "" {
		return fmt.Errorf("application name is required")
	}

	if len(config.Application.Pods) == 0 {
		return fmt.Errorf("at least one pod is required")
	}

	for _, pod := range config.Application.Pods {
		if pod.Name == "" {
			return fmt.Errorf("pod name is required")
		}

		if pod.Image == "" {
			return fmt.Errorf("image is required for pod '%s'", pod.Name)
		}

		if len(pod.ServicePorts) == 0 {
			return fmt.Errorf("at least one service port is required for pod '%s'", pod.Name)
		}
	}

	return nil
}

// reorderPods reorders the pods to prioritize application pods over infrastructure pods
func reorderPods(config *schema.NexlayerYAML) *schema.NexlayerYAML {
	// Create a new slice for the reordered pods
	reorderedPods := make([]schema.Pod, 0, len(config.Application.Pods))

	// First, add all application pods (non-infrastructure)
	for _, pod := range config.Application.Pods {
		if !isInfrastructurePod(pod.Name) {
			reorderedPods = append(reorderedPods, pod)
		}
	}

	// Then, add all infrastructure pods
	for _, pod := range config.Application.Pods {
		if isInfrastructurePod(pod.Name) {
			reorderedPods = append(reorderedPods, pod)
		}
	}

	// Update the config with the reordered pods
	config.Application.Pods = reorderedPods

	return config
}

// isInfrastructurePod determines if a pod is an infrastructure pod
func isInfrastructurePod(podName string) bool {
	// List of common infrastructure service names (case insensitive check)
	podNameLower := strings.ToLower(podName)

	infraServices := []string{
		"postgres", "mysql", "mariadb", "mongodb", "mongo",
		"redis", "memcached", "rabbitmq", "kafka",
		"elasticsearch", "kibana", "logstash",
		"prometheus", "grafana", "jaeger",
		"minio", "s3", "clickhouse", "influxdb",
		"cassandra", "zookeeper", "etcd", "consul",
		"nginx", "traefik", "haproxy", "envoy",
	}

	// Check if the pod name contains any of the infrastructure service names
	for _, infraService := range infraServices {
		if strings.Contains(podNameLower, infraService) {
			return true
		}
	}

	return false
}

// DetectAndConvert tries to detect a Docker Compose file in the given directory
// and convert it to a Nexlayer YAML if found
func DetectAndConvert(dir string, appName string) (*schema.NexlayerYAML, error) {
	fmt.Printf("🔍 Searching for Docker Compose files in directory: %s\n", dir)

	opts := ConvertOptions{
		ProjectDir:      dir,
		ApplicationName: appName,
	}

	// List of common Compose file names in order of preference
	composeFiles := []string{
		"docker-compose.yml",
		"docker-compose.yaml",
		"docker-compose.dev.yml",
		"docker-compose.prod.yml",
		"compose.yml",
		"compose.yaml",
	}

	for _, fileName := range composeFiles {
		composePath := filepath.Join(dir, fileName)
		fmt.Printf("🔍 Checking for compose file at: %s\n", composePath)

		if _, err := os.Stat(composePath); err == nil {
			// Found a Docker Compose file, convert it
			fmt.Printf("✅ Found Docker Compose file: %s\n", composePath)

			config, err := Convert(composePath, opts)
			if err != nil {
				fmt.Printf("❌ Error converting Docker Compose file: %v\n", err)
				return nil, err
			}

			fmt.Printf("✅ Successfully converted Docker Compose file with %d services\n", len(config.Application.Pods))

			// Print summary of converted pods
			fmt.Printf("✅ Converted Docker Compose to Nexlayer YAML with %d pods:\n", len(config.Application.Pods))
			for i, pod := range config.Application.Pods {
				fmt.Printf("  - Pod %d: %s (image: %s)\n", i+1, pod.Name, pod.Image)
			}

			return config, nil
		} else {
			fmt.Printf("❌ Compose file not found at: %s (error: %v)\n", composePath, err)
		}
	}

	// No Docker Compose file found
	fmt.Printf("❌ No Docker Compose file found in %s\n", dir)
	return nil, fmt.Errorf("no Docker Compose file found in %s", dir)
}

// ConvertToNexlayer converts a Docker Compose configuration to a Nexlayer YAML configuration
func ConvertToNexlayer(composeConfig DockerComposeConfig) (*schema.NexlayerYAML, error) {
	// Create a detector registry to help with project type detection
	registry := detection.NewDetectorRegistry()

	// Initialize the Nexlayer YAML configuration
	config := &schema.NexlayerYAML{
		Application: schema.Application{
			Name: filepath.Base(filepath.Dir(composeConfig.ConfigPath)),
			Pods: make([]schema.Pod, 0),
		},
	}

	// Detect if we're running in an AI-powered IDE
	projectInfo, err := registry.DetectProject(".")
	if err == nil && projectInfo != nil && projectInfo.LLMProvider != "" {
		// We'll add AI-specific metadata to pods after they're created
		// This is handled in ConvertToNexlayer
	}

	// Convert services to pods
	for serviceName, service := range composeConfig.Services {
		pod, err := convertServiceToPod(serviceName, service, composeConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to convert service %s: %w", serviceName, err)
		}
		config.Application.Pods = append(config.Application.Pods, *pod)
	}

	// Sort pods to ensure deterministic output
	sortPods(config.Application.Pods)

	// Add pod references
	config = addPodReferences(config, composeConfig)

	return config, nil
}

// sortPods sorts the pods to ensure deterministic output
func sortPods(pods []schema.Pod) {
	sort.Slice(pods, func(i, j int) bool {
		return pods[i].Name < pods[j].Name
	})
}

func ConvertFromFile(composeFilePath string) (*schema.NexlayerYAML, error) {
	// Read and parse Docker Compose file
	data, err := os.ReadFile(composeFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read Docker Compose file: %w", err)
	}

	var composeConfig DockerComposeConfig
	if err := yaml.Unmarshal(data, &composeConfig); err != nil {
		return nil, fmt.Errorf("failed to parse Docker Compose file: %w", err)
	}

	composeConfig.ConfigPath = composeFilePath

	// Create a detector registry to help with project type detection
	registry := detection.NewDetectorRegistry()

	// Initialize the Nexlayer YAML configuration
	nexlayerConfig := &schema.NexlayerYAML{
		Application: schema.Application{
			Name: filepath.Base(filepath.Dir(composeFilePath)),
			Pods: make([]schema.Pod, 0),
		},
	}

	// Detect if we're running in an AI-powered IDE
	projectInfo, err := registry.DetectProject(".")
	if err == nil && projectInfo != nil && projectInfo.LLMProvider != "" {
		// We'll add AI-specific metadata to pods after they're created
		// This is handled in ConvertToNexlayer
	}

	// Convert services to pods
	for serviceName, service := range composeConfig.Services {
		pod, err := convertServiceToPod(serviceName, service, composeConfig)
		if err != nil {
			return nil, fmt.Errorf("failed to convert service %s: %w", serviceName, err)
		}
		nexlayerConfig.Application.Pods = append(nexlayerConfig.Application.Pods, *pod)
	}

	// Sort pods to ensure deterministic output
	sortPods(nexlayerConfig.Application.Pods)

	// Add pod references
	nexlayerConfig = addPodReferences(nexlayerConfig, composeConfig)

	return nexlayerConfig, nil
}
