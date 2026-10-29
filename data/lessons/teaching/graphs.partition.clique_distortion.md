# How net models distort cut size

Graph partitioners and spectral methods work on ordinary graphs, so they replace each net by a clique or star. But the hypergraph cut charges a split net exactly once, while a clique charges it according to how its pins are split. Understanding that distortion explains why hypergraph partitioners beat graph partitioners on netlists, and how to weight cliques so the damage is small.

## Counting crossing clique edges

Model a $k$-pin net as a clique. If a partition puts $m$ of its pins on side $A$ and $k - m$ on side $B$, every pin on $A$ is joined to every pin on $B$, so $m(k - m)$ clique edges cross. The hyperedge itself is simply cut, cost 1.

For an 8-pin net split 4 and 4, 16 clique edges cross; split 1 and 7, only 7 cross. So a graph partitioner with unit clique edges would much rather split off one pin than split the net in half, even though both cost the same one cut net. The crossing count is maximized by the most even split, $\lfloor k/2 \rfloor \lceil k/2 \rceil$.

## Weighting to reduce the distortion

Weighting clique edges by $1/(k-1)$ makes the cost of splitting off a single pin exactly 1: $1 \cdot (k-1) / (k-1) = 1$. But an even split still costs more. A 4-pin net split 2 and 2 costs $4/3 \approx 1.33$; a 16-pin net split 8 and 8 costs $64/15 \approx 4.27$. The general cost $m(k-m)/(k-1)$ ranges from 1 up to about $k/4$.

No choice of edge weight can fix this, because the hyperedge cost is constant across splits while the clique cost depends on $m$. Weighting only rescales the curve. So graph models are always biased toward pulling large nets apart unevenly, and toward keeping them whole less strongly than the true objective does.

## Graph cut versus hypergraph cut

Over a whole netlist, the weighted graph cut is $\sum_e m_e(k_e - m_e)/(k_e - 1)$, which is at least the number of cut nets, with equality when every cut net has a single pin on one side. For a partition with 5 cut nets the graph cut might be 5.33, and a different partition with the same 5 cut nets could score higher or lower, so the two objectives can rank partitions differently.

This is the practical reason hMETIS, KaHyPar and PaToH keep nets as hyperedges: on circuit benchmarks they find substantially smaller cuts than graph partitioners applied to clique expansions. Graph models remain useful where the objective is itself a graph quantity, such as quadratic wirelength in placement, and in spectral methods, which need a Laplacian.

## Going deeper

Ihler, Wagner and Wagner, "Modeling hypergraphs by graphs with the same mincut properties" (Information Processing Letters, 1993); Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 2.
