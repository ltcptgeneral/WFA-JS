package tests

import (
	"bufio"
	"encoding/json"
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

			x := wfa.WFAlign(s1, s2, testPenalties, false)
			gotScore := x.Score

			if gotScore != -1*expectedScore {
				t.Errorf(`test: %s#%d, s1: %s, s2: %s, got: %d, expected: %d\n`, testName, idx, s1, s2, gotScore, expectedScore)
			}

			idx++
			bar.Add(1)
		}
	}
}
