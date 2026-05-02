import { useState } from 'react'
import { getApiErrorMessage, loginUser } from '../api'
import { AuthForm } from '../components/AuthForm'
import { AuthShell } from '../components/AuthShell'

const AUTH_TOKEN_STORAGE_KEY = 'music_recommender_token'

export function LoginPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')

  async function handleLogin() {
    setIsSubmitting(true)
    setMessage('')
    setError('')

    try {
      const response = await loginUser({ email, password })
      localStorage.setItem(AUTH_TOKEN_STORAGE_KEY, response.token)
      setMessage('Вы вошли. Токен сессии сохранен в этом браузере.')
      setPassword('')
    } catch (requestError) {
      setError(getApiErrorMessage(requestError, 'Не удалось войти.'))
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <AuthShell
      title="С возвращением в твое музыкальное пространство."
      subtitle="Войди, чтобы продолжить настраивать рекомендации."
    >
      <AuthForm
        title="Вход"
        submitLabel="Войти"
        switchText="Еще нет аккаунта?"
        switchHref="/register"
        switchLabel="Зарегистрироваться"
        passwordAutoComplete="current-password"
        email={email}
        password={password}
        isSubmitting={isSubmitting}
        message={message}
        error={error}
        onEmailChange={setEmail}
        onPasswordChange={setPassword}
        onSubmit={handleLogin}
      />
    </AuthShell>
  )
}
