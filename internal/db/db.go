package db

import (
	"context"
	"database/sql"
	"time"

	"news-terminal/internal/model"

	_ "modernc.org/sqlite"
)

type DB struct {
	conn    *sql.DB
	Queries *Queries
}

func NewDB(filepath string) (*DB, error) {
	conn, err := sql.Open("sqlite", filepath)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		return nil, err
	}

	database := &DB{
		conn:    conn,
		Queries: New(conn),
	}

	return database, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) UpsertSource(name string, refreshInterval int) (int64, error) {
	ctx := context.Background()
	err := db.Queries.UpsertSource(ctx, UpsertSourceParams{
		Name:            name,
		RefreshInterval: int64(refreshInterval),
	})
	if err != nil {
		return 0, err
	}

	id, err := db.Queries.GetSourceIDByName(ctx, name)
	return id, err
}

func (db *DB) UpdateSourceRefreshTime(sourceID int64, t time.Time) error {
	ctx := context.Background()
	return db.Queries.UpdateSourceRefreshTime(ctx, UpdateSourceRefreshTimeParams{
		LastRefreshedAt: sql.NullTime{Time: t, Valid: true},
		ID:              sourceID,
	})
}

func (db *DB) SaveStory(story model.Story, sourceID int64) error {
	ctx := context.Background()
	storyURL := story.StoryURL
	if storyURL == "" {
		storyURL = story.URL
	}
	return db.Queries.SaveStory(ctx, SaveStoryParams{
		SourceID:       sql.NullInt64{Int64: sourceID, Valid: true},
		Title:          story.Title,
		Url:            story.URL,
		StoryUrl:       storyURL,
		Points:         sql.NullInt64{Int64: int64(story.Points), Valid: true},
		Recap:          sql.NullString{String: story.Recap, Valid: story.Recap != ""},
		FullContentTui: sql.NullString{String: story.FullContentTUI, Valid: story.FullContentTUI != ""},
	})
}

func (db *DB) GetStoriesBySource(sourceID int64, limit, offset int) ([]model.Story, error) {
	ctx := context.Background()
	rows, err := db.Queries.GetStoriesBySource(ctx, GetStoriesBySourceParams{
		SourceID: sql.NullInt64{Int64: sourceID, Valid: true},
		Limit:    int64(limit),
		Offset:   int64(offset),
	})
	if err != nil {
		return nil, err
	}

	var stories []model.Story
	for _, row := range rows {
		stories = append(stories, model.Story{
			Title:          row.Title,
			URL:            row.Url,
			StoryURL:       row.StoryUrl,
			Points:         int(row.Points.Int64),
			Recap:          row.Recap.String,
			FullContentTui: row.FullContentTui.String,
		})
	}
	return stories, nil
}
