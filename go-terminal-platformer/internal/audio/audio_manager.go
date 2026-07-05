package audio

import (
	"log"
	"math"
	"os"
	"sync"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

// SoundType represents different types of sound effects
type SoundType int

const (
	// Player sounds
	SoundJump SoundType = iota
	SoundLand
	SoundCoinCollect
	SoundPowerUp
	SoundDamage
	SoundDeath
	SoundFireball

	// Enemy sounds
	SoundEnemyHit
	SoundEnemyDefeat
	SoundEnemySpawn

	// UI sounds
	SoundMenuSelect
	SoundMenuMove
	SoundPause
	SoundLevelComplete
	SoundGameOver
)

// MusicType represents different types of background music
type MusicType int

const (
	MusicMainMenu MusicType = iota
	MusicWorld1
	MusicWorld2
	MusicWorld3
	MusicBoss
	MusicGameOver
	MusicLevelComplete
)

// Sound represents a sound effect
type Sound struct {
	buffer   *beep.Buffer
	duration time.Duration
}

// Music represents background music
type Music struct {
	streamer beep.StreamSeekCloser
	ctrl     *beep.Ctrl
	loop     bool
	playing  bool
}

// AudioManager manages all audio playback
type AudioManager struct {
	enabled      bool
	musicVolume  float64
	sfxVolume    float64
	currentMusic *Music
	sounds       map[SoundType]*Sound
	musicPlayer  *MusicPlayer
	sfxPlayer    *SFXPlayer
	mu           sync.RWMutex
	initialized  bool
}

// NewAudioManager creates a new audio manager
func NewAudioManager() *AudioManager {
	return &AudioManager{
		enabled:     true,
		musicVolume: 0.7,
		sfxVolume:   0.8,
		sounds:      make(map[SoundType]*Sound),
	}
}

// Initialize initializes the audio system
func (am *AudioManager) Initialize() error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.initialized {
		return nil
	}

	// Check if audio is available
	am.enabled = am.checkAudioSupport()

	if !am.enabled {
		log.Println("Audio not available, running in silent mode")
		return nil
	}

	// Initialize speaker with sample rate
	sampleRate := beep.SampleRate(44100)
	err := speaker.Init(sampleRate, sampleRate.N(time.Second/10))
	if err != nil {
		log.Printf("Failed to initialize speaker: %v", err)
		am.enabled = false
		return nil
	}

	// Initialize audio players
	am.musicPlayer = NewMusicPlayer()
	am.sfxPlayer = NewSFXPlayer()

	// Load sound effects
	err = am.loadSounds()
	if err != nil {
		log.Printf("Failed to load sounds: %v", err)
		am.enabled = false
		return nil
	}

	am.initialized = true
	log.Println("Audio system initialized")
	return nil
}

// checkAudioSupport checks if audio libraries are available
func (am *AudioManager) checkAudioSupport() bool {
	// Check if audio libraries are available
	// This is platform-specific
	return true
}

// loadSounds loads all sound effects
func (am *AudioManager) loadSounds() error {
	soundFiles := map[SoundType]string{
		SoundJump:          "assets/audio/sfx/jump.wav",
		SoundLand:          "assets/audio/sfx/land.wav",
		SoundCoinCollect:   "assets/audio/sfx/coin.wav",
		SoundPowerUp:       "assets/audio/sfx/powerup.wav",
		SoundDamage:        "assets/audio/sfx/damage.wav",
		SoundDeath:         "assets/audio/sfx/death.wav",
		SoundFireball:      "assets/audio/sfx/fireball.wav",
		SoundEnemyHit:      "assets/audio/sfx/enemy_hit.wav",
		SoundEnemyDefeat:   "assets/audio/sfx/enemy_defeat.wav",
		SoundEnemySpawn:    "assets/audio/sfx/enemy_spawn.wav",
		SoundMenuSelect:    "assets/audio/sfx/menu_select.wav",
		SoundMenuMove:      "assets/audio/sfx/menu_move.wav",
		SoundPause:         "assets/audio/sfx/pause.wav",
		SoundLevelComplete: "assets/audio/sfx/level_complete.wav",
		SoundGameOver:      "assets/audio/sfx/game_over.wav",
	}

	for soundType, filepath := range soundFiles {
		sound, err := am.loadSound(filepath)
		if err != nil {
			log.Printf("Failed to load sound %s: %v", filepath, err)
			continue
		}
		am.sounds[soundType] = sound
	}

	return nil
}

// loadSound loads a single sound effect
func (am *AudioManager) loadSound(filepath string) (*Sound, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	streamer, format, err := wav.Decode(f)
	if err != nil {
		return nil, err
	}
	defer streamer.Close()

	buffer := beep.NewBuffer(format)
	buffer.Append(streamer)

	return &Sound{
		buffer:   buffer,
		duration: format.SampleRate.D(buffer.Len()),
	}, nil
}

// PlaySound plays a sound effect
func (am *AudioManager) PlaySound(soundType SoundType) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	if !am.enabled || !am.initialized {
		return
	}

	sound := am.sounds[soundType]
	if sound == nil {
		return
	}

	// Play sound asynchronously
	go am.sfxPlayer.Play(sound, am.sfxVolume)
}

// PlaySoundAt plays a sound effect with spatial audio (distance-based volume)
func (am *AudioManager) PlaySoundAt(soundType SoundType, x, y float64, listenerX, listenerY float64) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	if !am.enabled || !am.initialized {
		return
	}

	// Calculate distance-based volume
	distance := math.Sqrt(math.Pow(x-listenerX, 2) + math.Pow(y-listenerY, 2))
	maxDistance := 20.0

	if distance > maxDistance {
		return // Too far to hear
	}

	volume := am.sfxVolume * (1.0 - distance/maxDistance)

	sound := am.sounds[soundType]
	if sound == nil {
		return
	}

	go am.sfxPlayer.Play(sound, volume)
}

// PlayMusic starts playing background music
func (am *AudioManager) PlayMusic(musicType MusicType) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if !am.enabled || !am.initialized {
		return nil
	}

	// Stop current music
	am.stopMusicInternal()

	// Load music file
	musicFile := am.getMusicFile(musicType)
	music, err := am.loadMusic(musicFile)
	if err != nil {
		return err
	}

	music.loop = true
	am.currentMusic = music

	// Start playing
	go am.musicPlayer.Play(music, am.musicVolume)

	return nil
}

// loadMusic loads a music file
func (am *AudioManager) loadMusic(filepath string) (*Music, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	streamer, _, err := mp3.Decode(f)
	if err != nil {
		f.Close()
		return nil, err
	}

	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer), Paused: false}

	return &Music{
		streamer: streamer,
		ctrl:     ctrl,
		loop:     true,
		playing:  true,
	}, nil
}

// StopMusic stops the currently playing music
func (am *AudioManager) StopMusic() {
	am.mu.Lock()
	defer am.mu.Unlock()
	am.stopMusicInternal()
}

// stopMusicInternal stops music without locking (internal use)
func (am *AudioManager) stopMusicInternal() {
	if am.currentMusic != nil && am.currentMusic.playing {
		am.musicPlayer.Stop()
		if am.currentMusic.streamer != nil {
			am.currentMusic.streamer.Close()
		}
		am.currentMusic.playing = false
		am.currentMusic = nil
	}
}

// PauseMusic pauses the currently playing music
func (am *AudioManager) PauseMusic() {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.currentMusic != nil && am.currentMusic.playing {
		am.musicPlayer.Pause()
	}
}

// ResumeMusic resumes the paused music
func (am *AudioManager) ResumeMusic() {
	am.mu.Lock()
	defer am.mu.Unlock()

	if am.currentMusic != nil && !am.currentMusic.playing {
		am.musicPlayer.Resume()
	}
}

// getMusicFile returns the file path for a music type
func (am *AudioManager) getMusicFile(musicType MusicType) string {
	musicFiles := map[MusicType]string{
		MusicMainMenu:      "assets/audio/music/main_menu.mp3",
		MusicWorld1:        "assets/audio/music/world1.mp3",
		MusicWorld2:        "assets/audio/music/world2.mp3",
		MusicWorld3:        "assets/audio/music/world3.mp3",
		MusicBoss:          "assets/audio/music/boss.mp3",
		MusicGameOver:      "assets/audio/music/game_over.mp3",
		MusicLevelComplete: "assets/audio/music/level_complete.mp3",
	}
	return musicFiles[musicType]
}

// SetMusicVolume sets the music volume (0.0 to 1.0)
func (am *AudioManager) SetMusicVolume(volume float64) {
	am.mu.Lock()
	defer am.mu.Unlock()

	am.musicVolume = clamp(volume, 0.0, 1.0)
	if am.currentMusic != nil && am.currentMusic.playing {
		am.musicPlayer.SetVolume(am.musicVolume)
	}
}

// SetSFXVolume sets the sound effects volume (0.0 to 1.0)
func (am *AudioManager) SetSFXVolume(volume float64) {
	am.mu.Lock()
	defer am.mu.Unlock()

	am.sfxVolume = clamp(volume, 0.0, 1.0)
}

// SetMasterVolume sets both music and SFX volume (0.0 to 1.0)
func (am *AudioManager) SetMasterVolume(volume float64) {
	volume = clamp(volume, 0.0, 1.0)
	am.SetMusicVolume(volume)
	am.SetSFXVolume(volume)
}

// Mute mutes all audio
func (am *AudioManager) Mute() {
	am.mu.Lock()
	defer am.mu.Unlock()

	am.enabled = false
	am.stopMusicInternal()
}

// Unmute unmutes all audio
func (am *AudioManager) Unmute() {
	am.mu.Lock()
	defer am.mu.Unlock()

	am.enabled = true
}

// IsEnabled returns whether audio is enabled
func (am *AudioManager) IsEnabled() bool {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.enabled
}

// GetMusicVolume returns the current music volume
func (am *AudioManager) GetMusicVolume() float64 {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.musicVolume
}

// GetSFXVolume returns the current SFX volume
func (am *AudioManager) GetSFXVolume() float64 {
	am.mu.RLock()
	defer am.mu.RUnlock()
	return am.sfxVolume
}

// Shutdown cleans up audio resources
func (am *AudioManager) Shutdown() {
	am.mu.Lock()
	defer am.mu.Unlock()

	am.stopMusicInternal()

	if am.musicPlayer != nil {
		am.musicPlayer.Close()
	}

	if am.sfxPlayer != nil {
		am.sfxPlayer.Close()
	}

	am.initialized = false
	log.Println("Audio system shutdown")
}

// clamp restricts a value between min and max
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
