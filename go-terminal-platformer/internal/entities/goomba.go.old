package entities

import (
	"github.com/gdamore/tcell/v2"
)

// Goomba represents the most basic enemy type - a simple walker
type Goomba struct {
	*EnemyBase
}

// NewGoomba creates a new Goomba enemy
func NewGoomba(x, y float64) *Goomba {
	base := NewEnemyBase(EnemyTypeGoomba, x, y)
	
	// Goomba-specific configuration
	base.Health = 1
	base.MaxHealth = 1
	base.Damage = 1
	base.MoveSpeed = 60.0 // 60 pixels/second
	base.Sprite = 'G'
	base.Color = tcell.ColorMaroon
	base.SpriteStyle = tcell.StyleDefault.Foreground(tcell.ColorMaroon).Bold(true)
	
	goomba := &Goomba{
		EnemyBase: base,
	}
	
	// Start walking in a random direction
	goomba.Velocity.X = base.MoveSpeed
	goomba.FacingRight = true
	goomba.State = EnemyStatePatrol
	
	return goomba
}

// Update updates the Goomba's state
func (g *Goomba) Update(deltaTime float64) {
	if g.IsDead || !g.IsActive {
		return
	}
	
	// Simple patrol behavior
	g.updatePatrol(deltaTime)
	
	// Update base (applies physics and position updates)
	g.EnemyBase.Update(deltaTime)
}

// updatePatrol handles the patrol behavior
func (g *Goomba) updatePatrol(deltaTime float64) {
	// Walk back and forth
	if g.OnGround {
		if g.FacingRight {
			g.Velocity.X = g.MoveSpeed
		} else {
			g.Velocity.X = -g.MoveSpeed
		}
	}
	
	// Check patrol bounds (if set)
	if g.PatrolRight > 0 && g.Position.X >= g.PatrolRight {
		g.FacingRight = false
	} else if g.PatrolLeft > 0 && g.Position.X <= g.PatrolLeft {
		g.FacingRight = true
	}
}

// OnStomp is called when the player stomps on the Goomba
func (g *Goomba) OnStomp(player interface{}) {
	// Goomba dies instantly on stomp
	g.Die()
	
	// Player bounce would be handled by the game engine
	// This is just a placeholder for the interface
}

// OnHitByProjectile is called when hit by a projectile
func (g *Goomba) OnHitByProjectile(projectile interface{}) {
	// Goomba dies instantly from projectile
	g.Die()
}

// OnCollideWithPlayer is called when colliding with the player
func (g *Goomba) OnCollideWithPlayer(player interface{}) {
	if g.IsDead {
		return
	}
	
	// Goomba damages player on contact
	// Player damage logic would be handled by the game engine
}

// OnCollideWithWall is called when colliding with a wall
func (g *Goomba) OnCollideWithWall() {
	// Turn around when hitting a wall
	g.FacingRight = !g.FacingRight
}

// OnReachEdge is called when reaching a platform edge
func (g *Goomba) OnReachEdge() {
	// Turn around at edges (don't fall off)
	g.FacingRight = !g.FacingRight
}

// SetPatrolBounds sets the patrol boundaries for the Goomba
func (g *Goomba) SetPatrolBounds(left, right float64) {
	g.PatrolLeft = left
	g.PatrolRight = right
}

// GetSprite returns the current sprite character for rendering
func (g *Goomba) GetSprite() rune {
	return g.Sprite
}

// GetSpriteStyle returns the current sprite style for rendering
func (g *Goomba) GetSpriteStyle() tcell.Style {
	return g.SpriteStyle
}
