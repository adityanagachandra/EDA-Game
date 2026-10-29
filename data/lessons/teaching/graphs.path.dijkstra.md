# Dijkstra's shortest-path algorithm

Real routing costs are not uniform: vias cost more than wire, congested edges cost more than free ones, and wrong-way segments are penalized. Dijkstra's algorithm finds shortest paths with any non-negative edge costs, and it is the engine inside maze routers, global routers and many timing-driven tools. Knowing exactly what it settles, when, and at what cost per operation is basic equipment for a router designer.

## Settling vertices in order of distance

Dijkstra's algorithm keeps a tentative distance $d(v)$ for every vertex, starting with $d(s) = 0$ and $\infty$ elsewhere. It repeatedly **settles** the unsettled vertex with the smallest tentative distance, then relaxes its edges: for each neighbour $w$, set $d(w) = \min(d(w), d(v) + c(v,w))$.

Because costs are non-negative, a settled vertex's distance is final: any other route to it would pass through an unsettled vertex whose distance is already at least as large. With edges $v_1$-$v_4$ (2), $v_1$-$v_3$ (5), $v_4$-$v_6$ (3), $v_4$-$v_5$ (4), $v_2$-$v_6$ (4), the settling order from $v_1$ is $v_1$ (0), $v_4$ (2), then $v_3$ (5) and $v_6$ (5), then $v_5$ (6), then $v_2$ (9).

## Tentative distances during the run

Between settles, tentative distances are upper bounds that only decrease. After the source is settled, each neighbour's tentative distance is its direct edge cost; later settles may lower it. A vertex $v_6$ that first gets $d = 8$ through a direct edge may drop to 8 or lower when a cheaper two-edge route through $v_5$ is found.

The order of settling is fixed by the distances, with ties broken by a rule such as smaller label. A router's behaviour depends on that rule: when many equal-cost routes exist, the tie-break decides which one is found, which is why routers add small biases (prefer the current direction, prefer lower layers) rather than leaving ties to chance.

## Data structures and limits

With a binary heap, Dijkstra runs in $O((V + E)\log V)$; with a Fibonacci heap, $O(E + V \log V)$. On grid graphs with small integer costs, a bucket queue indexed by distance gives near-linear time, which many routers use. Each settle pops the heap, and each successful relaxation is a decrease-key or a lazy reinsert.

Two limits matter. Negative edge costs break the settled-is-final argument, so problems with negative costs need Bellman-Ford. And Dijkstra explores every direction equally: on a grid it settles a diamond of vertices around the source until the target is reached, which wastes work. A* fixes that by adding a heuristic to the priority, without giving up optimality when the heuristic is admissible.

## Going deeper

Dijkstra, "A note on two problems in connexion with graphs" (Numerische Mathematik, 1959); Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms" (4th ed., 2022), section 22.3.
