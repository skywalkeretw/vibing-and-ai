package entities

import (
	"math/rand"
	"time"

	"github.com/lukeroy/go-terminal-platformer/internal/engine"
)

// SpawnPoint represents a location where enemies can spawn
type SpawnPoint struct {
	Position   engine.Vector2
	EnemyType  EnemyType
	SpawnRate  float64 // Seconds between spawns
	MaxSpawns  int     // Maximum number of enemies this point can spawn (-1 for unlimited)
	Spawned    int     // Number of enemies spawned so far
	SpawnTimer float64 // Internal timer for spawning
	Active     bool    // Whether this spawn point is active
}

// NewSpawnPoint creates a new spawn point
func NewSpawnPoint(pos engine.Vector2, enemyType EnemyType, spawnRate float64, maxSpawns int) *SpawnPoint {
	return &SpawnPoint{
		Position:   pos,
		EnemyType:  enemyType,
		SpawnRate:  spawnRate,
		MaxSpawns:  maxSpawns,
		Spawned:    0,
		SpawnTimer: 0,
		Active:     true,
	}
}

// CanSpawn checks if this spawn point can spawn an enemy
func (sp *SpawnPoint) CanSpawn() bool {
	if !sp.Active {
		return false
	}
	if sp.MaxSpawns >= 0 && sp.Spawned >= sp.MaxSpawns {
		return false
	}
	return true
}

// EnemySpawner manages enemy spawning throughout the level
type EnemySpawner struct {
	enemies       []Enemy
	spawnPoints   []*SpawnPoint
	maxEnemies    int
	spawnInterval float64
	spawnTimer    float64
	physics       *engine.PhysicsEngine
	rng           *rand.Rand
	nextID        int
}

// NewEnemySpawner creates a new enemy spawner
func NewEnemySpawner(physics *engine.PhysicsEngine, maxEnemies int) *EnemySpawner {
	return &EnemySpawner{
		enemies:       make([]Enemy, 0),
		spawnPoints:   make([]*SpawnPoint, 0),
		maxEnemies:    maxEnemies,
		spawnInterval: 2.0, // Default 2 seconds between spawn attempts
		spawnTimer:    0,
		physics:       physics,
		rng:           rand.New(rand.NewSource(time.Now().UnixNano())),
		nextID:        1,
	}
}

// AddSpawnPoint adds a spawn point to the spawner
func (es *EnemySpawner) AddSpawnPoint(spawnPoint *SpawnPoint) {
	es.spawnPoints = append(es.spawnPoints, spawnPoint)
}

// RemoveSpawnPoint removes a spawn point from the spawner
func (es *EnemySpawner) RemoveSpawnPoint(spawnPoint *SpawnPoint) {
	for i, sp := range es.spawnPoints {
		if sp == spawnPoint {
			es.spawnPoints = append(es.spawnPoints[:i], es.spawnPoints[i+1:]...)
			break
		}
	}
}

// Update updates the spawner and all enemies
func (es *EnemySpawner) Update(deltaTime float64, players []Player) {
	// Update spawn timer
	es.spawnTimer += deltaTime

	// Try to spawn enemies at intervals
	if es.spawnTimer >= es.spawnInterval {
		es.trySpawn()
		es.spawnTimer = 0
	}

	// Update all spawn point timers
	for _, sp := range es.spawnPoints {
		if sp.Active && sp.CanSpawn() {
			sp.SpawnTimer += deltaTime
			if sp.SpawnTimer >= sp.SpawnRate {
				es.spawnAtPoint(sp)
				sp.SpawnTimer = 0
			}
		}
	}

	// Update all enemies
	for _, enemy := range es.enemies {
		if enemy.IsActive() {
			enemy.Update(deltaTime, players)
		}
	}

	// Remove dead enemies
	es.removeDeadEnemies()
}

// trySpawn attempts to spawn an enemy at a random spawn point
func (es *EnemySpawner) trySpawn() {
	// Check if we've reached max enemies
	if len(es.enemies) >= es.maxEnemies {
		return
	}

	// Find available spawn points
	availablePoints := make([]*SpawnPoint, 0)
	for _, sp := range es.spawnPoints {
		if sp.CanSpawn() {
			availablePoints = append(availablePoints, sp)
		}
	}

	if len(availablePoints) == 0 {
		return
	}

	// Select random spawn point
	spawnPoint := availablePoints[es.rng.Intn(len(availablePoints))]
	es.spawnAtPoint(spawnPoint)
}

// spawnAtPoint spawns an enemy at a specific spawn point
func (es *EnemySpawner) spawnAtPoint(spawnPoint *SpawnPoint) {
	if !spawnPoint.CanSpawn() {
		return
	}

	// Check if we've reached max enemies
	if len(es.enemies) >= es.maxEnemies {
		return
	}

	// Create enemy based on type
	enemy := es.createEnemy(spawnPoint.EnemyType, spawnPoint.Position)
	if enemy == nil {
		return
	}

	// Initialize enemy
	enemy.Initialize(spawnPoint.Position, es.physics)

	// Add to enemies list
	es.enemies = append(es.enemies, enemy)

	// Update spawn point
	spawnPoint.Spawned++
}

// createEnemy creates an enemy of the specified type
func (es *EnemySpawner) createEnemy(enemyType EnemyType, pos engine.Vector2) Enemy {
	enemy := NewEnemyBase(enemyType, pos.X, pos.Y)
	enemy.ID = es.nextID
	es.nextID++

	// Set default patrol range based on spawn position
	enemy.SetPatrolRange(pos.X-100, pos.X+100)

	return enemy
}

// removeDeadEnemies removes dead enemies from the list
func (es *EnemySpawner) removeDeadEnemies() {
	activeEnemies := make([]Enemy, 0)
	for _, enemy := range es.enemies {
		if enemy.IsAlive() {
			activeEnemies = append(activeEnemies, enemy)
		} else {
			// Remove from physics engine
			if enemyBase, ok := enemy.(*EnemyBase); ok {
				if enemyBase.PhysicsBody != nil {
					es.physics.RemoveBody(enemyBase.PhysicsBody)
				}
			}
		}
	}
	es.enemies = activeEnemies
}

// GetEnemies returns all active enemies
func (es *EnemySpawner) GetEnemies() []Enemy {
	return es.enemies
}

// GetEnemyCount returns the number of active enemies
func (es *EnemySpawner) GetEnemyCount() int {
	return len(es.enemies)
}

// GetAliveEnemyCount returns the number of alive enemies
func (es *EnemySpawner) GetAliveEnemyCount() int {
	count := 0
	for _, enemy := range es.enemies {
		if enemy.IsAlive() {
			count++
		}
	}
	return count
}

// Clear removes all enemies and spawn points
func (es *EnemySpawner) Clear() {
	// Remove all enemies from physics
	for _, enemy := range es.enemies {
		if enemyBase, ok := enemy.(*EnemyBase); ok {
			if enemyBase.PhysicsBody != nil {
				es.physics.RemoveBody(enemyBase.PhysicsBody)
			}
		}
	}

	es.enemies = make([]Enemy, 0)
	es.spawnPoints = make([]*SpawnPoint, 0)
	es.spawnTimer = 0
	es.nextID = 1
}

// SetMaxEnemies sets the maximum number of enemies
func (es *EnemySpawner) SetMaxEnemies(max int) {
	es.maxEnemies = max
}

// GetMaxEnemies returns the maximum number of enemies
func (es *EnemySpawner) GetMaxEnemies() int {
	return es.maxEnemies
}

// SetSpawnInterval sets the spawn interval
func (es *EnemySpawner) SetSpawnInterval(interval float64) {
	es.spawnInterval = interval
}

// GetSpawnInterval returns the spawn interval
func (es *EnemySpawner) GetSpawnInterval() float64 {
	return es.spawnInterval
}

// ActivateSpawnPoint activates a spawn point
func (es *EnemySpawner) ActivateSpawnPoint(index int) {
	if index >= 0 && index < len(es.spawnPoints) {
		es.spawnPoints[index].Active = true
	}
}

// DeactivateSpawnPoint deactivates a spawn point
func (es *EnemySpawner) DeactivateSpawnPoint(index int) {
	if index >= 0 && index < len(es.spawnPoints) {
		es.spawnPoints[index].Active = false
	}
}

// GetSpawnPointCount returns the number of spawn points
func (es *EnemySpawner) GetSpawnPointCount() int {
	return len(es.spawnPoints)
}

// GetSpawnPoint returns a spawn point by index
func (es *EnemySpawner) GetSpawnPoint(index int) *SpawnPoint {
	if index >= 0 && index < len(es.spawnPoints) {
		return es.spawnPoints[index]
	}
	return nil
}

// SpawnImmediate immediately spawns an enemy at a position
func (es *EnemySpawner) SpawnImmediate(enemyType EnemyType, pos engine.Vector2) Enemy {
	if len(es.enemies) >= es.maxEnemies {
		return nil
	}

	enemy := es.createEnemy(enemyType, pos)
	if enemy == nil {
		return nil
	}

	enemy.Initialize(pos, es.physics)
	es.enemies = append(es.enemies, enemy)

	return enemy
}

// KillAllEnemies kills all active enemies
func (es *EnemySpawner) KillAllEnemies() {
	for _, enemy := range es.enemies {
		if enemy.IsAlive() {
			enemy.Die()
		}
	}
}

// GetEnemiesInRadius returns all enemies within a radius of a point
func (es *EnemySpawner) GetEnemiesInRadius(center engine.Vector2, radius float64) []Enemy {
	enemiesInRadius := make([]Enemy, 0)
	radiusSquared := radius * radius

	for _, enemy := range es.enemies {
		if !enemy.IsAlive() {
			continue
		}

		distSquared := center.DistanceSquared(enemy.GetPosition())
		if distSquared <= radiusSquared {
			enemiesInRadius = append(enemiesInRadius, enemy)
		}
	}

	return enemiesInRadius
}

// GetEnemiesByType returns all enemies of a specific type
func (es *EnemySpawner) GetEnemiesByType(enemyType EnemyType) []Enemy {
	enemiesByType := make([]Enemy, 0)

	for _, enemy := range es.enemies {
		if enemy.GetType() == enemyType && enemy.IsAlive() {
			enemiesByType = append(enemiesByType, enemy)
		}
	}

	return enemiesByType
}
