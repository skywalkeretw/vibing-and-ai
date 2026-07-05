package particles

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
)

// ParticleType defines different types of particle effects
type ParticleType int

const (
	ParticleExplosion ParticleType = iota
	ParticleSparkle
	ParticleSmoke
	ParticleDust
	ParticleCoin
	ParticlePowerUp
	ParticleBlood
	ParticleJump
)

// Particle represents a single particle in the system
type Particle struct {
	X           float64
	Y           float64
	VelocityX   float64
	VelocityY   float64
	Char        rune
	Color       tcell.Color
	Lifetime    float64
	MaxLifetime float64
	Gravity     float64
	FadeOut     bool
}

// ParticleSystem manages all particles and emitters
type ParticleSystem struct {
	particles []*Particle
	emitters  []*Emitter
	pool      *ParticlePool
}

// Camera interface for visibility checking
type Camera interface {
	IsVisible(x, y, width, height float64) bool
	WorldToScreen(worldX, worldY float64) (int, int)
}

// Renderer interface for drawing particles
type Renderer interface {
	DrawChar(x, y int, char rune, fg, bg tcell.Color)
}

// NewParticleSystem creates a new particle system
func NewParticleSystem() *ParticleSystem {
	return &ParticleSystem{
		particles: make([]*Particle, 0, 1000),
		emitters:  make([]*Emitter, 0, 50),
		pool:      NewParticlePool(1000),
	}
}

// CreateParticle creates a new particle of the specified type
func (ps *ParticleSystem) CreateParticle(x, y float64, pType ParticleType) *Particle {
	p := ps.pool.Get()

	switch pType {
	case ParticleExplosion:
		p.Char = '*'
		p.Color = tcell.ColorOrange
		p.VelocityX = (rand.Float64()*2 - 1) * 5
		p.VelocityY = (rand.Float64()*2 - 1) * 5
		p.Gravity = 0.2
		p.Lifetime = 0.5
		p.MaxLifetime = 0.5
		p.FadeOut = true

	case ParticleSparkle:
		p.Char = '✦'
		p.Color = tcell.ColorYellow
		p.VelocityX = (rand.Float64()*2 - 1) * 2
		p.VelocityY = -rand.Float64() * 3
		p.Gravity = 0.1
		p.Lifetime = 0.8
		p.MaxLifetime = 0.8
		p.FadeOut = true

	case ParticleSmoke:
		p.Char = '░'
		p.Color = tcell.ColorGray
		p.VelocityX = (rand.Float64()*2 - 1) * 0.5
		p.VelocityY = -rand.Float64() * 2
		p.Gravity = -0.05
		p.Lifetime = 1.5
		p.MaxLifetime = 1.5
		p.FadeOut = true

	case ParticleDust:
		p.Char = '·'
		p.Color = tcell.ColorBrown
		p.VelocityX = (rand.Float64()*2 - 1) * 1
		p.VelocityY = -rand.Float64()
		p.Gravity = 0.3
		p.Lifetime = 0.3
		p.MaxLifetime = 0.3
		p.FadeOut = false

	case ParticleCoin:
		p.Char = '○'
		p.Color = tcell.ColorYellow
		p.VelocityX = 0
		p.VelocityY = -3
		p.Gravity = 0.2
		p.Lifetime = 1.0
		p.MaxLifetime = 1.0
		p.FadeOut = true

	case ParticlePowerUp:
		p.Char = '✦'
		p.Color = tcell.ColorGreen
		p.VelocityX = (rand.Float64()*2 - 1) * 2
		p.VelocityY = -rand.Float64() * 4
		p.Gravity = 0.1
		p.Lifetime = 1.0
		p.MaxLifetime = 1.0
		p.FadeOut = true

	case ParticleBlood:
		p.Char = '•'
		p.Color = tcell.ColorRed
		p.VelocityX = (rand.Float64()*2 - 1) * 3
		p.VelocityY = (rand.Float64()*2 - 1) * 3
		p.Gravity = 0.3
		p.Lifetime = 0.4
		p.MaxLifetime = 0.4
		p.FadeOut = true

	case ParticleJump:
		p.Char = '·'
		p.Color = tcell.ColorWhite
		p.VelocityX = (rand.Float64()*2 - 1) * 0.5
		p.VelocityY = -rand.Float64() * 0.5
		p.Gravity = 0.2
		p.Lifetime = 0.2
		p.MaxLifetime = 0.2
		p.FadeOut = false
	}

	p.X = x
	p.Y = y

	ps.particles = append(ps.particles, p)
	return p
}

// Update updates all particles and emitters
func (ps *ParticleSystem) Update(deltaTime float64) {
	// Update particles
	for i := len(ps.particles) - 1; i >= 0; i-- {
		p := ps.particles[i]

		// Update physics
		p.VelocityY += p.Gravity * deltaTime
		p.X += p.VelocityX * deltaTime
		p.Y += p.VelocityY * deltaTime

		// Update lifetime
		p.Lifetime -= deltaTime

		// Remove dead particles
		if p.Lifetime <= 0 {
			ps.pool.Return(p)
			ps.particles = append(ps.particles[:i], ps.particles[i+1:]...)
		}
	}

	// Update emitters
	for i := len(ps.emitters) - 1; i >= 0; i-- {
		e := ps.emitters[i]

		if !e.Active {
			continue
		}

		// Emit particles
		e.EmitTimer += deltaTime
		if e.EmitTimer >= e.EmitRate {
			e.EmitTimer = 0
			ps.CreateParticle(e.X, e.Y, e.ParticleType)
		}

		// Update duration
		if e.Duration > 0 {
			e.DurationTimer += deltaTime
			if e.DurationTimer >= e.Duration {
				e.Active = false
				ps.emitters = append(ps.emitters[:i], ps.emitters[i+1:]...)
			}
		}
	}
}

// Render renders all particles
func (ps *ParticleSystem) Render(renderer Renderer, camera Camera) {
	for _, p := range ps.particles {
		// Check if visible
		if !camera.IsVisible(p.X, p.Y, 1, 1) {
			continue
		}

		// Convert to screen coordinates
		screenX, screenY := camera.WorldToScreen(p.X, p.Y)

		// Apply fade out
		color := p.Color
		if p.FadeOut {
			alpha := p.Lifetime / p.MaxLifetime
			color = ps.applyAlpha(p.Color, alpha)
		}

		renderer.DrawChar(screenX, screenY, p.Char, color, tcell.ColorDefault)
	}
}

// applyAlpha applies alpha transparency by fading to darker colors
func (ps *ParticleSystem) applyAlpha(color tcell.Color, alpha float64) tcell.Color {
	if alpha < 0.3 {
		return tcell.ColorDarkGray
	} else if alpha < 0.6 {
		return tcell.ColorGray
	}
	return color
}

// CreateEmitter creates a new particle emitter
func (ps *ParticleSystem) CreateEmitter(x, y float64, pType ParticleType, duration float64) *Emitter {
	e := &Emitter{
		X:            x,
		Y:            y,
		Active:       true,
		EmitRate:     0.05,
		ParticleType: pType,
		Duration:     duration,
	}

	ps.emitters = append(ps.emitters, e)
	return e
}

// GetParticleCount returns the current number of active particles
func (ps *ParticleSystem) GetParticleCount() int {
	return len(ps.particles)
}

// GetEmitterCount returns the current number of active emitters
func (ps *ParticleSystem) GetEmitterCount() int {
	return len(ps.emitters)
}

// Clear removes all particles and emitters
func (ps *ParticleSystem) Clear() {
	// Return all particles to pool
	for _, p := range ps.particles {
		ps.pool.Return(p)
	}
	ps.particles = ps.particles[:0]
	ps.emitters = ps.emitters[:0]
}
