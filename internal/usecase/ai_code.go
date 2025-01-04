package usecase

import (
	"codetest/internal/entity"
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"regexp"
	"strings"
)

// aiCodeUseCase 处理与 AI 相关的用例
type aiCodeUseCase struct {
	client LLMClient
	logger Logger
}

// NewAiCode 创建新的 aiCodeUseCase
func NewAiCode(client LLMClient) AICodeUseCase {
	return &aiCodeUseCase{client: client}
}

// AIAnalysisCode 进行代码分析
func (uc *aiCodeUseCase) AIAnalysisCode(filename, code string) (string, entity.ParsedYAML, error) {
	response, err := uc.client.GetResponse(buildFileAnalysisPrompt(filename, code))
	if err != nil {
		return "", entity.ParsedYAML{}, err
	}

	response = cleanYAMLResponse(response)
	var parsedData entity.ParsedYAML
	if err = yaml.Unmarshal([]byte(response), &parsedData); err != nil {
		fmt.Println("Error parsing YAML:", err)
		return response, parsedData, nil
	}

	return response, parsedData, nil
}

// cleanYAMLResponse 清理 YAML 响应中的格式问题
func cleanYAMLResponse(response string) string {
	response = strings.TrimSpace(response)
	response = strings.TrimPrefix(response, "```yaml")
	response = strings.TrimSpace(strings.TrimSuffix(response, "```"))

	replacements := []struct {
		Old string
		New string
	}{
		{"structs: []\n", ""},
		{"structs: ''\n", ""},
		{"constants: ''\n", ""},
		{"constants: []\n", ""},
		{"interfaces: ''\n", ""},
		{"interfaces: []\n", ""},
		{"params: ''\n", ""},
		{"return_values: ''\n", ""},
		{"- []\n", ""},
	}

	for _, replacement := range replacements {
		response = strings.ReplaceAll(response, replacement.Old, replacement.New)
	}

	regex := regexp.MustCompile(`- (\w+): \*(.*)`)
	return regex.ReplaceAllString(response, `\$1: '*\$2'`)
}

// AIQuestion 处理问题并返回相关文件
func (uc *aiCodeUseCase) AIQuestion(summaryContent, question, helpInfo string) ([]string, error) {
	step1Response, err := uc.client.GetResponse(buildQuestionRelFilesPrompt(question, summaryContent))
	if err != nil {
		return nil, err
	}

	step1FileInfos, err := parseStep1FileInfos(step1Response)
	if err != nil {
		return nil, err
	}

	logFileInfo(uc.logger, step1FileInfos)

	for _, fileInfo := range step1FileInfos {
		if err := analyzeFile(uc.client, uc.logger, question, fileInfo); err != nil {
			return nil, err
		}
	}

	return summarizeFinalAnswer(uc.client, question, helpInfo, step1FileInfos)
}

// parseStep1FileInfos 从 YAML 响应中解析文件信息
func parseStep1FileInfos(response string) ([]*entity.Step1FileInfo, error) {
	response = strings.TrimSpace(strings.TrimPrefix(strings.TrimSuffix(response, "```"), "```yaml"))
	var fileInfos []*entity.Step1FileInfo
	if err := yaml.Unmarshal([]byte(response), &fileInfos); err != nil {
		fmt.Println("Error parsing Step1FileInfo YAML:", err)
		fmt.Println(response)
		return nil, err
	}
	return fileInfos, nil
}

// logFileInfo 记录文件信息
func logFileInfo(client Logger, fileInfos []*entity.Step1FileInfo) {
	client.LogDetail("----------问题关联的文件列表-------------")
	fmt.Println("----------问题关联的文件列表-------------")
	for _, info := range fileInfos {
		client.LogDetail(info.File)
		client.LogDetail(info.Why)
		fmt.Println(info.File, info.Why)
	}
}

// analyzeFile 分析指定文件的内容
func analyzeFile(client LLMClient, logger Logger, question string, fileInfo *entity.Step1FileInfo) error {
	fileContent, err := os.ReadFile(fileInfo.File)
	if err != nil {
		return err
	}

	prompt := buildQuestionRelFilesParsePrompt(question, "", fileInfo.File, string(fileContent))
	response, err := client.GetResponse(prompt)
	if err != nil {
		return err
	}

	logger.LogDetail("######## 局部答案分析结果 ########")
	logger.LogDetail(fileInfo.File)
	logger.LogDetail(response)
	fileInfo.ParseResult = response
	fmt.Println("分析", fileInfo.File, "完成")
	return nil
}

// summarizeFinalAnswer 总结最终答案
func summarizeFinalAnswer(client LLMClient, question, helpInfo string, fileInfos []*entity.Step1FileInfo) ([]string, error) {
	answerPromptBuilder := buildFinalAnswerPrompt(question, helpInfo)
	for _, info := range fileInfos {
		answerPromptBuilder.WriteString(info.ParseResult)
	}

	fmt.Println("开始总结答案")
	response, err := client.GetResponse(answerPromptBuilder.String())
	if err != nil {
		return nil, err
	}

	fmt.Println(response)
	return nil, nil
}
