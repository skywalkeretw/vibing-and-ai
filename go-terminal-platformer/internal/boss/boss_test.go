package boss

import (
	"testing"
)

// MockEntity for testing
type MockEntity struct {
	id     int
	x, y   float64
	w, h   float64
}

func (m *MockEntity) GetID() int {
	return m.id
}

func (m *MockEntity) GetPosition() Vector2 {
	return Vector2{X: m.x, Y: m.y}
}

func (m *MockEntity) GetBounds() (float64, float64, float64, float64) {
	return m.x, m.y, m.w, m.h
}

func TestVector2_Add(t *testing.T) {
	v1 := Vector2{X: 1, Y: 2}
	v2 := Vector2{X: 3, Y: 4}
	result := v1.Add(v2)

	if result.X != 4 || result.Y != 6 {
		t.Errorf("Expected (4, 6), got (%f, %f)", result.X, result.Y)
	}
}

func TestVector2_Multiply(t *testing.T) {
	v := Vector2{X: 2, Y: 3}
	result := v.Multiply(2.5)

	if result.X != 5 || result.Y != 7.5 {
		t.Errorf("Expected (5, 7.5), got (%f, %f)", result.X, result.Y)
	}
}

func TestVector2_Length(t *testing.T) {
	v := Vector2{X: 3, Y: 4}
	length := v.Length()

	if length != 5 {
		t.Errorf("Expected length 5, got %f", length)
	}
}

func TestRectangle_Center(t *testing.T) {
	rect := Rectangle{X: 10, Y: 20, Width: 100, Height: 80}
	center := rect.Center()

	expectedX := 60.0
	expectedY := 60.0

	if center.X != expectedX || center.Y != expectedY {
		t.Errorf("Expected center (%f, %f), got (%f, %f)", expectedX, expectedY, center.X, center.Y)
	}
}

func TestNewBossManager(t *testing.T) {
	bm := NewBossManager()

	if bm == nil {
		t.Fatal("NewBossManager returned nil")
	}

	if bm.IsActive() {
		t.Error("Boss manager should not be active initially")
	}

	if bm.GetCurrentBoss() != nil {
		t.Error("Should not have a current boss initially")
	}
}

func TestStartBossFight_GoombaKing(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}

	boss := bm.StartBossFight(BossGoombaKing, arena)

	if boss == nil {
		t.Fatal("StartBossFight returned nil boss")
	}

	if boss.BossType != BossGoombaKing {
		t.Errorf("Expected BossGoombaKing, got %v", boss.BossType)
	}

	if boss.MaxHealth != 10 {
		t.Errorf("Expected max health 10, got %d", boss.MaxHealth)
	}

	if boss.CurrentHealth != 10 {
		t.Errorf("Expected current health 10, got %d", boss.CurrentHealth)
	}

	if !boss.Vulnerable {
		t.Error("Boss should be vulnerable initially")
	}

	if boss.Defeated {
		t.Error("Boss should not be defeated initially")
	}

	if !bm.IsActive() {
		t.Error("Boss manager should be active after starting fight")
	}

	if !arena.Locked {
		t.Error("Arena should be locked during boss fight")
	}

	// Check boss is positioned at arena center
	expectedPos := arena.Bounds.Center()
	if boss.Position.X != expectedPos.X || boss.Position.Y != expectedPos.Y {
		t.Errorf("Expected boss at (%f, %f), got (%f, %f)",
			expectedPos.X, expectedPos.Y, boss.Position.X, boss.Position.Y)
	}
}

func TestStartBossFight_AllBossTypes(t *testing.T) {
	bossTypes := []BossType{
		BossGoombaKing,
		BossSandSerpent,
		BossIceGolem,
		BossFireDragon,
	}

	for _, bossType := range bossTypes {
		bm := NewBossManager()
		arena := &BossArena{
			Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
		}

		boss := bm.StartBossFight(bossType, arena)

		if boss == nil {
			t.Errorf("Failed to create boss of type %v", bossType)
			continue
		}

		if boss.BossType != bossType {
			t.Errorf("Expected boss type %v, got %v", bossType, boss.BossType)
		}

		if len(boss.Phases) != 3 {
			t.Errorf("Expected 3 phases for boss type %v, got %d", bossType, len(boss.Phases))
		}

		if len(boss.AttackPatterns) != 3 {
			t.Errorf("Expected 3 attack pattern sets for boss type %v, got %d", bossType, len(boss.AttackPatterns))
		}
	}
}

func TestBoss_TakeDamage(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	boss := bm.StartBossFight(BossGoombaKing, arena)
	attacker := &MockEntity{id: 1}

	initialHealth := boss.CurrentHealth
	hit := boss.TakeDamage(2, attacker)

	if !hit {
		t.Error("Damage should have been applied")
	}

	if boss.CurrentHealth != initialHealth-2 {
		t.Errorf("Expected health %d, got %d", initialHealth-2, boss.CurrentHealth)
	}

	if boss.FlashTimer <= 0 {
		t.Error("Flash timer should be set after taking damage")
	}
}

func TestBoss_TakeDamage_WhenInvulnerable(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	boss := bm.StartBossFight(BossGoombaKing, arena)
	attacker := &MockEntity{id: 1}

	boss.Vulnerable = false
	initialHealth := boss.CurrentHealth
	hit := boss.TakeDamage(2, attacker)

	if hit {
		t.Error("Damage should not be applied when invulnerable")
	}

	if boss.CurrentHealth != initialHealth {
		t.Error("Health should not change when invulnerable")
	}
}

func TestBoss_TakeDamage_WhenStunned(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	boss := bm.StartBossFight(BossGoombaKing, arena)
	attacker := &MockEntity{id: 1}

	boss.Stunned = true
	initialHealth := boss.CurrentHealth
	hit := boss.TakeDamage(2, attacker)

	if hit {
		t.Error("Damage should not be applied when stunned")
	}

	if boss.CurrentHealth != initialHealth {
		t.Error("Health should not change when stunned")
	}
}

func TestBoss_TakeDamage_ToDefeat(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	boss := bm.StartBossFight(BossGoombaKing, arena)
	attacker := &MockEntity{id: 1}

	// Deal enough damage to defeat boss
	boss.TakeDamage(10, attacker)

	if boss.CurrentHealth != 0 {
		t.Errorf("Expected health 0, got %d", boss.CurrentHealth)
	}

	if !boss.Defeated {
		t.Error("Boss should be defeated")
	}
}

func TestBoss_GetCurrentPhase(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	boss := bm.StartBossFight(BossGoombaKing, arena)

	phase := boss.GetCurrentPhase()

	if phase.HealthThreshold != 7 {
		t.Errorf("Expected health threshold 7, got %d", phase.HealthThreshold)
	}

	if phase.MoveSpeed != 100 {
		t.Errorf("Expected move speed 100, got %f", phase.MoveSpeed)
	}
}

func TestBoss_GetCurrentPattern(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	boss := bm.StartBossFight(BossGoombaKing, arena)

	pattern := boss.GetCurrentPattern()

	if pattern.Name != "Stomp" {
		t.Errorf("Expected pattern 'Stomp', got '%s'", pattern.Name)
	}
}

func TestBossManager_Update_TimerProgression(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	bm.StartBossFight(BossGoombaKing, arena)

	boss := bm.GetCurrentBoss()
	initialPhaseTimer := boss.PhaseTimer
	initialAttackTimer := boss.AttackTimer

	bm.Update(0.5)

	if boss.PhaseTimer <= initialPhaseTimer {
		t.Error("Phase timer should increase")
	}

	if boss.AttackTimer <= initialAttackTimer {
		t.Error("Attack timer should increase")
	}
}

func TestBossManager_Update_FlashTimer(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	boss := bm.StartBossFight(BossGoombaKing, arena)

	boss.FlashTimer = 0.5
	bm.Update(0.3)

	if boss.FlashTimer >= 0.5 {
		t.Error("Flash timer should decrease")
	}
}

func TestBossManager_Update_StunRecovery(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	boss := bm.StartBossFight(BossGoombaKing, arena)

	boss.Stunned = true
	boss.StunTimer = 0.5

	bm.Update(0.6)

	if boss.Stunned {
		t.Error("Boss should not be stunned after timer expires")
	}

	if boss.StunTimer > 0 {
		t.Error("Stun timer should be 0 or negative")
	}
}

func TestBossManager_PhaseTransition(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	boss := bm.StartBossFight(BossGoombaKing, arena)
	attacker := &MockEntity{id: 1}

	// Damage boss to trigger phase transition (health threshold is 7)
	boss.TakeDamage(4, attacker) // Health: 10 -> 6

	bm.Update(0.1)

	if boss.CurrentPhase != 1 {
		t.Errorf("Expected phase 1, got %d", boss.CurrentPhase)
	}

	if !boss.Stunned {
		t.Error("Boss should be stunned during phase transition")
	}

	if boss.StunTimer <= 0 {
		t.Error("Stun timer should be set during phase transition")
	}
}

func TestBossManager_AttackPatternExecution(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	boss := bm.StartBossFight(BossGoombaKing, arena)

	initialPattern := boss.CurrentPattern

	// Set attack timer to trigger attack
	pattern := boss.GetCurrentPattern()
	boss.AttackTimer = pattern.Cooldown + 0.1

	bm.Update(0.1)

	if boss.CurrentPattern == initialPattern {
		t.Error("Current pattern should have advanced")
	}

	if !boss.Attacking {
		t.Error("Boss should be attacking")
	}
}

func TestBossManager_HandleBossDefeat(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	bm.StartBossFight(BossGoombaKing, arena)

	bm.HandleBossDefeat()

	if arena.Locked {
		t.Error("Arena should be unlocked after boss defeat")
	}

	if bm.IsActive() {
		t.Error("Boss manager should not be active after defeat")
	}
}

func TestBossManager_IsBossDefeated(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	boss := bm.StartBossFight(BossGoombaKing, arena)

	if bm.IsBossDefeated() {
		t.Error("Boss should not be defeated initially")
	}

	boss.Defeated = true

	if !bm.IsBossDefeated() {
		t.Error("Boss should be defeated")
	}
}

func TestBossManager_EndBossFight(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	bm.StartBossFight(BossGoombaKing, arena)

	bm.EndBossFight()

	if bm.GetCurrentBoss() != nil {
		t.Error("Should not have a current boss after ending fight")
	}

	if bm.GetArena() != nil {
		t.Error("Should not have an arena after ending fight")
	}

	if bm.IsActive() {
		t.Error("Boss manager should not be active after ending fight")
	}
}

func TestBossManager_UpdateWhenInactive(t *testing.T) {
	bm := NewBossManager()

	// Should not panic when updating inactive manager
	bm.Update(0.1)

	if bm.IsActive() {
		t.Error("Boss manager should remain inactive")
	}
}

func TestBossManager_UpdateWhenDefeated(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	boss := bm.StartBossFight(BossGoombaKing, arena)

	boss.Defeated = true
	initialPhaseTimer := boss.PhaseTimer

	bm.Update(0.5)

	// Timers should not update when defeated
	if boss.PhaseTimer != initialPhaseTimer {
		t.Error("Phase timer should not update when boss is defeated")
	}
}

func TestBoss_GetBounds(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}
	boss := bm.StartBossFight(BossGoombaKing, arena)

	x, y, w, h := boss.GetBounds()

	if x != boss.Position.X || y != boss.Position.Y {
		t.Error("Bounds position should match boss position")
	}

	if w != boss.Width || h != boss.Height {
		t.Error("Bounds size should match boss size")
	}
}

func TestBossPhases_Configuration(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}

	testCases := []struct {
		bossType      BossType
		expectedPhases int
	}{
		{BossGoombaKing, 3},
		{BossSandSerpent, 3},
		{BossIceGolem, 3},
		{BossFireDragon, 3},
	}

	for _, tc := range testCases {
		boss := bm.StartBossFight(tc.bossType, arena)

		if len(boss.Phases) != tc.expectedPhases {
			t.Errorf("Boss type %v: expected %d phases, got %d",
				tc.bossType, tc.expectedPhases, len(boss.Phases))
		}

		// Verify phases are ordered by health threshold (descending)
		for i := 0; i < len(boss.Phases)-1; i++ {
			if boss.Phases[i].HealthThreshold <= boss.Phases[i+1].HealthThreshold {
				t.Errorf("Boss type %v: phases not ordered correctly", tc.bossType)
			}
		}

		bm.EndBossFight()
	}
}

func TestAttackPatterns_AllBosses(t *testing.T) {
	bm := NewBossManager()
	arena := &BossArena{
		Bounds: Rectangle{X: 0, Y: 0, Width: 400, Height: 300},
	}

	bossTypes := []BossType{
		BossGoombaKing,
		BossSandSerpent,
		BossIceGolem,
		BossFireDragon,
	}

	for _, bossType := range bossTypes {
		boss := bm.StartBossFight(bossType, arena)

		// Each phase should have attack patterns
		for phase := 0; phase < len(boss.Phases); phase++ {
			patterns := boss.AttackPatterns[phase]
			if len(patterns) == 0 {
				t.Errorf("Boss type %v phase %d has no attack patterns", bossType, phase)
			}

			// Verify each pattern has required fields
			for _, pattern := range patterns {
				if pattern.Name == "" {
					t.Errorf("Boss type %v phase %d has pattern with empty name", bossType, phase)
				}
				if pattern.Duration <= 0 {
					t.Errorf("Boss type %v phase %d pattern %s has invalid duration", bossType, phase, pattern.Name)
				}
				if pattern.Cooldown <= 0 {
					t.Errorf("Boss type %v phase %d pattern %s has invalid cooldown", bossType, phase, pattern.Name)
				}
			}
		}

		bm.EndBossFight()
	}
}
