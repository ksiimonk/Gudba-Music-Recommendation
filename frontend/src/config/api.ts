const DEFAULT_API_BASE_URL = ''

export const API_BASE_URL = normalizeBaseUrl(
  import.meta.env.VITE_API_BASE_URL ?? DEFAULT_API_BASE_URL,
)

export const API_ENDPOINTS = {
  auth: {
    register: '/api/v1/auth/register',
    login: '/api/v1/auth/login',
  },
  me: '/api/v1/me',
  tracks: '/api/v1/tracks',
  trackById: (id: number | string) => `/api/v1/tracks/${id}`,
  playlists: '/api/v1/playlists',
  playlistById: (id: number | string) => `/api/v1/playlists/${id}`,
  genres: '/api/v1/genres',
  artists: '/api/v1/artists',
  onboarding: '/api/v1/me/onboarding',
  profile: '/api/v1/me/profile',
} as const

function normalizeBaseUrl(value: string): string {
  return value.replace(/\/+$/, '')
}
