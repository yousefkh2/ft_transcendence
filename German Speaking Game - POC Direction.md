---
created: 2026-03-30
tags:
  - project
  - language-learning
  - game-design
  - poc
related:
  - "[[Clippings/German Learning App (Game)|German Learning App (Game)]]"
  - "[[Speaking Game - Scenario Families]]"
  - "[[Speaking Game - Hidden Prompt System]]"
---

# German Speaking Game - POC Direction

Related: [[Clippings/German Learning App (Game)|German Learning App (Game)]]

See also: [[Speaking Game - Scenario Families]] [[Speaking Game - Better Scenario Families]]

See also: [[Speaking Game - Hidden Prompt System]]

## Why This Idea Is Strong

The core problem is real:

- Speaking rooms often stall because nobody wants to keep organizing the topic flow.
- Learners need structure that gets them speaking immediately.
- The app should remove topic-planning friction and replace it with a cooperative objective.

The strongest version of the idea is not a single bomb-defusal clone. It is a **voice-first cooperative scenario engine** for speaking practice.

## Core Product Principle

The core question is:

**Can a structured co-op scenario make people speak more German, with more focus, and want another round?**

Everything should be judged by:

**Does this increase meaningful spoken German without hurting flow?**

## Main Product Direction

Start with **one core game**, not many mini-games.

Recommended foundation:

- One player knows the target state.
- One player can act in the world.
- They must speak German to succeed.
- The round is timed.
- Feedback comes after the round.

This is stronger than charades as a foundation because it produces:

- instructions
- clarification
- repair
- confirmation
- negotiation

That is closer to real speaking.

## Recommended MVP Scenario

### Arrange the Apartment

Two-player asymmetric co-op mission:

- `Player A: Mission Control`
  - Sees the target room setup.
  - Knows the objective state.
  - Cannot move objects.
- `Player B: On Site`
  - Sees the messy room.
  - Can move and place objects.
  - Cannot see the target setup.

They must speak German to make the room match the hidden target before time runs out.

### Why This Is the Best First Scenario

- Easy to build in 2D.
- High language density.
- Naturally trains:
  - prepositions
  - furniture/object vocabulary
  - imperatives
  - clarification questions
- Easy to vary for replayability.

## Language Focus

Do not try to teach broad German at first.

Train one domain well:

- `auf`
- `unter`
- `neben`
- `vor`
- `hinter`
- room objects
- simple commands

Examples:

- `Stell die Pflanze neben das Sofa.`
- `Leg das rote Buch auf den Tisch.`
- `Links oder rechts vom Sofa?`

## Voice Pipeline Decision

Voice is non-negotiable.

The correct proof order is:

1. Prove speech capture and transcription.
2. Prove transcript correction based on scenario/level prompt.
3. Prove flashcard generation.
4. Only then move into the real game loop.

### Current Technical POC

Project: [[German-Voice-Game-POC]]

Current pipeline:

- browser mic input
- transcription
- scenario-specific correction
- flashcard generation

This has already been scaffolded and tested locally.

## Important Design Decisions

### Non-negotiable

- web app
- voice-first
- post-round correction, not live correction
- one core game loop
- 2D point-and-click for the first real gameplay prototype
- scenario-based progression

### Flexible

- exact art style
- scoring model
- story wrapper
- whether the game later becomes a collection of scenario packs

## Replayability Strategy

Replayability should come from **variation inside one core loop**, not from building unrelated mini-games too early.

Good replay levers:

- role swap
- different room layouts
- different target states
- more distractor objects
- tighter time pressure
- more nuanced instructions
- new scenario packs later

## Why Not Start With Many Mini-Games

Too many mini-games too early will spread effort across:

- rules
- balancing
- UI
- art
- content
- testing

That makes it harder to prove whether the speaking concept itself works.

## Best Medium-Term Shape

If the concept works, the product can become:

- one engine
- many scenario packs
- each scenario pack tied to one language domain

Examples:

- apartment setup
- kitchen rush
- train station navigation
- lost luggage
- witness interview

For the broader expansion map, see [[Speaking Game - Scenario Families]].

## Cost Notes

Current assumption:

- transcription and correction can be cheap enough to support short multiplayer rounds
- post-round processing is much simpler and cheaper than live continuous correction

So the design should favor:

- record the whole round
- transcribe after the round
- generate targeted corrections after the round

## Current Conclusion

The strongest direction is:

**A voice-first cooperative speaking game with asymmetric information, starting with one scenario: Arrange the Apartment.**

This is a better foundation than:

- a static bomb-defusal clone
- a broad “teach all German” app
- many unrelated mini-games

It is focused enough to test, but extensible enough to grow.

## Next Questions

- How should the first playable apartment mission be scored?
- How long should one round last?
- Should both players be recorded separately or mixed?
- How should the transcript map back to player roles?
- What is the smallest real gameplay prototype after the voice pipeline?
