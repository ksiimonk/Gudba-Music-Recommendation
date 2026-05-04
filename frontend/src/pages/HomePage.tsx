import { useEffect, useMemo, useState } from 'react'
import { listPlaylists, listTracks } from '../api'
import { AppShell } from '../components/AppShell'
import { CoverTile } from '../components/CoverTile'
import { PlaylistCard } from '../components/PlaylistCard'
import { TrackRow } from '../components/TrackRow'
import type { Playlist, Track } from '../types'

type FeedTile = {
  title: string
  subtitle: string
  coverUrl: string
  href: string
}

const coverUrls = [
  'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=500',
  'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=500',
  'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=500',
  'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=500',
  'https://images.unsplash.com/photo-1508700115892-45ecd05ae2ad?w=500',
  'https://images.unsplash.com/photo-1524368535928-5b5e00ddc76b?w=500',
]

const fallbackPersonalTiles: FeedTile[] = [
  {
    title: 'Lo-Fi вечер',
    subtitle: 'Подборка по жанру',
    coverUrl: coverUrls[0],
    href: '/playlists',
  },
  {
    title: 'Электронная концентрация',
    subtitle: 'Для глубокой работы',
    coverUrl: coverUrls[2],
    href: '/playlists',
  },
  {
    title: 'Инди после полуночи',
    subtitle: 'Гитары и мягкий шум',
    coverUrl: coverUrls[5],
    href: '/playlists',
  },
  {
    title: 'Jazz Hop для прогулки',
    subtitle: 'Ритм без спешки',
    coverUrl: coverUrls[1],
    href: '/playlists',
  },
]

const fallbackQuickAccess: FeedTile[] = [
  {
    title: 'Любимые треки',
    subtitle: 'Недавно слушал',
    coverUrl: 'https://images.unsplash.com/photo-1516280440614-37939bbacd81?w=300',
    href: '/tracks',
  },
  {
    title: 'На подумать',
    subtitle: 'Продолжить',
    coverUrl: 'https://images.unsplash.com/photo-1516979187457-637abb4f9353?w=300',
    href: '/playlists',
  },
  {
    title: 'Ночной фокус',
    subtitle: 'Собрано для вечера',
    coverUrl: 'https://images.unsplash.com/photo-1504898770365-14faca6a7320?w=300',
    href: '/playlists',
  },
  {
    title: 'Pulse Nova',
    subtitle: 'Новый альбом',
    coverUrl: 'https://images.unsplash.com/photo-1492144534655-ae79c964c9d7?w=300',
    href: '/tracks',
  },
]

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
      return fallbackPersonalTiles
    }

    return playlists.slice(0, 6).map((playlist, index) => ({
      title: playlist.name,
      subtitle: playlist.description ?? `${playlist.track_count} треков`,
      coverUrl: getPlaylistCover(index),
      href: `/playlists/${playlist.id}`,
    }))
  }, [playlists])

  const quickAccess = useMemo(() => {
    const playlistItems = playlists.slice(0, 4).map((playlist, index) => ({
      title: playlist.name,
      subtitle: 'Плейлист',
      coverUrl: getPlaylistCover(index),
      href: `/playlists/${playlist.id}`,
    }))

    const trackItems = tracks.slice(0, 4).map((track) => ({
      title: track.title,
      subtitle: track.artist.name,
      coverUrl: track.cover_url || coverUrls[0],
      href: `/tracks/${track.id}`,
    }))

    const items = [...playlistItems, ...trackItems].slice(0, 8)

    return items.length > 0 ? items : fallbackQuickAccess
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
            <a className="active" href="/">
              Все
            </a>
            <a href="/tracks">Треки</a>
            <a href="/playlists">Плейлисты</a>
          </nav>
        </header>

        <section className="feed-section" aria-labelledby="quick-access-title">
          <div className="section-heading">
            <h2 id="quick-access-title">Быстрый доступ</h2>
          </div>
          <div className="quick-grid">
            {quickAccess.map((item) => (
              <a className="quick-card" href={item.href} key={`${item.href}-${item.title}`}>
                <img src={item.coverUrl} alt="" />
                <span>{item.title}</span>
              </a>
            ))}
          </div>
        </section>

        <section className="feed-section" aria-labelledby="personal-title">
          <div className="section-heading">
            <h2 id="personal-title">Только для тебя</h2>
            <a href="/playlists">Открыть все</a>
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
            <a href="/tracks">Все треки</a>
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
            <a href="/playlists">Медиатека</a>
          </div>

          {moodPlaylists.length === 0 ? (
            <div className="card-row">
              {fallbackPersonalTiles.map((item) => (
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
                  coverUrl={getPlaylistCover(index)}
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

function getPlaylistCover(index: number) {
  return coverUrls[index % coverUrls.length]
}
