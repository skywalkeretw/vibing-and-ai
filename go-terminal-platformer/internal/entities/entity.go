package entities

import (
	"github.com/lukeroy/go-terminal-platformer/internal/engine"
)

// Entity represents any game object that can be updated and rendered
type Entity interface {
	// Core lifecycle methods
	Update(deltaTime float64)
	Render(renderer Renderer)
	
	// Position management
	SetPosition(pos engine.Vector2)
	GetPosition() engine.Vector2
	
	// State management
	IsActive() bool
	SetActive(active bool)
	
	// Collision handling
	OnCollision(other Entity)
}

// Renderer interface for rendering entities
type Renderer interface {
	DrawSprite(x, y int, sprite *Sprite)
	DrawText(x, y int, text string, color int)
}

// Sprite represents a visual representation of an entity
type Sprite struct {
	Width  int
	Height int
	Data   [][]rune
	Color  int
}

// NewSprite creates a new sprite
func NewSprite(width, height int) *Sprite {
	data := make([][]rune, height)
	for i := range data {
		data[i] = make([]rune, width)
	}
	return &Sprite{
		Width:  width,
		Height: height,
		Data:   data,
	}
}

// FlipHorizontal returns a horizontally flipped copy of the sprite
func (s *Sprite) FlipHorizontal() *Sprite {
	flipped := NewSprite(s.Width, s.Height)
	flipped.Color = s.Color
	
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			flipped.Data[y][x] = s.Data[y][s.Width-1-x]
		}
	}
	
	return flipped
}

// PowerUpType represents different power-up types
type PowerUpType int

const (
	PowerUpNone PowerUpType = iota
	PowerUpFire
	PowerUpStar
	PowerUpMushroom
	PowerUpSpeedBoots
	PowerUpSuperJump
	PowerUpShield
)

// String returns the string representation of a power-up type
func (p PowerUpType) String() string {
	switch p {
	case PowerUpNone:
		return "None"
	case PowerUpFire:
		return "Fire"
	case PowerUpStar:
		return "Star"
	case PowerUpSpeedBoots:
		return "SpeedBoots"
	case PowerUpSuperJump:
		return "SuperJump"
	case PowerUpShield:
		return "Shield"
	default:
		return "Unknown"
	}
}
