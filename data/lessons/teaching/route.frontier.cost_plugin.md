# Cost functions as router plug-ins

The router's cost function is its most customer-specific behaviour: every team wants to penalise something different, from vias on critical nets to wiring over an analog block. Making the cost a plug-in at a defined seam lets an agent add or modify it as a Tier 2 change, while the core keeps search and legality fixed. The design question is what the seam passes and promises.

## Evaluating a plug-in cost

A cost plug-in receives an edge's context and returns a price. Example: $cost(e) = b_e(1 + \alpha\max(0, d_e + 1 - c_e)^2) + \beta h_e$, plus $\gamma$ for every via on a timing-critical net. With $\alpha = 1$, $\beta = 1$, $\gamma = 4$, a via edge with demand 5 and capacity 4 and history 3 costs $1 \cdot (1 + 4) + 3 + 4 = 12$, while a free wire edge with history 0 costs 1.

The core sums plug-in costs along each candidate and keeps the cheapest path, exactly as with a built-in cost. Squaring overflow is a typical Tier 2 idea: heavily overflowed edges become much more expensive, so demand spreads more aggressively at some cost in wirelength. The plug-in changes prices only; it cannot make an illegal path legal.

## The seam contract

For the core to trust a plug-in, the seam states a contract: finite, non-negative costs; deterministic results for the same context; no hidden global state; and only the context the seam passes. The core checks what it can, rejecting or clamping negative costs, because Dijkstra and A* lose optimality with them.

Performance is part of the contract. A plug-in that precomputes a cost map per iteration lets the maze kernel read a table instead of calling out per edge, which keeps the hot loop tight. The context the seam passes defines what is possible: to discourage vias on critical nets, the seam must expose whether an edge is a via and the net's slack or criticality.

## Proving a plug-in correct

A plug-in is accepted on evidence, not on reading its code. A differential run on reference designs, with and without it, must show that only the targeted nets' routes moved, that DRC stayed clean, and that global metrics (overflow, wirelength, timing) moved as intended. A plug-in that reduces overflow but raises detailed-routing violations is optimising a proxy that ignores something the next stage cares about.

When the same request recurs with different numbers, for example a different overflow exponent per design, the knob belongs in configuration: expose it as a plug-in parameter so the request becomes Tier 1. That is how the change surface grows deliberately, and how a validated plug-in eventually becomes a library skill.
