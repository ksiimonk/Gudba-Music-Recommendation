import { useEffect, useMemo, useState } from 'react'
import { getProfile, listArtists, listGenres, listPlaylists, listPlaylistRecommendations, listTrackRecommendations, listTracks } from '../api'
import { useAuth } from '../auth/AuthContext'
import { usePlayer } from '../context/PlayerContext'
import { AppShell } from '../components/AppShell'
import { PlaylistCard } from '../components/PlaylistCard'
import { TrackRow } from '../components/TrackRow'
import { ROUTES } from '../config/routes'
import { PLAYLIST_COVER_URLS, getPlaylistCoverUrl } from '../data/musicContent'
import type { Artist, Genre, Playlist, Track } from '../types'
import type { PlaylistRecommendation, TrackRecommendation } from '../api/recommendations'

export function HomePage() {
  const { isLoading: authLoading, token, user } = useAuth()
  const { playTrack, toggleLike, skipTrack } = usePlayer()

  const [genres, setGenres] = useState<Genre[]>([])
  const [artists, setArtists] = useState<Artist[]>([])
  const [playlists, setPlaylists] = useState<Playlist[]>([])
  const [tracks, setTracks] = useState<Track[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')
  const [hasProfile, setHasProfile] = useState(false)
  const [profileChecking, setProfileChecking] = useState(true)

  const [recommendations, setRecommendations] = useState<TrackRecommendation[]>([])
  const [recsLoading, setRecsLoading] = useState(false)
  const [recsError, setRecsError] = useState(false)

  const [playlistRecs, setPlaylistRecs] = useState<PlaylistRecommendation[]>([])
  const [plRecsLoading, setPlRecsLoading] = useState(false)

  useEffect(() => {
    let isMounted = true

    Promise.all([
      listGenres(),
      listArtists(),
      listPlaylists(),
      listTracks(),
    ])
      .then(([genresRes, artistsRes, playlistsRes, tracksRes]) => {
        if (isMounted) {
          setGenres(genresRes.genres)
          setArtists(artistsRes.artists)
          setPlaylists(playlistsRes.playlists)
          setTracks(tracksRes.tracks)
        }
      })
      .catch(() => {
        if (isMounted) setError('Не удалось загрузить данные. Проверь, запущен ли backend.')
      })
      .finally(() => {
        if (isMounted) setIsLoading(false)
      })

    return () => { isMounted = false }
  }, [])

  useEffect(() => {
    if (authLoading || !token) {
      setProfileChecking(false)
      setHasProfile(false)
      return
    }

    let isMounted = true

    getProfile(token)
      .then((res) => {
        if (isMounted) {
          const p = res.profile
          const has = p.favorite_genre_ids.length > 0 || p.favorite_artist_ids.length > 0 || p.starter_track_ids.length > 0
          setHasProfile(has)

          if (has) {
            setRecsLoading(true)
            setPlRecsLoading(true)

            listTrackRecommendations(token, 6)
              .then((recRes) => {
                if (isMounted) {
                  setRecommendations(recRes.recommendations)
                  setRecsError(false)
                }
              })
              .catch(() => {
                if (isMounted) setRecsError(true)
              })
              .finally(() => {
                if (isMounted) setRecsLoading(false)
              })

            listPlaylistRecommendations(token, 6)
              .then((plRes) => {
                if (isMounted) setPlaylistRecs(plRes.recommendations)
              })
              .catch(() => {})
              .finally(() => {
                if (isMounted) setPlRecsLoading(false)
              })
          }
        }
      })
      .catch(() => {
        if (isMounted) setHasProfile(false)
      })
      .finally(() => {
        if (isMounted) setProfileChecking(false)
      })

    return () => { isMounted = false }
  }, [authLoading, token])

  const quickAccess = useMemo(() => {
    const artistItems = artists.slice(0, 4).map((artist, index) => ({
      title: artist.name,
      subtitle: 'Артист',
      coverUrl: artist.image_url || PLAYLIST_COVER_URLS[index % PLAYLIST_COVER_URLS.length],
      href: ROUTES.tracks,
    }))

    const genreItems = genres.slice(0, 4).map((genre, index) => ({
      title: genre.name,
      subtitle: 'Жанр',
      coverUrl: PLAYLIST_COVER_URLS[(index + 2) % PLAYLIST_COVER_URLS.length],
      href: ROUTES.tracks,
    }))

    return [...artistItems, ...genreItems].slice(0, 8)
  }, [artists, genres])

  const popularTracks = useMemo(
    () => [...tracks].sort((a, b) => b.popularity_score - a.popularity_score).slice(0, 5),
    [tracks],
  )

  const moodPlaylists = useMemo(() => playlists.slice(0, 6), [playlists])

  function handlePlayRec(rec: TrackRecommendation) {
    playTrack(rec.track)
  }

  function handleLikeRec(rec: TrackRecommendation) {
    toggleLike(rec.track)
  }

  function handleSkipRec(rec: TrackRecommendation) {
    skipTrack(rec.track)
  }

  function renderOnboardingCTA() {
    if (authLoading || profileChecking) return null

    if (!user && !token) {
      return (
        <section className="feed-section onboarding-cta">
          <p>Войди, чтобы настроить рекомендации под твой вкус.</p>
          <a className="auth-link" href={ROUTES.login}>Войти</a>
        </section>
      )
    }

    if (token && !hasProfile) {
      return (
        <section className="feed-section onboarding-cta">
          <h2>Настрой свои музыкальные предпочтения</h2>
          <p>Выбери жанры, артистов и треки — мы подберём музыку специально для тебя.</p>
          <a className="auth-submit" href={ROUTES.onboarding} style={{ display: 'inline-block', textDecoration: 'none' }}>Настроить рекомендации</a>
        </section>
      )
    }

    return null
  }

  function renderRecommendations() {
    if (recsLoading) {
      return (
        <section className="feed-section feed-rec-fallback" aria-labelledby="rec-title">
          <div className="section-heading">
            <h2 id="rec-title">Только для тебя</h2>
          </div>
          <p className="page-state">Загружаем рекомендации...</p>
        </section>
      )
    }

    if (recsError || recommendations.length === 0) {
      return null
    }

    return (
      <section className="feed-section feed-rec" aria-labelledby="rec-title">
        <div className="section-heading">
          <h2 id="rec-title">Только для тебя</h2>
        </div>
        <div className="rec-track-list">
          {recommendations.map((rec) => (
            <div className="rec-track-row" key={rec.track.id}>
              <button className="rec-track-play" type="button" aria-label={`Слушать ${rec.track.title}`} onClick={() => handlePlayRec(rec)}>
                ▶
              </button>
              <img src={rec.track.cover_url || PLAYLIST_COVER_URLS[0]} alt="" />
              <div className="rec-track-info">
                <strong>{rec.track.title}</strong>
                <span>{rec.track.artist.name}</span>
                <span className="rec-explanation">{rec.explanation}</span>
              </div>
              <div className="rec-track-score">
                <div className="rec-score-bar">
                  <div className="rec-score-fill" style={{ width: `${Math.min(rec.score, 100)}%` }} />
                </div>
                <span>{Math.round(rec.score)}</span>
              </div>
              <div className="rec-track-actions">
                <button className="rec-like-btn" type="button" aria-label="Нравится" title="Нравится" onClick={() => handleLikeRec(rec)}>
                  ♥
                </button>
                <button className="rec-skip-btn" type="button" aria-label="Не нравится" title="Не нравится" onClick={() => handleSkipRec(rec)}>
                  ✕
                </button>
              </div>
            </div>
          ))}
        </div>
      </section>
    )
  }

  function renderPlaylistRecommendations() {
    if (plRecsLoading) return null

    if (playlistRecs.length === 0) return null

    return (
      <section className="feed-section" aria-labelledby="pl-recs-title">
        <div className="section-heading">
          <h2 id="pl-recs-title">Плейлисты под твой вкус</h2>
        </div>
        <div className="playlist-grid">
          {playlistRecs.map((plRec, i) => (
            <div className="playlist-rec-card" key={plRec.playlist.id}>
              <PlaylistCard
                playlist={plRec.playlist}
                coverUrl={getPlaylistCoverUrl(i)}
              />
              <span className="rec-explanation">{plRec.explanation}</span>
            </div>
          ))}
        </div>
      </section>
    )
  }

  if (isLoading) {
    return (
      <AppShell>
        <div className="home-page">
          <section className="feed-section">
            <p className="page-state">Загрузка данных...</p>
          </section>
        </div>
      </AppShell>
    )
  }

  return (
    <AppShell>
      <div className="home-page">
        <h1 style={{ fontSize: 32, fontWeight: 700, margin: '0 0 24px' }}>
          Добро пожаловать
        </h1>

        {renderOnboardingCTA()}

        {error && (
          <section className="feed-section">
            <p className="page-state">{error}</p>
          </section>
        )}

        {renderRecommendations()}

        {renderPlaylistRecommendations()}

        {quickAccess.length > 0 && (
          <section className="feed-section" aria-labelledby="quick-access-title">
            <div className="section-heading">
              <h2 id="quick-access-title">Быстрый доступ</h2>
              <a href={ROUTES.tracks}>Показать всё</a>
            </div>
            <div className="quick-grid">
              {quickAccess.map((item) => (
                <a
                  className="quick-card"
                  href={item.href}
                  key={`${item.href}-${item.title}`}
                >
                  <img src={item.coverUrl} alt="" />
                  <span>{item.title}</span>
                </a>
              ))}
            </div>
          </section>
        )}

        {popularTracks.length > 0 && (
          <section className="feed-section feed-tracks" aria-labelledby="popular-tracks-title">
            <div className="section-heading">
              <h2 id="popular-tracks-title">Популярные треки</h2>
              <a href={ROUTES.tracks}>Все треки</a>
            </div>

            <div className="track-list">
              {popularTracks.map((track, index) => (
                <TrackRow track={track} index={index + 1} key={track.id} />
              ))}
            </div>
          </section>
        )}

        {moodPlaylists.length > 0 && (
          <section className="feed-section" aria-labelledby="mood-playlists-title">
            <div className="section-heading">
              <h2 id="mood-playlists-title">Плейлисты для настроения</h2>
              <a href={ROUTES.playlists}>Медиатека</a>
            </div>

            <div className="playlist-grid">
              {moodPlaylists.map((playlist, index) => (
                <PlaylistCard
                  playlist={playlist}
                  coverUrl={getPlaylistCoverUrl(index)}
                  key={playlist.id}
                />
              ))}
            </div>
          </section>
        )}
      </div>
    </AppShell>
  )
}
