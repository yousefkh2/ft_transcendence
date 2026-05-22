# ft_transcendence

Real-time multiplayer language-learning game. The current first mode is
`Apartment Setup`: two players receive different information and must speak
German to arrange a room before time runs out.

## Stack

- Frontend: Vue + Vite
- Backend: Go
- Realtime: WebSockets
- Database: PostgreSQL
- Local runtime: Docker Compose

## First Run

```sh
cp .env.example .env
make up
```

Then open:

```text
Frontend: http://localhost:5173
Backend root: http://localhost:8080
Backend health: http://localhost:8080/health
Database reachability: http://localhost:8080/health/db
```

## Useful Commands

```sh
make up       # build and start all services
make down     # stop services
make logs     # follow logs
make ps       # show service status
make test     # run Go tests locally
make health   # check backend and database reachability
```

## Current Sprint Target

The infrastructure target is one-command startup for the team:

- Vue frontend starts from Docker Compose.
- Go backend starts from Docker Compose.
- PostgreSQL starts with a persistent named volume.
- Backend `/health` returns `ok`.
- Backend `/health/db` confirms the database is reachable.

The backend is intentionally minimal so the backend owner can build the real API
from a clean, explainable starting point.

## Database Schema

The draft schema lives at:

```text
backend/sql/migrations/001_initial_schema.sql
```

It is not auto-loaded by Docker Compose. Postgres init scripts only run on an
empty volume, which makes schema changes easy to miss during team development.
Treat the SQL file as reviewable migration input until the backend has a real
migration command.
