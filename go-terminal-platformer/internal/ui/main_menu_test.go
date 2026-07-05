package ui

import (
	"testing"
)

func TestNewMainMenu(t *testing.T) {
	menu := NewMainMenu()
	
	if menu == nil {
		t.Fatal("NewMainMenu returned nil")
	}
	
	if menu.title != "GO TERMINAL PLATFORMER" {
		t.Errorf("title = %s, expected 'GO TERMINAL PLATFORMER'", menu.title)
	}
	
	if menu.options == nil {
		t.Error("options not initialized")
	}
}

func TestMainMenu_Initialize(t *testing.T) {
	menu := NewMainMenu()
	menu.Initialize()
	
	if !menu.initialized {
		t.Error("menu not marked as initialized")
	}
	
	if len(menu.options) == 0 {
		t.Error("no default options created")
	}
	
	// Check default options exist
	expectedOptions := []string{"Start Game", "Continue", "Level Select", "Settings", "High Scores", "Exit"}
	if len(menu.options) != len(expectedOptions) {
		t.Errorf("expected %d options, got %d", len(expectedOptions), len(menu.options))
	}
	
	for i, expected := range expectedOptions {
		if i < len(menu.options) && menu.options[i].Text != expected {
			t.Errorf("option %d: expected %s, got %s", i, expected, menu.options[i].Text)
		}
	}
}

func TestMainMenu_GetType(t *testing.T) {
	menu := NewMainMenu()
	
	if menu.GetType() != ScreenMainMenu {
		t.Errorf("GetType() = %v, expected ScreenMainMenu", menu.GetType())
	}
}

func TestMainMenu_AddOption(t *testing.T) {
	menu := NewMainMenu()
	menu.Initialize()
	
	initialCount := len(menu.options)
	
	called := false
	menu.AddOption("Test Option", func() { called = true }, true)
	
	if len(menu.options) != initialCount+1 {
		t.Error("option not added")
	}
	
	lastOption := menu.options[len(menu.options)-1]
	if lastOption.Text != "Test Option" {
		t.Errorf("option text = %s, expected 'Test Option'", lastOption.Text)
	}
	
	if !lastOption.Enabled {
		t.Error("option should be enabled")
	}
	
	// Test action
	lastOption.Action()
	if !called {
		t.Error("option action not called")
	}
}

func TestMainMenu_SetOptionEnabled(t *testing.T) {
	menu := NewMainMenu()
	menu.Initialize()
	
	// Find "Continue" option (should be disabled by default)
	menu.SetOptionEnabled("Continue", true)
	
	for _, opt := range menu.options {
		if opt.Text == "Continue" {
			if !opt.Enabled {
				t.Error("Continue option should be enabled")
			}
			return
		}
	}
	
	t.Error("Continue option not found")
}

func TestMainMenu_HandleInput_Navigation(t *testing.T) {
	menu := NewMainMenu()
	menu.Initialize()
	
	initialIndex := menu.selectedIndex
	
	// Move down
	menu.HandleInput(MenuDown)
	if menu.selectedIndex <= initialIndex {
		t.Error("selection did not move down")
	}
	
	// Move up
	menu.HandleInput(MenuUp)
	if menu.selectedIndex != initialIndex {
		t.Error("selection did not move back up")
	}
}

func TestMainMenu_HandleInput_WrapAround(t *testing.T) {
	menu := NewMainMenu()
	menu.Initialize()
	
	// Move to last option
	menu.selectedIndex = len(menu.options) - 1
	lastIndex := menu.selectedIndex
	
	// Move down (should wrap to first)
	menu.HandleInput(MenuDown)
	if menu.selectedIndex != 0 {
		t.Error("selection did not wrap to first option")
	}
	
	// Move up (should wrap to last)
	menu.HandleInput(MenuUp)
	if menu.selectedIndex != lastIndex {
		t.Error("selection did not wrap to last option")
	}
}

func TestMainMenu_HandleInput_Select(t *testing.T) {
	menu := NewMainMenu()
	menu.Initialize()
	
	called := false
	menu.options[0].Action = func() { called = true }
	menu.options[0].Enabled = true
	menu.selectedIndex = 0
	
	menu.HandleInput(MenuSelect)
	
	if !called {
		t.Error("selected option action not called")
	}
}

func TestMainMenu_HandleInput_SkipDisabled(t *testing.T) {
	menu := NewMainMenu()
	menu.options = []MenuOption{
		{Text: "Option 1", Enabled: true},
		{Text: "Option 2", Enabled: false},
		{Text: "Option 3", Enabled: true},
	}
	menu.selectedIndex = 0
	
	// Move down should skip disabled option
	menu.HandleInput(MenuDown)
	
	if menu.selectedIndex == 1 {
		t.Error("selection should skip disabled option")
	}
	
	if menu.selectedIndex != 2 {
		t.Errorf("selection = %d, expected 2", menu.selectedIndex)
	}
}

func TestMainMenu_GetSelectedOption(t *testing.T) {
	menu := NewMainMenu()
	menu.Initialize()
	
	menu.selectedIndex = 0
	option := menu.GetSelectedOption()
	
	if option == nil {
		t.Fatal("GetSelectedOption returned nil")
	}
	
	if option.Text != menu.options[0].Text {
		t.Error("returned wrong option")
	}
}

func TestMainMenu_OnEnter(t *testing.T) {
	menu := NewMainMenu()
	menu.Initialize()
	
	// Set selection to middle
	menu.selectedIndex = 2
	
	// Disable first option
	menu.options[0].Enabled = false
	
	menu.OnEnter()
	
	// Should reset to first enabled option
	if menu.selectedIndex == 0 {
		t.Error("should not select disabled first option")
	}
}

func TestMainMenu_SetOptions(t *testing.T) {
	menu := NewMainMenu()
	
	newOptions := []MenuOption{
		{Text: "Custom 1", Enabled: true},
		{Text: "Custom 2", Enabled: true},
	}
	
	menu.SetOptions(newOptions)
	
	if len(menu.options) != len(newOptions) {
		t.Errorf("expected %d options, got %d", len(newOptions), len(menu.options))
	}
	
	if menu.selectedIndex != 0 {
		t.Error("selectedIndex should reset to 0")
	}
}
