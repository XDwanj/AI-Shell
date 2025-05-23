package openai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"AI-Shell/internal/config"
)

// Message 表示对话消息
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// Request 表示发送到 OpenAI API 的请求
type Request struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	MaxTokens   int       `json:"max_tokens"`
	Temperature float64   `json:"temperature"`
	Stream      bool      `json:"stream"`
}

// Choice 表示 API 响应中的选择
type Choice struct {
	Message struct {
		Content string `json:"content"`
	} `json:"message"`
}

// Response 表示从 OpenAI API 接收到的响应
type Response struct {
	Choices []Choice `json:"choices"`
}

// RequestResponse 包含请求和响应数据
type RequestResponse struct {
	Request  *Request  `json:"request"`
	Response *Response `json:"response"`
}

// Client OpenAI API 客户端
type Client struct {
	config *config.Config
	client *http.Client
}

// NewClient 创建新的 OpenAI 客户端
func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		client: &http.Client{
			Timeout: time.Second * 30,
		},
	}
}

// SendRequest 发送请求到 OpenAI API
func (c *Client) SendRequest(systemPrompt, userPrompt string) (*Response, error) {
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	reqBody := Request{
		Model:       c.config.Model,
		Messages:    messages,
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
		Stream:      false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("POST", c.config.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("API请求失败: 状态码 %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("API请求失败: %v", errResp)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &response, nil
}

// SendRequestWithData 发送请求到 OpenAI API 并返回请求和响应数据
func (c *Client) SendRequestWithData(systemPrompt, userPrompt string) (*RequestResponse, error) {
	messages := []Message{
		{Role: "system", Content: systemPrompt},
		{Role: "user", Content: userPrompt},
	}

	reqBody := Request{
		Model:       c.config.Model,
		Messages:    messages,
		MaxTokens:   c.config.MaxTokens,
		Temperature: c.config.Temperature,
		Stream:      false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %v", err)
	}

	req, err := http.NewRequest("POST", c.config.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.config.APIKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp map[string]any
		if err := json.NewDecoder(resp.Body).Decode(&errResp); err != nil {
			return nil, fmt.Errorf("API请求失败: 状态码 %d", resp.StatusCode)
		}
		return nil, fmt.Errorf("API请求失败: %v", errResp)
	}

	var response Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	return &RequestResponse{
		Request:  &reqBody,
		Response: &response,
	}, nil
}
