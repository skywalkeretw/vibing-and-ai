package save

import (
	"testing"
	"time"
)

func TestUpdateHighScore(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Create initial save data
	saveData := NewSaveData()
	sm.saves[1] = saveData
	
	// Test setting first high score
	isNew := sm.UpdateHighScore("world1-level1", 1000)
	if !isNew {
		t.Error("Expected new high score for first score")
	}
	
	if sm.saves[1].LevelScores["world1-level1"] != 1000 {
		t.Errorf("Expected score 1000, got %d", sm.saves[1].LevelScores["world1-level1"])
	}
	
	// Test updating with higher score
	isNew = sm.UpdateHighScore("world1-level1", 2000)
	if !isNew {
		t.Error("Expected new high score for higher score")
	}
	
	if sm.saves[1].LevelScores["world1-level1"] != 2000 {
		t.Errorf("Expected score 2000, got %d", sm.saves[1].LevelScores["world1-level1"])
	}
	
	// Test updating with lower score (should not update)
	isNew = sm.UpdateHighScore("world1-level1", 1500)
	if isNew {
		t.Error("Expected no new high score for lower score")
	}
	
	if sm.saves[1].LevelScores["world1-level1"] != 2000 {
		t.Errorf("Expected score to remain 2000, got %d", sm.saves[1].LevelScores["world1-level1"])
	}
}

func TestUpdateHighScoreNoSaveData(t *testing.T) {
	sm := NewSaveManager()
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// No save data in slot
	isNew := sm.UpdateHighScore("world1-level1", 1000)
	if isNew {
		t.Error("Expected false when no save data exists")
	}
}

func TestGetHighScore(t *testing.T) {
	sm := NewSaveManager()
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Test with no save data
	score := sm.GetHighScore("world1-level1")
	if score != 0 {
		t.Errorf("Expected score 0 for no save data, got %d", score)
	}
	
	// Create save data with scores
	saveData := NewSaveData()
	saveData.LevelScores["world1-level1"] = 5000
	saveData.LevelScores["world1-level2"] = 3000
	sm.saves[1] = saveData
	
	// Test getting existing score
	score = sm.GetHighScore("world1-level1")
	if score != 5000 {
		t.Errorf("Expected score 5000, got %d", score)
	}
	
	// Test getting non-existent level
	score = sm.GetHighScore("world2-level1")
	if score != 0 {
		t.Errorf("Expected score 0 for non-existent level, got %d", score)
	}
}

func TestUpdateWorldScore(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Create initial save data
	saveData := NewSaveData()
	sm.saves[1] = saveData
	
	// Test setting first world score
	isNew := sm.UpdateWorldScore(1, 10000)
	if !isNew {
		t.Error("Expected new high score for first score")
	}
	
	if sm.saves[1].WorldScores[1] != 10000 {
		t.Errorf("Expected score 10000, got %d", sm.saves[1].WorldScores[1])
	}
	
	// Test updating with higher score
	isNew = sm.UpdateWorldScore(1, 15000)
	if !isNew {
		t.Error("Expected new high score for higher score")
	}
	
	// Test updating with lower score
	isNew = sm.UpdateWorldScore(1, 12000)
	if isNew {
		t.Error("Expected no new high score for lower score")
	}
}

func TestGetWorldScore(t *testing.T) {
	sm := NewSaveManager()
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Test with no save data
	score := sm.GetWorldScore(1)
	if score != 0 {
		t.Errorf("Expected score 0 for no save data, got %d", score)
	}
	
	// Create save data with scores
	saveData := NewSaveData()
	saveData.WorldScores[1] = 20000
	saveData.WorldScores[2] = 15000
	sm.saves[1] = saveData
	
	// Test getting existing score
	score = sm.GetWorldScore(1)
	if score != 20000 {
		t.Errorf("Expected score 20000, got %d", score)
	}
	
	// Test getting non-existent world
	score = sm.GetWorldScore(3)
	if score != 0 {
		t.Errorf("Expected score 0 for non-existent world, got %d", score)
	}
}

func TestGetTotalScore(t *testing.T) {
	sm := NewSaveManager()
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Test with no save data
	score := sm.GetTotalScore()
	if score != 0 {
		t.Errorf("Expected score 0 for no save data, got %d", score)
	}
	
	// Create save data with total score
	saveData := NewSaveData()
	saveData.TotalScore = 50000
	sm.saves[1] = saveData
	
	score = sm.GetTotalScore()
	if score != 50000 {
		t.Errorf("Expected score 50000, got %d", score)
	}
}

func TestUpdateTotalScore(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Create save data
	saveData := NewSaveData()
	sm.saves[1] = saveData
	
	err := sm.UpdateTotalScore(25000)
	if err != nil {
		t.Fatalf("UpdateTotalScore failed: %v", err)
	}
	
	if sm.saves[1].TotalScore != 25000 {
		t.Errorf("Expected total score 25000, got %d", sm.saves[1].TotalScore)
	}
}

func TestAddToTotalScore(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Create save data
	saveData := NewSaveData()
	saveData.TotalScore = 10000
	sm.saves[1] = saveData
	
	err := sm.AddToTotalScore(5000)
	if err != nil {
		t.Fatalf("AddToTotalScore failed: %v", err)
	}
	
	if sm.saves[1].TotalScore != 15000 {
		t.Errorf("Expected total score 15000, got %d", sm.saves[1].TotalScore)
	}
	
	// Add more points
	err = sm.AddToTotalScore(2000)
	if err != nil {
		t.Fatalf("AddToTotalScore failed: %v", err)
	}
	
	if sm.saves[1].TotalScore != 17000 {
		t.Errorf("Expected total score 17000, got %d", sm.saves[1].TotalScore)
	}
}

func TestGetTopLevelScores(t *testing.T) {
	sm := NewSaveManager()
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Create save data with multiple level scores
	saveData := NewSaveData()
	saveData.LevelScores = map[string]int{
		"world1-level1": 5000,
		"world1-level2": 8000,
		"world1-level3": 3000,
		"world2-level1": 10000,
		"world2-level2": 6000,
	}
	sm.saves[1] = saveData
	
	// Get top 3 scores
	topScores := sm.GetTopLevelScores(3)
	if len(topScores) != 3 {
		t.Errorf("Expected 3 scores, got %d", len(topScores))
	}
	
	// Verify they're sorted by score descending
	if topScores[0].Score != 10000 {
		t.Errorf("Expected first score 10000, got %d", topScores[0].Score)
	}
	if topScores[1].Score != 8000 {
		t.Errorf("Expected second score 8000, got %d", topScores[1].Score)
	}
	if topScores[2].Score != 6000 {
		t.Errorf("Expected third score 6000, got %d", topScores[2].Score)
	}
	
	// Test requesting more than available
	topScores = sm.GetTopLevelScores(10)
	if len(topScores) != 5 {
		t.Errorf("Expected 5 scores (all available), got %d", len(topScores))
	}
}

func TestGetTopWorldScores(t *testing.T) {
	sm := NewSaveManager()
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Create save data with world scores
	saveData := NewSaveData()
	saveData.WorldScores = map[int]int{
		1: 15000,
		2: 20000,
		3: 10000,
		4: 25000,
	}
	sm.saves[1] = saveData
	
	// Get top 2 scores
	topScores := sm.GetTopWorldScores(2)
	if len(topScores) != 2 {
		t.Errorf("Expected 2 scores, got %d", len(topScores))
	}
	
	// Verify they're sorted by score descending
	if topScores[0].Score != 25000 {
		t.Errorf("Expected first score 25000, got %d", topScores[0].Score)
	}
	if topScores[1].Score != 20000 {
		t.Errorf("Expected second score 20000, got %d", topScores[1].Score)
	}
}

func TestGetAllLevelScores(t *testing.T) {
	sm := NewSaveManager()
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Test with no save data
	scores := sm.GetAllLevelScores()
	if len(scores) != 0 {
		t.Errorf("Expected empty map for no save data, got %d entries", len(scores))
	}
	
	// Create save data with scores
	saveData := NewSaveData()
	saveData.LevelScores = map[string]int{
		"world1-level1": 5000,
		"world1-level2": 8000,
	}
	sm.saves[1] = saveData
	
	scores = sm.GetAllLevelScores()
	if len(scores) != 2 {
		t.Errorf("Expected 2 scores, got %d", len(scores))
	}
	
	if scores["world1-level1"] != 5000 {
		t.Errorf("Expected score 5000 for world1-level1, got %d", scores["world1-level1"])
	}
	
	// Verify it's a copy (modifying shouldn't affect original)
	scores["world1-level3"] = 9999
	if _, exists := sm.saves[1].LevelScores["world1-level3"]; exists {
		t.Error("Modifying returned map should not affect original")
	}
}

func TestGetAllWorldScores(t *testing.T) {
	sm := NewSaveManager()
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Test with no save data
	scores := sm.GetAllWorldScores()
	if len(scores) != 0 {
		t.Errorf("Expected empty map for no save data, got %d entries", len(scores))
	}
	
	// Create save data with scores
	saveData := NewSaveData()
	saveData.WorldScores = map[int]int{
		1: 15000,
		2: 20000,
	}
	sm.saves[1] = saveData
	
	scores = sm.GetAllWorldScores()
	if len(scores) != 2 {
		t.Errorf("Expected 2 scores, got %d", len(scores))
	}
	
	if scores[1] != 15000 {
		t.Errorf("Expected score 15000 for world 1, got %d", scores[1])
	}
	
	// Verify it's a copy
	scores[3] = 9999
	if _, exists := sm.saves[1].WorldScores[3]; exists {
		t.Error("Modifying returned map should not affect original")
	}
}

func TestClearLevelScore(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Create save data with scores
	saveData := NewSaveData()
	saveData.LevelScores = map[string]int{
		"world1-level1": 5000,
		"world1-level2": 8000,
	}
	sm.saves[1] = saveData
	
	err := sm.ClearLevelScore("world1-level1")
	if err != nil {
		t.Fatalf("ClearLevelScore failed: %v", err)
	}
	
	if _, exists := sm.saves[1].LevelScores["world1-level1"]; exists {
		t.Error("Level score was not cleared")
	}
	
	// Verify other scores remain
	if sm.saves[1].LevelScores["world1-level2"] != 8000 {
		t.Error("Other level scores should remain unchanged")
	}
}

func TestClearAllScores(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Create save data with scores
	saveData := NewSaveData()
	saveData.LevelScores = map[string]int{
		"world1-level1": 5000,
		"world1-level2": 8000,
	}
	saveData.WorldScores = map[int]int{
		1: 15000,
		2: 20000,
	}
	saveData.TotalScore = 50000
	sm.saves[1] = saveData
	
	err := sm.ClearAllScores()
	if err != nil {
		t.Fatalf("ClearAllScores failed: %v", err)
	}
	
	if len(sm.saves[1].LevelScores) != 0 {
		t.Error("Level scores were not cleared")
	}
	if len(sm.saves[1].WorldScores) != 0 {
		t.Error("World scores were not cleared")
	}
	if sm.saves[1].TotalScore != 0 {
		t.Error("Total score was not cleared")
	}
}

func TestHighScorePersistence(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	// Create save data and set high score
	saveData := NewSaveData()
	sm.saves[1] = saveData
	sm.UpdateHighScore("world1-level1", 5000)
	
	// Load the save and verify score persisted
	loadedData, err := sm.LoadSave(1)
	if err != nil {
		t.Fatalf("LoadSave failed: %v", err)
	}
	
	if loadedData.LevelScores["world1-level1"] != 5000 {
		t.Errorf("Expected persisted score 5000, got %d", loadedData.LevelScores["world1-level1"])
	}
}

func TestMultipleHighScores(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	saveData := NewSaveData()
	sm.saves[1] = saveData
	
	// Set multiple high scores
	levels := map[string]int{
		"world1-level1": 5000,
		"world1-level2": 8000,
		"world2-level1": 10000,
		"world2-level2": 6000,
	}
	
	for level, score := range levels {
		sm.UpdateHighScore(level, score)
	}
	
	// Verify all scores were set
	for level, expectedScore := range levels {
		actualScore := sm.GetHighScore(level)
		if actualScore != expectedScore {
			t.Errorf("Level %s: expected score %d, got %d", level, expectedScore, actualScore)
		}
	}
}

func TestScoreUpdateTiming(t *testing.T) {
	tempDir := t.TempDir()
	
	sm := NewSaveManager()
	sm.saveDir = tempDir
	sm.currentSlot = 1
	sm.saves = make(map[int]*SaveData)
	
	saveData := NewSaveData()
	sm.saves[1] = saveData
	
	// Record timestamp before update
	timeBefore := time.Now()
	
	sm.UpdateHighScore("world1-level1", 5000)
	
	// Load and check timestamp
	loadedData, err := sm.LoadSave(1)
	if err != nil {
		t.Fatalf("LoadSave failed: %v", err)
	}
	
	if loadedData.Timestamp.Before(timeBefore) {
		t.Error("Timestamp should be updated when saving")
	}
}
