package server

import (
	"context"
	"database/sql"

	"news-terminal/internal/db"
)

type Server struct {
	db *db.DB
}

func NewServer(database *db.DB) *Server {
	return &Server{db: database}
}

func (s *Server) GetSources(ctx context.Context) ([]db.Source, error) {
	return s.db.GetSources(ctx)
}

func (s *Server) GetStoriesBySource(ctx context.Context, sourceID int64, limit, offset int64) ([]db.GetStoriesBySourceRow, error) {
	return s.db.GetStoriesBySource(ctx, db.GetStoriesBySourceParams{
		SourceID: sql.NullInt64{Int64: sourceID, Valid: true},
		Limit:    limit,
		Offset:   offset,
	})
}

func (s *Server) GetStoryDetail(ctx context.Context, storyID int64) (db.GetStoryDetailRow, error) {
	return s.db.GetStoryDetail(ctx, storyID)
}
