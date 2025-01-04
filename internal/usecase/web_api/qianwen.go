package web_api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Message 定义 API 消息结构
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// RequestBody 定义请求体
type RequestBody struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

// QwenClient 封装 Qwen 客户端
type QwenClient struct {
	client  *http.Client
	logFile *os.File
}

// NewQwenClient 创建新的 QwenClient
func NewQwenClient(apiKey string) *QwenClient {
	if apiKey == "" {
		apiKey = os.Getenv("DASHSCOPE_API_KEY")
	}

	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return nil
	}

	return &QwenClient{
		client:  &http.Client{},
		logFile: logFile,
	}
}

// LogDetail 记录详细日志
func (c *QwenClient) LogDetail(text string) {
	if _, err := c.logFile.WriteString(text + "\n"); err != nil {
		fmt.Printf("Failed to write log: %v\n", err)
	}
}

// GetResponse 调用 Qwen API 并返回回复
func (c *QwenClient) GetResponse(prompt string) (string, error) {
	// 构建请求体
	requestBody := RequestBody{
		Model: "qwen-plus",
		Messages: []Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: prompt},
		},
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	req, err := http.NewRequest("POST", "https://dashscope.aliyuncs.com/compatible-mode/v1/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("DASHSCOPE_API_KEY"))
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	var responseBody struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.Unmarshal(bodyText, &responseBody); err != nil {
		return "", fmt.Errorf("failed to unmarshal response body: %v", err)
	}

	if len(responseBody.Choices) == 0 {
		return "", fmt.Errorf("no choices in response")
	}

	return responseBody.Choices[0].Message.Content, nil
}
