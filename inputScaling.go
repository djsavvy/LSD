package lsd

import (
	"log"
	"math"
)

// Computes a Gaussian kernel, with length 'len' and standard deviation 'sigma', that is
// centered at 'mean'.
//
// For example, if 'mean' equals 0.5, the Gaussian will be centered in the middle point between the
// values in the 0th and 1st indices of the output slice.
func gaussianKernel(len int, sigma, mean float64) []float64 {
	if len < 0 {
		log.Fatalf("gaussianKernel: 'len' must be non-negative. Received value %v\n", len)
	}
	if sigma <= 0 {
		log.Fatalf("gaussianKernel: 'sigma' must be positive. Received value %v\n", sigma)
	}

	kernel := make([]float64, len)
	sum := float64(0)

	// Compute kernel
	for i := range kernel {
		val := (float64(i) - mean) / sigma
		kernel[i] = math.Exp(-0.5 * val * val)
		sum += kernel[i]
	}

	// Normalize kernel
	if sum >= 0.0 {
		for i := range kernel {
			kernel[i] /= sum
		}
	}

	return kernel
}
