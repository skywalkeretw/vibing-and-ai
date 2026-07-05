package engine

import (
	"fmt"
	"log"
	"time"

	"github.com/gdamore/tcell/v2"
)

// GameState represents the current state of the game
type GameState int

const (
	StateMenu GameState = iota
	StatePlaying
	StatePaused
	StateGameOver
	StateLevelComplete
	StateVictory
)

// String returns the string representation of the game state
func (gs GameState) String() string {
	switch gs {
	case StateMenu:
		return "Menu"
	case StatePlaying:
		return "Playing"
	case StatePaused:
		return "Paused"
	case StateGameOver:
		return "GameOver"
	case StateLevelComplete:
		return "LevelComplete"
	case StateVictory:
		return "Victory"
	default:
		return "Unknown"
	}
}

// Game represents the main game structure
type Game struct {
	screen     tcell.Screen
	state      GameState
	running    bool
	targetFPS  int
	deltaTime  float64
	lastUpdate time.Time

	// Debug
	showFPS    bool
	frameCount int
	fpsTimer   float64
	currentFPS int

	// Game systems (to be implemented in future issues)
	// renderer    *Renderer
	// inputMgr    *InputManager
	// physics     *PhysicsEngine
	// levelMgr    *LevelManager
}

// New creates a new Game instance
func New() (*Game, error) {
	return &Game{
		targetFPS: 30,
		showFPS:   true, // Enable FPS display for debugging
		state:     StateMenu,
	}, nil
}

// Initialize sets up the game and all its systems
func (g *Game) Initialize() error {
	// Initialize tcell screen
	screen, err := tcell.NewScreen()
	if err != nil {
		return fmt.Errorf("failed to create screen: %w", err)
	}

	if err := screen.Init(); err != nil {
		return fmt.Errorf("failed to initialize screen: %w", err)
	}

	// Configure screen
	screen.EnableMouse()
	screen.EnablePaste()
	screen.Clear()

	g.screen = screen
	g.running = true
	g.lastUpdate = time.Now()

	log.Println("Game initialized successfully")
	log.Printf("Target FPS: %d", g.targetFPS)
	log.Printf("Initial state: %s", g.state)

	return nil
}

// Run starts the main game loop
func (g *Game) Run() error {
	if g.screen == nil {
		return fmt.Errorf("game not initialized")
	}

	// Fixed timestep game loop
	ticker := time.NewTicker(time.Second / time.Duration(g.targetFPS))
	defer ticker.Stop()

	log.Println("Starting game loop...")

	for g.running {
		select {
		case <-ticker.C:
			g.Update()
			g.Render()
		}
	}

	log.Println("Game loop ended")
	return nil
}

// Update updates the game state
func (g *Game) Update() {
	// Calculate delta time
	now := time.Now()
	g.deltaTime = now.Sub(g.lastUpdate).Seconds()
	g.lastUpdate = now

	// Update FPS counter
	g.frameCount++
	g.fpsTimer += g.deltaTime
	if g.fpsTimer >= 1.0 {
		g.currentFPS = g.frameCount
		g.frameCount = 0
		g.fpsTimer = 0
	}

	// Poll for events
	g.handleInput()

	// Update based on current state
	switch g.state {
	case StateMenu:
		g.updateMenu()
	case StatePlaying:
		g.updatePlaying()
	case StatePaused:
		g.updatePaused()
	case StateGameOver:
		g.updateGameOver()
	case StateLevelComplete:
		g.updateLevelComplete()
	case StateVictory:
		g.updateVictory()
	}
}

// handleInput processes input events
func (g *Game) handleInput() {
	// Poll for events (non-blocking)
	ev := g.screen.PollEvent()
	if ev == nil {
		return
	}

	switch ev := ev.(type) {
	case *tcell.EventKey:
		// Handle key events
		switch ev.Key() {
		case tcell.KeyEscape:
			// ESC key - handle based on state
			switch g.state {
			case StateMenu:
				g.running = false
			case StatePlaying:
				g.ChangeState(StatePaused)
			case StatePaused:
				g.ChangeState(StatePlaying)
			default:
				g.ChangeState(StateMenu)
			}
		case tcell.KeyCtrlC:
			// Ctrl+C - always quit
			g.running = false
		case tcell.KeyRune:
			// Handle character keys
			switch ev.Rune() {
			case 'q', 'Q':
				g.running = false
			case 'p', 'P':
				// Toggle pause
				if g.state == StatePlaying {
					g.ChangeState(StatePaused)
				} else if g.state == StatePaused {
					g.ChangeState(StatePlaying)
				}
			case 'f', 'F':
				// Toggle FPS display
				g.showFPS = !g.showFPS
			case ' ':
				// Space - start game from menu
				if g.state == StateMenu {
					g.ChangeState(StatePlaying)
				}
			}
		}
	case *tcell.EventResize:
		// Handle terminal resize
		g.screen.Sync()
	}
}

// updateMenu updates the menu state
func (g *Game) updateMenu() {
	// Menu logic will be implemented in UI system issue
}

// updatePlaying updates the playing state
func (g *Game) updatePlaying() {
	// Game logic will be implemented in future issues
	// This will include:
	// - Physics updates
	// - Entity updates
	// - Collision detection
	// - Level management
}

// updatePaused updates the paused state
func (g *Game) updatePaused() {
	// Pause menu logic will be implemented in UI system issue
}

// updateGameOver updates the game over state
func (g *Game) updateGameOver() {
	// Game over logic will be implemented in UI system issue
}

// updateLevelComplete updates the level complete state
func (g *Game) updateLevelComplete() {
	// Level complete logic will be implemented in UI system issue
}

// updateVictory updates the victory state
func (g *Game) updateVictory() {
	// Victory logic will be implemented in UI system issue
}

// Render renders the current game state
func (g *Game) Render() {
	if g.screen == nil {
		return
	}

	// Clear screen
	g.screen.Clear()

	// Render based on current state
	switch g.state {
	case StateMenu:
		g.renderMenu()
	case StatePlaying:
		g.renderPlaying()
	case StatePaused:
		g.renderPaused()
	case StateGameOver:
		g.renderGameOver()
	case StateLevelComplete:
		g.renderLevelComplete()
	case StateVictory:
		g.renderVictory()
	}

	// Render FPS counter (if enabled)
	if g.showFPS {
		g.renderFPS()
	}

	// Show the screen
	g.screen.Show()
}

// renderMenu renders the main menu
func (g *Game) renderMenu() {
	width, height := g.screen.Size()
	
	// Title
	title := "GO TERMINAL PLATFORMER"
	g.drawText(width/2-len(title)/2, height/2-3, title, tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true))
	
	// Instructions
	instructions := []string{
		"Press SPACE to Start",
		"Press P to Pause",
		"Press ESC to Quit",
		"Press F to Toggle FPS",
	}
	
	for i, text := range instructions {
		g.drawText(width/2-len(text)/2, height/2+i, text, tcell.StyleDefault.Foreground(tcell.ColorWhite))
	}
}

// renderPlaying renders the playing state
func (g *Game) renderPlaying() {
	width, height := g.screen.Size()
	
	// Placeholder for actual game rendering
	text := "PLAYING - Press ESC to Pause"
	g.drawText(width/2-len(text)/2, height/2, text, tcell.StyleDefault.Foreground(tcell.ColorGreen))
	
	// Draw a simple ground line
	for x := 0; x < width; x++ {
		g.screen.SetContent(x, height-3, '═', nil, tcell.StyleDefault.Foreground(tcell.ColorBrown))
	}
}

// renderPaused renders the paused state
func (g *Game) renderPaused() {
	// First render the game state
	g.renderPlaying()
	
	// Then overlay pause menu
	width, height := g.screen.Size()
	
	// Draw semi-transparent overlay (using darker colors)
	pauseText := "PAUSED"
	g.drawText(width/2-len(pauseText)/2, height/2-1, pauseText, tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true))
	
	resumeText := "Press ESC or P to Resume"
	g.drawText(width/2-len(resumeText)/2, height/2+1, resumeText, tcell.StyleDefault.Foreground(tcell.ColorWhite))
}

// renderGameOver renders the game over state
func (g *Game) renderGameOver() {
	width, height := g.screen.Size()
	
	text := "GAME OVER"
	g.drawText(width/2-len(text)/2, height/2, text, tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true))
	
	instruction := "Press ESC to return to Menu"
	g.drawText(width/2-len(instruction)/2, height/2+2, instruction, tcell.StyleDefault.Foreground(tcell.ColorWhite))
}

// renderLevelComplete renders the level complete state
func (g *Game) renderLevelComplete() {
	width, height := g.screen.Size()
	
	text := "LEVEL COMPLETE!"
	g.drawText(width/2-len(text)/2, height/2, text, tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true))
}

// renderVictory renders the victory state
func (g *Game) renderVictory() {
	width, height := g.screen.Size()
	
	text := "VICTORY!"
	g.drawText(width/2-len(text)/2, height/2, text, tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true))
}

// renderFPS renders the FPS counter
func (g *Game) renderFPS() {
	fpsText := fmt.Sprintf("FPS: %d", g.currentFPS)
	g.drawText(1, 1, fpsText, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	
	deltaText := fmt.Sprintf("Delta: %.3fms", g.deltaTime*1000)
	g.drawText(1, 2, deltaText, tcell.StyleDefault.Foreground(tcell.ColorYellow))
	
	stateText := fmt.Sprintf("State: %s", g.state)
	g.drawText(1, 3, stateText, tcell.StyleDefault.Foreground(tcell.ColorYellow))
}

// drawText draws text at the specified position
func (g *Game) drawText(x, y int, text string, style tcell.Style) {
	for i, ch := range text {
		g.screen.SetContent(x+i, y, ch, nil, style)
	}
}

// ChangeState changes the game state
func (g *Game) ChangeState(newState GameState) {
	if g.state == newState {
		return
	}
	
	log.Printf("State transition: %s -> %s", g.state, newState)
	
	// Exit current state
	g.onStateExit(g.state)
	
	// Change state
	oldState := g.state
	g.state = newState
	
	// Enter new state
	g.onStateEnter(newState, oldState)
}

// onStateExit handles state exit logic
func (g *Game) onStateExit(state GameState) {
	switch state {
	case StateMenu:
		// Cleanup menu
	case StatePlaying:
		// Pause game systems
	case StatePaused:
		// Resume game systems
	case StateGameOver:
		// Cleanup game over screen
	case StateLevelComplete:
		// Cleanup level complete screen
	case StateVictory:
		// Cleanup victory screen
	}
}

// onStateEnter handles state entry logic
func (g *Game) onStateEnter(state GameState, fromState GameState) {
	switch state {
	case StateMenu:
		// Initialize menu
	case StatePlaying:
		// Start/resume game
	case StatePaused:
		// Show pause menu
	case StateGameOver:
		// Show game over screen
	case StateLevelComplete:
		// Show level complete screen
	case StateVictory:
		// Show victory screen
	}
}

// Shutdown cleans up resources and shuts down the game
func (g *Game) Shutdown() {
	log.Println("Shutting down game...")
	
	// Stop the game loop
	g.running = false
	
	// Clean up screen
	if g.screen != nil {
		g.screen.Fini()
		g.screen = nil
	}
	
	// Save game state (will be implemented in save system issue)
	
	log.Println("Game shutdown complete")
}

// Cleanup is an alias for Shutdown for defer compatibility
func (g *Game) Cleanup() {
	g.Shutdown()
}

// IsRunning returns whether the game is currently running
func (g *Game) IsRunning() bool {
	return g.running
}

// GetState returns the current game state
func (g *Game) GetState() GameState {
	return g.state
}

// GetDeltaTime returns the delta time for the current frame
func (g *Game) GetDeltaTime() float64 {
	return g.deltaTime
}

// GetFPS returns the current FPS
func (g *Game) GetFPS() int {
	return g.currentFPS
}

// SetShowFPS sets whether to show the FPS counter
func (g *Game) SetShowFPS(show bool) {
	g.showFPS = show
}
