package projectiles

import (
	"testing"
)

// MockEntity implements the Entity interface for testing
type MockEntity struct {
	id     int
	x, y   float64
	w, h   float64
}

func (m *MockEntity) GetID() int {
	return m.id
}

func (m *MockEntity) GetPosition() Vector2 {
	return Vector2{X: m.x, Y: m.y}
}

func (m *MockEntity) GetBounds() (float64, float64, float64, float64) {
	return m.x, m.y, m.w, m.h
}

// MockDamageable implements both Entity and Damageable interfaces
type MockDamageable struct {
	MockEntity
	damageTaken int
	damageSource Entity
}

func (m *MockDamageable) TakeDamage(damage int, source Entity) {
	m.damageTaken += damage
	m.damageSource = source
}

func TestVector2_Add(t *testing.T) {
	v1 := Vector2{X: 1, Y: 2}
	v2 := Vector2{X: 3, Y: 4}
	result := v1.Add(v2)

	if result.X != 4 || result.Y != 6 {
		t.Errorf("Expected (4, 6), got (%f, %f)", result.X, result.Y)
	}
}

func TestVector2_Multiply(t *testing.T) {
	v := Vector2{X: 2, Y: 3}
	result := v.Multiply(2.5)

	if result.X != 5 || result.Y != 7.5 {
		t.Errorf("Expected (5, 7.5), got (%f, %f)", result.X, result.Y)
	}
}

func TestVector2_Length(t *testing.T) {
	v := Vector2{X: 3, Y: 4}
	length := v.Length()

	if length != 5 {
		t.Errorf("Expected length 5, got %f", length)
	}
}

func TestVector2_Normalize(t *testing.T) {
	v := Vector2{X: 3, Y: 4}
	normalized := v.Normalize()

	expectedX := 0.6
	expectedY := 0.8

	if normalized.X != expectedX || normalized.Y != expectedY {
		t.Errorf("Expected (%f, %f), got (%f, %f)", expectedX, expectedY, normalized.X, normalized.Y)
	}
}

func TestVector2_Normalize_Zero(t *testing.T) {
	v := Vector2{X: 0, Y: 0}
	normalized := v.Normalize()

	if normalized.X != 0 || normalized.Y != 0 {
		t.Errorf("Expected (0, 0) for zero vector, got (%f, %f)", normalized.X, normalized.Y)
	}
}

func TestNewProjectileManager(t *testing.T) {
	pm := NewProjectileManager(50)

	if pm == nil {
		t.Fatal("NewProjectileManager returned nil")
	}

	if pm.maxActive != 50 {
		t.Errorf("Expected maxActive 50, got %d", pm.maxActive)
	}

	if pm.GetProjectileCount() != 0 {
		t.Errorf("Expected 0 projectiles, got %d", pm.GetProjectileCount())
	}

	if pm.gravity != 980.0 {
		t.Errorf("Expected gravity 980.0, got %f", pm.gravity)
	}
}

func TestSpawnProjectile_Fireball(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	pos := Vector2{X: 10, Y: 20}
	direction := Vector2{X: 1, Y: 0}

	proj := pm.SpawnProjectile(owner, pos, direction, ProjectileFireball)

	if proj == nil {
		t.Fatal("SpawnProjectile returned nil")
	}

	if proj.Owner != owner {
		t.Error("Owner not set correctly")
	}

	if proj.Position.X != 10 || proj.Position.Y != 20 {
		t.Errorf("Position not set correctly: (%f, %f)", proj.Position.X, proj.Position.Y)
	}

	if proj.Damage != 1 {
		t.Errorf("Expected damage 1, got %d", proj.Damage)
	}

	if proj.MaxLifetime != 3.0 {
		t.Errorf("Expected max lifetime 3.0, got %f", proj.MaxLifetime)
	}

	if proj.MaxBounces != 1 {
		t.Errorf("Expected max bounces 1, got %d", proj.MaxBounces)
	}

	if proj.AffectedByGravity {
		t.Error("Fireball should not be affected by gravity")
	}

	if !proj.Active {
		t.Error("Projectile should be active")
	}

	if pm.GetProjectileCount() != 1 {
		t.Errorf("Expected 1 projectile, got %d", pm.GetProjectileCount())
	}
}

func TestSpawnProjectile_Hammer(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	pos := Vector2{X: 10, Y: 20}
	direction := Vector2{X: 1, Y: -0.5}

	proj := pm.SpawnProjectile(owner, pos, direction, ProjectileHammer)

	if proj == nil {
		t.Fatal("SpawnProjectile returned nil")
	}

	if proj.MaxBounces != 0 {
		t.Errorf("Expected max bounces 0, got %d", proj.MaxBounces)
	}

	if !proj.AffectedByGravity {
		t.Error("Hammer should be affected by gravity")
	}

	if proj.MaxLifetime != 2.0 {
		t.Errorf("Expected max lifetime 2.0, got %f", proj.MaxLifetime)
	}
}

func TestSpawnProjectile_MaxActive(t *testing.T) {
	pm := NewProjectileManager(2)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	pos := Vector2{X: 10, Y: 20}
	direction := Vector2{X: 1, Y: 0}

	// Spawn 2 projectiles (max)
	proj1 := pm.SpawnProjectile(owner, pos, direction, ProjectileFireball)
	proj2 := pm.SpawnProjectile(owner, pos, direction, ProjectileFireball)

	if proj1 == nil || proj2 == nil {
		t.Fatal("Failed to spawn projectiles")
	}

	// Try to spawn a third (should fail)
	proj3 := pm.SpawnProjectile(owner, pos, direction, ProjectileFireball)

	if proj3 != nil {
		t.Error("Should not be able to spawn more than maxActive projectiles")
	}

	if pm.GetProjectileCount() != 2 {
		t.Errorf("Expected 2 projectiles, got %d", pm.GetProjectileCount())
	}
}

func TestUpdate_Movement(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	pos := Vector2{X: 10, Y: 20}
	direction := Vector2{X: 1, Y: 0}

	proj := pm.SpawnProjectile(owner, pos, direction, ProjectileFireball)
	initialX := proj.Position.X

	pm.Update(0.1) // 0.1 seconds

	// Fireball moves at 400 pixels/sec, so in 0.1s it should move 40 pixels
	expectedX := initialX + 40
	if proj.Position.X < expectedX-1 || proj.Position.X > expectedX+1 {
		t.Errorf("Expected X ~%f, got %f", expectedX, proj.Position.X)
	}
}

func TestUpdate_Gravity(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	pos := Vector2{X: 10, Y: 20}
	direction := Vector2{X: 1, Y: 0}

	proj := pm.SpawnProjectile(owner, pos, direction, ProjectileHammer)
	initialVelocityY := proj.Velocity.Y

	pm.Update(0.1) // 0.1 seconds

	// Gravity should increase Y velocity
	if proj.Velocity.Y <= initialVelocityY {
		t.Error("Gravity should have increased Y velocity")
	}
}

func TestUpdate_Lifetime(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	pos := Vector2{X: 10, Y: 20}
	direction := Vector2{X: 1, Y: 0}

	proj := pm.SpawnProjectile(owner, pos, direction, ProjectileFireball)

	if proj.Lifetime != 0 {
		t.Error("Initial lifetime should be 0")
	}

	pm.Update(0.5)

	if proj.Lifetime != 0.5 {
		t.Errorf("Expected lifetime 0.5, got %f", proj.Lifetime)
	}
}

func TestUpdate_LifetimeExpiry(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	pos := Vector2{X: 10, Y: 20}
	direction := Vector2{X: 1, Y: 0}

	pm.SpawnProjectile(owner, pos, direction, ProjectileFireball) // 3s lifetime

	if pm.GetProjectileCount() != 1 {
		t.Fatal("Expected 1 projectile")
	}

	// Update for longer than lifetime
	pm.Update(3.5)

	if pm.GetProjectileCount() != 0 {
		t.Errorf("Expected 0 projectiles after expiry, got %d", pm.GetProjectileCount())
	}
}

func TestCheckTerrainCollision_NoBounce(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	pos := Vector2{X: 10, Y: 93}
	direction := Vector2{X: 1, Y: 1}

	proj := pm.SpawnProjectile(owner, pos, direction, ProjectileHammer) // No bounces
	terrainY := 100.0

	hit := pm.CheckTerrainCollision(proj, terrainY)

	if !hit {
		t.Error("Expected terrain collision")
	}

	if proj.Active {
		t.Error("Projectile should be despawned after terrain hit with no bounces")
	}
}

func TestCheckTerrainCollision_WithBounce(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	pos := Vector2{X: 10, Y: 93}
	direction := Vector2{X: 1, Y: 1}

	proj := pm.SpawnProjectile(owner, pos, direction, ProjectileFireball) // 1 bounce
	initialVelocityY := proj.Velocity.Y
	terrainY := 100.0

	hit := pm.CheckTerrainCollision(proj, terrainY)

	if hit {
		t.Error("Should not despawn on first bounce")
	}

	if proj.Bounces != 1 {
		t.Errorf("Expected 1 bounce, got %d", proj.Bounces)
	}

	if proj.Velocity.Y >= 0 {
		t.Error("Velocity Y should be negative after bounce")
	}

	if proj.Velocity.Y >= initialVelocityY {
		t.Error("Velocity should be reduced after bounce")
	}
}

func TestCheckEntityCollision(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	target := &MockEntity{id: 2, x: 50, y: 50, w: 16, h: 16}
	pos := Vector2{X: 52, Y: 52}
	direction := Vector2{X: 1, Y: 0}

	proj := pm.SpawnProjectile(owner, pos, direction, ProjectileFireball)

	// Should collide (projectile is inside target bounds)
	if !pm.CheckEntityCollision(proj, target) {
		t.Error("Expected collision with target")
	}

	// Should not collide with owner
	if pm.CheckEntityCollision(proj, owner) {
		t.Error("Should not collide with owner")
	}
}

func TestCheckEntityCollision_NoOverlap(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	target := &MockEntity{id: 2, x: 100, y: 100, w: 16, h: 16}
	pos := Vector2{X: 10, Y: 20}
	direction := Vector2{X: 1, Y: 0}

	proj := pm.SpawnProjectile(owner, pos, direction, ProjectileFireball)

	if pm.CheckEntityCollision(proj, target) {
		t.Error("Should not collide when far apart")
	}
}

func TestOnProjectileHit(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	target := &MockDamageable{
		MockEntity: MockEntity{id: 2, x: 50, y: 50, w: 16, h: 16},
	}
	pos := Vector2{X: 10, Y: 20}
	direction := Vector2{X: 1, Y: 0}

	proj := pm.SpawnProjectile(owner, pos, direction, ProjectileFireball)

	pm.OnProjectileHit(proj, target)

	if target.damageTaken != 1 {
		t.Errorf("Expected damage 1, got %d", target.damageTaken)
	}

	if target.damageSource != owner {
		t.Error("Damage source should be projectile owner")
	}

	if proj.Active {
		t.Error("Projectile should be despawned after hit")
	}

	if !proj.HitEntities[target.GetID()] {
		t.Error("Target should be marked as hit")
	}
}

func TestClear(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	pos := Vector2{X: 10, Y: 20}
	direction := Vector2{X: 1, Y: 0}

	pm.SpawnProjectile(owner, pos, direction, ProjectileFireball)
	pm.SpawnProjectile(owner, pos, direction, ProjectileHammer)

	if pm.GetProjectileCount() != 2 {
		t.Fatal("Expected 2 projectiles")
	}

	pm.Clear()

	if pm.GetProjectileCount() != 0 {
		t.Errorf("Expected 0 projectiles after clear, got %d", pm.GetProjectileCount())
	}
}

func TestSetGetGravity(t *testing.T) {
	pm := NewProjectileManager(50)

	if pm.GetGravity() != 980.0 {
		t.Errorf("Expected default gravity 980.0, got %f", pm.GetGravity())
	}

	pm.SetGravity(500.0)

	if pm.GetGravity() != 500.0 {
		t.Errorf("Expected gravity 500.0, got %f", pm.GetGravity())
	}
}

func TestGetProjectiles(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	pos := Vector2{X: 10, Y: 20}
	direction := Vector2{X: 1, Y: 0}

	pm.SpawnProjectile(owner, pos, direction, ProjectileFireball)
	pm.SpawnProjectile(owner, pos, direction, ProjectileHammer)

	projectiles := pm.GetProjectiles()

	if len(projectiles) != 2 {
		t.Errorf("Expected 2 projectiles, got %d", len(projectiles))
	}
}

func TestMultipleProjectileTypes(t *testing.T) {
	pm := NewProjectileManager(50)
	owner := &MockEntity{id: 1, x: 10, y: 20, w: 16, h: 16}
	pos := Vector2{X: 10, Y: 20}
	direction := Vector2{X: 1, Y: 0}

	types := []ProjectileType{
		ProjectileFireball,
		ProjectileHammer,
		ProjectileIceball,
		ProjectileBullet,
	}

	for _, pType := range types {
		proj := pm.SpawnProjectile(owner, pos, direction, pType)
		if proj == nil {
			t.Errorf("Failed to spawn projectile of type %v", pType)
		}
	}

	if pm.GetProjectileCount() != len(types) {
		t.Errorf("Expected %d projectiles, got %d", len(types), pm.GetProjectileCount())
	}
}
