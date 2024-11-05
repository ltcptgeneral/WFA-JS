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

type Traceback byte

const (
	OpenIns Traceback = iota
	ExtdIns
	OpenDel
	ExtdDel
	Sub
	Ins
	Del
	End
)

// bitpacked wavefront values with 1 valid bit, 3 traceback bits, and 28 bits for the diag distance
// technically this restricts to solutions within 268 million score but that should be sufficient for most cases
type WavefrontValue uint32

// TODO: add 64 bit packed value in case more than 268 million score is needed

// PackWavefrontValue: packs a diag value and traceback into a WavefrontValue
func PackWavefrontValue(value uint32, traceback Traceback) WavefrontValue {
	valueBM := value & 0x0FFF_FFFF
	tracebackBM := uint32(traceback&0x0000_0007) << 28
	return WavefrontValue(0x8000_0000 | valueBM | tracebackBM)
}

// UnpackWavefrontValue: opens a WavefrontValue into a valid bool, diag value and traceback
func UnpackWavefrontValue(wf WavefrontValue) (bool, uint32, Traceback) {
	valueBM := uint32(wf & 0x0FFF_FFFF)
	tracebackBM := uint8(wf & 0x7000_0000 >> 28)
	validBM := wf&0x8000_0000 != 0
	return validBM, valueBM, Traceback(tracebackBM)
}

// Wavefront: stores a single wavefront, stores wavefront's lo value and hi is naturally lo + len(data)
type Wavefront struct { // since wavefronts store diag distance, they should never be negative, and traceback data can be stored as uint8
	data []WavefrontValue
	lo   int
}

// NewWavefront: returns a new wavefront with size accomodating lo and hi (inclusive)
func NewWavefront(lo int, hi int) *Wavefront {
	a := &Wavefront{}

	a.lo = lo
	size := a.TranslateIndex(hi)

	newData := make([]WavefrontValue, size+1)
	a.data = newData

	return a
}

// TranslateIndex: utility function for getting the data index given a diagonal
func (a *Wavefront) TranslateIndex(diagonal int) int {
	return diagonal - a.lo
}

// Get: returns WavefrontValue for given diagonal
func (a *Wavefront) Get(diagonal int) WavefrontValue {
	actualIdx := a.TranslateIndex(diagonal)
	if 0 <= actualIdx && actualIdx < len(a.data) { // idx is in the slice
		return a.data[actualIdx]
	} else { // idx is out of the slice
		return 0
	}
}

// Set: the diagonal to a WavefrontValue
func (a *Wavefront) Set(diagonal int, value WavefrontValue) {
	actualIdx := a.TranslateIndex(diagonal)

	/* in theory idx is always in bounds because the wavefront is preallocated
	if actualIdx < 0 || actualIdx >= len(a.data) {
		return
	}
	*/

	a.data[actualIdx] = value
}

// WavefrontComponent: each M/I/D wavefront matrix including the wavefront data, lo and hi
type WavefrontComponent struct {
	lo *PositiveSlice[int]        // lo for each wavefront
	hi *PositiveSlice[int]        // hi for each wavefront
	W  *PositiveSlice[*Wavefront] // wavefront diag distance and traceback for each wavefront
}

// NewWavefrontComponent: returns initialized WavefrontComponent
func NewWavefrontComponent(preallocateSize int) WavefrontComponent {
	// new wavefront component = {
	// lo = [0]
	// hi = [0]
	// W = []
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
		W: &PositiveSlice[*Wavefront]{
			defaultValue: &Wavefront{
				data: []WavefrontValue{0},
			},
		},
	}

	w.lo.Preallocate(preallocateSize)
	w.hi.Preallocate(preallocateSize)
	w.W.Preallocate(preallocateSize)

	return w
}

// GetVal: get value for wavefront=score, diag=k => returns ok, value, traceback
func (w *WavefrontComponent) GetVal(score int, k int) (bool, uint32, Traceback) {
	return UnpackWavefrontValue(w.W.Get(score).Get(k))
}

// SetVal: set value, traceback for wavefront=score, diag=k
func (w *WavefrontComponent) SetVal(score int, k int, val uint32, tb Traceback) {
	w.W.Get(score).Set(k, PackWavefrontValue(val, tb))
}

// GetLoHi: get lo and hi for wavefront=score
func (w *WavefrontComponent) GetLoHi(score int) (bool, int, int) {
	// if lo[score] and hi[score] are valid
	if w.lo.Valid(score) && w.hi.Valid(score) {
		// return lo[score] hi[score]
		return true, w.lo.Get(score), w.hi.Get(score)
	} else {
		return false, 0, 0
	}
}

// SetLoHi: set lo and hi for wavefront=score
func (w *WavefrontComponent) SetLoHi(score int, lo int, hi int) {
	// lo[score] = lo
	w.lo.Set(score, lo)
	// hi[score] = hi
	w.hi.Set(score, hi)

	// preemptively setup w.W
	b := NewWavefront(lo, hi)
	w.W.Set(score, b)
}
