# Min-cut placement

Min-cut placement was the first scalable approach: cut the region in two, partition the netlist into two halves that minimise the nets crossing, assign each half to a side and recurse. It reuses decades of partitioning algorithms (Kernighan-Lin, Fiduccia-Mattheyses, multilevel hMetis) and, in placers such as Capo, stayed competitive well into the 2000s.

## Bisecting the netlist and the region

At each level, a region and its cells are split by a cut line into two subregions, and the cells are partitioned into two sets with areas matching the subregions' capacities, within a balance tolerance. The partitioner minimises the cut, the number (or weight) of nets with cells on both sides, since every cut net must cross the line and will be long. For cells c3 and c5 on the left and c1, c2, c4, c6 on the right, a net {c3, c4, c6} is cut and a net {c1, c2} is not. Fiduccia-Mattheyses moves cells one at a time by gain, in linear time per pass, and multilevel partitioners coarsen the netlist, partition the small version and refine while uncoarsening.

## Levels of recursion

Each level doubles the number of regions. If every bisection splits cells evenly, after $L$ levels each region holds $n/2^L$ cells, so reaching at most $k$ cells per region takes the smallest $L$ with $n/2^L \le k$, that is $L = \lceil \log_2(n/k) \rceil$. A million cells down to regions of 10 cells takes $\lceil \log_2 10^5 \rceil = 17$ levels. Cut directions usually alternate, vertical then horizontal, or follow the region's aspect ratio so that regions stay roughly square. At the bottom, small regions are placed by enumeration or a simple end-case placer. With linear-time partitioning per level, the whole placement is $O(n \log n)$.

## Terminal propagation

Partitioning a region in isolation ignores the cells outside it. If a cell in the region connects to a cell already placed far to the right, it should go to the right half. Terminal propagation, due to Dunlop and Kernighan, fixes this: each external connection is replaced by a fixed terminal on the side of the cut nearer to it, and the partitioner treats it as a cell locked to that side. A net is cut if its members, including propagated terminals, end up on both sides. For a cell c with nets to a right-propagated terminal, to a left terminal and a left cell, and to three right cells, placing c on the right cuts the nets that also touch the left. Without propagation, min-cut placers produce long nets between regions; with it they track the global structure.

## Going deeper

Dunlop and Kernighan, "A procedure for placement of standard-cell VLSI circuits" (IEEE TCAD, 1985); Caldwell, Kahng and Markov, "Can recursive bisection alone produce routable placements?" (DAC 2000), on the Capo placer.
