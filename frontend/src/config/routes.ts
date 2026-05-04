export const ROUTES = {
  home: '/',
  login: '/login',
  register: '/register',
  onboarding: '/onboarding',
  tracks: '/tracks',
  playlists: '/playlists',
} as const

export function getTrackRoute(id: number | string) {
  return `${ROUTES.tracks}/${id}`
}

export function getPlaylistRoute(id: number | string) {
  return `${ROUTES.playlists}/${id}`
}

export function isNestedRoute(pathname: string, route: string) {
  return pathname.startsWith(`${route}/`)
}

export function readRouteId(pathname: string, route: string) {
  return Number(pathname.slice(`${route}/`.length))
}
