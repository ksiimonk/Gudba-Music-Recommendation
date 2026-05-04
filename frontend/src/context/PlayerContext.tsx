import { createContext, useContext, useState } from 'react'
import type { ReactNode } from 'react'
import type { Track } from '../types'

type PlayerContextValue = {
  currentTrack: Track | null
  isPlaying: boolean
  playTrack: (track: Track) => void
  pauseTrack: () => void
}

const PlayerContext = createContext<PlayerContextValue | undefined>(undefined)

export function PlayerProvider({ children }: { children: ReactNode }) {
  const [currentTrack, setCurrentTrack] = useState<Track | null>(null)
  const [isPlaying, setIsPlaying] = useState(false)

  function playTrack(track: Track) {
    setCurrentTrack(track)
    setIsPlaying(true)
  }

  function pauseTrack() {
    setIsPlaying(false)
  }

  return (
    <PlayerContext.Provider value={{ currentTrack, isPlaying, playTrack, pauseTrack }}>
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
