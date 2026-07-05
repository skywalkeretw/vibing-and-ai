package entities

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewPiranhaPlant(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	if plant == nil {
		t.Fatal("NewPiranhaPlant returned nil")
	}
	
	if plant.pipePosition.X != 100 || plant.pipePosition.Y != 200 {
		t.Errorf("Expected pipe position (100, 200), got (%f, %f)", plant.pipePosition.X, plant.pipePosition.Y)
	}
	
	if plant.Health != 2 {
		t.Errorf("Expected health 2, got %d", plant.Health)
	}
	
	if plant.MaxHealth != 2 {
		t.Errorf("Expected max health 2, got %d", plant.MaxHealth)
	}
	
	if plant.MoveSpeed != 0 {
		t.Errorf("Expected move speed 0 (stationary), got %f", plant.MoveSpeed)
	}
	
	if plant.plantState != PiranhaPlantStateHidden {
		t.Errorf("Expected initial state PiranhaPlantStateHidden, got %v", plant.plantState)
	}
	
	if plant.Sprite != 'P' {
		t.Errorf("Expected sprite 'P', got '%c'", plant.Sprite)
	}
	
	if plant.isEmerged {
		t.Error("Expected plant to start hidden")
	}
	
	if plant.emergeHeight != 0 {
		t.Errorf("Expected emerge height 0, got %f", plant.emergeHeight)
	}
	
	if plant.visibleDuration != 3.0 {
		t.Errorf("Expected visible duration 3.0, got %f", plant.visibleDuration)
	}
	
	if plant.hiddenDuration != 2.0 {
		t.Errorf("Expected hidden duration 2.0, got %f", plant.hiddenDuration)
	}
}

func TestPiranhaPlantHiddenState(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	if plant.plantState != PiranhaPlantStateHidden {
		t.Fatal("Expected initial state to be hidden")
	}
	
	// Update for less than hidden duration
	plant.Update(1.0)
	
	if plant.plantState != PiranhaPlantStateHidden {
		t.Error("Plant should remain hidden before timer expires")
	}
	
	if plant.emergeHeight != 0 {
		t.Errorf("Emerge height should be 0 when hidden, got %f", plant.emergeHeight)
	}
	
	if plant.isEmerged {
		t.Error("Plant should not be emerged when hidden")
	}
}

func TestPiranhaPlantEmergeCycle(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	// Wait for hidden duration to complete
	plant.Update(2.1) // Slightly more than 2 seconds
	
	if plant.plantState != PiranhaPlantStateEmerging {
		t.Errorf("Expected PiranhaPlantStateEmerging after hidden duration, got %v", plant.plantState)
	}
	
	// Continue updating to complete emergence
	for i := 0; i < 20; i++ {
		plant.Update(0.1)
		if plant.plantState == PiranhaPlantStateVisible {
			break
		}
	}
	
	if plant.plantState != PiranhaPlantStateVisible {
		t.Errorf("Expected PiranhaPlantStateVisible after emergence, got %v", plant.plantState)
	}
	
	if !plant.isEmerged {
		t.Error("Plant should be emerged when visible")
	}
	
	if plant.emergeHeight < plant.maxHeight-1.0 {
		t.Errorf("Emerge height should be near max (%f), got %f", plant.maxHeight, plant.emergeHeight)
	}
}

func TestPiranhaPlantHideCycle(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	// Force plant to visible state
	plant.plantState = PiranhaPlantStateVisible
	plant.isEmerged = true
	plant.emergeHeight = plant.maxHeight
	plant.emergeTimer = 0
	
	// Wait for visible duration
	plant.Update(3.1) // Slightly more than 3 seconds
	
	if plant.plantState != PiranhaPlantStateHiding {
		t.Errorf("Expected PiranhaPlantStateHiding after visible duration, got %v", plant.plantState)
	}
	
	// Continue updating to complete hiding
	for i := 0; i < 20; i++ {
		plant.Update(0.1)
		if plant.plantState == PiranhaPlantStateHidden {
			break
		}
	}
	
	if plant.plantState != PiranhaPlantStateHidden {
		t.Errorf("Expected PiranhaPlantStateHidden after hiding, got %v", plant.plantState)
	}
	
	if plant.isEmerged {
		t.Error("Plant should not be emerged when hidden")
	}
	
	if plant.emergeHeight > 1.0 {
		t.Errorf("Emerge height should be near 0, got %f", plant.emergeHeight)
	}
}

func TestPiranhaPlantStompImmunity(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	initialHealth := plant.Health
	
	// Try to stomp the plant
	plant.OnStomp(nil)
	
	// Health should not change (stomp immunity)
	if plant.Health != initialHealth {
		t.Errorf("Piranha Plant should be immune to stomps, health changed from %d to %d", initialHealth, plant.Health)
	}
	
	if plant.IsDead {
		t.Error("Piranha Plant should not die from stomp")
	}
}

func TestPiranhaPlantProjectileDefeat(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	if plant.Health != 2 {
		t.Fatal("Expected initial health 2")
	}
	
	// Hit with projectile
	plant.OnHitByProjectile(nil)
	
	if plant.Health != 1 {
		t.Errorf("Expected health 1 after first projectile hit, got %d", plant.Health)
	}
	
	if plant.IsDead {
		t.Error("Plant should not be dead after first hit")
	}
	
	// Wait for invulnerability timer to expire
	plant.Update(0.6) // Wait 0.6 seconds to clear invulnerability
	
	// Hit with second projectile
	plant.OnHitByProjectile(nil)
	
	if plant.Health != 0 {
		t.Errorf("Expected health 0 after second projectile hit, got %d", plant.Health)
	}
	
	if !plant.IsDead {
		t.Error("Plant should be dead after second projectile hit")
	}
}

func TestPiranhaPlantIsEmerged(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	if plant.IsEmerged() {
		t.Error("Plant should not be emerged initially")
	}
	
	// Force to visible state
	plant.plantState = PiranhaPlantStateVisible
	plant.isEmerged = true
	
	if !plant.IsEmerged() {
		t.Error("Plant should be emerged when in visible state")
	}
}

func TestPiranhaPlantGetEmergeHeight(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	height := plant.GetEmergeHeight()
	if height != 0 {
		t.Errorf("Expected initial emerge height 0, got %f", height)
	}
	
	// Set emerge height
	plant.emergeHeight = 16.0
	
	height = plant.GetEmergeHeight()
	if height != 16.0 {
		t.Errorf("Expected emerge height 16.0, got %f", height)
	}
}

func TestPiranhaPlantGetPlantState(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	state := plant.GetPlantState()
	if state != PiranhaPlantStateHidden {
		t.Errorf("Expected initial state PiranhaPlantStateHidden, got %v", state)
	}
	
	plant.plantState = PiranhaPlantStateVisible
	state = plant.GetPlantState()
	if state != PiranhaPlantStateVisible {
		t.Errorf("Expected state PiranhaPlantStateVisible, got %v", state)
	}
}

func TestPiranhaPlantSetPipePosition(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	plant.SetPipePosition(150, 250)
	
	if plant.pipePosition.X != 150 || plant.pipePosition.Y != 250 {
		t.Errorf("Expected pipe position (150, 250), got (%f, %f)", plant.pipePosition.X, plant.pipePosition.Y)
	}
	
	if plant.Position.X != 150 || plant.Position.Y != 250 {
		t.Errorf("Expected position (150, 250), got (%f, %f)", plant.Position.X, plant.Position.Y)
	}
}

func TestPiranhaPlantGetPipePosition(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	pos := plant.GetPipePosition()
	if pos.X != 100 || pos.Y != 200 {
		t.Errorf("Expected pipe position (100, 200), got (%f, %f)", pos.X, pos.Y)
	}
}

func TestPiranhaPlantSetDurations(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	plant.SetVisibleDuration(5.0)
	if plant.visibleDuration != 5.0 {
		t.Errorf("Expected visible duration 5.0, got %f", plant.visibleDuration)
	}
	
	plant.SetHiddenDuration(1.5)
	if plant.hiddenDuration != 1.5 {
		t.Errorf("Expected hidden duration 1.5, got %f", plant.hiddenDuration)
	}
}

func TestPiranhaPlantBlocksPipe(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	// Hidden state should not block pipe
	if plant.BlocksPipe() {
		t.Error("Hidden plant should not block pipe")
	}
	
	// Emerging state should block pipe
	plant.plantState = PiranhaPlantStateEmerging
	if !plant.BlocksPipe() {
		t.Error("Emerging plant should block pipe")
	}
	
	// Visible state should block pipe
	plant.plantState = PiranhaPlantStateVisible
	if !plant.BlocksPipe() {
		t.Error("Visible plant should block pipe")
	}
	
	// Hiding state should block pipe
	plant.plantState = PiranhaPlantStateHiding
	if !plant.BlocksPipe() {
		t.Error("Hiding plant should block pipe")
	}
	
	// Hidden state should not block pipe
	plant.plantState = PiranhaPlantStateHidden
	if plant.BlocksPipe() {
		t.Error("Hidden plant should not block pipe")
	}
}

func TestPiranhaPlantPositionUpdate(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	// Set emerge height
	plant.emergeHeight = 16.0
	
	// Update should adjust Y position based on emerge height
	plant.Update(0.016)
	
	expectedY := plant.pipePosition.Y - plant.emergeHeight
	if plant.Position.Y != expectedY {
		t.Errorf("Expected Y position %f, got %f", expectedY, plant.Position.Y)
	}
}

func TestPiranhaPlantGetSprite(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	sprite := plant.GetSprite()
	if sprite != 'P' {
		t.Errorf("Expected sprite 'P', got '%c'", sprite)
	}
}

func TestPiranhaPlantGetSpriteStyle(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	// Hidden state
	style := plant.GetSpriteStyle()
	if style == (tcell.Style{}) {
		t.Error("Expected non-empty sprite style")
	}
	
	// Emerged state should have different color
	plant.isEmerged = true
	emergedStyle := plant.GetSpriteStyle()
	if emergedStyle == (tcell.Style{}) {
		t.Error("Expected non-empty sprite style for emerged state")
	}
}

func TestPiranhaPlantDeadState(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	// Kill the plant
	plant.Die()
	
	initialState := plant.plantState
	
	// Update should not change state when dead
	plant.Update(5.0)
	
	if plant.plantState != initialState {
		t.Error("Dead plant should not change state")
	}
}

func TestPiranhaPlantInactiveState(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	// Deactivate the plant
	plant.Deactivate()
	
	initialState := plant.plantState
	
	// Update should not change state when inactive
	plant.Update(5.0)
	
	if plant.plantState != initialState {
		t.Error("Inactive plant should not change state")
	}
}

func TestPiranhaPlantFullCycle(t *testing.T) {
	plant := NewPiranhaPlant(100, 200)
	
	// Start hidden
	if plant.plantState != PiranhaPlantStateHidden {
		t.Fatal("Expected to start hidden")
	}
	
	// Complete hidden duration
	plant.Update(2.1)
	if plant.plantState != PiranhaPlantStateEmerging {
		t.Error("Should transition to emerging")
	}
	
	// Complete emergence
	for i := 0; i < 20; i++ {
		plant.Update(0.1)
		if plant.plantState == PiranhaPlantStateVisible {
			break
		}
	}
	if plant.plantState != PiranhaPlantStateVisible {
		t.Error("Should transition to visible")
	}
	
	// Complete visible duration
	plant.Update(3.1)
	if plant.plantState != PiranhaPlantStateHiding {
		t.Error("Should transition to hiding")
	}
	
	// Complete hiding
	for i := 0; i < 20; i++ {
		plant.Update(0.1)
		if plant.plantState == PiranhaPlantStateHidden {
			break
		}
	}
	if plant.plantState != PiranhaPlantStateHidden {
		t.Error("Should transition back to hidden")
	}
}

func TestLerpFunction(t *testing.T) {
	// Test lerp function
	result := lerp(0, 10, 0.5)
	if result != 5.0 {
		t.Errorf("Expected lerp(0, 10, 0.5) = 5.0, got %f", result)
	}
	
	result = lerp(0, 10, 0.0)
	if result != 0.0 {
		t.Errorf("Expected lerp(0, 10, 0.0) = 0.0, got %f", result)
	}
	
	result = lerp(0, 10, 1.0)
	if result != 10.0 {
		t.Errorf("Expected lerp(0, 10, 1.0) = 10.0, got %f", result)
	}
}
