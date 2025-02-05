// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package vars

// Global variables used across different commands
var (
	// APIEndpoint is the base URL for the Nexlayer API
	APIEndpoint = "https://app.staging.nexlayer.io"

	// Token is the authentication token
	Token string

	// APIURL is the base URL for the Nexlayer API
	APIURL = "https://app.staging.nexlayer.io"

	// AppID is the ID of the application to operate on
	AppID string

	// AppName is the name of the application
	AppName string

	// Namespace is the deployment namespace
	Namespace string

	// Domain is the custom domain to set
	Domain string

	// ConfigFile is the path to the YAML configuration file
	ConfigFile string

	// ServiceName is the name of the service
	ServiceName string

	// Service is the service identifier
	Service string

	// EnvVars are environment variables
	EnvVars []string

	// EnvPairs are key-value pairs for environment variables
	EnvPairs []string

	// URL is the API URL override
	URL string

	// Registry configuration
	RegistryType     string // Container registry type (ghcr, dockerhub, gcr, ecr, artifactory, gitlab)
	Registry         string // Container registry URL
	RegistryUsername string // Registry username
	RegistryRegion   string // Registry region (for ECR)
	RegistryProject  string // Registry project ID (for GCR)

	// Build configuration
	BuildContext string // Docker build context path
	ImageTag     string // Docker image tag

	// Graph configuration
	Depth        int    // Maximum depth to traverse when visualizing dependencies
	OutputFormat string // Format to use when visualizing dependencies
	OutputFile   string // File to write visualization output to
)
