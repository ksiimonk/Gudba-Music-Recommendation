type CoverTileProps = {
  title: string
  subtitle: string
  coverUrl: string
  href?: string
}

export function CoverTile({ title, subtitle, coverUrl, href }: CoverTileProps) {
  const content = (
    <>
      <img src={coverUrl} alt="" />
      <h3>{title}</h3>
      <p>{subtitle}</p>
    </>
  )

  if (href) {
    return (
      <a className="cover-tile" href={href}>
        {content}
      </a>
    )
  }

  return <article className="cover-tile">{content}</article>
}
