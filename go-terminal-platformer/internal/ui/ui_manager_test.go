package ui

import (
	"testing"

	"github.com/lukeroy/go-terminal-platformer/internal/input"
	"github.com/lukeroy/go-terminal-platformer/internal/renderer"
)

// MockScreen is a mock implementation of UIScreen for testing
type MockScreen struct {
	screenType    ScreenType
	initialized   bool
	entered       bool
	exited        bool
	updateCalled  bool
	renderCalled  bool
	inputReceived MenuAction
}

func (m *MockScreen) Initialize() {
	m.initialized = true
}

func (m *MockScreen) Update(deltaTime float64) {
	m.updateCalled = true
}

func (m *MockScreen) Render(renderer *renderer.Renderer) {
	m.renderCalled = true
}

func (m *MockScreen) HandleInput(action MenuAction) {
	m.inputReceived = action
}

func (m *MockScreen) OnEnter() {
	m.entered = true
}

func (m *MockScreen) OnExit() {
	m.exited = true
}

func (m *MockScreen) GetType() ScreenType {
	return m.screenType
}

func TestNewUIManager(t *testing.T) {
	// Skip test as it requires a real tcell screen
	t.Skip("Requires real tcell screen, skipping")
	
	inputMgr := input.New()
	
	uiMgr := NewUIManager(nil, inputMgr)
	
	if uiMgr == nil {
		t.Fatal("NewUIManager returned nil")
	}
	
	if uiMgr.screens == nil {
		t.Error("screens map not initialized")
	}
	
	if uiMgr.inputMgr != inputMgr {
		t.Error("input manager not set correctly")
	}
}

func TestUIManager_RegisterScreen(t *testing.T) {
	uiMgr := NewUIManager(nil, nil)
	
	mockScreen := &MockScreen{screenType: ScreenMainMenu}
	uiMgr.RegisterScreen(mockScreen)
	
	if !mockScreen.initialized {
		t.Error("screen not initialized on registration")
	}
	
	if _, ok := uiMgr.screens[ScreenMainMenu]; !ok {
		t.Error("screen not registered in screens map")
	}
}

func TestUIManager_ShowScreen(t *testing.T) {
	uiMgr := NewUIManager(nil, nil)
	
	mockScreen1 := &MockScreen{screenType: ScreenMainMenu}
	mockScreen2 := &MockScreen{screenType: ScreenPlaying}
	
	uiMgr.RegisterScreen(mockScreen1)
	uiMgr.RegisterScreen(mockScreen2)
	
	// Show first screen
	uiMgr.ShowScreen(ScreenMainMenu)
	
	if !mockScreen1.entered {
		t.Error("screen OnEnter not called")
	}
	
	if uiMgr.currentScreen != mockScreen1 {
		t.Error("current screen not set correctly")
	}
	
	// Show second screen
	uiMgr.ShowScreen(ScreenPlaying)
	
	if !mockScreen1.exited {
		t.Error("previous screen OnExit not called")
	}
	
	if !mockScreen2.entered {
		t.Error("new screen OnEnter not called")
	}
	
	if uiMgr.currentScreen != mockScreen2 {
		t.Error("current screen not updated")
	}
	
	if uiMgr.previousScreen != mockScreen1 {
		t.Error("previous screen not stored")
	}
}

func TestUIManager_ShowPreviousScreen(t *testing.T) {
	uiMgr := NewUIManager(nil, nil)
	
	mockScreen1 := &MockScreen{screenType: ScreenMainMenu}
	mockScreen2 := &MockScreen{screenType: ScreenPlaying}
	
	uiMgr.RegisterScreen(mockScreen1)
	uiMgr.RegisterScreen(mockScreen2)
	
	uiMgr.ShowScreen(ScreenMainMenu)
	uiMgr.ShowScreen(ScreenPlaying)
	
	// Reset flags
	mockScreen1.entered = false
	mockScreen2.exited = false
	
	uiMgr.ShowPreviousScreen()
	
	if uiMgr.currentScreen != mockScreen1 {
		t.Error("did not return to previous screen")
	}
	
	if !mockScreen2.exited {
		t.Error("current screen OnExit not called")
	}
	
	if !mockScreen1.entered {
		t.Error("previous screen OnEnter not called")
	}
}

func TestUIManager_Update(t *testing.T) {
	uiMgr := NewUIManager(nil, nil)
	
	mockScreen := &MockScreen{screenType: ScreenMainMenu}
	uiMgr.RegisterScreen(mockScreen)
	uiMgr.ShowScreen(ScreenMainMenu)
	
	mockScreen.updateCalled = false
	
	uiMgr.Update(0.016)
	
	if !mockScreen.updateCalled {
		t.Error("screen Update not called")
	}
}

func TestUIManager_GetCurrentScreenType(t *testing.T) {
	uiMgr := NewUIManager(nil, nil)
	
	mockScreen := &MockScreen{screenType: ScreenMainMenu}
	uiMgr.RegisterScreen(mockScreen)
	uiMgr.ShowScreen(ScreenMainMenu)
	
	if uiMgr.GetCurrentScreenType() != ScreenMainMenu {
		t.Error("GetCurrentScreenType returned wrong type")
	}
}

func TestUIManager_IsPlaying(t *testing.T) {
	uiMgr := NewUIManager(nil, nil)
	
	playingScreen := &MockScreen{screenType: ScreenPlaying}
	menuScreen := &MockScreen{screenType: ScreenMainMenu}
	
	uiMgr.RegisterScreen(playingScreen)
	uiMgr.RegisterScreen(menuScreen)
	
	uiMgr.ShowScreen(ScreenPlaying)
	if !uiMgr.IsPlaying() {
		t.Error("IsPlaying should return true when in playing screen")
	}
	
	uiMgr.ShowScreen(ScreenMainMenu)
	if uiMgr.IsPlaying() {
		t.Error("IsPlaying should return false when not in playing screen")
	}
}

func TestUIManager_IsPaused(t *testing.T) {
	uiMgr := NewUIManager(nil, nil)
	
	pauseScreen := &MockScreen{screenType: ScreenPause}
	playingScreen := &MockScreen{screenType: ScreenPlaying}
	
	uiMgr.RegisterScreen(pauseScreen)
	uiMgr.RegisterScreen(playingScreen)
	
	uiMgr.ShowScreen(ScreenPause)
	if !uiMgr.IsPaused() {
		t.Error("IsPaused should return true when in pause screen")
	}
	
	uiMgr.ShowScreen(ScreenPlaying)
	if uiMgr.IsPaused() {
		t.Error("IsPaused should return false when not in pause screen")
	}
}

func TestUIManager_IsInMenu(t *testing.T) {
	uiMgr := NewUIManager(nil, nil)
	
	menuScreen := &MockScreen{screenType: ScreenMainMenu}
	playingScreen := &MockScreen{screenType: ScreenPlaying}
	pauseScreen := &MockScreen{screenType: ScreenPause}
	
	uiMgr.RegisterScreen(menuScreen)
	uiMgr.RegisterScreen(playingScreen)
	uiMgr.RegisterScreen(pauseScreen)
	
	// Test menu screens
	menuScreenTypes := []ScreenType{
		ScreenMainMenu,
		ScreenPause,
		ScreenSettings,
		ScreenHighScores,
		ScreenLevelSelect,
	}
	
	for _, screenType := range menuScreenTypes {
		screen := &MockScreen{screenType: screenType}
		uiMgr.RegisterScreen(screen)
		uiMgr.ShowScreen(screenType)
		
		if !uiMgr.IsInMenu() {
			t.Errorf("IsInMenu should return true for %v", screenType)
		}
	}
	
	// Test non-menu screen
	uiMgr.ShowScreen(ScreenPlaying)
	if uiMgr.IsInMenu() {
		t.Error("IsInMenu should return false for playing screen")
	}
}

func TestUIManager_PauseResume(t *testing.T) {
	uiMgr := NewUIManager(nil, nil)
	
	playingScreen := &MockScreen{screenType: ScreenPlaying}
	pauseScreen := &MockScreen{screenType: ScreenPause}
	
	uiMgr.RegisterScreen(playingScreen)
	uiMgr.RegisterScreen(pauseScreen)
	
	// Start playing
	uiMgr.ShowScreen(ScreenPlaying)
	
	// Pause
	uiMgr.Pause()
	if uiMgr.GetCurrentScreenType() != ScreenPause {
		t.Error("Pause did not switch to pause screen")
	}
	
	// Resume
	uiMgr.Resume()
	if uiMgr.GetCurrentScreenType() != ScreenPlaying {
		t.Error("Resume did not switch back to playing screen")
	}
}

func TestScreenType_String(t *testing.T) {
	tests := []struct {
		screenType ScreenType
		expected   string
	}{
		{ScreenMainMenu, "MainMenu"},
		{ScreenPlaying, "Playing"},
		{ScreenPause, "Pause"},
		{ScreenGameOver, "GameOver"},
		{ScreenLevelComplete, "LevelComplete"},
		{ScreenSettings, "Settings"},
		{ScreenHighScores, "HighScores"},
		{ScreenLevelSelect, "LevelSelect"},
		{ScreenType(999), "Unknown"},
	}
	
	for _, tt := range tests {
		result := tt.screenType.String()
		if result != tt.expected {
			t.Errorf("ScreenType(%d).String() = %s, expected %s", tt.screenType, result, tt.expected)
		}
	}
}
