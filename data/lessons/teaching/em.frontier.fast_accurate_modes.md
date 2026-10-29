# Fast versus accurate EM modes

EM engines come in fast and accurate versions. Fast modes, closed-form fits, pattern-matched tables and surrogates, run in microseconds and drive continuous modeling loops. Accurate modes, 2D and 3D field solvers, cost seconds to hours. The engineering question is not which is better but where each is trustworthy, and when an agent must recurse from the fast answer to the accurate one.

## Closed forms, pattern libraries and full solvers

Fast modes for capacitance include closed-form fits such as Sakurai-Tamaru and calibrated pattern tables: a production extractor decomposes each net into cross-section and 3D patterns and looks up capacitance per unit length from tables built by running an accurate field solver over sweeps of width, spacing and layer for one process stack. Fast inductance modes assume an ideal or nearby return path. Fast S-parameter estimates use closed-form $Z_0$ and loss per unit length for uniform sections.

Accurate modes solve the field problem directly: 3D boundary-element capacitance, filament or PEEC inductance, full-wave FEM or FDTD. Each fast mode gives something up: unseen geometric patterns, the true return-current distribution, or discontinuities such as vias, pads and bends. Knowing exactly what a fast mode ignores is what tells an agent where it can mislead.

## Validity ranges and when to recurse

A fast mode is trusted only inside its calibrated domain. An agent should recurse to the accurate solver when the geometry falls outside that domain, when the pattern was never seen in calibration, or when the fast mode's error band could change a decision.

Example: a net with -2 ps slack has capacitance from a table with a stated 5% error. If 5% of that net's capacitance can move its delay by more than 2 ps, the sign of the slack is uncertain and the accurate extractor is worth calling; a net with -40 ps does not need it. Measured fast-versus-accurate error as a function of design features, such as metal density, tells the loop where recursion pays. The dangerous case is extrapolation: a fit evaluated far outside its range returns a plausible number with nothing in it to signal the error, so validity checks must be explicit.

## Error budgets in a modeling loop

A continuous modeling loop runs fast engines constantly and spends accurate compute where it matters. A workable policy: use the fast mode to explore, then confirm critical decisions with the accurate engine before committing. Watching the variance of a metric across perturbed fast runs, not only its mean, marks where the fast mode is unreliable.

Example: a loop's budget is 2% capacitance error on critical nets, and the fast mode's measured error on them is about 6%. The fast mode still steers exploration, but those nets get accurate solves before sign-off of the change. A surrogate model counts as a fast mode too and needs its error measured against the accurate solver on held-out designs, including worst cases. If the fast mode systematically underestimates coupling on long buses, crosstalk and Miller-affected delays look better than signoff will show, so the error propagates downstream into timing and noise.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 8 on extraction; Sakurai and Tamaru, "Simple formulas for two- and three-dimensional capacitances", 1983.
