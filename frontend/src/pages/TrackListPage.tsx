import { useEffect, useMemo, useState } from 'react'
import { getApiErrorMessage, listTracks } from '../api'
import { AppShell } from '../components/AppShell'
import { TrackRow } from '../components/TrackRow'
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
    <AppShell>
      <div className="tracks-page">
        <header className="page-header">
          <div>
            <span>Музыка</span>
            <h1>Треки</h1>
          </div>
          <a href="/playlists">Плейлисты</a>
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
              {tracks.map((track, index) => (
                <TrackRow track={track} index={index + 1} key={track.id} />
              ))}
            </div>
          )}
        </section>
      </div>
    </AppShell>
  )
}
