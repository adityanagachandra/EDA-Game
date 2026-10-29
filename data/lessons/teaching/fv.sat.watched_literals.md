# Two watched literals

Unit propagation dominates SAT run time, and the naive implementation rescans every clause that contains a newly falsified literal. The two-watched-literal scheme, introduced with the Chaff solver, watches only two literals per clause and does no work at all on backtracking. It is why modern solvers propagate millions of literals per second, and it is a model for any incremental engine that must react to local changes cheaply.

## Two watches per clause

A clause can become unit or conflicting only when all but at most one of its literals are false. So it is enough to watch two of its non-false literals: as long as both watched literals are non-false, the clause cannot be unit and need not be looked at. The solver keeps, for each literal, a watch list of the clauses watching it. When literal $\ell$ becomes false, the solver visits only the clauses on $\ell$'s watch list, not every clause containing $\ell$. Assigning $x_3 = 0$ falsifies the literal $x_3$, so the visited clauses are those watching $x_3$; clauses that contain $x_3$ but watch two other literals are untouched. To count visits from a watch table, list the clauses whose watch pair includes the falsified literal.

## Moving a watch

When a watched literal becomes false, the solver looks for a replacement: a literal of the clause that is not false and is not the other watch. If one exists, the watch moves there and the clause leaves this watch list. If none exists, the clause is unit (the other watch is unassigned and must be set true) or conflicting (the other watch is false too), and the watch stays. In this course the scan goes from the start of the clause. Example: $(x_1 \vee \neg x_2 \vee x_3 \vee x_4)$ watches $x_1$ and $\neg x_2$. When $x_1$ becomes false, the scan skips $x_1$ (false) and $\neg x_2$ (the other watch) and moves the watch to $x_3$: one move. When $x_4$ becomes false nothing happens, since it is not watched. When $x_3$ becomes false, no replacement remains, so the clause is unit and forces $\neg x_2$: no move. Total: 1 move.

## Why it beats occurrence lists

An occurrence-list solver visits every clause that contains the falsified literal; the watched scheme visits only the clauses watching it. If $\neg x_3$ occurs in 40 clauses but is watched in 6, assigning $x_3 = 1$ costs 6 visits instead of 40, and long clauses benefit most because only 2 of their literals are ever watched. The second advantage is backtracking: unassigning literals can only make watched literals non-false again, which never breaks the invariant, so the watch lists need no update when the solver backjumps. Counter-based schemes had to restore a counter per clause on every undo. Implementations also keep a blocking literal next to each watch-list entry, so a clause already satisfied by that literal is skipped without touching its memory, which matters for cache behaviour.

## Going deeper

Moskewicz, Madigan, Zhao, Zhang and Malik, "Chaff: Engineering an Efficient SAT Solver" (DAC 2001); Een and Sorensson, "An Extensible SAT-solver" (SAT 2003).
