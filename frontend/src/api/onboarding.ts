import { apiClient } from './client'
import { API_ENDPOINTS } from '../config/api'
import type { OnboardingRequest, UserProfile } from '../types'

type ProfileResponse = {
  profile: UserProfile
}

export function saveOnboarding(data: OnboardingRequest, token: string) {
  return apiClient.post<{ message: string }>(API_ENDPOINTS.onboarding, data, { token })
}

export function getProfile(token: string) {
  return apiClient.get<ProfileResponse>(API_ENDPOINTS.profile, { token })
}
