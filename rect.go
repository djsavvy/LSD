package lsd

// TODO: reduce code reuse between rect and LineSegment types

type rect struct {
	p1, p2           Point // Endpoints of the corresponding line segment
	width            float64
	centerX, centerY float64 // Center of the rectangle
	theta            float64 // Angle
	dx, dy           float64 // (dx, dy) is a vector oriented as the line segment
	anglePrecision   float64 // Tolerance angle
	pAligned         float64 // Probability of a point with angle within 'anglePrecision'
}

/* Since rect is made of primitive types, we do not need to define a method corresponding to rect_copy.
The assignment operator will suffice. */
