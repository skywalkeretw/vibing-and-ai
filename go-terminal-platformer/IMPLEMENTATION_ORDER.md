# Go Terminal Platformer - Complete Implementation Order

## Overview
This document provides the complete implementation order for all 31 GitHub issues, including the 7 enemy sub-issues and 4 world level creation issues.

## Issue Summary

### Original Issues (1-17)
- **Foundation**: #13, #2-5 (Project setup, game loop, renderer, input, physics)
- **Core Gameplay**: #6, #14, #9, #7, #16 (Player, camera, levels, enemies, projectiles)
- **Features**: #8, #15, #17, #10 (Power-ups, particles, bosses, UI)
- **Polish**: #11, #12 (Save system, audio)

### Enemy Sub-Issues (18-26)
- **#18**: Goomba Enemy (Basic Walker)
- **#19**: Koopa Enemy (Shell Enemy)
- **#20**: Piranha Plant Enemy (Stationary)
- **#21**: Lakitu Enemy (Flying)
- **#22**: Spiny Enemy (Spiked Walker)
- **#23**: Boo Enemy (Ghost)
- **#24**: Hammer Bro Enemy (Projectile Enemy)

### Level Creation Issues (27-30)
- **#27**: World 1 Levels (Grassland Plains) - 6 levels
- **#28**: World 2 Levels (Desert Dunes) - 6 levels
- **#29**: World 3 Levels (Ice Caverns) - 6 levels
- **#30**: World 4 Levels (Volcano Castle) - 6 levels

**Total**: 31 issues, 24 levels across 4 worlds

## Complete Implementation Order

### Phase 1: Foundation (Issues #13, #2-5)
**Critical Path - Must be completed in order**

1. **#13: Project Setup & Structure**
   - Initialize Go module
   - Set up directory structure
   - Create Makefile
   - **Blocks**: Everything

2. **#2: Core Game Loop and State Management**
   - Fixed timestep game loop (30 FPS)
   - State machine (Menu, Playing, Paused, GameOver)
   - **Blocks**: All gameplay systems

3. **#3: Terminal Renderer System** (Can parallel with #4)
   - tcell initialization
   - Drawing primitives
   - Layer-based rendering
   - Camera transformations
   - **Blocks**: All visual systems

4. **#4: Input Handler System** (Can parallel with #3)
   - Keyboard event polling
   - 2-player input support
   - **Blocks**: Player control, menus

5. **#5: Physics Engine and Collision Detection**
   - Gravity and velocity
   - AABB collision
   - Spatial grid optimization
   - **Blocks**: Player, enemies, projectiles

### Phase 2: Core Gameplay (Issues #6, #14, #9, #7)

6. **#6: Player Entity and Movement Mechanics**
   - Player physics and movement
   - State machine (idle, running, jumping, etc.)
   - 2-player support
   - **Blocks**: Camera, enemies, power-ups

7. **#14: Camera System with Smooth Scrolling**
   - Camera following with lerp
   - Split-screen for 2-player
   - **Blocks**: Level design, gameplay feel

8. **#9: Level System and Loader**
   - JSON level format
   - Level loader
   - Checkpoint system
   - **Blocks**: Level creation (#27-30)

9. **#7: Enemy Base System and AI**
   - Enemy base class
   - AI state machine
   - Pathfinding
   - Spawner system
   - **Blocks**: All enemy implementations (#18-26)

### Phase 3: Enemy Implementations (Issues #18-26)
**Complete after #7, order matters for dependencies**

10. **#18: Goomba Enemy** (First - simplest)
    - Basic walker
    - Reference implementation
    - **Blocks**: Level enemy placement

11. **#19: Koopa Enemy** (Second - needs Goomba reference)
    - Shell mechanics
    - Kick interaction
    - **Blocks**: Shell interactions in other enemies

12. **#20: Piranha Plant Enemy** (Third - needs level system)
    - Pipe emergence
    - Stomp immunity
    - **Blocks**: Pipe-based levels

13. **#22: Spiny Enemy** (Fourth - needs Koopa for shell)
    - Stomp immunity
    - Shell defeat mechanic
    - **Blocks**: Lakitu spawning

14. **#21: Lakitu Enemy** (Fifth - needs Spiny)
    - Flying behavior
    - Spiny spawning
    - **Blocks**: World 2-3 levels

15. **#23: Boo Enemy** (Sixth - needs power-up system)
    - Look-away chase
    - Star-only defeat
    - **Blocks**: World 3-4 levels

16. **#24: Hammer Bro Enemy** (Seventh - needs projectiles)
    - Hammer throwing
    - Platform standing
    - **Blocks**: World 3-4 levels

### Phase 4: Projectiles and Power-ups (Issues #16, #8)

17. **#16: Projectile System**
    - Projectile manager
    - Fireball and hammer projectiles
    - **Blocks**: Fire Flower, Hammer Bro

18. **#8: Power-Up System and Collectibles**
    - 6 power-up types
    - Coin collection
    - Block system
    - **Blocks**: Level design, gameplay variety

### Phase 5: Level Creation (Issues #27-30)
**Complete after enemies and power-ups are done**

19. **#27: World 1 Levels (Grassland Plains)**
    - 6 levels (tutorial to boss)
    - Goombas, Koopas, Piranha Plants
    - Goomba King boss
    - **Blocks**: World 2 progression

20. **#28: World 2 Levels (Desert Dunes)**
    - 6 levels (desert theme)
    - Lakitus, Spinies added
    - Sand Serpent boss
    - **Blocks**: World 3 progression

21. **#29: World 3 Levels (Ice Caverns)**
    - 6 levels (ice theme)
    - Boos, Hammer Bros added
    - Ice Golem boss
    - **Blocks**: World 4 progression

22. **#30: World 4 Levels (Volcano Castle)**
    - 6 levels (final world)
    - All enemy types
    - Fire Dragon final boss
    - **Blocks**: Game completion

### Phase 6: Visual Effects and Bosses (Issues #15, #17)

23. **#15: Particle Effects System**
    - Particle system with pooling
    - 5 particle types
    - **Blocks**: Visual polish

24. **#17: Boss Battle Mechanics System**
    - Boss manager
    - 4 boss types (Goomba King, Sand Serpent, Ice Golem, Fire Dragon)
    - Multi-phase system
    - **Blocks**: Boss levels (#27-30 level 6)

### Phase 7: UI and Menus (Issue #10)

25. **#10: UI System (Menus, HUD, Pause)**
    - Main menu
    - In-game HUD
    - Pause menu
    - Game over/level complete screens
    - **Blocks**: Complete game experience

### Phase 8: Persistence and Audio (Issues #11, #12)

26. **#11: Save System and Game State Persistence**
    - Save/load with 3 slots
    - Auto-save
    - High scores
    - **Blocks**: Progress tracking

27. **#12: Audio System (Sound Effects and Music)**
    - Sound effect playback
    - Background music
    - Volume control
    - **Blocks**: Audio experience

## Dependency Graph

```
#13 (Setup)
  ↓
#2 (Game Loop)
  ↓
#3 (Renderer) ← → #4 (Input)
  ↓
#5 (Physics)
  ↓
#6 (Player) → #14 (Camera)
  ↓           ↓
#9 (Levels) ← #7 (Enemy Base)
  ↓           ↓
  |         #18 (Goomba)
  |           ↓
  |         #19 (Koopa)
  |           ↓
  |         #20 (Piranha Plant)
  |           ↓
  |         #22 (Spiny)
  |           ↓
  |         #21 (Lakitu)
  |           ↓
  ↓         #23 (Boo) ← #8 (Power-ups)
  |           ↓
  |         #24 (Hammer Bro) ← #16 (Projectiles)
  |           ↓
  ↓         (All enemies done)
  |           ↓
#27 (World 1) → #28 (World 2) → #29 (World 3) → #30 (World 4)
  ↓
#15 (Particles) + #17 (Bosses)
  ↓
#10 (UI)
  ↓
#11 (Save) + #12 (Audio)
  ↓
COMPLETE
```

## Parallel Work Opportunities

### Can Work in Parallel:
- **#3 (Renderer)** and **#4 (Input)** after #2
- **#15 (Particles)** and **#16 (Projectiles)** after #7
- **#11 (Save)** and **#12 (Audio)** after #10

### Cannot Parallelize:
- Enemy implementations (#18-26) - must follow order due to dependencies
- Level creation (#27-30) - must follow world order
- Foundation issues (#13, #2, #5) - strict sequential order

## Enemy Implementation Order (Critical)

**Must implement in this exact order:**

1. **Goomba** (#18) - Base reference, simplest
2. **Koopa** (#19) - Needs Goomba reference for shell mechanics
3. **Piranha Plant** (#20) - Needs level system for pipes
4. **Spiny** (#22) - Needs Koopa for shell interaction
5. **Lakitu** (#21) - Needs Spiny for spawning
6. **Boo** (#23) - Needs power-up system for star defeat
7. **Hammer Bro** (#24) - Needs projectile system

## Level Creation Order (Critical)

**Must create in this exact order:**

1. **World 1** (#27) - Tutorial world, basic enemies
2. **World 2** (#28) - Desert world, adds Lakitu/Spiny
3. **World 3** (#29) - Ice world, adds Boo/Hammer Bro
4. **World 4** (#30) - Final world, all enemies

## Testing Checkpoints

### After Phase 1 (Foundation):
- ✅ Game loop runs at 30 FPS
- ✅ Terminal renders correctly
- ✅ Input responds to both players
- ✅ Physics calculates correctly

### After Phase 2 (Core Gameplay):
- ✅ Player moves and jumps
- ✅ Camera follows player
- ✅ Levels load from JSON
- ✅ Enemy base system works

### After Phase 3 (Enemies):
- ✅ All 7 enemy types functional
- ✅ Enemy AI works correctly
- ✅ Stomp and projectile defeat working

### After Phase 4 (Projectiles/Power-ups):
- ✅ Fireballs and hammers work
- ✅ All 6 power-ups functional
- ✅ Coins and blocks working

### After Phase 5 (Levels):
- ✅ All 24 levels playable
- ✅ Progressive difficulty curve
- ✅ 2-player co-op works in all levels

### After Phase 6 (Effects/Bosses):
- ✅ Particle effects enhance gameplay
- ✅ All 4 bosses functional
- ✅ Boss phases work correctly

### After Phase 7 (UI):
- ✅ Menus navigate correctly
- ✅ HUD displays all info
- ✅ Pause/resume works

### After Phase 8 (Persistence/Audio):
- ✅ Save/load works
- ✅ Audio plays correctly
- ✅ Settings persist

## Minimum Viable Product (MVP) Path

For fastest MVP, implement in this order:

1. **Foundation**: #13 → #2 → #3 → #4 → #5
2. **Core**: #6 → #14 → #9 → #7
3. **Basic Enemies**: #18 (Goomba) → #19 (Koopa) → #20 (Piranha Plant)
4. **Basic Features**: #16 (Projectiles) → #8 (Power-ups - Fire Flower only)
5. **World 1**: #27 (6 levels)
6. **Boss**: #17 (Goomba King only)
7. **UI**: #10 (Basic menus and HUD)
8. **Save**: #11

**MVP Deliverable**: Single-player World 1 with 3 enemy types, basic power-ups, and boss fight.

## Full Game Path

After MVP, continue with:

1. **Remaining Enemies**: #22 → #21 → #23 → #24
2. **Remaining Worlds**: #28 → #29 → #30
3. **Remaining Bosses**: #17 (complete all 4)
4. **Polish**: #15 (Particles) → #12 (Audio)

## Risk Mitigation

### High-Risk Issues:
- **#5 (Physics)**: Complex collision detection - allocate extra time
- **#7 (Enemy Base)**: Foundation for 7 enemies - must be solid
- **#9 (Level System)**: JSON parsing and tile system - test thoroughly
- **#17 (Bosses)**: Multi-phase mechanics - complex state management

### Mitigation Strategies:
- Write comprehensive unit tests for physics
- Create enemy base with extensibility in mind
- Validate level JSON format early
- Test boss phases independently

## Success Metrics

### Code Quality:
- All tests passing
- No race conditions (test with `-race`)
- No memory leaks
- Clean code following Go idioms

### Performance:
- Stable 30 FPS with 100+ entities
- No frame drops during boss fights
- Smooth camera movement
- Responsive input

### Gameplay:
- All 24 levels completable
- 2-player co-op functional
- All enemies behave correctly
- All power-ups work as designed
- All bosses beatable

### Polish:
- Particle effects enhance experience
- Audio adds to atmosphere
- UI is intuitive
- Save/load is reliable

## Estimated Timeline

- **Phase 1 (Foundation)**: 2-3 days
- **Phase 2 (Core Gameplay)**: 3-4 days
- **Phase 3 (Enemies)**: 4-5 days
- **Phase 4 (Projectiles/Power-ups)**: 2-3 days
- **Phase 5 (Levels)**: 5-7 days
- **Phase 6 (Effects/Bosses)**: 3-4 days
- **Phase 7 (UI)**: 2-3 days
- **Phase 8 (Persistence/Audio)**: 2-3 days

**Total Estimated Time**: 23-32 days for full implementation

**MVP Time**: 10-14 days

## Notes

- All issues are documented in GitHub with detailed specifications
- Each issue includes code examples and acceptance criteria
- Enemy implementations must follow the specified order
- Level creation requires all enemy types to be complete
- Test incrementally after each phase
- Use `gh issue view <number>` to see full issue details
- ASCII art is used for all sprites (no external assets needed)
- Game runs entirely in the terminal using tcell library
