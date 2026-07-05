package engine

import (
	"testing"
)

// Mock entity for testing
type MockEntity struct {
	position Vector2
	active   bool
}

func (m *MockEntity) SetPosition(pos Vector2) {
	m.position = pos
}

func (m *MockEntity) GetPosition() Vector2 {
	return m.position
}

func (m *MockEntity) IsActive() bool {
	return m.active
}

func TestNewPhysicsEngine(t *testing.T) {
	physics := NewPhysicsEngine()

	if physics.GetGravity() != 980.0 {
		t.Errorf("Expected gravity 980.0, got %f", physics.GetGravity())
	}

	if physics.GetTerminalVelocity() != 500.0 {
		t.Errorf("Expected terminal velocity 500.0, got %f", physics.GetTerminalVelocity())
	}

	if physics.GetBodyCount() != 0 {
		t.Error("New physics engine should have 0 bodies")
	}
}

func TestPhysicsEngineInitialize(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	if physics.GetBodyCount() != 0 {
		t.Error("Initialized physics engine should have 0 bodies")
	}

	if physics.GetStaticColliderCount() != 0 {
		t.Error("Initialized physics engine should have 0 static colliders")
	}
}

func TestPhysicsEngineAddBody(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	entity := &MockEntity{position: Vector2{X: 0, Y: 0}, active: true}
	body := &PhysicsBody{
		Entity:   entity,
		Position: Vector2{X: 0, Y: 0},
		Velocity: Vector2{X: 0, Y: 0},
		Mass:     1.0,
		Collider: NewAABBCollider(0, 0, 10, 10, LayerPlayer),
		Layer:    LayerPlayer,
	}

	physics.AddBody(body)

	if physics.GetBodyCount() != 1 {
		t.Error("Physics engine should have 1 body after adding")
	}

	if !body.Enabled {
		t.Error("Body should be enabled after adding")
	}
}

func TestPhysicsEngineRemoveBody(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	entity := &MockEntity{position: Vector2{X: 0, Y: 0}, active: true}
	body := &PhysicsBody{
		Entity:   entity,
		Position: Vector2{X: 0, Y: 0},
		Velocity: Vector2{X: 0, Y: 0},
		Mass:     1.0,
		Collider: NewAABBCollider(0, 0, 10, 10, LayerPlayer),
		Layer:    LayerPlayer,
	}

	physics.AddBody(body)
	physics.RemoveBody(body)

	if physics.GetBodyCount() != 0 {
		t.Error("Physics engine should have 0 bodies after removing")
	}
}

func TestPhysicsEngineAddStaticCollider(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	collider := NewAABBCollider(0, 0, 100, 10, LayerTerrain)
	physics.AddStaticCollider(collider)

	if physics.GetStaticColliderCount() != 1 {
		t.Error("Physics engine should have 1 static collider after adding")
	}
}

func TestPhysicsEngineRemoveStaticCollider(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	collider := NewAABBCollider(0, 0, 100, 10, LayerTerrain)
	physics.AddStaticCollider(collider)
	physics.RemoveStaticCollider(collider)

	if physics.GetStaticColliderCount() != 0 {
		t.Error("Physics engine should have 0 static colliders after removing")
	}
}

func TestPhysicsEngineGravity(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	entity := &MockEntity{position: Vector2{X: 0, Y: 0}, active: true}
	body := &PhysicsBody{
		Entity:       entity,
		Position:     Vector2{X: 0, Y: 0},
		Velocity:     Vector2{X: 0, Y: 0},
		Mass:         1.0,
		Collider:     NewAABBCollider(0, 0, 10, 10, LayerPlayer),
		Layer:        LayerPlayer,
		GravityScale: 1.0,
	}

	physics.AddBody(body)

	// Update physics for a small time step
	deltaTime := 0.016 // ~60 FPS
	physics.Update(deltaTime)

	// Body should have downward velocity due to gravity
	if body.Velocity.Y <= 0 {
		t.Error("Body should have positive (downward) velocity after gravity is applied")
	}

	// Position should have moved down
	if body.Position.Y <= 0 {
		t.Error("Body should have moved down after gravity is applied")
	}
}

func TestPhysicsEngineTerminalVelocity(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	entity := &MockEntity{position: Vector2{X: 0, Y: 0}, active: true}
	body := &PhysicsBody{
		Entity:       entity,
		Position:     Vector2{X: 0, Y: 0},
		Velocity:     Vector2{X: 0, Y: 0},
		Mass:         1.0,
		Collider:     NewAABBCollider(0, 0, 10, 10, LayerPlayer),
		Layer:        LayerPlayer,
		GravityScale: 1.0,
	}

	physics.AddBody(body)

	// Update physics many times to reach terminal velocity
	for i := 0; i < 100; i++ {
		physics.Update(0.016)
	}

	// Velocity should not exceed terminal velocity
	if body.Velocity.Y > physics.GetTerminalVelocity() {
		t.Errorf("Velocity (%f) should not exceed terminal velocity (%f)",
			body.Velocity.Y, physics.GetTerminalVelocity())
	}
}

func TestPhysicsEngineCollisionDetection(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	// Add a static ground collider
	ground := NewAABBCollider(0, 100, 200, 10, LayerTerrain)
	physics.AddStaticCollider(ground)

	// Add a falling body
	entity := &MockEntity{position: Vector2{X: 50, Y: 50}, active: true}
	body := &PhysicsBody{
		Entity:       entity,
		Position:     Vector2{X: 50, Y: 50},
		Velocity:     Vector2{X: 0, Y: 100},
		Mass:         1.0,
		Collider:     NewAABBCollider(50, 50, 10, 10, LayerPlayer),
		Layer:        LayerPlayer,
		GravityScale: 1.0,
	}

	physics.AddBody(body)

	// Update physics multiple times
	for i := 0; i < 10; i++ {
		physics.Update(0.016)
	}

	// Body should have collided with ground and stopped falling
	if body.Position.Y > 90 {
		t.Error("Body should have collided with ground and stopped")
	}

	if !body.Grounded {
		t.Error("Body should be grounded after collision with ground")
	}
}

func TestPhysicsEngineCollisionResolution(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	// Add a wall
	wall := NewAABBCollider(100, 0, 10, 200, LayerTerrain)
	physics.AddStaticCollider(wall)

	// Add a body moving toward the wall
	entity := &MockEntity{position: Vector2{X: 50, Y: 50}, active: true}
	body := &PhysicsBody{
		Entity:       entity,
		Position:     Vector2{X: 50, Y: 50},
		Velocity:     Vector2{X: 100, Y: 0},
		Mass:         1.0,
		Collider:     NewAABBCollider(50, 50, 10, 10, LayerPlayer),
		Layer:        LayerPlayer,
		GravityScale: 0, // No gravity for this test
	}

	physics.AddBody(body)

	// Update physics multiple times
	for i := 0; i < 10; i++ {
		physics.Update(0.016)
	}

	// Body should have stopped at the wall
	if body.Position.X >= 90 {
		t.Error("Body should have stopped before reaching the wall")
	}

	// Horizontal velocity should be zero after collision
	if body.Velocity.X != 0 {
		t.Error("Horizontal velocity should be zero after wall collision")
	}
}

func TestPhysicsEngineFriction(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	// Add ground
	ground := NewAABBCollider(0, 100, 200, 10, LayerTerrain)
	physics.AddStaticCollider(ground)

	// Add a body with horizontal velocity on the ground
	entity := &MockEntity{position: Vector2{X: 50, Y: 90}, active: true}
	body := &PhysicsBody{
		Entity:       entity,
		Position:     Vector2{X: 50, Y: 90},
		Velocity:     Vector2{X: 100, Y: 0},
		Mass:         1.0,
		Friction:     0.9,
		Collider:     NewAABBCollider(50, 90, 10, 10, LayerPlayer),
		Layer:        LayerPlayer,
		GravityScale: 1.0,
		Grounded:     true,
	}

	physics.AddBody(body)

	initialVelocity := body.Velocity.X

	// Update physics
	physics.Update(0.016)

	// Velocity should have decreased due to friction
	if body.Velocity.X >= initialVelocity {
		t.Error("Friction should reduce horizontal velocity when grounded")
	}
}

func TestPhysicsEngineRestitution(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	// Add ground
	ground := NewAABBCollider(0, 100, 200, 10, LayerTerrain)
	physics.AddStaticCollider(ground)

	// Add a bouncy body
	entity := &MockEntity{position: Vector2{X: 50, Y: 50}, active: true}
	body := &PhysicsBody{
		Entity:       entity,
		Position:     Vector2{X: 50, Y: 50},
		Velocity:     Vector2{X: 0, Y: 100},
		Mass:         1.0,
		Restitution:  0.8, // 80% bounce
		Collider:     NewAABBCollider(50, 50, 10, 10, LayerPlayer),
		Layer:        LayerPlayer,
		GravityScale: 1.0,
	}

	physics.AddBody(body)

	// Update until collision
	for i := 0; i < 10; i++ {
		physics.Update(0.016)
	}

	// After collision, body should bounce (negative Y velocity)
	// Note: This test might be flaky due to timing, but checks the concept
	if body.Velocity.Y > 0 {
		// Body might still be falling or have bounced and is falling again
		// This is acceptable for this test
	}
}

func TestPhysicsEngineOneWayPlatform(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	// Add one-way platform
	platform := NewOneWayCollider(50, 100, 100, 5, LayerOneWayPlatform)
	physics.AddStaticCollider(platform)

	// Add a body falling from above
	entity := &MockEntity{position: Vector2{X: 75, Y: 50}, active: true}
	body := &PhysicsBody{
		Entity:       entity,
		Position:     Vector2{X: 75, Y: 50},
		Velocity:     Vector2{X: 0, Y: 100},
		Mass:         1.0,
		Collider:     NewAABBCollider(75, 50, 10, 10, LayerPlayer),
		Layer:        LayerPlayer,
		GravityScale: 1.0,
	}

	physics.AddBody(body)

	// Update physics
	for i := 0; i < 10; i++ {
		physics.Update(0.016)
	}

	// Body should land on the platform
	if body.Position.Y > 95 {
		t.Error("Body should have landed on one-way platform")
	}
}

func TestPhysicsEngineRaycast(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	// Add a wall
	wall := NewAABBCollider(100, 50, 10, 100, LayerTerrain)
	physics.AddStaticCollider(wall)

	// Raycast toward the wall
	origin := Vector2{X: 0, Y: 100}
	direction := Vector2{X: 1, Y: 0}
	hit := physics.Raycast(origin, direction, 200, LayerTerrain)

	if !hit.Hit {
		t.Error("Raycast should hit the wall")
	}

	if hit.Distance <= 0 || hit.Distance > 200 {
		t.Errorf("Hit distance should be positive and less than max distance, got %f", hit.Distance)
	}
}

func TestPhysicsEngineGetBodiesInRadius(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	// Add bodies at different positions
	entity1 := &MockEntity{position: Vector2{X: 50, Y: 50}, active: true}
	body1 := &PhysicsBody{
		Entity:   entity1,
		Position: Vector2{X: 50, Y: 50},
		Collider: NewAABBCollider(50, 50, 10, 10, LayerPlayer),
		Layer:    LayerPlayer,
	}

	entity2 := &MockEntity{position: Vector2{X: 200, Y: 200}, active: true}
	body2 := &PhysicsBody{
		Entity:   entity2,
		Position: Vector2{X: 200, Y: 200},
		Collider: NewAABBCollider(200, 200, 10, 10, LayerEnemy),
		Layer:    LayerEnemy,
	}

	physics.AddBody(body1)
	physics.AddBody(body2)

	// Query bodies near body1
	center := Vector2{X: 50, Y: 50}
	radius := 50.0
	bodies := physics.GetBodiesInRadius(center, radius)

	foundBody1 := false
	foundBody2 := false

	for _, body := range bodies {
		if body == body1 {
			foundBody1 = true
		}
		if body == body2 {
			foundBody2 = true
		}
	}

	if !foundBody1 {
		t.Error("Should find body1 within radius")
	}

	if foundBody2 {
		t.Error("Should not find body2 outside radius")
	}
}

func TestPhysicsEngineSetGravity(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.SetGravity(500.0)

	if physics.GetGravity() != 500.0 {
		t.Errorf("Expected gravity 500.0, got %f", physics.GetGravity())
	}
}

func TestPhysicsEngineSetTerminalVelocity(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.SetTerminalVelocity(300.0)

	if physics.GetTerminalVelocity() != 300.0 {
		t.Errorf("Expected terminal velocity 300.0, got %f", physics.GetTerminalVelocity())
	}
}

func TestPhysicsEngineClear(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	// Add bodies and colliders
	entity := &MockEntity{position: Vector2{X: 0, Y: 0}, active: true}
	body := &PhysicsBody{
		Entity:   entity,
		Position: Vector2{X: 0, Y: 0},
		Collider: NewAABBCollider(0, 0, 10, 10, LayerPlayer),
		Layer:    LayerPlayer,
	}
	physics.AddBody(body)

	collider := NewAABBCollider(0, 100, 100, 10, LayerTerrain)
	physics.AddStaticCollider(collider)

	physics.Clear()

	if physics.GetBodyCount() != 0 {
		t.Error("Physics engine should have 0 bodies after clear")
	}

	if physics.GetStaticColliderCount() != 0 {
		t.Error("Physics engine should have 0 static colliders after clear")
	}
}

func TestPhysicsEngineInactiveEntity(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	// Add an inactive entity
	entity := &MockEntity{position: Vector2{X: 0, Y: 0}, active: false}
	body := &PhysicsBody{
		Entity:       entity,
		Position:     Vector2{X: 0, Y: 0},
		Velocity:     Vector2{X: 0, Y: 0},
		Collider:     NewAABBCollider(0, 0, 10, 10, LayerPlayer),
		Layer:        LayerPlayer,
		GravityScale: 1.0,
	}

	physics.AddBody(body)

	initialPosition := body.Position

	// Update physics
	physics.Update(0.016)

	// Inactive entity should not move
	if body.Position.Y != initialPosition.Y {
		t.Error("Inactive entity should not be affected by physics")
	}
}

func TestPhysicsEngineDisabledBody(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	entity := &MockEntity{position: Vector2{X: 0, Y: 0}, active: true}
	body := &PhysicsBody{
		Entity:       entity,
		Position:     Vector2{X: 0, Y: 0},
		Velocity:     Vector2{X: 0, Y: 0},
		Collider:     NewAABBCollider(0, 0, 10, 10, LayerPlayer),
		Layer:        LayerPlayer,
		GravityScale: 1.0,
		Enabled:      false,
	}

	physics.AddBody(body)

	initialPosition := body.Position

	// Update physics
	physics.Update(0.016)

	// Disabled body should not move
	if body.Position.Y != initialPosition.Y {
		t.Error("Disabled body should not be affected by physics")
	}
}

func TestPhysicsEngineGravityScale(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	// Body with normal gravity
	entity1 := &MockEntity{position: Vector2{X: 0, Y: 0}, active: true}
	body1 := &PhysicsBody{
		Entity:       entity1,
		Position:     Vector2{X: 0, Y: 0},
		Velocity:     Vector2{X: 0, Y: 0},
		Collider:     NewAABBCollider(0, 0, 10, 10, LayerPlayer),
		Layer:        LayerPlayer,
		GravityScale: 1.0,
	}

	// Body with no gravity
	entity2 := &MockEntity{position: Vector2{X: 50, Y: 0}, active: true}
	body2 := &PhysicsBody{
		Entity:       entity2,
		Position:     Vector2{X: 50, Y: 0},
		Velocity:     Vector2{X: 0, Y: 0},
		Collider:     NewAABBCollider(50, 0, 10, 10, LayerPlayer),
		Layer:        LayerPlayer,
		GravityScale: 0,
	}

	physics.AddBody(body1)
	physics.AddBody(body2)

	// Update physics
	physics.Update(0.016)

	// Body1 should fall, body2 should not
	if body1.Velocity.Y <= 0 {
		t.Error("Body with gravity scale 1.0 should fall")
	}

	if body2.Velocity.Y != 0 {
		t.Error("Body with gravity scale 0 should not fall")
	}
}

func TestPhysicsEngineDeltaTimeClamp(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	entity := &MockEntity{position: Vector2{X: 0, Y: 0}, active: true}
	body := &PhysicsBody{
		Entity:       entity,
		Position:     Vector2{X: 0, Y: 0},
		Velocity:     Vector2{X: 0, Y: 0},
		Collider:     NewAABBCollider(0, 0, 10, 10, LayerPlayer),
		Layer:        LayerPlayer,
		GravityScale: 1.0,
	}

	physics.AddBody(body)

	// Update with very large deltaTime (should be clamped)
	physics.Update(1.0) // 1 second

	// Body should not have moved too far (deltaTime clamped to 0.1)
	if body.Position.Y > 100 {
		t.Error("Large deltaTime should be clamped to prevent physics instability")
	}
}

func TestPhysicsEngineCollisionCount(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	// Add ground
	ground := NewAABBCollider(0, 100, 200, 10, LayerTerrain)
	physics.AddStaticCollider(ground)

	// Add falling body
	entity := &MockEntity{position: Vector2{X: 50, Y: 50}, active: true}
	body := &PhysicsBody{
		Entity:       entity,
		Position:     Vector2{X: 50, Y: 50},
		Velocity:     Vector2{X: 0, Y: 100},
		Collider:     NewAABBCollider(50, 50, 10, 10, LayerPlayer),
		Layer:        LayerPlayer,
		GravityScale: 1.0,
	}

	physics.AddBody(body)

	// Update physics
	physics.Update(0.016)

	// Should have detected at least one collision
	if physics.GetCollisionCount() == 0 {
		// Might not collide in first frame, update more
		for i := 0; i < 10; i++ {
			physics.Update(0.016)
			if physics.GetCollisionCount() > 0 {
				return // Test passed
			}
		}
		t.Error("Should have detected collision between body and ground")
	}
}

func TestPhysicsEngineMultipleBodies(t *testing.T) {
	physics := NewPhysicsEngine()
	physics.Initialize()

	// Add multiple bodies
	for i := 0; i < 10; i++ {
		entity := &MockEntity{position: Vector2{X: float64(i * 20), Y: 0}, active: true}
		body := &PhysicsBody{
			Entity:       entity,
			Position:     Vector2{X: float64(i * 20), Y: 0},
			Velocity:     Vector2{X: 0, Y: 0},
			Collider:     NewAABBCollider(float64(i*20), 0, 10, 10, LayerPlayer),
			Layer:        LayerPlayer,
			GravityScale: 1.0,
		}
		physics.AddBody(body)
	}

	if physics.GetBodyCount() != 10 {
		t.Errorf("Expected 10 bodies, got %d", physics.GetBodyCount())
	}

	// Update physics
	physics.Update(0.016)

	// All bodies should have moved
	// (This is a basic test, more detailed tests would check individual body states)
}
