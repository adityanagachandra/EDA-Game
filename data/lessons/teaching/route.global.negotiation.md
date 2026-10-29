# Negotiated-congestion global routing

Negotiated congestion, introduced as PathFinder by McMurchie and Ebeling for FPGA routing, is the dominant global-routing strategy. Every net is routed by a shortest-path search on costs that rise with present overuse and with a history of past overuse. Nets that have alternatives move away from contested resources, and nets without alternatives keep them.

## The PathFinder cost

For a resource $n$ (a gcell edge, a track, an FPGA wire), PathFinder charges $c_n = (b_n + h_n)\,p_n$. Here $b_n$ is the base cost (length or delay), $h_n$ the history cost and $p_n$ the present-congestion factor, $p_n = 1 + p_{fac}\max(0, occ_n + 1 - cap_n)$, where $occ_n$ counts the nets already using $n$ and the $+1$ accounts for the net being routed.

Example: $b_n = 3$, $h_n = 0.5$, capacity 1, occupancy 3 and $p_{fac} = 2$ give $p_n = 1 + 2 \cdot 3 = 7$ and $c_n = 3.5 \cdot 7 = 24.5$. An edge with spare capacity has $p_n = 1$ and costs only $b_n + h_n$. The multiplicative form means history matters most on edges that are also presently congested.

## Routing on negotiated costs

Each iteration rips up and reroutes every net (or every net on an overflowed resource) with a shortest-path search, A* or Dijkstra, on the current costs. Occupancies update as nets are routed, so later nets in the same iteration see the congestion earlier nets created. A short path through full edges can then cost more than a longer path through free ones.

Example: path A has four edges at or above capacity and path B has five free edges with some history. With $p_{fac} = 0.5$, A might cost 5 and B 8.5, so the net still takes A. As $p_{fac}$ grows in later iterations, A's cost rises and the net moves, or another net sharing A's edges moves instead. The net that has the cheaper alternative is the one that yields: that is the negotiation.

## Raising the price over iterations

PathFinder starts with a small $p_{fac}$, so early iterations allow sharing and let each net find its preferred region, then multiplies $p_{fac}$ by a factor (often between 1.3 and 2) every iteration. Sharing becomes ever more expensive until no resource is overused. Example: $p_{fac} = 1$ in iteration 1 doubling each iteration gives 4 in iteration 3, so an edge at capacity 3 with occupancy 3, $b = 1$ and $h = 1.5$ costs $2.5 \cdot 5 = 12.5$.

The schedule is a speed-quality knob: grow too fast and nets freeze into the first legal-looking configuration with long detours; grow too slowly and runtime balloons. The method resembles Lagrangian relaxation of the capacity constraints, with the present and history terms acting as adaptive multipliers. It is also the template for global routers such as NCTU-GR and the negotiation loop in OpenROAD's grt.

## Going deeper

L. McMurchie and C. Ebeling, "PathFinder: A Negotiation-Based Performance-Driven Router for FPGAs" (1995); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); the VPR routing documentation.
