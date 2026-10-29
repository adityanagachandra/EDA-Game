# The logic synthesis flow

Logic synthesis turns register-transfer-level code into a netlist of library cells that computes the same function between registers. It sits between elaboration, which produces a word-level design, and physical design, which places and routes the cells. A tool designer needs the stages and the representation each one uses, because almost every later quality problem (area, timing, congestion) traces back to a choice made here.

## From RTL to a Boolean network

Elaboration resolves parameters, generate blocks and hierarchy, leaving word-level operators (adders, comparators, multiplexers) and processes (always blocks). The front end then infers storage: every variable assigned on a clock edge becomes a flip-flop, and an incompletely assigned combinational variable becomes a latch. What remains between registers is combinational logic. Word-level operators are lowered (bit-blasted) into single-bit logic, giving a **Boolean network**: a directed acyclic graph whose nodes are Boolean functions of their fanins. In Yosys this is the path from RTLIL processes through `proc` to cells and finally to single-bit gates; in ABC the network becomes an and-inverter graph. From here on synthesis optimizes only the combinational network, with registers as its boundary.

## Technology-independent optimization, then mapping

The classic flow splits optimization in two. **Technology-independent** optimization restructures the network using a library-free cost: literal count, AIG node count, and logic depth. It runs two-level minimization of node functions, algebraic factoring and extraction, don't-care simplification, and AIG rewriting and balancing. **Technology mapping** then covers the optimized network with cells from a library (or with K-input lookup tables for FPGAs), choosing cells by area, delay and load. After mapping come sizing, buffering and local resynthesis. The split is a heuristic: the cost used before mapping is only a proxy for what the library will charge, and the mapper is biased by the structure it is handed.

## Why every stage is heuristic

Even minimizing a single two-level function is NP-hard: it contains set covering. Multi-level optimization has a vastly larger space (any circuit computing the function), so tools apply local, greedy transformations that improve a cost and preserve function. Two consequences matter for a tool designer. First, results depend on the order and mix of transformations, which is why scripts such as ABC's resyn2 exist and why agents can search over them. Second, small input changes can redirect many greedy decisions, so engineering-change flows patch an existing netlist rather than resynthesize. Correctness is never heuristic: every transformation must preserve the function, and combinational equivalence checking proves it.

## Going deeper

De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), the introductory chapters; Hachtel and Somenzi, "Logic Synthesis and Verification Algorithms" (1996), the introduction; the Yosys manual's description of the synth command.
