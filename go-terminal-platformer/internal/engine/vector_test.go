package engine

import (
	"math"
	"testing"
)

func TestVector2Add(t *testing.T) {
	v1 := Vector2{X: 1, Y: 2}
	v2 := Vector2{X: 3, Y: 4}
	result := v1.Add(v2)

	if result.X != 4 || result.Y != 6 {
		t.Errorf("Add failed: expected (4, 6), got (%f, %f)", result.X, result.Y)
	}
}

func TestVector2Subtract(t *testing.T) {
	v1 := Vector2{X: 5, Y: 7}
	v2 := Vector2{X: 2, Y: 3}
	result := v1.Subtract(v2)

	if result.X != 3 || result.Y != 4 {
		t.Errorf("Subtract failed: expected (3, 4), got (%f, %f)", result.X, result.Y)
	}
}

func TestVector2Multiply(t *testing.T) {
	v := Vector2{X: 2, Y: 3}
	result := v.Multiply(3)

	if result.X != 6 || result.Y != 9 {
		t.Errorf("Multiply failed: expected (6, 9), got (%f, %f)", result.X, result.Y)
	}
}

func TestVector2Divide(t *testing.T) {
	v := Vector2{X: 6, Y: 9}
	result := v.Divide(3)

	if result.X != 2 || result.Y != 3 {
		t.Errorf("Divide failed: expected (2, 3), got (%f, %f)", result.X, result.Y)
	}
}

func TestVector2DivideByZero(t *testing.T) {
	v := Vector2{X: 6, Y: 9}
	result := v.Divide(0)

	if result.X != 6 || result.Y != 9 {
		t.Errorf("Divide by zero should return original vector, got (%f, %f)", result.X, result.Y)
	}
}

func TestVector2Dot(t *testing.T) {
	v1 := Vector2{X: 1, Y: 2}
	v2 := Vector2{X: 3, Y: 4}
	result := v1.Dot(v2)

	expected := 1*3 + 2*4 // 11
	if result != float64(expected) {
		t.Errorf("Dot product failed: expected %f, got %f", float64(expected), result)
	}
}

func TestVector2Cross(t *testing.T) {
	v1 := Vector2{X: 1, Y: 2}
	v2 := Vector2{X: 3, Y: 4}
	result := v1.Cross(v2)

	expected := 1*4 - 2*3 // -2
	if result != float64(expected) {
		t.Errorf("Cross product failed: expected %f, got %f", float64(expected), result)
	}
}

func TestVector2Length(t *testing.T) {
	v := Vector2{X: 3, Y: 4}
	result := v.Length()

	expected := 5.0 // sqrt(3^2 + 4^2) = 5
	if math.Abs(result-expected) > 0.0001 {
		t.Errorf("Length failed: expected %f, got %f", expected, result)
	}
}

func TestVector2LengthSquared(t *testing.T) {
	v := Vector2{X: 3, Y: 4}
	result := v.LengthSquared()

	expected := 25.0 // 3^2 + 4^2 = 25
	if result != expected {
		t.Errorf("LengthSquared failed: expected %f, got %f", expected, result)
	}
}

func TestVector2Normalize(t *testing.T) {
	v := Vector2{X: 3, Y: 4}
	result := v.Normalize()

	expectedLength := 1.0
	if math.Abs(result.Length()-expectedLength) > 0.0001 {
		t.Errorf("Normalize failed: expected length %f, got %f", expectedLength, result.Length())
	}

	// Check direction is preserved
	expectedX := 3.0 / 5.0
	expectedY := 4.0 / 5.0
	if math.Abs(result.X-expectedX) > 0.0001 || math.Abs(result.Y-expectedY) > 0.0001 {
		t.Errorf("Normalize failed: expected (%f, %f), got (%f, %f)", expectedX, expectedY, result.X, result.Y)
	}
}

func TestVector2NormalizeZero(t *testing.T) {
	v := Vector2{X: 0, Y: 0}
	result := v.Normalize()

	if result.X != 0 || result.Y != 0 {
		t.Errorf("Normalize zero vector should return zero, got (%f, %f)", result.X, result.Y)
	}
}

func TestVector2Distance(t *testing.T) {
	v1 := Vector2{X: 0, Y: 0}
	v2 := Vector2{X: 3, Y: 4}
	result := v1.Distance(v2)

	expected := 5.0
	if math.Abs(result-expected) > 0.0001 {
		t.Errorf("Distance failed: expected %f, got %f", expected, result)
	}
}

func TestVector2DistanceSquared(t *testing.T) {
	v1 := Vector2{X: 0, Y: 0}
	v2 := Vector2{X: 3, Y: 4}
	result := v1.DistanceSquared(v2)

	expected := 25.0
	if result != expected {
		t.Errorf("DistanceSquared failed: expected %f, got %f", expected, result)
	}
}

func TestVector2Lerp(t *testing.T) {
	v1 := Vector2{X: 0, Y: 0}
	v2 := Vector2{X: 10, Y: 10}
	
	// Test at t=0.5 (midpoint)
	result := v1.Lerp(v2, 0.5)
	if result.X != 5 || result.Y != 5 {
		t.Errorf("Lerp at 0.5 failed: expected (5, 5), got (%f, %f)", result.X, result.Y)
	}

	// Test at t=0 (start)
	result = v1.Lerp(v2, 0)
	if result.X != 0 || result.Y != 0 {
		t.Errorf("Lerp at 0 failed: expected (0, 0), got (%f, %f)", result.X, result.Y)
	}

	// Test at t=1 (end)
	result = v1.Lerp(v2, 1)
	if result.X != 10 || result.Y != 10 {
		t.Errorf("Lerp at 1 failed: expected (10, 10), got (%f, %f)", result.X, result.Y)
	}
}

func TestVector2Clamp(t *testing.T) {
	v := Vector2{X: 15, Y: -5}
	min := Vector2{X: 0, Y: 0}
	max := Vector2{X: 10, Y: 10}
	result := v.Clamp(min, max)

	if result.X != 10 || result.Y != 0 {
		t.Errorf("Clamp failed: expected (10, 0), got (%f, %f)", result.X, result.Y)
	}
}

func TestVector2Rotate(t *testing.T) {
	v := Vector2{X: 1, Y: 0}
	result := v.Rotate(math.Pi / 2) // 90 degrees

	// Should be approximately (0, 1)
	if math.Abs(result.X) > 0.0001 || math.Abs(result.Y-1) > 0.0001 {
		t.Errorf("Rotate 90 degrees failed: expected (0, 1), got (%f, %f)", result.X, result.Y)
	}
}

func TestVector2Perpendicular(t *testing.T) {
	v := Vector2{X: 1, Y: 2}
	result := v.Perpendicular()

	// Perpendicular should be (-2, 1)
	if result.X != -2 || result.Y != 1 {
		t.Errorf("Perpendicular failed: expected (-2, 1), got (%f, %f)", result.X, result.Y)
	}

	// Dot product should be zero
	dot := v.Dot(result)
	if math.Abs(dot) > 0.0001 {
		t.Errorf("Perpendicular vectors should have dot product of 0, got %f", dot)
	}
}

func TestVector2Constants(t *testing.T) {
	zero := Zero()
	if zero.X != 0 || zero.Y != 0 {
		t.Errorf("Zero() failed: expected (0, 0), got (%f, %f)", zero.X, zero.Y)
	}

	one := One()
	if one.X != 1 || one.Y != 1 {
		t.Errorf("One() failed: expected (1, 1), got (%f, %f)", one.X, one.Y)
	}

	up := Up()
	if up.X != 0 || up.Y != -1 {
		t.Errorf("Up() failed: expected (0, -1), got (%f, %f)", up.X, up.Y)
	}

	down := Down()
	if down.X != 0 || down.Y != 1 {
		t.Errorf("Down() failed: expected (0, 1), got (%f, %f)", down.X, down.Y)
	}

	left := Left()
	if left.X != -1 || left.Y != 0 {
		t.Errorf("Left() failed: expected (-1, 0), got (%f, %f)", left.X, left.Y)
	}

	right := Right()
	if right.X != 1 || right.Y != 0 {
		t.Errorf("Right() failed: expected (1, 0), got (%f, %f)", right.X, right.Y)
	}
}

func TestRectangleEdges(t *testing.T) {
	r := Rectangle{X: 10, Y: 20, Width: 30, Height: 40}

	if r.Left() != 10 {
		t.Errorf("Left() failed: expected 10, got %f", r.Left())
	}
	if r.Right() != 40 {
		t.Errorf("Right() failed: expected 40, got %f", r.Right())
	}
	if r.Top() != 20 {
		t.Errorf("Top() failed: expected 20, got %f", r.Top())
	}
	if r.Bottom() != 60 {
		t.Errorf("Bottom() failed: expected 60, got %f", r.Bottom())
	}
}

func TestRectangleCenter(t *testing.T) {
	r := Rectangle{X: 10, Y: 20, Width: 30, Height: 40}
	center := r.Center()

	expectedX := 25.0 // 10 + 30/2
	expectedY := 40.0 // 20 + 40/2

	if center.X != expectedX || center.Y != expectedY {
		t.Errorf("Center() failed: expected (%f, %f), got (%f, %f)", expectedX, expectedY, center.X, center.Y)
	}
}

func TestRectangleContains(t *testing.T) {
	r := Rectangle{X: 10, Y: 10, Width: 20, Height: 20}

	// Point inside
	if !r.Contains(Vector2{X: 15, Y: 15}) {
		t.Error("Contains() failed: point (15, 15) should be inside")
	}

	// Point outside
	if r.Contains(Vector2{X: 5, Y: 5}) {
		t.Error("Contains() failed: point (5, 5) should be outside")
	}

	// Point on edge
	if !r.Contains(Vector2{X: 10, Y: 10}) {
		t.Error("Contains() failed: point (10, 10) should be on edge (inside)")
	}
}

func TestRectangleIntersects(t *testing.T) {
	r1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	r2 := Rectangle{X: 5, Y: 5, Width: 10, Height: 10}
	r3 := Rectangle{X: 20, Y: 20, Width: 10, Height: 10}

	// Overlapping rectangles
	if !r1.Intersects(r2) {
		t.Error("Intersects() failed: r1 and r2 should intersect")
	}

	// Non-overlapping rectangles
	if r1.Intersects(r3) {
		t.Error("Intersects() failed: r1 and r3 should not intersect")
	}
}

func TestRectangleIntersection(t *testing.T) {
	r1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	r2 := Rectangle{X: 5, Y: 5, Width: 10, Height: 10}

	intersection := r1.Intersection(r2)

	expectedX := 5.0
	expectedY := 5.0
	expectedWidth := 5.0
	expectedHeight := 5.0

	if intersection.X != expectedX || intersection.Y != expectedY ||
		intersection.Width != expectedWidth || intersection.Height != expectedHeight {
		t.Errorf("Intersection() failed: expected (%f, %f, %f, %f), got (%f, %f, %f, %f)",
			expectedX, expectedY, expectedWidth, expectedHeight,
			intersection.X, intersection.Y, intersection.Width, intersection.Height)
	}
}

func TestRectangleUnion(t *testing.T) {
	r1 := Rectangle{X: 0, Y: 0, Width: 10, Height: 10}
	r2 := Rectangle{X: 5, Y: 5, Width: 10, Height: 10}

	union := r1.Union(r2)

	expectedX := 0.0
	expectedY := 0.0
	expectedWidth := 15.0
	expectedHeight := 15.0

	if union.X != expectedX || union.Y != expectedY ||
		union.Width != expectedWidth || union.Height != expectedHeight {
		t.Errorf("Union() failed: expected (%f, %f, %f, %f), got (%f, %f, %f, %f)",
			expectedX, expectedY, expectedWidth, expectedHeight,
			union.X, union.Y, union.Width, union.Height)
	}
}

func TestRectangleExpand(t *testing.T) {
	r := Rectangle{X: 10, Y: 10, Width: 10, Height: 10}
	expanded := r.Expand(5)

	expectedX := 5.0
	expectedY := 5.0
	expectedWidth := 20.0
	expectedHeight := 20.0

	if expanded.X != expectedX || expanded.Y != expectedY ||
		expanded.Width != expectedWidth || expanded.Height != expectedHeight {
		t.Errorf("Expand() failed: expected (%f, %f, %f, %f), got (%f, %f, %f, %f)",
			expectedX, expectedY, expectedWidth, expectedHeight,
			expanded.X, expanded.Y, expanded.Width, expanded.Height)
	}
}

func TestRectangleTranslate(t *testing.T) {
	r := Rectangle{X: 10, Y: 10, Width: 10, Height: 10}
	translated := r.Translate(Vector2{X: 5, Y: -3})

	expectedX := 15.0
	expectedY := 7.0

	if translated.X != expectedX || translated.Y != expectedY {
		t.Errorf("Translate() failed: expected position (%f, %f), got (%f, %f)",
			expectedX, expectedY, translated.X, translated.Y)
	}

	// Size should remain the same
	if translated.Width != r.Width || translated.Height != r.Height {
		t.Error("Translate() should not change size")
	}
}

func TestHelperFunctions(t *testing.T) {
	// Test Min
	if Min(5, 10) != 5 {
		t.Error("Min(5, 10) should return 5")
	}
	if Min(10, 5) != 5 {
		t.Error("Min(10, 5) should return 5")
	}

	// Test Max
	if Max(5, 10) != 10 {
		t.Error("Max(5, 10) should return 10")
	}
	if Max(10, 5) != 10 {
		t.Error("Max(10, 5) should return 10")
	}

	// Test Abs
	if Abs(-5) != 5 {
		t.Error("Abs(-5) should return 5")
	}
	if Abs(5) != 5 {
		t.Error("Abs(5) should return 5")
	}
}
