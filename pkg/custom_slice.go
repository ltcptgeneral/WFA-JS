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
	if actualIdx < len(a.valid) { // idx is in the slice
		return a.valid[actualIdx]
	} else { // idx is out of the slice
		return false
	}
}

func (a *IntegerSlice[T]) Get(idx int) T {
	actualIdx := a.TranslateIndex(idx)
	if actualIdx < len(a.valid) { // idx is in the slice
		return a.data[actualIdx]
	} else { // idx is out of the slice
		return a.defaultValue
	}
}

func (a *IntegerSlice[T]) Set(idx int, value T) {
	actualIdx := a.TranslateIndex(idx)
	if actualIdx >= len(a.valid) { // idx is outside the slice
		// expand data array to actualIdx
		newData := make([]T, actualIdx+1)
		copy(newData, a.data)
		a.data = newData

		// expand valid array to actualIdx
		newValid := make([]bool, actualIdx+1)
		copy(newValid, a.valid)
		a.valid = newValid
	}

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
	if actualIdx >= 0 && actualIdx < len(a.valid) { // idx is in the slice
		return a.valid[actualIdx]
	} else { // idx is out of the slice
		return false
	}
}

func (a *PositiveSlice[T]) Get(idx int) T {
	actualIdx := a.TranslateIndex(idx)
	if actualIdx >= 0 && actualIdx < len(a.valid) { // idx is in the slice
		return a.data[actualIdx]
	} else { // idx is out of the slice
		return a.defaultValue
	}
}

func (a *PositiveSlice[T]) Set(idx int, value T) {
	actualIdx := a.TranslateIndex(idx)
	if actualIdx < 0 || actualIdx >= len(a.valid) { // idx is outside the slice
		// expand data array to actualIdx
		newData := make([]T, actualIdx+1)
		copy(newData, a.data)
		a.data = newData

		// expand valid array to actualIdx
		newValid := make([]bool, actualIdx+1)
		copy(newValid, a.valid)
		a.valid = newValid
	}

	a.data[actualIdx] = value
	a.valid[actualIdx] = true
}
