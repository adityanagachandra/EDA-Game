# Counterexample-guided agent repair

When a property fails, the prover returns a counterexample trace, which is unusually good input for an agent: short, precise and reproducible. An agentic repair loop reads the trace, localizes the cause, proposes a fix to the RTL, the environment model or the property, and re-runs the proofs. The engineering problem is making that loop trustworthy: the checks must be fixed, the edits tiered, and weakening the specification must never count as a repair.

## Reading the counterexample

BMC deepens one cycle at a time, so its first counterexample is a shortest one, with few irrelevant cycles. The agent first extracts the failing cycle and the signals in the assertion's cone of influence leading up to it, which turns a large waveform into a handful of events. Then it classifies the trace. A base-case trace from reset is a real failure of the design or of the environment model. A trace that drives inputs the interface specification forbids points to a missing assumption, not a design bug. An induction-step trace starts from an arbitrary state satisfying the assertions, so it may begin in an unreachable state and call for a helper invariant rather than an RTL change. Getting this classification right decides which artifact the fix should touch, and it is checkable: BMC aimed at the trace's start state tells whether that state is reachable.

## Proposing and checking a fix

The agent proposes a small edit and the harness re-runs the checks. The target assertion must now hold, and every previously passing assertion and cover must still pass: the whole property suite is the differential check against regressions, just as golden tests are for engine changes. A new counterexample is fed back as feedback for the next attempt. Covers matter here: a fix that makes a handshake assertion pass by stalling the interface forever satisfies every safety property and breaks the design, and a cover for normal traffic catches it. Bounded passes are not the end: if the original bug appeared at depth 14, BMC to depth 30 is evidence, not proof, so the loop should finish with an unbounded proof by k-induction or IC3 where one is achievable, and report the strength of the evidence when it is not.

## Guardrails on the repair loop

Some edits make failures disappear without fixing anything. Widening req |-> ##[1:2] ack to req |-> ##[1:200] ack, deleting an assertion or adding an assume that rules out the counterexample's inputs all weaken the specification. An assumption is legitimate only if the interface specification really forbids the behavior and the assumption is proved on the driving block. So the harness keeps the property suite, the covers and the assumptions under human control, outside the agent's edit scope, and runs in a sandbox that the agent cannot use to change its own oracle. Edits are tiered: environment constraints and helper invariants are small, reviewable changes, while RTL edits change the product and require a human-approved diff after the proofs pass. Every accepted repair is versioned and reversible, and the final proof artifacts, traces and invariants, are kept so that a reviewer can re-check them.

## Going deeper

Biere et al., "Bounded Model Checking" (Advances in Computers, 2003), on counterexample traces; the SymbiYosys documentation, on base-case and induction traces.
