import { isNavigationItemActive, navigationItems } from './navigation'

type BottomNavigationProps = {
  currentPath: string
}

export function BottomNavigation({ currentPath }: BottomNavigationProps) {
  return (
    <nav className="bottom-navigation" aria-label="Основная навигация">
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
  )
}
