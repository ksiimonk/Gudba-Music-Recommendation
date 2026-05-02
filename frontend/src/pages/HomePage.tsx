import { useEffect, useMemo, useState } from 'react'
import type { ReactNode } from 'react'
import { listPlaylists } from '../api'
import type { Playlist } from '../types'

type CoverCard = {
  title: string
  subtitle: string
  coverUrl: string
}

const coverUrls = [
  'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=500',
  'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=500',
  'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=500',
  'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=500',
  'https://images.unsplash.com/photo-1508700115892-45ecd05ae2ad?w=500',
  'https://images.unsplash.com/photo-1524368535928-5b5e00ddc76b?w=500',
]

const fallbackPlaylists: CoverCard[] = [
  {
    title: 'Lo-Fi вечер',
    subtitle: 'Подборка по жанру',
    coverUrl: coverUrls[0],
  },
  {
    title: 'Электронная концентрация',
    subtitle: 'Для глубокой работы',
    coverUrl: coverUrls[2],
  },
  {
    title: 'Инди после полуночи',
    subtitle: 'Гитары и мягкий шум',
    coverUrl: coverUrls[5],
  },
  {
    title: 'Jazz Hop для прогулки',
    subtitle: 'Ритм без спешки',
    coverUrl: coverUrls[1],
  },
]

const recentItems: CoverCard[] = [
  {
    title: 'Любимые треки',
    subtitle: 'Последнее прослушивание',
    coverUrl: 'https://images.unsplash.com/photo-1516280440614-37939bbacd81?w=300',
  },
  {
    title: 'На подумать',
    subtitle: 'Недавно слушал',
    coverUrl: 'https://images.unsplash.com/photo-1516979187457-637abb4f9353?w=300',
  },
  {
    title: 'Ночной фокус',
    subtitle: 'Продолжить',
    coverUrl: 'https://images.unsplash.com/photo-1504898770365-14faca6a7320?w=300',
  },
]

const recommendedAlbums: CoverCard[] = [
  {
    title: 'Golden Thread',
    subtitle: 'Новый альбом исполнителя',
    coverUrl: 'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=500',
  },
  {
    title: 'Pulse Nova',
    subtitle: 'Синтвейв и глубокий бас',
    coverUrl: 'https://images.unsplash.com/photo-1492144534655-ae79c964c9d7?w=500',
  },
  {
    title: 'Blue Note Trio',
    subtitle: 'Современный джаз',
    coverUrl: 'https://images.unsplash.com/photo-1415201364774-f6f0bb35f28f?w=500',
  },
  {
    title: 'Soft Glow',
    subtitle: 'Инди-поп релиз',
    coverUrl: 'https://images.unsplash.com/photo-1516280440614-37939bbacd81?w=500',
  },
]

export function HomePage() {
  const [playlists, setPlaylists] = useState<Playlist[]>([])

  useEffect(() => {
    let isMounted = true

    listPlaylists()
      .then((response) => {
        if (isMounted) {
          setPlaylists(response.playlists)
        }
      })
      .catch(() => {
        if (isMounted) {
          setPlaylists([])
        }
      })

    return () => {
      isMounted = false
    }
  }, [])

  const personalPlaylists = useMemo(() => {
    if (playlists.length === 0) {
      return fallbackPlaylists
    }

    return playlists.slice(0, 6).map((playlist, index) => ({
      title: playlist.name,
      subtitle: playlist.description ?? `${playlist.track_count} треков`,
      coverUrl: coverUrls[index % coverUrls.length],
    }))
  }, [playlists])

  const recentPlaylists = useMemo(
    () => [...personalPlaylists.slice(0, 5), ...recentItems].slice(0, 8),
    [personalPlaylists],
  )

  return (
    <main className="home-page">
      <header className="home-topbar">
        <a className="home-avatar" href="/login" aria-label="Профиль">
          S
        </a>
        <nav className="home-filters" aria-label="Фильтры главной">
          <a className="active" href="/">
            Все
          </a>
          <a href="/tracks">Музыка</a>
          <a href="/playlists">Плейлисты</a>
        </nav>
      </header>

      <section className="quick-grid" aria-label="Недавно прослушанные">
        {recentPlaylists.map((item) => (
          <a className="quick-card" href="/playlists" key={item.title}>
            <img src={item.coverUrl} alt="" />
            <span>{item.title}</span>
          </a>
        ))}
      </section>

      <section className="featured-release">
        <div className="artist-intro">
          <img src={recommendedAlbums[0].coverUrl} alt="" />
          <div>
            <span>Новый релиз исполнителя</span>
            <h1>{recommendedAlbums[0].title}</h1>
          </div>
        </div>

        <article className="release-card">
          <img src="https://images.unsplash.com/photo-1507838153414-b4b713384a76?w=500" alt="" />
          <div>
            <span>Альбом</span>
            <h2>Midnight Signals</h2>
            <p>Golden Thread, Soft Glow, Pulse Nova</p>
          </div>
          <button type="button" aria-label="Воспроизвести Midnight Signals">
            ▶
          </button>
        </article>
      </section>

      <HomeSection title="Только для тебя">
        {personalPlaylists.map((playlist) => (
          <CoverTile item={playlist} key={playlist.title} />
        ))}
      </HomeSection>

      <HomeSection title="Альбомы рекомендованных исполнителей">
        {recommendedAlbums.map((album) => (
          <CoverTile item={album} key={album.title} />
        ))}
      </HomeSection>

      <div className="mini-player" aria-label="Сейчас играет">
        <img src="https://images.unsplash.com/photo-1442975631134-6137411e9d4e?w=200" alt="" />
        <div>
          <strong>Night Drive</strong>
          <span>Ideal</span>
        </div>
        <button type="button" aria-label="Пауза">
          ❚❚
        </button>
      </div>

      <nav className="bottom-nav" aria-label="Основная навигация">
        <a className="active" href="/">
          <span>⌂</span>
          Главная
        </a>
        <a href="/tracks">
          <span>⌕</span>
          Поиск
        </a>
        <a href="/playlists">
          <span>▥</span>
          Моя медиатека
        </a>
      </nav>
    </main>
  )
}

function HomeSection({
  title,
  children,
}: {
  title: string
  children: ReactNode
}) {
  return (
    <section className="home-section">
      <h2>{title}</h2>
      <div className="card-row">{children}</div>
    </section>
  )
}

function CoverTile({ item }: { item: CoverCard }) {
  return (
    <article className="cover-tile">
      <img src={item.coverUrl} alt="" />
      <h3>{item.title}</h3>
      <p>{item.subtitle}</p>
    </article>
  )
}
