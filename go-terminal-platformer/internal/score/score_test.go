package score

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewScoreManager(t *testing.T) {
	sm := NewScoreManager()

	if sm == nil {
		t.Fatal("NewScoreManager returned nil")
	}

	if sm.multiplier != 1.0 {
		t.Errorf("Expected multiplier 1.0, got %f", sm.multiplier)
	}
}

func TestScoreManager_Initialize_SinglePlayer(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	if sm.player1Score == nil {
		t.Error("Player 1 score should be initialized")
	}

	if sm.player2Score != nil {
		t.Error("Player 2 score should not be initialized for single player")
	}

	if sm.highScores == nil {
		t.Error("High scores should be initialized")
	}
}

func TestScoreManager_Initialize_TwoPlayers(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(2)

	if sm.player1Score == nil {
		t.Error("Player 1 score should be initialized")
	}

	if sm.player2Score == nil {
		t.Error("Player 2 score should be initialized for two players")
	}
}

func TestScoreManager_AddScore_Basic(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	sm.AddScore(1, PointsCoin, ActionCoinCollect)

	score := sm.GetPlayerScore(1)
	if score.TotalScore != PointsCoin {
		t.Errorf("Expected total score %d, got %d", PointsCoin, score.TotalScore)
	}

	if score.Coins != 1 {
		t.Error("Coin count should be 1")
	}
}

func TestScoreManager_AddScore_WithMultiplier(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	// Build combo to get multiplier
	for i := 0; i < 5; i++ {
		sm.AddScore(1, PointsEnemy, ActionEnemyDefeat)
	}

	// Multiplier should be 1.5x at 5 combo
	if sm.multiplier != 1.5 {
		t.Errorf("Expected multiplier 1.5, got %f", sm.multiplier)
	}

	score := sm.GetPlayerScore(1)
	// All 5 at 1.0x (multiplier increases after combo count reaches threshold)
	expectedScore := 5 * PointsEnemy
	if score.TotalScore != expectedScore {
		t.Errorf("Expected total score %d, got %d", expectedScore, score.TotalScore)
	}

	// Add one more to see multiplier in effect
	sm.AddScore(1, PointsEnemy, ActionEnemyDefeat)
	expectedScore = (5 * PointsEnemy) + int(float64(PointsEnemy)*1.5)
	if score.TotalScore != expectedScore {
		t.Errorf("Expected total score with multiplier %d, got %d", expectedScore, score.TotalScore)
	}
}

func TestScoreManager_AddScore_EnemyKillTracking(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	sm.AddScore(1, PointsEnemy, ActionEnemyDefeat)
	sm.AddScore(1, PointsEnemy, ActionEnemyDefeat)

	score := sm.GetPlayerScore(1)
	if score.EnemiesKilled != 2 {
		t.Errorf("Expected 2 enemies killed, got %d", score.EnemiesKilled)
	}
}

func TestScoreManager_ComboSystem(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	// Test combo building
	sm.AddScore(1, PointsEnemy, ActionEnemyDefeat)
	if sm.comboCount != 1 {
		t.Errorf("Expected combo count 1, got %d", sm.comboCount)
	}

	// Build to 5 combo
	for i := 0; i < 4; i++ {
		sm.AddScore(1, PointsEnemy, ActionEnemyDefeat)
	}

	if sm.comboCount != 5 {
		t.Errorf("Expected combo count 5, got %d", sm.comboCount)
	}

	if sm.multiplier != 1.5 {
		t.Errorf("Expected multiplier 1.5 at 5 combo, got %f", sm.multiplier)
	}

	// Build to 10 combo
	for i := 0; i < 5; i++ {
		sm.AddScore(1, PointsEnemy, ActionEnemyDefeat)
	}

	if sm.comboCount != 10 {
		t.Errorf("Expected combo count 10, got %d", sm.comboCount)
	}

	if sm.multiplier != 2.0 {
		t.Errorf("Expected multiplier 2.0 at 10 combo, got %f", sm.multiplier)
	}
}

func TestScoreManager_ComboReset(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	// Build combo
	sm.AddScore(1, PointsEnemy, ActionEnemyDefeat)
	sm.AddScore(1, PointsEnemy, ActionEnemyDefeat)

	if sm.comboCount != 2 {
		t.Error("Combo should be 2")
	}

	// Simulate time passing
	sm.Update(2.5) // More than 2 second window

	if sm.comboCount != 0 {
		t.Errorf("Combo should reset after timeout, got %d", sm.comboCount)
	}

	if sm.multiplier != 1.0 {
		t.Errorf("Multiplier should reset to 1.0, got %f", sm.multiplier)
	}
}

func TestScoreManager_Update(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	sm.AddScore(1, PointsEnemy, ActionEnemyDefeat)
	initialTimer := sm.comboTimer

	sm.Update(0.5)

	if sm.comboTimer >= initialTimer {
		t.Error("Combo timer should decrease")
	}

	score := sm.GetPlayerScore(1)
	if score.ActionTime <= 0 {
		t.Error("Action time should increase")
	}
}

func TestScoreManager_CalculateTimeBonus(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	timeRemaining := 45.5
	bonus := sm.CalculateTimeBonus(1, timeRemaining)

	expectedBonus := 450 // 45 * 10
	if bonus != expectedBonus {
		t.Errorf("Expected time bonus %d, got %d", expectedBonus, bonus)
	}

	score := sm.GetPlayerScore(1)
	if score.TimeRemaining != timeRemaining {
		t.Errorf("Expected time remaining %f, got %f", timeRemaining, score.TimeRemaining)
	}
}

func TestScoreManager_OnLevelComplete(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	initialScore := sm.GetPlayerScore(1).TotalScore
	timeRemaining := 30.0

	sm.OnLevelComplete(1, timeRemaining, "World 1-1")

	score := sm.GetPlayerScore(1)
	expectedScore := initialScore + PointsLevelComplete + (30 * 10)

	if score.TotalScore != expectedScore {
		t.Errorf("Expected total score %d, got %d", expectedScore, score.TotalScore)
	}

	if score.LevelScore != 0 {
		t.Error("Level score should be reset after completion")
	}
}

func TestScoreManager_CheckExtraLife(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	score := sm.GetPlayerScore(1)

	// Add 9999 points (just below threshold)
	sm.AddScore(1, 9999, ActionEnemyDefeat)
	extraLife := sm.checkExtraLife(score, 9999)
	if extraLife {
		t.Error("Should not award extra life below 10000")
	}

	// Add 1 more point to cross threshold
	sm.AddScore(1, 1, ActionCoinCollect)
	extraLife = sm.checkExtraLife(score, 1)
	if !extraLife {
		t.Error("Should award extra life at 10000")
	}

	// Add another 10000 points
	sm.AddScore(1, 10000, ActionBossDefeat)
	extraLife = sm.checkExtraLife(score, 10000)
	if !extraLife {
		t.Error("Should award extra life at 20000")
	}
}

func TestScoreManager_GetPlayerScore(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(2)

	score1 := sm.GetPlayerScore(1)
	score2 := sm.GetPlayerScore(2)
	score3 := sm.GetPlayerScore(3)

	if score1 == nil {
		t.Error("Should return player 1 score")
	}

	if score2 == nil {
		t.Error("Should return player 2 score")
	}

	if score3 != nil {
		t.Error("Should return nil for invalid player ID")
	}
}

func TestScoreManager_ResetLevelScore(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	sm.AddScore(1, 500, ActionEnemyDefeat)
	score := sm.GetPlayerScore(1)

	if score.LevelScore != 500 {
		t.Error("Level score should be 500")
	}

	sm.ResetLevelScore(1)

	if score.LevelScore != 0 {
		t.Error("Level score should be reset to 0")
	}

	if score.TotalScore != 500 {
		t.Error("Total score should remain unchanged")
	}
}

func TestScoreManager_ResetAllScores(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(2)

	sm.AddScore(1, 1000, ActionEnemyDefeat)
	sm.AddScore(2, 2000, ActionEnemyDefeat)

	sm.ResetAllScores()

	score1 := sm.GetPlayerScore(1)
	score2 := sm.GetPlayerScore(2)

	if score1.TotalScore != 0 || score1.LevelScore != 0 {
		t.Error("Player 1 scores should be reset")
	}

	if score2.TotalScore != 0 || score2.LevelScore != 0 {
		t.Error("Player 2 scores should be reset")
	}

	if sm.comboCount != 0 || sm.multiplier != 1.0 {
		t.Error("Combo should be reset")
	}
}

func TestHighScoreTable_Initialize(t *testing.T) {
	hst := &HighScoreTable{}
	hst.Initialize()

	if hst.MaxEntries != 10 {
		t.Errorf("Expected max entries 10, got %d", hst.MaxEntries)
	}

	if len(hst.Entries) != 0 {
		t.Error("Entries should be empty initially")
	}
}

func TestHighScoreTable_IsHighScore(t *testing.T) {
	hst := &HighScoreTable{}
	hst.Initialize()

	// Empty table - any score is high score
	if !hst.IsHighScore(100) {
		t.Error("Any score should be high score in empty table")
	}

	// Fill table with 10 entries
	for i := 0; i < 10; i++ {
		hst.AddEntry(HighScoreEntry{
			PlayerName: "Test",
			Score:      (i + 1) * 100,
			Level:      "Test",
			Date:       time.Now(),
		})
	}

	// Score higher than lowest should be high score
	if !hst.IsHighScore(150) {
		t.Error("Score 150 should be high score (lowest is 100)")
	}

	// Score lower than lowest should not be high score
	if hst.IsHighScore(50) {
		t.Error("Score 50 should not be high score")
	}
}

func TestHighScoreTable_AddEntry(t *testing.T) {
	hst := &HighScoreTable{}
	hst.Initialize()

	entry := HighScoreEntry{
		PlayerName: "Player 1",
		Score:      5000,
		Level:      "World 1-1",
		Date:       time.Now(),
	}

	hst.AddEntry(entry)

	if len(hst.Entries) != 1 {
		t.Errorf("Expected 1 entry, got %d", len(hst.Entries))
	}

	if hst.Entries[0].Score != 5000 {
		t.Error("Entry not added correctly")
	}
}

func TestHighScoreTable_AddEntry_Sorting(t *testing.T) {
	hst := &HighScoreTable{}
	hst.Initialize()

	// Add entries in random order
	scores := []int{300, 100, 500, 200, 400}
	for _, score := range scores {
		hst.AddEntry(HighScoreEntry{
			PlayerName: "Test",
			Score:      score,
			Level:      "Test",
			Date:       time.Now(),
		})
	}

	// Verify sorted in descending order
	for i := 0; i < len(hst.Entries)-1; i++ {
		if hst.Entries[i].Score < hst.Entries[i+1].Score {
			t.Error("Entries not sorted correctly")
		}
	}

	if hst.Entries[0].Score != 500 {
		t.Error("Highest score should be first")
	}
}

func TestHighScoreTable_AddEntry_MaxEntries(t *testing.T) {
	hst := &HighScoreTable{}
	hst.Initialize()

	// Add 15 entries (more than max)
	for i := 0; i < 15; i++ {
		hst.AddEntry(HighScoreEntry{
			PlayerName: "Test",
			Score:      (i + 1) * 100,
			Level:      "Test",
			Date:       time.Now(),
		})
	}

	if len(hst.Entries) != 10 {
		t.Errorf("Expected 10 entries (max), got %d", len(hst.Entries))
	}

	// Verify lowest scores were removed
	if hst.Entries[len(hst.Entries)-1].Score != 600 {
		t.Error("Lowest scores should be removed")
	}
}

func TestHighScoreTable_GetRank(t *testing.T) {
	hst := &HighScoreTable{}
	hst.Initialize()

	scores := []int{500, 400, 300, 200, 100}
	for _, score := range scores {
		hst.AddEntry(HighScoreEntry{
			PlayerName: "Test",
			Score:      score,
			Level:      "Test",
			Date:       time.Now(),
		})
	}

	rank := hst.GetRank(500)
	if rank != 1 {
		t.Errorf("Expected rank 1 for score 500, got %d", rank)
	}

	rank = hst.GetRank(300)
	if rank != 3 {
		t.Errorf("Expected rank 3 for score 300, got %d", rank)
	}

	rank = hst.GetRank(999)
	if rank != 0 {
		t.Errorf("Expected rank 0 for non-existent score, got %d", rank)
	}
}

func TestHighScoreTable_Clear(t *testing.T) {
	hst := &HighScoreTable{}
	hst.Initialize()

	hst.AddEntry(HighScoreEntry{
		PlayerName: "Test",
		Score:      1000,
		Level:      "Test",
		Date:       time.Now(),
	})

	hst.Clear()

	if len(hst.Entries) != 0 {
		t.Error("Entries should be cleared")
	}
}

func TestHighScoreTable_SaveAndLoad(t *testing.T) {
	// Create temp directory for test
	homeDir, _ := os.UserHomeDir()
	testDir := filepath.Join(homeDir, ".go-terminal-platformer-test")
	os.Setenv("HOME", testDir)
	defer func() {
		os.RemoveAll(testDir)
		os.Unsetenv("HOME")
	}()

	hst := &HighScoreTable{}
	hst.Initialize()

	// Add test entries
	hst.AddEntry(HighScoreEntry{
		PlayerName: "Player 1",
		Score:      5000,
		Level:      "World 1-1",
		Date:       time.Now(),
	})
	hst.AddEntry(HighScoreEntry{
		PlayerName: "Player 2",
		Score:      3000,
		Level:      "World 1-2",
		Date:       time.Now(),
	})

	// Save
	err := hst.Save()
	if err != nil {
		t.Fatalf("Failed to save high scores: %v", err)
	}

	// Load
	loaded := LoadHighScores()

	if len(loaded.Entries) != 2 {
		t.Errorf("Expected 2 entries after load, got %d", len(loaded.Entries))
	}

	if loaded.Entries[0].Score != 5000 {
		t.Error("Loaded scores don't match saved scores")
	}
}

func TestScoreManager_TwoPlayerIndependence(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(2)

	sm.AddScore(1, 1000, ActionEnemyDefeat)
	sm.AddScore(2, 2000, ActionEnemyDefeat)

	score1 := sm.GetPlayerScore(1)
	score2 := sm.GetPlayerScore(2)

	if score1.TotalScore != 1000 {
		t.Error("Player 1 score incorrect")
	}

	if score2.TotalScore != 2000 {
		t.Error("Player 2 score incorrect")
	}

	// Scores should be independent
	if score1.TotalScore == score2.TotalScore {
		t.Error("Player scores should be independent")
	}
}

func TestScoreAction_AllTypes(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	actions := []struct {
		points int
		action ScoreAction
	}{
		{PointsCoin, ActionCoinCollect},
		{PointsEnemy, ActionEnemyDefeat},
		{PointsPowerUp, ActionPowerUpGet},
		{PointsLevelComplete, ActionLevelComplete},
		{PointsBoss, ActionBossDefeat},
		{PointsSecretArea, ActionSecretFound},
		{PointsCheckpoint, ActionCheckpoint},
	}

	for _, a := range actions {
		sm.AddScore(1, a.points, a.action)
	}

	score := sm.GetPlayerScore(1)
	expectedTotal := PointsCoin + PointsEnemy + PointsPowerUp + PointsLevelComplete +
		PointsBoss + PointsSecretArea + PointsCheckpoint

	if score.TotalScore < expectedTotal {
		t.Errorf("Expected at least %d total score, got %d", expectedTotal, score.TotalScore)
	}
}

func TestScoreManager_GetMultiplier(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	if sm.GetMultiplier() != 1.0 {
		t.Error("Initial multiplier should be 1.0")
	}

	// Build combo
	for i := 0; i < 5; i++ {
		sm.AddScore(1, PointsEnemy, ActionEnemyDefeat)
	}

	if sm.GetMultiplier() != 1.5 {
		t.Errorf("Expected multiplier 1.5, got %f", sm.GetMultiplier())
	}
}

func TestScoreManager_GetComboCount(t *testing.T) {
	sm := NewScoreManager()
	sm.Initialize(1)

	if sm.GetComboCount() != 0 {
		t.Error("Initial combo count should be 0")
	}

	sm.AddScore(1, PointsEnemy, ActionEnemyDefeat)
	sm.AddScore(1, PointsEnemy, ActionEnemyDefeat)

	if sm.GetComboCount() != 2 {
		t.Errorf("Expected combo count 2, got %d", sm.GetComboCount())
	}
}
