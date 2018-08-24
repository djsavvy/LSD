package lsd

// if 'defined' is false, then the value of 'value' is meaningless
type angle struct {
	defined bool
	value   float64
}

/*
Compute the direction of the level-line of 'img' at each pixel

Parameters:
'img': input grid of gray values at each pixel
'threshold': minimum norm for gradient to be considered defined
'bins': number of bins to use in the sorting algorithm

Return:
'angles': angle at each pixel
'norms': norm of gradient at each pixel
'descNorm': a channel returning coordinates for pixels roughly ordered by decreasing
gradient magnitude
*/
func levelLineAngles(img [][]float64, threshold float64, bins uint32) (angles [][]angle, norms [][]float64, descNorm <-chan Point) {

}
