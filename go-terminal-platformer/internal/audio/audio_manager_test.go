package audio

import (
	"testing"
)

func TestNewAudioManager(t *testing.T) {
	am := NewAudioManager()

	if am == nil {
		t.Fatal("NewAudioManager returned nil")
	}

	if am.musicVolume != 0.7 {
		t.Errorf("Expected default music volume 0.7, got %f", am.musicVolume)
	}

	if am.sfxVolume != 0.8 {
		t.Errorf("Expected default SFX volume 0.8, got %f", am.sfxVolume)
	}

	if !am.enabled {
		t.Error("Expected audio to be enabled by default")
	}

	if am.sounds == nil {
		t.Error("Expected sounds map to be initialized")
	}
}

func TestAudioManagerInitialize(t *testing.T) {
	am := NewAudioManager()

	// Initialize should not return an error even if audio is unavailable
	err := am.Initialize()
	if err != nil {
		t.Errorf("Initialize returned error: %v", err)
	}

	// Should be able to call Initialize multiple times
	err = am.Initialize()
	if err != nil {
		t.Errorf("Second Initialize returned error: %v", err)
	}
}

func TestAudioManagerVolumeControl(t *testing.T) {
	am := NewAudioManager()

	tests := []struct {
		name     string
		volume   float64
		expected float64
	}{
		{"Normal volume", 0.5, 0.5},
		{"Zero volume", 0.0, 0.0},
		{"Max volume", 1.0, 1.0},
		{"Above max (clamped)", 1.5, 1.0},
		{"Below min (clamped)", -0.5, 0.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am.SetMusicVolume(tt.volume)
			if am.GetMusicVolume() != tt.expected {
				t.Errorf("SetMusicVolume(%f): expected %f, got %f", tt.volume, tt.expected, am.GetMusicVolume())
			}

			am.SetSFXVolume(tt.volume)
			if am.GetSFXVolume() != tt.expected {
				t.Errorf("SetSFXVolume(%f): expected %f, got %f", tt.volume, tt.expected, am.GetSFXVolume())
			}
		})
	}
}

func TestAudioManagerSetMasterVolume(t *testing.T) {
	am := NewAudioManager()

	am.SetMasterVolume(0.6)

	if am.GetMusicVolume() != 0.6 {
		t.Errorf("Expected music volume 0.6, got %f", am.GetMusicVolume())
	}

	if am.GetSFXVolume() != 0.6 {
		t.Errorf("Expected SFX volume 0.6, got %f", am.GetSFXVolume())
	}
}

func TestAudioManagerMuteUnmute(t *testing.T) {
	am := NewAudioManager()

	if !am.IsEnabled() {
		t.Error("Expected audio to be enabled initially")
	}

	am.Mute()

	if am.IsEnabled() {
		t.Error("Expected audio to be disabled after Mute()")
	}

	am.Unmute()

	if !am.IsEnabled() {
		t.Error("Expected audio to be enabled after Unmute()")
	}
}

func TestAudioManagerPlaySoundWhenDisabled(t *testing.T) {
	am := NewAudioManager()
	am.Mute()

	// Should not panic when playing sound while disabled
	am.PlaySound(SoundJump)
	am.PlaySoundAt(SoundCoinCollect, 10, 10, 0, 0)
}

func TestAudioManagerShutdown(t *testing.T) {
	am := NewAudioManager()
	err := am.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Should not panic
	am.Shutdown()

	// Should be able to call Shutdown multiple times
	am.Shutdown()
}

func TestClamp(t *testing.T) {
	tests := []struct {
		name     string
		value    float64
		min      float64
		max      float64
		expected float64
	}{
		{"Within range", 0.5, 0.0, 1.0, 0.5},
		{"Below min", -0.5, 0.0, 1.0, 0.0},
		{"Above max", 1.5, 0.0, 1.0, 1.0},
		{"At min", 0.0, 0.0, 1.0, 0.0},
		{"At max", 1.0, 0.0, 1.0, 1.0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := clamp(tt.value, tt.min, tt.max)
			if result != tt.expected {
				t.Errorf("clamp(%f, %f, %f): expected %f, got %f", tt.value, tt.min, tt.max, tt.expected, result)
			}
		})
	}
}

func TestAudioManagerConcurrency(t *testing.T) {
	am := NewAudioManager()
	err := am.Initialize()
	if err != nil {
		t.Fatalf("Initialize failed: %v", err)
	}

	// Test concurrent access to volume controls
	done := make(chan bool)

	go func() {
		for i := 0; i < 100; i++ {
			am.SetMusicVolume(0.5)
			am.GetMusicVolume()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			am.SetSFXVolume(0.7)
			am.GetSFXVolume()
		}
		done <- true
	}()

	go func() {
		for i := 0; i < 100; i++ {
			am.Mute()
			am.Unmute()
			am.IsEnabled()
		}
		done <- true
	}()

	// Wait for all goroutines
	<-done
	<-done
	<-done

	am.Shutdown()
}

func TestSoundTypes(t *testing.T) {
	// Verify all sound types are defined
	soundTypes := []SoundType{
		SoundJump,
		SoundLand,
		SoundCoinCollect,
		SoundPowerUp,
		SoundDamage,
		SoundDeath,
		SoundFireball,
		SoundEnemyHit,
		SoundEnemyDefeat,
		SoundEnemySpawn,
		SoundMenuSelect,
		SoundMenuMove,
		SoundPause,
		SoundLevelComplete,
		SoundGameOver,
	}

	// Just verify they're all different values
	seen := make(map[SoundType]bool)
	for _, st := range soundTypes {
		if seen[st] {
			t.Errorf("Duplicate sound type value: %d", st)
		}
		seen[st] = true
	}
}

func TestMusicTypes(t *testing.T) {
	// Verify all music types are defined
	musicTypes := []MusicType{
		MusicMainMenu,
		MusicWorld1,
		MusicWorld2,
		MusicWorld3,
		MusicBoss,
		MusicGameOver,
		MusicLevelComplete,
	}

	// Just verify they're all different values
	seen := make(map[MusicType]bool)
	for _, mt := range musicTypes {
		if seen[mt] {
			t.Errorf("Duplicate music type value: %d", mt)
		}
		seen[mt] = true
	}
}
