package wfa

type Result struct {
	Score int
	CIGAR string
}

type Penalty struct {
	M int
	X int
	O int
	E int
}

type traceback byte

const (
	OpenIns traceback = iota
	ExtdIns
	OpenDel
	ExtdDel
	Sub
	Ins
	Del
	End
)

type WavefrontComponent struct {
	lo *PositiveSlice[int]                   // lo for each wavefront
	hi *PositiveSlice[int]                   // hi for each wavefront
	W  *PositiveSlice[*Wavefront[int]]       // wavefront diag distance for each wavefront
	A  *PositiveSlice[*Wavefront[traceback]] // compact CIGAR for backtrace for each wavefront
}

func NewWavefrontComponent(preallocateSize int) WavefrontComponent {
	// new wavefront component = {
	// lo = [0]
	// hi = [0]
	// W = []
	// A = []
	// }
	w := WavefrontComponent{
		lo: &PositiveSlice[int]{
			data:  []int{0},
			valid: []bool{true},
		},
		hi: &PositiveSlice[int]{
			data:  []int{0},
			valid: []bool{true},
		},
		W: &PositiveSlice[*Wavefront[int]]{
			defaultValue: &Wavefront[int]{
				data:  []int{0},
				valid: []bool{false},
			},
		},
		A: &PositiveSlice[*Wavefront[traceback]]{
			defaultValue: &Wavefront[traceback]{
				data:  []traceback{0},
				valid: []bool{false},
			},
		},
	}

	w.lo.Preallocate(preallocateSize)
	w.hi.Preallocate(preallocateSize)
	w.W.Preallocate(preallocateSize)
	w.A.Preallocate(preallocateSize)

	return w
}

// get value for wavefront=score, diag=k => returns ok, value
func (w *WavefrontComponent) GetVal(score int, k int) (bool, int) {
	return w.W.Valid(score) && w.W.Get(score).Valid(k), w.W.Get(score).Get(k)
}

// set value for wavefront=score, diag=k
func (w *WavefrontComponent) SetVal(score int, k int, val int) {
	w.W.Get(score).Set(k, val)
}

// get alignment traceback for wavefront=score, diag=k => returns ok, value
func (w *WavefrontComponent) GetTraceback(score int, k int) (bool, traceback) {
	return w.A.Valid(score) && w.A.Get(score).Valid(k), w.A.Get(score).Get(k)
}

// set alignment traceback for wavefront=score, diag=k
func (w *WavefrontComponent) SetTraceback(score int, k int, val traceback) {
	w.A.Get(score).Set(k, val)
}

// get hi for wavefront=score
func (w *WavefrontComponent) GetLoHi(score int) (bool, int, int) {
	// if lo[score] and hi[score] are valid
	if w.lo.Valid(score) && w.hi.Valid(score) {
		// return lo[score] hi[score]
		return true, w.lo.Get(score), w.hi.Get(score)
	} else {
		return false, 0, 0
	}
}

// set hi for wavefront=score
func (w *WavefrontComponent) SetLoHi(score int, lo int, hi int) {
	// lo[score] = lo
	w.lo.Set(score, lo)
	// hi[score] = hi
	w.hi.Set(score, hi)

	// preemptively setup w.A
	a := NewWavefront[traceback](lo, hi)
	w.A.Set(score, a)

	// preemptively setup w.W
	b := NewWavefront[int](lo, hi)
	w.W.Set(score, b)
}
