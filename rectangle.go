package ebimath

import "math"

// Rectangle represents a 2D rectangle with min and max vectors for bounds and an orientation angle.
type Rectangle struct {
	Min   Vector
	Max   Vector
	Angle float64
}

// Constructor for Rectangle
// -------------------------
// NewRectangle creates a new rectangle from two points.
func NewRectangle(x1, y1, x2, y2 float64) Rectangle {
	return Rectangle{
		Min: V(x1, y1),
		Max: V(x2, y2),
	}
}

// Basic Properties
// ----------------
// Width returns the width of the rectangle.
func (r Rectangle) Width() float64 {
	return r.Max.X - r.Min.X
}

// Height returns the height of the rectangle.
func (r Rectangle) Height() float64 {
	return r.Max.Y - r.Min.Y
}

// Center calculates and returns the center point of the rectangle.
func (r Rectangle) Center() Vector {
	return V((r.Min.X+r.Max.X)/2, (r.Min.Y+r.Max.Y)/2)
}

// Accessors for coordinates
// --------------------
func (r Rectangle) X1() float64 { return r.Min.X }
func (r Rectangle) Y1() float64 { return r.Min.Y }
func (r Rectangle) X2() float64 { return r.Max.X }
func (r Rectangle) Y2() float64 { return r.Max.Y }

// State Checks
// ------------
// IsEmpty checks if the rectangle has no area.
func (r Rectangle) IsEmpty() bool {
	return r.Min.X >= r.Max.X || r.Min.Y >= r.Max.Y
}

// Equals checks if two rectangles are equal.
func (r Rectangle) Equals(other Rectangle) bool {
	return r.Min.Equals(other.Min) && r.Max.Equals(other.Max)
}

// Containment and Intersection
// ---------------------------
// Contains checks if a point is within the rectangle.
func (r Rectangle) Contains(p Vector) bool {
	return r.Min.X <= p.X && p.X < r.Max.X &&
		r.Min.Y <= p.Y && p.Y < r.Max.Y
}

// ContainsRect checks if one rectangle is completely inside another.
func (r Rectangle) ContainsRect(other Rectangle) bool {
	return r.X1() <= other.X1() &&
		r.Y1() <= other.Y1() &&
		r.X2() >= other.X2() &&
		r.Y2() >= other.Y2()
}

// Intersects checks if two rectangles intersect, considering their rotation.
func (r Rectangle) Intersects(other Rectangle) bool {
	if r.Angle == 0 && other.Angle == 0 {
		return !r.IsEmpty() && !other.IsEmpty() &&
			r.Min.X < other.Max.X && other.Min.X < r.Max.X &&
			r.Min.Y < other.Max.Y && other.Min.Y < r.Max.Y
	}

	axes := []Vector{
		r.GetAxis(r.Angle, 0), r.GetAxis(r.Angle, 1),
		other.GetAxis(other.Angle, 0), other.GetAxis(other.Angle, 1),
	}

	for _, axis := range axes {
		if !r.OverlapOnAxis(other, axis) {
			return false
		}
	}
	return true
}

// IntersectsCircle checks if the rectangle intersects with a circle.
func (r Rectangle) IntersectsCircle(center ebimath.Vector, radius float64) bool {
	closestX := math.Max(r.Min.X, math.Min(center.X, r.Max.X))
	closestY := math.Max(r.Min.Y, math.Min(center.Y, r.Max.Y))
	dx := center.X - closestX
	dy := center.Y - closestY
	return dx*dx+dy*dy <= radius*radius
}

// Helper methods for intersection calculation
// -------------------------------------------
// GetAxis returns one of the two axes of the rectangle based on its angle.
func (r Rectangle) GetAxis(angle float64, index int) Vector {
	if index == 0 {
		return Vector{math.Cos(angle), math.Sin(angle)}
	}
	return Vector{-math.Sin(angle), math.Cos(angle)}
}

// GetCorners returns the four corners of the rectangle, considering rotation.
func (r Rectangle) GetCorners() [4]Vector {
	center := r.Center()
	halfWidth := r.Width() / 2
	halfHeight := r.Height() / 2
	cos, sin := math.Cos(r.Angle), math.Sin(r.Angle)

	return [4]Vector{
		{center.X + halfWidth*cos - halfHeight*sin, center.Y + halfWidth*sin + halfHeight*cos},
		{center.X - halfWidth*cos - halfHeight*sin, center.Y - halfWidth*sin + halfHeight*cos},
		{center.X - halfWidth*cos + halfHeight*sin, center.Y - halfWidth*sin - halfHeight*cos},
		{center.X + halfWidth*cos + halfHeight*sin, center.Y + halfWidth*sin - halfHeight*cos},
	}
}

// OverlapOnAxis checks if there's overlap along a specific axis.
func (r Rectangle) OverlapOnAxis(other Rectangle, axis Vector) bool {
	proj1 := r.ProjectOntoAxis(axis)
	proj2 := other.ProjectOntoAxis(axis)
	return proj1.Min <= proj2.Max && proj2.Min <= proj1.Max
}

// ProjectOntoAxis projects the rectangle onto an axis, returning the min and max values.
func (r Rectangle) ProjectOntoAxis(axis Vector) struct{ Min, Max float64 } {
	corners := r.GetCorners()
	min := corners[0].Dot(axis)
	max := min

	for i := 1; i < 4; i++ {
		p := corners[i].Dot(axis)
		if p < min {
			min = p
		} else if p > max {
			max = p
		}
	}

	return struct{ Min, Max float64 }{min, max}
}
