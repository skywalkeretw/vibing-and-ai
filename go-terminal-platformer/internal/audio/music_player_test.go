package audio

import (
	"testing"
)

func TestNewMusicPlayer(t *testing.T) {
	mp := NewMusicPlayer()

	if mp == nil {
		t.Fatal("NewMusicPlayer returned nil")
	}
}

func TestMusicPlayerStop(t *testing.T) {
	mp := NewMusicPlayer()

	// Should not panic when stopping with no music playing
	mp.Stop()
}

func TestMusicPlayerPause(t *testing.T) {
	mp := NewMusicPlayer()

	// Should not panic when pausing with no music playing
	mp.Pause()
}

func TestMusicPlayerResume(t *testing.T) {
	mp := NewMusicPlayer()

	// Should not panic when resuming with no music playing
	mp.Resume()
}

func TestMusicPlayerSetVolume(t *testing.T) {
	mp := NewMusicPlayer()

	// Should not panic when setting volume with no music playing
	mp.SetVolume(0.5)
	mp.SetVolume(0.0)
	mp.SetVolume(1.0)
}

func TestMusicPlayerIsPlaying(t *testing.T) {
	mp := NewMusicPlayer()

	if mp.IsPlaying() {
		t.Error("Expected IsPlaying to return false when no music is playing")
	}
}

func TestMusicPlayerClose(t *testing.T) {
	mp := NewMusicPlayer()

	// Should not panic
	mp.Close()

	// Should be able to call Close multiple times
	mp.Close()
}

func TestVolumeToDecibels(t *testing.T) {
	tests := []struct {
		name     string
		volume   float64
		expected float64
	}{
		{"Zero volume", 0.0, -10.0},
		{"Half volume", 0.5, -5.0},
		{"Full volume", 1.0, 0.0},
		{"Quarter volume", 0.25, -7.5},
		{"Three quarters volume", 0.75, -2.5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := volumeToDecibels(tt.volume)
			if result != tt.expected {
				t.Errorf("volumeToDecibels(%f): expected %f, got %f", tt.volume, tt.expected, result)
			}
		})
	}
}

func TestMusicPlayerConcurrency(t *testing.T) {
	mp := NewMusicPlayer()

	// Test concurrent access
	done := make(chan bool)

	go func() {
		for i := 0; i < 100; i++ {
			mp.SetVolume(0.5)
			mp.IsPlaying()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			mp.Pause()
			mp.Resume()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			mp.Stop()
		}
		done <- true
	}()

	// Wait for all goroutines
	<-done
	<-done
	<-done

	mp.Close()
}
