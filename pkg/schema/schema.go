// Copyright (c) 2025 Nexlayer. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package schema

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

// Schema represents the Nexlayer YAML schema
type Schema struct {
	Application struct {
		Name string `json:"name" yaml:"name"`
		URL  string `json:"url,omitempty" yaml:"url,omitempty"`
	} `json:"application" yaml:"application"`
	RegistryLogin *struct {
		Registry            string `json:"registry" yaml:"registry"`
		Username           string `json:"username" yaml:"username"`
		PersonalAccessToken string `json:"personalAccessToken" yaml:"personalAccessToken"`
	} `json:"registryLogin,omitempty" yaml:"registryLogin,omitempty"`
	Pods []struct {
		Name         string `json:"name" yaml:"name"`
		Type         string `json:"type,omitempty" yaml:"type,omitempty"`
		Path         string `json:"path,omitempty" yaml:"path,omitempty"`
		Image        string `json:"image" yaml:"image"`
		Volumes      []struct {
			Name      string `json:"name" yaml:"name"`
			Size      string `json:"size" yaml:"size"`
			MountPath string `json:"mountPath" yaml:"mountPath"`
		} `json:"volumes,omitempty" yaml:"volumes,omitempty"`
		Secrets []struct {
			Name      string `json:"name" yaml:"name"`
			Data      string `json:"data" yaml:"data"`
			MountPath string `json:"mountPath" yaml:"mountPath"`
			FileName  string `json:"fileName" yaml:"fileName"`
		} `json:"secrets,omitempty" yaml:"secrets,omitempty"`
		Vars         map[string]string `json:"vars,omitempty" yaml:"vars,omitempty"`
		ServicePorts []int            `json:"servicePorts" yaml:"servicePorts"`
		Entrypoint   string           `json:"entrypoint,omitempty" yaml:"entrypoint,omitempty"`
		Command      string           `json:"command,omitempty" yaml:"command,omitempty"`
	} `json:"pods" yaml:"pods"`
}

// LoadSchema loads the Nexlayer schema from the embedded JSON file
func LoadSchema() (*Schema, error) {
	schemaPath := filepath.Join("docs", "api-reference.json")
	data, err := ioutil.ReadFile(schemaPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema file: %w", err)
	}

	var schema Schema
	if err := json.Unmarshal(data, &schema); err != nil {
		return nil, fmt.Errorf("failed to parse schema: %w", err)
	}

	return &schema, nil
}
