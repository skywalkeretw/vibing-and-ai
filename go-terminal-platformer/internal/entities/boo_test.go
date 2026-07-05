package entities

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewBoo(t *testing.T) {
	boo := NewBoo(100, 200)
	
	if boo == nil {
		t.Fatal("NewBoo returned nil")
	}
	
	if boo.Position.X != 100 || boo.Position.Y != 200 {
		t.Errorf("Expected position (100, 200), got (%f, %f)", boo.Position.X, boo.Position.Y)
	}
	
	if boo.Health != 999 {
		t.Errorf("Expected health 999 (invulnerable), got %d", boo.Health)
	}
	
	if boo.MoveSpeed != 50.0 {
		t.Errorf("Expected move speed 50.0, got %f", boo.MoveSpeed)
	}
	
	if boo.Sprite != 'B' {
		t.Errorf("Expected sprite 'B', got '%c'", boo.Sprite)
	}
	
	if boo.State != EnemyStatePatrol {
		t.Errorf("Expected initial state EnemyStatePatrol, got %v", boo.State)
	}
	
	if !boo.IsActive {
		t.Error("Expected Boo to be active")
	}
	
	if boo.IsDead {
		t.Error("Expected Boo to be alive")
	}
	
	if !boo.IsFlying {
		t.Error("Expected Boo to be flying")
	}
}

func TestBooIsFlying(t *testing.T) {
	boo := NewBoo(100, 200)
	
	if !boo.IsFlyingEnemy() {
		t.Error("Boo should be a flying enemy")
	}
}

func TestBooInvulnerability(t *testing.T) {
	boo := NewBoo(100, 200)
	
	if !boo.IsInvulnerable() {
		t.Error("Boo should be invulnerable")
	}
	
	initialHealth := boo.Health
	
	// Try to damage Boo
	boo.TakeDamage(1)
	
	if boo.Health != initialHealth {
		t.Errorf("Boo should be invulnerable to normal damage, health changed from %d to %d", initialHealth, boo.Health)
	}
	
	if boo.IsDead {
		t.Error("Boo should not die from normal damage")
	}
}

func TestBooProjectileImmunity(t *testing.T) {
	boo := NewBoo(100, 200)
	
	initialHealth := boo.Health
	
	// Hit with projectile
	boo.OnHitByProjectile(nil)
	
	if boo.Health != initialHealth {
		t.Errorf("Boo should be immune to projectiles, health changed from %d to %d", initialHealth, boo.Health)
	}
	
	if boo.IsDead {
		t.Error("Boo should not die from projectile")
	}
}

func TestBooShellImmunity(t *testing.T) {
	boo := NewBoo(100, 200)
	
	initialHealth := boo.Health
	
	// Hit with shell
	boo.OnHitByShell(nil)
	
	if boo.Health != initialHealth {
		t.Errorf("Boo should be immune to shells, health changed from %d to %d", initialHealth, boo.Health)
	}
	
	if boo.IsDead {
		t.Error("Boo should not die from shell")
	}
}

func TestBooStompImmunity(t *testing.T) {
	boo := NewBoo(100, 200)
	
	initialHealth := boo.Health
	
	// Try to stomp
	boo.OnStomp(nil)
	
	if boo.Health != initialHealth {
		t.Errorf("Boo should be immune to stomps, health changed from %d to %d", initialHealth, boo.Health)
	}
	
	if boo.IsDead {
		t.Error("Boo should not die from stomp")
	}
}

func TestBooCanBeStomped(t *testing.T) {
	boo := NewBoo(100, 200)
	
	if boo.CanBeStomped() {
		t.Error("Boo should not be stompable")
	}
}

func TestBooStarDefeat(t *testing.T) {
	boo := NewBoo(100, 200)
	
	if boo.IsDead {
		t.Fatal("Boo should start alive")
	}
	
	// Hit with star power-up
	boo.OnStarContact(nil)
	
	if !boo.IsDead {
		t.Error("Boo should be dead after star contact")
	}
	
	if boo.State != EnemyStateDead {
		t.Errorf("Expected EnemyStateDead, got %v", boo.State)
	}
}

func TestBooPlayerFacingDetection(t *testing.T) {
	boo := NewBoo(150, 200)
	
	// Player at x=100, facing right (toward Boo at x=150)
	boo.UpdateWithPlayerFacing(0.016, 100, 200, true)
	
	if !boo.IsShy() {
		t.Error("Boo should be shy when player faces it (player left, facing right)")
	}
	
	if boo.Velocity.X != 0 || boo.Velocity.Y != 0 {
		t.Error("Boo should stop moving when shy")
	}
	
	// Player at x=100, facing left (away from Boo at x=150)
	boo.UpdateWithPlayerFacing(0.016, 100, 200, false)
	
	if boo.IsShy() {
		t.Error("Boo should not be shy when player faces away")
	}
	
	if boo.Velocity.X == 0 && boo.Velocity.Y == 0 {
		t.Error("Boo should be moving when not shy")
	}
}

func TestBooPlayerFacingDetectionReverse(t *testing.T) {
	boo := NewBoo(50, 200)
	
	// Player at x=100, facing left (toward Boo at x=50)
	boo.UpdateWithPlayerFacing(0.016, 100, 200, false)
	
	if !boo.IsShy() {
		t.Error("Boo should be shy when player faces it (player right, facing left)")
	}
	
	// Player at x=100, facing right (away from Boo at x=50)
	boo.UpdateWithPlayerFacing(0.016, 100, 200, true)
	
	if boo.IsShy() {
		t.Error("Boo should not be shy when player faces away")
	}
}

func TestBooChasePlayer(t *testing.T) {
	boo := NewBoo(100, 200)
	
	// Player at x=200, y=200, facing away (right, away from Boo on left)
	boo.UpdateWithPlayerFacing(0.016, 200, 200, true)
	
	if boo.IsShy() {
		t.Error("Boo should not be shy when player faces away")
	}
	
	if boo.Velocity.X <= 0 {
		t.Errorf("Boo should move right toward player, got velocity X=%f", boo.Velocity.X)
	}
	
	if boo.State != EnemyStateChase {
		t.Errorf("Expected EnemyStateChase when chasing, got %v", boo.State)
	}
}

func TestBooShyState(t *testing.T) {
	boo := NewBoo(150, 200)
	
	// Player facing Boo
	boo.UpdateWithPlayerFacing(0.016, 100, 200, true)
	
	if !boo.IsShy() {
		t.Error("Boo should be shy")
	}
	
	if boo.State != EnemyStateIdle {
		t.Errorf("Expected EnemyStateIdle when shy, got %v", boo.State)
	}
}

func TestBooSetShy(t *testing.T) {
	boo := NewBoo(100, 200)
	
	boo.SetShy(true)
	if !boo.IsShy() {
		t.Error("Expected Boo to be shy after SetShy(true)")
	}
	
	boo.SetShy(false)
	if boo.IsShy() {
		t.Error("Expected Boo to not be shy after SetShy(false)")
	}
}

func TestBooGetSprite(t *testing.T) {
	boo := NewBoo(100, 200)
	
	// Not shy
	boo.SetShy(false)
	sprite := boo.GetSprite()
	if sprite != 'B' {
		t.Errorf("Expected sprite 'B' when not shy, got '%c'", sprite)
	}
	
	// Shy
	boo.SetShy(true)
	sprite = boo.GetSprite()
	if sprite != 'b' {
		t.Errorf("Expected sprite 'b' when shy, got '%c'", sprite)
	}
}

func TestBooGetSpriteStyle(t *testing.T) {
	boo := NewBoo(100, 200)
	
	// Not shy - should be white and bold
	boo.SetShy(false)
	style := boo.GetSpriteStyle()
	if style == (tcell.Style{}) {
		t.Error("Expected non-empty sprite style when not shy")
	}
	
	// Shy - should be gray and not bold
	boo.SetShy(true)
	style = boo.GetSpriteStyle()
	if style == (tcell.Style{}) {
		t.Error("Expected non-empty sprite style when shy")
	}
}

func TestBooUpdate(t *testing.T) {
	boo := NewBoo(100, 200)
	
	// Basic update should not crash
	boo.Update(0.016)
	
	// Dead Boo should not update
	boo.Die()
	deadX := boo.Position.X
	boo.Update(0.1)
	
	if boo.Position.X != deadX {
		t.Error("Dead Boo should not move")
	}
}

func TestBooInactiveState(t *testing.T) {
	boo := NewBoo(100, 200)
	
	// Deactivate Boo
	boo.Deactivate()
	
	initialX := boo.Position.X
	
	// Update should not move inactive Boo
	boo.UpdateWithPlayerFacing(0.1, 200, 200, false)
	
	if boo.Position.X != initialX {
		t.Error("Inactive Boo should not move")
	}
}

func TestBooIsAlive(t *testing.T) {
	boo := NewBoo(100, 200)
	
	if !boo.IsAlive() {
		t.Error("New Boo should be alive")
	}
	
	boo.Die()
	
	if boo.IsAlive() {
		t.Error("Dead Boo should not be alive")
	}
}

func TestBooActivateDeactivate(t *testing.T) {
	boo := NewBoo(100, 200)
	
	if !boo.IsActive {
		t.Fatal("Boo should start active")
	}
	
	boo.Deactivate()
	
	if boo.IsActive {
		t.Error("Boo should be inactive after deactivation")
	}
	
	boo.Activate()
	
	if !boo.IsActive {
		t.Error("Boo should be active after activation")
	}
}

func TestBooGetBounds(t *testing.T) {
	boo := NewBoo(100, 200)
	
	x, y, width, height := boo.GetBounds()
	
	if x != 100 || y != 200 {
		t.Errorf("Expected bounds position (100, 200), got (%f, %f)", x, y)
	}
	
	if width != 16 || height != 16 {
		t.Errorf("Expected bounds size (16, 16), got (%f, %f)", width, height)
	}
}

func TestBooSetPosition(t *testing.T) {
	boo := NewBoo(100, 200)
	
	boo.SetPosition(150, 250)
	
	if boo.Position.X != 150 || boo.Position.Y != 250 {
		t.Errorf("Expected position (150, 250), got (%f, %f)", boo.Position.X, boo.Position.Y)
	}
}

func TestBooSetVelocity(t *testing.T) {
	boo := NewBoo(100, 200)
	
	boo.SetVelocity(50, 100)
	
	if boo.Velocity.X != 50 || boo.Velocity.Y != 100 {
		t.Errorf("Expected velocity (50, 100), got (%f, %f)", boo.Velocity.X, boo.Velocity.Y)
	}
}

func TestBooOnCollideWithPlayer(t *testing.T) {
	boo := NewBoo(100, 200)
	
	// This should not crash (damage handled by game engine)
	boo.OnCollideWithPlayer(nil)
	
	// Dead Boo should not damage player
	boo.Die()
	boo.OnCollideWithPlayer(nil)
}

func TestBooSetTargetPlayer(t *testing.T) {
	boo := NewBoo(100, 200)
	
	// This should not crash
	boo.SetTargetPlayer(nil)
	boo.SetTargetPlayer("mock_player")
}

func TestBooType(t *testing.T) {
	boo := NewBoo(100, 200)
	
	if boo.Type != EnemyTypeBoo {
		t.Errorf("Expected type EnemyTypeBoo, got %v", boo.Type)
	}
}

func TestBooChaseSpeed(t *testing.T) {
	boo := NewBoo(100, 200)
	
	// Player far away to the right, facing away (right)
	boo.UpdateWithPlayerFacing(0.016, 300, 200, true)
	
	// Calculate speed magnitude
	speed := boo.Velocity.X
	
	// Should be moving at approximately 50 pixels/second
	if speed < 45 || speed > 55 {
		t.Errorf("Expected chase speed around 50 pixels/second, got %f", speed)
	}
}

func TestBooFacingDirection(t *testing.T) {
	boo := NewBoo(100, 200)
	
	// Player to the right, facing away (right)
	boo.UpdateWithPlayerFacing(0.016, 200, 200, true)
	
	if !boo.FacingRight {
		t.Error("Boo should face right when chasing player to the right")
	}
	
	// Player to the left, facing away (left)
	boo.Position.X = 200
	boo.UpdateWithPlayerFacing(0.016, 100, 200, false)
	
	if boo.FacingRight {
		t.Error("Boo should face left when chasing player to the left")
	}
}

func TestBooVerticalChase(t *testing.T) {
	boo := NewBoo(100, 100)
	
	// Player below and to the right, facing away (right)
	boo.UpdateWithPlayerFacing(0.016, 150, 200, true)
	
	if boo.Velocity.Y <= 0 {
		t.Errorf("Boo should move down toward player below, got velocity Y=%f", boo.Velocity.Y)
	}
	
	// Player above and to the right, facing away (right)
	boo.Position.Y = 200
	boo.UpdateWithPlayerFacing(0.016, 150, 100, true)
	
	if boo.Velocity.Y >= 0 {
		t.Errorf("Boo should move up toward player above, got velocity Y=%f", boo.Velocity.Y)
	}
}

func TestBooDiagonalChase(t *testing.T) {
	boo := NewBoo(100, 100)
	
	// Player at diagonal position (right and down), facing away (right)
	boo.UpdateWithPlayerFacing(0.016, 200, 200, true)
	
	if boo.Velocity.X <= 0 {
		t.Error("Boo should move right toward diagonal player")
	}
	
	if boo.Velocity.Y <= 0 {
		t.Error("Boo should move down toward diagonal player")
	}
}

func TestBooShyStopsMovement(t *testing.T) {
	boo := NewBoo(100, 200)
	
	// Start chasing (player to right, facing right/away)
	boo.UpdateWithPlayerFacing(0.016, 200, 200, true)
	
	if boo.Velocity.X == 0 {
		t.Fatal("Boo should be moving when chasing")
	}
	
	// Player turns to face Boo (now facing left toward Boo)
	boo.UpdateWithPlayerFacing(0.016, 200, 200, false)
	
	if boo.Velocity.X != 0 || boo.Velocity.Y != 0 {
		t.Error("Boo should stop completely when shy")
	}
}

func TestBooMultipleUpdates(t *testing.T) {
	boo := NewBoo(100, 200)
	
	initialX := boo.Position.X
	
	// Chase for 1 second (60 frames), player to right facing away
	for i := 0; i < 60; i++ {
		boo.UpdateWithPlayerFacing(1.0/60.0, 200, 200, true)
	}
	
	// Should have moved approximately 50 pixels (50 pixels/second)
	distance := boo.Position.X - initialX
	
	// Allow some tolerance
	if distance < 45 || distance > 55 {
		t.Errorf("Expected Boo to move ~50 pixels in 1 second, moved %f", distance)
	}
}
