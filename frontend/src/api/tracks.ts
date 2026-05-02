import { apiClient } from './client'
import type { Track } from '../types'

type TracksResponse = {
  tracks: Track[]
}

type TrackResponse = {
  track: Track
}

export function listTracks() {
  return apiClient.get<TracksResponse>('/api/v1/tracks')
}

export function getTrackById(id: number) {
  return apiClient.get<TrackResponse>(`/api/v1/tracks/${id}`)
}
