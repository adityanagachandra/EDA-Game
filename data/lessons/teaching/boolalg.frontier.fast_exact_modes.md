# Fast and exact Boolean reasoning modes

Every Boolean engine has a fast mode and an exact mode. Random or guided simulation, structural hashing and cheap signatures run in microseconds and can only refute; SAT, BDDs and proof-producing checkers are complete but can take arbitrarily long. An agent running a continuous optimization loop must use the fast mode for almost everything and the exact mode exactly where a decision is committed. Getting the boundary wrong either wastes compute or ships a wrong circuit.

## Simulation as the fast mode

Simulating $t$ random vectors through an original circuit and a proposed rewrite costs a few word operations per gate per 64 vectors. A mismatch is a definite counterexample. A clean run is weak evidence: if the rewrite were wrong on $d$ of $2^n$ inputs, it would escape with probability $(1 - d/2^n)^t$.

For a 12-input block, a rewrite wrong on 2 inputs, and 4,096 vectors, the escape probability is $(1 - 2/4096)^{4096} \approx e^{-2} \approx 0.135$. One such bug in seven survives the run. For a 64-input block no simulation budget gets near certainty. So the fast mode can rank and filter candidates, but its "probably equal" is a hypothesis, not a verdict.

## Exact engines as the accurate mode

SAT on the miter, or a BDD comparison where BDDs stay small, decides equivalence on every input. SAT-sweeping engines combine the modes: simulation splits nodes into candidate equivalence classes, SAT is called only on pairs that survive, and every counterexample SAT finds is added to the simulation vectors, where it often splits many other classes at once.

Only an exact engine can produce the verdict "equal on every input". It is also the expensive step, so a well-built loop budgets it: limit conflicts per call, try the cheapest candidates first, and treat UNKNOWN as unresolved rather than as either answer.

## When the agent must recurse

The rule is simple: escalate to the exact engine before any transform is committed to the design or reported as equivalent. Ranking 200 candidate rewrites by estimated savings needs no proof, because a wrong rank costs only a little optimization; the 5 rewrites actually kept must each pass an exact check. A mismatch found by simulation needs no escalation at all, because it is already a definite answer.

Log the mode that produced each verdict. A downstream step that sees "equivalent (simulation, 4,096 vectors)" can demand a proof; one that sees an unlabeled "equivalent" cannot. Provenance is what keeps the fast loop honest when many agents and tools share one design database.

## Going deeper

Mishchenko, Chatterjee, Jiang and Brayton, "FRAIGs: a unifying representation for logic synthesis and verification", UC Berkeley technical report (2005). Kuehlmann, Paruthi, Krohm and Ganai, "Robust Boolean reasoning for equivalence checking and functional property verification", IEEE TCAD (2002).
