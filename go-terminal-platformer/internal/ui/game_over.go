package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/lukeroy/go-terminal-platformer/internal/renderer"
)

// GameOverScreen represents the game over screen
type GameOverScreen struct {
	finalScore    int
	highScore     int
	options       []MenuOption
	selectedIndex int
	animTime      float64
	initialized   bool
	isNewHighScore bool
}

// NewGameOverScreen creates a new game over screen
func NewGameOverScreen() *GameOverScreen {
	return &GameOverScreen{
		options: make([]MenuOption, 0),
	}
}

// Initialize sets up the game over screen
func (g *GameOverScreen) Initialize() {
	if g.initialized {
		return
	}

	// Default options
	g.options = []MenuOption{
		{Text: "Try Again", Enabled: true},
		{Text: "Main Menu", Enabled: true},
		{Text: "Exit", Enabled: true},
	}

	g.selectedIndex = 0
	g.initialized = true
}

// Update updates the game over screen
func (g *GameOverScreen) Update(deltaTime float64) {
	g.animTime += deltaTime
}

// Render renders the game over screen
func (g *GameOverScreen) Render(r *renderer.Renderer) {
	width, height := r.GetSize()
	centerX := width / 2
	centerY := height / 2

	// Animated "GAME OVER" text with flashing effect
	color := tcell.ColorRed
	if int(g.animTime*2)%2 == 0 {
		color = tcell.ColorYellow
	}

	// Draw large GAME OVER text
	gameOverLines := []string{
		"╔═══════════════════════════╗",
		"║      GAME OVER           ║",
		"╚═══════════════════════════╝",
	}

	for i, line := range gameOverLines {
		r.DrawStringCentered(centerY-8+i, line, color, tcell.ColorBlack)
	}

	// Final score
	scoreText := fmt.Sprintf("Final Score: %06d", g.finalScore)
	r.DrawStringCentered(centerY-3, scoreText, tcell.ColorWhite, tcell.ColorBlack)

	// High score or new high score message
	if g.isNewHighScore {
		// Animated new high score message
		highScoreColor := tcell.ColorGreen
		if int(g.animTime*3)%2 == 0 {
			highScoreColor = tcell.ColorYellow
		}
		r.DrawStringCentered(centerY-1, "★ NEW HIGH SCORE! ★", highScoreColor, tcell.ColorBlack)
	} else {
		highScoreText := fmt.Sprintf("High Score: %06d", g.highScore)
		r.DrawStringCentered(centerY-1, highScoreText, tcell.ColorGray, tcell.ColorBlack)
	}

	// Separator
	r.DrawStringCentered(centerY+1, "─────────────────────", tcell.ColorGray, tcell.ColorBlack)

	// Options
	optionsStartY := centerY + 3
	for i, option := range g.options {
		y := optionsStartY + i
		color := tcell.ColorWhite

		// Highlight selected option
		if i == g.selectedIndex {
			color = tcell.ColorGreen
			r.DrawStringScreen(centerX-len(option.Text)/2-2, y, ">", color, tcell.ColorBlack)
		}

		// Gray out disabled options
		if !option.Enabled {
			color = tcell.ColorGray
		}

		r.DrawStringCentered(y, option.Text, color, tcell.ColorBlack)
	}

	// Controls hint
	r.DrawStringCentered(height-3, "↑↓ Navigate  Enter Select", tcell.ColorGray, tcell.ColorBlack)

	// Draw decorative skulls
	g.drawDecorations(r, centerX, centerY)
}

// drawDecorations draws decorative elements
func (g *GameOverScreen) drawDecorations(r *renderer.Renderer, centerX, centerY int) {
	// Skull decorations on sides
	skull := "☠"
	skullColor := tcell.ColorDarkRed

	// Left side
	r.DrawStringScreen(centerX-15, centerY-5, skull, skullColor, tcell.ColorBlack)
	r.DrawStringScreen(centerX-15, centerY, skull, skullColor, tcell.ColorBlack)
	r.DrawStringScreen(centerX-15, centerY+5, skull, skullColor, tcell.ColorBlack)

	// Right side
	r.DrawStringScreen(centerX+15, centerY-5, skull, skullColor, tcell.ColorBlack)
	r.DrawStringScreen(centerX+15, centerY, skull, skullColor, tcell.ColorBlack)
	r.DrawStringScreen(centerX+15, centerY+5, skull, skullColor, tcell.ColorBlack)
}

// HandleInput handles game over screen input
func (g *GameOverScreen) HandleInput(action MenuAction) {
	switch action {
	case MenuUp:
		g.moveSelection(-1)
	case MenuDown:
		g.moveSelection(1)
	case MenuSelect:
		g.selectOption()
	case MenuBack:
		// Go to main menu (handled by second option)
		if len(g.options) > 1 && g.options[1].Text == "Main Menu" {
			g.options[1].Action()
		}
	}
}

// moveSelection moves the selection up or down
func (g *GameOverScreen) moveSelection(delta int) {
	if len(g.options) == 0 {
		return
	}

	// Move selection
	g.selectedIndex += delta

	// Wrap around
	if g.selectedIndex < 0 {
		g.selectedIndex = len(g.options) - 1
	} else if g.selectedIndex >= len(g.options) {
		g.selectedIndex = 0
	}

	// Skip disabled options
	attempts := 0
	for !g.options[g.selectedIndex].Enabled && attempts < len(g.options) {
		g.selectedIndex += delta
		if g.selectedIndex < 0 {
			g.selectedIndex = len(g.options) - 1
		} else if g.selectedIndex >= len(g.options) {
			g.selectedIndex = 0
		}
		attempts++
	}
}

// selectOption executes the selected option's action
func (g *GameOverScreen) selectOption() {
	if g.selectedIndex >= 0 && g.selectedIndex < len(g.options) {
		option := g.options[g.selectedIndex]
		if option.Enabled && option.Action != nil {
			option.Action()
		}
	}
}

// OnEnter is called when entering this screen
func (g *GameOverScreen) OnEnter() {
	g.selectedIndex = 0
	g.animTime = 0
	
	// Check if new high score
	g.isNewHighScore = g.finalScore > g.highScore
}

// OnExit is called when exiting this screen
func (g *GameOverScreen) OnExit() {
	// Nothing to clean up
}

// GetType returns the screen type
func (g *GameOverScreen) GetType() ScreenType {
	return ScreenGameOver
}

// SetFinalScore sets the final score
func (g *GameOverScreen) SetFinalScore(score int) {
	g.finalScore = score
}

// SetHighScore sets the high score
func (g *GameOverScreen) SetHighScore(score int) {
	g.highScore = score
}

// SetOptions sets the menu options
func (g *GameOverScreen) SetOptions(options []MenuOption) {
	g.options = options
	g.selectedIndex = 0
}

// AddOption adds a menu option
func (g *GameOverScreen) AddOption(text string, action func(), enabled bool) {
	g.options = append(g.options, MenuOption{
		Text:    text,
		Action:  action,
		Enabled: enabled,
	})
}

// GetSelectedIndex returns the currently selected index
func (g *GameOverScreen) GetSelectedIndex() int {
	return g.selectedIndex
}

// IsNewHighScore returns true if the final score is a new high score
func (g *GameOverScreen) IsNewHighScore() bool {
	return g.isNewHighScore
}
