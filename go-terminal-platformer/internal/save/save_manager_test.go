package save

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewSaveManager(t *testing.T) {
	sm := NewSaveManager()
	if sm == nil {
		t.Fatal("NewSaveManager returned nil")
	}
	if sm.saves == nil {
		t.Error("saves map not initialized")
	}
}

func TestSaveManagerInitialize(t *testing.T) {
	// Create temporary directory for testing
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	
	err := sm.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}
	
	// Check save directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		t.Error("Save directory was not created")
	}
	
	// Check config was created
	if sm.config == nil {
		t.Error("Config was not initialized")
	}
	
	// Check default slot is set
	if sm.currentSlot != 1 {
		t.Errorf("Expected default slot 1, got %d", sm.currentSlot)
	}
}

func TestSaveAndLoadGame(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.currentSlot = 1
	
	// Create test save data
	saveData := &SaveData{
		CurrentWorld:   2,
		CurrentLevel:   3,
		UnlockedLevels: map[string]bool{"world1-level1": true},
		TotalScore:     5000,
		TotalCoins:     50,
		TotalLives:     3,
		PlayTime:       time.Hour,
		LevelScores:    map[string]int{"world1-level1": 1000},
		WorldScores:    map[int]int{1: 5000},
		Achievements:   []string{"first_coin"},
	}
	
	// Save game
	err := sm.SaveGame(1, saveData)
	if err != nil {
		t.Fatalf("SaveGame failed: %v", err)
	}
	
	// Check file exists
	savePath := filepath.Join(tempDir, "save1.json")
	if _, err := os.Stat(savePath); os.IsNotExist(err) {
		t.Error("Save file was not created")
	}
	
	// Load game
	loadedData, err := sm.LoadSave(1)
	if err != nil {
		t.Fatalf("LoadSave failed: %v", err)
	}
	
	// Verify data
	if loadedData.CurrentWorld != saveData.CurrentWorld {
		t.Errorf("Expected world %d, got %d", saveData.CurrentWorld, loadedData.CurrentWorld)
	}
	if loadedData.CurrentLevel != saveData.CurrentLevel {
		t.Errorf("Expected level %d, got %d", saveData.CurrentLevel, loadedData.CurrentLevel)
	}
	if loadedData.TotalScore != saveData.TotalScore {
		t.Errorf("Expected score %d, got %d", saveData.TotalScore, loadedData.TotalScore)
	}
	if loadedData.TotalCoins != saveData.TotalCoins {
		t.Errorf("Expected coins %d, got %d", saveData.TotalCoins, loadedData.TotalCoins)
	}
	if loadedData.TotalLives != saveData.TotalLives {
		t.Errorf("Expected lives %d, got %d", saveData.TotalLives, loadedData.TotalLives)
	}
	if loadedData.PlayTime != saveData.PlayTime {
		t.Errorf("Expected playtime %v, got %v", saveData.PlayTime, loadedData.PlayTime)
	}
	if loadedData.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", loadedData.Version)
	}
}

func TestSaveGameInvalidSlot(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	
	saveData := NewSaveData()
	
	// Test invalid slots
	invalidSlots := []int{0, 4, -1, 100}
	for _, slot := range invalidSlots {
		err := sm.SaveGame(slot, saveData)
		if err == nil {
			t.Errorf("Expected error for invalid slot %d, got nil", slot)
		}
	}
}

func TestLoadSaveInvalidSlot(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	
	// Test invalid slots
	invalidSlots := []int{0, 4, -1, 100}
	for _, slot := range invalidSlots {
		_, err := sm.LoadSave(slot)
		if err == nil {
			t.Errorf("Expected error for invalid slot %d, got nil", slot)
		}
	}
}

func TestLoadSaveNonExistent(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	
	_, err := sm.LoadSave(1)
	if err == nil {
		t.Error("Expected error for non-existent save, got nil")
	}
}

func TestDeleteSave(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.saves = make(map[int]*SaveData)
	
	// Create a save
	saveData := NewSaveData()
	err := sm.SaveGame(1, saveData)
	if err != nil {
		t.Fatalf("SaveGame failed: %v", err)
	}
	
	// Delete the save
	err = sm.DeleteSave(1)
	if err != nil {
		t.Fatalf("DeleteSave failed: %v", err)
	}
	
	// Check file is deleted
	savePath := filepath.Join(tempDir, "save1.json")
	if _, err := os.Stat(savePath); !os.IsNotExist(err) {
		t.Error("Save file was not deleted")
	}
	
	// Check in-memory cache is cleared
	if _, exists := sm.saves[1]; exists {
		t.Error("Save was not removed from cache")
	}
}

func TestDeleteSaveInvalidSlot(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	
	invalidSlots := []int{0, 4, -1, 100}
	for _, slot := range invalidSlots {
		err := sm.DeleteSave(slot)
		if err == nil {
			t.Errorf("Expected error for invalid slot %d, got nil", slot)
		}
	}
}

func TestAutoSave(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.currentSlot = 2
	
	saveData := NewSaveData()
	saveData.TotalScore = 10000
	
	err := sm.AutoSave(saveData)
	if err != nil {
		t.Fatalf("AutoSave failed: %v", err)
	}
	
	// Verify it saved to the current slot
	loadedData, err := sm.LoadSave(2)
	if err != nil {
		t.Fatalf("LoadSave failed: %v", err)
	}
	
	if loadedData.TotalScore != 10000 {
		t.Errorf("Expected score 10000, got %d", loadedData.TotalScore)
	}
}

func TestSetCurrentSlot(t *testing.T) {
	sm := NewSaveManager()
	
	// Test valid slots
	for slot := 1; slot <= 3; slot++ {
		err := sm.SetCurrentSlot(slot)
		if err != nil {
			t.Errorf("SetCurrentSlot(%d) failed: %v", slot, err)
		}
		if sm.GetCurrentSlot() != slot {
			t.Errorf("Expected current slot %d, got %d", slot, sm.GetCurrentSlot())
		}
	}
	
	// Test invalid slots
	invalidSlots := []int{0, 4, -1, 100}
	for _, slot := range invalidSlots {
		err := sm.SetCurrentSlot(slot)
		if err == nil {
			t.Errorf("Expected error for invalid slot %d, got nil", slot)
		}
	}
}

func TestGetSaveSlotInfo(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.saves = make(map[int]*SaveData)
	
	// Test empty slot
	info := sm.GetSaveSlotInfo(1)
	if !info.Empty {
		t.Error("Expected empty slot")
	}
	if info.Slot != 1 {
		t.Errorf("Expected slot 1, got %d", info.Slot)
	}
	
	// Create a save
	saveData := &SaveData{
		CurrentWorld: 3,
		CurrentLevel: 5,
		TotalScore:   15000,
		PlayTime:     2 * time.Hour,
		Timestamp:    time.Now(),
	}
	sm.saves[1] = saveData
	
	// Test populated slot
	info = sm.GetSaveSlotInfo(1)
	if info.Empty {
		t.Error("Expected non-empty slot")
	}
	if info.World != 3 {
		t.Errorf("Expected world 3, got %d", info.World)
	}
	if info.Level != 5 {
		t.Errorf("Expected level 5, got %d", info.Level)
	}
	if info.TotalScore != 15000 {
		t.Errorf("Expected score 15000, got %d", info.TotalScore)
	}
}

func TestGetAllSaveSlotInfo(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.saves = make(map[int]*SaveData)
	
	// Create saves in slots 1 and 3
	sm.saves[1] = &SaveData{CurrentWorld: 1, TotalScore: 1000}
	sm.saves[3] = &SaveData{CurrentWorld: 3, TotalScore: 3000}
	
	allInfo := sm.GetAllSaveSlotInfo()
	
	if len(allInfo) != 3 {
		t.Errorf("Expected 3 slot infos, got %d", len(allInfo))
	}
	
	// Check slot 1
	if allInfo[0].Empty {
		t.Error("Slot 1 should not be empty")
	}
	if allInfo[0].World != 1 {
		t.Errorf("Expected world 1 in slot 1, got %d", allInfo[0].World)
	}
	
	// Check slot 2 (should be empty)
	if !allInfo[1].Empty {
		t.Error("Slot 2 should be empty")
	}
	
	// Check slot 3
	if allInfo[2].Empty {
		t.Error("Slot 3 should not be empty")
	}
	if allInfo[2].World != 3 {
		t.Errorf("Expected world 3 in slot 3, got %d", allInfo[2].World)
	}
}

func TestNewSaveData(t *testing.T) {
	data := NewSaveData()
	
	if data == nil {
		t.Fatal("NewSaveData returned nil")
	}
	
	if data.Version != "1.0.0" {
		t.Errorf("Expected version 1.0.0, got %s", data.Version)
	}
	if data.CurrentWorld != 1 {
		t.Errorf("Expected world 1, got %d", data.CurrentWorld)
	}
	if data.CurrentLevel != 1 {
		t.Errorf("Expected level 1, got %d", data.CurrentLevel)
	}
	if data.TotalLives != 5 {
		t.Errorf("Expected 5 lives, got %d", data.TotalLives)
	}
	if data.UnlockedLevels == nil {
		t.Error("UnlockedLevels map not initialized")
	}
	if data.LevelScores == nil {
		t.Error("LevelScores map not initialized")
	}
	if data.WorldScores == nil {
		t.Error("WorldScores map not initialized")
	}
	if data.Achievements == nil {
		t.Error("Achievements slice not initialized")
	}
}

func TestGetSaveDirectory(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	
	if sm.GetSaveDirectory() != tempDir {
		t.Errorf("Expected save directory %s, got %s", tempDir, sm.GetSaveDirectory())
	}
}

func TestMultipleSaveSlots(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	
	// Create different saves in each slot
	for slot := 1; slot <= 3; slot++ {
		saveData := NewSaveData()
		saveData.CurrentWorld = slot
		saveData.TotalScore = slot * 1000
		
		err := sm.SaveGame(slot, saveData)
		if err != nil {
			t.Fatalf("SaveGame(%d) failed: %v", slot, err)
		}
	}
	
	// Verify each slot has correct data
	for slot := 1; slot <= 3; slot++ {
		loadedData, err := sm.LoadSave(slot)
		if err != nil {
			t.Fatalf("LoadSave(%d) failed: %v", slot, err)
		}
		
		if loadedData.CurrentWorld != slot {
			t.Errorf("Slot %d: expected world %d, got %d", slot, slot, loadedData.CurrentWorld)
		}
		if loadedData.TotalScore != slot*1000 {
			t.Errorf("Slot %d: expected score %d, got %d", slot, slot*1000, loadedData.TotalScore)
		}
	}
}
