# A* routing

A* is the maze search used inside modern detailed and global routers. It expands cells in order of $f = g + h$, the cost so far plus an estimate of the cost to go. With the Manhattan distance as $h$ on a unit-cost grid, the estimate never overshoots, so A* still finds a shortest route while expanding far fewer cells than Lee's wavefront.

## Priority f equals g plus h

For a cell $v$, $g(v)$ is the cost of the best route found from the source to $v$, and $h(v)$ estimates the remaining cost to the target. A* keeps an open list ordered by $f(v) = g(v) + h(v)$ and always expands the smallest. On a grid with unit step cost the Manhattan distance $h(v) = |x_v - x_t| + |y_v - y_t|$ is admissible (it never exceeds the true remaining cost) and consistent ($h(u) \le 1 + h(v)$ for neighbours), so A* returns a shortest route.

Example: a cell reached by its shortest route of length $g = 8$ that is 5 columns and rows from the target has $f = 13$. If the best route to the target costs 12, that cell will never be expanded by a correct A*, while Lee would label it.

## Cells A* must expand

Let $C^*$ be the cost of the best route. With a consistent heuristic, $f$ never decreases along a shortest route, and A* expands every cell with $f(v) < C^*$ before it finishes; cells with $f(v) = C^*$ may or may not be expanded, depending on tie-breaking, and cells with $f(v) > C^*$ never are. Lee, by contrast, must label every reachable cell with $g(v) < C^*$.

So the guaranteed saving is the count of cells with $g < C^*$ minus the count with $f < C^*$. Example: with $C^* = 9$, Lee must label 24 cells while A* must expand only 10. The cells pruned are those whose Manhattan distance to the target makes them hopeless. On an open grid every cell in the pins' bounding box has $f = C^*$, so A* with good tie-breaking (prefer larger $g$) walks almost straight to the target.

## Heuristics with real costs

Real routers price vias, bends, wrong-way steps and congestion, so step costs are not all 1. The heuristic must stay admissible under those costs: Manhattan distance times the cheapest per-step cost, plus the minimum number of vias the target still needs times the via cost, is a common choice. An inadmissible heuristic (for example, weighting $h$ by more than 1) searches faster but may return a longer route; some routers accept that trade for non-critical nets.

A* also needs a clean priority queue, and on large grids the open and closed sets dominate memory. Detailed routers bound the search to the net's route guides or a box around the pins, which caps both. The search must also handle many sources and targets at once (a partial tree, a pin shape with several access points), which A* does by seeding all sources with $g = 0$ and using the distance to the nearest target as $h$.

## Going deeper

Hart, Nilsson and Raphael, "A Formal Basis for the Heuristic Determination of Minimum Cost Paths" (1968); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); the TritonRoute papers for A* in detailed routing.
