# Legalization with Abacus

Abacus improves on Tetris by letting cells already placed in a row shift to make room, choosing positions that minimise total squared displacement in that row. It does this with clusters: groups of abutting cells that move together to their collective optimum. Abacus became a standard legalizer because it is nearly as fast as Tetris and much better in quality.

## Clusters and their optimal position

Within one row, keep the cells in their x order and require no overlap. If cells $1..k$ abut as a cluster starting at $x_c$, cell $i$ sits at $x_c + o_i$, where $o_i$ is the total width of the cells before it in the cluster. Minimising $\sum_i e_i(x_c + o_i - x_i')^2$, with $e_i$ a weight such as cell area or 1, gives
$$x_c = \frac{\sum_i e_i (x_i' - o_i)}{\sum_i e_i},$$
then clamped to the row boundaries. Two cells wanting 23 and 26 with the first 7 wide form a cluster at $(23 + (26 - 7))/2 = 21$: the first cell moves left by 2 and the second right by 2. Abacus stores per cluster $e_c = \sum e_i$, $q_c = \sum e_i(x_i' - o_i)$ and total width $w_c$, so $x_c = q_c/e_c$.

## Adding cells and merging clusters

Cells are processed in x order. A cell that starts at or after the right end of the last cluster begins a new cluster at its desired position. Otherwise it joins the last cluster: $e_c \mathrel{+}= e_i$, $q_c \mathrel{+}= e_i(x_i' - w_c)$, $w_c \mathrel{+}= w_i$. The cluster then moves to $q_c/e_c$, clamped. If it now overlaps its left neighbour, the two merge: the left cluster absorbs the right one with $q \mathrel{+}= q_c - e_c w_{left}$, and the merged cluster is re-collapsed, possibly cascading left. Each merge reduces the number of clusters, so the total work in a row is linear in its cells. In a row with cells c3 and c4 overlapping near x = 23, they cluster at x = 22, and c4 ends at 26.

## Choosing rows and the full algorithm

The full Abacus algorithm processes cells globally in x order. For each cell it tries rows near its desired y, trial-places it at the end of each candidate row with the clustering procedure, computes the cost (the change in squared displacement, including the vertical move), and commits the cell to the cheapest row. Trial placements do not modify the row permanently, and the search stops early when the vertical distance alone exceeds the best cost found. Abacus legalizes in continuous x and then snaps clusters to the site grid. Weighting cells by area keeps large cells from moving much. Its limits are multi-row-height cells and complex rules such as edge spacing, which need extensions; modern legalizers built on it handle those with additional constraints.

## Going deeper

Spindler, Schlichtmann and Johannes, "Abacus: fast legalization of standard cell circuits with minimal movement" (ISPD 2008); Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed., 2022), section 4.4.
