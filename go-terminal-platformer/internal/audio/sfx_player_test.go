package audio

import (
	"testing"
)

func TestNewSFXPlayer(t *testing.T) {
	sp := NewSFXPlayer()

	if sp == nil {
		t.Fatal("NewSFXPlayer returned nil")
	}
}

func TestSFXPlayerPlayWithNilSound(t *testing.T) {
	sp := NewSFXPlayer()

	// Should not panic when playing nil sound
	sp.Play(nil, 0.5)
}

func TestSFXPlayerPlayWithPanNilSound(t *testing.T) {
	sp := NewSFXPlayer()

	// Should not panic when playing nil sound with pan
	sp.PlayWithPan(nil, 0.5, 0.0)
}

func TestSFXPlayerClose(t *testing.T) {
	sp := NewSFXPlayer()

	// Should not panic
	sp.Close()

	// Should be able to call Close multiple times
	sp.Close()
}

func TestSFXPlayerConcurrency(t *testing.T) {
	sp := NewSFXPlayer()

	// Test concurrent access
	done := make(chan bool)

	go func() {
		for i := 0; i < 100; i++ {
			sp.Play(nil, 0.5)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			sp.PlayWithPan(nil, 0.7, 0.5)
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			sp.PlayWithPan(nil, 0.3, -0.5)
		}
		done <- true
	}()

	// Wait for all goroutines
	<-done
	<-done
	<-done

	sp.Close()
}

func TestSFXPlayerPlayWithPanValues(t *testing.T) {
	sp := NewSFXPlayer()

	tests := []struct {
		name   string
		volume float64
		pan    float64
	}{
		{"Center pan", 0.5, 0.0},
		{"Left pan", 0.5, -1.0},
		{"Right pan", 0.5, 1.0},
		{"Slight left", 0.7, -0.3},
		{"Slight right", 0.7, 0.3},
		{"Zero volume", 0.0, 0.0},
		{"Max volume", 1.0, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Should not panic with various pan and volume values
			sp.PlayWithPan(nil, tt.volume, tt.pan)
		})
	}
}
