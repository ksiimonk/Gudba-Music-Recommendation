import { usePlayer } from '../context/PlayerContext'

export function MiniPlayer() {
  const { currentTrack, isPlaying, pauseTrack } = usePlayer()

  if (!currentTrack) {
    return null
  }

  const coverUrl = currentTrack.cover_url || 'https://images.unsplash.com/photo-1442975631134-6137411e9d4e?w=200'

  return (
    <section className="mini-player" aria-label="Сейчас играет">
      <img src={coverUrl} alt="" />
      <div className="mini-player-copy">
        <strong>{currentTrack.title}</strong>
        <span>{currentTrack.artist.name}</span>
      </div>
      <button type="button" aria-label={isPlaying ? 'Пауза' : 'Плей'} onClick={pauseTrack}>
        {isPlaying ? 'Ⅱ' : '▶'}
      </button>
    </section>
  )
}
