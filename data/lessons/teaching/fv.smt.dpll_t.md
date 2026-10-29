# DPLL(T)

Satisfiability modulo theories (SMT) solvers decide formulas that mix Boolean structure with atoms from theories such as bit-vectors, arrays, uninterpreted functions and arithmetic. The dominant architecture, DPLL(T), pairs a CDCL SAT solver that handles the Boolean structure with theory solvers that check whether the chosen atoms are consistent. Understanding the split explains both SMT performance and the interfaces that hardware model checkers rely on.

## Boolean abstraction

DPLL(T) replaces each theory atom by a Boolean variable. The formula $(x < 3 \vee y = 5) \wedge x > 4$ becomes $(b_1 \vee b_2) \wedge b_3$, with $b_1$ standing for $x < 3$, $b_2$ for $y = 5$ and $b_3$ for $x > 4$. The SAT solver searches over this abstraction as usual: decide, propagate, learn. A Boolean model is only a candidate, because the abstraction forgets what the atoms mean: $b_1 = b_3 = 1$ satisfies the Boolean formula but asks for $x < 3$ and $x > 4$ at once. So every candidate, or every partial assignment in an eager-checking solver, is passed to the theory solver. When the theory solver agrees, it extends the Boolean model to concrete values, such as $x = 5$ and $y = 5$, which is the model the SMT solver returns.

## Theory checks and explanations

The theory solver checks whether the asserted literals are consistent in the theory. If not, it returns an explanation: a small subset of the literals that is already inconsistent, such as $\{x < 3, x > 4\}$. Its negation, $(\neg b_1 \vee \neg b_3)$, is a valid clause, and the SAT solver learns it exactly like a conflict clause, then backjumps. Small explanations matter, because a short clause blocks every assignment containing those literals, not just the current one. Theory solvers can also propagate: with $x = y$ and $y = z$ asserted, the EUF solver tells the SAT solver that $x = z$ must hold before it is decided, giving a reason clause, just as unit propagation does. Because the SAT solver backtracks constantly, a theory solver must also retract assertions cheaply, keeping its own trail of changes.

## Eager versus lazy and theory combination

DPLL(T) is the lazy approach: the theory is consulted during search. The eager approach translates the whole problem to SAT first, as bit-blasting does for bit-vectors. Bit-vector solvers are largely eager, after aggressive word-level rewriting, because a bit-level SAT solver is the most effective engine for fixed-width arithmetic; arithmetic over unbounded integers and EUF are handled lazily. Formulas mixing theories, such as integer arithmetic with uninterpreted functions, are solved by Nelson-Oppen combination: each theory solver works on its own atoms, and they exchange equalities between shared variables until both agree or one finds a conflict, which is complete for stably infinite theories with disjoint signatures. For hardware model checking the decisive feature is incrementality, push and pop or assumptions, since BMC and IC3 issue many small related queries.

## Going deeper

Nieuwenhuis, Oliveras and Tinelli, "Solving SAT and SAT Modulo Theories: From an Abstract Davis-Putnam-Logemann-Loveland Procedure to DPLL(T)" (Journal of the ACM, 2006); Kroening and Strichman, "Decision Procedures" (2nd ed., 2016), chapters 3 and 10.
