import { MusicKeyLogo } from './MusicKeyLogo'
import { ROUTES } from '../config/routes'
import { isNavigationItemActive, navigationItems } from './navigation'

type SidebarNavigationProps = {
  currentPath: string
  isAuthLoading: boolean
  userEmail?: string
  onLogout: () => void
}

export function SidebarNavigation({
  currentPath,
  isAuthLoading,
  userEmail,
  onLogout,
}: SidebarNavigationProps) {
  return (
    <aside className="sidebar-navigation" aria-label="Основная навигация">
      <a className="sidebar-brand" href={ROUTES.home}>
        <MusicKeyLogo className="sidebar-brand-mark" />
        <span>Gudba Music</span>
      </a>

      <nav className="sidebar-links">
        {navigationItems.map((item) => {
          const isActive = isNavigationItemActive(currentPath, item.href)

          return (
            <a
              className={isActive ? 'active' : ''}
              href={item.href}
              aria-current={isActive ? 'page' : undefined}
              key={item.href}
            >
              <span>{item.icon}</span>
              {item.label}
            </a>
          )
        })}
      </nav>

      <section className="sidebar-library" aria-label="Библиотека">
        <span>Твоя библиотека</span>
        <strong>Рекомендации, треки и плейлисты</strong>
      </section>

      <section className="sidebar-auth" aria-label="Аккаунт">
        {userEmail ? (
          <>
            <div className="auth-user">
              <span className="auth-avatar">{getUserInitial(userEmail)}</span>
              <div>
                <strong>{userEmail}</strong>
                <small>Вы вошли</small>
              </div>
            </div>
            <button className="logout-button" type="button" onClick={onLogout}>
              Выйти
            </button>
          </>
        ) : (
          <>
            <span>{isAuthLoading ? 'Проверяем вход...' : 'Аккаунт'}</span>
            <a className="auth-link" href={ROUTES.login}>
              Войти
            </a>
            <a className="auth-link secondary" href={ROUTES.register}>
              Зарегистрироваться
            </a>
          </>
        )}
      </section>
    </aside>
  )
}

function getUserInitial(email: string) {
  return email.trim().charAt(0).toUpperCase() || 'U'
}
