import { useEffect, useState } from 'react'
import { getApiErrorMessage, loginUser, registerUser } from '../api'
import { useAuth } from '../auth/AuthContext'
import { AuthForm } from '../components/AuthForm'
import { AuthShell } from '../components/AuthShell'
import { ROUTES } from '../config/routes'

export function RegisterPage() {
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

  async function handleRegister() {
    setIsSubmitting(true)
    setMessage('')
    setError('')

    try {
      await registerUser({ email, password })
      const loginResponse = await loginUser({ email, password })
      await login(loginResponse.token)
      setPassword('')
      window.location.href = ROUTES.home
    } catch (requestError) {
      setError(
        getApiErrorMessage(requestError, 'Не удалось создать аккаунт.'),
      )
    } finally {
      setIsSubmitting(false)
    }
  }

  return (
    <AuthShell
      title="Миллионы треков. Рекомендации под твой вкус."
      subtitle="Создай аккаунт, чтобы начать собирать музыкальный профиль."
    >
      <AuthForm
        title="Регистрация"
        submitLabel="Зарегистрироваться"
        switchText="Уже есть аккаунт?"
        switchHref={ROUTES.login}
        switchLabel="Войти"
        passwordAutoComplete="new-password"
        email={email}
        password={password}
        isSubmitting={isSubmitting}
        message={message}
        error={error}
        onEmailChange={setEmail}
        onPasswordChange={setPassword}
        onSubmit={handleRegister}
      />
    </AuthShell>
  )
}
