package workflow_server

import (
	"codetest/internal/entity"
	"context"
	"fmt"
	"github.com/go-resty/resty/v2"
	"time"
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

type Response struct {
	Code    int         `json:"code"`
	CodeEn  string      `json:"code_en"`
	Doc     string      `json:"doc,omitempty"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

// Project represents the project entity.
type Project struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	Name            string    `gorm:"uniqueIndex;not null" json:"name"`
	Tags            []string  `gorm:"type:varchar[];default:NULL" json:"tags"`               // 标签 (使用 pq.StringArray 来支持 PostgreSQL 的数组类型)
	GitUrl          string    `gorm:"type:varchar(1024);not null" json:"git_url"`            // 代码地址
	Desc            string    `gorm:"type:varchar(4096);default:NULL" json:"desc"`           // 描述信息
	Language        string    `gorm:"type:varchar(64);not null" json:"language"`             // 编程语言
	LanguageVersion string    `gorm:"type:varchar(32);default:NULL" json:"language_version"` // 语言版本
	UserID          uint      `gorm:"default:NULL" json:"user_id"`                           // 用户 ID (可以为 NULL
	Status          string    `json:"status"`                                                // 状态
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// ApiClient 封装 API 客户端
type ApiClient struct {
	client      *resty.Client
	apiBasePath string
	token       string
	username    string
	password    string
}

// NewApiClient 创建一个新的 ApiClient
func NewApiClient(apiBasePath, username, password string) *ApiClient {
	return &ApiClient{
		client: resty.New().SetBaseURL(apiBasePath + "/api/v1/"),
		//apiBasePath: apiBasePath + "/api/v1/",
		username: username,
		password: password,
	}
}

func (a *ApiClient) Login(ctx context.Context) (string, error) {
	var response LoginResponse
	resp, err := a.client.R().
		SetContext(ctx).
		SetBody(LoginRequest{Username: a.username, Password: a.password}).
		SetResult(&response).
		Post("/auth/login")
	if err != nil {
		return "", fmt.Errorf("failed to send login request: %v", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("login failed: %s", resp.String())
	}

	a.token = response.Data.Token
	return response.Data.Token, nil
}

func (a *ApiClient) UploadCodeInfo(ctx context.Context, data entity.AICodeSnippet) (string, error) {
	resp, err := a.client.R().
		SetContext(ctx).
		SetAuthToken(a.token).
		SetBody(data).
		Post("/codes")
	if err != nil {
		return "", fmt.Errorf("failed to send upload request: %v", err)
	}
	if resp.IsError() {
		return "", fmt.Errorf("upload failed: %s", resp.String())
	}

	return resp.String(), nil
}

func (a *ApiClient) GetProjectByID(ctx context.Context, projectID uint) (*Project, error) {
	url := fmt.Sprintf("%s/projects/%d", a.apiBasePath, projectID)
	var response struct {
		Data    Project `json:"data"`
		Message string  `json:"message"`
	}
	resp, err := a.client.R().
		SetContext(ctx).
		SetAuthToken(a.token).
		SetResult(&response).
		Get(url)
	if err != nil {
		return nil, err
	}
	if resp.IsError() {
		return nil, fmt.Errorf("failed to fetch project: %s", resp.String())
	}
	return &response.Data, nil
}

func (a *ApiClient) UpdateProject(ctx context.Context, project *Project) error {
	url := fmt.Sprintf("/projects/%d", project.ID)
	response := Response{}
	resp, err := a.client.R().
		SetContext(ctx).
		SetAuthToken(a.token).
		SetBody(project).
		SetResult(&response).
		Put(url)
	if err != nil {
		return err
	}
	if resp.IsError() {
		return fmt.Errorf("failed to update project: %s", response.Message)
	}
	return nil
}
