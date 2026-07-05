package save

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	
	// Create default config
	sm.config = sm.createDefaultConfig()
	err := sm.SaveConfig()
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}
	
	// Load config
	loadedConfig, err := sm.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	
	// Verify default values
	if !loadedConfig.ColorMode {
		t.Error("Expected ColorMode to be true")
	}
	if loadedConfig.ShowFPS {
		t.Error("Expected ShowFPS to be false")
	}
	if !loadedConfig.SoundEffects {
		t.Error("Expected SoundEffects to be true")
	}
	if loadedConfig.MusicVolume != 0.7 {
		t.Errorf("Expected MusicVolume 0.7, got %f", loadedConfig.MusicVolume)
	}
	if loadedConfig.SFXVolume != 0.8 {
		t.Errorf("Expected SFXVolume 0.8, got %f", loadedConfig.SFXVolume)
	}
	if loadedConfig.Difficulty != DifficultyNormal {
		t.Errorf("Expected DifficultyNormal, got %v", loadedConfig.Difficulty)
	}
	if loadedConfig.PlayerCount != 1 {
		t.Errorf("Expected PlayerCount 1, got %d", loadedConfig.PlayerCount)
	}
}

func TestLoadConfigNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	
	_, err := sm.LoadConfig()
	if err == nil {
		t.Error("Expected error for non-existent config, got nil")
	}
}

func TestSaveConfig(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.config = sm.createDefaultConfig()
	
	// Modify config
	sm.config.ColorMode = false
	sm.config.MusicVolume = 0.5
	
	err := sm.SaveConfig()
	if err != nil {
		t.Fatalf("SaveConfig failed: %v", err)
	}
	
	// Check file exists
	configPath := filepath.Join(tempDir, "config.json")
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file was not created")
	}
	
	// Load and verify
	loadedConfig, err := sm.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	
	if loadedConfig.ColorMode {
		t.Error("Expected ColorMode to be false")
	}
	if loadedConfig.MusicVolume != 0.5 {
		t.Errorf("Expected MusicVolume 0.5, got %f", loadedConfig.MusicVolume)
	}
}

func TestGetConfig(t *testing.T) {
	sm := NewSaveManager()
	sm.config = sm.createDefaultConfig()
	
	config := sm.GetConfig()
	if config == nil {
		t.Fatal("GetConfig returned nil")
	}
	if config != sm.config {
		t.Error("GetConfig returned different config instance")
	}
}

func TestUpdateConfig(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.config = sm.createDefaultConfig()
	
	newConfig := sm.createDefaultConfig()
	newConfig.ColorMode = false
	newConfig.Difficulty = DifficultyHard
	
	err := sm.UpdateConfig(newConfig)
	if err != nil {
		t.Fatalf("UpdateConfig failed: %v", err)
	}
	
	// Verify config was updated
	if sm.config.ColorMode {
		t.Error("Expected ColorMode to be false")
	}
	if sm.config.Difficulty != DifficultyHard {
		t.Errorf("Expected DifficultyHard, got %v", sm.config.Difficulty)
	}
	
	// Verify it was saved
	loadedConfig, err := sm.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if loadedConfig.ColorMode {
		t.Error("Expected saved ColorMode to be false")
	}
}

func TestValidateConfig(t *testing.T) {
	sm := NewSaveManager()
	
	// Test nil config
	err := sm.validateConfig(nil)
	if err == nil {
		t.Error("Expected error for nil config")
	}
	
	// Test invalid music volume
	config := sm.createDefaultConfig()
	config.MusicVolume = 1.5
	err = sm.validateConfig(config)
	if err == nil {
		t.Error("Expected error for invalid music volume")
	}
	
	config.MusicVolume = -0.1
	err = sm.validateConfig(config)
	if err == nil {
		t.Error("Expected error for negative music volume")
	}
	
	// Test invalid SFX volume
	config = sm.createDefaultConfig()
	config.SFXVolume = 1.5
	err = sm.validateConfig(config)
	if err == nil {
		t.Error("Expected error for invalid SFX volume")
	}
	
	// Test invalid difficulty
	config = sm.createDefaultConfig()
	config.Difficulty = Difficulty(99)
	err = sm.validateConfig(config)
	if err == nil {
		t.Error("Expected error for invalid difficulty")
	}
	
	// Test invalid player count
	config = sm.createDefaultConfig()
	config.PlayerCount = 0
	err = sm.validateConfig(config)
	if err == nil {
		t.Error("Expected error for player count 0")
	}
	
	config.PlayerCount = 3
	err = sm.validateConfig(config)
	if err == nil {
		t.Error("Expected error for player count 3")
	}
	
	// Test valid config
	config = sm.createDefaultConfig()
	err = sm.validateConfig(config)
	if err != nil {
		t.Errorf("Valid config failed validation: %v", err)
	}
}

func TestDifficultyString(t *testing.T) {
	tests := []struct {
		difficulty Difficulty
		expected   string
	}{
		{DifficultyEasy, "Easy"},
		{DifficultyNormal, "Normal"},
		{DifficultyHard, "Hard"},
		{Difficulty(99), "Unknown"},
	}
	
	for _, tt := range tests {
		result := tt.difficulty.String()
		if result != tt.expected {
			t.Errorf("Difficulty(%d).String() = %s, expected %s", tt.difficulty, result, tt.expected)
		}
	}
}

func TestSetColorMode(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.config = sm.createDefaultConfig()
	
	err := sm.SetColorMode(false)
	if err != nil {
		t.Fatalf("SetColorMode failed: %v", err)
	}
	
	if sm.config.ColorMode {
		t.Error("Expected ColorMode to be false")
	}
	
	// Verify it was saved
	loadedConfig, err := sm.LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if loadedConfig.ColorMode {
		t.Error("Expected saved ColorMode to be false")
	}
}

func TestSetShowFPS(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.config = sm.createDefaultConfig()
	
	err := sm.SetShowFPS(true)
	if err != nil {
		t.Fatalf("SetShowFPS failed: %v", err)
	}
	
	if !sm.config.ShowFPS {
		t.Error("Expected ShowFPS to be true")
	}
}

func TestSetSoundEffects(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.config = sm.createDefaultConfig()
	
	err := sm.SetSoundEffects(false)
	if err != nil {
		t.Fatalf("SetSoundEffects failed: %v", err)
	}
	
	if sm.config.SoundEffects {
		t.Error("Expected SoundEffects to be false")
	}
}

func TestSetMusicVolume(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.config = sm.createDefaultConfig()
	
	// Test valid volume
	err := sm.SetMusicVolume(0.5)
	if err != nil {
		t.Fatalf("SetMusicVolume failed: %v", err)
	}
	if sm.config.MusicVolume != 0.5 {
		t.Errorf("Expected MusicVolume 0.5, got %f", sm.config.MusicVolume)
	}
	
	// Test invalid volumes
	err = sm.SetMusicVolume(1.5)
	if err == nil {
		t.Error("Expected error for volume > 1")
	}
	
	err = sm.SetMusicVolume(-0.1)
	if err == nil {
		t.Error("Expected error for volume < 0")
	}
}

func TestSetSFXVolume(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.config = sm.createDefaultConfig()
	
	// Test valid volume
	err := sm.SetSFXVolume(0.6)
	if err != nil {
		t.Fatalf("SetSFXVolume failed: %v", err)
	}
	if sm.config.SFXVolume != 0.6 {
		t.Errorf("Expected SFXVolume 0.6, got %f", sm.config.SFXVolume)
	}
	
	// Test invalid volumes
	err = sm.SetSFXVolume(1.5)
	if err == nil {
		t.Error("Expected error for volume > 1")
	}
	
	err = sm.SetSFXVolume(-0.1)
	if err == nil {
		t.Error("Expected error for volume < 0")
	}
}

func TestSetDifficulty(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.config = sm.createDefaultConfig()
	
	// Test valid difficulties
	difficulties := []Difficulty{DifficultyEasy, DifficultyNormal, DifficultyHard}
	for _, diff := range difficulties {
		err := sm.SetDifficulty(diff)
		if err != nil {
			t.Fatalf("SetDifficulty(%v) failed: %v", diff, err)
		}
		if sm.config.Difficulty != diff {
			t.Errorf("Expected difficulty %v, got %v", diff, sm.config.Difficulty)
		}
	}
	
	// Test invalid difficulty
	err := sm.SetDifficulty(Difficulty(99))
	if err == nil {
		t.Error("Expected error for invalid difficulty")
	}
}

func TestSetPlayerCount(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.config = sm.createDefaultConfig()
	
	// Test valid counts
	for count := 1; count <= 2; count++ {
		err := sm.SetPlayerCount(count)
		if err != nil {
			t.Fatalf("SetPlayerCount(%d) failed: %v", count, err)
		}
		if sm.config.PlayerCount != count {
			t.Errorf("Expected player count %d, got %d", count, sm.config.PlayerCount)
		}
	}
	
	// Test invalid counts
	invalidCounts := []int{0, 3, -1, 100}
	for _, count := range invalidCounts {
		err := sm.SetPlayerCount(count)
		if err == nil {
			t.Errorf("Expected error for player count %d", count)
		}
	}
}

func TestDefaultPlayerKeys(t *testing.T) {
	player1Keys := getDefaultPlayer1Keys()
	player2Keys := getDefaultPlayer2Keys()
	
	// Verify player 1 keys are set
	if player1Keys.Left == "" {
		t.Error("Player 1 Left key not set")
	}
	if player1Keys.Right == "" {
		t.Error("Player 1 Right key not set")
	}
	if player1Keys.Jump == "" {
		t.Error("Player 1 Jump key not set")
	}
	
	// Verify player 2 keys are set
	if player2Keys.Left == "" {
		t.Error("Player 2 Left key not set")
	}
	if player2Keys.Right == "" {
		t.Error("Player 2 Right key not set")
	}
	if player2Keys.Jump == "" {
		t.Error("Player 2 Jump key not set")
	}
	
	// Verify they're different
	if player1Keys.Left == player2Keys.Left {
		t.Error("Player 1 and 2 should have different Left keys")
	}
}

func TestCreateDefaultConfig(t *testing.T) {
	sm := NewSaveManager()
	config := sm.createDefaultConfig()
	
	if config == nil {
		t.Fatal("createDefaultConfig returned nil")
	}
	
	// Verify all fields are set to expected defaults
	if !config.ColorMode {
		t.Error("Expected ColorMode to be true")
	}
	if config.ShowFPS {
		t.Error("Expected ShowFPS to be false")
	}
	if !config.SoundEffects {
		t.Error("Expected SoundEffects to be true")
	}
	if config.MusicVolume != 0.7 {
		t.Errorf("Expected MusicVolume 0.7, got %f", config.MusicVolume)
	}
	if config.SFXVolume != 0.8 {
		t.Errorf("Expected SFXVolume 0.8, got %f", config.SFXVolume)
	}
	if config.Difficulty != DifficultyNormal {
		t.Errorf("Expected DifficultyNormal, got %v", config.Difficulty)
	}
	if config.PlayerCount != 1 {
		t.Errorf("Expected PlayerCount 1, got %d", config.PlayerCount)
	}
}
