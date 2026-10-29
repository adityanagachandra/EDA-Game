# LLM-to-HLS-C translation loops

Language models can translate a reference algorithm, written in Python, MATLAB or plain C, into synthesizable HLS-C with directives. The output is only useful inside a harness that compiles it, simulates it against the reference, synthesizes it and co-simulates the RTL, feeding each failure back to the model. Designing that loop, and knowing what each stage can and cannot catch, is the engineering problem.

## What makes generated code synthesizable

HLS accepts a subset of C and C++. Models trained mostly on software code readily produce constructs outside it: growing containers such as a vector with push_back, which have no static hardware size; recursion of data-dependent depth, which needs an unbounded stack; unbounded while-loops, which synthesize but have unknown worst-case latency and need an explicit trip-count bound; pointers whose aliasing the tool cannot resolve. Fixed-size arrays, bounded for-loops, streams and arbitrary-precision integer types are the safe vocabulary. The cheapest first check is to compile the code and run the C testbench against the reference outputs, which takes seconds and catches functional errors before any synthesis time is spent. A harness that front-loads cheap checks keeps the expensive tools for candidates that can actually succeed.

## The harness loop

The loop orders its checks by cost: C compile, C simulation, HLS synthesis, then C/RTL co-simulation, each requiring the previous one to pass. On a failure, the harness returns specific feedback: the compiler or HLS error message with the offending source lines, the first mismatching output vector, or the report's II violation reason and the loop it names. A C pass is not enough to accept the code. The C may rely on behavior the hardware does not reproduce, such as reading uninitialized memory or depending on overflow, and co-simulation is where those mismatches appear. When synthesis succeeds but misses its target, for example II = 8 instead of 1, the next prompt should carry the reported limiting resource or recurrence, so the model can restructure the loop instead of guessing.

## Measuring and bounding the system

The right metric is tool-grounded: the fraction of tasks whose output passes C simulation and co-simulation, together with QoR against a baseline implementation; text similarity to a human kernel rewards code that looks right rather than code that works. Tests must cover boundaries: a model that removes a bounds check never triggered by the testbench passes every test, so coverage of edge cases and randomized vectors is the defense. Number formats matter: replacing a 16-bit fixed-point reference with float changes both the results and the area, and the acceptance tolerance must come from the specification. When an external API model drives the loop, the IP-boundary principle applies: decide explicitly what design data may leave the site, and keep tools, vectors and results inside the sandbox.

## Going deeper

The agent_first_eda lessons on the harness loop, golden testing and sandboxing; Kastner, Matai and Neuendorffer, "Parallel Programming for FPGAs" (2018), on writing synthesizable HLS code.
