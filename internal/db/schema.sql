CREATE TABLE IF NOT EXISTS sources (
	id INTEGER PRIMARY KEY,
	name TEXT NOT NULL,
	last_refreshed_at DATETIME,
	refresh_interval INTEGER DEFAULT 14400,
	CONSTRAINT sources_name_key UNIQUE (name)
);

CREATE TABLE IF NOT EXISTS stories (
	id INTEGER PRIMARY KEY,
	source_id INTEGER,
	title TEXT NOT NULL,
	url TEXT NOT NULL,
	story_url TEXT NOT NULL,
	points INTEGER,
	published_at DATETIME,
	raw_content TEXT,
	recap TEXT,
	full_content_tui TEXT,
	created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (source_id) REFERENCES sources(id),
	CONSTRAINT stories_story_url_key UNIQUE (story_url)
);
