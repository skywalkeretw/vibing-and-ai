package entities

import (
	"testing"

	"github.com/lukeroy/go-terminal-platformer/internal/engine"
)

func TestNewPowerUp(t *testing.T) {
	powerUp := NewPowerUp()
	
	if powerUp == nil {
		t.Fatal("NewPowerUp() returned nil")
	}
	
	if !powerUp.IsActive() {
		t.Error("New power-up should be active")
	}
	
	if powerUp.IsCollected() {
		t.Error("New power-up should not be collected")
	}
}

func TestPowerUpInitialize(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	powerUp := NewPowerUp()
	pos := engine.Vector2{X: 100, Y: 100}
	powerUp.Initialize(PowerUpFire, pos, physics)
	
	if powerUp.GetType() != PowerUpFire {
		t.Errorf("Expected PowerUpFire, got %v", powerUp.GetType())
	}
	
	if powerUp.GetPosition() != pos {
		t.Errorf("Expected position %v, got %v", pos, powerUp.GetPosition())
	}
	
	if powerUp.physicsBody == nil {
		t.Error("Physics body should not be nil after initialization")
	}
	
	if powerUp.sprite == nil {
		t.Error("Sprite should not be nil after initialization")
	}
}

func TestPowerUpTypes(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	types := []PowerUpType{
		PowerUpFire,
		PowerUpStar,
		PowerUpMushroom,
		PowerUpSpeedBoots,
		PowerUpSuperJump,
		PowerUpShield,
	}
	
	for _, powerUpType := range types {
		powerUp := NewPowerUp()
		powerUp.Initialize(powerUpType, engine.Vector2{X: 0, Y: 0}, physics)
		
		if powerUp.GetType() != powerUpType {
			t.Errorf("Expected %v, got %v", powerUpType, powerUp.GetType())
		}
		
		if powerUp.sprite == nil {
			t.Errorf("Sprite should not be nil for %v", powerUpType)
		}
	}
}

func TestPowerUpCollection(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	powerUp := NewPowerUp()
	powerUp.Initialize(PowerUpFire, engine.Vector2{X: 100, Y: 100}, physics)
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	// Collect power-up
	powerUp.collect(player)
	
	if !powerUp.IsCollected() {
		t.Error("Power-up should be collected")
	}
	
	if powerUp.IsActive() {
		t.Error("Power-up should not be active after collection")
	}
	
	if player.GetPowerUp() != PowerUpFire {
		t.Errorf("Player should have PowerUpFire, got %v", player.GetPowerUp())
	}
}

func TestPowerUpUpdate(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	powerUp := NewPowerUp()
	powerUp.Initialize(PowerUpFire, engine.Vector2{X: 100, Y: 100}, physics)
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	players := []*PlayerEntity{player}
	
	// Update should trigger collection
	powerUp.Update(0.016, players)
	
	if !powerUp.IsCollected() {
		t.Error("Power-up should be collected after update with nearby player")
	}
}

func TestPowerUpFloating(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	powerUp := NewPowerUp()
	powerUp.Initialize(PowerUpFire, engine.Vector2{X: 100, Y: 100}, physics)
	
	initialOffset := powerUp.floatOffset
	
	// Update multiple times to see floating animation
	for i := 0; i < 10; i++ {
		powerUp.Update(0.016, []*PlayerEntity{})
	}
	
	if powerUp.floatOffset == initialOffset {
		t.Error("Float offset should change during updates")
	}
}

func TestPowerUpAnimation(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	powerUp := NewPowerUp()
	powerUp.Initialize(PowerUpFire, engine.Vector2{X: 100, Y: 100}, physics)
	
	initialFrame := powerUp.animFrame
	
	// Update animation multiple times
	for i := 0; i < 20; i++ {
		powerUp.updateAnimation(0.016)
	}
	
	if powerUp.animFrame == initialFrame {
		t.Error("Animation frame should change after multiple updates")
	}
}

func TestPowerUpDeactivate(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	powerUp := NewPowerUp()
	powerUp.Initialize(PowerUpFire, engine.Vector2{X: 100, Y: 100}, physics)
	
	powerUp.Deactivate()
	
	if powerUp.IsActive() {
		t.Error("Power-up should not be active after deactivation")
	}
	
	if powerUp.physicsBody.Enabled {
		t.Error("Physics body should be disabled after deactivation")
	}
}

func TestPowerUpSetPosition(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	powerUp := NewPowerUp()
	powerUp.Initialize(PowerUpFire, engine.Vector2{X: 100, Y: 100}, physics)
	
	newPos := engine.Vector2{X: 200, Y: 200}
	powerUp.SetPosition(newPos)
	
	if powerUp.GetPosition() != newPos {
		t.Errorf("Expected position %v, got %v", newPos, powerUp.GetPosition())
	}
	
	if powerUp.physicsBody.Position != newPos {
		t.Error("Physics body position should match power-up position")
	}
}

func TestPowerUpCollisionDetection(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	powerUp := NewPowerUp()
	powerUp.Initialize(PowerUpFire, engine.Vector2{X: 100, Y: 100}, physics)
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	// Should collide when at same position
	if !powerUp.checkCollision(player) {
		t.Error("Should detect collision when player is at same position")
	}
}

func TestPowerUpNoCollectionWhenInactive(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	powerUp := NewPowerUp()
	powerUp.Initialize(PowerUpFire, engine.Vector2{X: 100, Y: 100}, physics)
	powerUp.SetActive(false)
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	players := []*PlayerEntity{player}
	
	// Update should not trigger collection when inactive
	powerUp.Update(0.016, players)
	
	if powerUp.IsCollected() {
		t.Error("Inactive power-up should not be collected")
	}
}

func TestPowerUpMultiplePlayers(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	powerUp := NewPowerUp()
	powerUp.Initialize(PowerUpFire, engine.Vector2{X: 100, Y: 100}, physics)
	
	player1 := NewPlayer()
	player1.Initialize(1, engine.Vector2{X: 200, Y: 200}, physics)
	
	player2 := NewPlayer()
	player2.Initialize(2, engine.Vector2{X: 100, Y: 100}, physics)
	
	players := []*PlayerEntity{player1, player2}
	
	// Update should trigger collection by player2
	powerUp.Update(0.016, players)
	
	if !powerUp.IsCollected() {
		t.Error("Power-up should be collected by player2")
	}
	
	if player2.GetPowerUp() != PowerUpFire {
		t.Error("Player2 should have the power-up")
	}
	
	if player1.GetPowerUp() == PowerUpFire {
		t.Error("Player1 should not have the power-up")
	}
}
