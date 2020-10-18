package models

// DropTop will return the rest slice after dropping the top n elements
// if the slice has less elements then n that'll return empty slice
// if n < 0 it'll return empty slice.
func (ss PlayerColors) DropTop(n int) (drop PlayerColors) {
	if n < 0 || n >= len(ss) {
		return
	}

	// Copy ss, to make sure no memory is overlapping between input and
	// output. See issue #145.
	drop = make([]PlayerColor, len(ss)-n)
	copy(drop, ss[n:])

	return
}
