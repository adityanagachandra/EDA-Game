# Planarity and embedding in layout

A graph is planar if it can be drawn in the plane without crossings. Planarity is the theory behind single-layer routing, river routing and the need for vias, and graph drawing methods are close relatives of placement. A tool designer needs the quick necessary tests, the exact characterization, and the algorithmic facts, without confusing a necessary condition with a sufficient one.

## Euler's bound as a quick test

For a connected planar graph with $v$ vertices, $e$ edges and $f$ faces, Euler's formula gives $v - e + f = 2$. Each face is bounded by at least three edges and each edge borders at most two faces, so a simple planar graph with $v \ge 3$ has $e \le 3v - 6$. If it has no triangles (for example, if it is bipartite), each face has at least four edges and $e \le 2v - 4$.

A conflict graph with 14 vertices and 37 edges exceeds $3 \cdot 14 - 6 = 36$, so it is not planar. But a graph within the bound may still be non-planar: $K_{3,3}$ has 6 vertices and 9 edges, within $3v - 6 = 12$, and is non-planar (the triangle-free bound $2v - 4 = 8$ does catch it). The bound is necessary, never sufficient.

## Kuratowski and linear-time testing

Kuratowski's theorem characterizes planarity exactly: a graph is planar if and only if it contains no subdivision of $K_5$ or $K_{3,3}$. Wagner's equivalent version uses minors. These are the certificates a tool can report when a single-layer routing is impossible.

Algorithmically, planarity is easy: Hopcroft and Tarjan gave a linear-time test, and later algorithms (Boyer-Myrvold, for example) also output either a planar embedding or a Kuratowski subgraph in linear time. In river routing, where nets connect two sides of a channel on one layer, the planarity condition reduces to order: nets must appear in the same order on both sides. A net out of order must cross another, which on one layer is impossible.

## Drawing graphs and placing cells

Most netlists are far from planar, which is why chips stack many metal layers and why routing is about layer assignment and vias rather than planarity. Planarity ideas survive in pieces: each layer's wires form a planar drawing, power grids and clock spines are planar structures, and orthogonal drawings (edges as horizontal and vertical segments) require vertex degree at most 4 when vertices are points.

Graph drawing and placement share mathematics. Tutte's theorem says that pinning one face of a 3-connected planar graph on a convex polygon and placing every other vertex at the average of its neighbours gives a crossing-free drawing with convex faces. That barycentric system is a Laplacian solve with fixed vertices, the same system a quadratic placer solves with fixed I/O pins, and force-directed drawing is a direct ancestor of force-directed placement.

## Going deeper

Hopcroft and Tarjan, "Efficient planarity testing" (Journal of the ACM, 1974); Di Battista, Eades, Tamassia and Tollis, "Graph Drawing: Algorithms for the Visualization of Graphs" (1999).
