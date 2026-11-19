package recap

import (
	"fmt"
	"news-terminal/internal/model"
)

type Provider interface {
	GenerateSummaries(stories []model.Story) ([]string, error)
	GenerateFullContentTUI(title, url, rawContent string) (string, error)
}

func NewProvider(provider, ollamaEndpoint, ollamaModel, geminiAPIKey, geminiModel string) (Provider, error) {
	switch provider {
	case "gemini":
		return &GeminiProvider{
			APIKey: geminiAPIKey,
			Model:  geminiModel,
		}, nil
	case "ollama", "":
		return &OllamaProvider{
			Endpoint: ollamaEndpoint,
			Model:    ollamaModel,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported LLM provider: %s", provider)
	}
}
