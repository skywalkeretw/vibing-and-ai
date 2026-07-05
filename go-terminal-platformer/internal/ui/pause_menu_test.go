package ui

import (
	"testing"
)

func TestNewPauseMenu(t *testing.T) {
	menu := NewPauseMenu()
	
	if menu == nil {
		t.Fatal("NewPauseMenu returned nil")
	}
	
	if menu.options == nil {
		t.Error("options not initialized")
	}
}

func TestPauseMenu_Initialize(t *testing.T) {
	menu := NewPauseMenu()
	menu.Initialize()
	
	if !menu.initialized {
		t.Error("menu not marked as initialized")
	}
	
	if len(menu.options) == 0 {
		t.Error("no default options created")
	}
	
	// Check default options
	expectedOptions := []string{"Resume", "Restart Level", "Settings", "Quit to Menu"}
	if len(menu.options) != len(expectedOptions) {
		t.Errorf("expected %d options, got %d", len(expectedOptions), len(menu.options))
	}
	
	for i, expected := range expectedOptions {
		if i < len(menu.options) && menu.options[i].Text != expected {
			t.Errorf("option %d: expected %s, got %s", i, expected, menu.options[i].Text)
		}
	}
}

func TestPauseMenu_GetType(t *testing.T) {
	menu := NewPauseMenu()
	
	if menu.GetType() != ScreenPause {
		t.Errorf("GetType() = %v, expected ScreenPause", menu.GetType())
	}
}

func TestPauseMenu_HandleInput_Navigation(t *testing.T) {
	menu := NewPauseMenu()
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

func TestPauseMenu_HandleInput_Select(t *testing.T) {
	menu := NewPauseMenu()
	menu.Initialize()
	
	selectCalled := false
	menu.options[0].Action = func() { selectCalled = true }
	menu.options[0].Enabled = true
	menu.selectedIndex = 0
	
	menu.HandleInput(MenuSelect)
	
	if !selectCalled {
		t.Error("selected option action not called")
	}
}

func TestPauseMenu_HandleInput_Back(t *testing.T) {
	menu := NewPauseMenu()
	menu.Initialize()
	
	resumeCalled := false
	menu.options[0].Action = func() { resumeCalled = true }
	menu.options[0].Text = "Resume"
	menu.options[0].Enabled = true
	
	menu.HandleInput(MenuBack)
	
	if !resumeCalled {
		t.Error("resume action not called on back")
	}
}

func TestPauseMenu_OnEnter(t *testing.T) {
	menu := NewPauseMenu()
	menu.Initialize()
	
	menu.selectedIndex = 2
	menu.OnEnter()
	
	if menu.selectedIndex != 0 {
		t.Error("selectedIndex not reset to 0 on enter")
	}
}

func TestPauseMenu_AddOption(t *testing.T) {
	menu := NewPauseMenu()
	menu.Initialize()
	
	initialCount := len(menu.options)
	
	_ = false
	menu.AddOption("Test Option", func() {}, true)
	
	if len(menu.options) != initialCount+1 {
		t.Error("option not added")
	}
	
	lastOption := menu.options[len(menu.options)-1]
	if lastOption.Text != "Test Option" {
		t.Errorf("option text = %s, expected 'Test Option'", lastOption.Text)
	}
}

func TestPauseMenu_SetOptions(t *testing.T) {
	menu := NewPauseMenu()
	
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

func TestPauseMenu_GetSelectedIndex(t *testing.T) {
	menu := NewPauseMenu()
	menu.Initialize()
	
	menu.selectedIndex = 2
	
	if menu.GetSelectedIndex() != 2 {
		t.Errorf("GetSelectedIndex() = %d, expected 2", menu.GetSelectedIndex())
	}
}
