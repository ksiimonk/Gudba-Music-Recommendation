import type { FormEvent } from 'react'

type AuthFormProps = {
  title: string
  submitLabel: string
  switchText: string
  switchHref: string
  switchLabel: string
  passwordAutoComplete: 'current-password' | 'new-password'
  email: string
  password: string
  isSubmitting: boolean
  message: string
  error: string
  onEmailChange: (value: string) => void
  onPasswordChange: (value: string) => void
  onSubmit: () => void
}

export function AuthForm({
  title,
  submitLabel,
  switchText,
  switchHref,
  switchLabel,
  passwordAutoComplete,
  email,
  password,
  isSubmitting,
  message,
  error,
  onEmailChange,
  onPasswordChange,
  onSubmit,
}: AuthFormProps) {
  function handleSubmit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault()
    onSubmit()
  }

  return (
    <form className="auth-form" onSubmit={handleSubmit}>
      <h2>{title}</h2>

      <label>
        <span>Электронная почта</span>
        <input
          type="email"
          value={email}
          onChange={(event) => onEmailChange(event.target.value)}
          placeholder="you@example.com"
          autoComplete="email"
          required
        />
      </label>

      <label>
        <span>Пароль</span>
        <input
          type="password"
          value={password}
          onChange={(event) => onPasswordChange(event.target.value)}
          placeholder="Введите пароль"
          autoComplete={passwordAutoComplete}
          required
        />
      </label>

      {error && <p className="auth-feedback auth-feedback-error">{error}</p>}
      {message && (
        <p className="auth-feedback auth-feedback-success">{message}</p>
      )}

      <button className="auth-submit" type="submit" disabled={isSubmitting}>
        {isSubmitting ? 'Подождите...' : submitLabel}
      </button>

      <p className="auth-switch">
        {switchText} <a href={switchHref}>{switchLabel}</a>
      </p>
    </form>
  )
}
