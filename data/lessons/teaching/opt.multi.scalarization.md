# Weighted sums, epsilon constraints and hypervolume

To compute a single answer from a multi-objective problem, tools scalarize: combine the objectives with weights, or optimize one objective with the others bounded. Each method reaches different parts of the Pareto front. Weighted sums are simple and are what most flow scripts do, but they miss nonconvex parts of the front; epsilon constraints reach every Pareto point; hypervolume measures the quality of a whole front and is the standard yardstick for exploration methods.

## Weighted sums

A weighted sum minimizes $s = w f_1 + (1 - w) f_2$ with $w \in (0, 1)$. With D1 (delay 10, power 4), D2 (11, 5), D3 (11, 4) and D4 (8, 9) and $w = 0.1$, the scores are 4.6, 5.6, 4.7 and 8.9, so D1 wins at 4.6. With $w = 0.9$ they become 9.4, 10.4, 10.3 and 8.1, so D4 wins.

With positive weights the winner is always Pareto-optimal, since a dominated design scores worse than the design dominating it. But a weighted sum can only reach points on the convex hull of the front. A design in a dent, Pareto-optimal but above the line joining its neighbours, is never the minimum for any weight. Objectives must also be normalized before weighting, typically by baseline values, or the larger-scale objective dominates the score.

## Epsilon constraints

The epsilon-constraint method minimizes one objective subject to bounds on the others: minimize delay subject to power at most $\varepsilon$. With D1 (2, 5), D2 (10, 11), D3 (11, 4), D4 (4, 5), D5 (11, 10), D6 (10, 9) and $\varepsilon = 9$, the designs within the power budget are D1, D3, D4 and D6, and the smallest delay among them is 2 (D1).

Sweeping $\varepsilon$ traces the whole Pareto front, including points in nonconvex dents that no weight reaches. It also matches how engineers think: meet a power budget, then go as fast as possible. A caution: the optimum of an epsilon-constraint problem can be only weakly Pareto-optimal when there are ties in the minimized objective, so a second pass that minimizes the constrained objective among the tied designs is used to obtain a strictly Pareto-optimal answer.

## Hypervolume of a front

Hypervolume measures how much of objective space a set of designs dominates, up to a reference point R worse than every design. In two objectives it is the area of the union of rectangles $[f_1, R_1] \times [f_2, R_2]$. Sort the nondominated designs by $f_1$ (so $f_2$ decreases) and add
$$HV = \sum_i (R_1 - f_{1,i})(f_{2,i-1} - f_{2,i}), \qquad f_{2,0} = R_2.$$
Example with $R = (13, 13)$ and front (3, 8), (4, 1): the first adds $10 \cdot 5 = 50$, the second $9 \cdot 7 = 63$, total 113. Dominated designs add nothing.

A larger hypervolume means a front that is both closer to the ideal and better spread. It is one of the very few unary quality indicators that is strictly monotone with Pareto dominance, which is why exploration and multi-objective BO methods use it, and why its improvement serves as an acquisition function. The reference point must be fixed in advance, since moving it changes comparisons.

## Going deeper

Miettinen, "Nonlinear Multiobjective Optimization" (1999), chapters on the weighting and epsilon-constraint methods; Zitzler, Thiele, Laumanns, Fonseca and da Fonseca, "Performance assessment of multiobjective optimizers: an analysis and review" (IEEE Transactions on Evolutionary Computation, 2003).
