package wfa

// WFAlign takes strings s1, s2, penalties, and returns the score and CIGAR if doCIGAR is true
func WFAlign(s1 string, s2 string, penalties Penalty, doCIGAR bool) Result {
	n := len(s1)
	m := len(s2)
	A_k := m - n          // diagonal where both sequences end
	A_offset := uint64(m) // offset along a_k diagonal corresponding to end
	score := 0
	M := NewWavefrontComponent()
	M.SetLoHi(0, 0, 0)
	M.SetVal(0, 0, 0, End)
	I := NewWavefrontComponent()
	D := NewWavefrontComponent()

	for {
		WFExtend(M, s1, n, s2, m, score)
		ok, val, _ := M.GetVal(score, A_k)
		if ok && val >= A_offset { // exit when M_(s,a_k) >= A_offset, ie the wavefront has reached the end
			break
		}
		score = score + 1
		WFNext(M, I, D, score, penalties)
	}

	CIGAR := ""
	if doCIGAR { // if doCIGAR, then perform backtrace, otherwise just return the score
		CIGAR = WFBacktrace(M, I, D, score, penalties, A_k, A_offset, s1, s2)
	}

	return Result{
		Score: score,
		CIGAR: CIGAR,
	}
}

func WFExtend(M *WavefrontComponent, s1 string, n int, s2 string, m int, score int) {
	_, lo, hi := M.GetLoHi(score)
	for k := lo; k <= hi; k++ { // for each diagonal in current wavefront
		// v = M[score][k] - k
		// h = M[score][k]
		ok, uh, tb := M.GetVal(score, k)
		// exit early if M_(s,l) is invalid
		if !ok {
			continue
		}
		h := int(uh)
		v := h - k
		// in the paper, we do v++, h++, M_(s,k)++
		// however, note that h = M_(s,k) so instead we just do v++, h++ and set M_(s,k) at the end
		// this saves a some memory reads and writes
		for v < n && h < m && s1[v] == s2[h] { // extend diagonal for the next set of matches
			v++
			h++
		}
		M.SetVal(score, k, uint64(h), tb)
	}
}

func WFNext(M *WavefrontComponent, I *WavefrontComponent, D *WavefrontComponent, score int, penalties Penalty) {
	// get this score's lo, hi
	lo, hi := NextLoHi(M, I, D, score, penalties)

	for k := lo; k <= hi; k++ { // for each diagonal, extend the matrices for the next wavefronts
		NextI(M, I, score, k, penalties)
		NextD(M, D, score, k, penalties)
		NextM(M, I, D, score, k, penalties)
	}
}

func WFBacktrace(M *WavefrontComponent, I *WavefrontComponent, D *WavefrontComponent, score int, penalties Penalty, A_k int, A_offset uint64, s1 string, s2 string) string {
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

	CIGAR := ""
	for i := len(Ops) - 1; i > 0; i-- {
		CIGAR += UIntToString(Counts[i])
		CIGAR += string(Ops[i])
	}

	return CIGAR
}
