package recap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"news-terminal/internal/model"
)

type GeminiProvider struct {
	APIKey string
	Model  string
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

func (p *GeminiProvider) GenerateSummaries(stories []model.Story) ([]string, error) {
	var sb strings.Builder
	sb.WriteString("Summarize each of the following news articles concisely in one sentence. Return ONLY a valid JSON array of strings, where each string is the summary corresponding in order to the articles below.\n\n")
	for i, s := range stories {
		sb.WriteString(fmt.Sprintf("[%d] Title: %s (%s) URL: %s\n", i+1, s.Title, s.Source, s.URL))
	}

	reqBody := GeminiRequest{
		Contents: []GeminiContent{
			{
				Parts: []GeminiPart{
					{Text: sb.String()},
				},
			},
		},
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("https://generativelanguage.googleapis.com/v1beta/models/%s:generateContent?key=%s", p.Model, p.APIKey)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var geminiResp GeminiResponse
	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return nil, err
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		text := strings.TrimSpace(geminiResp.Candidates[0].Content.Parts[0].Text)
		if idx := strings.Index(text, "["); idx != -1 {
			if idxEnd := strings.LastIndex(text, "]"); idxEnd != -1 && idxEnd > idx {
				text = text[idx : idxEnd+1]
			}
		}

		var summaries []string
		if err := json.Unmarshal([]byte(text), &summaries); err != nil {
			return nil, fmt.Errorf("failed to parse JSON response: %v", err)
		}

		if len(summaries) != len(stories) {
			return nil, fmt.Errorf("expected %d summaries, got %d", len(stories), len(summaries))
		}

		return summaries, nil
	}

	return nil, fmt.Errorf("failed to parse Gemini response")
}
