# SAT and BDD engines as agent tools

In an agent-first flow, SAT solvers, BDD packages and simulators are not libraries buried inside one program; they are skills an agent calls through a typed tool interface, the way engines are exposed over MCP. The agent decides what to ask, the compiled engine does the search. For that to be safe, the interface must say exactly what goes in, what can come out, under what limits, and how much each kind of answer can be trusted. A Boolean engine with a sloppy interface turns timeouts into false proofs.

## Typed inputs and outputs

The input should be a well-defined format: CNF in DIMACS form, an AIG in AIGER form, or a netlist the tool converts itself. Accepting an AIG is better than demanding CNF, because the engine keeps the structure (gates, fanout, levels) it needs for simulation, structural hashing and choosing cut points. Every call also carries its resource limits (time, conflicts, memory) and its random seed, so a run, and especially an inconclusive run, can be reproduced.

The output has three cases, not two: SAT with a model, UNSAT with an optional proof, or UNKNOWN with the limit that was hit. A BDD tool's output must include the variable order, since node counts are meaningless without it. Collapsing UNKNOWN into UNSAT, or returning a bare boolean, is the most dangerous interface bug a Boolean tool can have.

## What the manifest advertises

The agent selects an engine from its manifest, a short description of what it does, what it is good at, its limits and failure modes, and its tunable parameters. Selection is a cheap text decision; the expensive search stays in the compiled engine. The design data is not in the manifest.

Manifests encode real knowledge. A BDD skill should warn that some functions blow up under every variable order; Bryant proved this for the middle bits of multipliers, so an agent asked to compare two 32 by 32 multipliers should pick a SAT- or hashing-based equivalence engine. Two SAT configurations tuned for structured industrial instances and for uniform random 3-SAT should say so, because miters of synthesized netlists are structured instances where learning and restarts pay off.

## Trusting the answer

Answers differ in how they are checked. A SAT model is its own certificate: the agent evaluates every clause under it, in linear time, before using it. An UNSAT answer has no short certificate in general, so signoff-grade flows ask the solver for a DRAT proof and validate it with an independent checker. UNKNOWN carries no information at all; it must never be recorded as equivalence or as a bug.

When a large equivalence check times out, the sound responses are to decompose the problem (sweep internal equivalences, introduce cut points, prove smaller miters) or to raise the budget, while the transform stays unapproved. Approving because no counterexample was found, or rejecting because no proof was found, both confuse the absence of an answer with an answer.

## Going deeper

Heule, Hunt and Wetzler, "Trimming while checking clausal proofs", FMCAD (2013), on DRAT. Biere, "AIGER format" documentation. Bryant, "On the complexity of VLSI implementations and graph representations of Boolean functions with application to integer multiplication" (1991).
