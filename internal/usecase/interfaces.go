package usecase

import "codetest/internal/entity"

type LLMClient interface {
	GetResponse(prompt string) (string, error)
}

type Logger interface {
	LogDetail(text string)
}

type AICodeUseCase interface {
	AIAnalysisCode(filename, code string) (string, entity.ParsedYAML, error)
	AIQuestion(summaryContent, question, helpInfo string) ([]string, error)
}
