import { apiClient } from './client'
import { API_ENDPOINTS } from '../config/api'
import type { Playlist, Track } from '../types'

export type TrackRecommendation = {
  track: Track
  score: number
  explanation: string
}

export type PlaylistRecommendation = {
  playlist: Playlist
  score: number
  explanation: string
}

type TrackRecsResponse = {
  recommendations: TrackRecommendation[]
}

type PlaylistRecsResponse = {
  recommendations: PlaylistRecommendation[]
}

export function listTrackRecommendations(token: string, limit = 10) {
  return apiClient.get<TrackRecsResponse>(
    `${API_ENDPOINTS.trackRecommendations}?limit=${limit}`,
    { token },
  )
}

export function listPlaylistRecommendations(token: string, limit = 6) {
  return apiClient.get<PlaylistRecsResponse>(
    `${API_ENDPOINTS.playlistRecommendations}?limit=${limit}`,
    { token },
  )
}
