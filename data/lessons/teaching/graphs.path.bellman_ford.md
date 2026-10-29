# Bellman-Ford and negative edge weights

Some EDA graphs carry negative weights: constraint graphs for layout compaction and retiming, timing graphs with hold constraints written as negative arcs, and residual graphs in min-cost flow. Dijkstra's algorithm fails on them, and a negative cycle means the problem has no solution at all. Bellman-Ford handles negative weights, detects negative cycles, and is the workhorse for difference constraints.

## Relaxing every arc in rounds

Bellman-Ford sets $d(s) = 0$ and $d(v) = \infty$ elsewhere, then performs rounds. In each round it relaxes every arc once: $d(v) = \min(d(v), d(u) + w(u,v))$. After round $k$, $d(v)$ is at most the length of the shortest path from $s$ to $v$ that uses at most $k$ arcs. A shortest simple path has at most $n - 1$ arcs, so $n - 1$ rounds suffice when there is no negative cycle.

Arcs $v_1 \to v_3$ (3), $v_3 \to v_5$ (2), $v_2 \to v_6$ (-4), $v_1 \to v_2$ (8) give $d(v_6) = 8 - 4 = 4$, found once $v_2$ has its distance. The algorithm costs $O(VE)$, slower than Dijkstra, but it never assumes that a vertex's distance is final early.

## Arc order and early termination

The round bound is a worst case. Within a round, a relaxation can use a value lowered earlier in the same round, so the arc order matters: if arcs happen to be listed in topological order, one round finishes the job; if they are listed in reverse, each round advances distances by only one arc. After round 1 an estimate may still be too high, and later rounds lower it.

So implementations stop as soon as a full round changes nothing, and they keep a queue of vertices whose distance changed so that only their outgoing arcs are relaxed next (the queue-based variant, often called SPFA). On sparse, nearly acyclic constraint graphs this makes Bellman-Ford fast in practice while keeping its guarantees.

## Negative cycles

If a cycle has negative total weight, going around it again always shortens the path, so shortest paths are undefined. Bellman-Ford detects this: if any arc can still be relaxed after $n - 1$ rounds, a negative cycle is reachable from the source. Following predecessor pointers from that arc leads into the cycle.

In a constraint graph a negative cycle means the constraints are contradictory. A cycle $v_1 \to v_2 \to v_3 \to v_4 \to v_1$ with weights -6, -6, 0 and $x$ has weight $x - 12$, so the system is feasible only for $x \ge 12$. A compactor that finds such a cycle reports the constraints on it, which tells the designer exactly which spacing rules or user constraints conflict.

## Going deeper

Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms" (4th ed., 2022), sections 22.1 and 22.4; Lengauer, "Combinatorial Algorithms for Integrated Circuit Layout" (1990), the chapter on compaction.
