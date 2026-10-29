# Breadth-first wavefront search on a grid

Lee's maze-routing algorithm, from 1961, is breadth-first search on a grid graph: a wavefront expands from the source one step at a time, labelling each cell with its distance, until it reaches the target. It always finds a shortest path if one exists. Every later router, from line-probe methods to A*-guided detailed routers, is measured against it, and its wavefront picture is the clearest way to see what a search spends its time on.

## Wavefront labels

Label the source 0. At step $k$, every unlabelled, unblocked neighbour of a cell labelled $k$ receives label $k + 1$. With unit edge costs this is breadth-first search, and a cell's label is its shortest-path distance from the source around the blockages.

On an open grid a cell's label equals its Manhattan distance from the source. Blockages bend the wavefront: a wall of blocked cells with a single gap forces every label beyond the wall to be counted through the gap. A cell at Manhattan distance 3 behind such a wall might receive label 5 or more. Each cell is labelled once and examined once, so the search costs time linear in the number of cells reached.

## Reaching the target and retracing

The expansion stops when the target receives a label; that label is the length of a shortest route. The route is recovered by **retrace**: from the target, step to any neighbour whose label is one less, until the source is reached. When several neighbours qualify, routers prefer the one that continues in the current direction, which minimizes bends and therefore vias.

If the wavefront dies out before reaching the target, no route exists with the current blockages, and the router must rip up other nets to make room. The difference between the label at the target and the Manhattan distance is the detour the blockages force, a quantity routers watch because long detours waste tracks and add delay.

## The cost of the wavefront

The number of cells at label exactly $k$ is the size of wavefront $k$. On an open grid it grows linearly with $k$ (up to $4k$ cells), so reaching a target at distance $D$ explores about $2D^2$ cells: quadratic in the route length, in every direction, including directions away from the target.

That waste is why practical routers bound the search to a box around the source and target, use bidirectional search, or guide it with A*. Lee's algorithm also handles only unit costs; non-uniform costs (vias, congestion, preferred direction) need Dijkstra's algorithm, of which Lee's search is the unit-cost special case. Its memory cost, one label per cell, motivated the compact coding schemes of early routers.

## Going deeper

Lee, "An algorithm for path connections and its applications" (IRE Trans. Electronic Computers, 1961); Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 5.
