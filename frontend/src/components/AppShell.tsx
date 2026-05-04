import type { ReactNode } from 'react'
import { BottomNavigation } from './BottomNavigation'
import { MiniPlayer } from './MiniPlayer'
import { SidebarNavigation } from './SidebarNavigation'

type AppShellProps = {
  children: ReactNode
}

export function AppShell({ children }: AppShellProps) {
  const currentPath = window.location.pathname

  return (
    <div className="app-shell">
      <SidebarNavigation currentPath={currentPath} />

      <main className="app-content">
        <header className="mobile-topbar">
          <a className="home-avatar" href="/login" aria-label="Профиль">
            S
          </a>
          <span>Gudba Music</span>
        </header>

        {children}
      </main>

      <MiniPlayer />
      <BottomNavigation currentPath={currentPath} />
    </div>
  )
}
