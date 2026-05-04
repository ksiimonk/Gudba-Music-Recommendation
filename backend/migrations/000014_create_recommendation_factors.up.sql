CREATE TABLE IF NOT EXISTS recommendation_factors (
    id SERIAL PRIMARY KEY,
    impression_id INTEGER NOT NULL REFERENCES recommendation_impressions(id) ON DELETE CASCADE,
    factor_name VARCHAR NOT NULL,
    factor_value NUMERIC,
    weight NUMERIC,
    contribution NUMERIC,
    created_at TIMESTAMPTZ DEFAULT now()
);
