# Notion Workspace Blueprint

## Goal

The Notion workspace should make the project easy to coordinate, not become another project to maintain.

The workspace has three jobs:

1. Make ownership clear.
2. Make blockers visible early.
3. Keep decisions and resources findable.

Do not create too many pages. Start with a small system that the whole team can actually use.

## Top-Level Workspace

Create one main Notion page:

```text
ft_transcendence - Team Hub
```

Inside it, create these pages:

```text
1. Project Dashboard
2. Task Board
3. Sprint Plan
4. Team Ownership
5. Architecture
6. Game Design Spec
7. API + WebSocket Contracts
8. Decisions Log
9. Blockers
10. Resources Dump
11. Evaluation Checklist
```

If the team only uses three pages every day, they should be:

- `Project Dashboard`
- `Task Board`
- `Blockers`

## 1. Project Dashboard

Purpose: the landing page everyone checks first.

Recommended sections:

```text
Current Sprint Goal
This Week's Priorities
Important Links
Active Blockers
Team Ownership
Next Meeting
14-Point Module Progress
```

### Current Sprint Goal

Example:

```text
Sprint 1 Goal:
Get the project skeleton running and prove that two browser tabs can join a room and exchange one game event.
```

### This Week's Priorities

Keep this to 3-5 items.

Example:

```text
- Docker Compose starts frontend, backend, and database.
- Go API exposes /health.
- Vue app has login/lobby/game placeholder pages.
- WebSocket POC: two tabs join the same room.
- Apartment Setup v1 spec is finalized.
```

### Important Links

Add links to:

- GitHub repository
- Figma or design reference, if any
- Architecture page
- Game Design Spec
- API + WebSocket Contracts
- Evaluation Checklist

## 2. Task Board

This should be the main Kanban board.

Create a Notion database called:

```text
Tasks
```

### Properties

Use these properties:

| Property | Type | Purpose |
|---|---|---|
| Task | Title | Short task name |
| Status | Select | Backlog, Ready, In Progress, Blocked, Review, Done |
| Owner | Person | One clear owner |
| Area | Select | Game Design, Frontend, Backend, Realtime, DevOps, QA, Docs |
| Priority | Select | P0, P1, P2 |
| Sprint | Select | Sprint 1, Sprint 2, Sprint 3... |
| Depends On | Relation | Link to blocking task |
| Blocks | Relation | Link to tasks blocked by this task |
| Due Date | Date | Only for real deadlines |
| Acceptance Criteria | Text | What must be true for Done |
| PR / Link | URL | GitHub PR, doc, design, etc. |

### Status Definitions

```text
Backlog:
Idea exists, not ready to work on.

Ready:
Task is clear enough to start.

In Progress:
Someone is actively working on it.

Blocked:
Cannot continue until another decision/task is done.

Review:
Work is done but needs review/testing.

Done:
Merged, tested, and documented if needed.
```

### Priority Definitions

```text
P0:
Blocks other people or required for evaluation.

P1:
Important for MVP but not blocking today.

P2:
Nice to have.
```

### Required Views

Create these views:

```text
Board by Status
Board by Area
My Tasks
Blocked Tasks
Sprint 1
Review Queue
```

`Blocked Tasks` is the most important view after the main board.

## 3. Sprint Plan

Purpose: avoid vague weekly work.

Create one page per sprint:

```text
Sprint 1 - Foundation
Sprint 2 - Multiplayer Loop
Sprint 3 - Playable Apartment
Sprint 4 - Auth + Stats + Polish
```

Each sprint page should contain:

```text
Sprint Goal
Must Finish
Nice To Have
Risks
Demo Target
Retrospective Notes
```

### Sprint 1 Example

```text
Sprint Goal:
Create the runnable skeleton and prove the real-time path.

Must Finish:
- Docker Compose skeleton.
- Go /health route.
- Vue app shell.
- WebSocket room POC.
- Apartment Setup v1 contract.

Demo Target:
Two browser tabs join one room. One tab sends "plant moved right_of sofa"; the other tab receives it.
```

## 4. Team Ownership

Purpose: prevent confusion about who owns what.

Use this table:

| Person | Primary Ownership | Secondary Ownership | Blocking Risk |
|---|---|---|---|
| Youssef | Game design spec, scenario rules, product direction | WebSocket/game event contract, playtesting | If game rules are unclear, frontend/realtime cannot finalize event payloads |
| Daniel | PM, QA, DevOps, evaluation readiness | Docker, README, task hygiene | If Docker/evaluation checklist is late, integration becomes chaotic |
| Gabriel | Backend API, auth, DB, stats | ORM/migrations, API docs | If auth/schema is late, frontend and stats cannot integrate |
| Stelio | WebSockets, lobby, sessions, game state sync | match lifecycle, spectator later | If realtime contract is late, the game cannot become multiplayer |
| Matteo | Vue frontend, game UI, drag/drop | design system, responsive polish | If UI prototype is late, game feel cannot be tested |

## 5. Architecture

Purpose: one canonical architecture reference.

Add:

- the architecture diagram image
- simplified MVP architecture
- explanation that backend is a modular monolith
- current tech decisions
- future optional pieces

### MVP Architecture

```text
Vue Frontend
  -> REST API for auth/profile/history
  -> WebSocket for lobby/gameplay

Go Backend
  -> internal modules
  -> game session manager
  -> pluggable game modes

PostgreSQL
  -> users
  -> matches
  -> stats
  -> reviews
```

### Optional Later

```text
Redis:
rate limits, temporary shared session state, WebSocket scaling

Object storage:
avatars, temporary audio files, uploaded files

AI API:
post-round transcription and review
```

## 6. Game Design Spec

Purpose: Youssef's main ownership page.

This page should become the source of truth for game rules.

Recommended sections:

```text
Game Summary
Player Roles
Round Flow
Objects
Relations
Valid Move Event
Target Conditions
Win/Loss Conditions
Scoring
Difficulty Levels
Sample Scenarios
Out of Scope
Open Questions
```

### Minimum Contract

This part must be decided early:

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

### First Scenario Data

```json
{
  "id": "apartment_a1_easy_001",
  "durationSeconds": 180,
  "objects": [
    "sofa",
    "coffee_table",
    "chair",
    "plant",
    "red_book",
    "blue_cup",
    "picture",
    "lamp"
  ],
  "targetConditions": [
    { "objectId": "plant", "relation": "right_of", "targetId": "sofa" },
    { "objectId": "red_book", "relation": "on", "targetId": "coffee_table" },
    { "objectId": "blue_cup", "relation": "under", "targetId": "coffee_table" },
    { "objectId": "picture", "relation": "above", "targetId": "sofa" }
  ]
}
```

## 7. API + WebSocket Contracts

Purpose: prevent frontend/realtime/backend mismatch.

Create two sections:

```text
REST API
WebSocket Messages
```

### REST API Examples

```text
POST /api/auth/register
POST /api/auth/login
GET  /api/me
GET  /api/matches
GET  /api/matches/:id
GET  /api/stats/me
```

### WebSocket Message Examples

```json
{
  "type": "room.join",
  "payload": {
    "roomCode": "ABCD"
  }
}
```

```json
{
  "type": "game.ready",
  "payload": {
    "ready": true
  }
}
```

```json
{
  "type": "game.object_moved",
  "payload": {
    "objectId": "plant",
    "relation": "right_of",
    "targetId": "sofa"
  }
}
```

```json
{
  "type": "game.state_patch",
  "payload": {
    "completedObjectives": ["objective_1"],
    "timeRemainingSeconds": 142
  }
}
```

## 8. Decisions Log

Purpose: stop the team from reopening the same decisions.

Create a database:

```text
Decisions
```

Properties:

| Property | Type |
|---|---|
| Decision | Title |
| Date | Date |
| Owner | Person |
| Status | Select: Proposed, Accepted, Rejected, Revisit Later |
| Context | Text |
| Decision | Text |
| Consequences | Text |

Example:

```text
Decision:
Use 2D Vue UI for MVP instead of Three.js.

Context:
The main risk is the language co-op loop, not graphics.

Consequences:
Lower rendering complexity. Three.js can remain a future stretch module.
```

## 9. Blockers

Purpose: make dependency problems impossible to miss.

Create a database:

```text
Blockers
```

Properties:

| Property | Type |
|---|---|
| Blocker | Title |
| Severity | Select: Critical, High, Medium, Low |
| Blocked Person | Person |
| Blocking Area | Select |
| Needed From | Person |
| Status | Select: Open, Being Resolved, Resolved |
| Due By | Date |
| Resolution | Text |

### Blocker Rule

If a task is blocked for more than one day, it must be listed here.

### Examples

```text
Blocker:
Frontend does not know final move event format.

Needed From:
Youssef + Stelio

Resolution:
Agree on objectId + relation + targetId message format.
```

```text
Blocker:
Backend schema for match history is not ready.

Needed From:
Gabriel

Resolution:
Create draft matches and match_players tables.
```

## 10. Resources Dump

Purpose: useful links without cluttering the core pages.

Create a database:

```text
Resources
```

Properties:

| Property | Type |
|---|---|
| Resource | Title |
| Type | Select: Article, Video, Docs, Example Repo, Tool, Book, Prompt |
| Area | Select: Go, Vue, WebSockets, Game Design, Docker, Security, AI, German |
| Added By | Person |
| Link | URL |
| Summary | Text |
| Useful For | Text |

Examples:

- 100 Go Mistakes
- Vue official docs
- Gorilla/WebSocket docs or chosen Go WebSocket library docs
- Docker Compose docs
- PostgreSQL schema design notes
- German prepositions reference
- Web game UX examples

Rule:

Every resource needs a one-sentence summary. Otherwise the dump becomes useless.

## 11. Evaluation Checklist

Purpose: make sure you pass the project requirements.

Sections:

```text
Mandatory Requirements
14-Point Modules
Security
Accessibility
Browser Testing
Docker Setup
README Requirements
Demo Script
Known Risks
```

### Demo Script

Write the evaluator demo before the app is finished.

Example:

```text
1. Start the project with Docker Compose.
2. Register two users.
3. User A creates a room.
4. User B joins the room.
5. Users start Apartment Setup.
6. Mission Control sees target checklist.
7. On Site moves objects.
8. WebSocket state updates are visible.
9. Round ends.
10. Match history and stats are shown.
```

## Blocking Factor Map

Use this to understand who can block whom.

```text
Youssef blocks:
- Matteo, if game objects/relations are undefined.
- Stelio, if move events and win rules are undefined.
- Gabriel, if match result/scoring data is undefined.

Daniel blocks:
- everyone, if project board, priorities, and deadlines are unclear.
- evaluation, if Docker/README/checklist are late.

Gabriel blocks:
- Matteo, if auth/profile/history APIs are unavailable.
- Daniel, if database/setup instructions are unclear.
- stats/gamification, if match schema is late.

Stelio blocks:
- Matteo, if WebSocket messages are unstable.
- gameplay demo, if room/session lifecycle is late.
- Gabriel, if match lifecycle events are unclear.

Matteo blocks:
- game testing, if the room UI is unavailable.
- Youssef, if game feel cannot be tested.
- demo, if pages are not connected.
```

## First Week Plan

### Youssef

Own:

- Game Design Spec
- first scenario data
- scoring rules
- player views
- out-of-scope list

Deliverable:

```text
Apartment Setup v1 spec ready for frontend/realtime/backend.
```

### Daniel

Own:

- Notion setup
- task board
- README skeleton
- Docker/evaluation checklist
- meeting notes

Deliverable:

```text
Everyone has tasks, the repo has a README skeleton, and evaluation requirements are tracked.
```

### Gabriel

Own:

- Go backend skeleton
- health route
- database schema draft
- auth route plan

Deliverable:

```text
Backend runs locally and exposes /health.
```

### Stelio

Own:

- WebSocket server POC
- room join/leave
- broadcast event
- session lifecycle draft

Deliverable:

```text
Two browser tabs can join one room and exchange one event.
```

### Matteo

Own:

- Vue app shell
- route structure
- lobby placeholder
- static apartment prototype
- draggable object experiment

Deliverable:

```text
Vue app shows a basic apartment room and one draggable object.
```

## Meeting Rhythm

Keep meetings short.

### Weekly Planning

Length: 45 minutes

Agenda:

```text
1. What did we finish?
2. What is blocked?
3. What is the sprint goal?
4. Who owns each P0/P1 task?
5. What is the next demo target?
```

### Twice-Weekly Check-In

Length: 15 minutes

Agenda:

```text
1. Any blocker?
2. Any contract changed?
3. Any task needs review?
```

### Rule

If a decision affects more than one person, write it in the Decisions Log.

## Main Advice

Do not let Notion become a place where ideas go to die.

Every important item should become one of:

- a task
- a blocker
- a decision
- a resource
- a spec update

If it is none of those, it probably does not belong in the project workspace.
