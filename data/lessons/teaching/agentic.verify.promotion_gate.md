# The promotion gate

The promotion gate is the rule a candidate change must satisfy before it becomes part of the library or the deployed configuration. It turns the evidence from builds, golden tests, differential reports, QoR baselines and determinism checks into one decision. Writing the gate down as an explicit rule, rather than a reviewer's impression, is what lets a harness process many agent changes a day without lowering the bar.

## What the gate checks

A typical gate has hard conditions, each of which must hold: the build succeeds with warnings treated as errors in the skill's own code; every golden test passes, including legality invariants; the differential report shows no unintended movers and no intended movers off by more than the tolerance. Any failure rejects the candidate regardless of how good its QoR looks. Given six candidates of which five build, four of those pass goldens and three of those have clean differentials, three remain in contention. A gate is a conjunction, so the order of checks does not change the answer, only the cost of reaching it.

## QoR and runtime thresholds

Candidates that pass the hard conditions are then checked against thresholds: no metric regresses beyond its threshold on any reference design (for example 0.5% wirelength, 2 ps WNS), and runtime and peak memory stay within a stated budget, such as at most 10% slower than the baseline. A plug-in that improves wirelength by 1% while doubling runtime fails, because the fast loop runs it thousands of times. The thresholds apply per design, not to the suite average, since a large gain on one design must not hide a regression on another.

## Determinism and review

Two further conditions complete the gate. The candidate is run twice with the same seed and must produce identical results, because a nondeterministic plug-in makes every later differential noisy. And the change must be in its declared tier: a Tier 2 plug-in whose diff touches core files is rejected and resubmitted as Tier 3. A candidate that passes everything is presented to a human as a diff with its evidence, and only approval promotes it. The gate makes the human's review fast, because the questions a test can answer are already answered.

## Going deeper

The agent-first EDA notes in this curriculum, sections on verification of engine changes and the skill flywheel.
