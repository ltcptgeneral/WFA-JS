package wfa

import (
	"strings"
)

func WFAlign(s1 string, s2 string, penalties Penalty, doCIGAR bool) Result {
	n := len(s1)
	m := len(s2)
	A_k := m - n
	A_offset := uint32(m)
	score := 0
	estimatedScore := (max(n, m) * max(penalties.M, penalties.X, penalties.O, penalties.E)) / 4
	M := NewWavefrontComponent(estimatedScore)
	M.SetLoHi(0, 0, 0)
	M.SetVal(0, 0, 0, End)
	I := NewWavefrontComponent(estimatedScore)
	D := NewWavefrontComponent(estimatedScore)

	for {
		WFExtend(M, s1, n, s2, m, score)
		ok, val, _ := M.GetVal(score, A_k)
		if ok && val >= A_offset {
			break
		}
		score = score + 1
		WFNext(M, I, D, score, penalties)
	}

	CIGAR := ""
	if doCIGAR {
		CIGAR = WFBacktrace(M, I, D, score, penalties, A_k, A_offset, s1, s2)
	}

	return Result{
		Score: score,
		CIGAR: CIGAR,
	}
}

func WFExtend(M *WavefrontComponent, s1 string, n int, s2 string, m int, score int) {
	_, lo, hi := M.GetLoHi(score)
	for k := lo; k <= hi; k++ {
		// v = M[score][k] - k
		// h = M[score][k]
		ok, hu, _ := M.GetVal(score, k)
		h := int(hu)
		v := h - k

		// exit early if v or h are invalid
		if !ok {
			continue
		}
		for v < n && h < m && s1[v] == s2[h] {
			_, val, tb := M.GetVal(score, k)
			M.SetVal(score, k, val+1, tb)
			v++
			h++
		}
	}
}

func WFNext(M *WavefrontComponent, I *WavefrontComponent, D *WavefrontComponent, score int, penalties Penalty) {
	// get this score's lo, hi
	lo, hi := NextLoHi(M, I, D, score, penalties)

	for k := lo; k <= hi; k++ {
		NextI(M, I, score, k, penalties)
		NextD(M, D, score, k, penalties)
		NextM(M, I, D, score, k, penalties)
	}
}

func WFBacktrace(M *WavefrontComponent, I *WavefrontComponent, D *WavefrontComponent, score int, penalties Penalty, A_k int, A_offset uint32, s1 string, s2 string) string {
	x := penalties.X
	o := penalties.O
	e := penalties.E

	tb_s := score
	tb_k := A_k
	done := false

	_, current_dist, current_traceback := M.GetVal(tb_s, tb_k)

	Ops := []rune{'~'}
	Counts := []uint{0}
	idx := 0

	for !done {
		switch current_traceback {
		case OpenIns:
			if Ops[idx] == 'I' {
				Counts[idx]++
			} else {
				Ops = append(Ops, 'I')
				Counts = append(Counts, 1)
				idx++
			}

			tb_s = tb_s - o - e
			tb_k = tb_k - 1
			_, current_dist, current_traceback = M.GetVal(tb_s, tb_k)
		case ExtdIns:
			if Ops[idx] == 'I' {
				Counts[idx]++
			} else {
				Ops = append(Ops, 'I')
				Counts = append(Counts, 1)
				idx++
			}

			tb_s = tb_s - e
			tb_k = tb_k - 1
			_, current_dist, current_traceback = I.GetVal(tb_s, tb_k)
		case OpenDel:
			if Ops[idx] == 'D' {
				Counts[idx]++
			} else {
				Ops = append(Ops, 'D')
				Counts = append(Counts, 1)
				idx++
			}

			tb_s = tb_s - o - e
			tb_k = tb_k + 1
			_, current_dist, current_traceback = M.GetVal(tb_s, tb_k)
		case ExtdDel:
			if Ops[idx] == 'D' {
				Counts[idx]++
			} else {
				Ops = append(Ops, 'D')
				Counts = append(Counts, 1)
				idx++
			}

			tb_s = tb_s - e
			tb_k = tb_k + 1
			_, current_dist, current_traceback = D.GetVal(tb_s, tb_k)
		case Sub:
			tb_s = tb_s - x
			// tb_k = tb_k;
			_, next_dist, next_traceback := M.GetVal(tb_s, tb_k)

			if int(current_dist-next_dist)-1 > 0 {
				Ops = append(Ops, 'M')
				Counts = append(Counts, uint(current_dist-next_dist)-1)
				idx++
			}

			if Ops[idx] == 'X' {
				Counts[idx]++
			} else {
				Ops = append(Ops, 'X')
				Counts = append(Counts, 1)
				idx++
			}

			current_dist = next_dist
			current_traceback = next_traceback
		case Ins:
			// tb_s = tb_s;
			// tb_k = tb_k;
			_, next_dist, next_traceback := I.GetVal(tb_s, tb_k)

			Ops = append(Ops, 'M')
			Counts = append(Counts, uint(current_dist-next_dist))
			idx++

			current_dist = next_dist
			current_traceback = next_traceback
		case Del:
			// tb_s = tb_s;
			// tb_k = tb_k;
			_, next_dist, next_traceback := D.GetVal(tb_s, tb_k)

			Ops = append(Ops, 'M')
			Counts = append(Counts, uint(current_dist-next_dist))
			idx++

			current_dist = next_dist
			current_traceback = next_traceback
		case End:
			Ops = append(Ops, 'M')
			Counts = append(Counts, uint(current_dist))
			idx++

			done = true
		}
	}

	CIGAR := strings.Builder{}
	for i := len(Ops) - 1; i > 0; i-- {
		CIGAR.WriteString(UIntToString(Counts[i]))
		CIGAR.WriteRune(Ops[i])
	}

	return CIGAR.String()
}
