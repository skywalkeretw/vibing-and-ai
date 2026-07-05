package levels

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/lukeroy/go-terminal-platformer/internal/engine"
	"github.com/lukeroy/go-terminal-platformer/internal/entities"
)

// TileType represents different types of tiles
type TileType int

const (
	TileEmpty TileType = iota
	TileGround
	TilePlatform
	TileOneWay
	TilePipe
	TileLava
	TileWater
	TileIce
	TileBlock
)

// Tile represents a single tile in the level
type Tile struct {
	TileType TileType
	Solid    bool
	OneWay   bool
	Sprite   rune
	Color    tcell.Color
}

// SpawnPoint represents an entity spawn location
type SpawnPoint struct {
	Type string  `json:"type"`
	X    float64 `json:"x"`
	Y    float64 `json:"y"`
}

// Checkpoint represents a checkpoint in the level
type Checkpoint struct {
	Position  engine.Vector2
	Activated bool
}

// Level represents a game level
type Level struct {
	ID       string `json:"id"`
	World    int    `json:"world"`
	LevelNum int    `json:"levelNum"`
	Name     string `json:"name"`

	// Dimensions
	Width  int `json:"width"`
	Height int `json:"height"`

	// Terrain
	Tiles        [][]Tile
	CollisionMap [][]bool

	// Entities
	SpawnPoints []*SpawnPoint
	Enemies     []entities.Entity
	PowerUps    []entities.Entity
	Coins       []entities.Entity
	Blocks      []entities.Entity

	// Checkpoints
	Checkpoints       []Checkpoint
	ActiveCheckpoint  int
	LastCheckpointPos engine.Vector2

	// Level bounds
	Bounds engine.Rectangle

	// Background
	Background *Background

	// Completion
	GoalPosition engine.Vector2
	Completed    bool

	// Tile size in pixels
	TileSize float64
}

// LevelData represents the JSON structure for level data
type LevelData struct {
	ID          string        `json:"id"`
	World       int           `json:"world"`
	LevelNum    int           `json:"levelNum"`
	Name        string        `json:"name"`
	Width       int           `json:"width"`
	Height      int           `json:"height"`
	Tiles       []string      `json:"tiles"`
	SpawnPoints []*SpawnPoint `json:"spawnPoints"`
	Checkpoints []struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"checkpoints"`
	Goal struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
	} `json:"goal"`
	Background string `json:"background"`
}

// LoadLevel loads a level from a JSON file
func LoadLevel(levelPath string, physics *engine.PhysicsEngine) (*Level, error) {
	// Read level file
	data, err := os.ReadFile(levelPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read level file: %w", err)
	}

	// Parse JSON
	var levelData LevelData
	err = json.Unmarshal(data, &levelData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse level JSON: %w", err)
	}

	// Create level
	level := &Level{
		ID:       levelData.ID,
		World:    levelData.World,
		LevelNum: levelData.LevelNum,
		Name:     levelData.Name,
		Width:    levelData.Width,
		Height:   levelData.Height,
		TileSize: 16.0, // 16 pixels per tile
	}

	// Parse tiles
	level.parseTiles(levelData.Tiles)

	// Create collision map
	if physics != nil {
		level.buildCollisionMap(physics)
	}

	// Load spawn points
	level.SpawnPoints = levelData.SpawnPoints

	// Load checkpoints
	level.Checkpoints = make([]Checkpoint, len(levelData.Checkpoints))
	for i, cp := range levelData.Checkpoints {
		level.Checkpoints[i] = Checkpoint{
			Position:  engine.Vector2{X: cp.X, Y: cp.Y},
			Activated: false,
		}
	}

	// Set goal
	level.GoalPosition = engine.Vector2{X: levelData.Goal.X, Y: levelData.Goal.Y}

	// Set level bounds
	level.Bounds = engine.Rectangle{
		X:      0,
		Y:      0,
		Width:  float64(level.Width) * level.TileSize,
		Height: float64(level.Height) * level.TileSize,
	}

	// Load background
	level.Background = LoadBackground(levelData.Background)

	return level, nil
}

// parseTiles converts string tile data to Tile structs
func (l *Level) parseTiles(tileData []string) {
	l.Height = len(tileData)
	if l.Height == 0 {
		return
	}

	// Find max width
	maxWidth := 0
	for _, row := range tileData {
		if len(row) > maxWidth {
			maxWidth = len(row)
		}
	}
	l.Width = maxWidth

	// Initialize tiles array
	l.Tiles = make([][]Tile, l.Height)
	for y := range l.Tiles {
		l.Tiles[y] = make([]Tile, l.Width)
	}

	// Parse each character
	for y, row := range tileData {
		for x, char := range row {
			if x < l.Width {
				l.Tiles[y][x] = l.charToTile(char)
			}
		}
		// Fill remaining columns with empty tiles
		for x := len(row); x < l.Width; x++ {
			l.Tiles[y][x] = l.charToTile(' ')
		}
	}
}

// charToTile converts a character to a Tile
func (l *Level) charToTile(char rune) Tile {
	switch char {
	case '=':
		return Tile{TileGround, true, false, '█', tcell.ColorBrown}
	case '-':
		return Tile{TilePlatform, true, true, '─', tcell.ColorGray}
	case '?':
		return Tile{TileBlock, false, false, '?', tcell.ColorYellow}
	case '#':
		return Tile{TilePipe, true, false, '║', tcell.ColorGreen}
	case '~':
		return Tile{TileWater, false, false, '≈', tcell.ColorBlue}
	case '*':
		return Tile{TileLava, false, false, '≈', tcell.ColorRed}
	case 'i':
		return Tile{TileIce, true, false, '▓', tcell.ColorLightCyan}
	default:
		return Tile{TileEmpty, false, false, ' ', tcell.ColorDefault}
	}
}

// buildCollisionMap creates collision data and adds colliders to physics engine
func (l *Level) buildCollisionMap(physics *engine.PhysicsEngine) {
	l.CollisionMap = make([][]bool, l.Height)

	for y := 0; y < l.Height; y++ {
		l.CollisionMap[y] = make([]bool, l.Width)
		for x := 0; x < l.Width; x++ {
			tile := l.Tiles[y][x]
			l.CollisionMap[y] = append(l.CollisionMap[y], tile.Solid)

			if tile.Solid {
				// Add static collider to physics
				colliderX := float64(x) * l.TileSize
				colliderY := float64(y) * l.TileSize

				if tile.OneWay {
					// One-way platform
					collider := engine.NewOneWayCollider(
						colliderX,
						colliderY,
						l.TileSize,
						l.TileSize,
						engine.LayerOneWayPlatform,
					)
					physics.AddStaticCollider(collider)
				} else {
					// Regular solid tile
					collider := engine.NewTileCollider(
						colliderX,
						colliderY,
						l.TileSize,
						true,
						engine.LayerTerrain,
					)
					physics.AddStaticCollider(collider)
				}
			}
		}
	}
}

// Update updates the level state
func (l *Level) Update(deltaTime float64, players []*entities.PlayerEntity) {
	// Update background parallax
	if l.Background != nil {
		l.Background.Update(deltaTime)
	}

	// Check checkpoint activation
	for i := range l.Checkpoints {
		if !l.Checkpoints[i].Activated {
			for _, player := range players {
				if player == nil {
					continue
				}
				playerPos := player.GetPosition()
				if playerPos.Distance(l.Checkpoints[i].Position) < 20 {
					l.Checkpoints[i].Activated = true
					l.ActiveCheckpoint = i
					l.LastCheckpointPos = l.Checkpoints[i].Position
					// TODO: Play checkpoint sound
				}
			}
		}
	}

	// Check level completion
	if !l.Completed {
		for _, player := range players {
			if player == nil {
				continue
			}
			playerPos := player.GetPosition()
			if playerPos.Distance(l.GoalPosition) < 20 {
				l.Completed = true
				// TODO: Trigger level complete event
			}
		}
	}
}

// Render renders the level
func (l *Level) Render(renderer interface{}, cameraX, cameraY, cameraWidth, cameraHeight int) {
	// Render background
	if l.Background != nil {
		l.Background.Render(renderer)
	}

	// Calculate visible tile range (camera culling)
	startX := engine.Max(0, float64(cameraX)/l.TileSize)
	endX := engine.Min(float64(l.Width), float64(cameraX+cameraWidth)/l.TileSize+1)
	startY := engine.Max(0, float64(cameraY)/l.TileSize)
	endY := engine.Min(float64(l.Height), float64(cameraY+cameraHeight)/l.TileSize+1)

	// Render tiles
	for y := int(startY); y < int(endY); y++ {
		for x := int(startX); x < int(endX); x++ {
			if y >= len(l.Tiles) || x >= len(l.Tiles[y]) {
				continue
			}
			tile := l.Tiles[y][x]
			if tile.TileType != TileEmpty {
				// TODO: Use actual renderer interface
				// renderer.DrawChar(int(float64(x)*l.TileSize), int(float64(y)*l.TileSize), tile.Sprite, tile.Color, tcell.ColorDefault)
			}
		}
	}

	// Render checkpoints
	for _, checkpoint := range l.Checkpoints {
		var sprite rune
		var color tcell.Color
		if checkpoint.Activated {
			sprite = '▼'
			color = tcell.ColorGreen
		} else {
			sprite = '▽'
			color = tcell.ColorGray
		}
		_ = sprite
		_ = color
		// TODO: Use actual renderer interface
		// renderer.DrawChar(int(checkpoint.Position.X), int(checkpoint.Position.Y), sprite, color, tcell.ColorDefault)
	}

	// Render goal
	// TODO: Use actual renderer interface
	// renderer.DrawChar(int(l.GoalPosition.X), int(l.GoalPosition.Y), '▓', tcell.ColorYellow, tcell.ColorDefault)
}

// GetTileAt returns the tile at the given world coordinates
func (l *Level) GetTileAt(x, y float64) *Tile {
	// Return nil for negative coordinates
	if x < 0 || y < 0 {
		return nil
	}
	
	tileX := int(x / l.TileSize)
	tileY := int(y / l.TileSize)

	if tileX >= l.Width || tileY >= len(l.Tiles) {
		return nil
	}
	
	if tileY < len(l.Tiles) && tileX < len(l.Tiles[tileY]) {
		return &l.Tiles[tileY][tileX]
	}

	return nil
}

// IsSolidAt checks if there's a solid tile at the given world coordinates
func (l *Level) IsSolidAt(x, y float64) bool {
	tile := l.GetTileAt(x, y)
	return tile != nil && tile.Solid
}

// GetSpawnPoint returns the spawn point for a given type
func (l *Level) GetSpawnPoint(spawnType string) *engine.Vector2 {
	for _, sp := range l.SpawnPoints {
		if sp.Type == spawnType {
			return &engine.Vector2{X: sp.X, Y: sp.Y}
		}
	}
	return nil
}

// GetPlayerSpawnPoint returns the player spawn position
func (l *Level) GetPlayerSpawnPoint() engine.Vector2 {
	sp := l.GetSpawnPoint("player")
	if sp != nil {
		return *sp
	}
	// Default spawn at top-left
	return engine.Vector2{X: 50, Y: 50}
}

// GetLastCheckpointPosition returns the position of the last activated checkpoint
func (l *Level) GetLastCheckpointPosition() engine.Vector2 {
	if l.ActiveCheckpoint >= 0 && l.ActiveCheckpoint < len(l.Checkpoints) {
		return l.Checkpoints[l.ActiveCheckpoint].Position
	}
	return l.GetPlayerSpawnPoint()
}

// Reset resets the level state
func (l *Level) Reset() {
	l.Completed = false
	l.ActiveCheckpoint = -1
	for i := range l.Checkpoints {
		l.Checkpoints[i].Activated = false
	}
}
