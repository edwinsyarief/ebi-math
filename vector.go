package ebimath

import (
	"fmt"
	"math"
)

// Vector is a 2-element structure that is used to represent positions,
// velocities, and other kinds numerical pairs.
//
// Its implementation as well as its API is inspired by Vector2 type
// of the Godot game engine. Where feasible, its adjusted to fit Go
// coding conventions better. Also, earlier versions of Godot used
// 32-bit values for X and Y; our vector uses 64-bit values.
//
// Since Go has no operator overloading, we implement scalar forms of
// operations with "f" suffix. So, Add() is used to add two vectors
// while Addf() is used to add scalar to the vector.
type Vector struct {
	X, Y float64
}

// Vector - Vec
func V(x, y float64) Vector {
	return Vector{X: x, Y: y}
}

func VInt(x, y int) Vector {
	return V(float64(x), float64(y))
}

func V2(v float64) Vector {
	return V(v, v)
}

func V2Int(v int) Vector {
	return V2(float64(v))
}

// RadianToVector converts a given angle into a normalized vector that encodes that direction.
func RadianToVector(angle float64) Vector {
	return Vector{X: math.Cos(angle), Y: math.Sin(angle)}
}

func (self Vector) Abs() Vector {
	return V(math.Abs(self.X), math.Abs(self.Y))
}

func (self Vector) ToInt() (int, int) {
	return int(self.X), int(self.Y)
}

func (self Vector) Apply(m Matrix) Vector {
	x, y := m.Apply(self.X, self.Y)
	return V(x, y)
}

// String returns a pretty-printed representation of a 2D vector object.
func (self Vector) String() string {
	return fmt.Sprintf("[%f, %f]", self.X, self.Y)
}

// IsZero reports whether v is a zero value vector.
// A zero value vector has X=0 and Y=0, created with Vec{}.
//
// The zero value vector has a property that its length is 0,
// but not all zero length vectors are zero value vectors.
func (self Vector) IsZero() bool {
	return self.X == 0 && self.Y == 0
}

// IsNormalizer reports whether the vector is normalized.
// A vector is considered to be normalized if its length is 1.
func (self Vector) IsNormalized() bool {
	return EqualsApproximately(self.LengthSquared(), 1)
}

// DistanceTo calculates the distance between the two vectors.
func (self Vector) DistanceTo(v2 Vector) float64 {
	return math.Sqrt((self.X-v2.X)*(self.X-v2.X) + (self.Y-v2.Y)*(self.Y-v2.Y))
}

func (self Vector) DistanceSquaredTo(v2 Vector) float64 {
	return ((self.X - v2.X) * (self.X - v2.X)) + ((self.Y - v2.Y) * (self.Y - v2.Y))
}

// Dot returns a dot-product of the two vectors.
func (self Vector) Dot(v2 Vector) float64 {
	return (self.X * v2.X) + (self.Y * v2.Y)
}

// Length reports the length of this vector (also known as magnitude).
func (self Vector) Length() float64 {
	return math.Sqrt(self.LengthSquared())
}

// LengthSquared returns the squared length of this vector.
//
// This function runs faster than Len(),
// so prefer it if you need to compare vectors
// or need the squared distance for some formula.
func (self Vector) LengthSquared() float64 {
	return self.Dot(self)
}

func (self Vector) RotateDegrees(degrees float64) Vector {
	if degrees < 0 || degrees >= 360 {
		degrees = Repeat(degrees, 360)
	}

	if degrees == 0 {
		return self
	}

	if degrees == 90 {
		return V(0-self.Y, self.X)
	}

	if degrees == 180 {
		return V(0-self.X, 0-self.Y)
	}

	if degrees == 270 {
		return V(self.Y, 0-self.X)
	}

	w := math.Sin(degrees * (Pi / 180))
	h := math.Cos(degrees * (Pi / 180))
	x := self.X
	y := self.Y
	self.X = h*x - w*y
	self.Y = w*x + h*y

	return self
}

func (self Vector) Rotate(angle float64) Vector {
	sine := math.Sin(angle)
	cosi := math.Cos(angle)
	return Vector{
		X: self.X*cosi - self.Y*sine,
		Y: self.X*sine + self.Y*cosi,
	}
}

func (self Vector) RotateAround(around Vector, angle float64) Vector {
	return V(
		math.Cos(angle)*(self.X-around.X)-math.Sin(angle)*(self.Y-around.Y)+around.X,
		math.Sin(angle)*(self.X-around.X)+math.Cos(angle)*(self.Y-around.Y)+around.Y)
}

func (self Vector) Angle() float64 {
	return math.Atan2(self.Y, self.X)
}

// AngleToPoint returns the angle from v towards the given point.
func (self Vector) AngleToPoint(other Vector) float64 {
	return other.Sub(self).Angle()
}

func (self Vector) DirectionTo(other Vector) Vector {
	return self.Sub(other).Normalized()
}

func (self Vector) VecTowards(other Vector, length float64) Vector {
	angle := self.AngleToPoint(other)
	result := Vector{X: math.Cos(angle), Y: math.Sin(angle)}
	return result.MulF(length)
}

func (self Vector) MoveTowards(other Vector, length float64) Vector {
	direction := other.Sub(self) // Not normalized
	dist := direction.Length()
	if dist <= length || dist < Epsilon {
		return other
	}
	return self.Add(direction.DivF(dist).MulF(length))
}

func (self Vector) EqualsApproximately(other Vector) bool {
	return EqualsApproximately(self.X, other.X) && EqualsApproximately(self.Y, other.Y)
}

func (self Vector) Equals(other Vector) bool {
	return self.Sub(other).LengthSquared() < 0.00000000009
}

func (self Vector) MoveInDirection(dist, dir float64) Vector {
	return Vector{
		X: self.X + dist*math.Cos(dir),
		Y: self.Y + dist*math.Sin(dir),
	}
}

func (self Vector) MulF(scalar float64) Vector {
	return Vector{
		X: self.X * scalar,
		Y: self.Y * scalar,
	}
}

func (self Vector) Mul(other Vector) Vector {
	return Vector{
		X: self.X * other.X,
		Y: self.Y * other.Y,
	}
}

func (self Vector) DivF(scalar float64) Vector {
	return Vector{
		X: self.X / scalar,
		Y: self.Y / scalar,
	}
}

func (self Vector) Div(other Vector) Vector {
	return Vector{
		X: self.X / other.X,
		Y: self.Y / other.Y,
	}
}

func (self Vector) Add(others ...Vector) Vector {
	for _, r := range others {
		self.X += r.X
		self.Y += r.Y
	}

	return self
}

func (self Vector) AddF(scalar float64) Vector {
	self.X += scalar
	self.Y += scalar

	return self
}

func (self Vector) Sub(others ...Vector) Vector {
	for _, r := range others {
		self.X -= r.X
		self.Y -= r.Y
	}

	return self
}

func (self Vector) SubF(scalar float64) Vector {
	self.X -= scalar
	self.Y -= scalar

	return self
}

// Normalized returns the vector scaled to unit length.
// Functionally equivalent to `v.Divf(v.Len())`.
//
// Special case: for zero value vectors it returns that unchanged.
func (self Vector) Normalized() Vector {
	l := self.LengthSquared()
	if l != 0 {
		return self.MulF(1 / math.Sqrt(l))
	}
	return self
}

func (self Vector) FastNormalize() Vector {
	h := math.Hypot(self.X, self.Y)
	if h == 0 {
		return V2(0)
	}

	return V(self.X/h, self.Y/h)
}

func (self Vector) ClampLength(limit float64) Vector {
	l := self.Length()
	if l > 0 && l > limit {
		self = self.DivF(l)
		self = self.MulF(limit)
	}
	return self
}

func (self Vector) Clamp(min, max float64) Vector {
	return V(Clamp(self.X, min, max), Clamp(self.Y, min, max))
}

// Negative applies unary minus (-) to the vector.
func (self Vector) Negative() Vector {
	return Vector{
		X: -self.X,
		Y: -self.Y,
	}
}

// CubicInterpolate interpolates between self (this vector) and b using
// preA and postB as handles.
// The t arguments specifies the interpolation progression (a value from 0 to 1).
// With t=0 it returns a, with t=1 it returns b.
func (self Vector) CubicInterpolate(preA, b, postB Vector, t float64) Vector {
	res := self
	res.X = CubicInterpolate(res.X, b.X, preA.X, postB.X, t)
	res.Y = CubicInterpolate(res.Y, b.Y, preA.Y, postB.Y, t)
	return res
}

// Lerp interpolates between two points by a normalized value.
// This function is commonly named "lerp".
func (self Vector) Lerp(other Vector, t float64) Vector {
	return Vector{
		X: Lerp(self.X, other.X, t),
		Y: Lerp(self.Y, other.Y, t),
	}
}

func (self Vector) Scale(scale Vector) Vector {
	return V(
		self.X*scale.X,
		self.Y*scale.Y,
	)
}

func (self Vector) ScaleF(scalar float64) Vector {
	return V(
		self.X*scalar,
		self.Y*scalar,
	)
}

func (self Vector) ToPoint() Point {
	return Pf(self.X, self.Y)
}

func (self Vector) Round() Vector {
	return V(math.Round(self.X), math.Round(self.Y))
}

func (self Vector) Floor() Vector {
	return V(math.Floor(self.X), math.Floor(self.Y))
}

func (self Vector) Ceil() Vector {
	return V(math.Ceil(self.X), math.Ceil(self.Y))
}

func (self Vector) Unit() Vector {
	l := self.Length()
	if l < 1e-8 || l == 1 {
		// If it's 0, then don't modify the vector
		return self
	}
	self.X, self.Y = self.X/l, self.Y/l
	return self
}

func (self Vector) Invert() Vector {
	self.X = -self.X
	self.Y = -self.Y

	return self
}

func (self Vector) Sign() Vector {
	self.X = Sign(self.X)
	self.Y = Sign(self.Y)
	return self
}

// Copy the value of another vector to a new one.
func (self Vector) Copy(other Vector) Vector {
	return V(other.X, other.Y)
}

// Clone this vector coordinates to a new vector with the same coordinates as this one.
func (self Vector) Clone() Vector {
	return V2(0).Copy(self)
}

// Orthogonal returns a new vector perpendicular / orthogonal from this one.
func (self Vector) Orthogonal() Vector {
	return V(self.Y, -self.X)
}

// Project this vector onto another one
func (self Vector) Project(other Vector) Vector {
	amt := self.Dot(other) / other.LengthSquared()
	return V(amt*other.X, amt*other.Y)
}

// ProjectN this vector onto a unit vector
func (self Vector) ProjectN(other Vector) Vector {
	amt := self.Dot(other)
	return V(amt*other.X, amt*other.Y)
}

// Reflect this vector on an arbitrary axis vector
func (self Vector) Reflect(other Vector) Vector {
	return self.Project(other).ScaleF(2).Sub(self)
}

// ReflectN this vector on an arbitrary axis unit vector
func (self Vector) ReflectN(axis Vector) Vector {
	return self.ProjectN(axis).ScaleF(2).Sub(self)
}

func (self Vector) AdjustPixel() Vector {
	return V(
		float64(AdjustDestinationPixel(float32(self.X))),
		float64(AdjustDestinationPixel(float32(self.Y))),
	)
}
