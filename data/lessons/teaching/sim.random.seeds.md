# Seeds and random stability

A constrained-random test is deterministic once its seed is fixed: the seed and the code determine every random value. That property is what makes random failures debuggable, and it is fragile in specific ways. SystemVerilog defines random stability rules so that local edits perturb only local stimulus. Regression management, failure reproduction and differential testing all depend on understanding exactly what a seed does and does not pin down.

## Seeds make runs reproducible

The simulator's random generators start from a seed, usually given on the command line. Rerunning the same compiled testbench and design with the same seed reproduces the same stimulus and, if the design is deterministic, the same failure. That is the whole basis of debugging random tests: the regression log must record each test's seed, command line and tool version. A seed is not portable across tools, since the standard does not fix the generator algorithm, nor across code changes that alter the order or number of random calls. A failure seen on seed 42 in one simulator says nothing about seed 42 in another. A regression that seeds from the clock and does not log it throws failures away.

## Random stability in SystemVerilog

SystemVerilog gives every thread and every object its own random number generator, seeded from the generator of its parent when it is created. This is **random stability**: the values an object draws depend on its own seed and its own sequence of `randomize()` calls, not on what other objects do. Adding a new random field or an extra `randomize()` call in one class changes that class's stimulus, while other objects continue to see the same values. Stability breaks when creation order changes: because children are seeded in the order they are created, adding a component or thread before others shifts the seeds of those created after it. Hierarchical seeding is a design choice that makes random stimulus robust to local edits, and a simulator implementer must follow it exactly.

## Managing seeds in regression

A regression balances stability and exploration. A fixed set of seeds gives a stable baseline: if a test that passed yesterday fails today on the same seed, the code change is the suspect. New seeds every night explore new stimulus and find bugs the fixed set never reaches. Coverage-based ranking keeps a small set of seeds that reach the same coverage as the full run, for fast checks on each change. One subtlety: after an RTL fix, a failing seed may pass simply because the stimulus changed and no longer reaches the scenario. Confirming a fix therefore needs evidence that the scenario was exercised, such as a cover property or a directed test, not only a passing seed.

## Going deeper

IEEE Std 1800-2017, clause 18.14 (random stability); Spear and Tumbush, "SystemVerilog for Verification" (3rd ed., 2012), the section on random stability.
