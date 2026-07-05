package save

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// SaveManager manages game saves, configuration, and high scores
type SaveManager struct {
	saveDir     string
	currentSlot int
	saves       map[int]*SaveData
	config      *Config
}

// SaveData represents a single save file
type SaveData struct {
	Version string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`

	// Progress
	CurrentWorld   int            `json:"current_world"`
	CurrentLevel   int            `json:"current_level"`
	UnlockedLevels map[string]bool `json:"unlocked_levels"`

	// Player Stats
	TotalScore int           `json:"total_score"`
	TotalCoins int           `json:"total_coins"`
	TotalLives int           `json:"total_lives"`
	PlayTime   time.Duration `json:"play_time"`

	// High Scores
	LevelScores map[string]int `json:"level_scores"`
	WorldScores map[int]int    `json:"world_scores"`

	// Achievements (future)
	Achievements []string `json:"achievements"`
}

// SaveSlotInfo provides information about a save slot
type SaveSlotInfo struct {
	Slot       int           `json:"slot"`
	Empty      bool          `json:"empty"`
	Timestamp  time.Time     `json:"timestamp,omitempty"`
	World      int           `json:"world,omitempty"`
	Level      int           `json:"level,omitempty"`
	TotalScore int           `json:"total_score,omitempty"`
	PlayTime   time.Duration `json:"play_time,omitempty"`
}

// NewSaveManager creates a new SaveManager instance
func NewSaveManager() *SaveManager {
	return &SaveManager{
		saves: make(map[int]*SaveData),
	}
}

// Initialize sets up the save system
func (sm *SaveManager) Initialize() error {
	// Get save directory
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	sm.saveDir = filepath.Join(homeDir, ".go-terminal-platformer")

	// Create save directory if not exists
	err = os.MkdirAll(sm.saveDir, 0755)
	if err != nil {
		return fmt.Errorf("failed to create save directory: %w", err)
	}

	// Load config
	sm.config, err = sm.LoadConfig()
	if err != nil {
		// Create default config
		sm.config = sm.createDefaultConfig()
		if err := sm.SaveConfig(); err != nil {
			log.Printf("Warning: failed to save default config: %v", err)
		}
	}

	// Load all save slots
	for i := 1; i <= 3; i++ {
		save, err := sm.LoadSave(i)
		if err == nil {
			sm.saves[i] = save
		}
	}

	// Default to slot 1
	sm.currentSlot = 1

	log.Printf("Save system initialized at: %s", sm.saveDir)
	return nil
}

// SaveGame saves game data to a specific slot
func (sm *SaveManager) SaveGame(slot int, data *SaveData) error {
	// Validate slot
	if slot < 1 || slot > 3 {
		return fmt.Errorf("invalid save slot: %d (must be 1-3)", slot)
	}

	// Set version and timestamp
	data.Version = "1.0.0"
	data.Timestamp = time.Now()

	// Initialize maps if nil
	if data.UnlockedLevels == nil {
		data.UnlockedLevels = make(map[string]bool)
	}
	if data.LevelScores == nil {
		data.LevelScores = make(map[string]int)
	}
	if data.WorldScores == nil {
		data.WorldScores = make(map[int]int)
	}
	if data.Achievements == nil {
		data.Achievements = []string{}
	}

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %w", err)
	}

	// Write to file
	savePath := filepath.Join(sm.saveDir, fmt.Sprintf("save%d.json", slot))
	err = os.WriteFile(savePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write save file: %w", err)
	}

	// Update in-memory cache
	sm.saves[slot] = data

	log.Printf("Game saved to slot %d", slot)
	return nil
}

// LoadSave loads game data from a specific slot
func (sm *SaveManager) LoadSave(slot int) (*SaveData, error) {
	// Validate slot
	if slot < 1 || slot > 3 {
		return nil, fmt.Errorf("invalid save slot: %d (must be 1-3)", slot)
	}

	// Read file
	savePath := filepath.Join(sm.saveDir, fmt.Sprintf("save%d.json", slot))
	jsonData, err := os.ReadFile(savePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("save slot %d is empty", slot)
		}
		return nil, fmt.Errorf("failed to read save file: %w", err)
	}

	// Unmarshal
	var data SaveData
	err = json.Unmarshal(jsonData, &data)
	if err != nil {
		return nil, fmt.Errorf("corrupted save file in slot %d: %w", slot, err)
	}

	// Validate version
	if data.Version != "1.0.0" {
		return nil, fmt.Errorf("incompatible save version: %s (expected 1.0.0)", data.Version)
	}

	// Initialize maps if nil (for backwards compatibility)
	if data.UnlockedLevels == nil {
		data.UnlockedLevels = make(map[string]bool)
	}
	if data.LevelScores == nil {
		data.LevelScores = make(map[string]int)
	}
	if data.WorldScores == nil {
		data.WorldScores = make(map[int]int)
	}
	if data.Achievements == nil {
		data.Achievements = []string{}
	}

	log.Printf("Game loaded from slot %d", slot)
	return &data, nil
}

// DeleteSave deletes a save file from a specific slot
func (sm *SaveManager) DeleteSave(slot int) error {
	if slot < 1 || slot > 3 {
		return fmt.Errorf("invalid save slot: %d (must be 1-3)", slot)
	}

	savePath := filepath.Join(sm.saveDir, fmt.Sprintf("save%d.json", slot))
	err := os.Remove(savePath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete save file: %w", err)
	}

	delete(sm.saves, slot)
	log.Printf("Save slot %d deleted", slot)
	return nil
}

// AutoSave performs an automatic save to the current slot
func (sm *SaveManager) AutoSave(data *SaveData) error {
	if sm.currentSlot < 1 || sm.currentSlot > 3 {
		return fmt.Errorf("invalid current slot: %d", sm.currentSlot)
	}

	return sm.SaveGame(sm.currentSlot, data)
}

// SetCurrentSlot sets the active save slot
func (sm *SaveManager) SetCurrentSlot(slot int) error {
	if slot < 1 || slot > 3 {
		return fmt.Errorf("invalid save slot: %d (must be 1-3)", slot)
	}
	sm.currentSlot = slot
	log.Printf("Current save slot set to %d", slot)
	return nil
}

// GetCurrentSlot returns the active save slot
func (sm *SaveManager) GetCurrentSlot() int {
	return sm.currentSlot
}

// GetSaveSlotInfo returns information about a save slot
func (sm *SaveManager) GetSaveSlotInfo(slot int) *SaveSlotInfo {
	if slot < 1 || slot > 3 {
		return &SaveSlotInfo{
			Slot:  slot,
			Empty: true,
		}
	}

	data := sm.saves[slot]
	if data == nil {
		return &SaveSlotInfo{
			Slot:  slot,
			Empty: true,
		}
	}

	return &SaveSlotInfo{
		Slot:       slot,
		Empty:      false,
		Timestamp:  data.Timestamp,
		World:      data.CurrentWorld,
		Level:      data.CurrentLevel,
		TotalScore: data.TotalScore,
		PlayTime:   data.PlayTime,
	}
}

// GetAllSaveSlotInfo returns information about all save slots
func (sm *SaveManager) GetAllSaveSlotInfo() []*SaveSlotInfo {
	info := make([]*SaveSlotInfo, 3)
	for i := 1; i <= 3; i++ {
		info[i-1] = sm.GetSaveSlotInfo(i)
	}
	return info
}

// GetSaveDirectory returns the save directory path
func (sm *SaveManager) GetSaveDirectory() string {
	return sm.saveDir
}

// NewSaveData creates a new SaveData with default values
func NewSaveData() *SaveData {
	return &SaveData{
		Version:        "1.0.0",
		Timestamp:      time.Now(),
		CurrentWorld:   1,
		CurrentLevel:   1,
		UnlockedLevels: make(map[string]bool),
		TotalScore:     0,
		TotalCoins:     0,
		TotalLives:     5,
		PlayTime:       0,
		LevelScores:    make(map[string]int),
		WorldScores:    make(map[int]int),
		Achievements:   []string{},
	}
}
