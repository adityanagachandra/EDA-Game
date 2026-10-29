# First-UIP learned clauses

When propagation hits a conflict, the solver has many valid clauses it could learn. The first-UIP scheme picks one with a precise recipe: resolve backwards from the conflict until exactly one literal of the conflict level remains. It is the scheme used by essentially every CDCL solver, because the clause it produces is short, asserting after backjumping, and cheap to compute from the trail. Computing it by hand is the core skill of conflict analysis.

## Unique implication points

Look at the implication graph restricted to the conflict level. A unique implication point (UIP) is a node of that level through which every path from the level's decision to the conflict passes. The decision itself is always a UIP, the last one in this sense; the first UIP is the one closest to the conflict. Cutting the graph just after the first UIP puts that node and all the earlier-level literals feeding the conflict side on the reason side of the cut. The learned clause contains the negations of the literals on the reason side that have an edge across the cut. It has exactly one literal of the conflict level, the negated UIP, plus literals of earlier levels. That single conflict-level literal is what makes the clause asserting: after undoing the conflict level, every other literal is false and the UIP's negation is forced.

## Resolving back to the first UIP

The UIP is found without drawing the graph. Start with the falsified clause. While it contains more than one literal of the conflict level, take the one whose variable was assigned most recently and resolve the clause with that variable's reason clause. Example at level 2: decide $x_3 = 1$; $x_5 = 1$ by $(\neg x_3 \vee x_5)$; $x_6 = 0$ by $(\neg x_5 \vee \neg x_6)$; $x_7 = 1$ by $(\neg x_5 \vee x_7)$; and the clause $(x_6 \vee \neg x_7 \vee \neg x_1)$ is falsified, with $x_1 = 1$ from level 1. Two level-2 literals, $x_6$ and $x_7$; the later is $x_7$, so resolve with its reason: $(x_6 \vee \neg x_1 \vee \neg x_5)$. Still two level-2 literals; resolve on $x_6$ with $(\neg x_5 \vee \neg x_6)$: $(\neg x_1 \vee \neg x_5)$. Only $\neg x_5$ remains at level 2, so $x_5$ is the first UIP and $(\neg x_1 \vee \neg x_5)$ is learned.

## Why the first UIP

Other cuts are just as sound. Resolving one more step in the example, on $x_5$ with $(\neg x_3 \vee x_5)$, gives the last-UIP clause $(\neg x_1 \vee \neg x_3)$, which here equals the all-decision clause. The falsified clause itself is also implied, but it has two conflict-level literals, so after backjumping it is not unit and asserts nothing. Experiments comparing learning schemes found first UIP consistently best: it stops as early as possible, so the clause tends to be short and to mention implied literals that summarize many decisions, and the learned clauses generalize to more of the search space. The resolution count is bounded by the number of literals at the conflict level, so analysis is linear in the trail. Solvers then shrink the clause further by clause minimization, which removes literals implied by others already in it.

## Going deeper

Zhang, Madigan, Moskewicz and Malik, "Efficient Conflict Driven Learning in a Boolean Satisfiability Solver" (ICCAD 2001); Biere et al. (eds.), "Handbook of Satisfiability" (2nd ed., 2021), the CDCL chapter.
