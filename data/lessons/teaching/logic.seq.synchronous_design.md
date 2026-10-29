# Synchronous design

Synchronous design is the discipline that makes digital design tractable: all state lives in flip-flops clocked by a common clock, and combinational logic between them settles within one period. That discipline is the assumption behind static timing analysis, cycle-based simulation, retiming, and sequential equivalence checking. Knowing its rules, and why each exists, explains what EDA tools can and cannot promise about a design that breaks them.

## The rules of the discipline

A synchronous circuit obeys three structural rules:

1. Every storage element is an edge-triggered flip-flop (or a clocked latch used in a disciplined two-phase scheme).
2. Every flip-flop's clock comes from the clock network, not from data logic.
3. Every cycle in the circuit passes through at least one flip-flop: there are no combinational loops.

Under these rules the circuit's behavior is a function from (state, inputs) to (next state, outputs), applied once per clock cycle, and glitches inside the combinational logic do not matter as long as it settles before the setup window. A tool can reason about each register-to-register path in isolation.

## What tools gain from it

Synchronous structure is what makes the main engines possible:

- **Static timing analysis** checks every path against one period without simulating any vectors, because only the path delays and the clock edges matter.
- **Cycle-based simulation** (Verilator) evaluates the logic once per cycle in topological order instead of tracking every gate-level event.
- **Sequential equivalence and retiming** treat the design as a finite-state machine with a well-defined next-state function.

A combinational loop breaks all three: the timing graph has a cycle, the simulator cannot levelize, and the logic may oscillate or hold state. Lint tools and synthesis report combinational loops as errors for this reason.

## Clocks, enables and resets done safely

Gating a clock with an AND gate saves power but can glitch the clock and adds skew. The safe form is a clock-gating cell, which latches the enable while the clock is low so the gated clock can only produce whole pulses; timers check its enable with a dedicated clock-gating check. The functionally equivalent alternative is a recirculating mux enable on the D input.

Resets have their own rule. An asynchronous reset may be asserted at any time, but its *release* must be synchronized to the clock, or flip-flops can leave reset in different cycles (or go metastable on the recovery and removal checks). The standard reset synchronizer is "asynchronous assert, synchronous deassert". Clock-domain crossings, multiple clocks and generated clocks are where the simple discipline ends and dedicated checks begin.

## Going deeper

Harris and Harris, "Digital Design and Computer Architecture" (2nd ed., 2012), section 3.3; Cummings and Mills, "Synchronous Resets? Asynchronous Resets? I am so confused!" (SNUG 2002).
