export type OnboardingRequest = {
  genre_ids: number[]
  artist_ids: number[]
  track_ids: number[]
  contexts: string[]
}

export type UserProfile = {
  user_id: number
  favorite_genre_ids: number[]
  favorite_artist_ids: number[]
  starter_track_ids: number[]
  contexts: string[]
  updated_at: string
}
