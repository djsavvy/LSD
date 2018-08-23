package lsd

import (
	"image"
	"log"
	"math"
)

// Scale the input image 'img' by a factor 'scale' using Gaussian sub-sampling.
// For example, scale=0.8 will give a result at 80% of the original size.
//
// The image is convolved with a Gaussian kernel
//
//     G(x,y) = \frac{1}{2\pi\sigma^2} e^{-\frac{x^2+y^2}{2\sigma^2}}
//
// before the sub-sampling to prevent aliasing.
//
// The standard deviation sigma given by:
// -  sigma = sigma_scale / scale,   if scale <  1.0
// -  sigma = sigma_scale,           if scale >= 1.0
//
// Note that if scale == 1.0, then no blurring is performed. This is to preserve
// statistics of the *a contrario* model. If blurring is performed without scaling,
// then some structures would be detected on a blurred white noise. [1]
//
// To be able to sub-sample at non-integer steps, some interpolation
// is needed. In this implementation, the interpolation is done by
// the Gaussian kernel, so both operations (filtering and sampling)
// are done at the same time. The Gaussian kernel is computed
// centered on the coordinates of the required sample. In this way,
// when applied, it gives directly the result of convolving the image
// with the kernel and interpolated to that particular position.
//
// A fast algorithm is done using the separability of the Gaussian
// kernel. Applying the 2D Gaussian kernel is equivalent to applying
// first a horizontal 1D Gaussian kernel and then a vertical 1D
// Gaussian kernel (or the other way round). The reason is that
//
//     G(x,y) = G(x) * G(y)
//
// where
//
//     G(x) = \frac{1}{\sqrt{2\pi}\sigma} e^{-\frac{x^2}{2\sigma^2}}.
//
// The algorithm first applies a combined Gaussian kernel and sampling
// in the x axis, and then the combined Gaussian kernel and sampling
// in the y axis.
//
// [1] http://www.ipol.im/pub/art/2012/gjmr-lsd/article.pdf
func gaussianSubSample(img image.Gray, scale, sigmaScale float64) [][]float64 {
	if scale <= 0 {
		log.Fatalf("gaussianSubSample: 'scale' must be positive. Received value %v\n", scale)
	}
	if sigmaScale <= 0 {
		log.Fatalf("gaussianSubSample: 'sigmaScale' must be positive. Received value %v\n", sigmaScale)
	}

	// Generate result
	oldXSize := img.Bounds().Size().X
	oldYSize := img.Bounds().Size().Y
	newXSize := int(math.Ceil(scale * float64(oldXSize)))
	newYSize := int(math.Ceil(scale * float64(oldYSize)))
	result := make([][]float64, newXSize)
	for i := range result {
		result[i] = make([]float64, newYSize)
	}

	// Deal with no scale/blur case
	if scale == 1.0 {
		for x := range result {
			for y := range result[x] {
				result[x][y] = float64(img.GrayAt(x, y).Y)
			}
		}
		return result
	}

	// Calculate sigma for the Gaussian kernel
	sigma := sigmaScale
	if scale < 1 {
		sigma /= scale
	}

	// Calculate length of the Gaussian kernel. Its length is selected to guarantee that
	// the first discarded term is at least 10^'prec' times smaller than the central value.
	prec := 3.0
	len := int(math.Ceil(sigma*math.Sqrt(2*prec*math.Ln10)))*2 + 1

	// Create temporary semi-scaled image
	// TODO: Refactor out code to create a 2D slice of float64s
	afterXScaling := make([][]float64, newXSize)
	for i := range afterXScaling {
		afterXScaling[i] = make([]float64, oldYSize)
	}

	// First subsampling: x-axis:
	/*
		newX		is the coordinate in the new image
		oldX		is the corresponding x-value in the original size image
		oldXpixel	is the integer value, the pixel coordinate of oldX

		Note that coordinate (0.0, 0.0) is in the center of pixel (0, 0),
		so the pixel with 'oldXpixel' value of 0 covers possible values of
		'oldX' from -0.5 to 0.5
	*/
	for newX := range afterXScaling {
		oldX := float64(newX) / scale
		oldXpixel := int(math.Floor(oldX + 0.5))

		// Kernel must be recomputed for each 'newX' because the fine offset
		// 'oldX' - 'oldXpixel' is different in each case
		kernel := gaussianKernel(len, sigma, oldX-float64(oldXpixel)+float64(len/2))

		for y := range afterXScaling[newX] {
			sum := float64(0)

			for i := range kernel {
				j := oldXpixel - (len / 2) + i

				// Symmetry boundary condition
				for j < 0 {
					j += oldXSize * 2
				}
				for j >= oldXSize*2 {
					j -= oldXSize * 2
				}
				if j >= oldXSize {
					j = 2*oldXSize - 1 - j
				}

				sum += float64(img.GrayAt(j, y).Y) * kernel[i]
			}
			afterXScaling[newX][y] = sum
		}
	}

	// Second subsampling: y-axis
	/*
		newY		is the coordinate in the new image
		oldY		is the corresponding y-value in the original size image
		oldYpixel	is the integer value, the pixel coordinate of oldY

		Note that coordinate (0.0, 0.0) is in the center of pixel (0, 0),
		so the pixel with 'oldYpixel' value of 0 covers possible values of 'oldY' from -0.5 to 0.5
	*/
	for newY := range result {
		oldY := float64(newY) / scale
		oldYpixel := int(math.Floor(oldY + 0.5))

		// Kernel must be recomputed for each 'newY' because the fine offset
		// 'oldY' - 'oldYpixel' is different in each case
		kernel := gaussianKernel(len, sigma, oldY-float64(oldYpixel)+float64(len/2))

		for x := range result {
			sum := float64(0)

			for i := range kernel {
				j := oldYpixel - (len / 2) + i

				// Symmetry boundary condition
				for j < 0 {
					j += oldYSize * 2
				}
				for j >= oldYSize*2 {
					j -= oldYSize * 2
				}
				if j >= oldYSize {
					j = 2*oldYSize - 1 - j
				}

				sum += afterXScaling[x][j] * kernel[i]
			}
			result[x][newY] = sum
		}
	}

	return result
}

// Converts an input image.Image to an image.Gray of the same dimensions. This is necessary
// since the LSD algorithm requires a grayscaled input.
func makeGrayscale(img image.Image) image.Gray {
	result := image.NewGray(img.Bounds())
	// According to the image package documentation:
	// "Looping over Y first and X second is more likely to result in better
	// memory access patterns than X first and Y second."
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			result.Set(x, y, img.At(x, y))
		}
	}
	return *result
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
