package lsd

// LineSegment defines a line segment detected by the LSD algorithm.
type LineSegment struct {
	X1, Y1, X2, Y2   float64 // Endpoints
	Width            float64
	AnglePrecision   float64 // In (0, 1) given by the angle tolerance divided by 180 degrees
	NegativeLog10NFA float64
	// TODO: Explain NFA values in documentation
}
