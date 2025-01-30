package ai

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "k8s-workload-analyzer/pkg/ai/prompts"
)

type GPTClient struct {
    apiKey string
}

func NewGPTClient(apiKey string) *GPTClient {
    return &GPTClient{apiKey: apiKey}
}

func summarizeYAML(yaml string) string {
    lines := strings.Split(yaml, "\n")
    summary := []string{}
    
    inContainer := false
    for _, line := range lines {
        // Focus on container section
        if strings.Contains(line, "containers:") {
            inContainer = true
            summary = append(summary, line)
            continue
        }
        
        // Include container-specific fields
        if inContainer && (
            strings.Contains(line, "name:") ||
            strings.Contains(line, "image:") ||
            strings.Contains(line, "resources:") ||
            strings.Contains(line, "limits:") ||
            strings.Contains(line, "requests:") ||
            strings.Contains(line, "securityContext:") ||
            strings.Contains(line, "volumeMounts:") ||
            strings.Contains(line, "ports:") ||
            strings.Contains(line, "livenessProbe:") ||
            strings.Contains(line, "readinessProbe:")) {
            summary = append(summary, line)
        }
        
        // Exit container section when indentation changes
        if inContainer && !strings.HasPrefix(line, "      ") {
            inContainer = false
        }
    }
    
    return strings.Join(summary, "\n")
}

func (c *GPTClient) AnalyzeWorkload(yaml string) (*WorkloadAnalysis, error) {
    // Summarize YAML before sending to GPT
    summarizedYAML := summarizeYAML(yaml)
    
    payload := map[string]interface{}{
        "model": "gpt-3.5-turbo",
        "messages": []map[string]string{
            {
                "role":    "system",
                "content": "You are a Kubernetes container expert. Focus on analyzing container configuration, resources, and best practices.",
            },
            {
                "role":    "user",
                "content": fmt.Sprintf(prompts.WorkloadAnalysisTemplate, summarizedYAML),
            },
        },
        "temperature": 0.1,
        "response_format": map[string]string{
            "type": "json_object",
        },
    }

    jsonData, err := json.Marshal(payload)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal request: %v", err)
    }

    req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %v", err)
    }

    req.Header.Set("Authorization", "Bearer "+c.apiKey)
    req.Header.Set("Content-Type", "application/json")

    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to make request: %v", err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        body, _ := io.ReadAll(resp.Body)
        return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
    }

    var result struct {
        Choices []struct {
            Message struct {
                Content string `json:"content"`
            } `json:"message"`
        } `json:"choices"`
        Error *struct {
            Message string `json:"message"`
        } `json:"error"`
    }

    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, fmt.Errorf("failed to decode response: %v", err)
    }

    if result.Error != nil {
        return nil, fmt.Errorf("API error: %s", result.Error.Message)
    }

    if len(result.Choices) == 0 {
        return nil, fmt.Errorf("no response from GPT")
    }

    // Improved content cleanup
    content := result.Choices[0].Message.Content
    content = strings.TrimSpace(content)
    
    // Remove any markdown or explanation text
    if idx := strings.Index(content, "{"); idx >= 0 {
        content = content[idx:]
        if lastIdx := strings.LastIndex(content, "}"); lastIdx >= 0 {
            content = content[:lastIdx+1]
        }
    }

    // Verify JSON structure
    if !strings.HasPrefix(content, "{") || !strings.HasSuffix(content, "}") {
        return nil, fmt.Errorf("invalid JSON response format: %s", content)
    }

    var analysis WorkloadAnalysis
    if err := json.Unmarshal([]byte(content), &analysis); err != nil {
        return nil, fmt.Errorf("failed to parse analysis (content: %s): %v", content, err)
    }

    return &analysis, nil
}