import { apiClient } from './client'
import { API_ENDPOINTS } from '../config/api'
import type { Genre } from '../types'

type GenresResponse = {
  genres: Genre[]
}

export function listGenres() {
  return apiClient.get<GenresResponse>(API_ENDPOINTS.genres)
}
