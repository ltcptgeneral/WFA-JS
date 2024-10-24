package wfa

import (
	"math"
	"unicode/utf8"
)

func SafeMin(valids []bool, values []int) (bool, int) {
	ok, idx := SafeArgMin(valids, values)
	return ok, values[idx]
}

func SafeMax(valids []bool, values []int) (bool, int) {
	ok, idx := SafeArgMax(valids, values)
	return ok, values[idx]
}

func SafeArgMax(valids []bool, values []int) (bool, int) {
	hasValid := false
	maxIndex := 0
	maxValue := math.MinInt
	for i := 0; i < len(valids); i++ {
		if valids[i] && values[i] > maxValue {
			hasValid = true
			maxIndex = i
			maxValue = values[i]
		}
	}
	if hasValid {
		return true, maxIndex
	} else {
		return false, 0
	}
}

func SafeArgMin(valids []bool, values []int) (bool, int) {
	hasValid := false
	minIndex := 0
	minValue := math.MaxInt
	for i := 0; i < len(valids); i++ {
		if valids[i] && values[i] < minValue {
			hasValid = true
			minIndex = i
			minValue = values[i]
		}
	}
	if hasValid {
		return true, minIndex
	} else {
		return false, 0
	}
}

func Reverse(s string) string {
	size := len(s)
	buf := make([]byte, size)
	for start := 0; start < size; {
		r, n := utf8.DecodeRuneInString(s[start:])
		start += n
		utf8.EncodeRune(buf[size-start:], r)
	}
	return string(buf)
}

func Splice(s string, c rune, idx int) string {
	return s[:idx] + string(c) + s[idx:]
}

func NextLo(M WavefrontComponent, I WavefrontComponent, D WavefrontComponent, score int, penalties Penalty) int {
	x := penalties.X
	o := penalties.O
	e := penalties.E

	a_ok, a := M.GetLo(score - x)
	b_ok, b := M.GetLo(score - o - e)
	c_ok, c := I.GetLo(score - e)
	d_ok, d := D.GetLo(score - e)

	ok, lo := SafeMin(
		[]bool{a_ok, b_ok, c_ok, d_ok},
		[]int{a, b, c, d},
	)
	lo--
	if ok {
		M.SetLo(score, lo)
		I.SetLo(score, lo)
		D.SetLo(score, lo)
	}
	return lo
}

func NextHi(M WavefrontComponent, I WavefrontComponent, D WavefrontComponent, score int, penalties Penalty) int {
	x := penalties.X
	o := penalties.O
	e := penalties.E

	a_ok, a := M.GetHi(score - x)
	b_ok, b := M.GetHi(score - o - e)
	c_ok, c := I.GetHi(score - e)
	d_ok, d := D.GetHi(score - e)

	ok, hi := SafeMax(
		[]bool{a_ok, b_ok, c_ok, d_ok},
		[]int{a, b, c, d},
	)
	hi++
	if ok {
		M.SetHi(score, hi)
		I.SetHi(score, hi)
		D.SetHi(score, hi)
	}
	return hi
}

func NextI(M WavefrontComponent, I WavefrontComponent, score int, k int, penalties Penalty) {
	o := penalties.O
	e := penalties.E

	a_ok, a := M.GetVal(score-o-e, k-1)
	b_ok, b := I.GetVal(score-e, k-1)

	ok, nextIVal := SafeMax([]bool{a_ok, b_ok}, []int{a, b})
	if ok {
		I.SetVal(score, k, nextIVal+1) // important that the +1 is here
	}

	ok, nextITraceback := SafeArgMax([]bool{a_ok, b_ok}, []int{a, b})
	if ok {
		I.SetTraceback(score, k, []traceback{OpenIns, ExtdIns}[nextITraceback])
	}
}

func NextD(M WavefrontComponent, D WavefrontComponent, score int, k int, penalties Penalty) {
	o := penalties.O
	e := penalties.E

	a_ok, a := M.GetVal(score-o-e, k+1)
	b_ok, b := D.GetVal(score-e, k+1)

	ok, nextDVal := SafeMax([]bool{a_ok, b_ok}, []int{a, b})
	if ok {
		D.SetVal(score, k, nextDVal) // nothing special
	}

	ok, nextDTraceback := SafeArgMax([]bool{a_ok, b_ok}, []int{a, b})
	if ok {
		D.SetTraceback(score, k, []traceback{OpenDel, ExtdDel}[nextDTraceback])
	}
}

func NextM(M WavefrontComponent, I WavefrontComponent, D WavefrontComponent, score int, k int, penalties Penalty) {
	x := penalties.X

	a_ok, a := M.GetVal(score-x, k)
	a++ // important to have +1 here
	b_ok, b := I.GetVal(score, k)
	c_ok, c := D.GetVal(score, k)

	ok, nextMVal := SafeMax([]bool{a_ok, b_ok, c_ok}, []int{a, b, c})
	if ok {
		M.SetVal(score, k, nextMVal)
	}

	ok, nextMTraceback := SafeArgMax([]bool{a_ok, b_ok, c_ok}, []int{a, b, c})
	if ok {
		M.SetTraceback(score, k, []traceback{Sub, Ins, Del}[nextMTraceback])
	}
}
