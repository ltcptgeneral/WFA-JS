import wf_align from "./wfa.js";

const p = {
	x: 4,
	o: 6,
	e: 2
};

// this should be score=24, Alignment=XDIX
console.time("time")
const {CIGAR, score} = wf_align("TCTTTACTCGCGCGTTGGAGAAATACAATAGT", "TCTATACTGCGCGTTTGGAGAAATAAAATAGT", p)
console.timeEnd("time")
console.log(`score: ${score}`);
console.log(`CIGAR: ${CIGAR}`);
