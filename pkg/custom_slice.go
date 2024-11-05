package wfa

import (
	"golang.org/x/exp/constraints"
)

type Wavefront[T constraints.Integer] struct { // since wavefronts store diag distance, they should never be negative, and traceback data can be stored as uint8
	data  []T
	valid []bool
	lo    int
}

func NewWavefront[T constraints.Integer](lo int, hi int) *Wavefront[T] {
	a := &Wavefront[T]{}

	a.lo = lo
	size := a.TranslateIndex(hi)

	newData := make([]T, size+1)
	a.data = newData

	newValid := make([]bool, size+1)
	a.valid = newValid

	return a
}

func (a *Wavefront[T]) TranslateIndex(idx int) int {
	return idx - a.lo
}

func (a *Wavefront[T]) Valid(idx int) bool {
	actualIdx := a.TranslateIndex(idx)
	return 0 <= actualIdx && actualIdx < len(a.data) && a.valid[actualIdx]
}

func (a *Wavefront[T]) Get(idx int) T {
	actualIdx := a.TranslateIndex(idx)
	if 0 <= actualIdx && actualIdx < len(a.data) { // idx is in the slice
		return a.data[actualIdx]
	} else { // idx is out of the slice
		return 0
	}
}

func (a *Wavefront[T]) Set(idx int, value T) {
	actualIdx := a.TranslateIndex(idx)

	/* in theory idx is always in bounds because the wavefront is preallocated
	if actualIdx < 0 || actualIdx >= len(a.data) {
		return
	}
	*/

	a.data[actualIdx] = value
	a.valid[actualIdx] = true
}

type PositiveSlice[T any] struct {
	data         []T
	valid        []bool
	defaultValue T
}

func (a *PositiveSlice[T]) TranslateIndex(idx int) int {
	return idx
}

func (a *PositiveSlice[T]) Valid(idx int) bool {
	actualIdx := a.TranslateIndex(idx)
	return 0 <= actualIdx && actualIdx < len(a.valid) && a.valid[actualIdx]
}

func (a *PositiveSlice[T]) Get(idx int) T {
	actualIdx := a.TranslateIndex(idx)
	if 0 <= actualIdx && actualIdx < len(a.valid) && a.valid[actualIdx] { // idx is in the slice
		return a.data[actualIdx]
	} else { // idx is out of the slice
		return a.defaultValue
	}
}

func (a *PositiveSlice[T]) Set(idx int, value T) {
	actualIdx := a.TranslateIndex(idx)
	if actualIdx < 0 || actualIdx >= len(a.valid) { // idx is outside the slice
		// expand data array to actualIdx
		newData := make([]T, 2*actualIdx+1)
		copy(newData, a.data)
		a.data = newData

		// expand valid array to actualIdx
		newValid := make([]bool, 2*actualIdx+1)
		copy(newValid, a.valid)
		a.valid = newValid
	}

	a.data[actualIdx] = value
	a.valid[actualIdx] = true
}

func (a *PositiveSlice[T]) Preallocate(hi int) {
	size := hi

	// expand data array to actualIdx
	newData := make([]T, size+1)
	a.data = newData

	// expand valid array to actualIdx
	newValid := make([]bool, size+1)
	a.valid = newValid
}
