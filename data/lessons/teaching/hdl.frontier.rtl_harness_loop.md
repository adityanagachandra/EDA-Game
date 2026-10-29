# Compile-and-simulate loops for agent-written RTL

Agents that write and refactor RTL work inside a harness: a loop that compiles, lints, simulates and synthesizes each candidate and feeds the results back. The harness, not the model, decides what counts as success, so its design determines whether an agent produces working hardware or code that merely passes the checks it was shown. This concept applies the edit-build-test loop to HDL engines concretely.

## The loop and its order of checks

A request arrives, the model proposes an RTL edit, and the harness writes it and runs checks in increasing order of cost: lint and compile first (seconds), then unit simulation tests, then synthesis checks and equivalence or formal checks. A failure at any stage goes back to the model and the loop repeats; when everything passes, the diff and the evidence go to a human reviewer.

Cheap failures first is the rule, because a width or syntax error found in seconds makes a ten-minute regression pointless. Iteration speed bounds how many attempts an agent gets, so the build step deserves engineering: if 3.5 of every 4 minutes are spent recompiling a Verilator model, a compiler cache, incremental builds of changed modules or an interpreted simulator for the inner loop pay off more than anything else. A fast strict engine (Verilator) in the inner loop and an accurate four-state simulator as a later gate is the fast-versus-accurate pattern applied to simulation.

## What feedback to return

The model can only fix what the feedback describes. For a compile failure, return the diagnostic verbatim: file, line, construct and error class. For a simulation failure, return the failing test, the first mismatching cycle, the expected and actual values and the signals involved, not the whole waveform database. A minimal counterexample is what a human would debug from, and it fits in a context window.

The harness also decides which diagnostics block acceptance. Synthesis findings matter even when simulation passes: an inferred latch means some path leaves an output unassigned, which the tests did not cover, so the harness should fail the candidate and return the latch report. Width-truncation lint is worth gating for the same reason. Hidden or held-out tests, and equivalence against a reference, catch an agent that special-cases the visible tests' inputs.

## Failure modes and acceptance

Generated RTL and testbenches can hang, through a combinational loop or a missing finish, or exhaust memory, so simulations run in a sandbox with time and memory limits that turn a hang into a reportable failure. An agent that oscillates between two edits, each fixing one test and breaking another, needs the loop to detect the cycle, stop and present both failures together.

Acceptance stays with a human, who reviews the diff with its evidence: tests passed, lint clean, synthesis report, equivalence result. Measuring the agent is part of the design too: the fraction of tasks whose final RTL compiles, passes hidden tests and is free of blocking lint and synthesis issues says far more than compile-pass rate or lines generated.

## Going deeper

The Verilator manual on lint and build options; the agent-first EDA notes in this curriculum on the harness loop and change tiers.
