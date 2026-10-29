# Integer programming and branch and bound

Branch and bound solves integer programs exactly by building a tree: each node is a subproblem with some variables fixed, its relaxation gives a bound, and any node whose bound cannot beat the best known solution is pruned. It is the core of every MILP solver used for EDA ILPs, and the same bound-and-prune idea drives exact covering, ATPG search (PODEM backtracks over input assignments) and exact floorplanning. Understanding pruning and gaps is how you read a solver log.

## Pruning by bound

Keep an incumbent, the best feasible integer solution found so far. Each open node carries a bound from its relaxation. For maximization the bound is an upper bound on any solution in that subtree; if it is at most the incumbent, the subtree cannot contain a strictly better solution and is pruned. For minimization the test flips: prune when the lower bound is at least the incumbent.

Example: a covering search for a minimum test set has an incumbent of 48 patterns and open nodes with lower bounds 50, 51, 52, 52, 46 and 49. Every node with bound at least 48 is pruned: 50, 51, 52, 52 and 49, five nodes. Only the node with bound 46 might hold a 46- or 47-pattern cover. A node is also pruned when its relaxation is infeasible, and when its relaxed solution is integral, in which case it may become the new incumbent.

## Branching on a fractional variable

When a node's relaxation has a fractional variable $x_k = f$, branch: one child adds $x_k \le \lfloor f \rfloor$, the other $x_k \ge \lceil f \rceil$; for binaries, $x_k = 0$ and $x_k = 1$. Each child's relaxation is solved, and its bound is never better than the parent's.

In the knapsack example (budget 14; upsizes worth 23, 5, 7 ps for 7, 12, 5 units) the root LP is 30.83 with $x_2 = 1/6$ fractional. Child $x_2 = 0$: the greedy LP over items 1 and 3 takes both whole, bound 30, and the solution is integral, so it becomes the incumbent. Child $x_2 = 1$: 12 units used, 2 left, take $2/7$ of upsize 1: $5 + 23 \cdot 2/7 = 11.57$. That is below the incumbent 30, so it is pruned. The search ends after three nodes. Branching on the most fractional variable is common, but strong branching, which tests several candidates, usually shrinks the tree more.

## Reading the optimality gap

At any moment the best bound over all open nodes and the incumbent bracket the optimum. Solvers report the relative gap, commonly $(\text{bound} - \text{incumbent})/|\text{incumbent}|$ for maximization (definitions vary slightly between solvers). A MILP global router with incumbent 547 and best bound 594 has gap $47/547 = 8.59\%$: no solution can be more than 8.59% better than the one in hand.

The gap is the certificate EDA flows use to stop early: a floorplanner may accept 1% and save most of the run time, since closing the last fraction of a gap often costs more than finding the first good solution. Two forces close it: better incumbents from heuristics (diving, rounding, local search) and better bounds from cuts and branching. A gap that stalls with the incumbent fixed is a sign of a weak relaxation, not of a bad heuristic.

## Going deeper

Wolsey, "Integer Programming" (2nd ed., 2020), chapter 7; Achterberg, "Constraint Integer Programming" (PhD thesis, TU Berlin, 2007), the design of the SCIP solver; Goel, "An implicit enumeration algorithm to generate tests for combinational logic circuits" (IEEE Transactions on Computers, 1981), for PODEM.
