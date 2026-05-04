import { createContext, useContext, useState } from 'react'
import type { ReactNode } from 'react'
import { trackEvent } from '../api'
import { getCurrentToken } from '../auth/tokenStorage'
import type { Track } from '../types'

type PlayerContextValue = {
  currentTrack: Track | null
  isPlaying: boolean
  playTrack: (track: Track) => void
  pauseTrack: () => void
  likeTrack: (track: Track) => void
  dislikeTrack: (track: Track) => void
  skipTrack: (track: Track) => void
}

const PlayerContext = createContext<PlayerContextValue | undefined>(undefined)

export function PlayerProvider({ children }: { children: ReactNode }) {
  const [currentTrack, setCurrentTrack] = useState<Track | null>(null)
  const [isPlaying, setIsPlaying] = useState(false)

  function playTrack(track: Track) {
    setCurrentTrack(track)
    setIsPlaying(true)

    const token = getCurrentToken()
    if (token) {
      trackEvent(track.id, 'play', token).catch(() => {})
    }
  }

  function pauseTrack() {
    setIsPlaying(false)
  }

  function likeTrack(track: Track) {
    const token = getCurrentToken()
    if (!token) {
      window.location.href = '/login'
      return
    }

    trackEvent(track.id, 'like', token).catch((err) => {
      console.error('like failed:', err)
    })
  }

  function dislikeTrack(track: Track) {
    const token = getCurrentToken()
    if (!token) {
      window.location.href = '/login'
      return
    }

    trackEvent(track.id, 'dislike', token).catch((err) => {
      console.error('dislike failed:', err)
    })
  }

  function skipTrack(track: Track) {
    const token = getCurrentToken()
    if (!token) {
      window.location.href = '/login'
      return
    }

    trackEvent(track.id, 'skip', token).catch((err) => {
      console.error('skip failed:', err)
    })
  }

  return (
    <PlayerContext.Provider value={{ currentTrack, isPlaying, playTrack, pauseTrack, likeTrack, dislikeTrack, skipTrack }}>
      {children}
    </PlayerContext.Provider>
  )
}

export function usePlayer() {
  const context = useContext(PlayerContext)

  if (!context) {
    throw new Error('usePlayer must be used inside PlayerProvider')
  }

  return context
}
