CREATE TABLE IF NOT EXISTS recommendation_requests (
    id SERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    request_type VARCHAR NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now()
);
