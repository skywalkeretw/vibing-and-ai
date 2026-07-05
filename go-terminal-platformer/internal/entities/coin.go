package entities

import (
	"github.com/lukeroy/go-terminal-platformer/internal/engine"
)

// Coin represents a collectible coin
type Coin struct {
	position    engine.Vector2
	velocity    engine.Vector2
	physicsBody *engine.PhysicsBody
	active      bool
	collected   bool
	value       int

	// Animation
	animFrame int
	animTime  float64

	// Physics reference
	physics *engine.PhysicsEngine
}

// NewCoin creates a new coin
func NewCoin() *Coin {
	return &Coin{
		active: true,
		value:  1,
	}
}

// Initialize sets up the coin with physics integration
func (c *Coin) Initialize(pos engine.Vector2, physics *engine.PhysicsEngine) {
	c.position = pos
	c.active = true
	c.value = 1
	c.physics = physics

	// Coins have simple physics
	c.physicsBody = &engine.PhysicsBody{
		Entity:       c,
		Position:     pos,
		Velocity:     engine.Vector2{X: 0, Y: 0},
		Acceleration: engine.Vector2{X: 0, Y: 0},
		Mass:         0.1,
		Friction:     0.95,
		Restitution:  0.5,
		Collider:     engine.NewAABBCollider(pos.X, pos.Y, 8, 8, engine.LayerCollectible),
		Layer:        engine.LayerCollectible,
		Enabled:      true,
		GravityScale: 1.0,
	}
	physics.AddBody(c.physicsBody)
}

// Update updates the coin state
func (c *Coin) Update(deltaTime float64, players []*PlayerEntity) {
	if !c.active || c.collected {
		return
	}

	// Check collection with all players
	for _, player := range players {
		if player != nil && player.IsActive() && c.checkCollision(player) {
			c.collect(player)
			return
		}
	}

	// Spin animation
	c.animTime += deltaTime
	if c.animTime >= 0.1 { // 10 FPS animation
		c.animTime = 0
		c.animFrame = (c.animFrame + 1) % 4
	}

	// Sync position from physics
	if c.physicsBody != nil {
		c.position = c.physicsBody.Position
	}
}

// checkCollision checks if the coin collides with a player
func (c *Coin) checkCollision(player *PlayerEntity) bool {
	if c.physicsBody == nil || player.physicsBody == nil {
		return false
	}
	return c.physicsBody.Collider.Intersects(player.physicsBody.Collider)
}

// collect handles collection by a player
func (c *Coin) collect(player *PlayerEntity) {
	c.collected = true
	c.active = false
	
	// Add coins to player
	player.AddCoins(c.value)
	
	// TODO: Play coin sound
	// TODO: Award 10 points
	
	// Remove from physics
	if c.physicsBody != nil && c.physics != nil {
		c.physics.RemoveBody(c.physicsBody)
	}
}

// Render renders the coin
func (c *Coin) Render(renderer Renderer) {
	if !c.active || c.collected {
		return
	}

	// Get sprite based on animation frame
	sprite := c.getSprite()
	renderer.DrawSprite(int(c.position.X), int(c.position.Y), sprite)
}

// getSprite returns the sprite for the current animation frame
func (c *Coin) getSprite() *Sprite {
	sprite := NewSprite(1, 1)
	
	// Spinning coin animation
	switch c.animFrame {
	case 0:
		sprite.Data[0][0] = 'O'
	case 1:
		sprite.Data[0][0] = 'o'
	case 2:
		sprite.Data[0][0] = '·'
	case 3:
		sprite.Data[0][0] = 'o'
	}
	
	sprite.Color = 4 // Yellow/Gold
	return sprite
}

// Entity interface implementation

// SetPosition sets the coin's position
func (c *Coin) SetPosition(pos engine.Vector2) {
	c.position = pos
	if c.physicsBody != nil {
		c.physicsBody.Position = pos
	}
}

// GetPosition returns the coin's position
func (c *Coin) GetPosition() engine.Vector2 {
	return c.position
}

// IsActive returns whether the coin is active
func (c *Coin) IsActive() bool {
	return c.active && !c.collected
}

// SetActive sets the coin's active state
func (c *Coin) SetActive(active bool) {
	c.active = active
}

// OnCollision handles collision with another entity
func (c *Coin) OnCollision(other Entity) {
	// Coins don't react to collisions with other entities
	// Collection is handled in Update()
}

// Getters

// GetValue returns the coin's value
func (c *Coin) GetValue() int {
	return c.value
}

// IsCollected returns whether the coin has been collected
func (c *Coin) IsCollected() bool {
	return c.collected
}

// GetPhysicsBody returns the physics body
func (c *Coin) GetPhysicsBody() *engine.PhysicsBody {
	return c.physicsBody
}

// SetValue sets the coin's value
func (c *Coin) SetValue(value int) {
	c.value = value
}

// Deactivate deactivates the coin
func (c *Coin) Deactivate() {
	c.active = false
	if c.physicsBody != nil {
		c.physicsBody.Enabled = false
	}
}
