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

  if (pathname === '/login') {
    return <LoginPage />
  }

  if (pathname === '/register') {
    return <RegisterPage />
  }

  if (pathname === '/tracks') {
    return <TrackListPage />
  }

  if (pathname.startsWith('/tracks/')) {
    return <TrackDetailPage trackId={Number(pathname.split('/')[2])} />
  }

  if (pathname === '/playlists') {
    return <PlaylistListPage />
  }

  if (pathname.startsWith('/playlists/')) {
    return <PlaylistDetailPage playlistId={Number(pathname.split('/')[2])} />
  }

  return <HomePage />
}

export default App
