package entities

import (
	"math"

	"github.com/lukeroy/go-terminal-platformer/internal/engine"
)

// PowerUp represents a collectible power-up item
type PowerUp struct {
	powerUpType PowerUpType
	position    engine.Vector2
	velocity    engine.Vector2
	physicsBody *engine.PhysicsBody
	active      bool
	collected   bool

	// Behavior
	floating    bool
	floatTime   float64
	floatOffset float64

	// Animation
	sprite    *Sprite
	animFrame int
	animTime  float64

	// Physics reference
	physics *engine.PhysicsEngine
}

// NewPowerUp creates a new power-up
func NewPowerUp() *PowerUp {
	return &PowerUp{
		active:   true,
		floating: true,
	}
}

// Initialize sets up the power-up with physics integration
func (p *PowerUp) Initialize(powerUpType PowerUpType, pos engine.Vector2, physics *engine.PhysicsEngine) {
	p.powerUpType = powerUpType
	p.position = pos
	p.active = true
	p.floating = true
	p.physics = physics

	// Create physics body
	p.physicsBody = &engine.PhysicsBody{
		Entity:       p,
		Position:     pos,
		Velocity:     engine.Vector2{X: 0, Y: 0},
		Acceleration: engine.Vector2{X: 0, Y: 0},
		Mass:         0.5,
		Friction:     0.9,
		Restitution:  0.3,
		Collider:     engine.NewAABBCollider(pos.X, pos.Y, 16, 16, engine.LayerCollectible),
		Layer:        engine.LayerCollectible,
		Enabled:      true,
		GravityScale: 0.5, // Light gravity for power-ups
	}
	physics.AddBody(p.physicsBody)

	// Set sprite based on type
	p.sprite = p.getSpriteForType(powerUpType)
}

// Update updates the power-up state
func (p *PowerUp) Update(deltaTime float64, players []*PlayerEntity) {
	if !p.active || p.collected {
		return
	}

	// Floating animation
	if p.floating {
		p.floatTime += deltaTime * 3
		p.floatOffset = math.Sin(p.floatTime) * 5
	}

	// Check collection with all players
	for _, player := range players {
		if player != nil && player.IsActive() && p.checkCollision(player) {
			p.collect(player)
			return
		}
	}

	// Update animation
	p.updateAnimation(deltaTime)

	// Sync position from physics
	if p.physicsBody != nil {
		p.position = p.physicsBody.Position
	}
}

// checkCollision checks if the power-up collides with a player
func (p *PowerUp) checkCollision(player *PlayerEntity) bool {
	if p.physicsBody == nil || player.physicsBody == nil {
		return false
	}
	return p.physicsBody.Collider.Intersects(player.physicsBody.Collider)
}

// collect handles collection by a player
func (p *PowerUp) collect(player *PlayerEntity) {
	p.collected = true
	p.active = false
	
	// Apply power-up effect to player
	player.ApplyPowerUp(p.powerUpType)
	
	// TODO: Play collection sound
	// TODO: Award points (100 points for power-up)
	
	// Remove from physics
	if p.physicsBody != nil && p.physics != nil {
		p.physics.RemoveBody(p.physicsBody)
	}
}

// updateAnimation updates the animation frame
func (p *PowerUp) updateAnimation(deltaTime float64) {
	p.animTime += deltaTime
	if p.animTime >= 0.15 { // ~6.67 FPS animation
		p.animTime = 0
		p.animFrame = (p.animFrame + 1) % 4
	}
}

// getSpriteForType returns the sprite for a power-up type
func (p *PowerUp) getSpriteForType(powerUpType PowerUpType) *Sprite {
	sprite := NewSprite(1, 1)
	
	switch powerUpType {
	case PowerUpFire:
		sprite.Data[0][0] = 'F'
		sprite.Color = 3 // Red/Orange
	case PowerUpStar:
		sprite.Data[0][0] = '★'
		sprite.Color = 4 // Yellow
	case PowerUpMushroom:
		sprite.Data[0][0] = '♠'
		sprite.Color = 2 // Green
	case PowerUpSpeedBoots:
		sprite.Data[0][0] = 'B'
		sprite.Color = 1 // Blue
	case PowerUpSuperJump:
		sprite.Data[0][0] = 'J'
		sprite.Color = 5 // Cyan
	case PowerUpShield:
		sprite.Data[0][0] = 'S'
		sprite.Color = 6 // White
	default:
		sprite.Data[0][0] = '?'
		sprite.Color = 7 // Gray
	}
	
	return sprite
}

// Render renders the power-up
func (p *PowerUp) Render(renderer Renderer) {
	if !p.active || p.collected {
		return
	}

	// Apply floating offset
	renderY := int(p.position.Y + p.floatOffset)
	renderer.DrawSprite(int(p.position.X), renderY, p.sprite)
}

// Entity interface implementation

// SetPosition sets the power-up's position
func (p *PowerUp) SetPosition(pos engine.Vector2) {
	p.position = pos
	if p.physicsBody != nil {
		p.physicsBody.Position = pos
	}
}

// GetPosition returns the power-up's position
func (p *PowerUp) GetPosition() engine.Vector2 {
	return p.position
}

// IsActive returns whether the power-up is active
func (p *PowerUp) IsActive() bool {
	return p.active && !p.collected
}

// SetActive sets the power-up's active state
func (p *PowerUp) SetActive(active bool) {
	p.active = active
}

// OnCollision handles collision with another entity
func (p *PowerUp) OnCollision(other Entity) {
	// Power-ups don't react to collisions with other entities
	// Collection is handled in Update()
}

// Getters

// GetType returns the power-up type
func (p *PowerUp) GetType() PowerUpType {
	return p.powerUpType
}

// IsCollected returns whether the power-up has been collected
func (p *PowerUp) IsCollected() bool {
	return p.collected
}

// GetPhysicsBody returns the physics body
func (p *PowerUp) GetPhysicsBody() *engine.PhysicsBody {
	return p.physicsBody
}

// Deactivate deactivates the power-up
func (p *PowerUp) Deactivate() {
	p.active = false
	if p.physicsBody != nil {
		p.physicsBody.Enabled = false
	}
}
