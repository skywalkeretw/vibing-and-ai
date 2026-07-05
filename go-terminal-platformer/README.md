# Go Terminal Platformer

A terminal-based 2D platformer game inspired by Super Mario Bros, built in Go with support for 1-2 player cooperative gameplay.

## 🎮 Game Features

- **2-Player Co-op**: Play solo or with a friend using WASD and Arrow keys
- **4 Unique Worlds**: Grassland, Desert, Ice, and Volcano themes
- **24 Hand-Crafted Levels**: 6 levels per world including epic boss battles
- **7 Enemy Types**: From basic Goombas to challenging Hammer Bros
- **6 Power-Ups**: Fire flowers, invulnerability stars, speed boots, and more
- **Dynamic Terminal Rendering**: Adapts to terminal size changes in real-time
- **Save System**: Auto-save progress and continue where you left off
- **ASCII Art Style**: Colorful ANSI graphics with fallback for monochrome terminals

## 📋 Requirements

- **Go**: Version 1.21 or higher
- **Terminal**: Minimum 100x30 characters (120x40 recommended)
- **OS**: Linux, macOS, or Windows
- **Color Support**: ANSI 256-color support (optional but recommended)

## 🚀 Quick Start

### Installation

```bash
# Clone the repository
cd vibing-and-ai/go-terminal-platformer

# Install dependencies
go mod download

# Build the game
go build -o platformer cmd/game/main.go

# Run the game
./platformer
```

### Development

```bash
# Run without building
go run cmd/game/main.go

# Run tests
go test ./...

# Run with race detection
go test -race ./...
```

## 🎯 Controls

### Player 1 (WASD)
- **W**: Jump
- **A**: Move Left
- **S**: Crouch/Fast Fall
- **D**: Move Right
- **Space**: Shoot (when power-up active)

### Player 2 (Arrow Keys)
- **↑**: Jump
- **←**: Move Left
- **↓**: Crouch/Fast Fall
- **→**: Move Right
- **Right Shift**: Shoot (when power-up active)

### General
- **P**: Pause
- **ESC**: Quit to menu
- **M**: Mute sound effects

## 📖 Documentation

- **[Game Design Document](GAME_DESIGN.md)**: Complete game specifications, mechanics, and technical details
- **[Contributing Guide](CONTRIBUTING.md)**: Guidelines for contributing to the project (coming soon)
- **[API Documentation](docs/API.md)**: Code documentation (coming soon)

## 🏗️ Project Structure

```
go-terminal-platformer/
├── README.md                 # This file
├── GAME_DESIGN.md           # Complete game design specification
├── go.mod                   # Go module definition
├── go.sum                   # Go dependencies
├── main.go                  # Entry point
├── cmd/
│   └── game/
│       └── main.go          # Game initialization
├── internal/
│   ├── engine/              # Core game engine
│   ├── entities/            # Game entities (player, enemies, etc.)
│   ├── levels/              # Level loading and management
│   ├── ui/                  # Menus and HUD
│   └── save/                # Save system and configuration
├── assets/
│   ├── levels/              # Level data files
│   └── sprites/             # ASCII art definitions
└── tests/                   # Unit and integration tests
```

## 🎨 Game Worlds

### World 1: Grassland Plains
Classic green hills and blue skies. Learn the basics and face the Goomba King.

### World 2: Desert Dunes
Navigate sandy terrain and quicksand. Battle the Sand Serpent boss.

### World 3: Frozen Tundra
Slip and slide on icy platforms. Defeat the mighty Ice Golem.

### World 4: Volcanic Inferno
Dodge lava and fire jets in the final challenge. Face the Fire Dragon.

## 🔧 Configuration

Settings can be configured in-game or by editing `~/.go-terminal-platformer/config.json`:

```json
{
  "color_mode": true,
  "sound_effects": true,
  "difficulty": "normal",
  "player_count": 1,
  "controls": {
    "player1": {
      "up": "w",
      "left": "a",
      "down": "s",
      "right": "d",
      "shoot": "space"
    },
    "player2": {
      "up": "up",
      "left": "left",
      "down": "down",
      "right": "right",
      "shoot": "rshift"
    }
  }
}
```

## 🐛 Known Issues

See the [GitHub Issues](../../issues?q=is%3Aissue+is%3Aopen+label%3Aproject%3Ago-terminal-platformer) page for current bugs and feature requests.

## 🤝 Contributing

This project is part of the "Vibing and AI" repository for AI-assisted development experiments. Contributions are welcome!

1. Check existing issues or create a new one
2. Fork the repository
3. Create a feature branch
4. Make your changes
5. Submit a pull request

All issues are written to be AI-agent friendly for autonomous development.

## 📝 Development Status

**Current Phase**: Design & Planning

- [x] Game design document completed
- [ ] Core engine implementation
- [ ] Player mechanics
- [ ] Enemy system
- [ ] Power-up system
- [ ] Level design
- [ ] UI/Menu system
- [ ] Save system
- [ ] Testing & polish

## 📜 License

This project is part of the Vibing and AI repository. See the main repository for license information.

## 🙏 Acknowledgments

- Inspired by Super Mario Bros (Nintendo)
- Built with [tcell](https://github.com/gdamore/tcell) for terminal rendering
- Part of the AI-assisted development experiments

## 📞 Support

For questions or issues:
1. Check the [Game Design Document](GAME_DESIGN.md)
2. Search [existing issues](../../issues?q=is%3Aissue+label%3Aproject%3Ago-terminal-platformer)
3. Create a new issue with the `project:go-terminal-platformer` label

---

**Note**: This is an experimental project for learning and exploration. The game is under active development.
