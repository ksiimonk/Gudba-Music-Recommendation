CREATE TABLE IF NOT EXISTS user_profiles (
    user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    favorite_genre_ids JSONB DEFAULT '[]',
    favorite_artist_ids JSONB DEFAULT '[]',
    starter_track_ids JSONB DEFAULT '[]',
    contexts JSONB DEFAULT '[]',
    updated_at TIMESTAMPTZ DEFAULT now()
);
