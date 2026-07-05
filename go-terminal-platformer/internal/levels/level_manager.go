package levels

import (
	"fmt"

	"github.com/lukeroy/go-terminal-platformer/internal/engine"
)

// LevelManager manages level loading and progression
type LevelManager struct {
	CurrentLevel    *Level
	CurrentWorld    int
	CurrentLevelNum int
	Physics         *engine.PhysicsEngine
}

// NewLevelManager creates a new level manager
func NewLevelManager(physics *engine.PhysicsEngine) *LevelManager {
	return &LevelManager{
		Physics:         physics,
		CurrentWorld:    1,
		CurrentLevelNum: 1,
	}
}

// LoadLevel loads a specific level by world and level number
func (lm *LevelManager) LoadLevel(world, levelNum int) error {
	levelPath := fmt.Sprintf("assets/levels/world%d/level%d.json", world, levelNum)
	level, err := LoadLevel(levelPath, lm.Physics)
	if err != nil {
		return fmt.Errorf("failed to load level %d-%d: %w", world, levelNum, err)
	}

	lm.CurrentLevel = level
	lm.CurrentWorld = world
	lm.CurrentLevelNum = levelNum
	return nil
}

// NextLevel loads the next level in sequence
func (lm *LevelManager) NextLevel() error {
	nextLevel := lm.CurrentLevelNum + 1
	nextWorld := lm.CurrentWorld

	// Check if we need to move to the next world
	if nextLevel > 6 {
		// Each world has 6 levels
		nextLevel = 1
		nextWorld++
	}

	// Check if game is complete (4 worlds total)
	if nextWorld > 4 {
		// Game complete
		return fmt.Errorf("game completed - no more levels")
	}

	return lm.LoadLevel(nextWorld, nextLevel)
}

// PreviousLevel loads the previous level
func (lm *LevelManager) PreviousLevel() error {
	prevLevel := lm.CurrentLevelNum - 1
	prevWorld := lm.CurrentWorld

	// Check if we need to move to the previous world
	if prevLevel < 1 {
		prevWorld--
		prevLevel = 6 // Last level of previous world
	}

	// Check if we're at the first level
	if prevWorld < 1 {
		return fmt.Errorf("already at first level")
	}

	return lm.LoadLevel(prevWorld, prevLevel)
}

// RestartLevel reloads the current level
func (lm *LevelManager) RestartLevel() error {
	if lm.CurrentLevel == nil {
		return fmt.Errorf("no level currently loaded")
	}
	return lm.LoadLevel(lm.CurrentWorld, lm.CurrentLevelNum)
}

// LoadWorldBoss loads the boss level for a specific world
func (lm *LevelManager) LoadWorldBoss(world int) error {
	// Boss is always level 6 in each world
	return lm.LoadLevel(world, 6)
}

// GetCurrentLevel returns the currently loaded level
func (lm *LevelManager) GetCurrentLevel() *Level {
	return lm.CurrentLevel
}

// GetCurrentWorld returns the current world number
func (lm *LevelManager) GetCurrentWorld() int {
	return lm.CurrentWorld
}

// GetCurrentLevelNum returns the current level number within the world
func (lm *LevelManager) GetCurrentLevelNum() int {
	return lm.CurrentLevelNum
}

// IsLevelComplete checks if the current level is completed
func (lm *LevelManager) IsLevelComplete() bool {
	if lm.CurrentLevel == nil {
		return false
	}
	return lm.CurrentLevel.Completed
}

// IsBossLevel checks if the current level is a boss level
func (lm *LevelManager) IsBossLevel() bool {
	return lm.CurrentLevelNum == 6
}

// IsGameComplete checks if all levels have been completed
func (lm *LevelManager) IsGameComplete() bool {
	return lm.CurrentWorld > 4
}

// GetLevelProgress returns the overall game progress as a percentage
func (lm *LevelManager) GetLevelProgress() float64 {
	// Total levels: 4 worlds * 6 levels = 24 levels
	totalLevels := 24.0
	completedLevels := float64((lm.CurrentWorld-1)*6 + lm.CurrentLevelNum - 1)
	return (completedLevels / totalLevels) * 100.0
}

// GetWorldProgress returns the progress within the current world as a percentage
func (lm *LevelManager) GetWorldProgress() float64 {
	return (float64(lm.CurrentLevelNum) / 6.0) * 100.0
}

// UnloadLevel unloads the current level and clears physics colliders
func (lm *LevelManager) UnloadLevel() {
	if lm.CurrentLevel != nil {
		// Clear physics colliders
		if lm.Physics != nil {
			lm.Physics.ClearStaticColliders()
		}
		lm.CurrentLevel = nil
	}
}

// ResetProgress resets the level manager to the first level
func (lm *LevelManager) ResetProgress() error {
	lm.UnloadLevel()
	lm.CurrentWorld = 1
	lm.CurrentLevelNum = 1
	return lm.LoadLevel(1, 1)
}
