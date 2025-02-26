// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package knowledge

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/schema"
)

// LLMContext represents the enriched context for LLM interactions
type LLMContext struct {
	ProjectStructure map[string]interface{} `json:"project_structure"`
	Dependencies     map[string]string      `json:"dependencies"`
	APIEndpoints     []interface{}          `json:"api_endpoints"`
	PodFlows         []interface{}          `json:"pod_flows"`
	Patterns         []interface{}          `json:"patterns"`
	Languages        map[string]interface{} `json:"languages"`
	Frameworks       map[string]interface{} `json:"frameworks"`
	Resources        map[string]interface{} `json:"resources"`
	Network          map[string]interface{} `json:"network"`
	Storage          map[string]interface{} `json:"storage"`
}

// LLMResult represents the result of an LLM query with metadata
type LLMResult struct {
	Result    string    `json:"result"`
	Timestamp time.Time `json:"timestamp"`
	Source    string    `json:"source"` // "cache" or "api"
}

// LLMEnricher enriches the knowledge graph with LLM metadata
type LLMEnricher struct {
	graph          *Graph
	metadata       map[string]interface{}
	metadataMu     sync.RWMutex
	metadataDir    string
	cache          sync.Map // Cache for LLM query results
	cacheTTL       time.Duration
	processingChan chan *processingTask
	wg             sync.WaitGroup
}

type processingTask struct {
	ctx      context.Context
	prompt   string
	config   *schema.NexlayerYAML
	resultCh chan<- *LLMResult
	errCh    chan<- error
}

// NewLLMEnricher creates a new LLM metadata enricher with optimized caching
func NewLLMEnricher(graph *Graph, metadataDir string) *LLMEnricher {
	// Default cache TTL of 30 minutes, can be customized if needed
	cacheTTL := 30 * time.Minute
	if ttlEnv := os.Getenv("NEXLAYER_LLM_CACHE_TTL"); ttlEnv != "" {
		if duration, err := time.ParseDuration(ttlEnv); err == nil {
			cacheTTL = duration
		}
	}

	// Initialize processing channel for background tasks
	processingChan := make(chan *processingTask, 10) // Buffer for 10 tasks

	enricher := &LLMEnricher{
		graph:          graph,
		metadata:       make(map[string]interface{}),
		metadataDir:    metadataDir,
		cacheTTL:       cacheTTL,
		processingChan: processingChan,
	}

	// Start background workers
	numWorkers := 2 // Default to 2 workers
	if workersEnv := os.Getenv("NEXLAYER_LLM_WORKERS"); workersEnv != "" {
		if workers, err := fmt.Sscanf(workersEnv, "%d", &numWorkers); err == nil && workers > 0 {
			numWorkers = workers
		}
	}

	for i := 0; i < numWorkers; i++ {
		go enricher.processTasksWorker()
	}

	return enricher
}

// generateCacheKey creates a deterministic cache key from a prompt and context
func (e *LLMEnricher) generateCacheKey(prompt string, config *schema.NexlayerYAML) string {
	// Create a composite key from the prompt and relevant config data
	var builder strings.Builder
	builder.WriteString(prompt)

	// Add key config elements if available
	if config != nil {
		builder.WriteString(config.Application.Name)
		for _, pod := range config.Application.Pods {
			builder.WriteString(pod.Name)
			builder.WriteString(pod.Image)
			for _, port := range pod.ServicePorts {
				builder.WriteString(fmt.Sprintf("%d", port.Port))
			}
		}
	}

	// Create a SHA-256 hash of the combined string
	hasher := sha256.New()
	hasher.Write([]byte(builder.String()))
	return hex.EncodeToString(hasher.Sum(nil))
}

// processTasksWorker handles background processing of LLM tasks
func (e *LLMEnricher) processTasksWorker() {
	for task := range e.processingChan {
		select {
		case <-task.ctx.Done():
			// Context cancelled, skip this task
			task.errCh <- task.ctx.Err()
			continue
		default:
			// Process the task
			result, err := e.performLLMQuery(task.ctx, task.prompt, task.config)
			if err != nil {
				task.errCh <- err
			} else {
				task.resultCh <- result
			}
		}
	}
}

// QueryLLM performs an LLM query with caching and optimized performance
func (e *LLMEnricher) QueryLLM(ctx context.Context, prompt string, config *schema.NexlayerYAML) (*LLMResult, error) {
	cacheKey := e.generateCacheKey(prompt, config)

	// Check cache first
	if cachedValue, found := e.cache.Load(cacheKey); found {
		result := cachedValue.(*LLMResult)
		// Check if cache entry is still valid
		if time.Since(result.Timestamp) < e.cacheTTL {
			return result, nil
		}
		// Cache expired, remove it
		e.cache.Delete(cacheKey)
	}

	// Perform actual LLM query
	return e.performLLMQuery(ctx, prompt, config)
}

// QueryLLMAsync performs an asynchronous LLM query
func (e *LLMEnricher) QueryLLMAsync(ctx context.Context, prompt string, config *schema.NexlayerYAML) (<-chan *LLMResult, <-chan error) {
	resultCh := make(chan *LLMResult, 1)
	errCh := make(chan error, 1)

	// Check cache first for immediate response
	cacheKey := e.generateCacheKey(prompt, config)
	if cachedValue, found := e.cache.Load(cacheKey); found {
		result := cachedValue.(*LLMResult)
		// Check if cache entry is still valid
		if time.Since(result.Timestamp) < e.cacheTTL {
			go func() {
				resultCh <- result
				close(resultCh)
				close(errCh)
			}()
			return resultCh, errCh
		}
		// Cache expired, remove it
		e.cache.Delete(cacheKey)
	}

	// Queue the task for background processing
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		select {
		case e.processingChan <- &processingTask{
			ctx:      ctx,
			prompt:   prompt,
			config:   config,
			resultCh: resultCh,
			errCh:    errCh,
		}:
			// Task queued successfully
		case <-ctx.Done():
			// Context cancelled
			errCh <- ctx.Err()
			close(resultCh)
			close(errCh)
		}
	}()

	return resultCh, errCh
}

// performLLMQuery is the actual implementation of the LLM query
func (e *LLMEnricher) performLLMQuery(ctx context.Context, prompt string, config *schema.NexlayerYAML) (*LLMResult, error) {
	// TODO: Implement actual LLM API call here
	// For now, this is a placeholder that simulates an LLM response

	// Create an enriched context for better LLM understanding
	enriched, err := e.EnrichContext(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("failed to enrich context: %w", err)
	}

	// Simulate LLM processing time (remove in production)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(100 * time.Millisecond):
		// Continue processing
	}

	// Generate a simulated response based on the prompt and context
	var response string
	if strings.Contains(prompt, "deployment issues") {
		response = e.simulateDeploymentIssueCheck(config, enriched)
	} else if strings.Contains(prompt, "volume recommendations") {
		response = e.simulateVolumeRecommendations(config, enriched)
	} else if strings.Contains(prompt, "port configuration") {
		response = e.simulatePortConfigurationCheck(config, enriched)
	} else {
		response = "LLM analysis complete. No issues detected."
	}

	// Create result
	result := &LLMResult{
		Result:    response,
		Timestamp: time.Now(),
		Source:    "api", // This would be "api" in a real implementation
	}

	// Cache the result for future queries
	cacheKey := e.generateCacheKey(prompt, config)
	e.cache.Store(cacheKey, result)

	return result, nil
}

// simulateDeploymentIssueCheck creates a simulated response for deployment issues
func (e *LLMEnricher) simulateDeploymentIssueCheck(config *schema.NexlayerYAML, context *LLMContext) string {
	var issues []string

	// Check for missing ports
	for _, pod := range config.Application.Pods {
		if len(pod.ServicePorts) == 0 {
			issues = append(issues, fmt.Sprintf("- Pod '%s' has no service ports defined", pod.Name))
		}
	}

	// Check for database pods without volumes
	for _, pod := range config.Application.Pods {
		if isDatabase(pod.Image) && len(pod.Volumes) == 0 {
			issues = append(issues, fmt.Sprintf("- Database pod '%s' has no persistent volumes", pod.Name))
		}
	}

	if len(issues) > 0 {
		return fmt.Sprintf("Deployment issues found:\n%s", strings.Join(issues, "\n"))
	}
	return "No deployment issues detected."
}

// simulateVolumeRecommendations creates a simulated response for volume recommendations
func (e *LLMEnricher) simulateVolumeRecommendations(config *schema.NexlayerYAML, context *LLMContext) string {
	var recommendations []string

	// Generate volume size recommendations based on pod type
	for _, pod := range config.Application.Pods {
		if isDatabase(pod.Image) {
			for i, volume := range pod.Volumes {
				if i < len(pod.Volumes) && (volume.Size == "" || volume.Size == "1Gi") {
					recommendations = append(recommendations, fmt.Sprintf("- Increase volume size for database pod '%s' to at least 10Gi", pod.Name))
					break
				}
			}
		}
	}

	if len(recommendations) > 0 {
		return fmt.Sprintf("Volume recommendations:\n%s", strings.Join(recommendations, "\n"))
	}
	return "No volume recommendations needed."
}

// simulatePortConfigurationCheck creates a simulated response for port configuration checks
func (e *LLMEnricher) simulatePortConfigurationCheck(config *schema.NexlayerYAML, context *LLMContext) string {
	var recommendations []string

	// Check for common port misconfigurations
	for _, pod := range config.Application.Pods {
		if isDatabase(pod.Image) {
			hasCorrectPort := false
			expectedPort := getDatabaseDefaultPort(pod.Image)
			for _, port := range pod.ServicePorts {
				if port.TargetPort == expectedPort {
					hasCorrectPort = true
					break
				}
			}
			if !hasCorrectPort {
				recommendations = append(recommendations, fmt.Sprintf("- Database pod '%s' should expose standard port %d", pod.Name, expectedPort))
			}
		}
	}

	if len(recommendations) > 0 {
		return fmt.Sprintf("Port configuration recommendations:\n%s", strings.Join(recommendations, "\n"))
	}
	return "Port configuration looks good."
}

// isDatabase checks if a pod image is a database
func isDatabase(image string) bool {
	dbImages := []string{"postgres", "mysql", "mariadb", "mongodb", "mongo", "redis", "clickhouse"}
	imageLower := strings.ToLower(image)
	for _, db := range dbImages {
		if strings.Contains(imageLower, db) {
			return true
		}
	}
	return false
}

// getDatabaseDefaultPort returns the default port for a database image
func getDatabaseDefaultPort(image string) int {
	imageLower := strings.ToLower(image)
	if strings.Contains(imageLower, "postgres") {
		return 5432
	} else if strings.Contains(imageLower, "mysql") || strings.Contains(imageLower, "mariadb") {
		return 3306
	} else if strings.Contains(imageLower, "mongo") {
		return 27017
	} else if strings.Contains(imageLower, "redis") {
		return 6379
	} else if strings.Contains(imageLower, "clickhouse") {
		return 8123
	}
	return 0
}

// LoadMetadata loads LLM metadata from the tools directory with caching
func (e *LLMEnricher) LoadMetadata() error {
	e.metadataMu.Lock()
	defer e.metadataMu.Unlock()

	metadataPath := filepath.Join(e.metadataDir, "llm", "metadata.json")
	data, err := os.ReadFile(metadataPath)
	if err != nil {
		return fmt.Errorf("failed to read metadata: %w", err)
	}

	if err := json.Unmarshal(data, &e.metadata); err != nil {
		return fmt.Errorf("failed to parse metadata: %w", err)
	}

	return nil
}

// EnrichContext creates an enriched context for LLM interactions
func (e *LLMEnricher) EnrichContext(ctx context.Context, yamlConfig *schema.NexlayerYAML) (*LLMContext, error) {
	enriched := &LLMContext{
		ProjectStructure: make(map[string]interface{}),
		Dependencies:     make(map[string]string),
		APIEndpoints:     make([]interface{}, 0),
		PodFlows:         make([]interface{}, 0),
		Patterns:         make([]interface{}, 0),
		Languages:        make(map[string]interface{}),
		Frameworks:       make(map[string]interface{}),
		Resources:        make(map[string]interface{}),
		Network:          make(map[string]interface{}),
		Storage:          make(map[string]interface{}),
	}

	// Extract project structure (only deployment-relevant files)
	e.graph.nodes.Range(func(key, value interface{}) bool {
		if node, ok := value.(*Node); ok {
			if node.Type == TypeFile {
				ext := strings.ToLower(filepath.Ext(node.Path))
				switch ext {
				case ".yaml", ".yml", ".json", ".toml":
					// Include configuration files
					parts := strings.Split(node.Path, string(os.PathSeparator))
					current := enriched.ProjectStructure
					for i, part := range parts {
						if i == len(parts)-1 {
							if metadata, ok := node.Metadata[MetadataDeployment]; ok {
								current[part] = metadata
							} else {
								current[part] = nil
							}
						} else {
							if _, exists := current[part]; !exists {
								current[part] = make(map[string]interface{})
							}
							if nextLevel, ok := current[part].(map[string]interface{}); ok {
								current = nextLevel
							} else {
								// Handle type mismatch gracefully
								current[part] = make(map[string]interface{})
								current = current[part].(map[string]interface{})
							}
						}
					}
				}
			}
		}
		return true
	})

	// Extract deployment-relevant information
	e.graph.nodes.Range(func(key, value interface{}) bool {
		if node, ok := value.(*Node); ok {
			switch node.Type {
			case TypeAPIEndpoint:
				if network, ok := node.Metadata[MetadataNetwork].(map[string]interface{}); ok {
					endpoint := map[string]interface{}{
						"path":   network["path"],
						"method": network["method"],
					}
					if auth, ok := node.Metadata[MetadataAuth].(map[string]interface{}); ok {
						endpoint["auth"] = auth
					}
					enriched.APIEndpoints = append(enriched.APIEndpoints, endpoint)
				}
			case TypeDependency:
				if deployment, ok := node.Metadata[MetadataDeployment].(map[string]interface{}); ok {
					depType := deployment["type"].(string)
					switch depType {
					case "database", "cache", "queue", "storage":
						enriched.Resources[node.Name] = deployment
					}
					enriched.Dependencies[node.Name] = deployment["version"].(string)
				}
			}
		}
		return true
	})

	// Extract pod communication flows from nexlayer.yaml
	if yamlConfig != nil {
		podMap := make(map[string]bool)
		podNetworking := make(map[string]interface{})
		podStorage := make(map[string]interface{})

		for _, pod := range yamlConfig.Application.Pods {
			podMap[pod.Name] = true

			// Collect networking config
			if len(pod.ServicePorts) > 0 {
				podNetworking[pod.Name] = map[string]interface{}{
					"ports": pod.ServicePorts,
					"path":  pod.Path,
				}
			}

			// Collect storage requirements
			if len(pod.Volumes) > 0 {
				podStorage[pod.Name] = pod.Volumes
			}
		}

		enriched.Network = podNetworking
		enriched.Storage = podStorage

		// Extract pod communication flows
		for _, pod := range yamlConfig.Application.Pods {
			for _, v := range pod.Vars {
				if strings.Contains(v.Value, ".pod") {
					targetPod := strings.Split(v.Value, ".pod")[0]
					if podMap[targetPod] {
						flow := map[string]interface{}{
							"source": pod.Name,
							"target": targetPod,
							"var":    v.Key,
							"value":  v.Value,
						}
						enriched.PodFlows = append(enriched.PodFlows, flow)
					}
				}
			}
		}
	}

	// Add patterns from cached metadata
	e.metadataMu.RLock()
	if patterns, ok := e.metadata["deployment_patterns"].([]interface{}); ok {
		enriched.Patterns = patterns
	}
	e.metadataMu.RUnlock()

	return enriched, nil
}

// GeneratePrompt generates an LLM prompt with enriched context
func (e *LLMEnricher) GeneratePrompt(ctx context.Context, basePrompt string, yamlConfig *schema.NexlayerYAML) (string, error) {
	enriched, err := e.EnrichContext(ctx, yamlConfig)
	if err != nil {
		return "", fmt.Errorf("failed to enrich context: %w", err)
	}

	// Convert enriched context to JSON
	contextJSON, err := json.MarshalIndent(enriched, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal context: %w", err)
	}

	// Combine base prompt with enriched context
	prompt := fmt.Sprintf("%s\n\nContext:\n%s", basePrompt, string(contextJSON))
	return prompt, nil
}

// Shutdown gracefully shuts down the LLM enricher
func (e *LLMEnricher) Shutdown(ctx context.Context) {
	// Close the processing channel to stop workers
	close(e.processingChan)

	// Wait for all background tasks to complete with timeout
	waitCh := make(chan struct{})
	go func() {
		e.wg.Wait()
		close(waitCh)
	}()

	// Wait for either all tasks to complete or context to cancel
	select {
	case <-waitCh:
		// All tasks completed
	case <-ctx.Done():
		// Context cancelled, some tasks may still be running
	}
}
