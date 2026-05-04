import { DEFAULT_MINI_PLAYER } from '../data/musicContent'

type MiniPlayerProps = {
  title?: string
  artist?: string
  coverUrl?: string
}

export function MiniPlayer({
  title = DEFAULT_MINI_PLAYER.title,
  artist = DEFAULT_MINI_PLAYER.artist,
  coverUrl = DEFAULT_MINI_PLAYER.coverUrl,
}: MiniPlayerProps) {
  return (
    <section className="mini-player" aria-label="Сейчас играет">
      <img src={coverUrl} alt="" />
      <div className="mini-player-copy">
        <strong>{title}</strong>
        <span>{artist}</span>
      </div>
      <button type="button" aria-label="Пауза">
        Ⅱ
      </button>
    </section>
  )
}
