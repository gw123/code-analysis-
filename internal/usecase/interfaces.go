package usecase

import (
	"codetest/internal/entity"
	"context"
)

type LLMClient interface {
	GetResponse(prompt string) (string, error)
}

type Logger interface {
	LogDetail(text string)
}

type AICodeUseCase interface {
	AIAnalysisCode(filename, code string) (string, entity.ParsedYAML, error)
	AIQuestion(summaryContent, question, helpInfo string) ([]string, error)
	UploadCodeInfo(ctx context.Context, data entity.AICodeSnippet) error
}

type ApiClient interface {
	Login(ctx context.Context) (string, error)
	UploadCodeInfo(ctx context.Context, data entity.AICodeSnippet) (string, error)
}
