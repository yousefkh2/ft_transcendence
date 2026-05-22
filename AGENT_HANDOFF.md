# Agent Handoff

## Project Context

Project: `ft_transcendence`

Current product direction:

- real-time multiplayer language-learning game
- first game mode: `Apartment Setup`
- two-player asymmetric co-op
- voice-first product direction, but voice is not part of Sprint 1
- 2D MVP, not Three.js
- fallback later: `Snake`, using the same session/realtime platform

Core idea:

Two players have different information and must speak German to arrange a room before time runs out.

## Current Team Shape

The team is now 3 people.

### Youssef

Role:

- Product Owner
- Game Designer
- Developer

Owns:

- product vision
- game rules
- Apartment Setup spec
- scenarios
- scoring
- feature priorities
- gameplay validation

Availability:

- high availability for the next two weeks
- uncertain availability after starting a job
- should avoid becoming the long-term implementation bottleneck

### Daniel

Role:

- Scrum Master
- DevOps
- QA
- Developer

Owns:

- meetings
- task board
- sprint planning
- blockers
- Docker setup
- evaluation checklist
- README structure
- test coordination

### Gabriel

Role:

- Technical Lead
- Backend Lead
- Developer

Owns:

- architecture quality
- Go backend
- database schema
- authentication
- API structure
- match history
- critical code review

Realtime and frontend ownership needs to be reassigned now that Mateo is out.
Until the team decides otherwise, keep the next frontend/realtime targets small
and integration-focused.

## Process Decision

Use Scrum-lite.

Do not use full corporate Scrum.
Do not use pure Kanban.
Do not use Waterfall.

Reason:

This project is uncertain, integration-heavy, and time-limited. The main risk is not lack of effort. The main risk is that everyone builds against different assumptions.

Rule:

```text
Every week, prove one thing works.
```

Sprint 1 target:

```text
Two browser tabs join the same room, receive different roles, and one game.object_moved event produces a state update in the other tab.
```

Sprint 2 learning slice:

```text
Two browser tabs receive role-specific UI, the backend starts the round when both players join, and both tabs display the same backend-owned timer and objective checklist.
```

## Current App State

The repo started as a docs-only workspace. A minimal Sprint 1 proof-of-concept has now been added.

Current runnable app:

- Docker Compose stack with frontend, backend, and PostgreSQL
- Vue/Vite frontend shell
- Go HTTP server
- static browser POC page
- minimal WebSocket endpoint
- room join
- role assignment
- hardcoded Apartment scenario
- `game.object_moved` events
- backend objective validation
- `game.state_patch` broadcast to both tabs
- `game.round_ended` after all four objectives are completed
- backend-owned round timer
- `game.timer_tick` broadcast every second after both players join
- role-specific frontend panels
- objective checklist that updates from backend state

This is not production architecture. It is a learning/proof step.

## Current Files of Interest

Main handoff:

- `AGENT_HANDOFF.md`

Specs/docs:

- `Apartment Setup v1 Spec.md`
- `Sprint 1 POC Notes.md`
- `Architecture and Module Plan.md`
- `transcendence.md`

Code:

- `docker-compose.yml`
- `.env.example`
- `Makefile`
- `go.mod`
- `backend/Dockerfile`
- `backend/cmd/api/main.go`
- `frontend/Dockerfile`
- `frontend/package.json`
- `frontend/sprint1-poc/index.html`
- `.gitignore`

## Current Run Command

Preferred team run command:

```sh
cp .env.example .env
make up
```

Then open:

```text
Frontend: http://localhost:5173
Backend POC: http://localhost:8080
Backend health: http://localhost:8080/health
Database reachability: http://localhost:8080/health/db
```

Local backend-only fallback:

```sh
PORT=8081 GOCACHE=/home/yousef/Desktop/TRANSENDENCE/.gocache go run ./backend/cmd/api
```

## Demo Script

1. Open `http://localhost:8080` in Tab A.
2. Open `http://localhost:8080` in Tab B.
3. Click `Join Room` in both tabs.
4. Confirm one tab receives `mission_control`.
5. Confirm the other tab receives `on_site`.
6. Confirm both tabs receive `game.round_started`.
7. Confirm both tabs show the same live timer.
8. Confirm the action buttons are enabled only on the `on_site` tab.
9. In the `on_site` tab, click:

```text
plant right_of sofa
```

10. Confirm both tabs receive a `game.state_patch`.
11. Confirm the objective checklist updates on both tabs.
12. Click the remaining valid move buttons.
13. Confirm both tabs receive `game.round_ended`.

## Verification Already Done

Go compile/test:

```sh
GOCACHE=/home/yousef/Desktop/TRANSENDENCE/.gocache go test ./...
```

Result:

```text
?    transcendence/backend/cmd/api    [no test files]
```

Health endpoint:

```sh
curl -s http://localhost:8080/health
```

Result:

```text
ok
```

## Important Technical Note

The first POC used handwritten WebSocket frame handling because package downloads were initially restricted.

That was useful as a before-picture, but it immediately produced protocol bugs:

- `8082` fixed the UI so only `on_site` can click action buttons.
- `8083` fixed one frame-length issue in the handwritten parser.
- A second disconnect showed the handwritten parser was still the wrong place to spend time.

The current POC now uses:

```text
github.com/gorilla/websocket v1.5.3
```

This is a better meeting lesson:

```text
We felt the pain of hand-rolling WebSocket protocol details, then replaced it with a maintained library.
```

The team can still choose a different Go WebSocket library later, but the POC should no longer use handwritten frame parsing.

## Sprint 2 State

Sprint 2 learning slice has been added.

New backend behavior:

- first player joins and waits
- second player joins and starts the round
- server broadcasts `game.round_started` to both players
- server starts a timer goroutine for the room
- server broadcasts `game.timer_tick` every second
- server rejects object moves before the round starts or after it ends
- server marks the round ended after all objectives or timer expiration

New frontend behavior:

- role panel explains Mission Control vs On Site
- action buttons unlock only for On Site after the round starts
- timer displays backend-owned time
- objective checklist updates when `game.state_patch` arrives

Core principle:

```text
The browser requests actions. The server owns game truth.
```

## MVP Architecture Position

For MVP:

- live round state can live in the Go session manager
- PostgreSQL stores durable data like users, matches, stats, and history
- Redis is not required for the first playable version

Redis can be added later if needed for:

- matchmaking queues
- player presence
- rate limiting
- shared session metadata
- scaling WebSockets across multiple backend instances

Current stance:

Design so Redis is easy to add later, but do not add it before the team has a concrete use case.

## Current Product Scope for Apartment v1

Objects:

```text
sofa
coffee_table
chair
plant
red_book
blue_cup
picture
lamp
```

Relations:

```text
on
under
left_of
right_of
above
behind
```

Hardcoded target conditions:

```json
[
  { "id": "objective_1", "objectId": "plant", "relation": "right_of", "targetId": "sofa" },
  { "id": "objective_2", "objectId": "red_book", "relation": "on", "targetId": "coffee_table" },
  { "id": "objective_3", "objectId": "blue_cup", "relation": "under", "targetId": "coffee_table" },
  { "id": "objective_4", "objectId": "picture", "relation": "above", "targetId": "sofa" }
]
```

Out of scope for v1:

- rotation
- physics
- inventory
- live grammar correction
- AI-generated scenarios
- advanced scoring
- spectator mode
- tournament mode
- 3D rendering
- Redis as a required dependency

## Recommended Next Steps

Immediate next steps before the next team meeting:

1. Send the repo after committing the canonical Docker setup.
2. Ask every teammate to run `cp .env.example .env && make up`.
3. Ask the team to confirm who owns realtime/frontend after Mateo's departure.
4. Youssef should test the two-tab POC manually.
5. Write down what feels unclear or awkward in the flow.
6. Bring the demo script to the meeting.
7. Ask the team to accept or adjust the Sprint 1/Sprint 2 success target.
8. Ask Gabriel/Daniel whether to keep Gorilla WebSocket for the first real backend package split.

Recommended implementation steps after team alignment:

1. Confirm whether to keep Gorilla WebSocket or choose another maintained library.
2. Move room/session code into proper packages:

```text
backend/internal/realtime
backend/internal/game
```

3. Replace action buttons with simple drag/drop zones in the Vue app.
4. Add PostgreSQL schema/migrations when auth/match history starts.
5. Add README details for contribution workflow and evaluation startup.

## What To Tell The Team

Short version:

```text
I built a tiny Sprint 1/Sprint 2 POC to feel the flow before we over-plan it.
It proves browser -> WebSocket -> backend session -> game validation -> timer/state broadcast -> both browsers.
It is intentionally not final architecture.
The repo now starts with Docker Compose, so everyone should be able to run the same frontend/backend/database stack.
Now we should decide who owns frontend/realtime and whether to keep Gorilla WebSocket.
```

## Last Known Server State

At the time this handoff was updated, Docker Compose had been started with:

```sh
make up
```

Do not assume it is still running in a later session. Re-run `make up` if needed.
