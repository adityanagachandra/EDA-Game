# Layout decomposition as graph coloring

To print a layer with several masks, every shape must be assigned to one mask so that shapes too close for one mask land on different masks. This is graph coloring: shapes are vertices, too-close pairs are edges, masks are colours. Decomposition engines in physical verification and colour-aware routers solve it, and its failures, odd cycles and dense clusters, show up as native conflicts that the layout must fix.

## Building the conflict graph

Every pair of shapes whose separation is below the same-mask minimum spacing $d$ becomes an edge of the conflict graph. Separation means the distance a mask can resolve: edge to edge for wires, centre to centre for small vias. Edges can join non-neighbours. Vias at positions 0, 30, 70 and 90 nm with $d = 50$ nm give edges (0, 30), (30, 70) and (70, 90), at distances 30, 40 and 20 nm; the pairs (0, 70), (0, 90) and (30, 90) are 70, 90 and 60 nm apart and do not conflict. The graph has three edges. Building the graph efficiently uses spatial indices or a sweep so that only nearby pairs are tested.

## Masks as colours

A decomposition with $k$ masks is a proper $k$-colouring of the conflict graph, so the fewest masks needed is the chromatic number. Two masks suffice exactly when the graph is bipartite, meaning it has no odd cycle. Three shapes that all conflict with each other form a triangle and need three masks; four mutually conflicting shapes need four. Two-colourability is checkable in linear time by breadth-first search, but deciding three-colourability is NP-complete, so triple-patterning decomposers use heuristics, integer programming or SAT on the small components that remain after splitting the graph into connected pieces.

## Native conflicts and stitches

When a graph is not two-colourable but the layer is LELE, some edges must remain inside one mask, and each is a violation to fix. The fewest such edges is found by choosing the best two-colouring, which is a maximum-cut problem. A wire can sometimes be split into two pieces on different masks with a small overlap, a stitch, which breaks an odd cycle at the cost of a yield-sensitive overlap region. Otherwise the layout must change: spread the shapes or reroute. Colour-aware routers and placers avoid creating odd cycles in the first place, which is cheaper than fixing them after the fact.

## Going deeper

Kahng, Park, Xu and Yao, "Layout decomposition for double patterning lithography", ICCAD (2008); Yu, Yuan, Gao and Pan, "Layout decomposition for triple patterning lithography", ICCAD (2011).
