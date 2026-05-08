import { ROUTES } from '../config/routes'

export type NavigationItem = {
  href: string
  label: string
  icon: string
}

export const navigationItems: NavigationItem[] = [
  {
    href: ROUTES.home,
    label: 'Главная',
    icon: '⌂',
  },
  {
    href: ROUTES.tracks,
    label: 'Треки',
    icon: '♪',
  },
  {
    href: ROUTES.playlists,
    label: 'Плейлисты',
    icon: '▣',
  },
  {
    href: ROUTES.favorites,
    label: 'Избранное',
    icon: '♥',
  },
]

export function isNavigationItemActive(pathname: string, href: string) {
  if (href === ROUTES.home) {
    return pathname === ROUTES.home
  }

  return pathname === href || pathname.startsWith(`${href}/`)
}
