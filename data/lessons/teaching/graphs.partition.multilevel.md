# Multilevel partitioning

Multilevel partitioning is the method behind hMETIS, KaHyPar and PaToH, and the reason partitioning million-cell netlists is fast and good. It coarsens the netlist into a hierarchy of smaller ones, partitions the smallest, then projects the solution back level by level while refining it with FM. Each level contributes moves at a different scale, which flat local search cannot make.

## Coarsen, partition, uncoarsen

In the **coarsening** phase, the netlist is repeatedly contracted by matching or clustering, giving a sequence $H_0, H_1, \dots, H_L$ of ever smaller hypergraphs. In the **initial partitioning** phase, the coarsest $H_L$, often a few hundred vertices, is partitioned by any method, typically many FM runs from random or greedy starts. In the **uncoarsening** phase, the partition is projected to $H_{L-1}$ (each fine vertex takes its cluster's side) and refined with FM, then projected again, down to $H_0$.

Projection preserves the cut exactly, because merged nets carry summed weights. Refinement then only improves it, so each level hands down a solution at least as good as the coarser one.

## Depth of the hierarchy

If each level shrinks the vertex count by a factor $r$, reaching a stop size $T$ from $N$ vertices takes the smallest $L$ with $N / r^L \le T$, that is $L = \lceil \log(N/T) / \log r \rceil$. From $10^6$ cells with $r = 1.6$ and $T = 500$, $L = 17$ levels.

The total work is a geometric series: level $i$ costs time proportional to its size, $N / r^i$, so all levels together cost about $N \cdot r / (r - 1)$, a small constant times one finest-level pass. That is why multilevel partitioning stays near-linear. Coarsening too slowly (r close to 1) multiplies levels and work; too aggressively makes clusters that hide structure the partitioner needed to see.

## Why it works, and variants

Moving one coarse vertex relocates a whole cluster, a move flat FM could make only through a long chain of negative-gain steps. Refinement at each finer level then repairs the boundary at a finer granularity. The result is consistently better cuts than flat FM in comparable time.

Variants refine the scheme. **V-cycles** re-coarsen while respecting the current partition and refine again. **n-level** coarsening (KaHyPar) contracts one pair per level for very gradual refinement. **Direct k-way** refinement avoids the weakness of recursive bisection, whose early cuts ignore later splits. Flow-based refinement, as in KaHyPar's max-flow min-cut step, improves the boundary between two blocks optimally within a region.

## Going deeper

Karypis, Aggarwal, Kumar and Shekhar, "Multilevel hypergraph partitioning: applications in VLSI domain" (IEEE TVLSI, 1999); Schlag et al., "High-quality hypergraph partitioning" (ACM Journal of Experimental Algorithmics).
