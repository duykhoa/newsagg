CREATE TABLE IF NOT EXISTS sources (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL UNIQUE,
	last_refreshed_at DATETIME,
	refresh_interval INTEGER
);

CREATE TABLE IF NOT EXISTS stories (
	id INTEGER PRIMARY KEY,
	source_id INTEGER,
	title TEXT NOT NULL,
	url TEXT NOT NULL,
	story_url TEXT UNIQUE NOT NULL,
	points INTEGER,
	published_at DATETIME,
	raw_content TEXT,
	recap TEXT,
	full_content_tui TEXT,
	FOREIGN KEY (source_id) REFERENCES sources(id)
);

CREATE INDEX IF NOT EXISTS idx_stories_source_id ON stories(source_id);
