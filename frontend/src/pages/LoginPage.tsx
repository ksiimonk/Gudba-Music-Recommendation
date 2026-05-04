import { useEffect, useState } from 'react'
import { getApiErrorMessage, loginUser } from '../api'
import { useAuth } from '../auth/AuthContext'
import { AuthForm } from '../components/AuthForm'
import { AuthShell } from '../components/AuthShell'
import { ROUTES } from '../config/routes'

export function LoginPage() {
  const { isLoading, login, user } = useAuth()
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')

  useEffect(() => {
    if (!isLoading && user) {
      window.location.href = ROUTES.home
    }
  }, [isLoading, user])

  async function handleLogin() {
    setIsSubmitting(true)
    setMessage('')
    setError('')

    try {
      const response = await loginUser({ email, password })
      await login(response.token)
      setPassword('')
      window.location.href = ROUTES.home
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
        switchHref={ROUTES.register}
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
