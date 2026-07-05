package renderer

import (
	"math"
)

// Camera represents the viewport into the game world
type Camera struct {
	x      int // World X position of camera
	y      int // World Y position of camera
	width  int // Camera viewport width
	height int // Camera viewport height

	// Camera behavior
	followTarget interface{} // Entity to follow (will be *entities.Player later)
	smoothing    float64     // Camera smoothing factor (0.0 = instant, 1.0 = very smooth)
	
	// Bounds
	bounds Rectangle
	
	// Offset for look-ahead
	offsetX int
	offsetY int
}

// Rectangle represents a rectangular area
type Rectangle struct {
	X      int
	Y      int
	Width  int
	Height int
}

// NewCamera creates a new camera with the given viewport dimensions
func NewCamera(width, height int) *Camera {
	return &Camera{
		x:         0,
		y:         0,
		width:     width,
		height:    height,
		smoothing: 0.1, // Default smoothing
		bounds: Rectangle{
			X:      0,
			Y:      0,
			Width:  1000, // Default large bounds
			Height: 1000,
		},
	}
}

// Update updates the camera position
func (c *Camera) Update(deltaTime float64) {
	// Camera update logic will be implemented when we have entities
	// For now, this is a placeholder
}

// SetPosition sets the camera position directly
func (c *Camera) SetPosition(x, y int) {
	c.x = x
	c.y = y
	c.clampToBounds()
}

// Move moves the camera by the given delta
func (c *Camera) Move(dx, dy int) {
	c.x += dx
	c.y += dy
	c.clampToBounds()
}

// CenterOn centers the camera on the given world coordinates
func (c *Camera) CenterOn(worldX, worldY int) {
	c.x = worldX - c.width/2
	c.y = worldY - c.height/2
	c.clampToBounds()
}

// WorldToScreen converts world coordinates to screen coordinates
func (c *Camera) WorldToScreen(worldX, worldY int) (int, int) {
	screenX := worldX - c.x
	screenY := worldY - c.y
	return screenX, screenY
}

// ScreenToWorld converts screen coordinates to world coordinates
func (c *Camera) ScreenToWorld(screenX, screenY int) (int, int) {
	worldX := screenX + c.x
	worldY := screenY + c.y
	return worldX, worldY
}

// IsVisible checks if a rectangle in world coordinates is visible
func (c *Camera) IsVisible(worldX, worldY, width, height int) bool {
	// Check if rectangle overlaps with camera viewport
	return worldX+width >= c.x &&
		worldX <= c.x+c.width &&
		worldY+height >= c.y &&
		worldY <= c.y+c.height
}

// IsPointVisible checks if a point in world coordinates is visible
func (c *Camera) IsPointVisible(worldX, worldY int) bool {
	return worldX >= c.x &&
		worldX < c.x+c.width &&
		worldY >= c.y &&
		worldY < c.y+c.height
}

// SetBounds sets the camera movement bounds
func (c *Camera) SetBounds(x, y, width, height int) {
	c.bounds = Rectangle{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
	c.clampToBounds()
}

// GetBounds returns the camera bounds
func (c *Camera) GetBounds() Rectangle {
	return c.bounds
}

// clampToBounds ensures the camera stays within bounds
func (c *Camera) clampToBounds() {
	// Don't scroll past left edge
	if c.x < c.bounds.X {
		c.x = c.bounds.X
	}
	
	// Don't scroll past top edge
	if c.y < c.bounds.Y {
		c.y = c.bounds.Y
	}
	
	// Don't scroll past right edge
	if c.x+c.width > c.bounds.X+c.bounds.Width {
		c.x = c.bounds.X + c.bounds.Width - c.width
		if c.x < c.bounds.X {
			c.x = c.bounds.X
		}
	}
	
	// Don't scroll past bottom edge
	if c.y+c.height > c.bounds.Y+c.bounds.Height {
		c.y = c.bounds.Y + c.bounds.Height - c.height
		if c.y < c.bounds.Y {
			c.y = c.bounds.Y
		}
	}
}

// GetPosition returns the camera's world position
func (c *Camera) GetPosition() (int, int) {
	return c.x, c.y
}

// GetSize returns the camera's viewport size
func (c *Camera) GetSize() (int, int) {
	return c.width, c.height
}

// SetSize sets the camera's viewport size
func (c *Camera) SetSize(width, height int) {
	c.width = width
	c.height = height
	c.clampToBounds()
}

// SetSmoothing sets the camera smoothing factor
func (c *Camera) SetSmoothing(smoothing float64) {
	c.smoothing = clamp(smoothing, 0.0, 1.0)
}

// GetSmoothing returns the camera smoothing factor
func (c *Camera) GetSmoothing() float64 {
	return c.smoothing
}

// SetOffset sets the camera offset (for look-ahead)
func (c *Camera) SetOffset(offsetX, offsetY int) {
	c.offsetX = offsetX
	c.offsetY = offsetY
}

// GetOffset returns the camera offset
func (c *Camera) GetOffset() (int, int) {
	return c.offsetX, c.offsetY
}

// LerpToPosition smoothly moves the camera towards a target position
func (c *Camera) LerpToPosition(targetX, targetY int, deltaTime float64) {
	// Calculate the difference
	dx := float64(targetX - c.x)
	dy := float64(targetY - c.y)
	
	// Apply smoothing
	moveX := dx * c.smoothing
	moveY := dy * c.smoothing
	
	// Update position
	c.x += int(math.Round(moveX))
	c.y += int(math.Round(moveY))
	
	c.clampToBounds()
}

// Shake applies a screen shake effect (to be called each frame during shake)
func (c *Camera) Shake(intensity int) (int, int) {
	// Simple shake: random offset within intensity range
	// This would need a random number generator in practice
	// For now, return zero offset
	return 0, 0
}

// GetViewportRect returns the camera's viewport as a rectangle in world coordinates
func (c *Camera) GetViewportRect() Rectangle {
	return Rectangle{
		X:      c.x,
		Y:      c.y,
		Width:  c.width,
		Height: c.height,
	}
}

// Helper function to clamp a float64 value
func clamp(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Helper function to clamp an int value
func clampInt(value, min, max int) int {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}
