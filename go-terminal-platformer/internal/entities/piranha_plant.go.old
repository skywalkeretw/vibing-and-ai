package entities

import (
	"github.com/gdamore/tcell/v2"
)

// PiranhaPlantState represents the specific states for Piranha Plant
type PiranhaPlantState int

const (
	PiranhaPlantStateHidden PiranhaPlantState = iota
	PiranhaPlantStateEmerging
	PiranhaPlantStateVisible
	PiranhaPlantStateHiding
)

// PiranhaPlant represents a stationary enemy that emerges from pipes
type PiranhaPlant struct {
	*EnemyBase
	
	// Piranha-specific state
	plantState    PiranhaPlantState
	pipePosition  Vector2
	emergeTimer   float64
	isEmerged     bool
	emergeHeight  float64
	maxHeight     float64
	
	// Timing constants
	visibleDuration float64
	hiddenDuration  float64
}

// NewPiranhaPlant creates a new Piranha Plant enemy
func NewPiranhaPlant(x, y float64) *PiranhaPlant {
	base := NewEnemyBase(EnemyTypePiranhaPlant, x, y)
	
	// Piranha Plant-specific configuration
	base.Health = 2
	base.MaxHealth = 2
	base.Damage = 1
	base.MoveSpeed = 0 // Stationary
	base.Sprite = 'P'
	base.Color = tcell.ColorRed
	base.SpriteStyle = tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	
	plant := &PiranhaPlant{
		EnemyBase:       base,
		plantState:      PiranhaPlantStateHidden,
		pipePosition:    Vector2{X: x, Y: y},
		emergeTimer:     0,
		isEmerged:       false,
		emergeHeight:    0,
		maxHeight:       32.0, // Full emergence height
		visibleDuration: 3.0,  // 3 seconds visible
		hiddenDuration:  2.0,  // 2 seconds hidden
	}
	
	// Start hidden in pipe
	plant.Position.Y = y // Start at pipe position
	
	return plant
}

// Update updates the Piranha Plant's state
func (p *PiranhaPlant) Update(deltaTime float64) {
	if p.IsDead || !p.IsActive {
		return
	}
	
	// Update base
	p.EnemyBase.Update(deltaTime)
	
	// Update emerge/hide cycle
	p.emergeTimer += deltaTime
	
	// State machine for emerge/hide cycle
	switch p.plantState {
	case PiranhaPlantStateHidden:
		p.updateHidden(deltaTime)
	case PiranhaPlantStateEmerging:
		p.updateEmerging(deltaTime)
	case PiranhaPlantStateVisible:
		p.updateVisible(deltaTime)
	case PiranhaPlantStateHiding:
		p.updateHiding(deltaTime)
	}
	
	// Update position based on emergence height
	p.Position.Y = p.pipePosition.Y - p.emergeHeight
}

// updateHidden handles the hidden state
func (p *PiranhaPlant) updateHidden(deltaTime float64) {
	p.emergeHeight = 0
	p.isEmerged = false
	
	// Check if it's time to emerge
	if p.emergeTimer >= p.hiddenDuration {
		p.emergeTimer = 0
		p.plantState = PiranhaPlantStateEmerging
	}
}

// updateEmerging handles the emerging animation
func (p *PiranhaPlant) updateEmerging(deltaTime float64) {
	// Smoothly emerge from pipe
	p.emergeHeight = lerp(p.emergeHeight, p.maxHeight, deltaTime*5.0)
	
	// Check if fully emerged
	if p.emergeHeight >= p.maxHeight-1.0 {
		p.emergeHeight = p.maxHeight
		p.isEmerged = true
		p.plantState = PiranhaPlantStateVisible
		p.emergeTimer = 0
	}
}

// updateVisible handles the visible state
func (p *PiranhaPlant) updateVisible(deltaTime float64) {
	p.emergeHeight = p.maxHeight
	p.isEmerged = true
	
	// Check if it's time to hide
	if p.emergeTimer >= p.visibleDuration {
		p.emergeTimer = 0
		p.plantState = PiranhaPlantStateHiding
	}
}

// updateHiding handles the hiding animation
func (p *PiranhaPlant) updateHiding(deltaTime float64) {
	// Smoothly hide into pipe
	p.emergeHeight = lerp(p.emergeHeight, 0, deltaTime*5.0)
	
	// Check if fully hidden
	if p.emergeHeight <= 1.0 {
		p.emergeHeight = 0
		p.isEmerged = false
		p.plantState = PiranhaPlantStateHidden
		p.emergeTimer = 0
	}
}

// OnStomp is called when the player tries to stomp on the Piranha Plant
func (p *PiranhaPlant) OnStomp(player interface{}) {
	// Piranha Plant cannot be stomped - it damages the player instead
	// The player damage logic would be handled by the game engine
	// This method intentionally does not call TakeDamage on the plant
}

// OnHitByProjectile is called when hit by a projectile
func (p *PiranhaPlant) OnHitByProjectile(projectile interface{}) {
	// Projectile is the only way to defeat Piranha Plant
	p.TakeDamage(1)
}

// OnCollideWithPlayer is called when colliding with the player
func (p *PiranhaPlant) OnCollideWithPlayer(player interface{}) {
	if p.IsDead || !p.isEmerged {
		return
	}
	
	// Piranha Plant damages player on contact
	// Player damage logic would be handled by the game engine
}

// IsEmerged returns whether the Piranha Plant is currently emerged
func (p *PiranhaPlant) IsEmerged() bool {
	return p.isEmerged
}

// GetEmergeHeight returns the current emergence height
func (p *PiranhaPlant) GetEmergeHeight() float64 {
	return p.emergeHeight
}

// GetPlantState returns the current Piranha Plant state
func (p *PiranhaPlant) GetPlantState() PiranhaPlantState {
	return p.plantState
}

// SetPipePosition sets the pipe position for the Piranha Plant
func (p *PiranhaPlant) SetPipePosition(x, y float64) {
	p.pipePosition.X = x
	p.pipePosition.Y = y
	p.Position.X = x
	p.Position.Y = y
}

// GetPipePosition returns the pipe position
func (p *PiranhaPlant) GetPipePosition() Vector2 {
	return p.pipePosition
}

// SetVisibleDuration sets how long the plant stays visible
func (p *PiranhaPlant) SetVisibleDuration(duration float64) {
	p.visibleDuration = duration
}

// SetHiddenDuration sets how long the plant stays hidden
func (p *PiranhaPlant) SetHiddenDuration(duration float64) {
	p.hiddenDuration = duration
}

// BlocksPipe returns whether the plant is currently blocking pipe entry
func (p *PiranhaPlant) BlocksPipe() bool {
	// Plant blocks pipe when emerged or emerging
	return p.plantState == PiranhaPlantStateEmerging ||
		p.plantState == PiranhaPlantStateVisible ||
		p.plantState == PiranhaPlantStateHiding
}

// GetSprite returns the current sprite character for rendering
func (p *PiranhaPlant) GetSprite() rune {
	return p.Sprite
}

// GetSpriteStyle returns the current sprite style for rendering
func (p *PiranhaPlant) GetSpriteStyle() tcell.Style {
	// Change color based on state for visual feedback
	if p.isEmerged {
		return tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	}
	return tcell.StyleDefault.Foreground(tcell.ColorDarkRed)
}

// lerp performs linear interpolation between two values
func lerp(start, end, t float64) float64 {
	return start + (end-start)*t
}
