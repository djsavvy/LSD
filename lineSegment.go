package lsd

// LineSegment defines a line segment detected by the LSD algorithm.
type LineSegment struct {
	X1, Y1, X2, Y2 float64 // Endpoints
	Width          float64
	AnglePrecision float64 // In (0, 1) given by the angle tolerance divided by 180 degrees
	Confidence     float64 // The negative base-10 logarithm of the NFA (Number of False Alarms) value.
	/* NFA of a line segment L is the expected number of line segments as rare as L found under the stochastic
	*a contrario* model H_0 defined by random, unstructured data (i.i.d. uniformly distributed
	level-line angles for each pixel).  [1]

	[1] http://www.ipol.im/pub/art/2012/gjmr-lsd/ */
}
