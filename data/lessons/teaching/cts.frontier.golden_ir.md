# Golden and differential IR checks

Before a PDN solver change or a PDN ECO is accepted, it must pass golden and differential checks. Golden tests compare against known-correct voltages on reference grids; differential checks compare a new result against the previous one and prove that only the intended things moved. These checks are what let an agent or a developer change the solver or the grid quickly without silently breaking signoff.

## Golden tests for a solver

A golden IR test is a fixed grid and load set with known-correct node voltages. Good golden cases include hand-solvable resistor ladders, small meshes solved by a dense direct solver, and grids with a single supply connection, which are barely positive definite and expose boundary and preconditioner bugs. Golden data should include branch currents for EM, since currents come from small voltage differences across small resistances and amplify errors. The golden data must live outside anything an optimising agent can edit, or the oracle can be changed to match the answer.

## The differential comparison

Compare node drops with an absolute tolerance well below the IR margin, plus a relative one for large drops. For reference and new drops (in mV) of 38.1 and 33.6, 23.1 and 24.3, 27.8 and 29.9, 34.5 and 34.0, 6.0 and 7.5, 38.8 and 40.0, 19.7 and 15.7, with a 2 mV tolerance, the differences are 4.5, 1.2, 2.1, 0.5, 1.5, 1.2 and 4.0 mV, so three nodes fail and the change cannot be promoted. Exact bit equality is the wrong test, because summation order changes floating-point results; determinism, the same result for the same inputs and seeds, is the right requirement.

## What a passing check must also show

A preconditioner change must leave voltages unchanged within tolerance and improve runtime or iterations; a doubled iteration count on one test is a performance regression even if voltages match. A PDN ECO meant for one region must improve that region and move the rest negligibly. A dynamic-solver change needs dynamic golden data, reference waveforms of droop over time for known current profiles, not only static voltages. When two tools disagree on a production block though both pass their golden tests, the first suspects are the inputs: extracted resistance, pad models and current assignment.
