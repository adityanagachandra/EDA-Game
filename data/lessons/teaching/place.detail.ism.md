# Independent-set matching

Independent-set matching picks a set of same-size cells that share no nets, so their wirelength costs are independent of each other, and reassigns them among their own locations by solving a minimum-cost bipartite matching. One matching can move many cells at once to a jointly optimal arrangement, which pairwise swaps cannot find.

## Why the cells must be independent

If two chosen cells share a net, the cost of putting one in a slot depends on where the other goes, and the assignment cost no longer decomposes into a sum of independent terms. Choosing cells that share no net, an independent set in the graph where cells are adjacent if they share a net, makes each cell's cost depend only on its own slot. Then the cost of putting cell $i$ in slot $j$ is just the HPWL of $i$'s nets with $i$ at $j$, computed once per pair. Cells are also chosen with the same width (and height), so any assignment of cells to the set's locations is legal. In practice the set is drawn from a local window, so the slots are near each other.

## Solving the matching

With $n$ cells and their $n$ current locations as slots, build an $n \times n$ cost matrix $c_{ij}$ and find the permutation minimising $\sum_i c_{i\pi(i)}$, the linear assignment problem. The Hungarian algorithm solves it in $O(n^3)$; for $n = 3$ there are just six permutations. For three cells A, B and C, the identity assignment might cost 237 um and the best permutation, sending A to B's spot, B to C's and C to A's, 194 um, a gain of 43 um. Because the costs are exact HPWL of independent nets, the gain predicted by the matching is exactly the gain realised.

## ISM in a detailed placer

A detailed placer runs ISM repeatedly: pick a seed cell, grow an independent set of similar cells in a window around it, typically up to a few dozen cells, solve the matching, apply it, and continue. Larger sets find better arrangements but cost more per solve and are harder to keep independent. ISM is used in NTUplace3's detailed placer and in GPU detailed placers such as ABCDPlace, where many independent sets are solved concurrently, since sets that share no nets can be processed in parallel. Combined with global swap and local reordering, it forms the core move set of modern wirelength-driven detailed placement.

## Going deeper

Chen et al., "NTUplace3: an analytical placer for large-scale mixed-size designs" (IEEE TCAD, 2008); Lin et al., "ABCDPlace" (IEEE TCAD, 2020).
