import type { ReactNode } from 'react'
import { ROUTES } from '../config/routes'

type AppLayoutProps = {
  children: ReactNode
}

export function AppLayout({ children }: AppLayoutProps) {
  return (
    <div className="app-layout">
      <header className="app-header">
        <a className="brand" href={ROUTES.home}>
          Gudba Music
        </a>
        <nav className="main-nav" aria-label="Main navigation">
          <a href={ROUTES.tracks}>Tracks</a>
          <a href={ROUTES.playlists}>Playlists</a>
          <a href={ROUTES.login}>Login</a>
          <a href={ROUTES.register}>Register</a>
        </nav>
      </header>

      <main className="app-main">{children}</main>
    </div>
  )
}
