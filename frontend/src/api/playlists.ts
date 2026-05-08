import { apiClient } from './client'
import { API_ENDPOINTS } from '../config/api'
import type { Playlist } from '../types'

type PlaylistsResponse = {
  playlists: Playlist[]
}

type PlaylistResponse = {
  playlist: Playlist
}

export function listPlaylists() {
  return apiClient.get<PlaylistsResponse>(API_ENDPOINTS.playlists)
}

export function getPlaylistById(id: number) {
  return apiClient.get<PlaylistResponse>(API_ENDPOINTS.playlistById(id))
}

export function getFavorites(token: string) {
  return apiClient.get<PlaylistResponse>('/api/v1/me/favorites', { token })
}
