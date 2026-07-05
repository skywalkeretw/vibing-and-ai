package entities

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewSpiny(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	if spiny == nil {
		t.Fatal("NewSpiny returned nil")
	}
	
	if spiny.Position.X != 100 || spiny.Position.Y != 200 {
		t.Errorf("Expected position (100, 200), got (%f, %f)", spiny.Position.X, spiny.Position.Y)
	}
	
	if spiny.Health != 1 {
		t.Errorf("Expected health 1, got %d", spiny.Health)
	}
	
	if spiny.MaxHealth != 1 {
		t.Errorf("Expected max health 1, got %d", spiny.MaxHealth)
	}
	
	if spiny.MoveSpeed != 70.0 {
		t.Errorf("Expected move speed 70.0, got %f", spiny.MoveSpeed)
	}
	
	if spiny.Sprite != 'S' {
		t.Errorf("Expected sprite 'S', got '%c'", spiny.Sprite)
	}
	
	if spiny.State != EnemyStatePatrol {
		t.Errorf("Expected initial state EnemyStatePatrol, got %v", spiny.State)
	}
	
	if !spiny.IsActive {
		t.Error("Expected Spiny to be active")
	}
	
	if spiny.IsDead {
		t.Error("Expected Spiny to be alive")
	}
}

func TestSpinyPatrolBehavior(t *testing.T) {
	spiny := NewSpiny(100, 200)
	spiny.OnGround = true
	spiny.FacingRight = true
	
	// Update patrol behavior
	spiny.Update(0.016) // ~60 FPS
	
	if spiny.Velocity.X <= 0 {
		t.Errorf("Expected positive velocity when facing right, got %f", spiny.Velocity.X)
	}
	
	// Test facing left
	spiny.FacingRight = false
	spiny.Update(0.016)
	
	if spiny.Velocity.X >= 0 {
		t.Errorf("Expected negative velocity when facing left, got %f", spiny.Velocity.X)
	}
}

func TestSpinyStompImmunity(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	initialHealth := spiny.Health
	
	// Try to stomp the Spiny
	spiny.OnStomp(nil)
	
	// Health should not change (stomp immunity)
	if spiny.Health != initialHealth {
		t.Errorf("Spiny should be immune to stomps, health changed from %d to %d", initialHealth, spiny.Health)
	}
	
	if spiny.IsDead {
		t.Error("Spiny should not die from stomp")
	}
}

func TestSpinyCanBeStomped(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	if spiny.CanBeStomped() {
		t.Error("Spiny should not be stompable (has spikes)")
	}
}

func TestSpinyProjectileDefeat(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	if spiny.IsDead {
		t.Fatal("Spiny should start alive")
	}
	
	// Hit with projectile
	spiny.OnHitByProjectile(nil)
	
	if !spiny.IsDead {
		t.Error("Spiny should be dead after projectile hit")
	}
	
	if spiny.State != EnemyStateDead {
		t.Errorf("Expected EnemyStateDead, got %v", spiny.State)
	}
}

func TestSpinyShellDefeat(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	if spiny.IsDead {
		t.Fatal("Spiny should start alive")
	}
	
	// Hit with kicked shell
	spiny.OnHitByShell(nil)
	
	if !spiny.IsDead {
		t.Error("Spiny should be dead after shell hit")
	}
	
	if spiny.State != EnemyStateDead {
		t.Errorf("Expected EnemyStateDead, got %v", spiny.State)
	}
}

func TestSpinyWallCollision(t *testing.T) {
	spiny := NewSpiny(100, 200)
	spiny.FacingRight = true
	
	// Hit a wall
	spiny.OnCollideWithWall()
	
	if spiny.FacingRight {
		t.Error("Spiny should turn around after hitting a wall")
	}
	
	// Hit another wall
	spiny.OnCollideWithWall()
	
	if !spiny.FacingRight {
		t.Error("Spiny should turn around again after hitting another wall")
	}
}

func TestSpinyEdgeDetection(t *testing.T) {
	spiny := NewSpiny(100, 200)
	spiny.FacingRight = true
	
	// Reach an edge
	spiny.OnReachEdge()
	
	if spiny.FacingRight {
		t.Error("Spiny should turn around at edge")
	}
	
	// Reach another edge
	spiny.OnReachEdge()
	
	if !spiny.FacingRight {
		t.Error("Spiny should turn around at edge again")
	}
}

func TestSpinyPatrolBounds(t *testing.T) {
	spiny := NewSpiny(100, 200)
	spiny.OnGround = true
	spiny.FacingRight = true
	
	// Set patrol bounds
	spiny.SetPatrolBounds(50, 150)
	
	if spiny.PatrolLeft != 50 {
		t.Errorf("Expected patrol left 50, got %f", spiny.PatrolLeft)
	}
	
	if spiny.PatrolRight != 150 {
		t.Errorf("Expected patrol right 150, got %f", spiny.PatrolRight)
	}
	
	// Move to right boundary
	spiny.Position.X = 151
	spiny.Update(0.016)
	
	if spiny.FacingRight {
		t.Error("Expected Spiny to turn left at right boundary")
	}
	
	// Move to left boundary
	spiny.Position.X = 49
	spiny.Update(0.016)
	
	if !spiny.FacingRight {
		t.Error("Expected Spiny to turn right at left boundary")
	}
}

func TestSpinyUpdate(t *testing.T) {
	spiny := NewSpiny(100, 200)
	spiny.OnGround = true
	
	initialX := spiny.Position.X
	
	// Update should move the Spiny
	spiny.Update(0.1) // 100ms
	
	if spiny.Position.X == initialX {
		t.Error("Expected Spiny position to change after update")
	}
	
	// Dead Spiny should not update
	spiny.Die()
	deadX := spiny.Position.X
	spiny.Update(0.1)
	
	if spiny.Position.X != deadX {
		t.Error("Dead Spiny should not move")
	}
}

func TestSpinyInactiveState(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	// Deactivate the Spiny
	spiny.Deactivate()
	
	initialX := spiny.Position.X
	
	// Update should not move inactive Spiny
	spiny.Update(0.1)
	
	if spiny.Position.X != initialX {
		t.Error("Inactive Spiny should not move")
	}
}

func TestSpinyGetSprite(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	sprite := spiny.GetSprite()
	if sprite != 'S' {
		t.Errorf("Expected sprite 'S', got '%c'", sprite)
	}
}

func TestSpinyGetSpriteStyle(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	style := spiny.GetSpriteStyle()
	if style == (tcell.Style{}) {
		t.Error("Expected non-empty sprite style")
	}
}

func TestSpinyIsAlive(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	if !spiny.IsAlive() {
		t.Error("New Spiny should be alive")
	}
	
	spiny.Die()
	
	if spiny.IsAlive() {
		t.Error("Dead Spiny should not be alive")
	}
}

func TestSpinyActivateDeactivate(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	if !spiny.IsActive {
		t.Fatal("Spiny should start active")
	}
	
	spiny.Deactivate()
	
	if spiny.IsActive {
		t.Error("Spiny should be inactive after deactivation")
	}
	
	spiny.Activate()
	
	if !spiny.IsActive {
		t.Error("Spiny should be active after activation")
	}
}

func TestSpinyGetBounds(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	x, y, width, height := spiny.GetBounds()
	
	if x != 100 || y != 200 {
		t.Errorf("Expected bounds position (100, 200), got (%f, %f)", x, y)
	}
	
	if width != 16 || height != 16 {
		t.Errorf("Expected bounds size (16, 16), got (%f, %f)", width, height)
	}
}

func TestSpinySetPosition(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	spiny.SetPosition(150, 250)
	
	if spiny.Position.X != 150 || spiny.Position.Y != 250 {
		t.Errorf("Expected position (150, 250), got (%f, %f)", spiny.Position.X, spiny.Position.Y)
	}
}

func TestSpinySetVelocity(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	spiny.SetVelocity(50, 100)
	
	if spiny.Velocity.X != 50 || spiny.Velocity.Y != 100 {
		t.Errorf("Expected velocity (50, 100), got (%f, %f)", spiny.Velocity.X, spiny.Velocity.Y)
	}
}

func TestSpinyOnCollideWithPlayer(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	// This should not crash (player damage is handled by game engine)
	spiny.OnCollideWithPlayer(nil)
	
	// Dead Spiny should not damage player
	spiny.Die()
	spiny.OnCollideWithPlayer(nil)
}

func TestSpinyMovementSpeed(t *testing.T) {
	spiny := NewSpiny(100, 200)
	spiny.OnGround = true
	spiny.FacingRight = true
	
	// Update for 1 second
	for i := 0; i < 60; i++ {
		spiny.Update(1.0 / 60.0)
	}
	
	// Should have moved approximately 70 pixels (70 pixels/second)
	distance := spiny.Position.X - 100
	
	// Allow some tolerance for floating point calculations
	if distance < 65 || distance > 75 {
		t.Errorf("Expected Spiny to move ~70 pixels in 1 second, moved %f", distance)
	}
}

func TestSpinyType(t *testing.T) {
	spiny := NewSpiny(100, 200)
	
	if spiny.Type != EnemyTypeSpiny {
		t.Errorf("Expected type EnemyTypeSpiny, got %v", spiny.Type)
	}
}

func TestSpinyVsGoombaSpeed(t *testing.T) {
	spiny := NewSpiny(100, 200)
	goomba := NewGoomba(100, 200)
	
	// Spiny should be faster than Goomba
	if spiny.MoveSpeed <= goomba.MoveSpeed {
		t.Errorf("Spiny should be faster than Goomba. Spiny: %f, Goomba: %f", spiny.MoveSpeed, goomba.MoveSpeed)
	}
}
