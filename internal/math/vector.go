package math

import (
	"math"
)

// Vector2D represents a 2D vector with X and Y coordinates
type Vector2D struct {
	X float64
	Y float64
}

// NewVector2D creates a new Vector2D
func NewVector2D(x, y float64) Vector2D {
	return Vector2D{X: x, Y: y}
}

// Add returns the sum of two vectors
func (v Vector2D) Add(other Vector2D) Vector2D {
	return Vector2D{X: v.X + other.X, Y: v.Y + other.Y}
}

// Sub returns the difference of two vectors
func (v Vector2D) Sub(other Vector2D) Vector2D {
	return Vector2D{X: v.X - other.X, Y: v.Y - other.Y}
}

// Mul returns the vector multiplied by a scalar
func (v Vector2D) Mul(scalar float64) Vector2D {
	return Vector2D{X: v.X * scalar, Y: v.Y * scalar}
}

// Distance returns the distance between two vectors
func (v Vector2D) Distance(other Vector2D) float64 {
	dx := v.X - other.X
	dy := v.Y - other.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Length returns the length of the vector
func (v Vector2D) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y)
}

// Normalize returns a normalized vector (length = 1)
func (v Vector2D) Normalize() Vector2D {
	length := v.Length()
	if length == 0 {
		return Vector2D{X: 0, Y: 0}
	}
	return Vector2D{X: v.X / length, Y: v.Y / length}
}

// Dot returns the dot product of two vectors
func (v Vector2D) Dot(other Vector2D) float64 {
	return v.X*other.X + v.Y*other.Y
}

// Angle returns the angle of the vector in radians
func (v Vector2D) Angle() float64 {
	return math.Atan2(v.Y, v.X)
}
