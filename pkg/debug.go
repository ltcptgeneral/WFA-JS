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
