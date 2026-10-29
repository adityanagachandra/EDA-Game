# Global routing as an integer program

Global routing has an exact formulation as an integer linear program: pick one candidate tree per net so that no edge exceeds its capacity and total cost is minimal. Real instances are far too large to solve exactly, but the formulation is the reference that negotiation, Lagrangian methods and flow relaxations approximate, and it is how you reason about what an optimal router would do.

## Variables: one per candidate tree

For each net $i$, generate a set of candidate Steiner trees $T_{i1}, \dots, T_{ik_i}$ (for example, pattern routes of a FLUTE topology). Introduce a binary variable $x_{ij} = 1$ if net $i$ uses tree $j$. The program has $\sum_i k_i$ variables. Example: six nets with 3, 6, 6, 5, 2 and 5 candidates give 27 variables.

The candidate set is a modelling choice. With all possible trees the program is exact but has exponentially many variables; with a few candidates per net it is small but can only pick among them. Column generation adds candidate trees on demand, guided by dual prices from the LP relaxation, which is how exact-leaning routers keep the set small.

## Constraints and objective

Two families of constraints define a legal choice. Each net picks exactly one tree: $\sum_j x_{ij} = 1$ for every net $i$. Each edge respects capacity: $\sum_{i,j : e \in T_{ij}} x_{ij} \le c_e$ for every edge $e$. The objective minimises total cost, usually wirelength: $\min \sum_{i,j} w_{ij} x_{ij}$.

So there are as many selection constraints as nets and as many capacity constraints as edges. Example: 316 nets on a 20 by 12 gcell grid give $316 + 448 = 764$ constraints. A softer version replaces hard capacity with overflow variables penalised in the objective, which is always feasible and lets the solver report where capacity must be violated.

## Solving a tiny instance

For three nets with two candidates each, there are $2^3 = 8$ combinations. Check each against every capacity, discard the infeasible ones, and take the cheapest of the rest. That brute force is what branch and bound does cleverly: it relaxes $x_{ij} \in \{0, 1\}$ to $0 \le x_{ij} \le 1$, solves the LP, and branches on fractional variables, pruning subtrees whose LP bound exceeds the best feasible cost.

The program is NP-hard in general, and full-chip instances have millions of variables, so production routers do not solve it. They use it as a reference: negotiated congestion behaves like Lagrangian relaxation of the capacity constraints, LP relaxation with randomised rounding gives provable bounds, and ILP on small windows is used for local refinement.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); P. Raghavan and C. D. Thompson, "Randomized Rounding: A Technique for Provably Good Algorithms and Algorithmic Proofs" (1987).
