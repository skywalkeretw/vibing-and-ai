package game

import (
	"testing"
	"time"
)

// Mock implementations for testing

type mockGame struct {
	paused  bool
	resumed bool
}

func (m *mockGame) Pause() {
	m.paused = true
	m.resumed = false
}

func (m *mockGame) Resume() {
	m.resumed = true
	m.paused = false
}

type mockLevelManager struct {
	levelComplete  bool
	bossLevel      bool
	bossDefeated   bool
	hasTimeLimit   bool
	timeRemaining  float64
	nextLevel      string
	loadedLevel    string
	respawnedPlayer int
}

func (m *mockLevelManager) IsLevelComplete() bool {
	return m.levelComplete
}

func (m *mockLevelManager) IsBossLevel() bool {
	return m.bossLevel
}

func (m *mockLevelManager) IsBossDefeated() bool {
	return m.bossDefeated
}

func (m *mockLevelManager) HasTimeLimit() bool {
	return m.hasTimeLimit
}

func (m *mockLevelManager) GetTimeRemaining() float64 {
	return m.timeRemaining
}

func (m *mockLevelManager) GetNextLevel(currentLevel string) string {
	return m.nextLevel
}

func (m *mockLevelManager) LoadLevel(levelName string) error {
	m.loadedLevel = levelName
	return nil
}

func (m *mockLevelManager) RespawnPlayer(playerID int) error {
	m.respawnedPlayer = playerID
	return nil
}

type mockSaveManager struct {
	savedProgress   bool
	savedCompletion bool
}

func (m *mockSaveManager) SaveProgress(level string, lives map[int]int, scores map[int]int) error {
	m.savedProgress = true
	return nil
}

func (m *mockSaveManager) SaveGameCompletion(scores map[int]int) error {
	m.savedCompletion = true
	return nil
}

type mockUIManager struct {
	currentScreen string
	input         InputType
}

func (m *mockUIManager) ShowMainMenu() {
	m.currentScreen = "main_menu"
}

func (m *mockUIManager) HideMainMenu() {
	m.currentScreen = ""
}

func (m *mockUIManager) ShowHUD() {
	m.currentScreen = "hud"
}

func (m *mockUIManager) HideHUD() {
	m.currentScreen = ""
}

func (m *mockUIManager) ShowPauseMenu() {
	m.currentScreen = "pause_menu"
}

func (m *mockUIManager) HidePauseMenu() {
	m.currentScreen = ""
}

func (m *mockUIManager) ShowLevelTransition(levelName string) {
	m.currentScreen = "level_transition"
}

func (m *mockUIManager) ShowLevelCompleteScreen(scores map[int]int, timeRemaining float64) {
	m.currentScreen = "level_complete"
}

func (m *mockUIManager) HideLevelCompleteScreen() {
	m.currentScreen = ""
}

func (m *mockUIManager) ShowGameOverScreen(reason string, scores map[int]int) {
	m.currentScreen = "game_over"
}

func (m *mockUIManager) HideGameOverScreen() {
	m.currentScreen = ""
}

func (m *mockUIManager) ShowVictoryScreen(scores map[int]int) {
	m.currentScreen = "victory"
}

func (m *mockUIManager) HideVictoryScreen() {
	m.currentScreen = ""
}

func (m *mockUIManager) ShowSettings() {
	m.currentScreen = "settings"
}

func (m *mockUIManager) ShowLevelSelect() {
	m.currentScreen = "level_select"
}

func (m *mockUIManager) ShowHighScoreEntry(playerID int) {
	m.currentScreen = "high_score_entry"
}

func (m *mockUIManager) GetInput() InputType {
	return m.input
}

type mockScoreManager struct {
	scores       map[int]int
	playerScores map[int]PlayerScore
}

func (m *mockScoreManager) OnLevelComplete(playerID int, timeRemaining float64) {
	if m.scores == nil {
		m.scores = make(map[int]int)
	}
	m.scores[playerID] += 1000
}

func (m *mockScoreManager) GetScores() map[int]int {
	if m.scores == nil {
		return make(map[int]int)
	}
	return m.scores
}

func (m *mockScoreManager) GetPlayerScore(playerID int) PlayerScore {
	if m.playerScores == nil {
		return PlayerScore{PlayerID: playerID, TotalScore: 0}
	}
	return m.playerScores[playerID]
}

func (m *mockScoreManager) IsHighScore(score int) bool {
	return score > 10000
}

// Tests

func TestNewGameStateManager(t *testing.T) {
	gsm := NewGameStateManager()

	if gsm == nil {
		t.Fatal("NewGameStateManager returned nil")
	}

	if gsm.currentState != StateMainMenu {
		t.Errorf("Expected initial state to be StateMainMenu, got %v", gsm.currentState)
	}

	if gsm.livesRemaining == nil {
		t.Error("Expected livesRemaining to be initialized")
	}
}

func TestInitialize(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}

	gsm.Initialize(game)

	if gsm.currentState != StateMainMenu {
		t.Errorf("Expected state to be StateMainMenu after Initialize, got %v", gsm.currentState)
	}

	if gsm.livesRemaining[1] != 5 {
		t.Errorf("Expected player 1 to have 5 lives, got %d", gsm.livesRemaining[1])
	}

	if gsm.livesRemaining[2] != 5 {
		t.Errorf("Expected player 2 to have 5 lives, got %d", gsm.livesRemaining[2])
	}
}

func TestChangeState(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}
	ui := &mockUIManager{}

	gsm.Initialize(game)
	gsm.SetUIManager(ui)

	// Test transition to playing
	gsm.ChangeState(StatePlaying)

	if gsm.GetCurrentState() != StatePlaying {
		t.Errorf("Expected state to be StatePlaying, got %v", gsm.GetCurrentState())
	}

	if gsm.GetPreviousState() != StateMainMenu {
		t.Errorf("Expected previous state to be StateMainMenu, got %v", gsm.GetPreviousState())
	}

	if !game.resumed {
		t.Error("Expected game to be resumed")
	}

	if ui.currentScreen != "hud" {
		t.Errorf("Expected UI to show HUD, got %s", ui.currentScreen)
	}

	// Test transition to paused
	gsm.ChangeState(StatePaused)

	if gsm.GetCurrentState() != StatePaused {
		t.Errorf("Expected state to be StatePaused, got %v", gsm.GetCurrentState())
	}

	if !game.paused {
		t.Error("Expected game to be paused")
	}

	if ui.currentScreen != "pause_menu" {
		t.Errorf("Expected UI to show pause menu, got %s", ui.currentScreen)
	}
}

func TestChangeStateNoOp(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}

	gsm.Initialize(game)
	gsm.ChangeState(StatePlaying)

	previousTimer := gsm.GetStateTimer()
	time.Sleep(10 * time.Millisecond)

	// Try to change to same state
	gsm.ChangeState(StatePlaying)

	// Timer should not have reset
	if gsm.GetStateTimer() < previousTimer {
		t.Error("State timer should not reset when changing to same state")
	}
}

func TestStateCallbacks(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}

	gsm.Initialize(game)

	callbackCalled := false
	var oldState, newState GameState

	gsm.SetOnStateChange(func(old, new GameState) {
		callbackCalled = true
		oldState = old
		newState = new
	})

	gsm.ChangeState(StatePlaying)

	if !callbackCalled {
		t.Error("State change callback was not called")
	}

	if oldState != StateMainMenu {
		t.Errorf("Expected old state to be StateMainMenu, got %v", oldState)
	}

	if newState != StatePlaying {
		t.Errorf("Expected new state to be StatePlaying, got %v", newState)
	}
}

func TestCheckWinConditions(t *testing.T) {
	tests := []struct {
		name          string
		levelComplete bool
		bossLevel     bool
		bossDefeated  bool
		expectedState GameState
	}{
		{
			name:          "Level complete",
			levelComplete: true,
			expectedState: StateLevelComplete,
		},
		{
			name:          "Boss defeated",
			bossLevel:     true,
			bossDefeated:  true,
			expectedState: StateLevelComplete,
		},
		{
			name:          "No win condition",
			levelComplete: false,
			bossLevel:     false,
			expectedState: StatePlaying,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gsm := NewGameStateManager()
			game := &mockGame{}
			levelMgr := &mockLevelManager{
				levelComplete: tt.levelComplete,
				bossLevel:     tt.bossLevel,
				bossDefeated:  tt.bossDefeated,
			}

			gsm.Initialize(game)
			gsm.SetLevelManager(levelMgr)
			gsm.ChangeState(StatePlaying)

			gsm.Update(0.016) // One frame

			if gsm.GetCurrentState() != tt.expectedState {
				t.Errorf("Expected state %v, got %v", tt.expectedState, gsm.GetCurrentState())
			}
		})
	}
}

func TestCheckLoseConditions(t *testing.T) {
	tests := []struct {
		name          string
		lives         map[int]int
		hasTimeLimit  bool
		timeRemaining float64
		expectedState GameState
		expectedReason string
	}{
		{
			name:          "No lives",
			lives:         map[int]int{1: 0, 2: 0},
			expectedState: StateGameOver,
			expectedReason: "Out of Lives",
		},
		{
			name:          "Time expired",
			lives:         map[int]int{1: 3, 2: 3},
			hasTimeLimit:  true,
			timeRemaining: 0,
			expectedState: StateGameOver,
			expectedReason: "Time's Up!",
		},
		{
			name:          "Still playing",
			lives:         map[int]int{1: 3, 2: 3},
			hasTimeLimit:  true,
			timeRemaining: 100,
			expectedState: StatePlaying,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gsm := NewGameStateManager()
			game := &mockGame{}
			levelMgr := &mockLevelManager{
				hasTimeLimit:  tt.hasTimeLimit,
				timeRemaining: tt.timeRemaining,
			}

			gsm.Initialize(game)
			gsm.SetLevelManager(levelMgr)

			for playerID, lives := range tt.lives {
				gsm.SetLives(playerID, lives)
			}

			gsm.ChangeState(StatePlaying)
			gsm.Update(0.016) // One frame

			if gsm.GetCurrentState() != tt.expectedState {
				t.Errorf("Expected state %v, got %v", tt.expectedState, gsm.GetCurrentState())
			}

			if tt.expectedReason != "" && gsm.GetGameOverReason() != tt.expectedReason {
				t.Errorf("Expected game over reason %q, got %q", tt.expectedReason, gsm.GetGameOverReason())
			}
		})
	}
}

func TestOnPlayerDeath(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}
	levelMgr := &mockLevelManager{}

	gsm.Initialize(game)
	gsm.SetLevelManager(levelMgr)

	initialLives := gsm.GetLives(1)

	gsm.OnPlayerDeath(1)

	if gsm.GetLives(1) != initialLives-1 {
		t.Errorf("Expected lives to decrease by 1, got %d", gsm.GetLives(1))
	}

	if levelMgr.respawnedPlayer != 1 {
		t.Errorf("Expected player 1 to be respawned, got player %d", levelMgr.respawnedPlayer)
	}
}

func TestOnPlayerDeathGameOver(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}
	levelMgr := &mockLevelManager{}

	gsm.Initialize(game)
	gsm.SetLevelManager(levelMgr)
	gsm.SetLives(1, 1)
	gsm.SetLives(2, 0)
	gsm.ChangeState(StatePlaying)

	gsm.OnPlayerDeath(1)

	if gsm.GetCurrentState() != StateGameOver {
		t.Errorf("Expected state to be StateGameOver, got %v", gsm.GetCurrentState())
	}
}

func TestAddLife(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}

	gsm.Initialize(game)

	initialLives := gsm.GetLives(1)
	gsm.AddLife(1)

	if gsm.GetLives(1) != initialLives+1 {
		t.Errorf("Expected lives to increase by 1, got %d", gsm.GetLives(1))
	}
}

func TestLevelTransition(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}
	levelMgr := &mockLevelManager{
		nextLevel: "level-2",
	}
	ui := &mockUIManager{}

	gsm.Initialize(game)
	gsm.SetLevelManager(levelMgr)
	gsm.SetUIManager(ui)
	gsm.SetCurrentLevel("level-1")

	gsm.ChangeState(StateLevelTransition)

	if !gsm.IsTransitioning() {
		t.Error("Expected to be transitioning")
	}

	// Update for transition duration
	gsm.Update(2.1)

	if gsm.GetCurrentState() != StatePlaying {
		t.Errorf("Expected state to be StatePlaying after transition, got %v", gsm.GetCurrentState())
	}

	if !gsm.IsTransitioning() {
		// Transitioning flag should be cleared
	}
}

func TestHandleLevelComplete(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}
	levelMgr := &mockLevelManager{
		timeRemaining: 100,
		nextLevel:     "level-2",
	}
	saveMgr := &mockSaveManager{}
	scoreMgr := &mockScoreManager{
		scores: map[int]int{1: 5000},
	}
	ui := &mockUIManager{}

	gsm.Initialize(game)
	gsm.SetLevelManager(levelMgr)
	gsm.SetSaveManager(saveMgr)
	gsm.SetScoreManager(scoreMgr)
	gsm.SetUIManager(ui)
	gsm.SetCurrentLevel("level-1")

	gsm.ChangeState(StateLevelComplete)

	if !saveMgr.savedProgress {
		t.Error("Expected progress to be saved")
	}

	if ui.currentScreen != "level_complete" {
		t.Errorf("Expected UI to show level complete screen, got %s", ui.currentScreen)
	}
}

func TestHandleGameOver(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}
	scoreMgr := &mockScoreManager{
		scores: map[int]int{1: 5000},
	}
	ui := &mockUIManager{}

	gsm.Initialize(game)
	gsm.SetScoreManager(scoreMgr)
	gsm.SetUIManager(ui)

	gsm.ChangeState(StateGameOver)

	if ui.currentScreen != "game_over" {
		t.Errorf("Expected UI to show game over screen, got %s", ui.currentScreen)
	}

	// Test auto-return to main menu
	gsm.Update(3.1)

	if gsm.GetCurrentState() != StateMainMenu {
		t.Errorf("Expected state to be StateMainMenu after timeout, got %v", gsm.GetCurrentState())
	}
}

func TestHandleVictory(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}
	saveMgr := &mockSaveManager{}
	scoreMgr := &mockScoreManager{
		scores: map[int]int{1: 15000},
	}
	ui := &mockUIManager{}

	gsm.Initialize(game)
	gsm.SetSaveManager(saveMgr)
	gsm.SetScoreManager(scoreMgr)
	gsm.SetUIManager(ui)

	gsm.ChangeState(StateVictory)

	if !saveMgr.savedCompletion {
		t.Error("Expected game completion to be saved")
	}

	if ui.currentScreen != "victory" {
		t.Errorf("Expected UI to show victory screen, got %s", ui.currentScreen)
	}
}

func TestProceedToNextLevel(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}
	levelMgr := &mockLevelManager{
		nextLevel: "level-2",
	}
	ui := &mockUIManager{
		input: InputContinue,
	}

	gsm.Initialize(game)
	gsm.SetLevelManager(levelMgr)
	gsm.SetUIManager(ui)
	gsm.SetCurrentLevel("level-1")

	// Set next level
	gsm.nextLevel = "level-2"

	gsm.ChangeState(StateLevelComplete)
	gsm.Update(0.016)

	if gsm.GetCurrentState() != StateLevelTransition {
		t.Errorf("Expected state to be StateLevelTransition, got %v", gsm.GetCurrentState())
	}

	if gsm.GetCurrentLevel() != "level-2" {
		t.Errorf("Expected current level to be level-2, got %s", gsm.GetCurrentLevel())
	}
}

func TestProceedToNextLevelNoMore(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}
	levelMgr := &mockLevelManager{
		nextLevel: "",
	}
	ui := &mockUIManager{
		input: InputContinue,
	}

	gsm.Initialize(game)
	gsm.SetLevelManager(levelMgr)
	gsm.SetUIManager(ui)
	gsm.SetCurrentLevel("level-final")

	// No next level
	gsm.nextLevel = ""

	gsm.ChangeState(StateLevelComplete)
	gsm.Update(0.016)

	if gsm.GetCurrentState() != StateMainMenu {
		t.Errorf("Expected state to be StateMainMenu when no more levels, got %v", gsm.GetCurrentState())
	}
}

func TestStateQueries(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}

	gsm.Initialize(game)

	// Test IsPlaying
	gsm.ChangeState(StatePlaying)
	if !gsm.IsPlaying() {
		t.Error("Expected IsPlaying to return true")
	}

	// Test IsPaused
	gsm.ChangeState(StatePaused)
	if !gsm.IsPaused() {
		t.Error("Expected IsPaused to return true")
	}

	// Test IsTransitioning
	gsm.ChangeState(StateLevelTransition)
	if !gsm.IsTransitioning() {
		t.Error("Expected IsTransitioning to return true")
	}
}

func TestReset(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}

	gsm.Initialize(game)
	gsm.SetCurrentLevel("level-1")
	gsm.SetLives(1, 3)
	gsm.ChangeState(StatePlaying)

	gsm.Reset()

	if gsm.GetCurrentState() != StateMainMenu {
		t.Errorf("Expected state to be StateMainMenu after reset, got %v", gsm.GetCurrentState())
	}

	if gsm.GetCurrentLevel() != "" {
		t.Errorf("Expected current level to be empty after reset, got %s", gsm.GetCurrentLevel())
	}

	if gsm.GetLives(1) != 5 {
		t.Errorf("Expected lives to be reset to 5, got %d", gsm.GetLives(1))
	}
}

func TestGameStateString(t *testing.T) {
	tests := []struct {
		state    GameState
		expected string
	}{
		{StateMainMenu, "MainMenu"},
		{StatePlaying, "Playing"},
		{StatePaused, "Paused"},
		{StateLevelTransition, "LevelTransition"},
		{StateLevelComplete, "LevelComplete"},
		{StateGameOver, "GameOver"},
		{StateVictory, "Victory"},
		{StateSettings, "Settings"},
		{StateLevelSelect, "LevelSelect"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.state.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.state.String())
			}
		})
	}
}

func TestWinConditionString(t *testing.T) {
	tests := []struct {
		condition WinCondition
		expected  string
	}{
		{WinLevelComplete, "LevelComplete"},
		{WinBossDefeated, "BossDefeated"},
		{WinAllLevelsComplete, "AllLevelsComplete"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.condition.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.condition.String())
			}
		})
	}
}

func TestLoseConditionString(t *testing.T) {
	tests := []struct {
		condition LoseCondition
		expected  string
	}{
		{LoseNoLives, "NoLives"},
		{LoseTimeExpired, "TimeExpired"},
		{LoseBothPlayersDead, "BothPlayersDead"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.condition.String() != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, tt.condition.String())
			}
		})
	}
}

func TestConcurrentAccess(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}

	gsm.Initialize(game)

	done := make(chan bool)

	// Concurrent reads
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_ = gsm.GetCurrentState()
				_ = gsm.GetLives(1)
				_ = gsm.IsPlaying()
			}
			done <- true
		}()
	}

	// Concurrent writes
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				gsm.AddLife(1)
				gsm.Update(0.016)
			}
			done <- true
		}()
	}

	// Wait for all goroutines
	for i := 0; i < 20; i++ {
		<-done
	}
}

func TestAllLevelsComplete(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}
	levelMgr := &mockLevelManager{
		levelComplete: true,
		nextLevel:     "", // No next level
	}

	gsm.Initialize(game)
	gsm.SetLevelManager(levelMgr)
	gsm.SetCurrentLevel("final-level")
	gsm.ChangeState(StatePlaying)

	gsm.Update(0.016)

	// Should transition to victory when no next level
	if gsm.GetCurrentState() != StateLevelComplete {
		t.Errorf("Expected state to be StateLevelComplete, got %v", gsm.GetCurrentState())
	}
}

func TestMultipleStateTransitions(t *testing.T) {
	gsm := NewGameStateManager()
	game := &mockGame{}
	ui := &mockUIManager{}

	gsm.Initialize(game)
	gsm.SetUIManager(ui)

	states := []GameState{
		StatePlaying,
		StatePaused,
		StatePlaying,
		StateSettings,
		StateMainMenu,
		StateLevelSelect,
		StatePlaying,
	}

	for _, state := range states {
		gsm.ChangeState(state)
		if gsm.GetCurrentState() != state {
			t.Errorf("Expected state %v, got %v", state, gsm.GetCurrentState())
		}
	}
}
