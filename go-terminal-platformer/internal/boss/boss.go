package boss

import (
	"math"
)

// BossType defines different types of bosses
type BossType int

const (
	BossGoombaKing BossType = iota
	BossSandSerpent
	BossIceGolem
	BossFireDragon
)

// Vector2 represents a 2D vector
type Vector2 struct {
	X float64
	Y float64
}

// Add adds two vectors
func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{X: v.X + other.X, Y: v.Y + other.Y}
}

// Multiply multiplies a vector by a scalar
func (v Vector2) Multiply(scalar float64) Vector2 {
	return Vector2{X: v.X * scalar, Y: v.Y * scalar}
}

// Length returns the length of the vector
func (v Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Rectangle represents a rectangular area
type Rectangle struct {
	X      float64
	Y      float64
	Width  float64
	Height float64
}

// Center returns the center point of the rectangle
func (r Rectangle) Center() Vector2 {
	return Vector2{
		X: r.X + r.Width/2,
		Y: r.Y + r.Height/2,
	}
}

// Entity interface for boss interactions
type Entity interface {
	GetID() int
	GetPosition() Vector2
	GetBounds() (float64, float64, float64, float64)
}

// AttackPattern defines a boss attack pattern
type AttackPattern struct {
	Name        string
	Duration    float64
	Cooldown    float64
	Projectiles int
	Damage      int
}

// BossPhase defines a phase in the boss fight
type BossPhase struct {
	HealthThreshold int
	MoveSpeed       float64
	AttackSpeed     float64
	Patterns        []AttackPattern
	Vulnerable      bool
}

// Boss represents a boss enemy
type Boss struct {
	ID            int
	BossType      BossType
	Position      Vector2
	Velocity      Vector2
	Width         float64
	Height        float64
	MaxHealth     int
	CurrentHealth int
	CurrentPhase  int
	Phases        []BossPhase
	Vulnerable    bool
	Attacking     bool
	Stunned       bool
	Defeated      bool
	PhaseTimer    float64
	AttackTimer   float64
	StunTimer     float64
	FlashTimer    float64
	CurrentPattern int
	AttackPatterns map[int][]AttackPattern
}

// GetID returns the boss ID
func (b *Boss) GetID() int {
	return b.ID
}

// GetPosition returns the boss position
func (b *Boss) GetPosition() Vector2 {
	return b.Position
}

// GetBounds returns the boss bounds
func (b *Boss) GetBounds() (float64, float64, float64, float64) {
	return b.Position.X, b.Position.Y, b.Width, b.Height
}

// TakeDamage applies damage to the boss
func (b *Boss) TakeDamage(damage int, attacker Entity) bool {
	if !b.Vulnerable || b.Stunned || b.Defeated {
		return false
	}

	b.CurrentHealth -= damage
	if b.CurrentHealth < 0 {
		b.CurrentHealth = 0
	}

	if b.CurrentHealth == 0 {
		b.Defeated = true
	}

	b.FlashTimer = 0.2

	return true
}

// GetCurrentPhase returns the current phase configuration
func (b *Boss) GetCurrentPhase() BossPhase {
	if b.CurrentPhase >= len(b.Phases) {
		return b.Phases[len(b.Phases)-1]
	}
	return b.Phases[b.CurrentPhase]
}

// GetCurrentPattern returns the current attack pattern
func (b *Boss) GetCurrentPattern() AttackPattern {
	patterns := b.AttackPatterns[b.CurrentPhase]
	if len(patterns) == 0 {
		return AttackPattern{}
	}
	return patterns[b.CurrentPattern%len(patterns)]
}

// BossArena represents the boss battle arena
type BossArena struct {
	Bounds    Rectangle
	Entrances []Vector2
	Exits     []Vector2
	Locked    bool
}

// BossManager manages boss battles
type BossManager struct {
	currentBoss *Boss
	arena       *BossArena
	active      bool
	nextID      int
}

// NewBossManager creates a new boss manager
func NewBossManager() *BossManager {
	return &BossManager{
		nextID: 1,
	}
}

// StartBossFight starts a boss battle
func (bm *BossManager) StartBossFight(bossType BossType, arena *BossArena) *Boss {
	bm.arena = arena
	bm.arena.Locked = true

	bm.currentBoss = bm.createBoss(bossType)

	// Position boss in arena center
	centerPos := arena.Bounds.Center()
	bm.currentBoss.Position = centerPos

	bm.active = true

	return bm.currentBoss
}

// createBoss creates a boss of the specified type
func (bm *BossManager) createBoss(bossType BossType) *Boss {
	boss := &Boss{
		ID:             bm.nextID,
		BossType:       bossType,
		MaxHealth:      10,
		CurrentHealth:  10,
		CurrentPhase:   0,
		Vulnerable:     true,
		Width:          32,
		Height:         32,
		AttackPatterns: make(map[int][]AttackPattern),
	}
	bm.nextID++

	switch bossType {
	case BossGoombaKing:
		boss.Phases = []BossPhase{
			{HealthThreshold: 7, MoveSpeed: 100, AttackSpeed: 2.0, Vulnerable: true},
			{HealthThreshold: 3, MoveSpeed: 150, AttackSpeed: 1.5, Vulnerable: true},
			{HealthThreshold: 0, MoveSpeed: 200, AttackSpeed: 1.0, Vulnerable: true},
		}
		boss.AttackPatterns[0] = []AttackPattern{
			{Name: "Stomp", Duration: 1.0, Cooldown: 3.0, Damage: 1},
			{Name: "Charge", Duration: 2.0, Cooldown: 5.0, Damage: 2},
		}
		boss.AttackPatterns[1] = []AttackPattern{
			{Name: "Stomp", Duration: 0.8, Cooldown: 2.0, Damage: 1},
			{Name: "Charge", Duration: 1.5, Cooldown: 3.0, Damage: 2},
			{Name: "Summon", Duration: 2.0, Cooldown: 8.0, Projectiles: 3},
		}
		boss.AttackPatterns[2] = []AttackPattern{
			{Name: "Rage", Duration: 3.0, Cooldown: 2.0, Damage: 3},
			{Name: "Charge", Duration: 1.0, Cooldown: 2.0, Damage: 2},
			{Name: "Summon", Duration: 1.5, Cooldown: 5.0, Projectiles: 5},
		}

	case BossSandSerpent:
		boss.Phases = []BossPhase{
			{HealthThreshold: 7, MoveSpeed: 120, AttackSpeed: 2.5, Vulnerable: false},
			{HealthThreshold: 3, MoveSpeed: 150, AttackSpeed: 2.0, Vulnerable: false},
			{HealthThreshold: 0, MoveSpeed: 180, AttackSpeed: 1.5, Vulnerable: false},
		}
		boss.AttackPatterns[0] = []AttackPattern{
			{Name: "Burrow", Duration: 2.0, Cooldown: 4.0},
			{Name: "SandBlast", Duration: 1.5, Cooldown: 3.0, Projectiles: 3, Damage: 1},
		}
		boss.AttackPatterns[1] = []AttackPattern{
			{Name: "Burrow", Duration: 1.5, Cooldown: 3.0},
			{Name: "SandBlast", Duration: 1.0, Cooldown: 2.0, Projectiles: 5, Damage: 1},
			{Name: "Sandstorm", Duration: 3.0, Cooldown: 8.0, Damage: 2},
		}
		boss.AttackPatterns[2] = []AttackPattern{
			{Name: "Burrow", Duration: 1.0, Cooldown: 2.0},
			{Name: "SandBlast", Duration: 0.8, Cooldown: 1.5, Projectiles: 7, Damage: 1},
			{Name: "Sandstorm", Duration: 2.5, Cooldown: 5.0, Damage: 2},
		}

	case BossIceGolem:
		boss.Phases = []BossPhase{
			{HealthThreshold: 7, MoveSpeed: 80, AttackSpeed: 3.0, Vulnerable: true},
			{HealthThreshold: 3, MoveSpeed: 100, AttackSpeed: 2.5, Vulnerable: true},
			{HealthThreshold: 0, MoveSpeed: 120, AttackSpeed: 2.0, Vulnerable: true},
		}
		boss.AttackPatterns[0] = []AttackPattern{
			{Name: "IcePunch", Duration: 1.2, Cooldown: 3.5, Damage: 2},
			{Name: "IceSpike", Duration: 1.5, Cooldown: 4.0, Projectiles: 3, Damage: 1},
		}
		boss.AttackPatterns[1] = []AttackPattern{
			{Name: "IcePunch", Duration: 1.0, Cooldown: 2.5, Damage: 2},
			{Name: "IceSpike", Duration: 1.2, Cooldown: 3.0, Projectiles: 5, Damage: 1},
			{Name: "Blizzard", Duration: 3.0, Cooldown: 8.0, Damage: 3},
		}
		boss.AttackPatterns[2] = []AttackPattern{
			{Name: "IcePunch", Duration: 0.8, Cooldown: 2.0, Damage: 2},
			{Name: "IceSpike", Duration: 1.0, Cooldown: 2.0, Projectiles: 7, Damage: 1},
			{Name: "Blizzard", Duration: 2.5, Cooldown: 5.0, Damage: 3},
		}

	case BossFireDragon:
		boss.Phases = []BossPhase{
			{HealthThreshold: 7, MoveSpeed: 150, AttackSpeed: 2.0, Vulnerable: true},
			{HealthThreshold: 3, MoveSpeed: 180, AttackSpeed: 1.5, Vulnerable: true},
			{HealthThreshold: 0, MoveSpeed: 220, AttackSpeed: 1.0, Vulnerable: true},
		}
		boss.AttackPatterns[0] = []AttackPattern{
			{Name: "Fireball", Duration: 1.0, Cooldown: 2.5, Projectiles: 1, Damage: 2},
			{Name: "Dive", Duration: 2.0, Cooldown: 5.0, Damage: 3},
		}
		boss.AttackPatterns[1] = []AttackPattern{
			{Name: "Fireball", Duration: 0.8, Cooldown: 2.0, Projectiles: 3, Damage: 2},
			{Name: "Dive", Duration: 1.5, Cooldown: 3.5, Damage: 3},
			{Name: "FlameBreath", Duration: 2.5, Cooldown: 7.0, Damage: 4},
		}
		boss.AttackPatterns[2] = []AttackPattern{
			{Name: "Fireball", Duration: 0.6, Cooldown: 1.5, Projectiles: 5, Damage: 2},
			{Name: "Dive", Duration: 1.0, Cooldown: 2.5, Damage: 3},
			{Name: "FlameBreath", Duration: 2.0, Cooldown: 4.0, Damage: 4},
		}
	}

	return boss
}

// Update updates the boss battle
func (bm *BossManager) Update(deltaTime float64) {
	if !bm.active || bm.currentBoss == nil {
		return
	}

	boss := bm.currentBoss

	if boss.Defeated {
		return
	}

	// Update timers
	boss.PhaseTimer += deltaTime
	boss.AttackTimer += deltaTime

	if boss.FlashTimer > 0 {
		boss.FlashTimer -= deltaTime
	}

	if boss.Stunned {
		boss.StunTimer -= deltaTime
		if boss.StunTimer <= 0 {
			boss.Stunned = false
			boss.Vulnerable = boss.GetCurrentPhase().Vulnerable
		}
		return
	}

	// Check phase transition
	bm.checkPhaseTransition()

	// Execute attack pattern
	currentPattern := boss.GetCurrentPattern()
	if boss.AttackTimer >= currentPattern.Cooldown {
		bm.executeAttackPattern()
		boss.AttackTimer = 0
	}
}

// checkPhaseTransition checks if boss should transition to next phase
func (bm *BossManager) checkPhaseTransition() {
	boss := bm.currentBoss
	currentPhase := boss.GetCurrentPhase()

	if boss.CurrentHealth <= currentPhase.HealthThreshold && boss.CurrentPhase < len(boss.Phases)-1 {
		boss.CurrentPhase++
		bm.onPhaseTransition()
	}
}

// onPhaseTransition handles phase transition
func (bm *BossManager) onPhaseTransition() {
	boss := bm.currentBoss

	// Brief invulnerability during transition
	boss.Vulnerable = false
	boss.Stunned = true
	boss.StunTimer = 2.0

	// Reset attack timer
	boss.AttackTimer = 0
	boss.CurrentPattern = 0
}

// executeAttackPattern executes the current attack pattern
func (bm *BossManager) executeAttackPattern() {
	boss := bm.currentBoss
	patterns := boss.AttackPatterns[boss.CurrentPhase]

	if len(patterns) == 0 {
		return
	}

	boss.CurrentPattern++
	boss.Attacking = true
}

// HandleBossDefeat handles boss defeat
func (bm *BossManager) HandleBossDefeat() {
	if bm.arena != nil {
		bm.arena.Locked = false
	}

	bm.active = false
}

// IsActive returns whether a boss fight is active
func (bm *BossManager) IsActive() bool {
	return bm.active
}

// GetCurrentBoss returns the current boss
func (bm *BossManager) GetCurrentBoss() *Boss {
	return bm.currentBoss
}

// GetArena returns the current arena
func (bm *BossManager) GetArena() *BossArena {
	return bm.arena
}

// IsBossDefeated returns whether the current boss is defeated
func (bm *BossManager) IsBossDefeated() bool {
	return bm.currentBoss != nil && bm.currentBoss.Defeated
}

// EndBossFight ends the current boss fight
func (bm *BossManager) EndBossFight() {
	bm.HandleBossDefeat()
	bm.currentBoss = nil
	bm.arena = nil
}
