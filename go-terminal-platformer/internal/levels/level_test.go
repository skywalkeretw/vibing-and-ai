package levels

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/lukeroy/go-terminal-platformer/internal/engine"
	"github.com/lukeroy/go-terminal-platformer/internal/entities"
)

func TestLoadLevel(t *testing.T) {
	// Create a temporary test level file
	tmpDir := t.TempDir()
	levelPath := filepath.Join(tmpDir, "test_level.json")

	levelJSON := `{
		"id": "test-level",
		"world": 1,
		"levelNum": 1,
		"name": "Test Level",
		"width": 50,
		"height": 20,
		"tiles": [
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"                                                  ",
			"=================================================="
		],
		"spawnPoints": [
			{"type": "player", "x": 50, "y": 100}
		],
		"checkpoints": [
			{"x": 200, "y": 100}
		],
		"goal": {"x": 400, "y": 100},
		"background": "grassland"
	}`

	err := os.WriteFile(levelPath, []byte(levelJSON), 0644)
	if err != nil {
		t.Fatalf("Failed to create test level file: %v", err)
	}

	// Test loading without physics engine
	level, err := LoadLevel(levelPath, nil)
	if err != nil {
		t.Fatalf("Failed to load level: %v", err)
	}

	// Verify level properties
	if level.ID != "test-level" {
		t.Errorf("Expected ID 'test-level', got '%s'", level.ID)
	}
	if level.World != 1 {
		t.Errorf("Expected World 1, got %d", level.World)
	}
	if level.LevelNum != 1 {
		t.Errorf("Expected LevelNum 1, got %d", level.LevelNum)
	}
	if level.Name != "Test Level" {
		t.Errorf("Expected Name 'Test Level', got '%s'", level.Name)
	}
	if level.Width != 50 {
		t.Errorf("Expected Width 50, got %d", level.Width)
	}
	if level.Height != 20 {
		t.Errorf("Expected Height 20, got %d", level.Height)
	}

	// Verify spawn points
	if len(level.SpawnPoints) != 1 {
		t.Errorf("Expected 1 spawn point, got %d", len(level.SpawnPoints))
	}

	// Verify checkpoints
	if len(level.Checkpoints) != 1 {
		t.Errorf("Expected 1 checkpoint, got %d", len(level.Checkpoints))
	}

	// Verify goal position
	if level.GoalPosition.X != 400 || level.GoalPosition.Y != 100 {
		t.Errorf("Expected goal at (400, 100), got (%.0f, %.0f)", level.GoalPosition.X, level.GoalPosition.Y)
	}

	// Verify background
	if level.Background == nil {
		t.Error("Expected background to be loaded")
	}
}

func TestParseTiles(t *testing.T) {
	level := &Level{}

	tileData := []string{
		"  =  ",
		"  -  ",
		"  ?  ",
		"=====",
	}

	level.parseTiles(tileData)

	// Check dimensions
	if level.Height != 4 {
		t.Errorf("Expected height 4, got %d", level.Height)
	}
	if level.Width != 5 {
		t.Errorf("Expected width 5, got %d", level.Width)
	}

	// Check specific tiles
	groundTile := level.Tiles[0][2]
	if groundTile.TileType != TileGround {
		t.Errorf("Expected TileGround at (2,0), got %v", groundTile.TileType)
	}
	if !groundTile.Solid {
		t.Error("Expected ground tile to be solid")
	}

	platformTile := level.Tiles[1][2]
	if platformTile.TileType != TilePlatform {
		t.Errorf("Expected TilePlatform at (2,1), got %v", platformTile.TileType)
	}
	if !platformTile.OneWay {
		t.Error("Expected platform tile to be one-way")
	}

	blockTile := level.Tiles[2][2]
	if blockTile.TileType != TileBlock {
		t.Errorf("Expected TileBlock at (2,2), got %v", blockTile.TileType)
	}

	emptyTile := level.Tiles[0][0]
	if emptyTile.TileType != TileEmpty {
		t.Errorf("Expected TileEmpty at (0,0), got %v", emptyTile.TileType)
	}
}

func TestCharToTile(t *testing.T) {
	level := &Level{}

	tests := []struct {
		char     rune
		expected TileType
		solid    bool
		oneWay   bool
	}{
		{'=', TileGround, true, false},
		{'-', TilePlatform, true, true},
		{'?', TileBlock, false, false},
		{'#', TilePipe, true, false},
		{'~', TileWater, false, false},
		{'*', TileLava, false, false},
		{'i', TileIce, true, false},
		{' ', TileEmpty, false, false},
	}

	for _, tt := range tests {
		tile := level.charToTile(tt.char)
		if tile.TileType != tt.expected {
			t.Errorf("char '%c': expected type %v, got %v", tt.char, tt.expected, tile.TileType)
		}
		if tile.Solid != tt.solid {
			t.Errorf("char '%c': expected solid %v, got %v", tt.char, tt.solid, tile.Solid)
		}
		if tile.OneWay != tt.oneWay {
			t.Errorf("char '%c': expected oneWay %v, got %v", tt.char, tt.oneWay, tile.OneWay)
		}
	}
}

func TestGetTileAt(t *testing.T) {
	level := &Level{
		TileSize: 16.0,
	}

	tileData := []string{
		"===",
		"   ",
		"===",
	}
	level.parseTiles(tileData)

	// Test valid coordinates
	tile := level.GetTileAt(16, 0)
	if tile == nil {
		t.Error("Expected tile at (16, 0), got nil")
	} else if tile.TileType != TileGround {
		t.Errorf("Expected TileGround at (16, 0), got %v", tile.TileType)
	}

	// Test empty tile
	tile = level.GetTileAt(16, 16)
	if tile == nil {
		t.Error("Expected tile at (16, 16), got nil")
	} else if tile.TileType != TileEmpty {
		t.Errorf("Expected TileEmpty at (16, 16), got %v", tile.TileType)
	}

	// Test out of bounds
	tile = level.GetTileAt(-10, 0)
	if tile != nil {
		t.Error("Expected nil for out of bounds coordinates, got tile")
	}

	tile = level.GetTileAt(1000, 1000)
	if tile != nil {
		t.Error("Expected nil for out of bounds coordinates, got tile")
	}
}

func TestIsSolidAt(t *testing.T) {
	level := &Level{
		TileSize: 16.0,
	}

	tileData := []string{
		"===",
		"   ",
		"===",
	}
	level.parseTiles(tileData)

	// Test solid tile
	if !level.IsSolidAt(16, 0) {
		t.Error("Expected solid tile at (16, 0)")
	}

	// Test empty tile
	if level.IsSolidAt(16, 16) {
		t.Error("Expected non-solid tile at (16, 16)")
	}

	// Test out of bounds
	if level.IsSolidAt(-10, 0) {
		t.Error("Expected non-solid for out of bounds")
	}
}

func TestCheckpointActivation(t *testing.T) {
	level := &Level{
		Checkpoints: []Checkpoint{
			{Position: engine.Vector2{X: 100, Y: 100}, Activated: false},
			{Position: engine.Vector2{X: 200, Y: 100}, Activated: false},
		},
		ActiveCheckpoint: -1,
	}

	// Create mock player
	player := entities.NewPlayer()
	player.SetPosition(engine.Vector2{X: 105, Y: 105})

	// Update level
	level.Update(0.016, []*entities.PlayerEntity{player})

	// Check if first checkpoint was activated
	if !level.Checkpoints[0].Activated {
		t.Error("Expected first checkpoint to be activated")
	}
	if level.ActiveCheckpoint != 0 {
		t.Errorf("Expected active checkpoint 0, got %d", level.ActiveCheckpoint)
	}

	// Move player to second checkpoint
	player.SetPosition(engine.Vector2{X: 205, Y: 105})
	level.Update(0.016, []*entities.PlayerEntity{player})

	// Check if second checkpoint was activated
	if !level.Checkpoints[1].Activated {
		t.Error("Expected second checkpoint to be activated")
	}
	if level.ActiveCheckpoint != 1 {
		t.Errorf("Expected active checkpoint 1, got %d", level.ActiveCheckpoint)
	}
}

func TestLevelCompletion(t *testing.T) {
	level := &Level{
		GoalPosition: engine.Vector2{X: 500, Y: 100},
		Completed:    false,
	}

	// Create mock player far from goal
	player := entities.NewPlayer()
	player.SetPosition(engine.Vector2{X: 100, Y: 100})

	// Update level
	level.Update(0.016, []*entities.PlayerEntity{player})

	// Level should not be completed
	if level.Completed {
		t.Error("Expected level to not be completed")
	}

	// Move player to goal
	player.SetPosition(engine.Vector2{X: 505, Y: 105})
	level.Update(0.016, []*entities.PlayerEntity{player})

	// Level should be completed
	if !level.Completed {
		t.Error("Expected level to be completed")
	}
}

func TestGetPlayerSpawnPoint(t *testing.T) {
	level := &Level{
		SpawnPoints: []*SpawnPoint{
			{Type: "player", X: 100, Y: 200},
			{Type: "enemy", X: 300, Y: 200},
		},
	}

	spawnPoint := level.GetPlayerSpawnPoint()
	if spawnPoint.X != 100 || spawnPoint.Y != 200 {
		t.Errorf("Expected spawn point (100, 200), got (%.0f, %.0f)", spawnPoint.X, spawnPoint.Y)
	}

	// Test with no player spawn point
	level.SpawnPoints = []*SpawnPoint{
		{Type: "enemy", X: 300, Y: 200},
	}

	spawnPoint = level.GetPlayerSpawnPoint()
	// Should return default spawn point
	if spawnPoint.X != 50 || spawnPoint.Y != 50 {
		t.Errorf("Expected default spawn point (50, 50), got (%.0f, %.0f)", spawnPoint.X, spawnPoint.Y)
	}
}

func TestGetLastCheckpointPosition(t *testing.T) {
	level := &Level{
		Checkpoints: []Checkpoint{
			{Position: engine.Vector2{X: 100, Y: 100}, Activated: true},
			{Position: engine.Vector2{X: 200, Y: 100}, Activated: true},
		},
		ActiveCheckpoint: 1,
		SpawnPoints: []*SpawnPoint{
			{Type: "player", X: 50, Y: 50},
		},
	}

	pos := level.GetLastCheckpointPosition()
	if pos.X != 200 || pos.Y != 100 {
		t.Errorf("Expected checkpoint position (200, 100), got (%.0f, %.0f)", pos.X, pos.Y)
	}

	// Test with no active checkpoint
	level.ActiveCheckpoint = -1
	pos = level.GetLastCheckpointPosition()
	if pos.X != 50 || pos.Y != 50 {
		t.Errorf("Expected spawn position (50, 50), got (%.0f, %.0f)", pos.X, pos.Y)
	}
}

func TestReset(t *testing.T) {
	level := &Level{
		Checkpoints: []Checkpoint{
			{Position: engine.Vector2{X: 100, Y: 100}, Activated: true},
			{Position: engine.Vector2{X: 200, Y: 100}, Activated: true},
		},
		ActiveCheckpoint: 1,
		Completed:        true,
	}

	level.Reset()

	// Check if level was reset
	if level.Completed {
		t.Error("Expected level to not be completed after reset")
	}
	if level.ActiveCheckpoint != -1 {
		t.Errorf("Expected active checkpoint -1, got %d", level.ActiveCheckpoint)
	}
	for i, cp := range level.Checkpoints {
		if cp.Activated {
			t.Errorf("Expected checkpoint %d to not be activated after reset", i)
		}
	}
}
