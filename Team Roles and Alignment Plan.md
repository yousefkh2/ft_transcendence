# Team Roles and Alignment Plan

## Context

Project: `ft_transcendence`

Direction: real-time multiplayer language-learning game.

Stack:

- Backend: Go
- Frontend: Vue.js
- Database: PostgreSQL
- Real-time: WebSockets
- Deployment: Docker Compose

First game mode: `Apartment Setup`

Fallback game mode: `Snake`, using the same auth, lobby, WebSocket, match history, and stats architecture.

## Team

| Person | Primary Role | Main Ownership |
|---|---|---|
| Youssef | Product Owner / Game Design Spec | Game concept, rules, scenario data, scoring, language-learning loop, evaluator story |
| Daniel | Project Manager / QA / DevOps | planning, task board, deadlines, testing checklist, Docker/run instructions, evaluation readiness |
| Gabriel | Backend Lead | Go API, auth, users, database, ORM/migrations, profiles, match history, stats |
| Stelio | Realtime / Game Server Lead | WebSocket hub, lobby rooms, matchmaking, game sessions, role assignment, state sync |
| Matteo | Frontend Lead | Vue app, UI structure, game screen, drag/drop room, responsive layout, review/results screens |

Everyone contributes code, testing, review, and documentation. The ownership above means each area has a final decision-maker and a person responsible for keeping it unblocked.

## Product Decision

The MVP should not start with Three.js.

The first version should be a 2D game because the main risk is not graphics. The main risk is whether the game loop creates useful spoken German between two players.

The game should be:

- real-time
- two-player
- asymmetric
- voice-first
- simple to manipulate
- semantically validated
- reviewed after the round

## First Game: Apartment Setup

### Core Concept

Two players must arrange a room before time runs out.

Player A: `Mission Control`

- sees the target room layout
- sees the required checklist
- cannot move objects
- gives German instructions

Player B: `On Site`

- sees the messy room
- can move objects
- cannot see the target layout
- asks clarification questions and follows instructions

### Example Round

Mission Control sees:

- The plant is next to the sofa.
- The red book is on the coffee table.
- The picture is above the sofa.
- The blue cup is under the table.

On Site sees:

- messy room
- draggable objects
- valid placement zones

Expected speech:

- `Stell die Pflanze neben das Sofa.`
- `Links oder rechts?`
- `Rechts vom Sofa.`
- `Leg das rote Buch auf den Couchtisch.`
- `Nein, das andere Buch.`
- `Ja, genau.`

## MVP Game Mechanics

Version 1 should be intentionally small.

Scope:

- one room
- 6 to 8 movable objects
- 4 target conditions
- 3-minute timer
- drag and drop only
- no rotation
- no physics
- no inventory
- no live grammar correction

Core objects:

- sofa
- coffee table
- chair
- plant
- red book
- blue cup
- picture
- lamp

Core relations:

- `on`
- `under`
- `left_of`
- `right_of`
- `above`
- `behind`

The game should validate semantic relations, not exact pixels.

Example target condition:

```json
{
  "object": "plant",
  "relation": "right_of",
  "target": "sofa"
}
```

## Architecture Rule

The game must be modular.

The platform should provide:

- auth
- users
- lobby
- WebSocket transport
- game sessions
- match history
- stats
- achievements
- Docker deployment

Each game mode should provide only:

- initial state
- role assignment
- valid events
- state update rules
- win/loss validation
- final scoring

This lets the team keep the same platform if `Apartment Setup` needs to be simplified or replaced by `Snake`.

## Transcendence Points Plan

Target 14 points first.

| Module | Type | Points | Owner |
|---|---:|---:|---|
| Use a framework for both frontend and backend | Major | 2 | Gabriel + Matteo |
| Real-time features using WebSockets | Major | 2 | Stelio |
| User interaction system | Major | 2 | Gabriel + Matteo |
| Standard user management and authentication | Major | 2 | Gabriel |
| Web-based game where users play against each other | Major | 2 | Youssef + Stelio + Matteo |
| Remote players, real-time | Major | 2 | Stelio |
| Game statistics and match history | Minor | 1 | Gabriel |
| Gamification system | Minor | 1 | Youssef + Matteo |
| **Total** |  | **14** |  |

Bonus/stretch modules should only be considered after these are stable.

Best stretch candidates:

- voice/speech integration
- tournament system
- spectator mode
- custom design system
- monitoring with Prometheus/Grafana

## Suggested Repository Structure

```text
.
  backend/
    cmd/api/
    internal/
      auth/
      users/
      friends/
      chat/
      lobby/
      realtime/
      game/
        modes/
          apartment/
          snake/
      stats/
      review/
      db/
      config/
  frontend/
    src/
      api/
      stores/
      components/
      games/
        apartment/
        snake/
      pages/
  docs/
  docker-compose.yml
  README.md
  .env.example
```

## Meetup Plan

Date: 2026-04-27

Goal: align the team before Youssef travels on 2026-04-28.

### 1. Confirm Product Direction

Decision to confirm:

- We are building a real-time multiplayer language-learning game.
- First game is `Apartment Setup`.
- Graphics are 2D for MVP.
- `Snake` is the fallback game mode.
- The architecture must support multiple game modes.

### 2. Confirm Team Ownership

Each person should accept one primary area:

- Youssef: game design spec
- Daniel: PM / QA / DevOps
- Gabriel: backend/API/database
- Stelio: WebSockets/realtime/game sessions
- Matteo: frontend/game UI

### 3. Confirm First Sprint Deliverables

Each owner should leave the meeting with a concrete first task.

| Person | First Deliverable |
|---|---|
| Youssef | `Apartment Setup v1` game spec with objects, relations, targets, scoring, and example scenarios |
| Daniel | task board, team calendar, Docker/evaluation checklist, README skeleton |
| Gabriel | Go backend skeleton, database schema draft, auth route plan |
| Stelio | WebSocket proof of concept: two browsers join one room and exchange events |
| Matteo | Vue app shell and static apartment room prototype with draggable objects |

### 4. Define Integration Contract

The frontend and realtime backend should agree early on WebSocket message shapes.

Example:

```json
{
  "type": "game.object_moved",
  "matchId": "match_123",
  "payload": {
    "objectId": "plant",
    "relation": "right_of",
    "targetId": "sofa"
  }
}
```

### 5. Define Done for Week 1

By the end of the first week, the team should aim for:

- Docker Compose starts frontend, backend, and database
- user can open Vue app
- Go API has health route
- two browser tabs can join the same room
- frontend can send a movement event over WebSocket
- backend can broadcast the event
- static apartment room exists
- first scenario JSON exists

## Youssef Pre-Travel Handoff

Before travel, Youssef should unblock the team by delivering:

1. Final `Apartment Setup v1` rules.
2. Object list and relation list.
3. 3 example scenarios as JSON.
4. Scoring formula.
5. Win/loss conditions.
6. First version of post-round language review requirements.
7. Clear answer on what is out of scope.

This is the most valuable contribution before travel because it lets frontend, backend, and realtime development move independently.

## Youssef Travel Window

Expected travel starts: 2026-04-28

Expected away period: about two weeks.

During this time:

- Daniel owns weekly coordination.
- Gabriel and Stelio align on API/WebSocket contracts.
- Matteo builds the apartment UI against mock data first.
- Youssef reviews async and answers game-design questions when possible.

The team should avoid blocking on perfect game design. If a rule is unclear, use the simplest version and document the assumption.

## Immediate Next Step for Youssef

Create the `Apartment Setup v1` spec with:

- object list
- valid placement zones
- target condition format
- 3 sample scenarios
- scoring
- win/loss rules
- player views
- out-of-scope list

This file should become the contract between game design, frontend, realtime, and backend.
