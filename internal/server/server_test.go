package server

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"

	"news-terminal/internal/db"
	"news-terminal/internal/model"
)

func setupTestDB(t *testing.T) *db.DB {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}

	_, err = conn.Exec(`
		CREATE TABLE sources (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL UNIQUE,
			last_refreshed_at DATETIME,
			refresh_interval INTEGER
		);
		CREATE TABLE stories (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			source_id INTEGER,
			title TEXT NOT NULL,
			url TEXT NOT NULL,
			story_url TEXT NOT NULL UNIQUE,
			points INTEGER,
			published_at DATETIME,
			raw_content TEXT,
			recap TEXT,
			full_content_tui TEXT,
			FOREIGN KEY(source_id) REFERENCES sources(id)
		);
	`)
	if err != nil {
		t.Fatalf("failed to create schema: %v", err)
	}

	return db.NewDB(conn)
}

func TestServer_GetSources(t *testing.T) {
	database := setupTestDB(t)
	srv := NewServer(database)

	ctx := context.Background()
	_, err := database.UpsertSource("Hacker News", 14400)
	if err != nil {
		t.Fatalf("failed to insert source: %v", err)
	}

	sources, err := srv.GetSources(ctx)
	if err != nil {
		t.Fatalf("GetSources failed: %v", err)
	}

	if len(sources) != 1 || sources[0].Name != "Hacker News" {
		t.Errorf("unexpected sources: %+v", sources)
	}
}

func TestServer_GetStoriesBySource(t *testing.T) {
	database := setupTestDB(t)
	srv := NewServer(database)

	sourceID, err := database.UpsertSource("Hacker News", 14400)
	if err != nil {
		t.Fatalf("failed to insert source: %v", err)
	}

	story := model.Story{
		Title:          "Test Story",
		URL:            "https://example.com",
		StoryURL:       "https://news.ycombinator.com/item?id=1",
		Points:         100,
		Recap:          "Test recap",
		FullContentTui: "Full content",
	}
	if err := database.SaveStory(story, sourceID); err != nil {
		t.Fatalf("failed to save story: %v", err)
	}

	ctx := context.Background()
	stories, err := srv.GetStoriesBySource(ctx, sourceID, 10, 0)
	if err != nil {
		t.Fatalf("GetStoriesBySource failed: %v", err)
	}

	if len(stories) != 1 || stories[0].Title != "Test Story" {
		t.Errorf("unexpected stories: %+v", stories)
	}
}
