package engine

import (
	"sync"
)

// Entity interface that physics bodies can reference
type Entity interface {
	SetPosition(pos Vector2)
	GetPosition() Vector2
	IsActive() bool
}

// PhysicsEngine manages all physics simulation including gravity, collisions, and movement
type PhysicsEngine struct {
	gravity          float64
	terminalVelocity float64
	bodies           []*PhysicsBody
	staticColliders  []Collider
	spatialGrid      *SpatialGrid
	collisions       []*Collision
	mutex            sync.RWMutex
}

// PhysicsBody represents a dynamic physics object
type PhysicsBody struct {
	Entity       Entity
	Position     Vector2
	Velocity     Vector2
	Acceleration Vector2
	Mass         float64
	Friction     float64
	Restitution  float64 // Bounciness (0 = no bounce, 1 = perfect bounce)
	Grounded     bool
	Collider     Collider
	Layer        CollisionLayer
	Enabled      bool
	GravityScale float64 // Multiplier for gravity (0 = no gravity, 1 = normal)
}

// NewPhysicsEngine creates a new physics engine
func NewPhysicsEngine() *PhysicsEngine {
	return &PhysicsEngine{
		gravity:          980.0, // pixels per second squared (similar to real gravity)
		terminalVelocity: 500.0, // maximum falling speed
		bodies:           make([]*PhysicsBody, 0),
		staticColliders:  make([]Collider, 0),
		spatialGrid:      NewSpatialGrid(64), // 64x64 pixel cells
		collisions:       make([]*Collision, 0),
	}
}

// Initialize sets up the physics engine
func (pe *PhysicsEngine) Initialize() {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()

	pe.bodies = make([]*PhysicsBody, 0)
	pe.staticColliders = make([]Collider, 0)
	pe.spatialGrid.Clear()
	pe.collisions = make([]*Collision, 0)
}

// Update performs one physics simulation step
func (pe *PhysicsEngine) Update(deltaTime float64) {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()

	// Clamp deltaTime to prevent physics instability
	if deltaTime > 0.1 {
		deltaTime = 0.1
	}

	// Update all physics bodies
	for _, body := range pe.bodies {
		if body.Enabled && body.Entity.IsActive() {
			pe.updateBody(body, deltaTime)
		}
	}

	// Rebuild spatial grid for dynamic objects
	pe.rebuildSpatialGrid()

	// Detect collisions
	pe.detectCollisions()

	// Resolve collisions
	pe.resolveCollisions()
}

// updateBody updates a single physics body
func (pe *PhysicsEngine) updateBody(body *PhysicsBody, deltaTime float64) {
	// Apply gravity
	if !body.Grounded && body.GravityScale > 0 {
		body.Acceleration.Y += pe.gravity * body.GravityScale * deltaTime
	}

	// Update velocity from acceleration
	body.Velocity.X += body.Acceleration.X * deltaTime
	body.Velocity.Y += body.Acceleration.Y * deltaTime

	// Apply friction when grounded
	if body.Grounded && body.Friction > 0 {
		body.Velocity.X *= (1.0 - body.Friction * deltaTime)
		// Stop very small velocities to prevent sliding
		if Abs(body.Velocity.X) < 1.0 {
			body.Velocity.X = 0
		}
	}

	// Clamp to terminal velocity
	if body.Velocity.Y > pe.terminalVelocity {
		body.Velocity.Y = pe.terminalVelocity
	}

	// Update position
	body.Position.X += body.Velocity.X * deltaTime
	body.Position.Y += body.Velocity.Y * deltaTime

	// Update collider position
	if aabb, ok := body.Collider.(*AABBCollider); ok {
		aabb.SetPosition(body.Position.X, body.Position.Y)
	}

	// Reset acceleration for next frame
	body.Acceleration = Vector2{0, 0}

	// Reset grounded state (will be set by collision detection)
	body.Grounded = false

	// Update entity position
	body.Entity.SetPosition(body.Position)
}

// rebuildSpatialGrid rebuilds the spatial grid with current body positions
func (pe *PhysicsEngine) rebuildSpatialGrid() {
	pe.spatialGrid.Clear()

	// Insert static colliders
	for _, collider := range pe.staticColliders {
		pe.spatialGrid.Insert(collider)
	}

	// Insert dynamic body colliders
	for _, body := range pe.bodies {
		if body.Enabled && body.Entity.IsActive() {
			pe.spatialGrid.Insert(body.Collider)
		}
	}
}

// detectCollisions finds all collisions between bodies and colliders
func (pe *PhysicsEngine) detectCollisions() {
	pe.collisions = make([]*Collision, 0)

	for _, body := range pe.bodies {
		if !body.Enabled || !body.Entity.IsActive() {
			continue
		}

		// Query spatial grid for nearby colliders
		nearbyColliders := pe.spatialGrid.Query(body.Collider.GetBounds())

		for _, collider := range nearbyColliders {
			// Skip self-collision
			if collider == body.Collider {
				continue
			}

			// Check if layers should collide
			if !ShouldCheckCollision(body.Layer, collider.GetLayer()) {
				continue
			}

			// Check for intersection
			if body.Collider.Intersects(collider) {
				// Handle one-way platforms specially
				if collider.GetType() == ColliderOneWay {
					if oneWay, ok := collider.(*OneWayCollider); ok {
						bodyBounds := body.Collider.GetBounds()
						bodyBottom := bodyBounds.Y + bodyBounds.Height
						if !oneWay.ShouldCollideFromAbove(bodyBottom, body.Velocity.Y) {
							continue
						}
					}
				}

				collision := NewCollision(body.Collider, collider)
				pe.collisions = append(pe.collisions, collision)
			}
		}
	}
}

// resolveCollisions resolves all detected collisions
func (pe *PhysicsEngine) resolveCollisions() {
	for _, collision := range pe.collisions {
		// Find the body that owns this collider
		var body *PhysicsBody
		for _, b := range pe.bodies {
			if b.Collider == collision.ColliderA {
				body = b
				break
			}
		}

		if body == nil || !body.Enabled {
			continue
		}

		normal := collision.GetNormal()
		penetration := collision.GetPenetration()

		// Separate the body from the collision
		separation := normal.Multiply(penetration)
		body.Position = body.Position.Add(separation)

		// Update collider position
		if aabb, ok := body.Collider.(*AABBCollider); ok {
			aabb.SetPosition(body.Position.X, body.Position.Y)
		}

		// Calculate relative velocity along the normal
		velocityAlongNormal := body.Velocity.Dot(normal)

		// Only resolve if moving into the collision
		if velocityAlongNormal < 0 {
			// Apply restitution (bounce)
			impulse := -(1 + body.Restitution) * velocityAlongNormal
			impulseVector := normal.Multiply(impulse)
			body.Velocity = body.Velocity.Add(impulseVector)
		}

		// Check if grounded (collision from below)
		if normal.Y < -0.5 {
			body.Grounded = true
			body.Velocity.Y = Max(0, body.Velocity.Y) // Prevent downward velocity when grounded
		}

		// Stop horizontal velocity on wall collision
		if Abs(normal.X) > 0.5 {
			body.Velocity.X = 0
		}
	}
}

// AddBody adds a physics body to the simulation
func (pe *PhysicsEngine) AddBody(body *PhysicsBody) {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()

	body.Enabled = true
	if body.GravityScale == 0 {
		body.GravityScale = 1.0 // Default gravity scale
	}
	pe.bodies = append(pe.bodies, body)
}

// RemoveBody removes a physics body from the simulation
func (pe *PhysicsEngine) RemoveBody(body *PhysicsBody) {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()

	for i, b := range pe.bodies {
		if b == body {
			pe.bodies = append(pe.bodies[:i], pe.bodies[i+1:]...)
			break
		}
	}
}

// AddStaticCollider adds a static collider to the simulation
func (pe *PhysicsEngine) AddStaticCollider(collider Collider) {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()

	pe.staticColliders = append(pe.staticColliders, collider)
	pe.spatialGrid.Insert(collider)
}

// RemoveStaticCollider removes a static collider from the simulation
func (pe *PhysicsEngine) RemoveStaticCollider(collider Collider) {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()

	for i, c := range pe.staticColliders {
		if c == collider {
			pe.staticColliders = append(pe.staticColliders[:i], pe.staticColliders[i+1:]...)
			pe.spatialGrid.Remove(collider)
			break
		}
	}
}

// Raycast performs a raycast and returns the first hit
func (pe *PhysicsEngine) Raycast(origin, direction Vector2, maxDistance float64, layerMask CollisionLayer) *RaycastHit {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()

	// Collect all colliders
	allColliders := make([]Collider, 0, len(pe.staticColliders)+len(pe.bodies))
	allColliders = append(allColliders, pe.staticColliders...)
	
	for _, body := range pe.bodies {
		if body.Enabled && body.Entity.IsActive() {
			allColliders = append(allColliders, body.Collider)
		}
	}

	return Raycast(origin, direction, maxDistance, allColliders, layerMask)
}

// GetBodiesInRadius returns all physics bodies within a radius of a point
func (pe *PhysicsEngine) GetBodiesInRadius(center Vector2, radius float64) []*PhysicsBody {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()

	colliders := pe.spatialGrid.QueryRadius(center, radius)
	bodies := make([]*PhysicsBody, 0)

	for _, collider := range colliders {
		for _, body := range pe.bodies {
			if body.Collider == collider && body.Enabled && body.Entity.IsActive() {
				bodies = append(bodies, body)
				break
			}
		}
	}

	return bodies
}

// SetGravity sets the gravity value
func (pe *PhysicsEngine) SetGravity(gravity float64) {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()
	pe.gravity = gravity
}

// GetGravity returns the current gravity value
func (pe *PhysicsEngine) GetGravity() float64 {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	return pe.gravity
}

// SetTerminalVelocity sets the terminal velocity
func (pe *PhysicsEngine) SetTerminalVelocity(velocity float64) {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()
	pe.terminalVelocity = velocity
}

// GetTerminalVelocity returns the current terminal velocity
func (pe *PhysicsEngine) GetTerminalVelocity() float64 {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	return pe.terminalVelocity
}

// GetBodyCount returns the number of active physics bodies
func (pe *PhysicsEngine) GetBodyCount() int {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	return len(pe.bodies)
}

// GetStaticColliderCount returns the number of static colliders
func (pe *PhysicsEngine) GetStaticColliderCount() int {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	return len(pe.staticColliders)
}

// GetCollisionCount returns the number of collisions detected in the last frame
func (pe *PhysicsEngine) GetCollisionCount() int {
	pe.mutex.RLock()
	defer pe.mutex.RUnlock()
	return len(pe.collisions)
}

// ClearStaticColliders removes all static colliders
func (pe *PhysicsEngine) ClearStaticColliders() {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()

	pe.staticColliders = make([]Collider, 0)
	pe.spatialGrid.Clear()
}

// Clear removes all bodies and colliders
func (pe *PhysicsEngine) Clear() {
	pe.mutex.Lock()
	defer pe.mutex.Unlock()

	pe.bodies = make([]*PhysicsBody, 0)
	pe.staticColliders = make([]Collider, 0)
	pe.spatialGrid.Clear()
	pe.collisions = make([]*Collision, 0)
}
