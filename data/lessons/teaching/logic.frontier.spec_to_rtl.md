# Spec-to-RTL agent loops

Language models can now write plausible RTL from a plain-language specification. Plausible is not correct. An agent-first front-end flow wraps generation in a harness loop: the model proposes RTL, tools compile, lint, simulate and formally check it against golden references, and failures go back to the model. This concept applies the harness loop to RTL generation for the circuits of this domain.

## The loop for RTL

The harness loop for RTL has the same shape as for engine code: a request arrives ("a 4-deep synchronous FIFO with full and empty flags"), the model proposes RTL, and the harness runs cheap checks first and expensive ones later. Parse and elaborate (a Verilog front end such as Yosys or Verilator), lint (latches, width mismatches, combinational loops, multiple drivers), simulate against a test bench, check properties formally, and synthesize to see area and timing. Each failure is turned into feedback the model can act on: the line of a lint warning, the first failing cycle of a waveform, a formal counterexample trace.

The loop ends with a reviewable diff and its evidence, not with "the model says it works".

## Golden references for small circuits

For the circuits in this domain, golden references are unusually strong. A combinational block with up to about 20 inputs can be checked exhaustively in simulation. An adder or multiplier can be checked against the language's own arithmetic operator by a combinational equivalence checker, regardless of architecture. An FSM can be checked against a reference machine by sequential equivalence from reset, or against temporal properties ("full and empty are never both asserted").

The agent should prefer checks that are complete over checks that sample: a passing random test bench says little about a corner such as the most negative two's-complement value, whereas an equivalence proof covers every input.

## Where the loop fails

The loop inherits the weaknesses of its checks. If the specification is ambiguous (overlapping or non-overlapping detection, which reset value, signed or unsigned comparison), the model may pick one reading, the test bench another, and the loop will either thrash or converge on the test bench's reading. If the test bench is weak, the model can pass it without implementing the specification.

So the harness should surface ambiguities as questions to the human rather than resolve them silently, keep the golden reference independent of the model's own output (a model-written test bench checking model-written RTL proves little), and record each confirmed specification decision so the next request starts from it.

## Going deeper

The Yosys, Verilator and SymbiYosys documentation; Thakur et al., "VeriGen: A Large Language Model for Verilog Code Generation" (ACM TODAES, 2024); Liu et al., "VerilogEval: Evaluating Large Language Models for Verilog Code Generation" (ICCAD, 2023).
