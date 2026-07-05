package levels

import (
	"github.com/lukeroy/go-terminal-platformer/internal/engine"
	"github.com/lukeroy/go-terminal-platformer/internal/entities"
)

// BlockType represents the type of block
type BlockType int

const (
	BlockSolid BlockType = iota
	BlockQuestion
	BlockBrick
	BlockInvisible
)

// String returns the string representation of a block type
func (bt BlockType) String() string {
	switch bt {
	case BlockSolid:
		return "Solid"
	case BlockQuestion:
		return "Question"
	case BlockBrick:
		return "Brick"
	case BlockInvisible:
		return "Invisible"
	default:
		return "Unknown"
	}
}

// Block represents a level block that can contain items
type Block struct {
	position     engine.Vector2
	blockType    BlockType
	contents     entities.PowerUpType
	coinCount    int
	hit          bool
	hitAnimation float64
	hitOffset    float64
	destroyed    bool

	// Collision
	collider *engine.AABBCollider

	// Rendering
	sprite rune
	color  int
}

// NewBlock creates a new block
func NewBlock(pos engine.Vector2, blockType BlockType) *Block {
	block := &Block{
		position:  pos,
		blockType: blockType,
		contents:  entities.PowerUpNone,
		coinCount: 0,
		hit:       false,
		destroyed: false,
	}

	// Create collider
	block.collider = engine.NewAABBCollider(pos.X, pos.Y, 16, 16, engine.LayerTerrain)

	// Set sprite based on type
	block.setSprite()

	return block
}

// NewBlockWithPowerUp creates a block containing a power-up
func NewBlockWithPowerUp(pos engine.Vector2, blockType BlockType, powerUp entities.PowerUpType) *Block {
	block := NewBlock(pos, blockType)
	block.contents = powerUp
	return block
}

// NewBlockWithCoins creates a block containing coins
func NewBlockWithCoins(pos engine.Vector2, blockType BlockType, coinCount int) *Block {
	block := NewBlock(pos, blockType)
	block.coinCount = coinCount
	return block
}

// setSprite sets the sprite based on block type
func (b *Block) setSprite() {
	switch b.blockType {
	case BlockSolid:
		b.sprite = '█'
		b.color = 7 // Gray
	case BlockQuestion:
		if b.hit {
			b.sprite = '□'
			b.color = 7 // Gray (used block)
		} else {
			b.sprite = '?'
			b.color = 4 // Yellow
		}
	case BlockBrick:
		b.sprite = '▒'
		b.color = 3 // Brown/Red
	case BlockInvisible:
		b.sprite = ' '
		b.color = 0 // Transparent
	}
}

// Update updates the block state
func (b *Block) Update(deltaTime float64) {
	if b.destroyed {
		return
	}

	// Update hit animation
	if b.hitAnimation > 0 {
		b.hitAnimation -= deltaTime
		
		// Animate block moving up and down
		progress := 1.0 - (b.hitAnimation / 0.3)
		if progress < 0.5 {
			// Moving up
			b.hitOffset = -10 * (progress * 2)
		} else {
			// Moving down
			b.hitOffset = -10 * (2 - progress*2)
		}

		if b.hitAnimation <= 0 {
			b.hitAnimation = 0
			b.hitOffset = 0
		}
	}
}

// Hit handles the block being hit by a player
func (b *Block) Hit(player entities.Player, powerUpMgr *entities.PowerUpManager) bool {
	if b.destroyed {
		return false
	}

	// Solid blocks can't be hit
	if b.blockType == BlockSolid {
		return false
	}

	// Question blocks can only be hit once
	if b.blockType == BlockQuestion && b.hit {
		return false
	}

	// Brick blocks can be destroyed if player has power-up
	if b.blockType == BlockBrick {
		if player.GetPowerUp() != entities.PowerUpNone {
			b.Destroy(powerUpMgr)
			return true
		}
		// Otherwise just bump it
		b.hitAnimation = 0.3
		return false
	}

	// Mark as hit
	b.hit = true
	b.hitAnimation = 0.3
	b.setSprite()

	// Spawn contents
	spawnPos := engine.Vector2{X: b.position.X, Y: b.position.Y - 20}
	spawnVelocity := engine.Vector2{X: 0, Y: -200}

	if b.contents != entities.PowerUpNone {
		// Spawn power-up
		powerUpMgr.SpawnPowerUpWithVelocity(b.contents, spawnPos, spawnVelocity)
		b.contents = entities.PowerUpNone
	} else if b.coinCount > 0 {
		// Spawn coin
		powerUpMgr.SpawnCoinWithVelocity(spawnPos, spawnVelocity)
		b.coinCount--
		
		// If more coins remain, block can be hit again
		if b.coinCount > 0 {
			b.hit = false
		}
	} else {
		// Empty block, spawn single coin
		powerUpMgr.SpawnCoinWithVelocity(spawnPos, spawnVelocity)
	}

	return true
}

// Destroy destroys the block
func (b *Block) Destroy(powerUpMgr *entities.PowerUpManager) {
	b.destroyed = true
	
	// Spawn debris/particles (TODO: implement particle system)
	// For now, spawn a coin as reward
	spawnPos := engine.Vector2{X: b.position.X, Y: b.position.Y}
	powerUpMgr.SpawnCoin(spawnPos)
}

// CheckHitFromBelow checks if a player hit the block from below
func (b *Block) CheckHitFromBelow(player entities.Player) bool {
	if b.destroyed {
		return false
	}

	playerPos := player.GetPosition()
	playerBody := player.GetPhysicsBody()
	
	if playerBody == nil {
		return false
	}

	// Check if player is below the block
	if playerPos.Y > b.position.Y {
		// Check if player's top collides with block's bottom
		playerBounds := playerBody.Collider.GetBounds()
		blockBounds := b.collider.GetBounds()
		
		playerTop := playerBounds.Y
		blockBottom := blockBounds.Y + blockBounds.Height
		
		// Check vertical overlap
		if playerTop <= blockBottom && playerTop >= blockBottom-8 {
			// Check horizontal overlap
			if playerBounds.X < blockBounds.X+blockBounds.Width &&
				playerBounds.X+playerBounds.Width > blockBounds.X {
				// Check if player is moving upward
				if playerBody.Velocity.Y < 0 {
					return true
				}
			}
		}
	}

	return false
}

// Render renders the block
func (b *Block) Render(renderer entities.Renderer) {
	if b.destroyed {
		return
	}

	// Apply hit animation offset
	renderY := int(b.position.Y + b.hitOffset)
	
	sprite := entities.NewSprite(1, 1)
	sprite.Data[0][0] = b.sprite
	sprite.Color = b.color
	
	renderer.DrawSprite(int(b.position.X), renderY, sprite)
}

// Getters

// GetPosition returns the block's position
func (b *Block) GetPosition() engine.Vector2 {
	return b.position
}

// GetType returns the block type
func (b *Block) GetType() BlockType {
	return b.blockType
}

// IsHit returns whether the block has been hit
func (b *Block) IsHit() bool {
	return b.hit
}

// IsDestroyed returns whether the block is destroyed
func (b *Block) IsDestroyed() bool {
	return b.destroyed
}

// GetCollider returns the block's collider
func (b *Block) GetCollider() *engine.AABBCollider {
	return b.collider
}

// GetContents returns the power-up contents
func (b *Block) GetContents() entities.PowerUpType {
	return b.contents
}

// GetCoinCount returns the remaining coin count
func (b *Block) GetCoinCount() int {
	return b.coinCount
}

// Setters

// SetContents sets the power-up contents
func (b *Block) SetContents(powerUp entities.PowerUpType) {
	b.contents = powerUp
}

// SetCoinCount sets the coin count
func (b *Block) SetCoinCount(count int) {
	b.coinCount = count
}

// Reset resets the block to its initial state
func (b *Block) Reset() {
	b.hit = false
	b.hitAnimation = 0
	b.hitOffset = 0
	b.destroyed = false
	b.setSprite()
}
