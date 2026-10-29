# Fast and accurate PDN solvers in a loop

PDN analysis exists in fast and accurate forms. Fast estimates run inside floorplanning, placement and ECO loops; accurate solves run at signoff. An agent that uses both must know what the fast form trades away and when to recurse to the accurate one. This is the PDN instance of the fast-and-accurate engine pattern, and the decision rule is the part an engine designer must build.

## What makes a PDN solver fast

A fast solver coarsens or reduces the grid (merging nodes, using model order reduction for upper layers), uses average or region-level currents instead of per-instance waveforms, and stops an iterative method such as conjugate gradient at a loose tolerance. Each choice cuts runtime by large factors and adds error: coarse grids blur local hot spots, averaged currents hide simultaneity, and loose tolerances leave a residual. All forms share the same core, a symmetric positive-definite conductance matrix built by KCL stamping, which is why both can use the same solver family. A learned surrogate that predicts drop from density maps is another fast form, with its own risk: it is reliable only near its training data.

## When to recurse to the accurate solver

Recurse when the fast answer cannot decide the question. A fast worst drop of 38 mV against a 40 mV budget, with a typical fast-model error of 4 mV, is inconclusive. Recurse where the engines are known to disagree: if they agree within 1 mV almost everywhere but differ by 9 mV near macro edges, macro-edge regions should go to the accurate solver, and the gap should be recorded to improve the fast model. Recurse for dynamic questions that averaged currents cannot answer, and for inputs outside a surrogate's training distribution. Otherwise, the fast solver screens every trial and the accurate solver confirms the winners.

## Where the choice lives in the engine

Solver mode, tolerance and coarsening level are change surface: parameters or a policy selected per run, not edits to the solver core. The grid model and stamping stay in the core. The loop's selection rule can use more than the margin: high variance of fast estimates across nearby candidate edits signals a sensitive region where small modelling errors can reorder candidates, which is where accurate runs earn their cost. Speed claims for a fast mode are meaningful only together with its error against the accurate solver on the same design and corner.

## Going deeper

Saad, "Iterative Methods for Sparse Linear Systems" (2nd ed., 2003).
