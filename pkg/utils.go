package wfa

import (
	"math"
	"strings"

	"golang.org/x/exp/constraints"
)

// convert an unsigned into to string
func UIntToString(num uint) string { // num assumed to be positive
	var builder strings.Builder

	for num > 0 {
		digit := num % 10
		builder.WriteRune(rune('0' + digit))
		num /= 10
	}

	// Reverse the string as we built it in reverse order
	str := []rune(builder.String())
	for i, j := 0, len(str)-1; i < j; i, j = i+1, j-1 {
		str[i], str[j] = str[j], str[i]
	}

	return string(str)
}

// decode runlength encoded string such as CIGARs
func RunLengthDecode(encoded string) string {
	decoded := strings.Builder{}
	length := len(encoded)
	i := 0

	for i < length {
		// If the current character is a digit, we need to extract the run length
		runLength := 0
		for i < length && encoded[i] >= '0' && encoded[i] <= '9' {
			runLength = runLength*10 + int(encoded[i]-'0')
			i++
		}

		// The next character will be the character to repeat
		if i < length {
			char := encoded[i]
			for j := 0; j < runLength; j++ {
				decoded.WriteByte(char)
			}
			i++ // Move past the character
		}
	}

	return decoded.String()
}

// given the min index, return the item in values at that index
func SafeMin[T constraints.Integer](values []T, idx int) T {
	return values[idx]
}

// given the max index, return the item in values at that index
func SafeMax[T constraints.Integer](values []T, idx int) T {
	return values[idx]
}

// given array of values and corresponding array of valid flags, find the min of value which is valid or return false if there does not exist any
func SafeArgMin[T constraints.Integer](valids []bool, values []T) (bool, int) {
	hasValid := false
	minIndex := 0
	minValue := math.MaxInt
	for i := 0; i < len(valids); i++ {
		if valids[i] && int(values[i]) < minValue {
			hasValid = true
			minIndex = i
			minValue = int(values[i])
		}
	}
	if hasValid {
		return true, minIndex
	} else {
		return false, 0
	}
}

// given array of values and corresponding array of valid flags, find the max of value which is valid or return false if there does not exist any
func SafeArgMax[T constraints.Integer](valids []bool, values []T) (bool, int) {
	hasValid := false
	maxIndex := 0
	maxValue := math.MinInt
	for i := range valids {
		if valids[i] && int(values[i]) > maxValue {
			hasValid = true
			maxIndex = i
			maxValue = int(values[i])
		}
	}
	return hasValid, maxIndex
}

// set the lext lo and hi bounds for wavefronts M, I, D
func NextLoHi(M *WavefrontComponent, I *WavefrontComponent, D *WavefrontComponent, score int, penalties Penalty) (int, int) {
	x := penalties.X
	o := penalties.O
	e := penalties.E

	a_ok, a_lo, a_hi := M.GetLoHi(score - x)
	b_ok, b_lo, b_hi := M.GetLoHi(score - o - e)
	c_ok, c_lo, c_hi := I.GetLoHi(score - e)
	d_ok, d_lo, d_hi := D.GetLoHi(score - e)

	ok_lo, idx := SafeArgMin(
		[]bool{a_ok, b_ok, c_ok, d_ok},
		[]int{a_lo, b_lo, c_lo, d_lo},
	)
	lo := SafeMin([]int{a_lo, b_lo, c_lo, d_lo}, idx) - 1

	ok_hi, idx := SafeArgMax(
		[]bool{a_ok, b_ok, c_ok, d_ok},
		[]int{a_hi, b_hi, c_hi, d_hi},
	)
	hi := SafeMax([]int{a_hi, b_hi, c_hi, d_hi}, idx) + 1

	if ok_lo && ok_hi {
		M.SetLoHi(score, lo, hi)
		I.SetLoHi(score, lo, hi)
		D.SetLoHi(score, lo, hi)
	}
	return lo, hi
}

// set the traceback and diag value for the next I wavefront
func NextI(M *WavefrontComponent, I *WavefrontComponent, score int, k int, penalties Penalty) {
	o := penalties.O
	e := penalties.E

	a_ok, a, _ := M.GetVal(score-o-e, k-1)
	b_ok, b, _ := I.GetVal(score-e, k-1)

	ok, nextITraceback := SafeArgMax([]bool{a_ok, b_ok}, []uint64{a, b})
	nextIVal := SafeMax([]uint64{a, b}, nextITraceback) + 1 // important that the +1 is here
	if ok {
		I.SetVal(score, k, nextIVal, []Traceback{OpenIns, ExtdIns}[nextITraceback])
	}
}

// set the traceback and diag value for the next D wavefront
func NextD(M *WavefrontComponent, D *WavefrontComponent, score int, k int, penalties Penalty) {
	o := penalties.O
	e := penalties.E

	a_ok, a, _ := M.GetVal(score-o-e, k+1)
	b_ok, b, _ := D.GetVal(score-e, k+1)

	ok, nextDTraceback := SafeArgMax([]bool{a_ok, b_ok}, []uint64{a, b})
	nextDVal := SafeMax([]uint64{a, b}, nextDTraceback)
	if ok {
		D.SetVal(score, k, nextDVal, []Traceback{OpenDel, ExtdDel}[nextDTraceback])
	}
}

// set the traceback and diag value for the next M wavefront
func NextM(M *WavefrontComponent, I *WavefrontComponent, D *WavefrontComponent, score int, k int, penalties Penalty) {
	x := penalties.X

	a_ok, a, _ := M.GetVal(score-x, k)
	a++ // important to have +1 here
	b_ok, b, _ := I.GetVal(score, k)
	c_ok, c, _ := D.GetVal(score, k)

	ok, nextMTraceback := SafeArgMax([]bool{a_ok, b_ok, c_ok}, []uint64{a, b, c})
	nextMVal := SafeMax([]uint64{a, b, c}, nextMTraceback)
	if ok {
		M.SetVal(score, k, nextMVal, []Traceback{Sub, Ins, Del}[nextMTraceback])
	}
}
