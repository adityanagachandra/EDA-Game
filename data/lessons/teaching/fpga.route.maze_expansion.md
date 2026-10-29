# Wavefront expansion on the RR graph

FPGA routers connect each net by a wavefront search over the RR graph: Dijkstra's algorithm with node costs, sped up by a directed (A*) estimate toward the sink. Multi-sink nets are routed sink by sink, with the partial tree reused as a free source. This search is the inner loop of PathFinder, so its cost model and efficiency set the router's quality and speed.

## Shortest paths with node costs

FPGA routers price nodes, not edges: a path's cost is the sum of the costs of the wires it uses. Dijkstra's algorithm adapts directly. Seed the priority queue with the wires the source's output pin reaches, at their own cost; pop the cheapest wire, and for each neighbor reachable through a switch, offer its tentative cost as the popped cost plus the neighbor's cost. Stop when a wire that reaches the sink's input pin is popped. With source wires $w_5$ (cost 6) and $w_4$ (cost 3), and $w_4 \to w_8 \to w_3$ costing $3 + 4 + 3 = 10$, the route to a sink reachable from $w_3$ costs 10. Because all costs are positive, the first time the sink is popped its cost is optimal.

## Routing multi-sink nets

A net with several sinks is routed one sink at a time. After the first sink is reached, every node on the path joins the net's routing tree, and the next search starts from the whole tree with those nodes at zero cost, because the net already owns them. The branch to sink 2 then costs only the new wires it adds. If the route to sink 1 uses $w_3$ and $w_1$ at cost 3, and sink 2 is reached by a new branch costing 7, the tree costs 10. This incremental, greedy Steiner construction is not optimal, but it reuses wiring naturally. VPR orders sinks, for example by criticality, because the first sinks routed shape the tree for the rest.

## Directed search

A plain wavefront expands in all directions. A directed search orders the queue by $f = g + \alpha h$, where $g$ is the path cost so far and $h$ estimates the remaining cost, for example remaining distance times the cheapest cost per tile. With $\alpha = 1$ and an $h$ that never overestimates, A* still returns an optimal path while expanding far fewer nodes. VPR's astar_fac sets $\alpha$, often above 1, which overestimates and pulls the search hard toward the sink: much faster, slightly suboptimal. With nodes at $(g, \text{distance}) = (7, 2)$, $(38, 1)$ and $(12, 10)$, and $h = $ distance, the $f$ values are 9, 39 and 22, so the first node is expanded next. Directed search is what makes routing large FPGAs tractable.

## Going deeper

McMurchie and Ebeling, "PathFinder: A Negotiation-Based Performance-Driven Router for FPGAs", FPGA (1995); Swartz, Betz and Rose, "A Fast Routability-Driven Router for FPGAs", FPGA (1998).
