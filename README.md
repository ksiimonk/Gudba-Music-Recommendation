# Gudba-Music-Recommendation

## Local runner

Use the PowerShell runner from the repository root:

```powershell
.\scripts\dev.ps1 help
```

Common flow:

```powershell
.\scripts\dev.ps1 doctor
.\scripts\dev.ps1 reset-db
.\scripts\dev.ps1 test
.\scripts\dev.ps1 frontend-build
.\scripts\dev.ps1 backend
```

In another terminal, with the backend still running:

```powershell
.\scripts\dev.ps1 api-smoke
```

Frontend:

```powershell
.\scripts\dev.ps1 frontend
```

## What the commands do

- `doctor` checks Go, Node, npm, and whether PostgreSQL from `backend/.env` is reachable.
- `migrate` creates the configured database if possible, applies `backend/migrations/*.up.sql`, and loads demo seed data.
- `reset-db` drops the `public` schema, recreates it, then applies all migrations and seed data from scratch.
- `migration-status` shows applied and pending migrations.
- `test` runs backend tests without scanning local cache folders.
- `backend` starts the Go API.
- `frontend` starts Vite.
- `frontend-build` builds the frontend.
- `api-smoke` checks `/healthz`, `/api/v1/tracks`, `/api/v1/tracks/1`, `/api/v1/playlists`, and `/api/v1/playlists/1`.

The backend uses `backend/.env`. By default it expects PostgreSQL at `127.0.0.1:5432` and database `music_recommender`.

Use `reset-db` when you changed migrations or seed data and want to verify the project on a clean database.
