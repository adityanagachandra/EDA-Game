# Coarsening by heavy-edge matching

Multilevel partitioners shrink the netlist level by level, contracting groups of strongly connected cells into single vertices. The simplest and most widely used grouping is a matching: pair up adjacent vertices and contract each pair. Which pairs are chosen decides how much connection weight disappears inside clusters, and therefore how good the coarse problem is.

## Greedy heavy-edge matching

Visit the vertices in some order, often random. When the current vertex is unmatched, match it with the unmatched neighbour joined by the heaviest edge; if no unmatched neighbour remains, leave it single. The result is a matching, a set of disjoint pairs, found in one pass over the edges.

With edges $v_1$-$v_2$ (9), $v_1$-$v_6$ (6), $v_3$-$v_5$ (6), $v_2$-$v_5$ (6), visited in the order $v_6, v_3, v_5, v_1, v_2$: $v_6$ matches $v_1$ (its only neighbour), $v_3$ matches $v_5$, then $v_5$ and $v_1$ are already matched, and $v_2$ finds no unmatched neighbour. The heaviest edge $v_1$-$v_2$ was missed because $v_6$ took $v_1$ first, which shows how the visit order matters.

## Contracting matched pairs

Each matched pair becomes one coarse vertex whose weight is the sum of the two areas. Edges between two coarse vertices merge, with weights added. The coarse graph has $n - |M|$ vertices for a matching of $|M|$ pairs; a matching that pairs most vertices roughly halves the size per level.

With 8 vertices and 2 matched pairs, the coarse graph has 6 vertices, a weak reduction. Real coarseners therefore try to match as many vertices as possible, and hypergraph coarseners often use clustering schemes (each vertex joins its best neighbouring cluster) that shrink faster than strict matchings. Vertex-weight caps keep any cluster from growing so large that balance becomes impossible.

## Why heavy edges

Any edge inside a coarse vertex can never be cut by a partition of the coarse graph. The weight left exposed is the original total minus the matched weight. If the total edge weight is 48 and the matching hides 24, any coarse partition can cut at most 24.

Heavy-edge matching greedily maximizes the hidden weight, which lowers the cut available at the coarse level and keeps strongly connected pairs together, as a good fine partition would. For hypergraphs the analogue is a rating function, such as the sum over shared nets of $w(e)/(|e| - 1)$, which favours pairs sharing many small nets. Good coarsening is what lets the coarse solution be close to a good fine solution, so most quality differences between multilevel tools come from it.

## Going deeper

Karypis and Kumar, "A fast and high quality multilevel scheme for partitioning irregular graphs" (SIAM J. Scientific Computing, 1998); Karypis, Aggarwal, Kumar and Shekhar, "Multilevel hypergraph partitioning: applications in VLSI domain" (IEEE TVLSI, 1999).
