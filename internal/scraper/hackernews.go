package scraper

import (
	"encoding/json"
	"fmt"
	"net/http"
	"news-terminal/internal/model"
)

type HNItem struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Score int    `json:"score"`
}

func FetchHackerNews(limit int) ([]model.Story, error) {
	resp, err := http.Get("https://hacker-news.firebaseio.com/v0/topstories.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var ids []int
	if err := json.NewDecoder(resp.Body).Decode(&ids); err != nil {
		return nil, err
	}

	var stories []model.Story
	for i := 0; i < limit && i < len(ids); i++ {
		itemURL := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", ids[i])
		resp, err := http.Get(itemURL)
		if err != nil {
			continue
		}

		var item HNItem
		if err := json.NewDecoder(resp.Body).Decode(&item); err == nil && item.Title != "" {
			stories = append(stories, model.Story{
				Title:  item.Title,
				URL:    item.URL,
				Source: "Hacker News",
				Points: item.Score,
			})
		}
		resp.Body.Close()
	}
	return stories, nil
}
