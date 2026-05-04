export type NavigationItem = {
  href: string
  label: string
  icon: string
}

export const navigationItems: NavigationItem[] = [
  {
    href: '/',
    label: 'Главная',
    icon: '⌂',
  },
  {
    href: '/tracks',
    label: 'Треки',
    icon: '♪',
  },
  {
    href: '/playlists',
    label: 'Плейлисты',
    icon: '▣',
  },
]

export function isNavigationItemActive(pathname: string, href: string) {
  if (href === '/') {
    return pathname === '/'
  }

  return pathname === href || pathname.startsWith(`${href}/`)
}
