package particles

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

// MockCamera implements the Camera interface for testing
type MockCamera struct {
	x      float64
	y      float64
	width  int
	height int
}

func (m *MockCamera) IsVisible(x, y, width, height float64) bool {
	return x+width >= m.x &&
		x <= m.x+float64(m.width) &&
		y+height >= m.y &&
		y <= m.y+float64(m.height)
}

func (m *MockCamera) WorldToScreen(worldX, worldY float64) (int, int) {
	return int(worldX - m.x), int(worldY - m.y)
}

// MockRenderer implements the Renderer interface for testing
type MockRenderer struct {
	drawCalls []DrawCall
}

type DrawCall struct {
	x    int
	y    int
	char rune
	fg   tcell.Color
	bg   tcell.Color
}

func (m *MockRenderer) DrawChar(x, y int, char rune, fg, bg tcell.Color) {
	m.drawCalls = append(m.drawCalls, DrawCall{x, y, char, fg, bg})
}

func TestNewParticleSystem(t *testing.T) {
	ps := NewParticleSystem()

	if ps == nil {
		t.Fatal("NewParticleSystem returned nil")
	}

	if ps.particles == nil {
		t.Error("particles slice is nil")
	}

	if ps.emitters == nil {
		t.Error("emitters slice is nil")
	}

	if ps.pool == nil {
		t.Error("pool is nil")
	}

	if ps.GetParticleCount() != 0 {
		t.Errorf("Expected 0 particles, got %d", ps.GetParticleCount())
	}

	if ps.GetEmitterCount() != 0 {
		t.Errorf("Expected 0 emitters, got %d", ps.GetEmitterCount())
	}
}

func TestCreateParticle_Explosion(t *testing.T) {
	ps := NewParticleSystem()
	p := ps.CreateParticle(10, 20, ParticleExplosion)

	if p == nil {
		t.Fatal("CreateParticle returned nil")
	}

	if p.X != 10 {
		t.Errorf("Expected X=10, got %f", p.X)
	}

	if p.Y != 20 {
		t.Errorf("Expected Y=20, got %f", p.Y)
	}

	if p.Char != '*' {
		t.Errorf("Expected char '*', got '%c'", p.Char)
	}

	if p.Color != tcell.ColorOrange {
		t.Errorf("Expected color Orange, got %v", p.Color)
	}

	if p.Lifetime != 0.5 {
		t.Errorf("Expected lifetime 0.5, got %f", p.Lifetime)
	}

	if !p.FadeOut {
		t.Error("Expected FadeOut to be true")
	}

	if ps.GetParticleCount() != 1 {
		t.Errorf("Expected 1 particle, got %d", ps.GetParticleCount())
	}
}

func TestCreateParticle_Sparkle(t *testing.T) {
	ps := NewParticleSystem()
	p := ps.CreateParticle(5, 10, ParticleSparkle)

	if p.Char != '✦' {
		t.Errorf("Expected char '✦', got '%c'", p.Char)
	}

	if p.Color != tcell.ColorYellow {
		t.Errorf("Expected color Yellow, got %v", p.Color)
	}

	if p.Lifetime != 0.8 {
		t.Errorf("Expected lifetime 0.8, got %f", p.Lifetime)
	}
}

func TestCreateParticle_Smoke(t *testing.T) {
	ps := NewParticleSystem()
	p := ps.CreateParticle(5, 10, ParticleSmoke)

	if p.Char != '░' {
		t.Errorf("Expected char '░', got '%c'", p.Char)
	}

	if p.Color != tcell.ColorGray {
		t.Errorf("Expected color Gray, got %v", p.Color)
	}

	if p.Gravity != -0.05 {
		t.Errorf("Expected gravity -0.05, got %f", p.Gravity)
	}
}

func TestCreateParticle_Dust(t *testing.T) {
	ps := NewParticleSystem()
	p := ps.CreateParticle(5, 10, ParticleDust)

	if p.Char != '·' {
		t.Errorf("Expected char '·', got '%c'", p.Char)
	}

	if p.Color != tcell.ColorBrown {
		t.Errorf("Expected color Brown, got %v", p.Color)
	}

	if p.FadeOut {
		t.Error("Expected FadeOut to be false for dust")
	}
}

func TestCreateParticle_Coin(t *testing.T) {
	ps := NewParticleSystem()
	p := ps.CreateParticle(5, 10, ParticleCoin)

	if p.Char != '○' {
		t.Errorf("Expected char '○', got '%c'", p.Char)
	}

	if p.VelocityX != 0 {
		t.Errorf("Expected VelocityX 0, got %f", p.VelocityX)
	}

	if p.VelocityY != -3 {
		t.Errorf("Expected VelocityY -3, got %f", p.VelocityY)
	}
}

func TestUpdate_ParticlePhysics(t *testing.T) {
	ps := NewParticleSystem()
	p := ps.CreateParticle(10, 20, ParticleExplosion)

	initialX := p.X
	initialY := p.Y
	initialVelocityY := p.VelocityY
	initialLifetime := p.Lifetime

	ps.Update(0.1)

	// Position should change based on velocity
	if p.X == initialX && p.VelocityX != 0 {
		t.Error("X position should have changed")
	}

	if p.Y == initialY && p.VelocityY != 0 {
		t.Error("Y position should have changed")
	}

	// Velocity should change based on gravity
	if p.VelocityY == initialVelocityY && p.Gravity != 0 {
		t.Error("VelocityY should have changed due to gravity")
	}

	// Lifetime should decrease
	if p.Lifetime >= initialLifetime {
		t.Error("Lifetime should have decreased")
	}
}

func TestUpdate_ParticleRemoval(t *testing.T) {
	ps := NewParticleSystem()
	ps.CreateParticle(10, 20, ParticleDust) // Short lifetime (0.3s)

	if ps.GetParticleCount() != 1 {
		t.Fatal("Expected 1 particle")
	}

	// Update for longer than particle lifetime
	ps.Update(0.5)

	if ps.GetParticleCount() != 0 {
		t.Errorf("Expected 0 particles after expiry, got %d", ps.GetParticleCount())
	}
}

func TestCreateEmitter(t *testing.T) {
	ps := NewParticleSystem()
	e := ps.CreateEmitter(15, 25, ParticleSmoke, 1.0)

	if e == nil {
		t.Fatal("CreateEmitter returned nil")
	}

	if e.X != 15 {
		t.Errorf("Expected X=15, got %f", e.X)
	}

	if e.Y != 25 {
		t.Errorf("Expected Y=25, got %f", e.Y)
	}

	if !e.Active {
		t.Error("Expected emitter to be active")
	}

	if e.ParticleType != ParticleSmoke {
		t.Errorf("Expected ParticleSmoke, got %v", e.ParticleType)
	}

	if ps.GetEmitterCount() != 1 {
		t.Errorf("Expected 1 emitter, got %d", ps.GetEmitterCount())
	}
}

func TestUpdate_EmitterEmission(t *testing.T) {
	ps := NewParticleSystem()
	ps.CreateEmitter(10, 20, ParticleSmoke, 0) // Infinite duration

	initialCount := ps.GetParticleCount()

	// Update multiple times to trigger emission
	for i := 0; i < 10; i++ {
		ps.Update(0.01)
	}

	if ps.GetParticleCount() <= initialCount {
		t.Error("Emitter should have created particles")
	}
}

func TestUpdate_EmitterDuration(t *testing.T) {
	ps := NewParticleSystem()
	ps.CreateEmitter(10, 20, ParticleSmoke, 0.1) // 0.1s duration

	if ps.GetEmitterCount() != 1 {
		t.Fatal("Expected 1 emitter")
	}

	// Update for longer than emitter duration
	ps.Update(0.2)

	if ps.GetEmitterCount() != 0 {
		t.Errorf("Expected 0 emitters after duration, got %d", ps.GetEmitterCount())
	}
}

func TestRender(t *testing.T) {
	ps := NewParticleSystem()
	ps.CreateParticle(10, 20, ParticleExplosion)

	camera := &MockCamera{x: 0, y: 0, width: 80, height: 24}
	renderer := &MockRenderer{}

	ps.Render(renderer, camera)

	if len(renderer.drawCalls) != 1 {
		t.Errorf("Expected 1 draw call, got %d", len(renderer.drawCalls))
	}

	if len(renderer.drawCalls) > 0 {
		call := renderer.drawCalls[0]
		if call.x != 10 || call.y != 20 {
			t.Errorf("Expected draw at (10, 20), got (%d, %d)", call.x, call.y)
		}
		if call.char != '*' {
			t.Errorf("Expected char '*', got '%c'", call.char)
		}
	}
}

func TestRender_CullingOffScreen(t *testing.T) {
	ps := NewParticleSystem()
	ps.CreateParticle(100, 100, ParticleExplosion) // Outside camera view

	camera := &MockCamera{x: 0, y: 0, width: 80, height: 24}
	renderer := &MockRenderer{}

	ps.Render(renderer, camera)

	if len(renderer.drawCalls) != 0 {
		t.Errorf("Expected 0 draw calls for off-screen particle, got %d", len(renderer.drawCalls))
	}
}

func TestApplyAlpha(t *testing.T) {
	ps := NewParticleSystem()

	tests := []struct {
		alpha    float64
		expected tcell.Color
	}{
		{0.2, tcell.ColorDarkGray},
		{0.5, tcell.ColorGray},
		{0.8, tcell.ColorOrange},
		{1.0, tcell.ColorOrange},
	}

	for _, tt := range tests {
		result := ps.applyAlpha(tcell.ColorOrange, tt.alpha)
		if result != tt.expected {
			t.Errorf("applyAlpha(%f) = %v, expected %v", tt.alpha, result, tt.expected)
		}
	}
}

func TestClear(t *testing.T) {
	ps := NewParticleSystem()

	// Create some particles and emitters
	ps.CreateParticle(10, 20, ParticleExplosion)
	ps.CreateParticle(15, 25, ParticleSparkle)
	ps.CreateEmitter(5, 10, ParticleSmoke, 1.0)

	if ps.GetParticleCount() != 2 {
		t.Fatal("Expected 2 particles")
	}

	if ps.GetEmitterCount() != 1 {
		t.Fatal("Expected 1 emitter")
	}

	ps.Clear()

	if ps.GetParticleCount() != 0 {
		t.Errorf("Expected 0 particles after clear, got %d", ps.GetParticleCount())
	}

	if ps.GetEmitterCount() != 0 {
		t.Errorf("Expected 0 emitters after clear, got %d", ps.GetEmitterCount())
	}
}

func TestMultipleParticleTypes(t *testing.T) {
	ps := NewParticleSystem()

	types := []ParticleType{
		ParticleExplosion,
		ParticleSparkle,
		ParticleSmoke,
		ParticleDust,
		ParticleCoin,
		ParticlePowerUp,
		ParticleBlood,
		ParticleJump,
	}

	for _, pType := range types {
		p := ps.CreateParticle(10, 20, pType)
		if p == nil {
			t.Errorf("Failed to create particle of type %v", pType)
		}
	}

	if ps.GetParticleCount() != len(types) {
		t.Errorf("Expected %d particles, got %d", len(types), ps.GetParticleCount())
	}
}

func TestParticlePoolIntegration(t *testing.T) {
	ps := NewParticleSystem()

	// Create and destroy particles to test pooling
	for i := 0; i < 100; i++ {
		ps.CreateParticle(float64(i), float64(i), ParticleDust)
	}

	if ps.GetParticleCount() != 100 {
		t.Fatalf("Expected 100 particles, got %d", ps.GetParticleCount())
	}

	// Update to remove all particles (dust has 0.3s lifetime)
	ps.Update(0.5)

	if ps.GetParticleCount() != 0 {
		t.Errorf("Expected 0 particles after expiry, got %d", ps.GetParticleCount())
	}

	// Create new particles - should reuse from pool
	for i := 0; i < 50; i++ {
		ps.CreateParticle(float64(i), float64(i), ParticleDust)
	}

	if ps.GetParticleCount() != 50 {
		t.Errorf("Expected 50 particles, got %d", ps.GetParticleCount())
	}
}
