// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package detection

import (
	"fmt"
	"gopkg.in/yaml.v2"
)

// YAMLConfig represents the nexlayer.yaml structure
type YAMLConfig struct {
	Application struct {
		Name string `yaml:"name"`
		Type string `yaml:"type"`
		Pods []PodConfig `yaml:"pods"`
	} `yaml:"application"`
}

// PodConfig represents a pod configuration in nexlayer.yaml
type PodConfig struct {
	Name         string   `yaml:"name"`
	Type         string   `yaml:"type"`
	Image        string   `yaml:"image"`
	ServicePorts []int    `yaml:"servicePorts"`
	Env          []string `yaml:"env,omitempty"`
	BuildConfig  *struct {
		Context    string `yaml:"context"`
		Dockerfile string `yaml:"dockerfile"`
	} `yaml:"buildConfig,omitempty"`
}

// GenerateYAML creates a nexlayer.yaml based on project detection
func GenerateYAML(info *ProjectInfo) (string, error) {
	config := YAMLConfig{}
	config.Application.Name = info.Name
	config.Application.Type = string(info.Type)

	// Configure pods based on project type
	switch info.Type {
	case TypeNextjs:
		config.Application.Pods = []PodConfig{
			{
				Name:  "frontend",
				Type:  "frontend",
				Image: "node:18-alpine",
				ServicePorts: []int{3000},
				BuildConfig: &struct {
					Context    string `yaml:"context"`
					Dockerfile string `yaml:"dockerfile"`
				}{
					Context:    ".",
					Dockerfile: "Dockerfile",
				},
			},
		}
	case TypeReact:
		config.Application.Pods = []PodConfig{
			{
				Name:  "frontend",
				Type:  "frontend",
				Image: "node:18-alpine",
				ServicePorts: []int{3000},
				BuildConfig: &struct {
					Context    string `yaml:"context"`
					Dockerfile string `yaml:"dockerfile"`
				}{
					Context:    ".",
					Dockerfile: "Dockerfile",
				},
			},
		}
	case TypeNode:
		config.Application.Pods = []PodConfig{
			{
				Name:  "api",
				Type:  "backend",
				Image: "node:18-alpine",
				ServicePorts: []int{8000},
				Env: []string{
					"NODE_ENV=production",
				},
				BuildConfig: &struct {
					Context    string `yaml:"context"`
					Dockerfile string `yaml:"dockerfile"`
				}{
					Context:    ".",
					Dockerfile: "Dockerfile",
				},
			},
		}
	case TypePython:
		config.Application.Pods = []PodConfig{
			{
				Name:  "api",
				Type:  "backend",
				Image: "python:3.9",
				ServicePorts: []int{8000},
				Env: []string{
					"PYTHONUNBUFFERED=1",
				},
				BuildConfig: &struct {
					Context    string `yaml:"context"`
					Dockerfile string `yaml:"dockerfile"`
				}{
					Context:    ".",
					Dockerfile: "Dockerfile",
				},
			},
		}
	case TypeGo:
		config.Application.Pods = []PodConfig{
			{
				Name:  "api",
				Type:  "backend",
				Image: "golang:1.21-alpine",
				ServicePorts: []int{8080},
				BuildConfig: &struct {
					Context    string `yaml:"context"`
					Dockerfile string `yaml:"dockerfile"`
				}{
					Context:    ".",
					Dockerfile: "Dockerfile",
				},
			},
		}
	case TypeMERN:
		config.Application.Pods = []PodConfig{
			{
				Name:  "frontend",
				Type:  "frontend",
				Image: "node:18-alpine",
				ServicePorts: []int{3000},
				BuildConfig: &struct {
					Context    string `yaml:"context"`
					Dockerfile string `yaml:"dockerfile"`
				}{
					Context:    "frontend",
					Dockerfile: "Dockerfile",
				},
			},
			{
				Name:  "api",
				Type:  "backend",
				Image: "node:18-alpine",
				ServicePorts: []int{8000},
				Env: []string{
					"NODE_ENV=production",
					"MONGODB_URI=mongodb://db:27017/app",
				},
				BuildConfig: &struct {
					Context    string `yaml:"context"`
					Dockerfile string `yaml:"dockerfile"`
				}{
					Context:    "backend",
					Dockerfile: "Dockerfile",
				},
			},
			{
				Name:  "db",
				Type:  "database",
				Image: "mongo:latest",
				ServicePorts: []int{27017},
			},
		}
	case TypePERN:
		config.Application.Pods = []PodConfig{
			{
				Name:  "frontend",
				Type:  "frontend",
				Image: "node:18-alpine",
				ServicePorts: []int{3000},
				BuildConfig: &struct {
					Context    string `yaml:"context"`
					Dockerfile string `yaml:"dockerfile"`
				}{
					Context:    "frontend",
					Dockerfile: "Dockerfile",
				},
			},
			{
				Name:  "api",
				Type:  "backend",
				Image: "node:18-alpine",
				ServicePorts: []int{8000},
				Env: []string{
					"NODE_ENV=production",
					"DATABASE_URL=postgresql://postgres:postgres@db:5432/app",
				},
				BuildConfig: &struct {
					Context    string `yaml:"context"`
					Dockerfile string `yaml:"dockerfile"`
				}{
					Context:    "backend",
					Dockerfile: "Dockerfile",
				},
			},
			{
				Name:  "db",
				Type:  "database",
				Image: "postgres:latest",
				ServicePorts: []int{5432},
				Env: []string{
					"POSTGRES_USER=postgres",
					"POSTGRES_PASSWORD=postgres",
					"POSTGRES_DB=app",
				},
			},
		}
	default:
		// Basic template for unknown project types
		config.Application.Pods = []PodConfig{
			{
				Name:  "app",
				Type:  "frontend",
				Image: "nginx:latest",
				ServicePorts: []int{80},
				BuildConfig: &struct {
					Context    string `yaml:"context"`
					Dockerfile string `yaml:"dockerfile"`
				}{
					Context:    ".",
					Dockerfile: "Dockerfile",
				},
			},
		}
	}

	// Marshal to YAML
	yamlBytes, err := yaml.Marshal(&config)
	if err != nil {
		return "", fmt.Errorf("failed to generate YAML: %w", err)
	}

	return string(yamlBytes), nil
}
