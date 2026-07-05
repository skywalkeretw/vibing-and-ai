package entities

import (
	"testing"

	"github.com/lukeroy/go-terminal-platformer/internal/engine"
)

func TestNewEnemySpawner(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)

	if spawner == nil {
		t.Fatal("NewEnemySpawner returned nil")
	}

	if spawner.maxEnemies != 10 {
		t.Errorf("Expected max enemies 10, got %d", spawner.maxEnemies)
	}

	if spawner.physics != physics {
		t.Error("Physics engine not set correctly")
	}

	if len(spawner.enemies) != 0 {
		t.Error("Should start with no enemies")
	}

	if len(spawner.spawnPoints) != 0 {
		t.Error("Should start with no spawn points")
	}

	if spawner.rng == nil {
		t.Error("Random number generator should be initialized")
	}
}

func TestNewSpawnPoint(t *testing.T) {
	pos := engine.Vector2{X: 100, Y: 200}
	sp := NewSpawnPoint(pos, EnemyGoomba, 5.0, 3)

	if sp == nil {
		t.Fatal("NewSpawnPoint returned nil")
	}

	if sp.Position != pos {
		t.Errorf("Expected position %v, got %v", pos, sp.Position)
	}

	if sp.EnemyType != EnemyGoomba {
		t.Errorf("Expected enemy type Goomba, got %v", sp.EnemyType)
	}

	if sp.SpawnRate != 5.0 {
		t.Errorf("Expected spawn rate 5.0, got %f", sp.SpawnRate)
	}

	if sp.MaxSpawns != 3 {
		t.Errorf("Expected max spawns 3, got %d", sp.MaxSpawns)
	}

	if sp.Spawned != 0 {
		t.Error("Should start with 0 spawned")
	}

	if !sp.Active {
		t.Error("Should be active by default")
	}
}

func TestSpawnPointCanSpawn(t *testing.T) {
	sp := NewSpawnPoint(engine.Vector2{X: 100, Y: 200}, EnemyGoomba, 5.0, 2)

	// Should be able to spawn initially
	if !sp.CanSpawn() {
		t.Error("Should be able to spawn initially")
	}

	// Spawn one
	sp.Spawned = 1
	if !sp.CanSpawn() {
		t.Error("Should still be able to spawn")
	}

	// Spawn two (max)
	sp.Spawned = 2
	if sp.CanSpawn() {
		t.Error("Should not be able to spawn after reaching max")
	}

	// Test inactive spawn point
	sp.Spawned = 0
	sp.Active = false
	if sp.CanSpawn() {
		t.Error("Inactive spawn point should not be able to spawn")
	}
}

func TestSpawnPointUnlimitedSpawns(t *testing.T) {
	sp := NewSpawnPoint(engine.Vector2{X: 100, Y: 200}, EnemyGoomba, 5.0, -1)

	sp.Spawned = 100
	if !sp.CanSpawn() {
		t.Error("Unlimited spawn point should always be able to spawn")
	}
}

func TestEnemySpawnerAddSpawnPoint(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)
	sp := NewSpawnPoint(engine.Vector2{X: 100, Y: 200}, EnemyGoomba, 5.0, 3)

	spawner.AddSpawnPoint(sp)

	if len(spawner.spawnPoints) != 1 {
		t.Errorf("Expected 1 spawn point, got %d", len(spawner.spawnPoints))
	}

	if spawner.spawnPoints[0] != sp {
		t.Error("Spawn point not added correctly")
	}
}

func TestEnemySpawnerRemoveSpawnPoint(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)
	sp1 := NewSpawnPoint(engine.Vector2{X: 100, Y: 200}, EnemyGoomba, 5.0, 3)
	sp2 := NewSpawnPoint(engine.Vector2{X: 200, Y: 200}, EnemyKoopa, 5.0, 3)

	spawner.AddSpawnPoint(sp1)
	spawner.AddSpawnPoint(sp2)

	spawner.RemoveSpawnPoint(sp1)

	if len(spawner.spawnPoints) != 1 {
		t.Errorf("Expected 1 spawn point, got %d", len(spawner.spawnPoints))
	}

	if spawner.spawnPoints[0] != sp2 {
		t.Error("Wrong spawn point removed")
	}
}

func TestEnemySpawnerSpawnImmediate(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)
	pos := engine.Vector2{X: 100, Y: 200}

	enemy := spawner.SpawnImmediate(EnemyGoomba, pos)

	if enemy == nil {
		t.Fatal("SpawnImmediate returned nil")
	}

	if enemy.GetType() != EnemyGoomba {
		t.Errorf("Expected Goomba, got %v", enemy.GetType())
	}

	if len(spawner.enemies) != 1 {
		t.Errorf("Expected 1 enemy, got %d", len(spawner.enemies))
	}

	if spawner.GetEnemyCount() != 1 {
		t.Errorf("Expected enemy count 1, got %d", spawner.GetEnemyCount())
	}
}

func TestEnemySpawnerSpawnImmediateMaxLimit(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 2)

	// Spawn up to max
	spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 100, Y: 200})
	spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 200, Y: 200})

	// Try to spawn beyond max
	enemy := spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 300, Y: 200})

	if enemy != nil {
		t.Error("Should not spawn beyond max enemies")
	}

	if len(spawner.enemies) != 2 {
		t.Errorf("Expected 2 enemies, got %d", len(spawner.enemies))
	}
}

func TestEnemySpawnerUpdate(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)
	sp := NewSpawnPoint(engine.Vector2{X: 100, Y: 200}, EnemyGoomba, 0.1, -1)
	spawner.AddSpawnPoint(sp)

	players := []Player{
		&MockPlayer{position: engine.Vector2{X: 300, Y: 200}, alive: true},
	}

	// Update multiple times to trigger spawn
	for i := 0; i < 10; i++ {
		spawner.Update(0.05, players) // 0.05 * 10 = 0.5 seconds
	}

	// Should have spawned at least one enemy
	if len(spawner.enemies) == 0 {
		t.Error("Should have spawned enemies")
	}
}

func TestEnemySpawnerRemoveDeadEnemies(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)

	// Spawn some enemies
	enemy1 := spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 100, Y: 200})
	enemy2 := spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 200, Y: 200})
	enemy3 := spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 300, Y: 200})

	if len(spawner.enemies) != 3 {
		t.Fatalf("Expected 3 enemies, got %d", len(spawner.enemies))
	}

	// Kill one enemy
	enemy2.Die()

	// Update to trigger cleanup
	players := []Player{}
	spawner.Update(0.016, players)

	// Should have removed dead enemy
	if len(spawner.enemies) != 2 {
		t.Errorf("Expected 2 enemies after cleanup, got %d", len(spawner.enemies))
	}

	// Verify remaining enemies are alive
	for _, enemy := range spawner.enemies {
		if !enemy.IsAlive() {
			t.Error("Dead enemy not removed")
		}
	}

	// Verify the correct enemies remain
	found1 := false
	found3 := false
	for _, enemy := range spawner.enemies {
		if enemy == enemy1 {
			found1 = true
		}
		if enemy == enemy3 {
			found3 = true
		}
	}

	if !found1 || !found3 {
		t.Error("Wrong enemies removed")
	}
}

func TestEnemySpawnerGetAliveEnemyCount(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)

	enemy1 := spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 100, Y: 200})
	enemy2 := spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 200, Y: 200})
	spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 300, Y: 200})

	if spawner.GetAliveEnemyCount() != 3 {
		t.Errorf("Expected 3 alive enemies, got %d", spawner.GetAliveEnemyCount())
	}

	// Kill one
	enemy1.Die()

	if spawner.GetAliveEnemyCount() != 2 {
		t.Errorf("Expected 2 alive enemies, got %d", spawner.GetAliveEnemyCount())
	}

	// Kill another
	enemy2.Die()

	if spawner.GetAliveEnemyCount() != 1 {
		t.Errorf("Expected 1 alive enemy, got %d", spawner.GetAliveEnemyCount())
	}
}

func TestEnemySpawnerKillAllEnemies(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)

	spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 100, Y: 200})
	spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 200, Y: 200})
	spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 300, Y: 200})

	spawner.KillAllEnemies()

	if spawner.GetAliveEnemyCount() != 0 {
		t.Errorf("Expected 0 alive enemies, got %d", spawner.GetAliveEnemyCount())
	}

	// All enemies should be dead
	for _, enemy := range spawner.enemies {
		if enemy.IsAlive() {
			t.Error("Enemy should be dead")
		}
	}
}

func TestEnemySpawnerClear(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)

	spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 100, Y: 200})
	spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 200, Y: 200})

	sp := NewSpawnPoint(engine.Vector2{X: 100, Y: 200}, EnemyGoomba, 5.0, 3)
	spawner.AddSpawnPoint(sp)

	spawner.Clear()

	if len(spawner.enemies) != 0 {
		t.Error("Should have no enemies after clear")
	}

	if len(spawner.spawnPoints) != 0 {
		t.Error("Should have no spawn points after clear")
	}

	if spawner.nextID != 1 {
		t.Error("Next ID should be reset")
	}
}

func TestEnemySpawnerGetEnemiesInRadius(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)

	spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 100, Y: 200})
	spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 150, Y: 200})
	spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 300, Y: 200})

	center := engine.Vector2{X: 100, Y: 200}
	enemies := spawner.GetEnemiesInRadius(center, 60)

	// Should find 2 enemies within radius 60
	if len(enemies) != 2 {
		t.Errorf("Expected 2 enemies in radius, got %d", len(enemies))
	}
}

func TestEnemySpawnerGetEnemiesByType(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)

	spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 100, Y: 200})
	spawner.SpawnImmediate(EnemyKoopa, engine.Vector2{X: 200, Y: 200})
	spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 300, Y: 200})

	goombas := spawner.GetEnemiesByType(EnemyGoomba)
	if len(goombas) != 2 {
		t.Errorf("Expected 2 Goombas, got %d", len(goombas))
	}

	koopas := spawner.GetEnemiesByType(EnemyKoopa)
	if len(koopas) != 1 {
		t.Errorf("Expected 1 Koopa, got %d", len(koopas))
	}

	lakitus := spawner.GetEnemiesByType(EnemyLakitu)
	if len(lakitus) != 0 {
		t.Errorf("Expected 0 Lakitus, got %d", len(lakitus))
	}
}

func TestEnemySpawnerSetMaxEnemies(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)

	spawner.SetMaxEnemies(20)

	if spawner.GetMaxEnemies() != 20 {
		t.Errorf("Expected max enemies 20, got %d", spawner.GetMaxEnemies())
	}
}

func TestEnemySpawnerSetSpawnInterval(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)

	spawner.SetSpawnInterval(5.0)

	if spawner.GetSpawnInterval() != 5.0 {
		t.Errorf("Expected spawn interval 5.0, got %f", spawner.GetSpawnInterval())
	}
}

func TestEnemySpawnerActivateDeactivateSpawnPoint(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)
	sp := NewSpawnPoint(engine.Vector2{X: 100, Y: 200}, EnemyGoomba, 5.0, 3)
	spawner.AddSpawnPoint(sp)

	if !sp.Active {
		t.Error("Spawn point should be active initially")
	}

	spawner.DeactivateSpawnPoint(0)

	if sp.Active {
		t.Error("Spawn point should be deactivated")
	}

	spawner.ActivateSpawnPoint(0)

	if !sp.Active {
		t.Error("Spawn point should be activated")
	}
}

func TestEnemySpawnerGetSpawnPoint(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)
	sp1 := NewSpawnPoint(engine.Vector2{X: 100, Y: 200}, EnemyGoomba, 5.0, 3)
	sp2 := NewSpawnPoint(engine.Vector2{X: 200, Y: 200}, EnemyKoopa, 5.0, 3)

	spawner.AddSpawnPoint(sp1)
	spawner.AddSpawnPoint(sp2)

	retrieved := spawner.GetSpawnPoint(0)
	if retrieved != sp1 {
		t.Error("Wrong spawn point retrieved")
	}

	retrieved = spawner.GetSpawnPoint(1)
	if retrieved != sp2 {
		t.Error("Wrong spawn point retrieved")
	}

	retrieved = spawner.GetSpawnPoint(2)
	if retrieved != nil {
		t.Error("Should return nil for invalid index")
	}
}

func TestEnemySpawnerGetSpawnPointCount(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)

	if spawner.GetSpawnPointCount() != 0 {
		t.Error("Should start with 0 spawn points")
	}

	spawner.AddSpawnPoint(NewSpawnPoint(engine.Vector2{X: 100, Y: 200}, EnemyGoomba, 5.0, 3))
	spawner.AddSpawnPoint(NewSpawnPoint(engine.Vector2{X: 200, Y: 200}, EnemyKoopa, 5.0, 3))

	if spawner.GetSpawnPointCount() != 2 {
		t.Errorf("Expected 2 spawn points, got %d", spawner.GetSpawnPointCount())
	}
}

func TestEnemySpawnerSpawnAtPointRespectMaxSpawns(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)
	sp := NewSpawnPoint(engine.Vector2{X: 100, Y: 200}, EnemyGoomba, 0.1, 2)
	spawner.AddSpawnPoint(sp)

	// Spawn twice (max)
	spawner.spawnAtPoint(sp)
	spawner.spawnAtPoint(sp)

	if sp.Spawned != 2 {
		t.Errorf("Expected 2 spawned, got %d", sp.Spawned)
	}

	// Try to spawn again (should fail)
	initialCount := len(spawner.enemies)
	spawner.spawnAtPoint(sp)

	if len(spawner.enemies) != initialCount {
		t.Error("Should not spawn beyond max spawns")
	}
}

func TestEnemySpawnerUniqueIDs(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	spawner := NewEnemySpawner(physics, 10)

	enemy1 := spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 100, Y: 200})
	enemy2 := spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 200, Y: 200})
	enemy3 := spawner.SpawnImmediate(EnemyGoomba, engine.Vector2{X: 300, Y: 200})

	// Cast to EnemyBase to access ID
	base1, ok1 := enemy1.(*EnemyBase)
	base2, ok2 := enemy2.(*EnemyBase)
	base3, ok3 := enemy3.(*EnemyBase)

	if !ok1 || !ok2 || !ok3 {
		t.Fatal("Failed to cast enemies to EnemyBase")
	}

	// IDs should be unique and sequential
	if base1.ID == base2.ID || base2.ID == base3.ID || base1.ID == base3.ID {
		t.Error("Enemy IDs should be unique")
	}

	if base1.ID != 1 || base2.ID != 2 || base3.ID != 3 {
		t.Error("Enemy IDs should be sequential starting from 1")
	}
}
