package scraper

import (
	"fmt"
	"net/http"
	"strings"

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

		linkSel := s.Find("a.u-url")
		title := strings.TrimSpace(linkSel.Text())
		url, _ := linkSel.Attr("href")
		scoreSel := s.Find(".voters a.upvoter, .score")
		pointsStr := strings.TrimSpace(scoreSel.Text())
		if pointsStr == "" {
			if titleAttr, exists := scoreSel.Attr("title"); exists {
				pointsStr = titleAttr
			}
		}
		var points int
		fmt.Sscanf(pointsStr, "%d", &points)

		stories = append(stories, model.Story{
			Title:  title,
			URL:    url,
			Source: "Lobste.rs",
			Points: points,
		})
	})

	return stories, nil
}
