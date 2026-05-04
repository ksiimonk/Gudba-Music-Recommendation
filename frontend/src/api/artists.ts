import { apiClient } from './client'
import { API_ENDPOINTS } from '../config/api'
import type { Artist } from '../types'

type ArtistsResponse = {
  artists: Artist[]
}

export function listArtists() {
  return apiClient.get<ArtistsResponse>(API_ENDPOINTS.artists)
}
