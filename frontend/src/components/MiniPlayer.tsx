import { usePlayer } from '../context/PlayerContext'
import type { Track } from '../types'

function fmt(s: number): string {
  if (!s || !isFinite(s)) return '0:00'
  const m = Math.floor(s / 60)
  const sec = Math.floor(s % 60)
  return `${m}:${sec.toString().padStart(2, '0')}`
}

export function MiniPlayer() {
  const {
    currentTrack, isPlaying, isLoading, progress, duration,
    pauseTrack, playTrack, toggleLike, isFavorited,
  } = usePlayer()

  if (!currentTrack) return null

  const track: Track = currentTrack
  const coverUrl = track.cover_url || 'https://images.unsplash.com/photo-1442975631134-6137411e9d4e?w=200'
  const fav = isFavorited(track.id)
  const pct = duration > 0 ? (progress / duration) * 100 : 0

  function handlePlayPause() {
    if (isLoading) return
    if (isPlaying) pauseTrack()
    else playTrack(track)
  }

  function handleLike() { toggleLike(track) }

  return (
    <section className="mini-player" aria-label="Сейчас играет">
      <div className="player-track-info">
        <img src={coverUrl} alt="" />
        <div className="player-track-copy">
          <strong>{track.title}</strong>
          <span>{track.artist.name}</span>
        </div>
      </div>

      <div className="player-controls">
        <div className="player-buttons">
          <button className="player-btn" type="button" aria-label="Предыдущий" disabled>⏮</button>
          <button
            className="player-btn play-btn"
            type="button"
            aria-label={isLoading ? 'Загрузка' : isPlaying ? 'Пауза' : 'Воспроизвести'}
            onClick={handlePlayPause}
          >
            {isLoading ? '⏳' : isPlaying ? '⏸' : '▶'}
          </button>
          <button className="player-btn" type="button" aria-label="Следующий" disabled>⏭</button>
        </div>
        <div className="player-progress">
          <span>{fmt(progress)}</span>
          <div className="progress-bar">
            <div className="progress-fill" style={{ width: `${pct}%` }} />
          </div>
          <span>{fmt(duration)}</span>
        </div>
      </div>

      <div className="player-actions">
        <button
          className={`player-like-btn ${fav ? 'liked' : ''}`}
          type="button"
          aria-label={fav ? 'Убрать из избранного' : 'Добавить в избранное'}
          onClick={handleLike}
        >
          {fav ? '♥' : '♡'}
        </button>
        <div className="player-volume">
          <span style={{ fontSize: 16 }}>🔊</span>
          <input type="range" min="0" max="100" defaultValue="60" aria-label="Громкость" />
        </div>
      </div>
    </section>
  )
}
