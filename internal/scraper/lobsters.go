package scraper

import (
	"net/http"
	"news-terminal/internal/model"

	"github.com/PuerkitoBio/goquery"
)

func FetchLobsters(limit int) ([]model.Story, error) {
	resp, err := http.Get("https://lobste.rs/")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var stories []model.Story
	doc.Find(".story").Each(func(i int, s *goquery.Selection) {
		if i >= limit {
			return
		}

		title := s.Find(".link").Text()
		url, _ := s.Find(".link").Attr("href")

		stories = append(stories, model.Story{
			Title:  title,
			URL:    url,
			Source: "Lobste.rs",
		})
	})

	return stories, nil
}
