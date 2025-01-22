package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"

	"github.com/patrickmn/go-cache"
)

const (
	bufferSize = 64 * 1024 // 64KB buffer
)

var (
	templateCache = cache.New(5*time.Minute, 10*time.Minute)
)

type PluginMetadata struct {
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Usage       string `json:"usage"`
}

type NexlayerTemplate struct {
	Name        string                 `yaml:"name" json:"name"`
	Version     string                 `yaml:"version" json:"version"`
	Type        string                 `yaml:"type" json:"type"`
	Environment map[string]interface{} `yaml:"environment" json:"environment"`
	Resources   []Resource             `yaml:"resources" json:"resources"`
}

type Resource struct {
	Name       string                 `yaml:"name" json:"name"`
	Type       string                 `yaml:"type" json:"type"`
	Properties map[string]interface{} `yaml:"properties" json:"properties"`
}

type LintError struct {
	Field    string
	Message  string
	Fix      func() error
	Resource string
}

func main() {
	// Handle metadata request
	if len(os.Args) > 1 && os.Args[1] == "--describe" {
		metadata := PluginMetadata{
			Name:        "lint",
			Version:     "1.0.0",
			Description: "A plugin for validating Nexlayer YAML/JSON templates",
			Usage:       "nexlayer lint [filePath] [--fix]",
		}
		json.NewEncoder(os.Stdout).Encode(metadata)
		return
	}

	// Parse flags
	fixFlag := flag.Bool("fix", false, "Automatically fix common issues")
	flag.Parse()

	// Get file path argument
	args := flag.Args()
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Error: Please provide a template file path")
		os.Exit(1)
	}
	filePath := args[0]

	// Read and parse template
	template, err := parseTemplate(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing template: %v\n", err)
		os.Exit(1)
	}

	// Run linting checks
	errors := lintTemplate(template)

	// Handle errors
	if len(errors) > 0 {
		fmt.Println("Found", len(errors), "issue(s):")
		for i, err := range errors {
			fmt.Printf("%d. %s: %s\n", i+1, err.Field, err.Message)
			if *fixFlag && err.Fix != nil {
				if fixErr := err.Fix(); fixErr != nil {
					fmt.Printf("   Failed to auto-fix: %v\n", fixErr)
				} else {
					fmt.Println("   ✓ Auto-fixed")
				}
			}
		}
		if !*fixFlag && containsFixable(errors) {
			fmt.Println("\nTip: Run with --fix to attempt automatic fixes")
		}
		if *fixFlag {
			// Save fixed template
			if err := saveTemplate(filePath, template); err != nil {
				fmt.Fprintf(os.Stderr, "Error saving fixed template: %v\n", err)
				os.Exit(1)
			}
		}
		os.Exit(1)
	}

	fmt.Println("✓ Template validation passed!")
}

func parseTemplate(filePath string) (*NexlayerTemplate, error) {
	// Check cache first
	if cached, found := templateCache.Get(filePath); found {
		return cached.(*NexlayerTemplate), nil
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	// Use buffered reader for better performance
	reader := bufio.NewReaderSize(file, bufferSize)
	data, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var template NexlayerTemplate
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".yaml", ".yml":
		err = yaml.UnmarshalStrict(data, &template)
	case ".json":
		decoder := json.NewDecoder(bytes.NewReader(data))
		decoder.DisallowUnknownFields()
		err = decoder.Decode(&template)
	default:
		return nil, fmt.Errorf("unsupported file format: %s", ext)
	}

	if err != nil {
		return nil, &ParseError{
			File: filePath,
			Err:  err,
		}
	}

	// Cache the parsed template
	templateCache.Set(filePath, &template, cache.DefaultExpiration)
	return &template, nil
}

// ParseError provides detailed parsing error information
type ParseError struct {
	File string
	Err  error
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error in %s: %v", e.File, e.Err)
}

// ConcurrentLint performs template validation concurrently
func ConcurrentLint(templates []string) []LintError {
	var (
		wg            sync.WaitGroup
		mu            sync.Mutex
		errors        []LintError
		workers       = runtime.GOMAXPROCS(0)
		templatesChan = make(chan string, len(templates))
	)

	// Fill templates channel
	for _, t := range templates {
		templatesChan <- t
	}
	close(templatesChan)

	// Start workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for template := range templatesChan {
				tmpl, err := parseTemplate(template)
				if err != nil {
					mu.Lock()
					errors = append(errors, LintError{
						Field:   template,
						Message: err.Error(),
					})
					mu.Unlock()
					continue
				}

				errs := lintTemplate(tmpl)
				if len(errs) > 0 {
					mu.Lock()
					errors = append(errors, errs...)
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	return errors
}

func lintTemplate(template *NexlayerTemplate) []LintError {
	var errors []LintError

	// Validate template structure
	if err := validateTemplateStructure(template); err != nil {
		errors = append(errors, LintError{
			Message: fmt.Sprintf("invalid template structure: %v", err),
		})
	}

	// Validate resources concurrently
	var (
		wg            sync.WaitGroup
		mu            sync.Mutex
		workers       = runtime.GOMAXPROCS(0)
		resourcesChan = make(chan Resource, len(template.Resources))
	)

	// Fill resources channel
	for _, resource := range template.Resources {
		resourcesChan <- resource
	}
	close(resourcesChan)

	// Start resource validation workers
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for resource := range resourcesChan {
				if errs := validateResource(resource); len(errs) > 0 {
					mu.Lock()
					errors = append(errors, errs...)
					mu.Unlock()
				}
			}
		}()
	}

	wg.Wait()
	return errors
}

func validateTemplateStructure(template *NexlayerTemplate) error {
	if template.Version == "" {
		return errors.New("missing template version")
	}

	if len(template.Resources) == 0 {
		return errors.New("template must contain at least one resource")
	}

	return nil
}

func validateResource(resource Resource) []LintError {
	var errors []LintError

	if resource.Name == "" {
		errors = append(errors, LintError{
			Resource: resource.Type,
			Message:  "resource name is required",
		})
	}

	if resource.Type == "" {
		errors = append(errors, LintError{
			Resource: resource.Name,
			Message:  "resource type is required",
		})
	}

	// Validate resource properties based on type
	if err := validateResourceProperties(resource); err != nil {
		errors = append(errors, LintError{
			Resource: resource.Name,
			Message:  fmt.Sprintf("invalid properties: %v", err),
		})
	}

	return errors
}

func validateResourceProperties(resource Resource) error {
	validator, ok := resourceValidators[resource.Type]
	if !ok {
		return fmt.Errorf("unknown resource type: %s", resource.Type)
	}

	return validator(resource.Properties)
}

var resourceValidators = map[string]func(map[string]interface{}) error{
	"AWS::Lambda::Function": validateLambdaFunction,
	"AWS::S3::Bucket":       validateS3Bucket,
	// Add more resource validators as needed
}

func validateLambdaFunction(props map[string]interface{}) error {
	required := []string{"Runtime", "Handler", "Code"}
	for _, prop := range required {
		if _, ok := props[prop]; !ok {
			return fmt.Errorf("missing required property: %s", prop)
		}
	}
	return nil
}

func validateS3Bucket(props map[string]interface{}) error {
	// Add S3 bucket specific validation
	return nil
}

func saveTemplate(filePath string, template *NexlayerTemplate) error {
	var data []byte
	var err error

	switch strings.ToLower(filepath.Ext(filePath)) {
	case ".yaml", ".yml":
		data, err = yaml.Marshal(template)
	case ".json":
		data, err = json.MarshalIndent(template, "", "  ")
	default:
		return fmt.Errorf("unsupported file format: %s", filepath.Ext(filePath))
	}

	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	return ioutil.WriteFile(filePath, data, 0644)
}

func containsFixable(errors []LintError) bool {
	for _, err := range errors {
		if err.Fix != nil {
			return true
		}
	}
	return false
}
