package entities

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

// Boo represents a ghost enemy with shy behavior
type Boo struct {
	*EnemyBase
	isShy        bool
	targetPlayer interface{} // Reference to player being chased
}

// NewBoo creates a new Boo enemy
func NewBoo(x, y float64) *Boo {
	base := NewEnemyBase(EnemyTypeBoo, x, y)
	
	// Boo-specific configuration
	base.Health = 999 // Effectively invulnerable (except to star)
	base.MaxHealth = 999
	base.Damage = 1
	base.MoveSpeed = 50.0 // 50 pixels/second (slow chase)
	base.Sprite = 'B'
	base.Color = tcell.ColorWhite
	base.SpriteStyle = tcell.StyleDefault.Foreground(tcell.ColorWhite).Bold(true)
	base.IsFlying = true // Boo is a flying ghost
	
	boo := &Boo{
		EnemyBase: base,
		isShy:     false,
	}
	
	boo.State = EnemyStatePatrol
	
	return boo
}

// Update updates the Boo's state
func (b *Boo) Update(deltaTime float64) {
	if b.IsDead || !b.IsActive {
		return
	}
	
	// Boo behavior is handled by UpdateWithPlayerFacing
	// This method just updates base state
	b.EnemyBase.Update(deltaTime)
}

// UpdateWithPlayerFacing updates Boo based on player's facing direction
func (b *Boo) UpdateWithPlayerFacing(deltaTime float64, playerX, playerY float64, playerFacingRight bool) {
	if b.IsDead || !b.IsActive {
		return
	}
	
	// Determine if player is looking at Boo
	isPlayerLookingAtBoo := b.isPlayerLookingAt(playerX, playerFacingRight)
	
	b.isShy = isPlayerLookingAtBoo
	
	if b.isShy {
		// Stop moving and become shy
		b.Velocity.X = 0
		b.Velocity.Y = 0
		b.State = EnemyStateIdle
	} else {
		// Chase player slowly
		b.chasePlayer(playerX, playerY)
		b.State = EnemyStateChase
	}
	
	// Update base (applies position updates)
	b.EnemyBase.Update(deltaTime)
}

// isPlayerLookingAt determines if the player is facing the Boo
func (b *Boo) isPlayerLookingAt(playerX float64, playerFacingRight bool) bool {
	// Determine which side of the player Boo is on
	booIsOnRight := b.Position.X > playerX
	
	// Player is looking at Boo if:
	// - Boo is on the right AND player is facing right
	// - Boo is on the left AND player is facing left
	return booIsOnRight == playerFacingRight
}

// chasePlayer moves Boo toward the player
func (b *Boo) chasePlayer(playerX, playerY float64) {
	// Calculate direction to player
	deltaX := playerX - b.Position.X
	deltaY := playerY - b.Position.Y
	
	// Calculate distance
	distance := math.Sqrt(deltaX*deltaX + deltaY*deltaY)
	
	if distance > 0 {
		// Normalize and apply speed
		b.Velocity.X = (deltaX / distance) * b.MoveSpeed
		b.Velocity.Y = (deltaY / distance) * b.MoveSpeed
		
		// Update facing direction
		b.FacingRight = deltaX > 0
	}
}

// OnStomp is called when the player tries to stomp on Boo
func (b *Boo) OnStomp(player interface{}) {
	// Boo cannot be stomped - player would take damage instead
	// Damage logic handled by game engine
}

// OnHitByProjectile is called when hit by a projectile
func (b *Boo) OnHitByProjectile(projectile interface{}) {
	// Boo is invulnerable to projectiles
	// No damage taken
}

// OnHitByShell is called when hit by a kicked Koopa shell
func (b *Boo) OnHitByShell(shell interface{}) {
	// Boo is invulnerable to shells
	// No damage taken
}

// TakeDamage overrides base to make Boo invulnerable to normal damage
func (b *Boo) TakeDamage(damage int) {
	// Boo is invulnerable to normal damage
	// Only star power-up can defeat Boo
}

// OnStarContact is called when player with star power-up touches Boo
func (b *Boo) OnStarContact(player interface{}) {
	// Only way to defeat Boo
	b.Die()
}

// OnCollideWithPlayer is called when colliding with the player
func (b *Boo) OnCollideWithPlayer(player interface{}) {
	if b.IsDead {
		return
	}
	
	// Boo damages player on contact (unless player has star)
	// Damage logic handled by game engine
}

// SetTargetPlayer sets the player to chase
func (b *Boo) SetTargetPlayer(player interface{}) {
	b.targetPlayer = player
}

// IsShy returns whether Boo is currently in shy state
func (b *Boo) IsShy() bool {
	return b.isShy
}

// SetShy sets the shy state (useful for testing)
func (b *Boo) SetShy(shy bool) {
	b.isShy = shy
}

// CanBeStomped returns whether Boo can be stomped
func (b *Boo) CanBeStomped() bool {
	return false // Boo cannot be stomped
}

// IsInvulnerable returns whether Boo is invulnerable to normal attacks
func (b *Boo) IsInvulnerable() bool {
	return true // Boo is invulnerable except to star
}

// GetSprite returns the current sprite character for rendering
func (b *Boo) GetSprite() rune {
	if b.isShy {
		// Could return a different sprite when shy (e.g., covering face)
		return 'b' // lowercase when shy
	}
	return b.Sprite
}

// GetSpriteStyle returns the current sprite style for rendering
func (b *Boo) GetSpriteStyle() tcell.Style {
	if b.isShy {
		// Dimmer when shy
		return tcell.StyleDefault.Foreground(tcell.ColorGray).Bold(false)
	}
	return b.SpriteStyle
}

// IsFlyingEnemy returns whether this enemy is flying
func (b *Boo) IsFlyingEnemy() bool {
	return b.IsFlying
}
