# Prim-Dijkstra routing trees

A minimum spanning tree minimises total wirelength; a shortest-path tree gives every sink the shortest path from the source. Timing-driven routing wants something in between. The Prim-Dijkstra construction of Alpert, Hu, Huang and Kahng blends the two with one parameter, trading wirelength for source-to-sink path length.

## Prim's end: c equals zero

Grow a tree from the source. At each step, over every tree pin $u$ and outside pin $v$, add the edge minimising $c \cdot pl(u) + d(u, v)$, where $pl(u)$ is the tree path length from the source to $u$ and $d$ is the Manhattan distance. With $c = 0$ the key is the edge length alone: this is Prim's algorithm, and the result is a rectilinear minimum spanning tree.

Example: source (8, 4) and sinks (7, 0), (0, 7), (1, 2) give an MST of length 19. Its radius, the longest source-to-sink tree path, can be long, because the MST happily chains sinks one after another away from the source.

## Dijkstra's end: c equals one

With $c = 1$ the key is $pl(u) + d(u, v)$, the length of the source path the new pin would get. This is Dijkstra's algorithm, and the result is a shortest-path tree: every sink's tree path equals the best achievable through the pins, so on a Manhattan plane the radius is minimal.

The price is wirelength. Each sink tends to connect almost directly toward the source, so sinks close to each other do not share wire. Example: a net whose MST is 44 might have a shortest-path tree of 79. Neither extreme is right for a critical net with many sinks: the MST is too slow, and the shortest-path tree loads the driver with too much capacitance.

## The trade-off in between

For $0 < c < 1$ the construction interpolates. Raising $c$ penalises attaching a sink to a tree pin that is already far from the source, so radius falls and wirelength rises; lowering it does the opposite. Example: with $c = 0.75$, a net might keep the MST's wirelength 36 and radius 19 if its sinks are already well placed, while other nets change shape sharply.

In practice tools pick $c$ per net from its criticality, then refine the tree into a Steiner tree and buffer it. The construction is $O(n^2)$ for $n$ pins, trivial next to routing, so it can be rerun whenever slack changes. A later refinement (PD-II) post-processes the tree to recover wirelength without lengthening critical paths.

## Going deeper

C. J. Alpert, T. C. Hu, J. H. Huang and A. B. Kahng, "A Direct Combination of the Prim and Dijkstra Constructions for Improved Performance-Driven Global Routing" (1993); C. J. Alpert et al., "Prim-Dijkstra Revisited" (2018).
