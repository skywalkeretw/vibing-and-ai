package camera

import (
	"testing"
)

func TestNewSplitScreenCamera_Horizontal(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	splitCam := NewSplitScreenCamera(80, 24, SplitHorizontal, bounds)

	if splitCam == nil {
		t.Fatal("NewSplitScreenCamera returned nil")
	}

	if splitCam.mode != SplitHorizontal {
		t.Errorf("Expected mode SplitHorizontal, got %v", splitCam.mode)
	}

	// Check camera 1 dimensions (top half)
	width1, height1 := splitCam.camera1.GetSize()
	if width1 != 80 {
		t.Errorf("Camera1 width: expected 80, got %d", width1)
	}
	if height1 != 12 {
		t.Errorf("Camera1 height: expected 12 (24/2), got %d", height1)
	}

	// Check camera 2 dimensions (bottom half)
	width2, height2 := splitCam.camera2.GetSize()
	if width2 != 80 {
		t.Errorf("Camera2 width: expected 80, got %d", width2)
	}
	if height2 != 12 {
		t.Errorf("Camera2 height: expected 12 (24/2), got %d", height2)
	}
}

func TestNewSplitScreenCamera_Vertical(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	splitCam := NewSplitScreenCamera(80, 24, SplitVertical, bounds)

	if splitCam == nil {
		t.Fatal("NewSplitScreenCamera returned nil")
	}

	if splitCam.mode != SplitVertical {
		t.Errorf("Expected mode SplitVertical, got %v", splitCam.mode)
	}

	// Check camera 1 dimensions (left half)
	width1, height1 := splitCam.camera1.GetSize()
	if width1 != 40 {
		t.Errorf("Camera1 width: expected 40 (80/2), got %d", width1)
	}
	if height1 != 24 {
		t.Errorf("Camera1 height: expected 24, got %d", height1)
	}

	// Check camera 2 dimensions (right half)
	width2, height2 := splitCam.camera2.GetSize()
	if width2 != 40 {
		t.Errorf("Camera2 width: expected 40 (80/2), got %d", width2)
	}
	if height2 != 24 {
		t.Errorf("Camera2 height: expected 24, got %d", height2)
	}
}

func TestSplitScreenCamera_Update(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	splitCam := NewSplitScreenCamera(80, 24, SplitHorizontal, bounds)

	player1 := &MockPlayer{x: 100, y: 50}
	player2 := &MockPlayer{x: 200, y: 100}

	splitCam.camera1.SetMode(CameraModeFollow)
	splitCam.camera2.SetMode(CameraModeFollow)

	splitCam.Update(0.016, player1, player2)

	// Check camera 1 followed player 1
	x1, y1 := splitCam.camera1.GetPosition()
	expectedX1 := 100 - 80/2.0
	expectedY1 := 50 - 12/2.0

	if x1 != expectedX1 {
		t.Errorf("Camera1 x: expected %f, got %f", expectedX1, x1)
	}
	if y1 != expectedY1 {
		t.Errorf("Camera1 y: expected %f, got %f", expectedY1, y1)
	}

	// Check camera 2 followed player 2
	x2, y2 := splitCam.camera2.GetPosition()
	expectedX2 := 200 - 80/2.0
	expectedY2 := 100 - 12/2.0

	if x2 != expectedX2 {
		t.Errorf("Camera2 x: expected %f, got %f", expectedX2, x2)
	}
	if y2 != expectedY2 {
		t.Errorf("Camera2 y: expected %f, got %f", expectedY2, y2)
	}
}

func TestSplitScreenCamera_UpdateWithNilPlayers(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	splitCam := NewSplitScreenCamera(80, 24, SplitHorizontal, bounds)

	splitCam.camera1.SetPosition(10, 10)
	splitCam.camera2.SetPosition(20, 20)

	// Should not crash with nil players
	splitCam.Update(0.016, nil, nil)

	x1, y1 := splitCam.camera1.GetPosition()
	if x1 != 10 || y1 != 10 {
		t.Error("Camera1 should not move with nil player")
	}

	x2, y2 := splitCam.camera2.GetPosition()
	if x2 != 20 || y2 != 20 {
		t.Error("Camera2 should not move with nil player")
	}
}

func TestSplitScreenCamera_GetCameras(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	splitCam := NewSplitScreenCamera(80, 24, SplitHorizontal, bounds)

	cam1 := splitCam.GetCamera1()
	cam2 := splitCam.GetCamera2()

	if cam1 == nil {
		t.Error("GetCamera1 returned nil")
	}

	if cam2 == nil {
		t.Error("GetCamera2 returned nil")
	}

	if cam1 == cam2 {
		t.Error("Camera1 and Camera2 should be different instances")
	}
}

func TestSplitScreenCamera_GetMode(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}

	splitCamH := NewSplitScreenCamera(80, 24, SplitHorizontal, bounds)
	if splitCamH.GetMode() != SplitHorizontal {
		t.Errorf("Expected SplitHorizontal, got %v", splitCamH.GetMode())
	}

	splitCamV := NewSplitScreenCamera(80, 24, SplitVertical, bounds)
	if splitCamV.GetMode() != SplitVertical {
		t.Errorf("Expected SplitVertical, got %v", splitCamV.GetMode())
	}
}

func TestSplitScreenCamera_WorldToScreen1(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	splitCam := NewSplitScreenCamera(80, 24, SplitHorizontal, bounds)

	splitCam.camera1.SetPosition(100, 50)

	screenX, screenY := splitCam.WorldToScreen1(150, 75)

	expectedX := 150 - 100
	expectedY := 75 - 50

	if screenX != expectedX {
		t.Errorf("Expected screen x %d, got %d", expectedX, screenX)
	}

	if screenY != expectedY {
		t.Errorf("Expected screen y %d, got %d", expectedY, screenY)
	}
}

func TestSplitScreenCamera_WorldToScreen2_Horizontal(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	splitCam := NewSplitScreenCamera(80, 24, SplitHorizontal, bounds)

	splitCam.camera2.SetPosition(100, 50)

	screenX, screenY := splitCam.WorldToScreen2(150, 75)

	// In horizontal split, camera 2 is below camera 1
	// So Y coordinate should be offset by camera1's height
	expectedX := 150 - 100
	expectedY := (75 - 50) + 12 // +12 because camera1 height is 12

	if screenX != expectedX {
		t.Errorf("Expected screen x %d, got %d", expectedX, screenX)
	}

	if screenY != expectedY {
		t.Errorf("Expected screen y %d, got %d", expectedY, screenY)
	}
}

func TestSplitScreenCamera_WorldToScreen2_Vertical(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	splitCam := NewSplitScreenCamera(80, 24, SplitVertical, bounds)

	splitCam.camera2.SetPosition(100, 50)

	screenX, screenY := splitCam.WorldToScreen2(150, 75)

	// In vertical split, camera 2 is to the right of camera 1
	// So X coordinate should be offset by camera1's width
	expectedX := (150 - 100) + 40 // +40 because camera1 width is 40
	expectedY := 75 - 50

	if screenX != expectedX {
		t.Errorf("Expected screen x %d, got %d", expectedX, screenX)
	}

	if screenY != expectedY {
		t.Errorf("Expected screen y %d, got %d", expectedY, screenY)
	}
}

func TestSplitScreenCamera_IsVisible(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	splitCam := NewSplitScreenCamera(80, 24, SplitHorizontal, bounds)

	splitCam.camera1.SetPosition(100, 50)
	splitCam.camera2.SetPosition(200, 100)

	// Test visibility in camera 1
	if !splitCam.IsVisible1(120, 55, 10, 10) {
		t.Error("Entity should be visible in camera 1")
	}

	if splitCam.IsVisible1(50, 55, 10, 10) {
		t.Error("Entity should not be visible in camera 1")
	}

	// Test visibility in camera 2
	if !splitCam.IsVisible2(220, 105, 10, 10) {
		t.Error("Entity should be visible in camera 2")
	}

	if splitCam.IsVisible2(150, 105, 10, 10) {
		t.Error("Entity should not be visible in camera 2")
	}
}

func TestSplitScreenCamera_Shake(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	splitCam := NewSplitScreenCamera(80, 24, SplitHorizontal, bounds)

	splitCam.Shake(2.0, 0.5)

	// Both cameras should have shake applied
	if splitCam.camera1.shakeAmount != 2.0 {
		t.Errorf("Camera1 shake amount: expected 2.0, got %f", splitCam.camera1.shakeAmount)
	}

	if splitCam.camera1.shakeTime != 0.5 {
		t.Errorf("Camera1 shake time: expected 0.5, got %f", splitCam.camera1.shakeTime)
	}

	if splitCam.camera2.shakeAmount != 2.0 {
		t.Errorf("Camera2 shake amount: expected 2.0, got %f", splitCam.camera2.shakeAmount)
	}

	if splitCam.camera2.shakeTime != 0.5 {
		t.Errorf("Camera2 shake time: expected 0.5, got %f", splitCam.camera2.shakeTime)
	}
}

func TestSplitScreenCamera_ShakeIndividual(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	splitCam := NewSplitScreenCamera(80, 24, SplitHorizontal, bounds)

	// Shake only camera 1
	splitCam.ShakeCamera1(1.5, 0.3)

	if splitCam.camera1.shakeAmount != 1.5 {
		t.Errorf("Camera1 shake amount: expected 1.5, got %f", splitCam.camera1.shakeAmount)
	}

	if splitCam.camera2.shakeAmount != 0 {
		t.Errorf("Camera2 should not have shake, got %f", splitCam.camera2.shakeAmount)
	}

	// Shake only camera 2
	splitCam.ShakeCamera2(2.5, 0.4)

	if splitCam.camera2.shakeAmount != 2.5 {
		t.Errorf("Camera2 shake amount: expected 2.5, got %f", splitCam.camera2.shakeAmount)
	}
}

func TestSplitScreenCamera_SetBounds(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	splitCam := NewSplitScreenCamera(80, 24, SplitHorizontal, bounds)

	newBounds := Rect{X: 10, Y: 20, Width: 800, Height: 400}
	splitCam.SetBounds(newBounds)

	bounds1 := splitCam.camera1.GetBounds()
	bounds2 := splitCam.camera2.GetBounds()

	if bounds1 != newBounds {
		t.Errorf("Camera1 bounds not updated correctly")
	}

	if bounds2 != newBounds {
		t.Errorf("Camera2 bounds not updated correctly")
	}
}

func TestSplitScreenCamera_SetMode(t *testing.T) {
	bounds := Rect{X: 0, Y: 0, Width: 1000, Height: 500}
	splitCam := NewSplitScreenCamera(80, 24, SplitHorizontal, bounds)

	// Initially horizontal
	if splitCam.GetMode() != SplitHorizontal {
		t.Error("Initial mode should be SplitHorizontal")
	}

	// Switch to vertical - this is a no-op in current implementation
	// but we test that it doesn't crash
	splitCam.SetMode(SplitVertical)

	// Setting same mode should be no-op
	splitCam.SetMode(SplitVertical)
}
