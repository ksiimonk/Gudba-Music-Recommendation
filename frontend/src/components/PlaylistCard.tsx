import type { Playlist } from '../types'

type PlaylistCardProps = {
  playlist: Playlist
  coverUrl: string
}

export function PlaylistCard({ playlist, coverUrl }: PlaylistCardProps) {
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
