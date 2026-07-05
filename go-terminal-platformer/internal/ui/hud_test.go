package ui

import (
	"testing"

	"github.com/lukeroy/go-terminal-platformer/internal/entities"
	"github.com/lukeroy/go-terminal-platformer/internal/levels"
)

func TestNewHUD(t *testing.T) {
	hud := NewHUD()
	
	if hud == nil {
		t.Fatal("NewHUD returned nil")
	}
	
	if hud.showFPS {
		t.Error("showFPS should be false by default")
	}
}

func TestHUD_GetType(t *testing.T) {
	hud := NewHUD()
	
	if hud.GetType() != ScreenPlaying {
		t.Errorf("GetType() = %v, expected ScreenPlaying", hud.GetType())
	}
}

func TestHUD_SetPlayer1(t *testing.T) {
	hud := NewHUD()
	player := &entities.PlayerEntity{}
	
	hud.SetPlayer1(player)
	
	if hud.player1 != player {
		t.Error("player1 not set correctly")
	}
}

func TestHUD_SetPlayer2(t *testing.T) {
	hud := NewHUD()
	player := &entities.PlayerEntity{}
	
	hud.SetPlayer2(player)
	
	if hud.player2 != player {
		t.Error("player2 not set correctly")
	}
}

func TestHUD_SetLevel(t *testing.T) {
	hud := NewHUD()
	level := &levels.Level{
		World:    1,
		LevelNum: 1,
		Name:     "Test Level",
	}
	
	hud.SetLevel(level)
	
	if hud.level != level {
		t.Error("level not set correctly")
	}
}

func TestHUD_SetShowFPS(t *testing.T) {
	hud := NewHUD()
	
	hud.SetShowFPS(true)
	if !hud.showFPS {
		t.Error("showFPS not set to true")
	}
	
	hud.SetShowFPS(false)
	if hud.showFPS {
		t.Error("showFPS not set to false")
	}
}

func TestHUD_SetFPS(t *testing.T) {
	hud := NewHUD()
	
	hud.SetFPS(60)
	if hud.fps != 60 {
		t.Errorf("fps = %d, expected 60", hud.fps)
	}
}

func TestHUD_Update(t *testing.T) {
	hud := NewHUD()
	
	initialTime := hud.gameTime
	hud.Update(0.016)
	
	if hud.gameTime <= initialTime {
		t.Error("gameTime not updated")
	}
	
	expectedTime := initialTime + 0.016
	if hud.gameTime != expectedTime {
		t.Errorf("gameTime = %f, expected %f", hud.gameTime, expectedTime)
	}
}

func TestHUD_GetGameTime(t *testing.T) {
	hud := NewHUD()
	
	hud.gameTime = 123.45
	
	if hud.GetGameTime() != 123.45 {
		t.Errorf("GetGameTime() = %f, expected 123.45", hud.GetGameTime())
	}
}

func TestHUD_ResetGameTime(t *testing.T) {
	hud := NewHUD()
	
	hud.gameTime = 123.45
	hud.ResetGameTime()
	
	if hud.gameTime != 0 {
		t.Errorf("gameTime = %f, expected 0", hud.gameTime)
	}
}

func TestHUD_OnEnter(t *testing.T) {
	hud := NewHUD()
	
	hud.gameTime = 123.45
	hud.OnEnter()
	
	if hud.gameTime != 0 {
		t.Error("gameTime not reset on enter")
	}
}

func TestHUD_GetPowerUpIcon(t *testing.T) {
	hud := NewHUD()
	
	tests := []struct {
		powerUp  entities.PowerUpType
		expected string
	}{
		{entities.PowerUpFire, "[F]"},
		{entities.PowerUpStar, "[★]"},
		{entities.PowerUpSpeedBoots, "[S]"},
		{entities.PowerUpSuperJump, "[J]"},
		{entities.PowerUpShield, "[◊]"},
		{entities.PowerUpMushroom, "[M]"},
		{entities.PowerUpNone, ""},
	}
	
	for _, tt := range tests {
		result := hud.getPowerUpIcon(tt.powerUp)
		if result != tt.expected {
			t.Errorf("getPowerUpIcon(%v) = %s, expected %s", tt.powerUp, result, tt.expected)
		}
	}
}

func TestHUD_FormatTime(t *testing.T) {
	hud := NewHUD()
	
	tests := []struct {
		seconds  float64
		expected string
	}{
		{0, "Time: 00:00"},
		{30, "Time: 00:30"},
		{60, "Time: 01:00"},
		{90, "Time: 01:30"},
		{125, "Time: 02:05"},
		{3661, "Time: 61:01"},
	}
	
	for _, tt := range tests {
		result := hud.formatTime(tt.seconds)
		if result != tt.expected {
			t.Errorf("formatTime(%f) = %s, expected %s", tt.seconds, result, tt.expected)
		}
	}
}
