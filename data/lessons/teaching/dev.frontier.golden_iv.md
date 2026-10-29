# Golden I-V tests for model changes

Every compact-model change, whether a refit, a plug-in edit or a new model version, must prove that it changed only what it was meant to change. Golden I-V tests compare the candidate against a frozen, trusted baseline at fixed bias points; differential tests check that the differences land only where intended. Designing the suite, choosing the error metric and tolerance, and reading the results are concrete verification skills for an agentic modeling loop.

## What a golden suite contains

A golden suite must exercise every output the simulator uses, in every region: $I_D$-$V_{GS}$ families at low and high drain bias (subthreshold, threshold, strong inversion), $I_D$-$V_{DS}$ families, C-V curves, derivatives such as $g_m$ and $g_{ds}$, and several temperatures and geometries. It should include the Gummel symmetry test around $V_{DS} = 0$, which catches asymmetric or non-smooth behavior that breaks analog and RF simulation.

The golden reference values come from the previously validated model version, run through the same simulator on fixed inputs. Generating them from the code under test would prove nothing.

Benchmark circuits complete the suite: ring oscillators, inverter chains and a few standard cells. They integrate current and charge over real transitions, so they catch changes that point tests miss.

A suite certifies only what it exercises. Coverage of temperature, geometry and bias is the first thing to review when a change passes tests but misbehaves later.

## Checking points against a tolerance

Compare each candidate current with its golden value by relative error,
$$e = \frac{|I_{cand} - I_{gold}|}{I_{gold}},$$
and fail any point with $e$ above the tolerance. Equivalently, compare $\log I$ for subthreshold points.

Example with a 1% tolerance: golden 120.0 µA, candidate 120.6 µA gives 0.5%, a pass. Golden 0.250 µA, candidate 0.254 µA gives 1.6%, a fail, though its absolute error, 4 nA, is far smaller than the first point's 600 nA.

Absolute error would ignore the subthreshold region entirely, since currents span many decades. A tolerance of exactly zero is also a poor choice: compiler and platform differences change the last bits of floating point, so tests would fail with no real change. Tolerances should be tight enough to catch meaningful shifts and loose enough to absorb roundoff, well below typical fitting error.

## Differential checks

A golden test says whether results changed; a differential test says whether they changed only where intended. For an edit that adds a gate-induced drain leakage term, the expected difference is confined to negative $V_{GS}$ at high $V_{DS}$; every other region should match the baseline within tolerance. A difference anywhere else is a bug, however good the new term looks.

Some failures hide from point tests. A change to the charge or capacitance model can pass every DC I-V point but shift a ring oscillator's frequency, because delay depends on charge as well as current. A change to a mobility temperature coefficient can pass a suite that covers only nominal temperature and corrupt hot and cold corners.

The same structure applies at library level: a differential characterization of a few cells, comparing delay and power tables before and after, shows whether a model change propagates to timing as expected. Only after golden and differential checks pass is a change shown to a human as a reviewable diff.

## Going deeper

The repository's docs/eda/agent-first-concepts.md, section 6; McAndrew, "Validation of MOSFET model source-drain symmetry", IEEE Transactions on Electron Devices (2006).
