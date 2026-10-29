# Unit propagation

Unit propagation, also called Boolean constraint propagation, is the inference step that every SAT solver runs after every assignment. When all literals of a clause but one are false, that last literal must be true. Repeating this until nothing changes either extends the assignment or exposes a conflict. Solvers spend most of their run time here, so its rules and its cost model shape the design of the whole solver.

## Unit clauses

Under a partial assignment, a clause is satisfied if some literal is true, conflicting if every literal is false, unit if every literal but one is false and that one is unassigned, and unresolved otherwise. Only unit clauses force anything. Take $(x_1 \vee \neg x_2 \vee x_3)$ with $x_1 = 0$ and $x_3 = 0$: two literals are false and $\neg x_2$ is unassigned, so the clause is unit and forces $x_2 = 0$. With only $x_1 = 0$ it is unresolved, because two literals are still free. With $x_3 = 1$ it is satisfied and plays no further role until backtracking. To count the unit clauses of a CNF under an assignment, scan each clause, skip it if any literal is true, and mark it unit when exactly one literal is unassigned.

## Propagating to a fixpoint

Propagation repeats: set the forced literal of every unit clause, which may make new clauses unit, until no unit clause remains or some clause is conflicting. Take $(x_1)$, $(\neg x_1 \vee x_2)$, $(\neg x_2 \vee \neg x_3)$, $(x_3 \vee x_4 \vee x_5)$ and $(\neg x_4 \vee x_5)$. The unit clause sets $x_1 = 1$; then $(\neg x_1 \vee x_2)$ forces $x_2 = 1$; then $(\neg x_2 \vee \neg x_3)$ forces $x_3 = 0$. Now $(x_3 \vee x_4 \vee x_5)$ has two free literals and $(\neg x_4 \vee x_5)$ has two, so the fixpoint assigns three variables. The fixpoint does not depend on the order in which unit clauses are processed, as long as no conflict arises; if one does, the solver stops and analyzes it. Propagation alone decides satisfiability for Horn formulas, but in general it is only a cheap, sound inference that leaves the rest to search.

## Propagation after a decision

Inside search, propagation runs after every decision, and the literals it sets are the decision's implications. Decide $x_1 = 0$ on $(x_1 \vee x_2)$, $(\neg x_2 \vee x_3)$, $(x_1 \vee \neg x_3 \vee x_4)$ and $(\neg x_4 \vee \neg x_5)$. The first clause forces $x_2 = 1$, the second $x_3 = 1$, the third (with $x_1 = 0$ and $x_3 = 1$) forces $x_4 = 1$, and the last forces $x_5 = 0$: four implied variables from one decision. This chain is what the implication graph records, and the reason clause of each implied literal is the clause that became unit. Because propagation dominates run time, solvers avoid rescanning every clause: the two-watched-literal scheme visits a clause only when one of two chosen literals becomes false.

## Going deeper

Biere et al. (eds.), "Handbook of Satisfiability" (2nd ed., 2021), the chapter on CDCL solvers; Moskewicz, Madigan, Zhao, Zhang and Malik, "Chaff: Engineering an Efficient SAT Solver" (DAC 2001).
