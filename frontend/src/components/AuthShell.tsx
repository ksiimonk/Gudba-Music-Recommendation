import type { ReactNode } from 'react'
import { ROUTES } from '../config/routes'
import { MusicKeyLogo } from './MusicKeyLogo'

type AuthShellProps = {
  title: string
  subtitle: string
  children: ReactNode
}

export function AuthShell({ title, subtitle, children }: AuthShellProps) {
  return (
    <main className="auth-shell">
      <a className="auth-brand" href={ROUTES.home}>
        <MusicKeyLogo className="auth-brand-mark" label="Home" />
        <span>Gudba Music</span>
      </a>

      <section className="auth-panel">
        <div className="auth-copy">
          <MusicKeyLogo className="auth-logo" />
          <h1>{title}</h1>
          <p>{subtitle}</p>
        </div>

        {children}
      </section>
    </main>
  )
}
