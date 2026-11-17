package main

import (
	"fmt"
	"news-terminal/internal/recap"
	"news-terminal/internal/scraper"
)

func main() {
	fmt.Println("--- Fetching News ---")

	hnStories, err := scraper.FetchHackerNews(5)
	if err != nil {
		fmt.Println("Error fetching HN:", err)
	}

	lbStories, err := scraper.FetchLobsters(5)
	if err != nil {
		fmt.Println("Error fetching Lobsters:", err)
	}

	allStories := append(hnStories, lbStories...)

	for _, s := range allStories {
		recap.GenerateSummary(&s)
		fmt.Printf("[%s] %s\n", s.Source, s.Title)
		fmt.Printf("    Link: %s\n", s.URL)
		fmt.Printf("    Recap: %s\n\n", s.Recap)
	}
}
