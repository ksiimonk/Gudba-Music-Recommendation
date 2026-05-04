import { MusicKeyLogo } from './MusicKeyLogo'
import { isNavigationItemActive, navigationItems } from './navigation'

type SidebarNavigationProps = {
  currentPath: string
}

export function SidebarNavigation({ currentPath }: SidebarNavigationProps) {
  return (
    <aside className="sidebar-navigation" aria-label="Основная навигация">
      <a className="sidebar-brand" href="/">
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
    </aside>
  )
}
