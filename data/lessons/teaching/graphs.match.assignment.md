# Weighted assignment of cells to slots

When every pairing has a cost, matching becomes the assignment problem: place $n$ cells into $n$ slots, one each, at minimum total cost. Detailed placers solve it on windows of a few cells to reorder them optimally, and pin and track assigners solve it at larger scale. It is solvable exactly in polynomial time, so a tool designer should know when to use it and how much better it is than greedy choice.

## The assignment problem

Given an $n \times n$ cost matrix $C$, where $C_{ij}$ is the cost of putting cell $i$ in slot $j$, find a permutation $\pi$ minimizing $\sum_i C_{i\pi(i)}$. For three cells with costs (3, 1, 6), (1, 7, 7) and (1, 9, 3), the best assignment is cell 1 to slot 2, cell 2 to slot 1 and cell 3 to slot 3, for $1 + 1 + 3 = 5$.

There are $n!$ permutations, 24 for $n = 4$ and about 3.6 million for $n = 10$, so enumeration is only for tiny windows. In detailed placement, costs are the wirelength of each cell's nets with the cell at that slot, holding all other cells fixed, which is exact when no two window cells share a net and an approximation otherwise.

## Polynomial-time solution

The Hungarian algorithm solves the assignment problem in $O(n^3)$. It maintains dual prices on rows and columns so that reduced costs $C_{ij} - u_i - v_j$ stay non-negative, and grows a matching on zero-reduced-cost edges with augmenting paths, adjusting prices when stuck. At the end, the matching uses only zero reduced costs, which proves optimality.

Equivalently, the assignment problem is a min-cost perfect matching, a min-cost flow with unit capacities from a source through cells and slots to a sink. For four cells with costs (5, 6, 9, 1), (1, 4, 3, 8), (7, 5, 2, 5), (7, 2, 7, 9), the optimum takes the four cheapest compatible entries 1, 1, 2, 2 for a total of 6. In practice, solvers such as the Jonker-Volgenant algorithm are faster than the textbook version.

## The gap to greedy assignment

A greedy pass lets each cell in turn take its cheapest free slot. It runs in $O(n^2)$ but can be badly suboptimal: the first cell may take a slot that a later cell needed far more. Sometimes the greedy result is already optimal, as when each cell's cheapest slot is different; often it is not, and the gap grows with window size.

This is the case for exact assignment inside heuristics: a detailed placer that slides a window across a row and solves each window's assignment exactly improves wirelength monotonically, with each step provably optimal for its window. The window size trades quality against $O(n^3)$ cost; windows of a few to a few dozen cells are typical, and independent windows can be solved in parallel.

## Going deeper

Kuhn, "The Hungarian method for the assignment problem" (Naval Research Logistics Quarterly, 1955); Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter on detailed placement.
