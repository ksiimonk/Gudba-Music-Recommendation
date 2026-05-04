import { useEffect, useMemo, useState } from 'react'
import { listPlaylists, listTracks } from '../api'
import { AppShell } from '../components/AppShell'
import { CoverTile } from '../components/CoverTile'
import { PlaylistCard } from '../components/PlaylistCard'
import { TrackRow } from '../components/TrackRow'
import { ROUTES, getPlaylistRoute, getTrackRoute } from '../config/routes'
import {
  FALLBACK_PERSONAL_TILES,
  FALLBACK_QUICK_ACCESS,
  PLAYLIST_COVER_URLS,
  getPlaylistCoverUrl,
} from '../data/musicContent'
import type { Playlist, Track } from '../types'

export function HomePage() {
  const [playlists, setPlaylists] = useState<Playlist[]>([])
  const [tracks, setTracks] = useState<Track[]>([])

  useEffect(() => {
    let isMounted = true

    Promise.all([listPlaylists(), listTracks()])
      .then(([playlistResponse, trackResponse]) => {
        if (isMounted) {
          setPlaylists(playlistResponse.playlists)
          setTracks(trackResponse.tracks)
        }
      })
      .catch(() => {
        if (isMounted) {
          setPlaylists([])
          setTracks([])
        }
      })

    return () => {
      isMounted = false
    }
  }, [])

  const personalTiles = useMemo(() => {
    if (playlists.length === 0) {
      return FALLBACK_PERSONAL_TILES
    }

    return playlists.slice(0, 6).map((playlist, index) => ({
      title: playlist.name,
      subtitle: playlist.description ?? `${playlist.track_count} треков`,
      coverUrl: getPlaylistCoverUrl(index),
      href: getPlaylistRoute(playlist.id),
    }))
  }, [playlists])

  const quickAccess = useMemo(() => {
    const playlistItems = playlists.slice(0, 4).map((playlist, index) => ({
      title: playlist.name,
      subtitle: 'Плейлист',
      coverUrl: getPlaylistCoverUrl(index),
      href: getPlaylistRoute(playlist.id),
    }))

    const trackItems = tracks.slice(0, 4).map((track) => ({
      title: track.title,
      subtitle: track.artist.name,
      coverUrl: track.cover_url || PLAYLIST_COVER_URLS[0],
      href: getTrackRoute(track.id),
    }))

    const items = [...playlistItems, ...trackItems].slice(0, 8)

    return items.length > 0 ? items : FALLBACK_QUICK_ACCESS
  }, [playlists, tracks])

  const popularTracks = useMemo(() => tracks.slice(0, 5), [tracks])
  const moodPlaylists = useMemo(() => playlists.slice(0, 6), [playlists])

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

        <section className="feed-section" aria-labelledby="personal-title">
          <div className="section-heading">
            <h2 id="personal-title">Только для тебя</h2>
            <a href={ROUTES.playlists}>Открыть все</a>
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

        <section className="feed-section" aria-labelledby="mood-playlists-title">
          <div className="section-heading">
            <h2 id="mood-playlists-title">Плейлисты для настроения</h2>
            <a href={ROUTES.playlists}>Медиатека</a>
          </div>

          {moodPlaylists.length === 0 ? (
            <div className="card-row">
              {FALLBACK_PERSONAL_TILES.map((item) => (
                <CoverTile
                  title={item.title}
                  subtitle={item.subtitle}
                  coverUrl={item.coverUrl}
                  href={item.href}
                  key={item.title}
                />
              ))}
            </div>
          ) : (
            <div className="playlist-grid">
              {moodPlaylists.map((playlist, index) => (
                <PlaylistCard
                  playlist={playlist}
                  coverUrl={getPlaylistCoverUrl(index)}
                  key={playlist.id}
                />
              ))}
            </div>
          )}
        </section>
      </div>
    </AppShell>
  )
}
