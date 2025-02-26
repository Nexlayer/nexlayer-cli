// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

// Package ai provides AI-powered enhancements for Nexlayer configurations
package ai

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/Nexlayer/nexlayer-cli/pkg/core/schema"
	"github.com/Nexlayer/nexlayer-cli/pkg/detection"
	"github.com/Nexlayer/nexlayer-cli/pkg/knowledge"
)

// EnhancementIssue represents a detected issue in the configuration
type EnhancementIssue struct {
	Type        string   // Type of issue: "error", "warning", "suggestion"
	Field       string   // The field with the issue
	Message     string   // Description of the issue
	Suggestions []string // Suggestions to fix the issue
}

// EnhancementResult contains the results of AI enhancement
type EnhancementResult struct {
	Issues          []EnhancementIssue         // Issues detected in the configuration
	Suggestions     []string                   // General suggestions for improvement
	ImprovedConfig  *schema.NexlayerYAML       // Enhanced configuration (if any)
	Comments        map[string]string          // Comments for specific fields
	DetectionResult *detection.DetectionResult // Detection results used for enhancement
	EnhancementTime time.Duration              // How long enhancement took
}

// Enhancer provides AI-powered enhancements for Nexlayer configurations
type Enhancer struct {
	llmEnricher     *knowledge.LLMEnricher
	detectionMgr    *detection.DetectionManager
	cache           sync.Map             // Cache for enhancement results
	cacheTTL        time.Duration        // How long to cache results
	analysisTimeout time.Duration        // Timeout for analysis operations
	enhancementChan chan enhancementTask // Channel for background enhancements
	wg              sync.WaitGroup       // WaitGroup for background tasks
}

type enhancementTask struct {
	ctx          context.Context
	config       *schema.NexlayerYAML
	detectionDir string
	resultCh     chan<- *EnhancementResult
	errCh        chan<- error
}

// NewEnhancer creates a new AI Enhancer with the provided LLM enricher
func NewEnhancer(llmEnricher *knowledge.LLMEnricher, detectionMgr *detection.DetectionManager) *Enhancer {
	// Create a channel for background enhancement tasks
	enhancementChan := make(chan enhancementTask, 5) // Buffer for 5 tasks

	e := &Enhancer{
		llmEnricher:     llmEnricher,
		detectionMgr:    detectionMgr,
		cacheTTL:        10 * time.Minute, // Cache results for 10 minutes
		analysisTimeout: 30 * time.Second, // 30 second timeout for analysis
		enhancementChan: enhancementChan,
	}

	// Start worker goroutines
	numWorkers := 2
	for i := 0; i < numWorkers; i++ {
		go e.enhancementWorker()
	}

	return e
}

// enhancementWorker processes enhancement tasks in the background
func (e *Enhancer) enhancementWorker() {
	for task := range e.enhancementChan {
		select {
		case <-task.ctx.Done():
			// Context cancelled
			task.errCh <- task.ctx.Err()
			continue
		default:
			// Process the task
			result, err := e.performEnhancement(task.ctx, task.config, task.detectionDir)
			if err != nil {
				task.errCh <- err
			} else {
				task.resultCh <- result
			}
		}
	}
}

// EnhanceAsync enhances a Nexlayer YAML configuration asynchronously using AI
func (e *Enhancer) EnhanceAsync(ctx context.Context, config *schema.NexlayerYAML, detectionDir string) (<-chan *EnhancementResult, <-chan error) {
	resultCh := make(chan *EnhancementResult, 1)
	errCh := make(chan error, 1)

	// Generate a cache key from the configuration
	cacheKey := generateCacheKey(config)

	// Check if we have a cached result
	if cachedValue, found := e.cache.Load(cacheKey); found {
		cachedResult := cachedValue.(*EnhancementResult)
		// Check if cache is still valid
		if time.Since(cachedResult.DetectionResult.CompletedAt) < e.cacheTTL {
			// Return cached result
			go func() {
				resultCh <- cachedResult
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
		case e.enhancementChan <- enhancementTask{
			ctx:          ctx,
			config:       config,
			detectionDir: detectionDir,
			resultCh:     resultCh,
			errCh:        errCh,
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

// Enhance enhances a Nexlayer YAML configuration using AI (synchronous version)
func (e *Enhancer) Enhance(ctx context.Context, config *schema.NexlayerYAML, detectionDir string) (*EnhancementResult, error) {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(ctx, e.analysisTimeout)
	defer cancel()

	resultCh, errCh := e.EnhanceAsync(ctx, config, detectionDir)

	// Wait for either a result or an error
	select {
	case result := <-resultCh:
		return result, nil
	case err := <-errCh:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// performEnhancement performs the actual enhancement of a configuration
func (e *Enhancer) performEnhancement(ctx context.Context, config *schema.NexlayerYAML, detectionDir string) (*EnhancementResult, error) {
	startTime := time.Now()

	// Get detection results to inform enhancements
	var detectionResult *detection.DetectionResult
	if e.detectionMgr != nil {
		detectionCh := e.detectionMgr.DetectAsync(ctx, detectionDir)
		detectionResult = <-detectionCh
	}

	// Initial enhancement result
	result := &EnhancementResult{
		Issues:          make([]EnhancementIssue, 0),
		Suggestions:     make([]string, 0),
		ImprovedConfig:  config,
		Comments:        make(map[string]string),
		DetectionResult: detectionResult,
	}

	// If we have an LLM enricher, use it for advanced analysis
	if e.llmEnricher != nil {
		// Analyze deployment issues
		deploymentIssues, err := e.analyzeDeploymentIssues(ctx, config)
		if err != nil {
			log.Printf("Warning: AI deployment analysis failed: %v", err)
		} else {
			for _, issue := range deploymentIssues {
				result.Issues = append(result.Issues, issue)
			}
		}

		// Generate volume size recommendations
		volumeRecs, err := e.analyzeVolumeRecommendations(ctx, config)
		if err != nil {
			log.Printf("Warning: AI volume analysis failed: %v", err)
		} else {
			for _, rec := range volumeRecs {
				result.Suggestions = append(result.Suggestions, rec)
			}
		}

		// Generate port configuration recommendations
		portRecs, err := e.analyzePortConfigurations(ctx, config)
		if err != nil {
			log.Printf("Warning: AI port analysis failed: %v", err)
		} else {
			for _, rec := range portRecs {
				result.Suggestions = append(result.Suggestions, rec)
			}
		}

		// Generate network comments
		if detectionResult != nil && detectionResult.ProjectInfo != nil {
			networkComments := e.generateNetworkComments(config, detectionResult)
			for field, comment := range networkComments {
				result.Comments[field] = comment
			}
		}
	} else {
		// Without LLM, perform basic analysis
		result = e.performBasicAnalysis(config)
	}

	// Record enhancement duration
	result.EnhancementTime = time.Since(startTime)

	// Cache the result
	cacheKey := generateCacheKey(config)
	e.cache.Store(cacheKey, result)

	return result, nil
}

// analyzeDeploymentIssues uses LLM to identify deployment issues
func (e *Enhancer) analyzeDeploymentIssues(ctx context.Context, config *schema.NexlayerYAML) ([]EnhancementIssue, error) {
	prompt := "Analyze the following Nexlayer configuration for deployment issues. " +
		"Focus on missing ports, improper volume configurations, and pod dependencies."

	llmResult, err := e.llmEnricher.QueryLLM(ctx, prompt, config)
	if err != nil {
		return nil, fmt.Errorf("LLM query failed: %w", err)
	}

	// Parse the LLM result into structured issues
	issues := parseIssuesFromLLMResponse(llmResult.Result)
	return issues, nil
}

// analyzeVolumeRecommendations uses LLM to generate volume recommendations
func (e *Enhancer) analyzeVolumeRecommendations(ctx context.Context, config *schema.NexlayerYAML) ([]string, error) {
	prompt := "Recommend optimal volume sizes for the pods in this Nexlayer configuration. " +
		"Consider the type of service, expected data growth, and performance needs."

	llmResult, err := e.llmEnricher.QueryLLM(ctx, prompt, config)
	if err != nil {
		return nil, fmt.Errorf("LLM query failed: %w", err)
	}

	// Extract recommendations from the LLM response
	recommendations := parseRecommendationsFromLLMResponse(llmResult.Result)
	return recommendations, nil
}

// analyzePortConfigurations uses LLM to analyze port configurations
func (e *Enhancer) analyzePortConfigurations(ctx context.Context, config *schema.NexlayerYAML) ([]string, error) {
	prompt := "Analyze the port configurations for the pods in this Nexlayer configuration. " +
		"Recommend standard ports for services and identify any potential port conflicts."

	llmResult, err := e.llmEnricher.QueryLLM(ctx, prompt, config)
	if err != nil {
		return nil, fmt.Errorf("LLM query failed: %w", err)
	}

	// Extract recommendations from the LLM response
	recommendations := parseRecommendationsFromLLMResponse(llmResult.Result)
	return recommendations, nil
}

// generateNetworkComments generates comments for network-related fields
func (e *Enhancer) generateNetworkComments(config *schema.NexlayerYAML, detectionResult *detection.DetectionResult) map[string]string {
	comments := make(map[string]string)

	// Add comments for pod dependencies
	podDependencies := make(map[string][]string)
	for _, pod := range config.Application.Pods {
		for _, envVar := range pod.Vars {
			for _, otherPod := range config.Application.Pods {
				if otherPod.Name != pod.Name && strings.Contains(envVar.Value, otherPod.Name+".pod") {
					podDependencies[pod.Name] = append(podDependencies[pod.Name], otherPod.Name)
					break
				}
			}
		}
	}

	// Add comments about pod dependencies
	for podName, deps := range podDependencies {
		if len(deps) > 0 {
			comment := fmt.Sprintf("Pod %s depends on: %s", podName, strings.Join(deps, ", "))
			comments[fmt.Sprintf("pods.%s.dependencies", podName)] = comment
		}
	}

	// Add comments about database pods
	for _, pod := range config.Application.Pods {
		if isDatabase(pod.Image) {
			// Check if volumes are properly configured
			if len(pod.Volumes) == 0 {
				comments[fmt.Sprintf("pods.%s.volumes", pod.Name)] = "Database pod should have persistent storage configured"
			} else {
				// Check volume size
				for i, volume := range pod.Volumes {
					if volume.Size == "" || volume.Size == "1Gi" {
						comments[fmt.Sprintf("pods.%s.volumes[%d].size", pod.Name, i)] = "Consider increasing volume size for database pod"
					}
				}
			}
		}
	}

	return comments
}

// performBasicAnalysis performs basic analysis without LLM
func (e *Enhancer) performBasicAnalysis(config *schema.NexlayerYAML) *EnhancementResult {
	result := &EnhancementResult{
		Issues:      make([]EnhancementIssue, 0),
		Suggestions: make([]string, 0),
		Comments:    make(map[string]string),
	}

	// Check for missing ports
	for _, pod := range config.Application.Pods {
		if len(pod.ServicePorts) == 0 {
			result.Issues = append(result.Issues, EnhancementIssue{
				Type:    "error",
				Field:   fmt.Sprintf("pods.%s.servicePorts", pod.Name),
				Message: fmt.Sprintf("Pod '%s' has no service ports defined", pod.Name),
				Suggestions: []string{
					fmt.Sprintf("Add appropriate service ports for pod '%s'", pod.Name),
				},
			})
		}
	}

	// Check database pods for volumes
	for _, pod := range config.Application.Pods {
		if isDatabase(pod.Image) && len(pod.Volumes) == 0 {
			result.Issues = append(result.Issues, EnhancementIssue{
				Type:    "warning",
				Field:   fmt.Sprintf("pods.%s.volumes", pod.Name),
				Message: fmt.Sprintf("Database pod '%s' has no persistent volumes", pod.Name),
				Suggestions: []string{
					fmt.Sprintf("Add a persistent volume for database pod '%s'", pod.Name),
					fmt.Sprintf("Recommended size: 10Gi for production databases"),
				},
			})
		}
	}

	// Check for proper service port naming
	for _, pod := range config.Application.Pods {
		for i, port := range pod.ServicePorts {
			if port.Name == "" {
				result.Issues = append(result.Issues, EnhancementIssue{
					Type:    "warning",
					Field:   fmt.Sprintf("pods.%s.servicePorts[%d].name", pod.Name, i),
					Message: fmt.Sprintf("Service port #%d on pod '%s' has no name", i, pod.Name),
					Suggestions: []string{
						fmt.Sprintf("Add a descriptive name for the port (e.g., '%s-port-%d')", pod.Name, port.Port),
					},
				})
			}
		}
	}

	return result
}

// parseIssuesFromLLMResponse parses issues from an LLM response
func parseIssuesFromLLMResponse(response string) []EnhancementIssue {
	issues := make([]EnhancementIssue, 0)

	// Simple parsing: look for lines that start with "-" or "*"
	for _, line := range strings.Split(response, "\n") {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			// Remove the prefix
			line = strings.TrimPrefix(line, "- ")
			line = strings.TrimPrefix(line, "* ")

			// Create an issue
			issue := EnhancementIssue{
				Type:        "suggestion",
				Message:     line,
				Suggestions: []string{},
			}

			// Extract field if it's in the format "Field: message"
			if parts := strings.SplitN(line, ":", 2); len(parts) == 2 {
				fieldName := strings.TrimSpace(parts[0])
				message := strings.TrimSpace(parts[1])
				if fieldName != "" && message != "" {
					issue.Field = fieldName
					issue.Message = message
				}
			}

			issues = append(issues, issue)
		}
	}

	return issues
}

// parseRecommendationsFromLLMResponse parses recommendations from an LLM response
func parseRecommendationsFromLLMResponse(response string) []string {
	recommendations := make([]string, 0)

	// Split by lines and filter for actionable recommendations
	for _, line := range strings.Split(response, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Check if line starts with a recommendation marker
		if strings.HasPrefix(line, "- ") || strings.HasPrefix(line, "* ") {
			// Remove the prefix
			line = strings.TrimPrefix(line, "- ")
			line = strings.TrimPrefix(line, "* ")
			recommendations = append(recommendations, line)
		} else if !strings.HasPrefix(line, "#") && !strings.HasPrefix(line, "//") {
			// Include non-empty, non-comment lines
			recommendations = append(recommendations, line)
		}
	}

	return recommendations
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

// generateCacheKey generates a cache key for a configuration
func generateCacheKey(config *schema.NexlayerYAML) string {
	if config == nil {
		return "nil-config"
	}

	// Simple hash based on application name and number of pods
	key := fmt.Sprintf("%s-%d", config.Application.Name, len(config.Application.Pods))
	return key
}

// Shutdown gracefully shuts down the enhancer
func (e *Enhancer) Shutdown(ctx context.Context) {
	// Close the enhancement channel to stop workers
	close(e.enhancementChan)

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
