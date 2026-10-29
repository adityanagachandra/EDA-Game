# Big-M and disjunctive constraints

Many EDA rules are either-or: block i is left of, right of, below or above block j; a cell sits in one of two regions; a wire uses one of two layers. Integer programs encode such disjunctions with binaries and a large constant M that switches a constraint off. Big-M models are easy to write and are how ILP floorplanners and legalizers are usually stated, but a careless M ruins the relaxation, so choosing the smallest valid M is a real skill.

## Choosing the smallest valid M

The constraint $x_i + w_i \le x_j + Mp$ forces block i left of block j when $p = 0$ and must be vacuous when $p = 1$. Vacuous means it holds for every feasible placement, so M must cover the worst case: $M \ge \max (x_i + w_i - x_j)$ over the allowed ranges. That maximum is at the largest $x_i$ and the smallest $x_j$.

Example: block i is 21 um wide with fence $0 \le x_i \le 70$, and block j's fence requires $x_j \ge 2$. Then $M \ge 70 + 21 - 2 = 89$ um. Any larger M is still valid, but each extra unit loosens the LP relaxation. With no fences, on a chip of width W, the bound is W. Deriving M from the actual variable bounds, region by region, is cheap and often halves solve times.

## Counting binaries in disjunctive models

Pairwise non-overlap for n rectangles needs one disjunction per pair. A standard encoding uses two binaries $(p_{ij}, q_{ij})$ per pair to pick one of four relations, left-of, right-of, below or above, with four big-M constraints of which three are switched off. There are $\binom{n}{2}$ pairs, so $2\binom{n}{2} = n(n-1)$ binaries.

For 8 blocks that is 28 pairs and 56 binaries; with a rotation binary per block, $56 + 8 = 64$. For 50 blocks it is 2450 binaries, plus 4900 big-M constraints. The count grows quadratically, and every pair's binaries interact through the coordinates, so the tree grows fast. Exact ILP floorplanning is therefore practical only for small block counts, and larger floorplanners use sequence pairs or Polish expressions with annealing, or fix relative orders from a constraint graph and solve an LP for the coordinates.

## Why big-M relaxations are weak

In the LP relaxation $p$ may be fractional, and a small fraction of a large M buys a large violation. Take $w_i = 38$ um, $M = 350$ um, $x_i = 50$, $x_j = 80$: the blocks overlap by $50 + 38 - 80 = 8$ um, yet $p = 8/350 = 0.023$ satisfies the relaxed constraint. The LP pays almost nothing to ignore non-overlap, so its bound is close to the bound with no non-overlap constraints at all.

The remedies are the standard ones: the smallest valid M, symmetry breaking (fix the relative order of identical blocks), valid inequalities that link the binaries to the block dimensions and the die size, and, best of all, a different formulation. Indicator constraints, which solvers branch on directly, avoid M entirely. Fixing the relative-position binaries from a heuristic and then solving the LP for coordinates is how many legalizers use this model in practice.

## Going deeper

Sutanthavibul, Shragowitz and Rosen, "An analytical approach to floorplan design and optimization" (IEEE TCAD, 1991), the classic MILP floorplanning model; Wolsey, "Integer Programming" (2nd ed., 2020), chapter 1 on formulations.
