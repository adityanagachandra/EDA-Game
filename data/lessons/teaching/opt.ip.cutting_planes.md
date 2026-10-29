# Cutting planes and branch and cut

A cutting plane is an inequality that every integer solution satisfies but the current LP solution violates. Adding it shrinks the relaxation toward the integer hull without losing any integer point, so bounds improve and the branch-and-bound tree shrinks. Modern MILP solvers are branch-and-cut codes, and knapsack-like rows from legalization, covering and capacity constraints in EDA models are exactly where cover cuts and Chvatal-Gomory cuts earn their keep.

## Cover inequalities

For a binary knapsack row $\sum a_i x_i \le b$, a cover is a set C whose weights sum to more than b. No feasible integer point can set every variable in C to 1, so $\sum_{i \in C} x_i \le |C| - 1$ is valid.

Example: $4x_1 + 8x_2 + 5x_3 + 2x_4 \le 12$. The set $\{x_2, x_3\}$ weighs 13, more than 12, so it is a cover and $x_2 + x_3 \le 1$ holds for every integer solution. The LP point $x = (0.2, 0.9, 0.2, 1)$ satisfies the knapsack row ($0.8 + 7.2 + 1 + 2 = 11$) but has $x_2 + x_3 = 1.1$, violating the cut by 0.1. Adding the cut removes this fractional point. Minimal covers, where dropping any item leaves a non-cover, give the strongest cuts, and lifting extends them to variables outside the cover.

## Chvatal-Gomory rounding

For nonnegative integer variables, any valid inequality can be strengthened by scaling and rounding. Multiply $5x + 6y \le 14$ by $u = 1/3$: $1.667x + 2y \le 4.667$. Rounding the left coefficients down keeps the inequality valid because $x, y \ge 0$: $x + 2y \le 4.667$. Now the left side is an integer for integer x and y, so the right side can be rounded down too: $x + 2y \le 4$.

That is a Chvatal-Gomory cut. It is valid for every integer point but can cut off fractional ones, such as $x = 0$, $y = 2.33$ allowed by the original row. Gomory's fractional cuts, generated from a row of the optimal simplex tableau, are a special case, and every valid inequality for a pure integer program can be obtained by repeating this procedure finitely often, although in practice solvers add only a few rounds.

## Validity and branch and cut

A cut must never remove a feasible integer point. Checking a proposed cut on a small instance is direct: enumerate the integer points that satisfy the constraints and count those that violate the cut. For $5x + 2y \le 8$ with $0 \le x, y \le 3$, the feasible points are $(0,0)$ to $(0,3)$, $(1,0)$ and $(1,1)$, six in all, and none has $x + y > 4$, so $x + y \le 4$ is valid; it is also useless, since no LP point violates it either. A cut is only worth adding if it is both valid and violated.

Branch and cut adds violated cuts at the root and at nodes, re-solves the LP, and branches when cuts stop helping. The trade-off is LP size versus bound quality, so solvers manage a cut pool and drop inactive cuts. Cuts derived from problem structure, such as capacity cuts in routing, often beat generic ones.

## Going deeper

Wolsey, "Integer Programming" (2nd ed., 2020), chapters 8 and 9; Nemhauser and Wolsey, "Integer and Combinatorial Optimization" (1988), for cover inequalities and Chvatal-Gomory closure.
