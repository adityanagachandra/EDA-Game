# Horizontal and vertical constraint graphs

Channel routing has two kinds of constraint. Horizontal constraints say two trunks that overlap cannot share a track. Vertical constraints say that at a column with net a on top and net b at the bottom, a's trunk must lie above b's, or their vertical branches would cross on the same layer. The horizontal and vertical constraint graphs encode both, and together they bound and drive every channel router.

## The vertical constraint graph

In a two-layer channel with trunks on one layer and branches on the other, column $c$ with net $a$ at the top and a different net $b$ at the bottom has two vertical branches in the same column. To avoid overlapping, $a$'s trunk must be above $b$'s. The vertical constraint graph (VCG) has a directed edge $a \to b$ for every such column; repeated pairs give one edge.

Example: top pins 2, 3, 3, 0, 4 and bottom pins 0, 1, 4, 2, 1 give edges $3 \to 1$, $3 \to 4$ and $4 \to 1$: three distinct edges (columns with only one pin or with the same net on both sides give none).

## Longest chain as a lower bound

Each VCG edge forces the source net strictly above the target, so a directed path through $k$ nets forces $k$ distinct tracks. If the VCG is acyclic, the number of vertices on its longest path is a lower bound on tracks for any dogleg-free routing, computed by longest-path dynamic programming in topological order.

Example: a chain $3 \to 2 \to 5$ forces 3 tracks, even if the channel density is only 2. The tracks needed are at least $\max(d, \ell)$ for density $d$ and longest chain $\ell$. Neither bound is tight in general, because the two constraint types interact: a net forced low by the VCG may collide horizontally with nets that are also forced low.

## The horizontal constraint graph

The horizontal constraint graph (HCG) is undirected, with an edge between two nets whose spans share a column, so they cannot share a track. It is an interval graph: its largest clique is exactly the channel density, and without vertical constraints it can be colored with that many colors (tracks) by the left-edge algorithm.

Example: spans [4, 9], [3, 8], [4, 8], [2, 3] and [6, 9] give 7 overlapping pairs, and the largest clique is 4 (column 6 to 8). A channel router assigns tracks so that HCG neighbours differ and VCG edges point downward, which is a coloring of the HCG constrained by an order from the VCG.

## Going deeper

A. Hashimoto and J. Stevens, "Wire Routing by Optimizing Channel Assignment within Large Apertures" (1971); Sherwani, "Algorithms for VLSI Physical Design Automation".
