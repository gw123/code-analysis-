// cmd/cmd.go
package cmd

import (
	code "codetest"
	"codetest/internal/usecase"
	"codetest/internal/usecase/repo"
	"codetest/internal/usecase/web_api"
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var (
	dir       string
	apiToken  string
	outputDir string
)

// analyzeCmd 定义了分析命令
var analyzeCmd = &cobra.Command{
	Use:   "analyze",
	Short: "Analyze code in the specified directory using AI",
	RunE: func(cmd *cobra.Command, args []string) error {
		return run(dir, apiToken)
	},
}

// init 函数用于设置分析命令的参数
func init() {
	rootCmd.AddCommand(analyzeCmd)

	analyzeCmd.Flags().StringVarP(&dir, "dir", "d", ".", "Directory to analyze (required)")
	analyzeCmd.Flags().StringVarP(&apiToken, "token", "t", "", "API token for AI analysis (required)")
	analyzeCmd.Flags().StringVarP(&outputDir, "output-dir", "o", "./result", "总结文件输出地方")

}

// run 主要逻辑
func run(directory, token string) error {
	llmClient := web_api.NewChatGPTClient(token)
	codeSummaryRepo := repo.NewCodeSummaryRepo(outputDir)
	aiCode := usecase.NewAiCode(llmClient)

	var count int
	// 遍历目录并处理每个文件
	err := code.WalkDir(directory, func(path string) {
		if err := processFile("", path, aiCode, codeSummaryRepo); err == nil {
			count++
		}
	})

	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}
	fmt.Printf("Processed %d files\n", count)
	return nil
}

// 处理单个文件
func processFile(projectName, path string, aiClient usecase.AICodeUseCase, repo *repo.CodeSummary) error {
	fmt.Println("Processing file:", path)

	fileContent, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("Failed to read file %s: %v", path, err)
	}

	rawAiResponse, yamlResult, err := aiClient.AIAnalysisCode(path, string(fileContent))
	if err != nil {
		return fmt.Errorf("AI analysis failed for %s: %v", path, err)
	}

	if err := repo.SaveAIResult(projectName, path, rawAiResponse); err != nil {
		return fmt.Errorf("Failed to save AI result for %s: %v", path, err)
	}

	if err := repo.UpdateSummaryFile(projectName, path, &yamlResult); err != nil {
		return fmt.Errorf("Failed to update summary for %s: %v", path, err)
	}

	return nil
}
