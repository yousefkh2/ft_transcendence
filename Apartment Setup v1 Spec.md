# Apartment Setup v1 Spec

## Purpose

This document defines the first playable version of `Apartment Setup`.

The goal of v1 is not to build the full language-learning product. The goal is to prove the core multiplayer loop:

Two players with different information must communicate in German to arrange a room before time runs out.

## Owner

Youssef owns:

- product vision
- game rules
- scenario design
- scoring rules
- gameplay validation
- feature priorities

Youssef should not be a blocker for day-to-day implementation once this v1 contract is accepted. If a detail is missing, the team should choose the simplest option that preserves this spec.

## Sprint 1 Success Target

Sprint 1 is successful when:

- two browser tabs can join the same room
- each tab receives a different role
- one tab can send a `game.object_moved` event
- the backend accepts the event
- the other tab receives a game state update

The first demo does not need polished UI, authentication, scoring, voice, or a full room layout.

## Core Game Loop

1. Two players join a room.
2. The system assigns roles.
3. The round starts with a fixed duration.
4. `Mission Control` sees the target room conditions.
5. `On Site` sees the messy room and can move objects.
6. Players speak German to solve the room.
7. `On Site` sends object movement events.
8. The backend validates semantic relations.
9. The round ends when all objectives are complete or the timer reaches zero.
10. The result screen shows success/failure and completed objectives.

## Roles

### Mission Control

Mission Control sees:

- target room conditions
- objective checklist
- timer

Mission Control cannot:

- move objects
- see the exact messy room state in v1

Mission Control's job:

- give German instructions
- answer clarification questions
- guide On Site toward the target state

Example phrases:

- `Stell die Pflanze rechts neben das Sofa.`
- `Leg das rote Buch auf den Couchtisch.`
- `Das Bild soll ueber dem Sofa haengen.`

### On Site

On Site sees:

- current room state
- draggable objects
- valid placement targets
- timer

On Site cannot:

- see the target conditions
- see Mission Control's checklist

On Site's job:

- move objects
- ask clarification questions in German
- confirm completed actions

Example phrases:

- `Links oder rechts vom Sofa?`
- `Meinst du das rote Buch?`
- `Die Pflanze steht jetzt neben dem Sofa.`

## Round Settings

For v1:

```json
{
  "durationSeconds": 180,
  "playerCount": 2,
  "gameMode": "apartment",
  "difficulty": "easy"
}
```

## Object List

Use exactly these objects for the first playable version:

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

These IDs should be stable across frontend, backend, scenario data, and tests.

## Allowed Relations

Use exactly these semantic relations for v1:

```text
on
under
left_of
right_of
above
behind
```

The game should validate relations semantically, not by exact pixels.

Example:

`plant right_of sofa` is valid if the plant is placed in a right-side zone near the sofa. The player should not need pixel-perfect placement.

## First Scenario

Use one hardcoded scenario first.

```json
{
  "id": "apartment_easy_001",
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
    {
      "id": "objective_1",
      "objectId": "plant",
      "relation": "right_of",
      "targetId": "sofa"
    },
    {
      "id": "objective_2",
      "objectId": "red_book",
      "relation": "on",
      "targetId": "coffee_table"
    },
    {
      "id": "objective_3",
      "objectId": "blue_cup",
      "relation": "under",
      "targetId": "coffee_table"
    },
    {
      "id": "objective_4",
      "objectId": "picture",
      "relation": "above",
      "targetId": "sofa"
    }
  ]
}
```

## WebSocket Event Contract Draft

### Join Room

```json
{
  "type": "room.join",
  "payload": {
    "roomCode": "ABCD"
  }
}
```

### Role Assigned

```json
{
  "type": "room.role_assigned",
  "payload": {
    "role": "mission_control",
    "playerId": "player_1"
  }
}
```

Valid role values:

```text
mission_control
on_site
```

### Round Started

```json
{
  "type": "game.round_started",
  "payload": {
    "matchId": "match_123",
    "scenarioId": "apartment_easy_001",
    "durationSeconds": 180,
    "role": "on_site"
  }
}
```

### Object Moved

Sent by `On Site`.

```json
{
  "type": "game.object_moved",
  "payload": {
    "matchId": "match_123",
    "objectId": "plant",
    "relation": "right_of",
    "targetId": "sofa"
  }
}
```

### State Patch

Sent by backend after validating a move.

```json
{
  "type": "game.state_patch",
  "payload": {
    "matchId": "match_123",
    "completedObjectives": ["objective_1"],
    "timeRemainingSeconds": 142
  }
}
```

### Round Ended

```json
{
  "type": "game.round_ended",
  "payload": {
    "matchId": "match_123",
    "result": "success",
    "completedObjectives": [
      "objective_1",
      "objective_2",
      "objective_3",
      "objective_4"
    ],
    "score": 1.0
  }
}
```

Valid result values:

```text
success
failure
```

## Scoring v1

Keep scoring simple:

```text
score = completedObjectives / totalObjectives
```

Examples:

```text
4 / 4 = 1.00
3 / 4 = 0.75
2 / 4 = 0.50
1 / 4 = 0.25
0 / 4 = 0.00
```

Win condition:

```text
All objectives completed before the timer reaches zero.
```

Failure condition:

```text
Timer reaches zero before all objectives are completed.
```

No time bonus in v1. Add it later only after the basic loop works.

## Validation Rule

The backend should treat the submitted relation as the gameplay truth for the first technical prototype.

For the first polished playable version, the frontend should map drag/drop zones to semantic relations.

Example:

- dropping `plant` into the right-side zone of `sofa` sends `plant right_of sofa`
- dropping `red_book` onto `coffee_table` sends `red_book on coffee_table`

The backend then checks whether the submitted relation matches one of the target conditions.

## Data Ownership

### Backend Owns

- room membership
- role assignment
- match/session lifecycle
- timer authority
- objective validation
- completed objectives
- final result
- saved match history

### Frontend Owns

- visual room layout
- drag/drop interaction
- mapping drop zones to semantic relations
- role-specific screen display
- result/review screens

### Game Design Owns

- object list
- allowed relations
- scenario data
- target conditions
- scoring rules
- language focus

## MVP Architecture Position

For the MVP:

- live round state can live in the Go session manager
- PostgreSQL stores durable data like users, matches, stats, and history
- Redis is not required for the first playable version

Design should leave space for Redis later if needed for:

- matchmaking queues
- player presence
- rate limiting
- shared session metadata
- scaling WebSockets across multiple backend instances

## Out of Scope for v1

Do not build these in v1:

- rotation
- physics
- inventory
- object damage
- penalties for touching wrong objects
- live grammar correction
- AI-generated scenarios
- advanced scoring
- time bonus
- hints
- spectator mode
- tournament mode
- 3D rendering
- Redis as a required dependency

## Open Questions for Team

These should be answered in the next meeting:

1. Which Go WebSocket library will Mateo use?
2. Which Go backend framework will Gabriel use?
3. Does the backend own the countdown timer from Sprint 1, or only after the WebSocket proof works?
4. Should `roomCode` be human-readable, like `ABCD`, or UUID-based for the first version?
5. Should role assignment be random in v1, or should players choose roles?
6. What is the minimum Docker Compose skeleton Daniel wants for Sprint 1?
7. What does Mateo need from this spec to create the first room UI and drag/drop prototype?

## Meeting Decision Needed

The team should accept or adjust this sentence:

```text
Sprint 1 is successful when two browser tabs can join the same room, receive different roles, and send one game.object_moved event that produces a state update in the other tab.
```

Once this is accepted, each owner can start without waiting for the full product to be designed.

## Team Process: Why Scrum-Lite

The team should use Scrum-lite for this project.

This does not mean heavy process, long meetings, or ceremony for its own sake. It means using short sprints to force integration and visible progress.

The main project risk is not that nobody works. The main risk is that everyone works on different assumptions:

- backend builds one session model
- realtime builds a different room flow
- frontend expects a different event contract
- game design asks for rules that the technical flow does not support

For a real-time multiplayer game, isolated progress is dangerous. The important proof is whether the pieces connect:

```text
frontend -> websocket -> backend -> game session -> state update -> frontend
```

Scrum-lite helps because every sprint has one concrete demo target.

Example:

```text
Sprint 1:
Two tabs join one room and exchange one game event.

Sprint 2:
Roles, timer, and hardcoded Apartment scenario work.

Sprint 3:
Drag/drop UI sends semantic object events.

Sprint 4:
Match result is saved to PostgreSQL and shown in history.
```

The team should keep the process small:

- one weekly planning meeting
- one midweek async check-in
- one end-of-week demo/review
- one simple task board
- one clear owner per task
- no daily standup unless people are blocked

The purpose is simple:

```text
Turn vague ambition into weekly proof.
```

## Updated 4-Person Ownership

With a 4-person team, ownership should be clear and lightweight.

### Youssef: Product Owner / Game Designer / Developer

Owns:

- product vision
- game rules
- Apartment Setup spec
- scenarios
- scoring
- feature priorities
- gameplay validation

Youssef should focus on keeping the game clear and unblocked. Since Youssef is starting a job soon, he should avoid becoming the main implementation bottleneck.

### Daniel: Scrum Master / DevOps / QA / Developer

Owns:

- meetings
- task board
- sprint planning
- blockers
- Docker setup
- evaluation checklist
- README structure
- test coordination

Daniel's job is to keep the team moving and make blockers visible early.

### Gabriel: Technical Lead / Backend Lead / Developer

Owns:

- architecture quality
- Go backend
- database schema
- authentication
- API structure
- match history
- critical code review

Gabriel's job is to keep the backend coherent and prevent technical drift.

### Mateo: Realtime + Frontend Lead / Developer

Owns:

- WebSockets
- lobby rooms
- matchmaking/session lifecycle
- role assignment
- live game state sync
- Vue app structure
- game UI
- drag/drop interaction
- responsive layout
- review/results screens

Mateo's area is broad, so the first frontend and realtime targets should stay intentionally simple.

## Final Call on Process

Use:

```text
Scrum-lite
```

Do not use full corporate Scrum.

Do not use pure Kanban.

Do not use Waterfall.

The project is uncertain, integration-heavy, and time-limited. Scrum-lite gives the team rhythm and accountability without adding unnecessary ceremony.

The rule is:

```text
Every week, prove one thing works.
```

If the team cannot demo it, the team does not really have it.

## Youssef Availability Note

Youssef has high availability for the next two weeks.

After that, Youssef starts a job and his availability is uncertain. The team should plan around this honestly instead of assuming the same contribution level across the whole project.

This means:

- use the next two weeks to lock product direction
- finalize the Apartment Setup v1 rules early
- define event contracts before implementation spreads
- remove game-design ambiguity while Youssef is highly available
- avoid making Youssef the long-term implementation bottleneck

Recommended approach:

```text
Weeks 1-2:
Youssef moves fast on product/game design, scenario specs, acceptance criteria, and validation.

After Week 2:
Youssef remains Product Owner/Game Designer, but the team should rely on written specs and clear priorities rather than real-time availability.
```

Youssef can still contribute as a developer, but his core responsibility should be keeping the game concept coherent and the scope controlled.
