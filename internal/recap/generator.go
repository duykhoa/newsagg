package recap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"news-terminal/internal/model"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

type GeminiRequest struct {
	Contents []GeminiContent `json:"contents"`
}

type GeminiContent struct {
	Parts []GeminiPart `json:"parts"`
}

type GeminiPart struct {
	Text string `json:"text"`
}

type GeminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func GenerateSummaries(provider string, ollamaEndpoint, ollamaModel string, geminiAPIKey, geminiModel string, stories []model.Story, chunkSize int) error {
	for i := range stories {
		s := &stories[i]
		prompt := fmt.Sprintf("Summarize this news article concisely in one sentence:\n\nTitle: %s (%s)\nURL: %s", s.Title, s.Source, s.URL)

		if provider == "gemini" {
			reqBody := GeminiRequest{
				Contents: []GeminiContent{
					{
						Parts: []GeminiPart{
							{Text: prompt},
						},
					},
				},
			}

			bodyBytes, err := json.Marshal(reqBody)
			if err != nil {
				s.Recap = "Failed to marshal request"
				continue
			}

			url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", geminiModel, geminiAPIKey)
			resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
			if err != nil {
				s.Recap = "Error connecting to Gemini"
				continue
			}
			defer resp.Body.Close()

			var geminiResp GeminiResponse
			if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err == nil && len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
				s.Recap = strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text)
			} else {
				s.Recap = "Failed to parse Gemini response"
			}
		} else {
			reqBody := OllamaRequest{
				Model:  ollamaModel,
				Prompt: prompt,
				Stream: false,
			}

			bodyBytes, err := json.Marshal(reqBody)
			if err != nil {
				s.Recap = "Failed to marshal request"
				continue
			}

			resp, err := http.Post(ollamaEndpoint+"/api/generate", "application/json", bytes.NewBuffer(bodyBytes))
			if err != nil {
				s.Recap = "Error connecting to Ollama"
				continue
			}
			defer resp.Body.Close()

			var ollamaResp OllamaResponse
			if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err == nil {
				s.Recap = strings.TrimSpace(ollamaResp.Response)
			} else {
				s.Recap = "Failed to parse Ollama response"
			}
		}
	}

	return nil
}
