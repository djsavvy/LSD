package lsd

import (
	"image"
	"log"
	"math"
)

// Converts an input image.Image to an *image.Gray of the same dimensions. This is necessary
// since the LSD algorithm requires a grayscaled input.
func makeGrayscale(img image.Image) *image.Gray {
	result := image.NewGray(img.Bounds())
	// According to the image package documentation:
	// "Looping over Y first and X second is more likely to result in better
	// memory access patterns than X first and Y second."
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			result.Set(x, y, img.At(x, y))
		}
	}
	return result
}

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
