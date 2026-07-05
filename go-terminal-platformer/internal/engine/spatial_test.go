package engine

import (
	"testing"
)

func TestNewSpatialGrid(t *testing.T) {
	grid := NewSpatialGrid(64)

	if grid.GetCellSize() != 64 {
		t.Errorf("Expected cell size 64, got %d", grid.GetCellSize())
	}

	if grid.GetCellCount() != 0 {
		t.Error("New grid should have 0 cells")
	}
}

func TestSpatialGridInsert(t *testing.T) {
	grid := NewSpatialGrid(64)
	collider := NewAABBCollider(10, 10, 20, 20, LayerPlayer)

	grid.Insert(collider)

	if grid.GetCellCount() == 0 {
		t.Error("Grid should have at least one cell after insert")
	}
}

func TestSpatialGridQuery(t *testing.T) {
	grid := NewSpatialGrid(64)
	c1 := NewAABBCollider(10, 10, 20, 20, LayerPlayer)
	c2 := NewAABBCollider(100, 100, 20, 20, LayerEnemy)

	grid.Insert(c1)
	grid.Insert(c2)

	// Query area containing c1
	bounds := Rectangle{X: 0, Y: 0, Width: 50, Height: 50}
	results := grid.Query(bounds)

	found := false
	for _, collider := range results {
		if collider == c1 {
			found = true
			break
		}
	}

	if !found {
		t.Error("Query should find c1 in the specified bounds")
	}

	// c2 should not be in results
	for _, collider := range results {
		if collider == c2 {
			t.Error("Query should not find c2 outside the specified bounds")
		}
	}
}

func TestSpatialGridQueryPoint(t *testing.T) {
	grid := NewSpatialGrid(64)
	collider := NewAABBCollider(10, 10, 20, 20, LayerPlayer)

	grid.Insert(collider)

	// Query point inside collider's cell
	point := Vector2{X: 15, Y: 15}
	results := grid.QueryPoint(point)

	found := false
	for _, c := range results {
		if c == collider {
			found = true
			break
		}
	}

	if !found {
		t.Error("QueryPoint should find collider in the same cell")
	}

	// Query point in different cell
	point = Vector2{X: 200, Y: 200}
	results = grid.QueryPoint(point)

	for _, c := range results {
		if c == collider {
			t.Error("QueryPoint should not find collider in different cell")
		}
	}
}

func TestSpatialGridQueryRadius(t *testing.T) {
	grid := NewSpatialGrid(64)
	c1 := NewAABBCollider(50, 50, 20, 20, LayerPlayer)
	c2 := NewAABBCollider(200, 200, 20, 20, LayerEnemy)

	grid.Insert(c1)
	grid.Insert(c2)

	// Query with radius that includes c1 but not c2
	center := Vector2{X: 60, Y: 60}
	radius := 50.0
	results := grid.QueryRadius(center, radius)

	foundC1 := false
	foundC2 := false

	for _, collider := range results {
		if collider == c1 {
			foundC1 = true
		}
		if collider == c2 {
			foundC2 = true
		}
	}

	if !foundC1 {
		t.Error("QueryRadius should find c1 within radius")
	}

	if foundC2 {
		t.Error("QueryRadius should not find c2 outside radius")
	}
}

func TestSpatialGridRemove(t *testing.T) {
	grid := NewSpatialGrid(64)
	collider := NewAABBCollider(10, 10, 20, 20, LayerPlayer)

	grid.Insert(collider)
	grid.Remove(collider)

	// Query should not find the removed collider
	bounds := Rectangle{X: 0, Y: 0, Width: 50, Height: 50}
	results := grid.Query(bounds)

	for _, c := range results {
		if c == collider {
			t.Error("Query should not find removed collider")
		}
	}
}

func TestSpatialGridUpdate(t *testing.T) {
	grid := NewSpatialGrid(64)
	collider := NewAABBCollider(10, 10, 20, 20, LayerPlayer)

	grid.Insert(collider)

	// Move collider to new position
	collider.SetPosition(100, 100)
	grid.Update(collider)

	// Query old position - should not find it
	oldBounds := Rectangle{X: 0, Y: 0, Width: 50, Height: 50}
	results := grid.Query(oldBounds)

	for _, c := range results {
		if c == collider {
			t.Error("Query should not find collider at old position after update")
		}
	}

	// Query new position - should find it
	newBounds := Rectangle{X: 90, Y: 90, Width: 50, Height: 50}
	results = grid.Query(newBounds)

	found := false
	for _, c := range results {
		if c == collider {
			found = true
			break
		}
	}

	if !found {
		t.Error("Query should find collider at new position after update")
	}
}

func TestSpatialGridClear(t *testing.T) {
	grid := NewSpatialGrid(64)
	c1 := NewAABBCollider(10, 10, 20, 20, LayerPlayer)
	c2 := NewAABBCollider(100, 100, 20, 20, LayerEnemy)

	grid.Insert(c1)
	grid.Insert(c2)

	grid.Clear()

	if grid.GetCellCount() != 0 {
		t.Error("Grid should have 0 cells after clear")
	}

	// Query should return empty results
	bounds := Rectangle{X: 0, Y: 0, Width: 200, Height: 200}
	results := grid.Query(bounds)

	if len(results) != 0 {
		t.Error("Query should return empty results after clear")
	}
}

func TestSpatialGridMultipleCells(t *testing.T) {
	grid := NewSpatialGrid(64)
	
	// Create a large collider that spans multiple cells
	largeCollider := NewAABBCollider(10, 10, 100, 100, LayerPlayer)
	grid.Insert(largeCollider)

	// Query different areas that should all find the large collider
	testBounds := []Rectangle{
		{X: 0, Y: 0, Width: 50, Height: 50},       // Top-left
		{X: 80, Y: 0, Width: 50, Height: 50},      // Top-right
		{X: 0, Y: 80, Width: 50, Height: 50},      // Bottom-left
		{X: 80, Y: 80, Width: 50, Height: 50},     // Bottom-right
	}

	for i, bounds := range testBounds {
		results := grid.Query(bounds)
		found := false
		for _, c := range results {
			if c == largeCollider {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Query %d should find large collider spanning multiple cells", i)
		}
	}
}

func TestSpatialGridNoDuplicates(t *testing.T) {
	grid := NewSpatialGrid(64)
	
	// Create a collider that spans multiple cells
	collider := NewAABBCollider(10, 10, 100, 100, LayerPlayer)
	grid.Insert(collider)

	// Query an area that overlaps multiple cells containing the same collider
	bounds := Rectangle{X: 0, Y: 0, Width: 150, Height: 150}
	results := grid.Query(bounds)

	// Count occurrences of the collider
	count := 0
	for _, c := range results {
		if c == collider {
			count++
		}
	}

	if count != 1 {
		t.Errorf("Query should return each collider only once, got %d occurrences", count)
	}
}

func TestSpatialGridGetTotalColliders(t *testing.T) {
	grid := NewSpatialGrid(64)
	c1 := NewAABBCollider(10, 10, 20, 20, LayerPlayer)
	c2 := NewAABBCollider(100, 100, 20, 20, LayerEnemy)

	grid.Insert(c1)
	grid.Insert(c2)

	total := grid.GetTotalColliders()
	if total < 2 {
		t.Errorf("Expected at least 2 collider references, got %d", total)
	}
}

func TestSpatialGridDebugGetCellsForBounds(t *testing.T) {
	grid := NewSpatialGrid(64)
	bounds := Rectangle{X: 0, Y: 0, Width: 100, Height: 100}

	cells := grid.DebugGetCellsForBounds(bounds)

	if len(cells) == 0 {
		t.Error("DebugGetCellsForBounds should return at least one cell")
	}

	// Verify cells are within expected range
	for _, cell := range cells {
		if cell.X < 0 || cell.Y < 0 {
			t.Errorf("Cell coordinates should be non-negative, got (%d, %d)", cell.X, cell.Y)
		}
	}
}

func TestSpatialGridConcurrency(t *testing.T) {
	grid := NewSpatialGrid(64)
	collider := NewAABBCollider(10, 10, 20, 20, LayerPlayer)

	// Test that concurrent reads don't panic
	done := make(chan bool)
	
	// Writer goroutine
	go func() {
		for i := 0; i < 100; i++ {
			grid.Insert(collider)
			grid.Remove(collider)
		}
		done <- true
	}()

	// Reader goroutine
	go func() {
		for i := 0; i < 100; i++ {
			bounds := Rectangle{X: 0, Y: 0, Width: 50, Height: 50}
			grid.Query(bounds)
		}
		done <- true
	}()

	// Wait for both goroutines
	<-done
	<-done
}

func TestSpatialGridEdgeCases(t *testing.T) {
	grid := NewSpatialGrid(64)

	// Test with zero-size collider
	zeroCollider := NewAABBCollider(10, 10, 0, 0, LayerPlayer)
	grid.Insert(zeroCollider)

	bounds := Rectangle{X: 0, Y: 0, Width: 50, Height: 50}
	results := grid.Query(bounds)

	found := false
	for _, c := range results {
		if c == zeroCollider {
			found = true
			break
		}
	}

	if !found {
		t.Error("Should be able to insert and query zero-size collider")
	}

	// Test with negative coordinates
	negCollider := NewAABBCollider(-50, -50, 20, 20, LayerPlayer)
	grid.Insert(negCollider)

	bounds = Rectangle{X: -100, Y: -100, Width: 100, Height: 100}
	results = grid.Query(bounds)

	found = false
	for _, c := range results {
		if c == negCollider {
			found = true
			break
		}
	}

	if !found {
		t.Error("Should be able to insert and query collider with negative coordinates")
	}
}

func TestSpatialGridPerformance(t *testing.T) {
	grid := NewSpatialGrid(64)

	// Insert many colliders
	colliders := make([]Collider, 1000)
	for i := 0; i < 1000; i++ {
		x := float64(i % 100 * 10)
		y := float64(i / 100 * 10)
		colliders[i] = NewAABBCollider(x, y, 10, 10, LayerPlayer)
		grid.Insert(colliders[i])
	}

	// Query should be fast even with many colliders
	bounds := Rectangle{X: 0, Y: 0, Width: 100, Height: 100}
	results := grid.Query(bounds)

	// Should find some colliders but not all (spatial partitioning working)
	if len(results) == 0 {
		t.Error("Query should find some colliders")
	}

	if len(results) >= 1000 {
		t.Error("Query should not return all colliders (spatial partitioning not working)")
	}
}
