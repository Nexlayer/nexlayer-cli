package initcmd

// AIProvider is the interface for template generation
type AIProvider interface {
	GenerateTemplate(projectName string) (string, error)
}

var defaultProvider AIProvider

// SetAIProvider sets the AI provider for template generation
func SetAIProvider(provider AIProvider) {
	defaultProvider = provider
}

// GetAIProvider returns the current AI provider
func GetAIProvider() AIProvider {
	return defaultProvider
}
