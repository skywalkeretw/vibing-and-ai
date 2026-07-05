# Go Terminal Platformer - Game Design Document

## Project Overview

**Project Name:** go-terminal-platformer  
**Genre:** Terminal-based 2D Platformer (Mario-style)  
**Language:** Go  
**Target Platforms:** Cross-platform (Linux, macOS, Windows)

## Core Game Mechanics

### Player Controls

#### Player 1 (WASD)
- **W**: Jump
- **A**: Move Left
- **S**: Crouch/Fast Fall
- **D**: Move Right
- **Space**: Shoot Projectile (when power-up active)

#### Player 2 (Arrow Keys)
- **↑**: Jump
- **←**: Move Left
- **↓**: Crouch/Fast Fall
- **→**: Move Right
- **Right Shift**: Shoot Projectile (when power-up active)

### Movement Physics
- **Running Speed**: 8 characters/second
- **Jump Height**: 5 characters
- **Jump Duration**: ~0.5 seconds
- **Gravity**: Realistic falling acceleration
- **Momentum**: Slight slide when changing direction

### Lives System
- **Starting Lives**: 5 per player
- **Game Over**: When lives reach 0
- **Extra Lives**: Collectible in levels (1-UP mushroom)

## Multiplayer System

### Mode
- **Type**: Cooperative (Co-op)
- **Players**: 1-2 simultaneous players
- **Screen**: Shared full terminal screen
- **Camera**: Follows both players (centers between them)
- **Respawn**: Dead player respawns at checkpoint or with partner

### Screen Management
- **Dynamic Resize**: Game adapts to terminal size changes
- **Minimum Size**: 100x30 characters (width x height)
- **Recommended**: 120x40 characters
- **Split Handling**: If players separate too far, camera prioritizes player 1

## World & Level Structure

### World Overview
- **Total Worlds**: 4
- **Levels per World**: 6 (5 regular + 1 boss)
- **Total Levels**: 24

### World Themes

#### World 1: Grassland Plains
- **Theme**: Green grass, blue sky, clouds
- **Hazards**: Small pits, basic enemies
- **Colors**: Green, blue, brown
- **Boss**: Giant Goomba King

#### World 2: Desert Dunes
- **Theme**: Sand, cacti, pyramids
- **Hazards**: Quicksand, sandstorms, larger pits
- **Colors**: Yellow, orange, tan
- **Boss**: Sand Serpent

#### World 3: Frozen Tundra
- **Theme**: Ice, snow, glaciers
- **Hazards**: Slippery ice, falling icicles, freezing water
- **Colors**: White, light blue, cyan
- **Boss**: Ice Golem

#### World 4: Volcanic Inferno
- **Theme**: Lava, fire, dark caves
- **Hazards**: Lava pools, fire jets, crumbling platforms
- **Colors**: Red, orange, dark gray
- **Boss**: Fire Dragon

### Level Design
- **Hand-Crafted**: All levels manually designed
- **Replayable**: Players can replay any unlocked level
- **Length**: 200-400 characters wide per level
- **Vertical Space**: 25-35 characters height
- **Secrets**: Hidden areas with extra coins/power-ups

## Enemy System

### Enemy Types (7 Total)

#### 1. Goomba (Basic Walker)
- **Behavior**: Walks left/right, turns at edges
- **Defeat**: Jump on head or projectile
- **Damage**: 1 hit point to player
- **Worlds**: 1, 2

#### 2. Koopa (Shell Enemy)
- **Behavior**: Walks, becomes shell when jumped on
- **Shell**: Can be kicked to defeat other enemies
- **Defeat**: Jump twice or projectile
- **Worlds**: 1, 3

#### 3. Piranha Plant (Stationary)
- **Behavior**: Emerges from pipes periodically
- **Defeat**: Projectile only (can't jump on)
- **Damage**: 1 hit point
- **Worlds**: 1, 2, 4

#### 4. Lakitu (Flying)
- **Behavior**: Flies above player, drops Spinies
- **Defeat**: Multiple projectiles (3 hits)
- **Damage**: 1 hit point
- **Worlds**: 2, 3

#### 5. Spiny (Spiked Walker)
- **Behavior**: Walks, cannot be jumped on
- **Defeat**: Projectile or shell only
- **Damage**: 1 hit point
- **Worlds**: 2, 3, 4

#### 6. Boo (Ghost)
- **Behavior**: Chases player when not looking, stops when facing
- **Defeat**: Invulnerability star only
- **Damage**: 1 hit point
- **Worlds**: 3, 4

#### 7. Hammer Bro (Projectile Enemy)
- **Behavior**: Stands on platforms, throws hammers
- **Defeat**: Jump or projectile (2 hits)
- **Damage**: 1 hit point
- **Worlds**: 3, 4

### Boss Mechanics

#### Common Boss Features
- **Health**: 10 hit points
- **Phases**: 3 phases (changes pattern at 7 HP and 3 HP)
- **Arena**: Enclosed boss room
- **Defeat**: Projectiles or jump attacks (varies by boss)

#### Boss Patterns
- **Phase 1**: Slow, predictable attacks
- **Phase 2**: Faster movement, more projectiles
- **Phase 3**: Aggressive, complex patterns

## Power-Up System

### Power-Up Types (6 Total)

#### 1. Fire Flower (Projectile)
- **Effect**: Shoot fireballs
- **Ammo**: 10 shots
- **Duration**: Until ammo depleted or hit by enemy
- **Visual**: Player changes color (red tint)

#### 2. Star (Invulnerability)
- **Effect**: Immune to damage, defeats enemies on contact
- **Duration**: 10 seconds
- **Visual**: Player flashes rainbow colors
- **Speed**: 1.5x movement speed

#### 3. Mushroom (Extra Life)
- **Effect**: +1 life
- **Instant**: Consumed immediately
- **Visual**: Red mushroom with white spots

#### 4. Speed Boots
- **Effect**: 1.5x running speed
- **Duration**: 30 seconds
- **Visual**: Player has speed lines

#### 5. Super Jump
- **Effect**: 2x jump height
- **Duration**: 30 seconds
- **Visual**: Player has spring icon

#### 6. Shield
- **Effect**: Absorbs 1 hit without losing life
- **Duration**: Until hit or level complete
- **Visual**: Bubble around player

### Power-Up Rules
- **No Stacking**: Only one power-up active at a time
- **Replacement**: New power-up replaces current one
- **Loss**: Lost when hit by enemy (except shield)
- **Blocks**: Hidden in ? blocks or floating in air

## Collectibles

### Coins
- **Value**: 1 point each
- **Extra Life**: 100 coins = 1 life
- **Locations**: Scattered throughout levels, in blocks
- **Visual**: Yellow spinning coin

### Checkpoints
- **Function**: Respawn point when player dies
- **Visual**: Flag pole or checkpoint banner
- **Activation**: Touch to activate
- **Persistence**: Saved in level progress

## Game Progression

### Save System
- **Auto-Save**: After completing each level
- **Save Data**: 
  - Current world and level
  - Lives remaining per player
  - Coins collected
  - High scores
  - Unlocked levels
- **Save Location**: `~/.go-terminal-platformer/save.json`

### Checkpoint System
- **In-Level**: 2-3 checkpoints per level
- **Respawn**: Player respawns at last checkpoint
- **Co-op**: Both players respawn together
- **Lives**: Checkpoint respawn costs 1 life

### Level Completion
- **Goal**: Reach end flag pole
- **Both Players**: Both must reach goal in co-op
- **Bonus**: Time bonus (faster = more points)
- **Unlock**: Next level unlocked on completion

### Scoring System
- **Coins**: 10 points each
- **Enemies**: 100 points each
- **Time Bonus**: (Time Remaining × 10) points
- **Level Complete**: 1000 points
- **Boss Defeat**: 5000 points
- **High Score**: Tracked per level and total

## UI/Menu System

### Main Menu
```
╔════════════════════════════════════════╗
║   GO TERMINAL PLATFORMER               ║
║                                        ║
║   [1] Start Game                       ║
║   [2] Continue                         ║
║   [3] Level Select                     ║
║   [4] Settings                         ║
║   [5] High Scores                      ║
║   [6] Exit                             ║
║                                        ║
║   Press number to select               ║
╚════════════════════════════════════════╝
```

### Settings Menu
- **Color Mode**: On/Off (ASCII only vs colored)
- **Difficulty**: Easy/Normal/Hard
  - Easy: 7 lives, more power-ups
  - Normal: 5 lives, standard
  - Hard: 3 lives, fewer power-ups
- **Sound Effects**: On/Off (terminal beep)
- **Player Count**: 1 or 2 players
- **Controls**: View/customize key bindings
- **Screen Size**: Display current terminal size

### HUD (In-Game)
```
P1: ♥♥♥♥♥  Coins: 045  Score: 12450  [FIRE]    World 1-3    P2: ♥♥♥♥♥
```

### Pause Menu
- **Resume**: Continue game
- **Restart Level**: Start level from beginning
- **Settings**: Quick settings access
- **Quit to Menu**: Return to main menu

## Technical Specifications

### Terminal Requirements
- **Minimum Size**: 100x30 characters
- **Recommended**: 120x40 characters
- **Color Support**: ANSI 256-color (fallback to 16-color)
- **Unicode**: UTF-8 support for better graphics

### Performance
- **Target FPS**: 30 frames per second
- **Input Polling**: 60Hz (16.67ms)
- **Physics Update**: 30Hz (33.33ms)
- **Render Update**: 30Hz (33.33ms)

### Go Libraries (Recommended)
- **Terminal UI**: `github.com/gdamore/tcell/v2`
- **Input Handling**: Built into tcell
- **Save System**: `encoding/json` (standard library)
- **Configuration**: `github.com/spf13/viper`
- **Logging**: `github.com/sirupsen/logrus`

### Project Structure
```
go-terminal-platformer/
├── README.md
├── GAME_DESIGN.md
├── go.mod
├── go.sum
├── main.go
├── cmd/
│   └── game/
│       └── main.go
├── internal/
│   ├── engine/
│   │   ├── game.go
│   │   ├── renderer.go
│   │   ├── input.go
│   │   └── physics.go
│   ├── entities/
│   │   ├── player.go
│   │   ├── enemy.go
│   │   ├── powerup.go
│   │   └── projectile.go
│   ├── levels/
│   │   ├── level.go
│   │   ├── loader.go
│   │   └── collision.go
│   ├── ui/
│   │   ├── menu.go
│   │   ├── hud.go
│   │   └── settings.go
│   └── save/
│       ├── save.go
│       └── config.go
├── assets/
│   ├── levels/
│   │   ├── world1/
│   │   ├── world2/
│   │   ├── world3/
│   │   └── world4/
│   └── sprites/
│       └── ascii_art.go
└── tests/
    ├── engine_test.go
    ├── physics_test.go
    └── collision_test.go
```

## ASCII Art Style

### Character Set
- **Player**: `@` or `☺`
- **Enemies**: Various ASCII characters (`G`, `K`, `P`, etc.)
- **Blocks**: `█`, `▓`, `▒`, `░`
- **Coins**: `○` or `◎`
- **Power-ups**: `♦`, `★`, `↑`, `◊`
- **Platforms**: `═`, `─`, `━`
- **Pipes**: `║`, `╔`, `╗`, `╚`, `╝`

### Color Scheme (ANSI)
- **Player 1**: Cyan/Blue
- **Player 2**: Magenta/Purple
- **Enemies**: Red/Yellow
- **Coins**: Yellow/Gold
- **Power-ups**: Bright colors (rainbow)
- **Background**: Dark gray/black
- **Platforms**: Brown/gray
- **Sky**: Blue gradient

## Sound Effects (Terminal Beep)

### Events
- **Jump**: Short beep (200ms)
- **Coin Collect**: High beep (100ms)
- **Enemy Defeat**: Medium beep (150ms)
- **Power-up Get**: Ascending beeps (3 tones)
- **Hit/Damage**: Low beep (300ms)
- **Level Complete**: Victory jingle (5 beeps)
- **Game Over**: Descending beeps (4 tones)

## Development Phases

### Phase 1: Core Engine (Issues #1-5)
- Game loop and rendering
- Input handling
- Physics engine
- Collision detection
- Terminal management

### Phase 2: Player Mechanics (Issues #6-10)
- Player entity
- Movement and jumping
- Animation system
- Player state management
- Multiplayer player handling

### Phase 3: Enemy System (Issues #11-15)
- Enemy base class
- Individual enemy types
- Enemy AI behaviors
- Boss mechanics
- Enemy spawning system

### Phase 4: Power-ups & Collectibles (Issues #16-20)
- Power-up system
- Collectible items (coins)
- Power-up effects
- Inventory management
- Item spawning

### Phase 5: Level System (Issues #21-25)
- Level loader
- Level data format
- World progression
- Checkpoint system
- Level design tools

### Phase 6: UI/Menus (Issues #26-30)
- Main menu
- Settings menu
- HUD system
- Pause menu
- Level select screen

### Phase 7: Save System (Issues #31-33)
- Save/load functionality
- Configuration management
- High score tracking

### Phase 8: Polish & Testing (Issues #34-40)
- Sound effects
- Visual effects
- Performance optimization
- Bug fixes
- Playtesting
- Documentation
- Build system

## Success Criteria

### Minimum Viable Product (MVP)
- ✓ Single player functional
- ✓ 1 world (6 levels) playable
- ✓ 3 enemy types working
- ✓ Basic power-ups (fire, star)
- ✓ Save/load system
- ✓ Main menu and HUD

### Full Release
- ✓ 2-player co-op functional
- ✓ All 4 worlds (24 levels) complete
- ✓ All 7 enemy types + 4 bosses
- ✓ All 6 power-ups working
- ✓ Complete UI/menu system
- ✓ Settings and configuration
- ✓ High score tracking
- ✓ Sound effects
- ✓ Cross-platform builds

## Future Enhancements (Post-Release)
- Level editor
- Custom level support
- Online leaderboards
- Speedrun timer
- Achievement system
- Additional worlds/DLC
- Competitive multiplayer mode
