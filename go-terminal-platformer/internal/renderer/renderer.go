package renderer

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// RenderLayer represents different rendering layers
type RenderLayer int

const (
	LayerBackground RenderLayer = iota
	LayerTerrain
	LayerEntities
	LayerEffects
	LayerUI
	LayerDebug
)

// Renderer handles all terminal rendering operations
type Renderer struct {
	screen    tcell.Screen
	width     int
	height    int
	colorMode bool
	camera    *Camera
}

// New creates a new Renderer instance
func New(screen tcell.Screen) *Renderer {
	width, height := screen.Size()
	
	return &Renderer{
		screen:    screen,
		width:     width,
		height:    height,
		colorMode: true, // Assume color support by default
		camera:    NewCamera(width, height),
	}
}

// Initialize sets up the renderer
func (r *Renderer) Initialize() error {
	if r.screen == nil {
		return fmt.Errorf("screen is nil")
	}

	// Get screen dimensions
	r.width, r.height = r.screen.Size()

	// Detect color support
	r.colorMode = r.detectColorSupport()

	// Initialize camera with screen dimensions
	r.camera = NewCamera(r.width, r.height)

	return nil
}

// detectColorSupport checks if the terminal supports colors
func (r *Renderer) detectColorSupport() bool {
	// tcell handles color detection internally
	// We'll assume color support is available
	return true
}

// Clear clears the entire screen
func (r *Renderer) Clear() {
	r.screen.Clear()
}

// DrawChar draws a single character at the specified world coordinates
func (r *Renderer) DrawChar(worldX, worldY int, char rune, fg, bg tcell.Color) {
	// Apply camera offset to convert world to screen coordinates
	screenX, screenY := r.camera.WorldToScreen(worldX, worldY)

	// Check if within screen bounds
	if screenX >= 0 && screenX < r.width && screenY >= 0 && screenY < r.height {
		style := tcell.StyleDefault.Foreground(fg).Background(bg)
		r.screen.SetContent(screenX, screenY, char, nil, style)
	}
}

// DrawCharScreen draws a single character at screen coordinates (no camera transform)
func (r *Renderer) DrawCharScreen(screenX, screenY int, char rune, fg, bg tcell.Color) {
	if screenX >= 0 && screenX < r.width && screenY >= 0 && screenY < r.height {
		style := tcell.StyleDefault.Foreground(fg).Background(bg)
		r.screen.SetContent(screenX, screenY, char, nil, style)
	}
}

// DrawString draws a string at the specified world coordinates
func (r *Renderer) DrawString(worldX, worldY int, text string, fg, bg tcell.Color) {
	for i, ch := range text {
		r.DrawChar(worldX+i, worldY, ch, fg, bg)
	}
}

// DrawStringScreen draws a string at screen coordinates (no camera transform)
func (r *Renderer) DrawStringScreen(screenX, screenY int, text string, fg, bg tcell.Color) {
	for i, ch := range text {
		r.DrawCharScreen(screenX+i, screenY, ch, fg, bg)
	}
}

// DrawStringCentered draws a centered string at the specified screen Y coordinate
func (r *Renderer) DrawStringCentered(screenY int, text string, fg, bg tcell.Color) {
	screenX := (r.width - len(text)) / 2
	r.DrawStringScreen(screenX, screenY, text, fg, bg)
}

// DrawBox draws a box with borders at screen coordinates
func (r *Renderer) DrawBox(x, y, width, height int, fg, bg tcell.Color, filled bool) {
	// Draw corners
	r.DrawCharScreen(x, y, '┌', fg, bg)
	r.DrawCharScreen(x+width-1, y, '┐', fg, bg)
	r.DrawCharScreen(x, y+height-1, '└', fg, bg)
	r.DrawCharScreen(x+width-1, y+height-1, '┘', fg, bg)

	// Draw horizontal borders
	for i := 1; i < width-1; i++ {
		r.DrawCharScreen(x+i, y, '─', fg, bg)
		r.DrawCharScreen(x+i, y+height-1, '─', fg, bg)
	}

	// Draw vertical borders
	for i := 1; i < height-1; i++ {
		r.DrawCharScreen(x, y+i, '│', fg, bg)
		r.DrawCharScreen(x+width-1, y+i, '│', fg, bg)
	}

	// Fill interior if requested
	if filled {
		for dy := 1; dy < height-1; dy++ {
			for dx := 1; dx < width-1; dx++ {
				r.DrawCharScreen(x+dx, y+dy, ' ', fg, bg)
			}
		}
	}
}

// DrawRect draws a filled rectangle at world coordinates
func (r *Renderer) DrawRect(worldX, worldY, width, height int, char rune, fg, bg tcell.Color) {
	for dy := 0; dy < height; dy++ {
		for dx := 0; dx < width; dx++ {
			r.DrawChar(worldX+dx, worldY+dy, char, fg, bg)
		}
	}
}

// DrawSprite draws a sprite at world coordinates
func (r *Renderer) DrawSprite(worldX, worldY int, sprite *Sprite) {
	if sprite == nil {
		return
	}

	for dy := 0; dy < sprite.Height; dy++ {
		for dx := 0; dx < sprite.Width; dx++ {
			if dy < len(sprite.Data) && dx < len(sprite.Data[dy]) {
				cell := sprite.Data[dy][dx]
				if cell.Char != 0 { // Skip empty cells
					r.DrawChar(worldX+dx, worldY+dy, cell.Char, cell.FG, cell.BG)
				}
			}
		}
	}
}

// DrawSpriteScreen draws a sprite at screen coordinates (no camera transform)
func (r *Renderer) DrawSpriteScreen(screenX, screenY int, sprite *Sprite) {
	if sprite == nil {
		return
	}

	for dy := 0; dy < sprite.Height; dy++ {
		for dx := 0; dx < sprite.Width; dx++ {
			if dy < len(sprite.Data) && dx < len(sprite.Data[dy]) {
				cell := sprite.Data[dy][dx]
				if cell.Char != 0 { // Skip empty cells
					r.DrawCharScreen(screenX+dx, screenY+dy, cell.Char, cell.FG, cell.BG)
				}
			}
		}
	}
}

// Show displays the rendered content (flips buffers)
func (r *Renderer) Show() {
	r.screen.Show()
}

// HandleResize handles terminal resize events
func (r *Renderer) HandleResize() {
	r.width, r.height = r.screen.Size()
	
	// Update camera dimensions
	if r.camera != nil {
		r.camera.width = r.width
		r.camera.height = r.height
	}
	
	// Sync screen
	r.screen.Sync()
}

// GetSize returns the current screen dimensions
func (r *Renderer) GetSize() (int, int) {
	return r.width, r.height
}

// GetCamera returns the renderer's camera
func (r *Renderer) GetCamera() *Camera {
	return r.camera
}

// SetCamera sets the renderer's camera
func (r *Renderer) SetCamera(camera *Camera) {
	r.camera = camera
}

// IsColorMode returns whether color mode is enabled
func (r *Renderer) IsColorMode() bool {
	return r.colorMode
}

// SetColorMode sets the color mode
func (r *Renderer) SetColorMode(enabled bool) {
	r.colorMode = enabled
}

// DrawLine draws a line between two points (Bresenham's algorithm)
func (r *Renderer) DrawLine(x0, y0, x1, y1 int, char rune, fg, bg tcell.Color) {
	dx := abs(x1 - x0)
	dy := abs(y1 - y0)
	sx := 1
	if x0 > x1 {
		sx = -1
	}
	sy := 1
	if y0 > y1 {
		sy = -1
	}
	err := dx - dy

	for {
		r.DrawChar(x0, y0, char, fg, bg)
		
		if x0 == x1 && y0 == y1 {
			break
		}
		
		e2 := 2 * err
		if e2 > -dy {
			err -= dy
			x0 += sx
		}
		if e2 < dx {
			err += dx
			y0 += sy
		}
	}
}

// DrawHorizontalLine draws a horizontal line at world coordinates
func (r *Renderer) DrawHorizontalLine(worldX, worldY, length int, char rune, fg, bg tcell.Color) {
	for i := 0; i < length; i++ {
		r.DrawChar(worldX+i, worldY, char, fg, bg)
	}
}

// DrawVerticalLine draws a vertical line at world coordinates
func (r *Renderer) DrawVerticalLine(worldX, worldY, length int, char rune, fg, bg tcell.Color) {
	for i := 0; i < length; i++ {
		r.DrawChar(worldX, worldY+i, char, fg, bg)
	}
}

// FillScreen fills the entire screen with a character
func (r *Renderer) FillScreen(char rune, fg, bg tcell.Color) {
	style := tcell.StyleDefault.Foreground(fg).Background(bg)
	for y := 0; y < r.height; y++ {
		for x := 0; x < r.width; x++ {
			r.screen.SetContent(x, y, char, nil, style)
		}
	}
}

// Helper function
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
