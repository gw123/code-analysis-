// repo/repo.go
package repo

import (
	"codetest/internal/entity"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CodeSummary 结构体用于封装文件操作
type CodeSummary struct {
	OutputDir string
}

// NewCodeSummaryRepo 返回一个新的 CodeSummary 实例
func NewCodeSummaryRepo(outputDir string) *CodeSummary {
	return &CodeSummary{
		OutputDir: outputDir,
	}
}

// SaveAIResult 保存 AI 分析结果到文件
func (r *CodeSummary) SaveAIResult(projectName, path, rawAiResponse string) error {
	resultPath := filepath.Join(r.OutputDir, strings.ReplaceAll(path, "/", "|")+".yaml")
	fmt.Println("resultPath = ", resultPath)
	if err := os.WriteFile(resultPath, []byte(rawAiResponse), 0644); err != nil {
		return fmt.Errorf("error writing result file: %v", err)
	}
	return nil
}

// UpdateSummaryFile 更新总结文件
func (r *CodeSummary) UpdateSummaryFile(projectName, path string, yamlResult *entity.ParsedYAML) error {
	var strBuilder strings.Builder

	strBuilder.WriteString(fmt.Sprintf("文件名: %s\n", path))
	strBuilder.WriteString(fmt.Sprintf("功能: %s\n", yamlResult.FunctionDescription))
	strBuilder.WriteString(fmt.Sprintf("包名: %s\n", yamlResult.FileInfo.PackageName))
	strBuilder.WriteString("依赖导入项目: ")
	strBuilder.WriteString(strings.Join(yamlResult.FileInfo.Imports, ","))
	strBuilder.WriteString("\n---\n")

	// 追加写入总结文件
	summaryFilePath := filepath.Join(r.OutputDir, "all.md")
	file, err := os.OpenFile(summaryFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("failed to open summary file: %v", err)
	}
	defer file.Close()

	if _, err := file.WriteString(strBuilder.String()); err != nil {
		return fmt.Errorf("failed to write to summary file: %v", err)
	}
	return nil
}
