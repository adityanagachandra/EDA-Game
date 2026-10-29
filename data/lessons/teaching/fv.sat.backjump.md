# Non-chronological backtracking

After learning a clause, a CDCL solver does not simply undo its last decision. It jumps back to the level at which the learned clause first becomes unit, often many levels earlier, discarding decisions that had nothing to do with the conflict. This non-chronological backtracking, or backjumping, is what lets CDCL escape regions of the search that DPLL would grind through. Computing the backjump level from a learned clause is quick, and it is an exact rule.

## The second-highest level

A first-UIP clause has exactly one literal at the conflict level $c$ and its other literals at lower levels. The solver backjumps to the highest of those lower levels, the second-highest level appearing in the clause. There, every literal except the UIP literal is still false and the UIP literal is unassigned, so the clause is unit and propagation asserts it immediately. Example: a conflict at level 4 learns a clause whose literals sit at levels 4, 1 and 2. The backjump level is 2: levels 3 and 4 are undone, and the clause forces its level-4 literal's negation at level 2. If the learned clause has a single literal, it is a unit clause true regardless of any decision, so the solver backjumps to level 0 and the literal becomes a permanent fact.

## From trail to backjump level

In practice the clause comes from conflict analysis, so the two steps chain: compute the first-UIP clause from the trail, then read off its levels. In the first-UIP example, the trail at level 2 produced $(\neg x_1 \vee \neg x_5)$, where $x_1$ was set at level 1 and $x_5$ at level 2. The second-highest level is 1, so the solver undoes level 2 and, at level 1, propagates $x_5 = 0$ from the learned clause. That propagation may cascade into further implications at level 1 or even an immediate new conflict, which is analyzed the same way. Note what happened to the decision $x_3 = 1$: it is gone, and the solver will choose its next decision afresh with the new clause in place, instead of mechanically trying $x_3 = 0$ as DPLL would.

## How much search is skipped

Chronological backtracking would return to level $c - 1$ and flip that decision. Backjumping to level $b$ undoes $c - b$ levels instead, so it skips $(c - 1) - b$ levels that chronological backtracking would have kept. A conflict at level 6 whose learned clause has literals at levels 6, 2 and 1 backjumps to level 2, skipping levels 5, 4 and 3: three levels whose decisions were irrelevant to this conflict. Those decisions may be remade later, and phase saving (re-using each variable's last value) makes the re-descent cheap when they were fine. Some solvers, on very long trails, use chronological backtracking selectively to avoid repeating expensive propagation, but the learned clause and the asserting level are computed the same way.

## Going deeper

Marques-Silva and Sakallah, "GRASP" (IEEE Transactions on Computers, 1999); Nadel and Ryvchin, "Chronological Backtracking" (SAT 2018).
