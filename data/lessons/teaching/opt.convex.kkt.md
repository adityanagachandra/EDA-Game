# The KKT conditions

The Karush-Kuhn-Tucker conditions are the optimality conditions of constrained optimization: stationarity of the Lagrangian, feasibility, nonnegative multipliers and complementary slackness. For convex problems they are sufficient, so solving them solves the problem. They explain why Lagrange multipliers price constraints, why inactive constraints carry zero multipliers, and how sizing, placement and budgeting steps project a desired solution onto a feasible set.

## Stationarity and the constrained point

For minimize $f(x)$ subject to $g(x) \le 0$, the KKT conditions are: $\nabla f(x) + \lambda \nabla g(x) = 0$, $g(x) \le 0$, $\lambda \ge 0$ and $\lambda g(x) = 0$. Take a sizing step that wants sizes near $a = 9$ and $b = 6$ under an area budget: minimize $(x - 9)^2 + 4(y - 6)^2$ subject to $x + y \le 13$.

Stationarity gives $2(x - 9) + \lambda = 0$ and $8(y - 6) + \lambda = 0$, so $x = 9 - \lambda/2$ and $y = 6 - \lambda/8$. The unconstrained point $(9, 6)$ uses 15 units, more than 13, so the budget must be tight: $x + y = 13$. Substituting, $15 - \lambda(1/2 + 1/8) = 13$, so $\lambda = 2/0.625 = 3.2$, giving $x = 7.4$ and $y = 5.6$. The more heavily weighted coordinate y moves less, since deviations there cost four times as much.

## The multiplier as a price

The multiplier solves the active constraint once the variables are written in terms of it. For minimize $p(x - a)^2 + q(y - b)^2$ subject to $x + y \le c$ with $a + b > c$, stationarity gives $x = a - \lambda/(2p)$ and $y = b - \lambda/(2q)$, so $\lambda = (a + b - c)/(1/(2p) + 1/(2q))$.

With $p = 1$, $q = 4$, $a = 10$, $b = 9$, $c = 5$: $\lambda = 14/0.625 = 22.4$. The multiplier is the marginal value of the constraint: relaxing the budget from c to $c + \delta$ lowers the optimal cost by about $\lambda\delta$. That is the same shadow-price meaning that LP duals have, and in Lagrangian-relaxation sizing the multipliers on timing arcs measure how much each arc's delay costs at the optimum. A large multiplier marks a constraint worth relaxing, for example a block budget worth renegotiating.

## Complementary slackness and inactive constraints

Complementary slackness, $\lambda g(x) = 0$, says either the constraint is tight or its multiplier is zero. If the unconstrained minimizer already satisfies the constraint, it is the solution and $\lambda = 0$. For minimize $4(x - 9)^2 + (y - 10)^2$ subject to $x + y \le 19$, the unconstrained point $(9, 10)$ uses exactly 19, feasible, so $\lambda = 0$; the constraint is tight but not binding, a degenerate case.

The practical recipe is a case split: try $\lambda = 0$, check feasibility; if infeasible, assume the constraint active and solve for $\lambda$, then check $\lambda \ge 0$. A negative multiplier means the guess was wrong. With many constraints this case split is what active-set QP solvers automate, and interior-point methods instead relax $\lambda g = 0$ to $\lambda g = -\mu$ and drive $\mu$ to zero. For convex problems any point satisfying all four KKT conditions is globally optimal.

## Going deeper

Boyd and Vandenberghe, "Convex Optimization" (2004), chapter 5; Nocedal and Wright, "Numerical Optimization" (2nd ed., 2006), chapter 12 on the theory of constrained optimization.
