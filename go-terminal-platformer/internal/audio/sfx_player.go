package audio

import (
	"sync"

	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
)

// SFXPlayer handles sound effect playback
type SFXPlayer struct {
	mu sync.Mutex
}

// NewSFXPlayer creates a new SFX player
func NewSFXPlayer() *SFXPlayer {
	return &SFXPlayer{}
}

// Play plays a sound effect with the specified volume
func (sp *SFXPlayer) Play(sound *Sound, volume float64) {
	if sound == nil || sound.buffer == nil {
		return
	}

	sp.mu.Lock()
	defer sp.mu.Unlock()

	// Create a new streamer from the buffer
	streamer := sound.buffer.Streamer(0, sound.buffer.Len())

	// Apply volume control
	volumeCtrl := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   volumeToDecibels(volume),
		Silent:   volume == 0,
	}

	// Play the sound effect
	// beep.Seq will play the sound and then call Done when finished
	speaker.Play(beep.Seq(volumeCtrl, beep.Callback(func() {
		// Sound finished playing
	})))
}

// PlayWithPan plays a sound effect with volume and panning
// pan: -1.0 (left) to 1.0 (right), 0.0 is center
func (sp *SFXPlayer) PlayWithPan(sound *Sound, volume float64, pan float64) {
	if sound == nil || sound.buffer == nil {
		return
	}

	sp.mu.Lock()
	defer sp.mu.Unlock()

	// Create a new streamer from the buffer
	streamer := sound.buffer.Streamer(0, sound.buffer.Len())

	// Apply volume control
	volumeCtrl := &effects.Volume{
		Streamer: streamer,
		Base:     2,
		Volume:   volumeToDecibels(volume),
		Silent:   volume == 0,
	}

	// Apply panning
	// Note: beep doesn't have built-in panning, so we'll just use volume for now
	// In a full implementation, you'd adjust left/right channel volumes separately

	// Play the sound effect
	speaker.Play(beep.Seq(volumeCtrl, beep.Callback(func() {
		// Sound finished playing
	})))
}

// Close cleans up the SFX player
func (sp *SFXPlayer) Close() {
	// Nothing to clean up for SFX player
	// Individual sounds are cleaned up automatically after playing
}
