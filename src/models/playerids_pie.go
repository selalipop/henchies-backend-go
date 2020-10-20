package models

import (
	"github.com/elliotchance/pie/pie/util"
	"math/rand"
)

// Contains returns true if the element exists in the slice.
//
// When using slices of pointers it will only compare by address, not value.
func (ss PlayerIds) Contains(lookingFor PlayerId) bool {
	for _, s := range ss {
		if lookingFor == s {
			return true
		}
	}

	return false
}

// Filter will return a new slice containing only the elements that return
// true from the condition. The returned slice may contain zero elements (nil).
//
// FilterNot works in the opposite way of Filter.
func (ss PlayerIds) Filter(condition func(PlayerId) bool) (ss2 PlayerIds) {
	for _, s := range ss {
		if condition(s) {
			ss2 = append(ss2, s)
		}
	}
	return
}

// Shuffle returns shuffled slice by your rand.Source
func (ss PlayerIds) Shuffle(source rand.Source) PlayerIds {
	n := len(ss)

	// Avoid the extra allocation.
	if n < 2 {
		return ss
	}

	// go 1.10+ provides rnd.Shuffle. However, to support older versions we copy
	// the algorithm directly from the go source: src/math/rand/rand.go below,
	// with some adjustments:
	shuffled := make([]PlayerId, n)
	copy(shuffled, ss)

	rnd := rand.New(source)

	util.Shuffle(rnd, n, func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	return shuffled
}
