package entities

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewGoomba(t *testing.T) {
	goomba := NewGoomba(100, 200)
	
	if goomba == nil {
		t.Fatal("NewGoomba returned nil")
	}
	
	if goomba.Position.X != 100 || goomba.Position.Y != 200 {
		t.Errorf("Expected position (100, 200), got (%f, %f)", goomba.Position.X, goomba.Position.Y)
	}
	
	if goomba.Health != 1 {
		t.Errorf("Expected health 1, got %d", goomba.Health)
	}
	
	if goomba.MaxHealth != 1 {
		t.Errorf("Expected max health 1, got %d", goomba.MaxHealth)
	}
	
	if goomba.MoveSpeed != 60.0 {
		t.Errorf("Expected move speed 60.0, got %f", goomba.MoveSpeed)
	}
	
	if goomba.Sprite != 'G' {
		t.Errorf("Expected sprite 'G', got '%c'", goomba.Sprite)
	}
	
	if goomba.State != EnemyStatePatrol {
		t.Errorf("Expected initial state EnemyStatePatrol, got %v", goomba.State)
	}
	
	if !goomba.IsActive {
		t.Error("Expected Goomba to be active")
	}
	
	if goomba.IsDead {
		t.Error("Expected Goomba to be alive")
	}
}

func TestGoombaPatrolBehavior(t *testing.T) {
	goomba := NewGoomba(100, 200)
	goomba.OnGround = true
	goomba.FacingRight = true
	
	// Update patrol behavior
	goomba.Update(0.016) // ~60 FPS
	
	if goomba.Velocity.X <= 0 {
		t.Errorf("Expected positive velocity when facing right, got %f", goomba.Velocity.X)
	}
	
	// Test facing left
	goomba.FacingRight = false
	goomba.Update(0.016)
	
	if goomba.Velocity.X >= 0 {
		t.Errorf("Expected negative velocity when facing left, got %f", goomba.Velocity.X)
	}
}

func TestGoombaStompDefeat(t *testing.T) {
	goomba := NewGoomba(100, 200)
	
	if goomba.IsDead {
		t.Fatal("Goomba should start alive")
	}
	
	// Stomp the Goomba
	goomba.OnStomp(nil)
	
	if !goomba.IsDead {
		t.Error("Goomba should be dead after stomp")
	}
	
	if goomba.State != EnemyStateDead {
		t.Errorf("Expected EnemyStateDead, got %v", goomba.State)
	}
}

func TestGoombaProjectileDefeat(t *testing.T) {
	goomba := NewGoomba(100, 200)
	
	if goomba.IsDead {
		t.Fatal("Goomba should start alive")
	}
	
	// Hit with projectile
	goomba.OnHitByProjectile(nil)
	
	if !goomba.IsDead {
		t.Error("Goomba should be dead after projectile hit")
	}
	
	if goomba.State != EnemyStateDead {
		t.Errorf("Expected EnemyStateDead, got %v", goomba.State)
	}
}

func TestGoombaWallCollision(t *testing.T) {
	goomba := NewGoomba(100, 200)
	goomba.FacingRight = true
	
	// Hit a wall
	goomba.OnCollideWithWall()
	
	if goomba.FacingRight {
		t.Error("Goomba should turn around after hitting a wall")
	}
	
	// Hit another wall
	goomba.OnCollideWithWall()
	
	if !goomba.FacingRight {
		t.Error("Goomba should turn around again after hitting another wall")
	}
}

func TestGoombaEdgeDetection(t *testing.T) {
	goomba := NewGoomba(100, 200)
	goomba.FacingRight = true
	
	// Reach an edge
	goomba.OnReachEdge()
	
	if goomba.FacingRight {
		t.Error("Goomba should turn around at edge")
	}
	
	// Reach another edge
	goomba.OnReachEdge()
	
	if !goomba.FacingRight {
		t.Error("Goomba should turn around at edge again")
	}
}

func TestGoombaPatrolBounds(t *testing.T) {
	goomba := NewGoomba(100, 200)
	goomba.OnGround = true
	goomba.FacingRight = true
	
	// Set patrol bounds
	goomba.SetPatrolBounds(50, 150)
	
	if goomba.PatrolLeft != 50 {
		t.Errorf("Expected patrol left 50, got %f", goomba.PatrolLeft)
	}
	
	if goomba.PatrolRight != 150 {
		t.Errorf("Expected patrol right 150, got %f", goomba.PatrolRight)
	}
	
	// Move to right boundary
	goomba.Position.X = 151
	goomba.Update(0.016)
	
	if goomba.FacingRight {
		t.Error("Expected Goomba to turn left at right boundary")
	}
	
	// Move to left boundary
	goomba.Position.X = 49
	goomba.Update(0.016)
	
	if !goomba.FacingRight {
		t.Error("Expected Goomba to turn right at left boundary")
	}
}

func TestGoombaUpdate(t *testing.T) {
	goomba := NewGoomba(100, 200)
	goomba.OnGround = true
	
	initialX := goomba.Position.X
	
	// Update should move the Goomba
	goomba.Update(0.1) // 100ms
	
	if goomba.Position.X == initialX {
		t.Error("Expected Goomba position to change after update")
	}
	
	// Dead Goomba should not update
	goomba.Die()
	deadX := goomba.Position.X
	goomba.Update(0.1)
	
	if goomba.Position.X != deadX {
		t.Error("Dead Goomba should not move")
	}
}

func TestGoombaInactiveState(t *testing.T) {
	goomba := NewGoomba(100, 200)
	
	// Deactivate the Goomba
	goomba.Deactivate()
	
	initialX := goomba.Position.X
	
	// Update should not move inactive Goomba
	goomba.Update(0.1)
	
	if goomba.Position.X != initialX {
		t.Error("Inactive Goomba should not move")
	}
}

func TestGoombaGetSprite(t *testing.T) {
	goomba := NewGoomba(100, 200)
	
	sprite := goomba.GetSprite()
	if sprite != 'G' {
		t.Errorf("Expected sprite 'G', got '%c'", sprite)
	}
}

func TestGoombaGetSpriteStyle(t *testing.T) {
	goomba := NewGoomba(100, 200)
	
	style := goomba.GetSpriteStyle()
	if style == (tcell.Style{}) {
		t.Error("Expected non-empty sprite style")
	}
}

func TestGoombaIsAlive(t *testing.T) {
	goomba := NewGoomba(100, 200)
	
	if !goomba.IsAlive() {
		t.Error("New Goomba should be alive")
	}
	
	goomba.Die()
	
	if goomba.IsAlive() {
		t.Error("Dead Goomba should not be alive")
	}
}

func TestGoombaActivateDeactivate(t *testing.T) {
	goomba := NewGoomba(100, 200)
	
	if !goomba.IsActive {
		t.Fatal("Goomba should start active")
	}
	
	goomba.Deactivate()
	
	if goomba.IsActive {
		t.Error("Goomba should be inactive after deactivation")
	}
	
	goomba.Activate()
	
	if !goomba.IsActive {
		t.Error("Goomba should be active after activation")
	}
}

func TestGoombaGetBounds(t *testing.T) {
	goomba := NewGoomba(100, 200)
	
	x, y, width, height := goomba.GetBounds()
	
	if x != 100 || y != 200 {
		t.Errorf("Expected bounds position (100, 200), got (%f, %f)", x, y)
	}
	
	if width != 16 || height != 16 {
		t.Errorf("Expected bounds size (16, 16), got (%f, %f)", width, height)
	}
}

func TestGoombaSetPosition(t *testing.T) {
	goomba := NewGoomba(100, 200)
	
	goomba.SetPosition(150, 250)
	
	if goomba.Position.X != 150 || goomba.Position.Y != 250 {
		t.Errorf("Expected position (150, 250), got (%f, %f)", goomba.Position.X, goomba.Position.Y)
	}
}

func TestGoombaSetVelocity(t *testing.T) {
	goomba := NewGoomba(100, 200)
	
	goomba.SetVelocity(50, 100)
	
	if goomba.Velocity.X != 50 || goomba.Velocity.Y != 100 {
		t.Errorf("Expected velocity (50, 100), got (%f, %f)", goomba.Velocity.X, goomba.Velocity.Y)
	}
}

func TestGoombaOnCollideWithPlayer(t *testing.T) {
	goomba := NewGoomba(100, 200)
	
	// This should not crash (player damage is handled by game engine)
	goomba.OnCollideWithPlayer(nil)
	
	// Dead Goomba should not damage player
	goomba.Die()
	goomba.OnCollideWithPlayer(nil)
}

func TestGoombaMovementSpeed(t *testing.T) {
	goomba := NewGoomba(100, 200)
	goomba.OnGround = true
	goomba.FacingRight = true
	
	// Update for 1 second
	for i := 0; i < 60; i++ {
		goomba.Update(1.0 / 60.0)
	}
	
	// Should have moved approximately 60 pixels (60 pixels/second)
	distance := goomba.Position.X - 100
	
	// Allow some tolerance for floating point calculations
	if distance < 55 || distance > 65 {
		t.Errorf("Expected Goomba to move ~60 pixels in 1 second, moved %f", distance)
	}
}

func TestGoombaType(t *testing.T) {
	goomba := NewGoomba(100, 200)
	
	if goomba.Type != EnemyTypeGoomba {
		t.Errorf("Expected type EnemyTypeGoomba, got %v", goomba.Type)
	}
}
