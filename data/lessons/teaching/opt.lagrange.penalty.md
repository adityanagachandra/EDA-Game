# Quadratic penalty methods

A penalty method replaces a hard constraint with a cost that grows with the violation, then increases the penalty weight until the violation is small. Analytical placement is a penalty method at heart: density is a penalty added to wirelength, and its weight rises over hundreds of iterations until cells are spread. The method is simple and robust but never exactly feasible at a finite weight, and large weights make the problem stiff.

## The penalized minimizer

Take a cell that wirelength pulls toward $x = a$ but a fence requires $x \le c$ with $a > c$. The quadratic-penalty problem is
$$\min_x\ (x - a)^2 + \rho \max(0, x - c)^2.$$
The minimizer lies beyond the fence, where the penalty is active, so setting $2(x - a) + 2\rho(x - c) = 0$ gives
$$x(\rho) = \frac{a + \rho c}{1 + \rho}.$$
With $a = 20$, $c = 13$ and $\rho = 50$: $x = (20 + 650)/51 = 13.137$ um. The penalized solution is a weighted average of where the objective wants to be and where the constraint allows, with the weight $\rho$ pulling toward feasibility. In a placer the same balance plays out between each cell's wirelength pull and the density gradient pushing it out of overfull bins.

## How the violation shrinks

The remaining violation is
$$x(\rho) - c = \frac{a - c}{1 + \rho}.$$
With $a = 42$, $c = 24$ and $\rho = 20$: $18/21 = 0.857$ um. Doubling the weight roughly halves the violation, so the method approaches feasibility only like $1/\rho$ and is never exactly feasible at finite $\rho$.

The penalty also estimates the multiplier. At the penalized minimizer $2\rho(x(\rho) - c) = 2(a - x(\rho))$, and as $\rho \to \infty$ this tends to $2(a - c)$, the true KKT multiplier of the fence constraint. That observation is the seed of the augmented Lagrangian method. In placement the overflow, the fraction of cell area in overfull bins, plays the role of the violation, and placers stop raising the density weight once overflow falls below a target such as 10%, leaving the remainder to legalization.

## Choosing the penalty weight

To bring the violation down to $\varepsilon$, solve $(a - c)/(1 + \rho) \le \varepsilon$: $\rho \ge (a - c)/\varepsilon - 1$. With $a = 46$, $c = 11$ and $\varepsilon = 0.1$ um, $\rho \ge 349$.

Large weights cost conditioning. The penalized problem's curvature grows with $\rho$, so a gradient method's stable step shrinks like $1/\rho$, and an ill-conditioned mix of a soft wirelength term and a stiff penalty slows convergence. Placers therefore raise the weight gradually, multiplying it by a factor slightly above 1 each iteration (ePlace adapts the factor to the change in wirelength), so each problem is warm-started from the last. The alternative is to keep $\rho$ moderate and add a multiplier estimate, which is the augmented Lagrangian method.

## Going deeper

Nocedal and Wright, "Numerical Optimization" (2nd ed., 2006), chapter 17 on penalty and augmented Lagrangian methods; Lu et al., "ePlace" (ACM TODAES, 2015), for the density-penalty schedule.
