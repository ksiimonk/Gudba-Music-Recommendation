import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import { AuthProvider } from './auth/AuthContext'
import { PlayerProvider } from './context/PlayerContext'
import { ErrorBoundary } from './components/ErrorBoundary'
import './index.css'
import App from './App.tsx'

createRoot(document.getElementById('root')!).render(
  <StrictMode>
    <ErrorBoundary>
      <AuthProvider>
        <PlayerProvider>
          <App />
        </PlayerProvider>
      </AuthProvider>
    </ErrorBoundary>
  </StrictMode>,
)
