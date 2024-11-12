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

// bitpacked wavefront lo/hi values with 32 bits each
type WavefrontLoHi uint64

func PackWavefrontLoHi(lo int, hi int) WavefrontLoHi {
	loBM := int64(int32(lo)) & 0x0000_0000_FFFF_FFFF
	hiBM := int64(int64(hi) << 32)
	return WavefrontLoHi(hiBM | loBM)
}

func UnpackWavefrontLoHi(lohi WavefrontLoHi) (int, int) {
	loBM := int(int32(lohi & 0x0000_0000_FFFF_FFFF))
	hiBM := int(int32(lohi & 0xFFFF_FFFF_0000_0000 >> 32))
	return loBM, hiBM
}

// bitpacked wavefront values with 1 valid bit, 3 traceback bits, and 28 bits for the diag distance
// technically this restricts to alignments with less than 268 million characters but that should be sufficient for most cases
type WavefrontValue uint32

// TODO: add 64 bit packed value in case more than 268 million characters are needed

// PackWavefrontValue: packs a diag value and traceback into a WavefrontValue
func PackWavefrontValue(value uint32, traceback Traceback) WavefrontValue {
	validBM := uint32(0x8000_0000)
	tracebackBM := uint32(traceback&0x0000_0007) << 28
	valueBM := value & 0x0FFF_FFFF
	return WavefrontValue(validBM | tracebackBM | valueBM)
}

// UnpackWavefrontValue: opens a WavefrontValue into a valid bool, diag value and traceback
func UnpackWavefrontValue(wfv WavefrontValue) (bool, uint32, Traceback) {
	validBM := wfv&0x8000_0000 != 0
	tracebackBM := uint8(wfv & 0x7000_0000 >> 28)
	valueBM := uint32(wfv & 0x0FFF_FFFF)
	return validBM, valueBM, Traceback(tracebackBM)
}

// Wavefront: stores a single wavefront, stores wavefront's lo value and hi is naturally lo + len(data)
type Wavefront struct { // since wavefronts store diag distance, they should never be negative, and traceback data can be stored as uint8
	data []WavefrontValue
	lohi WavefrontLoHi
}

// NewWavefront: returns a new wavefront with size accomodating lo and hi (inclusive)
func NewWavefront(lo int, hi int) *Wavefront {
	a := &Wavefront{}

	a.lohi = PackWavefrontLoHi(lo, hi)
	size := hi - lo

	newData := make([]WavefrontValue, size+1)
	a.data = newData

	return a
}

// TranslateIndex: utility function for getting the data index given a diagonal
func (a *Wavefront) TranslateIndex(diagonal int) int {
	lo := int(int32(a.lohi & 0x0000_0000_FFFF_FFFF))
	return diagonal - lo
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
	W *PositiveSlice[*Wavefront] // wavefront diag distance and traceback for each wavefront
}

// NewWavefrontComponent: returns initialized WavefrontComponent
func NewWavefrontComponent() *WavefrontComponent {
	// new wavefront component = {
	// lo = [0]
	// hi = [0]
	// W = []
	// }
	w := &WavefrontComponent{
		W: &PositiveSlice[*Wavefront]{
			defaultValue: &Wavefront{
				data: []WavefrontValue{0},
			},
		},
	}

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
	lo, hi := UnpackWavefrontLoHi(w.W.Get(score).lohi)
	return w.W.Valid(score), lo, hi
}

// SetLoHi: set lo and hi for wavefront=score
func (w *WavefrontComponent) SetLoHi(score int, lo int, hi int) {
	b := NewWavefront(lo, hi)
	w.W.Set(score, b)
}
