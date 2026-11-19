package recap

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"news-terminal/internal/model"
)

type OllamaProvider struct {
	Endpoint string
	Model    string
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
}

func (p *OllamaProvider) GenerateSummaries(stories []model.Story) ([]string, error) {
	var sb strings.Builder
	sb.WriteString("Summarize each of the following news articles concisely in one sentence. Return ONLY a valid JSON array of strings, where each string is the summary corresponding in order to the articles below.\n\n")
	for i, s := range stories {
		sb.WriteString(fmt.Sprintf("[%d] Title: %s (%s) URL: %s\n", i+1, s.Title, s.Source, s.URL))
	}

	reqBody := OllamaRequest{
		Model:  p.Model,
		Prompt: sb.String(),
		Stream: false,
	}

	bodyBytes, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(p.Endpoint+"/api/generate", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ollamaResp OllamaResponse
	if err := json.NewDecoder(resp.Body).Decode(&ollamaResp); err != nil {
		return nil, err
	}

	text := strings.TrimSpace(ollamaResp.Response)
	if idx := strings.Index(text, "["); idx != -1 {
		if idxEnd := strings.LastIndex(text, "]"); idxEnd != -1 && idxEnd > idx {
			text = text[idx : idxEnd+1]
		}
	}

	var summaries []string
	if err := json.Unmarshal([]byte(text), &summaries); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v, text: %s", err, text)
	}

	if len(summaries) != len(stories) {
		return nil, fmt.Errorf("expected %d summaries, got %d", len(stories), len(summaries))
	}

	return summaries, nil
}
