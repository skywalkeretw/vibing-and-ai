package renderer

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewCamera(t *testing.T) {
	camera := NewCamera(80, 24)
	
	if camera == nil {
		t.Fatal("Camera is nil")
	}
	
	if camera.width != 80 {
		t.Errorf("Expected width 80, got %d", camera.width)
	}
	
	if camera.height != 24 {
		t.Errorf("Expected height 24, got %d", camera.height)
	}
	
	if camera.x != 0 || camera.y != 0 {
		t.Errorf("Expected initial position (0,0), got (%d,%d)", camera.x, camera.y)
	}
}

func TestCameraSetPosition(t *testing.T) {
	camera := NewCamera(80, 24)
	camera.SetBounds(0, 0, 1000, 1000)
	
	camera.SetPosition(100, 50)
	
	x, y := camera.GetPosition()
	if x != 100 || y != 50 {
		t.Errorf("Expected position (100,50), got (%d,%d)", x, y)
	}
}

func TestCameraMove(t *testing.T) {
	camera := NewCamera(80, 24)
	camera.SetBounds(0, 0, 1000, 1000)
	camera.SetPosition(100, 100)
	
	camera.Move(10, -5)
	
	x, y := camera.GetPosition()
	if x != 110 || y != 95 {
		t.Errorf("Expected position (110,95), got (%d,%d)", x, y)
	}
}

func TestCameraCenterOn(t *testing.T) {
	camera := NewCamera(80, 24)
	camera.SetBounds(0, 0, 1000, 1000)
	
	camera.CenterOn(100, 100)
	
	x, y := camera.GetPosition()
	expectedX := 100 - 80/2 // 60
	expectedY := 100 - 24/2 // 88
	
	if x != expectedX || y != expectedY {
		t.Errorf("Expected position (%d,%d), got (%d,%d)", expectedX, expectedY, x, y)
	}
}

func TestCameraWorldToScreen(t *testing.T) {
	camera := NewCamera(80, 24)
	camera.SetPosition(100, 50)
	
	screenX, screenY := camera.WorldToScreen(150, 75)
	
	if screenX != 50 || screenY != 25 {
		t.Errorf("Expected screen coords (50,25), got (%d,%d)", screenX, screenY)
	}
}

func TestCameraScreenToWorld(t *testing.T) {
	camera := NewCamera(80, 24)
	camera.SetPosition(100, 50)
	
	worldX, worldY := camera.ScreenToWorld(50, 25)
	
	if worldX != 150 || worldY != 75 {
		t.Errorf("Expected world coords (150,75), got (%d,%d)", worldX, worldY)
	}
}

func TestCameraIsVisible(t *testing.T) {
	camera := NewCamera(80, 24)
	camera.SetPosition(0, 0)
	
	tests := []struct {
		name     string
		x, y, w, h int
		expected bool
	}{
		{"Inside viewport", 10, 10, 5, 5, true},
		{"Partially visible", 75, 20, 10, 10, true},
		{"Outside right", 100, 10, 5, 5, false},
		{"Outside bottom", 10, 30, 5, 5, false},
		{"Outside left", -10, 10, 5, 5, false},
		{"Outside top", 10, -10, 5, 5, false},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := camera.IsVisible(tt.x, tt.y, tt.w, tt.h)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v for (%d,%d,%d,%d)", 
					tt.expected, result, tt.x, tt.y, tt.w, tt.h)
			}
		})
	}
}

func TestCameraBounds(t *testing.T) {
	camera := NewCamera(80, 24)
	camera.SetBounds(0, 0, 200, 100)
	
	// Try to move beyond right bound
	camera.SetPosition(200, 0)
	x, y := camera.GetPosition()
	
	// Camera should be clamped so that right edge doesn't exceed bounds
	expectedX := 200 - 80 // 120
	if x != expectedX {
		t.Errorf("Expected x to be clamped to %d, got %d", expectedX, x)
	}
	
	// Try to move beyond bottom bound
	camera.SetPosition(0, 100)
	x, y = camera.GetPosition()
	
	expectedY := 100 - 24 // 76
	if y != expectedY {
		t.Errorf("Expected y to be clamped to %d, got %d", expectedY, y)
	}
	
	// Try to move beyond left bound
	camera.SetPosition(-10, 0)
	x, y = camera.GetPosition()
	
	if x != 0 {
		t.Errorf("Expected x to be clamped to 0, got %d", x)
	}
	
	// Try to move beyond top bound
	camera.SetPosition(0, -10)
	x, y = camera.GetPosition()
	
	if y != 0 {
		t.Errorf("Expected y to be clamped to 0, got %d", y)
	}
}

func TestNewSprite(t *testing.T) {
	sprite := NewSprite(5, 3)
	
	if sprite == nil {
		t.Fatal("Sprite is nil")
	}
	
	if sprite.Width != 5 {
		t.Errorf("Expected width 5, got %d", sprite.Width)
	}
	
	if sprite.Height != 3 {
		t.Errorf("Expected height 3, got %d", sprite.Height)
	}
	
	if len(sprite.Data) != 3 {
		t.Errorf("Expected 3 rows, got %d", len(sprite.Data))
	}
	
	if len(sprite.Data[0]) != 5 {
		t.Errorf("Expected 5 columns, got %d", len(sprite.Data[0]))
	}
}

func TestNewSpriteFromString(t *testing.T) {
	lines := []string{
		"ABC",
		"DEF",
	}
	
	sprite := NewSpriteFromString(lines, tcell.ColorWhite, tcell.ColorBlack)
	
	if sprite.Width != 3 {
		t.Errorf("Expected width 3, got %d", sprite.Width)
	}
	
	if sprite.Height != 2 {
		t.Errorf("Expected height 2, got %d", sprite.Height)
	}
	
	if sprite.Data[0][0].Char != 'A' {
		t.Errorf("Expected 'A', got '%c'", sprite.Data[0][0].Char)
	}
	
	if sprite.Data[1][2].Char != 'F' {
		t.Errorf("Expected 'F', got '%c'", sprite.Data[1][2].Char)
	}
}

func TestSpriteSetGetCell(t *testing.T) {
	sprite := NewSprite(3, 3)
	
	sprite.SetCell(1, 1, 'X', tcell.ColorRed, tcell.ColorBlue)
	
	cell := sprite.GetCell(1, 1)
	
	if cell.Char != 'X' {
		t.Errorf("Expected 'X', got '%c'", cell.Char)
	}
	
	if cell.FG != tcell.ColorRed {
		t.Errorf("Expected red foreground")
	}
	
	if cell.BG != tcell.ColorBlue {
		t.Errorf("Expected blue background")
	}
}

func TestSpriteClone(t *testing.T) {
	original := NewSprite(2, 2)
	original.SetCell(0, 0, 'A', tcell.ColorWhite, tcell.ColorBlack)
	original.SetCell(1, 1, 'B', tcell.ColorRed, tcell.ColorBlue)
	
	clone := original.Clone()
	
	if clone.Width != original.Width {
		t.Error("Clone width doesn't match")
	}
	
	if clone.Height != original.Height {
		t.Error("Clone height doesn't match")
	}
	
	if clone.GetCell(0, 0).Char != 'A' {
		t.Error("Clone cell (0,0) doesn't match")
	}
	
	if clone.GetCell(1, 1).Char != 'B' {
		t.Error("Clone cell (1,1) doesn't match")
	}
	
	// Modify clone and ensure original is unchanged
	clone.SetCell(0, 0, 'Z', tcell.ColorYellow, tcell.ColorGreen)
	
	if original.GetCell(0, 0).Char == 'Z' {
		t.Error("Modifying clone affected original")
	}
}

func TestSpriteFlipHorizontal(t *testing.T) {
	lines := []string{
		"ABC",
		"DEF",
	}
	
	sprite := NewSpriteFromString(lines, tcell.ColorWhite, tcell.ColorBlack)
	flipped := sprite.FlipHorizontal()
	
	if flipped.GetCell(0, 0).Char != 'C' {
		t.Errorf("Expected 'C' at (0,0), got '%c'", flipped.GetCell(0, 0).Char)
	}
	
	if flipped.GetCell(2, 0).Char != 'A' {
		t.Errorf("Expected 'A' at (2,0), got '%c'", flipped.GetCell(2, 0).Char)
	}
	
	if flipped.GetCell(0, 1).Char != 'F' {
		t.Errorf("Expected 'F' at (0,1), got '%c'", flipped.GetCell(0, 1).Char)
	}
}

func TestSpriteFlipVertical(t *testing.T) {
	lines := []string{
		"ABC",
		"DEF",
	}
	
	sprite := NewSpriteFromString(lines, tcell.ColorWhite, tcell.ColorBlack)
	flipped := sprite.FlipVertical()
	
	if flipped.GetCell(0, 0).Char != 'D' {
		t.Errorf("Expected 'D' at (0,0), got '%c'", flipped.GetCell(0, 0).Char)
	}
	
	if flipped.GetCell(0, 1).Char != 'A' {
		t.Errorf("Expected 'A' at (0,1), got '%c'", flipped.GetCell(0, 1).Char)
	}
}

func TestCommonSprites(t *testing.T) {
	// Test that common sprites can be created without panicking
	sprites := []struct {
		name   string
		sprite *Sprite
	}{
		{"Player", PlayerSprite()},
		{"Goomba", GoombaSprite()},
		{"Coin", CoinSprite()},
		{"Block", BlockSprite()},
		{"Brick", BrickSprite()},
		{"Pipe", PipeSprite()},
		{"Mushroom", MushroomSprite()},
		{"FireFlower", FireFlowerSprite()},
		{"Star", StarSprite()},
		{"Cloud", CloudSprite()},
		{"Bush", BushSprite()},
		{"Flag", FlagSprite()},
	}
	
	for _, tt := range sprites {
		t.Run(tt.name, func(t *testing.T) {
			if tt.sprite == nil {
				t.Errorf("%s sprite is nil", tt.name)
			}
			
			if tt.sprite.Width <= 0 {
				t.Errorf("%s sprite has invalid width: %d", tt.name, tt.sprite.Width)
			}
			
			if tt.sprite.Height <= 0 {
				t.Errorf("%s sprite has invalid height: %d", tt.name, tt.sprite.Height)
			}
		})
	}
}

func TestRectangle(t *testing.T) {
	rect := Rectangle{
		X:      10,
		Y:      20,
		Width:  30,
		Height: 40,
	}
	
	if rect.X != 10 {
		t.Errorf("Expected X=10, got %d", rect.X)
	}
	
	if rect.Y != 20 {
		t.Errorf("Expected Y=20, got %d", rect.Y)
	}
	
	if rect.Width != 30 {
		t.Errorf("Expected Width=30, got %d", rect.Width)
	}
	
	if rect.Height != 40 {
		t.Errorf("Expected Height=40, got %d", rect.Height)
	}
}
