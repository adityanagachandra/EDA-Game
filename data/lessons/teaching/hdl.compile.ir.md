# Intermediate representations for hardware

An intermediate representation (IR) is the form a compiler transforms between parsing and output. Hardware compilers use several, from process-level IRs that still contain branches, through word-level netlists, to bit-level and-inverter graphs. Choosing what each level keeps (words, names, processes, registers) decides which optimizations are easy, which is the central design decision of any synthesis or simulation tool.

## Why compilers lower the AST

An AST mirrors source syntax: there are dozens of statement and expression forms, and many ways to write the same hardware. An IR gives a small, uniform vocabulary, such as cells connected by nets or SSA operations with explicit operands, that passes can analyze and rewrite without caring how the source was written. Most hardware IRs are graphs: operations or cells are nodes, and the values or nets connecting them are edges.

The level of the IR fits the optimization. Logic optimization by rewriting works best on an and-inverter graph, a bit-level graph of two-input AND nodes with inverted edges, which is what ABC rewrites. Datapath optimization needs word-level operators, such as a 32-bit add or a multiplexer, which a bit-level graph has already dissolved.

## Levels and progressive lowering

Real flows lower through several levels. Yosys keeps procedural code as processes in RTLIL until the proc passes convert them into multiplexers and flip-flops; then word-level cells are optimized; then techmap lowers them to single-bit gates; then ABC works on an AIG. MLIR-based compilers such as CIRCT make this explicit as progressive lowering through dialects. Each step is small, so it is easy to test and to reason about.

Two kinds of information must survive lowering. Registers must stay explicit, because they separate cycles: combinational optimization works within a cycle, and moving state (retiming) is a deliberate, separate transformation. Source names and locations must travel as attributes, so that a timing report or a lint warning on a netlist cell can be traced back to an RTL line.

## SSA, verifiers and text formats

Many modern IRs use static single assignment: every value has exactly one definition, so use-def chains are explicit and dataflow optimizations are simple. A combinational netlist is already close to SSA, since every net has exactly one driver; registers break cycles the way phi functions merge values at loop headers.

Two engineering properties make an IR practical. A textual, round-trippable format, such as RTLIL text or MLIR assembly, lets an engineer dump a pass's input and output, diff them and replay the pass alone, which is how regression tests for compiler passes are written. A verifier that checks well-formedness after every pass catches a malformed IR at the pass that created it rather than many passes later. Both matter even more when passes are written or modified automatically.

## Going deeper

Wolf, "Yosys Manual", chapters on RTLIL and the synthesis flow; Lattner et al., "MLIR: Scaling Compiler Infrastructure for Domain Specific Computation" (CGO 2021).
