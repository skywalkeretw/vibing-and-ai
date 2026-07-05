package entities

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewKoopa(t *testing.T) {
	koopa := NewKoopa(100, 200)
	
	if koopa == nil {
		t.Fatal("NewKoopa returned nil")
	}
	
	if koopa.Position.X != 100 || koopa.Position.Y != 200 {
		t.Errorf("Expected position (100, 200), got (%f, %f)", koopa.Position.X, koopa.Position.Y)
	}
	
	if koopa.Health != 2 {
		t.Errorf("Expected health 2, got %d", koopa.Health)
	}
	
	if koopa.MaxHealth != 2 {
		t.Errorf("Expected max health 2, got %d", koopa.MaxHealth)
	}
	
	if koopa.MoveSpeed != 80.0 {
		t.Errorf("Expected move speed 80.0, got %f", koopa.MoveSpeed)
	}
	
	if koopa.koopaState != KoopaStateWalking {
		t.Errorf("Expected initial state KoopaStateWalking, got %v", koopa.koopaState)
	}
	
	if koopa.Sprite != 'K' {
		t.Errorf("Expected sprite 'K', got '%c'", koopa.Sprite)
	}
	
	if !koopa.IsActive {
		t.Error("Expected Koopa to be active")
	}
	
	if koopa.IsDead {
		t.Error("Expected Koopa to be alive")
	}
}

func TestKoopaWalkingState(t *testing.T) {
	koopa := NewKoopa(100, 200)
	koopa.OnGround = true
	koopa.FacingRight = true
	
	// Update walking state
	koopa.Update(0.016) // ~60 FPS
	
	if koopa.koopaState != KoopaStateWalking {
		t.Errorf("Expected KoopaStateWalking, got %v", koopa.koopaState)
	}
	
	if koopa.Velocity.X <= 0 {
		t.Errorf("Expected positive velocity when facing right, got %f", koopa.Velocity.X)
	}
	
	// Test facing left
	koopa.FacingRight = false
	koopa.Update(0.016)
	
	if koopa.Velocity.X >= 0 {
		t.Errorf("Expected negative velocity when facing left, got %f", koopa.Velocity.X)
	}
}

func TestKoopaStompToShell(t *testing.T) {
	koopa := NewKoopa(100, 200)
	
	if koopa.koopaState != KoopaStateWalking {
		t.Fatal("Expected initial state to be walking")
	}
	
	// First stomp should convert to shell
	koopa.OnStomp(nil)
	
	if koopa.koopaState != KoopaStateShell {
		t.Errorf("Expected KoopaStateShell after first stomp, got %v", koopa.koopaState)
	}
	
	if koopa.Health != 1 {
		t.Errorf("Expected health 1 after becoming shell, got %d", koopa.Health)
	}
	
	if koopa.Sprite != 'o' {
		t.Errorf("Expected shell sprite 'o', got '%c'", koopa.Sprite)
	}
	
	if koopa.Velocity.X != 0 {
		t.Errorf("Expected shell to be stationary, got velocity %f", koopa.Velocity.X)
	}
	
	if koopa.IsDead {
		t.Error("Expected Koopa to still be alive as shell")
	}
}

func TestKoopaStompShellDeath(t *testing.T) {
	koopa := NewKoopa(100, 200)
	
	// First stomp to shell
	koopa.OnStomp(nil)
	
	if koopa.IsDead {
		t.Fatal("Koopa should not be dead after first stomp")
	}
	
	// Second stomp should kill
	koopa.OnStomp(nil)
	
	if !koopa.IsDead {
		t.Error("Expected Koopa to be dead after second stomp")
	}
	
	if koopa.State != EnemyStateDead {
		t.Errorf("Expected EnemyStateDead, got %v", koopa.State)
	}
}

func TestKoopaKickMechanics(t *testing.T) {
	koopa := NewKoopa(100, 200)
	
	// Convert to shell first
	koopa.OnStomp(nil)
	
	if koopa.koopaState != KoopaStateShell {
		t.Fatal("Expected shell state")
	}
	
	// Kick the shell
	koopa.OnKick(nil)
	
	if koopa.koopaState != KoopaStateKicked {
		t.Errorf("Expected KoopaStateKicked after kick, got %v", koopa.koopaState)
	}
	
	if koopa.kickVelocity.X == 0 {
		t.Error("Expected non-zero kick velocity")
	}
	
	// Verify kick speed is 300 pixels/second
	expectedSpeed := 300.0
	actualSpeed := koopa.kickVelocity.X
	if actualSpeed != expectedSpeed && actualSpeed != -expectedSpeed {
		t.Errorf("Expected kick speed %f or %f, got %f", expectedSpeed, -expectedSpeed, actualSpeed)
	}
}

func TestKoopaWallCollision(t *testing.T) {
	koopa := NewKoopa(100, 200)
	koopa.FacingRight = true
	
	// Test walking state wall collision
	koopa.OnCollideWithWall()
	
	if koopa.FacingRight {
		t.Error("Expected Koopa to turn around after wall collision")
	}
	
	// Test kicked shell wall collision
	koopa.OnStomp(nil) // Convert to shell
	koopa.OnKick(nil)  // Kick the shell
	
	originalVelocity := koopa.kickVelocity.X
	koopa.OnCollideWithWall()
	
	if koopa.kickVelocity.X == originalVelocity {
		t.Error("Expected kicked shell to bounce back after wall collision")
	}
	
	if koopa.kickVelocity.X != -originalVelocity {
		t.Errorf("Expected velocity to reverse, got %f (original: %f)", koopa.kickVelocity.X, originalVelocity)
	}
}

func TestKoopaProjectileDefeat(t *testing.T) {
	koopa := NewKoopa(100, 200)
	
	// Projectile should defeat Koopa in walking state
	koopa.OnHitByProjectile(nil)
	
	if !koopa.IsDead {
		t.Error("Expected Koopa to be dead after projectile hit")
	}
	
	// Test projectile on shell
	koopa2 := NewKoopa(100, 200)
	koopa2.OnStomp(nil) // Convert to shell
	koopa2.OnHitByProjectile(nil)
	
	if !koopa2.IsDead {
		t.Error("Expected shell to be defeated by projectile")
	}
}

func TestKoopaIsShell(t *testing.T) {
	koopa := NewKoopa(100, 200)
	
	if koopa.IsShell() {
		t.Error("Walking Koopa should not be considered a shell")
	}
	
	koopa.OnStomp(nil) // Convert to shell
	
	if !koopa.IsShell() {
		t.Error("Koopa in shell state should be considered a shell")
	}
	
	koopa.OnKick(nil) // Kick the shell
	
	if !koopa.IsShell() {
		t.Error("Kicked shell should still be considered a shell")
	}
}

func TestKoopaIsKicked(t *testing.T) {
	koopa := NewKoopa(100, 200)
	
	if koopa.IsKicked() {
		t.Error("Walking Koopa should not be kicked")
	}
	
	koopa.OnStomp(nil) // Convert to shell
	
	if koopa.IsKicked() {
		t.Error("Stationary shell should not be kicked")
	}
	
	koopa.OnKick(nil) // Kick the shell
	
	if !koopa.IsKicked() {
		t.Error("Shell should be in kicked state after kick")
	}
}

func TestKoopaPatrolBounds(t *testing.T) {
	koopa := NewKoopa(100, 200)
	koopa.OnGround = true
	koopa.FacingRight = true
	
	// Set patrol bounds
	koopa.SetPatrolBounds(50, 150)
	
	if koopa.PatrolLeft != 50 {
		t.Errorf("Expected patrol left 50, got %f", koopa.PatrolLeft)
	}
	
	if koopa.PatrolRight != 150 {
		t.Errorf("Expected patrol right 150, got %f", koopa.PatrolRight)
	}
	
	// Move to right boundary
	koopa.Position.X = 151
	koopa.Update(0.016)
	
	if koopa.FacingRight {
		t.Error("Expected Koopa to turn left at right boundary")
	}
	
	// Move to left boundary
	koopa.Position.X = 49
	koopa.Update(0.016)
	
	if !koopa.FacingRight {
		t.Error("Expected Koopa to turn right at left boundary")
	}
}

func TestKoopaUpdate(t *testing.T) {
	koopa := NewKoopa(100, 200)
	koopa.OnGround = true
	
	initialX := koopa.Position.X
	
	// Update should move the Koopa
	koopa.Update(0.1) // 100ms
	
	if koopa.Position.X == initialX {
		t.Error("Expected Koopa position to change after update")
	}
	
	// Dead Koopa should not update
	koopa.Die()
	deadX := koopa.Position.X
	koopa.Update(0.1)
	
	if koopa.Position.X != deadX {
		t.Error("Dead Koopa should not move")
	}
}

func TestKoopaGetters(t *testing.T) {
	koopa := NewKoopa(100, 200)
	
	sprite := koopa.GetSprite()
	if sprite != 'K' {
		t.Errorf("Expected sprite 'K', got '%c'", sprite)
	}
	
	style := koopa.GetSpriteStyle()
	if style == (tcell.Style{}) {
		t.Error("Expected non-empty sprite style")
	}
	
	state := koopa.GetKoopaState()
	if state != KoopaStateWalking {
		t.Errorf("Expected KoopaStateWalking, got %v", state)
	}
}

func TestKoopaShellUpdate(t *testing.T) {
	koopa := NewKoopa(100, 200)
	koopa.OnStomp(nil) // Convert to shell
	
	// Shell should remain stationary
	koopa.Update(0.1)
	
	if koopa.Velocity.X != 0 {
		t.Errorf("Shell should be stationary, got velocity %f", koopa.Velocity.X)
	}
}

func TestKoopaKickedUpdate(t *testing.T) {
	koopa := NewKoopa(100, 200)
	koopa.OnStomp(nil) // Convert to shell
	koopa.OnKick(nil)  // Kick the shell
	
	initialX := koopa.Position.X
	
	// Kicked shell should move
	koopa.Update(0.1)
	
	if koopa.Position.X == initialX {
		t.Error("Kicked shell should move")
	}
	
	// Velocity should match kick velocity
	if koopa.Velocity.X != koopa.kickVelocity.X {
		t.Errorf("Expected velocity %f, got %f", koopa.kickVelocity.X, koopa.Velocity.X)
	}
}
