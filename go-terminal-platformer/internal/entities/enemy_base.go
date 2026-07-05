package entities

import (
	"math"

	"github.com/gdamore/tcell/v2"
	"github.com/lukeroy/go-terminal-platformer/internal/engine"
)

// Direction represents facing direction
type Direction int

const (
	DirectionLeft  Direction = -1
	DirectionRight Direction = 1
)

// EnemyState represents the current state of an enemy
type EnemyState int

const (
	EnemyStateIdle EnemyState = iota
	EnemyStatePatrol
	EnemyStateChase
	EnemyStateAttack
	EnemyStateHurt
	EnemyStateDead
	EnemyStateStunned
)

// AIState represents the AI behavior state
type AIState int

const (
	AIPatrol AIState = iota
	AIChase
	AIAttack
	AIFlee
)

// EnemyType represents the type of enemy
type EnemyType int

const (
	EnemyGoomba EnemyType = iota
	EnemyKoopa
	EnemyPiranha
	EnemyLakitu
	EnemySpiny
	EnemyBoo
	EnemyHammerBro
)

// String returns the string representation of EnemyType
func (et EnemyType) String() string {
	switch et {
	case EnemyGoomba:
		return "Goomba"
	case EnemyKoopa:
		return "Koopa"
	case EnemyPiranha:
		return "Piranha"
	case EnemyLakitu:
		return "Lakitu"
	case EnemySpiny:
		return "Spiny"
	case EnemyBoo:
		return "Boo"
	case EnemyHammerBro:
		return "HammerBro"
	default:
		return "Unknown"
	}
}

// Enemy interface defines the contract for all enemy types
type Enemy interface {
	Initialize(pos engine.Vector2, physics *engine.PhysicsEngine)
	Update(deltaTime float64, players []Player)
	TakeDamage(damage int)
	Die()
	GetPosition() engine.Vector2
	SetPosition(pos engine.Vector2)
	GetCollider() engine.Collider
	IsActive() bool
	IsAlive() bool
	GetType() EnemyType
	OnCollision(other interface{})
}

// Player interface for enemy interaction (minimal interface to avoid circular dependency)
type Player interface {
	GetPosition() engine.Vector2
	TakeDamage(damage int)
	IsAlive() bool
	GetPhysicsBody() *engine.PhysicsBody
}

// EnemyBase provides the base functionality for all enemies
type EnemyBase struct {
	// Identity
	ID       int
	Type     EnemyType
	Active   bool
	IsDead   bool

	// Physics
	PhysicsBody *engine.PhysicsBody
	Position    engine.Vector2
	Velocity    engine.Vector2

	// Stats
	Health    int
	MaxHealth int
	Damage    int
	MoveSpeed float64

	// State
	State         EnemyState
	AIState       AIState
	Facing        Direction
	OnGround      bool
	IsFlying      bool
	StateTimer    float64
	InvulnTimer   float64
	HurtTimer     float64

	// AI
	Target         Player
	PatrolPoints   []engine.Vector2
	CurrentPoint   int
	DetectRange    float64
	AttackRange    float64
	PatrolLeft     float64
	PatrolRight    float64
	TurnAtEdge     bool
	TurnAtWall     bool

	// Rendering
	Sprite      rune
	Color       tcell.Color
	SpriteStyle tcell.Style
	AnimFrame   int
	AnimTime    float64

	// Physics reference
	physics *engine.PhysicsEngine
}

// NewEnemyBase creates a new enemy base with default values
func NewEnemyBase(enemyType EnemyType, x, y float64) *EnemyBase {
	enemy := &EnemyBase{
		Type:        enemyType,
		Active:      true,
		IsDead:      false,
		Position:    engine.Vector2{X: x, Y: y},
		Velocity:    engine.Vector2{X: 0, Y: 0},
		State:       EnemyStatePatrol,
		AIState:     AIPatrol,
		Facing:      DirectionLeft,
		OnGround:    false,
		DetectRange: 200.0,
		AttackRange: 32.0,
		TurnAtEdge:  true,
		TurnAtWall:  true,
	}

	// Set stats based on enemy type
	enemy.setStatsForType(enemyType)

	return enemy
}

// setStatsForType sets enemy stats based on type
func (e *EnemyBase) setStatsForType(enemyType EnemyType) {
	switch enemyType {
	case EnemyGoomba:
		e.Health = 1
		e.MaxHealth = 1
		e.Damage = 1
		e.MoveSpeed = 60.0
		e.IsFlying = false
		e.Sprite = 'G'
		e.Color = tcell.ColorBrown

	case EnemyKoopa:
		e.Health = 2
		e.MaxHealth = 2
		e.Damage = 1
		e.MoveSpeed = 80.0
		e.IsFlying = false
		e.Sprite = 'K'
		e.Color = tcell.ColorGreen

	case EnemyPiranha:
		e.Health = 2
		e.MaxHealth = 2
		e.Damage = 1
		e.MoveSpeed = 0.0 // Stationary
		e.IsFlying = false
		e.Sprite = 'P'
		e.Color = tcell.ColorRed

	case EnemyLakitu:
		e.Health = 3
		e.MaxHealth = 3
		e.Damage = 1
		e.MoveSpeed = 100.0
		e.IsFlying = true
		e.Sprite = 'L'
		e.Color = tcell.ColorYellow

	case EnemySpiny:
		e.Health = 1
		e.MaxHealth = 1
		e.Damage = 1
		e.MoveSpeed = 70.0
		e.IsFlying = false
		e.Sprite = 'S'
		e.Color = tcell.ColorRed

	case EnemyBoo:
		e.Health = 2
		e.MaxHealth = 2
		e.Damage = 1
		e.MoveSpeed = 90.0
		e.IsFlying = true
		e.Sprite = 'B'
		e.Color = tcell.ColorWhite

	case EnemyHammerBro:
		e.Health = 3
		e.MaxHealth = 3
		e.Damage = 1
		e.MoveSpeed = 50.0
		e.IsFlying = false
		e.Sprite = 'H'
		e.Color = tcell.ColorDarkGreen

	default:
		e.Health = 1
		e.MaxHealth = 1
		e.Damage = 1
		e.MoveSpeed = 60.0
		e.IsFlying = false
		e.Sprite = '?'
		e.Color = tcell.ColorWhite
	}

	e.SpriteStyle = tcell.StyleDefault.Foreground(e.Color)
}

// Initialize sets up the enemy with physics integration
func (e *EnemyBase) Initialize(pos engine.Vector2, physics *engine.PhysicsEngine) {
	e.Position = pos
	e.physics = physics

	// Create collider based on enemy type
	var collider engine.Collider
	if e.IsFlying {
		// Flying enemies have smaller colliders
		collider = engine.NewAABBCollider(pos.X, pos.Y, 12, 12, engine.LayerEnemy)
	} else {
		// Ground enemies have standard colliders
		collider = engine.NewAABBCollider(pos.X, pos.Y, 16, 16, engine.LayerEnemy)
	}

	// Set gravity scale based on enemy type
	gravityScale := 1.0
	if e.IsFlying {
		gravityScale = 0.0
	}

	// Create physics body
	e.PhysicsBody = &engine.PhysicsBody{
		Entity:       e,
		Position:     pos,
		Velocity:     engine.Vector2{X: 0, Y: 0},
		Acceleration: engine.Vector2{X: 0, Y: 0},
		Mass:         1.0,
		Friction:     0.8,
		Restitution:  0.0,
		Grounded:     false,
		Collider:     collider,
		Layer:        engine.LayerEnemy,
		Enabled:      true,
		GravityScale: gravityScale,
	}

	// Add to physics engine
	physics.AddBody(e.PhysicsBody)

	// Set default patrol range if not set
	if e.PatrolLeft == 0 && e.PatrolRight == 0 {
		e.PatrolLeft = pos.X - 100
		e.PatrolRight = pos.X + 100
	}
}

// Update updates the enemy's state
func (e *EnemyBase) Update(deltaTime float64, players []Player) {
	if e.IsDead || !e.IsActive() {
		return
	}

	// Update timers
	if e.StateTimer > 0 {
		e.StateTimer -= deltaTime
	}
	if e.InvulnTimer > 0 {
		e.InvulnTimer -= deltaTime
	}
	if e.HurtTimer > 0 {
		e.HurtTimer -= deltaTime
		if e.HurtTimer <= 0 && e.State == EnemyStateHurt {
			e.State = EnemyStatePatrol
		}
	}

	// Update AI
	e.updateAI(deltaTime, players)

	// Update state
	e.updateState(deltaTime)

	// Update animation
	e.updateAnimation(deltaTime)

	// Sync position from physics
	if e.PhysicsBody != nil {
		e.Position = e.PhysicsBody.Position
		e.Velocity = e.PhysicsBody.Velocity
		e.OnGround = e.PhysicsBody.Grounded
	}
}

// updateAI handles AI behavior
func (e *EnemyBase) updateAI(deltaTime float64, players []Player) {
	if e.State == EnemyStateDead || e.State == EnemyStateHurt {
		return
	}

	// Find nearest player
	nearestPlayer := e.findNearestPlayer(players)
	if nearestPlayer == nil {
		return
	}

	distance := e.Position.Distance(nearestPlayer.GetPosition())

	switch e.AIState {
	case AIPatrol:
		e.patrol(deltaTime)

		// Check if player in detection range
		if distance < e.DetectRange {
			e.AIState = AIChase
			e.Target = nearestPlayer
			e.State = EnemyStateChase
		}

	case AIChase:
		if distance > e.DetectRange*1.5 {
			// Lost player
			e.AIState = AIPatrol
			e.Target = nil
			e.State = EnemyStatePatrol
		} else if distance < e.AttackRange {
			e.AIState = AIAttack
			e.State = EnemyStateAttack
		} else {
			e.chaseTarget(deltaTime)
		}

	case AIAttack:
		if distance > e.AttackRange*1.5 {
			e.AIState = AIChase
			e.State = EnemyStateChase
		} else {
			e.attack(deltaTime)
		}

	case AIFlee:
		// Flee behavior (for specific enemy types)
		e.flee(deltaTime)
	}
}

// patrol implements patrol behavior
func (e *EnemyBase) patrol(deltaTime float64) {
	if e.MoveSpeed == 0 {
		// Stationary enemy (like Piranha Plant)
		e.PhysicsBody.Velocity.X = 0
		return
	}

	if len(e.PatrolPoints) > 0 {
		// Move toward current patrol point
		target := e.PatrolPoints[e.CurrentPoint]
		direction := target.Subtract(e.Position).Normalize()
		e.PhysicsBody.Velocity.X = direction.X * e.MoveSpeed

		// Update facing
		if direction.X < 0 {
			e.Facing = DirectionLeft
		} else if direction.X > 0 {
			e.Facing = DirectionRight
		}

		// Check if reached point
		if e.Position.Distance(target) < 10 {
			e.CurrentPoint = (e.CurrentPoint + 1) % len(e.PatrolPoints)
		}
	} else {
		// Simple back-and-forth patrol
		e.PhysicsBody.Velocity.X = e.MoveSpeed * float64(e.Facing)

		// Check if should turn
		if e.shouldTurn() {
			e.Facing = -e.Facing
		}
	}
}

// chaseTarget implements chase behavior
func (e *EnemyBase) chaseTarget(deltaTime float64) {
	if e.Target == nil {
		return
	}

	// Move toward target
	direction := e.Target.GetPosition().Subtract(e.Position).Normalize()

	if e.IsFlying {
		// Flying enemies can move in both X and Y
		e.PhysicsBody.Velocity.X = direction.X * e.MoveSpeed * 1.5
		e.PhysicsBody.Velocity.Y = direction.Y * e.MoveSpeed * 1.5
	} else {
		// Ground enemies only move horizontally
		e.PhysicsBody.Velocity.X = direction.X * e.MoveSpeed * 1.5
	}

	// Update facing
	if direction.X < 0 {
		e.Facing = DirectionLeft
	} else if direction.X > 0 {
		e.Facing = DirectionRight
	}
}

// attack implements attack behavior (base version)
func (e *EnemyBase) attack(deltaTime float64) {
	// Default attack: just move toward player
	e.chaseTarget(deltaTime)
}

// flee implements flee behavior
func (e *EnemyBase) flee(deltaTime float64) {
	if e.Target == nil {
		return
	}

	// Move away from target
	direction := e.Position.Subtract(e.Target.GetPosition()).Normalize()
	e.PhysicsBody.Velocity.X = direction.X * e.MoveSpeed * 1.2

	// Update facing
	if direction.X < 0 {
		e.Facing = DirectionLeft
	} else if direction.X > 0 {
		e.Facing = DirectionRight
	}
}

// shouldTurn checks if enemy should turn around
func (e *EnemyBase) shouldTurn() bool {
	// Check patrol bounds
	if e.Facing == DirectionLeft && e.Position.X <= e.PatrolLeft {
		return true
	}
	if e.Facing == DirectionRight && e.Position.X >= e.PatrolRight {
		return true
	}

	// Check for edge (if enabled)
	if e.TurnAtEdge && e.OnGround {
		// Raycast down to check for edge
		rayOrigin := e.Position
		rayOrigin.X += float64(e.Facing) * 16 // Check ahead
		rayDirection := engine.Vector2{X: 0, Y: 1}
		hit := e.physics.Raycast(rayOrigin, rayDirection, 20, engine.LayerTerrain)
		if hit == nil {
			// No ground ahead, turn around
			return true
		}
	}

	// Check for wall (if enabled)
	if e.TurnAtWall {
		// Raycast forward to check for wall
		rayOrigin := e.Position
		rayDirection := engine.Vector2{X: float64(e.Facing), Y: 0}
		hit := e.physics.Raycast(rayOrigin, rayDirection, 20, engine.LayerTerrain)
		if hit != nil {
			// Wall ahead, turn around
			return true
		}
	}

	return false
}

// updateState updates enemy state
func (e *EnemyBase) updateState(deltaTime float64) {
	switch e.State {
	case EnemyStateIdle:
		e.PhysicsBody.Velocity.X = 0

	case EnemyStateDead:
		// Death animation/cleanup handled elsewhere
		e.PhysicsBody.Velocity.X = 0

	case EnemyStateStunned:
		e.PhysicsBody.Velocity.X = 0
		if e.StateTimer <= 0 {
			e.State = EnemyStatePatrol
		}
	}
}

// updateAnimation updates animation frame
func (e *EnemyBase) updateAnimation(deltaTime float64) {
	e.AnimTime += deltaTime
	if e.AnimTime >= 0.2 { // 5 FPS animation
		e.AnimFrame = (e.AnimFrame + 1) % 4
		e.AnimTime = 0
	}
}

// findNearestPlayer finds the nearest alive player
func (e *EnemyBase) findNearestPlayer(players []Player) Player {
	var nearest Player
	minDist := math.MaxFloat64

	for _, player := range players {
		if player.IsAlive() {
			dist := e.Position.Distance(player.GetPosition())
			if dist < minDist {
				minDist = dist
				nearest = player
			}
		}
	}

	return nearest
}

// TakeDamage applies damage to the enemy
func (e *EnemyBase) TakeDamage(damage int) {
	if e.InvulnTimer > 0 || e.IsDead {
		return
	}

	e.Health -= damage
	e.State = EnemyStateHurt
	e.HurtTimer = 0.3 // 0.3 seconds hurt state
	e.InvulnTimer = 0.5 // 0.5 seconds of invulnerability

	// Knockback
	if e.PhysicsBody != nil {
		e.PhysicsBody.Velocity.Y = -150
	}

	if e.Health <= 0 {
		e.Die()
	}
}

// Die handles enemy death
func (e *EnemyBase) Die() {
	e.IsDead = true
	e.State = EnemyStateDead
	
	if e.PhysicsBody != nil {
		e.PhysicsBody.Velocity = engine.Vector2{0, 0}
		e.PhysicsBody.Enabled = false
	}

	// Award points (handled by game manager)
	// Spawn coin/power-up (10% chance, handled by game manager)
}

// OnCollision handles collision with other entities
func (e *EnemyBase) OnCollision(other interface{}) {
	// This will be called by the collision system
	// Specific enemy types can override this
}

// GetPosition returns the enemy's position (implements Entity interface)
func (e *EnemyBase) GetPosition() engine.Vector2 {
	return e.Position
}

// SetPosition sets the enemy's position (implements Entity interface)
func (e *EnemyBase) SetPosition(pos engine.Vector2) {
	e.Position = pos
	if e.PhysicsBody != nil {
		e.PhysicsBody.Position = pos
	}
}

// GetCollider returns the enemy's collider
func (e *EnemyBase) GetCollider() engine.Collider {
	if e.PhysicsBody != nil {
		return e.PhysicsBody.Collider
	}
	return nil
}

// IsActive returns whether the enemy is active (implements Entity interface)
func (e *EnemyBase) IsActive() bool {
	return e.Active && !e.IsDead
}

// IsAlive returns whether the enemy is alive
func (e *EnemyBase) IsAlive() bool {
	return !e.IsDead && e.Active
}

// GetType returns the enemy type
func (e *EnemyBase) GetType() EnemyType {
	return e.Type
}

// Deactivate deactivates the enemy
func (e *EnemyBase) Deactivate() {
	e.Active = false
	if e.PhysicsBody != nil {
		e.PhysicsBody.Enabled = false
	}
}

// Activate activates the enemy
func (e *EnemyBase) Activate() {
	e.Active = true
	if e.PhysicsBody != nil {
		e.PhysicsBody.Enabled = true
	}
}

// SetPatrolPoints sets custom patrol points
func (e *EnemyBase) SetPatrolPoints(points []engine.Vector2) {
	e.PatrolPoints = points
	e.CurrentPoint = 0
}

// SetPatrolRange sets the patrol range for simple back-and-forth patrol
func (e *EnemyBase) SetPatrolRange(left, right float64) {
	e.PatrolLeft = left
	e.PatrolRight = right
}

// GetPhysicsBody returns the physics body
func (e *EnemyBase) GetPhysicsBody() *engine.PhysicsBody {
	return e.PhysicsBody
}

// GetSprite returns the current sprite character
func (e *EnemyBase) GetSprite() rune {
	return e.Sprite
}

// GetSpriteStyle returns the sprite style
func (e *EnemyBase) GetSpriteStyle() tcell.Style {
	return e.SpriteStyle
}