package entities

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewHammerBro(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	if hammerBro == nil {
		t.Fatal("NewHammerBro returned nil")
	}
	
	if hammerBro.Position.X != 100 || hammerBro.Position.Y != 200 {
		t.Errorf("Expected position (100, 200), got (%f, %f)", hammerBro.Position.X, hammerBro.Position.Y)
	}
	
	if hammerBro.Health != 2 {
		t.Errorf("Expected health 2, got %d", hammerBro.Health)
	}
	
	if hammerBro.MaxHealth != 2 {
		t.Errorf("Expected max health 2, got %d", hammerBro.MaxHealth)
	}
	
	if hammerBro.MoveSpeed != 40.0 {
		t.Errorf("Expected move speed 40.0, got %f", hammerBro.MoveSpeed)
	}
	
	if hammerBro.Sprite != 'H' {
		t.Errorf("Expected sprite 'H', got '%c'", hammerBro.Sprite)
	}
	
	if hammerBro.throwCooldown != 2.0 {
		t.Errorf("Expected throw cooldown 2.0, got %f", hammerBro.throwCooldown)
	}
	
	if hammerBro.State != EnemyStateIdle {
		t.Errorf("Expected initial state EnemyStateIdle, got %v", hammerBro.State)
	}
	
	if !hammerBro.IsActive {
		t.Error("Expected Hammer Bro to be active")
	}
	
	if hammerBro.IsDead {
		t.Error("Expected Hammer Bro to be alive")
	}
}

func TestHammerBroTwoHitDefeat(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	if hammerBro.Health != 2 {
		t.Fatalf("Expected initial health 2, got %d", hammerBro.Health)
	}
	
	// First hit
	hammerBro.OnStomp(nil)
	if hammerBro.Health != 1 {
		t.Errorf("Expected health 1 after first stomp, got %d", hammerBro.Health)
	}
	if hammerBro.IsDead {
		t.Error("Hammer Bro should not be dead after first hit")
	}
	
	// Wait for invulnerability to expire
	hammerBro.Update(0.6)
	
	// Second hit (should kill)
	hammerBro.OnStomp(nil)
	if hammerBro.Health != 0 {
		t.Errorf("Expected health 0 after second stomp, got %d", hammerBro.Health)
	}
	if !hammerBro.IsDead {
		t.Error("Hammer Bro should be dead after second hit")
	}
}

func TestHammerBroProjectileDamage(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	// First hit
	hammerBro.OnHitByProjectile(nil)
	if hammerBro.Health != 1 {
		t.Errorf("Expected health 1 after projectile hit, got %d", hammerBro.Health)
	}
	
	// Wait for invulnerability
	hammerBro.Update(0.6)
	
	// Second hit
	hammerBro.OnHitByProjectile(nil)
	if !hammerBro.IsDead {
		t.Error("Hammer Bro should be dead after two projectile hits")
	}
}

func TestHammerBroShellDamage(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	// First hit
	hammerBro.OnHitByShell(nil)
	if hammerBro.Health != 1 {
		t.Errorf("Expected health 1 after shell hit, got %d", hammerBro.Health)
	}
	
	// Wait for invulnerability
	hammerBro.Update(0.6)
	
	// Second hit
	hammerBro.OnHitByShell(nil)
	if !hammerBro.IsDead {
		t.Error("Hammer Bro should be dead after two shell hits")
	}
}

func TestHammerBroCanBeStomped(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	if !hammerBro.CanBeStomped() {
		t.Error("Hammer Bro should be stompable")
	}
}

func TestHammerBroThrowTimer(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	if hammerBro.GetThrowTimer() != 0.0 {
		t.Errorf("Expected initial throw timer 0.0, got %f", hammerBro.GetThrowTimer())
	}
	
	// Update for 1 second
	hammerBro.Update(1.0)
	
	if hammerBro.GetThrowTimer() != 1.0 {
		t.Errorf("Expected throw timer 1.0, got %f", hammerBro.GetThrowTimer())
	}
	
	// Update for another 1 second (total 2.0, should trigger throw and reset)
	hammerBro.Update(1.0)
	
	if hammerBro.GetThrowTimer() != 0.0 {
		t.Errorf("Expected throw timer to reset to 0.0 after 2 seconds, got %f", hammerBro.GetThrowTimer())
	}
}

func TestHammerBroResetThrowTimer(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	hammerBro.Update(1.5)
	
	if hammerBro.GetThrowTimer() != 1.5 {
		t.Fatalf("Expected throw timer 1.5, got %f", hammerBro.GetThrowTimer())
	}
	
	hammerBro.ResetThrowTimer()
	
	if hammerBro.GetThrowTimer() != 0.0 {
		t.Errorf("Expected throw timer to reset to 0.0, got %f", hammerBro.GetThrowTimer())
	}
}

func TestHammerBroThrowHammer(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	throwCount := 0
	var lastX, lastY, lastVX, lastVY float64
	
	// Set up spawner callback
	hammerBro.SetHammerSpawner(func(x, y, vx, vy float64) {
		throwCount++
		lastX = x
		lastY = y
		lastVX = vx
		lastVY = vy
	})
	
	// Manually trigger throw
	hammerBro.ThrowHammer()
	
	if throwCount != 1 {
		t.Errorf("Expected 1 hammer throw, got %d", throwCount)
	}
	
	// Check that hammer spawns near Hammer Bro
	if lastX < 90 || lastX > 120 {
		t.Errorf("Expected hammer spawn X near 100, got %f", lastX)
	}
	
	if lastY < 190 || lastY > 200 {
		t.Errorf("Expected hammer spawn Y near 195, got %f", lastY)
	}
	
	// Check horizontal velocity
	if lastVX == 0 {
		t.Error("Expected non-zero X velocity for hammer")
	}
	
	// Check upward arc
	if lastVY >= 0 {
		t.Errorf("Expected negative Y velocity (upward), got %f", lastVY)
	}
}

func TestHammerBroPeriodicThrow(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	throwCount := 0
	hammerBro.SetHammerSpawner(func(x, y, vx, vy float64) {
		throwCount++
	})
	
	// Update for 2 seconds to trigger throw
	hammerBro.Update(2.0)
	
	if throwCount != 1 {
		t.Errorf("Expected 1 hammer throw after 2 seconds, got %d", throwCount)
	}
	
	// Update for another 2 seconds
	hammerBro.Update(2.0)
	
	if throwCount != 2 {
		t.Errorf("Expected 2 hammer throws total, got %d", throwCount)
	}
}

func TestHammerBroThrowHammerAtPlayer(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	hammerBro.FacingRight = true
	
	lastVX := 0.0
	hammerBro.SetHammerSpawner(func(x, y, vx, vy float64) {
		lastVX = vx
	})
	
	// Player to the right
	hammerBro.ThrowHammerAtPlayer(200, 200)
	
	if lastVX <= 0 {
		t.Errorf("Expected positive X velocity when player is to the right, got %f", lastVX)
	}
	
	if !hammerBro.FacingRight {
		t.Error("Hammer Bro should face right when throwing at player to the right")
	}
	
	// Player to the left
	hammerBro.ThrowHammerAtPlayer(50, 200)
	
	if lastVX >= 0 {
		t.Errorf("Expected negative X velocity when player is to the left, got %f", lastVX)
	}
	
	if hammerBro.FacingRight {
		t.Error("Hammer Bro should face left when throwing at player to the left")
	}
}

func TestHammerBroUpdateWithPlayerPosition(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	throwCount := 0
	hammerBro.SetHammerSpawner(func(x, y, vx, vy float64) {
		throwCount++
	})
	
	// Update with player position for 2 seconds
	hammerBro.UpdateWithPlayerPosition(2.0, 200, 200)
	
	if throwCount != 1 {
		t.Errorf("Expected 1 hammer throw after 2 seconds, got %d", throwCount)
	}
}

func TestHammerBroIsStanding(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	if !hammerBro.IsStanding() {
		t.Error("Hammer Bro should start standing")
	}
}

func TestHammerBroOnCollideWithWall(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	hammerBro.Velocity.X = 40.0
	
	// Hit a wall
	hammerBro.OnCollideWithWall()
	
	if hammerBro.Velocity.X != 0 {
		t.Errorf("Expected velocity X to be 0 after wall collision, got %f", hammerBro.Velocity.X)
	}
	
	if !hammerBro.IsStanding() {
		t.Error("Hammer Bro should be standing after wall collision")
	}
}

func TestHammerBroUpdate(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	initialTimer := hammerBro.GetThrowTimer()
	
	// Update should increment throw timer
	hammerBro.Update(0.5)
	
	if hammerBro.GetThrowTimer() <= initialTimer {
		t.Error("Expected throw timer to increase after update")
	}
	
	// Dead Hammer Bro should not update
	hammerBro.Die()
	deadTimer := hammerBro.GetThrowTimer()
	hammerBro.Update(1.0)
	
	if hammerBro.GetThrowTimer() != deadTimer {
		t.Error("Dead Hammer Bro should not update throw timer")
	}
}

func TestHammerBroInactiveState(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	// Deactivate Hammer Bro
	hammerBro.Deactivate()
	
	initialTimer := hammerBro.GetThrowTimer()
	
	// Update should not affect inactive Hammer Bro
	hammerBro.Update(1.0)
	
	if hammerBro.GetThrowTimer() != initialTimer {
		t.Error("Inactive Hammer Bro should not update")
	}
}

func TestHammerBroGetSprite(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	sprite := hammerBro.GetSprite()
	if sprite != 'H' {
		t.Errorf("Expected sprite 'H', got '%c'", sprite)
	}
}

func TestHammerBroGetSpriteStyle(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	style := hammerBro.GetSpriteStyle()
	if style == (tcell.Style{}) {
		t.Error("Expected non-empty sprite style")
	}
}

func TestHammerBroIsAlive(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	if !hammerBro.IsAlive() {
		t.Error("New Hammer Bro should be alive")
	}
	
	hammerBro.Die()
	
	if hammerBro.IsAlive() {
		t.Error("Dead Hammer Bro should not be alive")
	}
}

func TestHammerBroActivateDeactivate(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	if !hammerBro.IsActive {
		t.Fatal("Hammer Bro should start active")
	}
	
	hammerBro.Deactivate()
	
	if hammerBro.IsActive {
		t.Error("Hammer Bro should be inactive after deactivation")
	}
	
	hammerBro.Activate()
	
	if !hammerBro.IsActive {
		t.Error("Hammer Bro should be active after activation")
	}
}

func TestHammerBroGetBounds(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	x, y, width, height := hammerBro.GetBounds()
	
	if x != 100 || y != 200 {
		t.Errorf("Expected bounds position (100, 200), got (%f, %f)", x, y)
	}
	
	if width != 16 || height != 16 {
		t.Errorf("Expected bounds size (16, 16), got (%f, %f)", width, height)
	}
}

func TestHammerBroSetPosition(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	hammerBro.SetPosition(150, 250)
	
	if hammerBro.Position.X != 150 || hammerBro.Position.Y != 250 {
		t.Errorf("Expected position (150, 250), got (%f, %f)", hammerBro.Position.X, hammerBro.Position.Y)
	}
}

func TestHammerBroSetVelocity(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	hammerBro.SetVelocity(50, 100)
	
	if hammerBro.Velocity.X != 50 || hammerBro.Velocity.Y != 100 {
		t.Errorf("Expected velocity (50, 100), got (%f, %f)", hammerBro.Velocity.X, hammerBro.Velocity.Y)
	}
}

func TestHammerBroOnCollideWithPlayer(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	// This should not crash (Hammer Bro doesn't directly damage player)
	hammerBro.OnCollideWithPlayer(nil)
	
	// Dead Hammer Bro should not interact
	hammerBro.Die()
	hammerBro.OnCollideWithPlayer(nil)
}

func TestHammerBroSetTargetPlayer(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	// This should not crash
	hammerBro.SetTargetPlayer(nil)
	hammerBro.SetTargetPlayer("mock_player")
}

func TestHammerBroType(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	if hammerBro.Type != EnemyTypeHammerBro {
		t.Errorf("Expected type EnemyTypeHammerBro, got %v", hammerBro.Type)
	}
}

func TestHammerBroNoSpawnerSet(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	// Should not crash when spawner is not set
	hammerBro.ThrowHammer()
	hammerBro.Update(2.0)
}

func TestHammerBroMultipleThrows(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	throwCount := 0
	hammerBro.SetHammerSpawner(func(x, y, vx, vy float64) {
		throwCount++
	})
	
	// Simulate 6 seconds (should throw 3 times)
	for i := 0; i < 6; i++ {
		hammerBro.Update(1.0)
	}
	
	if throwCount != 3 {
		t.Errorf("Expected 3 hammer throws in 6 seconds, got %d", throwCount)
	}
}

func TestHammerBroDamageDoesNotAffectThrowTimer(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	hammerBro.Update(1.0)
	timerBeforeDamage := hammerBro.GetThrowTimer()
	
	// Take damage
	hammerBro.OnStomp(nil)
	
	// Timer should not be affected by damage
	if hammerBro.GetThrowTimer() != timerBeforeDamage {
		t.Error("Throw timer should not be affected by taking damage")
	}
}

func TestHammerBroThrowCooldown(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	if hammerBro.throwCooldown != 2.0 {
		t.Errorf("Expected throw cooldown 2.0 seconds, got %f", hammerBro.throwCooldown)
	}
}

func TestHammerBroMoveSpeed(t *testing.T) {
	hammerBro := NewHammerBro(100, 200)
	
	if hammerBro.MoveSpeed != 40.0 {
		t.Errorf("Expected move speed 40.0 pixels/second, got %f", hammerBro.MoveSpeed)
	}
}
