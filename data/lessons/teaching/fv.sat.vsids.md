# The VSIDS decision heuristic

Which variable a CDCL solver decides next matters as much as how it learns. VSIDS (variable state independent decaying sum) answers with a score per variable that is bumped whenever the variable takes part in a conflict and decays over time, so the search focuses on the variables of recent conflicts. It is cheap, needs no look-ahead, and with small variations it is the default heuristic in most industrial and open-source solvers.

## Bumping with a growing increment

Each variable has an activity, initially 0. On each conflict, every variable in the learned clause is bumped: its activity grows by the current increment. Then the increment is divided by a decay factor $d < 1$, so later bumps count more than earlier ones. With $d = 0.5$ the increments are 1, 2, 4, 8, ... Suppose three conflicts learn clauses over $\{x_1, x_2\}$, then $\{x_2, x_3\}$, then $\{x_2\}$. The increments are 1, 2 and 4, so $x_1$ ends at 1, $x_3$ at 2, and $x_2$ at $1 + 2 + 4 = 7$. With $d = 0.8$ the increments are 1, 1.25, 1.5625, ..., a gentler bias toward the present. Real solvers use decay factors between about 0.8 and 0.95, and rescale all activities when the increment grows too large for floating point.

## Picking the next decision

At each decision the solver picks the unassigned variable with the highest activity. Assigned variables are skipped, however active, because they already have values. In the example above, with $x_2$ assigned, the candidates are $x_1$ (activity 1) and $x_3$ (activity 2), so the solver decides on $x_3$. The value is usually taken from phase saving, the variable's most recent value, rather than from the score. To make this fast, solvers keep unassigned variables in a binary heap ordered by activity: bumping a variable is a heap update, finding the best is constant time, and variables re-enter the heap when backjumping unassigns them. The effect is that after a conflict the solver tends to branch on the variables involved in it, which keeps the search local to the part of the formula currently causing trouble.

## Explicit decay gives the same ranking

The original formulation decays explicitly: on each conflict, multiply every activity by $d$, then add 1 to each learned-clause variable. On the same three conflicts with $d = 0.5$: after the first, $x_1 = x_2 = 1$. After the second, decay to 0.5 and 0.5, then bump $x_2$ and $x_3$: $x_1 = 0.5$, $x_2 = 1.5$, $x_3 = 1$. After the third, decay to 0.25, 0.75 and 0.5, then bump $x_2$: $x_2 = 1.75$. These are exactly the growing-increment scores (1, 7, 2) multiplied by $d^2 = 0.25$, and in general by $d^{k-1}$ after $k$ conflicts. Because all scores are scaled by the same factor, the ranking is identical. The growing increment is preferred because it touches only the bumped variables instead of every variable on every conflict.

## Going deeper

Moskewicz et al., "Chaff: Engineering an Efficient SAT Solver" (DAC 2001); Een and Sorensson, "An Extensible SAT-solver" (SAT 2003); Biere and Froehlich, "Evaluating CDCL Variable Scoring Schemes" (SAT 2015).
