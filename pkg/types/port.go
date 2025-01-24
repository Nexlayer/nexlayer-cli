package types

// PortConfig defines port configuration
type PortConfig struct {
	Container int    `json:"container" yaml:"container"`
	Host      int    `json:"host" yaml:"host"`
	Protocol  string `json:"protocol,omitempty" yaml:"protocol,omitempty"`
}
