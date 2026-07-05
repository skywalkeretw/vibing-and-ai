package particles

import (
	"github.com/gdamore/tcell/v2"
)

// Explosion creates an explosion particle effect
func (ps *ParticleSystem) Explosion(x, y float64) {
	for i := 0; i < 20; i++ {
		ps.CreateParticle(x, y, ParticleExplosion)
	}
}

// CoinCollect creates a coin collection particle effect
func (ps *ParticleSystem) CoinCollect(x, y float64) {
	for i := 0; i < 5; i++ {
		ps.CreateParticle(x, y, ParticleSparkle)
	}
}

// PowerUpActivate creates a power-up activation particle effect
func (ps *ParticleSystem) PowerUpActivate(x, y float64) {
	for i := 0; i < 15; i++ {
		ps.CreateParticle(x, y, ParticlePowerUp)
	}
}

// JumpDust creates dust particles when jumping
func (ps *ParticleSystem) JumpDust(x, y float64) {
	for i := 0; i < 3; i++ {
		ps.CreateParticle(x, y, ParticleDust)
	}
}

// LandingDust creates dust particles when landing
func (ps *ParticleSystem) LandingDust(x, y float64) {
	for i := 0; i < 5; i++ {
		ps.CreateParticle(x, y, ParticleDust)
	}
}

// EnemyDefeat creates particles when an enemy is defeated
func (ps *ParticleSystem) EnemyDefeat(x, y float64) {
	for i := 0; i < 10; i++ {
		ps.CreateParticle(x, y, ParticleBlood)
	}
}

// WalkDust creates small dust particles while walking
func (ps *ParticleSystem) WalkDust(x, y float64) {
	ps.CreateParticle(x, y, ParticleDust)
}

// SmokeTrail creates a continuous smoke trail emitter
func (ps *ParticleSystem) SmokeTrail(x, y float64, duration float64) *Emitter {
	return ps.CreateEmitter(x, y, ParticleSmoke, duration)
}

// SparkleEffect creates a sparkle effect emitter
func (ps *ParticleSystem) SparkleEffect(x, y float64, duration float64) *Emitter {
	return ps.CreateEmitter(x, y, ParticleSparkle, duration)
}

// FireEffect creates a fire effect with multiple particle types
func (ps *ParticleSystem) FireEffect(x, y float64) {
	// Create explosion particles for fire
	for i := 0; i < 3; i++ {
		ps.CreateParticle(x, y, ParticleExplosion)
	}
	// Create smoke particles
	for i := 0; i < 2; i++ {
		ps.CreateParticle(x, y, ParticleSmoke)
	}
}

// WaterSplash creates a water splash effect
func (ps *ParticleSystem) WaterSplash(x, y float64) {
	for i := 0; i < 8; i++ {
		p := ps.CreateParticle(x, y, ParticleSparkle)
		// Override color for water effect
		p.Color = tcell.ColorBlue
		p.Char = '·'
	}
}

// ImpactEffect creates an impact particle effect
func (ps *ParticleSystem) ImpactEffect(x, y float64) {
	for i := 0; i < 6; i++ {
		ps.CreateParticle(x, y, ParticleDust)
	}
}

// HealEffect creates a healing particle effect
func (ps *ParticleSystem) HealEffect(x, y float64) {
	for i := 0; i < 10; i++ {
		p := ps.CreateParticle(x, y, ParticleSparkle)
		// Override color for heal effect
		p.Color = tcell.ColorGreen
	}
}

// DamageEffect creates a damage particle effect
func (ps *ParticleSystem) DamageEffect(x, y float64) {
	for i := 0; i < 5; i++ {
		ps.CreateParticle(x, y, ParticleBlood)
	}
}
