import { useState } from 'react'
import { getApiErrorMessage, registerUser } from '../api'
import { AuthForm } from '../components/AuthForm'
import { AuthShell } from '../components/AuthShell'

export function RegisterPage() {
  const [email, setEmail] = useState('')
  const [password, setPassword] = useState('')
  const [isSubmitting, setIsSubmitting] = useState(false)
  const [message, setMessage] = useState('')
  const [error, setError] = useState('')

  async function handleRegister() {
    setIsSubmitting(true)
    setMessage('')
    setError('')

    try {
      const response = await registerUser({ email, password })
      setMessage(`${response.user.email} зарегистрирован. Теперь можно войти.`)
      setPassword('')
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
        switchHref="/login"
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
