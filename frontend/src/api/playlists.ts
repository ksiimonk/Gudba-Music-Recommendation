import { apiClient } from './client'
import type { Playlist } from '../types'

type PlaylistsResponse = {
  playlists: Playlist[]
}

type PlaylistResponse = {
  playlist: Playlist
}

export function listPlaylists() {
  return apiClient.get<PlaylistsResponse>('/api/v1/playlists')
}

export function getPlaylistById(id: number) {
  return apiClient.get<PlaylistResponse>(`/api/v1/playlists/${id}`)
}
