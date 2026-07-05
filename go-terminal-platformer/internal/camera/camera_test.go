package camera

import (
	"math"
	"testing"
)

// MockPlayer implements the Player interface for testing
type MockPlayer struct {
	x         float64
	y         float64
	velocityX float64
	velocityY float64
}

func (m *MockPlayer) GetPosition() (float64, float64) {
	return m.x, m.y
}

func (m *MockPlayer) GetVelocity() (float64, float64) {
	return m.velocityX, m.velocityY
}

func TestNewCamera(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)

	if camera == nil {
		t.Fatal("NewCamera returned nil")
	}

	if camera.width != 80 {
		t.Errorf("Expected width 80, got %d", camera.width)
	}

	if camera.height != 24 {
		t.Errorf("Expected height 24, got %d", camera.height)
	}

	if camera.smoothing != 0.1 {
		t.Errorf("Expected smoothing 0.1, got %f", camera.smoothing)
	}

	if camera.mode != CameraModeLerp {
		t.Errorf("Expected mode CameraModeLerp, got %v", camera.mode)
	}

	if camera.lookAhead != 5.0 {
		t.Errorf("Expected lookAhead 5.0, got %f", camera.lookAhead)
	}
}

func TestCameraUpdate_Follow(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)
	camera.SetMode(CameraModeFollow)

	player := &MockPlayer{x: 100, y: 50}
	camera.Update(0.016, player)

	// Camera should center on player instantly in follow mode
	expectedX := 100 - 80/2.0
	expectedY := 50 - 24/2.0

	if camera.x != expectedX {
		t.Errorf("Expected camera x %f, got %f", expectedX, camera.x)
	}

	if camera.y != expectedY {
		t.Errorf("Expected camera y %f, got %f", expectedY, camera.y)
	}
}

func TestCameraUpdate_Lerp(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)
	camera.SetMode(CameraModeLerp)
	camera.SetSmoothing(0.5) // 50% smoothing for easier testing

	player := &MockPlayer{x: 100, y: 50}
	camera.Update(0.016, player)

	// Camera should move towards player but not reach it instantly
	expectedTargetX := 100 - 80/2.0
	expectedTargetY := 50 - 24/2.0

	// With 0.5 smoothing, camera should be halfway to target
	expectedX := 0 + (expectedTargetX-0)*0.5
	expectedY := 0 + (expectedTargetY-0)*0.5

	if math.Abs(camera.x-expectedX) > 0.01 {
		t.Errorf("Expected camera x ~%f, got %f", expectedX, camera.x)
	}

	if math.Abs(camera.y-expectedY) > 0.01 {
		t.Errorf("Expected camera y ~%f, got %f", expectedY, camera.y)
	}
}

func TestCameraUpdate_LookAhead(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)
	camera.SetMode(CameraModeLookAhead)
	camera.SetSmoothing(1.0) // Instant for easier testing
	camera.SetLookAhead(10.0)

	// Player moving right
	player := &MockPlayer{x: 100, y: 50, velocityX: 5.0}
	camera.Update(0.016, player)

	expectedX := (100 - 80/2.0) + 10.0 // Base position + look-ahead

	if math.Abs(camera.x-expectedX) > 0.01 {
		t.Errorf("Expected camera x with look-ahead ~%f, got %f", expectedX, camera.x)
	}

	// Player moving left
	player.velocityX = -5.0
	camera.SetPosition(0, 0) // Reset
	camera.Update(0.016, player)

	expectedX = (100 - 80/2.0) - 10.0 // Base position - look-ahead

	if math.Abs(camera.x-expectedX) > 0.01 {
		t.Errorf("Expected camera x with negative look-ahead ~%f, got %f", expectedX, camera.x)
	}
}

func TestCameraUpdate_Fixed(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)
	camera.SetMode(CameraModeFixed)
	camera.SetPosition(50, 25)

	player := &MockPlayer{x: 100, y: 50}
	camera.Update(0.016, player)

	// Camera should not move in fixed mode
	if camera.x != 50 {
		t.Errorf("Expected camera x to stay at 50, got %f", camera.x)
	}

	if camera.y != 25 {
		t.Errorf("Expected camera y to stay at 25, got %f", camera.y)
	}
}

func TestCameraUpdate_NilPlayer(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)
	camera.SetPosition(10, 10)

	camera.Update(0.016, nil)

	// Camera should not crash and should stay at current position
	if camera.x != 10 {
		t.Errorf("Expected camera x to stay at 10, got %f", camera.x)
	}
}

func TestCameraConstrainToBounds(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 200, Height: 100}
	camera := NewCamera(80, 24, bounds)
	camera.SetMode(CameraModeFollow)

	// Test left boundary
	player := &MockPlayer{x: 10, y: 50}
	camera.Update(0.016, player)
	if camera.x < bounds.X {
		t.Errorf("Camera x %f exceeded left boundary %f", camera.x, bounds.X)
	}

	// Test right boundary
	player.x = 190
	camera.Update(0.016, player)
	if camera.x+float64(camera.width) > bounds.X+bounds.Width {
		t.Errorf("Camera exceeded right boundary")
	}

	// Test top boundary
	player.x = 100
	player.y = 5
	camera.Update(0.016, player)
	if camera.y < bounds.Y {
		t.Errorf("Camera y %f exceeded top boundary %f", camera.y, bounds.Y)
	}

	// Test bottom boundary
	player.y = 95
	camera.Update(0.016, player)
	if camera.y+float64(camera.height) > bounds.Y+bounds.Height {
		t.Errorf("Camera exceeded bottom boundary")
	}
}

func TestWorldToScreen(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)
	camera.SetPosition(100, 50)

	screenX, screenY := camera.WorldToScreen(150, 75)

	expectedX := int(math.Round(150 - 100))
	expectedY := int(math.Round(75 - 50))

	if screenX != expectedX {
		t.Errorf("Expected screen x %d, got %d", expectedX, screenX)
	}

	if screenY != expectedY {
		t.Errorf("Expected screen y %d, got %d", expectedY, screenY)
	}
}

func TestScreenToWorld(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)
	camera.SetPosition(100, 50)

	worldX, worldY := camera.ScreenToWorld(10, 5)

	expectedX := 10.0 + 100
	expectedY := 5.0 + 50

	if worldX != expectedX {
		t.Errorf("Expected world x %f, got %f", expectedX, worldX)
	}

	if worldY != expectedY {
		t.Errorf("Expected world y %f, got %f", expectedY, worldY)
	}
}

func TestIsVisible(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)
	camera.SetPosition(100, 50)

	tests := []struct {
		name     string
		x, y     float64
		w, h     float64
		expected bool
	}{
		{"Fully visible", 120, 60, 10, 10, true},
		{"Partially visible left", 95, 60, 10, 10, true},
		{"Partially visible right", 175, 60, 10, 10, true},
		{"Partially visible top", 120, 45, 10, 10, true},
		{"Partially visible bottom", 120, 69, 10, 10, true},
		{"Not visible left", 50, 60, 10, 10, false},
		{"Not visible right", 200, 60, 10, 10, false},
		{"Not visible top", 120, 20, 10, 10, false},
		{"Not visible bottom", 120, 100, 10, 10, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := camera.IsVisible(tt.x, tt.y, tt.w, tt.h)
			if result != tt.expected {
				t.Errorf("IsVisible(%f, %f, %f, %f) = %v, expected %v",
					tt.x, tt.y, tt.w, tt.h, result, tt.expected)
			}
		})
	}
}

func TestShake(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)
	camera.SetPosition(100, 50)

	camera.Shake(2.0, 0.5)

	if camera.shakeAmount != 2.0 {
		t.Errorf("Expected shake amount 2.0, got %f", camera.shakeAmount)
	}

	if camera.shakeTime != 0.5 {
		t.Errorf("Expected shake time 0.5, got %f", camera.shakeTime)
	}

	// Update to apply shake
	player := &MockPlayer{x: 100, y: 50}
	originalX := camera.x
	camera.Update(0.016, player)

	// Position should have changed due to shake
	if camera.x == originalX {
		t.Error("Camera position should have changed due to shake")
	}

	// Shake time should decrease
	if camera.shakeTime >= 0.5 {
		t.Errorf("Shake time should have decreased, got %f", camera.shakeTime)
	}
}

func TestShakeExpires(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)
	camera.SetMode(CameraModeFixed)
	camera.SetPosition(100, 50)

	camera.Shake(2.0, 0.1)

	player := &MockPlayer{x: 100, y: 50}

	// Update with enough time to expire shake
	camera.Update(0.2, player)

	if camera.shakeTime > 0 {
		t.Errorf("Shake time should be 0 or less after expiring, got %f", camera.shakeTime)
	}
}

func TestGetPosition(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)
	camera.SetPosition(123.45, 67.89)

	x, y := camera.GetPosition()

	if x != 123.45 {
		t.Errorf("Expected x 123.45, got %f", x)
	}

	if y != 67.89 {
		t.Errorf("Expected y 67.89, got %f", y)
	}
}

func TestGetSize(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)

	width, height := camera.GetSize()

	if width != 80 {
		t.Errorf("Expected width 80, got %d", width)
	}

	if height != 24 {
		t.Errorf("Expected height 24, got %d", height)
	}
}

func TestSetGetMode(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)

	modes := []CameraMode{
		CameraModeFollow,
		CameraModeLerp,
		CameraModeLookAhead,
		CameraModeFixed,
	}

	for _, mode := range modes {
		camera.SetMode(mode)
		if camera.GetMode() != mode {
			t.Errorf("Expected mode %v, got %v", mode, camera.GetMode())
		}
	}
}

func TestSetGetSmoothing(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)

	// Test normal value
	camera.SetSmoothing(0.5)
	if camera.GetSmoothing() != 0.5 {
		t.Errorf("Expected smoothing 0.5, got %f", camera.GetSmoothing())
	}

	// Test clamping to 0
	camera.SetSmoothing(-0.5)
	if camera.GetSmoothing() != 0 {
		t.Errorf("Expected smoothing to be clamped to 0, got %f", camera.GetSmoothing())
	}

	// Test clamping to 1
	camera.SetSmoothing(1.5)
	if camera.GetSmoothing() != 1 {
		t.Errorf("Expected smoothing to be clamped to 1, got %f", camera.GetSmoothing())
	}
}

func TestSetGetLookAhead(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)

	camera.SetLookAhead(15.0)
	if camera.GetLookAhead() != 15.0 {
		t.Errorf("Expected look-ahead 15.0, got %f", camera.GetLookAhead())
	}
}

func TestSetGetBounds(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)

	newBounds := Rect{X: 10, Y: 20, Width: 800, Height: 400}
	camera.SetBounds(newBounds)

	retrievedBounds := camera.GetBounds()

	if retrievedBounds.X != newBounds.X ||
		retrievedBounds.Y != newBounds.Y ||
		retrievedBounds.Width != newBounds.Width ||
		retrievedBounds.Height != newBounds.Height {
		t.Errorf("Bounds mismatch: expected %+v, got %+v", newBounds, retrievedBounds)
	}
}

func TestSetPosition(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	camera := NewCamera(80, 24, bounds)

	camera.SetPosition(200, 100)

	if camera.x != 200 {
		t.Errorf("Expected x 200, got %f", camera.x)
	}

	if camera.y != 100 {
		t.Errorf("Expected y 100, got %f", camera.y)
	}

	if camera.targetX != 200 {
		t.Errorf("Expected targetX 200, got %f", camera.targetX)
	}

	if camera.targetY != 100 {
		t.Errorf("Expected targetY 100, got %f", camera.targetY)
	}
}
