import { useEffect, useState } from 'react'
import { getFavorites } from '../api'
import { AppShell } from '../components/AppShell'
import { TrackRow } from '../components/TrackRow'
import { useAuth } from '../auth/AuthContext'
import { ROUTES } from '../config/routes'
import type { Playlist } from '../types'

export function FavoritesPage() {
  const { token, isLoading: authLoading } = useAuth()
  const [playlist, setPlaylist] = useState<Playlist | null>(null)
  const [pageLoading, setPageLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let isMounted = true

    if (authLoading) return

    if (!token) {
      window.location.href = ROUTES.login
      return
    }

    getFavorites(token)
      .then((res) => {
        if (isMounted) setPlaylist(res.playlist)
      })
      .catch(() => {
        if (isMounted) setError('Не удалось загрузить избранное.')
      })
      .finally(() => {
        if (isMounted) setPageLoading(false)
      })

    return () => { isMounted = false }
  }, [token, authLoading])

  return (
    <AppShell>
      <div className="playlist-detail-page">
        <header className="page-header">
          <div>
            <span>Плейлист</span>
            <h1>{playlist?.name ?? 'Мои любимые треки'}</h1>
          </div>
          <a href={ROUTES.home}>На главную</a>
        </header>

        {(authLoading || pageLoading) && <p className="page-state">Загружаем избранное...</p>}
        {error && <p className="page-state page-state-error">{error}</p>}

        {playlist && !error && (
          <section className="playlist-tracks">
            {playlist.tracks && playlist.tracks.length > 0 ? (
              playlist.tracks.map((track, index) => (
                <TrackRow key={track.id} track={track} index={index + 1} showLike />
              ))
            ) : (
              !pageLoading && (
                <p className="page-state">
                  Здесь пока нет треков. Отмечай ♥ понравившиеся треки, и они появятся здесь.
                </p>
              )
            )}
          </section>
        )}
      </div>
    </AppShell>
  )
}
