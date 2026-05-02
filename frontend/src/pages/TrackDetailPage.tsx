import { useEffect, useState } from 'react'
import { getApiErrorMessage, getTrackById } from '../api'
import type { Track } from '../types'

type TrackDetailPageProps = {
  trackId: number
}

export function TrackDetailPage({ trackId }: TrackDetailPageProps) {
  const isInvalidTrackId = !Number.isFinite(trackId) || trackId <= 0
  const [track, setTrack] = useState<Track | null>(null)
  const [isLoading, setIsLoading] = useState(!isInvalidTrackId)
  const [error, setError] = useState('')

  useEffect(() => {
    let isMounted = true

    if (isInvalidTrackId) {
      return () => {
        isMounted = false
      }
    }

    getTrackById(trackId)
      .then((response) => {
        if (isMounted) {
          setTrack(response.track)
          setError('')
        }
      })
      .catch((requestError) => {
        if (isMounted) {
          setError(
            getApiErrorMessage(
              requestError,
              'Не удалось загрузить трек. Проверь, что backend запущен.',
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
  }, [isInvalidTrackId, trackId])

  const visibleError = isInvalidTrackId ? 'Некорректный id трека.' : error

  return (
    <main className="track-detail-page">
      <header className="tracks-header">
        <a className="back-link" href="/tracks" aria-label="Назад к трекам">
          ←
        </a>
        <div>
          <span>Трек</span>
          <h1>{track?.title ?? 'Детали трека'}</h1>
        </div>
      </header>

      {isLoading && <p className="page-state">Загружаем трек...</p>}
      {visibleError && (
        <p className="page-state page-state-error">{visibleError}</p>
      )}

      {track && !visibleError && (
        <>
          <section className="track-detail-hero">
            <img src={track.cover_url} alt="" />
            <div>
              <span>{track.artist.name}</span>
              <h2>{track.title}</h2>
              <p>{track.genres.map((genre) => genre.name).join(' • ')}</p>
            </div>
          </section>

          <section className="track-actions-panel">
            <button type="button">▶ Воспроизвести</button>
            {track.spotify_url && (
              <a href={track.spotify_url} target="_blank" rel="noreferrer">
                Открыть в Spotify
              </a>
            )}
          </section>

          <section className="track-meta-grid" aria-label="Информация о треке">
            <MetaItem label="Длительность" value={formatDuration(track.duration_ms)} />
            <MetaItem label="Популярность" value={`${track.popularity_score}/100`} />
            <MetaItem label="Исполнитель" value={track.artist.name} />
          </section>

          {track.artist.bio && (
            <section className="artist-note">
              <h2>Об исполнителе</h2>
              <p>{track.artist.bio}</p>
            </section>
          )}
        </>
      )}
    </main>
  )
}

function MetaItem({ label, value }: { label: string; value: string }) {
  return (
    <article className="meta-item">
      <span>{label}</span>
      <strong>{value}</strong>
    </article>
  )
}

function formatDuration(durationMs: number) {
  const totalSeconds = Math.round(durationMs / 1000)
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60

  return `${minutes}:${seconds.toString().padStart(2, '0')}`
}
