package tests

import (
	"bufio"
	"encoding/json"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"testing"
	wfa "wfa/pkg"

	"github.com/schollz/progressbar/v3"
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

func randRange(min, max int) uint32 {
	return uint32(rand.IntN(max-min) + min)
}

func TestWavefrontPacking(t *testing.T) {
	for range 1000 {
		val := randRange(0, 1000)
		tb := wfa.Traceback(randRange(0, 7))
		v := wfa.PackWavefrontValue(val, tb)

		valid, gotVal, gotTB := wfa.UnpackWavefrontValue(v)

		if !valid || gotVal != val || gotTB != tb {
			t.Errorf(`test WavefrontPack/Unpack, val: %d, tb: %d, packedval: %x, gotok: %t, gotval: %d, gottb: %d\n`, val, tb, v, valid, gotVal, gotTB)
		}
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

			sequences.Scan()
			s1 := sequences.Text()
			s1 = s1[1:]

			sequences.Scan()
			s2 := sequences.Text()
			s2 = s2[1:]

			x := wfa.WFAlign(s1, s2, testPenalties, true)
			gotScore := x.Score

			if gotScore != -1*expectedScore {
				t.Errorf(`test: %s#%d, s1: %s, s2: %s, got: %d, expected: %d\n`, testName, idx, s1, s2, gotScore, expectedScore)
			}

			idx++
			bar.Add(1)
		}
	}
}
