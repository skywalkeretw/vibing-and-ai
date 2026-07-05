package entities

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// HammerBro represents an enemy that throws hammers at the player
type HammerBro struct {
	*EnemyBase
	throwTimer      float64
	throwCooldown   float64
	targetPlayer    interface{} // Reference to player being targeted
	hammerSpawner   func(x, y, vx, vy float64) // Callback to spawn hammer projectiles
	movementTimer   float64
	isStanding      bool
}

// NewHammerBro creates a new Hammer Bro enemy
func NewHammerBro(x, y float64) *HammerBro {
	base := NewEnemyBase(EnemyTypeHammerBro, x, y)
	
	// Hammer Bro-specific configuration
	base.Health = 2
	base.MaxHealth = 2
	base.Damage = 0 // Hammers do the damage
	base.MoveSpeed = 40.0 // 40 pixels/second (slow movement)
	base.Sprite = 'H'
	base.Color = tcell.ColorGreen
	base.SpriteStyle = tcell.StyleDefault.Foreground(tcell.ColorGreen).Bold(true)
	
	hammerBro := &HammerBro{
		EnemyBase:     base,
		throwTimer:    0.0,
		throwCooldown: 2.0, // Throw every 2 seconds
		movementTimer: 0.0,
		isStanding:    true,
	}
	
	hammerBro.State = EnemyStateIdle
	
	return hammerBro
}

// Update updates the Hammer Bro's state
func (h *HammerBro) Update(deltaTime float64) {
	if h.IsDead || !h.IsActive {
		return
	}
	
	// Update throw timer
	h.throwTimer += deltaTime
	
	// Update movement timer
	h.movementTimer += deltaTime
	
	// Throw hammer periodically
	if h.throwTimer >= h.throwCooldown {
		h.ThrowHammer()
		h.throwTimer = 0.0
	}
	
	// Occasional random movement (every 3-5 seconds)
	if h.movementTimer >= 3.0 && rand.Float64() < 0.3 {
		h.randomMovement()
		h.movementTimer = 0.0
	}
	
	// Update base (applies physics and position updates)
	h.EnemyBase.Update(deltaTime)
}

// UpdateWithPlayerPosition updates Hammer Bro with player position for targeting
func (h *HammerBro) UpdateWithPlayerPosition(deltaTime float64, playerX, playerY float64) {
	if h.IsDead || !h.IsActive {
		return
	}
	
	// Update throw timer
	h.throwTimer += deltaTime
	
	// Update movement timer
	h.movementTimer += deltaTime
	
	// Throw hammer periodically toward player
	if h.throwTimer >= h.throwCooldown {
		h.ThrowHammerAtPlayer(playerX, playerY)
		h.throwTimer = 0.0
	}
	
	// Occasional random movement
	if h.movementTimer >= 3.0 && rand.Float64() < 0.3 {
		h.randomMovement()
		h.movementTimer = 0.0
	}
	
	// Update base
	h.EnemyBase.Update(deltaTime)
}

// randomMovement applies random movement to Hammer Bro
func (h *HammerBro) randomMovement() {
	// Random chance to move left, right, or stay still
	roll := rand.Float64()
	
	if roll < 0.33 {
		// Move left
		h.Velocity.X = -h.MoveSpeed
		h.FacingRight = false
		h.isStanding = false
	} else if roll < 0.66 {
		// Move right
		h.Velocity.X = h.MoveSpeed
		h.FacingRight = true
		h.isStanding = false
	} else {
		// Stand still
		h.Velocity.X = 0
		h.isStanding = true
	}
}

// ThrowHammer spawns a hammer projectile
func (h *HammerBro) ThrowHammer() {
	if h.hammerSpawner != nil {
		// Throw hammer in facing direction with arc
		vx := 200.0
		if !h.FacingRight {
			vx = -200.0
		}
		vy := -150.0 // Upward arc
		
		// Spawn hammer slightly in front of Hammer Bro
		spawnX := h.Position.X
		if h.FacingRight {
			spawnX += 10.0
		} else {
			spawnX -= 10.0
		}
		spawnY := h.Position.Y - 5.0
		
		h.hammerSpawner(spawnX, spawnY, vx, vy)
	}
}

// ThrowHammerAtPlayer spawns a hammer projectile aimed at player
func (h *HammerBro) ThrowHammerAtPlayer(playerX, playerY float64) {
	if h.hammerSpawner != nil {
		// Calculate direction to player
		deltaX := playerX - h.Position.X
		deltaY := playerY - h.Position.Y
		
		// Determine throw direction
		vx := 200.0
		if deltaX < 0 {
			vx = -200.0
			h.FacingRight = false
		} else {
			h.FacingRight = true
		}
		
		// Add upward arc (higher if player is far)
		vy := -150.0
		if deltaY < -50 {
			vy = -200.0 // Higher arc for elevated targets
		}
		
		// Spawn hammer
		spawnX := h.Position.X
		if h.FacingRight {
			spawnX += 10.0
		} else {
			spawnX -= 10.0
		}
		spawnY := h.Position.Y - 5.0
		
		h.hammerSpawner(spawnX, spawnY, vx, vy)
	}
}

// SetHammerSpawner sets the callback function for spawning hammer projectiles
func (h *HammerBro) SetHammerSpawner(spawner func(x, y, vx, vy float64)) {
	h.hammerSpawner = spawner
}

// SetTargetPlayer sets the player to target
func (h *HammerBro) SetTargetPlayer(player interface{}) {
	h.targetPlayer = player
}

// OnStomp is called when the player stomps on Hammer Bro
func (h *HammerBro) OnStomp(player interface{}) {
	// Hammer Bro can be stomped (takes 1 damage)
	h.TakeDamage(1)
}

// OnHitByProjectile is called when hit by a projectile
func (h *HammerBro) OnHitByProjectile(projectile interface{}) {
	// Hammer Bro takes damage from projectiles
	h.TakeDamage(1)
}

// OnHitByShell is called when hit by a kicked Koopa shell
func (h *HammerBro) OnHitByShell(shell interface{}) {
	// Hammer Bro takes damage from shells
	h.TakeDamage(1)
}

// OnCollideWithPlayer is called when colliding with the player
func (h *HammerBro) OnCollideWithPlayer(player interface{}) {
	if h.IsDead {
		return
	}
	
	// Hammer Bro doesn't directly damage player (hammers do)
	// But collision could still be handled by game engine
}

// OnCollideWithWall is called when colliding with a wall
func (h *HammerBro) OnCollideWithWall() {
	// Stop movement when hitting a wall
	h.Velocity.X = 0
	h.isStanding = true
}

// GetThrowTimer returns the current throw timer value
func (h *HammerBro) GetThrowTimer() float64 {
	return h.throwTimer
}

// ResetThrowTimer resets the throw timer
func (h *HammerBro) ResetThrowTimer() {
	h.throwTimer = 0.0
}

// IsStanding returns whether Hammer Bro is standing still
func (h *HammerBro) IsStanding() bool {
	return h.isStanding
}

// CanBeStomped returns whether Hammer Bro can be stomped
func (h *HammerBro) CanBeStomped() bool {
	return true // Hammer Bro can be stomped
}

// GetSprite returns the current sprite character for rendering
func (h *HammerBro) GetSprite() rune {
	return h.Sprite
}

// GetSpriteStyle returns the current sprite style for rendering
func (h *HammerBro) GetSpriteStyle() tcell.Style {
	return h.SpriteStyle
}
