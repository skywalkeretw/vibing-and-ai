package entities

import (
	"math/rand"

	"github.com/lukeroy/go-terminal-platformer/internal/engine"
)

// PowerUpManager manages all power-ups and coins in the game
type PowerUpManager struct {
	powerUps []*PowerUp
	coins    []*Coin
	physics  *engine.PhysicsEngine
}

// NewPowerUpManager creates a new power-up manager
func NewPowerUpManager(physics *engine.PhysicsEngine) *PowerUpManager {
	return &PowerUpManager{
		powerUps: make([]*PowerUp, 0),
		coins:    make([]*Coin, 0),
		physics:  physics,
	}
}

// Update updates all power-ups and coins
func (pm *PowerUpManager) Update(deltaTime float64, players []*PlayerEntity) {
	// Update all power-ups
	for _, powerUp := range pm.powerUps {
		if powerUp.IsActive() {
			powerUp.Update(deltaTime, players)
		}
	}

	// Update all coins
	for _, coin := range pm.coins {
		if coin.IsActive() {
			coin.Update(deltaTime, players)
		}
	}

	// Remove collected items
	pm.removeCollected()
}

// removeCollected removes collected power-ups and coins
func (pm *PowerUpManager) removeCollected() {
	// Remove collected power-ups
	activePowerUps := make([]*PowerUp, 0)
	for _, powerUp := range pm.powerUps {
		if !powerUp.IsCollected() {
			activePowerUps = append(activePowerUps, powerUp)
		}
	}
	pm.powerUps = activePowerUps

	// Remove collected coins
	activeCoins := make([]*Coin, 0)
	for _, coin := range pm.coins {
		if !coin.IsCollected() {
			activeCoins = append(activeCoins, coin)
		}
	}
	pm.coins = activeCoins
}

// SpawnPowerUp spawns a specific power-up at a position
func (pm *PowerUpManager) SpawnPowerUp(powerUpType PowerUpType, pos engine.Vector2) *PowerUp {
	powerUp := NewPowerUp()
	powerUp.Initialize(powerUpType, pos, pm.physics)
	pm.powerUps = append(pm.powerUps, powerUp)
	return powerUp
}

// SpawnCoin spawns a coin at a position
func (pm *PowerUpManager) SpawnCoin(pos engine.Vector2) *Coin {
	coin := NewCoin()
	coin.Initialize(pos, pm.physics)
	pm.coins = append(pm.coins, coin)
	return coin
}

// SpawnCoinWithValue spawns a coin with a specific value
func (pm *PowerUpManager) SpawnCoinWithValue(pos engine.Vector2, value int) *Coin {
	coin := NewCoin()
	coin.Initialize(pos, pm.physics)
	coin.SetValue(value)
	pm.coins = append(pm.coins, coin)
	return coin
}

// SpawnRandomPowerUp spawns a random power-up at a position
func (pm *PowerUpManager) SpawnRandomPowerUp(pos engine.Vector2) *PowerUp {
	// Weighted random selection
	roll := rand.Float64()

	var powerUpType PowerUpType
	if roll < 0.3 {
		powerUpType = PowerUpFire // 30% chance
	} else if roll < 0.5 {
		powerUpType = PowerUpStar // 20% chance
	} else if roll < 0.6 {
		powerUpType = PowerUpMushroom // 10% chance
	} else if roll < 0.75 {
		powerUpType = PowerUpSpeedBoots // 15% chance
	} else if roll < 0.9 {
		powerUpType = PowerUpSuperJump // 15% chance
	} else {
		powerUpType = PowerUpShield // 10% chance
	}

	return pm.SpawnPowerUp(powerUpType, pos)
}

// SpawnPowerUpWithVelocity spawns a power-up with initial velocity
func (pm *PowerUpManager) SpawnPowerUpWithVelocity(powerUpType PowerUpType, pos engine.Vector2, velocity engine.Vector2) *PowerUp {
	powerUp := pm.SpawnPowerUp(powerUpType, pos)
	if powerUp.physicsBody != nil {
		powerUp.physicsBody.Velocity = velocity
	}
	return powerUp
}

// SpawnCoinWithVelocity spawns a coin with initial velocity
func (pm *PowerUpManager) SpawnCoinWithVelocity(pos engine.Vector2, velocity engine.Vector2) *Coin {
	coin := pm.SpawnCoin(pos)
	if coin.physicsBody != nil {
		coin.physicsBody.Velocity = velocity
	}
	return coin
}

// Render renders all power-ups and coins
func (pm *PowerUpManager) Render(renderer Renderer) {
	// Render all power-ups
	for _, powerUp := range pm.powerUps {
		if powerUp.IsActive() {
			powerUp.Render(renderer)
		}
	}

	// Render all coins
	for _, coin := range pm.coins {
		if coin.IsActive() {
			coin.Render(renderer)
		}
	}
}

// Clear removes all power-ups and coins
func (pm *PowerUpManager) Clear() {
	// Remove all power-ups from physics
	for _, powerUp := range pm.powerUps {
		if powerUp.physicsBody != nil {
			pm.physics.RemoveBody(powerUp.physicsBody)
		}
	}

	// Remove all coins from physics
	for _, coin := range pm.coins {
		if coin.physicsBody != nil {
			pm.physics.RemoveBody(coin.physicsBody)
		}
	}

	pm.powerUps = make([]*PowerUp, 0)
	pm.coins = make([]*Coin, 0)
}

// GetPowerUpCount returns the number of active power-ups
func (pm *PowerUpManager) GetPowerUpCount() int {
	count := 0
	for _, powerUp := range pm.powerUps {
		if powerUp.IsActive() {
			count++
		}
	}
	return count
}

// GetCoinCount returns the number of active coins
func (pm *PowerUpManager) GetCoinCount() int {
	count := 0
	for _, coin := range pm.coins {
		if coin.IsActive() {
			count++
		}
	}
	return count
}

// GetPowerUpsInRadius returns all power-ups within a radius of a point
func (pm *PowerUpManager) GetPowerUpsInRadius(center engine.Vector2, radius float64) []*PowerUp {
	result := make([]*PowerUp, 0)
	radiusSquared := radius * radius

	for _, powerUp := range pm.powerUps {
		if powerUp.IsActive() {
			distSquared := center.DistanceSquared(powerUp.GetPosition())
			if distSquared <= radiusSquared {
				result = append(result, powerUp)
			}
		}
	}

	return result
}

// GetCoinsInRadius returns all coins within a radius of a point
func (pm *PowerUpManager) GetCoinsInRadius(center engine.Vector2, radius float64) []*Coin {
	result := make([]*Coin, 0)
	radiusSquared := radius * radius

	for _, coin := range pm.coins {
		if coin.IsActive() {
			distSquared := center.DistanceSquared(coin.GetPosition())
			if distSquared <= radiusSquared {
				result = append(result, coin)
			}
		}
	}

	return result
}

// GetPowerUpsByType returns all power-ups of a specific type
func (pm *PowerUpManager) GetPowerUpsByType(powerUpType PowerUpType) []*PowerUp {
	result := make([]*PowerUp, 0)

	for _, powerUp := range pm.powerUps {
		if powerUp.IsActive() && powerUp.GetType() == powerUpType {
			result = append(result, powerUp)
		}
	}

	return result
}

// RemovePowerUp removes a specific power-up
func (pm *PowerUpManager) RemovePowerUp(powerUp *PowerUp) {
	for i, p := range pm.powerUps {
		if p == powerUp {
			// Remove from physics
			if p.physicsBody != nil {
				pm.physics.RemoveBody(p.physicsBody)
			}
			// Remove from list
			pm.powerUps = append(pm.powerUps[:i], pm.powerUps[i+1:]...)
			break
		}
	}
}

// RemoveCoin removes a specific coin
func (pm *PowerUpManager) RemoveCoin(coin *Coin) {
	for i, c := range pm.coins {
		if c == coin {
			// Remove from physics
			if c.physicsBody != nil {
				pm.physics.RemoveBody(c.physicsBody)
			}
			// Remove from list
			pm.coins = append(pm.coins[:i], pm.coins[i+1:]...)
			break
		}
	}
}

// DeactivateAll deactivates all power-ups and coins
func (pm *PowerUpManager) DeactivateAll() {
	for _, powerUp := range pm.powerUps {
		powerUp.Deactivate()
	}
	for _, coin := range pm.coins {
		coin.Deactivate()
	}
}

// ActivateAll activates all power-ups and coins
func (pm *PowerUpManager) ActivateAll() {
	for _, powerUp := range pm.powerUps {
		powerUp.SetActive(true)
		if powerUp.physicsBody != nil {
			powerUp.physicsBody.Enabled = true
		}
	}
	for _, coin := range pm.coins {
		coin.SetActive(true)
		if coin.physicsBody != nil {
			coin.physicsBody.Enabled = true
		}
	}
}

// GetAllPowerUps returns all power-ups (active and inactive)
func (pm *PowerUpManager) GetAllPowerUps() []*PowerUp {
	return pm.powerUps
}

// GetAllCoins returns all coins (active and inactive)
func (pm *PowerUpManager) GetAllCoins() []*Coin {
	return pm.coins
}

// GetActivePowerUps returns only active power-ups
func (pm *PowerUpManager) GetActivePowerUps() []*PowerUp {
	result := make([]*PowerUp, 0)
	for _, powerUp := range pm.powerUps {
		if powerUp.IsActive() {
			result = append(result, powerUp)
		}
	}
	return result
}

// GetActiveCoins returns only active coins
func (pm *PowerUpManager) GetActiveCoins() []*Coin {
	result := make([]*Coin, 0)
	for _, coin := range pm.coins {
		if coin.IsActive() {
			result = append(result, coin)
		}
	}
	return result
}
