package engine

import (
	"math"
	"testing"
)

func TestAABBColliderGetBounds(t *testing.T) {
	collider := NewAABBCollider(10, 20, 30, 40, LayerPlayer)
	bounds := collider.GetBounds()

	if bounds.X != 10 || bounds.Y != 20 || bounds.Width != 30 || bounds.Height != 40 {
		t.Errorf("GetBounds() failed: expected (10, 20, 30, 40), got (%f, %f, %f, %f)",
			bounds.X, bounds.Y, bounds.Width, bounds.Height)
	}
}

func TestAABBColliderIntersects(t *testing.T) {
	c1 := NewAABBCollider(0, 0, 10, 10, LayerPlayer)
	c2 := NewAABBCollider(5, 5, 10, 10, LayerEnemy)
	c3 := NewAABBCollider(20, 20, 10, 10, LayerEnemy)

	// Overlapping colliders
	if !c1.Intersects(c2) {
		t.Error("Intersects() failed: c1 and c2 should intersect")
	}

	// Non-overlapping colliders
	if c1.Intersects(c3) {
		t.Error("Intersects() failed: c1 and c3 should not intersect")
	}
}

func TestAABBColliderGetType(t *testing.T) {
	collider := NewAABBCollider(0, 0, 10, 10, LayerPlayer)
	if collider.GetType() != ColliderAABB {
		t.Error("GetType() should return ColliderAABB")
	}
}

func TestAABBColliderGetLayer(t *testing.T) {
	collider := NewAABBCollider(0, 0, 10, 10, LayerPlayer)
	if collider.GetLayer() != LayerPlayer {
		t.Error("GetLayer() should return LayerPlayer")
	}
}

func TestAABBColliderSetPosition(t *testing.T) {
	collider := NewAABBCollider(0, 0, 10, 10, LayerPlayer)
	collider.SetPosition(50, 60)

	if collider.X != 50 || collider.Y != 60 {
		t.Errorf("SetPosition() failed: expected (50, 60), got (%f, %f)", collider.X, collider.Y)
	}
}

func TestTileColliderGetBounds(t *testing.T) {
	collider := NewTileCollider(10, 20, 16, true, LayerTerrain)
	bounds := collider.GetBounds()

	if bounds.X != 10 || bounds.Y != 20 || bounds.Width != 16 || bounds.Height != 16 {
		t.Errorf("GetBounds() failed: expected (10, 20, 16, 16), got (%f, %f, %f, %f)",
			bounds.X, bounds.Y, bounds.Width, bounds.Height)
	}
}

func TestTileColliderIntersects(t *testing.T) {
	solidTile := NewTileCollider(0, 0, 16, true, LayerTerrain)
	nonSolidTile := NewTileCollider(0, 0, 16, false, LayerTerrain)
	player := NewAABBCollider(5, 5, 10, 10, LayerPlayer)

	// Solid tile should intersect
	if !solidTile.Intersects(player) {
		t.Error("Solid tile should intersect with player")
	}

	// Non-solid tile should not intersect
	if nonSolidTile.Intersects(player) {
		t.Error("Non-solid tile should not intersect")
	}
}

func TestTileColliderGetType(t *testing.T) {
	collider := NewTileCollider(0, 0, 16, true, LayerTerrain)
	if collider.GetType() != ColliderTile {
		t.Error("GetType() should return ColliderTile")
	}
}

func TestOneWayColliderGetBounds(t *testing.T) {
	collider := NewOneWayCollider(10, 20, 30, 5, LayerOneWayPlatform)
	bounds := collider.GetBounds()

	if bounds.X != 10 || bounds.Y != 20 || bounds.Width != 30 || bounds.Height != 5 {
		t.Errorf("GetBounds() failed: expected (10, 20, 30, 5), got (%f, %f, %f, %f)",
			bounds.X, bounds.Y, bounds.Width, bounds.Height)
	}
}

func TestOneWayColliderGetType(t *testing.T) {
	collider := NewOneWayCollider(0, 0, 10, 5, LayerOneWayPlatform)
	if collider.GetType() != ColliderOneWay {
		t.Error("GetType() should return ColliderOneWay")
	}
}

func TestOneWayColliderShouldCollideFromAbove(t *testing.T) {
	platform := NewOneWayCollider(0, 100, 50, 5, LayerOneWayPlatform)

	// Entity above platform, moving down - should collide
	if !platform.ShouldCollideFromAbove(95, 10) {
		t.Error("Should collide when entity is above and moving down")
	}

	// Entity below platform - should not collide
	if platform.ShouldCollideFromAbove(110, 10) {
		t.Error("Should not collide when entity is below platform")
	}

	// Entity moving up - should not collide
	if platform.ShouldCollideFromAbove(95, -10) {
		t.Error("Should not collide when entity is moving up")
	}
}

func TestNewCollision(t *testing.T) {
	c1 := NewAABBCollider(0, 0, 10, 10, LayerPlayer)
	c2 := NewAABBCollider(5, 5, 10, 10, LayerEnemy)

	collision := NewCollision(c1, c2)

	if collision.ColliderA != c1 || collision.ColliderB != c2 {
		t.Error("NewCollision() failed to set colliders correctly")
	}

	// Check that collision data was calculated
	if collision.Penetration == 0 {
		t.Error("Collision penetration should be calculated")
	}
}

func TestCollisionCalculateData(t *testing.T) {
	// Collision from left
	c1 := NewAABBCollider(0, 0, 10, 10, LayerPlayer)
	c2 := NewAABBCollider(8, 0, 10, 10, LayerEnemy)

	collision := NewCollision(c1, c2)

	// Should have horizontal collision
	if math.Abs(collision.Normal.Y) > 0.1 {
		t.Error("Collision should be horizontal")
	}

	// Penetration should be 2 (10 - 8)
	if math.Abs(collision.Penetration-2) > 0.1 {
		t.Errorf("Expected penetration of 2, got %f", collision.Penetration)
	}
}

func TestCollisionGetters(t *testing.T) {
	c1 := NewAABBCollider(0, 0, 10, 10, LayerPlayer)
	c2 := NewAABBCollider(5, 5, 10, 10, LayerEnemy)

	collision := NewCollision(c1, c2)

	normal := collision.GetNormal()
	if normal.X == 0 && normal.Y == 0 {
		t.Error("GetNormal() should return non-zero normal")
	}

	penetration := collision.GetPenetration()
	if penetration == 0 {
		t.Error("GetPenetration() should return non-zero penetration")
	}

	contactPoint := collision.GetContactPoint()
	if contactPoint.X == 0 && contactPoint.Y == 0 {
		t.Error("GetContactPoint() should return valid contact point")
	}
}

func TestAABBIntersects(t *testing.T) {
	r1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	r2 := Rectangle{X: 5, Y: 5, Width: 10, Height: 10}
	r3 := Rectangle{X: 20, Y: 20, Width: 10, Height: 10}

	// Overlapping rectangles
	if !AABBIntersects(r1, r2) {
		t.Error("AABBIntersects() failed: r1 and r2 should intersect")
	}

	// Non-overlapping rectangles
	if AABBIntersects(r1, r3) {
		t.Error("AABBIntersects() failed: r1 and r3 should not intersect")
	}

	// Edge touching (should intersect)
	r4 := Rectangle{X: 10, Y: 0, Width: 10, Height: 10}
	if !AABBIntersects(r1, r4) {
		t.Error("AABBIntersects() failed: edge touching should intersect")
	}
}

func TestGetAABBPenetration(t *testing.T) {
	// Horizontal collision
	r1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	r2 := Rectangle{X: 8, Y: 0, Width: 10, Height: 10}

	penetration := GetAABBPenetration(r1, r2)

	// Should push r1 to the left
	if penetration.X >= 0 {
		t.Error("Horizontal penetration should be negative (push left)")
	}
	if penetration.Y != 0 {
		t.Error("Vertical penetration should be zero for horizontal collision")
	}

	// Vertical collision
	r3 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	r4 := Rectangle{X: 0, Y: 8, Width: 10, Height: 10}

	penetration = GetAABBPenetration(r3, r4)

	// Should push r3 up
	if penetration.Y >= 0 {
		t.Error("Vertical penetration should be negative (push up)")
	}
	if penetration.X != 0 {
		t.Error("Horizontal penetration should be zero for vertical collision")
	}
}

func TestShouldCheckCollision(t *testing.T) {
	// Player should collide with terrain
	if !ShouldCheckCollision(LayerPlayer, LayerTerrain) {
		t.Error("Player should collide with terrain")
	}

	// Player should collide with enemies
	if !ShouldCheckCollision(LayerPlayer, LayerEnemy) {
		t.Error("Player should collide with enemies")
	}

	// Player should collide with collectibles
	if !ShouldCheckCollision(LayerPlayer, LayerCollectible) {
		t.Error("Player should collide with collectibles")
	}

	// Player should collide with one-way platforms
	if !ShouldCheckCollision(LayerPlayer, LayerOneWayPlatform) {
		t.Error("Player should collide with one-way platforms")
	}

	// Enemies should collide with terrain
	if !ShouldCheckCollision(LayerEnemy, LayerTerrain) {
		t.Error("Enemies should collide with terrain")
	}

	// Enemies should collide with projectiles
	if !ShouldCheckCollision(LayerEnemy, LayerProjectile) {
		t.Error("Enemies should collide with projectiles")
	}

	// Projectiles should collide with terrain
	if !ShouldCheckCollision(LayerProjectile, LayerTerrain) {
		t.Error("Projectiles should collide with terrain")
	}

	// Enemies should not collide with other enemies
	if ShouldCheckCollision(LayerEnemy, LayerEnemy) {
		t.Error("Enemies should not collide with other enemies")
	}

	// Collectibles should not collide with terrain
	if ShouldCheckCollision(LayerCollectible, LayerTerrain) {
		t.Error("Collectibles should not collide with terrain")
	}
}

func TestRaycast(t *testing.T) {
	colliders := []Collider{
		NewAABBCollider(50, 50, 20, 20, LayerTerrain),
		NewAABBCollider(100, 50, 20, 20, LayerEnemy),
	}

	// Raycast that hits the first collider
	origin := Vector2{X: 0, Y: 60}
	direction := Vector2{X: 1, Y: 0}
	hit := Raycast(origin, direction, 100, colliders, LayerTerrain)

	if !hit.Hit {
		t.Error("Raycast should hit the terrain collider")
	}

	if hit.Collider != colliders[0] {
		t.Error("Raycast should hit the first collider")
	}

	if hit.Distance <= 0 || hit.Distance > 100 {
		t.Errorf("Hit distance should be positive and less than max distance, got %f", hit.Distance)
	}

	// Raycast that misses
	origin = Vector2{X: 0, Y: 0}
	direction = Vector2{X: 0, Y: -1}
	hit = Raycast(origin, direction, 100, colliders, LayerTerrain)

	if hit.Hit {
		t.Error("Raycast should miss when pointing away from colliders")
	}

	// Raycast with layer mask that excludes all colliders
	origin = Vector2{X: 0, Y: 60}
	direction = Vector2{X: 1, Y: 0}
	hit = Raycast(origin, direction, 100, colliders, LayerCollectible)

	if hit.Hit {
		t.Error("Raycast should miss when layer mask excludes all colliders")
	}
}

func TestRaycastHitNormal(t *testing.T) {
	colliders := []Collider{
		NewAABBCollider(50, 50, 20, 20, LayerTerrain),
	}

	// Raycast from left
	origin := Vector2{X: 0, Y: 60}
	direction := Vector2{X: 1, Y: 0}
	hit := Raycast(origin, direction, 100, colliders, LayerTerrain)

	if !hit.Hit {
		t.Fatal("Raycast should hit")
	}

	// Normal should point left (opposite of ray direction)
	if hit.Normal.X >= 0 {
		t.Errorf("Normal should point left, got (%f, %f)", hit.Normal.X, hit.Normal.Y)
	}
}

func TestRaycastMultipleColliders(t *testing.T) {
	colliders := []Collider{
		NewAABBCollider(50, 50, 20, 20, LayerTerrain),
		NewAABBCollider(30, 50, 20, 20, LayerTerrain),
	}

	// Raycast that could hit both, should hit the closest
	origin := Vector2{X: 0, Y: 60}
	direction := Vector2{X: 1, Y: 0}
	hit := Raycast(origin, direction, 100, colliders, LayerTerrain)

	if !hit.Hit {
		t.Fatal("Raycast should hit")
	}

	// Should hit the second collider (closer)
	if hit.Collider != colliders[1] {
		t.Error("Raycast should hit the closest collider")
	}
}

func TestCollisionLayers(t *testing.T) {
	// Test layer constants are unique
	layers := []CollisionLayer{
		LayerPlayer,
		LayerEnemy,
		LayerTerrain,
		LayerProjectile,
		LayerCollectible,
		LayerOneWayPlatform,
	}

	for i, layer1 := range layers {
		for j, layer2 := range layers {
			if i != j && layer1 == layer2 {
				t.Errorf("Collision layers %d and %d have the same value", i, j)
			}
		}
	}

	// Test that layers are powers of 2 (bitwise flags)
	for i, layer := range layers {
		if layer&(layer-1) != 0 && layer != 0 {
			t.Errorf("Layer %d is not a power of 2", i)
		}
	}
}
