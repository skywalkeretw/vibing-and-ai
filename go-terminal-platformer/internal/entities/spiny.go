package entities

import (
	"github.com/gdamore/tcell/v2"
)

// Spiny represents a spiked enemy that cannot be stomped
type Spiny struct {
	*EnemyBase
}

// NewSpiny creates a new Spiny enemy
func NewSpiny(x, y float64) *Spiny {
	base := NewEnemyBase(EnemyTypeSpiny, x, y)
	
	// Spiny-specific configuration
	base.Health = 1
	base.MaxHealth = 1
	base.Damage = 1
	base.MoveSpeed = 70.0 // 70 pixels/second (slightly faster than Goomba)
	base.Sprite = 'S'
	base.Color = tcell.ColorRed
	base.SpriteStyle = tcell.StyleDefault.Foreground(tcell.ColorRed).Bold(true)
	
	spiny := &Spiny{
		EnemyBase: base,
	}
	
	// Start walking in a random direction
	spiny.Velocity.X = base.MoveSpeed
	spiny.FacingRight = true
	spiny.State = EnemyStatePatrol
	
	return spiny
}

// Update updates the Spiny's state
func (s *Spiny) Update(deltaTime float64) {
	if s.IsDead || !s.IsActive {
		return
	}
	
	// Simple patrol behavior (same as Goomba)
	s.updatePatrol(deltaTime)
	
	// Update base (applies physics and position updates)
	s.EnemyBase.Update(deltaTime)
}

// updatePatrol handles the patrol behavior
func (s *Spiny) updatePatrol(deltaTime float64) {
	// Walk back and forth
	if s.OnGround {
		if s.FacingRight {
			s.Velocity.X = s.MoveSpeed
		} else {
			s.Velocity.X = -s.MoveSpeed
		}
	}
	
	// Check patrol bounds (if set)
	if s.PatrolRight > 0 && s.Position.X >= s.PatrolRight {
		s.FacingRight = false
	} else if s.PatrolLeft > 0 && s.Position.X <= s.PatrolLeft {
		s.FacingRight = true
	}
}

// OnStomp is called when the player tries to stomp on the Spiny
func (s *Spiny) OnStomp(player interface{}) {
	// Spiny cannot be stomped - spikes damage the player instead
	// The player damage logic would be handled by the game engine
	// This method intentionally does not call Die() on the Spiny
}

// OnHitByProjectile is called when hit by a projectile
func (s *Spiny) OnHitByProjectile(projectile interface{}) {
	// Spiny can be defeated by projectiles
	s.Die()
}

// OnHitByShell is called when hit by a kicked Koopa shell
func (s *Spiny) OnHitByShell(shell interface{}) {
	// Spiny can be defeated by kicked shells
	s.Die()
}

// OnCollideWithPlayer is called when colliding with the player
func (s *Spiny) OnCollideWithPlayer(player interface{}) {
	if s.IsDead {
		return
	}
	
	// Spiny damages player on contact
	// Player damage logic would be handled by the game engine
}

// OnCollideWithWall is called when colliding with a wall
func (s *Spiny) OnCollideWithWall() {
	// Turn around when hitting a wall
	s.FacingRight = !s.FacingRight
}

// OnReachEdge is called when reaching a platform edge
func (s *Spiny) OnReachEdge() {
	// Turn around at edges (don't fall off)
	s.FacingRight = !s.FacingRight
}

// SetPatrolBounds sets the patrol boundaries for the Spiny
func (s *Spiny) SetPatrolBounds(left, right float64) {
	s.PatrolLeft = left
	s.PatrolRight = right
}

// CanBeStomped returns whether the Spiny can be stomped
func (s *Spiny) CanBeStomped() bool {
	return false // Spiny has spikes and cannot be stomped
}

// GetSprite returns the current sprite character for rendering
func (s *Spiny) GetSprite() rune {
	return s.Sprite
}

// GetSpriteStyle returns the current sprite style for rendering
func (s *Spiny) GetSpriteStyle() tcell.Style {
	return s.SpriteStyle
}
