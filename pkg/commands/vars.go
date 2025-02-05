// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package commands

// Package commands contains variables used across different CLI commands.

// CI command variables
var (
	// Stack represents the application stack for CI commands.
	// It is used to specify the technology stack being deployed.
	Stack string
	// Registry specifies the container image registry URL.
	// It is used to define where images are stored and retrieved.
	Registry string
	// ImageName is the name of the container image.
	// It is used to identify the image within the registry.
	ImageName string
	// ImageTag is the tag for the container image.
	// It is used to specify the version of the image.
	ImageTag string
	// BuildContext specifies the build context directory.
	// It is used to define the root directory for building images.
	BuildContext string
	// Token is the authentication token for accessing the registry.
	// It is used to authenticate API requests to the registry.
	Token string
)

// Service command variables
var (
	// AppName represents the application name for service commands.
	// It is used to identify the application being managed.
	AppName string
	// Service specifies the service name within the application.
	// It is used to target specific services for operations.
	Service string
	// OutputFormat defines the format for command output.
	// It is used to specify how results are displayed (e.g., JSON, YAML).
	OutputFormat string
	// OutputFile specifies the file path for output.
	// It is used to redirect command output to a file.
	OutputFile string
	// EnvPairs contains environment variable key-value pairs.
	// It is used to pass configuration settings to services.
	EnvPairs []string
)
