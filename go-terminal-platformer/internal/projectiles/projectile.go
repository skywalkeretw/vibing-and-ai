package projectiles

import (
	"math"
)

// ProjectileType defines different types of projectiles
type ProjectileType int

const (
	ProjectileFireball ProjectileType = iota
	ProjectileHammer
	ProjectileIceball
	ProjectileBullet
)

// Vector2 represents a 2D vector
type Vector2 struct {
	X float64
	Y float64
}

// Add adds two vectors
func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Multiply multiplies a vector by a scalar
func (v Vector2) Multiply(scalar float64) Vector2 {
	return Vector2{X: v.X * scalar, Y: v.Y * scalar}
}

// Length returns the length of the vector
func (v Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize returns a normalized version of the vector
func (v Vector2) Normalize() Vector2 {
	length := v.Length()
	if length == 0 {
		return Vector2{X: 0, Y: 0}
	}
	return Vector2{X: v.X / length, Y: v.Y / length}
}

// Entity interface for projectile owners and targets
type Entity interface {
	GetID() int
	GetPosition() Vector2
	GetBounds() (float64, float64, float64, float64) // x, y, width, height
}

// Damageable interface for entities that can take damage
type Damageable interface {
	TakeDamage(damage int, source Entity)
}

// Projectile represents a single projectile in the game
type Projectile struct {
	ID             int
	Owner          Entity
	Position       Vector2
	Velocity       Vector2
	ProjectileType ProjectileType
	Damage         int
	Lifetime       float64
	MaxLifetime    float64
	Bounces        int
	MaxBounces     int
	Active         bool
	HitEntities    map[int]bool
	Width          float64
	Height         float64
	AffectedByGravity bool
}

// ProjectileManager manages all projectiles in the game
type ProjectileManager struct {
	projectiles []*Projectile
	pool        *ProjectilePool
	maxActive   int
	nextID      int
	gravity     float64
}

// NewProjectileManager creates a new projectile manager
func NewProjectileManager(maxActive int) *ProjectileManager {
	return &ProjectileManager{
		projectiles: make([]*Projectile, 0, maxActive),
		pool:        NewProjectilePool(maxActive),
		maxActive:   maxActive,
		nextID:      1,
		gravity:     980.0, // pixels/sec^2
	}
}

// SpawnProjectile spawns a new projectile
func (pm *ProjectileManager) SpawnProjectile(owner Entity, pos Vector2, direction Vector2, projType ProjectileType) *Projectile {
	if len(pm.projectiles) >= pm.maxActive {
		return nil
	}

	proj := pm.pool.Get()
	proj.ID = pm.nextID
	pm.nextID++
	proj.Owner = owner
	proj.Position = pos
	proj.ProjectileType = projType
	proj.Active = true
	proj.Lifetime = 0
	proj.Bounces = 0

	// Normalize direction
	direction = direction.Normalize()

	// Set properties based on type
	switch projType {
	case ProjectileFireball:
		proj.Velocity = direction.Multiply(400.0) // 400 pixels/sec
		proj.Damage = 1
		proj.MaxLifetime = 3.0
		proj.MaxBounces = 1
		proj.Width = 8
		proj.Height = 8
		proj.AffectedByGravity = false

	case ProjectileHammer:
		proj.Velocity = direction.Multiply(300.0) // 300 pixels/sec
		proj.Damage = 1
		proj.MaxLifetime = 2.0
		proj.MaxBounces = 0
		proj.Width = 8
		proj.Height = 8
		proj.AffectedByGravity = true

	case ProjectileIceball:
		proj.Velocity = direction.Multiply(350.0) // 350 pixels/sec
		proj.Damage = 1
		proj.MaxLifetime = 3.0
		proj.MaxBounces = 1
		proj.Width = 8
		proj.Height = 8
		proj.AffectedByGravity = false

	case ProjectileBullet:
		proj.Velocity = direction.Multiply(500.0) // 500 pixels/sec
		proj.Damage = 1
		proj.MaxLifetime = 2.0
		proj.MaxBounces = 0
		proj.Width = 4
		proj.Height = 4
		proj.AffectedByGravity = false
	}

	pm.projectiles = append(pm.projectiles, proj)
	return proj
}

// Update updates all projectiles
func (pm *ProjectileManager) Update(deltaTime float64) {
	for i := len(pm.projectiles) - 1; i >= 0; i-- {
		proj := pm.projectiles[i]

		if !proj.Active {
			pm.removeProjectile(i)
			continue
		}

		// Update lifetime
		proj.Lifetime += deltaTime
		if proj.Lifetime >= proj.MaxLifetime {
			pm.despawnProjectile(proj)
			pm.removeProjectile(i)
			continue
		}

		// Apply gravity if affected
		if proj.AffectedByGravity {
			proj.Velocity.Y += pm.gravity * deltaTime
		}

		// Update position
		proj.Position.X += proj.Velocity.X * deltaTime
		proj.Position.Y += proj.Velocity.Y * deltaTime
	}
}

// CheckTerrainCollision checks if a projectile collides with terrain at given position
func (pm *ProjectileManager) CheckTerrainCollision(proj *Projectile, terrainY float64) bool {
	// Simple ground collision check
	if proj.Position.Y+proj.Height >= terrainY {
		if proj.Bounces < proj.MaxBounces {
			proj.Bounces++
			proj.Velocity.Y *= -0.7 // Bounce with energy loss
			proj.Position.Y = terrainY - proj.Height
			return false // Don't despawn, bounced
		}
		pm.despawnProjectile(proj)
		return true // Despawn
	}
	return false
}

// CheckEntityCollision checks if a projectile collides with an entity
func (pm *ProjectileManager) CheckEntityCollision(proj *Projectile, entity Entity) bool {
	if entity == proj.Owner {
		return false
	}

	if proj.HitEntities[entity.GetID()] {
		return false
	}

	// Get entity bounds
	ex, ey, ew, eh := entity.GetBounds()

	// AABB collision detection
	if proj.Position.X < ex+ew &&
		proj.Position.X+proj.Width > ex &&
		proj.Position.Y < ey+eh &&
		proj.Position.Y+proj.Height > ey {
		return true
	}

	return false
}

// OnProjectileHit handles projectile hitting an entity
func (pm *ProjectileManager) OnProjectileHit(proj *Projectile, target Entity) {
	// Mark entity as hit
	proj.HitEntities[target.GetID()] = true

	// Apply damage if target is damageable
	if damageable, ok := target.(Damageable); ok {
		damageable.TakeDamage(proj.Damage, proj.Owner)
	}

	// Despawn projectile after hit
	pm.despawnProjectile(proj)
}

// despawnProjectile marks a projectile as inactive
func (pm *ProjectileManager) despawnProjectile(proj *Projectile) {
	proj.Active = false
}

// removeProjectile removes a projectile from the active list and returns it to pool
func (pm *ProjectileManager) removeProjectile(index int) {
	proj := pm.projectiles[index]
	pm.pool.Return(proj)
	pm.projectiles = append(pm.projectiles[:index], pm.projectiles[index+1:]...)
}

// GetProjectiles returns all active projectiles
func (pm *ProjectileManager) GetProjectiles() []*Projectile {
	return pm.projectiles
}

// GetProjectileCount returns the number of active projectiles
func (pm *ProjectileManager) GetProjectileCount() int {
	return len(pm.projectiles)
}

// Clear removes all projectiles
func (pm *ProjectileManager) Clear() {
	for _, proj := range pm.projectiles {
		pm.pool.Return(proj)
	}
	pm.projectiles = pm.projectiles[:0]
}

// SetGravity sets the gravity value
func (pm *ProjectileManager) SetGravity(gravity float64) {
	pm.gravity = gravity
}

// GetGravity returns the current gravity value
func (pm *ProjectileManager) GetGravity() float64 {
	return pm.gravity
}
