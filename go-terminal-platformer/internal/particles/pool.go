package particles

// ParticlePool manages a pool of reusable particles for performance
type ParticlePool struct {
	pool []*Particle
}

// NewParticlePool creates a new particle pool with the specified size
func NewParticlePool(size int) *ParticlePool {
	pool := make([]*Particle, 0, size)
	for i := 0; i < size; i++ {
		pool = append(pool, &Particle{})
	}
	return &ParticlePool{pool: pool}
}

// Get retrieves a particle from the pool
func (pp *ParticlePool) Get() *Particle {
	if len(pp.pool) > 0 {
		p := pp.pool[len(pp.pool)-1]
		pp.pool = pp.pool[:len(pp.pool)-1]
		return p
	}
	// Pool exhausted, create new particle
	return &Particle{}
}

// Return returns a particle to the pool
func (pp *ParticlePool) Return(p *Particle) {
	// Reset particle state
	p.X = 0
	p.Y = 0
	p.VelocityX = 0
	p.VelocityY = 0
	p.Lifetime = 0
	p.MaxLifetime = 0
	p.Gravity = 0
	p.FadeOut = false

	pp.pool = append(pp.pool, p)
}

// Size returns the current number of available particles in the pool
func (pp *ParticlePool) Size() int {
	return len(pp.pool)
}
