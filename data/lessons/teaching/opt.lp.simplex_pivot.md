# Standard form and the simplex pivot

The simplex method is the algebra behind the vertex walk. It rewrites inequalities as equations with slack variables, represents a vertex as a basic solution, and moves to an adjacent vertex by a pivot that swaps one basic variable out and one nonbasic variable in. Every LP solver inside an EDA flow, and every branch-and-bound node of an ILP router, runs pivots like these, so the ratio test and reduced costs are worth knowing by hand.

## Standard form and basic solutions

Standard form writes $Ax \le b$ as $Ax + s = b$ with slacks $s \ge 0$. With $n$ original variables and $m$ inequality rows there are $n + m$ variables and $m$ equations. A basic solution chooses $m$ basic variables, sets the other $n$ to zero, and solves for the basics. When the basics are nonnegative the point is a vertex of the feasible region.

So a problem has at most $\binom{n + m}{m}$ basic solutions. A routing LP with 3 variables and 3 constraints has at most $\binom{6}{3} = 20$, and the count explodes for realistic sizes: $\binom{200}{100}$ is about $9 \times 10^{58}$. The simplex method is practical because it visits only a short sequence of improving bases, typically a small multiple of $m$ pivots, even though the worst case over specially built instances is exponential.

## The minimum-ratio test

A pivot raises one nonbasic variable, the entering variable, from zero. Each basic variable changes along the entering column: row $i$ has $x_{B_i} = b_i - a_i\theta$ where $a_i$ is the column coefficient and $\theta$ the entering value. Rows with $a_i \le 0$ never decrease, so they impose no limit. Rows with $a_i > 0$ hit zero at $\theta = b_i/a_i$. The entering variable rises to the smallest such ratio, and that row's basic variable leaves.

Example: column coefficients 2, 0, 4, 3 with right-hand sides 17, 19, 14, 8. The ratios for positive coefficients are 8.5, 3.5 and 2.67, so $\theta = 8/3 \approx 2.67$ and row 4's basic variable leaves the basis. If no coefficient is positive the entering variable can grow forever, and the LP is unbounded in that direction.

## Reduced costs and the objective step

A reduced cost $d_j$ says how much the objective changes per unit of nonbasic variable $x_j$ while the basics adjust to stay feasible. In a maximization, a positive reduced cost means $x_j$ is worth bringing in; when every reduced cost is nonpositive the current vertex is optimal.

After the pivot the objective changes by $d_j \theta$. Suppose the current objective is 16, $x_e$ has reduced cost 3, and its column has coefficients 4 and 4 in rows with right-hand sides 20 and 15. The ratio test gives $\theta = \min(20/4, 15/4) = 3.75$, and the new objective is $16 + 3 \cdot 3.75 = 27.25$. When the limiting ratio is zero (a degenerate vertex) the basis changes but the point and the objective do not. Pricing rules such as steepest edge choose the entering variable to cut the number of pivots.

## Going deeper

Bertsimas and Tsitsiklis, "Introduction to Linear Optimization" (1997), chapter 3; Vanderbei, "Linear Programming: Foundations and Extensions" (5th ed., 2020), for dictionaries and pivoting rules.
