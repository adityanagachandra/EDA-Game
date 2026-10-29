# Double-patterning coloring conflicts

When features are packed tighter than one exposure can print, double patterning splits a layer onto two masks. Two features closer than the same-mask spacing must go on different masks, which makes mask assignment a graph 2-coloring problem. A layout is decomposable exactly when its conflict graph is bipartite, so double-patterning checks reduce to finding odd cycles.

## Building the conflict graph

Each feature is a vertex, and an edge joins two features closer than the same-mask spacing $S$, measured as Euclidean distance. With $S = 27$ nm, five features may produce two edges. Building the graph is a spacing check with a larger distance: candidate pairs from a sweep, then exact distances. Edges mean must-differ constraints. Features farther than $S$ from everything form isolated vertices, free to go on either mask, which the decomposer uses to balance density between the masks.

## Odd cycles

A graph is 2-colourable if and only if it has no odd cycle. A breadth-first search that alternates colours finds a same-colour edge exactly in a component that contains an odd cycle. Three features each closer than $S$ to the other two form a triangle, which no assignment can satisfy. Each component is checked separately, and each non-bipartite component is a double-patterning violation. Fixes are a layout change that removes an edge, or a stitch that splits one feature into two overlapping parts on different masks, where the process allows stitches.

## Counting mask assignments

A connected bipartite component has exactly two colourings: fix one vertex's mask and every other vertex is forced. Components are independent, so a graph with $c$ components, all bipartite, has $2^c$ assignments, counting isolated features as components; one odd component makes the count zero. Seven features in two components give 4 assignments. Decomposers use this freedom to balance pattern density between masks and to avoid stitches. Triple patterning raises the problem to 3-coloring, which is NP-complete in general, and changes the algorithms used.

## Going deeper

Kahng, Park, Xu and Yao, "Layout decomposition for double patterning lithography" (ICCAD, 2008).
