package ui

import (
	"github.com/gdamore/tcell/v2"
	"github.com/lukeroy/go-terminal-platformer/internal/renderer"
)

// MenuOption represents a menu option
type MenuOption struct {
	Text    string
	Action  func()
	Enabled bool
}

// MainMenu represents the main menu screen
type MainMenu struct {
	options       []MenuOption
	selectedIndex int
	title         string
	subtitle      string
	initialized   bool
}

// NewMainMenu creates a new main menu
func NewMainMenu() *MainMenu {
	return &MainMenu{
		title:    "GO TERMINAL PLATFORMER",
		subtitle: "A Classic Platformer Adventure",
		options:  make([]MenuOption, 0),
	}
}

// Initialize sets up the main menu
func (m *MainMenu) Initialize() {
	if m.initialized {
		return
	}

	// Default options (will be set by game)
	m.options = []MenuOption{
		{Text: "Start Game", Enabled: true},
		{Text: "Continue", Enabled: false},
		{Text: "Level Select", Enabled: false},
		{Text: "Settings", Enabled: true},
		{Text: "High Scores", Enabled: true},
		{Text: "Exit", Enabled: true},
	}

	m.selectedIndex = 0
	m.initialized = true
}

// Update updates the main menu
func (m *MainMenu) Update(deltaTime float64) {
	// Main menu is mostly static, no updates needed
}

// Render renders the main menu
func (m *MainMenu) Render(r *renderer.Renderer) {
	width, height := r.GetSize()
	centerX := width / 2
	startY := height / 3

	// Draw title
	m.drawTitle(r, centerX, startY-8)

	// Draw subtitle
	subtitleX := centerX - len(m.subtitle)/2
	r.DrawStringScreen(subtitleX, startY-5, m.subtitle, tcell.ColorGray, tcell.ColorBlack)

	// Draw options
	for i, option := range m.options {
		y := startY + i*2
		color := tcell.ColorWhite

		// Highlight selected option
		if i == m.selectedIndex {
			color = tcell.ColorGreen
			r.DrawStringScreen(centerX-len(option.Text)/2-2, y, ">", color, tcell.ColorBlack)
			r.DrawStringScreen(centerX+len(option.Text)/2+2, y, "<", color, tcell.ColorBlack)
		}

		// Gray out disabled options
		if !option.Enabled {
			color = tcell.ColorGray
		}

	r.DrawStringCentered(y, option.Text, color, tcell.ColorBlack)
	}

	// Draw controls hint
	controlsY := height - 3
r.DrawStringCentered(controlsY, "↑↓ Navigate  Enter Select  Esc Exit", tcell.ColorGray, tcell.ColorBlack)

	// Draw decorative border
	m.drawBorder(r, width, height)
}

// drawTitle draws the game title with ASCII art style
func (m *MainMenu) drawTitle(r *renderer.Renderer, centerX, y int) {
	// Draw title with color gradient effect
	titleLines := []string{
		"╔═══════════════════════════════════╗",
		"║  GO TERMINAL PLATFORMER          ║",
		"╚═══════════════════════════════════╝",
	}

	for i, line := range titleLines {
		r.DrawStringCentered(y+i, line, tcell.ColorYellow, tcell.ColorBlack)
	}
}

// drawBorder draws a decorative border around the screen
func (m *MainMenu) drawBorder(r *renderer.Renderer, width, height int) {
	borderColor := tcell.ColorDarkCyan

	// Top border
	for x := 0; x < width; x++ {
	r.DrawCharScreen(x, 0, '═', borderColor, tcell.ColorBlack)
	}

	// Bottom border
	for x := 0; x < width; x++ {
	r.DrawCharScreen(x, height-1, '═', borderColor, tcell.ColorBlack)
	}

	// Left border
	for y := 0; y < height; y++ {
	r.DrawCharScreen(0, y, '║', borderColor, tcell.ColorBlack)
	}

	// Right border
	for y := 0; y < height; y++ {
	r.DrawCharScreen(width-1, y, '║', borderColor, tcell.ColorBlack)
	}

	// Corners
	r.DrawCharScreen(0, 0, '╔', borderColor, tcell.ColorBlack)
	r.DrawCharScreen(width-1, 0, '╗', borderColor, tcell.ColorBlack)
	r.DrawCharScreen(0, height-1, '╚', borderColor, tcell.ColorBlack)
	r.DrawCharScreen(width-1, height-1, '╝', borderColor, tcell.ColorBlack)
}

// HandleInput handles menu input
func (m *MainMenu) HandleInput(action MenuAction) {
	switch action {
	case MenuUp:
		m.moveSelection(-1)
	case MenuDown:
		m.moveSelection(1)
	case MenuSelect:
		m.selectOption()
	case MenuBack:
		// Exit game (handled by game manager)
		if len(m.options) > 0 {
			// Find and execute exit option
			for _, opt := range m.options {
				if opt.Text == "Exit" && opt.Enabled {
					opt.Action()
					break
				}
			}
		}
	}
}

// moveSelection moves the selection up or down
func (m *MainMenu) moveSelection(delta int) {
	if len(m.options) == 0 {
		return
	}

	// Move selection
	m.selectedIndex += delta

	// Wrap around
	if m.selectedIndex < 0 {
		m.selectedIndex = len(m.options) - 1
	} else if m.selectedIndex >= len(m.options) {
		m.selectedIndex = 0
	}

	// Skip disabled options
	attempts := 0
	for !m.options[m.selectedIndex].Enabled && attempts < len(m.options) {
		m.selectedIndex += delta
		if m.selectedIndex < 0 {
			m.selectedIndex = len(m.options) - 1
		} else if m.selectedIndex >= len(m.options) {
			m.selectedIndex = 0
		}
		attempts++
	}
}

// selectOption executes the selected option's action
func (m *MainMenu) selectOption() {
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.options) {
		option := m.options[m.selectedIndex]
		if option.Enabled && option.Action != nil {
			option.Action()
		}
	}
}

// OnEnter is called when entering this screen
func (m *MainMenu) OnEnter() {
	// Reset selection to first enabled option
	m.selectedIndex = 0
	for i, opt := range m.options {
		if opt.Enabled {
			m.selectedIndex = i
			break
		}
	}
}

// OnExit is called when exiting this screen
func (m *MainMenu) OnExit() {
	// Nothing to clean up
}

// GetType returns the screen type
func (m *MainMenu) GetType() ScreenType {
	return ScreenMainMenu
}

// SetOptions sets the menu options
func (m *MainMenu) SetOptions(options []MenuOption) {
	m.options = options
	m.selectedIndex = 0
	
	// Find first enabled option
	for i, opt := range m.options {
		if opt.Enabled {
			m.selectedIndex = i
			break
		}
	}
}

// AddOption adds a menu option
func (m *MainMenu) AddOption(text string, action func(), enabled bool) {
	m.options = append(m.options, MenuOption{
		Text:    text,
		Action:  action,
		Enabled: enabled,
	})
}

// SetOptionEnabled enables or disables an option by text
func (m *MainMenu) SetOptionEnabled(text string, enabled bool) {
	for i := range m.options {
		if m.options[i].Text == text {
			m.options[i].Enabled = enabled
			break
		}
	}
}

// GetSelectedIndex returns the currently selected index
func (m *MainMenu) GetSelectedIndex() int {
	return m.selectedIndex
}

// GetSelectedOption returns the currently selected option
func (m *MainMenu) GetSelectedOption() *MenuOption {
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.options) {
		return &m.options[m.selectedIndex]
	}
	return nil
}
