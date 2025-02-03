package ebimath

import (
	"fmt"
	"math"
)

// Vector represents a 2D vector with X and Y components using float64 for precision.
// This structure is inspired by Godot's Vector2 but adapted for Go conventions.
type Vector struct {
	X, Y float64
}

// Constructor functions for Vector
// --------------------------------
// V creates a new Vector with given x and y coordinates.
func V(x, y float64) Vector {
	return Vector{X: x, Y: y}
}

// VInt converts integer coordinates to a Vector.
func VInt(x, y int) Vector {
	return V(float64(x), float64(y))
}

// V2 creates a Vector where both X and Y are set to the same value.
func V2(v float64) Vector {
	return V(v, v)
}

// V2Int creates a Vector where both X and Y are set to the integer value converted to float64.
func V2Int(v int) Vector {
	return V2(float64(v))
}

// Basic Vector Operations
// -----------------------
// String returns a string representation of the Vector.
func (self Vector) String() string {
	return fmt.Sprintf("[%f, %f]", self.X, self.Y)
}

// IsZero checks if the Vector is at the origin (0, 0).
func (self Vector) IsZero() bool {
	return self.X == 0 && self.Y == 0
}

// IsNormalized checks if the vector's length is approximately 1.
func (self Vector) IsNormalized() bool {
	return EqualsApproximately(self.LengthSquared(), 1)
}

// Abs returns a new Vector with the absolute values of X and Y.
func (self Vector) Abs() Vector {
	return V(math.Abs(self.X), math.Abs(self.Y))
}

// ToInt converts the Vector to integer coordinates.
func (self Vector) ToInt() (int, int) {
	return int(self.X), int(self.Y)
}

// Vector Transformations
// ----------------------
// Apply applies a matrix transformation to this Vector.
func (self Vector) Apply(m Matrix) Vector {
	x, y := m.Apply(self.X, self.Y)
	return V(x, y)
}

// Rotate rotates the Vector by degrees.
func (self Vector) RotateDegrees(degrees float64) Vector {
	degrees = Repeat(degrees, 360) // Normalize angle
	radians := degrees * (Pi / 180)
	sine, cosine := math.Sin(radians), math.Cos(radians)
	return V(self.X*cosine-self.Y*sine, self.X*sine+self.Y*cosine)
}

// Rotate rotates the Vector by the given angle in radians.
func (self Vector) Rotate(angle float64) Vector {
	sine, cosine := math.Sin(angle), math.Cos(angle)
	return V(self.X*cosine-self.Y*sine, self.X*sine+self.Y*cosine)
}

// RotateAround rotates this Vector around another Vector by an angle in radians.
func (self Vector) RotateAround(around Vector, angle float64) Vector {
	return V(
		math.Cos(angle)*(self.X-around.X)-math.Sin(angle)*(self.Y-around.Y)+around.X,
		math.Sin(angle)*(self.X-around.X)+math.Cos(angle)*(self.Y-around.Y)+around.Y,
	)
}

// Vector Calculations
// -------------------
// DistanceTo calculates the Euclidean distance to another Vector.
func (self Vector) DistanceTo(v2 Vector) float64 {
	return math.Sqrt(self.DistanceSquaredTo(v2))
}

// DistanceSquaredTo computes the squared distance to another Vector, which is faster for comparisons.
func (self Vector) DistanceSquaredTo(v2 Vector) float64 {
	dx, dy := self.X-v2.X, self.Y-v2.Y
	return dx*dx + dy*dy
}

// Dot computes the dot product between this Vector and another.
func (self Vector) Dot(v2 Vector) float64 {
	return self.X*v2.X + self.Y*v2.Y
}

// Length returns the magnitude of the Vector.
func (self Vector) Length() float64 {
	return math.Sqrt(self.LengthSquared())
}

// LengthSquared returns the square of the Vector's length, useful for comparisons.
func (self Vector) LengthSquared() float64 {
	return self.Dot(self)
}

// Additional Vector Operations
// ----------------------------
// Angle returns the angle of the Vector from the positive X-axis in radians.
func (self Vector) Angle() float64 {
	return math.Atan2(self.Y, self.X)
}

// AngleToPoint returns the angle from this Vector towards another point.
func (self Vector) AngleToPoint(other Vector) float64 {
	return other.Sub(self).Angle()
}

// DirectionTo returns a normalized vector pointing from this Vector to another.
func (self Vector) DirectionTo(other Vector) Vector {
	return self.Sub(other).Normalized()
}

// VecTowards calculates a Vector of given length towards another point.
func (self Vector) VecTowards(other Vector, length float64) Vector {
	angle := self.AngleToPoint(other)
	return V(math.Cos(angle), math.Sin(angle)).MulF(length)
}

// MoveTowards moves the Vector towards another Vector by a maximum distance.
func (self Vector) MoveTowards(other Vector, length float64) Vector {
	direction := other.Sub(self)
	dist := direction.Length()
	if dist <= length || dist < Epsilon {
		return other
	}
	return self.Add(direction.DivF(dist).MulF(length))
}

// Vector Math Operations
// ----------------------
// Add adds one or more Vectors to this Vector.
func (self Vector) Add(others ...Vector) Vector {
	for _, r := range others {
		self.X += r.X
		self.Y += r.Y
	}
	return self
}

// AddF adds a scalar to both components of the Vector.
func (self Vector) AddF(scalar float64) Vector {
	self.X += scalar
	self.Y += scalar
	return self
}

// Sub subtracts one or more Vectors from this Vector.
func (self Vector) Sub(others ...Vector) Vector {
	for _, r := range others {
		self.X -= r.X
		self.Y -= r.Y
	}
	return self
}

// SubF subtracts a scalar from both components of the Vector.
func (self Vector) SubF(scalar float64) Vector {
	self.X -= scalar
	self.Y -= scalar
	return self
}

// Mul multiplies the Vector by another Vector component-wise.
func (self Vector) Mul(other Vector) Vector {
	self.X *= other.X
	self.Y *= other.Y
	return self
}

// MulF multiplies the Vector by a scalar.
func (self Vector) MulF(scalar float64) Vector {
	return V(self.X*scalar, self.Y*scalar)
}

// Div divides the Vector by another Vector component-wise.
func (self Vector) Div(other Vector) Vector {
	self.X /= other.X
	self.Y /= other.Y
	return self
}

// DivF divides the Vector by a scalar.
func (self Vector) DivF(scalar float64) Vector {
	return V(self.X/scalar, self.Y/scalar)
}

// Scale scales the Vector by another Vector component-wise.
func (self Vector) Scale(other Vector) Vector {
	self.X *= other.X
	self.Y *= other.Y
	return self
}

// ScaleF scales the Vector by a scalar.
func (self Vector) ScaleF(scalar float64) Vector {
	self.X *= scalar
	self.Y *= scalar
	return self
}

// Unit returns a unit vector in the same direction as this Vector.
func (self Vector) Unit() Vector {
	l := self.Length()
	if l != 0 {
		return self.DivF(l)
	}
	return self
}

// Normalized returns a unit vector in the same direction as this Vector.
func (self Vector) Normalized() Vector {
	l := self.LengthSquared()
	if l != 0 {
		return self.MulF(1 / math.Sqrt(l))
	}
	return self
}

// ClampLength ensures the Vector's length does not exceed a given limit.
func (self Vector) ClampLength(limit float64) Vector {
	l := self.Length()
	if l > 0 && l > limit {
		return self.DivF(l).MulF(limit)
	}
	return self
}

// Round returns a new Vector with each component rounded to the nearest integer.
func (self Vector) Round() Vector {
	return V(math.Round(self.X), math.Round(self.Y))
}

// MoveInDirection moves the Vector in the direction of the angle by a given distance.
func (self Vector) MoveInDirection(angle, distance float64) Vector {
	return self.Add(V(math.Cos(angle), math.Sin(angle)).MulF(distance))
}

// Equals checks if two Vectors are equal within a small tolerance.
func (self Vector) Equals(other Vector) bool {
	return self.Sub(other).LengthSquared() < 0.00000000009
}
