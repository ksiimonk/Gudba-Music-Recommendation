CREATE TABLE track_genres (
    track_id BIGINT NOT NULL REFERENCES tracks (id) ON DELETE CASCADE,
    genre_id BIGINT NOT NULL REFERENCES genres (id) ON DELETE CASCADE,
    PRIMARY KEY (track_id, genre_id)
);

CREATE INDEX idx_track_genres_genre_id ON track_genres (genre_id);