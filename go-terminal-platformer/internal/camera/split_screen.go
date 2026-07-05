package camera

// SplitMode defines how the screen is split for multiplayer
type SplitMode int

const (
	SplitHorizontal SplitMode = iota
	SplitVertical
)

// SplitScreenCamera manages two cameras for split-screen multiplayer
type SplitScreenCamera struct {
	camera1 *Camera
	camera2 *Camera
	mode    SplitMode
}

// NewSplitScreenCamera creates a new split-screen camera system
func NewSplitScreenCamera(width, height int, mode SplitMode, bounds Rect) *SplitScreenCamera {
	var cam1, cam2 *Camera

	if mode == SplitHorizontal {
		// Split screen horizontally (top/bottom)
		cam1 = NewCamera(width, height/2, bounds)
		cam2 = NewCamera(width, height/2, bounds)
	} else {
		// Split screen vertically (left/right)
		cam1 = NewCamera(width/2, height, bounds)
		cam2 = NewCamera(width/2, height, bounds)
	}

	return &SplitScreenCamera{
		camera1: cam1,
		camera2: cam2,
		mode:    mode,
	}
}

// Update updates both cameras based on their respective players
func (sc *SplitScreenCamera) Update(deltaTime float64, player1, player2 Player) {
	if player1 != nil {
		sc.camera1.Update(deltaTime, player1)
	}
	if player2 != nil {
		sc.camera2.Update(deltaTime, player2)
	}
}

// GetCamera1 returns the first camera
func (sc *SplitScreenCamera) GetCamera1() *Camera {
	return sc.camera1
}

// GetCamera2 returns the second camera
func (sc *SplitScreenCamera) GetCamera2() *Camera {
	return sc.camera2
}

// GetMode returns the split mode
func (sc *SplitScreenCamera) GetMode() SplitMode {
	return sc.mode
}

// SetMode changes the split mode and resizes cameras accordingly
func (sc *SplitScreenCamera) SetMode(mode SplitMode) {
	if sc.mode == mode {
		return
	}

	sc.mode = mode

	// Get current dimensions from camera1
	width, height := sc.camera1.GetSize()
	bounds := sc.camera1.GetBounds()

	// Recalculate based on original full screen size
	if mode == SplitHorizontal {
		// If switching to horizontal, we need to adjust
		if sc.mode == SplitVertical {
			width = width * 2
			height = height / 2
		}
		sc.camera1.width = width
		sc.camera1.height = height / 2
		sc.camera2.width = width
		sc.camera2.height = height / 2
	} else {
		// If switching to vertical
		if sc.mode == SplitHorizontal {
			width = width / 2
			height = height * 2
		}
		sc.camera1.width = width / 2
		sc.camera1.height = height
		sc.camera2.width = width / 2
		sc.camera2.height = height
	}

	// Ensure bounds are still valid
	sc.camera1.SetBounds(bounds)
	sc.camera2.SetBounds(bounds)
}

// WorldToScreen1 converts world coordinates to screen coordinates for camera 1
func (sc *SplitScreenCamera) WorldToScreen1(worldX, worldY float64) (int, int) {
	return sc.camera1.WorldToScreen(worldX, worldY)
}

// WorldToScreen2 converts world coordinates to screen coordinates for camera 2
// The coordinates are adjusted based on split mode
func (sc *SplitScreenCamera) WorldToScreen2(worldX, worldY float64) (int, int) {
	screenX, screenY := sc.camera2.WorldToScreen(worldX, worldY)

	// Adjust coordinates based on split mode
	if sc.mode == SplitHorizontal {
		// Camera 2 is below camera 1
		screenY += sc.camera1.height
	} else {
		// Camera 2 is to the right of camera 1
		screenX += sc.camera1.width
	}

	return screenX, screenY
}

// IsVisible1 checks if an entity is visible in camera 1's viewport
func (sc *SplitScreenCamera) IsVisible1(x, y, width, height float64) bool {
	return sc.camera1.IsVisible(x, y, width, height)
}

// IsVisible2 checks if an entity is visible in camera 2's viewport
func (sc *SplitScreenCamera) IsVisible2(x, y, width, height float64) bool {
	return sc.camera2.IsVisible(x, y, width, height)
}

// Shake applies screen shake to both cameras
func (sc *SplitScreenCamera) Shake(amount, duration float64) {
	sc.camera1.Shake(amount, duration)
	sc.camera2.Shake(amount, duration)
}

// ShakeCamera1 applies screen shake only to camera 1
func (sc *SplitScreenCamera) ShakeCamera1(amount, duration float64) {
	sc.camera1.Shake(amount, duration)
}

// ShakeCamera2 applies screen shake only to camera 2
func (sc *SplitScreenCamera) ShakeCamera2(amount, duration float64) {
	sc.camera2.Shake(amount, duration)
}

// SetBounds updates the bounds for both cameras
func (sc *SplitScreenCamera) SetBounds(bounds Rect) {
	sc.camera1.SetBounds(bounds)
	sc.camera2.SetBounds(bounds)
}
