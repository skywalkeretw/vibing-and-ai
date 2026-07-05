package entities

import (
	"testing"

	"github.com/lukeroy/go-terminal-platformer/internal/engine"
)

func TestNewPowerUpManager(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	if manager == nil {
		t.Fatal("NewPowerUpManager() returned nil")
	}
	
	if manager.GetPowerUpCount() != 0 {
		t.Error("New manager should have 0 power-ups")
	}
	
	if manager.GetCoinCount() != 0 {
		t.Error("New manager should have 0 coins")
	}
}

func TestPowerUpManagerSpawnPowerUp(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	pos := engine.Vector2{X: 100, Y: 100}
	
	powerUp := manager.SpawnPowerUp(PowerUpFire, pos)
	
	if powerUp == nil {
		t.Fatal("SpawnPowerUp() returned nil")
	}
	
	if manager.GetPowerUpCount() != 1 {
		t.Errorf("Expected 1 power-up, got %d", manager.GetPowerUpCount())
	}
	
	if powerUp.GetType() != PowerUpFire {
		t.Errorf("Expected PowerUpFire, got %v", powerUp.GetType())
	}
}

func TestPowerUpManagerSpawnCoin(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	pos := engine.Vector2{X: 100, Y: 100}
	
	coin := manager.SpawnCoin(pos)
	
	if coin == nil {
		t.Fatal("SpawnCoin() returned nil")
	}
	
	if manager.GetCoinCount() != 1 {
		t.Errorf("Expected 1 coin, got %d", manager.GetCoinCount())
	}
}

func TestPowerUpManagerSpawnCoinWithValue(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	pos := engine.Vector2{X: 100, Y: 100}
	
	coin := manager.SpawnCoinWithValue(pos, 5)
	
	if coin.GetValue() != 5 {
		t.Errorf("Expected value 5, got %d", coin.GetValue())
	}
}

func TestPowerUpManagerSpawnRandomPowerUp(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	pos := engine.Vector2{X: 100, Y: 100}
	
	// Spawn multiple random power-ups
	for i := 0; i < 10; i++ {
		powerUp := manager.SpawnRandomPowerUp(pos)
		if powerUp == nil {
			t.Fatal("SpawnRandomPowerUp() returned nil")
		}
		if powerUp.GetType() == PowerUpNone {
			t.Error("Random power-up should not be PowerUpNone")
		}
	}
	
	if manager.GetPowerUpCount() != 10 {
		t.Errorf("Expected 10 power-ups, got %d", manager.GetPowerUpCount())
	}
}

func TestPowerUpManagerUpdate(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	powerUp := manager.SpawnPowerUp(PowerUpFire, engine.Vector2{X: 100, Y: 100})
	coin := manager.SpawnCoin(engine.Vector2{X: 100, Y: 100})
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	players := []*PlayerEntity{player}
	
	// Update should trigger collection
	manager.Update(0.016, players)
	
	if !powerUp.IsCollected() {
		t.Error("Power-up should be collected")
	}
	
	if !coin.IsCollected() {
		t.Error("Coin should be collected")
	}
}

func TestPowerUpManagerRemoveCollected(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	_ = manager.SpawnPowerUp(PowerUpFire, engine.Vector2{X: 100, Y: 100})
	_ = manager.SpawnCoin(engine.Vector2{X: 100, Y: 100})
	
	player := NewPlayer()
	player.Initialize(1, engine.Vector2{X: 100, Y: 100}, physics)
	
	players := []*PlayerEntity{player}
	
	// Update to collect items
	manager.Update(0.016, players)
	
	// Items should be removed from manager
	if manager.GetPowerUpCount() != 0 {
		t.Errorf("Expected 0 power-ups after collection, got %d", manager.GetPowerUpCount())
	}
	
	if manager.GetCoinCount() != 0 {
		t.Errorf("Expected 0 coins after collection, got %d", manager.GetCoinCount())
	}
}

func TestPowerUpManagerClear(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	manager.SpawnPowerUp(PowerUpFire, engine.Vector2{X: 100, Y: 100})
	manager.SpawnCoin(engine.Vector2{X: 150, Y: 100})
	
	manager.Clear()
	
	if manager.GetPowerUpCount() != 0 {
		t.Error("Manager should have 0 power-ups after clear")
	}
	
	if manager.GetCoinCount() != 0 {
		t.Error("Manager should have 0 coins after clear")
	}
}

func TestPowerUpManagerGetPowerUpsInRadius(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	manager.SpawnPowerUp(PowerUpFire, engine.Vector2{X: 100, Y: 100})
	manager.SpawnPowerUp(PowerUpStar, engine.Vector2{X: 150, Y: 100})
	manager.SpawnPowerUp(PowerUpMushroom, engine.Vector2{X: 300, Y: 100})
	
	center := engine.Vector2{X: 100, Y: 100}
	powerUps := manager.GetPowerUpsInRadius(center, 100)
	
	if len(powerUps) != 2 {
		t.Errorf("Expected 2 power-ups in radius, got %d", len(powerUps))
	}
}

func TestPowerUpManagerGetCoinsInRadius(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	manager.SpawnCoin(engine.Vector2{X: 100, Y: 100})
	manager.SpawnCoin(engine.Vector2{X: 150, Y: 100})
	manager.SpawnCoin(engine.Vector2{X: 300, Y: 100})
	
	center := engine.Vector2{X: 100, Y: 100}
	coins := manager.GetCoinsInRadius(center, 100)
	
	if len(coins) != 2 {
		t.Errorf("Expected 2 coins in radius, got %d", len(coins))
	}
}

func TestPowerUpManagerGetPowerUpsByType(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	manager.SpawnPowerUp(PowerUpFire, engine.Vector2{X: 100, Y: 100})
	manager.SpawnPowerUp(PowerUpFire, engine.Vector2{X: 150, Y: 100})
	manager.SpawnPowerUp(PowerUpStar, engine.Vector2{X: 200, Y: 100})
	
	firePowerUps := manager.GetPowerUpsByType(PowerUpFire)
	
	if len(firePowerUps) != 2 {
		t.Errorf("Expected 2 fire power-ups, got %d", len(firePowerUps))
	}
}

func TestPowerUpManagerRemovePowerUp(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	powerUp := manager.SpawnPowerUp(PowerUpFire, engine.Vector2{X: 100, Y: 100})
	
	if manager.GetPowerUpCount() != 1 {
		t.Error("Should have 1 power-up before removal")
	}
	
	manager.RemovePowerUp(powerUp)
	
	if manager.GetPowerUpCount() != 0 {
		t.Error("Should have 0 power-ups after removal")
	}
}

func TestPowerUpManagerRemoveCoin(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	coin := manager.SpawnCoin(engine.Vector2{X: 100, Y: 100})
	
	if manager.GetCoinCount() != 1 {
		t.Error("Should have 1 coin before removal")
	}
	
	manager.RemoveCoin(coin)
	
	if manager.GetCoinCount() != 0 {
		t.Error("Should have 0 coins after removal")
	}
}

func TestPowerUpManagerDeactivateAll(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	manager.SpawnPowerUp(PowerUpFire, engine.Vector2{X: 100, Y: 100})
	manager.SpawnCoin(engine.Vector2{X: 150, Y: 100})
	
	manager.DeactivateAll()
	
	for _, powerUp := range manager.GetAllPowerUps() {
		if powerUp.IsActive() {
			t.Error("Power-up should be deactivated")
		}
	}
	
	for _, coin := range manager.GetAllCoins() {
		if coin.IsActive() {
			t.Error("Coin should be deactivated")
		}
	}
}

func TestPowerUpManagerActivateAll(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	manager.SpawnPowerUp(PowerUpFire, engine.Vector2{X: 100, Y: 100})
	manager.SpawnCoin(engine.Vector2{X: 150, Y: 100})
	
	manager.DeactivateAll()
	manager.ActivateAll()
	
	for _, powerUp := range manager.GetAllPowerUps() {
		if !powerUp.IsActive() {
			t.Error("Power-up should be activated")
		}
	}
	
	for _, coin := range manager.GetAllCoins() {
		if !coin.IsActive() {
			t.Error("Coin should be activated")
		}
	}
}

func TestPowerUpManagerSpawnWithVelocity(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	pos := engine.Vector2{X: 100, Y: 100}
	velocity := engine.Vector2{X: 50, Y: -100}
	
	powerUp := manager.SpawnPowerUpWithVelocity(PowerUpFire, pos, velocity)
	
	if powerUp.physicsBody.Velocity != velocity {
		t.Errorf("Expected velocity %v, got %v", velocity, powerUp.physicsBody.Velocity)
	}
	
	coin := manager.SpawnCoinWithVelocity(pos, velocity)
	
	if coin.physicsBody.Velocity != velocity {
		t.Errorf("Expected velocity %v, got %v", velocity, coin.physicsBody.Velocity)
	}
}

func TestPowerUpManagerGetActivePowerUps(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	powerUp1 := manager.SpawnPowerUp(PowerUpFire, engine.Vector2{X: 100, Y: 100})
	manager.SpawnPowerUp(PowerUpStar, engine.Vector2{X: 150, Y: 100})
	
	powerUp1.SetActive(false)
	
	activePowerUps := manager.GetActivePowerUps()
	
	if len(activePowerUps) != 1 {
		t.Errorf("Expected 1 active power-up, got %d", len(activePowerUps))
	}
}

func TestPowerUpManagerGetActiveCoins(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()
	
	manager := NewPowerUpManager(physics)
	
	coin1 := manager.SpawnCoin(engine.Vector2{X: 100, Y: 100})
	manager.SpawnCoin(engine.Vector2{X: 150, Y: 100})
	
	coin1.SetActive(false)
	
	activeCoins := manager.GetActiveCoins()
	
	if len(activeCoins) != 1 {
		t.Errorf("Expected 1 active coin, got %d", len(activeCoins))
	}
}
