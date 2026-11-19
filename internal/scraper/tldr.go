package scraper

import (
	"net/http"
	"strings"

	"news-terminal/internal/model"

	"github.com/PuerkitoBio/goquery"
)

func FetchTLDR(limit int) ([]model.Story, error) {
	resp, err := http.Get("https://tldr.tech/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var stories []model.Story
	doc.Find("h3 a").Each(func(i int, s *goquery.Selection) {
		if i >= limit {
			return
		}

		title := strings.TrimSpace(s.Text())
		url, exists := s.Attr("href")
		if !exists || title == "" {
			return
		}

		if !strings.HasPrefix(url, "http") {
			url = "https://tldr.tech" + url
		}

		stories = append(stories, model.Story{
			Title:    title,
			URL:      url,
			StoryURL: url,
			Source:   "TLDR",
			Points:   0,
		})
	})

	return stories, nil
}
