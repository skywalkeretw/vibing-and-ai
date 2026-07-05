package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/lukeroy/go-terminal-platformer/internal/renderer"
)

// PauseMenu represents the pause menu screen
type PauseMenu struct {
	options       []MenuOption
	selectedIndex int
	initialized   bool
}

// NewPauseMenu creates a new pause menu
func NewPauseMenu() *PauseMenu {
	return &PauseMenu{
		options: make([]MenuOption, 0),
	}
}

// Initialize sets up the pause menu
func (p *PauseMenu) Initialize() {
	if p.initialized {
		return
	}

	// Default options
	p.options = []MenuOption{
		{Text: "Resume", Enabled: true},
		{Text: "Restart Level", Enabled: true},
		{Text: "Settings", Enabled: true},
		{Text: "Quit to Menu", Enabled: true},
	}

	p.selectedIndex = 0
	p.initialized = true
}

// Update updates the pause menu
func (p *PauseMenu) Update(deltaTime float64) {
	// Pause menu is static
}

// Render renders the pause menu
func (p *PauseMenu) Render(r *renderer.Renderer) {
	width, height := r.GetSize()
	centerX := width / 2
	centerY := height / 2

	// Draw semi-transparent overlay
	p.renderOverlay(r, width, height)

	// Draw pause box
	boxWidth := 30
	boxHeight := 12
	boxX := centerX - boxWidth/2
	boxY := centerY - boxHeight/2

	p.drawBox(r, boxX, boxY, boxWidth, boxHeight)

	// Title
	r.DrawStringCentered(boxY+2, "PAUSED", tcell.ColorYellow, tcell.ColorBlack)

	// Options
	optionsStartY := boxY + 5
	for i, option := range p.options {
		y := optionsStartY + i
		color := tcell.ColorWhite

		// Highlight selected option
		if i == p.selectedIndex {
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
	r.DrawStringCentered(boxY+boxHeight-2, "↑↓ Navigate  Enter Select  Esc Resume", tcell.ColorGray, tcell.ColorBlack)
}

// renderOverlay renders a semi-transparent overlay
func (p *PauseMenu) renderOverlay(r *renderer.Renderer, width, height int) {
	// Draw a dimmed background using dots or lighter characters
	overlayChar := '░'
	overlayColor := tcell.ColorGray

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Skip every other character for semi-transparency effect
			if (x+y)%2 == 0 {
				r.DrawCharScreen(x, y, overlayChar, overlayColor, tcell.ColorBlack)
			}
		}
	}
}

// drawBox draws a box with borders
func (p *PauseMenu) drawBox(r *renderer.Renderer, x, y, width, height int) {
	boxColor := tcell.ColorWhite
	bgColor := tcell.ColorBlack

	// Fill background
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			r.DrawCharScreen(x+dx, y+dy, ' ', bgColor, tcell.ColorBlack)
		}
	}

	// Draw borders
	// Top and bottom
	for dx := 0; dx < width; dx++ {
		r.DrawCharScreen(x+dx, y, '═', boxColor, tcell.ColorBlack)
		r.DrawCharScreen(x+dx, y+height-1, '═', boxColor, tcell.ColorBlack)
	}

	// Left and right
	for dy := 0; dy < height; dy++ {
		r.DrawCharScreen(x, y+dy, '║', boxColor, tcell.ColorBlack)
		r.DrawCharScreen(x+width-1, y+dy, '║', boxColor, tcell.ColorBlack)
	}

	// Corners
	r.DrawCharScreen(x, y, '╔', boxColor, tcell.ColorBlack)
	r.DrawCharScreen(x+width-1, y, '╗', boxColor, tcell.ColorBlack)
	r.DrawCharScreen(x, y+height-1, '╚', boxColor, tcell.ColorBlack)
	r.DrawCharScreen(x+width-1, y+height-1, '╝', boxColor, tcell.ColorBlack)
}

// HandleInput handles pause menu input
func (p *PauseMenu) HandleInput(action MenuAction) {
	switch action {
	case MenuUp:
		p.moveSelection(-1)
	case MenuDown:
		p.moveSelection(1)
	case MenuSelect:
		p.selectOption()
	case MenuBack:
		// Resume game (handled by first option)
		if len(p.options) > 0 && p.options[0].Text == "Resume" {
			p.options[0].Action()
		}
	}
}

// moveSelection moves the selection up or down
func (p *PauseMenu) moveSelection(delta int) {
	if len(p.options) == 0 {
		return
	}

	// Move selection
	p.selectedIndex += delta

	// Wrap around
	if p.selectedIndex < 0 {
		p.selectedIndex = len(p.options) - 1
	} else if p.selectedIndex >= len(p.options) {
		p.selectedIndex = 0
	}

	// Skip disabled options
	attempts := 0
	for !p.options[p.selectedIndex].Enabled && attempts < len(p.options) {
		p.selectedIndex += delta
		if p.selectedIndex < 0 {
			p.selectedIndex = len(p.options) - 1
		} else if p.selectedIndex >= len(p.options) {
			p.selectedIndex = 0
		}
		attempts++
	}
}

// selectOption executes the selected option's action
func (p *PauseMenu) selectOption() {
	if p.selectedIndex >= 0 && p.selectedIndex < len(p.options) {
		option := p.options[p.selectedIndex]
		if option.Enabled && option.Action != nil {
			option.Action()
		}
	}
}

// OnEnter is called when entering this screen
func (p *PauseMenu) OnEnter() {
	// Reset to first option
	p.selectedIndex = 0
}

// OnExit is called when exiting this screen
func (p *PauseMenu) OnExit() {
	// Nothing to clean up
}

// GetType returns the screen type
func (p *PauseMenu) GetType() ScreenType {
	return ScreenPause
}

// SetOptions sets the menu options
func (p *PauseMenu) SetOptions(options []MenuOption) {
	p.options = options
	p.selectedIndex = 0
}

// AddOption adds a menu option
func (p *PauseMenu) AddOption(text string, action func(), enabled bool) {
	p.options = append(p.options, MenuOption{
		Text:    text,
		Action:  action,
		Enabled: enabled,
	})
}

// GetSelectedIndex returns the currently selected index
func (p *PauseMenu) GetSelectedIndex() int {
	return p.selectedIndex
}
