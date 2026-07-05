package entities

import (
	"github.com/gdamore/tcell/v2"
)

// Vector2 represents a 2D vector for position and velocity
type Vector2 struct {
	X float64
	Y float64
}

// EnemyState represents the current state of an enemy
type EnemyState int

const (
	EnemyStateIdle EnemyState = iota
	EnemyStatePatrol
	EnemyStateChase
	EnemyStateAttack
	EnemyStateDead
	EnemyStateStunned
)

// EnemyType represents the type of enemy
type EnemyType int

const (
	EnemyTypeGoomba EnemyType = iota
	EnemyTypeKoopa
	EnemyTypePiranhaPlant
	EnemyTypeLakitu
	EnemyTypeSpiny
	EnemyTypeBoo
	EnemyTypeHammerBro
)

// EnemyBase provides the base functionality for all enemies
type EnemyBase struct {
	// Identity
	ID       int
	Type     EnemyType
	IsActive bool
	IsDead   bool

	// Transform
	Position Vector2
	Velocity Vector2
	Size     Vector2

	// Stats
	Health    int
	MaxHealth int
	Damage    int
	MoveSpeed float64

	// State
	State         EnemyState
	FacingRight   bool
	OnGround      bool
	IsFlying      bool
	StateTimer    float64
	InvulnTimer   float64

	// Rendering
	Sprite      rune
	Color       tcell.Color
	SpriteStyle tcell.Style

	// AI
	PatrolLeft  float64
	PatrolRight float64
	DetectRange float64
}

// NewEnemyBase creates a new enemy base with default values
func NewEnemyBase(enemyType EnemyType, x, y float64) *EnemyBase {
	return &EnemyBase{
		Type:        enemyType,
		IsActive:    true,
		IsDead:      false,
		Position:    Vector2{X: x, Y: y},
		Velocity:    Vector2{X: 0, Y: 0},
		Size:        Vector2{X: 16, Y: 16},
		Health:      1,
		MaxHealth:   1,
		Damage:      1,
		MoveSpeed:   80.0,
		State:       EnemyStatePatrol,
		FacingRight: true,
		OnGround:    false,
		DetectRange: 100.0,
	}
}

// Update updates the enemy's state
func (e *EnemyBase) Update(deltaTime float64) {
	if e.IsDead || !e.IsActive {
		return
	}

	// Update timers
	if e.StateTimer > 0 {
		e.StateTimer -= deltaTime
	}
	if e.InvulnTimer > 0 {
		e.InvulnTimer -= deltaTime
	}

	// Apply gravity (not for flying enemies)
	if !e.OnGround && !e.IsFlying {
		e.Velocity.Y += 980.0 * deltaTime // Gravity
	}

	// Update position
	e.Position.X += e.Velocity.X * deltaTime
	e.Position.Y += e.Velocity.Y * deltaTime
}

// TakeDamage applies damage to the enemy
func (e *EnemyBase) TakeDamage(damage int) {
	if e.InvulnTimer > 0 || e.IsDead {
		return
	}

	e.Health -= damage
	if e.Health <= 0 {
		e.Die()
	} else {
		e.InvulnTimer = 0.5 // 0.5 seconds of invulnerability
	}
}

// Die handles enemy death
func (e *EnemyBase) Die() {
	e.IsDead = true
	e.State = EnemyStateDead
	e.Velocity.X = 0
	e.Velocity.Y = 0
}

// OnStomp is called when the player stomps on the enemy
func (e *EnemyBase) OnStomp(player interface{}) {
	e.TakeDamage(1)
}

// OnHitByProjectile is called when hit by a projectile
func (e *EnemyBase) OnHitByProjectile(projectile interface{}) {
	e.TakeDamage(1)
}

// OnCollideWithPlayer is called when colliding with the player
func (e *EnemyBase) OnCollideWithPlayer(player interface{}) {
	// Default behavior: damage the player
	// This will be overridden by specific enemy types
}

// GetBounds returns the enemy's bounding box
func (e *EnemyBase) GetBounds() (x, y, width, height float64) {
	return e.Position.X, e.Position.Y, e.Size.X, e.Size.Y
}

// SetPosition sets the enemy's position
func (e *EnemyBase) SetPosition(x, y float64) {
	e.Position.X = x
	e.Position.Y = y
}

// SetVelocity sets the enemy's velocity
func (e *EnemyBase) SetVelocity(x, y float64) {
	e.Velocity.X = x
	e.Velocity.Y = y
}

// IsAlive returns whether the enemy is alive
func (e *EnemyBase) IsAlive() bool {
	return !e.IsDead && e.IsActive
}

// Deactivate deactivates the enemy
func (e *EnemyBase) Deactivate() {
	e.IsActive = false
}

// Activate activates the enemy
func (e *EnemyBase) Activate() {
	e.IsActive = true
}
