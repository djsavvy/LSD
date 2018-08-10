package lsd

// LineSegment defines a line segment detected by the LSD algorithm.
type LineSegment struct {
	x1, y1, x2, y2   float64 // Endpoints
	width            float64
	anglePrecision   float64 // In (0, 1) given by the angle tolerance divided by 180 degrees
	negativeLog10NFA float64
	// TODO: Explain NFA values in documentation
}
