package db

import (
	"context"
	"database/sql"
	"time"

	"news-terminal/internal/model"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

type Queries struct {
	db DBTX
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db: tx,
	}
}

type DB struct {
	*Queries
	Conn *sql.DB
}

func NewDB(conn *sql.DB) *DB {
	return &DB{
		Queries: New(conn),
		Conn:    conn,
	}
}

func (db *DB) UpsertSource(name string, refreshInterval int) (int64, error) {
	ctx := context.Background()
	err := db.Queries.UpsertSource(ctx, UpsertSourceParams{
		Name:            name,
		RefreshInterval: sql.NullInt64{Int64: int64(refreshInterval), Valid: true},
	})
	if err != nil {
		return 0, err
	}
	return db.Queries.GetSourceIDByName(ctx, name)
}

func (db *DB) SaveStory(story model.Story, sourceID int64) error {
	ctx := context.Background()
	return db.Queries.SaveStory(ctx, SaveStoryParams{
		SourceID:       sql.NullInt64{Int64: sourceID, Valid: true},
		Title:          story.Title,
		Url:            story.URL,
		StoryUrl:       story.StoryURL,
		Points:         sql.NullInt64{Int64: int64(story.Points), Valid: true},
		Recap:          sql.NullString{String: story.Recap, Valid: story.Recap != ""},
		FullContentTui: sql.NullString{String: story.FullContentTui, Valid: story.FullContentTui != ""},
	})
}

func (db *DB) UpdateSourceRefreshTime(sourceID int64, t time.Time) error {
	ctx := context.Background()
	return db.Queries.UpdateSourceRefreshTime(ctx, UpdateSourceRefreshTimeParams{
		LastRefreshedAt: sql.NullTime{Time: t, Valid: true},
		ID:              sourceID,
	})
}
