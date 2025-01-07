// cmd/cmd.go
package cmd

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path/filepath"

	"codetest"
	"codetest/internal/entity"
	"codetest/internal/usecase"
	"codetest/internal/usecase/repo"
	"codetest/internal/usecase/web_api"
	workflow_server "codetest/internal/usecase/workflow-server"

	"github.com/spf13/cobra"
)

var (
	dir             string
	openAIToken     string
	outputDir       string
	projectName     string
	projectID       int
	language        string
	languageVersion string
	username        string
	password        string
	apiBasePath     string
	configFile      string // 新增：配置文件路径

)

// Config 配置结构体，用于映射 YAML 文件
type Config struct {
	ProjectName     string `yaml:"project_name"`
	Language        string `yaml:"language"`
	LanguageVersion string `yaml:"language_version"`
	ApiBasePath     string `yaml:"api_base_path"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	OutputDir       string `yaml:"output_dir"`
	OpenAIToken     string `yaml:"openai_token"`
	Dir             string `yaml:"dir"`
	ProjectID       int    `yaml:"project_id"`
}

// analyzeCmd 定义了分析命令
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze code in the specified directory using AI",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 读取配置文件
		if err := loadConfig(configFile); err != nil {
			return err
		}

		// 确保必需的参数存在
		if projectName == "" || language == "" || languageVersion == "" {
			return fmt.Errorf("projectName, language, and languageVersion are required")
		}
		if username == "" || password == "" || apiBasePath == "" {
			return fmt.Errorf("username, password, and apiBasePath are required for authentication")
		}
		return run(dir, openAIToken)
	},
}

// init 函数用于设置分析命令的参数
func init() {
	rootCmd.AddCommand(analyzeCmd)

	// 必需的参数
	analyzeCmd.Flags().StringVarP(&dir, "dir", "d", ".", "Directory to analyze (required)")
	analyzeCmd.Flags().StringVarP(&openAIToken, "token", "t", "", "API token for AI analysis (required)")
	analyzeCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "./result", "Directory to save analysis results")

	// 新增的参数
	analyzeCmd.Flags().StringVarP(&projectName, "project-name", "p", "", "Project name (required)")
	analyzeCmd.Flags().IntVarP(&projectID, "project-id", "i", 0, "Project ID (required)")
	analyzeCmd.Flags().StringVarP(&language, "language", "l", "", "Programming language (required)")
	analyzeCmd.Flags().StringVarP(&languageVersion, "language-version", "v", "", "Programming language version (required)")
	analyzeCmd.Flags().StringVarP(&username, "username", "u", "", "Username for authentication (required)")
	analyzeCmd.Flags().StringVarP(&password, "password", "w", "", "Password for authentication (required)")
	analyzeCmd.Flags().StringVarP(&apiBasePath, "api-base-path", "a", "", "Base API URL for the server (required)")
	analyzeCmd.Flags().StringVarP(&configFile, "config", "c", "./config/config.yaml", "Path to the YAML configuration file")
}

// 加载配置文件
func loadConfig(configFile string) error {
	if configFile == "" {
		return nil
	}

	file, err := os.Open(configFile)
	if err != nil {
		return fmt.Errorf("failed to open config file: %v", err)
	}
	defer file.Close()

	var config Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&config); err != nil {
		return fmt.Errorf("failed to decode config file: %v", err)
	}

	// 使用配置文件中的值，覆盖命令行传递的参数（如果有的话）
	if projectName == "" {
		projectName = config.ProjectName
	}
	if language == "" {
		language = config.Language
	}
	if languageVersion == "" {
		languageVersion = config.LanguageVersion
	}
	if apiBasePath == "" {
		apiBasePath = config.ApiBasePath
	}
	if username == "" {
		username = config.Username
	}
	if password == "" {
		password = config.Password
	}
	if outputDir == "" {
		outputDir = config.OutputDir
	}
	if openAIToken == "" {
		openAIToken = config.OpenAIToken
	}
	if dir == "" {
		dir = config.Dir
	}

	if projectID == 0 {
		projectID = config.ProjectID
	}

	return nil
}

// run 主要逻辑
func run(directory, token string) error {
	// 创建 API 客户端
	llmClient := web_api.NewChatGPTClient(token)
	codeSummaryRepo := repo.NewCodeSummaryRepo(outputDir)
	apiClient := workflow_server.NewApiClient(apiBasePath, username, password) // 使用新的身份认证参数
	loginRes, err := apiClient.Login(context.Background())
	if err != nil {
		log.Printf("Failed to login: %v\n", err)
		return err
	}
	log.Println("Login Success:", loginRes)
	aiCode := usecase.NewAiCode(llmClient, apiClient)

	var count int
	// 遍历目录并处理每个文件
	err = code.WalkDir(directory, func(path string) {
		if err := processFile(path, aiCode, codeSummaryRepo); err == nil {
			count++
		}
	})

	// 上报项目的汇总信息
	{
		// 获取项目详情
		project, err := apiClient.GetProjectByID(context.Background(), uint(projectID))
		if err != nil {
			log.Printf("Failed to get project details: %v\n", err)
			return err
		}
		summaryFilePath := filepath.Join(outputDir, "summary.md")
		summary, err := os.ReadFile(summaryFilePath)
		if err != nil {
			log.Printf("Failed to read summary file: %v\n", err)
			return err
		}

		project.Desc = string(summary)
		if err := apiClient.UpdateProject(context.Background(), project); err != nil {
			log.Printf("Failed to update project: %v\n", err)
			return err
		}
	}

	if err != nil {
		log.Printf("Error during directory traversal: %v\n", err)
		return err
	}
	fmt.Printf("Processed %d files\n", count)
	return nil
}

// 处理单个文件
func processFile(path string, aiClient usecase.AICodeUseCase, repo *repo.CodeSummary) error {
	fmt.Println("Processing file:", path)

	// 读取文件内容
	fileContent, err := os.ReadFile(path)
	if err != nil {
		log.Printf("Failed to read file %s: %v\n", path, err)
		return fmt.Errorf("failed to read file %s: %v", path, err)
	}

	// 调用 AI 进行分析
	rawAiResponse, yamlResult, err := aiClient.AIAnalysisCode(path, string(fileContent))
	if err != nil {
		log.Printf("AI analysis failed for %s: %v\n", path, err)
		return fmt.Errorf("AI analysis failed for %s: %v", path, err)
	}

	// 保存 AI 分析结果
	if err := repo.SaveAIResult(projectName, path, rawAiResponse); err != nil {
		log.Printf("Failed to save AI result for %s: %v\n", path, err)
		return fmt.Errorf("failed to save AI result for %s: %v", path, err)
	}

	// 更新总结文件
	if err := repo.UpdateSummaryFile(projectName, path, &yamlResult); err != nil {
		log.Printf("Failed to update summary for %s: %v\n", path, err)
		return fmt.Errorf("failed to update summary for %s: %v", path, err)
	}

	// 上传代码信息到远程
	err = aiClient.UploadCodeInfo(context.Background(), entity.AICodeSnippet{
		ProjectName:     projectName,
		FilePath:        path,
		FileName:        filepath.Base(path),
		FileType:        filepath.Ext(path),
		CodeRaw:         string(fileContent),
		Desc:            yamlResult.FileDescription,
		Snippet:         rawAiResponse,
		Language:        language,
		LanguageVersion: languageVersion,
		Tags: []string{
			"lang:" + language,
			"langVersion:" + languageVersion,
			"project:" + projectName,
		},
	})
	if err != nil {
		log.Printf("Failed to upload code info for %s: %v\n", path, err)
		return err
	}

	return nil
}
