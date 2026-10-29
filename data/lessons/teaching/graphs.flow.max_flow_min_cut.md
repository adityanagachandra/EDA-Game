# Maximum flow and minimum cut

The max-flow min-cut theorem turns a cut problem, which looks combinatorial, into a flow problem with fast exact algorithms. In EDA it computes minimum separating cuts between fixed terminals, refines partition boundaries optimally in tools such as KaHyPar, assigns resources in bipartite settings, and bounds routing capacity across a cut. A designer needs to compute flows and cuts and to read one from the other.

## Flows and their value

A flow network is a directed graph with a source $s$, a sink $t$ and capacities $c(u,v) \ge 0$. A flow assigns $0 \le f(u,v) \le c(u,v)$ to each arc, with inflow equal to outflow at every vertex except $s$ and $t$. Its value is the net flow out of $s$.

The Ford-Fulkerson method increases the flow along **augmenting paths**, $s$-$t$ paths with spare residual capacity, until none remains. The Edmonds-Karp rule always augments along a shortest path (fewest arcs) and runs in $O(VE^2)$; push-relabel and Dinic's algorithm are faster in practice. With arcs $s \to a$ (2), $s \to b$ (5), $s \to c$ (3), $a \to t$ (6), $b \to c$ (5), $c \to t$ (3), at most 2 units can go through $a$ and at most 3 can reach $t$ through $c$, so the maximum flow is 5.

## Cuts and their capacity

An $s$-$t$ cut splits the vertices into $S \ni s$ and $T \ni t$. Its capacity is the sum of capacities of arcs from $S$ to $T$; arcs from $T$ back to $S$ do not count. Every flow must cross every cut, so every cut's capacity is an upper bound on the maximum flow.

For source side $\{s, a, d\}$ in a network with arcs $s \to b$ (6), $s \to c$ (3), $a \to c$ (7), $d \to t$ (8) leaving it, the capacity is 24, even if the maximum flow is only 14. The max-flow min-cut theorem says the smallest such bound is achieved: the maximum flow value equals the minimum cut capacity.

## Reading the minimum cut

When the algorithm stops, the vertices reachable from $s$ in the residual graph form the source side of a minimum cut: every arc leaving that set is saturated, and every arc entering it carries no flow, so its capacity equals the flow value.

Sensitivity follows from the cut. Widening an arc helps only if it crosses every minimum cut, and then by at most the added capacity, until another cut becomes binding. So an arc not on any minimum cut can be widened with no gain at all. In EDA terms: adding tracks across a congested boundary helps only on the boundary that actually limits the flow. Max-flow min-cut is also the basis of flow-based partitioning, where repeated max-flow computations with growing terminal sets find balanced min cuts.

## Going deeper

Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms" (4th ed., 2022), chapter 24; Ford and Fulkerson, "Flows in Networks" (1962).
