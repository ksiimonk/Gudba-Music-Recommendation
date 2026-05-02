import { apiClient } from './client'
import type {
  AuthCredentials,
  LoginResponse,
  RegisterResponse,
} from '../types'

export function registerUser(credentials: AuthCredentials) {
  return apiClient.post<RegisterResponse>('/api/v1/auth/register', credentials)
}

export function loginUser(credentials: AuthCredentials) {
  return apiClient.post<LoginResponse>('/api/v1/auth/login', credentials)
}
