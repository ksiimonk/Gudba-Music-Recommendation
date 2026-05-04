CREATE TABLE IF NOT EXISTS events (
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entity_type VARCHAR NOT NULL,
    entity_id INTEGER NOT NULL,
    event_type VARCHAR NOT NULL,
    session_id VARCHAR,
    metadata JSONB,
    created_at TIMESTAMPTZ DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_events_user_id ON events(user_id, created_at DESC);
