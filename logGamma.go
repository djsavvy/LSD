package lsd

import (
	"math"
)

// UseApproxLogGamma denotes whether to use the approximations to LogGamma
// found in the original LSD source code.
var UseApproxLogGamma = false

// Computes the natural logarithm of the absolute value of the Gamma function.
// If 'UseApproxLogGamma' is set to true, calls 'approxLogGamma(x)'. Otherwise,
// calls the Go standard library's builtin method.
func logGamma(x float64) float64 {
	if UseApproxLogGamma {
		return approxLogGamma(x)
	}
	// Sign value is sign of the imaginary part
	value, _ := math.Lgamma(x)
	return value
}

// Computes the natural logarithm of the absolute value of the Gamma function,
// using the approximations to LogGamma found in the original LSD source code.
// When x>15 use the Windschitl Approximation; otherwise use the Lanczos approximation.
func approxLogGamma(x float64) float64 {
	if x > 15 {
		return logGammaWindschitlApprox(x)
	}
	return logGammaLanczosApprox(x)
}

// Computes the natural logarithm of the absolute value of
// the gamma function of x using the Lanczos approximation.
//
// The formula used is
//
//   \log\Gamma(x) = \log\left( \sum_{n=0}^{N} q_n x^n \right)
//                   + (x+0.5) \log(x+5.5) - (x+5.5) - \sum_{n=0}^{N} \log(x+n)
//
// Note that q0 = 75122.6331530, q1 = 80916.6278952, q2 = 36308.2951477,
// q3 = 8687.24529705, q4 = 1168.92649479, q5 = 83.8676043424, q6 = 2.50662827511.
//
// According to [1], we only need coefficients through q6 to get an approximation of the Gamma
// function with an error of less than 2e-10 for any positive real number.
//
// [1] http://www.rskey.org/gamma.htm
func logGammaLanczosApprox(x float64) float64 {
	a := (x+0.5)*math.Log(x+5.5) - math.Log(x+5.5)
	b := float64(0.0)
	// Coefficients of the polynomial q used in the Lanczos Approximation
	var qPolynomial = [7]float64{75122.6331530, 80916.6278952, 36308.2951477,
		8687.24529705, 1168.92649479, 83.8676043424, 2.50662827511}

	for i := 0; i < 7; i++ {
		a -= math.Log(x + float64(i))
		b += qPolynomial[i] * math.Pow(x, float64(i))
	}
	return a + math.Log(b)
}

// Computes the natural logarithm of the absolute value of
// the gamma function of x using Windschitl method.
// See http://www.rskey.org/gamma.htm
//
// The formula used is
//
//     \log\Gamma(x) = 0.5\log(2\pi) + (x-0.5)\log(x) - x
//                   + 0.5x\log\left( x\sinh(1/x) + \frac{1}{810x^6} \right).
//
// This formula is a good approximation when x > 15.
func logGammaWindschitlApprox(x float64) float64 {
	result := 0.918938533204673 + (x-0.5)*math.Log(x) - x
	result += 0.5 * x * math.Log(x*math.Sinh(1.0/x)+1.0/(810.0*math.Pow(x, 6.0)))
	return result
}
