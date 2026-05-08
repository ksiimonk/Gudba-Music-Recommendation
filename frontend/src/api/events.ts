import { apiClient } from './client'
import { API_ENDPOINTS } from '../config/api'

export function trackEvent(trackId: number, eventType: string, token: string) {
  const path = `${API_ENDPOINTS.trackById(trackId)}/${eventType}`
  return apiClient.post<{ message: string }>(path, undefined, { token })
}

export function toggleLikeTrack(trackId: number, token: string) {
  const path = `${API_ENDPOINTS.trackById(trackId)}/like`
  return apiClient.post<{ liked: boolean }>(path, undefined, { token })
}


