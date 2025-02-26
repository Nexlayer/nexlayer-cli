// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package detection

import (
	"context"
	"sync"
	"time"
)

// DetectionResult represents the result of multiple asynchronous detections
type DetectionResult struct {
	ProjectTypes map[string]float64     // Map of project type to confidence level
	Editors      map[string]float64     // Map of editors to confidence level
	ProjectInfo  *ProjectInfo           // Consolidated project info
	Metadata     map[string]interface{} // Additional metadata from detection
	CompletedAt  time.Time              // When the detection completed
	mutex        sync.RWMutex           // For thread-safe access
}

// DetectionTask represents a detection operation that can run asynchronously
type DetectionTask struct {
	Detector Detector
	Weight   float64
	Timeout  time.Duration
}

// DetectionManager handles asynchronous detection operations
type DetectionManager struct {
	registry      *DetectorRegistry
	tasks         []DetectionTask
	results       *DetectionResult
	confidenceMap map[string]float64 // Default confidence thresholds
	cacheTTL      time.Duration      // How long to cache results
	lastRun       time.Time
	mutex         sync.RWMutex
}

// NewDetectionManager creates a new detection manager with default settings
func NewDetectionManager(registry *DetectorRegistry) *DetectionManager {
	return &DetectionManager{
		registry: registry,
		tasks:    make([]DetectionTask, 0),
		results: &DetectionResult{
			ProjectTypes: make(map[string]float64),
			Editors:      make(map[string]float64),
			Metadata:     make(map[string]interface{}),
		},
		confidenceMap: map[string]float64{
			"editor":      0.9, // 90% confidence for editor detection
			"projectType": 0.7, // 70% confidence for project type
			"language":    0.8, // 80% confidence for language detection
		},
		cacheTTL: 5 * time.Minute, // Cache results for 5 minutes by default
	}
}

// RegisterDetectionTask adds a new detection task to the manager
func (dm *DetectionManager) RegisterDetectionTask(detector Detector, weight float64, timeout time.Duration) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.tasks = append(dm.tasks, DetectionTask{
		Detector: detector,
		Weight:   weight,
		Timeout:  timeout,
	})
}

// SetConfidenceThreshold sets the confidence threshold for a specific detection type
func (dm *DetectionManager) SetConfidenceThreshold(detectionType string, threshold float64) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.confidenceMap[detectionType] = threshold
}

// SetCacheTTL sets how long detection results should be cached
func (dm *DetectionManager) SetCacheTTL(ttl time.Duration) {
	dm.mutex.Lock()
	defer dm.mutex.Unlock()

	dm.cacheTTL = ttl
}

// GetConfidenceThreshold gets the confidence threshold for a detection type
func (dm *DetectionManager) GetConfidenceThreshold(detectionType string) float64 {
	dm.mutex.RLock()
	defer dm.mutex.RUnlock()

	if threshold, ok := dm.confidenceMap[detectionType]; ok {
		return threshold
	}
	return 0.7 // Default threshold
}

// DetectAsync runs all registered detection tasks asynchronously
func (dm *DetectionManager) DetectAsync(ctx context.Context, dir string) <-chan *DetectionResult {
	resultCh := make(chan *DetectionResult, 1)

	// Check if we have recent cached results
	dm.mutex.RLock()
	cacheValid := time.Since(dm.lastRun) < dm.cacheTTL
	dm.mutex.RUnlock()

	if cacheValid {
		// Return cached results
		go func() {
			dm.mutex.RLock()
			resultCopy := copyDetectionResult(dm.results)
			dm.mutex.RUnlock()
			resultCh <- resultCopy
			close(resultCh)
		}()
		return resultCh
	}

	// Start detection in a goroutine
	go func() {
		result := &DetectionResult{
			ProjectTypes: make(map[string]float64),
			Editors:      make(map[string]float64),
			Metadata:     make(map[string]interface{}),
		}

		var wg sync.WaitGroup
		resultsMutex := &sync.Mutex{}

		// Execute each detection task in its own goroutine
		for _, task := range dm.tasks {
			wg.Add(1)
			go func(t DetectionTask) {
				defer wg.Done()

				// Create a timeout context for this task
				taskCtx, cancel := context.WithTimeout(ctx, t.Timeout)
				defer cancel()

				// Run the detector
				info, err := t.Detector.Detect(taskCtx, dir)
				if err != nil || info == nil {
					return
				}

				// Safely update the result
				resultsMutex.Lock()
				defer resultsMutex.Unlock()

				// Update project type with confidence
				if info.Type != "" {
					// Apply the weight to the confidence
					confidence := info.Confidence * t.Weight
					if existing, ok := result.ProjectTypes[info.Type]; ok {
						// Take the max confidence if we already have this type
						if confidence > existing {
							result.ProjectTypes[info.Type] = confidence
						}
					} else {
						result.ProjectTypes[info.Type] = confidence
					}
				}

				// Update editor info
				if info.Editor != "" {
					confidence := info.Confidence * t.Weight
					if existing, ok := result.Editors[info.Editor]; ok {
						if confidence > existing {
							result.Editors[info.Editor] = confidence
						}
					} else {
						result.Editors[info.Editor] = confidence
					}
				}

				// Merge metadata
				for k, v := range info.Metadata {
					result.Metadata[k] = v
				}
			}(task)
		}

		// Wait for all tasks to complete
		wg.Wait()

		// Determine the best project info from the results
		bestType := ""
		bestTypeConfidence := 0.0
		for t, confidence := range result.ProjectTypes {
			if confidence > bestTypeConfidence {
				bestType = t
				bestTypeConfidence = confidence
			}
		}

		bestEditor := ""
		bestEditorConfidence := 0.0
		for e, confidence := range result.Editors {
			if confidence > bestEditorConfidence {
				bestEditor = e
				bestEditorConfidence = confidence
			}
		}

		// Only set project type and editor if confidence exceeds threshold
		projectInfo := &ProjectInfo{
			Metadata: result.Metadata,
		}

		if bestTypeConfidence >= dm.GetConfidenceThreshold("projectType") {
			projectInfo.Type = bestType
		}

		if bestEditorConfidence >= dm.GetConfidenceThreshold("editor") {
			projectInfo.Editor = bestEditor
		}

		// Include LLM information if available
		if llmProvider, ok := result.Metadata["llm_provider"].(string); ok {
			projectInfo.LLMProvider = llmProvider
		}
		if llmModel, ok := result.Metadata["llm_model"].(string); ok {
			projectInfo.LLMModel = llmModel
		}

		// Set the result completion time
		result.CompletedAt = time.Now()
		result.ProjectInfo = projectInfo

		// Cache the results
		dm.mutex.Lock()
		dm.results = result
		dm.lastRun = time.Now()
		dm.mutex.Unlock()

		// Return results to the caller
		resultCh <- copyDetectionResult(result)
		close(resultCh)
	}()

	return resultCh
}

// DetectProject performs detection and returns a project info
// This is a synchronous version of DetectAsync for backward compatibility
func (dm *DetectionManager) DetectProject(ctx context.Context, dir string) (*ProjectInfo, error) {
	resultCh := dm.DetectAsync(ctx, dir)
	result := <-resultCh

	if result != nil && result.ProjectInfo != nil {
		return result.ProjectInfo, nil
	}
	return nil, nil
}

// copyDetectionResult creates a deep copy of a detection result to avoid race conditions
func copyDetectionResult(src *DetectionResult) *DetectionResult {
	if src == nil {
		return nil
	}

	src.mutex.RLock()
	defer src.mutex.RUnlock()

	dst := &DetectionResult{
		ProjectTypes: make(map[string]float64),
		Editors:      make(map[string]float64),
		Metadata:     make(map[string]interface{}),
		CompletedAt:  src.CompletedAt,
	}

	// Copy maps
	for k, v := range src.ProjectTypes {
		dst.ProjectTypes[k] = v
	}

	for k, v := range src.Editors {
		dst.Editors[k] = v
	}

	for k, v := range src.Metadata {
		dst.Metadata[k] = v
	}

	// Copy project info
	if src.ProjectInfo != nil {
		dst.ProjectInfo = &ProjectInfo{
			Type:        src.ProjectInfo.Type,
			Path:        src.ProjectInfo.Path,
			Framework:   src.ProjectInfo.Framework,
			Language:    src.ProjectInfo.Language,
			Editor:      src.ProjectInfo.Editor,
			LLMProvider: src.ProjectInfo.LLMProvider,
			LLMModel:    src.ProjectInfo.LLMModel,
			Confidence:  src.ProjectInfo.Confidence,
			Metadata:    make(map[string]interface{}),
		}

		for k, v := range src.ProjectInfo.Metadata {
			dst.ProjectInfo.Metadata[k] = v
		}
	}

	return dst
}

// RegisterDefaultTasks registers all the default detection tasks
func (dm *DetectionManager) RegisterDefaultTasks() {
	// Register editor detectors with 60 second timeout
	dm.RegisterDetectionTask(NewVSCodeDetector(), 1.0, 60*time.Second)
	dm.RegisterDetectionTask(NewJetBrainsDetector(), 1.0, 60*time.Second)
	dm.RegisterDetectionTask(NewLLMDetector(), 1.0, 60*time.Second)

	// Register language/framework detectors with 30 second timeout
	dm.RegisterDetectionTask(NewGoDetector(), 0.8, 30*time.Second)
	dm.RegisterDetectionTask(NewNodeJSDetector(), 0.8, 30*time.Second)
	dm.RegisterDetectionTask(NewPythonDetector(), 0.8, 30*time.Second)
	dm.RegisterDetectionTask(NewDockerDetector(), 0.9, 30*time.Second)
}
