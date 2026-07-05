package entities

import (
	"testing"

	"github.com/lukeroy/go-terminal-platformer/internal/engine"
	"github.com/lukeroy/go-terminal-platformer/internal/input"
)

func TestNewPlayerEntity(t *testing.T) {
	player := NewPlayer()
	
	if player == nil {
		t.Fatal("NewPlayer() returned nil")
	}
	
	if !player.IsActive() {
		t.Error("New player should be active")
	}
	
	if player.GetLives() != 5 {
		t.Errorf("Expected 5 lives, got %d", player.GetLives())
	}
	
	if player.GetState() != StateIdle {
		t.Errorf("Expected StateIdle, got %v", player.GetState())
	}
}

func TestPlayerInitialize(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	startPos := engine.Vector2{X: 100, Y: 100}
	player.Initialize(1, startPos, physics)
	
	if player.GetID() != 1 {
		t.Errorf("Expected ID 1, got %d", player.GetID())
	}
	
	pos := player.GetPosition()
	if pos.X != startPos.X || pos.Y != startPos.Y {
		t.Errorf("Expected position %v, got %v", startPos, pos)
	}
	
	if player.GetFacing() != DirectionRight {
		t.Errorf("Expected DirectionRight, got %v", player.GetFacing())
	}
	
	if player.sprite == nil {
		t.Error("Player sprite should not be nil after initialization")
	}
}

func TestPlayerMovement(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	// Test move left
	player.moveLeft()
	if player.GetFacing() != DirectionLeft {
		t.Error("Player should face left after moveLeft()")
	}
	if player.physicsBody.Velocity.X >= 0 {
		t.Error("Player velocity should be negative when moving left")
	}
	
	// Test move right
	player.moveRight()
	if player.GetFacing() != DirectionRight {
		t.Error("Player should face right after moveRight()")
	}
	if player.physicsBody.Velocity.X <= 0 {
		t.Error("Player velocity should be positive when moving right")
	}
}

func TestPlayerJump(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	// Set player as grounded
	player.grounded = true
	player.physicsBody.Grounded = true
	
	// Test jump
	player.jump()
	
	if player.GetState() != StateJumping {
		t.Errorf("Expected StateJumping, got %v", player.GetState())
	}
	
	if player.physicsBody.Velocity.Y >= 0 {
		t.Error("Player velocity should be negative (upward) when jumping")
	}
	
	if player.jumpTimeLeft <= 0 {
		t.Error("Jump time should be set when jumping")
	}
}

func TestPlayerCrouch(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	// Set player as grounded
	player.grounded = true
	player.physicsBody.Grounded = true
	
	// Test crouch
	player.crouch()
	
	if player.GetState() != StateCrouching {
		t.Errorf("Expected StateCrouching, got %v", player.GetState())
	}
	
	if !player.crouching {
		t.Error("Player crouching flag should be true")
	}
	
	if player.physicsBody.Velocity.X != 0 {
		t.Error("Player horizontal velocity should be 0 when crouching")
	}
}

func TestPlayerTakeDamage(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	initialLives := player.GetLives()
	
	// Test taking damage
	player.TakeDamage()
	
	if player.GetLives() != initialLives-1 {
		t.Errorf("Expected %d lives, got %d", initialLives-1, player.GetLives())
	}
	
	if !player.IsInvulnerable() {
		t.Error("Player should be invulnerable after taking damage")
	}
	
	if player.GetState() != StateHurt {
		t.Errorf("Expected StateHurt, got %v", player.GetState())
	}
	
	// Test invulnerability prevents damage
	player.TakeDamage()
	if player.GetLives() != initialLives-1 {
		t.Error("Player should not take damage while invulnerable")
	}
}

func TestPlayerDeath(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	// Reduce lives to 1
	for player.GetLives() > 1 {
		player.lives--
	}
	
	// Take final damage
	player.invulnerable = false
	player.TakeDamage()
	
	if !player.IsDead() {
		t.Error("Player should be dead after losing all lives")
	}
	
	if player.GetState() != StateDead {
		t.Errorf("Expected StateDead, got %v", player.GetState())
	}
}

func TestPlayerPowerUps(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	// Test fire power-up
	player.ApplyPowerUp(PowerUpFire)
	if player.GetPowerUp() != PowerUpFire {
		t.Errorf("Expected PowerUpFire, got %v", player.GetPowerUp())
	}
	if player.GetAmmo() != 10 {
		t.Errorf("Expected 10 ammo, got %d", player.GetAmmo())
	}
	
	// Test speed boots
	player.ApplyPowerUp(PowerUpSpeedBoots)
	if player.moveSpeed != 360.0 {
		t.Errorf("Expected move speed 360.0, got %f", player.moveSpeed)
	}
	
	// Test super jump
	player.ApplyPowerUp(PowerUpSuperJump)
	if player.jumpForce != -600.0 {
		t.Errorf("Expected jump force -600.0, got %f", player.jumpForce)
	}
	
	// Test remove power-up
	player.RemovePowerUp()
	if player.GetPowerUp() != PowerUpNone {
		t.Errorf("Expected PowerUpNone, got %v", player.GetPowerUp())
	}
	if player.moveSpeed != 240.0 {
		t.Errorf("Expected move speed 240.0, got %f", player.moveSpeed)
	}
}

func TestPlayerShield(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	initialLives := player.GetLives()
	
	// Apply shield
	player.ApplyPowerUp(PowerUpShield)
	
	// Take damage with shield
	player.TakeDamage()
	
	if player.GetLives() != initialLives {
		t.Error("Player should not lose life when shield absorbs damage")
	}
	
	if player.GetPowerUp() != PowerUpNone {
		t.Error("Shield should be removed after absorbing damage")
	}
}

func TestPlayerUpdate(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	// Test update with no input
	player.Update(0.016) // ~60 FPS
	
	if player.GetState() != StateIdle {
		t.Errorf("Expected StateIdle with no input, got %v", player.GetState())
	}
	
	// Test update with movement input
	player.inputActions = []input.InputAction{input.ActionMoveRight}
	player.grounded = true
	player.physicsBody.Grounded = true
	player.Update(0.016)
	
	if player.GetState() != StateRunning {
		t.Errorf("Expected StateRunning with movement input, got %v", player.GetState())
	}
}

func TestPlayerStateTransitions(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	// Test idle to running
	player.grounded = true
	player.physicsBody.Grounded = true
	player.physicsBody.Velocity.X = 100
	player.updateState()
	if player.GetState() != StateRunning {
		t.Errorf("Expected StateRunning, got %v", player.GetState())
	}
	
	// Test running to jumping
	player.grounded = false
	player.physicsBody.Grounded = false
	player.physicsBody.Velocity.Y = -100
	player.updateState()
	if player.GetState() != StateJumping {
		t.Errorf("Expected StateJumping, got %v", player.GetState())
	}
	
	// Test jumping to falling
	player.physicsBody.Velocity.Y = 100
	player.updateState()
	if player.GetState() != StateFalling {
		t.Errorf("Expected StateFalling, got %v", player.GetState())
	}
	
	// Test falling to idle
	player.grounded = true
	player.physicsBody.Grounded = true
	player.physicsBody.Velocity.X = 0
	player.updateState()
	if player.GetState() != StateIdle {
		t.Errorf("Expected StateIdle, got %v", player.GetState())
	}
}

func TestPlayerAnimation(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	initialFrame := player.animFrame
	
	// Update animation multiple times
	for i := 0; i < 20; i++ {
		player.updateAnimation(0.016)
	}
	
	if player.animFrame == initialFrame {
		t.Error("Animation frame should have changed after multiple updates")
	}
}

func TestPlayerTwoPlayers(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player1 := NewPlayer()
	player1.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	player2 := NewPlayer()
	player2.Initialize(2, engine.Vector2{X: 150, Y: 100}, physics)
	
	if player1.GetID() == player2.GetID() {
		t.Error("Players should have different IDs")
	}
	
	// Test independent movement
	player1.moveLeft()
	player2.moveRight()
	
	if player1.GetFacing() == player2.GetFacing() {
		t.Error("Players should be able to face different directions")
	}
	
	// Test independent state
	player1.grounded = true
	player1.physicsBody.Grounded = true
	player1.jump()
	
	if player1.GetState() == player2.GetState() {
		t.Error("Players should be able to have different states")
	}
}

func TestPlayerAddLife(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	initialLives := player.GetLives()
	player.AddLife()
	
	if player.GetLives() != initialLives+1 {
		t.Errorf("Expected %d lives, got %d", initialLives+1, player.GetLives())
	}
}

func TestPlayerSpriteFlip(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	// Get sprite facing right
	spriteRight := player.getSpriteForState()
	
	// Change facing to left
	player.facing = DirectionLeft
	spriteLeft := spriteRight.FlipHorizontal()
	
	if spriteLeft == spriteRight {
		t.Error("Flipped sprite should be different from original")
	}
}
