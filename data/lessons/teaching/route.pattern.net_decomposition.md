# Decomposing multipin nets into two-pin nets

Routers search for paths between two points, but most nets have more pins. Global routers therefore decompose each multipin net into two-pin connections, usually along a rectilinear Steiner or spanning tree, and route those. The choice of decomposition fixes how much demand the net places on the grid before any search starts.

## Spanning-tree decomposition

Build a rectilinear minimum spanning tree over the pins with Manhattan distances: Prim's algorithm adds, at each step, the outside pin closest to any pin already in the tree. Its $n - 1$ edges become the two-pin connections. Example: pins at (5, 15), (4, 4), (17, 9) and (7, 5) give an MST of length 30, and each of its 3 edges is routed separately (by patterns, then maze if needed).

MST decomposition is simple and deterministic, but it never uses Steiner points, places where a tree branches away from a pin. On a rectilinear grid that costs wire: the MST can be up to 1.5 times as long as the optimal Steiner tree (Hwang's bound). Global routers usually start from a Steiner tree instead.

## Steiner points save wire

For three pins the optimal rectilinear Steiner tree is easy: put one Steiner point at the median x and median y of the pins and connect each pin to it. Its length equals the half-perimeter of the pins' bounding box, which is optimal for three pins. Then decompose into three two-pin connections, each from a pin to the Steiner point.

Example: pins at (0, 0), (10, 2) and (4, 9). The MST is $12 + 13 = 25$; the Steiner point is (4, 2) and the tree has length $10 + 9 = 19$, the half-perimeter, saving 6. For more pins the problem is NP-hard, and routers use fast heuristics or lookup methods such as FLUTE that are optimal or near optimal for small degrees.

## Why not a star from the driver

A naive decomposition connects the driver to every sink directly. Its total length is $\sum_i d(s, t_i)$, often far more than an MST, because sinks close to each other each get their own long wire back to the driver. Example: a driver at one corner and four sinks clustered at the opposite corner pays the long distance four times.

The star's one virtue is that every sink gets a shortest path from the driver, which minimises source-to-sink delay when wire resistance dominates. That trade between total wirelength and source-to-sink path length is exactly what timing-driven tree construction (Prim-Dijkstra trade-offs) controls. For ordinary nets, total length is what loads the routing grid, so Steiner or spanning decompositions win.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); F. K. Hwang, "On Steiner Minimal Trees with Rectilinear Distance" (1976).
