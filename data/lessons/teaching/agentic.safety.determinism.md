# Determinism and seedability

An engine is deterministic if the same inputs, configuration and seed always produce the same outputs. In an agent-first flow that property is not a nicety. Differential tests, QoR baselines, reproducible bug reports and the agent's own reasoning about cause and effect all assume it. A nondeterministic engine makes every comparison noisy and every claim unverifiable.

## Same inputs, same outputs

Determinism means that every source of variation is an explicit input. Randomized algorithms such as simulated annealing, random restarts or randomized rounding take a seed from the configuration and record it in the run's provenance. Parallel code produces the same result regardless of thread count or scheduling. Nothing depends on wall-clock time, the host name, memory addresses or the order a hash table happens to iterate in. With these rules, an agent can rerun a result exactly, and a change in output can be attributed to a change in input.

## Sources of nondeterminism

The common sources in engines are well known. Floating-point reductions in parallel: addition is not associative, so summing partial results in completion order gives slightly different totals, which can flip a tie-break and change a whole placement. Iteration over unordered containers keyed by pointers, whose order varies between runs. Uninitialized memory. Time-based decisions, such as "stop refining after 60 seconds", which make results depend on machine load. Threads racing to claim work, such as two routing threads committing overlapping nets in either order. The fixes are ordered reductions, containers keyed by stable identifiers, deterministic work partitioning and iteration-based limits instead of time-based ones.

## Why differential tests need determinism

A differential test attributes every change in output to the change in code. If the baseline itself varies run to run, the differential shows movement everywhere and the test can only pass with loose tolerances that hide real regressions. Determinism makes the noise floor zero for a fixed seed: two runs of the baseline are identical, so any movement is caused by the change. Seed variation is still studied deliberately, by running several fixed seeds, but it becomes a controlled experiment rather than an accident. A promotion gate therefore runs each candidate twice with the same seed and rejects it if the outputs differ.

## Going deeper

Demmel and Nguyen, "Fast Reproducible Floating-Point Summation" (ARITH 2013); Goldberg, "What Every Computer Scientist Should Know About Floating-Point Arithmetic" (1991).
