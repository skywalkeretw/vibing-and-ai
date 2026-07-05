package entities

import (
	"math"

	"github.com/lukeroy/go-terminal-platformer/internal/engine"
	"github.com/lukeroy/go-terminal-platformer/internal/input"
)

// PlayerState represents the current state of the player
type PlayerState int

const (
	StateIdle PlayerState = iota
	StateRunning
	StateJumping
	StateFalling
	StateCrouching
	StateHurt
	StateDead
)

// String returns the string representation of a player state
func (ps PlayerState) String() string {
	switch ps {
	case StateIdle:
		return "Idle"
	case StateRunning:
		return "Running"
	case StateJumping:
		return "Jumping"
	case StateFalling:
		return "Falling"
	case StateCrouching:
		return "Crouching"
	case StateHurt:
		return "Hurt"
	case StateDead:
		return "Dead"
	default:
		return "Unknown"
	}
}

// PlayerEntity represents a player character
type PlayerEntity struct {
	id          int
	position    engine.Vector2
	velocity    engine.Vector2
	physicsBody *engine.PhysicsBody

	// Movement
	moveSpeed    float64
	jumpForce    float64
	maxJumpTime  float64
	jumpTimeLeft float64

	// State
	state     PlayerState
	facing    Direction
	grounded  bool
	crouching bool
	active    bool

	// Combat
	lives        int
	invulnerable bool
	invulnTime   float64
	powerUp      PowerUpType
	ammo         int
	coins        int

	// Animation
	sprite    *Sprite
	animFrame int
	animTime  float64

	// Input
	inputActions []input.InputAction

	// Physics reference
	physics *engine.PhysicsEngine
}

// NewPlayer creates a new player entity
func NewPlayer() *PlayerEntity {
	return &PlayerEntity{
		active: true,
		lives:  5,
		state:  StateIdle,
		facing: DirectionRight,
	}
}

// Initialize sets up the player with starting values
func (p *PlayerEntity) Initialize(id int, startPos engine.Vector2, physics *engine.PhysicsEngine) {
	p.id = id
	p.position = startPos
	p.moveSpeed = 240.0 // 8 chars/sec at 30 FPS
	p.jumpForce = -400.0
	p.maxJumpTime = 0.3
	p.lives = 5
	p.state = StateIdle
	p.facing = DirectionRight
	p.physics = physics
	p.active = true

	// Create physics body
	p.physicsBody = &engine.PhysicsBody{
		Entity:       p,
		Position:     startPos,
		Velocity:     engine.Vector2{X: 0, Y: 0},
		Acceleration: engine.Vector2{X: 0, Y: 0},
		Mass:         1.0,
		Friction:     0.8,
		Restitution:  0.0,
		Collider:     engine.NewAABBCollider(startPos.X, startPos.Y, 16, 32, engine.LayerPlayer),
		Layer:        engine.LayerPlayer,
		Enabled:      true,
		GravityScale: 1.0,
	}
	physics.AddBody(p.physicsBody)

	// Load sprite based on player ID
	if id == 1 {
		p.sprite = p.createPlayer1Sprite()
	} else {
		p.sprite = p.createPlayer2Sprite()
	}
}

// Update updates the player state
func (p *PlayerEntity) Update(deltaTime float64) {
	if p.state == StateDead {
		return
	}

	// Update invulnerability
	if p.invulnerable {
		p.invulnTime -= deltaTime
		if p.invulnTime <= 0 {
			p.invulnerable = false
		}
	}

	// Process input
	p.handleInput(deltaTime)

	// Update state
	p.updateState()

	// Update animation
	p.updateAnimation(deltaTime)

	// Sync with physics body
	p.position = p.physicsBody.Position
	p.grounded = p.physicsBody.Grounded
}

// UpdateWithInput updates the player with input actions
func (p *PlayerEntity) UpdateWithInput(deltaTime float64, inputActions []input.InputAction) {
	p.inputActions = inputActions
	p.Update(deltaTime)
}

// handleInput processes input actions
func (p *PlayerEntity) handleInput(deltaTime float64) {
	// Reset horizontal velocity if no movement input
	hasMovementInput := false

	for _, action := range p.inputActions {
		switch action {
		case input.ActionMoveLeft:
			p.moveLeft()
			hasMovementInput = true
		case input.ActionMoveRight:
			p.moveRight()
			hasMovementInput = true
		case input.ActionJump:
			p.jump()
		case input.ActionCrouch:
			p.crouch()
		case input.ActionShoot:
			p.shoot()
		}
	}

	// Stop horizontal movement if no input
	if !hasMovementInput && p.grounded && p.state != StateCrouching {
		p.physicsBody.Velocity.X = 0
	}

	// Variable jump height
	if p.state == StateJumping && p.jumpTimeLeft > 0 {
		if p.hasAction(input.ActionJump) {
			p.physicsBody.Velocity.Y = p.jumpForce
			p.jumpTimeLeft -= deltaTime
		} else {
			p.jumpTimeLeft = 0
		}
	}

	// Reset crouching if not holding crouch
	if p.crouching && !p.hasAction(input.ActionCrouch) {
		p.crouching = false
	}
}

// hasAction checks if an action is in the input actions
func (p *PlayerEntity) hasAction(action input.InputAction) bool {
	for _, a := range p.inputActions {
		if a == action {
			return true
		}
	}
	return false
}

// moveLeft moves the player left
func (p *PlayerEntity) moveLeft() {
	if p.state != StateCrouching {
		p.physicsBody.Velocity.X = -p.moveSpeed
		p.facing = DirectionLeft
	}
}

// moveRight moves the player right
func (p *PlayerEntity) moveRight() {
	if p.state != StateCrouching {
		p.physicsBody.Velocity.X = p.moveSpeed
		p.facing = DirectionRight
	}
}

// jump makes the player jump
func (p *PlayerEntity) jump() {
	if p.grounded && p.state != StateCrouching {
		p.physicsBody.Velocity.Y = p.jumpForce
		p.jumpTimeLeft = p.maxJumpTime
		p.state = StateJumping
		// TODO: Play jump sound
	}
}

// crouch makes the player crouch
func (p *PlayerEntity) crouch() {
	if p.grounded {
		p.crouching = true
		p.state = StateCrouching
		p.physicsBody.Velocity.X = 0
	} else {
		// Fast fall
		p.physicsBody.Velocity.Y += 200
	}
}

// shoot fires a projectile
func (p *PlayerEntity) shoot() {
	if p.powerUp == PowerUpFire && p.ammo > 0 {
		p.createProjectile()
		p.ammo--
		if p.ammo <= 0 {
			p.powerUp = PowerUpNone
		}
	}
}

// createProjectile creates a projectile (placeholder)
func (p *PlayerEntity) createProjectile() {
	// TODO: Implement projectile creation when projectile system is ready
	// This will be implemented in issue #16
}

// updateState updates the player's state based on current conditions
func (p *PlayerEntity) updateState() {
	if p.state == StateHurt {
		return // Wait for hurt animation to finish
	}

	if !p.grounded {
		if p.physicsBody.Velocity.Y < 0 {
			p.state = StateJumping
		} else {
			p.state = StateFalling
		}
	} else if p.crouching {
		p.state = StateCrouching
	} else if math.Abs(p.physicsBody.Velocity.X) > 10 {
		p.state = StateRunning
	} else {
		p.state = StateIdle
	}
}

// updateAnimation updates the animation frame
func (p *PlayerEntity) updateAnimation(deltaTime float64) {
	p.animTime += deltaTime

	frameRate := 0.1 // 10 FPS for animations
	if p.animTime >= frameRate {
		p.animTime = 0
		p.animFrame++

		// Get frame count for current state
		frameCount := p.getFrameCount(p.state)
		if p.animFrame >= frameCount {
			p.animFrame = 0
			
			// Exit hurt state after animation completes
			if p.state == StateHurt {
				p.state = StateIdle
			}
		}
	}
}

// getFrameCount returns the number of animation frames for a state
func (p *PlayerEntity) getFrameCount(state PlayerState) int {
	switch state {
	case StateIdle:
		return 4
	case StateRunning:
		return 6
	case StateJumping, StateFalling:
		return 1
	case StateCrouching:
		return 1
	case StateHurt:
		return 3
	}
	return 1
}

// TakeDamage applies damage to the player
func (p *PlayerEntity) TakeDamage() {
	if p.invulnerable || p.state == StateDead {
		return
	}

	// Check for shield power-up
	if p.powerUp == PowerUpShield {
		p.powerUp = PowerUpNone
		return
	}

	p.lives--
	p.invulnerable = true
	p.invulnTime = 2.0 // 2 seconds invulnerability
	p.state = StateHurt

	// Knockback
	p.physicsBody.Velocity.Y = -200
	if p.facing == DirectionRight {
		p.physicsBody.Velocity.X = -150
	} else {
		p.physicsBody.Velocity.X = 150
	}

	if p.lives <= 0 {
		p.Die()
	}
}

// Die kills the player
func (p *PlayerEntity) Die() {
	p.state = StateDead
	p.physicsBody.Velocity = engine.Vector2{X: 0, Y: 0}
	// TODO: Trigger death animation and respawn logic
}

// AddLife adds a life to the player
func (p *PlayerEntity) AddLife() {
	p.lives++
}

// AddCoins adds coins to the player
func (p *PlayerEntity) AddCoins(amount int) {
	p.coins += amount
	
	// Every 100 coins = 1 extra life
	if p.coins >= 100 {
		livesToAdd := p.coins / 100
		p.lives += livesToAdd
		p.coins = p.coins % 100
	}
}

// GetCoins returns the player's coin count
func (p *PlayerEntity) GetCoins() int {
	return p.coins
}

// ApplyPowerUp applies a power-up to the player
func (p *PlayerEntity) ApplyPowerUp(powerUp PowerUpType) {
	// Remove current power-up effects
	p.RemovePowerUp()
	
	p.powerUp = powerUp

	switch powerUp {
	case PowerUpFire:
		p.ammo = 10
	case PowerUpStar:
		p.invulnerable = true
		p.invulnTime = 10.0
		p.moveSpeed = 360.0 // 1.5x speed
	case PowerUpSpeedBoots:
		p.moveSpeed = 360.0
	case PowerUpSuperJump:
		p.jumpForce = -600.0
	case PowerUpShield:
		// Shield handled in TakeDamage
	}
}

// RemovePowerUp removes the current power-up
func (p *PlayerEntity) RemovePowerUp() {
	p.powerUp = PowerUpNone
	p.moveSpeed = 240.0
	p.jumpForce = -400.0
	p.ammo = 0
}

// Render renders the player
func (p *PlayerEntity) Render(renderer Renderer) {
	// Get sprite for current state
	sprite := p.getSpriteForState()

	// Flip sprite based on facing direction
	if p.facing == DirectionLeft {
		sprite = sprite.FlipHorizontal()
	}

	// Flash when invulnerable
	if p.invulnerable && int(p.invulnTime*10)%2 == 0 {
		return // Skip rendering for flashing effect
	}

	renderer.DrawSprite(int(p.position.X), int(p.position.Y), sprite)
}

// getSpriteForState returns the appropriate sprite for the current state
func (p *PlayerEntity) getSpriteForState() *Sprite {
	// For now, return the base sprite
	// TODO: Implement different sprites for different states
	return p.sprite
}

// createPlayer1Sprite creates the sprite for player 1 (blue/cyan)
func (p *PlayerEntity) createPlayer1Sprite() *Sprite {
	sprite := NewSprite(2, 2)
	sprite.Data[0][0] = '☺'
	sprite.Data[0][1] = ' '
	sprite.Data[1][0] = '║'
	sprite.Data[1][1] = ' '
	sprite.Color = 1 // Blue/Cyan color code
	return sprite
}

// createPlayer2Sprite creates the sprite for player 2 (magenta/purple)
func (p *PlayerEntity) createPlayer2Sprite() *Sprite {
	sprite := NewSprite(2, 2)
	sprite.Data[0][0] = '☻'
	sprite.Data[0][1] = ' '
	sprite.Data[1][0] = '║'
	sprite.Data[1][1] = ' '
	sprite.Color = 2 // Magenta/Purple color code
	return sprite
}

// Entity interface implementation

// SetPosition sets the player's position
func (p *PlayerEntity) SetPosition(pos engine.Vector2) {
	p.position = pos
	if p.physicsBody != nil {
		p.physicsBody.Position = pos
	}
}

// GetPosition returns the player's position
func (p *PlayerEntity) GetPosition() engine.Vector2 {
	return p.position
}

// IsActive returns whether the player is active
func (p *PlayerEntity) IsActive() bool {
	return p.active
}

// SetActive sets the player's active state
func (p *PlayerEntity) SetActive(active bool) {
	p.active = active
}

// OnCollision handles collision with another entity
func (p *PlayerEntity) OnCollision(other Entity) {
	// TODO: Implement collision handling with enemies, collectibles, etc.
}

// Getters

// GetID returns the player's ID
func (p *PlayerEntity) GetID() int {
	return p.id
}

// GetState returns the player's current state
func (p *PlayerEntity) GetState() PlayerState {
	return p.state
}

// GetFacing returns the player's facing direction
func (p *PlayerEntity) GetFacing() Direction {
	return p.facing
}

// GetLives returns the player's remaining lives
func (p *PlayerEntity) GetLives() int {
	return p.lives
}

// GetPowerUp returns the player's current power-up
func (p *PlayerEntity) GetPowerUp() PowerUpType {
	return p.powerUp
}

// GetAmmo returns the player's remaining ammo
func (p *PlayerEntity) GetAmmo() int {
	return p.ammo
}

// IsGrounded returns whether the player is on the ground
func (p *PlayerEntity) IsGrounded() bool {
	return p.grounded
}

// IsInvulnerable returns whether the player is invulnerable
func (p *PlayerEntity) IsInvulnerable() bool {
	return p.invulnerable
}

// IsDead returns whether the player is dead
func (p *PlayerEntity) IsDead() bool {
	return p.state == StateDead
}
