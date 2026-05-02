import { ApiError } from './client'

type ErrorPayload = {
  message?: string
}

export function getApiErrorMessage(error: unknown, fallback: string) {
  if (error instanceof ApiError) {
    const payload = error.data as ErrorPayload
    return payload.message ?? fallback
  }

  if (error instanceof Error) {
    return error.message
  }

  return fallback
}
