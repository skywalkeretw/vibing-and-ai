package ui

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/lukeroy/go-terminal-platformer/internal/renderer"
)

// LevelCompleteScreen represents the level complete screen
type LevelCompleteScreen struct {
	levelScore     int
	timeBonus      int
	coinsCollected int
	totalScore     int
	nextLevel      string
	showContinue   bool
	animTime       float64
	statsRevealed  int // Number of stats revealed in animation
	initialized    bool
}

// NewLevelCompleteScreen creates a new level complete screen
func NewLevelCompleteScreen() *LevelCompleteScreen {
	return &LevelCompleteScreen{
		showContinue: true,
	}
}

// Initialize sets up the level complete screen
func (l *LevelCompleteScreen) Initialize() {
	if l.initialized {
		return
	}

	l.initialized = true
}

// Update updates the level complete screen
func (l *LevelCompleteScreen) Update(deltaTime float64) {
	l.animTime += deltaTime

	// Reveal stats progressively
	if l.animTime > 0.5 && l.statsRevealed < 1 {
		l.statsRevealed = 1
	}
	if l.animTime > 1.0 && l.statsRevealed < 2 {
		l.statsRevealed = 2
	}
	if l.animTime > 1.5 && l.statsRevealed < 3 {
		l.statsRevealed = 3
	}
	if l.animTime > 2.0 && l.statsRevealed < 4 {
		l.statsRevealed = 4
	}
}

// Render renders the level complete screen
func (l *LevelCompleteScreen) Render(r *renderer.Renderer) {
	width, height := r.GetSize()
	centerX := width / 2
	startY := height / 3

	// Title with animation
	titleColor := tcell.ColorGreen
	if int(l.animTime*3)%2 == 0 {
		titleColor = tcell.ColorYellow
	}

	titleLines := []string{
		"╔═══════════════════════════════╗",
		"║    LEVEL COMPLETE!           ║",
		"╚═══════════════════════════════╝",
	}

	for i, line := range titleLines {
		r.DrawStringCentered(startY+i, line, titleColor, tcell.ColorBlack)
	}

	// Stats section
	statsY := startY + 5

	// Level Score
	if l.statsRevealed >= 1 {
		scoreText := fmt.Sprintf("Level Score:    %06d", l.levelScore)
		r.DrawStringCentered(statsY, scoreText, tcell.ColorWhite, tcell.ColorBlack)
	}

	// Time Bonus
	if l.statsRevealed >= 2 {
		timeBonusText := fmt.Sprintf("Time Bonus:     %06d", l.timeBonus)
		r.DrawStringCentered(statsY+1, timeBonusText, tcell.ColorYellow, tcell.ColorBlack)
	}

	// Coins
	if l.statsRevealed >= 3 {
		coinsText := fmt.Sprintf("Coins:          %06d", l.coinsCollected*100)
		r.DrawStringCentered(statsY+2, coinsText, tcell.ColorYellow, tcell.ColorBlack)
	}

	// Separator
	if l.statsRevealed >= 4 {
		r.DrawStringCentered(statsY+3, "─────────────────────────", tcell.ColorGray, tcell.ColorBlack)

		// Total Score
		totalText := fmt.Sprintf("Total Score:    %06d", l.totalScore)
		r.DrawStringCentered(statsY+4, totalText, tcell.ColorGreen, tcell.ColorBlack)
	}

	// Next level info
	if l.showContinue && l.statsRevealed >= 4 {
		nextLevelY := statsY + 7
		r.DrawStringCentered(nextLevelY, fmt.Sprintf("Next: %s", l.nextLevel), tcell.ColorWhite, tcell.ColorBlack)

		// Blinking continue prompt
		if int(l.animTime*2)%2 == 0 {
			r.DrawStringCentered(nextLevelY+2, "Press Enter to Continue", tcell.ColorTeal, tcell.ColorBlack)
		}
	} else if !l.showContinue && l.statsRevealed >= 4 {
		// No more levels
		r.DrawStringCentered(statsY+7, "Congratulations!", tcell.ColorGreen, tcell.ColorBlack)
		r.DrawStringCentered(statsY+8, "You completed all levels!", tcell.ColorWhite, tcell.ColorBlack)
		
		if int(l.animTime*2)%2 == 0 {
			r.DrawStringCentered(statsY+10, "Press Enter to Return to Menu", tcell.ColorTeal, tcell.ColorBlack)
		}
	}

	// Draw stars decoration
	l.drawStars(r, centerX, startY)
}

// drawStars draws decorative stars
func (l *LevelCompleteScreen) drawStars(r *renderer.Renderer, centerX, startY int) {
	star := "★"
	starColor := tcell.ColorYellow

	// Animated stars around the title
	if int(l.animTime*4)%4 == 0 {
		r.DrawStringScreen(centerX-18, startY+1, star, starColor, tcell.ColorBlack)
		r.DrawStringScreen(centerX+18, startY+1, star, starColor, tcell.ColorBlack)
	} else if int(l.animTime*4)%4 == 1 {
		r.DrawStringScreen(centerX-16, startY, star, starColor, tcell.ColorBlack)
		r.DrawStringScreen(centerX+16, startY, star, starColor, tcell.ColorBlack)
	} else if int(l.animTime*4)%4 == 2 {
		r.DrawStringScreen(centerX-18, startY+2, star, starColor, tcell.ColorBlack)
		r.DrawStringScreen(centerX+18, startY+2, star, starColor, tcell.ColorBlack)
	} else {
		r.DrawStringScreen(centerX-16, startY+1, star, starColor, tcell.ColorBlack)
		r.DrawStringScreen(centerX+16, startY+1, star, starColor, tcell.ColorBlack)
	}
}

// HandleInput handles level complete screen input
func (l *LevelCompleteScreen) HandleInput(action MenuAction) {
	// Only accept input after all stats are revealed
	if l.statsRevealed < 4 {
		return
	}

	switch action {
	case MenuSelect, MenuBack:
		// Continue to next level or return to menu
		// This will be handled by the game manager
	}
}

// OnEnter is called when entering this screen
func (l *LevelCompleteScreen) OnEnter() {
	l.animTime = 0
	l.statsRevealed = 0
}

// OnExit is called when exiting this screen
func (l *LevelCompleteScreen) OnExit() {
	// Nothing to clean up
}

// GetType returns the screen type
func (l *LevelCompleteScreen) GetType() ScreenType {
	return ScreenLevelComplete
}

// SetLevelScore sets the level score
func (l *LevelCompleteScreen) SetLevelScore(score int) {
	l.levelScore = score
}

// SetTimeBonus sets the time bonus
func (l *LevelCompleteScreen) SetTimeBonus(bonus int) {
	l.timeBonus = bonus
}

// SetCoinsCollected sets the number of coins collected
func (l *LevelCompleteScreen) SetCoinsCollected(coins int) {
	l.coinsCollected = coins
}

// SetTotalScore sets the total score
func (l *LevelCompleteScreen) SetTotalScore(score int) {
	l.totalScore = score
}

// SetNextLevel sets the next level name
func (l *LevelCompleteScreen) SetNextLevel(levelName string) {
	l.nextLevel = levelName
}

// SetShowContinue sets whether to show continue prompt
func (l *LevelCompleteScreen) SetShowContinue(show bool) {
	l.showContinue = show
}

// CalculateTotalScore calculates the total score from components
func (l *LevelCompleteScreen) CalculateTotalScore() {
	l.totalScore = l.levelScore + l.timeBonus + (l.coinsCollected * 100)
}

// IsAnimationComplete returns true if the stats animation is complete
func (l *LevelCompleteScreen) IsAnimationComplete() bool {
	return l.statsRevealed >= 4
}
