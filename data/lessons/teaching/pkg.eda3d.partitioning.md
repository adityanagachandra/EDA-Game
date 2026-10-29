# Partitioning a design across dies

A 3D or multi-die design must be partitioned: each cell or block is assigned to a die so that the vertical connections stay within budget while area, power and timing balance. This is a hypergraph partitioning problem with die-specific constraints, solved by the same min-cut machinery used in placement, Fiduccia-Mattheyses moves and multilevel coarsening, but with costs that reflect TSVs, bumps or bond pads.

## Cut nets between two dies

Model the netlist as a hypergraph: cells are vertices and each net is a hyperedge joining its cells. Given an assignment of cells to dies, a net is cut when its cells sit on more than one die, and each cut net needs at least one vertical connection. For nets $\{a,b\}$, $\{b,c,d\}$, $\{a,e\}$ and $\{d,e\}$ with $a, b$ on die 0 and $c, d, e$ on die 1, the cut nets are $\{b,c,d\}$, $\{a,e\}$, so the cut is 2 (the net $\{d,e\}$ is internal to die 1). Minimizing the cut subject to area balance is NP-hard, which is why practical tools use multilevel heuristics.

## Vertical connections in a stack

In a stack of more than two dies, a net touching die 0 and die 2 must pass through die 1, so it needs a connection at each interface it crosses. Number the dies by height; a net spanning dies from $d_{min}$ to $d_{max}$ needs $d_{max} - d_{min}$ vertical connections in a simple model with one connection per interface. The total is the sum over nets. This makes the objective sensitive to which dies communicate: putting two heavily connected blocks on non-adjacent dies costs twice. Face-to-face bonding of the two most connected dies, with TSVs only to the others, is a common response.

## Move gain

Iterative improvement moves one cell at a time to the other die and keeps moves that lower the cut. The gain of moving a cell is the cut before the move minus the cut after. Moving cell $e$ from die 1 to die 0 in the example affects only its nets: $\{a,e\}$ becomes internal to die 0, saving one, while $\{d,e\}$ becomes cut, costing one, so the gain is 0. Fiduccia and Mattheyses compute such gains for all cells, keep them in bucket lists for constant-time selection, and lock moved cells to avoid cycling. Area balance limits which moves are legal, and 3D tools add timing and thermal terms to the gain.

## Going deeper

Fiduccia and Mattheyses, "A linear-time heuristic for improving network partitions", DAC (1982); Lim, "Design for High Performance, Low Power, and Reliable 3D Integrated Circuits" (2013).
