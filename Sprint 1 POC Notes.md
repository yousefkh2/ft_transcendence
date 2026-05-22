# Sprint 1 POC Notes

## What Exists

This is a minimal proof for the first sprint target.

It includes:

- a Go HTTP server
- a minimal WebSocket endpoint at `/ws`
- a static browser test page
- room joining
- role assignment
- hardcoded `Apartment Setup` scenario data
- `game.object_moved` events
- backend objective validation
- `game.state_patch` broadcasts
- `game.round_ended` when all four objectives are complete
- backend-owned round timer
- `game.timer_tick` broadcasts
- role-specific frontend panels
- objective checklist display

This is not the final architecture. It is a learning/proof step.

## Run Command

```sh
GOCACHE=/home/yousef/Desktop/TRANSENDENCE/.gocache go run ./backend/cmd/api
```

If port `8080` is already busy:

```sh
PORT=8081 GOCACHE=/home/yousef/Desktop/TRANSENDENCE/.gocache go run ./backend/cmd/api
```

If `8081` is also occupied by an older run:

```sh
PORT=8082 GOCACHE=/home/yousef/Desktop/TRANSENDENCE/.gocache go run ./backend/cmd/api
```

If testing the build with extended WebSocket frame parsing fixed:

```sh
PORT=8083 GOCACHE=/home/yousef/Desktop/TRANSENDENCE/.gocache go run ./backend/cmd/api
```

Current recommended run command:

```sh
cp .env.example .env
make up
```

This version uses `github.com/gorilla/websocket` instead of the handwritten WebSocket parser.

It also includes the Sprint 2 learning slice:

- round starts only after the second player joins
- backend owns the timer
- backend broadcasts `game.timer_tick`
- frontend displays role-specific panels
- frontend displays an objective checklist instead of only raw JSON

Then open:

```text
Frontend shell: http://localhost:5173
Backend POC: http://localhost:8080
```

## Demo Script

1. Open `http://localhost:8080` in Tab A.
2. Open `http://localhost:8080` in Tab B.
3. Click `Join Room` in both tabs.
4. Confirm one tab receives `mission_control`.
5. Confirm the other tab receives `on_site`.
6. Confirm both tabs receive `game.round_started`.
7. Confirm both tabs show the same live timer.
8. The action buttons should only be enabled in the `on_site` tab.
9. In the `on_site` tab, click:

```text
plant right_of sofa
```

10. Confirm both tabs receive a `game.state_patch` with:

```json
{
  "completedObjectives": ["objective_1"]
}
```

11. Confirm the checklist updates on both tabs.
12. Click the remaining valid moves from the `on_site` tab.
13. Confirm both tabs receive `game.round_ended`.

## Sprint 1 Learning Goal

The point is to feel the real-time path:

```text
browser -> websocket -> backend room -> game validation -> broadcast -> other browser
```

This gives the team a concrete before-picture before adding real frontend structure, database persistence, auth, Docker, Redis, or a maintained WebSocket package.

## Sprint 2 Learning Goal

Sprint 2 turns the raw event path into a tiny game-state screen.

The main principle:

```text
The server owns important game state. The browser displays it.
```

The backend now decides:

- when the round starts
- how much time remains
- whether an object move completes an objective
- when the round ends

The frontend now displays:

- role-specific instructions
- live timer
- objective checklist
- disabled actions after objectives are complete

This is still not a final Vue game UI. It is the next proof layer between raw JSON and real gameplay.

## Known Shortcuts

The current POC intentionally cuts corners:

- no Vue yet
- no PostgreSQL yet
- no auth yet
- no Docker yet
- no real drag/drop yet
- no reconnect handling
- no persistence after server restart

Some previous shortcuts are now resolved:

- WebSocket handling uses Gorilla WebSocket
- backend-owned countdown timer exists
- role-specific UI panels exist
- objective checklist exists

## Debug Note

An early version disconnected the `on_site` tab after the first action because the handwritten WebSocket parser rejected frames using the 64-bit extended payload length format.

A later test still disconnected after a second action, which showed the handwritten parser was not worth debugging further.

The POC now uses `github.com/gorilla/websocket` for upgrade, read, and write handling. This is closer to the real implementation direction and keeps the Sprint 1 proof focused on game/session flow instead of WebSocket protocol details.

## Replace Later

Before this becomes production project code, the team should move the current single-file backend into proper packages.

The real decision for the backend/realtime owner is now:

```text
Do we keep Gorilla WebSocket or choose a different maintained Go WebSocket library?
```

## Next Implementation Steps

Recommended order:

1. Decide the final Go WebSocket library.
2. Move room/session code into `backend/internal/realtime` and `backend/internal/game`.
3. Add a Vue frontend scaffold.
4. Replace action buttons with simple drag/drop zones.
5. Add PostgreSQL only when match history/auth work starts.
6. Add Docker Compose once the frontend/backend structure stabilizes.
