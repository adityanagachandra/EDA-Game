# Routing as constrained graph search

Routing turns a placed netlist into metal: every net's pins must be joined by wires and vias, with no shorts, no design-rule violations and no more wires through any region than it can hold. For a tool designer, routing is constrained graph search. Knowing exactly which graph, which constraint and which objective is what separates a router you can reason about from a pile of heuristics.

## The routing graph

A router models the layout as a graph $G = (V, E)$. In detailed routing, a vertex is a grid point on a routing track of one metal layer. Wire edges join neighbouring points along a track, and via edges join the same point on adjacent layers. Blockages (macros, power straps, pre-routed nets) delete vertices or edges. Each edge carries a cost: its length, plus extra for vias, wrong-way wiring or congestion.

Routing one two-pin net is then a shortest-path problem: on a grid with unit costs, breadth-first search (Lee's algorithm) solves it, and with general non-negative costs Dijkstra or A* does. A net with $k > 2$ pins needs a tree that may branch at points that are not pins. Finding the cheapest such tree is the Steiner tree problem in graphs, which is NP-hard, so routers decompose multipin nets into two-pin connections or build trees heuristically.

## Capacity couples the nets

If each net could use any edge, routing would split into independent searches. It cannot: a track segment holds one wire, and a gcell edge holds as many wires as it has tracks. This capacity constraint couples all nets. Formally, choose for each net $i$ a tree $T_i$ so that for every edge $e$, the number of trees using $e$ is at most $cap(e)$, and minimise $\sum_i cost(T_i)$.

Routing nets one at a time on shortest paths is a greedy heuristic for that problem. Early nets take the best tracks, later nets detour or fail, and a different order can succeed where this one failed. A failure under one order therefore says the instance is congested, not that it is infeasible. Rip-up and reroute and negotiated congestion exist to undo the damage of a bad order.

## Legality first, then cost

A routed layout is useful only if it can be manufactured, so legality is a constraint, not a cost: every net connected, no shorts, every spacing, width, area and via rule satisfied. Among legal routings, the objective is usually a weighted sum of wirelength and via count, with timing, crosstalk and yield terms for nets that need them. During search, rule violations are priced with penalties so the router can pass through a bad state on the way to a good one, but the final check is pass or fail.

Scale forces a decomposition. A modern block has millions of nets and billions of grid points, so production flows route coarsely first (global routing on gcells, a few thousand tracks folded into each edge) and then in detail inside the corridor each net was given. Pin locations come from placement, so routing quality is bounded by placement quality: a hot spot of pins creates congestion no router can remove.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapters 5 to 7; Sherwani, "Algorithms for VLSI Physical Design Automation".
