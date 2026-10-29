# DPLL

DPLL, named for Davis, Putnam, Logemann and Loveland, is the backtracking search at the heart of every complete SAT solver. It picks a variable, tries a value, propagates, and backtracks on a conflict. Modern CDCL solvers add learning and non-chronological backtracking, but their skeleton is still DPLL. Tracing it by hand on small formulas is the best way to see why plain backtracking repeats work and what learning removes.

## Split, propagate, backtrack

DPLL is the Shannon expansion turned into search: $F$ is satisfiable if $F|_{x=0}$ or $F|_{x=1}$ is. The convention used here: after every assignment run unit propagation; report SAT as soon as every clause is satisfied; otherwise branch on the lowest-index unassigned variable, try 0 first, and on a conflict undo back to the most recent decision whose other value is untried (chronological backtracking). A decision is every branching assignment, including the second value tried after a backtrack; forced values are not decisions. On $(x_1 \vee x_2)$, $(x_1 \vee \neg x_2)$, $(\neg x_1 \vee x_3)$: decide $x_1 = 0$ (decision 1); the first clause forces $x_2 = 1$ and the second then conflicts. Backtrack and set $x_1 = 1$ (decision 2); propagation forces $x_3 = 1$ and every clause is satisfied. Two decisions, answer SAT.

## Counting decisions on a satisfiable formula

Take $C_1 = (x_1 \vee x_2 \vee x_3)$, $C_2 = (x_1 \vee \neg x_2)$, $C_3 = (\neg x_3 \vee x_4)$, $C_4 = (x_1 \vee x_2 \vee \neg x_4)$. Decide $x_1 = 0$ (decision 1). $C_2$ forces $x_2 = 0$, then $C_1$ forces $x_3 = 1$, $C_3$ forces $x_4 = 1$, and $C_4$ has every literal false: conflict. Backtrack: $x_1 = 1$ (decision 2) satisfies $C_1$, $C_2$ and $C_4$, but $C_3$ is unresolved and nothing is unit. Branch on $x_2 = 0$ (decision 3): $C_3$ is still open. Branch on $x_3 = 0$ (decision 4): $C_3$ is satisfied, so every clause is, and the solver reports SAT after 4 decisions. Two of those decisions were on variables that did not matter once $x_1 = 1$; real solvers waste effort the same way, which is one reason decision heuristics prefer variables in recent conflicts.

## Proving unsatisfiability

To answer UNSAT, DPLL must refute both values of every branch: the search tree is explored completely and every leaf ends in a conflict. With $(x_1 \vee x_2)$, $(x_1 \vee \neg x_2)$, $(\neg x_1 \vee x_2)$, $(\neg x_1 \vee \neg x_2)$: decide $x_1 = 0$ (decision 1); the first two clauses force $x_2 = 1$ and $x_2 = 0$, a conflict. Decide $x_1 = 1$ (decision 2); the last two clauses conflict the same way. With both values of the first variable refuted, the formula is UNSAT after 2 decisions. In general the tree can be exponential, and a DPLL refutation corresponds to a tree-like resolution proof, which is provably weaker than the general resolution proofs that clause learning can produce. That gap is why every competitive solver is CDCL.

## Going deeper

Davis, Logemann and Loveland, "A Machine Program for Theorem-Proving" (Communications of the ACM, 1962); Biere et al. (eds.), "Handbook of Satisfiability" (2nd ed., 2021), the chapters on complete algorithms and CDCL.
