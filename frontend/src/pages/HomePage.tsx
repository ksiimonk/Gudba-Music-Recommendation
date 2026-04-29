import { AppLayout } from '../components/AppLayout'
import { StatusCard } from '../components/StatusCard'

export function HomePage() {
  return (
    <AppLayout>
      <section className="home-hero">
        <div>
          <p className="eyebrow">Music recommender</p>
          <h1>Find tracks and playlists that fit your taste.</h1>
        </div>
      </section>

      <section className="status-grid" aria-label="Project status">
        <StatusCard label="Backend" value="Go API" />
        <StatusCard label="Database" value="PostgreSQL" />
        <StatusCard label="Frontend" value="React + Vite" />
      </section>
    </AppLayout>
  )
}
