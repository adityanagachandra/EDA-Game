# Co-simulation as the golden test

In an HLS flow the C model is the specification, and C/RTL co-simulation is the test that checks the generated hardware against it. For agent-driven flows, whether an agent writes the C, tunes the directives or edits the engine itself, co-simulation is the golden check and the most trustworthy reward signal. Knowing what it proves, and what it cannot, is essential before trusting any automated result.

## What co-simulation checks

Co-simulation runs the original C testbench twice: once calling the C function, and once driving the generated RTL through adapters that translate arguments into the RTL's interface protocols, and it compares the outputs. Because the C model is what the RTL was generated from, and it runs fast enough to label many vectors, it is the golden reference. A pass shows that on the tested vectors the outputs match, the handshakes complete and the RTL terminates. It also measures what the HLS report only estimates: the actual latency and interval in cycles on those inputs, including stalls and variable-latency loops. It does not show that the design meets timing after implementation, nor that it is correct on untested inputs; those need static timing analysis and better coverage.

## Golden and differential tests for engine changes

When the HLS engine itself changes, a golden kernel suite with known-good outputs guards it. A good suite exercises distinct features: recurrences, partitioned arrays, streams, branches, dataflow regions, multicycle operators. For a new scheduling heuristic the differential test is: every kernel still passes co-simulation, and the QoR differences in latency, II, registers and area are the intended ones. Byte-identical RTL is the wrong criterion there, because a scheduler change should change schedules; it is the right criterion for a behavior-preserving refactor, such as restructuring the RTL emitter, where any difference is a regression. A binding change that keeps all co-simulations passing but raises area by 9% on two kernels must still be flagged: golden tests guard correctness, the QoR baseline guards quality.

## Where the golden test is weak

A golden test is only as strong as its stimulus. A testbench that never stalls a stream output never exercises the hold-while-not-ready path, so back-pressure bugs slip through; the test must randomize ready and valid. Ten hand-picked vectors miss corner cases that randomized and coverage-directed vectors, labelled by the C model, catch. Floating-point comparison needs care: a permitted reassociation changes rounding, so the comparison must use a tolerance justified by the numerical specification, while integer paths stay bit-exact. For learned generators, co-simulation is the right reward: a reward based only on the HLS report's latency can be maximized by code that computes the wrong function, which is reward hacking, while a co-simulation-gated reward grounds the score in correct behavior.

## Going deeper

The C/RTL co-simulation chapter of the AMD Vitis HLS user guide (UG1399); the agent_first_eda lessons on golden and differential testing and on tool-grounded reinforcement learning.
