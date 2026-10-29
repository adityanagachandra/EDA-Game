# Learning-based routing

Routers make many heuristic choices: which nets to route first, which to rip up, how to price congested edges, which pin-access pattern to try. Learning can make those choices, while the router's path search and design-rule model keep routes legal. The search-theoretic cost of a learned heuristic is precise: an inadmissible estimate can make A* return longer paths, by a bounded amount.

## What a router can learn

A negotiated-congestion global router repeatedly rips up and reroutes nets with history costs on overused edges. A learned model can predict which nets to rip up and in what order, or predict a congestion cost map that the router adds to its edge costs before the first iteration. A detailed router can use a model to rank pin-access patterns by predicted violation risk. In every case the model biases the search and the router still finds paths, checks rules and accounts overflow. A wrong prediction then costs runtime or wirelength, not correctness, which is why these are the safe places for learning.

## Heuristics and optimality in A*

A* returns an optimal path when its heuristic never overestimates the remaining cost. A learned heuristic can overestimate, and two bounds describe the damage. Weighted A*, $f = g + w h$ with an admissible $h$, returns a path costing at most $w C^*$: with $w = 1.5$ and an optimal cost of 40 track units, at most 60. If a nonnegative heuristic overestimates the true remaining cost by at most $\epsilon$ everywhere, A* returns a path costing at most $C^* + \epsilon$. Inflated heuristics expand fewer nodes, so a learned cost model that speeds up the search and lengthens wires has made exactly this trade.

## Evaluating a learned router policy

A learned net ordering trained on one family of designs may fail on a new floorplan style, such as a design with many macros, where congestion patterns differ: that is distribution shift, and it must be measured on held-out design types. The baseline must be the same router with its standard heuristic, tuned, on the same designs and runtime budget; a disabled or untuned baseline inflates the gain. Before promotion, a differential test checks that routes stay legal and connected and that overflow, wirelength and runtime changes against the baseline are within agreed bounds. The policy belongs behind a cost or ordering seam so it can be rolled back.

## Going deeper

Pearl, "Heuristics: Intelligent Search Strategies for Computer Problem Solving" (1984), on weighted and epsilon-admissible search; McMurchie and Ebeling, "PathFinder: A Negotiation-Based Performance-Driven Router for FPGAs" (FPGA 1995).
