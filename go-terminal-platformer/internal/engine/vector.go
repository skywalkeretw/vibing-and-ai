package engine

import (
	"math"
)

// Vector2 represents a 2D vector with X and Y components
type Vector2 struct {
	X, Y float64
}

// NewVector2 creates a new Vector2
func NewVector2(x, y float64) Vector2 {
	return Vector2{X: x, Y: y}
}

// Add returns the sum of two vectors
func (v Vector2) Add(other Vector2) Vector2 {
	return Vector2{
		X: v.X + other.X,
		Y: v.Y + other.Y,
	}
}

// Subtract returns the difference of two vectors
func (v Vector2) Subtract(other Vector2) Vector2 {
	return Vector2{
		X: v.X - other.X,
		Y: v.Y - other.Y,
	}
}

// Multiply returns the vector scaled by a scalar
func (v Vector2) Multiply(scalar float64) Vector2 {
	return Vector2{
		X: v.X * scalar,
		Y: v.Y * scalar,
	}
}

// Divide returns the vector divided by a scalar
func (v Vector2) Divide(scalar float64) Vector2 {
	if scalar == 0 {
		return v
	}
	return Vector2{
		X: v.X / scalar,
		Y: v.Y / scalar,
	}
}

// Dot returns the dot product of two vectors
func (v Vector2) Dot(other Vector2) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Cross returns the cross product (z-component) of two 2D vectors
func (v Vector2) Cross(other Vector2) float64 {
	return v.X*other.Y - v.Y*other.X
}

// Length returns the magnitude of the vector
func (v Vector2) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// LengthSquared returns the squared magnitude (faster than Length)
func (v Vector2) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y
}

// Normalize returns a unit vector in the same direction
func (v Vector2) Normalize() Vector2 {
	length := v.Length()
	if length == 0 {
		return Vector2{0, 0}
	}
	return v.Divide(length)
}

// Distance returns the distance between two points
func (v Vector2) Distance(other Vector2) float64 {
	return v.Subtract(other).Length()
}

// DistanceSquared returns the squared distance (faster than Distance)
func (v Vector2) DistanceSquared(other Vector2) float64 {
	return v.Subtract(other).LengthSquared()
}

// Lerp performs linear interpolation between two vectors
func (v Vector2) Lerp(other Vector2, t float64) Vector2 {
	return Vector2{
		X: v.X + (other.X-v.X)*t,
		Y: v.Y + (other.Y-v.Y)*t,
	}
}

// Clamp clamps the vector's components between min and max
func (v Vector2) Clamp(min, max Vector2) Vector2 {
	return Vector2{
		X: clampFloat(v.X, min.X, max.X),
		Y: clampFloat(v.Y, min.Y, max.Y),
	}
}

// Rotate rotates the vector by the given angle in radians
func (v Vector2) Rotate(angle float64) Vector2 {
	cos := math.Cos(angle)
	sin := math.Sin(angle)
	return Vector2{
		X: v.X*cos - v.Y*sin,
		Y: v.X*sin + v.Y*cos,
	}
}

// Perpendicular returns a vector perpendicular to this one
func (v Vector2) Perpendicular() Vector2 {
	return Vector2{X: -v.Y, Y: v.X}
}

// Zero returns a zero vector
func Zero() Vector2 {
	return Vector2{0, 0}
}

// One returns a vector with both components set to 1
func One() Vector2 {
	return Vector2{1, 1}
}

// Up returns an up vector (0, -1)
func Up() Vector2 {
	return Vector2{0, -1}
}

// Down returns a down vector (0, 1)
func Down() Vector2 {
	return Vector2{0, 1}
}

// Left returns a left vector (-1, 0)
func Left() Vector2 {
	return Vector2{-1, 0}
}

// Right returns a right vector (1, 0)
func Right() Vector2 {
	return Vector2{1, 0}
}

// Rectangle represents an axis-aligned bounding box
type Rectangle struct {
	X, Y, Width, Height float64
}

// NewRectangle creates a new Rectangle
func NewRectangle(x, y, width, height float64) Rectangle {
	return Rectangle{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// Left returns the left edge of the rectangle
func (r Rectangle) Left() float64 {
	return r.X
}

// Right returns the right edge of the rectangle
func (r Rectangle) Right() float64 {
	return r.X + r.Width
}

// Top returns the top edge of the rectangle
func (r Rectangle) Top() float64 {
	return r.Y
}

// Bottom returns the bottom edge of the rectangle
func (r Rectangle) Bottom() float64 {
	return r.Y + r.Height
}

// Center returns the center point of the rectangle
func (r Rectangle) Center() Vector2 {
	return Vector2{
		X: r.X + r.Width/2,
		Y: r.Y + r.Height/2,
	}
}

// Contains checks if a point is inside the rectangle
func (r Rectangle) Contains(point Vector2) bool {
	return point.X >= r.X &&
		point.X <= r.X+r.Width &&
		point.Y >= r.Y &&
		point.Y <= r.Y+r.Height
}

// Intersects checks if two rectangles intersect
func (r Rectangle) Intersects(other Rectangle) bool {
	return r.X < other.X+other.Width &&
		r.X+r.Width > other.X &&
		r.Y < other.Y+other.Height &&
		r.Y+r.Height > other.Y
}

// Intersection returns the intersection rectangle of two rectangles
func (r Rectangle) Intersection(other Rectangle) Rectangle {
	x := math.Max(r.X, other.X)
	y := math.Max(r.Y, other.Y)
	width := math.Min(r.Right(), other.Right()) - x
	height := math.Min(r.Bottom(), other.Bottom()) - y

	if width < 0 || height < 0 {
		return Rectangle{} // No intersection
	}

	return Rectangle{X: x, Y: y, Width: width, Height: height}
}

// Union returns the smallest rectangle that contains both rectangles
func (r Rectangle) Union(other Rectangle) Rectangle {
	x := math.Min(r.X, other.X)
	y := math.Min(r.Y, other.Y)
	width := math.Max(r.Right(), other.Right()) - x
	height := math.Max(r.Bottom(), other.Bottom()) - y

	return Rectangle{X: x, Y: y, Width: width, Height: height}
}

// Expand expands the rectangle by the given amount in all directions
func (r Rectangle) Expand(amount float64) Rectangle {
	return Rectangle{
		X:      r.X - amount,
		Y:      r.Y - amount,
		Width:  r.Width + amount*2,
		Height: r.Height + amount*2,
	}
}

// Translate moves the rectangle by the given offset
func (r Rectangle) Translate(offset Vector2) Rectangle {
	return Rectangle{
		X:      r.X + offset.X,
		Y:      r.Y + offset.Y,
		Width:  r.Width,
		Height: r.Height,
	}
}

// Helper function to clamp a float value
func clampFloat(value, min, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}
	return value
}

// Min returns the minimum of two float64 values
func Min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

// Max returns the maximum of two float64 values
func Max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// Abs returns the absolute value of a float64
func Abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
