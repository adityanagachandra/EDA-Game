# PathFinder negotiated-congestion routing

PathFinder routes all nets at once by negotiation. Every net is routed with the cheapest path under current prices, nodes may be shared temporarily, and the price of a shared node rises until nets that have alternatives move away. It is the basis of nearly every FPGA router, VPR's and nextpnr's among them, because FPGA routing resources are discrete and fixed.

## The node cost

PathFinder prices node $n$ as $c_n = (b_n + h_n) \cdot p_n$. The base cost $b_n$ reflects the resource (often its delay or its type); the history cost $h_n$ accumulates past congestion; and the present cost $p_n = 1 + p_{fac} \cdot \max(0, \mathrm{occ}_n + 1 - \mathrm{cap}_n)$ charges for the overuse that adding one more net would create. A wire with $b = 3$, $h = 0.5$, occupancy 0 and capacity 2 costs 3.5: no overuse. A wire with $b = 2$, $h = 2.25$, occupancy 2 and capacity 1, at $p_{fac} = 1$, has $p = 1 + 1 \times 2 = 3$ and costs $4.25 \times 3 = 12.75$. Prices, not prohibitions, keep the search over the whole graph open.

## Negotiation between paths

Each iteration rips up and reroutes every net (or every net touching congestion) with its cheapest path under the current prices. A net choosing between two paths sums the node costs along each and takes the cheaper, even if it is congested. Path A through nodes costing 6.25, 6.25 and 6.25 totals 18.75; path B through nodes costing 2, 2.25 and 11 totals 15.25, so the net takes B although B's last node is overused. Over iterations, nets with cheap alternatives leave contested nodes, and the net that most needs a node keeps it. This is a market: the present cost is a price that rises with demand, and the history cost remembers which resources were persistently scarce.

## Convergence and legality

A routing is legal when no node is overused: $\sum_n \max(0, \mathrm{occ}_n - \mathrm{cap}_n) = 0$. With nodes at occupancy over capacity 4/2, 0/2, 3/1, 4/1 and 4/1, the total overuse is $2 + 0 + 2 + 3 + 3 = 10$, so another iteration is needed. PathFinder iterates until overuse reaches zero or an iteration limit is hit, in which case the design is declared unroutable at this channel width. The present-cost factor $p_{fac}$ grows each iteration, so sharing is cheap at first, which lets nets explore, and prohibitive later, which forces a legal solution. Convergence is not guaranteed in theory, but in practice it is reliable, and overuse per iteration is the metric routers report.

## Going deeper

McMurchie and Ebeling, "PathFinder: A Negotiation-Based Performance-Driven Router for FPGAs", FPGA (1995); Betz, Rose and Marquardt (1999), chapter 6.
