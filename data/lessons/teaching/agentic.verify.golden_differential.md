# Golden and differential testing

An agent can write a plug-in in seconds, so the bottleneck of an agent-first engine is proving that the change is correct. Two kinds of test carry that proof. Golden tests check known inputs against known-good outputs. Differential tests compare the changed engine against a baseline on reference designs and check that only the intended things moved. Every engine ships with both, and nothing is promoted without passing them.

## Golden tests

A golden test is a fixed input with an expected output that someone has verified: a small placement that must legalize to a known result, a routed block that must be DRC-clean with a known wirelength, a timing graph with known slacks, an extraction test structure with field-solver capacitances. Goldens check invariants (legal, DRC-clean, connected), exact results where the engine is deterministic, and metric values within stated tolerances where it is not. They are small so they run in seconds inside the harness loop. A golden suite answers the question "is the engine still correct on cases we understand?"

## Differential checks against a baseline

A differential check runs the baseline engine and the changed engine on the same reference designs with the same seeds and compares results object by object. For each object (endpoint, net, region, block) it computes the change and classifies it: intended and correct, intended but by the wrong amount, or unintended. Suppose a router cost plug-in should reduce vias only on nets near an analog block. The differential lists via counts per net before and after; with a tolerance of 0, any net outside the target set whose count changed is a failure, and so is a target net that did not improve. The pass condition is zero unintended movers and zero wrong-amount movers, plus aggregate QoR within its threshold.

## What the pair proves and what it cannot

Together, goldens and differentials prove that the change did what was intended on the reference set and broke nothing visible there. They cannot prove correctness on designs unlike the references: a provider tested only on designs with logic depth above 5 has never exercised its depth condition. They also depend on determinism, because if the baseline itself varies run to run, the differential shows noise everywhere. So reference designs are chosen to exercise each condition of the change, seeds are fixed, and the suite grows each time a bug escapes.

## Going deeper

McKeeman, "Differential Testing for Software" (1998); the OpenROAD-flow-scripts documentation on its regression metrics.
