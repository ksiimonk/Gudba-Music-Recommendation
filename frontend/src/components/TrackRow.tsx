import { getTrackRoute } from '../config/routes'
import { FALLBACK_TRACK_COVER_URL } from '../data/musicContent'
import { usePlayer } from '../context/PlayerContext'
import type { Track } from '../types'

type TrackRowProps = {
  track: Track
  index?: number
  showLike?: boolean
}

export function TrackRow({ track, index, showLike }: TrackRowProps) {
  const genres = track.genres.map((genre) => genre.name).join(', ')
  const { playTrack, toggleLike, isFavorited } = usePlayer()
  const fav = isFavorited(track.id)

  function handlePlay(event: React.MouseEvent) {
    event.preventDefault()
    event.stopPropagation()
    playTrack(track)
  }

  function handleLike(event: React.MouseEvent) {
    event.preventDefault()
    event.stopPropagation()
    toggleLike(track)
  }

  function handleNavigate() {
    window.location.href = getTrackRoute(track.id)
  }

  return (
    <div
      className={`track-row ${index !== undefined ? 'track-row-numbered' : ''}`}
      onClick={handleNavigate}
      role="button"
      tabIndex={0}
      onKeyDown={(e) => { if (e.key === 'Enter') handleNavigate() }}
    >
      {index !== undefined && <span className="track-index">{index}</span>}
      <button className="track-play-button" type="button" aria-label={`Слушать ${track.title}`} onClick={handlePlay}>
        ▶
      </button>
      <img src={track.cover_url || FALLBACK_TRACK_COVER_URL} alt="" />
      <div className="track-row-copy">
        <strong>{track.title}</strong>
        <span>{track.artist.name}</span>
        {genres && <small>{genres}</small>}
      </div>
      {showLike && (
        <button
          className={`track-like-button ${fav ? 'liked' : ''}`}
          type="button"
          aria-label={fav ? 'Убрать из избранного' : 'Добавить в избранное'}
          onClick={handleLike}
        >
          {fav ? '♥' : '♡'}
        </button>
      )}
      <time>{formatDuration(track.duration_ms)}</time>
    </div>
  )
}

export function formatDuration(durationMs: number) {
  const totalSeconds = Math.round(durationMs / 1000)
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60

  return `${minutes}:${seconds.toString().padStart(2, '0')}`
}
