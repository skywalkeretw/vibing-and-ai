package audio

import (
	"sync"

	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

// MusicPlayer handles background music playback
type MusicPlayer struct {
	current *Music
	volume  *effects.Volume
	mu      sync.RWMutex
}

// NewMusicPlayer creates a new music player
func NewMusicPlayer() *MusicPlayer {
	return &MusicPlayer{}
}

// Play starts playing music
func (mp *MusicPlayer) Play(music *Music, volume float64) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.current != nil {
		mp.stopInternal()
	}

	mp.current = music

	// Create volume control
	mp.volume = &effects.Volume{
		Streamer: music.ctrl,
		Base:     2,
		Volume:   volumeToDecibels(volume),
		Silent:   volume == 0,
	}

	speaker.Play(mp.volume)
	music.playing = true
}

// Stop stops the currently playing music
func (mp *MusicPlayer) Stop() {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.stopInternal()
}

// stopInternal stops music without locking (internal use)
func (mp *MusicPlayer) stopInternal() {
	if mp.current != nil {
		speaker.Clear()
		if mp.current.ctrl != nil {
			mp.current.ctrl.Paused = true
		}
		mp.current.playing = false
		mp.current = nil
	}
}

// Pause pauses the currently playing music
func (mp *MusicPlayer) Pause() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.current != nil && mp.current.ctrl != nil {
		speaker.Lock()
		mp.current.ctrl.Paused = true
		speaker.Unlock()
		mp.current.playing = false
	}
}

// Resume resumes the paused music
func (mp *MusicPlayer) Resume() {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.current != nil && mp.current.ctrl != nil {
		speaker.Lock()
		mp.current.ctrl.Paused = false
		speaker.Unlock()
		mp.current.playing = true
	}
}

// SetVolume sets the music volume
func (mp *MusicPlayer) SetVolume(volume float64) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	if mp.volume != nil {
		speaker.Lock()
		mp.volume.Volume = volumeToDecibels(volume)
		mp.volume.Silent = volume == 0
		speaker.Unlock()
	}
}

// IsPlaying returns whether music is currently playing
func (mp *MusicPlayer) IsPlaying() bool {
	mp.mu.RLock()
	defer mp.mu.RUnlock()

	return mp.current != nil && mp.current.playing
}

// Close cleans up the music player
func (mp *MusicPlayer) Close() {
	mp.Stop()
}

// volumeToDecibels converts a linear volume (0.0-1.0) to decibels
func volumeToDecibels(volume float64) float64 {
	if volume <= 0 {
		return -10 // Silent
	}
	// Convert linear volume to decibels
	// 0.0 = -10dB (silent), 1.0 = 0dB (full volume)
	return (volume - 1.0) * 10
}
