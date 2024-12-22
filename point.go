package ebimath

type Point struct {
	X, Y int
}

func P(x, y int) Point {
	return Point{x, y}
}

func Pf(x, y float64) Point {
	return Point{
		FastFloor[float64, int](x),
		FastFloor[float64, int](y),
	}
}
