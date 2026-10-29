# Minimum spanning trees for net estimation

Before a net is routed, placers and timers need an estimate of its wirelength and topology. The rectilinear minimum spanning tree over the pins, with Manhattan distances, is the classic estimate: exact algorithms are fast, it is never more than 50% longer than the best Steiner tree, and it gives a tree topology for delay estimation. It is also the starting point of most Steiner heuristics.

## Rectilinear distance and the MST

Pins sit at integer coordinates, and wires run horizontally and vertically, so the distance between pins $p$ and $q$ is $|x_p - x_q| + |y_p - y_q|$. The rectilinear MST is a minimum spanning tree of the complete graph on the pins with these distances. Each tree edge is then realized as an L-shaped wire.

For pins $(5, 3)$, $(8, 8)$, $(11, 1)$, $(7, 0)$, the shortest pairwise distances are 5 from $(5,3)$ to $(7,0)$, 5 from $(7,0)$ to $(11,1)$, and 8 from $(5,3)$ to $(8,8)$, giving a tree of length 18. The half-perimeter of the bounding box is 14, a lower bound that no tree can beat.

## Prim's algorithm on pins

Prim's algorithm grows a tree from one pin, at each step adding the shortest edge from the tree to a pin not yet in it. Keeping, for each outside pin, its distance to the nearest tree pin makes each step $O(n)$ and the whole run $O(n^2)$, which is optimal for the complete graph. Kruskal's algorithm, sorting all $\binom{n}{2}$ edges, gives the same tree when lengths are distinct.

The order of edges added is not sorted: from pin 1, the steps might add lengths 4, 6, 7, 10, 3, the last being short because its pin was far from the early tree but close to a late one. For nets with many pins, rectilinear MSTs can be built in $O(n \log n)$ using the fact that each pin needs candidate edges only to its nearest neighbour in each of eight octants.

## What the MST estimate is good for

Placement objectives usually use half-perimeter wirelength because it is cheap and smooth enough, but HPWL underestimates nets with four or more pins. The MST is a better length estimate for multi-pin nets and gives an actual topology, which a timer needs for Elmore delay estimates before routing.

The MST is never shorter than the rectilinear Steiner minimum tree, and Hwang's theorem bounds it at $3/2$ of the RSMT. On typical nets the gap is much smaller. Global routers often start from a Steiner tree (from FLUTE in OpenROAD) rather than an MST, but MSTs remain the basis of fast estimators, of Steiner heuristics that improve an MST by adding points, and of clock and buffer tree construction.

## Going deeper

Prim, "Shortest connection networks and some generalizations" (Bell System Technical Journal, 1957); Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 5.
