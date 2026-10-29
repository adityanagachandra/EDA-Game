# Open simulators and formal tools

Open EDA has a full verification toolbox: Verilator and cocotb for fast digital simulation with Python testbenches, SymbiYosys for formal proofs, and ngspice, Xyce and Xschem for transistor-level work. Knowing which tool answers which question, and what each one's model leaves out, keeps a tool builder from trusting a result that the tool was never able to produce.

## Digital simulation with Verilator and cocotb

Verilator compiles synthesizable Verilog and SystemVerilog into an optimized C++ model that is evaluated cycle by cycle, which makes it one of the fastest open simulators for large designs. Its default model is two-state: unknown values become 0 or 1, optionally randomized, so a test that relies on seeing X propagate from an uninitialized register will not see it. Event-driven simulators such as Icarus Verilog keep four-state values at lower speed. cocotb is a framework for writing testbenches in Python as coroutines that await clock edges and signal changes, driving a simulator underneath. Pairing cocotb with Verilator gives Python's libraries for stimulus and checking with a compiled model's speed. For a regression of thousands of random tests, that combination keeps run time manageable.

## Formal checks with SymbiYosys

SymbiYosys (sby) is a front end that reads a design and its assertions through Yosys and runs SAT- or SMT-based engines on it. In BMC mode it unrolls the design for a fixed depth and searches for a counterexample; a clean BMC run of depth 20 proves only that no violation exists within 20 cycles of reset. Prove mode attempts an unbounded proof, for example by k-induction; success means the property holds in every reachable state. Cover mode searches for a trace that reaches each cover statement, which proves a behavior is possible and catches over-constrained assumptions that make a proof vacuous. A failed proof produces a trace that can be viewed in a waveform viewer and replayed in simulation.

## Circuit simulation with ngspice, Xyce and Xschem

Questions about analog behavior, such as whether a comparator settles within 2 ns or how much a ring oscillator's frequency shifts at the slow corner, need circuit simulation with the PDK's device models. ngspice is the widely used open SPICE-class simulator, and the open PDKs ship models for it. Xyce, from Sandia National Laboratories, is a SPICE-compatible simulator designed for parallel execution on large circuits, which helps with big post-layout netlists. Xschem is schematic capture: draw the circuit with PDK symbols, then export a SPICE netlist for ngspice or Xyce, and later compare it against the layout netlist with Netgen. None of the digital tools above can answer these questions, because their models have no voltages, currents or device physics.

## Going deeper

The Verilator, cocotb and SymbiYosys manuals; the ngspice and Xyce user guides; the Xschem documentation.
