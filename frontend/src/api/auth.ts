import { apiClient } from './client'
import { API_ENDPOINTS } from '../config/api'
import type {
  AuthCredentials,
  LoginResponse,
  RegisterResponse,
  User,
} from '../types'

export function registerUser(credentials: AuthCredentials) {
  return apiClient.post<RegisterResponse>(
    API_ENDPOINTS.auth.register,
    credentials,
  )
}

export function loginUser(credentials: AuthCredentials) {
  return apiClient.post<LoginResponse>(API_ENDPOINTS.auth.login, credentials)
}

type CurrentUserResponse = {
  user: User
}

export async function getCurrentUser(token: string) {
  const response = await apiClient.get<CurrentUserResponse>(API_ENDPOINTS.me, {
    token,
  })

  return response.user
}
