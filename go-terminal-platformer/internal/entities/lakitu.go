package entities

import (
	"math"

	"github.com/gdamore/tcell/v2"
)

// Lakitu represents a flying enemy that drops Spinies
type Lakitu struct {
	*EnemyBase
	dropTimer      float64
	targetHeight   float64
	minHeight      float64
	maxHeight      float64
	targetPlayer   interface{} // Reference to player being followed
	spinySpawner   func(x, y float64) // Callback to spawn Spinies
}

// NewLakitu creates a new Lakitu enemy
func NewLakitu(x, y float64) *Lakitu {
	base := NewEnemyBase(EnemyTypeLakitu, x, y)
	
	// Lakitu-specific configuration
	base.Health = 3
	base.MaxHealth = 3
	base.Damage = 0 // Lakitu doesn't directly damage player (Spinies do)
	base.MoveSpeed = 100.0 // 100 pixels/second horizontal movement
	base.Sprite = 'L'
	base.Color = tcell.ColorYellow
	base.SpriteStyle = tcell.StyleDefault.Foreground(tcell.ColorYellow).Bold(true)
	
	lakitu := &Lakitu{
		EnemyBase:    base,
		dropTimer:    0.0,
		targetHeight: 65.0, // Default height above player (middle of 50-80 range)
		minHeight:    50.0,
		maxHeight:    80.0,
	}
	
	lakitu.IsFlying = true // Lakitu is a flying enemy
	lakitu.State = EnemyStatePatrol
	
	return lakitu
}

// Update updates the Lakitu's state
func (l *Lakitu) Update(deltaTime float64) {
	if l.IsDead || !l.IsActive {
		return
	}
	
	// Update flying behavior
	l.updateFlying(deltaTime)
	
	// Update Spiny drop timer
	l.updateSpinyDrop(deltaTime)
	
	// Update base (applies position updates, but not gravity for flying enemies)
	l.EnemyBase.Update(deltaTime)
}

// updateFlying handles the flying and player following behavior
func (l *Lakitu) updateFlying(deltaTime float64) {
	// If we have a target player, follow them
	if l.targetPlayer != nil {
		// This would be implemented by the game engine
		// For now, we just maintain horizontal movement
		// The game engine would set velocity based on player position
	}
	
	// Lakitu maintains altitude (doesn't fall due to gravity)
	// This is handled by the IsFlying flag in EnemyBase
}

// updateSpinyDrop handles the periodic Spiny spawning
func (l *Lakitu) updateSpinyDrop(deltaTime float64) {
	l.dropTimer += deltaTime
	
	// Drop a Spiny every 5 seconds
	if l.dropTimer >= 5.0 {
		l.DropSpiny()
		l.dropTimer = 0.0
	}
}

// DropSpiny spawns a Spiny enemy below Lakitu
func (l *Lakitu) DropSpiny() {
	if l.spinySpawner != nil {
		// Spawn Spiny 20 pixels below Lakitu
		spawnX := l.Position.X
		spawnY := l.Position.Y + 20.0
		l.spinySpawner(spawnX, spawnY)
	}
}

// SetSpinySpawner sets the callback function for spawning Spinies
func (l *Lakitu) SetSpinySpawner(spawner func(x, y float64)) {
	l.spinySpawner = spawner
}

// SetTargetPlayer sets the player to follow
func (l *Lakitu) SetTargetPlayer(player interface{}) {
	l.targetPlayer = player
}

// FollowPlayer updates velocity to follow a player at given position
func (l *Lakitu) FollowPlayer(playerX, playerY float64) {
	// Follow player horizontally
	deltaX := playerX - l.Position.X
	
	// Only move if player is more than 20 pixels away
	if math.Abs(deltaX) > 20.0 {
		if deltaX > 0 {
			l.Velocity.X = l.MoveSpeed
			l.FacingRight = true
		} else {
			l.Velocity.X = -l.MoveSpeed
			l.FacingRight = false
		}
	} else {
		// Hover in place
		l.Velocity.X = 0
	}
	
	// Maintain height above player
	targetY := playerY - l.targetHeight
	deltaY := targetY - l.Position.Y
	
	// Smooth vertical movement
	if math.Abs(deltaY) > 5.0 {
		if deltaY > 0 {
			l.Velocity.Y = 50.0 // Move down
		} else {
			l.Velocity.Y = -50.0 // Move up
		}
	} else {
		l.Velocity.Y = 0
	}
}

// OnStomp is called when the player tries to stomp on Lakitu
func (l *Lakitu) OnStomp(player interface{}) {
	// Lakitu cannot be stomped (flying enemy)
	// Player would bounce off or take damage
}

// OnHitByProjectile is called when hit by a projectile
func (l *Lakitu) OnHitByProjectile(projectile interface{}) {
	// Lakitu takes damage from projectiles
	l.TakeDamage(1)
	
	// Die if health reaches 0
	if l.Health <= 0 {
		l.Die()
	}
}

// OnHitByShell is called when hit by a kicked Koopa shell
func (l *Lakitu) OnHitByShell(shell interface{}) {
	// Lakitu takes damage from shells
	l.TakeDamage(1)
	
	// Die if health reaches 0
	if l.Health <= 0 {
		l.Die()
	}
}

// OnCollideWithPlayer is called when colliding with the player
func (l *Lakitu) OnCollideWithPlayer(player interface{}) {
	// Lakitu doesn't directly damage player
	// Damage comes from the Spinies it drops
}

// OnCollideWithWall is called when colliding with a wall
func (l *Lakitu) OnCollideWithWall() {
	// Flying enemies can pass through walls or turn around
	l.FacingRight = !l.FacingRight
}

// SetHeightRange sets the min and max height above player
func (l *Lakitu) SetHeightRange(min, max float64) {
	l.minHeight = min
	l.maxHeight = max
	l.targetHeight = (min + max) / 2.0
}

// GetDropTimer returns the current drop timer value
func (l *Lakitu) GetDropTimer() float64 {
	return l.dropTimer
}

// ResetDropTimer resets the Spiny drop timer
func (l *Lakitu) ResetDropTimer() {
	l.dropTimer = 0.0
}

// CanBeStomped returns whether Lakitu can be stomped
func (l *Lakitu) CanBeStomped() bool {
	return false // Lakitu is flying and cannot be stomped
}

// GetSprite returns the current sprite character for rendering
func (l *Lakitu) GetSprite() rune {
	return l.Sprite
}

// GetSpriteStyle returns the current sprite style for rendering
func (l *Lakitu) GetSpriteStyle() tcell.Style {
	return l.SpriteStyle
}

// IsFlyingEnemy returns whether this enemy is flying
func (l *Lakitu) IsFlyingEnemy() bool {
	return l.IsFlying
}
