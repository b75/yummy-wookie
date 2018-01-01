package gameutil

import (
	"math"
)

func VectorLength(x, y float64) float64 {
	return math.Sqrt(x*x + y*y)
}
