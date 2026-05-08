import {
  ROUTES,
  isNestedRoute,
  readRouteId,
} from './config/routes'
import { FavoritesPage } from './pages/FavoritesPage'
import { HomePage } from './pages/HomePage'
import { LoginPage } from './pages/LoginPage'
import { OnboardingPage } from './pages/OnboardingPage'
import { PlaylistDetailPage } from './pages/PlaylistDetailPage'
import { PlaylistListPage } from './pages/PlaylistListPage'
import { RegisterPage } from './pages/RegisterPage'
import { TrackDetailPage } from './pages/TrackDetailPage'
import { TrackListPage } from './pages/TrackListPage'
import './App.css'

function App() {
  const pathname = window.location.pathname

  if (pathname === ROUTES.login) {
    return <LoginPage />
  }

  if (pathname === ROUTES.register) {
    return <RegisterPage />
  }

  if (pathname === ROUTES.onboarding) {
    return <OnboardingPage />
  }

  if (pathname === ROUTES.tracks) {
    return <TrackListPage />
  }

  if (isNestedRoute(pathname, ROUTES.tracks)) {
    return <TrackDetailPage trackId={readRouteId(pathname, ROUTES.tracks)} />
  }

  if (pathname === ROUTES.playlists) {
    return <PlaylistListPage />
  }

  if (isNestedRoute(pathname, ROUTES.playlists)) {
    return (
      <PlaylistDetailPage
        playlistId={readRouteId(pathname, ROUTES.playlists)}
      />
    )
  }

  if (pathname === ROUTES.favorites) {
    return <FavoritesPage />
  }

  return <HomePage />
}

export default App
