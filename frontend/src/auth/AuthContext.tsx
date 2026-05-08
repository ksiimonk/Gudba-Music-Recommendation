import { createContext, useContext, useEffect, useState } from 'react'
import type { ReactNode } from 'react'
import { getCurrentUser } from '../api'
import { setOnUnauthorized } from '../api/client'
import type { User } from '../types'
import { clearToken, getCurrentToken, saveToken } from './tokenStorage'

type AuthContextValue = {
  user: User | null
  token: string | null
  isLoading: boolean
  login: (token: string) => Promise<void>
  logout: () => void
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined)

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null)
  const [token, setToken] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const storedToken = getCurrentToken()

    if (!storedToken) {
      setIsLoading(false)
      return
    }

    setToken(storedToken)

    getCurrentUser(storedToken)
      .then((currentUser) => {
        setUser(currentUser)
      })
      .catch(() => {
        clearToken()
        setToken(null)
        setUser(null)
      })
      .finally(() => {
        setIsLoading(false)
      })
  }, [])

  useEffect(() => {
    setOnUnauthorized(() => {
      clearToken()
      setToken(null)
      setUser(null)
    })

    return () => {
      setOnUnauthorized(null)
    }
  }, [])

  async function login(nextToken: string) {
    saveToken(nextToken)
    setToken(nextToken)

    try {
      const currentUser = await getCurrentUser(nextToken)
      setUser(currentUser)
    } catch (error) {
      clearToken()
      setToken(null)
      setUser(null)
      throw error
    }
  }

  function logout() {
    clearToken()
    setToken(null)
    setUser(null)
  }

  return (
    <AuthContext.Provider value={{ user, token, isLoading, login, logout }}>
      {children}
    </AuthContext.Provider>
  )
}

export function useAuth() {
  const context = useContext(AuthContext)

  if (!context) {
    throw new Error('useAuth must be used inside AuthProvider')
  }

  return context
}
