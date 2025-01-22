package version

// Formatted with gofmt -s
// Version is the current version of the Nexlayer CLI
const Version = "v0.1.0-alpha.1"

// GetVersion returns the current version of the CLI
func GetVersion() string {
	return Version
}
