---
created: 2026-04-09
tags:
  - project
  - language-learning
  - game-design
  - mvp
  - 2-month-sprint
related:
  - "[[German Speaking Game - Project Hub]]"
  - "[[Apartment Setup - Game Design Blueprint]]"
  - "[[German-Voice-Game-POC]]"
---

# German Speaking Game - 2-Month MVP Design

> **📋 Project Hub:** Start with [[German Speaking Game - Project Hub]] for the complete project overview.
>
> **This document:** The complete 2-month MVP build plan. This is what you should build.

## What I Learned From Keep Talking and Nobody Explodes

After reading the manual, here are the core patterns that make it work:

### 1. Information Asymmetry
- **Defuser:** Sees the bomb (visual), has the manual (reference), takes action
- **Expert:** Has the knowledge (manual), but can't see the specific problem
- **Communication is mandatory:** Neither can succeed alone

### 2. Reference-Based Problem Solving
- Expert uses manual to lookup rules based on Defuser's descriptions
- Rules are conditional: "If X, then do Y. Otherwise, if Z..."
- Forces precise description: "Is there a lit indicator CAR?" vs "Is there a light?"

### 3. Multiple Independent Modules
- Each module is a self-contained puzzle
- Can solve in any order
- Varied mechanics keep it fresh

### 4. Time Pressure
- Countdown creates urgency
- Strikes create tension but not instant failure
- Forces efficient communication

### 5. Reference Information Flow
```
Defuser describes → Expert looks up → Expert instructs → Defuser acts → Repeat
```

## What We'll Adapt (Not Copy)

### ✅ We'll Use:
- Information asymmetry (both players need each other)
- Reference-based problem solving
- Time pressure
- Precise description requirements
- Multiple independent tasks

### ❌ We Won't Use:
- Bomb theme (too stressful, not language-focused)
- Binary pass/fail (language learning needs nuance)
- 3D bomb visualization (we're 2D icon-based)

## The MVP: One Core Scenario

### Scenario: The Restaurant Kitchen

**Why This Works:**
- Universal vocabulary (food, cooking, kitchen items)
- Natural instruction language (put, take, cut, cook)
- Visual and simple (2D kitchen with icons)
- Asymmetric information is natural (recipe vs. ingredients)

---

## Game Mechanics

### Player Roles

**Player 1: The Head Chef (Der Küchenchef)**
- Has: The recipe book (sees target dishes and steps)
- Cannot see: The kitchen layout or ingredient locations
- Must: Instruct Player 2 to cook the correct dishes

**Player 2: The Line Cook (Der Koch)**
- Has: The kitchen (sees ingredients, stations, tools)
- Cannot see: The recipes or what orders need to be cooked
- Must: Follow Player 1's instructions to cook

### The Information Gap

Neither can succeed alone:
- Chef knows WHAT to cook, but not WHERE things are
- Cook knows WHERE everything is, but not WHAT to cook
- Both must speak German to bridge the gap

### Time Pressure

- Orders come in on a ticket
- Timer counts down (2-5 minutes per round)
- Wrong dishes = wasted time (not instant strike)
- Complete orders = score points

---

## Visual Design (2D Icons)

### Kitchen Layout

```
┌─────────────────────────────────────────────┐
│                    ⏱️ 3:45                  │
├─────────────────┬───────────────────────────┤
│                 │                           │
│  🍳 STOVE       │   🧊 COUNTER             │
│  [🍳] [🍳] [🍳] │   [🧅] [🧄] [🍅] [🥬]    │
│                 │                           │
├─────────────────┼───────────────────────────┤
│                 │                           │
│  🔪 CUTTING     │   🥫 PANTRY              │
│  [              │   [🥩] [🐟] [🍗] [🦐]    │
│                 │                           │
├─────────────────┼───────────────────────────┤
│                 │                           │
│  🥣 MIXING      │   🧂 SPICES              │
│  [              │   [🧂] [🌶️] [🧄] [🍋]    │
│                 │                           │
└─────────────────┴───────────────────────────┘
```

**What Player 2 (Cook) sees:**
- The 2D kitchen grid with stations
- Ingredient icons at each station
- Can drag and drop ingredients between stations

**What Player 1 (Chef) sees:**
- Order ticket: "1x Tomato Soup, 2x Grilled Fish"
- Recipe book with steps:
  - "Tomato Soup: Cut onion and garlic → Put in pot → Add tomatoes → Cook 2 min"
  - "Grilled Fish: Season fish with salt and pepper → Put on stove → Cook 3 min"

**Neither sees:**
- What the other player sees

---

## Mutual Information Dependency

This is the key improvement over one-way asymmetry:

### Phase 1: Chef directs (0-1 min)
- Chef has recipe, Cook has kitchen access
- Chef: "Schneide die Zwiebel und den Knoblauch"
- Cook: "Wo ist der Knoblauch?"
- Chef: "Der Knoblauch ist in der Gewürzschar"

### Phase 2: Cook reports back (1-2 min)
- New information appears only Cook can see
- Cook notices: "Der Fisch ist angebrannt!"
- Cook must communicate: "Der Fisch ist schwarz!"
- Chef must adapt: "Nimm einen neuen Fisch!"

### Phase 3: Combined problem solving (2-3 min)
- Both have partial information about a new problem
- Example: "Der Herd funktioniert nicht" (only Cook sees stove is off)
- Chef sees recipe requires stove, must work around with Cook
- Together: "Können wir den Ofen benutzen?"

---

## Language Targets

### Vocabulary Domains

**Kitchen Verbs:**
- schneiden (cut)
- kochen (cook)
- braten (fry/grill)
- mischen (mix)
- würzen (season)
- nehmen (take)
- legen (put/place)
- stellen (stand/put)

**Prepositions:**
- in die Pfanne (into the pan)
- auf den Tisch (on the table)
- neben den Herd (next to the stove)
- unter das Dach (under the lid)

**Sequence Words:**
- zuerst (first)
- dann (then)
- danach (after that)
- zum Schluss (finally)

**Clarification:**
- Wo ist...? (where is...?)
- Welcher...? (which...?)
- Wie viel...? (how much...?)
- Meinst du...? (do you mean...?)

---

## Technical Implementation

### Tech Stack

**Frontend:**
- React/Next.js
- Tailwind CSS for layout
- Unicode emoji for icons (no external assets needed)

**Backend:**
- Node.js/Express
- WebRTC for peer-to-peer audio (or WebSockets via server)
- Whisper API for transcription

**Data:**
- JSON files for scenarios
- Simple session storage (no DB needed for MVP)

### Core Components

```javascript
// Kitchen grid (2D array of stations)
const kitchenLayout = [
  { id: 'stove', icon: '🍳', x: 0, y: 0, slots: 3 },
  { id: 'counter', icon: '🧊', x: 1, y: 0, items: ['🧅', '🧄', '🍅'] },
  { id: 'cutting', icon: '🔪', x: 0, y: 1, slots: 1 },
  { id: 'pantry', icon: '🥫', x: 1, y: 1, items: ['🥩', '🐟', '🍗'] },
  { id: 'mixing', icon: '🥣', x: 0, y: 2, slots: 1 },
  { id: 'spices', icon: '🧂', x: 1, y: 2, items: ['🧂', '🌶️', '🍋'] }
];

// Recipe structure
const recipe = {
  name: 'Tomatensuppe',
  ingredients: ['🍅', '🧅', '🧄'],
  steps: [
    'Schneide die Zwiebel und den Knoblauch',
    'Lege sie in den Topf',
    'Füge die Tomaten hinzu',
    'Koche alles 2 Minuten'
  ],
  time: 120 // seconds
};

// Session state
const session = {
  role: 'chef' | 'cook',
  view: 'kitchen' | 'recipe',
  inventory: [],
  completedSteps: [],
  score: 0,
  timeRemaining: 180
};
```

### Voice Pipeline

Already built in your POC - reuse it:
1. Record audio chunks during game
2. Transcribe after round ends
3. Generate feedback ( corrections, vocabulary learned )

---

## 2-Month Build Plan

### Month 1: Core Gameplay

**Week 1: Foundation**
- [ ] Set up Next.js project
- [ ] Build 2D kitchen grid with emoji icons
- [ ] Implement drag-and-drop for ingredients
- [ ] Basic state management (session, inventory, timer)

**Week 2: Two-Player Sync**
- [ ] WebRTC or WebSocket connection
- [ ] Role assignment (Chef vs Cook)
- [ ] Separate views for each role
- [ ] Basic lobby/join room flow

**Week 3: Scenario Logic**
- [ ] Recipe book component
- [ ] Order ticket system
- [ ] Cooking mechanics (place ingredient → timer → complete)
- [ ] Scoring system

**Week 4: Voice Integration**
- [ ] Reuse your existing voice pipeline
- [ ] Record during gameplay
- [ ] Transcribe after round
- [ ] Basic feedback screen

### Month 2: Polish & One More Scenario

**Week 5: First Scenario Polish**
- [ ] Kitchen scenario refinement
- [ ] Difficulty levels (easy/hard recipes)
- [ ] Visual polish (animations, feedback)
- [ ] Test with German learners

**Week 6: Second Scenario**
- [ ] Choose: Apartment Setup OR The Heist
- [ ] Build with reusable components
- [ ] Test variety of language domains

**Week 7: Session Flow**
- [ ] Lobby screen (enter room, choose role)
- [ ] Post-game feedback screen
- [ ] Simple history (sessions played)
- [ ] Basic onboarding/tutorial

**Week 8: Launch Prep**
- [ ] Bug fixes
- [ ] Performance optimization
- [ ] Deploy to production
- [ ] Test with real users

---

## Success Criteria

### MVP Must-Haves:
- ✅ Two players can connect and play
- ✅ Information asymmetry forces German communication
- ✅ Voice capture and post-round transcription works
- ✅ One complete scenario (Kitchen)
- ✅ Playable 5-minute round
- ✅ Fun enough to want to play again

### Nice-to-Haves (if time permits):
- Second scenario
- Simple profile tracking (what vocabulary was practiced)
- Difficulty adjustment
- Replays/recordings

### Deliberately Out of Scope:
- User profiles and adaptive learning
- LLM scenario generation
- More than 2 scenarios
- Complex analytics
- 3D graphics

---

## What Makes This Doable in 2 Months

### Simplicity Wins:
1. **One core scenario** - Kitchen is universal and visual
2. **Emoji icons** - No assets to create or source
3. **Reuse POC work** - Voice pipeline already exists
4. **No adaptive complexity** - Static scenarios for now
5. **2D grid** - Simple coordinate system, no physics

### Risk Mitigation:
- Week 4 checkpoint: Is core gameplay fun? If not, pivot
- Week 6 checkpoint: If second scenario is too much, stop at one
- Week 7 checkpoint: Polish over features

### The "Good Enough" Bar:
- Doesn't need to look professional
- Doesn't need to be perfectly balanced
- Doesn't need to cover all German grammar
- **Just needs to prove: cooperative German speaking is fun**

---

## Next Step

**Which scenario do you want to build first?**

1. **The Kitchen** (described above)
2. **Apartment Setup** (from your original docs)
3. **The Heist** (Hacker/Thief navigation)

I recommend starting with whichever feels most intuitive to YOU, because you'll be building it quickly and need to stay motivated.

Once you decide, we can break down Week 1 into specific tasks.
