package models

import (
	"github.com/elliotchance/pie/pie/util"
	"math/rand"
)

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

// Shuffle returns shuffled slice by your rand.Source
func (ss PlayerColors) Shuffle(source rand.Source) PlayerColors {
	n := len(ss)

	// Avoid the extra allocation.
	if n < 2 {
		return ss
	}

	// go 1.10+ provides rnd.Shuffle. However, to support older versions we copy
	// the algorithm directly from the go source: src/math/rand/rand.go below,
	// with some adjustments:
	shuffled := make([]PlayerColor, n)
	copy(shuffled, ss)

	rnd := rand.New(source)

	util.Shuffle(rnd, n, func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}
