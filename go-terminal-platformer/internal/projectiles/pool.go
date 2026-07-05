package projectiles

// ProjectilePool manages a pool of reusable projectiles for performance
type ProjectilePool struct {
	pool []*Projectile
}

// NewProjectilePool creates a new projectile pool with the specified size
func NewProjectilePool(size int) *ProjectilePool {
	pool := make([]*Projectile, 0, size)
	for i := 0; i < size; i++ {
		pool = append(pool, &Projectile{
			HitEntities: make(map[int]bool),
		})
	}
	return &ProjectilePool{pool: pool}
}

// Get retrieves a projectile from the pool
func (pp *ProjectilePool) Get() *Projectile {
	if len(pp.pool) > 0 {
		proj := pp.pool[len(pp.pool)-1]
		pp.pool = pp.pool[:len(pp.pool)-1]
		return proj
	}
	// Pool exhausted, create new projectile
	return &Projectile{
		HitEntities: make(map[int]bool),
	}
}

// Return returns a projectile to the pool
func (pp *ProjectilePool) Return(proj *Projectile) {
	// Reset projectile state
	proj.ID = 0
	proj.Owner = nil
	proj.Position = Vector2{X: 0, Y: 0}
	proj.Velocity = Vector2{X: 0, Y: 0}
	proj.Lifetime = 0
	proj.MaxLifetime = 0
	proj.Bounces = 0
	proj.MaxBounces = 0
	proj.Active = false
	proj.Damage = 0
	proj.Width = 0
	proj.Height = 0
	proj.AffectedByGravity = false

	// Clear hit entities map
	for k := range proj.HitEntities {
		delete(proj.HitEntities, k)
	}

	pp.pool = append(pp.pool, proj)
}

// Size returns the current number of available projectiles in the pool
func (pp *ProjectilePool) Size() int {
	return len(pp.pool)
}
