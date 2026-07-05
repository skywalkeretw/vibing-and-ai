package particles

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestExplosion(t *testing.T) {
	ps := NewParticleSystem()
	ps.Explosion(10, 20)

	if ps.GetParticleCount() != 20 {
		t.Errorf("Expected 20 particles for explosion, got %d", ps.GetParticleCount())
	}

	// Verify all particles are explosion type
	for _, p := range ps.particles {
		if p.Char != '*' {
			t.Errorf("Expected explosion char '*', got '%c'", p.Char)
		}
	}
}

func TestCoinCollect(t *testing.T) {
	ps := NewParticleSystem()
	ps.CoinCollect(5, 10)

	if ps.GetParticleCount() != 5 {
		t.Errorf("Expected 5 particles for coin collect, got %d", ps.GetParticleCount())
	}
}

func TestPowerUpActivate(t *testing.T) {
	ps := NewParticleSystem()
	ps.PowerUpActivate(15, 25)

	if ps.GetParticleCount() != 15 {
		t.Errorf("Expected 15 particles for power-up, got %d", ps.GetParticleCount())
	}
}

func TestJumpDust(t *testing.T) {
	ps := NewParticleSystem()
	ps.JumpDust(8, 12)

	if ps.GetParticleCount() != 3 {
		t.Errorf("Expected 3 particles for jump dust, got %d", ps.GetParticleCount())
	}
}

func TestLandingDust(t *testing.T) {
	ps := NewParticleSystem()
	ps.LandingDust(8, 12)

	if ps.GetParticleCount() != 5 {
		t.Errorf("Expected 5 particles for landing dust, got %d", ps.GetParticleCount())
	}
}

func TestEnemyDefeat(t *testing.T) {
	ps := NewParticleSystem()
	ps.EnemyDefeat(20, 30)

	if ps.GetParticleCount() != 10 {
		t.Errorf("Expected 10 particles for enemy defeat, got %d", ps.GetParticleCount())
	}
}

func TestWalkDust(t *testing.T) {
	ps := NewParticleSystem()
	ps.WalkDust(5, 10)

	if ps.GetParticleCount() != 1 {
		t.Errorf("Expected 1 particle for walk dust, got %d", ps.GetParticleCount())
	}
}

func TestSmokeTrail(t *testing.T) {
	ps := NewParticleSystem()
	emitter := ps.SmokeTrail(10, 20, 1.0)

	if emitter == nil {
		t.Fatal("SmokeTrail returned nil emitter")
	}

	if ps.GetEmitterCount() != 1 {
		t.Errorf("Expected 1 emitter, got %d", ps.GetEmitterCount())
	}

	if emitter.ParticleType != ParticleSmoke {
		t.Error("Expected smoke particle type")
	}
}

func TestSparkleEffect(t *testing.T) {
	ps := NewParticleSystem()
	emitter := ps.SparkleEffect(15, 25, 2.0)

	if emitter == nil {
		t.Fatal("SparkleEffect returned nil emitter")
	}

	if emitter.ParticleType != ParticleSparkle {
		t.Error("Expected sparkle particle type")
	}
}

func TestFireEffect(t *testing.T) {
	ps := NewParticleSystem()
	ps.FireEffect(10, 20)

	// Fire effect creates 3 explosion + 2 smoke = 5 particles
	if ps.GetParticleCount() != 5 {
		t.Errorf("Expected 5 particles for fire effect, got %d", ps.GetParticleCount())
	}
}

func TestWaterSplash(t *testing.T) {
	ps := NewParticleSystem()
	ps.WaterSplash(10, 20)

	if ps.GetParticleCount() != 8 {
		t.Errorf("Expected 8 particles for water splash, got %d", ps.GetParticleCount())
	}

	// Verify particles have water properties
	for _, p := range ps.particles {
		if p.Color != tcell.ColorBlue {
			t.Errorf("Expected blue color for water, got %v", p.Color)
		}
		if p.Char != '·' {
			t.Errorf("Expected '·' char for water, got '%c'", p.Char)
		}
	}
}

func TestImpactEffect(t *testing.T) {
	ps := NewParticleSystem()
	ps.ImpactEffect(5, 10)

	if ps.GetParticleCount() != 6 {
		t.Errorf("Expected 6 particles for impact, got %d", ps.GetParticleCount())
	}
}

func TestHealEffect(t *testing.T) {
	ps := NewParticleSystem()
	ps.HealEffect(15, 25)

	if ps.GetParticleCount() != 10 {
		t.Errorf("Expected 10 particles for heal effect, got %d", ps.GetParticleCount())
	}

	// Verify particles have heal properties
	for _, p := range ps.particles {
		if p.Color != tcell.ColorGreen {
			t.Errorf("Expected green color for heal, got %v", p.Color)
		}
	}
}

func TestDamageEffect(t *testing.T) {
	ps := NewParticleSystem()
	ps.DamageEffect(20, 30)

	if ps.GetParticleCount() != 5 {
		t.Errorf("Expected 5 particles for damage effect, got %d", ps.GetParticleCount())
	}
}

func TestMultiplePresets(t *testing.T) {
	ps := NewParticleSystem()

	ps.Explosion(10, 20)
	ps.CoinCollect(15, 25)
	ps.PowerUpActivate(20, 30)

	expectedTotal := 20 + 5 + 15
	if ps.GetParticleCount() != expectedTotal {
		t.Errorf("Expected %d total particles, got %d", expectedTotal, ps.GetParticleCount())
	}
}
