import { getTrackRoute } from '../config/routes'
import { FALLBACK_TRACK_COVER_URL } from '../data/musicContent'
import type { Track } from '../types'

type TrackRowProps = {
  track: Track
  index?: number
}

export function TrackRow({ track, index }: TrackRowProps) {
  const genres = track.genres.map((genre) => genre.name).join(', ')

  return (
    <a
      className={`track-row ${index !== undefined ? 'track-row-numbered' : ''}`}
      href={getTrackRoute(track.id)}
    >
      {index !== undefined && <span className="track-index">{index}</span>}
      <img src={track.cover_url || FALLBACK_TRACK_COVER_URL} alt="" />
      <div className="track-row-copy">
        <strong>{track.title}</strong>
        <span>{track.artist.name}</span>
        {genres && <small>{genres}</small>}
      </div>
      <time>{formatDuration(track.duration_ms)}</time>
    </a>
  )
}

export function formatDuration(durationMs: number) {
  const totalSeconds = Math.round(durationMs / 1000)
  const minutes = Math.floor(totalSeconds / 60)
  const seconds = totalSeconds % 60

  return `${minutes}:${seconds.toString().padStart(2, '0')}`
}
