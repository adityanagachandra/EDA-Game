# Safety and liveness properties

Every linear-time property is either about something bad never happening, about something good eventually happening, or a combination of the two. This split into safety and liveness decides which engines can check a property, what its counterexamples look like, and whether environment fairness must be modelled. A verification engineer should classify every assertion before choosing how to prove it.

## Safety: nothing bad happens

A safety property says that something bad never happens: the FIFO never overflows, two grants are never active together ($G \neg(g_1 \wedge g_2)$), the state machine never enters an illegal state. Its defining feature is that every violation has a finite witness: a prefix of the run ending at the bad moment, after which no continuation can repair it. So counterexamples to safety are finite traces, and bounded model checking finds every reachable violation once its depth reaches the trace length. Bounded response properties are safety too: req |-> ##[1:4] ack fails at a definite cycle, four cycles after an unanswered request. Most hardware assertions are safety properties, and the fast engines (BMC, k-induction, IC3) are built for them.

## Liveness: something good eventually happens

A liveness property says that something good eventually happens: every request is eventually acknowledged ($G(\text{req} \rightarrow F\, \text{ack})$), the arbiter grants each client infinitely often ($G\, F\, \text{grant}_B$). No finite prefix can violate it, because the good event could still come later. A violation is an infinite run on which the good event never happens, and in a finite-state system such a run can always be written as a lasso: a prefix followed by a loop repeated forever without the good event. Model checkers report liveness counterexamples in exactly that form. Every linear-time property is the intersection of a safety property and a liveness property, so a flow can route the two parts to different engines.

## Fairness and checking liveness

Liveness checks often fail for reasons the real system rules out: a counterexample in which the environment never toggles an input, or a scheduler that never runs one process. Fairness constraints restrict the check to runs the environment can really produce, for example runs in which the input is 1 infinitely often. Engines check liveness in several ways. Automata-based model checking searches for a reachable fair cycle. BMC with loop detection finds short lassos but, like BMC for safety, proves nothing beyond its bound. The liveness-to-safety translation instruments the design to guess a loop start, record the state, and flag a bad state when the run returns to it without the good event, so a safety engine such as IC3 can decide the liveness property. Teams also replace liveness by bounded response (\"ack within 16 cycles\"), which is safety but fails correct designs that are slower than the bound.

## Going deeper

Alpern and Schneider, "Defining Liveness" (Information Processing Letters, 1985); Biere, Artho and Schuppan, "Liveness Checking as Safety Checking" (Electronic Notes in Theoretical Computer Science, 2002).
