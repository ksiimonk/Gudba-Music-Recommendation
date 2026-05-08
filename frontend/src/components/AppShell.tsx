import { useEffect, useRef, type ReactNode } from 'react'
import { useAuth } from '../auth/AuthContext'
import { ROUTES } from '../config/routes'
import { usePlayer } from '../context/PlayerContext'
import { BottomNavigation } from './BottomNavigation'
import { MiniPlayer } from './MiniPlayer'
import { SidebarNavigation } from './SidebarNavigation'

type AppShellProps = { children: ReactNode }

export function AppShell({ children }: AppShellProps) {
  const { isLoading, logout, user, token } = useAuth()
  const { loadFavorites } = usePlayer()
  const currentPath = window.location.pathname
  const loaded = useRef(false)

  useEffect(() => {
    if (token && !loaded.current) {
      loaded.current = true
      loadFavorites(token)
    }
    if (!token) {
      loaded.current = false
    }
  }, [token])

  function handleLogout() {
    logout()
    window.location.href = ROUTES.login
  }

  return (
    <div className="app-shell">
      <SidebarNavigation
        currentPath={currentPath}
        isAuthLoading={isLoading}
        userEmail={user?.email}
        onLogout={handleLogout}
      />

      <main className="app-content">
        <header className="mobile-topbar">
          {user ? (
            <div className="mobile-user">
              <span className="auth-avatar">{getUserInitial(user.email)}</span>
              <strong>{user.email}</strong>
            </div>
          ) : (
            <a className="home-avatar" href={ROUTES.login} aria-label="Профиль">
              S
            </a>
          )}
          <span>Gudba Music</span>
          {user ? (
            <button className="logout-button" type="button" onClick={handleLogout}>
              Выйти
            </button>
          ) : (
            <div className="mobile-auth-actions">
              <a href={ROUTES.login}>Вход</a>
              <a href={ROUTES.register}>Регистрация</a>
            </div>
          )}
        </header>

        {children}
      </main>

      <MiniPlayer />
      <BottomNavigation currentPath={currentPath} />
    </div>
  )
}

function getUserInitial(email: string) {
  return email.trim().charAt(0).toUpperCase() || 'U'
}
