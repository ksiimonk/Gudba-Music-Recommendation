CREATE TABLE artists (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    image_url TEXT,
    bio TEXT,
    spotify_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX idx_artists_name ON artists (name);
