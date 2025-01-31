package ebimath

import (
	"cmp"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Matrix = ebiten.GeoM

const Pi = 3.141592653589793
const Epsilon = 1e-9

func EqualsApproximately[T float32 | float64](a, b T) bool {
	// Check for exact equality first, required to handle "infinity" values.
	if a == b {
		return true
	}
	// Then check for approximate equality.
	tolerance := Epsilon * Abs(float64(a))
	if tolerance < Epsilon {
		tolerance = Epsilon
	}
	return Abs(float64(a-b)) < tolerance
}

func Repeat(t, length float64) float64 {
	return Clamp(t-math.Floor(t/length)*length, 0, length)
}

func CubicInterpolate(from, to, pre, post, t float64) float64 {
	return 0.5 *
		((from * 2.0) +
			(-pre+to)*t +
			(2.0*pre-5.0*from+4.0*to-post)*(t*t) +
			(-pre+3.0*from-3.0*to+post)*(t*t*t))
}

// ToDegrees is a helper function to easily convert radians to degrees for human readability.
func ToDegrees(radians float64) float64 {
	return radians / math.Pi * 180
}

// ToRadians is a helper function to easily convert degrees to radians (which is what the rotation-oriented functions in Tetra3D use).
func ToRadians(degrees float64) float64 {
	return math.Pi * degrees / 180
}

func Lerp[T float64 | float32 | int | int16 | int32 | int64](from, to, t T) T {
	return from + ((to - from) * t)
}

func Clamp[T cmp.Ordered](value, min, max T) T {
	if value <= min {
		return min
	}
	if value >= max {
		return max
	}
	return value
}

func FastFloor[T float64 | float32, U float64 | float32 | int](value T) U {
	return U((value + 32768.0) - 32768)
}

func ClampTowardsZero[T float64 | float32 | int | int16 | int32 | int64](value, clampReference T) T {
	if clampReference > 0 {
		return min(value, clampReference)
	}
	return max(value, clampReference)
}

func Abs[T float64 | float32 | int | int8 | int16 | int32 | int64](value T) T {
	if value >= 0 {
		return value
	}
	return -value
}

// It doesn't take zero into account, but this is intentional.
func Sign[T float64 | float32 | int | int8 | int16 | int32 | int64](value T) T {
	if value >= 0 {
		return +1
	}
	return -1
}

func Max(v1, v2 float64) float64 {
	if v1 > v2 {
		return v1
	}

	return v2
}

func Min(v1, v2 float64) float64 {
	if v1 < v2 {
		return v1
	}

	return v2
}

func AngleToVector(angleRadians float64, length float64) Vector {
	return V(math.Cos(float64(angleRadians))*length, math.Sin(float64(angleRadians))*length)
}

func AdjustDestinationPixel(x float32) float32 {
	// Avoid the center of the pixel, which is problematic (#929, #1171).
	// Instead, align the vertices with about 1/3 pixels.
	//
	// The intention here is roughly this code:
	//
	//     float32(math.Floor((float64(x)+1.0/6.0)*3) / 3)
	//
	// The actual implementation is more optimized than the above implementation.
	ix := float32(int(x))
	if x < 0 && x != ix {
		ix -= 1
	}
	frac := x - ix
	switch {
	case frac < 3.0/16.0:
		return ix
	case frac < 8.0/16.0:
		return ix + 5.0/16.0
	case frac < 13.0/16.0:
		return ix + 11.0/16.0
	default:
		return ix + 1
	}
}
