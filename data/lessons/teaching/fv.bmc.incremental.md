# Incremental SAT across BMC depths

A BMC run solves a sequence of closely related SAT problems, one per depth, each containing the previous one. Rebuilding the solver every time throws away everything it learned. Incremental SAT keeps one solver alive, adds clauses as the bound grows and uses assumptions for the parts that change. It is how production BMC engines run, and the same interface underlies k-induction and IC3.

## Reusing the solver across depths

The depth-$(k+1)$ problem contains every clause of the depth-$k$ problem except the property check: the reset clauses and the transition clauses of frames 0 to $k$ are unchanged, and one new frame is added. So an incremental engine keeps one solver, adds the new frame's clauses, and moves the bad-state check to the new frame. Clauses learned while solving depth $k$ were derived from clauses that are still present, so they remain valid and keep pruning the search at depth $k + 1$; VSIDS scores and saved phases carry over as well, which is always sound because they are only heuristics. The only requirement is that nothing specific to the old query, such as \"bad at frame $k$\", was added as a permanent clause, because a lemma derived from it would be false for the new query.

## Assumptions and activation literals

Solvers such as MiniSat accept assumptions: literals that hold for one solve call only, treated as decisions at the lowest levels, so no learned clause depends on them unconditionally. BMC passes \"bad at frame $k$\" as an assumption. An equivalent pattern guards a clause with an activation literal: add $(\neg a_k \vee \text{bad}_k)$, assume $a_k$ to enable the check, and later assume or add $\neg a_k$, which satisfies the guarded clause and disables it permanently. When a query is UNSAT under assumptions, the solver can return the subset of assumptions its proof used, an unsat core over the assumptions. It tells the engine which frames or properties the proof needed, which abstraction and IC3 engines exploit to decide what to keep.

## What carries over between depths

Incrementality changes speed, never answers: rebuilding from scratch at every depth gives the same SAT and UNSAT results, only more slowly. Results must also be read the same way: UNSAT at depth $k$ means no counterexample of length up to $k$, not that the property holds forever. Learned clauses over the transition relation are reusable across depths and even across properties of the same design, which is why multi-property engines share one solver. Clauses learned under an assumption are kept only in a conditional form, with the negated assumption literal inside them. Over a long run, learned clauses still have to be reduced like in any CDCL solver. Open-source engines that work this way include ABC's bmc3 and the smtbmc engine used by SymbiYosys, which keeps an SMT solver session across depths.

## Going deeper

Een and Sorensson, "Temporal Induction by Incremental SAT Solving" (BMC Workshop, 2003); Een and Sorensson, "An Extensible SAT-solver" (SAT 2003), the section on incremental solving with assumptions.
