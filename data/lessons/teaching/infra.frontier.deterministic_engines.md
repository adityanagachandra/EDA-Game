# Deterministic, seedable engines

An engine is deterministic when the same inputs, parameters and seed always give the same result, whatever the thread count or machine. For a human user this is convenience. For an agent it is a prerequisite: differential tests, debugging by replay and tool-grounded rewards all assume that a difference in output was caused by a difference in input. This concept covers how to build such engines and what to record with their results.

## Same seed, same result

Every random choice should come from a generator created from an explicit seed and passed to the code that needs it, never from the time of day or a global generator that threads consume in arrival order. Tie-breaking must use stable keys: comparing pointer addresses gives different answers on each run because allocation addresses change. Iterating a hash map to decide an order is another trap, since iteration order can vary; sort by a stable key first. For parallel code, fix the partition of work and the order in which partial results are merged, so that 8 threads and 16 threads produce the same floating-point sums. A placer run with seed 7 that reports a different HPWL on 8 and 16 threads has failed this property, even though each run on its own looks fine.

## Recording what a result depends on

A result is reproducible only if everything it depends on is recorded. A result manifest stores the engine version and build, every parameter, the seed, the thread count and content hashes of every input file. Hashes matter because a path can point to a file that was edited after the run; the hash identifies the exact bytes. When a flow runs five seeds and keeps the best, the manifest must store the chosen run's seed, not just the list. If two results differ while their manifests are identical, the engine has a hidden source of nondeterminism, something it depends on but does not record, and that is a bug to find, not noise to accept.

## Why agents need determinism

A differential test compares a baseline and a candidate. If baseline reruns already differ, the comparison cannot attribute any difference to the change. An RL reward computed by running a nondeterministic engine on a model's output is noisy, so the model learns more slowly and may be rewarded for luck. A failure that vanishes on rerun cannot be bisected or fixed. Determinism does not remove seed sensitivity: a placer may legitimately give different QoR for different seeds. When it does, an agent comparing two parameter settings should run several seeds per setting and compare distributions, not single runs; determinism guarantees that each of those runs is itself reproducible.

## Going deeper

The agent-first notes in this curriculum on determinism and testing; literature on deterministic parallel placement and routing in OpenROAD and elsewhere.
