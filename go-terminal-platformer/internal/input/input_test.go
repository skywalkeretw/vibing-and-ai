package input

import (
	"testing"
	"time"

	"github.com/gdamore/tcell/v2"
)

func TestNew(t *testing.T) {
	im := New()
	
	if im == nil {
		t.Fatal("InputManager is nil")
	}
	
	if im.keyStates == nil {
		t.Error("keyStates map is nil")
	}
	
	if im.specialKeys == nil {
		t.Error("specialKeys map is nil")
	}
	
	if im.eventQueue == nil {
		t.Error("eventQueue is nil")
	}
}

func TestInputActionString(t *testing.T) {
	tests := []struct {
		action   InputAction
		expected string
	}{
		{ActionNone, "None"},
		{ActionMoveLeft, "MoveLeft"},
		{ActionMoveRight, "MoveRight"},
		{ActionJump, "Jump"},
		{ActionCrouch, "Crouch"},
		{ActionShoot, "Shoot"},
		{ActionPause, "Pause"},
		{ActionQuit, "Quit"},
		{ActionMenuUp, "MenuUp"},
		{ActionMenuDown, "MenuDown"},
		{ActionMenuSelect, "MenuSelect"},
		{ActionMenuBack, "MenuBack"},
	}
	
	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.action.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestKeyState(t *testing.T) {
	state := KeyState{
		Pressed:      true,
		JustPressed:  false,
		JustReleased: false,
		HoldTime:     1.5,
	}
	
	if !state.Pressed {
		t.Error("Expected Pressed to be true")
	}
	
	if state.JustPressed {
		t.Error("Expected JustPressed to be false")
	}
	
	if state.JustReleased {
		t.Error("Expected JustReleased to be false")
	}
	
	if state.HoldTime != 1.5 {
		t.Errorf("Expected HoldTime 1.5, got %f", state.HoldTime)
	}
}

func TestPlayerControls(t *testing.T) {
	controls := PlayerControls{
		Up:    'w',
		Down:  's',
		Left:  'a',
		Right: 'd',
		Jump:  'w',
		Shoot: ' ',
	}
	
	if controls.Up != 'w' {
		t.Errorf("Expected Up='w', got '%c'", controls.Up)
	}
	
	if controls.Down != 's' {
		t.Errorf("Expected Down='s', got '%c'", controls.Down)
	}
	
	if controls.Left != 'a' {
		t.Errorf("Expected Left='a', got '%c'", controls.Left)
	}
	
	if controls.Right != 'd' {
		t.Errorf("Expected Right='d', got '%c'", controls.Right)
	}
	
	if controls.Jump != 'w' {
		t.Errorf("Expected Jump='w', got '%c'", controls.Jump)
	}
	
	if controls.Shoot != ' ' {
		t.Errorf("Expected Shoot=' ', got '%c'", controls.Shoot)
	}
}

func TestSetPlayerControls(t *testing.T) {
	im := New()
	
	customControls := PlayerControls{
		Up:    'i',
		Down:  'k',
		Left:  'j',
		Right: 'l',
		Jump:  'i',
		Shoot: 'u',
	}
	
	im.SetPlayerControls(1, customControls)
	
	retrieved := im.GetPlayerControls(1)
	
	if retrieved.Up != 'i' {
		t.Errorf("Expected Up='i', got '%c'", retrieved.Up)
	}
	
	if retrieved.Left != 'j' {
		t.Errorf("Expected Left='j', got '%c'", retrieved.Left)
	}
}

func TestGetPlayerControls(t *testing.T) {
	im := New()
	
	// Set default controls
	im.player1Keys = PlayerControls{
		Up:    'w',
		Down:  's',
		Left:  'a',
		Right: 'd',
		Jump:  'w',
		Shoot: ' ',
	}
	
	im.player2Keys = PlayerControls{
		Up:    '↑',
		Down:  '↓',
		Left:  '←',
		Right: '→',
		Jump:  '↑',
		Shoot: '⇧',
	}
	
	p1Controls := im.GetPlayerControls(1)
	if p1Controls.Left != 'a' {
		t.Errorf("Expected Player 1 Left='a', got '%c'", p1Controls.Left)
	}
	
	p2Controls := im.GetPlayerControls(2)
	if p2Controls.Left != '←' {
		t.Errorf("Expected Player 2 Left='←', got '%c'", p2Controls.Left)
	}
}

func TestIsKeyPressed(t *testing.T) {
	im := New()
	
	// Simulate key press
	im.keyStates['a'] = KeyState{
		Pressed:     true,
		JustPressed: true,
		HoldTime:    0,
	}
	
	if !im.IsKeyPressed('a') {
		t.Error("Expected 'a' to be pressed")
	}
	
	if im.IsKeyPressed('b') {
		t.Error("Expected 'b' to not be pressed")
	}
}

func TestIsKeyJustPressed(t *testing.T) {
	im := New()
	
	// Simulate key just pressed
	im.keyStates['a'] = KeyState{
		Pressed:     true,
		JustPressed: true,
		HoldTime:    0,
	}
	
	if !im.IsKeyJustPressed('a') {
		t.Error("Expected 'a' to be just pressed")
	}
	
	// After update, JustPressed should be false
	im.keyStates['a'] = KeyState{
		Pressed:     true,
		JustPressed: false,
		HoldTime:    0.1,
	}
	
	if im.IsKeyJustPressed('a') {
		t.Error("Expected 'a' to not be just pressed after update")
	}
}

func TestIsKeyJustReleased(t *testing.T) {
	im := New()
	
	// Simulate key just released
	im.keyStates['a'] = KeyState{
		Pressed:      false,
		JustReleased: true,
		HoldTime:     0,
	}
	
	if !im.IsKeyJustReleased('a') {
		t.Error("Expected 'a' to be just released")
	}
}

func TestGetKeyHoldTime(t *testing.T) {
	im := New()
	
	im.keyStates['a'] = KeyState{
		Pressed:  true,
		HoldTime: 2.5,
	}
	
	holdTime := im.GetKeyHoldTime('a')
	if holdTime != 2.5 {
		t.Errorf("Expected hold time 2.5, got %f", holdTime)
	}
}

func TestIsSpecialKeyPressed(t *testing.T) {
	im := New()
	
	// Simulate special key press
	im.specialKeys[tcell.KeyEscape] = KeyState{
		Pressed:     true,
		JustPressed: true,
		HoldTime:    0,
	}
	
	if !im.IsSpecialKeyPressed(tcell.KeyEscape) {
		t.Error("Expected Escape to be pressed")
	}
	
	if im.IsSpecialKeyPressed(tcell.KeyEnter) {
		t.Error("Expected Enter to not be pressed")
	}
}

func TestIsSpecialKeyJustPressed(t *testing.T) {
	im := New()
	
	im.specialKeys[tcell.KeyEscape] = KeyState{
		Pressed:     true,
		JustPressed: true,
		HoldTime:    0,
	}
	
	if !im.IsSpecialKeyJustPressed(tcell.KeyEscape) {
		t.Error("Expected Escape to be just pressed")
	}
}

func TestUpdate(t *testing.T) {
	im := New()
	
	// Set initial state
	im.keyStates['a'] = KeyState{
		Pressed:     true,
		JustPressed: true,
		HoldTime:    0,
	}
	
	// Update with delta time
	im.Update(0.016) // ~60 FPS
	
	// JustPressed should be reset
	if im.keyStates['a'].JustPressed {
		t.Error("Expected JustPressed to be reset after update")
	}
	
	// HoldTime should increase
	if im.keyStates['a'].HoldTime < 0.016 {
		t.Errorf("Expected HoldTime to increase, got %f", im.keyStates['a'].HoldTime)
	}
}

func TestUpdateHoldTime(t *testing.T) {
	im := New()
	
	im.keyStates['a'] = KeyState{
		Pressed:  true,
		HoldTime: 0,
	}
	
	// Simulate multiple updates
	for i := 0; i < 10; i++ {
		im.Update(0.016)
	}
	
	// HoldTime should have accumulated
	holdTime := im.GetKeyHoldTime('a')
	expectedMin := 0.16 // 10 * 0.016
	if holdTime < expectedMin {
		t.Errorf("Expected hold time >= %f, got %f", expectedMin, holdTime)
	}
}

func TestClearKeyStates(t *testing.T) {
	im := New()
	
	// Add some key states
	im.keyStates['a'] = KeyState{Pressed: true}
	im.keyStates['b'] = KeyState{Pressed: true}
	im.specialKeys[tcell.KeyEscape] = KeyState{Pressed: true}
	
	// Clear
	im.ClearKeyStates()
	
	if len(im.keyStates) != 0 {
		t.Errorf("Expected keyStates to be empty, got %d entries", len(im.keyStates))
	}
	
	if len(im.specialKeys) != 0 {
		t.Errorf("Expected specialKeys to be empty, got %d entries", len(im.specialKeys))
	}
}

func TestShutdown(t *testing.T) {
	im := New()
	im.running = true
	
	im.Shutdown()
	
	if im.IsRunning() {
		t.Error("Expected input manager to not be running after shutdown")
	}
}

func TestIsRunning(t *testing.T) {
	im := New()
	
	if im.IsRunning() {
		t.Error("Expected input manager to not be running initially")
	}
	
	im.running = true
	
	if !im.IsRunning() {
		t.Error("Expected input manager to be running")
	}
}

func TestHasQueuedEvents(t *testing.T) {
	im := New()
	
	if im.HasQueuedEvents() {
		t.Error("Expected no queued events initially")
	}
	
	// Add an event (mock)
	im.eventQueue <- &tcell.EventKey{}
	
	if !im.HasQueuedEvents() {
		t.Error("Expected queued events after adding one")
	}
}

func TestGetQueuedEventCount(t *testing.T) {
	im := New()
	
	count := im.GetQueuedEventCount()
	if count != 0 {
		t.Errorf("Expected 0 queued events, got %d", count)
	}
	
	// Add events
	im.eventQueue <- &tcell.EventKey{}
	im.eventQueue <- &tcell.EventKey{}
	
	count = im.GetQueuedEventCount()
	if count != 2 {
		t.Errorf("Expected 2 queued events, got %d", count)
	}
}

func TestGetPlayerInputPlayer1(t *testing.T) {
	im := New()
	
	// Set Player 1 controls
	im.player1Keys = PlayerControls{
		Up:    'w',
		Down:  's',
		Left:  'a',
		Right: 'd',
		Jump:  'w',
		Shoot: ' ',
	}
	
	// Simulate pressing 'a' (left)
	im.keyStates['a'] = KeyState{Pressed: true}
	
	actions := im.GetPlayerInput(1)
	
	if len(actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(actions))
	}
	
	if actions[0] != ActionMoveLeft {
		t.Errorf("Expected ActionMoveLeft, got %s", actions[0])
	}
}

func TestGetPlayerInputPlayer2(t *testing.T) {
	im := New()
	
	// Simulate pressing left arrow
	im.specialKeys[tcell.KeyLeft] = KeyState{Pressed: true}
	
	actions := im.GetPlayerInput(2)
	
	if len(actions) != 1 {
		t.Errorf("Expected 1 action, got %d", len(actions))
	}
	
	if actions[0] != ActionMoveLeft {
		t.Errorf("Expected ActionMoveLeft, got %s", actions[0])
	}
}

func TestGetPlayerInputMultipleActions(t *testing.T) {
	im := New()
	
	im.player1Keys = PlayerControls{
		Up:    'w',
		Down:  's',
		Left:  'a',
		Right: 'd',
		Jump:  'w',
		Shoot: ' ',
	}
	
	// Simulate pressing multiple keys
	im.keyStates['a'] = KeyState{Pressed: true}
	im.keyStates['w'] = KeyState{Pressed: false, JustPressed: true}
	
	actions := im.GetPlayerInput(1)
	
	// Should have both MoveLeft and Jump
	if len(actions) != 2 {
		t.Errorf("Expected 2 actions, got %d", len(actions))
	}
	
	hasLeft := false
	hasJump := false
	for _, action := range actions {
		if action == ActionMoveLeft {
			hasLeft = true
		}
		if action == ActionJump {
			hasJump = true
		}
	}
	
	if !hasLeft {
		t.Error("Expected ActionMoveLeft in actions")
	}
	
	if !hasJump {
		t.Error("Expected ActionJump in actions")
	}
}

func TestConcurrentAccess(t *testing.T) {
	im := New()
	
	// Test concurrent reads and writes
	done := make(chan bool)
	
	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			im.keyStates['a'] = KeyState{Pressed: true}
			time.Sleep(time.Microsecond)
		}
		done <- true
	}()
	
	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			_ = im.IsKeyPressed('a')
			time.Sleep(time.Microsecond)
		}
		done <- true
	}()
	
	// Wait for both goroutines
	<-done
	<-done
}
