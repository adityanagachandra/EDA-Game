# Seams and tiers in a placer and router

An agent-first FPGA placer or router is built in four layers: an immutable core that guarantees legality, strategy plug-ins for heuristics, parameters for tuning, and observation hooks for live state. Most customer requests then land as configuration or plug-ins, never as core edits. Deciding where each piece of VPR or nextpnr belongs is the design skill this concept practices.

## The four layers of a placer and a router

In a router, the core is the routing-resource graph's legality (a net uses only existing edges, and in the final solution no node exceeds its capacity) plus the wavefront kernel. Strategy is the cost function, the rip-up order and the net ordering. Parameters are base costs per segment type, $p_{fac}$ and its multiplier, $h_{fac}$, and the A* factor. Hooks are per-iteration callbacks exposing overuse and critical-path delay. In an annealing placer, the core is legalization and macro handling, such as carry chains moved as one; strategy is the move set and the cost terms; parameters are InnerNum, $\lambda$ and the criticality exponent; and hooks expose the cost and acceptance rate per temperature. The core is small and rarely touched.

## Placing a request in a change tier

A request goes to the smallest safe tier. Making one region's long wires 20% more expensive is Tier 1 when base costs per resource are already a parameter table. Penalizing routes that use more than two long wires in a row needs new logic, but if a CostProvider seam is called per node expansion with the partial path's context, it is a Tier 2 plug-in. A soft goal such as keeping two blocks within three columns belongs in a Tier 2 cost term, not in a Tier 3 edit to the legality check, which would turn a preference into a constraint and could make placements infeasible. For a router seam to serve both timing-driven and congestion-aware plug-ins, it must pass the node, its type and location, its occupancy, capacity and history, and the connection's criticality.

## Seams in hot loops

Router cost functions run millions of times per iteration, so a seam there must be cheap: plain data, no hidden state, no allocation, and deterministic for a given seed, so the hot kernel can still be compiled and inlined around it. After a Tier 2 plug-in lands, the differential test must show that only the targeted nets changed, that legality holds, and that QoR elsewhere is within seed noise. When every new placement objective has required editing the annealer's inner loop, the objective seam is missing; the fix is to move cost terms behind an interface so that such requests become plug-ins. Determinism is not optional, because golden and differential tests compare runs, and nondeterminism makes a regression indistinguishable from noise.

## Going deeper

The VPR developer documentation on router and placer options.
