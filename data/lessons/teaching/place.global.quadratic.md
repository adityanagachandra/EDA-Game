# Quadratic placement

Quadratic placement replaces each net by two-pin springs and minimises the sum of squared spring lengths. Because the objective is a convex quadratic, its minimum satisfies a linear system that can be solved in near-linear time. GORDIAN, FastPlace, Kraftwerk2 and SimPL are all built on this core, with different net models and spreading methods.

## One movable cell

With two-pin connections of weight $w_j$ from a movable cell to fixed pins at $x_j$, the cost in x is $\sum_j w_j(x - x_j)^2$. Setting the derivative $2\sum_j w_j(x - x_j)$ to zero gives the weighted mean:
$$x^* = \frac{\sum_j w_j x_j}{\sum_j w_j}.$$
For pins at 74 (weight 1) and 1 (weight 3), $x^* = (74 + 3)/4 = 19.25$. The x and y problems separate, so $y^*$ is the weighted mean of the y coordinates. Every cell in a quadratic placement ends up at the weighted mean of everything it connects to, which is both its strength (a closed form) and its weakness (cells pile up in the middle).

## Two or three movable cells

With several movable cells, each cell's optimality condition involves its neighbours. For movable a and b connected with weight 1, a also tied to a fixed pin at 46 and b to one at 45 (weight 1 each), the conditions are $2x_a - x_b = 46$ and $-x_a + 2x_b = 45$, so $x_a = 45.67$ and $x_b = 45.33$. For a chain of three cells, the system is tridiagonal and can be solved by elimination in linear time. In general the conditions form $Ax = b$, where A collects connection weights and b collects weighted fixed-pin coordinates. Solving this system is the whole of one quadratic placement iteration.

## Why quadratic placement needs help

The quadratic solution minimises squared lengths, not HPWL, so it over-penalises long connections and pulls cells inward; the B2B net model corrects the mismatch by re-weighting edges. More seriously, the solution ignores density: highly connected cells cluster near the centre with massive overlap. Quadratic placers therefore alternate the solve with spreading. GORDIAN used partitioning-based constraints; FastPlace used cell shifting; Kraftwerk added density-derived forces; SimPL computes a spread-out target with look-ahead legalization and pulls cells toward it with anchor pseudo-nets. The solve stays cheap because A is sparse, symmetric and positive definite, so preconditioned conjugate gradient handles millions of cells.

## Going deeper

Kleinhans, Sigl, Johannes and Antreich, "GORDIAN: VLSI placement by quadratic programming and slicing optimization" (IEEE TCAD, 1991); Kim, Lee and Markov, "SimPL: an effective placement algorithm" (ICCAD 2010).
