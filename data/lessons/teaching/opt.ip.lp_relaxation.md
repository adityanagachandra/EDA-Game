# LP relaxation and the integrality gap

An integer program becomes an LP when its integrality requirements are dropped. That LP relaxation is the engine of every MILP solver: it bounds the optimum, guides branching and suggests roundings. In EDA, relaxations bound timing-repair selection, test compaction, global routing and floorplanning ILPs. The distance between the relaxed and the integer optimum, the integrality gap, measures how much work the solver still faces and how good a formulation is.

## Solving the relaxed knapsack

A repair pass choosing whole upsizes under an area budget is a 0-1 knapsack: maximize $\sum v_i x_i$ subject to $\sum w_i x_i \le W$, $x_i \in \{0, 1\}$, where $v_i$ is the TNS recovered and $w_i$ the area. Relax to $0 \le x_i \le 1$. The relaxed problem is solved greedily: sort by $v_i / w_i$, take items whole while they fit, then take the fraction of the next one that fills the budget.

Example: budget 14; upsize 1 recovers 23 ps for 7 units (3.29 ps per unit), upsize 3 recovers 7 ps for 5 units (1.4), upsize 2 recovers 5 ps for 12 units (0.42). Take upsize 1 (7 units left), then upsize 3 (2 units left), then $2/12$ of upsize 2: $23 + 7 + 5 \cdot 2/12 = 30.83$ ps. Only one variable ends fractional, a general property of LP vertices with one coupling constraint.

## The integer optimum

The integer optimum must pick whole items. For small instances enumerate the feasible subsets; for large ones a solver branches. In the example the subsets that fit in 14 units are $\{1, 3\}$ worth 30, $\{1\}$ worth 23, $\{2\}$ worth 5, $\{3\}$ worth 7, and the empty set. $\{1, 2\}$ needs 19 units and $\{2, 3\}$ needs 17, so both are infeasible. The integer optimum is 30 ps.

The relaxation bound, 30.83, is never below the integer optimum for a maximization, because every integer solution is also feasible for the LP. A bound is only useful if it is tight. Greedy by ratio is not optimal for the integer problem in general: an item with a lower ratio can pack the budget better, which is why rounding the LP is a heuristic and branch and bound is needed for proof.

## Integrality gap and formulation strength

The integrality gap of an instance is the relaxed optimum minus the integer optimum (or their ratio). Here it is $30.83 - 30 = 0.83$ ps. A small gap means the LP nearly describes the integer hull and branch and bound will prune quickly. A large gap means many nodes.

Formulations of the same problem can have very different gaps. Big-M disjunctions for floorplan non-overlap are notoriously weak: the LP can set binaries to tiny fractions and ignore overlap entirely. Adding valid inequalities (cuts), using tighter M values, or reformulating with more variables (for example, path or pattern variables in routing and covering) shrinks the gap. A useful habit when building an EDA ILP is to solve the relaxation of a small instance first and compare: the gap predicts how the model will scale long before the full run does.

## Going deeper

Wolsey, "Integer Programming" (2nd ed., 2020), chapters 1 and 2 on formulations and relaxations; Martello and Toth, "Knapsack Problems" (1990).
