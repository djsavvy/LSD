package lsd

import (
	"math"
)

// Compares float64s by relative error.
//
// The resulting rounding error after floating point computations
// depend on the specific operations done. The same number computed by
// different algorithms could present different rounding errors. For a
// useful comparison, an estimation of the relative rounding error
// should be considered and compared to a factor times EPS. The factor
// should be related to the cumulated rounding error in the chain of
// computation. Here, as a simplification, a fixed factor is used.
func equalWithinError(a, b float64) bool {
	/* Trivial case */
	if a == b {
		return true
	}

	diffAbs := math.Abs(a - b)
	aAbs := math.Abs(a)
	bAbs := math.Abs(b)
	maxAbs := math.Max(aAbs, bAbs)

	const relativeErrorFactor = 100.0
	// Smallest positive normalized float64
	const minNormalizedFloat64 = 2.22507385850720138309023271733240406421921598046233e-308
	// Difference between 1.0 and the minimum float64 greater than 1
	const epsilonFloat64 = 2.22044604925031308084726333618164062e-16

	// 'minNormalizedFloat64' is the smallest normalized number, thus, the smallest
	// number whose relative error is bounded by 'epsilonFloat64'. For
	// smaller numbers, the same quantization steps as for 'minNormalizedFloat64'
	// are used. Then, for smaller numbers, a meaningful "relative"
	// error should be computed by dividing the difference by 'minNormalizedFloat64'.
	if maxAbs < minNormalizedFloat64 {
		maxAbs = minNormalizedFloat64
	}

	return (diffAbs / maxAbs) <= (relativeErrorFactor * epsilonFloat64)
}

// Computes Euclidean distance between points (x1, y1) and (x2, y2)
func dist(x1, x2, y1, y2 float64) float64 {
	return math.Hypot(x2-x1, y2-y1)
}

// Compute absolute value angle difference
func absAngleDiff(a, b float64) float64 {
	return math.Abs(signedAngleDiff(a, b))
}

// Compute signed angle difference
func signedAngleDiff(a, b float64) float64 {
	diff := a - b
	for diff <= -math.Pi {
		diff += 2 * math.Pi
	}
	for diff > math.Pi {
		diff -= 2 * math.Pi
	}
	return diff
}
