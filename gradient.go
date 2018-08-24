package lsd

// if 'norm' is zero, then the value of 'angle' is meaningless
type gradient struct {
	norm  bool
	angle float64
}

/*
Compute the magnitude and direction of the level-line of 'img' at each pixel

Parameters:
'img': input grid of gray values at each pixel
'threshold': minimum norm for gradient to be considered defined
'bins': number of bins to use in the sorting algorithm

Return:
'grads': gradient at each pixel
'descNorm': a channel returning coordinates for pixels roughly ordered by decreasing
gradient magnitude
*/
func computeGradients(img [][]float64, threshold float64, bins uint32) (grads [][]gradient, descNorm <-chan Point) {

}
