package input

import (
	"sync"

	"github.com/gdamore/tcell/v2"
)

// InputAction represents a game action
type InputAction int

const (
	ActionNone InputAction = iota
	ActionMoveLeft
	ActionMoveRight
	ActionJump
	ActionCrouch
	ActionShoot
	ActionPause
	ActionQuit
	ActionMenuUp
	ActionMenuDown
	ActionMenuSelect
	ActionMenuBack
)

// String returns the string representation of an input action
func (ia InputAction) String() string {
	switch ia {
	case ActionNone:
		return "None"
	case ActionMoveLeft:
		return "MoveLeft"
	case ActionMoveRight:
		return "MoveRight"
	case ActionJump:
		return "Jump"
	case ActionCrouch:
		return "Crouch"
	case ActionShoot:
		return "Shoot"
	case ActionPause:
		return "Pause"
	case ActionQuit:
		return "Quit"
	case ActionMenuUp:
		return "MenuUp"
	case ActionMenuDown:
		return "MenuDown"
	case ActionMenuSelect:
		return "MenuSelect"
	case ActionMenuBack:
		return "MenuBack"
	default:
		return "Unknown"
	}
}

// KeyState tracks the state of a key
type KeyState struct {
	Pressed      bool
	JustPressed  bool
	JustReleased bool
	HoldTime     float64
}

// PlayerControls defines the key bindings for a player
type PlayerControls struct {
	Up    rune
	Down  rune
	Left  rune
	Right rune
	Jump  rune
	Shoot rune
}

// InputManager handles all keyboard input
type InputManager struct {
	screen      tcell.Screen
	keyStates   map[rune]KeyState
	specialKeys map[tcell.Key]KeyState
	player1Keys PlayerControls
	player2Keys PlayerControls
	eventQueue  chan tcell.Event
	running     bool
	mu          sync.RWMutex
}

// New creates a new InputManager
func New() *InputManager {
	return &InputManager{
		keyStates:   make(map[rune]KeyState),
		specialKeys: make(map[tcell.Key]KeyState),
		eventQueue:  make(chan tcell.Event, 100),
	}
}

// Initialize sets up the input manager with default key bindings
func (im *InputManager) Initialize(screen tcell.Screen) {
	im.screen = screen

	// Default Player 1 controls (WASD + Space)
	im.player1Keys = PlayerControls{
		Up:    'w',
		Down:  's',
		Left:  'a',
		Right: 'd',
		Jump:  'w',
		Shoot: ' ',
	}

	// Default Player 2 controls (Arrow keys + RShift)
	im.player2Keys = PlayerControls{
		Up:    '↑', // We'll handle arrow keys specially
		Down:  '↓',
		Left:  '←',
		Right: '→',
		Jump:  '↑',
		Shoot: '⇧', // RShift
	}

	im.running = true
	go im.pollEvents()
}

// pollEvents continuously polls for input events in the background
func (im *InputManager) pollEvents() {
	for im.running {
		if im.screen == nil {
			break
		}
		ev := im.screen.PollEvent()
		if ev != nil {
			select {
			case im.eventQueue <- ev:
			default:
				// Queue full, skip event
			}
		}
	}
}

// Update processes queued events and updates key states
func (im *InputManager) Update(deltaTime float64) {
	im.mu.Lock()
	defer im.mu.Unlock()

	// Reset just pressed/released flags
	for key, state := range im.keyStates {
		state.JustPressed = false
		state.JustReleased = false
		if state.Pressed {
			state.HoldTime += deltaTime
		}
		im.keyStates[key] = state
	}

	for key, state := range im.specialKeys {
		state.JustPressed = false
		state.JustReleased = false
		if state.Pressed {
			state.HoldTime += deltaTime
		}
		im.specialKeys[key] = state
	}

	// Process all queued events
	for len(im.eventQueue) > 0 {
		select {
		case ev := <-im.eventQueue:
			im.processEvent(ev)
		default:
			return
		}
	}
}

// processEvent processes a single event
func (im *InputManager) processEvent(ev tcell.Event) {
	switch ev := ev.(type) {
	case *tcell.EventKey:
		im.handleKeyEvent(ev)
	case *tcell.EventResize:
		// Resize events are handled by the renderer
	}
}

// handleKeyEvent processes keyboard events
func (im *InputManager) handleKeyEvent(ev *tcell.EventKey) {
	// Handle special keys (arrow keys, escape, etc.)
	if ev.Key() != tcell.KeyRune {
		key := ev.Key()
		state := im.specialKeys[key]

		if !state.Pressed {
			state.JustPressed = true
			state.HoldTime = 0
		}
		state.Pressed = true
		state.JustReleased = false

		im.specialKeys[key] = state
		return
	}

	// Handle rune keys
	ch := ev.Rune()
	state := im.keyStates[ch]

	if !state.Pressed {
		state.JustPressed = true
		state.HoldTime = 0
	}
	state.Pressed = true
	state.JustReleased = false

	im.keyStates[ch] = state
}

// GetPlayerInput returns the current input actions for a player
func (im *InputManager) GetPlayerInput(playerNum int) []InputAction {
	im.mu.RLock()
	defer im.mu.RUnlock()

	controls := im.player1Keys
	if playerNum == 2 {
		controls = im.player2Keys
	}

	actions := []InputAction{}

	// Handle Player 1 (WASD)
	if playerNum == 1 {
		if im.IsKeyPressed(controls.Left) {
			actions = append(actions, ActionMoveLeft)
		}
		if im.IsKeyPressed(controls.Right) {
			actions = append(actions, ActionMoveRight)
		}
		if im.IsKeyJustPressed(controls.Jump) {
			actions = append(actions, ActionJump)
		}
		if im.IsKeyPressed(controls.Down) {
			actions = append(actions, ActionCrouch)
		}
		if im.IsKeyJustPressed(controls.Shoot) {
			actions = append(actions, ActionShoot)
		}
	} else {
		// Handle Player 2 (Arrow keys)
		if im.IsSpecialKeyPressed(tcell.KeyLeft) {
			actions = append(actions, ActionMoveLeft)
		}
		if im.IsSpecialKeyPressed(tcell.KeyRight) {
			actions = append(actions, ActionMoveRight)
		}
		if im.IsSpecialKeyJustPressed(tcell.KeyUp) {
			actions = append(actions, ActionJump)
		}
		if im.IsSpecialKeyPressed(tcell.KeyDown) {
			actions = append(actions, ActionCrouch)
		}
		// Right Shift is handled as a rune '⇧' or we can use Enter as alternative
		if im.IsSpecialKeyJustPressed(tcell.KeyEnter) {
			actions = append(actions, ActionShoot)
		}
	}

	return actions
}

// IsKeyPressed returns true if the key is currently pressed
func (im *InputManager) IsKeyPressed(key rune) bool {
	return im.keyStates[key].Pressed
}

// IsKeyJustPressed returns true if the key was just pressed this frame
func (im *InputManager) IsKeyJustPressed(key rune) bool {
	return im.keyStates[key].JustPressed
}

// IsKeyJustReleased returns true if the key was just released this frame
func (im *InputManager) IsKeyJustReleased(key rune) bool {
	return im.keyStates[key].JustReleased
}

// GetKeyHoldTime returns how long the key has been held
func (im *InputManager) GetKeyHoldTime(key rune) float64 {
	return im.keyStates[key].HoldTime
}

// IsSpecialKeyPressed returns true if the special key is currently pressed
func (im *InputManager) IsSpecialKeyPressed(key tcell.Key) bool {
	return im.specialKeys[key].Pressed
}

// IsSpecialKeyJustPressed returns true if the special key was just pressed this frame
func (im *InputManager) IsSpecialKeyJustPressed(key tcell.Key) bool {
	return im.specialKeys[key].JustPressed
}

// IsSpecialKeyJustReleased returns true if the special key was just released this frame
func (im *InputManager) IsSpecialKeyJustReleased(key tcell.Key) bool {
	return im.specialKeys[key].JustReleased
}

// GetSpecialKeyHoldTime returns how long the special key has been held
func (im *InputManager) GetSpecialKeyHoldTime(key tcell.Key) float64 {
	return im.specialKeys[key].HoldTime
}

// SetPlayerControls sets custom key bindings for a player
func (im *InputManager) SetPlayerControls(playerNum int, controls PlayerControls) {
	im.mu.Lock()
	defer im.mu.Unlock()

	if playerNum == 1 {
		im.player1Keys = controls
	} else if playerNum == 2 {
		im.player2Keys = controls
	}
}

// GetPlayerControls returns the current key bindings for a player
func (im *InputManager) GetPlayerControls(playerNum int) PlayerControls {
	im.mu.RLock()
	defer im.mu.RUnlock()

	if playerNum == 1 {
		return im.player1Keys
	}
	return im.player2Keys
}

// ClearKeyStates clears all key states
func (im *InputManager) ClearKeyStates() {
	im.mu.Lock()
	defer im.mu.Unlock()

	im.keyStates = make(map[rune]KeyState)
	im.specialKeys = make(map[tcell.Key]KeyState)
}

// Shutdown stops the input manager
func (im *InputManager) Shutdown() {
	im.running = false
	close(im.eventQueue)
}

// IsRunning returns whether the input manager is running
func (im *InputManager) IsRunning() bool {
	return im.running
}

// HasQueuedEvents returns true if there are events in the queue
func (im *InputManager) HasQueuedEvents() bool {
	return len(im.eventQueue) > 0
}

// GetQueuedEventCount returns the number of events in the queue
func (im *InputManager) GetQueuedEventCount() int {
	return len(im.eventQueue)
}
