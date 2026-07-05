package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/lukeroy/go-terminal-platformer/internal/input"
	"github.com/lukeroy/go-terminal-platformer/internal/renderer"
)

// ScreenType represents different UI screen types
type ScreenType int

const (
	ScreenMainMenu ScreenType = iota
	ScreenPlaying
	ScreenPause
	ScreenGameOver
	ScreenLevelComplete
	ScreenSettings
	ScreenHighScores
	ScreenLevelSelect
)

// String returns the string representation of a screen type
func (st ScreenType) String() string {
	switch st {
	case ScreenMainMenu:
		return "MainMenu"
	case ScreenPlaying:
		return "Playing"
	case ScreenPause:
		return "Pause"
	case ScreenGameOver:
		return "GameOver"
	case ScreenLevelComplete:
		return "LevelComplete"
	case ScreenSettings:
		return "Settings"
	case ScreenHighScores:
		return "HighScores"
	case ScreenLevelSelect:
		return "LevelSelect"
	default:
		return "Unknown"
	}
}

// MenuAction represents menu navigation actions
type MenuAction int

const (
	MenuUp MenuAction = iota
	MenuDown
	MenuLeft
	MenuRight
	MenuSelect
	MenuBack
	MenuNone
)

// UIScreen interface defines the contract for all UI screens
type UIScreen interface {
	Initialize()
	Update(deltaTime float64)
	Render(renderer *renderer.Renderer)
	HandleInput(action MenuAction)
	OnEnter()
	OnExit()
	GetType() ScreenType
}

// UIManager manages all UI screens and transitions
type UIManager struct {
	currentScreen UIScreen
	previousScreen UIScreen
	screens       map[ScreenType]UIScreen
	inputMgr      *input.InputManager
	renderer      *renderer.Renderer
	transitioning bool
	transitionTime float64
}

// NewUIManager creates a new UI manager
func NewUIManager(renderer *renderer.Renderer, inputMgr *input.InputManager) *UIManager {
	return &UIManager{
		screens:  make(map[ScreenType]UIScreen),
		renderer: renderer,
		inputMgr: inputMgr,
	}
}

// Initialize sets up the UI manager
func (um *UIManager) Initialize() {
	// Screens will be registered by the game
	// Start with main menu if available
	if mainMenu, ok := um.screens[ScreenMainMenu]; ok {
		um.ShowScreen(ScreenMainMenu)
		mainMenu.Initialize()
	}
}

// RegisterScreen registers a UI screen
func (um *UIManager) RegisterScreen(screen UIScreen) {
	um.screens[screen.GetType()] = screen
	screen.Initialize()
}

// ShowScreen transitions to a specific screen
func (um *UIManager) ShowScreen(screenType ScreenType) {
	screen, ok := um.screens[screenType]
	if !ok {
		return
	}

	// Exit current screen
	if um.currentScreen != nil {
		um.currentScreen.OnExit()
		um.previousScreen = um.currentScreen
	}

	// Enter new screen
	um.currentScreen = screen
	um.currentScreen.OnEnter()
	um.transitioning = true
	um.transitionTime = 0.3 // 300ms transition
}

// ShowPreviousScreen returns to the previous screen
func (um *UIManager) ShowPreviousScreen() {
	if um.previousScreen != nil {
		um.ShowScreen(um.previousScreen.GetType())
	}
}

// Update updates the current screen
func (um *UIManager) Update(deltaTime float64) {
	if um.currentScreen == nil {
		return
	}

	// Update transition
	if um.transitioning {
		um.transitionTime -= deltaTime
		if um.transitionTime <= 0 {
			um.transitioning = false
		}
	}

	// Process input and convert to menu actions
	action := um.getMenuAction()
	if action != MenuNone {
		um.currentScreen.HandleInput(action)
	}

	// Update current screen
	um.currentScreen.Update(deltaTime)
}

// Render renders the current screen
func (um *UIManager) Render() {
	if um.currentScreen == nil {
		return
	}

	um.currentScreen.Render(um.renderer)

	// Render transition effect if transitioning
	if um.transitioning {
		um.renderTransition()
	}
}

// getMenuAction converts input to menu actions
func (um *UIManager) getMenuAction() MenuAction {
	if um.inputMgr == nil {
		return MenuNone
	}

	// Check for menu navigation keys
	if um.inputMgr.IsSpecialKeyJustPressed(tcell.KeyUp) || um.inputMgr.IsKeyJustPressed('w') {
		return MenuUp
	}
	if um.inputMgr.IsSpecialKeyJustPressed(tcell.KeyDown) || um.inputMgr.IsKeyJustPressed('s') {
		return MenuDown
	}
	if um.inputMgr.IsSpecialKeyJustPressed(tcell.KeyLeft) || um.inputMgr.IsKeyJustPressed('a') {
		return MenuLeft
	}
	if um.inputMgr.IsSpecialKeyJustPressed(tcell.KeyRight) || um.inputMgr.IsKeyJustPressed('d') {
		return MenuRight
	}
	if um.inputMgr.IsSpecialKeyJustPressed(tcell.KeyEnter) || um.inputMgr.IsKeyJustPressed(' ') {
		return MenuSelect
	}
	if um.inputMgr.IsSpecialKeyJustPressed(tcell.KeyEscape) {
		return MenuBack
	}

	return MenuNone
}

// renderTransition renders a transition effect
func (um *UIManager) renderTransition() {
	// Simple fade effect (could be enhanced)
	// For now, just a placeholder
}

// GetCurrentScreen returns the current screen
func (um *UIManager) GetCurrentScreen() UIScreen {
	return um.currentScreen
}

// GetCurrentScreenType returns the current screen type
func (um *UIManager) GetCurrentScreenType() ScreenType {
	if um.currentScreen == nil {
		return ScreenMainMenu
	}
	return um.currentScreen.GetType()
}

// IsPlaying returns true if currently in playing screen
func (um *UIManager) IsPlaying() bool {
	return um.GetCurrentScreenType() == ScreenPlaying
}

// IsPaused returns true if currently in pause screen
func (um *UIManager) IsPaused() bool {
	return um.GetCurrentScreenType() == ScreenPause
}

// IsInMenu returns true if currently in any menu screen
func (um *UIManager) IsInMenu() bool {
	screenType := um.GetCurrentScreenType()
	return screenType == ScreenMainMenu ||
		screenType == ScreenPause ||
		screenType == ScreenSettings ||
		screenType == ScreenHighScores ||
		screenType == ScreenLevelSelect
}

// Pause shows the pause menu
func (um *UIManager) Pause() {
	if um.IsPlaying() {
		um.ShowScreen(ScreenPause)
	}
}

// Resume returns to playing screen
func (um *UIManager) Resume() {
	if um.IsPaused() {
		um.ShowScreen(ScreenPlaying)
	}
}

// TextAlignment represents text alignment options
type TextAlignment int

const (
	AlignLeft TextAlignment = iota
	AlignCenter
	AlignRight
)
