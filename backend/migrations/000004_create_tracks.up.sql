CREATE TABLE tracks (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    artist_id BIGINT NOT NULL REFERENCES artists (id) ON DELETE CASCADE,
    duration_ms INTEGER NOT NULL,
    spotify_url TEXT,
    cover_url TEXT,
    popularity_score INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_tracks_artist_id ON tracks (artist_id);
CREATE INDEX idx_tracks_popularity_score ON tracks (popularity_score DESC);
