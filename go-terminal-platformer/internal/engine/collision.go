package engine

import (
	"math"
)

// ColliderType represents the type of collider
type ColliderType int

const (
	ColliderAABB ColliderType = iota
	ColliderTile
	ColliderOneWay
)

// CollisionLayer represents collision layers using bitwise flags
type CollisionLayer int

const (
	LayerPlayer CollisionLayer = 1 << iota
	LayerEnemy
	LayerTerrain
	LayerProjectile
	LayerCollectible
	LayerOneWayPlatform
)

// Collider interface defines the contract for all collider types
type Collider interface {
	GetBounds() Rectangle
	Intersects(other Collider) bool
	GetType() ColliderType
	GetLayer() CollisionLayer
}

// AABBCollider represents an axis-aligned bounding box collider
type AABBCollider struct {
	X, Y          float64
	Width, Height float64
	Layer         CollisionLayer
}

// NewAABBCollider creates a new AABB collider
func NewAABBCollider(x, y, width, height float64, layer CollisionLayer) *AABBCollider {
	return &AABBCollider{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		Layer:  layer,
	}
}

// GetBounds returns the bounding rectangle
func (c *AABBCollider) GetBounds() Rectangle {
	return Rectangle{
		X:      c.X,
		Y:      c.Y,
		Width:  c.Width,
		Height: c.Height,
	}
}

// Intersects checks if this collider intersects with another
func (c *AABBCollider) Intersects(other Collider) bool {
	return AABBIntersects(c.GetBounds(), other.GetBounds())
}

// GetType returns the collider type
func (c *AABBCollider) GetType() ColliderType {
	return ColliderAABB
}

// GetLayer returns the collision layer
func (c *AABBCollider) GetLayer() CollisionLayer {
	return c.Layer
}

// SetPosition updates the collider position
func (c *AABBCollider) SetPosition(x, y float64) {
	c.X = x
	c.Y = y
}

// TileCollider represents a tile-based collider
type TileCollider struct {
	X, Y     float64
	TileSize float64
	Solid    bool
	Layer    CollisionLayer
}

// NewTileCollider creates a new tile collider
func NewTileCollider(x, y, tileSize float64, solid bool, layer CollisionLayer) *TileCollider {
	return &TileCollider{
		X:        x,
		Y:        y,
		TileSize: tileSize,
		Solid:    solid,
		Layer:    layer,
	}
}

// GetBounds returns the bounding rectangle
func (c *TileCollider) GetBounds() Rectangle {
	return Rectangle{
		X:      c.X,
		Y:      c.Y,
		Width:  c.TileSize,
		Height: c.TileSize,
	}
}

// Intersects checks if this collider intersects with another
func (c *TileCollider) Intersects(other Collider) bool {
	if !c.Solid {
		return false
	}
	return AABBIntersects(c.GetBounds(), other.GetBounds())
}

// GetType returns the collider type
func (c *TileCollider) GetType() ColliderType {
	return ColliderTile
}

// GetLayer returns the collision layer
func (c *TileCollider) GetLayer() CollisionLayer {
	return c.Layer
}

// OneWayCollider represents a one-way platform collider
type OneWayCollider struct {
	X, Y          float64
	Width, Height float64
	Layer         CollisionLayer
}

// NewOneWayCollider creates a new one-way platform collider
func NewOneWayCollider(x, y, width, height float64, layer CollisionLayer) *OneWayCollider {
	return &OneWayCollider{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
		Layer:  layer,
	}
}

// GetBounds returns the bounding rectangle
func (c *OneWayCollider) GetBounds() Rectangle {
	return Rectangle{
		X:      c.X,
		Y:      c.Y,
		Width:  c.Width,
		Height: c.Height,
	}
}

// Intersects checks if this collider intersects with another
// One-way platforms only collide from above
func (c *OneWayCollider) Intersects(other Collider) bool {
	return AABBIntersects(c.GetBounds(), other.GetBounds())
}

// GetType returns the collider type
func (c *OneWayCollider) GetType() ColliderType {
	return ColliderOneWay
}

// GetLayer returns the collision layer
func (c *OneWayCollider) GetLayer() CollisionLayer {
	return c.Layer
}

// ShouldCollideFromAbove checks if an entity should collide with this one-way platform
func (c *OneWayCollider) ShouldCollideFromAbove(entityBottom, entityVelocityY float64) bool {
	// Only collide if entity is moving downward and its bottom is above or at the platform top
	tolerance := 5.0 // 5 pixel tolerance
	return entityVelocityY > 0 && entityBottom <= c.Y+tolerance
}

// Collision represents a collision between two colliders
type Collision struct {
	ColliderA    Collider
	ColliderB    Collider
	Normal       Vector2
	Penetration  float64
	ContactPoint Vector2
}

// NewCollision creates a new collision
func NewCollision(a, b Collider) *Collision {
	collision := &Collision{
		ColliderA: a,
		ColliderB: b,
	}
	collision.calculateCollisionData()
	return collision
}

// calculateCollisionData calculates the collision normal and penetration
func (c *Collision) calculateCollisionData() {
	boundsA := c.ColliderA.GetBounds()
	boundsB := c.ColliderB.GetBounds()

	// Calculate overlap on each axis
	overlapX := Min(boundsA.Right(), boundsB.Right()) - Max(boundsA.Left(), boundsB.Left())
	overlapY := Min(boundsA.Bottom(), boundsB.Bottom()) - Max(boundsA.Top(), boundsB.Top())

	// Determine collision normal and penetration based on smallest overlap
	if overlapX < overlapY {
		c.Penetration = overlapX
		if boundsA.Center().X < boundsB.Center().X {
			c.Normal = Vector2{-1, 0} // Collision from left
		} else {
			c.Normal = Vector2{1, 0} // Collision from right
		}
	} else {
		c.Penetration = overlapY
		if boundsA.Center().Y < boundsB.Center().Y {
			c.Normal = Vector2{0, -1} // Collision from top
		} else {
			c.Normal = Vector2{0, 1} // Collision from bottom
		}
	}

	// Calculate contact point (center of overlap)
	c.ContactPoint = Vector2{
		X: Max(boundsA.Left(), boundsB.Left()) + overlapX/2,
		Y: Max(boundsA.Top(), boundsB.Top()) + overlapY/2,
	}
}

// GetNormal returns the collision normal
func (c *Collision) GetNormal() Vector2 {
	return c.Normal
}

// GetPenetration returns the penetration depth
func (c *Collision) GetPenetration() float64 {
	return c.Penetration
}

// GetContactPoint returns the contact point
func (c *Collision) GetContactPoint() Vector2 {
	return c.ContactPoint
}

// AABBIntersects checks if two axis-aligned bounding boxes intersect
func AABBIntersects(a, b Rectangle) bool {
	return a.X < b.X+b.Width &&
		a.X+a.Width > b.X &&
		a.Y < b.Y+b.Height &&
		a.Y+a.Height > b.Y
}

// GetAABBPenetration calculates the minimum translation vector to separate two AABBs
func GetAABBPenetration(a, b Rectangle) Vector2 {
	// Calculate overlap on each axis
	overlapX := Min(a.Right(), b.Right()) - Max(a.Left(), b.Left())
	overlapY := Min(a.Bottom(), b.Bottom()) - Max(a.Top(), b.Top())

	// Return the minimum translation vector
	if overlapX < overlapY {
		if a.Center().X < b.Center().X {
			return Vector2{-overlapX, 0}
		}
		return Vector2{overlapX, 0}
	} else {
		if a.Center().Y < b.Center().Y {
			return Vector2{0, -overlapY}
		}
		return Vector2{0, overlapY}
	}
}

// ShouldCheckCollision checks if two layers should collide
func ShouldCheckCollision(layerA, layerB CollisionLayer) bool {
	// Define collision matrix
	// Players collide with: Terrain, Enemies, Collectibles, OneWayPlatforms
	// Enemies collide with: Terrain, Players, Projectiles
	// Projectiles collide with: Terrain, Enemies
	// Collectibles collide with: Players

	if layerA == LayerPlayer {
		return layerB == LayerTerrain ||
			layerB == LayerEnemy ||
			layerB == LayerCollectible ||
			layerB == LayerOneWayPlatform
	}

	if layerA == LayerEnemy {
		return layerB == LayerTerrain ||
			layerB == LayerPlayer ||
			layerB == LayerProjectile
	}

	if layerA == LayerProjectile {
		return layerB == LayerTerrain ||
			layerB == LayerEnemy
	}

	if layerA == LayerCollectible {
		return layerB == LayerPlayer
	}

	if layerA == LayerTerrain {
		return layerB == LayerPlayer ||
			layerB == LayerEnemy ||
			layerB == LayerProjectile
	}

	if layerA == LayerOneWayPlatform {
		return layerB == LayerPlayer
	}

	return false
}

// RaycastHit represents the result of a raycast
type RaycastHit struct {
	Collider Collider
	Point    Vector2
	Normal   Vector2
	Distance float64
	Hit      bool
}

// Raycast performs a raycast from origin in direction for maxDistance
func Raycast(origin, direction Vector2, maxDistance float64, colliders []Collider, layerMask CollisionLayer) *RaycastHit {
	direction = direction.Normalize()
	hit := &RaycastHit{Hit: false}
	closestDistance := maxDistance

	for _, collider := range colliders {
		// Check if collider is in the layer mask
		if collider.GetLayer()&layerMask == 0 {
			continue
		}

		bounds := collider.GetBounds()
		
		// Ray-AABB intersection test
		tMin := (bounds.Left() - origin.X) / direction.X
		tMax := (bounds.Right() - origin.X) / direction.X

		if tMin > tMax {
			tMin, tMax = tMax, tMin
		}

		tyMin := (bounds.Top() - origin.Y) / direction.Y
		tyMax := (bounds.Bottom() - origin.Y) / direction.Y

		if tyMin > tyMax {
			tyMin, tyMax = tyMax, tyMin
		}

		if tMin > tyMax || tyMin > tMax {
			continue // No intersection
		}

		tMin = math.Max(tMin, tyMin)
		tMax = math.Min(tMax, tyMax)

		if tMin < 0 {
			tMin = tMax
		}

		if tMin >= 0 && tMin <= maxDistance && tMin < closestDistance {
			closestDistance = tMin
			hit.Hit = true
			hit.Collider = collider
			hit.Distance = tMin
			hit.Point = origin.Add(direction.Multiply(tMin))

			// Calculate normal based on which face was hit
			center := bounds.Center()
			diff := hit.Point.Subtract(center)
			absX := math.Abs(diff.X)
			absY := math.Abs(diff.Y)

			if absX > absY {
				if diff.X > 0 {
					hit.Normal = Vector2{1, 0}
				} else {
					hit.Normal = Vector2{-1, 0}
				}
			} else {
				if diff.Y > 0 {
					hit.Normal = Vector2{0, 1}
				} else {
					hit.Normal = Vector2{0, -1}
				}
			}
		}
	}

	return hit
}
