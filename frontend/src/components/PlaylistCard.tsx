import { getPlaylistRoute } from '../config/routes'
import type { Playlist } from '../types'

type PlaylistCardProps = { playlist: Playlist; coverUrl: string }

export function PlaylistCard({ playlist, coverUrl }: PlaylistCardProps) {
  return (
    <a className="playlist-card" href={getPlaylistRoute(playlist.id)}>
      <article>
        <img src={coverUrl} alt="" />
        <h3>{playlist.name}</h3>
        <p>{playlist.description ?? 'Персональная подборка'}</p>
        <span>{playlist.track_count} треков</span>
      </article>
    </a>
  )
}
