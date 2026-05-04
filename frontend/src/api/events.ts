import { apiClient } from './client'
import { API_ENDPOINTS } from '../config/api'

export function trackEvent(trackId: number, eventType: string, token: string) {
  const path = `${API_ENDPOINTS.trackById(trackId)}/${eventType}`
  return apiClient.post<{ message: string }>(path, undefined, { token })
}

export function playlistEvent(playlistId: number, eventType: string, token: string) {
  const path = `${API_ENDPOINTS.playlistById(playlistId)}/${eventType}`
  return apiClient.post<{ message: string }>(path, undefined, { token })
}
