# Two-coloring for double patterning

At advanced nodes, some layers are printed with two exposures, each on its own mask. Features closer than the same-mask spacing must go on different masks. Deciding the masks is 2-coloring a conflict graph, which is easy when it is possible, and the cases where it is not possible (odd cycles) are exactly what layout and decomposition tools must find and fix.

## The conflict graph and its components

Make a vertex for each layout feature and an edge between two features closer than the minimum same-mask spacing. A valid decomposition assigns mask A or mask B to every feature so that every edge joins different masks: a proper 2-coloring.

Components are independent: features in different connected components never constrain each other. With features $F_1$ to $F_8$ and edges $F_1$-$F_6$, $F_2$-$F_6$, $F_5$-$F_6$, $F_3$-$F_8$, the components are $\{F_1, F_2, F_5, F_6\}$, $\{F_3, F_8\}$, $\{F_4\}$ and $\{F_7\}$: four components, each colored separately. Decomposition tools split large layouts this way to keep each subproblem small.

## Bipartite means decomposable

A graph is 2-colorable exactly when it is bipartite, which is exactly when it has no odd cycle. Breadth-first search decides this in linear time: color the start A, its neighbours B, theirs A, and report an odd cycle if an edge ever joins two vertices of the same color.

If every component is bipartite, each component has exactly two colorings (swap A and B), so a layout with $k$ components has $2^k$ valid mask assignments; a single odd cycle anywhere makes the count 0. With 2 components and no odd cycle, there are 4 assignments. The freedom matters: tools use it to balance mask density or align the masks of repeated cells, choosing one coloring per component.

## Odd cycles and their repair

An odd cycle, such as three features each too close to the other two, cannot be decomposed. Every assignment leaves at least one violated edge on every odd cycle; finding the fewest violations over a whole component is the maximum-cut problem in disguise, NP-hard in general, though small components are solved exactly by search and planar conflict graphs by special algorithms.

Each remaining violation needs a fix: a **stitch** splits one feature into two overlapping parts on different masks, breaking the cycle; a layout change moves a feature to increase spacing; or a third mask is added (triple patterning, whose 3-coloring is NP-complete in general). Design rules increasingly forbid odd-cycle configurations up front, so that legal layouts are decomposable by construction.

## Going deeper

Kahng, Park, Xu and Yao, "Layout decomposition for double patterning lithography" (Proc. ICCAD, 2008).
