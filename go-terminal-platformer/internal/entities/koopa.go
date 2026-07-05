package entities

import (
	"github.com/gdamore/tcell/v2"
)

// KoopaState represents the specific states for Koopa
type KoopaState int

const (
	KoopaStateWalking KoopaState = iota
	KoopaStateShell
	KoopaStateKicked
)

// Koopa represents a turtle enemy with shell mechanics
type Koopa struct {
	*EnemyBase
	
	// Koopa-specific state
	koopaState   KoopaState
	shellTimer   float64
	kickVelocity Vector2
	kickedBy     interface{} // Reference to who kicked the shell
}

// NewKoopa creates a new Koopa enemy
func NewKoopa(x, y float64) *Koopa {
	base := NewEnemyBase(EnemyTypeKoopa, x, y)
	
	// Koopa-specific configuration
	base.Health = 2
	base.MaxHealth = 2
	base.Damage = 1
	base.MoveSpeed = 80.0 // 80 pixels/second
	base.Sprite = 'K'
	base.Color = tcell.ColorGreen
	base.SpriteStyle = tcell.StyleDefault.Foreground(tcell.ColorGreen)
	
	koopa := &Koopa{
		EnemyBase:  base,
		koopaState: KoopaStateWalking,
		shellTimer: 0,
	}
	
	// Start walking in a random direction
	koopa.Velocity.X = base.MoveSpeed
	koopa.FacingRight = true
	
	return koopa
}

// Update updates the Koopa's state
func (k *Koopa) Update(deltaTime float64) {
	if k.IsDead || !k.IsActive {
		return
	}
	
	// Update shell timer
	if k.shellTimer > 0 {
		k.shellTimer -= deltaTime
	}
	
	// Update based on Koopa state
	switch k.koopaState {
	case KoopaStateWalking:
		k.updateWalking(deltaTime)
	case KoopaStateShell:
		k.updateShell(deltaTime)
	case KoopaStateKicked:
		k.updateKicked(deltaTime)
	}
	
	// Update base (applies physics and position updates)
	k.EnemyBase.Update(deltaTime)
}

// updateWalking handles the walking phase behavior
func (k *Koopa) updateWalking(deltaTime float64) {
	// Check patrol bounds (if set) before updating velocity
	if k.PatrolRight > 0 && k.Position.X >= k.PatrolRight {
		k.FacingRight = false
	} else if k.PatrolLeft > 0 && k.Position.X <= k.PatrolLeft {
		k.FacingRight = true
	}
	
	// Patrol behavior - walk back and forth
	if k.OnGround {
		if k.FacingRight {
			k.Velocity.X = k.MoveSpeed
		} else {
			k.Velocity.X = -k.MoveSpeed
		}
	}
}

// updateShell handles the shell phase behavior
func (k *Koopa) updateShell(deltaTime float64) {
	// Shell is stationary
	k.Velocity.X = 0
	
	// Update sprite to shell
	k.Sprite = 'o'
	k.SpriteStyle = tcell.StyleDefault.Foreground(tcell.ColorGreen)
}

// updateKicked handles the kicked shell behavior
func (k *Koopa) updateKicked(deltaTime float64) {
	// Kicked shell moves at high speed
	k.Velocity.X = k.kickVelocity.X
	k.Velocity.Y = 0 // Shell stays on ground
	
	// Keep shell sprite
	k.Sprite = 'o'
	k.SpriteStyle = tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
}

// OnStomp is called when the player stomps on the Koopa
func (k *Koopa) OnStomp(player interface{}) {
	if k.IsDead {
		return
	}
	
	switch k.koopaState {
	case KoopaStateWalking:
		// First stomp: become shell
		k.becomeShell()
	case KoopaStateShell, KoopaStateKicked:
		// Second stomp: die
		k.Die()
	}
}

// OnKick is called when the player kicks the shell
func (k *Koopa) OnKick(player interface{}) {
	if k.koopaState != KoopaStateShell {
		return
	}
	
	// Determine kick direction based on player position
	// This is a placeholder - actual implementation would check player position
	direction := 1.0
	if k.Position.X > 0 { // Simplified logic
		direction = -1.0
	}
	
	k.kickVelocity.X = direction * 300.0 // 300 pixels/second when kicked
	k.koopaState = KoopaStateKicked
	k.kickedBy = player
}

// OnCollideWithPlayer is called when colliding with the player
func (k *Koopa) OnCollideWithPlayer(player interface{}) {
	if k.IsDead {
		return
	}
	
	switch k.koopaState {
	case KoopaStateWalking:
		// Walking Koopa damages player
		// Player damage logic would be handled by the game engine
	case KoopaStateShell:
		// Shell can be kicked
		k.OnKick(player)
	case KoopaStateKicked:
		// Kicked shell damages player if not the one who kicked it
		if k.kickedBy != player {
			// Player damage logic would be handled by the game engine
		}
	}
}

// OnCollideWithEnemy is called when the kicked shell collides with another enemy
func (k *Koopa) OnCollideWithEnemy(enemy interface{}) {
	if k.koopaState != KoopaStateKicked {
		return
	}
	
	// Kicked shell defeats other enemies
	// This would be handled by the collision system
	// The enemy parameter would have a TakeDamage or Die method called
}

// OnCollideWithWall is called when colliding with a wall
func (k *Koopa) OnCollideWithWall() {
	switch k.koopaState {
	case KoopaStateWalking:
		// Turn around when hitting a wall
		k.FacingRight = !k.FacingRight
	case KoopaStateKicked:
		// Bounce back when kicked shell hits wall
		k.kickVelocity.X = -k.kickVelocity.X
	}
}

// OnHitByProjectile is called when hit by a projectile
func (k *Koopa) OnHitByProjectile(projectile interface{}) {
	// Projectile defeats Koopa in any state
	k.Die()
}

// becomeShell transforms the Koopa into a shell
func (k *Koopa) becomeShell() {
	k.koopaState = KoopaStateShell
	k.Health = 1
	k.Velocity.X = 0
	k.shellTimer = 0
	k.Sprite = 'o'
	k.SpriteStyle = tcell.StyleDefault.Foreground(tcell.ColorGreen)
}

// Die handles Koopa death
func (k *Koopa) Die() {
	k.EnemyBase.Die()
	k.koopaState = KoopaStateWalking // Reset state
}

// GetKoopaState returns the current Koopa-specific state
func (k *Koopa) GetKoopaState() KoopaState {
	return k.koopaState
}

// IsShell returns whether the Koopa is in shell form
func (k *Koopa) IsShell() bool {
	return k.koopaState == KoopaStateShell || k.koopaState == KoopaStateKicked
}

// IsKicked returns whether the shell is currently kicked
func (k *Koopa) IsKicked() bool {
	return k.koopaState == KoopaStateKicked
}

// SetPatrolBounds sets the patrol boundaries for the Koopa
func (k *Koopa) SetPatrolBounds(left, right float64) {
	k.PatrolLeft = left
	k.PatrolRight = right
}

// GetSprite returns the current sprite character for rendering
func (k *Koopa) GetSprite() rune {
	return k.Sprite
}

// GetSpriteStyle returns the current sprite style for rendering
func (k *Koopa) GetSpriteStyle() tcell.Style {
	return k.SpriteStyle
}
