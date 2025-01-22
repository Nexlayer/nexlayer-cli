package ai

// AIClient defines the interface that all AI providers must implement
type AIClient interface {
    // Suggest receives a prompt and returns a suggestion or an error
    Suggest(prompt string) (string, error)
    // GetProvider returns the name of the AI provider
    GetProvider() string
    // GetModel returns the model being used
    GetModel() string
}

// Config holds the AI configuration
type Config struct {
    Provider string
    APIKey   string
}
