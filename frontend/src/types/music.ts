export type Artist = {
  id: number
  name: string
  image_url?: string
  bio?: string
  spotify_url?: string
  created_at: string
  updated_at: string
}

export type Genre = {
  id: number
  name: string
  created_at?: string
}

export type Track = {
  id: number
  title: string
  duration_ms: number
  spotify_url?: string
  cover_url?: string
  popularity_score: number
  artist: Artist
  genres: Genre[]
  created_at: string
  updated_at: string
}

export type Playlist = {
  id: number
  user_id: number
  name: string
  description?: string
  is_public: boolean
  track_count: number
  tracks?: Track[]
  created_at: string
  updated_at: string
}
