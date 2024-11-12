package tests

import (
	"bufio"
	"encoding/json"
	"log"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"testing"
	wfa "wfa/pkg"

	"github.com/schollz/progressbar/v3"
	"golang.org/x/exp/constraints"
)

const testJsonPath = "tests.json"
const testSequences = "sequences"

type TestPenalty struct {
	M int `json:"m"`
	X int `json:"x"`
	O int `json:"o"`
	E int `json:"e"`
}

type TestCase struct {
	Penalties TestPenalty `json:"penalties"`
	Solutions string      `json:"solutions"`
}

func randRange[T constraints.Integer](min, max int) T {
	return T(rand.IntN(max-min) + min)
}

func TestWavefrontPacking(t *testing.T) {
	for range 1000 {
		val := randRange[uint32](0, 1000)
		tb := wfa.Traceback(randRange[uint32](0, 7))
		v := wfa.PackWavefrontValue(val, tb)

		valid, gotVal, gotTB := wfa.UnpackWavefrontValue(v)

		if !valid || gotVal != val || gotTB != tb {
			t.Errorf(`test WavefrontPack/Unpack, val: %d, tb: %d, packedval: %x, gotok: %t, gotval: %d, gottb: %d\n`, val, tb, v, valid, gotVal, gotTB)
		}
	}
}

func TestLoHiPacking(t *testing.T) {
	for range 1000 {
		lo := randRange[int](-1000, 1000)
		hi := randRange[int](-1000, 1000)
		v := wfa.PackWavefrontLoHi(lo, hi)

		gotLo, gotHi := wfa.UnpackWavefrontLoHi(v)

		if gotLo != lo || gotHi != hi {
			t.Errorf(`test WavefrontPack/Unpack, lo: %d, hi: %d, packedval: %x, gotlo: %d, gothi: %d`, lo, hi, v, gotLo, gotHi)
		}
	}
}

func GetScoreFromCIGAR(CIGAR string, penalties wfa.Penalty) int {
	unpackedCIGAR := wfa.RunLengthDecode(CIGAR)
	previousOp := '~'
	score := 0
	for _, Op := range unpackedCIGAR {
		if Op == 'M' {
			score = score + penalties.M
		} else if Op == 'X' {
			score = score + penalties.X
		} else if (Op == 'I' && previousOp != 'I') || (Op == 'D' && previousOp != 'D') {
			score = score + penalties.O + penalties.E
		} else if (Op == 'I' && previousOp == 'I') || (Op == 'D' && previousOp == 'D') {
			score = score + penalties.E
		}
		previousOp = Op
	}
	return score
}

func CheckCIGARCorrectness(s1 string, s2 string, CIGAR string) bool {
	unpackedCIGAR := wfa.RunLengthDecode(CIGAR)
	i := 0
	j := 0

	s1Aligned := strings.Builder{}
	alignment := strings.Builder{}
	s2Aligned := strings.Builder{}

	for c := 0; c < len(unpackedCIGAR); c++ {
		Op := unpackedCIGAR[c]
		if Op == 'M' {
			s1Aligned.WriteByte(s1[i])
			alignment.WriteRune('|')
			s2Aligned.WriteByte(s2[j])
			i++
			j++
		} else if Op == 'X' {
			s1Aligned.WriteByte(s1[i])
			alignment.WriteRune(' ')
			s2Aligned.WriteByte(s2[j])
			i++
			j++
		} else if Op == 'I' {

			s1Aligned.WriteRune('-')
			alignment.WriteRune(' ')
			s2Aligned.WriteByte(s2[j])

			j++
		} else if Op == 'D' {
			s1Aligned.WriteByte(s1[i])
			alignment.WriteRune('|')
			s2Aligned.WriteRune('-')
			i++
		}
	}

	if i == len(s1) && j == len(s2) {
		return true
	} else {
		log.Printf("\n%s\n%s\n%s\n i=%d, j=%d, |s1|=%d, |s2|=%d\n", s1Aligned.String(), alignment.String(), s2Aligned.String(), i, j, len(s1), len(s2))
		return false
	}
}

func TestWFA(t *testing.T) {
	content, _ := os.ReadFile(testJsonPath)

	var testMap map[string]TestCase
	json.Unmarshal(content, &testMap)

	for k, v := range testMap {
		testName := k

		testPenalties := wfa.Penalty{
			M: v.Penalties.M,
			X: v.Penalties.X,
			O: v.Penalties.O,
			E: v.Penalties.E,
		}

		sequencesFile, _ := os.Open(testSequences)
		sequences := bufio.NewScanner(sequencesFile)
		solutionsFile, _ := os.Open(v.Solutions)
		solutions := bufio.NewScanner(solutionsFile)

		bar := progressbar.Default(305, k)

		idx := 0

		for solutions.Scan() {
			solution := solutions.Text()

			expectedScore, _ := strconv.Atoi(strings.Split(solution, "\t")[0])
			expectedCIGAR := strings.Split(solution, "\t")[1]

			sequences.Scan()
			s1 := sequences.Text()
			s1 = s1[1:]

			sequences.Scan()
			s2 := sequences.Text()
			s2 = s2[1:]

			x := wfa.WFAlign(s1, s2, testPenalties, true)
			gotScore := x.Score
			gotCIGAR := x.CIGAR

			if gotScore != -1*expectedScore {
				t.Errorf(`test: %s#%d, s1: %s, s2: %s, got: %d, expected: %d`, testName, idx, s1, s2, gotScore, expectedScore)
				os.Exit(1)
			}

			if gotCIGAR != expectedCIGAR {
				checkScore := GetScoreFromCIGAR(gotCIGAR, testPenalties)
				CIGARCorrectness := CheckCIGARCorrectness(s1, s2, gotCIGAR)
				if checkScore != gotScore && checkScore != -1*expectedScore { // nonequivalent alignment
					t.Errorf(`test: %s#%d, s1: %s, s2: %s, got: [%s], expected: [%s]`, testName, idx, s1, s2, gotCIGAR, expectedCIGAR)
					t.Errorf(`test: %s#%d, recalculated score: %d`, testName, idx, checkScore)
					os.Exit(1)
				}
				if !CIGARCorrectness {
					t.Errorf(`test: %s#%d, s1: %s, s2: %s, got: [%s], expected: [%s]`, testName, idx, s1, s2, gotCIGAR, expectedCIGAR)
					os.Exit(1)
				}
			}

			idx++
			bar.Add(1)
		}
	}
}
