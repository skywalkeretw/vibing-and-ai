package particles

import (
	"testing"
)

func TestNewEmitter(t *testing.T) {
	e := NewEmitter(10, 20, ParticleSmoke, 0.05, 1.0)

	if e == nil {
		t.Fatal("NewEmitter returned nil")
	}

	if e.X != 10 {
		t.Errorf("Expected X=10, got %f", e.X)
	}

	if e.Y != 20 {
		t.Errorf("Expected Y=20, got %f", e.Y)
	}

	if !e.Active {
		t.Error("Expected emitter to be active")
	}

	if e.EmitRate != 0.05 {
		t.Errorf("Expected EmitRate=0.05, got %f", e.EmitRate)
	}

	if e.ParticleType != ParticleSmoke {
		t.Errorf("Expected ParticleSmoke, got %v", e.ParticleType)
	}

	if e.Duration != 1.0 {
		t.Errorf("Expected Duration=1.0, got %f", e.Duration)
	}
}

func TestEmitter_SetPosition(t *testing.T) {
	e := NewEmitter(10, 20, ParticleSmoke, 0.05, 1.0)

	e.SetPosition(30, 40)

	if e.X != 30 {
		t.Errorf("Expected X=30, got %f", e.X)
	}

	if e.Y != 40 {
		t.Errorf("Expected Y=40, got %f", e.Y)
	}
}

func TestEmitter_Stop(t *testing.T) {
	e := NewEmitter(10, 20, ParticleSmoke, 0.05, 1.0)

	if !e.Active {
		t.Fatal("Emitter should start active")
	}

	e.Stop()

	if e.Active {
		t.Error("Emitter should be inactive after Stop")
	}
}

func TestEmitter_Start(t *testing.T) {
	e := NewEmitter(10, 20, ParticleSmoke, 0.05, 1.0)
	e.Stop()

	if e.Active {
		t.Fatal("Emitter should be inactive")
	}

	e.Start()

	if !e.Active {
		t.Error("Emitter should be active after Start")
	}

	if e.DurationTimer != 0 {
		t.Error("DurationTimer should be reset to 0 on Start")
	}
}

func TestEmitter_IsActive(t *testing.T) {
	e := NewEmitter(10, 20, ParticleSmoke, 0.05, 1.0)

	if !e.IsActive() {
		t.Error("Expected IsActive to return true")
	}

	e.Stop()

	if e.IsActive() {
		t.Error("Expected IsActive to return false after Stop")
	}
}

func TestEmitter_SetGetEmitRate(t *testing.T) {
	e := NewEmitter(10, 20, ParticleSmoke, 0.05, 1.0)

	if e.GetEmitRate() != 0.05 {
		t.Errorf("Expected EmitRate=0.05, got %f", e.GetEmitRate())
	}

	e.SetEmitRate(0.1)

	if e.GetEmitRate() != 0.1 {
		t.Errorf("Expected EmitRate=0.1 after set, got %f", e.GetEmitRate())
	}
}

func TestEmitter_InfiniteEmitter(t *testing.T) {
	e := NewEmitter(10, 20, ParticleSmoke, 0.05, 0) // 0 duration = infinite

	if e.Duration != 0 {
		t.Errorf("Expected Duration=0 for infinite emitter, got %f", e.Duration)
	}
}

func TestEmitter_MultipleStopStart(t *testing.T) {
	e := NewEmitter(10, 20, ParticleSmoke, 0.05, 1.0)

	// Stop and start multiple times
	for i := 0; i < 5; i++ {
		e.Stop()
		if e.IsActive() {
			t.Errorf("Iteration %d: Expected inactive after Stop", i)
		}

		e.Start()
		if !e.IsActive() {
			t.Errorf("Iteration %d: Expected active after Start", i)
		}
	}
}
