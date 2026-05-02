import { useEffect, useState } from 'react'
import { getApiErrorMessage, listPlaylists } from '../api'
import type { Playlist } from '../types'

const playlistCovers = [
  'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=500',
  'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=500',
  'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=500',
  'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=500',
  'https://images.unsplash.com/photo-1524368535928-5b5e00ddc76b?w=500',
  'https://images.unsplash.com/photo-1516280440614-37939bbacd81?w=500',
]

export function PlaylistListPage() {
  const [playlists, setPlaylists] = useState<Playlist[]>([])
  const [isLoading, setIsLoading] = useState(true)
  const [error, setError] = useState('')

  useEffect(() => {
    let isMounted = true

    listPlaylists()
      .then((response) => {
        if (isMounted) {
          setPlaylists(response.playlists)
          setError('')
        }
      })
      .catch((requestError) => {
        if (isMounted) {
          setError(
            getApiErrorMessage(
              requestError,
              'Не удалось загрузить плейлисты. Проверь, что backend запущен.',
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

  return (
    <main className="playlists-page">
      <header className="tracks-header">
        <a className="back-link" href="/" aria-label="На главную">
          ←
        </a>
        <div>
          <span>Медиатека</span>
          <h1>Плейлисты</h1>
        </div>
      </header>

      <section className="playlist-library-hero">
        <div>
          <span>Подборки по жанрам, настроению и артистам</span>
          <h2>Выбирай плейлист и запускай поток треков.</h2>
        </div>
      </section>

      <section className="track-list-section">
        <div className="section-heading">
          <h2>Все плейлисты</h2>
          <span>
            {playlists.length > 0 ? `${playlists.length} подборок` : ''}
          </span>
        </div>

        {isLoading && <p className="page-state">Загружаем плейлисты...</p>}
        {error && <p className="page-state page-state-error">{error}</p>}

        {!isLoading && !error && (
          <div className="playlist-grid">
            {playlists.map((playlist, index) => (
              <PlaylistCard
                coverUrl={playlistCovers[index % playlistCovers.length]}
                playlist={playlist}
                key={playlist.id}
              />
            ))}
          </div>
        )}
      </section>

      <nav className="bottom-nav" aria-label="Основная навигация">
        <a href="/">
          <span>⌂</span>
          Главная
        </a>
        <a href="/tracks">
          <span>⌕</span>
          Поиск
        </a>
        <a className="active" href="/playlists">
          <span>▥</span>
          Моя медиатека
        </a>
      </nav>
    </main>
  )
}

function PlaylistCard({
  playlist,
  coverUrl,
}: {
  playlist: Playlist
  coverUrl: string
}) {
  return (
    <a className="playlist-card" href={`/playlists/${playlist.id}`}>
      <img src={coverUrl} alt="" />
      <div>
        <h3>{playlist.name}</h3>
        <p>{playlist.description ?? 'Персональная подборка'}</p>
        <span>{playlist.track_count} треков</span>
      </div>
    </a>
  )
}
