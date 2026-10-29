# Formal apps: X-propagation, connectivity and CDC

Formal apps package property checking for specific, recurring questions, so that engineers who never write a temporal property still get exhaustive answers. The properties are generated from a specification or from the design's structure, and the engines underneath are the BMC, induction and IC3 engines of this domain. X-propagation, connectivity and clock-domain-crossing checks are three of the most widely used.

## X-propagation checks

Registers without reset start with unknown values, and RTL simulation represents them as X. Simulation is unreliable about X: an if statement whose condition is X takes the else branch, which can hide the fact that real hardware might go either way, while other constructs can propagate X pessimistically. A formal X-propagation app instead treats every unknown bit as a free value and asks whether a signal that should always be known, such as a control output or a valid-qualified data bus, can differ between two resolutions of the unknowns. Technically this is a miter: two copies of the design with the same inputs and independent initial values, compared at the checked points. A difference is reported with a trace; no difference proves the X harmless. Checks are written relative to qualifiers: unknown data under a low valid bit is usually by design.

## Connectivity checking

SoC integration wires thousands of signals between IP blocks: interrupts, resets, clocks, debug and test buses. A connectivity app reads a connection specification, often a spreadsheet, and turns each row into a property: the destination equals the source, possibly after a fixed number of register stages, and possibly only when a condition holds, such as test_mode == 0 or a mux select set to a given value. Each check involves only the cone of logic between two points, a few muxes, buffers and flops, so it is fast even on a full chip. Formal search explores every configuration, including debug and test modes that directed simulation rarely sets, so a failing check often reveals a path where some unrelated register can override a connection that the specification says is fixed.

## Clock-domain crossing formal

Signals crossing between asynchronous clock domains can go metastable. Structural CDC analysis finds every crossing and checks that it uses an approved synchronizer, typically two flip-flops in the destination domain with no combinational logic in between. Formal CDC goes further and checks functional rules. A multi-bit value synchronized bit by bit is safe only if at most one bit changes per source cycle, as with gray-coded FIFO pointers, which is an assertion on the source logic. Two synchronized bits that reconverge downstream can resolve in different cycles and briefly form an illegal combination. To check that the design tolerates synchronizer uncertainty, formal tools inject metastability effects, a random extra cycle of latency or a random value at each synchronizer, and re-run the functional properties. Metastability itself cannot be proved away; the checks show that its consequences are contained.

## Going deeper

Seligman, Schubert and Kumar, "Formal Verification: An Essential Toolkit for Modern VLSI Design" (2nd ed., 2023); Cummings, "Clock Domain Crossing (CDC) Design and Verification Techniques Using SystemVerilog" (SNUG 2008).
