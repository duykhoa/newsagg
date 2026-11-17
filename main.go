package main

import (
	"fmt"
	"strings"

	"news-terminal/internal/config"
	"news-terminal/internal/recap"
	"news-terminal/internal/scraper"
)

func main() {
	fmt.Println("--- Loading Config ---")
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Println("Error loading config, using defaults:", err)
	}

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

	fmt.Println("--- Generating Recaps ---")
	if err := recap.GenerateSummaries(cfg.Provider, cfg.OllamaEndpoint, cfg.OllamaModel, cfg.GeminiAPIKey, cfg.GeminiModel, allStories, 3); err != nil {
		fmt.Println("Error generating summaries:", err)
	}

	for i, s := range allStories {
		fmt.Printf("[%s] (%d points) %s\n", s.Source, s.Points, s.Title)
		fmt.Printf("  Link:  %s\n", s.URL)
		fmt.Printf("  Recap: %s\n", s.Recap)
		if i < len(allStories)-1 {
			fmt.Println(strings.Repeat("-", 60))
		}
	}
}
