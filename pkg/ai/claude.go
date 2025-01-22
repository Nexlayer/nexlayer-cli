package ai

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type ClaudeClient struct {
    apiKey string
}

func (c *ClaudeClient) GetProvider() string {
    return "claude"
}

func (c *ClaudeClient) GetModel() string {
    return "claude-2"
}

func (c *ClaudeClient) Suggest(prompt string) (string, error) {
    reqBody, err := json.Marshal(map[string]interface{}{
        "model":     "claude-2",
        "prompt":    prompt,
        "max_tokens_to_sample": 1000,
    })
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %w", err)
    }

    req, err := http.NewRequest("POST", "https://api.anthropic.com/v1/complete", bytes.NewBuffer(reqBody))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("X-API-Key", c.apiKey)

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return "", fmt.Errorf("failed to make request: %w", err)
    }
    defer resp.Body.Close()

    body, err := io.ReadAll(resp.Body)
    if err != nil {
        return "", fmt.Errorf("failed to read response: %w", err)
    }

    var response struct {
        Completion string `json:"completion"`
    }

    if err := json.Unmarshal(body, &response); err != nil {
        return "", fmt.Errorf("failed to parse response: %w", err)
    }

    return response.Completion, nil
}
