# Equivalence checking as the golden test

When an agent proposes a new synthesis transformation, a script or a rule, the first question is never whether it improves QoR. It is whether it preserves function. Combinational equivalence checking answers that with a proof, so it is the golden test in any harness that lets models edit or extend a synthesis engine.

## A proof, not a sample

The harness runs every changed engine on a benchmark suite and checks each output network against its input with a miter and SAT. UNSAT on every benchmark means the transformation preserved function on those designs. Random simulation is not enough: it finds bugs quickly but cannot prove their absence, and synthesis bugs often live in rare corners. A satisfiable miter is a counterexample, a specific input on which the outputs differ. That is the most useful feedback the model can get, far better than a pass or fail bit, and it becomes a permanent regression test.

## Handling failures, timeouts and scope

One counterexample out of a hundred benchmarks proves a bug, so the change is rejected, not accepted at 99 percent. A timeout is not a pass: the change stays unproven until the check succeeds, for example by partitioning outputs, sweeping internal equivalences first, or using a stronger engine. Checking each transformation step separately gives small, fast miters and pinpoints the failing step, while an end-to-end check guards their composition. Combinational checking covers combinational passes only; a change that moves registers needs sequential equivalence, and one that uses environment don't cares must be checked modulo those constraints.

## Beyond correctness

Equivalence is necessary, not sufficient. A differential QoR check shows the change moves area and delay in the intended direction on the intended designs and nowhere else. Resource limits and determinism checks guard against run-time blowups and nondeterminism on unseen designs. The checker itself should be independent of the optimizer, built on a separate code path, so a shared bug cannot make both wrong in the same way. Together these make the harness loop safe: the model proposes, the tools verify, and a human reviews only changes that are already proven correct and shown useful.

## Going deeper

Mishchenko, Chatterjee, Brayton and Een, "Improvements to combinational equivalence checking", ICCAD (2006); the golden and differential testing concept in this curriculum's agentic domain.
