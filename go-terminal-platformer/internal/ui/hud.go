package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/lukeroy/go-terminal-platformer/internal/entities"
	"github.com/lukeroy/go-terminal-platformer/internal/levels"
	"github.com/lukeroy/go-terminal-platformer/internal/renderer"
)

// HUD represents the in-game heads-up display
type HUD struct {
	player1  *entities.PlayerEntity
	player2  *entities.PlayerEntity
	level    *levels.Level
	showFPS  bool
	fps      int
	gameTime float64
}

// NewHUD creates a new HUD
func NewHUD() *HUD {
	return &HUD{
		showFPS: false,
	}
}

// Initialize sets up the HUD
func (h *HUD) Initialize() {
	// HUD is initialized with game data
}

// Update updates the HUD
func (h *HUD) Update(deltaTime float64) {
	h.gameTime += deltaTime
}

// Render renders the HUD
func (h *HUD) Render(r *renderer.Renderer) {
	width, _ := r.GetSize()

	// Player 1 info (left side)
	if h.player1 != nil {
		h.renderPlayerInfo(r, h.player1, 2, 1, true)
	}

	// Player 2 info (right side) if active
	if h.player2 != nil && !h.player2.IsDead() {
		h.renderPlayerInfo(r, h.player2, width-20, 1, false)
	}

	// Center info

	// World and level
	if h.level != nil {
		worldLevel := fmt.Sprintf("World %d-%d", h.level.World, h.level.LevelNum)
		r.DrawStringCentered(1, worldLevel, tcell.ColorWhite, tcell.ColorBlack)
		
		// Level name
		r.DrawStringCentered(2, h.level.Name, tcell.ColorGray, tcell.ColorBlack)
	}

	// FPS (if debug mode)
	if h.showFPS {
		fpsText := fmt.Sprintf("FPS: %d", h.fps)
		r.DrawStringScreen(width-10, 1, fpsText, tcell.ColorYellow, tcell.ColorBlack)
	}

	// Game time
	timeText := h.formatTime(h.gameTime)
	r.DrawStringCentered(3, timeText, tcell.ColorTeal, tcell.ColorBlack)
}

// renderPlayerInfo renders player information
func (h *HUD) renderPlayerInfo(r *renderer.Renderer, player *entities.PlayerEntity, x, y int, isP1 bool) {
	// Player label
	label := "P1"
	color := tcell.ColorBlue
	if !isP1 {
		label = "P2"
		color = tcell.ColorPurple
	}
	r.DrawStringScreen(x, y, label, color, tcell.ColorBlack)

	// Lives
	livesText := ""
	lives := player.GetLives()
	for i := 0; i < lives && i < 10; i++ {
		livesText += "♥"
	}
	if lives > 10 {
		livesText = fmt.Sprintf("♥x%d", lives)
	}
	r.DrawStringScreen(x, y+1, livesText, tcell.ColorRed, tcell.ColorBlack)

	// Coins
	coins := player.GetCoins()
	coinsText := fmt.Sprintf("○ %03d", coins)
	r.DrawStringScreen(x, y+2, coinsText, tcell.ColorYellow, tcell.ColorBlack)

	// Score (if player has score tracking)
	// Note: Score system not yet implemented in PlayerEntity
	// This is a placeholder for future implementation
	scoreText := "Score: 000000"
	r.DrawStringScreen(x, y+3, scoreText, tcell.ColorWhite, tcell.ColorBlack)

	// Power-up indicator
	powerUp := player.GetPowerUp()
	if powerUp != entities.PowerUpNone {
		powerUpText := h.getPowerUpIcon(powerUp)
		r.DrawStringScreen(x, y+4, powerUpText, tcell.ColorGreen, tcell.ColorBlack)

		// Show ammo for fire flower
		if powerUp == entities.PowerUpFire {
			ammo := player.GetAmmo()
			ammoText := fmt.Sprintf("x%d", ammo)
			r.DrawStringScreen(x+4, y+4, ammoText, tcell.ColorWhite, tcell.ColorBlack)
		}
	}

	// Invulnerability indicator
	if player.IsInvulnerable() {
		r.DrawStringScreen(x, y+5, "[★]", tcell.ColorYellow, tcell.ColorBlack)
	}
}

// getPowerUpIcon returns the icon for a power-up
func (h *HUD) getPowerUpIcon(powerUp entities.PowerUpType) string {
	switch powerUp {
	case entities.PowerUpFire:
		return "[F]"
	case entities.PowerUpStar:
		return "[★]"
	case entities.PowerUpSpeedBoots:
		return "[S]"
	case entities.PowerUpSuperJump:
		return "[J]"
	case entities.PowerUpShield:
		return "[◊]"
	case entities.PowerUpMushroom:
		return "[M]"
	default:
		return ""
	}
}

// formatTime formats game time as MM:SS
func (h *HUD) formatTime(seconds float64) string {
	totalSeconds := int(seconds)
	minutes := totalSeconds / 60
	secs := totalSeconds % 60
	return fmt.Sprintf("Time: %02d:%02d", minutes, secs)
}

// HandleInput handles HUD input (none for HUD)
func (h *HUD) HandleInput(action MenuAction) {
	// HUD doesn't handle input
}

// OnEnter is called when entering this screen
func (h *HUD) OnEnter() {
	h.gameTime = 0
}

// OnExit is called when exiting this screen
func (h *HUD) OnExit() {
	// Nothing to clean up
}

// GetType returns the screen type
func (h *HUD) GetType() ScreenType {
	return ScreenPlaying
}

// SetPlayer1 sets player 1
func (h *HUD) SetPlayer1(player *entities.PlayerEntity) {
	h.player1 = player
}

// SetPlayer2 sets player 2
func (h *HUD) SetPlayer2(player *entities.PlayerEntity) {
	h.player2 = player
}

// SetLevel sets the current level
func (h *HUD) SetLevel(level *levels.Level) {
	h.level = level
}

// SetShowFPS enables or disables FPS display
func (h *HUD) SetShowFPS(show bool) {
	h.showFPS = show
}

// SetFPS sets the current FPS value
func (h *HUD) SetFPS(fps int) {
	h.fps = fps
}

// GetGameTime returns the current game time
func (h *HUD) GetGameTime() float64 {
	return h.gameTime
}

// ResetGameTime resets the game time to zero
func (h *HUD) ResetGameTime() {
	h.gameTime = 0
}
