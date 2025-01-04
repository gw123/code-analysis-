package web_api

import (
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	"os"
)

// ChatGPTClient 结构体封装 ChatGPT 客户端
type ChatGPTClient struct {
	client  *openai.Client
	logFile *os.File
}

// NewChatGPTClient 创建新的 ChatGPTClient
func NewChatGPTClient(apiKey string) *ChatGPTClient {
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}
	cfg := openai.DefaultConfig(apiKey)
	cfg.BaseURL = "https://api.chatanywhere.tech/v1"

	logFile, err := os.OpenFile("log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Failed to open log file:", err)
		return nil
	}

	return &ChatGPTClient{
		client:  openai.NewClientWithConfig(cfg),
		logFile: logFile,
	}
}

// LogDetail 记录详细日志
func (c *ChatGPTClient) LogDetail(text string) {
	if _, err := c.logFile.WriteString(text + "\n"); err != nil {
		fmt.Printf("Failed to write log: %v\n", err)
	}
}

// GetResponse 调用 ChatGPT API 并返回回复
func (c *ChatGPTClient) GetResponse(prompt string) (string, error) {
	ctx := context.Background()
	req := openai.ChatCompletionRequest{
		Temperature: 0,
		Model:       openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			},
		},
	}

	resp, err := c.client.CreateChatCompletion(ctx, req)
	if err != nil {
		return "", fmt.Errorf("ChatGPT request failed: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}
