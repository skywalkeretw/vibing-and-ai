package save

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/gdamore/tcell/v2"
)

// Difficulty represents game difficulty levels
type Difficulty int

const (
	DifficultyEasy Difficulty = iota
	DifficultyNormal
	DifficultyHard
)

// String returns the string representation of the difficulty
func (d Difficulty) String() string {
	switch d {
	case DifficultyEasy:
		return "Easy"
	case DifficultyNormal:
		return "Normal"
	case DifficultyHard:
		return "Hard"
	default:
		return "Unknown"
	}
}

// PlayerControls defines key bindings for a player
type PlayerControls struct {
	Left      string `json:"left"`
	Right     string `json:"right"`
	Jump      string `json:"jump"`
	Run       string `json:"run"`
	Crouch    string `json:"crouch"`
	FirePower string `json:"fire_power"`
}

// Config represents game configuration and settings
type Config struct {
	// Display
	ColorMode bool `json:"color_mode"`
	ShowFPS   bool `json:"show_fps"`

	// Audio
	SoundEffects bool    `json:"sound_effects"`
	MusicVolume  float64 `json:"music_volume"`
	SFXVolume    float64 `json:"sfx_volume"`

	// Gameplay
	Difficulty  Difficulty `json:"difficulty"`
	PlayerCount int        `json:"player_count"`

	// Controls
	Player1Keys PlayerControls `json:"player1_keys"`
	Player2Keys PlayerControls `json:"player2_keys"`
}

// LoadConfig loads the configuration from disk
func (sm *SaveManager) LoadConfig() (*Config, error) {
	configPath := filepath.Join(sm.saveDir, "config.json")
	jsonData, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("config file does not exist")
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	err = json.Unmarshal(jsonData, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Validate config
	if err := sm.validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to disk
func (sm *SaveManager) SaveConfig() error {
	if sm.config == nil {
		return fmt.Errorf("no config to save")
	}

	// Validate before saving
	if err := sm.validateConfig(sm.config); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	jsonData, err := json.MarshalIndent(sm.config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	configPath := filepath.Join(sm.saveDir, "config.json")
	err = os.WriteFile(configPath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

// GetConfig returns the current configuration
func (sm *SaveManager) GetConfig() *Config {
	return sm.config
}

// UpdateConfig updates the configuration and saves it
func (sm *SaveManager) UpdateConfig(config *Config) error {
	if err := sm.validateConfig(config); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}

	sm.config = config
	return sm.SaveConfig()
}

// createDefaultConfig creates a default configuration
func (sm *SaveManager) createDefaultConfig() *Config {
	return &Config{
		ColorMode:    true,
		ShowFPS:      false,
		SoundEffects: true,
		MusicVolume:  0.7,
		SFXVolume:    0.8,
		Difficulty:   DifficultyNormal,
		PlayerCount:  1,
		Player1Keys:  getDefaultPlayer1Keys(),
		Player2Keys:  getDefaultPlayer2Keys(),
	}
}

// validateConfig validates configuration values
func (sm *SaveManager) validateConfig(config *Config) error {
	if config == nil {
		return fmt.Errorf("config is nil")
	}

	// Validate volumes
	if config.MusicVolume < 0 || config.MusicVolume > 1 {
		return fmt.Errorf("music volume must be between 0 and 1")
	}
	if config.SFXVolume < 0 || config.SFXVolume > 1 {
		return fmt.Errorf("sfx volume must be between 0 and 1")
	}

	// Validate difficulty
	if config.Difficulty < DifficultyEasy || config.Difficulty > DifficultyHard {
		return fmt.Errorf("invalid difficulty level")
	}

	// Validate player count
	if config.PlayerCount < 1 || config.PlayerCount > 2 {
		return fmt.Errorf("player count must be 1 or 2")
	}

	return nil
}

// getDefaultPlayer1Keys returns default key bindings for player 1
func getDefaultPlayer1Keys() PlayerControls {
	return PlayerControls{
		Left:      tcell.KeyNames[tcell.KeyLeft],
		Right:     tcell.KeyNames[tcell.KeyRight],
		Jump:      tcell.KeyNames[tcell.KeyUp],
		Run:       "Rune[x]",
		Crouch:    tcell.KeyNames[tcell.KeyDown],
		FirePower: "Rune[z]",
	}
}

// getDefaultPlayer2Keys returns default key bindings for player 2
func getDefaultPlayer2Keys() PlayerControls {
	return PlayerControls{
		Left:      "Rune[a]",
		Right:     "Rune[d]",
		Jump:      "Rune[w]",
		Run:       "Rune[g]",
		Crouch:    "Rune[s]",
		FirePower: "Rune[f]",
	}
}

// SetColorMode updates the color mode setting
func (sm *SaveManager) SetColorMode(enabled bool) error {
	sm.config.ColorMode = enabled
	return sm.SaveConfig()
}

// SetShowFPS updates the FPS display setting
func (sm *SaveManager) SetShowFPS(enabled bool) error {
	sm.config.ShowFPS = enabled
	return sm.SaveConfig()
}

// SetSoundEffects updates the sound effects setting
func (sm *SaveManager) SetSoundEffects(enabled bool) error {
	sm.config.SoundEffects = enabled
	return sm.SaveConfig()
}

// SetMusicVolume updates the music volume
func (sm *SaveManager) SetMusicVolume(volume float64) error {
	if volume < 0 || volume > 1 {
		return fmt.Errorf("volume must be between 0 and 1")
	}
	sm.config.MusicVolume = volume
	return sm.SaveConfig()
}

// SetSFXVolume updates the sound effects volume
func (sm *SaveManager) SetSFXVolume(volume float64) error {
	if volume < 0 || volume > 1 {
		return fmt.Errorf("volume must be between 0 and 1")
	}
	sm.config.SFXVolume = volume
	return sm.SaveConfig()
}

// SetDifficulty updates the difficulty setting
func (sm *SaveManager) SetDifficulty(difficulty Difficulty) error {
	if difficulty < DifficultyEasy || difficulty > DifficultyHard {
		return fmt.Errorf("invalid difficulty level")
	}
	sm.config.Difficulty = difficulty
	return sm.SaveConfig()
}

// SetPlayerCount updates the player count
func (sm *SaveManager) SetPlayerCount(count int) error {
	if count < 1 || count > 2 {
		return fmt.Errorf("player count must be 1 or 2")
	}
	sm.config.PlayerCount = count
	return sm.SaveConfig()
}
