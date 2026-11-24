package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	_ "modernc.org/sqlite"

	"news-terminal/internal/config"
	"news-terminal/internal/db"
	"news-terminal/internal/model"
	"news-terminal/internal/recap"
	"news-terminal/internal/scraper"
	"news-terminal/internal/server"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Printf("Error loading config, using defaults: %v", err)
	}

	conn, err := sql.Open("sqlite", cfg.DBName)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer conn.Close()

	schemaBytes, err := os.ReadFile("internal/db/schema.sql")
	if err != nil {
		log.Fatalf("Failed to read schema file: %v", err)
	}

	if _, err := conn.Exec(string(schemaBytes)); err != nil {
		log.Fatalf("Failed to apply database schema: %v", err)
	}

	dbInstance := db.NewDB(conn)
	srv := server.NewServer(dbInstance)

	ctx := context.Background()
	sources, err := srv.GetSources(ctx)
	if err != nil || len(sources) == 0 {
		log.Println("Performing initial startup sync...")
		runInitialSync(dbInstance, cfg)
	}

	log.Println("Server ready with initial sync complete.")
}

func runInitialSync(db *db.DB, cfg *config.Config) {
	syncSources := []struct {
		name    string
		fetchFn func(int) ([]model.Story, error)
	}{
		{"Hacker News", scraper.FetchHackerNews},
		{"Lobsters", scraper.FetchLobsters},
		{"TLDR", scraper.FetchTLDR},
	}

	for _, src := range syncSources {
		stories, err := src.fetchFn(cfg.BatchSize)
		if err != nil {
			log.Printf("Failed to fetch from %s: %v", src.name, err)
			continue
		}

		sourceID, err := db.UpsertSource(src.name, 14400)
		if err != nil {
			log.Printf("Failed to upsert source %s: %v", src.name, err)
			continue
		}

		_ = recap.GenerateSummaries(cfg.Provider, cfg.OllamaEndpoint, cfg.OllamaModel, cfg.GeminiAPIKey, cfg.GeminiModel, stories, cfg.BatchSize)
		for _, s := range stories {
			_ = db.SaveStory(s, sourceID)
		}
		_ = db.UpdateSourceRefreshTime(sourceID, time.Now())
	}
}
