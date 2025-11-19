package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	openRouterBaseURL = "https://openrouter.ai/api/v1"
	defaultModel      = "deepseek/deepseek-chat"
	maxRetries        = 3
	retryDelay        = time.Second * 2
)

// OpenRouterClient handles communication with OpenRouter API
type OpenRouterClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Tool represents a function tool definition
type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

// ToolFunction defines a function tool
type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Parameters  map[string]interface{} `json:"parameters"`
}

// ToolCall represents a tool call in the response
type ToolCall struct {
	ID       string `json:"id"`
	Type     string `json:"type"`
	Function struct {
		Name      string `json:"name"`
		Arguments string `json:"arguments"`
	} `json:"function"`
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Tools       []Tool    `json:"tools,omitempty"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role      string     `json:"role"`
			Content   string     `json:"content"`
			ToolCalls []ToolCall `json:"tool_calls,omitempty"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error,omitempty"`
}

// NewOpenRouterClient creates a new OpenRouter client
func NewOpenRouterClient(apiKey string) *OpenRouterClient {
	return &OpenRouterClient{
		apiKey:  apiKey,
		baseURL: openRouterBaseURL,
		httpClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// Chat sends a chat completion request
func (c *OpenRouterClient) Chat(ctx context.Context, messages []Message, model string) (*ChatResponse, error) {
	if model == "" {
		model = defaultModel
	}

	req := ChatRequest{
		Model:       model,
		Messages:    messages,
		Temperature: 0.7,
		MaxTokens:   4096,
	}

	return c.sendChatRequest(ctx, req)
}

// ChatWithTools sends a chat completion request with tool support
func (c *OpenRouterClient) ChatWithTools(ctx context.Context, messages []Message, tools []Tool, model string) (*ChatResponse, error) {
	if model == "" {
		model = defaultModel
	}

	req := ChatRequest{
		Model:       model,
		Messages:    messages,
		Tools:       tools,
		Temperature: 0.7,
		MaxTokens:   4096,
	}

	return c.sendChatRequest(ctx, req)
}

// sendChatRequest sends a chat request with retry logic
func (c *OpenRouterClient) sendChatRequest(ctx context.Context, chatReq ChatRequest) (*ChatResponse, error) {
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("[OpenRouter] Retry attempt %d/%d after error: %v", attempt+1, maxRetries, lastErr)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(retryDelay * time.Duration(attempt)):
			}
		}

		resp, err := c.doRequest(ctx, chatReq)
		if err == nil {
			if resp.Error != nil {
				return nil, fmt.Errorf("OpenRouter API error: %s (type: %s, code: %s)", resp.Error.Message, resp.Error.Type, resp.Error.Code)
			}
			return resp, nil
		}

		lastErr = err

		// Don't retry on context errors
		if ctx.Err() != nil {
			return nil, ctx.Err()
		}
	}

	return nil, fmt.Errorf("failed after %d retries: %w", maxRetries, lastErr)
}

// doRequest performs the actual HTTP request
func (c *OpenRouterClient) doRequest(ctx context.Context, chatReq ChatRequest) (*ChatResponse, error) {
	reqBody, err := json.Marshal(chatReq)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("[OpenRouter] Request: model=%s, messages=%d, tools=%d", chatReq.Model, len(chatReq.Messages), len(chatReq.Tools))

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/chat/completions", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("HTTP-Referer", "https://fitness-coach-app.com")
	req.Header.Set("X-Title", "Fitness Coach AI")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	log.Printf("[OpenRouter] Response: status=%d, body_length=%d", resp.StatusCode, len(respBody))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(respBody, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if len(chatResp.Choices) > 0 {
		log.Printf("[OpenRouter] Success: content_length=%d, finish_reason=%s, tokens=%d",
			len(chatResp.Choices[0].Message.Content),
			chatResp.Choices[0].FinishReason,
			chatResp.Usage.TotalTokens)
	}

	return &chatResp, nil
}
