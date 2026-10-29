# The cut size of a partition

Partitioning splits a netlist into blocks so that few nets cross between them. It drives floorplanning, min-cut placement, multi-FPGA mapping, parallel simulation and the clustering steps inside modern placers. The objective is the cut: how many nets, or how much net weight, spans more than one block. Computing it exactly, and knowing which of its variants a tool optimizes, is the starting point for every partitioning algorithm.

## Cut nets in a bipartition

A bipartition assigns each cell to side $A$ or side $B$. A net is **cut** if it has pins on both sides. The cut size is the number of cut nets: a net counts once however many of its pins cross. For nets $n_1 = \{c_1, c_6\}$, $n_2 = \{c_2, c_3, c_5\}$, $n_3 = \{c_1, c_5\}$ and the partition $A = \{c_2, c_5\}$, $B = \{c_1, c_3, c_6\}$, net $n_1$ lies entirely in $B$, while $n_2$ and $n_3$ are cut, giving a cut size of 2.

Checking a net needs only the side of its first pin and whether any other pin differs, so the whole cut costs one pass over the pins, $O(P)$. Because each cut net becomes at least one wire crossing the block boundary, the cut is a proxy for inter-block wiring and for pins needed on each block.

## Weighted cuts

Nets carry weights when they are not equally important: a bus bit, a timing-critical net, or a coarse net standing for several merged nets. The weighted cut is $\sum_{e \text{ cut}} w(e)$. With weights 2, 4 and 3 on three cut nets and an uncut net of weight 5, the weighted cut is $2 + 4 + 3 = 9$.

Weights are also how multilevel partitioners stay exact: when coarsening merges two nets with identical pin sets, it keeps one net with the summed weight, so the coarse cut equals the fine cut of any projected partition. Timing-driven partitioning raises weights on critical nets so the optimizer avoids cutting them, at the price of cutting more of the others.

## Cut metrics for k-way partitions

With $k > 2$ blocks a net can touch several blocks, and two metrics differ. The **cut-net** metric counts each net touching more than one block once. The **connectivity** metric, $\sum_e (\lambda_e - 1)$ where $\lambda_e$ is the number of blocks net $e$ touches, charges a net spanning three blocks twice.

Connectivity models reality better: a net spread over three blocks needs wires into at least two other blocks, and in parallel computing it measures communication volume. Over the same 3-way partition, a design might have 6 cut nets but a connectivity cost of 8 because two nets touch all three blocks. Tools such as KaHyPar and PaToH optimize connectivity; hMETIS supports both, and the sum-of-external-degrees metric is a third variant. A designer must check which one a tool reports before comparing numbers.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 2; Schlag et al., "High-quality hypergraph partitioning" (ACM Journal of Experimental Algorithmics), for the connectivity metric in KaHyPar.
