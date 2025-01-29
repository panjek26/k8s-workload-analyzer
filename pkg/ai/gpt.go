package ai

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "io"
)

type GPTClient struct {
    apiKey  string
    baseURL string
}

type gptRequest struct {
    Model       string    `json:"model"`
    Messages    []message `json:"messages"`
    Temperature float64   `json:"temperature"`
}

type message struct {
    Role    string `json:"role"`
    Content string `json:"content"`
}

func NewGPTClient(apiKey string) *GPTClient {
    return &GPTClient{
        apiKey:  apiKey,
        baseURL: "https://api.openai.com/v1/chat/completions",
    }
}

func (c *GPTClient) AnalyzeWorkload(yaml, namespace, workloadType, workloadName string) (*WorkloadAnalysis, error) {
    prompt := fmt.Sprintf(`Analyze this Kubernetes %s named "%s" in namespace "%s" and provide security, performance, and best practices recommendations:

%s

Please format your response as JSON with the following structure:
{
    "analysis": "overall analysis text",
    "recommendations": [
        {
            "category": "security|performance|best-practices",
            "description": "issue description",
            "severity": "high|medium|low",
            "suggested_action": "how to fix"
        }
    ]
}`, workloadType, workloadName, namespace, yaml)

    payload := gptRequest{
        Model: "gpt-3.5-turbo-0125",  // Changed from "gpt-4" to "gpt-3.5-turbo"
        Messages: []message{
            {
                Role:    "user",
                Content: prompt,
            },
        },
        Temperature: 0.3,
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }

    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("API request failed with status: %d, body: %s", resp.StatusCode, string(body))
    }

    var gptResponse struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&gptResponse); err != nil {
        return nil, fmt.Errorf("failed to decode GPT response: %v", err)
    }

    if len(gptResponse.Choices) == 0 {
        return nil, fmt.Errorf("no response from GPT")
    }

    var result WorkloadAnalysis
    if err := json.Unmarshal([]byte(gptResponse.Choices[0].Message.Content), &result); err != nil {
        return nil, fmt.Errorf("failed to parse GPT analysis: %v", err)
    }

    return &result, nil
}