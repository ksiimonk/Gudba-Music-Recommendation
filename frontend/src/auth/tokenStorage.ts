import { STORAGE_KEYS } from '../config/storage'

export function getCurrentToken() {
  return localStorage.getItem(STORAGE_KEYS.authToken)
}

export function saveToken(token: string) {
  localStorage.setItem(STORAGE_KEYS.authToken, token)
}

export function clearToken() {
  localStorage.removeItem(STORAGE_KEYS.authToken)
}
