-- Seed data: genres, artists, tracks, playlists, playlist_tracks
-- Dependencies: creates a seed owner user if user id=1 does not exist
-- Order: genres → artists → tracks → track_genres → playlists → playlist_tracks

-- Genres
INSERT INTO genres (name) VALUES
    ('Lo-Fi Hip Hop'),
    ('Chill Hop'),
    ('Jazz Hop'),
    ('Electronic'),
    ('Ambient'),
    ('Jazz'),
    ('Indie Rock'),
    ('Indie Pop'),
    ('Hip Hop'),
    ('Neo-Soul')
ON CONFLICT (name) DO NOTHING;

-- Artists
INSERT INTO artists (name, bio, image_url, spotify_url) VALUES
    -- Lo-Fi / Chill Hop
    ('Ideal', 'Lo-fi producer crafting late-night beats for study and relaxation.', 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=400', 'https://open.spotify.com/artist/lofi1'),
    ('Soul Drift', 'Ambient lo-fi artist blending soul samples with dusty drums.', 'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=400', 'https://open.spotify.com/artist/lofi2'),
    ('Chillwave', 'Chill hop producer known for dreamy textures and mellow flows.', 'https://images.unsplash.com/photo-1514320291840-2e0a9bf2a9ae?w=400', 'https://open.spotify.com/artist/lofi3'),
    ('Mellow Key', 'Japanese lo-fi beatmaker with a distinct warm and nostalgic sound.', 'https://images.unsplash.com/photo-1487180144351-b8472da7d491?w=400', 'https://open.spotify.com/artist/lofi4'),
    ('Rainy Keys', 'Piano-driven lo-fi beats for rainy evenings and quiet mornings.', 'https://images.unsplash.com/photo-1520523839897-bd0b52f945a0?w=400', 'https://open.spotify.com/artist/lofi5'),
    -- Hip Hop
    ('Flow State', 'Conscious hip-hop with jazz-influenced production and lyrical depth.', 'https://images.unsplash.com/photo-1508700115892-45ecd05ae2ad?w=400', 'https://open.spotify.com/artist/hiphop1'),
    ('Rhyme Lab', 'Experimental hip-hop blending abstract beats with vivid storytelling.', 'https://images.unsplash.com/photo-1571330735066-03aaa9429d89?w=400', 'https://open.spotify.com/artist/hiphop2'),
    ('Bars & Beats', 'East Coast boom-bap rapper with soulful, classic feel.', 'https://images.unsplash.com/photo-1546527868-ccb7ee7dfa6a?w=400', 'https://open.spotify.com/artist/hiphop3'),
    ('Grasp', 'West Coast-influenced rapper with laid-back vibes and heavy bass.', 'https://images.unsplash.com/photo-1516450360452-9312f5e86fc7?w=400', 'https://open.spotify.com/artist/hiphop4'),
    -- Electronic
    ('Beat Surge', 'Electronic producer creating high-energy tracks for focus and movement.', 'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=400', 'https://open.spotify.com/artist/electronic1'),
    ('Pulse Nova', 'Deep house and synthwave artist with cinematic soundscapes.', 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=400', 'https://open.spotify.com/artist/electronic2'),
    ('Synth Wave', 'Retro-futuristic electronic with driving synths and pulsing rhythms.', 'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=400', 'https://open.spotify.com/artist/electronic3'),
    -- Jazz
    ('Blue Note Trio', 'Contemporary jazz trio pushing the boundaries of modern improvisation.', 'https://images.unsplash.com/photo-1415201364774-f6f0bb35f28f?w=400', 'https://open.spotify.com/artist/jazz1'),
    ('Echoes Quartet', 'Post-bop jazz ensemble blending tradition with modern sensibility.', 'https://images.unsplash.com/photo-1507838153414-b4b713384a76?w=400', 'https://open.spotify.com/artist/jazz2'),
    -- Rock
    ('Midnight Echo', 'Indie rock band with reverb-drenched guitars and atmospheric vocals.', 'https://images.unsplash.com/photo-1524368535928-5b5e00ddc76b?w=400', 'https://open.spotify.com/artist/rock1'),
    ('Velvet Riot', 'High-energy indie rock with driving hooks and anthemic choruses.', 'https://images.unsplash.com/photo-1498038432885-c6f3f1b912ee?w=400', 'https://open.spotify.com/artist/rock2'),
    ('The Wanderers', 'Indie rock collective known for expansive arrangements and emotional depth.', 'https://images.unsplash.com/photo-1501612780327-ac450701b08c?w=400', 'https://open.spotify.com/artist/rock3'),
    -- Indie Pop / Neo-Soul
    ('Soft Glow', 'Dreamy indie pop with lush production and intimate vocals.', 'https://images.unsplash.com/photo-1516280440614-37939bbacd81?w=400', 'https://open.spotify.com/artist/indie1'),
    ('Haze & Glass', 'Atmospheric indie pop blending electronic textures with acoustic warmth.', 'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=400', 'https://open.spotify.com/artist/indie2'),
    ('Golden Thread', 'Neo-soul artist blending R&B, soul, and contemporary jazz influences.', 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=400', 'https://open.spotify.com/artist/neosoul1')
ON CONFLICT (name) DO NOTHING;

-- Tracks
-- Artist IDs: 1=Ideal, 2=Soul Drift, 3=Chillwave, 4=Mellow Key, 5=Rainy Keys,
--              6=Flow State, 7=Rhyme Lab, 8=Bars & Beats, 9=Grasp,
--              10=Beat Surge, 11=Pulse Nova, 12=Synth Wave,
--              13=Blue Note Trio, 14=Echoes Quartet,
--              15=Midnight Echo, 16=Velvet Riot, 17=The Wanderers,
--              18=Soft Glow, 19=Haze & Glass, 20=Golden Thread

-- Lo-Fi Hip Hop tracks
INSERT INTO tracks (title, artist_id, duration_ms, popularity_score, spotify_url, cover_url) VALUES
    ('Night Drive', 1, 187000, 89, 'https://open.spotify.com/track/t1', 'https://images.unsplash.com/photo-1504898770365-14faca6a7320?w=600'),
    ('Coffee Shop', 1, 214000, 92, 'https://open.spotify.com/track/t2', 'https://images.unsplash.com/photo-1442975631134-6137411e9d4e?w=600'),
    ('Rainy Afternoon', 1, 195000, 85, 'https://open.spotify.com/track/t3', 'https://images.unsplash.com/photo-1428592953211-077101b2021b?w=600'),
    ('Midnight Study', 2, 223000, 87, 'https://open.spotify.com/track/t4', 'https://images.unsplash.com/photo-1516979187457-637abb4f9353?w=600'),
    ('Sunset Soul', 2, 198000, 83, 'https://open.spotify.com/track/t5', 'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=600'),
    ('Dream State', 3, 206000, 79, 'https://open.spotify.com/track/t6', 'https://images.unsplash.com/photo-1492144534655-ae79c964c9d7?w=600'),
    ('Chill Vibes', 3, 234000, 81, 'https://open.spotify.com/track/t7', 'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=600'),
    ('Tokyo Nights', 4, 189000, 88, 'https://open.spotify.com/track/t8', 'https://images.unsplash.com/photo-1480796927426-f609979314bd?w=600'),
    ('Quiet Hours', 4, 212000, 84, 'https://open.spotify.com/track/t9', 'https://images.unsplash.com/photo-1512436991641-6745cdb1723f?w=600'),
    ('Piano in the Dark', 5, 201000, 86, 'https://open.spotify.com/track/t10', 'https://images.unsplash.com/photo-1520523839897-bd0b52f945a0?w=600');

-- Hip Hop tracks
INSERT INTO tracks (title, artist_id, duration_ms, popularity_score, spotify_url, cover_url) VALUES
    ('City Lights', 6, 241000, 76, 'https://open.spotify.com/track/t11', 'https://images.unsplash.com/photo-1477959858617-67f85cf4f1df?w=600'),
    ('Flow & Fire', 6, 198000, 72, 'https://open.spotify.com/track/t12', 'https://images.unsplash.com/photo-1571330735066-03aaa9429d89?w=600'),
    ('Deep Thought', 7, 265000, 68, 'https://open.spotify.com/track/t13', 'https://images.unsplash.com/photo-1508700115892-45ecd05ae2ad?w=600'),
    ('The Lab', 7, 229000, 71, 'https://open.spotify.com/track/t14', 'https://images.unsplash.com/photo-1546527868-ccb7ee7dfa6a?w=600'),
    ('Golden Era', 8, 214000, 80, 'https://open.spotify.com/track/t15', 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=600'),
    ('Boom Bap Nights', 8, 187000, 77, 'https://open.spotify.com/track/t16', 'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=600'),
    ('West Coast Wind', 9, 203000, 74, 'https://open.spotify.com/track/t17', 'https://images.unsplash.com/photo-1504898770365-14faca6a7320?w=600'),
    ('Bass & Soul', 9, 219000, 73, 'https://open.spotify.com/track/t18', 'https://images.unsplash.com/photo-1442975631134-6137411e9d4e?w=600');

-- Electronic tracks
INSERT INTO tracks (title, artist_id, duration_ms, popularity_score, spotify_url, cover_url) VALUES
    ('Surge', 10, 312000, 75, 'https://open.spotify.com/track/t19', 'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=600'),
    ('Electric Dreams', 10, 287000, 78, 'https://open.spotify.com/track/t20', 'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=600'),
    ('Digital Rain', 11, 341000, 70, 'https://open.spotify.com/track/t21', 'https://images.unsplash.com/photo-1428592953211-077101b2021b?w=600'),
    ('Neon Horizon', 11, 298000, 72, 'https://open.spotify.com/track/t22', 'https://images.unsplash.com/photo-1492144534655-ae79c964c9d7?w=600'),
    ('Retro Wave', 12, 264000, 82, 'https://open.spotify.com/track/t23', 'https://images.unsplash.com/photo-1477959858617-67f85cf4f1df?w=600'),
    ('Cyber City', 12, 319000, 79, 'https://open.spotify.com/track/t24', 'https://images.unsplash.com/photo-1480796927426-f609979314bd?w=600');

-- Jazz tracks
INSERT INTO tracks (title, artist_id, duration_ms, popularity_score, spotify_url, cover_url) VALUES
    ('Blue Horizons', 13, 356000, 65, 'https://open.spotify.com/track/t25', 'https://images.unsplash.com/photo-1415201364774-f6f0bb35f28f?w=600'),
    ('Standards', 13, 412000, 62, 'https://open.spotify.com/track/t26', 'https://images.unsplash.com/photo-1507838153414-b4b713384a76?w=600'),
    ('After Hours', 14, 389000, 60, 'https://open.spotify.com/track/t27', 'https://images.unsplash.com/photo-1516979187457-637abb4f9353?w=600'),
    ('Echoes of Now', 14, 334000, 58, 'https://open.spotify.com/track/t28', 'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=600');

-- Indie Rock tracks
INSERT INTO tracks (title, artist_id, duration_ms, popularity_score, spotify_url, cover_url) VALUES
    ('Echo Chamber', 15, 278000, 71, 'https://open.spotify.com/track/t29', 'https://images.unsplash.com/photo-1524368535928-5b5e00ddc76b?w=600'),
    ('Static Dreams', 15, 301000, 68, 'https://open.spotify.com/track/t30', 'https://images.unsplash.com/photo-1498038432885-c6f3f1b912ee?w=600'),
    ('Reverb Road', 15, 245000, 74, 'https://open.spotify.com/track/t31', 'https://images.unsplash.com/photo-1501612780327-ac450701b08c?w=600'),
    ('Run Wild', 16, 223000, 76, 'https://open.spotify.com/track/t32', 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=600'),
    ('Break the Frame', 16, 198000, 73, 'https://open.spotify.com/track/t33', 'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=600'),
    ('Open Sky', 17, 312000, 69, 'https://open.spotify.com/track/t34', 'https://images.unsplash.com/photo-1504898770365-14faca6a7320?w=600'),
    ('Wanderlust', 17, 287000, 66, 'https://open.spotify.com/track/t35', 'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=600');

-- Indie Pop / Neo-Soul tracks
INSERT INTO tracks (title, artist_id, duration_ms, popularity_score, spotify_url, cover_url) VALUES
    ('Soft Light', 18, 209000, 80, 'https://open.spotify.com/track/t36', 'https://images.unsplash.com/photo-1516280440614-37939bbacd81?w=600'),
    ('Morning Glow', 18, 234000, 77, 'https://open.spotify.com/track/t37', 'https://images.unsplash.com/photo-1512436991641-6745cdb1723f?w=600'),
    ('Glass House', 19, 218000, 75, 'https://open.spotify.com/track/t38', 'https://images.unsplash.com/photo-1442975631134-6137411e9d4e?w=600'),
    ('Haze', 19, 241000, 72, 'https://open.spotify.com/track/t39', 'https://images.unsplash.com/photo-1492144534655-ae79c964c9d7?w=600'),
    ('Golden', 20, 256000, 78, 'https://open.spotify.com/track/t40', 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=600'),
    ('Soul Shine', 20, 223000, 81, 'https://open.spotify.com/track/t41', 'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=600');

-- Track genres (track_id → genre_id)
-- Genres: 1=Lo-Fi Hip Hop, 2=Chill Hop, 3=Jazz Hop, 4=Electronic, 5=Ambient, 6=Jazz, 7=Indie Rock, 8=Indie Pop, 9=Hip Hop, 10=Neo-Soul
-- Tracks 1-10: Lo-Fi artists (mostly 1=Lo-Fi Hip Hop, some 2=Chill Hop)
INSERT INTO track_genres (track_id, genre_id) VALUES
    -- Lo-Fi tracks
    (1, 1), (1, 2),
    (2, 1), (2, 2),
    (3, 1),
    (4, 1), (4, 3),
    (5, 2), (5, 3),
    (6, 2),
    (7, 2),
    (8, 1), (8, 2),
    (9, 1),
    (10, 1), (10, 5),
    -- Hip Hop tracks
    (11, 9), (11, 3),
    (12, 9),
    (13, 9), (13, 3),
    (14, 9),
    (15, 9), (15, 3),
    (16, 9),
    (17, 9),
    (18, 9), (18, 10),
    -- Electronic tracks
    (19, 4),
    (20, 4),
    (21, 4), (21, 5),
    (22, 4),
    (23, 4), (23, 5),
    (24, 4),
    -- Jazz tracks
    (25, 6), (25, 3),
    (26, 6),
    (27, 6), (27, 3),
    (28, 6),
    -- Indie Rock tracks
    (29, 7),
    (30, 7), (30, 5),
    (31, 7),
    (32, 7),
    (33, 7),
    (34, 7),
    (35, 7),
    -- Indie Pop / Neo-Soul tracks
    (36, 8), (36, 5),
    (37, 8),
    (38, 8), (38, 4),
    (39, 8), (39, 5),
    (40, 10), (40, 8),
    (41, 10), (41, 8);

-- Seed owner for public demo playlists. Password is a local fixture hash.
INSERT INTO users (id, email, password_hash) VALUES
    (1, 'seed@music.local', '$2a$10$npUgIfKQZ/otVbHBxtuU8eaEOJjcDKyYXvo7AumcDSFHUf0cAuiVC')
ON CONFLICT (id) DO NOTHING;

SELECT setval(pg_get_serial_sequence('users', 'id'), GREATEST((SELECT MAX(id) FROM users), 1), true);

-- Playlists (user_id=1 is the seeded public playlist owner)
INSERT INTO playlists (user_id, name, description, is_public) VALUES
    (1, 'Lo-Fi Hip Hop Mix', 'Chill beats to study, relax, and focus.', true),
    (1, 'Late Night Jazz', 'Smooth jazz for winding down after dark.', true),
    (1, 'Electronic Study', 'High-energy electronic for deep work sessions.', true),
    (1, 'Indie Rock Revival', 'Driving guitar sounds for active minds.', true),
    (1, 'Neo-Soul & Chill', 'Smooth neo-soul and laid-back R&B vibes.', true),
    (1, 'Favorites', 'Your most-loved tracks in one place.', false),
    (1, 'Discover: Hip Hop', 'Explore the best in conscious and experimental hip-hop.', true);

-- Playlist tracks (playlist_id → track_id, position)
-- Playlist 1: Lo-Fi Hip Hop Mix (tracks 1-5, 8-10)
INSERT INTO playlist_tracks (playlist_id, track_id, position) VALUES
    (1, 1, 1), (1, 2, 2), (1, 3, 3), (1, 4, 4), (1, 5, 5), (1, 8, 6), (1, 9, 7), (1, 10, 8);

-- Playlist 2: Late Night Jazz (tracks 25-28)
INSERT INTO playlist_tracks (playlist_id, track_id, position) VALUES
    (2, 25, 1), (2, 26, 2), (2, 27, 3), (2, 28, 4);

-- Playlist 3: Electronic Study (tracks 19-24)
INSERT INTO playlist_tracks (playlist_id, track_id, position) VALUES
    (3, 19, 1), (3, 20, 2), (3, 21, 3), (3, 22, 4), (3, 23, 5), (3, 24, 6);

-- Playlist 4: Indie Rock Revival (tracks 29-35)
INSERT INTO playlist_tracks (playlist_id, track_id, position) VALUES
    (4, 29, 1), (4, 30, 2), (4, 31, 3), (4, 32, 4), (4, 33, 5), (4, 34, 6), (4, 35, 7);

-- Playlist 5: Neo-Soul & Chill (tracks 40-41, 36-39)
INSERT INTO playlist_tracks (playlist_id, track_id, position) VALUES
    (5, 40, 1), (5, 41, 2), (5, 36, 3), (5, 37, 4), (5, 38, 5), (5, 39, 6);

-- Playlist 6: Favorites (curated selection)
INSERT INTO playlist_tracks (playlist_id, track_id, position) VALUES
    (6, 2, 1), (6, 8, 2), (6, 15, 3), (6, 23, 4), (6, 36, 5), (6, 41, 6);

-- Playlist 7: Discover: Hip Hop (tracks 11-18)
INSERT INTO playlist_tracks (playlist_id, track_id, position) VALUES
    (7, 11, 1), (7, 12, 2), (7, 13, 3), (7, 14, 4), (7, 15, 5), (7, 16, 6), (7, 17, 7), (7, 18, 8);
