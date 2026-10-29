# Speedup accounting for surrogates

A surrogate is a fast model that stands in for an expensive evaluation, such as SPICE, routing or signoff STA, inside a search loop. Its value is the end-to-end speedup of the loop, not its own inference speed. That speedup must include the evaluations that still run for verification, the cost of building the surrogate, and the part of the flow it does not touch.

## Screening plus verification

Suppose a loop scores $N = 1000$ candidate sizings, SPICE takes 60 minutes each and the surrogate 1 minute, and the top 5 percent by surrogate score are re-run in SPICE before anything is accepted. Without the surrogate the loop costs $1000 \times 60 = 60{,}000$ minutes. With it the loop costs $1000 \times 1 + 0.05 \times 1000 \times 60 = 4000$ minutes, a speedup of 15. The verification term dominates: as the surrogate gets faster the speedup approaches $1/0.05 = 20$, never more. The fraction sent to verification is a policy decision, and shrinking it trades safety for speed.

## Amortizing the build cost

The surrogate itself costs labels and training. If building it took 200 machine-hours (12,000 minutes) and it serves 4 exploration campaigns, each campaign carries 3000 minutes of build cost, and the loop above costs $4000 + 3000 = 7000$ minutes, a speedup of $60{,}000/7000 = 8.6$. Served once, the same surrogate gives $60{,}000/16{,}000 = 3.75$. A reported speedup that leaves out data generation and training answers a different question: how fast the model runs once it exists. A tool designer needs the amortized number, which depends on how often the model is reused before the design, node or tool changes enough to invalidate it.

## The whole-flow bound

Accelerating one step speeds up the flow only in proportion to that step's share of runtime, as Amdahl's law states. If a step is a fraction $p$ of flow runtime and becomes $s$ times faster, the flow speedup is $1/((1 - p) + p/s)$. With $p = 0.6$ and $s = 10$ the flow is $1/(0.4 + 0.06) = 2.17$ times faster, and no surrogate can exceed $1/0.4 = 2.5$. Claims of orders-of-magnitude speedup for one analysis step are therefore compatible with modest flow-level gains. The step a surrogate replaces must also not be the one whose result the next stage depends on exactly.

## Going deeper

Forrester, Sobester and Keane, "Engineering Design via Surrogate Modelling" (2008), chapters 1 and 2.
