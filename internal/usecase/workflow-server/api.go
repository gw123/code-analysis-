package workflow_server

import (
	"bytes"
	"codetest/internal/entity"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Data struct {
		Token    string `json:"token"`
		UserID   int    `json:"user_id"`
		Username string `json:"username"`
	} `json:"data"`
}

type ErrorResponse struct {
	Code  string `json:"code"`
	Error string `json:"message"`
}

// ApiClient 封装 API 客户端
type ApiClient struct {
	client      *http.Client
	apiBasePath string
	token       string
	username    string
	password    string
}

// NewApiClient 创建一个新的 ApiClient
func NewApiClient(apiBasePath, username, password string) *ApiClient {
	return &ApiClient{
		client:      &http.Client{},
		apiBasePath: apiBasePath,
		username:    username,
		password:    password,
	}
}

// Login 进行登录并获取 token
func (a *ApiClient) Login(ctx context.Context) (string, error) {
	loginReq := LoginRequest{
		Username: a.username,
		Password: a.password,
	}

	// 将请求体转换为 JSON
	jsonData, err := json.Marshal(loginReq)
	if err != nil {
		return "", fmt.Errorf("failed to marshal login request body: %v", err)
	}

	// 创建新的 HTTP 请求
	req, err := http.NewRequest("POST", a.apiBasePath+"/auth/login", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create login request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送登录请求
	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send login request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read login response body: %v", err)
	}

	// 如果返回的 HTTP 状态码不是 200，表示登录失败
	if resp.StatusCode != http.StatusOK {
		var errResp ErrorResponse
		if err := json.Unmarshal(bodyText, &errResp); err != nil {
			return "", fmt.Errorf("failed to parse error response: %v", err)
		}
		return "", fmt.Errorf("login failed: %s", errResp.Error)
	}

	// 解析登录成功的响应
	var loginResp LoginResponse
	if err := json.Unmarshal(bodyText, &loginResp); err != nil {
		return "", fmt.Errorf("failed to parse login response: %v", err)
	}

	// 返回 Token
	a.token = loginResp.Data.Token
	fmt.Println("Login successful. Token:", loginResp.Data.Token)
	return loginResp.Data.Token, nil
}

// CreateAiCode 发送 POST 请求
func (a *ApiClient) UploadCodeInfo(ctx context.Context, data entity.AICodeSnippet) (string, error) {
	// 将请求体数据转换为 JSON
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body: %v", err)
	}

	// 创建新的 HTTP 请求
	req, err := http.NewRequest("POST", a.apiBasePath+"/codes", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v", err)
	}

	// 设置请求头
	req.Header.Set("Authorization", "Bearer "+a.token)
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := a.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %v", err)
	}

	return string(bodyText), nil
}
