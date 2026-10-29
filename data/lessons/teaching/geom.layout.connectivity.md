# Connectivity extraction with union-find

Layout-versus-schematic and every net-aware check start by deciding which shapes are electrically connected: shapes that overlap or abut on a layer, joined across layers by vias and contacts. The result is a set of nets, the connected components of a graph whose vertices are shapes. Building that graph needs geometric overlap queries; computing components needs a disjoint-set (union-find) structure, which processes the overlaps in near-linear time.

## Components on one layer

On one layer, two shapes are connected when they overlap or share a boundary segment of positive length; touching only at a corner is ambiguous and is usually treated as not connected (and often flagged by a rule). The nets are the connected components of the graph with an edge for every connected pair. Example: shapes $[0, 4] \times [0, 2]$, $[4, 6] \times [0, 5]$, $[10, 12] \times [0, 3]$ and $[6, 8] \times [5, 7]$: the first two share an edge, the second and fourth touch only at the corner $(6, 5)$, and the third is alone, so there are 3 nets. A sweep that reports touching and overlapping pairs produces the edges, and union-find merges them without ever building the graph explicitly.

## Across layers through vias

Connectivity across layers is defined by the process: a via on the via layer connects the metal shapes above and below that it overlaps. The extractor finds, for each via, the M1 shapes and M2 shapes it overlaps (two window queries), and unions them. A via that lands on only one layer connects nothing and is often a violation in itself. Example: three M1 wires and three M2 wires form six shapes; two vias each join one M1 wire to one M2 wire, merging six components into four nets. Stacked vias, contacts to diffusion and poly, and wells extend the same idea through the whole layer stack, defined by a connectivity table in the rule deck. The result labels every shape with a net id.

## Union-find in practice

A disjoint-set forest keeps each component as a tree with a representative root. Union links two roots; find walks to the root, compressing the path on the way. With union by rank or size and path compression, a sequence of $m$ operations on $n$ elements costs $O(m\,\alpha(n))$, where $\alpha$ is the inverse Ackermann function, below 5 for any realistic $n$. A union call merges only when its two shapes are in different sets, so the number of successful merges is $n$ minus the number of components. Example: 8 shapes with union calls (1, 2), (3, 4), (2, 1), (2, 3), (5, 6) merge 4 times, since (2, 1) repeats a connection, leaving 4 nets. Incremental extraction keeps the forest and updates it as edits arrive, though deletions require rebuilding affected components.

## Going deeper

Tarjan, "Efficiency of a good but not linear set union algorithm" (Journal of the ACM, 1975); the KLayout documentation on LVS and net tracing; Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms" (4th ed., 2022), chapter on disjoint sets.
