package wfa

type IntegerSlice[T any] struct {
	data         []T
	valid        []bool
	defaultValue T
}

func (a *IntegerSlice[T]) TranslateIndex(idx int) int {
	if idx >= 0 { // 0 -> 0, 1 -> 2, 2 -> 4, 3 -> 6, ...
		return 2 * idx
	} else { // -1 -> 1, -2 -> 3, -3 -> 5, ...
		return (-2 * idx) - 1
	}
}

func (a *IntegerSlice[T]) Valid(idx int) bool {
	actualIdx := a.TranslateIndex(idx)
	return 0 <= actualIdx && actualIdx < len(a.valid) && a.valid[actualIdx]
}

func (a *IntegerSlice[T]) Get(idx int) T {
	actualIdx := a.TranslateIndex(idx)
	if 0 <= actualIdx && actualIdx < len(a.valid) && a.valid[actualIdx] { // idx is in the slice
		return a.data[actualIdx]
	} else { // idx is out of the slice
		return a.defaultValue
	}
}

func (a *IntegerSlice[T]) Set(idx int, value T) {
	actualIdx := a.TranslateIndex(idx)
	if actualIdx >= len(a.valid) { // idx is outside the slice
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

func (a *IntegerSlice[T]) Preallocate(lo int, hi int) {
	actualLo := a.TranslateIndex(lo)
	actualHi := a.TranslateIndex(hi)
	size := max(actualHi, actualLo)

	// expand data array to actualIdx
	newData := make([]T, size+1)
	a.data = newData

	// expand valid array to actualIdx
	newValid := make([]bool, size+1)
	a.valid = newValid
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
