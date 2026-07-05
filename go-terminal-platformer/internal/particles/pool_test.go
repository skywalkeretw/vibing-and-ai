package particles

import (
	"testing"
)

func TestNewParticlePool(t *testing.T) {
	pool := NewParticlePool(100)

	if pool == nil {
		t.Fatal("NewParticlePool returned nil")
	}

	if pool.Size() != 100 {
		t.Errorf("Expected pool size 100, got %d", pool.Size())
	}
}

func TestParticlePool_Get(t *testing.T) {
	pool := NewParticlePool(10)
	initialSize := pool.Size()

	p := pool.Get()

	if p == nil {
		t.Fatal("Get returned nil particle")
	}

	if pool.Size() != initialSize-1 {
		t.Errorf("Expected pool size %d after Get, got %d", initialSize-1, pool.Size())
	}
}

func TestParticlePool_Return(t *testing.T) {
	pool := NewParticlePool(10)
	p := pool.Get()
	initialSize := pool.Size()

	// Set some values on the particle
	p.X = 100
	p.Y = 200
	p.VelocityX = 5
	p.VelocityY = 10
	p.Lifetime = 1.0

	pool.Return(p)

	if pool.Size() != initialSize+1 {
		t.Errorf("Expected pool size %d after Return, got %d", initialSize+1, pool.Size())
	}

	// Verify particle was reset
	if p.X != 0 || p.Y != 0 {
		t.Error("Particle position should be reset to 0")
	}

	if p.VelocityX != 0 || p.VelocityY != 0 {
		t.Error("Particle velocity should be reset to 0")
	}

	if p.Lifetime != 0 {
		t.Error("Particle lifetime should be reset to 0")
	}
}

func TestParticlePool_GetWhenEmpty(t *testing.T) {
	pool := NewParticlePool(1)

	// Get all particles from pool
	p1 := pool.Get()
	if p1 == nil {
		t.Fatal("First Get returned nil")
	}

	// Pool should be empty now
	if pool.Size() != 0 {
		t.Errorf("Expected pool size 0, got %d", pool.Size())
	}

	// Get when pool is empty should create new particle
	p2 := pool.Get()
	if p2 == nil {
		t.Fatal("Get on empty pool returned nil")
	}

	// Pool should still be empty
	if pool.Size() != 0 {
		t.Errorf("Expected pool size 0 after empty Get, got %d", pool.Size())
	}
}

func TestParticlePool_ReuseParticles(t *testing.T) {
	pool := NewParticlePool(5)

	// Get a particle
	p1 := pool.Get()
	p1.X = 123.45

	// Return it
	pool.Return(p1)

	// Get another particle - should be the same one (reused)
	p2 := pool.Get()

	if p2 != p1 {
		t.Error("Expected to get the same particle instance back")
	}

	// But it should be reset
	if p2.X != 0 {
		t.Error("Reused particle should have been reset")
	}
}

func TestParticlePool_MultipleGetReturn(t *testing.T) {
	pool := NewParticlePool(10)
	initialSize := pool.Size()

	particles := make([]*Particle, 5)

	// Get 5 particles
	for i := 0; i < 5; i++ {
		particles[i] = pool.Get()
	}

	if pool.Size() != initialSize-5 {
		t.Errorf("Expected pool size %d, got %d", initialSize-5, pool.Size())
	}

	// Return all particles
	for _, p := range particles {
		pool.Return(p)
	}

	if pool.Size() != initialSize {
		t.Errorf("Expected pool size back to %d, got %d", initialSize, pool.Size())
	}
}

func TestParticlePool_Size(t *testing.T) {
	sizes := []int{0, 1, 10, 100, 1000}

	for _, size := range sizes {
		pool := NewParticlePool(size)
		if pool.Size() != size {
			t.Errorf("Expected pool size %d, got %d", size, pool.Size())
		}
	}
}
