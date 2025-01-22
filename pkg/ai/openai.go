package ai

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
)

type OpenAIClient struct {
    apiKey string
    model  string
}

func (c *OpenAIClient) GetProvider() string {
    return "openai"
}

func (c *OpenAIClient) GetModel() string {
    return c.model
}

func (c *OpenAIClient) Suggest(prompt string) (string, error) {
    reqBody, err := json.Marshal(map[string]interface{}{
        "model": c.model,
        "messages": []map[string]string{
            {"role": "user", "content": prompt},
        },
    })
    if err != nil {
        return "", fmt.Errorf("failed to marshal request: %w", err)
    }

    req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(reqBody))
    if err != nil {
        return "", fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer "+c.apiKey)

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
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
    }

    if err := json.Unmarshal(body, &response); err != nil {
        return "", fmt.Errorf("failed to parse response: %w", err)
    }

    if len(response.Choices) == 0 {
        return "", fmt.Errorf("no suggestions received")
    }

    return response.Choices[0].Message.Content, nil
}
