# Seams and certificates in a formal engine

An agent that customizes or extends a formal engine needs to know which parts it may touch. The four-layer model, an immutable core, pluggable policies, parameters and observation hooks, maps cleanly onto SAT solvers and model checkers, and formal engines have an advantage most engines lack: their answers come with certificates that an independent checker can verify. Those certificates are the golden tests that make agent-driven engine changes safe to promote.

## The four layers of a formal engine

In a CDCL SAT solver the immutable core is what soundness depends on: unit propagation, resolution-based conflict analysis, and proof logging. The policy layer holds strategies that cannot affect correctness: the decision heuristic (VSIDS or a variant), the restart schedule (Luby, geometric or dynamic), the clause-deletion policy and phase selection. Parameters tune those policies: decay factors, the restart unit, deletion thresholds. Observation hooks report progress, such as conflicts per second, learned-clause LBD and trail depth. In an IC3 engine the core is the relative-inductiveness query, the frame invariants and the convergence check; the generalization strategy (plain MIC, ctgDown) and obligation ordering are policies; depth limits and time budgets are parameters; frame sizes and clause counts are hooks that let an agent decide whether a run is converging or should be stopped and handed to another engine.

## Which tier a request lands in

Each customer request lands at the smallest tier that can satisfy it. Raising a block's BMC depth from 30 to 60, or changing the restart unit, is Tier 1: configuration, no recompile. A new decision heuristic that prefers variables in the current property's cone of influence is Tier 2: a plug-in at the decision-policy seam, leaving the sound core untouched. A change to how learned clauses are derived is Tier 3: a core edit that touches the correctness argument, so it needs a human-approved diff and full regression, because a mistake there produces wrong UNSAT answers. Customer-specific front-end behavior, such as translating SVA into monitors in a house style, belongs in a plug-in at the input boundary that produces the netlist-plus-property interface every engine already consumes. A seam is well placed when requests of a given kind essentially never need Tier 3.

## Certificates as golden checks

Formal answers can be validated without trusting the engine that produced them. An UNSAT answer from a SAT solver comes with a DRAT proof, checked by an independent proof checker; a counterexample from a model checker is replayed in a separate simulator to confirm that the assertion really fails; an IC3 proof returns an inductive invariant, re-checked with three SAT queries for initiation, consecution and implication of the property. These checks are the golden tests for an engine change: run the regression suite, validate every answer's certificate, and compare run time and solved-instance counts against the baseline. Differential testing between two solvers adds a further check: when they disagree on an instance, validating the SAT model by evaluation and the UNSAT answer by its proof decides which engine is wrong. Changes run sandboxed and are promoted only after every certificate checks.

## Going deeper

Wetzler, Heule and Hunt, "DRAT-trim: Efficient Checking and Trimming Using Expressive Clausal Proofs" (SAT 2014); Brummayer, Lonsing and Biere, "Automated Testing and Debugging of SAT and QBF Solvers" (SAT 2010).
