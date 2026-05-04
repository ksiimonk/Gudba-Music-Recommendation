import {
  ROUTES,
  isNestedRoute,
  readRouteId,
} from './config/routes'
import { HomePage } from './pages/HomePage'
import { LoginPage } from './pages/LoginPage'
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

  return <HomePage />
}

export default App
