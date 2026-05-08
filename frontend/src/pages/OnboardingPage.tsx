import { useEffect, useState } from 'react'
import { getApiErrorMessage, getProfile, listArtists, listGenres, listTracks, saveOnboarding } from '../api'
import { useAuth } from '../auth/AuthContext'
import { ROUTES } from '../config/routes'
import type { Artist, Genre, Track } from '../types'

type Step = 1 | 2 | 3 | 4

const CONTEXTS = [
  { key: 'focus', label: 'Фокус' },
  { key: 'relax', label: 'Отдых' },
  { key: 'workout', label: 'Тренировка' },
  { key: 'evening', label: 'Вечер' },
  { key: 'study', label: 'Учёба' },
  { key: 'walk', label: 'Прогулка' },
]

export function OnboardingPage() {
  const { isLoading: authLoading, token, user } = useAuth()

  const [step, setStep] = useState<Step>(1)
  const [loading, setLoading] = useState(true)
  const [submitting, setSubmitting] = useState(false)
  const [error, setError] = useState('')

  const [genres, setGenres] = useState<Genre[]>([])
  const [artists, setArtists] = useState<Artist[]>([])
  const [tracks, setTracks] = useState<Track[]>([])

  const [selectedGenres, setSelectedGenres] = useState<Set<number>>(new Set())
  const [selectedArtists, setSelectedArtists] = useState<Set<number>>(new Set())
  const [selectedTracks, setSelectedTracks] = useState<Set<number>>(new Set())
  const [selectedContexts, setSelectedContexts] = useState<Set<string>>(new Set())

  useEffect(() => {
    if (authLoading || !token) return

    let isMounted = true
    setLoading(true)

    Promise.all([
      listGenres(),
      listArtists(),
      listTracks(),
      getProfile(token),
    ])
      .then(([g, a, t, p]) => {
        if (!isMounted) return
        setGenres(g.genres)
        setArtists(a.artists)
        setTracks(t.tracks)

        const profile = p.profile
        if (profile) {
          setSelectedGenres(new Set(profile.favorite_genre_ids))
          setSelectedArtists(new Set(profile.favorite_artist_ids))
          setSelectedTracks(new Set(profile.starter_track_ids))
          setSelectedContexts(new Set(profile.contexts))
        }
      })
      .catch(() => {
        if (isMounted) setError('Не удалось загрузить данные. Проверь, запущен ли backend.')
      })
      .finally(() => {
        if (isMounted) setLoading(false)
      })

    return () => { isMounted = false }
  }, [authLoading, token])

  if (authLoading) {
    return <OnboardingShell><p className="page-state">Проверка авторизации...</p></OnboardingShell>
  }

  if (!user || !token) {
    return (
      <OnboardingShell>
        <div className="onboarding-auth-prompt">
          <h2>Настрой рекомендации</h2>
          <p>Войди, чтобы выбрать любимые жанры, артистов и треки.</p>
          <a className="auth-link" href={ROUTES.login}>Войти</a>
          <a className="auth-link secondary" href={ROUTES.register}>Зарегистрироваться</a>
        </div>
      </OnboardingShell>
    )
  }

  if (loading) {
    return <OnboardingShell><p className="page-state">Загрузка...</p></OnboardingShell>
  }

  if (error) {
    return (
      <OnboardingShell>
        <p className="page-state">{error}</p>
        <button className="auth-submit" type="button" onClick={() => window.location.reload()}>Повторить</button>
      </OnboardingShell>
    )
  }

  function toggleGenre(id: number) {
    setSelectedGenres((prev) => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id); else next.add(id)
      return next
    })
  }

  function toggleArtist(id: number) {
    setSelectedArtists((prev) => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id); else next.add(id)
      return next
    })
  }

  function toggleTrack(id: number) {
    setSelectedTracks((prev) => {
      const next = new Set(prev)
      if (next.has(id)) next.delete(id); else next.add(id)
      return next
    })
  }

  function toggleContext(key: string) {
    setSelectedContexts((prev) => {
      const next = new Set(prev)
      if (next.has(key)) next.delete(key); else next.add(key)
      return next
    })
  }

  function canProceed(): boolean {
    switch (step) {
      case 1: return selectedGenres.size > 0
      case 2: return true
      case 3: return true
      case 4: return true
    }
  }

  function handleNext() {
    if (step < 4) setStep((step + 1) as Step)
  }

  function handleBack() {
    if (step > 1) setStep((step - 1) as Step)
  }

  async function handleSubmit() {
    if (!token) return
    setSubmitting(true)
    setError('')

    try {
      await saveOnboarding(
        {
          genre_ids: Array.from(selectedGenres),
          artist_ids: Array.from(selectedArtists),
          track_ids: Array.from(selectedTracks),
          contexts: Array.from(selectedContexts),
        },
        token,
      )
      window.location.href = ROUTES.home
    } catch (requestError) {
      setError(getApiErrorMessage(requestError, 'Не удалось сохранить предпочтения.'))
    } finally {
      setSubmitting(false)
    }
  }

  return (
    <OnboardingShell>
      <div className="onboarding-wizard">
        <div className="onboarding-progress">
          <div className="onboarding-progress-bar">
            <div className="onboarding-progress-fill" style={{ width: `${(step / 4) * 100}%` }} />
          </div>
          <span className="onboarding-step-label">{step} из 4</span>
        </div>

        {step === 1 && (
          <div className="onboarding-step">
            <h2>Выбери любимые жанры</h2>
            <p>Это поможет нам понять, какая музыка тебе нравится.</p>
            <div className="onboarding-chips">
              {genres.map((genre) => (
                <button
                  key={genre.id}
                  type="button"
                  className={`onboarding-chip ${selectedGenres.has(genre.id) ? 'selected' : ''}`}
                  onClick={() => toggleGenre(genre.id)}
                >
                  {genre.name}
                </button>
              ))}
            </div>
          </div>
        )}

        {step === 2 && (
          <div className="onboarding-step">
            <h2>Выбери любимых артистов</h2>
            <p>Отметь тех, кого ты чаще всего слушаешь.</p>
            <div className="onboarding-artist-grid">
              {artists.slice(0, 20).map((artist) => (
                <button
                  key={artist.id}
                  type="button"
                  className={`onboarding-artist-card ${selectedArtists.has(artist.id) ? 'selected' : ''}`}
                  onClick={() => toggleArtist(artist.id)}
                >
                  <img src={artist.image_url || 'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=200'} alt="" />
                  <span>{artist.name}</span>
                </button>
              ))}
            </div>
          </div>
        )}

        {step === 3 && (
          <div className="onboarding-step">
            <h2>Отметь треки, которые нравятся</h2>
            <p>Выбери несколько треков для старта.</p>
            <div className="onboarding-track-list">
              {tracks.slice(0, 30).map((track) => (
                <button
                  key={track.id}
                  type="button"
                  className={`onboarding-track-row ${selectedTracks.has(track.id) ? 'selected' : ''}`}
                  onClick={() => toggleTrack(track.id)}
                >
                  <img src={track.cover_url || 'https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=200'} alt="" />
                  <div className="onboarding-track-info">
                    <strong>{track.title}</strong>
                    <span>{track.artist.name}</span>
                  </div>
                  {selectedTracks.has(track.id) && <span className="onboarding-check">✓</span>}
                </button>
              ))}
            </div>
          </div>
        )}

        {step === 4 && (
          <div className="onboarding-step">
            <h2>Где ты слушаешь музыку?</h2>
            <p>Выбери подходящие ситуации.</p>
            <div className="onboarding-chips">
              {CONTEXTS.map((ctx) => (
                <button
                  key={ctx.key}
                  type="button"
                  className={`onboarding-chip ${selectedContexts.has(ctx.key) ? 'selected' : ''}`}
                  onClick={() => toggleContext(ctx.key)}
                >
                  {ctx.label}
                </button>
              ))}
            </div>
          </div>
        )}

        <div className="onboarding-nav">
          {step > 1 && (
            <button className="auth-submit secondary" type="button" onClick={handleBack} disabled={submitting}>
              Назад
            </button>
          )}
          <div className="onboarding-nav-right">
            {step < 4 ? (
              <button className="auth-submit" type="button" onClick={handleNext} disabled={!canProceed()}>
                Далее
              </button>
            ) : (
              <button className="auth-submit" type="button" onClick={handleSubmit} disabled={submitting}>
                {submitting ? 'Сохраняем...' : 'Готово'}
              </button>
            )}
          </div>
        </div>

        {error && <p className="onboarding-error">{error}</p>}
      </div>
    </OnboardingShell>
  )
}

function OnboardingShell({ children }: { children: React.ReactNode }) {
  return (
    <div className="onboarding-page">
      <header className="onboarding-header">
        <a className="sidebar-brand" href={ROUTES.home}>
          <span>Gudba Music</span>
        </a>
      </header>
      <main className="onboarding-content">
        {children}
      </main>
    </div>
  )
}
