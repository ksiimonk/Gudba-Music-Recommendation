CREATE TABLE IF NOT EXISTS recommendation_impressions (
    id SERIAL PRIMARY KEY,
    request_id BIGINT NOT NULL REFERENCES recommendation_requests(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    entity_type VARCHAR NOT NULL,
    entity_id BIGINT NOT NULL,
    rank INTEGER NOT NULL,
    score NUMERIC NOT NULL,
    explanation VARCHAR,
    created_at TIMESTAMPTZ DEFAULT now()
);
