# Behavioral versus structural modeling

Hardware description languages let a designer say what a circuit does (behavioral style) or what it is made of (structural style), with continuous-assignment dataflow in between. The distinction matters to tool builders because each style asks a different amount of work from the compiler: structural code is already a netlist, while behavioral code must be inferred into one. Real designs mix the styles, and synthesizer outputs are always structural.

## Structural descriptions are netlists

A structural description instantiates components and wires their ports together. In Verilog, and g1 (y, a, b); or g2 (z, y, c); is structural: two gate primitives and three nets, with no statement of behavior other than what the components themselves do. Module instances work the same way at any level, from a full adder to a processor core.

The natural data structure is a hypergraph. Instances are nodes and nets are hyperedges connecting every pin they touch; a net driving four loads is one hyperedge of five pins. A synthesized gate-level netlist written as Verilog is exactly this: instances of library cells joined by wires, with nothing procedural. Because the function is implied entirely by the components and their connections, a structural netlist needs no inference: a tool reads it and has the graph.

## Behavioral and dataflow descriptions

A behavioral description states function without choosing components. An always block with if statements, case statements and loops is behavioral; so is any testbench, which uses delays, file I/O and other constructs that only simulation supports. Between the two sits dataflow: assign y = (a & b) | c; describes a combinational equation that continuously drives a net, naming no gates and having no procedural sequencing.

For synthesis, behavioral code must be turned into components. A for loop over the bits of a vector, with a constant bound, is unrolled into replicated logic. A clocked always block becomes registers plus their next-state logic. The payoff is productivity: the designer writes what is computed, and a change of function is a change of a few statements rather than of a netlist.

## What each style costs a tool

Inference is the price of behavioral code. The synthesizer applies rules (a clocked block yields flip-flops, an incompletely assigned combinational block yields latches, an array with a matching access pattern yields a RAM) to deduce components. Structural code skips this step, which is why designers hand-instantiate cells whose exact form matters, such as clock-gating cells, synchronizers or pads, inside otherwise behavioral RTL.

Simulation cost runs the other way. A gate-level netlist simulates far slower than its RTL, because each gate is a separate process with its own events instead of one word-level operation. Equivalence checking connects the two worlds: to compare behavioral RTL with its synthesized netlist, the checker first reduces the RTL to a logic graph (often an and-inverter graph) using simple, trusted rules, then compares functions at matched registers with SAT or BDD methods.

## Going deeper

Thomas and Moorby, "The Verilog Hardware Description Language", chapters on structural and behavioral modeling; the Yosys manual, chapter on the Verilog front end.
