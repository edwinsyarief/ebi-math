package ebimath

// Point represents a point in 2D space with integer coordinates.
type Point struct {
	X, Y int
}

// Constructor functions for Point
// -------------------------------
// P creates a new Point with given integer x and y coordinates.
func P(x, y int) Point {
	return Point{X: x, Y: y}
}

// Pf converts floating-point coordinates to a Point with integer coordinates.
func Pf(x, y float64) Point {
	return Point{
		X: FastFloor[float64, int](x),
		Y: FastFloor[float64, int](y),
	}
}
