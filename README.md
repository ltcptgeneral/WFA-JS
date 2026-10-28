# Using WFA-JS

Download `wfa.js` and `wfa.wasm`from [releases](https://git.tronnet.net/alu/WFA-JS/releases) to your project. Add to your script:

```
import wfaInit from "./wfa.js"
const wfa = await wfaInit("<path to wasm>")
let result = wfa.wfAlign(...)
let operations = wfa.DecodeCIGAR(result.CIGAR)
```

Where `<path to wasm>` is the path from the site root ie. `./scripts/wfa.wasm`. This will depend on your project structure.