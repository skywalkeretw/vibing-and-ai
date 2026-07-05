package entities

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewLakitu(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	if lakitu == nil {
		t.Fatal("NewLakitu returned nil")
	}
	
	if lakitu.Position.X != 100 || lakitu.Position.Y != 200 {
		t.Errorf("Expected position (100, 200), got (%f, %f)", lakitu.Position.X, lakitu.Position.Y)
	}
	
	if lakitu.Health != 3 {
		t.Errorf("Expected health 3, got %d", lakitu.Health)
	}
	
	if lakitu.MaxHealth != 3 {
		t.Errorf("Expected max health 3, got %d", lakitu.MaxHealth)
	}
	
	if lakitu.MoveSpeed != 100.0 {
		t.Errorf("Expected move speed 100.0, got %f", lakitu.MoveSpeed)
	}
	
	if lakitu.Sprite != 'L' {
		t.Errorf("Expected sprite 'L', got '%c'", lakitu.Sprite)
	}
	
	if lakitu.State != EnemyStatePatrol {
		t.Errorf("Expected initial state EnemyStatePatrol, got %v", lakitu.State)
	}
	
	if !lakitu.IsActive {
		t.Error("Expected Lakitu to be active")
	}
	
	if lakitu.IsDead {
		t.Error("Expected Lakitu to be alive")
	}
	
	if !lakitu.IsFlying {
		t.Error("Expected Lakitu to be flying")
	}
}

func TestLakituIsFlying(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	if !lakitu.IsFlying {
		t.Error("Lakitu should be a flying enemy")
	}
}

func TestLakituStompImmunity(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	initialHealth := lakitu.Health
	
	// Try to stomp Lakitu
	lakitu.OnStomp(nil)
	
	// Health should not change (stomp immunity)
	if lakitu.Health != initialHealth {
		t.Errorf("Lakitu should be immune to stomps, health changed from %d to %d", initialHealth, lakitu.Health)
	}
	
	if lakitu.IsDead {
		t.Error("Lakitu should not die from stomp")
	}
}

func TestLakituCanBeStomped(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	if lakitu.CanBeStomped() {
		t.Error("Lakitu should not be stompable (flying enemy)")
	}
}

func TestLakituProjectileDamage(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	if lakitu.Health != 3 {
		t.Fatalf("Expected initial health 3, got %d", lakitu.Health)
	}
	
	// First hit
	lakitu.OnHitByProjectile(nil)
	if lakitu.Health != 2 {
		t.Errorf("Expected health 2 after first hit, got %d", lakitu.Health)
	}
	if lakitu.IsDead {
		t.Error("Lakitu should not be dead after first hit")
	}
	
	// Wait for invulnerability to expire
	lakitu.Update(0.6)
	
	// Second hit
	lakitu.OnHitByProjectile(nil)
	if lakitu.Health != 1 {
		t.Errorf("Expected health 1 after second hit, got %d", lakitu.Health)
	}
	if lakitu.IsDead {
		t.Error("Lakitu should not be dead after second hit")
	}
	
	// Wait for invulnerability to expire
	lakitu.Update(0.6)
	
	// Third hit (should kill)
	lakitu.OnHitByProjectile(nil)
	if lakitu.Health != 0 {
		t.Errorf("Expected health 0 after third hit, got %d", lakitu.Health)
	}
	if !lakitu.IsDead {
		t.Error("Lakitu should be dead after third hit")
	}
}

func TestLakituShellDamage(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	if lakitu.Health != 3 {
		t.Fatalf("Expected initial health 3, got %d", lakitu.Health)
	}
	
	// First hit
	lakitu.OnHitByShell(nil)
	if lakitu.Health != 2 {
		t.Errorf("Expected health 2 after first hit, got %d", lakitu.Health)
	}
	
	// Wait for invulnerability to expire
	lakitu.Update(0.6)
	
	// Second hit
	lakitu.OnHitByShell(nil)
	if lakitu.Health != 1 {
		t.Errorf("Expected health 1 after second hit, got %d", lakitu.Health)
	}
	
	// Wait for invulnerability to expire
	lakitu.Update(0.6)
	
	// Third hit (should kill)
	lakitu.OnHitByShell(nil)
	if !lakitu.IsDead {
		t.Error("Lakitu should be dead after third shell hit")
	}
}

func TestLakituSpinyDropTimer(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	if lakitu.GetDropTimer() != 0.0 {
		t.Errorf("Expected initial drop timer 0.0, got %f", lakitu.GetDropTimer())
	}
	
	// Update for 2.5 seconds
	lakitu.Update(2.5)
	
	if lakitu.GetDropTimer() != 2.5 {
		t.Errorf("Expected drop timer 2.5, got %f", lakitu.GetDropTimer())
	}
	
	// Update for another 2.5 seconds (total 5.0, should trigger drop and reset)
	lakitu.Update(2.5)
	
	if lakitu.GetDropTimer() != 0.0 {
		t.Errorf("Expected drop timer to reset to 0.0 after 5 seconds, got %f", lakitu.GetDropTimer())
	}
}

func TestLakituSpinySpawning(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	spawnCount := 0
	var lastSpawnX, lastSpawnY float64
	
	// Set up spawner callback
	lakitu.SetSpinySpawner(func(x, y float64) {
		spawnCount++
		lastSpawnX = x
		lastSpawnY = y
	})
	
	// Update for 5 seconds to trigger spawn
	lakitu.Update(5.0)
	
	if spawnCount != 1 {
		t.Errorf("Expected 1 Spiny spawn, got %d", spawnCount)
	}
	
	if lastSpawnX != 100 {
		t.Errorf("Expected spawn X at 100, got %f", lastSpawnX)
	}
	
	if lastSpawnY != 220 {
		t.Errorf("Expected spawn Y at 220 (20 pixels below Lakitu), got %f", lastSpawnY)
	}
	
	// Update for another 5 seconds
	lakitu.Update(5.0)
	
	if spawnCount != 2 {
		t.Errorf("Expected 2 Spiny spawns total, got %d", spawnCount)
	}
}

func TestLakituDropSpinyManual(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	spawnCount := 0
	
	lakitu.SetSpinySpawner(func(x, y float64) {
		spawnCount++
	})
	
	// Manually trigger drop
	lakitu.DropSpiny()
	
	if spawnCount != 1 {
		t.Errorf("Expected 1 Spiny spawn, got %d", spawnCount)
	}
}

func TestLakituFollowPlayer(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	// Player far to the right
	lakitu.FollowPlayer(200, 100)
	
	if lakitu.Velocity.X <= 0 {
		t.Errorf("Expected positive velocity when player is to the right, got %f", lakitu.Velocity.X)
	}
	
	if !lakitu.FacingRight {
		t.Error("Expected Lakitu to face right when following player to the right")
	}
	
	// Player far to the left
	lakitu.FollowPlayer(0, 100)
	
	if lakitu.Velocity.X >= 0 {
		t.Errorf("Expected negative velocity when player is to the left, got %f", lakitu.Velocity.X)
	}
	
	if lakitu.FacingRight {
		t.Error("Expected Lakitu to face left when following player to the left")
	}
	
	// Player close (within 20 pixels)
	lakitu.Position.X = 100
	lakitu.FollowPlayer(110, 100)
	
	if lakitu.Velocity.X != 0 {
		t.Errorf("Expected zero velocity when player is close, got %f", lakitu.Velocity.X)
	}
}

func TestLakituVerticalPositioning(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	// Player below Lakitu
	lakitu.FollowPlayer(100, 300)
	
	// Lakitu should move down to maintain height above player
	if lakitu.Velocity.Y <= 0 {
		t.Errorf("Expected positive Y velocity when player is below, got %f", lakitu.Velocity.Y)
	}
	
	// Player above Lakitu
	lakitu.Position.Y = 200
	lakitu.FollowPlayer(100, 100)
	
	// Lakitu should move up to maintain height above player
	if lakitu.Velocity.Y >= 0 {
		t.Errorf("Expected negative Y velocity when player is above, got %f", lakitu.Velocity.Y)
	}
}

func TestLakituSetHeightRange(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	lakitu.SetHeightRange(40, 60)
	
	if lakitu.minHeight != 40 {
		t.Errorf("Expected min height 40, got %f", lakitu.minHeight)
	}
	
	if lakitu.maxHeight != 60 {
		t.Errorf("Expected max height 60, got %f", lakitu.maxHeight)
	}
	
	expectedTarget := (40.0 + 60.0) / 2.0
	if lakitu.targetHeight != expectedTarget {
		t.Errorf("Expected target height %f, got %f", expectedTarget, lakitu.targetHeight)
	}
}

func TestLakituResetDropTimer(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	lakitu.Update(3.0)
	
	if lakitu.GetDropTimer() != 3.0 {
		t.Fatalf("Expected drop timer 3.0, got %f", lakitu.GetDropTimer())
	}
	
	lakitu.ResetDropTimer()
	
	if lakitu.GetDropTimer() != 0.0 {
		t.Errorf("Expected drop timer to reset to 0.0, got %f", lakitu.GetDropTimer())
	}
}

func TestLakituUpdate(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	initialTimer := lakitu.GetDropTimer()
	
	// Update should increment drop timer
	lakitu.Update(1.0)
	
	if lakitu.GetDropTimer() <= initialTimer {
		t.Error("Expected drop timer to increase after update")
	}
	
	// Dead Lakitu should not update
	lakitu.Die()
	deadTimer := lakitu.GetDropTimer()
	lakitu.Update(1.0)
	
	if lakitu.GetDropTimer() != deadTimer {
		t.Error("Dead Lakitu should not update drop timer")
	}
}

func TestLakituInactiveState(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	// Deactivate Lakitu
	lakitu.Deactivate()
	
	initialTimer := lakitu.GetDropTimer()
	
	// Update should not affect inactive Lakitu
	lakitu.Update(1.0)
	
	if lakitu.GetDropTimer() != initialTimer {
		t.Error("Inactive Lakitu should not update")
	}
}

func TestLakituGetSprite(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	sprite := lakitu.GetSprite()
	if sprite != 'L' {
		t.Errorf("Expected sprite 'L', got '%c'", sprite)
	}
}

func TestLakituGetSpriteStyle(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	style := lakitu.GetSpriteStyle()
	if style == (tcell.Style{}) {
		t.Error("Expected non-empty sprite style")
	}
}

func TestLakituIsAlive(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	if !lakitu.IsAlive() {
		t.Error("New Lakitu should be alive")
	}
	
	lakitu.Die()
	
	if lakitu.IsAlive() {
		t.Error("Dead Lakitu should not be alive")
	}
}

func TestLakituActivateDeactivate(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	if !lakitu.IsActive {
		t.Fatal("Lakitu should start active")
	}
	
	lakitu.Deactivate()
	
	if lakitu.IsActive {
		t.Error("Lakitu should be inactive after deactivation")
	}
	
	lakitu.Activate()
	
	if !lakitu.IsActive {
		t.Error("Lakitu should be active after activation")
	}
}

func TestLakituGetBounds(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	x, y, width, height := lakitu.GetBounds()
	
	if x != 100 || y != 200 {
		t.Errorf("Expected bounds position (100, 200), got (%f, %f)", x, y)
	}
	
	if width != 16 || height != 16 {
		t.Errorf("Expected bounds size (16, 16), got (%f, %f)", width, height)
	}
}

func TestLakituSetPosition(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	lakitu.SetPosition(150, 250)
	
	if lakitu.Position.X != 150 || lakitu.Position.Y != 250 {
		t.Errorf("Expected position (150, 250), got (%f, %f)", lakitu.Position.X, lakitu.Position.Y)
	}
}

func TestLakituSetVelocity(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	lakitu.SetVelocity(50, 100)
	
	if lakitu.Velocity.X != 50 || lakitu.Velocity.Y != 100 {
		t.Errorf("Expected velocity (50, 100), got (%f, %f)", lakitu.Velocity.X, lakitu.Velocity.Y)
	}
}

func TestLakituOnCollideWithPlayer(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	// This should not crash (Lakitu doesn't directly damage player)
	lakitu.OnCollideWithPlayer(nil)
}

func TestLakituOnCollideWithWall(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	lakitu.FacingRight = true
	
	// Hit a wall
	lakitu.OnCollideWithWall()
	
	if lakitu.FacingRight {
		t.Error("Lakitu should turn around after hitting a wall")
	}
}

func TestLakituSetTargetPlayer(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	// This should not crash
	lakitu.SetTargetPlayer(nil)
	lakitu.SetTargetPlayer("mock_player")
}

func TestLakituType(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	if lakitu.Type != EnemyTypeLakitu {
		t.Errorf("Expected type EnemyTypeLakitu, got %v", lakitu.Type)
	}
}

func TestLakituMultipleSpinyDrops(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	spawnCount := 0
	lakitu.SetSpinySpawner(func(x, y float64) {
		spawnCount++
	})
	
	// Simulate 15 seconds (should spawn 3 Spinies)
	for i := 0; i < 15; i++ {
		lakitu.Update(1.0)
	}
	
	if spawnCount != 3 {
		t.Errorf("Expected 3 Spiny spawns in 15 seconds, got %d", spawnCount)
	}
}

func TestLakituNoSpawnerSet(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	// Should not crash when spawner is not set
	lakitu.DropSpiny()
	lakitu.Update(5.0)
}

func TestLakituDamageDoesNotAffectDropTimer(t *testing.T) {
	lakitu := NewLakitu(100, 200)
	
	lakitu.Update(2.0)
	timerBeforeDamage := lakitu.GetDropTimer()
	
	// Take damage
	lakitu.OnHitByProjectile(nil)
	
	// Timer should not be affected by damage
	if lakitu.GetDropTimer() != timerBeforeDamage {
		t.Error("Drop timer should not be affected by taking damage")
	}
}
