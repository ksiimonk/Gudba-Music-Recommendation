-- Rollback seed data (reverse dependency order)
DELETE FROM playlist_tracks WHERE playlist_id IN (1, 2, 3, 4, 5, 6, 7);
DELETE FROM playlists WHERE user_id = 1;
DELETE FROM track_genres WHERE track_id BETWEEN 1 AND 41;
DELETE FROM tracks WHERE artist_id BETWEEN 1 AND 20;
DELETE FROM artists WHERE id BETWEEN 1 AND 20;
DELETE FROM genres WHERE id BETWEEN 1 AND 10;
DELETE FROM users WHERE id = 1 AND email = 'seed@music.local';
