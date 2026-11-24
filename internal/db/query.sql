-- name: UpsertSource :exec
INSERT INTO sources (name, refresh_interval) VALUES (?, ?)
ON CONFLICT (name) DO UPDATE SET refresh_interval = excluded.refresh_interval;

-- name: GetSourceIDByName :one
SELECT id FROM sources WHERE name = ?;

-- name: UpdateSourceRefreshTime :exec
UPDATE sources SET last_refreshed_at = ? WHERE id = ?;

-- name: SaveStory :exec
INSERT INTO stories (source_id, title, url, story_url, points, recap, full_content_tui)
VALUES (?, ?, ?, ?, ?, ?, ?)
ON CONFLICT (story_url) DO UPDATE SET points = excluded.points, recap = COALESCE(excluded.recap, stories.recap);

-- name: GetStoriesBySource :many
SELECT title, url, story_url, points, recap, full_content_tui FROM stories WHERE source_id = ? ORDER BY id DESC LIMIT ? OFFSET ?;

-- name: GetSources :many
SELECT id, name, last_refreshed_at, refresh_interval FROM sources;

-- name: GetStoryDetail :one
SELECT id, source_id, title, url, story_url, points, published_at, raw_content, recap, full_content_tui FROM stories WHERE id = ?;
