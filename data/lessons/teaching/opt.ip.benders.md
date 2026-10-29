# Benders decomposition

Benders decomposition splits a mixed-integer program into a master problem over the complicating integer variables and an LP subproblem over the continuous rest. The subproblem's dual returns cuts that teach the master about the continuous part. It suits EDA models with discrete structural choices and many continuous details, such as macro positions with wire flows, or buffer locations with timing LPs, and it is the formal version of fix-the-structure-then-optimize-the-details loops tools already run.

## Master and subproblem

Write the problem as minimize $f^T y + c^T x$ subject to $Wx \ge h - Gy$, $x \ge 0$, y integer. The integer y couples everything; for fixed y the rest is an LP in x. Benders keeps y in a master problem and replaces the LP by an auxiliary variable $\theta$ that underestimates its optimal cost.

The master solves minimize $f^T y + \theta$ subject to the cuts found so far. For its proposed $y$, the subproblem LP is solved. If it is infeasible, an extreme ray of its dual yields a feasibility cut excluding that y and others like it. If it is feasible, its optimal dual gives an optimality cut. When the continuous part splits into independent blocks, such as one routing region each, the subproblem becomes many small LPs solved in parallel, each returning its own cut. Benders partitions variables; Lagrangian relaxation instead dualizes the coupling constraints.

## Evaluating an optimality cut

By LP duality the subproblem cost for fixed y is $\max\{u^T(h - Gy) : W^T u \le c, u \ge 0\}$. Any dual-feasible u therefore gives a lower bound valid for every y: $\theta \ge u^T(h - Gy)$. The optimal u at the current y makes the bound tight there.

Example: $h = (5, 13)$, $G = (1, 4)^T$ with a scalar y, and an earlier subproblem returned $u = (0, 3)$. At $y = 1$ the cut says $\theta \ge 0 \cdot (5 - 1) + 3 \cdot (13 - 4) = 27$. At $y = 2$ it says $\theta \ge 3 \cdot (13 - 8) = 15$. Each cut is a linear function of y, so the master's view of the subproblem cost is the maximum of all cuts: a piecewise-linear underestimate that improves with every iteration. The common slips are adding $Gy$ instead of subtracting it and dropping a dual component.

## Bounds, convergence and fit

Each iteration gives two bounds for a minimization. The master's optimum, with only some cuts, is a lower bound. The master's y evaluated exactly, $f^T y$ plus the subproblem cost, is an upper bound. Stop when they meet within a tolerance: lower bound 812 and upper bound 830 put the best solution within $18/812 \approx 2.2\%$ of optimal.

Convergence can be slow, because each cut may raise the lower bound a little and the master, an integer program, is re-solved with more rows every time. Remedies include stronger (Pareto-optimal) cuts, several cuts per iteration, and branch-and-Benders-cut, which adds cuts inside one search tree instead of restarting. Benders is a poor fit when fixing the integers does not leave an easy, decomposable continuous part.

## Going deeper

Benders, "Partitioning procedures for solving mixed-variables programming problems" (Numerische Mathematik, 1962); Rahmaniani, Crainic, Gendreau and Rei, "The Benders decomposition algorithm: a literature review" (European Journal of Operational Research, 2017).
