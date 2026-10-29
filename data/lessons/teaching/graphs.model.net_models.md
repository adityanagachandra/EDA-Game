# Clique and star models of a net

Many algorithms, from spectral partitioning to quadratic placement, need an ordinary weighted graph rather than a hypergraph. The two standard ways to turn a $k$-pin net into edges are the clique model (an edge between every pair of pins) and the star model (a new centre vertex joined to every pin). The choice changes the edge count, the memory footprint and the weights the objective sees, so it is a design decision with consequences.

## Clique and star expansions

The **clique model** replaces a $k$-pin net by the complete graph $K_k$ on its pins, with $\binom{k}{2} = k(k-1)/2$ edges. The **star model** adds one auxiliary vertex for the net and joins it to each pin, giving $k$ edges and one extra vertex. Two-pin nets are a single edge in either model.

Edge counts matter because nets are skewed. A netlist with nets of sizes 3, 2 and 8 has $3 + 1 + 28 = 32$ clique edges, and the one 8-pin net accounts for 28 of them. A 1,000-pin enable net would add almost 500,000 clique edges on its own, against 1,000 star edges. This is why tools use stars for large nets, or drop nets above a size threshold from connectivity-driven steps.

## Weighting clique edges

Expanding every net into a unit-weight clique overweights large nets: a 10-pin net contributes 45 edges and would dominate a 2-pin net's single edge. So clique edges are weighted by net size. Two common choices are $w = 1/(k-1)$, which makes the total weight of a $k$-pin net $k/2$, and $w = 2/k$, which makes it $k-1$, the number of edges of a spanning tree.

When several nets contain the same pair of cells, their edges are merged by adding weights. For cells $c_1, c_4$ on a 5-pin net, a 3-pin net and another 3-pin net, the merged weight under $1/(k-1)$ is $1/4 + 1/2 + 1/2 = 1.25$. This merged weight is what a spectral partitioner or a quadratic placer actually sees between those two cells.

## Choosing a model

The star model is equivalent to the clique for quadratic objectives when the star centre is a free vertex at the optimum: eliminating the centre from a star with edge weights $w$ gives a clique with weights $w/k$ on its pins. That equivalence is why quadratic placers can use stars for large nets without changing the solution, while saving memory and solver time.

For cut-based objectives no graph model is exact. A clique charges a split net more when its pins are split evenly than when one pin is split off, while the hypergraph charges 1 either way. Hybrid schemes are common: cliques for small nets, stars above a threshold of a few pins, and the hybrid Bound2Bound model in analytical placers, which connects each pin to the net's two extreme pins and reproduces half-perimeter wirelength exactly at the current placement.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 4; Spindler, Schlichtmann and Johannes, "Kraftwerk2" (IEEE TCAD, 2008) for the Bound2Bound model.
