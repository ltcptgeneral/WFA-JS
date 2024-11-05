package wfa

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
		CIGAR = WFBacktrace(M, I, D, score, penalties, A_k, s1, s2)
	}

	return Result{
		Score: score,
		CIGAR: CIGAR,
	}
}

func WFExtend(M WavefrontComponent, s1 string, n int, s2 string, m int, score int) {
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

func WFNext(M WavefrontComponent, I WavefrontComponent, D WavefrontComponent, score int, penalties Penalty) {
	// get this score's lo, hi
	lo, hi := NextLoHi(M, I, D, score, penalties)

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
	_, _, current_traceback := M.GetVal(tb_s, tb_k)
	done := false

	for !done {
		CIGAR_rev = CIGAR_rev + traceback_CIGAR[current_traceback]
		switch current_traceback {
		case OpenIns:
			tb_s = tb_s - o - e
			tb_k = tb_k - 1
			_, _, current_traceback = M.GetVal(tb_s, tb_k)
		case ExtdIns:
			tb_s = tb_s - e
			tb_k = tb_k - 1
			_, _, current_traceback = I.GetVal(tb_s, tb_k)
		case OpenDel:
			tb_s = tb_s - o - e
			tb_k = tb_k + 1
			_, _, current_traceback = M.GetVal(tb_s, tb_k)
		case ExtdDel:
			tb_s = tb_s - e
			tb_k = tb_k + 1
			_, _, current_traceback = D.GetVal(tb_s, tb_k)
		case Sub:
			tb_s = tb_s - x
			// tb_k = tb_k;
			_, _, current_traceback = M.GetVal(tb_s, tb_k)
		case Ins:
			// tb_s = tb_s;
			// tb_k = tb_k;
			_, _, current_traceback = I.GetVal(tb_s, tb_k)
		case Del:
			// tb_s = tb_s;
			// tb_k = tb_k;
			_, _, current_traceback = D.GetVal(tb_s, tb_k)
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
