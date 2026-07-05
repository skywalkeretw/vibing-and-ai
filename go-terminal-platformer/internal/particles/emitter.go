package particles

// Emitter continuously emits particles at a specified rate
type Emitter struct {
	X             float64
	Y             float64
	Active        bool
	EmitRate      float64
	EmitTimer     float64
	ParticleType  ParticleType
	Duration      float64
	DurationTimer float64
}

// NewEmitter creates a new particle emitter
func NewEmitter(x, y float64, pType ParticleType, emitRate, duration float64) *Emitter {
	return &Emitter{
		X:            x,
		Y:            y,
		Active:       true,
		EmitRate:     emitRate,
		ParticleType: pType,
		Duration:     duration,
	}
}

// SetPosition updates the emitter's position
func (e *Emitter) SetPosition(x, y float64) {
	e.X = x
	e.Y = y
}

// Stop deactivates the emitter
func (e *Emitter) Stop() {
	e.Active = false
}

// Start activates the emitter
func (e *Emitter) Start() {
	e.Active = true
	e.DurationTimer = 0
}

// IsActive returns whether the emitter is currently active
func (e *Emitter) IsActive() bool {
	return e.Active
}

// SetEmitRate changes the emission rate
func (e *Emitter) SetEmitRate(rate float64) {
	e.EmitRate = rate
}

// GetEmitRate returns the current emission rate
func (e *Emitter) GetEmitRate() float64 {
	return e.EmitRate
}
