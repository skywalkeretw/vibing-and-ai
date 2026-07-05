package entities

import (
	"testing"

	"github.com/lukeroy/go-terminal-platformer/internal/engine"
)

func TestNewCoin(t *testing.T) {
	coin := NewCoin()
	
	if coin == nil {
		t.Fatal("NewCoin() returned nil")
	}
	
	if !coin.IsActive() {
		t.Error("New coin should be active")
	}
	
	if coin.IsCollected() {
		t.Error("New coin should not be collected")
	}
	
	if coin.GetValue() != 1 {
		t.Errorf("Expected value 1, got %d", coin.GetValue())
	}
}

func TestCoinInitialize(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	coin := NewCoin()
	pos := engine.Vector2{X: 100, Y: 100}
	coin.Initialize(pos, physics)
	
	if coin.GetPosition() != pos {
		t.Errorf("Expected position %v, got %v", pos, coin.GetPosition())
	}
	
	if coin.physicsBody == nil {
		t.Error("Physics body should not be nil after initialization")
	}
}

func TestCoinCollection(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	coin := NewCoin()
	coin.Initialize(engine.Vector2{X: 100, Y: 100}, physics)
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	initialCoins := player.GetCoins()
	
	// Collect coin
	coin.collect(player)
	
	if !coin.IsCollected() {
		t.Error("Coin should be collected")
	}
	
	if coin.IsActive() {
		t.Error("Coin should not be active after collection")
	}
	
	if player.GetCoins() != initialCoins+1 {
		t.Errorf("Player should have %d coins, got %d", initialCoins+1, player.GetCoins())
	}
}

func TestCoinUpdate(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	coin := NewCoin()
	coin.Initialize(engine.Vector2{X: 100, Y: 100}, physics)
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	players := []*PlayerEntity{player}
	
	// Update should trigger collection
	coin.Update(0.016, players)
	
	if !coin.IsCollected() {
		t.Error("Coin should be collected after update with nearby player")
	}
}

func TestCoinAnimation(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	coin := NewCoin()
	coin.Initialize(engine.Vector2{X: 100, Y: 100}, physics)
	
	initialFrame := coin.animFrame
	
	// Update animation multiple times
	for i := 0; i < 15; i++ {
		coin.Update(0.016, []*PlayerEntity{})
	}
	
	if coin.animFrame == initialFrame {
		t.Error("Animation frame should change after multiple updates")
	}
}

func TestCoinSetValue(t *testing.T) {
	coin := NewCoin()
	
	coin.SetValue(5)
	
	if coin.GetValue() != 5 {
		t.Errorf("Expected value 5, got %d", coin.GetValue())
	}
}

func TestCoinDeactivate(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	coin := NewCoin()
	coin.Initialize(engine.Vector2{X: 100, Y: 100}, physics)
	
	coin.Deactivate()
	
	if coin.IsActive() {
		t.Error("Coin should not be active after deactivation")
	}
	
	if coin.physicsBody.Enabled {
		t.Error("Physics body should be disabled after deactivation")
	}
}

func TestCoinSetPosition(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	coin := NewCoin()
	coin.Initialize(engine.Vector2{X: 100, Y: 100}, physics)
	
	newPos := engine.Vector2{X: 200, Y: 200}
	coin.SetPosition(newPos)
	
	if coin.GetPosition() != newPos {
		t.Errorf("Expected position %v, got %v", newPos, coin.GetPosition())
	}
	
	if coin.physicsBody.Position != newPos {
		t.Error("Physics body position should match coin position")
	}
}

func TestCoinCollisionDetection(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	coin := NewCoin()
	coin.Initialize(engine.Vector2{X: 100, Y: 100}, physics)
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	// Should collide when at same position
	if !coin.checkCollision(player) {
		t.Error("Should detect collision when player is at same position")
	}
}

func TestCoinNoCollectionWhenInactive(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	coin := NewCoin()
	coin.Initialize(engine.Vector2{X: 100, Y: 100}, physics)
	coin.SetActive(false)
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	players := []*PlayerEntity{player}
	
	// Update should not trigger collection when inactive
	coin.Update(0.016, players)
	
	if coin.IsCollected() {
		t.Error("Inactive coin should not be collected")
	}
}

func TestCoinExtraLife(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	initialLives := player.GetLives()
	
	// Collect 100 coins
	for i := 0; i < 100; i++ {
		coin := NewCoin()
		coin.Initialize(engine.Vector2{X: 100, Y: 100}, physics)
		coin.collect(player)
	}
	
	if player.GetLives() != initialLives+1 {
		t.Errorf("Player should have %d lives after 100 coins, got %d", initialLives+1, player.GetLives())
	}
	
	if player.GetCoins() != 0 {
		t.Errorf("Player should have 0 coins after getting extra life, got %d", player.GetCoins())
	}
}

func TestCoinMultipleExtraLives(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	initialLives := player.GetLives()
	
	// Collect 250 coins (should give 2 extra lives)
	for i := 0; i < 250; i++ {
		coin := NewCoin()
		coin.Initialize(engine.Vector2{X: 100, Y: 100}, physics)
		coin.collect(player)
	}
	
	if player.GetLives() != initialLives+2 {
		t.Errorf("Player should have %d lives after 250 coins, got %d", initialLives+2, player.GetLives())
	}
	
	if player.GetCoins() != 50 {
		t.Errorf("Player should have 50 coins remaining, got %d", player.GetCoins())
	}
}

func TestCoinMultiplePlayers(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	coin := NewCoin()
	coin.Initialize(engine.Vector2{X: 100, Y: 100}, physics)
	
	player1 := NewPlayer()
	player1.Initialize(1, engine.Vector2{X: 200, Y: 200}, physics)
	
	player2 := NewPlayer()
	player2.Initialize(2, engine.Vector2{X: 100, Y: 100}, physics)
	
	players := []*PlayerEntity{player1, player2}
	
	// Update should trigger collection by player2
	coin.Update(0.016, players)
	
	if !coin.IsCollected() {
		t.Error("Coin should be collected by player2")
	}
	
	if player2.GetCoins() != 1 {
		t.Error("Player2 should have 1 coin")
	}
	
	if player1.GetCoins() != 0 {
		t.Error("Player1 should have 0 coins")
	}
}
