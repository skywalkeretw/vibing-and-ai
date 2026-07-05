package score

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// ScoreAction defines different types of scoring actions
type ScoreAction int

const (
	ActionCoinCollect ScoreAction = iota
	ActionEnemyDefeat
	ActionPowerUpGet
	ActionLevelComplete
	ActionBossDefeat
	ActionSecretFound
	Action1UP
	ActionCheckpoint
)

// Point values for different actions
const (
	PointsCoin          = 10
	PointsEnemy         = 100
	PointsBoss          = 5000
	PointsPowerUp       = 50
	PointsSecretArea    = 500
	PointsLevelComplete = 1000
	Points1UP           = 200
	PointsCheckpoint    = 100
)

// PlayerScore tracks score for a single player
type PlayerScore struct {
	PlayerID      int
	CurrentScore  int
	LevelScore    int
	TotalScore    int
	Coins         int
	EnemiesKilled int
	TimeRemaining float64
	LastAction    ScoreAction
	ActionTime    float64
}

// HighScoreEntry represents a single high score entry
type HighScoreEntry struct {
	PlayerName string    `json:"player_name"`
	Score      int       `json:"score"`
	Level      string    `json:"level"`
	Date       time.Time `json:"date"`
}

// HighScoreTable manages high scores
type HighScoreTable struct {
	Entries    []HighScoreEntry `json:"entries"`
	MaxEntries int              `json:"max_entries"`
}

// ScoreManager manages scoring for all players
type ScoreManager struct {
	player1Score *PlayerScore
	player2Score *PlayerScore
	highScores   *HighScoreTable
	multiplier   float64
	comboTimer   float64
	comboCount   int
}

// NewScoreManager creates a new score manager
func NewScoreManager() *ScoreManager {
	return &ScoreManager{
		multiplier: 1.0,
	}
}

// Initialize initializes the score manager
func (sm *ScoreManager) Initialize(numPlayers int) {
	sm.player1Score = &PlayerScore{PlayerID: 1}

	if numPlayers == 2 {
		sm.player2Score = &PlayerScore{PlayerID: 2}
	}

	sm.highScores = LoadHighScores()
	sm.multiplier = 1.0
	sm.comboTimer = 0
	sm.comboCount = 0
}

// AddScore adds points to a player's score
func (sm *ScoreManager) AddScore(playerID int, points int, action ScoreAction) {
	score := sm.getPlayerScore(playerID)
	if score == nil {
		return
	}

	// Apply multiplier
	finalPoints := int(float64(points) * sm.multiplier)

	// Add to scores
	score.CurrentScore += finalPoints
	score.LevelScore += finalPoints
	score.TotalScore += finalPoints

	// Track action
	score.LastAction = action
	score.ActionTime = 0

	// Update specific counters
	switch action {
	case ActionCoinCollect:
		score.Coins++
	case ActionEnemyDefeat:
		score.EnemiesKilled++
	}

	// Update combo
	sm.updateCombo(action)

	// Check for extra life
	sm.checkExtraLife(score, finalPoints)
}

// updateCombo updates the combo system
func (sm *ScoreManager) updateCombo(action ScoreAction) {
	// Only certain actions contribute to combo
	if action == ActionEnemyDefeat || action == ActionCoinCollect {
		sm.comboCount++
		sm.comboTimer = 2.0 // 2 second combo window

		// Increase multiplier based on combo
		if sm.comboCount >= 10 {
			sm.multiplier = 2.0
		} else if sm.comboCount >= 5 {
			sm.multiplier = 1.5
		} else {
			sm.multiplier = 1.0
		}
	}
}

// resetCombo resets the combo system
func (sm *ScoreManager) resetCombo() {
	sm.comboCount = 0
	sm.multiplier = 1.0
}

// Update updates the score manager
func (sm *ScoreManager) Update(deltaTime float64) {
	// Update combo timer
	if sm.comboTimer > 0 {
		sm.comboTimer -= deltaTime
		if sm.comboTimer <= 0 {
			sm.resetCombo()
		}
	}

	// Update player score timers
	if sm.player1Score != nil {
		sm.player1Score.ActionTime += deltaTime
	}
	if sm.player2Score != nil {
		sm.player2Score.ActionTime += deltaTime
	}
}

// CalculateTimeBonus calculates time bonus for level completion
func (sm *ScoreManager) CalculateTimeBonus(playerID int, timeRemaining float64) int {
	// Time bonus: 10 points per second remaining
	bonus := int(timeRemaining) * 10

	score := sm.getPlayerScore(playerID)
	if score != nil {
		score.TimeRemaining = timeRemaining
	}

	return bonus
}

// OnLevelComplete handles level completion scoring
func (sm *ScoreManager) OnLevelComplete(playerID int, timeRemaining float64, levelName string) {
	score := sm.getPlayerScore(playerID)
	if score == nil {
		return
	}

	// Award level completion bonus
	sm.AddScore(playerID, PointsLevelComplete, ActionLevelComplete)

	// Calculate and award time bonus
	timeBonus := sm.CalculateTimeBonus(playerID, timeRemaining)
	sm.AddScore(playerID, timeBonus, ActionLevelComplete)

	// Check if new high score
	sm.checkHighScore(score, levelName)

	// Reset level score for next level
	score.LevelScore = 0
}

// checkHighScore checks if score qualifies as high score
func (sm *ScoreManager) checkHighScore(score *PlayerScore, levelName string) bool {
	if sm.highScores.IsHighScore(score.TotalScore) {
		entry := HighScoreEntry{
			PlayerName: fmt.Sprintf("Player %d", score.PlayerID),
			Score:      score.TotalScore,
			Level:      levelName,
			Date:       time.Now(),
		}
		sm.highScores.AddEntry(entry)
		sm.highScores.Save()
		return true
	}
	return false
}

// checkExtraLife checks if player earned an extra life
func (sm *ScoreManager) checkExtraLife(score *PlayerScore, pointsAdded int) bool {
	// Award extra life every 10000 points
	threshold := 10000
	previousThreshold := (score.TotalScore - pointsAdded) / threshold
	currentThreshold := score.TotalScore / threshold

	return currentThreshold > previousThreshold
}

// getPlayerScore returns the score for a player
func (sm *ScoreManager) getPlayerScore(playerID int) *PlayerScore {
	if playerID == 1 {
		return sm.player1Score
	} else if playerID == 2 {
		return sm.player2Score
	}
	return nil
}

// GetPlayerScore returns the score for a player (public)
func (sm *ScoreManager) GetPlayerScore(playerID int) *PlayerScore {
	return sm.getPlayerScore(playerID)
}

// GetMultiplier returns the current combo multiplier
func (sm *ScoreManager) GetMultiplier() float64 {
	return sm.multiplier
}

// GetComboCount returns the current combo count
func (sm *ScoreManager) GetComboCount() int {
	return sm.comboCount
}

// GetHighScores returns the high score table
func (sm *ScoreManager) GetHighScores() *HighScoreTable {
	return sm.highScores
}

// ResetLevelScore resets the level score for a player
func (sm *ScoreManager) ResetLevelScore(playerID int) {
	score := sm.getPlayerScore(playerID)
	if score != nil {
		score.LevelScore = 0
	}
}

// ResetAllScores resets all scores
func (sm *ScoreManager) ResetAllScores() {
	if sm.player1Score != nil {
		sm.player1Score.CurrentScore = 0
		sm.player1Score.LevelScore = 0
		sm.player1Score.TotalScore = 0
		sm.player1Score.Coins = 0
		sm.player1Score.EnemiesKilled = 0
	}
	if sm.player2Score != nil {
		sm.player2Score.CurrentScore = 0
		sm.player2Score.LevelScore = 0
		sm.player2Score.TotalScore = 0
		sm.player2Score.Coins = 0
		sm.player2Score.EnemiesKilled = 0
	}
	sm.resetCombo()
}

// Initialize initializes the high score table
func (hst *HighScoreTable) Initialize() {
	hst.MaxEntries = 10
	hst.Entries = make([]HighScoreEntry, 0, hst.MaxEntries)
}

// IsHighScore checks if a score qualifies as a high score
func (hst *HighScoreTable) IsHighScore(score int) bool {
	if len(hst.Entries) < hst.MaxEntries {
		return true
	}
	return score > hst.Entries[len(hst.Entries)-1].Score
}

// AddEntry adds a new high score entry
func (hst *HighScoreTable) AddEntry(entry HighScoreEntry) {
	hst.Entries = append(hst.Entries, entry)

	// Sort by score (descending)
	sort.Slice(hst.Entries, func(i, j int) bool {
		return hst.Entries[i].Score > hst.Entries[j].Score
	})

	// Keep only top entries
	if len(hst.Entries) > hst.MaxEntries {
		hst.Entries = hst.Entries[:hst.MaxEntries]
	}
}

// GetEntries returns all high score entries
func (hst *HighScoreTable) GetEntries() []HighScoreEntry {
	return hst.Entries
}

// GetRank returns the rank of a score (1-based, 0 if not in table)
func (hst *HighScoreTable) GetRank(score int) int {
	for i, entry := range hst.Entries {
		if entry.Score == score {
			return i + 1
		}
	}
	return 0
}

// Save saves the high score table to disk
func (hst *HighScoreTable) Save() error {
	data, err := json.MarshalIndent(hst, "", "  ")
	if err != nil {
		return err
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	saveDir := filepath.Join(homeDir, ".go-terminal-platformer")
	if err := os.MkdirAll(saveDir, 0755); err != nil {
		return err
	}

	savePath := filepath.Join(saveDir, "highscores.json")
	return os.WriteFile(savePath, data, 0644)
}

// LoadHighScores loads the high score table from disk
func LoadHighScores() *HighScoreTable {
	hst := &HighScoreTable{}
	hst.Initialize()

	homeDir, err := os.UserHomeDir()
	if err != nil {
		return hst
	}

	savePath := filepath.Join(homeDir, ".go-terminal-platformer", "highscores.json")
	data, err := os.ReadFile(savePath)
	if err != nil {
		return hst
	}

	json.Unmarshal(data, hst)
	return hst
}

// Clear clears all high scores
func (hst *HighScoreTable) Clear() {
	hst.Entries = make([]HighScoreEntry, 0, hst.MaxEntries)
}
