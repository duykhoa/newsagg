package worker

import (
	"context"
	"log"
	"time"

	"news-terminal/internal/config"
	"news-terminal/internal/db"
	"news-terminal/internal/recap"
	"news-terminal/internal/scraper"
)

type Worker struct {
	db  *db.DB
	cfg *config.Config
}

func NewWorker(database *db.DB, cfg *config.Config) *Worker {
	return &Worker{
		db:  database,
		cfg: cfg,
	}
}

func (w *Worker) Start(ctx context.Context, interval time.Duration) {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				w.RunRefreshCycle()
			}
		}
	}()
}

func (w *Worker) RunRefreshCycle() {
	log.Println("Starting background refresh cycle...")

	sourceID, err := w.db.UpsertSource("Hacker News", 14400)
	if err != nil {
		log.Printf("Failed to upsert source: %v\n", err)
		return
	}

	stories, err := scraper.FetchHackerNews(w.cfg.BatchSize)
	if err != nil {
		log.Printf("Failed to fetch stories: %v\n", err)
		return
	}

	provider, err := recap.NewProvider(w.cfg.Provider, w.cfg.OllamaEndpoint, w.cfg.OllamaModel, w.cfg.GeminiAPIKey, w.cfg.GeminiModel)
	if err != nil {
		log.Printf("Failed to init LLM provider: %v\n", err)
		return
	}

	for _, story := range stories {
		if story.URL != "" {
			if raw, err := scraper.FetchRawContent(story.URL); err == nil {
				if tuiContent, err := provider.GenerateFullContentTUI(story.Title, story.URL, raw); err == nil {
					story.FullContentTUI = tuiContent
				}
			}
		}

		if err := w.db.SaveStory(story, sourceID); err != nil {
			log.Printf("Failed to save story %s: %v\n", story.URL, err)
		}
	}

	_ = w.db.UpdateSourceRefreshTime(sourceID, time.Now())
	log.Println("Background refresh cycle completed.")
}
