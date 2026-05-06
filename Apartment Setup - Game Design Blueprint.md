---
created: 2026-03-30
tags:
  - project
  - game-design
  - language-learning
  - scenario
related:
  - "[[German Speaking Game - Project Hub]]"
  - "[[German Speaking Game - 2-Month MVP Design]]"
---

# Apartment Setup - Game Design Blueprint

> **📋 This is Scenario 2** - Build [[German Speaking Game - 2-Month MVP Design|The Kitchen]] first, then this if time permits.
>
> **🎯 Project Hub:** See [[German Speaking Game - Project Hub]] for all project files.

Related:

- [[German Speaking Game - POC Direction]]
- [[Clippings/German Learning App (Game)|German Learning App (Game)]]

## Purpose

`Apartment Setup` is the first real gameplay scenario for the German speaking game.

Its job is not to be a complete product. Its job is to prove one thing:

**Can asymmetric co-op gameplay make learners speak more German with more focus and better recall than an unstructured speaking session?**

## Design Goal

The round should feel:

- cooperative
- urgent
- slightly chaotic
- easy to understand
- language-dense

It should produce natural utterances like:

- instructions
- clarification questions
- confirmations
- corrections of misunderstanding
- short negotiation

## Core Loop

1. Two players join a round.
2. One player gets the target room state.
3. The other player gets the messy room and can manipulate objects.
4. They speak German to make the messy room match the target state.
5. Time runs out or all objectives are completed.
6. The game shows mission results.
7. The system generates post-game transcript review, corrections, and flashcards.

## Roles

### Player A: Mission Control

- Sees the target room layout.
- Sees the checklist of required placements.
- Cannot directly interact with the room.
- Must guide in German.

### Player B: On Site

- Sees the current messy room.
- Can drag, place, rotate, open, and close allowed objects.
- Cannot see the target layout.
- Must ask questions in German when uncertain.

## Why This Scenario Works

It naturally trains:

- prepositions of place
- furniture/object vocabulary
- colors and attributes
- imperative forms
- clarification and confirmation language

It is also cheap to build:

- one 2D room
- a fixed object set
- simple drag-and-drop interactions
- easy state validation

## Scene and Object Set

### Base Scene

One living room.

### Core Objects

- sofa
- coffee table
- chair
- floor lamp
- wall picture
- rug
- plant
- bookshelf
- red book
- blue cup

### Allowed Actions

- pick up
- place
- rotate
- hang
- open
- close

For the first version, keep actions minimal. Most rounds can work with only:

- drag object
- drop object
- hang picture on wall anchors

## Win Condition

The round is won when all target conditions are satisfied before the timer ends.

### Example Target Conditions

- The plant is next to the sofa.
- The red book is on the coffee table.
- The picture hangs above the sofa.
- The lamp stands behind the chair.

## Fail States

There should be only a few fail states.

### Main Fail State

- Timer reaches zero before all objectives are complete.

### Secondary Fail State

- Too many incorrect placements remain unresolved at the end of the round.

For MVP, avoid punitive fail states like:

- object destruction
- penalties for touching the wrong item
- random hazards

Those add stress without increasing language quality.

## Round Structure

### 1. Lobby

- Create room
- Invite partner
- Pick difficulty
- Ready up

### 2. Briefing

- 10 to 15 seconds
- Mission Control sees the target state
- On Site sees the messy room
- Short reminder: speak German

### 3. Active Round

- 3 to 5 minute timer
- Voice active
- No grammar interruption
- On Site manipulates objects
- Mission Control gives directions

### 4. End Screen

- success or failure
- objectives completed
- time used
- speaking summary

### 5. Debrief

- transcript snippets
- targeted corrections
- flashcards
- replay option
- swap roles option

## Mechanics

### Interaction Model

Use a 2D point-and-click room.

On Site can:

- click an object
- drag it
- drop it on a valid surface or anchor

Mission Control can:

- inspect the target room
- inspect checklist entries
- maybe zoom in if needed

### Validation Logic

Each object is defined by:

- id
- type
- position
- anchor/surface
- orientation if relevant

The game compares current object relations against target object relations.

The validation should be semantic, not pixel-perfect.

Example:

- `plant next to sofa` should allow a small range
- `picture above sofa` should match wall anchor zones

## Scoring

Keep scoring simple and legible.

### Mission Score

Mission score should come from:

- objectives completed
- time remaining
- number of hints used

Possible formula:

- 70% objective completion
- 20% time bonus
- 10% hint penalty

### Language Score

Language score should come from:

- speaking participation by both players
- number of meaningful turns
- number of clarification exchanges
- use of target vocabulary

Do not claim this is a full language proficiency score. It is only a scenario score.

### Star Rating

For the player-facing view, stars are clearer than raw math.

- 1 star: mission mostly failed, but some communication happened
- 2 stars: mission completed or nearly completed
- 3 stars: mission completed efficiently with strong communication

## Difficulty Tiers

Difficulty should come from information complexity, not control complexity.

### A1

- 4 target conditions
- few objects
- no distractors
- only basic prepositions
- clear visual contrast

Example focus:

- auf
- unter
- neben
- vor

### A2

- 6 target conditions
- more objects
- color and size distinctions
- a few distractors
- more article/case opportunities

Example focus:

- hinter
- zwischen
- links von
- rechts von

### B1

- 8 target conditions
- more ambiguity
- multiple similar objects
- indirect phrasing
- more need for clarification

Example focus:

- precise descriptions
- natural command phrasing
- confirmation and repair strategies

## What Counts as Good Language Output

The game should reward:

- clear instructions
- useful clarification
- confirmation language
- recovery from misunderstanding

Examples:

- `Stell die Pflanze neben das Sofa.`
- `Meinst du links oder rechts?`
- `Nein, das andere Buch.`
- `So?`
- `Ja, genau.`

This is better than rewarding only perfect grammar.

## Hint System

Hints should exist, but be limited.

Recommended:

- each team gets 1 or 2 hints
- using a hint reduces mission score slightly

Hint examples:

- highlight the target object category
- narrow the placement zone
- show one checklist line in simplified German

Do not let hints replace speaking.

## What the Game Must Not Do

- no live grammar popups during the round
- no forced tutoring behavior between players
- no over-detailed correction wall after the round
- no heavy reflex gameplay
- no too-precise placement requirements

## Post-Game Correction Design

This is where gameplay and learning connect.

The game should not send the whole transcript to the model with a vague prompt like:

- `correct the German`

Instead, the correction should be shaped by the scenario.

### Scenario Review Inputs

The model should receive:

- transcript
- player role labels
- difficulty tier
- scenario id
- target language focus
- list of target vocabulary and grammar structures

### Apartment Setup Review Focus

For this scenario, prioritize:

- prepositions of place
- article/case mistakes that affect placement phrases
- imperative forms
- useful object vocabulary

Ignore:

- minor off-focus mistakes
- filler words
- tiny phrasing imperfections that do not matter

### Output Limits

Return only:

- 3 to 5 priority corrections
- 2 positive examples
- 3 to 5 flashcards

That keeps the review useful instead of overwhelming.

## Mapping Gameplay to Correction

Gameplay events should help decide what to review.

### Examples

If the team struggled with object placement:

- prioritize corrections involving prepositions

If the same object was confused repeatedly:

- prioritize flashcards for that object vocabulary

If the team kept repairing misunderstanding:

- surface one or two strong clarification phrases as positive examples

This means review is not just transcript-based. It is **transcript + gameplay context**.

## Example Review Prompt Logic

For `Apartment Setup` at `A1`:

- only review basic location language
- ignore advanced grammar
- prioritize corrections like:
  - `auf dem Tisch`
  - `neben dem Sofa`
  - `unter dem Stuhl`

For `B1`:

- include more natural phrasing improvements
- include article/case precision where it clearly improves fluency

## Sample Flashcards

### Vocabulary

- Front: `plant`
- Back: `die Pflanze`
- Hint: `living room object`

### Phrase

- Front: `next to the sofa`
- Back: `neben dem Sofa`
- Hint: `location phrase`

### Action

- Front: `Place the red book on the table.`
- Back: `Leg das rote Buch auf den Tisch.`
- Hint: `imperative + object placement`

## Replayability Plan

Replayability should come from variation in the same scenario.

### Variation Levers

- swap roles
- random object positions
- random target placements
- different object subsets
- more distractors
- altered time limit
- different language focus within the same room

This lets one room support many rounds before needing a new scenario.

## Telemetry to Measure Whether the Scenario Works

For the first real tests, track:

- total speaking time per player
- number of speaking turns
- number of clarification turns
- mission completion rate
- average time to completion
- average number of corrections surfaced
- replay rate

The most important signals are:

- did both players speak a lot
- did they want another round
- did the review feel useful

## MVP Version of Apartment Setup

The first playable version should be very small.

### Scope

- one room
- 6 to 8 objects
- one difficulty tier at first
- one timer
- one success state
- one post-game review format

### Hard Constraints

- no 3D
- no story mode
- no inventory system
- no complicated physics
- no AI during the round beyond recording

## Future Extensions

If this scenario works, it becomes the template for the rest of the game.

Future scenario packs could reuse the same structure:

- Kitchen Rush
- Train Station Navigation
- Lost Luggage Sorting
- Witness Interview

## Final Design Summary

`Apartment Setup` should be:

- the first true gameplay scenario
- asymmetric
- voice-first
- cooperative
- simple to manipulate
- semantically validated
- reviewed after the round with scenario-specific corrections

If this scenario feels fun and produces a lot of meaningful speech, then the concept is real.
