//go:build debug

package wfa

import (
	"fmt"
	"math"
)

func (w *WavefrontComponent) String(score int) string {
	traceback_str := []string{"OI", "EI", "OD", "ED", "SB", "IN", "DL", "EN"}
	s := "<"
	min_lo := math.MaxInt
	max_hi := math.MinInt

	for i := 0; i <= score; i++ {
		valid := w.W.Valid(i)
		lo, hi := UnpackWavefrontLoHi(w.W.Get(i).lohi)
		if valid && lo < min_lo {
			min_lo = lo
		}
		if valid && hi > max_hi {
			max_hi = hi
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
		lo, hi := UnpackWavefrontLoHi(w.W.Get(i).lohi)
		for k := min_lo; k <= max_hi; k++ {
			valid, val, _ := UnpackWavefrontValue(w.W.Get(i).Get(k))
			if valid {
				s = s + fmt.Sprintf("%02d", val)
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
			valid, _, tb := UnpackWavefrontValue(w.W.Get(i).Get(k))
			if valid {
				s = s + traceback_str[tb]
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
