# Shadow prices and complementary slackness

A dual value is a price: how much the optimal objective improves per extra unit of a constraint's right-hand side. In EDA that is the value of one more routing track on an edge, one more unit of power budget, or one more picosecond on a path. Complementary slackness links primal and dual: a constraint with slack is worth nothing at the margin. Together these turn an LP solution into a sensitivity report that tells the designer where resources are binding.

## Slack at the optimum

The slack of constraint $i$ at a point is $r_i - a_i^T x$, the unused part of its right-hand side. At the optimum some constraints are tight (slack 0) and others are not.

Take maximize $7x + 11y$ subject to (1) $2x + 3y \le 11$, (2) $3x + 5y \le 18$, (3) $x \le 5$, with optimum $(1, 3)$. Constraint (1): $2 + 9 = 11$, slack 0. Constraint (2): $3 + 15 = 18$, slack 0. Constraint (3): $x = 1$ against 5, slack 4. Constraints (1) and (2) are binding; (3) is not. In a global-routing LP, the binding constraints are the edges whose capacity is fully used: the congestion bottlenecks. A non-binding constraint can move a little without changing the optimum, which is the first step of sensitivity analysis.

## Complementary slackness

Complementary slackness says that at optimal primal $x$ and dual $y$, for every constraint, $y_i \cdot \text{slack}_i = 0$, and for every variable, $x_j$ times its dual slack is 0. A constraint with positive slack has $y_i = 0$. A constraint with $y_i > 0$ is tight.

In the example, (3) has slack 4, so $y_3 = 0$. The tight constraints carry the duals: from $2y_1 + 3y_2 = 7$ and $3y_1 + 5y_2 = 11$ we get $y_1 = 2$, $y_2 = 1$. Check: $11 \cdot 2 + 18 \cdot 1 + 5 \cdot 0 = 40 = 7 \cdot 1 + 11 \cdot 3$, so primal and dual objectives agree. These conditions are the LP case of the KKT conditions, and in primal-dual interior-point methods they are the quantity driven to zero.

## Shadow prices and their range

The dual value $y_i$ is the shadow price of constraint $i$: raising $r_i$ by $\delta$ raises the optimum by $y_i \delta$, as long as the same constraints stay tight. In the LP above $y_1 = 2$. Raising $r_1$ from 11 to 11.5 moves the vertex to $(3.5, 1.5)$, and the optimum rises from 40 to 41, exactly $2 \cdot 0.5$.

The price holds only over a range. Raise $r_1$ to 12 and the same two tight constraints would put the vertex at $(6, 0)$, which violates $x \le 5$: constraint (3) becomes binding, the basis changes, and so does the price. Solvers report this range. For an EDA engineer the shadow price ranks investments: a routing edge with a large dual is the one where an extra track, a wider channel or a layer change pays most. A zero dual means that buying more of that resource does not help at all, however congested it looks.

## Going deeper

Bertsimas and Tsitsiklis, "Introduction to Linear Optimization" (1997), chapter 5 on sensitivity analysis; Vanderbei, "Linear Programming" (5th ed., 2020), chapter on duality and complementary slackness.
