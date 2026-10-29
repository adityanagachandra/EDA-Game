# The Lagrangian dual function

The Lagrangian dual function minimizes the Lagrangian over the primal variables for fixed multipliers. It is always concave, and every value it takes is a lower bound on the primal minimum, for convex and nonconvex problems alike. Maximizing it gives the dual problem. This is the machinery behind Lagrangian-relaxation sizing, subgradient methods and the quality certificates that tell a heuristic how close to optimal it is.

## Evaluating the dual function

For minimize $f(x)$ subject to $g(x) \le 0$, the Lagrangian is $L(x, \lambda) = f(x) + \lambda g(x)$ with $\lambda \ge 0$, and the dual function is $q(\lambda) = \min_x L(x, \lambda)$ (written $g(\lambda)$ in many texts).

Take minimize $4x^2 + y^2$ subject to $x + y \ge 6$, a power model with a speed requirement. Write the constraint as $6 - x - y \le 0$, so $L = 4x^2 + y^2 + \lambda(6 - x - y)$. For fixed $\lambda$ the minimum over x and y is at $x = \lambda/8$, $y = \lambda/2$, and substituting gives
$$q(\lambda) = 6\lambda - \left(\tfrac{1}{16} + \tfrac{1}{4}\right)\lambda^2 = 6\lambda - 0.3125\lambda^2.$$
At $\lambda = 2.5$: $q = 15 - 1.953 = 13.047$. Every feasible design costs at least 13.047. The inner minimization is unconstrained and separates by variable, which is the point: the hard coupling is gone.

## Maximizing the dual

The best lower bound is $\max_{\lambda \ge 0} q(\lambda)$. For $q(\lambda) = c\lambda - s\lambda^2$ the maximum is at $\lambda^* = c/(2s)$ with value $c^2/(4s)$. With $\min 3x^2 + 4y^2$ subject to $x + y \ge 5$: $s = 1/12 + 1/16 = 0.1458$, so $\lambda^* = 5/0.2917 = 17.14$ and $q^* = 25/0.5833 = 42.86$.

Check against the primal: the optimum puts $x = \lambda^*/6 = 2.857$ and $y = \lambda^*/8 = 2.143$, summing to 5, with cost $3 \cdot 8.163 + 4 \cdot 4.592 = 42.86$. The values agree: zero duality gap. That holds for convex problems that have a strictly feasible point (Slater's condition). For integer or other nonconvex problems the dual maximum can be strictly below the primal minimum, and that duality gap is a property of the problem, not a flaw in the solver.

## Bounding a heuristic solution

Weak duality, $q(\lambda) \le f(x)$ for every feasible x and every $\lambda \ge 0$, turns any dual value into a certificate. If a heuristic returns a feasible point with cost $f(x)$, it is within $f(x) - q(\lambda)$ of optimal, whatever $\lambda$ was used.

For $\min x^2 + 4y^2$ subject to $x + y \ge 9$ a heuristic returns $(4, 6)$, feasible with cost $16 + 144 = 160$. With $\lambda = 2.5$: $q = 9 \cdot 2.5 - (1/4 + 1/16) \cdot 6.25 = 22.5 - 1.953 = 20.55$, so the point is within 139.45 of optimal: a weak certificate because $\lambda$ was poor. At the best multiplier $\lambda^* = 14.4$ the bound rises to 64.8, which is the true optimum. Lagrangian-relaxation sizers report exactly this gap between the best feasible sizing and the best dual bound to show convergence.

## Going deeper

Boyd and Vandenberghe, "Convex Optimization" (2004), chapter 5 on duality; Bertsekas, "Nonlinear Programming" (3rd ed., 2016), chapter 6.
