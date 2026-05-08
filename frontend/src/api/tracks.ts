import { apiClient } from './client'
import { API_ENDPOINTS } from '../config/api'
import type { Track } from '../types'

type TracksResponse = {
  tracks: Track[]
}

type TrackResponse = {
  track: Track
}

export function listTracks() {
  return apiClient.get<TracksResponse>(API_ENDPOINTS.tracks)
}

export function getTrackById(id: number) {
  return apiClient.get<TrackResponse>(API_ENDPOINTS.trackById(id))
}
