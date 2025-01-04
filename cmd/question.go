package cmd

import (
	"codetest/internal/usecase"
	"codetest/internal/usecase/web_api"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var summaryFilePath string

// questionNodeCmd 定义了 file 节点的命令
//
//	go run entry/main.go question 请帮我分析一下这个项目主要是干什么的 -t sk-bZf8kBSewGAYqJ85CgDZzJtGyBO1AcBdA6OKdy0ntNkUtob6 -s /home/gw123/go/src/github.com/mytoolzone/task-mini-program/result/all.md
var questionNodeCmd = &cobra.Command{
	Use:   "question [question]",
	Short: "Ask a question and get an AI-generated answer about the file node usage",
	Args:  cobra.ExactArgs(1), // 确保用户输入一个问题
	RunE: func(cmd *cobra.Command, args []string) error {
		// 获取用户输入的问题
		question := args[0]
		if question == "" {
			return fmt.Errorf("question cannot be empty")
		}
		return runFileNode(apiToken, question)
	},
}

// init 函数用于设置 file-node 命令的参数
func init() {
	rootCmd.AddCommand(questionNodeCmd) // 将子命令添加到根命令
	questionNodeCmd.Flags().StringVarP(&apiToken, "token", "t", "", "API token for AI analysis (required)")
	questionNodeCmd.Flags().StringVarP(&summaryFilePath, "summary-dir", "s", "./result/all.md", "总结文件输出地方")
}

// runFileNode 主要逻辑
func runFileNode(token, question string) error {
	llmClient := web_api.NewChatGPTClient(token)
	//codeSummaryRepo := repo.NewCodeSummaryRepo(outputDir)
	aiCode := usecase.NewAiCode(llmClient)

	summary, err := os.ReadFile(summaryFilePath)
	if err != nil {
		fmt.Println("os.ReadFile(path) Error:", err)
		return err
	}
	// 调用 AI 客户端以获取答案
	answer, err := aiCode.AIQuestion(string(summary), question, "")
	if err != nil {
		return fmt.Errorf("error: %v", err)
	}

	// 打印生成的答案
	fmt.Println("AI 回复结果：")
	fmt.Println(answer)
	return nil
}
