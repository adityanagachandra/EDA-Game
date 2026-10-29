# A* search

A* is Dijkstra's algorithm with a sense of direction. It orders the search by the cost so far plus an estimate of the cost to go, so a router expands vertices toward the target instead of in every direction. With an admissible estimate such as Manhattan distance on a unit-cost grid it still finds a shortest path. Every modern detailed router uses A* or a variant of it in its maze search.

## Priority f = g + h

A* keeps the cost from the source, $g(v)$, as Dijkstra does, and adds a heuristic $h(v)$ estimating the remaining cost to the target. It expands the open vertex with the smallest $f(v) = g(v) + h(v)$. On a grid with unit step cost, the Manhattan distance $h(v) = |x_v - x_t| + |y_v - y_t|$ is the natural heuristic: no route can be shorter.

A vertex at $(3, 3)$ reached with $g = 10$, routing to a target at $(7, 8)$, has $h = 4 + 5 = 9$ and $f = 19$. With $h = 0$ everywhere, A* is exactly Dijkstra. With a perfect heuristic, it expands only vertices on a shortest path.

## Which vertex is expanded next

A* expands the smallest $f$, not the smallest $g$. A vertex near the source with small $g$ but far from the target may wait, while a vertex with larger $g$ that is close to the target goes first. In an open list holding $(3, 2)$ with $g = 3$ and $(6, 12)$ with $g = 18$, routing to $(8, 14)$, the first has $f = 3 + 17 = 20$ and the second $f = 18 + 4 = 22$, so $(3, 2)$ is expanded first.

Ties on $f$ are common on grids, and breaking them toward larger $g$ (deeper vertices, closer to the target) sharply reduces expansions. The number of expansions, not the path length, is A*'s runtime, so tie-breaking and heuristic quality are the main levers a router designer has.

## Admissibility and weighted A*

A heuristic is **admissible** if it never overestimates the true remaining cost; then A* returns an optimal path. It is **consistent** if $h(u) \le c(u,v) + h(v)$ for every edge; then no vertex needs to be expanded twice. Manhattan distance is both on a unit-cost grid. If via costs or congestion raise edge costs above 1, it stays admissible; if some edges cost less than 1, it can overestimate and lose optimality.

Weighted A* uses $f = g + w h$ with $w > 1$. It expands far fewer vertices and returns a path of cost at most $w$ times the optimum when $h$ is admissible. A weighted search with $w = 2$ that returns cost 29 proves the optimum is at least 14.5. Routers trade this bounded suboptimality for speed in early iterations and tighten $w$ later.

## Going deeper

Hart, Nilsson and Raphael, "A formal basis for the heuristic determination of minimum cost paths" (IEEE Trans. SSC, 1968); the TritonRoute papers on its A*-based detailed router.
