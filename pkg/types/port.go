// Copyright (c) 2025 Nexlayer. All rights reserved.n// Use of this source code is governed by an MIT-stylen// license that can be found in the LICENSE file.nn
package types

// PortConfig defines port configuration
type PortConfig struct {
	Container int    `json:"container" yaml:"container"`
	Host      int    `json:"host" yaml:"host"`
	Protocol  string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}
