# Hypergraph min-cut through a flow network

A partitioner that wants the minimum number of nets separating two groups of cells needs a minimum cut in a hypergraph, where cutting a net costs its weight once however many pins cross. Lawler's construction turns this into an ordinary max-flow problem with exactly the right costs. It is the core of flow-based balanced partitioning and of the flow refinement in KaHyPar.

## Lawler's network

Build a directed network with one vertex per cell. For each net $e$ add two vertices $e_{in}$ and $e_{out}$ and a **bridge** arc $e_{in} \to e_{out}$ with capacity $w(e)$. For every pin of $e$ on cell $v$, add arcs $v \to e_{in}$ and $e_{out} \to v$ of infinite capacity.

For $n$ cells, $m$ nets and $P$ pins, the network has $n + 2m$ vertices and $m + 2P$ arcs. With 7 cells, 5 nets and 15 pins: 17 vertices and 35 arcs. The construction is linear in the netlist size, so a max-flow on it costs about as much as a few FM passes for typical flow algorithms on such sparse graphs.

## Why the minimum cut removes whole nets

A finite cut can never contain an infinite arc, so it consists only of bridges. Cutting bridge $e$ costs $w(e)$ once, which is exactly the hypergraph cost of cutting net $e$. Conversely, any partition of the cells induces a cut of bridges equal to the weight of the nets it cuts. So the minimum $s$-$t$ cut in Lawler's network equals the minimum weight of nets separating cell $s$ from cell $t$.

The graph clique model would charge split nets by pin pairs, the wrong objective. Lawler's network charges exactly once per net. For small netlists the answer can be checked by enumeration: with $s$ and $t$ fixed, try all $2^{n-2}$ assignments of the other cells and take the minimum weighted cut.

## Balanced flow-based partitioning

A minimum $s$-$t$ cut has no balance constraint and often separates almost nothing from $s$. Flow-based balanced bipartitioning (Yang and Wong) repeats the computation: after each max-flow, if the source side is too small, merge a vertex adjacent to the cut into the source, and rerun incrementally from the current flow; similarly on the sink side. The cut grows monotonically until the sides meet the balance window.

Modern partitioners use flows locally. KaHyPar's flow refinement takes a region around the boundary between two blocks, contracts everything outside it into the source and sink, and computes a min cut there, which is optimal within the region. This combines the global optimality of flow with the scalability of multilevel schemes, and it is one reason current tools beat FM-only refinement.

## Going deeper

Lawler, "Cutsets and partitions of hypergraphs" (Networks, 1973); Yang and Wong, "Efficient network flow based min-cut balanced partitioning" (IEEE TCAD, 1996).
