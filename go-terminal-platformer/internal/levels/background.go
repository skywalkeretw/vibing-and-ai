package levels

import (
	"github.com/gdamore/tcell/v2"
)

// BackgroundLayer represents a single parallax layer
type BackgroundLayer struct {
	Sprite      [][]rune
	ScrollSpeed float64
	Offset      float64
	Color       tcell.Color
	Width       int
	Height      int
}

// Background represents a multi-layer parallax background
type Background struct {
	Layers []BackgroundLayer
	Name   string
}

// NewBackground creates a new background with the given layers
func NewBackground(name string) *Background {
	return &Background{
		Name:   name,
		Layers: make([]BackgroundLayer, 0),
	}
}

// AddLayer adds a parallax layer to the background
func (b *Background) AddLayer(sprite [][]rune, scrollSpeed float64, color tcell.Color) {
	height := len(sprite)
	width := 0
	if height > 0 {
		width = len(sprite[0])
	}

	layer := BackgroundLayer{
		Sprite:      sprite,
		ScrollSpeed: scrollSpeed,
		Offset:      0,
		Color:       color,
		Width:       width,
		Height:      height,
	}
	b.Layers = append(b.Layers, layer)
}

// Update updates all background layers
func (b *Background) Update(deltaTime float64) {
	for i := range b.Layers {
		b.Layers[i].Offset += b.Layers[i].ScrollSpeed * deltaTime
		
		// Wrap offset for seamless scrolling
		if b.Layers[i].Width > 0 {
			layerWidth := float64(b.Layers[i].Width)
			for b.Layers[i].Offset >= layerWidth {
				b.Layers[i].Offset -= layerWidth
			}
			for b.Layers[i].Offset < 0 {
				b.Layers[i].Offset += layerWidth
			}
		}
	}
}

// Render renders all background layers
func (b *Background) Render(renderer interface{}) {
	// TODO: Implement rendering when renderer interface is available
	// For each layer:
	//   1. Calculate the offset position
	//   2. Render the sprite at the offset
	//   3. Render a second copy for seamless wrapping
	_ = renderer
}

// LoadBackground loads a predefined background by name
func LoadBackground(name string) *Background {
	bg := NewBackground(name)

	switch name {
	case "grassland":
		bg.loadGrasslandBackground()
	case "desert":
		bg.loadDesertBackground()
	case "ice":
		bg.loadIceBackground()
	case "volcano":
		bg.loadVolcanoBackground()
	default:
		bg.loadDefaultBackground()
	}

	return bg
}

// loadGrasslandBackground creates a grassland-themed background
func (b *Background) loadGrasslandBackground() {
	// Far layer (clouds) - slowest
	farLayer := [][]rune{
		[]rune("    ☁️      ☁️         ☁️    "),
		[]rune("  ☁️    ☁️       ☁️         "),
	}
	b.AddLayer(farLayer, 0.2, tcell.ColorWhite)

	// Mid layer (hills) - medium speed
	midLayer := [][]rune{
		[]rune("      /\\        /\\      "),
		[]rune("     /  \\      /  \\     "),
		[]rune("    /    \\    /    \\    "),
	}
	b.AddLayer(midLayer, 0.5, tcell.ColorGreen)

	// Near layer (trees) - fastest
	nearLayer := [][]rune{
		[]rune("  ♣   ♣    ♣   ♣  "),
		[]rune("  ║   ║    ║   ║  "),
	}
	b.AddLayer(nearLayer, 0.8, tcell.ColorDarkGreen)
}

// loadDesertBackground creates a desert-themed background
func (b *Background) loadDesertBackground() {
	// Far layer (sun/sky) - slowest
	farLayer := [][]rune{
		[]rune("        ☀️          "),
		[]rune("                   "),
	}
	b.AddLayer(farLayer, 0.1, tcell.ColorYellow)

	// Mid layer (dunes) - medium speed
	midLayer := [][]rune{
		[]rune("    ∿∿∿    ∿∿∿    "),
		[]rune("  ∿∿   ∿∿∿    ∿∿  "),
	}
	b.AddLayer(midLayer, 0.4, tcell.ColorOrange)

	// Near layer (cacti) - fastest
	nearLayer := [][]rune{
		[]rune("  🌵    🌵    🌵  "),
	}
	b.AddLayer(nearLayer, 0.7, tcell.ColorGreen)
}

// loadIceBackground creates an ice-themed background
func (b *Background) loadIceBackground() {
	// Far layer (snow) - slowest
	farLayer := [][]rune{
		[]rune("  ❄️   ❄️    ❄️   ❄️  "),
		[]rune("    ❄️    ❄️    ❄️    "),
	}
	b.AddLayer(farLayer, 0.3, tcell.ColorWhite)

	// Mid layer (ice formations) - medium speed
	midLayer := [][]rune{
		[]rune("   /\\  /\\   /\\   "),
		[]rune("  /  \\/  \\ /  \\  "),
	}
	b.AddLayer(midLayer, 0.5, tcell.ColorLightCyan)

	// Near layer (icicles) - fastest
	nearLayer := [][]rune{
		[]rune(" | | |  | | |  | "),
	}
	b.AddLayer(nearLayer, 0.8, tcell.ColorLightCyan)
}

// loadVolcanoBackground creates a volcano-themed background
func (b *Background) loadVolcanoBackground() {
	// Far layer (smoke) - slowest
	farLayer := [][]rune{
		[]rune("  ☁️  ☁️   ☁️  ☁️  "),
		[]rune("    ☁️   ☁️   ☁️   "),
	}
	b.AddLayer(farLayer, 0.2, tcell.ColorGray)

	// Mid layer (lava flows) - medium speed
	midLayer := [][]rune{
		[]rune("  ≈≈≈  ≈≈≈  ≈≈≈  "),
		[]rune(" ≈≈≈≈≈≈≈≈≈≈≈≈≈≈ "),
	}
	b.AddLayer(midLayer, 0.6, tcell.ColorRed)

	// Near layer (rocks) - fastest
	nearLayer := [][]rune{
		[]rune(" ▓ ▓  ▓ ▓  ▓ ▓ "),
	}
	b.AddLayer(nearLayer, 0.9, tcell.ColorDarkRed)
}

// loadDefaultBackground creates a simple default background
func (b *Background) loadDefaultBackground() {
	// Single layer with simple pattern
	defaultLayer := [][]rune{
		[]rune("                    "),
		[]rune("                    "),
	}
	b.AddLayer(defaultLayer, 0.3, tcell.ColorDefault)
}

// SetScrollSpeed sets the scroll speed for a specific layer
func (b *Background) SetScrollSpeed(layerIndex int, speed float64) {
	if layerIndex >= 0 && layerIndex < len(b.Layers) {
		b.Layers[layerIndex].ScrollSpeed = speed
	}
}

// GetLayerCount returns the number of layers in the background
func (b *Background) GetLayerCount() int {
	return len(b.Layers)
}

// Reset resets all layer offsets to zero
func (b *Background) Reset() {
	for i := range b.Layers {
		b.Layers[i].Offset = 0
	}
}
