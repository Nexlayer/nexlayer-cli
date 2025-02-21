// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package api

import "time"

// APIResponse is a generic response type for all API responses
type APIResponse[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}

// DeploymentResponse represents the response from starting a deployment
type DeploymentResponse struct {
	Namespace string `json:"namespace"`
	URL       string `json:"url"`
}

// Deployment represents a deployment in the system
type Deployment struct {
	Namespace    string      `json:"namespace"`
	TemplateID   string      `json:"templateId"`
	TemplateName string      `json:"templateName"`
	Status       string      `json:"status"`
	URL          string      `json:"url"`
	CustomDomain string      `json:"customDomain"`
	Version      string      `json:"version"`
	CreatedAt    time.Time   `json:"createdAt"`
	LastUpdated  time.Time   `json:"lastUpdated"`
	PodStatuses  []PodStatus `json:"podStatuses"`
}

// PodStatus represents the status of a pod in a deployment
type PodStatus struct {
	Name      string    `json:"name"`
	Type      string    `json:"type"`
	Status    string    `json:"status"`
	Ready     bool      `json:"ready"`
	Restarts  int       `json:"restarts"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"createdAt"`
}

// Domain represents a custom domain configuration
type Domain struct {
	Domain        string    `json:"domain"`
	ApplicationID string    `json:"applicationId"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"createdAt"`
	SSLEnabled    bool      `json:"sslEnabled"`
}

// DomainResponse represents the response from saving a custom domain
type DomainResponse struct {
	Domain string `json:"domain"`
	URL    string `json:"url"`
}

// APIError represents an error returned by the API
type APIError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
	ErrorCode  string `json:"error"`
}

// Error implements the error interface
func (e *APIError) Error() string {
	return e.Message
}
