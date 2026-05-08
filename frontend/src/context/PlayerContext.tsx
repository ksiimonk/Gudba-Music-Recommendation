import { createContext, useContext, useState, useCallback, useEffect, useRef } from 'react'
import type { ReactNode } from 'react'
import { getFavorites, toggleLikeTrack, trackEvent } from '../api'
import { getCurrentToken } from '../auth/tokenStorage'
import type { Track } from '../types'

const TRACK_KEY = 'gudba_player_track'
const PLAYING_KEY = 'gudba_player_playing'
const TIME_KEY = 'gudba_player_time'

function loadJSON<T>(key: string): T | null {
  try { const r = localStorage.getItem(key); return r ? JSON.parse(r) : null }
  catch { return null }
}
function saveJSON(key: string, v: unknown) {
  if (v !== undefined && v !== null) localStorage.setItem(key, JSON.stringify(v))
  else localStorage.removeItem(key)
}

type PlayerContextValue = {
  currentTrack: Track | null
  isPlaying: boolean
  isLoading: boolean
  progress: number
  duration: number
  playTrack: (track: Track) => void
  pauseTrack: () => void
  toggleLike: (track: Track) => void
  dislikeTrack: (track: Track) => void
  skipTrack: (track: Track) => void
  isFavorited: (trackId: number) => boolean
  loadFavorites: (token: string) => Promise<void>
}

const PlayerContext = createContext<PlayerContextValue | undefined>(undefined)

export function PlayerProvider({ children }: { children: ReactNode }) {
  const [currentTrack, setCurrentTrack] = useState<Track | null>(() => loadJSON<Track>(TRACK_KEY))
  const [isPlaying, setIsPlaying] = useState(() => localStorage.getItem(PLAYING_KEY) === 'true')
  const [isLoading, setIsLoading] = useState(false)
  const [progress, setProgress] = useState(0)
  const [duration, setDuration] = useState(0)
  const [favoritedIds, setFavoritedIds] = useState<Set<number>>(new Set())
  const audioRef = useRef<HTMLAudioElement | null>(null)

  useEffect(() => { saveJSON(TRACK_KEY, currentTrack) }, [currentTrack])
  useEffect(() => {
    localStorage.setItem(PLAYING_KEY, isPlaying ? 'true' : 'false')
  }, [isPlaying])

  useEffect(() => {
    if (!currentTrack) {
      setIsPlaying(false)
      return
    }

    const storedUrl = currentTrack.preview_url
    if (!storedUrl) {
      setIsPlaying(false)
      return
    }

    const savedTime = loadJSON<number>(TIME_KEY) ?? 0
    const audio = new Audio(storedUrl)
    audio.currentTime = savedTime
    audioRef.current = audio

    const onEnded = () => { setIsPlaying(false); setProgress(0) }
    const onWaiting = () => setIsLoading(true)
    const onCanPlay = () => setIsLoading(false)
    const onTimeUpdate = () => setProgress(audio.currentTime)

    audio.addEventListener('loadedmetadata', () => setDuration(audio.duration))
    audio.addEventListener('timeupdate', onTimeUpdate)
    audio.addEventListener('ended', onEnded)
    audio.addEventListener('waiting', onWaiting)
    audio.addEventListener('canplay', onCanPlay)

    if (isPlaying) {
      audio.play().catch(() => setIsPlaying(false))
    }

    return () => {
      saveJSON(TIME_KEY, audio.currentTime)
      audio.pause()
      audio.src = ''
      audio.removeEventListener('timeupdate', onTimeUpdate)
      audio.removeEventListener('ended', onEnded)
      audio.removeEventListener('waiting', onWaiting)
      audio.removeEventListener('canplay', onCanPlay)
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [])

  const playTrack = useCallback((track: Track) => {
    // Stop previous
    if (audioRef.current) {
      saveJSON(TIME_KEY, audioRef.current.currentTime)
      audioRef.current.pause()
      audioRef.current.src = ''
    }

    setCurrentTrack(track)
    setProgress(0)

    if (!track.preview_url) {
      setIsPlaying(false)
      return
    }

    setIsPlaying(true)

    const token = getCurrentToken()
    if (token) {
      trackEvent(track.id, 'play', token).catch(() => {})
    }

    const audio = new Audio(track.preview_url)
    audioRef.current = audio

    audio.addEventListener('loadedmetadata', () => setDuration(audio.duration))
    audio.addEventListener('timeupdate', () => setProgress(audio.currentTime))
    audio.addEventListener('ended', () => { setIsPlaying(false); setProgress(0) })
    audio.addEventListener('waiting', () => setIsLoading(true))
    audio.addEventListener('canplay', () => setIsLoading(false))

    audio.play().catch(() => setIsPlaying(false))
  }, [])

  const pauseTrack = useCallback(() => {
    setIsPlaying(false)
    if (audioRef.current) {
      saveJSON(TIME_KEY, audioRef.current.currentTime)
      audioRef.current.pause()
    }
  }, [])

  const toggleLike = useCallback((track: Track) => {
    const token = getCurrentToken()
    if (!token) { window.location.href = '/login'; return }

    toggleLikeTrack(track.id, token).then((res) => {
      setFavoritedIds((prev) => {
        const next = new Set(prev)
        if (res.liked) next.add(track.id)
        else next.delete(track.id)
        return next
      })
    }).catch((err) => console.error('toggle like failed:', err))
  }, [])

  function isFavorited(trackId: number): boolean {
    return favoritedIds.has(trackId)
  }

  const loadFavorites = useCallback(async (token: string) => {
    try {
      const res = await getFavorites(token)
      const ids = new Set<number>()
      if (res.playlist?.tracks) {
        for (const t of res.playlist.tracks) ids.add(t.id)
      }
      setFavoritedIds(ids)
    } catch {
      setFavoritedIds(new Set())
    }
  }, [])

  function dislikeTrack(track: Track) {
    const token = getCurrentToken()
    if (!token) { window.location.href = '/login'; return }
    trackEvent(track.id, 'dislike', token).catch((err) => console.error('dislike failed:', err))
  }

  function skipTrack(track: Track) {
    const token = getCurrentToken()
    if (!token) { window.location.href = '/login'; return }
    trackEvent(track.id, 'skip', token).catch((err) => console.error('skip failed:', err))
  }

  return (
    <PlayerContext.Provider value={{
      currentTrack, isPlaying, isLoading, progress, duration, playTrack, pauseTrack,
      toggleLike, dislikeTrack, skipTrack, isFavorited, loadFavorites,
    }}>
      {children}
    </PlayerContext.Provider>
  )
}

export function usePlayer() {
  const context = useContext(PlayerContext)
  if (!context) throw new Error('usePlayer must be used inside PlayerProvider')
  return context
}
