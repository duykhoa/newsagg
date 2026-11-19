package recap

import (
	"news-terminal/internal/model"
)

func GenerateSummaries(providerName string, ollamaEndpoint, ollamaModel string, geminiAPIKey, geminiModel string, stories []model.Story, batchSize int) error {
	provider, err := NewProvider(providerName, ollamaEndpoint, ollamaModel, geminiAPIKey, geminiModel)
	if err != nil {
		return err
	}

	if batchSize <= 0 {
		batchSize = 25
	}

	for i := 0; i < len(stories); i += batchSize {
		end := i + batchSize
		if end > len(stories) {
			end = len(stories)
		}

		batch := stories[i:end]
		summaries, err := provider.GenerateSummaries(batch)
		if err != nil {
			for j := i; j < end; j++ {
				stories[j].Recap = "Failed to generate summary: " + err.Error()
			}
			continue
		}

		for j, summary := range summaries {
			if i+j < len(stories) {
				stories[i+j].Recap = summary
			}
		}
	}

	return nil
}
