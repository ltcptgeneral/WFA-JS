package wfa

func WFAlign(s1 string, s2 string, penalties Penalty, doCIGAR bool) Result {
	n := len(s1)
	m := len(s2)
	A_k := m - n
	A_offset := m
	score := 0
	M := NewWavefrontComponent()
	M.SetVal(0, 0, 0)
	M.SetHi(0, 0)
	M.SetLo(0, 0)
	M.SetTraceback(0, 0, End)
	I := NewWavefrontComponent()
	D := NewWavefrontComponent()

	for {
		WFExtend(M, s1, n, s2, m, score)
		ok, val := M.GetVal(score, A_k)
		if ok && val >= A_offset {
			break
		}
		score = score + 1
		WFNext(M, I, D, score, penalties)
	}

	CIGAR := ""
	if doCIGAR {
		CIGAR = WFBacktrace(M, I, D, score, penalties, A_k, s1, s2)
	}

	return Result{
		Score: score,
		CIGAR: CIGAR,
	}
}

func WFExtend(M WavefrontComponent, s1 string, n int, s2 string, m int, score int) {
	_, lo := M.GetLo(score)
	_, hi := M.GetHi(score)
	for k := lo; k <= hi; k++ {
		// v = M[score][k] - k
		// h = M[score][k]
		ok, h := M.GetVal(score, k)
		v := h - k

		// exit early if v or h are invalid
		if !ok {
			continue
		}
		for v < n && h < m && s1[v] == s2[h] {
			_, val := M.GetVal(score, k)
			M.SetVal(score, k, val+1)
			v++
			h++
		}
	}
}

func WFNext(M WavefrontComponent, I WavefrontComponent, D WavefrontComponent, score int, penalties Penalty) {
	// get this score's lo
	lo := NextLo(M, I, D, score, penalties)

	// get this score's hi
	hi := NextHi(M, I, D, score, penalties)

	for k := lo; k <= hi; k++ {
		NextI(M, I, score, k, penalties)
		NextD(M, D, score, k, penalties)
		NextM(M, I, D, score, k, penalties)
	}
}

func WFBacktrace(M WavefrontComponent, I WavefrontComponent, D WavefrontComponent, score int, penalties Penalty, A_k int, s1 string, s2 string) string {
	traceback_CIGAR := []string{"I", "I", "D", "D", "X", "", "", ""}
	x := penalties.X
	o := penalties.O
	e := penalties.E
	CIGAR_rev := ""
	tb_s := score
	tb_k := A_k
	_, current_traceback := M.GetTraceback(tb_s, tb_k)
	done := false

	for !done {
		CIGAR_rev = CIGAR_rev + traceback_CIGAR[current_traceback]
		switch current_traceback {
		case OpenIns:
			tb_s = tb_s - o - e
			tb_k = tb_k - 1
			_, current_traceback = M.GetTraceback(tb_s, tb_k)
		case ExtdIns:
			tb_s = tb_s - e
			tb_k = tb_k - 1
			_, current_traceback = I.GetTraceback(tb_s, tb_k)
		case OpenDel:
			tb_s = tb_s - o - e
			tb_k = tb_k + 1
			_, current_traceback = M.GetTraceback(tb_s, tb_k)
		case ExtdDel:
			tb_s = tb_s - e
			tb_k = tb_k + 1
			_, current_traceback = D.GetTraceback(tb_s, tb_k)
		case Sub:
			tb_s = tb_s - x
			// tb_k = tb_k;
			_, current_traceback = M.GetTraceback(tb_s, tb_k)
		case Ins:
			// tb_s = tb_s;
			// tb_k = tb_k;
			_, current_traceback = I.GetTraceback(tb_s, tb_k)
		case Del:
			// tb_s = tb_s;
			// tb_k = tb_k;
			_, current_traceback = D.GetTraceback(tb_s, tb_k)
		case End:
			done = true
		}
	}

	CIGAR_part := Reverse(CIGAR_rev)
	c := 0
	i := 0
	j := 0
	for i < len(s1) && j < len(s2) {
		if s1[i] == s2[j] {
			//CIGAR_part.splice(c, 0, "M")
			CIGAR_part = Splice(CIGAR_part, 'M', c)
			c++
			i++
			j++
		} else if CIGAR_part[c] == 'X' {
			c++
			i++
			j++
		} else if CIGAR_part[c] == 'I' {
			c++
			j++
		} else if CIGAR_part[c] == 'D' {
			c++
			i++
		}
	}

	return CIGAR_part
}
