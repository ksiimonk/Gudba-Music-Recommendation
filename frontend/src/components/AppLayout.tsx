import type { ReactNode } from 'react'

type AppLayoutProps = {
  children: ReactNode
}

export function AppLayout({ children }: AppLayoutProps) {
  return (
    <div className="app-layout">
      <header className="app-header">
        <a className="brand" href="/">
          Gudba Music
        </a>
        <nav className="main-nav" aria-label="Main navigation">
          <a href="/tracks">Tracks</a>
          <a href="/playlists">Playlists</a>
          <a href="/login">Login</a>
        </nav>
      </header>

      <main className="app-main">{children}</main>
    </div>
  )
}
