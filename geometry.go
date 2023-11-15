package rc

import "math"

type Distance = float64

type Point2D struct {
	X float64
	Y float64
}

func Distance2D(a, b Point2D) (s Distance) {
	dx := math.Abs(a.X - b.X)
	dy := math.Abs(a.Y - b.Y)
	return math.Sqrt(dx*dx + dy*dy)
}
