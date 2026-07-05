package game

import (
	"log"
	"sync"
)

// GameState represents the current state of the game
type GameState int

const (
	StateMainMenu GameState = iota
	StatePlaying
	StatePaused
	StateLevelTransition
	StateLevelComplete
	StateGameOver
	StateVictory
	StateSettings
	StateLevelSelect
)

// String returns the string representation of a GameState
func (gs GameState) String() string {
	switch gs {
	case StateMainMenu:
		return "MainMenu"
	case StatePlaying:
		return "Playing"
	case StatePaused:
		return "Paused"
	case StateLevelTransition:
		return "LevelTransition"
	case StateLevelComplete:
		return "LevelComplete"
	case StateGameOver:
		return "GameOver"
	case StateVictory:
		return "Victory"
	case StateSettings:
		return "Settings"
	case StateLevelSelect:
		return "LevelSelect"
	default:
		return "Unknown"
	}
}

// WinCondition represents different ways to win
type WinCondition int

const (
	WinLevelComplete WinCondition = iota
	WinBossDefeated
	WinAllLevelsComplete
)

// String returns the string representation of a WinCondition
func (wc WinCondition) String() string {
	switch wc {
	case WinLevelComplete:
		return "LevelComplete"
	case WinBossDefeated:
		return "BossDefeated"
	case WinAllLevelsComplete:
		return "AllLevelsComplete"
	default:
		return "Unknown"
	}
}

// LoseCondition represents different ways to lose
type LoseCondition int

const (
	LoseNoLives LoseCondition = iota
	LoseTimeExpired
	LoseBothPlayersDead
)

// String returns the string representation of a LoseCondition
func (lc LoseCondition) String() string {
	switch lc {
	case LoseNoLives:
		return "NoLives"
	case LoseTimeExpired:
		return "TimeExpired"
	case LoseBothPlayersDead:
		return "BothPlayersDead"
	default:
		return "Unknown"
	}
}

// LevelManager interface for level operations
type LevelManager interface {
	IsLevelComplete() bool
	IsBossLevel() bool
	IsBossDefeated() bool
	HasTimeLimit() bool
	GetTimeRemaining() float64
	GetNextLevel(currentLevel string) string
	LoadLevel(levelName string) error
	RespawnPlayer(playerID int) error
}

// SaveManager interface for save operations
type SaveManager interface {
	SaveProgress(level string, lives map[int]int, scores map[int]int) error
	SaveGameCompletion(scores map[int]int) error
}

// UIManager interface for UI operations
type UIManager interface {
	ShowMainMenu()
	HideMainMenu()
	ShowHUD()
	HideHUD()
	ShowPauseMenu()
	HidePauseMenu()
	ShowLevelTransition(levelName string)
	ShowLevelCompleteScreen(scores map[int]int, timeRemaining float64)
	HideLevelCompleteScreen()
	ShowGameOverScreen(reason string, scores map[int]int)
	HideGameOverScreen()
	ShowVictoryScreen(scores map[int]int)
	HideVictoryScreen()
	ShowSettings()
	ShowLevelSelect()
	ShowHighScoreEntry(playerID int)
	GetInput() InputType
}

// ScoreManager interface for score operations
type ScoreManager interface {
	OnLevelComplete(playerID int, timeRemaining float64)
	GetScores() map[int]int
	GetPlayerScore(playerID int) PlayerScore
	IsHighScore(score int) bool
}

// PlayerScore represents a player's score details
type PlayerScore struct {
	PlayerID   int
	TotalScore int
	Combo      int
}

// InputType represents different input types
type InputType int

const (
	InputNone InputType = iota
	InputContinue
	InputRetry
	InputQuit
)

// GameInterface represents the game operations needed by state manager
type GameInterface interface {
	Pause()
	Resume()
}

// GameStateManager manages the overall game state and flow
type GameStateManager struct {
	mu              sync.RWMutex
	currentState    GameState
	previousState   GameState
	stateTimer      float64
	transitionTimer float64
	transitioning   bool

	// Systems
	game         GameInterface
	levelManager LevelManager
	saveManager  SaveManager
	uiManager    UIManager
	scoreManager ScoreManager

	// State data
	currentLevel   string
	nextLevel      string
	livesRemaining map[int]int
	gameOverReason string

	// Callbacks
	onStateChange func(oldState, newState GameState)
}

// NewGameStateManager creates a new game state manager
func NewGameStateManager() *GameStateManager {
	return &GameStateManager{
		currentState:   StateMainMenu,
		previousState:  StateMainMenu,
		livesRemaining: make(map[int]int),
	}
}

// Initialize sets up the game state manager with required systems
func (gsm *GameStateManager) Initialize(game GameInterface) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()

	gsm.game = game
	gsm.currentState = StateMainMenu
	gsm.previousState = StateMainMenu
	gsm.livesRemaining = make(map[int]int)
	gsm.livesRemaining[1] = 5 // Default lives for player 1
	gsm.livesRemaining[2] = 5 // Default lives for player 2
	gsm.stateTimer = 0
	gsm.transitionTimer = 0
	gsm.transitioning = false
}

// SetLevelManager sets the level manager
func (gsm *GameStateManager) SetLevelManager(lm LevelManager) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()
	gsm.levelManager = lm
}

// SetSaveManager sets the save manager
func (gsm *GameStateManager) SetSaveManager(sm SaveManager) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()
	gsm.saveManager = sm
}

// SetUIManager sets the UI manager
func (gsm *GameStateManager) SetUIManager(um UIManager) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()
	gsm.uiManager = um
}

// SetScoreManager sets the score manager
func (gsm *GameStateManager) SetScoreManager(sm ScoreManager) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()
	gsm.scoreManager = sm
}

// SetOnStateChange sets a callback for state changes
func (gsm *GameStateManager) SetOnStateChange(callback func(oldState, newState GameState)) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()
	gsm.onStateChange = callback
}

// ChangeState transitions to a new game state
func (gsm *GameStateManager) ChangeState(newState GameState) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()

	if gsm.currentState == newState {
		return
	}

	// Exit current state
	gsm.onStateExit(gsm.currentState)

	// Update states
	oldState := gsm.currentState
	gsm.previousState = gsm.currentState
	gsm.currentState = newState
	gsm.stateTimer = 0

	// Enter new state
	gsm.onStateEnter(newState)

	log.Printf("State changed: %v -> %v", oldState, gsm.currentState)

	// Call callback if set
	if gsm.onStateChange != nil {
		gsm.onStateChange(oldState, newState)
	}
}

// onStateEnter handles entering a new state
func (gsm *GameStateManager) onStateEnter(state GameState) {
	switch state {
	case StateMainMenu:
		if gsm.uiManager != nil {
			gsm.uiManager.ShowMainMenu()
		}
		// audio.PlayMusic("main_menu")

	case StatePlaying:
		if gsm.uiManager != nil {
			gsm.uiManager.ShowHUD()
		}
		if gsm.game != nil {
			gsm.game.Resume()
		}
		// audio.PlayMusic("level_music")

	case StatePaused:
		if gsm.uiManager != nil {
			gsm.uiManager.ShowPauseMenu()
		}
		if gsm.game != nil {
			gsm.game.Pause()
		}

	case StateLevelTransition:
		gsm.transitioning = true
		gsm.transitionTimer = 2.0 // 2 second transition
		if gsm.uiManager != nil {
			gsm.uiManager.ShowLevelTransition(gsm.nextLevel)
		}

	case StateLevelComplete:
		gsm.handleLevelComplete()

	case StateGameOver:
		gsm.handleGameOver()

	case StateVictory:
		gsm.handleVictory()

	case StateSettings:
		if gsm.uiManager != nil {
			gsm.uiManager.ShowSettings()
		}

	case StateLevelSelect:
		if gsm.uiManager != nil {
			gsm.uiManager.ShowLevelSelect()
		}
	}
}

// onStateExit handles exiting a state
func (gsm *GameStateManager) onStateExit(state GameState) {
	switch state {
	case StateMainMenu:
		if gsm.uiManager != nil {
			gsm.uiManager.HideMainMenu()
		}

	case StatePlaying:
		if gsm.uiManager != nil {
			gsm.uiManager.HideHUD()
		}

	case StatePaused:
		if gsm.uiManager != nil {
			gsm.uiManager.HidePauseMenu()
		}
		if gsm.game != nil {
			gsm.game.Resume()
		}

	case StateLevelTransition:
		gsm.transitioning = false

	case StateLevelComplete:
		if gsm.uiManager != nil {
			gsm.uiManager.HideLevelCompleteScreen()
		}

	case StateGameOver:
		if gsm.uiManager != nil {
			gsm.uiManager.HideGameOverScreen()
		}

	case StateVictory:
		if gsm.uiManager != nil {
			gsm.uiManager.HideVictoryScreen()
		}
	}
}

// Update updates the state manager
func (gsm *GameStateManager) Update(deltaTime float64) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()

	gsm.stateTimer += deltaTime

	switch gsm.currentState {
	case StatePlaying:
		gsm.checkWinConditions()
		gsm.checkLoseConditions()

	case StateLevelTransition:
		gsm.updateLevelTransition(deltaTime)

	case StateLevelComplete:
		// Wait for player input to continue
		if gsm.uiManager != nil && gsm.uiManager.GetInput() == InputContinue {
			gsm.proceedToNextLevel()
		}

	case StateGameOver:
		// Auto-return to main menu after 3 seconds
		if gsm.stateTimer > 3.0 {
			gsm.currentState = StateMainMenu
			gsm.onStateEnter(StateMainMenu)
		}
	}
}

// checkWinConditions checks if any win conditions are met
func (gsm *GameStateManager) checkWinConditions() {
	if gsm.levelManager == nil {
		return
	}

	// Check if level goal reached
	if gsm.levelManager.IsLevelComplete() {
		gsm.onWinCondition(WinLevelComplete)
		return
	}

	// Check if boss defeated
	if gsm.levelManager.IsBossLevel() && gsm.levelManager.IsBossDefeated() {
		gsm.onWinCondition(WinBossDefeated)
		return
	}
}

// checkLoseConditions checks if any lose conditions are met
func (gsm *GameStateManager) checkLoseConditions() {
	// Check if all players out of lives
	allDead := true
	for _, lives := range gsm.livesRemaining {
		if lives > 0 {
			allDead = false
			break
		}
	}

	if allDead {
		gsm.onLoseCondition(LoseNoLives)
		return
	}

	// Check if time expired (if level has time limit)
	if gsm.levelManager != nil && gsm.levelManager.HasTimeLimit() && gsm.levelManager.GetTimeRemaining() <= 0 {
		gsm.onLoseCondition(LoseTimeExpired)
		return
	}
}

// onWinCondition handles a win condition being met
func (gsm *GameStateManager) onWinCondition(condition WinCondition) {
	log.Printf("Win condition met: %v", condition)

	switch condition {
	case WinLevelComplete:
		gsm.currentState = StateLevelComplete
		gsm.onStateEnter(StateLevelComplete)

	case WinBossDefeated:
		gsm.currentState = StateLevelComplete
		gsm.onStateEnter(StateLevelComplete)

	case WinAllLevelsComplete:
		gsm.currentState = StateVictory
		gsm.onStateEnter(StateVictory)
	}
}

// onLoseCondition handles a lose condition being met
func (gsm *GameStateManager) onLoseCondition(condition LoseCondition) {
	log.Printf("Lose condition met: %v", condition)

	switch condition {
	case LoseNoLives:
		gsm.gameOverReason = "Out of Lives"

	case LoseTimeExpired:
		gsm.gameOverReason = "Time's Up!"

	case LoseBothPlayersDead:
		gsm.gameOverReason = "Game Over"
	}

	gsm.currentState = StateGameOver
	gsm.onStateEnter(StateGameOver)
}

// handleLevelComplete handles level completion
func (gsm *GameStateManager) handleLevelComplete() {
	if gsm.levelManager == nil {
		return
	}

	// Calculate scores and bonuses
	timeRemaining := gsm.levelManager.GetTimeRemaining()
	if gsm.scoreManager != nil {
		for playerID := range gsm.livesRemaining {
			gsm.scoreManager.OnLevelComplete(playerID, timeRemaining)
		}
	}

	// Save progress
	if gsm.saveManager != nil && gsm.scoreManager != nil {
		gsm.saveManager.SaveProgress(gsm.currentLevel, gsm.livesRemaining, gsm.scoreManager.GetScores())
	}

	// Show level complete screen
	if gsm.uiManager != nil && gsm.scoreManager != nil {
		gsm.uiManager.ShowLevelCompleteScreen(gsm.scoreManager.GetScores(), timeRemaining)
	}

	// Play victory music
	// audio.PlayMusic("level_complete")

	// Determine next level
	gsm.nextLevel = gsm.levelManager.GetNextLevel(gsm.currentLevel)
}

// handleGameOver handles game over
func (gsm *GameStateManager) handleGameOver() {
	// Show game over screen
	if gsm.uiManager != nil && gsm.scoreManager != nil {
		gsm.uiManager.ShowGameOverScreen(gsm.gameOverReason, gsm.scoreManager.GetScores())
	}

	// Play game over music
	// audio.PlayMusic("game_over")

	// Check for high scores
	if gsm.uiManager != nil && gsm.scoreManager != nil {
		for playerID := range gsm.livesRemaining {
			score := gsm.scoreManager.GetPlayerScore(playerID)
			if gsm.scoreManager.IsHighScore(score.TotalScore) {
				gsm.uiManager.ShowHighScoreEntry(playerID)
			}
		}
	}
}

// handleVictory handles game victory
func (gsm *GameStateManager) handleVictory() {
	// Show victory screen
	if gsm.uiManager != nil && gsm.scoreManager != nil {
		gsm.uiManager.ShowVictoryScreen(gsm.scoreManager.GetScores())
	}

	// Play victory music
	// audio.PlayMusic("victory")

	// Save completion
	if gsm.saveManager != nil && gsm.scoreManager != nil {
		gsm.saveManager.SaveGameCompletion(gsm.scoreManager.GetScores())
	}

	// Unlock extras (if any)
	gsm.unlockExtras()
}

// proceedToNextLevel proceeds to the next level
func (gsm *GameStateManager) proceedToNextLevel() {
	if gsm.nextLevel == "" {
		gsm.currentState = StateMainMenu
		gsm.onStateEnter(StateMainMenu)
		return
	}

	gsm.currentLevel = gsm.nextLevel
	gsm.currentState = StateLevelTransition
	gsm.onStateEnter(StateLevelTransition)
}

// updateLevelTransition updates the level transition
func (gsm *GameStateManager) updateLevelTransition(deltaTime float64) {
	gsm.transitionTimer -= deltaTime

	if gsm.transitionTimer <= 0 {
		// Load next level
		if gsm.levelManager != nil {
			if err := gsm.levelManager.LoadLevel(gsm.currentLevel); err != nil {
				log.Printf("Error loading level %s: %v", gsm.currentLevel, err)
				gsm.currentState = StateMainMenu
				gsm.onStateEnter(StateMainMenu)
				return
			}
		}
		gsm.currentState = StatePlaying
		gsm.onStateEnter(StatePlaying)
	}
}

// unlockExtras unlocks any extras after game completion
func (gsm *GameStateManager) unlockExtras() {
	// TODO: Implement extras unlocking
	log.Println("Unlocking extras...")
}

// OnPlayerDeath handles a player death
func (gsm *GameStateManager) OnPlayerDeath(playerID int) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()

	gsm.livesRemaining[playerID]--
	log.Printf("Player %d died. Lives remaining: %d", playerID, gsm.livesRemaining[playerID])

	if gsm.livesRemaining[playerID] <= 0 {
		// Player is out of lives
		gsm.checkLoseConditions()
	} else {
		// Respawn player at checkpoint
		if gsm.levelManager != nil {
			if err := gsm.levelManager.RespawnPlayer(playerID); err != nil {
				log.Printf("Error respawning player %d: %v", playerID, err)
			}
		}
	}
}

// AddLife adds a life to a player
func (gsm *GameStateManager) AddLife(playerID int) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()

	gsm.livesRemaining[playerID]++
	log.Printf("Player %d gained a life. Lives: %d", playerID, gsm.livesRemaining[playerID])
}

// GetLives returns the number of lives for a player
func (gsm *GameStateManager) GetLives(playerID int) int {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()

	return gsm.livesRemaining[playerID]
}

// SetLives sets the number of lives for a player
func (gsm *GameStateManager) SetLives(playerID int, lives int) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()

	gsm.livesRemaining[playerID] = lives
}

// IsPlaying returns true if the game is in playing state
func (gsm *GameStateManager) IsPlaying() bool {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()

	return gsm.currentState == StatePlaying
}

// IsPaused returns true if the game is paused
func (gsm *GameStateManager) IsPaused() bool {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()

	return gsm.currentState == StatePaused
}

// IsTransitioning returns true if a level transition is in progress
func (gsm *GameStateManager) IsTransitioning() bool {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()

	return gsm.transitioning
}

// GetCurrentState returns the current game state
func (gsm *GameStateManager) GetCurrentState() GameState {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()

	return gsm.currentState
}

// GetPreviousState returns the previous game state
func (gsm *GameStateManager) GetPreviousState() GameState {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()

	return gsm.previousState
}

// GetStateTimer returns the time spent in the current state
func (gsm *GameStateManager) GetStateTimer() float64 {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()

	return gsm.stateTimer
}

// GetCurrentLevel returns the current level name
func (gsm *GameStateManager) GetCurrentLevel() string {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()

	return gsm.currentLevel
}

// SetCurrentLevel sets the current level name
func (gsm *GameStateManager) SetCurrentLevel(level string) {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()

	gsm.currentLevel = level
}

// GetGameOverReason returns the reason for game over
func (gsm *GameStateManager) GetGameOverReason() string {
	gsm.mu.RLock()
	defer gsm.mu.RUnlock()

	return gsm.gameOverReason
}

// Reset resets the state manager to initial state
func (gsm *GameStateManager) Reset() {
	gsm.mu.Lock()
	defer gsm.mu.Unlock()

	gsm.currentState = StateMainMenu
	gsm.previousState = StateMainMenu
	gsm.stateTimer = 0
	gsm.transitionTimer = 0
	gsm.transitioning = false
	gsm.currentLevel = ""
	gsm.nextLevel = ""
	gsm.gameOverReason = ""
	gsm.livesRemaining = make(map[int]int)
	gsm.livesRemaining[1] = 5
	gsm.livesRemaining[2] = 5
}
