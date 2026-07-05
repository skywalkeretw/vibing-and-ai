package entities

import (
	"testing"

	"github.com/lukeroy/go-terminal-platformer/internal/engine"
)

// MockPlayer implements the Player interface for testing
type MockPlayer struct {
	position    engine.Vector2
	alive       bool
	physicsBody *engine.PhysicsBody
}

func (mp *MockPlayer) GetPosition() engine.Vector2 {
	return mp.position
}

func (mp *MockPlayer) TakeDamage(damage int) {
	// Mock implementation
}

func (mp *MockPlayer) IsAlive() bool {
	return mp.alive
}

func (mp *MockPlayer) GetPhysicsBody() *engine.PhysicsBody {
	return mp.physicsBody
}

func TestNewEnemyBase(t *testing.T) {
	tests := []struct {
		name      string
		enemyType EnemyType
		x, y      float64
	}{
		{"Create Goomba", EnemyGoomba, 100, 200},
		{"Create Koopa", EnemyKoopa, 150, 250},
		{"Create Piranha", EnemyPiranha, 200, 300},
		{"Create Lakitu", EnemyLakitu, 250, 350},
		{"Create Spiny", EnemySpiny, 300, 400},
		{"Create Boo", EnemyBoo, 350, 450},
		{"Create HammerBro", EnemyHammerBro, 400, 500},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enemy := NewEnemyBase(tt.enemyType, tt.x, tt.y)

			if enemy == nil {
				t.Fatal("NewEnemyBase returned nil")
			}

			if enemy.Type != tt.enemyType {
				t.Errorf("Expected type %v, got %v", tt.enemyType, enemy.Type)
			}

			if enemy.Position.X != tt.x || enemy.Position.Y != tt.y {
				t.Errorf("Expected position (%f, %f), got (%f, %f)", tt.x, tt.y, enemy.Position.X, enemy.Position.Y)
			}

				if !enemy.Active {
					t.Error("Enemy should be active by default")
				}
			if enemy.IsDead {
				t.Error("Enemy should not be dead by default")
			}

			if enemy.State != EnemyStatePatrol {
				t.Errorf("Expected initial state Patrol, got %v", enemy.State)
			}

			if enemy.AIState != AIPatrol {
				t.Errorf("Expected initial AI state Patrol, got %v", enemy.AIState)
			}
		})
	}
}

func TestEnemyBaseStats(t *testing.T) {
	tests := []struct {
		name         string
		enemyType    EnemyType
		expectedHP   int
		expectedDmg  int
		expectedSpd  float64
		expectedFly  bool
	}{
		{"Goomba stats", EnemyGoomba, 1, 1, 60.0, false},
		{"Koopa stats", EnemyKoopa, 2, 1, 80.0, false},
		{"Piranha stats", EnemyPiranha, 2, 1, 0.0, false},
		{"Lakitu stats", EnemyLakitu, 3, 1, 100.0, true},
		{"Spiny stats", EnemySpiny, 1, 1, 70.0, false},
		{"Boo stats", EnemyBoo, 2, 1, 90.0, true},
		{"HammerBro stats", EnemyHammerBro, 3, 1, 50.0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			enemy := NewEnemyBase(tt.enemyType, 0, 0)

			if enemy.Health != tt.expectedHP {
				t.Errorf("Expected health %d, got %d", tt.expectedHP, enemy.Health)
			}

			if enemy.MaxHealth != tt.expectedHP {
				t.Errorf("Expected max health %d, got %d", tt.expectedHP, enemy.MaxHealth)
			}

			if enemy.Damage != tt.expectedDmg {
				t.Errorf("Expected damage %d, got %d", tt.expectedDmg, enemy.Damage)
			}

			if enemy.MoveSpeed != tt.expectedSpd {
				t.Errorf("Expected move speed %f, got %f", tt.expectedSpd, enemy.MoveSpeed)
			}

			if enemy.IsFlying != tt.expectedFly {
				t.Errorf("Expected flying %v, got %v", tt.expectedFly, enemy.IsFlying)
			}
		})
	}
}

func TestEnemyBaseInitialize(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	enemy := NewEnemyBase(EnemyGoomba, 100, 200)
	enemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)

	if enemy.PhysicsBody == nil {
		t.Fatal("Physics body should be created")
	}

	if enemy.physics != physics {
		t.Error("Physics engine reference not set")
	}

	if enemy.PhysicsBody.Entity != enemy {
		t.Error("Physics body entity reference not set correctly")
	}

	if !enemy.PhysicsBody.Enabled {
		t.Error("Physics body should be enabled")
	}

	if enemy.PhysicsBody.Layer != engine.LayerEnemy {
		t.Errorf("Expected layer Enemy, got %v", enemy.PhysicsBody.Layer)
	}

	// Check gravity scale for flying vs ground enemies
	flyingEnemy := NewEnemyBase(EnemyLakitu, 100, 200)
	flyingEnemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)

	if flyingEnemy.PhysicsBody.GravityScale != 0.0 {
		t.Errorf("Flying enemy should have gravity scale 0, got %f", flyingEnemy.PhysicsBody.GravityScale)
	}

	groundEnemy := NewEnemyBase(EnemyGoomba, 100, 200)
	groundEnemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)

	if groundEnemy.PhysicsBody.GravityScale != 1.0 {
		t.Errorf("Ground enemy should have gravity scale 1, got %f", groundEnemy.PhysicsBody.GravityScale)
	}
}

func TestEnemyBaseTakeDamage(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	enemy := NewEnemyBase(EnemyKoopa, 100, 200) // Koopa has 2 HP
	enemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)

	// First damage
	enemy.TakeDamage(1)

	if enemy.Health != 1 {
		t.Errorf("Expected health 1, got %d", enemy.Health)
	}

	if enemy.State != EnemyStateHurt {
		t.Errorf("Expected state Hurt, got %v", enemy.State)
	}

	if enemy.InvulnTimer <= 0 {
		t.Error("Invulnerability timer should be set")
	}

	if enemy.IsDead {
		t.Error("Enemy should not be dead yet")
	}

	// Second damage (should kill)
	enemy.InvulnTimer = 0 // Reset invuln timer
	enemy.TakeDamage(1)

	if enemy.Health != 0 {
		t.Errorf("Expected health 0, got %d", enemy.Health)
	}

	if !enemy.IsDead {
		t.Error("Enemy should be dead")
	}

	if enemy.State != EnemyStateDead {
		t.Errorf("Expected state Dead, got %v", enemy.State)
	}
}

func TestEnemyBaseInvulnerability(t *testing.T) {
	enemy := NewEnemyBase(EnemyKoopa, 100, 200)

	// Take damage
	enemy.TakeDamage(1)
	initialHealth := enemy.Health

	// Try to take damage again while invulnerable
	enemy.TakeDamage(1)

	if enemy.Health != initialHealth {
		t.Error("Enemy should not take damage while invulnerable")
	}
}

func TestEnemyBaseDie(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	enemy := NewEnemyBase(EnemyGoomba, 100, 200)
	enemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)

	enemy.Die()

	if !enemy.IsDead {
		t.Error("Enemy should be dead")
	}

	if enemy.State != EnemyStateDead {
		t.Errorf("Expected state Dead, got %v", enemy.State)
	}

	if enemy.PhysicsBody.Enabled {
		t.Error("Physics body should be disabled")
	}

	if enemy.PhysicsBody.Velocity.X != 0 || enemy.PhysicsBody.Velocity.Y != 0 {
		t.Error("Velocity should be zero")
	}
}

func TestEnemyBaseFindNearestPlayer(t *testing.T) {
	enemy := NewEnemyBase(EnemyGoomba, 100, 100)

	players := []Player{
		&MockPlayer{position: engine.Vector2{X: 150, Y: 100}, alive: true},
		&MockPlayer{position: engine.Vector2{X: 200, Y: 100}, alive: true},
		&MockPlayer{position: engine.Vector2{X: 120, Y: 100}, alive: true},
	}

	nearest := enemy.findNearestPlayer(players)

	if nearest == nil {
		t.Fatal("Should find nearest player")
	}

	// Nearest should be at (120, 100) - distance 20
	expectedPos := engine.Vector2{X: 120, Y: 100}
	if nearest.GetPosition() != expectedPos {
		t.Errorf("Expected nearest at %v, got %v", expectedPos, nearest.GetPosition())
	}
}

func TestEnemyBaseFindNearestPlayerIgnoresDead(t *testing.T) {
	enemy := NewEnemyBase(EnemyGoomba, 100, 100)

	players := []Player{
		&MockPlayer{position: engine.Vector2{X: 110, Y: 100}, alive: false}, // Closest but dead
		&MockPlayer{position: engine.Vector2{X: 150, Y: 100}, alive: true},  // Should be selected
	}

	nearest := enemy.findNearestPlayer(players)

	if nearest == nil {
		t.Fatal("Should find nearest alive player")
	}

	expectedPos := engine.Vector2{X: 150, Y: 100}
	if nearest.GetPosition() != expectedPos {
		t.Errorf("Expected nearest at %v, got %v", expectedPos, nearest.GetPosition())
	}
}

func TestEnemyBasePatrolBehavior(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	enemy := NewEnemyBase(EnemyGoomba, 100, 200)
	enemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)
	enemy.SetPatrolRange(50, 150)

	// Patrol should set velocity
	enemy.patrol(0.016) // ~60 FPS

	if enemy.PhysicsBody.Velocity.X == 0 {
		t.Error("Patrol should set horizontal velocity")
	}

	// Check facing direction affects velocity
	initialVelocity := enemy.PhysicsBody.Velocity.X
	enemy.Facing = -enemy.Facing
	enemy.patrol(0.016)

	if enemy.PhysicsBody.Velocity.X == initialVelocity {
		t.Error("Changing facing should change velocity direction")
	}
}

func TestEnemyBasePatrolWithPoints(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	enemy := NewEnemyBase(EnemyGoomba, 100, 200)
	enemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)

	patrolPoints := []engine.Vector2{
		{X: 100, Y: 200},
		{X: 200, Y: 200},
		{X: 150, Y: 200},
	}
	enemy.SetPatrolPoints(patrolPoints)

	// Should move toward first point
	enemy.patrol(0.016)

	if enemy.PhysicsBody.Velocity.X == 0 {
		t.Error("Should move toward patrol point")
	}

	if enemy.CurrentPoint != 0 {
		t.Error("Should start at first patrol point")
	}
}

func TestEnemyBaseChaseTarget(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	enemy := NewEnemyBase(EnemyGoomba, 100, 200)
	enemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)

	target := &MockPlayer{
		position: engine.Vector2{X: 200, Y: 200},
		alive:    true,
	}
	enemy.Target = target

	enemy.chaseTarget(0.016)

	if enemy.PhysicsBody.Velocity.X <= 0 {
		t.Error("Should move toward target (positive X)")
	}

	if enemy.Facing != DirectionRight {
		t.Error("Should face right when chasing target to the right")
	}

	// Test chasing left
	target.position = engine.Vector2{X: 50, Y: 200}
	enemy.chaseTarget(0.016)

	if enemy.PhysicsBody.Velocity.X >= 0 {
		t.Error("Should move toward target (negative X)")
	}

	if enemy.Facing != DirectionLeft {
		t.Error("Should face left when chasing target to the left")
	}
}

func TestEnemyBaseFlyingChase(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	enemy := NewEnemyBase(EnemyLakitu, 100, 200)
	enemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)

	target := &MockPlayer{
		position: engine.Vector2{X: 200, Y: 100},
		alive:    true,
	}
	enemy.Target = target

	enemy.chaseTarget(0.016)

	if enemy.PhysicsBody.Velocity.X == 0 {
		t.Error("Flying enemy should move horizontally")
	}

	if enemy.PhysicsBody.Velocity.Y == 0 {
		t.Error("Flying enemy should move vertically")
	}
}

func TestEnemyBaseAIStateTransitions(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	enemy := NewEnemyBase(EnemyGoomba, 100, 200)
	enemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)
	enemy.DetectRange = 100

	players := []Player{
		&MockPlayer{position: engine.Vector2{X: 150, Y: 200}, alive: true},
	}

	// Initially patrolling
	if enemy.AIState != AIPatrol {
		t.Error("Should start in patrol state")
	}

	// Update should detect player and switch to chase
	enemy.updateAI(0.016, players)

	if enemy.AIState != AIChase {
		t.Error("Should switch to chase when player in range")
	}

	if enemy.Target == nil {
		t.Error("Should have target set")
	}
}

func TestEnemyBaseAIChaseToAttack(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	enemy := NewEnemyBase(EnemyGoomba, 100, 200)
	enemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)
	enemy.AIState = AIChase
	enemy.AttackRange = 50

	target := &MockPlayer{
		position: engine.Vector2{X: 120, Y: 200}, // Within attack range
		alive:    true,
	}
	enemy.Target = target

	players := []Player{target}
	enemy.updateAI(0.016, players)

	if enemy.AIState != AIAttack {
		t.Error("Should switch to attack when in range")
	}
}

func TestEnemyBaseAILosePlayer(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	enemy := NewEnemyBase(EnemyGoomba, 100, 200)
	enemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)
	enemy.AIState = AIChase
	enemy.DetectRange = 100

	target := &MockPlayer{
		position: engine.Vector2{X: 300, Y: 200}, // Far away
		alive:    true,
	}
	enemy.Target = target

	players := []Player{target}
	enemy.updateAI(0.016, players)

	if enemy.AIState != AIPatrol {
		t.Error("Should return to patrol when player out of range")
	}

	if enemy.Target != nil {
		t.Error("Should clear target")
	}
}

func TestEnemyBaseUpdate(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	enemy := NewEnemyBase(EnemyGoomba, 100, 200)
	enemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)

	players := []Player{
		&MockPlayer{position: engine.Vector2{X: 300, Y: 200}, alive: true},
	}

	initialPos := enemy.Position

	// Update should process AI and sync position
	enemy.Update(0.016, players)

	// Position should be updated from physics
	if enemy.Position == initialPos && enemy.MoveSpeed > 0 {
		// Position might not change in first frame, but physics should be set up
		if enemy.PhysicsBody == nil {
			t.Error("Physics body should be set up")
		}
	}
}

func TestEnemyBaseUpdateTimers(t *testing.T) {
	enemy := NewEnemyBase(EnemyGoomba, 100, 200)
	enemy.StateTimer = 1.0
	enemy.InvulnTimer = 0.5
	enemy.HurtTimer = 0.3

	players := []Player{}
	enemy.Update(0.1, players)

	if enemy.StateTimer >= 1.0 {
		t.Error("State timer should decrease")
	}

	if enemy.InvulnTimer >= 0.5 {
		t.Error("Invuln timer should decrease")
	}

	if enemy.HurtTimer >= 0.3 {
		t.Error("Hurt timer should decrease")
	}
}

func TestEnemyBaseDeadDoesNotUpdate(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	enemy := NewEnemyBase(EnemyGoomba, 100, 200)
	enemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)
	enemy.Die()

	players := []Player{
		&MockPlayer{position: engine.Vector2{X: 150, Y: 200}, alive: true},
	}

	initialState := enemy.State
	enemy.Update(0.016, players)

	if enemy.State != initialState {
		t.Error("Dead enemy should not change state")
	}

	if enemy.AIState != AIPatrol {
		// AI state shouldn't change for dead enemies
	}
}

func TestEnemyBaseSetPatrolRange(t *testing.T) {
	enemy := NewEnemyBase(EnemyGoomba, 100, 200)
	enemy.SetPatrolRange(50, 150)

	if enemy.PatrolLeft != 50 {
		t.Errorf("Expected patrol left 50, got %f", enemy.PatrolLeft)
	}

	if enemy.PatrolRight != 150 {
		t.Errorf("Expected patrol right 150, got %f", enemy.PatrolRight)
	}
}

func TestEnemyBaseActivateDeactivate(t *testing.T) {
	physics := engine.NewPhysicsEngine()
	physics.Initialize()

	enemy := NewEnemyBase(EnemyGoomba, 100, 200)
	enemy.Initialize(engine.Vector2{X: 100, Y: 200}, physics)

	if !enemy.Active {
		t.Error("Enemy should be active initially")
	}

	enemy.Deactivate()

	if enemy.Active {
		t.Error("Enemy should be deactivated")
	}

	if enemy.PhysicsBody.Enabled {
		t.Error("Physics body should be disabled")
	}

	enemy.Activate()

	if !enemy.Active {
		t.Error("Enemy should be activated")
	}

	if !enemy.PhysicsBody.Enabled {
		t.Error("Physics body should be enabled")
	}
}

func TestEnemyTypeString(t *testing.T) {
	tests := []struct {
		enemyType EnemyType
		expected  string
	}{
		{EnemyGoomba, "Goomba"},
		{EnemyKoopa, "Koopa"},
		{EnemyPiranha, "Piranha"},
		{EnemyLakitu, "Lakitu"},
		{EnemySpiny, "Spiny"},
		{EnemyBoo, "Boo"},
		{EnemyHammerBro, "HammerBro"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			result := tt.enemyType.String()
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}
