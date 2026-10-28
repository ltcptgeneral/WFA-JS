package main

import (
	"unsafe"
	wfa "wfa/pkg"
)

// inputs from JS to WASM
var inputs = map[unsafe.Pointer][]byte{}

// output from WASM to JS
var output string

func main() {
	select {}
}

//go:wasmexport alloc
func alloc(size uint32) unsafe.Pointer {
	if size == 0 {
		size = 1
	}
	buf := make([]byte, size)
	ptr := unsafe.Pointer(&buf[0])
	inputs[ptr] = buf
	return ptr
}

//go:wasmexport free
func free(ptr unsafe.Pointer) {
	delete(inputs, ptr)
}

//go:wasmexport outPtr
func outPtr() unsafe.Pointer {
	return unsafe.Pointer(unsafe.StringData(output))
}

//go:wasmexport outLen
func outLen() uint32 {
	return uint32(len(output))
}

//go:wasmexport wfAlign
func wfAlign(s1 string, s2 string, m int32, x int32, o int32, e int32, doCIGAR bool) int32 {
	penalties := wfa.Penalty{
		M: int(m),
		X: int(x),
		O: int(o),
		E: int(e),
	}

	result := wfa.WFAlign(s1, s2, penalties, doCIGAR)
	output = result.CIGAR

	return int32(result.Score)
}

//go:wasmexport DecodeCIGAR
func DecodeCIGAR(CIGAR string) {
	output = wfa.RunLengthDecode(CIGAR)
}
