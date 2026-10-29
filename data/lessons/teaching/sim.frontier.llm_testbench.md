# LLM-generated testbenches checked by simulation

Language models can write testbenches, assertions and reference models quickly. The question is how to trust them. A generated testbench that compiles and passes tells you almost nothing, because a testbench that checks nothing also passes. Simulation supplies the means to measure what matters: does the testbench run, does it reach the intended states, and does it fail when the design is wrong? Mutation testing turns the last question into a number.

## Compile, run and actually check

A harness accepts a generated testbench in stages. First it must compile: errors such as a call to a UVM method that does not exist go straight back to the model as feedback, the cheapest loop. Then it must run to completion on the original RTL without failures. Neither step shows that it checks anything. Common hollow patterns: expected values computed from the DUT's own internal signals, so the DUT is checked against itself; assertions whose antecedents never occur, so every attempt passes vacuously; scoreboards that compare only the first transaction or only a valid flag; stimulus that never leaves reset or never reaches the interesting state. Coverage on the antecedents and on the intended scenarios exposes the reach problems. Checking strength needs a separate measurement.

## Mutation score

**Mutation testing** injects small deliberate bugs (mutants) into the RTL: an operator flipped, a constant changed, a condition inverted. Each mutant is simulated with the testbench; the mutant is **killed** if some test fails. Some mutants are **equivalent**: they do not change observable behavior, so no test can kill them, and they are excluded. The **mutation score** is killed over killable: with 60 mutants, 4 equivalent, and 42 killed, the score is $42 / (60 - 4) = 75\%$. Two testbenches with identical code coverage can have very different scores, because coverage measures what ran and mutation measures what was checked. The original RTL must still pass, otherwise a testbench that always fails would score 100%.

## Failure modes and safe use

A surviving mutant has two common explanations: the tests never exercise the mutated logic (a stimulus gap) or they exercise it but never check the affected output (a checking gap), and both are actionable. Generated reference models deserve the most suspicion. A model produced from the same misreading of the specification as the RTL will agree with it and hide the bug, so a generated model is treated as untrusted until it agrees with an independent source, such as directed cases derived by hand from the specification or an existing golden model. Hallucinated APIs, plausible but wrong protocol timing, and over-constrained stimulus are typical defects. In the harness loop they are caught by compiling, by coverage of intended scenarios, and by mutation score, rather than by reading the code.
