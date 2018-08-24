package lsd

// Take in pixels and their corresponding gradient values, and output
// the pixels roughly in decreasing order by gradient magnitude, up to a precision given by 'maxGrad'/'numBins'.
// This operation is performed in linear time by splitting the pixels into bins.
func sortPixels(grads [][]gradient, maxGrad float64, numBins uint32) <-chan pixel {
	sorted := make(chan pixel, len(grads)*len(grads[0]))

	// Create intermediate bins for sorting
	bins := make([][]pixel, numBins)
	for i := range bins {
		bins[i] = make([]pixel, 0)
	}

	// Note that zero-norm points are discarded down the line in the
	// original LSD algorithm; to save memory, we simply discard them here.
	// Sort the pixels into bins
	for x := range grads {
		for y := range grads[x] {
			if grads[x][y].norm == 0 {
				continue
			}
			bin := uint32(grads[x][y].norm * float64(numBins) / maxGrad)
			if bin >= uint32(numBins) {
				bin = numBins - 1
			}
			bins[bin] = append(bins[bin], pixel{x, y})
		}
	}

	go func() {
		for i := numBins - 1; i >= 0; i-- {
			for _, p := range bins[i] {
				sorted <- p
			}
		}
		close(sorted)
	}()
	return sorted
}
