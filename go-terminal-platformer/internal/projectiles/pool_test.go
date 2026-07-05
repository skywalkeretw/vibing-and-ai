package projectiles

import (
	"testing"
)

func TestNewProjectilePool(t *testing.T) {
	pool := NewProjectilePool(100)

	if pool == nil {
		t.Fatal("NewProjectilePool returned nil")
	}

	if pool.Size() != 100 {
		t.Errorf("Expected pool size 100, got %d", pool.Size())
	}
}

func TestProjectilePool_Get(t *testing.T) {
	pool := NewProjectilePool(10)
	initialSize := pool.Size()

	proj := pool.Get()

	if proj == nil {
		t.Fatal("Get returned nil projectile")
	}

	if pool.Size() != initialSize-1 {
		t.Errorf("Expected pool size %d after Get, got %d", initialSize-1, pool.Size())
	}

	if proj.HitEntities == nil {
		t.Error("HitEntities map should be initialized")
	}
}

func TestProjectilePool_Return(t *testing.T) {
	pool := NewProjectilePool(10)
	proj := pool.Get()
	initialSize := pool.Size()

	// Set some values on the projectile
	proj.ID = 123
	proj.Position = Vector2{X: 100, Y: 200}
	proj.Velocity = Vector2{X: 5, Y: 10}
	proj.Lifetime = 1.0
	proj.Active = true
	proj.Damage = 5
	proj.HitEntities[1] = true
	proj.HitEntities[2] = true

	pool.Return(proj)

	if pool.Size() != initialSize+1 {
		t.Errorf("Expected pool size %d after Return, got %d", initialSize+1, pool.Size())
	}

	// Verify projectile was reset
	if proj.ID != 0 {
		t.Error("ID should be reset to 0")
	}

	if proj.Position.X != 0 || proj.Position.Y != 0 {
		t.Error("Position should be reset to (0, 0)")
	}

	if proj.Velocity.X != 0 || proj.Velocity.Y != 0 {
		t.Error("Velocity should be reset to (0, 0)")
	}

	if proj.Lifetime != 0 {
		t.Error("Lifetime should be reset to 0")
	}

	if proj.Active {
		t.Error("Active should be reset to false")
	}

	if proj.Damage != 0 {
		t.Error("Damage should be reset to 0")
	}

	if len(proj.HitEntities) != 0 {
		t.Error("HitEntities map should be cleared")
	}
}

func TestProjectilePool_GetWhenEmpty(t *testing.T) {
	pool := NewProjectilePool(1)

	// Get all projectiles from pool
	proj1 := pool.Get()
	if proj1 == nil {
		t.Fatal("First Get returned nil")
	}

	// Pool should be empty now
	if pool.Size() != 0 {
		t.Errorf("Expected pool size 0, got %d", pool.Size())
	}

	// Get when pool is empty should create new projectile
	proj2 := pool.Get()
	if proj2 == nil {
		t.Fatal("Get on empty pool returned nil")
	}

	// Pool should still be empty
	if pool.Size() != 0 {
		t.Errorf("Expected pool size 0 after empty Get, got %d", pool.Size())
	}

	// Should be different instances
	if proj1 == proj2 {
		t.Error("Should create new projectile when pool is empty")
	}
}

func TestProjectilePool_ReuseProjectiles(t *testing.T) {
	pool := NewProjectilePool(5)

	// Get a projectile
	proj1 := pool.Get()
	proj1.ID = 999

	// Return it
	pool.Return(proj1)

	// Get another projectile - should be the same one (reused)
	proj2 := pool.Get()

	if proj2 != proj1 {
		t.Error("Expected to get the same projectile instance back")
	}

	// But it should be reset
	if proj2.ID != 0 {
		t.Error("Reused projectile should have been reset")
	}
}

func TestProjectilePool_MultipleGetReturn(t *testing.T) {
	pool := NewProjectilePool(10)
	initialSize := pool.Size()

	projectiles := make([]*Projectile, 5)

	// Get 5 projectiles
	for i := 0; i < 5; i++ {
		projectiles[i] = pool.Get()
	}

	if pool.Size() != initialSize-5 {
		t.Errorf("Expected pool size %d, got %d", initialSize-5, pool.Size())
	}

	// Return all projectiles
	for _, proj := range projectiles {
		pool.Return(proj)
	}

	if pool.Size() != initialSize {
		t.Errorf("Expected pool size back to %d, got %d", initialSize, pool.Size())
	}
}

func TestProjectilePool_Size(t *testing.T) {
	sizes := []int{0, 1, 10, 100, 1000}

	for _, size := range sizes {
		pool := NewProjectilePool(size)
		if pool.Size() != size {
			t.Errorf("Expected pool size %d, got %d", size, pool.Size())
		}
	}
}

func TestProjectilePool_ResetAllFields(t *testing.T) {
	pool := NewProjectilePool(1)
	proj := pool.Get()

	// Set all fields
	proj.ID = 123
	proj.Owner = &MockEntity{id: 1}
	proj.Position = Vector2{X: 100, Y: 200}
	proj.Velocity = Vector2{X: 50, Y: 75}
	proj.ProjectileType = ProjectileFireball
	proj.Damage = 10
	proj.Lifetime = 2.5
	proj.MaxLifetime = 5.0
	proj.Bounces = 3
	proj.MaxBounces = 5
	proj.Active = true
	proj.Width = 16
	proj.Height = 16
	proj.AffectedByGravity = true
	proj.HitEntities[1] = true

	pool.Return(proj)

	// Verify all fields are reset
	if proj.ID != 0 {
		t.Error("ID not reset")
	}
	if proj.Owner != nil {
		t.Error("Owner not reset")
	}
	if proj.Position.X != 0 || proj.Position.Y != 0 {
		t.Error("Position not reset")
	}
	if proj.Velocity.X != 0 || proj.Velocity.Y != 0 {
		t.Error("Velocity not reset")
	}
	if proj.Damage != 0 {
		t.Error("Damage not reset")
	}
	if proj.Lifetime != 0 {
		t.Error("Lifetime not reset")
	}
	if proj.MaxLifetime != 0 {
		t.Error("MaxLifetime not reset")
	}
	if proj.Bounces != 0 {
		t.Error("Bounces not reset")
	}
	if proj.MaxBounces != 0 {
		t.Error("MaxBounces not reset")
	}
	if proj.Active {
		t.Error("Active not reset")
	}
	if proj.Width != 0 {
		t.Error("Width not reset")
	}
	if proj.Height != 0 {
		t.Error("Height not reset")
	}
	if proj.AffectedByGravity {
		t.Error("AffectedByGravity not reset")
	}
	if len(proj.HitEntities) != 0 {
		t.Error("HitEntities not cleared")
	}
}
