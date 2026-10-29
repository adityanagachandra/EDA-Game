# Logic synthesis

Logic synthesis translates RTL into a netlist of cells from a target library, optimizing for timing, area and power under the constraints in an SDC file. It is the first stage where the design meets the technology, and the first where timing can be computed from real cell delays. Tools include Synopsys Design Compiler and Cadence Genus commercially, and Yosys with ABC in open source.

## Elaboration and generic logic

Synthesis first elaborates the RTL: it resolves parameters and hierarchy, infers registers from clocked processes, infers multiplexers from conditional assignments, and recognizes arithmetic operators and memories. The result is a technology-independent Boolean network: registers plus combinational logic expressed as generic gates or as an and-inverter graph. No library cell has been chosen yet. Inference rules matter: an incompletely specified combinational assignment infers a latch, which is a frequent and costly surprise. Elaboration also flattens or preserves hierarchy as directed, resolves generate loops, and records which RTL signals became which registers, information that equivalence checking and debugging rely on later.

## Optimization and technology mapping

Technology-independent optimization restructures the Boolean network to reduce its size and depth: it removes redundancy, shares common subexpressions and rebalances logic, while preserving the function. Technology mapping then covers the network with cells from the Liberty library, choosing for each region the cell (or cells) that implement it at the lowest cost: area under a delay constraint, or delay under an area budget. After mapping, the tool sizes gates, inserts buffers and may restructure critical paths, using static timing analysis on the mapped netlist as its guide. Mapping quality depends directly on the library: more cell functions and drive strengths give the mapper more options.

## Correctness and what synthesis cannot see

Synthesis must preserve the function, and logic equivalence checking proves that it did, by matching registers between RTL and netlist and checking that each register's next-state logic is equivalent. Synthesis also works with limited physical information. Classic flows estimate wire delay from fanout with a wire-load model, which can be badly wrong in deep-submicron processes where wires dominate delay. Physically aware synthesis runs a quick placement to estimate wirelength instead. A netlist that met timing in synthesis but fails after placement is the classic symptom of that gap. Both effects show why synthesis and placement are increasingly run as one optimization.

## Going deeper

De Micheli, "Synthesis and Optimization of Digital Circuits" (1994); Hachtel and Somenzi, "Logic Synthesis and Verification Algorithms" (1996); the Yosys and ABC documentation.
