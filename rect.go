package lsd

// TODO: reduce code reuse between rect and LineSegment types

type rect struct {
	p1, p2         Point // Endpoints of the line segment going through 'weightedCenter' and aligned with the vector <'dx', 'dy'>; also, midpoints of opposite sides of the rectangle
	width          float64
	weightedCenter Point   // Weighted sum of the coordinates of all the pixels in the region, weighted by the norm of the gradient
	theta          float64 // Angle
	dx, dy         float64 // <'dx', 'dy'> is the rectangle's primary axis
	anglePrecision float64 // Tolerance angle
	pAligned       float64 // Probability of a point in the rectangle having angle value within 'anglePrecision' of 'theta'
}

/* Since rect is made of primitive types, we do not need to define a method corresponding to rect_copy.
The assignment operator will suffice. */
