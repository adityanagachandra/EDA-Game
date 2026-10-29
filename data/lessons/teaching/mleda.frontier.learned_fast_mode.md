# Learned models as fast modes in agent loops

In an agent-first flow, engines come in fast and accurate versions, and an agent runs the fast ones continuously while recursing into the accurate ones where it matters. A learned predictor is one way to build a fast mode. Whether it helps depends on its recursion rate, the share of candidates it is unsure about, more than on its own speed, and an agent should watch that rate and stop using the fast mode when it stops paying.

## The loop's cost

An agent tuning placement evaluates $N$ candidates. Each goes through the learned fast mode at $t_s$; a fraction $f$ whose predictions are uncertain is re-run accurately at $t_a$; the final $m$ chosen candidates are always verified accurately. The loop costs $N t_s + (fN + m) t_a$. With $N = 500$, $t_a = 60$ minutes, $t_s = 3$ minutes, $f = 0.2$ and $m = 10$, that is $1500 + 110 \times 60 = 8100$ minutes, 135 hours. The uncertainty threshold that sets $f$ is a policy parameter, and the final verification is not optional: it is the gate on every learned output that is kept.

## Speedup and what dominates it

Running the accurate engine on all 500 candidates costs 30,000 minutes, so the fast mode gives a speedup of $30{,}000/8100 = 3.7$. Making the fast mode ten times faster would save only 1350 minutes; halving the recursion fraction would save 3000. The design choices that matter are therefore calibration (an uncertainty estimate that sends the right candidates to the accurate engine) and coverage of the training data (so fewer candidates are uncertain), not raw inference speed. These are the same observations the agent can read from its hooks: recursion rate, and the disagreement between fast and accurate results on recursed candidates.

## When to abandon the fast mode

The loop is slower than running the accurate engine on everything when $N t_s + (fN + m) t_a > N t_a$, that is when $f > (N t_a - N t_s - m t_a)/(N t_a)$. In the example the break-even fraction is $(30{,}000 - 1500 - 600)/30{,}000 = 93$ percent. Long before that, a rising recursion rate signals distribution shift: the design has left the fast mode's training distribution. A well-designed agent then routes this design to the accurate engine, logs the event, and lets the on-site learning loop collect the accurate results as new training data for the next version of the fast mode.

## Going deeper

Forrester, Sobester and Keane, "Engineering Design via Surrogate Modelling" (2008), on multi-fidelity use of surrogates; the OpenROAD documentation for the fast and full engines in its flow.
