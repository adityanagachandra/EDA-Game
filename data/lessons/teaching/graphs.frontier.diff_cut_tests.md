# Differential tests on cut size

When an agent edits a partitioner, the harness decides whether the change is safe. For graph engines the decisive evidence is cheap and exact: recompute the cut, check legality, and compare against a baseline on fixed inputs. Designing these golden and differential tests well is what lets model-made changes be promoted quickly without letting subtle bookkeeping bugs or quality regressions through.

## Recompute, never trust the counter

A partitioner maintains its cut incrementally, move by move. That counter is exactly what a refactor of the gain-update code can break. So the test recomputes the cut from first principles: read the returned partition, scan every net, and count those with pins on both sides. If the engine reports 6 and the recount gives 4, the engine has a bug, because a cut is an exact integer; there is no tolerance.

The same test checks legality: every cell is assigned exactly once, and every block's area lies inside the balance window. A tiny cut from an illegal partition is worthless, and an illegal result that slipped through would mislead every downstream stage.

## Golden and differential checks

Golden tests pin known-good outputs for fixed inputs: tiny netlists with hand-computable cuts, and edge cases such as nets with duplicate pins, single-pin nets, zero-weight nets and fixed cells. They catch off-by-one errors that averages over large designs hide.

Differential tests run the baseline and the changed engine on the same inputs, seed, thread count and build flags, and compare. A change meant to be neutral must produce identical partitions. A change meant to improve quality must be judged per design as well as on average: a mean improvement of 0.4 percent can hide a 6 percent regression on one benchmark. Runtime and memory belong in the comparison too, since a change that doubles runtime is a regression even if every cut matches.

## Gates that fit heuristics

A heuristic change legitimately helps some designs and hurts others, so a gate requiring the new cut to be at most the old cut on every design rejects almost every useful change. Better gates use tolerances and distributions: no design worse by more than a set percentage, the median improved, and outliers flagged for a person to judge.

Reproducibility makes the evidence trustworthy. Inside the sandbox the test pins inputs, seed, thread count and build, and has no network or clock dependence. Multithreaded partitioners need deterministic modes for testing, since thread timing otherwise changes move order. With those controls, a differential report shows exactly what a change did, which is what a human reviewer needs to approve promotion.

## Going deeper

The CI and QoR metric checks in the OpenROAD-flow-scripts repository; the KaHyPar repository's regression tests.
