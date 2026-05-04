type MiniPlayerProps = {
  title?: string
  artist?: string
  coverUrl?: string
}

const defaultCoverUrl =
  'https://images.unsplash.com/photo-1442975631134-6137411e9d4e?w=200'

export function MiniPlayer({
  title = 'Night Drive',
  artist = 'Ideal',
  coverUrl = defaultCoverUrl,
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
