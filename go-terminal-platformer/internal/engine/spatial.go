package engine

import (
	"fmt"
	"sync"
)

// SpatialGrid provides efficient broad-phase collision detection using a grid-based spatial partitioning system
type SpatialGrid struct {
	cellSize int
	cells    map[string][]Collider
	mutex    sync.RWMutex
}

// NewSpatialGrid creates a new spatial grid with the given cell size
func NewSpatialGrid(cellSize int) *SpatialGrid {
	return &SpatialGrid{
		cellSize: cellSize,
		cells:    make(map[string][]Collider),
	}
}

// Clear removes all colliders from the grid
func (sg *SpatialGrid) Clear() {
	sg.mutex.Lock()
	defer sg.mutex.Unlock()
	sg.cells = make(map[string][]Collider)
}

// Insert adds a collider to the grid
func (sg *SpatialGrid) Insert(collider Collider) {
	sg.mutex.Lock()
	defer sg.mutex.Unlock()

	bounds := collider.GetBounds()
	cells := sg.getCellsForBounds(bounds)

	for _, cellKey := range cells {
		sg.cells[cellKey] = append(sg.cells[cellKey], collider)
	}
}

// Remove removes a collider from the grid
func (sg *SpatialGrid) Remove(collider Collider) {
	sg.mutex.Lock()
	defer sg.mutex.Unlock()

	bounds := collider.GetBounds()
	cells := sg.getCellsForBounds(bounds)

	for _, cellKey := range cells {
		if colliders, exists := sg.cells[cellKey]; exists {
			// Remove the collider from this cell
			for i, c := range colliders {
				if c == collider {
					sg.cells[cellKey] = append(colliders[:i], colliders[i+1:]...)
					break
				}
			}
			// Clean up empty cells
			if len(sg.cells[cellKey]) == 0 {
				delete(sg.cells, cellKey)
			}
		}
	}
}

// Update updates a collider's position in the grid
func (sg *SpatialGrid) Update(collider Collider) {
	// Remove and re-insert is simpler and works well for dynamic objects
	sg.Remove(collider)
	sg.Insert(collider)
}

// Query returns all colliders that could potentially intersect with the given bounds
func (sg *SpatialGrid) Query(bounds Rectangle) []Collider {
	sg.mutex.RLock()
	defer sg.mutex.RUnlock()

	cells := sg.getCellsForBounds(bounds)
	colliders := make([]Collider, 0)
	seen := make(map[Collider]bool)

	for _, cellKey := range cells {
		if cellColliders, exists := sg.cells[cellKey]; exists {
			for _, collider := range cellColliders {
				if !seen[collider] {
					colliders = append(colliders, collider)
					seen[collider] = true
				}
			}
		}
	}

	return colliders
}

// QueryPoint returns all colliders in the cell containing the given point
func (sg *SpatialGrid) QueryPoint(point Vector2) []Collider {
	sg.mutex.RLock()
	defer sg.mutex.RUnlock()

	cellKey := sg.getCellKey(int(point.X)/sg.cellSize, int(point.Y)/sg.cellSize)
	
	if colliders, exists := sg.cells[cellKey]; exists {
		result := make([]Collider, len(colliders))
		copy(result, colliders)
		return result
	}

	return []Collider{}
}

// QueryRadius returns all colliders within a radius of the given point
func (sg *SpatialGrid) QueryRadius(center Vector2, radius float64) []Collider {
	bounds := Rectangle{
		X:      center.X - radius,
		Y:      center.Y - radius,
		Width:  radius * 2,
		Height: radius * 2,
	}
	
	candidates := sg.Query(bounds)
	result := make([]Collider, 0)
	radiusSquared := radius * radius

	for _, collider := range candidates {
		colliderBounds := collider.GetBounds()
		colliderCenter := colliderBounds.Center()
		
		if center.DistanceSquared(colliderCenter) <= radiusSquared {
			result = append(result, collider)
		}
	}

	return result
}

// getCellsForBounds returns all cell keys that the bounds overlap
func (sg *SpatialGrid) getCellsForBounds(bounds Rectangle) []string {
	minX := int(bounds.X) / sg.cellSize
	minY := int(bounds.Y) / sg.cellSize
	maxX := int(bounds.X+bounds.Width) / sg.cellSize
	maxY := int(bounds.Y+bounds.Height) / sg.cellSize

	cells := make([]string, 0, (maxX-minX+1)*(maxY-minY+1))

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			cells = append(cells, sg.getCellKey(x, y))
		}
	}

	return cells
}

// getCellKey generates a unique key for a cell at the given coordinates
func (sg *SpatialGrid) getCellKey(x, y int) string {
	return fmt.Sprintf("%d,%d", x, y)
}

// GetCellSize returns the cell size of the grid
func (sg *SpatialGrid) GetCellSize() int {
	return sg.cellSize
}

// GetCellCount returns the number of active cells in the grid
func (sg *SpatialGrid) GetCellCount() int {
	sg.mutex.RLock()
	defer sg.mutex.RUnlock()
	return len(sg.cells)
}

// GetTotalColliders returns the total number of collider references in the grid
// Note: This counts duplicates (same collider in multiple cells)
func (sg *SpatialGrid) GetTotalColliders() int {
	sg.mutex.RLock()
	defer sg.mutex.RUnlock()
	
	total := 0
	for _, colliders := range sg.cells {
		total += len(colliders)
	}
	return total
}

// DebugGetCellsForBounds is a debug helper that returns cell coordinates for bounds
func (sg *SpatialGrid) DebugGetCellsForBounds(bounds Rectangle) []struct{ X, Y int } {
	minX := int(bounds.X) / sg.cellSize
	minY := int(bounds.Y) / sg.cellSize
	maxX := int(bounds.X+bounds.Width) / sg.cellSize
	maxY := int(bounds.Y+bounds.Height) / sg.cellSize

	cells := make([]struct{ X, Y int }, 0)

	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			cells = append(cells, struct{ X, Y int }{X: x, Y: y})
		}
	}

	return cells
}
