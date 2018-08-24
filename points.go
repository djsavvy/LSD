package lsd

import (
	"math"
)

// Point represents a point in the 2-D space of the image.
type Point struct {
	X, Y float64
}

// Computes Euclidean distance between two points.
func dist(p1, p2 Point) float64 {
	return math.Hypot(p2.X-p1.X, p2.Y-p1.Y)
}

// TODO: combine pixel and point into an interface
type pixel struct {
	x, y int
}
