import { useEffect, useMemo, useState } from 'react'
import { getApiErrorMessage, listTracks } from '../api'
import type { Track } from '../types'

export function TrackListPage() {
  const [tracks, setTracks] = useState<Track[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let isMounted = true

    listTracks()
      .then((response) => {
        if (isMounted) {
          setTracks(response.tracks)
          setError('')
        }
      })
      .catch((requestError) => {
        if (isMounted) {
          setError(
            getApiErrorMessage(
              requestError,
              'Не удалось загрузить треки. Проверь, что backend запущен.',
            ),
          )
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

  const popularTracks = useMemo(() => tracks.slice(0, 5), [tracks])

  return (
    <main className="tracks-page">
      <header className="tracks-header">
        <a className="back-link" href="/" aria-label="На главную">
          ←
        </a>
        <div>
          <span>Музыка</span>
          <h1>Треки</h1>
        </div>
      </header>

      {popularTracks.length > 0 && (
        <section className="track-hero" aria-label="Популярные треки">
          <div>
            <span>Популярное сейчас</span>
            <h2>{popularTracks[0].title}</h2>
            <p>{popularTracks[0].artist.name}</p>
          </div>
          <a href={`/tracks/${popularTracks[0].id}`}>Открыть</a>
        </section>
      )}

      <section className="track-list-section">
        <div className="section-heading">
          <h2>Все треки</h2>
          <span>{tracks.length > 0 ? `${tracks.length} треков` : ''}</span>
        </div>

        {isLoading && <p className="page-state">Загружаем треки...</p>}
        {error && <p className="page-state page-state-error">{error}</p>}

        {!isLoading && !error && (
          <div className="track-list">
            {tracks.map((track) => (
              <TrackRow track={track} key={track.id} />
            ))}
          </div>
        )}
      </section>

      <nav className="bottom-nav" aria-label="Основная навигация">
        <a href="/">
          <span>⌂</span>
          Главная
        </a>
        <a className="active" href="/tracks">
          <span>⌕</span>
          Поиск
        </a>
        <a href="/playlists">
          <span>▥</span>
          Моя медиатека
        </a>
      </nav>
    </main>
  )
}

function TrackRow({ track }: { track: Track }) {
  const genres = track.genres.map((genre) => genre.name).join(', ')

  return (
    <a className="track-row" href={`/tracks/${track.id}`}>
      <img src={track.cover_url} alt="" />
      <div>
        <strong>{track.title}</strong>
        <span>{track.artist.name}</span>
        {genres && <small>{genres}</small>}
      </div>
      <time>{formatDuration(track.duration_ms)}</time>
    </a>
  )
}

function formatDuration(durationMs: number) {
  const totalSeconds = Math.round(durationMs / 1000)
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60

  return `${minutes}:${seconds.toString().padStart(2, '0')}`
}
