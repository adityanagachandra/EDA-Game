# Scan-chain ordering as a traveling salesman path

Scan design threads every flip-flop of a block into shift-register chains for test. The logical order of a chain is free, so physical design reorders it after placement to minimize the wire that connects consecutive flops. That is a traveling-salesman path problem on the flop locations, NP-hard in general, and handled with construction heuristics followed by local improvement such as 2-opt.

## Chain length as a path length

A scan chain starts at a scan-in port, visits each flop once, and ends at the last flop (or a scan-out port). Its wirelength, estimated with Manhattan distances, is the sum of the hops: $d(\text{port}, f_{\pi(1)}) + \sum_k d(f_{\pi(k)}, f_{\pi(k+1)})$ for the order $\pi$.

For flops at $(3, 7)$, $(9, 10)$, $(5, 10)$, $(1, 3)$, $(7, 1)$ and the order 4, 3, 2, 1, 5 from a port at $(0, 0)$, the hops are 4, 11, 4, 9, 10, for a total of 38. Minimizing this over all orders is a Hamiltonian path problem with a fixed start, a variant of the traveling salesman problem.

## Nearest-neighbour construction

The nearest-neighbour heuristic starts at the port and always hops to the closest unvisited flop. It runs in $O(n^2)$, or near-linear with spatial indexing, and gives reasonable chains when flops are clustered. Its failure mode is stranding: it leaves distant flops until the end and pays a long final hop.

Nearest neighbour has no constant worst-case guarantee for TSP. For a closed tour with metric distances, Christofides' algorithm guarantees at most $3/2$ of optimal, and doubling a minimum spanning tree guarantees 2. Practical scan stitching uses fast construction followed by local improvement, and respects constraints such as clock domains, chain-length balance and lockup latches, which partition the flops into separate ordering problems.

## 2-opt improvement

A 2-opt move reverses a contiguous segment of the order. For a path, only the hop into the segment and the hop out of it change, since the hops inside a reversed segment keep their lengths. So the change is computed in constant time: old length of those two hops minus new length.

Reversing positions 4 and 5 of the order f1, f4, f6, f5, f2, f3 replaces hops f6 to f5 and f2 to f3 with f6 to f2 and f5 to f3. If the sum is unchanged the move is neutral; if it shrinks, keep it. Repeating improving 2-opt moves until none remains gives a 2-opt local optimum, typically within a few percent of optimal for Euclidean or Manhattan instances. Or-opt and Lin-Kernighan moves extend the idea to larger neighbourhoods.

## Going deeper

Applegate, Bixby, Chvátal and Cook, "The Traveling Salesman Problem: A Computational Study" (2006); Bushnell and Agrawal, "Essentials of Electronic Testing" (2000), chapter on scan design.
