package main

import (
	"syscall/js"
	wfa "wfa/pkg"
)

func main() {
	c := make(chan bool)
	js.Global().Get("wfa").Set("wfAlign", js.FuncOf(wfAlign))
	js.Global().Get("wfa").Set("DecodeCIGAR", js.FuncOf(DecodeCIGAR))
	<-c
}

func wfAlign(this js.Value, args []js.Value) interface{} {
	if len(args) != 4 {
		resultMap := map[string]interface{}{
			"ok":    false,
			"error": "invalid number of args, requires 4: s1, s2, penalties, doCIGAR",
		}
		return js.ValueOf(resultMap)
	}

	if args[0].Type() != js.TypeString {
		resultMap := map[string]interface{}{
			"ok":    false,
			"error": "s1 should be a string",
		}
		return js.ValueOf(resultMap)
	}

	s1 := args[0].String()

	if args[1].Type() != js.TypeString {
		resultMap := map[string]interface{}{
			"ok":    false,
			"error": "s2 should be a string",
		}
		return js.ValueOf(resultMap)
	}

	s2 := args[1].String()

	if args[2].Type() != js.TypeObject {
		resultMap := map[string]interface{}{
			"ok":    false,
			"error": "penalties should be a map with key values m, x, o, e",
		}
		return js.ValueOf(resultMap)
	}

	if args[2].Get("m").IsUndefined() || args[2].Get("x").IsUndefined() || args[2].Get("o").IsUndefined() || args[2].Get("e").IsUndefined() {
		resultMap := map[string]interface{}{
			"ok":    false,
			"error": "penalties should be a map with key values m, x, o, e",
		}
		return js.ValueOf(resultMap)
	}

	m := args[2].Get("m").Int()
	x := args[2].Get("x").Int()
	o := args[2].Get("o").Int()
	e := args[2].Get("e").Int()

	penalties := wfa.Penalty{
		M: m,
		X: x,
		O: o,
		E: e,
	}

	if args[3].Type() != js.TypeBoolean {
		resultMap := map[string]interface{}{
			"ok":    false,
			"error": "doCIGAR should be a boolean",
		}
		return js.ValueOf(resultMap)
	}

	doCIGAR := args[3].Bool()

	// Call the actual func.
	result := wfa.WFAlign(s1, s2, penalties, doCIGAR)
	resultMap := map[string]interface{}{
		"ok":    true,
		"score": result.Score,
		"CIGAR": result.CIGAR,
		"error": "",
	}

	return js.ValueOf(resultMap)
}

func DecodeCIGAR(this js.Value, args []js.Value) interface{} {
	if len(args) != 1 {
		println("invalid number of args, requires 1: CIGAR")
		return nil
	}

	if args[0].Type() != js.TypeString {
		println("CIGAR should be a string")
		return nil
	}

	CIGAR := args[0].String()

	decoded := wfa.RunLengthDecode(CIGAR)

	return js.ValueOf(decoded)
}
