# Gudba Music Recommendation

Дипломный проект — персонализированная система рекомендаций музыкальных треков с объяснением результатов.

**Стек:** Go + Gin + PostgreSQL / React + Vite + TypeScript

---

## Требования

| Компонент | Версия |
|-----------|--------|
| Go | 1.23+ |
| Node.js | 22+ |
| npm | 9+ |
| PostgreSQL | 14+ |

---

## Быстрый старт

### 1. Настройка PostgreSQL

Убедись, что PostgreSQL запущен. По умолчанию конфигурация ожидает:

```
Хост: 127.0.0.1
Порт: 5432
Пользователь: postgres
Пароль: admin
База: music_recommender
```

Параметры подключения задаются в `backend/.env`:

```
PG_URL=postgres://postgres:admin@127.0.0.1:5432/music_recommender?sslmode=disable
```

### 2. Установка зависимостей

```powershell
npm --prefix frontend install
```

Go-модули скачиваются автоматически при первой сборке.

### 3. Запуск через PowerShell-runner

Все команды выполняются из **корня репозитория**:

```powershell
.\scripts\dev.ps1 help
```

### 4. Полный цикл запуска

**Терминал 1 — База данных и бэкенд:**

```powershell
.\scripts\dev.ps1 doctor       # Проверка инструментов
.\scripts\dev.ps1 reset-db     # Сброс БД + миграции + seed
.\scripts\dev.ps1 test         # Прогон backend-тестов
.\scripts\dev.ps1 backend      # Запуск API на http://127.0.0.1:8080
```

**Терминал 2 — Smoke-тест (при запущенном backend):**

```powershell
.\scripts\dev.ps1 api-smoke    # Полный цикл: health, треки, auth, онбординг, рекомендации
```

**Терминал 3 — Фронтенд:**

```powershell
.\scripts\dev.ps1 frontend     # Vite dev server на http://localhost:5173
.\scripts\dev.ps1 frontend-build  # Production-сборка
```

---

## Команды dev.ps1

| Команда | Описание |
|---------|----------|
| `help` | Показать справку |
| `doctor` | Проверить Go, Node, npm и доступность PostgreSQL |
| `migrate` | Создать БД (если нет) и накатить миграции + seed |
| `reset-db` | Сбросить public schema → миграции → seed |
| `migration-status` | Показать статус миграций |
| `test` | Запустить backend-тесты |
| `backend` | Запустить Go API |
| `frontend` | Запустить Vite dev server |
| `frontend-build` | Собрать frontend в production |
| `api-smoke` | Сквозной smoke-тест API (весь user flow) |
| `check` | doctor + test + frontend-build + reset-db |

---

## API эндпоинты

### Публичные (без auth)

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/healthz` | Проверка здоровья |
| GET | `/api/v1/genres` | Список жанров |
| GET | `/api/v1/artists` | Список артистов |
| GET | `/api/v1/tracks` | Список треков |
| GET | `/api/v1/tracks/:id` | Детали трека |
| GET | `/api/v1/playlists` | Список плейлистов (публичные) |
| GET | `/api/v1/playlists/:id` | Детали плейлиста |
| POST | `/api/v1/auth/register` | Регистрация |
| POST | `/api/v1/auth/login` | Вход |
| GET | `/swagger/` | Swagger UI |

### Защищённые (требуют Bearer token)

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/v1/me` | Текущий пользователь |
| POST | `/api/v1/me/onboarding` | Сохранить предпочтения |
| GET | `/api/v1/me/profile` | Профиль пользователя |
| GET | `/api/v1/recommendations/tracks` | Рекомендации треков |
| GET | `/api/v1/recommendations/playlists` | Рекомендации плейлистов |
| POST | `/api/v1/events` | Событие (generic) |
| POST | `/api/v1/tracks/:id/play` | Прослушивание |
| POST | `/api/v1/tracks/:id/like` | Лайк |
| POST | `/api/v1/tracks/:id/dislike` | Дизлайк |
| POST | `/api/v1/tracks/:id/skip` | Пропуск |
| POST | `/api/v1/playlists/:id/open` | Открытие плейлиста |
| GET | `/api/v1/admin/stats` | Статистика системы |
| GET | `/api/v1/admin/recommendation-metrics` | Метрики рекомендаций |

---

## Демо-сценарий

1. Открыть http://localhost:5173
2. Зарегистрироваться (email + пароль)
3. Система автоматически входит и редиректит на главную
4. Нажать "Настроить рекомендации"
5. Пройти 4 шага: жанры → артисты → треки → контексты
6. После онбординга — персонализированные рекомендации на главной
7. Оценивать треки (♥ / ✕) — влияет на будущие рекомендации

Системный пользователь для тестирования:
- Email: `seed@music.local`
- Пароль: `plain-password`

---

## Структура проекта

```
├── backend/
│   ├── cmd/
│   │   ├── app/          # Точка входа API
│   │   └── migrate/      # Утилита миграций
│   ├── config/           # Конфигурация
│   ├── internal/
│   │   ├── app/          # DI и запуск
│   │   ├── auth/         # JWT
│   │   ├── controller/http/  # Хендлеры + роутер
│   │   ├── entity/       # Сущности
│   │   ├── repository/   # PostgreSQL
│   │   └── usecase/      # Бизнес-логика
│   ├── migrations/       # SQL-миграции (14 шт)
│   └── .env              # Конфиг
├── frontend/
│   └── src/
│       ├── api/          # HTTP-клиенты
│       ├── auth/         # AuthContext
│       ├── components/   # UI-компоненты
│       ├── config/       # Роуты, API endpoints
│       ├── context/      # PlayerContext
│       ├── data/         # Fallback-данные
│       ├── pages/        # Страницы
│       └── types/        # TypeScript-типы
├── docs/                 # Документация диплома
└── scripts/
    └── dev.ps1           # PowerShell-runner
```

---

## Troubleshooting

### Backend не может подключиться к БД

```powershell
.\scripts\dev.ps1 doctor
# Проверить: запущен ли PostgreSQL, совпадает ли пароль в backend/.env
```

### Ошибка "could not connect to server"

1. Убедись, что PostgreSQL запущен
2. Проверь `backend/.env` — строка `PG_URL` должна быть корректной
3. Если база `music_recommender` не создана — выполни `.\scripts\dev.ps1 migrate`

### Миграции падают

```powershell
.\scripts\dev.ps1 reset-db    # Полный сброс и применение с нуля
```

### Frontend не может достучаться до API

По умолчанию Vite proxy перенаправляет `/api` запросы на `http://127.0.0.1:8080` (см. `vite.config.ts`). Если backend на другом хосте/порту, задай переменную окружения:

```powershell
$env:VITE_API_BASE_URL="http://127.0.0.1:8080"
```

### Сборка frontend падает

```powershell
npm --prefix frontend install  # Обновить зависимости
.\scripts\dev.ps1 frontend-build
```
