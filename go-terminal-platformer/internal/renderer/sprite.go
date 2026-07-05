package renderer

import (
	"github.com/gdamore/tcell/v2"
)

// SpriteCell represents a single cell in a sprite
type SpriteCell struct {
	Char rune
	FG   tcell.Color
	BG   tcell.Color
}

// Sprite represents a 2D sprite made of characters
type Sprite struct {
	Width  int
	Height int
	Data   [][]SpriteCell
}

// NewSprite creates a new sprite with the given dimensions
func NewSprite(width, height int) *Sprite {
	data := make([][]SpriteCell, height)
	for i := range data {
		data[i] = make([]SpriteCell, width)
	}
	
	return &Sprite{
		Width:  width,
		Height: height,
		Data:   data,
	}
}

// NewSpriteFromString creates a sprite from a string representation
// Each line represents a row, each character a cell
func NewSpriteFromString(lines []string, fg, bg tcell.Color) *Sprite {
	if len(lines) == 0 {
		return NewSprite(0, 0)
	}
	
	height := len(lines)
	width := 0
	for _, line := range lines {
		if len(line) > width {
			width = len(line)
		}
	}
	
	sprite := NewSprite(width, height)
	
	for y, line := range lines {
		for x, ch := range line {
			sprite.Data[y][x] = SpriteCell{
				Char: ch,
				FG:   fg,
				BG:   bg,
			}
		}
	}
	
	return sprite
}

// SetCell sets a cell in the sprite
func (s *Sprite) SetCell(x, y int, char rune, fg, bg tcell.Color) {
	if x >= 0 && x < s.Width && y >= 0 && y < s.Height {
		s.Data[y][x] = SpriteCell{
			Char: char,
			FG:   fg,
			BG:   bg,
		}
	}
}

// GetCell gets a cell from the sprite
func (s *Sprite) GetCell(x, y int) SpriteCell {
	if x >= 0 && x < s.Width && y >= 0 && y < s.Height {
		return s.Data[y][x]
	}
	return SpriteCell{}
}

// Clone creates a copy of the sprite
func (s *Sprite) Clone() *Sprite {
	clone := NewSprite(s.Width, s.Height)
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			clone.Data[y][x] = s.Data[y][x]
		}
	}
	return clone
}

// FlipHorizontal flips the sprite horizontally
func (s *Sprite) FlipHorizontal() *Sprite {
	flipped := NewSprite(s.Width, s.Height)
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			flipped.Data[y][s.Width-1-x] = s.Data[y][x]
		}
	}
	return flipped
}

// FlipVertical flips the sprite vertically
func (s *Sprite) FlipVertical() *Sprite {
	flipped := NewSprite(s.Width, s.Height)
	for y := 0; y < s.Height; y++ {
		for x := 0; x < s.Width; x++ {
			flipped.Data[s.Height-1-y][x] = s.Data[y][x]
		}
	}
	return flipped
}

// Common sprite definitions for the game

// PlayerSprite creates a basic player sprite
func PlayerSprite() *Sprite {
	lines := []string{
		" O ",
		"/|\\",
		"/ \\",
	}
	return NewSpriteFromString(lines, tcell.ColorRed, tcell.ColorDefault)
}

// GoombaSprite creates a basic Goomba enemy sprite
func GoombaSprite() *Sprite {
	lines := []string{
		" oo ",
		"(  )",
		" MM ",
	}
	return NewSpriteFromString(lines, tcell.ColorBrown, tcell.ColorDefault)
}

// CoinSprite creates a coin sprite
func CoinSprite() *Sprite {
	lines := []string{
		"($)",
	}
	return NewSpriteFromString(lines, tcell.ColorYellow, tcell.ColorDefault)
}

// BlockSprite creates a block sprite
func BlockSprite() *Sprite {
	lines := []string{
		"[?]",
	}
	return NewSpriteFromString(lines, tcell.ColorOrange, tcell.ColorDefault)
}

// BrickSprite creates a brick sprite
func BrickSprite() *Sprite {
	lines := []string{
		"###",
	}
	return NewSpriteFromString(lines, tcell.ColorMaroon, tcell.ColorDefault)
}

// PipeSprite creates a pipe sprite (top part)
func PipeSprite() *Sprite {
	lines := []string{
		"╔═╗",
		"║ ║",
		"║ ║",
	}
	return NewSpriteFromString(lines, tcell.ColorGreen, tcell.ColorDefault)
}

// MushroomSprite creates a mushroom power-up sprite
func MushroomSprite() *Sprite {
	lines := []string{
		" ● ",
		"(_)",
	}
	return NewSpriteFromString(lines, tcell.ColorRed, tcell.ColorDefault)
}

// FireFlowerSprite creates a fire flower power-up sprite
func FireFlowerSprite() *Sprite {
	lines := []string{
		"*o*",
		" | ",
	}
	return NewSpriteFromString(lines, tcell.ColorOrange, tcell.ColorDefault)
}

// StarSprite creates a star power-up sprite
func StarSprite() *Sprite {
	lines := []string{
		" * ",
		"***",
		" * ",
	}
	return NewSpriteFromString(lines, tcell.ColorYellow, tcell.ColorDefault)
}

// CloudSprite creates a cloud sprite
func CloudSprite() *Sprite {
	lines := []string{
		" ≈≈≈ ",
		"≈≈≈≈≈",
	}
	return NewSpriteFromString(lines, tcell.ColorWhite, tcell.ColorDefault)
}

// BushSprite creates a bush sprite
func BushSprite() *Sprite {
	lines := []string{
		"∩∩∩",
	}
	return NewSpriteFromString(lines, tcell.ColorGreen, tcell.ColorDefault)
}

// FlagSprite creates a flag sprite
func FlagSprite() *Sprite {
	lines := []string{
		"▶",
		"|",
		"|",
		"|",
	}
	return NewSpriteFromString(lines, tcell.ColorRed, tcell.ColorDefault)
}
