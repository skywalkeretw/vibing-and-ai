package camera

import (
	"math"
	"math/rand"
)

// CameraMode defines different camera behavior modes
type CameraMode int

const (
	CameraModeFollow CameraMode = iota
	CameraModeLerp
	CameraModeLookAhead
	CameraModeFixed
)

// Rect represents a rectangular boundary
type Rect struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// Camera represents the game camera that follows the player
type Camera struct {
	x           float64
	y           float64
	width       int
	height      int
	targetX     float64
	targetY     float64
	smoothing   float64
	bounds      Rect
	shakeAmount float64
	shakeTime   float64
	mode        CameraMode
	lookAhead   float64
}

// Player interface defines what the camera needs from a player
type Player interface {
	GetPosition() (float64, float64)
	GetVelocity() (float64, float64)
}

// NewCamera creates a new camera with the specified dimensions and level bounds
func NewCamera(width, height int, levelBounds Rect) *Camera {
	return &Camera{
		x:         0,
		y:         0,
		width:     width,
		height:    height,
		smoothing: 0.1,
		bounds:    levelBounds,
		mode:      CameraModeLerp,
		lookAhead: 5.0,
	}
}

// Update updates the camera position based on the target player
func (c *Camera) Update(deltaTime float64, target Player) {
	if target == nil {
		return
	}

	// Get player position and velocity
	playerX, playerY := target.GetPosition()
	velocityX, _ := target.GetVelocity()

	// Calculate target position (center on player)
	c.targetX = playerX - float64(c.width)/2
	c.targetY = playerY - float64(c.height)/2

	switch c.mode {
	case CameraModeFollow:
		// Instant follow
		c.x = c.targetX
		c.y = c.targetY

	case CameraModeLerp:
		// Apply smoothing (lerp)
		c.x += (c.targetX - c.x) * c.smoothing
		c.y += (c.targetY - c.y) * c.smoothing

	case CameraModeLookAhead:
		// Apply smoothing with look-ahead
		c.x += (c.targetX - c.x) * c.smoothing
		c.y += (c.targetY - c.y) * c.smoothing

		// Apply look-ahead based on player direction
		if velocityX > 0 {
			c.x += c.lookAhead
		} else if velocityX < 0 {
			c.x -= c.lookAhead
		}

	case CameraModeFixed:
		// Camera doesn't move (but shake still applies)
	}

	// Apply screen shake
	if c.shakeTime > 0 {
		c.x += (rand.Float64()*2 - 1) * c.shakeAmount
		c.y += (rand.Float64()*2 - 1) * c.shakeAmount
		c.shakeTime -= deltaTime
		if c.shakeTime < 0 {
			c.shakeTime = 0
		}
	}

	// Constrain to level bounds
	c.constrainToBounds()
}

// constrainToBounds ensures the camera stays within level boundaries
func (c *Camera) constrainToBounds() {
	// Don't scroll past level boundaries
	if c.x < c.bounds.X {
		c.x = c.bounds.X
	}
	if c.y < c.bounds.Y {
		c.y = c.bounds.Y
	}
	if c.x+float64(c.width) > c.bounds.X+c.bounds.Width {
		c.x = c.bounds.X + c.bounds.Width - float64(c.width)
	}
	if c.y+float64(c.height) > c.bounds.Y+c.bounds.Height {
		c.y = c.bounds.Y + c.bounds.Height - float64(c.height)
	}
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *Camera) WorldToScreen(worldX, worldY float64) (int, int) {
	screenX := int(math.Round(worldX - c.x))
	screenY := int(math.Round(worldY - c.y))
	return screenX, screenY
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *Camera) ScreenToWorld(screenX, screenY int) (float64, float64) {
	worldX := float64(screenX) + c.x
	worldY := float64(screenY) + c.y
	return worldX, worldY
}

// IsVisible checks if an entity at the given position and size is visible in the camera viewport
func (c *Camera) IsVisible(x, y, width, height float64) bool {
	return x+width >= c.x &&
		x <= c.x+float64(c.width) &&
		y+height >= c.y &&
		y <= c.y+float64(c.height)
}

// Shake applies a screen shake effect
func (c *Camera) Shake(amount, duration float64) {
	c.shakeAmount = amount
	c.shakeTime = duration
}

// GetPosition returns the current camera position
func (c *Camera) GetPosition() (float64, float64) {
	return c.x, c.y
}

// GetSize returns the camera viewport size
func (c *Camera) GetSize() (int, int) {
	return c.width, c.height
}

// SetMode sets the camera behavior mode
func (c *Camera) SetMode(mode CameraMode) {
	c.mode = mode
}

// GetMode returns the current camera mode
func (c *Camera) GetMode() CameraMode {
	return c.mode
}

// SetSmoothing sets the camera smoothing factor (0.0 = no smoothing, 1.0 = instant)
func (c *Camera) SetSmoothing(smoothing float64) {
	if smoothing < 0 {
		smoothing = 0
	}
	if smoothing > 1 {
		smoothing = 1
	}
	c.smoothing = smoothing
}

// GetSmoothing returns the current smoothing factor
func (c *Camera) GetSmoothing() float64 {
	return c.smoothing
}

// SetLookAhead sets the look-ahead distance
func (c *Camera) SetLookAhead(distance float64) {
	c.lookAhead = distance
}

// GetLookAhead returns the current look-ahead distance
func (c *Camera) GetLookAhead() float64 {
	return c.lookAhead
}

// SetBounds updates the camera bounds
func (c *Camera) SetBounds(bounds Rect) {
	c.bounds = bounds
}

// GetBounds returns the current camera bounds
func (c *Camera) GetBounds() Rect {
	return c.bounds
}

// SetPosition directly sets the camera position (useful for initialization)
func (c *Camera) SetPosition(x, y float64) {
	c.x = x
	c.y = y
	c.targetX = x
	c.targetY = y
}
