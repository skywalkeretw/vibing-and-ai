package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/lukeroy/go-terminal-platformer/internal/renderer"
)

// SettingOption represents a configurable setting
type SettingOption struct {
	Name     string
	Value    interface{}
	Values   []interface{}
	OnChange func(value interface{})
}

// SettingsMenu represents the settings menu screen
type SettingsMenu struct {
	options       []SettingOption
	selectedIndex int
	initialized   bool
	modified      bool
}

// NewSettingsMenu creates a new settings menu
func NewSettingsMenu() *SettingsMenu {
	return &SettingsMenu{
		options: make([]SettingOption, 0),
	}
}

// Initialize sets up the settings menu
func (s *SettingsMenu) Initialize() {
	if s.initialized {
		return
	}

	// Default settings
	s.options = []SettingOption{
		{
			Name:   "Music Volume",
			Value:  50,
			Values: []interface{}{0, 25, 50, 75, 100},
		},
		{
			Name:   "SFX Volume",
			Value:  75,
			Values: []interface{}{0, 25, 50, 75, 100},
		},
		{
			Name:   "Show FPS",
			Value:  false,
			Values: []interface{}{false, true},
		},
		{
			Name:   "Difficulty",
			Value:  "Normal",
			Values: []interface{}{"Easy", "Normal", "Hard"},
		},
		{
			Name:   "Screen Shake",
			Value:  true,
			Values: []interface{}{false, true},
		},
		{
			Name:   "Particles",
			Value:  "High",
			Values: []interface{}{"Off", "Low", "Medium", "High"},
		},
	}

	s.selectedIndex = 0
	s.initialized = true
}

// Update updates the settings menu
func (s *SettingsMenu) Update(deltaTime float64) {
	// Settings menu is mostly static
}

// Render renders the settings menu
func (s *SettingsMenu) Render(r *renderer.Renderer) {
	width, height := r.GetSize()
	centerX := width / 2
	startY := height / 4

	// Title
	r.DrawStringCentered(startY-2, "SETTINGS", tcell.ColorYellow, tcell.ColorBlack)

	// Draw settings options
	for i, option := range s.options {
		y := startY + i*2
		color := tcell.ColorWhite

		// Highlight selected option
		if i == s.selectedIndex {
			color = tcell.ColorGreen
		}

		// Option name (left aligned)
		nameX := centerX - 15
		r.DrawStringScreen(nameX, y, option.Name, color, tcell.ColorBlack)

		// Option value (right aligned with arrows)
		valueX := centerX + 5
		valueStr := s.formatValue(option.Value)
		
		if i == s.selectedIndex {
			// Show arrows for selected option
			r.DrawStringScreen(valueX-2, y, "<", color, tcell.ColorBlack)
			r.DrawStringScreen(valueX, y, valueStr, color, tcell.ColorBlack)
			r.DrawStringScreen(valueX+len(valueStr)+2, y, ">", color, tcell.ColorBlack)
		} else {
			r.DrawStringScreen(valueX, y, valueStr, color, tcell.ColorBlack)
		}
	}

	// Modified indicator
	if s.modified {
		r.DrawStringCentered(startY+len(s.options)*2+2, "* Settings Modified *", tcell.ColorYellow, tcell.ColorBlack)
	}

	// Controls hint
	controlsY := height - 3
	r.DrawStringCentered(controlsY, "↑↓ Navigate  ←→ Change  Enter Save  Esc Back", tcell.ColorGray, tcell.ColorBlack)

	// Draw border
	s.drawBorder(r, width, height)
}

// formatValue formats a setting value for display
func (s *SettingsMenu) formatValue(value interface{}) string {
	switch v := value.(type) {
	case bool:
		if v {
			return "ON"
		}
		return "OFF"
	case int:
		return fmt.Sprintf("%d", v)
	case string:
		return v
	default:
		return fmt.Sprintf("%v", v)
	}
}

// drawBorder draws a decorative border
func (s *SettingsMenu) drawBorder(r *renderer.Renderer, width, height int) {
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

// HandleInput handles settings menu input
func (s *SettingsMenu) HandleInput(action MenuAction) {
	switch action {
	case MenuUp:
		s.moveSelection(-1)
	case MenuDown:
		s.moveSelection(1)
	case MenuLeft:
		s.changeValue(-1)
	case MenuRight:
		s.changeValue(1)
	case MenuSelect:
		s.saveSettings()
	case MenuBack:
		// Return to previous screen
	}
}

// moveSelection moves the selection up or down
func (s *SettingsMenu) moveSelection(delta int) {
	if len(s.options) == 0 {
		return
	}

	s.selectedIndex += delta

	// Wrap around
	if s.selectedIndex < 0 {
		s.selectedIndex = len(s.options) - 1
	} else if s.selectedIndex >= len(s.options) {
		s.selectedIndex = 0
	}
}

// changeValue changes the selected setting's value
func (s *SettingsMenu) changeValue(delta int) {
	if s.selectedIndex < 0 || s.selectedIndex >= len(s.options) {
		return
	}

	option := &s.options[s.selectedIndex]
	if len(option.Values) == 0 {
		return
	}

	// Find current value index
	currentIndex := -1
	for i, v := range option.Values {
		if v == option.Value {
			currentIndex = i
			break
		}
	}

	if currentIndex == -1 {
		currentIndex = 0
	}

	// Move to next/previous value
	currentIndex += delta

	// Wrap around
	if currentIndex < 0 {
		currentIndex = len(option.Values) - 1
	} else if currentIndex >= len(option.Values) {
		currentIndex = 0
	}

	// Update value
	option.Value = option.Values[currentIndex]
	s.modified = true

	// Call onChange callback if set
	if option.OnChange != nil {
		option.OnChange(option.Value)
	}
}

// saveSettings saves the current settings
func (s *SettingsMenu) saveSettings() {
	// This would typically save to a config file
	// For now, just mark as not modified
	s.modified = false
}

// OnEnter is called when entering this screen
func (s *SettingsMenu) OnEnter() {
	s.selectedIndex = 0
	s.modified = false
}

// OnExit is called when exiting this screen
func (s *SettingsMenu) OnExit() {
	// Could prompt to save if modified
}

// GetType returns the screen type
func (s *SettingsMenu) GetType() ScreenType {
	return ScreenSettings
}

// SetOptions sets the settings options
func (s *SettingsMenu) SetOptions(options []SettingOption) {
	s.options = options
	s.selectedIndex = 0
}

// AddOption adds a setting option
func (s *SettingsMenu) AddOption(name string, value interface{}, values []interface{}, onChange func(interface{})) {
	s.options = append(s.options, SettingOption{
		Name:     name,
		Value:    value,
		Values:   values,
		OnChange: onChange,
	})
}

// GetOption returns a setting option by name
func (s *SettingsMenu) GetOption(name string) *SettingOption {
	for i := range s.options {
		if s.options[i].Name == name {
			return &s.options[i]
		}
	}
	return nil
}

// GetValue returns the value of a setting by name
func (s *SettingsMenu) GetValue(name string) interface{} {
	option := s.GetOption(name)
	if option != nil {
		return option.Value
	}
	return nil
}

// SetValue sets the value of a setting by name
func (s *SettingsMenu) SetValue(name string, value interface{}) {
	option := s.GetOption(name)
	if option != nil {
		option.Value = value
		s.modified = true
		if option.OnChange != nil {
			option.OnChange(value)
		}
	}
}

// IsModified returns true if settings have been modified
func (s *SettingsMenu) IsModified() bool {
	return s.modified
}
