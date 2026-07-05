package engine

import (
	"testing"
	"time"
)

func TestNewGame(t *testing.T) {
	game, err := New()
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	if game == nil {
		t.Fatal("Game is nil")
	}

	if game.targetFPS != 30 {
		t.Errorf("Expected targetFPS to be 30, got %d", game.targetFPS)
	}

	if game.state != StateMenu {
		t.Errorf("Expected initial state to be StateMenu, got %v", game.state)
	}
}

func TestGameStateString(t *testing.T) {
	tests := []struct {
		state    GameState
		expected string
	}{
		{StateMenu, "Menu"},
		{StatePlaying, "Playing"},
		{StatePaused, "Paused"},
		{StateGameOver, "GameOver"},
		{StateLevelComplete, "LevelComplete"},
		{StateVictory, "Victory"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if got := tt.state.String(); got != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestChangeState(t *testing.T) {
	game, err := New()
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Test state transition
	game.ChangeState(StatePlaying)
	if game.state != StatePlaying {
		t.Errorf("Expected state to be StatePlaying, got %v", game.state)
	}

	// Test same state (should not change)
	game.ChangeState(StatePlaying)
	if game.state != StatePlaying {
		t.Errorf("Expected state to remain StatePlaying, got %v", game.state)
	}

	// Test another transition
	game.ChangeState(StatePaused)
	if game.state != StatePaused {
		t.Errorf("Expected state to be StatePaused, got %v", game.state)
	}
}

func TestDeltaTimeCalculation(t *testing.T) {
	game, err := New()
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Initialize last update time
	game.lastUpdate = time.Now()

	// Wait a bit
	time.Sleep(50 * time.Millisecond)

	// Simulate update
	now := time.Now()
	game.deltaTime = now.Sub(game.lastUpdate).Seconds()
	game.lastUpdate = now

	// Delta time should be approximately 0.05 seconds (50ms)
	if game.deltaTime < 0.04 || game.deltaTime > 0.06 {
		t.Errorf("Expected deltaTime to be around 0.05, got %f", game.deltaTime)
	}
}

func TestGameGetters(t *testing.T) {
	game, err := New()
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Test GetState
	if game.GetState() != StateMenu {
		t.Errorf("Expected GetState to return StateMenu, got %v", game.GetState())
	}

	// Test GetDeltaTime
	game.deltaTime = 0.033
	if game.GetDeltaTime() != 0.033 {
		t.Errorf("Expected GetDeltaTime to return 0.033, got %f", game.GetDeltaTime())
	}

	// Test GetFPS
	game.currentFPS = 30
	if game.GetFPS() != 30 {
		t.Errorf("Expected GetFPS to return 30, got %d", game.GetFPS())
	}

	// Test IsRunning
	game.running = true
	if !game.IsRunning() {
		t.Error("Expected IsRunning to return true")
	}

	game.running = false
	if game.IsRunning() {
		t.Error("Expected IsRunning to return false")
	}
}

func TestSetShowFPS(t *testing.T) {
	game, err := New()
	if err != nil {
		t.Fatalf("Failed to create game: %v", err)
	}

	// Default should be true
	if !game.showFPS {
		t.Error("Expected showFPS to be true by default")
	}

	// Test setting to false
	game.SetShowFPS(false)
	if game.showFPS {
		t.Error("Expected showFPS to be false after SetShowFPS(false)")
	}

	// Test setting to true
	game.SetShowFPS(true)
	if !game.showFPS {
		t.Error("Expected showFPS to be true after SetShowFPS(true)")
	}
}
