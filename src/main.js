import wf_align from "./wfa.js";
import fs from "fs";
import ProgressBar from "progress";

let data = fs.readFileSync("./tests/tests.json");
data = JSON.parse(data);
const sequences = fs.readFileSync("./tests/sequences").toString().split("\n");
//const total = sequences.length;
const total = 500; // skip the later tests because of memory usage

for (const test_name of Object.keys(data)) {
	const test = data[test_name];
	const penalties = test.penalties;
	const solutions = fs.readFileSync(test.solutions).toString().split("\n");
	const bar = new ProgressBar(":bar :current/:total", { total: total / 2 });
	console.log(`test: ${test_name}`);
	let correct = 0;
	let j = 0;
	for (let i = 0; i < total; i += 2) {
		const s1 = sequences[i].replace(">");
		const s2 = sequences[i + 1].replace("<");
		const { CIGAR, score } = wf_align(s1, s2, penalties);
		const solution_score = Number(solutions[j].split("\t")[0]);
		if (solution_score === -score) {
			correct += 1;
		}
		j += 1;
		bar.tick();
	}
	console.log(`correct: ${correct}\ntotal: ${total / 2}\n`);
}