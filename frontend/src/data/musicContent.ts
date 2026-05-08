import { ROUTES } from '../config/routes'

export type FeedTile = {
  title: string
  subtitle: string
  coverUrl: string
  href: string
}

export const PLAYLIST_COVER_URLS = [
  'https://images.unsplash.com/photo-1493225457124-a3eb161ffa5f?w=500',
  'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=500',
  'https://images.unsplash.com/photo-1470225620780-dba8ba36b745?w=500',
  'https://images.unsplash.com/photo-1459749411175-04bf5292ceea?w=500',
  'https://images.unsplash.com/photo-1508700115892-45ecd05ae2ad?w=500',
  'https://images.unsplash.com/photo-1524368535928-5b5e00ddc76b?w=500',
] as const

export const FALLBACK_TRACK_COVER_URL =
  'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=300'

export const FALLBACK_PLAYLIST_COVER_URL =
  'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=600'

export const DEFAULT_MINI_PLAYER = {
  title: 'Night Drive',
  artist: 'Ideal',
  coverUrl: 'https://images.unsplash.com/photo-1442975631134-6137411e9d4e?w=200',
} as const

export const FALLBACK_PERSONAL_TILES: FeedTile[] = [
  {
    title: 'Lo-Fi вечер',
    subtitle: 'Подборка по жанру',
    coverUrl: PLAYLIST_COVER_URLS[0],
    href: ROUTES.playlists,
  },
  {
    title: 'Электронная концентрация',
    subtitle: 'Для глубокой работы',
    coverUrl: PLAYLIST_COVER_URLS[2],
    href: ROUTES.playlists,
  },
  {
    title: 'Инди после полуночи',
    subtitle: 'Гитары и мягкий шум',
    coverUrl: PLAYLIST_COVER_URLS[5],
    href: ROUTES.playlists,
  },
  {
    title: 'Jazz Hop для прогулки',
    subtitle: 'Ритм без спешки',
    coverUrl: PLAYLIST_COVER_URLS[1],
    href: ROUTES.playlists,
  },
]

export const FALLBACK_QUICK_ACCESS: FeedTile[] = [
  {
    title: 'Любимые треки',
    subtitle: 'Недавно слушал',
    coverUrl: 'https://images.unsplash.com/photo-1516280440614-37939bbacd81?w=300',
    href: ROUTES.tracks,
  },
  {
    title: 'На подумать',
    subtitle: 'Продолжить',
    coverUrl: 'https://images.unsplash.com/photo-1516979187457-637abb4f9353?w=300',
    href: ROUTES.playlists,
  },
  {
    title: 'Ночной фокус',
    subtitle: 'Собрано для вечера',
    coverUrl: 'https://images.unsplash.com/photo-1504898770365-14faca6a7320?w=300',
    href: ROUTES.playlists,
  },
  {
    title: 'Pulse Nova',
    subtitle: 'Новый альбом',
    coverUrl: 'https://images.unsplash.com/photo-1492144534655-ae79c964c9d7?w=300',
    href: ROUTES.tracks,
  },
]

export function getPlaylistCoverUrl(index: number) {
  return PLAYLIST_COVER_URLS[index % PLAYLIST_COVER_URLS.length]
}
