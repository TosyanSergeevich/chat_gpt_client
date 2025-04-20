package chatgpt

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Client struct {
	apiKey      string
	model       string
	maxTokens   int
	temperature float64
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ImageMessage struct {
	Role    string          `json:"role"`
	Content []ContentObject `json:"content"`
}

type ContentObject struct {
	Type     string   `json:"type"`
	Text     string   `json:"text,omitempty"`
	ImageURL ImageURL `json:"image_url,omitempty"`
}

type ImageURL struct {
	URL string `json:"url"`
}

type ChatCompletionRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
}

type ChatCompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewClient(apiKey, model string, maxTokens int, temperature float64) *Client {
	return &Client{
		apiKey:      apiKey,
		model:       model,
		maxTokens:   maxTokens,
		temperature: temperature,
	}
}

func (c *Client) SendMessage(messages []Message) (string, error) {
	requestBody := ChatCompletionRequest{
		Model:       c.model,
		Messages:    messages,
		MaxTokens:   c.maxTokens,
		Temperature: c.temperature,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	// Log the request for debugging
	log.Printf("Sending request to ChatGPT: %s", string(jsonData))

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Log the response for debugging
	log.Printf("ChatGPT response status: %d", resp.StatusCode)
	log.Printf("ChatGPT response body: %s", string(body))

	// Check for API errors
	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResponse); err == nil {
			return "", fmt.Errorf("API error: %s (type: %s)", errorResponse.Error.Message, errorResponse.Error.Type)
		}
		return "", fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var response ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v, body: %s", err, string(body))
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from ChatGPT, full response: %s", string(body))
	}

	return response.Choices[0].Message.Content, nil
}

func (c *Client) SendImageMessage(messages []ImageMessage) (string, error) {
	jsonData, err := json.Marshal(map[string]interface{}{
		"model":       c.model,
		"messages":    messages,
		"max_tokens":  c.maxTokens,
		"temperature": c.temperature,
	})
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}

	// Log the request for debugging
	log.Printf("Sending image request to ChatGPT: %s", string(jsonData))

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}

	// Log the response for debugging
	log.Printf("ChatGPT image response status: %d", resp.StatusCode)
	log.Printf("ChatGPT image response body: %s", string(body))

	// Check for API errors
	if resp.StatusCode != http.StatusOK {
		var errorResponse struct {
			Error struct {
				Message string `json:"message"`
				Type    string `json:"type"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResponse); err == nil {
			return "", fmt.Errorf("API error: %s (type: %s)", errorResponse.Error.Message, errorResponse.Error.Type)
		}
		return "", fmt.Errorf("API error: status %d, body: %s", resp.StatusCode, string(body))
	}

	var response ChatCompletionResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v, body: %s", err, string(body))
	}

	if len(response.Choices) == 0 {
		return "", fmt.Errorf("no response from ChatGPT, full response: %s", string(body))
	}

	return response.Choices[0].Message.Content, nil
}
