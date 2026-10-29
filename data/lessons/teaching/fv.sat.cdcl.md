# Conflict-driven clause learning

Conflict-driven clause learning turned SAT from a theoretical hardness example into an industrial engine that routinely solves formulas with millions of clauses from hardware verification. CDCL keeps the DPLL skeleton of decide and propagate, but on every conflict it analyzes the implication graph, learns a new clause, and jumps back past irrelevant decisions. Knowing the loop, why the learned clauses are sound, and how solvers keep them under control is the base for every SAT-based tool.

## The CDCL loop

A CDCL solver repeats one loop. Decide: pick an unassigned variable and a value, opening a new decision level. Propagate: run unit propagation to a fixpoint. If no conflict arises and every variable is assigned, report SAT with the current assignment. If a conflict arises at level 0, report UNSAT, because level-0 literals are forced by the formula alone. Otherwise analyze the conflict: walk the implication graph back to the first unique implication point and derive a learned clause. Add the clause, backjump to the second-highest decision level in it, and let propagation use it: after the backjump the learned clause is unit, so it immediately asserts the negation of the UIP literal. Restarts occasionally clear the trail back to level 0 without forgetting learned clauses. Each piece (decision heuristic, analysis, backjump, restart, deletion) is its own concept in this subdomain.

## Why learned clauses are sound and useful

Conflict analysis is a sequence of resolution steps. Starting from the falsified clause, it resolves with the reason clause of a literal assigned at the conflict level, and repeats. Every resolvent is implied by the clauses it came from, so the learned clause is implied by the original formula: adding it never removes a model, and the answer cannot change. What it changes is the search. DPLL that fails below $x_1 = 1$ and $x_4 = 0$ will happily make the same assignment again in another branch and rediscover the same conflict. A learned clause such as $(\neg x_1 \vee x_4)$ blocks that combination everywhere in the remaining search. Because every learned clause is a resolution consequence, the sequence of learned clauses ending in the empty clause is a proof of unsatisfiability that an independent checker can verify by reverse unit propagation, which is what the DRAT proof format records.

## Managing the learned-clause database

A solver can learn a clause per conflict and run millions of conflicts, so learned clauses would soon dominate memory and slow propagation, since every clause sits on two watch lists. Solvers therefore delete learned clauses periodically. They can delete only learned clauses: those are redundant consequences, whereas deleting an original clause could change the problem. To choose what to keep, modern solvers measure the literal block distance (LBD), the number of distinct decision levels among a clause's literals when it was learned; clauses with small LBD, often called glue clauses, tend to keep propagating and are kept, while others are kept only if recently active. A solver asked for an auditable UNSAT answer also logs its learned clauses and deletions as a DRAT proof, so that a small, separately trusted checker can confirm the result without trusting the solver.

## Going deeper

Marques-Silva, Lynce and Malik, "Conflict-Driven Clause Learning SAT Solvers", in Biere et al. (eds.), "Handbook of Satisfiability" (2nd ed., 2021); Audemard and Simon, "Predicting Learnt Clauses Quality in Modern SAT Solvers" (IJCAI 2009); Een and Sorensson, "An Extensible SAT-solver" (SAT 2003), the MiniSat paper.
