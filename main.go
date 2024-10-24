package main

import (
	"fmt"
	"syscall/js"
	wfa "wfa/pkg"
)

func main() {
	c := make(chan bool)
	js.Global().Set("wfAlign", js.FuncOf(wfAlign))
	<-c
}

func wfAlign(this js.Value, args []js.Value) interface{} {
	if len(args) != 4 {
		fmt.Println("invalid number of args, requires 4: s1, s2, penalties, doCIGAR")
		return nil
	}

	if args[0].Type() != js.TypeString {
		fmt.Println("s1 should be a string")
		return nil
	}

	s1 := args[0].String()

	if args[1].Type() != js.TypeString {
		fmt.Println("s2 should be a string")
		return nil
	}

	s2 := args[1].String()

	if args[2].Type() != js.TypeObject {
		fmt.Println("penalties should be a map with key values m, x, o, e")
		return nil
	}

	if args[2].Get("m").IsUndefined() || args[2].Get("x").IsUndefined() || args[2].Get("o").IsUndefined() || args[2].Get("e").IsUndefined() {
		fmt.Println("penalties should be a map with key values m, x, o, e")
		return nil
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
		fmt.Println("doCIGAR should be a boolean")
		return nil
	}

	doCIGAR := args[3].Bool()

	// Call the actual func.
	result := wfa.WFAlign(s1, s2, penalties, doCIGAR)
	resultMap := map[string]interface{}{
		"score": result.Score,
		"CIGAR": result.CIGAR,
	}

	return js.ValueOf(resultMap)
}
