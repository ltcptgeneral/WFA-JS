package wfa

type PositiveSlice[T any] struct {
	data         []T
	valid        []bool
	defaultValue T
}

func (a *PositiveSlice[T]) Valid(idx int) bool {
	return 0 <= idx && idx < len(a.valid) && a.valid[idx]
}

func (a *PositiveSlice[T]) Get(idx int) T {
	if 0 <= idx && idx < len(a.valid) && a.valid[idx] { // idx is in the slice
		return a.data[idx]
	} else { // idx is out of the slice
		return a.defaultValue
	}
}

func (a *PositiveSlice[T]) Set(idx int, value T) {
	if idx >= len(a.valid) { // idx is outside the slice
		// expand data array to 2*idx
		newData := make([]T, 2*idx+1)
		copy(newData, a.data)
		a.data = newData

		// expand valid array to 2*idx
		newValid := make([]bool, 2*idx+1)
		copy(newValid, a.valid)
		a.valid = newValid
	}

	a.data[idx] = value
	a.valid[idx] = true
}
