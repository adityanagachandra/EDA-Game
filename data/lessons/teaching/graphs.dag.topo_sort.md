# Topological sorting of a DAG

A combinational netlist, cut at its registers, is a directed acyclic graph. A topological order lists its vertices so that every arc goes forward, which is exactly the order in which arrival times, logic values or signal probabilities can be computed in one pass. Timers, simulators and synthesis tools all start by sorting the graph, and they must do it in linear time and update it as the netlist changes.

## Kahn's algorithm

Kahn's algorithm keeps the in-degree of every vertex. Vertices with in-degree 0 (the primary inputs and register outputs) are ready. Repeatedly remove a ready vertex, append it to the order, and decrement the in-degree of each of its successors; any successor that reaches 0 becomes ready. If vertices remain when no vertex is ready, the graph has a cycle.

With arcs $v_1 \to v_2$, $v_1 \to v_5$, $v_2 \to v_3$, $v_4 \to v_3$, $v_5 \to v_3$ and the smallest-label ready vertex removed first, the ready set starts as $\{v_1, v_4\}$. The order is $v_1, v_2, v_4, v_5, v_3$. Each arc is touched once, so the cost is $O(V + E)$.

## Depth-first search order

The second standard method runs depth-first search and records each vertex when it finishes (all its descendants are done). Reversing the finishing order gives a topological order, because in a DAG a vertex always finishes after everything reachable from it. A back arc found during the search, an arc to a vertex still on the recursion stack, proves a cycle.

Both methods are linear. Kahn's algorithm is natural for levelized processing and for parallel evaluation, since every vertex in the ready set can be processed at once. The DFS method is natural when the tool already traverses cones, for example when it sorts only the fan-in cone of one output. Real timers keep a level for every vertex, computed with Kahn's algorithm, and sort dirty vertices by level during incremental updates.

## How free the order is

A DAG usually has many topological orders. The constraints come only from ancestors and descendants: a vertex must follow all of its ancestors and precede all of its descendants. So the earliest position a vertex $v$ can take is $|\text{anc}(v)| + 1$, and the latest is $n - |\text{desc}(v)|$. Everything else is free.

That freedom is useful. Simulators choose orders that keep related vertices close in memory; parallel timers choose orders that expose independent vertices; synthesis tools choose orders that make cut enumeration cache-friendly. It also means a tool must not depend on one particular order for correctness, only on the order being topological, which is a property a differential test can check directly.

## Going deeper

Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms" (4th ed., 2022), section 20.4; Kahn, "Topological sorting of large networks" (CACM, 1962).
