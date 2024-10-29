package wfa

import (
	"fmt"
	"math"
)

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
	lo *PositiveSlice[int]                      // lo for each wavefront
	hi *PositiveSlice[int]                      // hi for each wavefront
	W  *PositiveSlice[*IntegerSlice[int]]       // wavefront diag distance for each wavefront
	A  *PositiveSlice[*IntegerSlice[traceback]] // compact CIGAR for backtrace for each wavefront
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
		W: &PositiveSlice[*IntegerSlice[int]]{
			defaultValue: &IntegerSlice[int]{
				data:  []int{},
				valid: []bool{},
			},
		},
		A: &PositiveSlice[*IntegerSlice[traceback]]{
			defaultValue: &IntegerSlice[traceback]{
				data:  []traceback{},
				valid: []bool{},
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
	w.A.Set(score, &IntegerSlice[traceback]{})
	w.A.Get(score).Preallocate(lo, hi)

	// preemptively setup w.W
	w.W.Set(score, &IntegerSlice[int]{})
	w.W.Get(score).Preallocate(lo, hi)
}

func (w *WavefrontComponent) String(score int) string {
	traceback_str := []string{"OI", "EI", "OD", "ED", "SB", "IN", "DL", "EN"}
	s := "<"
	min_lo := math.MaxInt
	max_hi := math.MinInt

	for i := 0; i <= score; i++ {
		if w.lo.Valid(i) && w.lo.Get(i) < min_lo {
			min_lo = w.lo.Get(i)
		}
		if w.hi.Valid(i) && w.hi.Get(i) > max_hi {
			max_hi = w.hi.Get(i)
		}
	}

	for k := min_lo; k <= max_hi; k++ {
		s = s + fmt.Sprintf("%02d", k)
		if k < max_hi {
			s = s + "|"
		}
	}

	s = s + ">\t<"

	for k := min_lo; k <= max_hi; k++ {
		s = s + fmt.Sprintf("%02d", k)
		if k < max_hi {
			s = s + "|"
		}
	}

	s = s + ">\n"

	for i := 0; i <= score; i++ {
		s = s + "["
		lo := w.lo.Get(i)
		hi := w.hi.Get(i)
		// print out wavefront matrix
		for k := min_lo; k <= max_hi; k++ {
			if w.W.Valid(i) && w.W.Get(i).Valid(k) {
				s = s + fmt.Sprintf("%02d", w.W.Get(i).Get(k))
			} else if k < lo || k > hi {
				s = s + "--"
			} else {
				s = s + "  "
			}

			if k < max_hi {
				s = s + "|"
			}
		}
		s = s + "]\t["
		// print out traceback matrix
		for k := min_lo; k <= max_hi; k++ {
			if w.A.Valid(i) && w.A.Get(i).Valid(k) {
				s = s + traceback_str[w.A.Get(i).Get(k)]
			} else if k < lo || k > hi {
				s = s + "--"
			} else {
				s = s + "  "
			}

			if k < max_hi {
				s = s + "|"
			}
		}
		s = s + "]\n"
	}
	return s
}
