# Change tiers in a Boolean reasoning engine

When an agent or a user asks for a change in a SAT solver or BDD package, the first decision is where the change belongs. The four-layer model, immutable core, policy, parameters and observation hooks, says what can be changed cheaply and what needs human approval and full regression. Boolean engines are an ideal case study because their cores are small and correctness-critical (a wrong UNSAT is a false proof) while their behavior is dominated by tunable heuristics.

## Core, policy, parameters and hooks

In a CDCL solver the core is clause semantics, unit propagation, conflict analysis correctness and proof logging: anything whose bug produces a wrong answer. Policy is the decision heuristic, the restart strategy and the clause-deletion scheme. Parameters are activity decay factors, restart intervals and conflict limits. Hooks report conflict counts, learned-clause statistics and trail depth to whoever is watching.

In a BDD package the core is the unique table and the canonicity invariants: no duplicate nodes, no node whose children are equal. Variable reordering methods, sifting or window permutation, are policy; node limits, cache sizes and reordering thresholds are parameters; a callback on garbage collection is a hook.

## Landing a request in a tier

Try the smallest tier first. "Restart more often on our bounded model checking jobs" is Tier 1, a parameter change with no recompile. "Use a variable order that interleaves the bits of the two operand buses", when the package offers only static and sifting orders, is Tier 2, a new policy at an existing seam. "Store clauses differently to support a new proof format" touches clause storage and proof logging, the core, so it is Tier 3: a human-approved diff and the full regression suite.

Passing golden tests is necessary, not sufficient. A plug-in heuristic that gives correct answers but runs 40% slower must also be compared against the baseline on the benchmark suite before it is promoted. Correctness is checked by golden tests; quality is decided by differential QoR and runtime.

## Seams for Boolean engines

A seam is a narrow interface the core calls at a well-defined point. A decision-heuristic seam in a SAT solver exposes read-only search state, variable activities, the assignment trail and the decision level, and returns a literal; it never gets write access to the clause database, which would let a plug-in break soundness. A variable-order seam in a BDD package accepts and returns a permutation, while the core rebuilds the diagram and keeps it canonical.

Seams must be coarse and deterministic. Propagation is the hot loop, executed millions of times per second, so a plug-in call there would block inlining and dominate runtime; calls belong at decisions, restarts or reorderings. A plug-in that seeds a random generator from the clock breaks reproducibility, and with it every golden and differential test.

## Going deeper

Een and Sorensson, "An extensible SAT-solver", SAT (2003), on MiniSat's design. Somenzi, CUDD: CU Decision Diagram Package documentation, on reordering and the unique table.
