import { useEffect, useMemo, useState } from 'react'
import { listArtists, listGenres, listPlaylists, listTracks } from '../api'
import { AppShell } from '../components/AppShell'
import { CoverTile } from '../components/CoverTile'
import { PlaylistCard } from '../components/PlaylistCard'
import { TrackRow } from '../components/TrackRow'
import { ROUTES } from '../config/routes'
import { PLAYLIST_COVER_URLS, getPlaylistCoverUrl } from '../data/musicContent'
import type { Artist, Genre, Playlist, Track } from '../types'

export function HomePage() {
  const [genres, setGenres] = useState<Genre[]>([])
  const [artists, setArtists] = useState<Artist[]>([])
  const [playlists, setPlaylists] = useState<Playlist[]>([])
  const [tracks, setTracks] = useState<Track[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

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
        if (isMounted) {
          setError('Не удалось загрузить данные. Проверь, запущен ли backend.')
        }
      })
      .finally(() => {
        if (isMounted) {
          setIsLoading(false)
        }
      })

    return () => {
      isMounted = false
    }
  }, [])

  const personalTiles = useMemo(() => {
    return genres.slice(0, 6).map((genre, index) => ({
      title: genre.name,
      subtitle: 'Подборка',
      coverUrl: PLAYLIST_COVER_URLS[index % PLAYLIST_COVER_URLS.length],
      href: ROUTES.tracks,
    }))
  }, [genres])

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

  if (isLoading) {
    return (
      <AppShell>
        <div className="home-page">
          <header className="feed-header">
            <div>
              <span>Сегодня в Gudba</span>
              <h1>Загружаем...</h1>
            </div>
          </header>
          <section className="feed-section">
            <div className="section-heading"><h2>Популярные треки</h2></div>
            <p className="page-state">Загрузка данных...</p>
          </section>
        </div>
      </AppShell>
    )
  }

  return (
    <AppShell>
      <div className="home-page">
        <header className="feed-header">
          <div>
            <span>Сегодня в Gudba</span>
            <h1>Музыка под твой день</h1>
          </div>
          <nav className="home-filters" aria-label="Фильтры главной">
            <a className="active" href={ROUTES.home}>
              Все
            </a>
            <a href={ROUTES.tracks}>Треки</a>
            <a href={ROUTES.playlists}>Плейлисты</a>
          </nav>
        </header>

        {error && (
          <section className="feed-section">
            <p className="page-state">{error}</p>
          </section>
        )}

        {quickAccess.length > 0 && (
          <section className="feed-section" aria-labelledby="quick-access-title">
            <div className="section-heading">
              <h2 id="quick-access-title">Быстрый доступ</h2>
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

        {personalTiles.length > 0 && (
          <section className="feed-section" aria-labelledby="personal-title">
            <div className="section-heading">
              <h2 id="personal-title">Только для тебя</h2>
            </div>
            <div className="card-row">
              {personalTiles.map((item) => (
                <CoverTile
                  title={item.title}
                  subtitle={item.subtitle}
                  coverUrl={item.coverUrl}
                  href={item.href}
                  key={`${item.href}-${item.title}`}
                />
              ))}
            </div>
          </section>
        )}

        <section className="feed-section feed-tracks" aria-labelledby="popular-tracks-title">
          <div className="section-heading">
            <h2 id="popular-tracks-title">Популярные треки</h2>
            <a href={ROUTES.tracks}>Все треки</a>
          </div>

          {popularTracks.length === 0 ? (
            <p className="page-state">Треки появятся после запуска backend.</p>
          ) : (
            <div className="track-list">
              {popularTracks.map((track, index) => (
                <TrackRow track={track} index={index + 1} key={track.id} />
              ))}
            </div>
          )}
        </section>

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
