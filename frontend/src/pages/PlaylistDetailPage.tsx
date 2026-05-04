import { useEffect, useState } from 'react'
import { getApiErrorMessage, getPlaylistById } from '../api'
import { AppShell } from '../components/AppShell'
import { TrackRow } from '../components/TrackRow'
import type { Playlist, Track } from '../types'

type PlaylistDetailPageProps = {
  playlistId: number
}

const fallbackCover =
  'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=600'

export function PlaylistDetailPage({ playlistId }: PlaylistDetailPageProps) {
  const isInvalidPlaylistId = !Number.isFinite(playlistId) || playlistId <= 0
  const [playlist, setPlaylist] = useState<Playlist | null>(null)
  const [isLoading, setIsLoading] = useState(!isInvalidPlaylistId)
  const [error, setError] = useState('')

  useEffect(() => {
    let isMounted = true

    if (isInvalidPlaylistId) {
      return () => {
        isMounted = false
      }
    }

    getPlaylistById(playlistId)
      .then((response) => {
        if (isMounted) {
          setPlaylist(response.playlist)
          setError('')
        }
      })
      .catch((requestError) => {
        if (isMounted) {
          setError(
            getApiErrorMessage(
              requestError,
              'Не удалось загрузить плейлист. Проверь, что backend запущен.',
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
  }, [isInvalidPlaylistId, playlistId])

  const visibleError = isInvalidPlaylistId
    ? 'Некорректный id плейлиста.'
    : error
  const tracks = playlist?.tracks ?? []
  const coverUrl = getPlaylistCover(tracks)

  return (
    <AppShell>
      <div className="playlist-detail-page">
        <header className="page-header">
          <div>
            <span>Плейлист</span>
            <h1>{playlist?.name ?? 'Детали плейлиста'}</h1>
          </div>
          <a href="/playlists">Назад к плейлистам</a>
        </header>

        {isLoading && <p className="page-state">Загружаем плейлист...</p>}
        {visibleError && (
          <p className="page-state page-state-error">{visibleError}</p>
        )}

        {playlist && !visibleError && (
          <>
            <section className="playlist-detail-hero">
              <div className="playlist-cover-stack">
                <img src={coverUrl} alt="" />
              </div>
              <div>
                <span>Публичная подборка</span>
                <h2>{playlist.name}</h2>
                <p>
                  {playlist.description ??
                    'Собрано для музыкального настроения.'}
                </p>
                <strong>{tracks.length} треков</strong>
              </div>
            </section>

            <section className="track-actions-panel">
              <button type="button">▶ Воспроизвести</button>
              <a href="/tracks">Открыть все треки</a>
            </section>

            <section className="playlist-tracks-section">
              <div className="section-heading">
                <h2>Треки в плейлисте</h2>
                <span>{tracks.length > 0 ? `${tracks.length} треков` : ''}</span>
              </div>

              {tracks.length === 0 ? (
                <p className="page-state">
                  В этом плейлисте пока нет треков.
                </p>
              ) : (
                <div className="track-list">
                  {tracks.map((track, index) => (
                    <TrackRow track={track} index={index + 1} key={track.id} />
                  ))}
                </div>
              )}
            </section>
          </>
        )}
      </div>
    </AppShell>
  )
}

function getPlaylistCover(tracks: Track[]) {
  return tracks.find((track) => track.cover_url)?.cover_url ?? fallbackCover
}
