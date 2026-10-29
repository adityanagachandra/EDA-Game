# Linear programming duality

Every LP has a dual LP whose optimum equals the primal's. The dual is not a curiosity: its variables are prices on the primal's constraints, its feasible points are certificates that bound the primal, and its structure drives Lagrangian relaxation in gate sizing, Benders cuts, shadow prices on routing capacity and the stopping tests in every solver. Learning to write and use the dual is a core skill for reading EDA optimization papers.

## Weak duality gives bounds

For the primal maximize $c^T x$ subject to $Ax \le b$, $x \ge 0$, the dual is minimize $b^T y$ subject to $A^T y \ge c$, $y \ge 0$. One dual variable per primal constraint, one dual constraint per primal variable.

Weak duality: for any primal-feasible $x$ and dual-feasible $y$, $c^T x \le (A^T y)^T x = y^T A x \le y^T b$. So every dual-feasible point certifies an upper bound on the primal maximum. Take maximize $8x + 6y$ subject to $2x + 2y \le 18$ and $x + 6y \le 34$. The point $y_1 = 2$, $y_2 = 5$ satisfies $2y_1 + y_2 = 9 \ge 8$ and $2y_1 + 6y_2 = 34 \ge 6$, so it is dual feasible and certifies the bound $18 \cdot 2 + 34 \cdot 5 = 206$. A better dual point gives a tighter bound.

## Strong duality

Strong duality says that if the primal has a finite optimum, so does the dual, and the two optimal values are equal. No gap remains for LPs. So the dual optimum can be found by solving the primal, and vice versa.

For maximize $7x + 5y$ subject to $3x + y \le 16$ and $3x + 6y \le 21$, the vertices are $(0, 0)$, $(5.33, 0)$, $(5, 1)$ and $(0, 3.5)$, with objective values 0, 37.3, 40 and 17.5. The primal optimum is 40, so the dual minimum is 40 too. Solvers exploit this: a primal-dual pair with equal objectives proves optimality, which is how an LP code knows it can stop and how a test can check a solver's answer without trusting it. When the primal is unbounded the dual is infeasible, and the reverse case is also possible: the dual can be unbounded when the primal is infeasible.

## Optimal dual values at a vertex

At an optimal vertex the dual variables of the tight constraints solve $A_B^T y = c$, where $A_B$ holds the rows of the tight constraints, and slack constraints get $y = 0$. Nonnegative solutions confirm optimality: $c$ lies in the cone of the tight constraints' normals.

Example: maximize $6x + 5y$ subject to (1) $3x + y \le 7$ and (2) $3x + 4y \le 19$. The vertex $(1, 4)$ makes both tight. Solving $3y_1 + 3y_2 = 6$ and $y_1 + 4y_2 = 5$ gives $y_1 = 1$, $y_2 = 1$, both nonnegative, so $(1, 4)$ is optimal with $z = 26 = 7 \cdot 1 + 19 \cdot 1$. In Lagrangian-relaxation gate sizing the multipliers on timing arcs play exactly this role: they are the dual values of the timing constraints, and their optimality conditions become flow conservation on the timing graph.

## Going deeper

Bertsimas and Tsitsiklis, "Introduction to Linear Optimization" (1997), chapter 4; Chen, Chu and Wong, "Fast and exact simultaneous gate and wire sizing by Lagrangian relaxation" (IEEE TCAD, 1999), for duality applied to sizing.
