# Clock tree topology generation

Zero-skew merging needs a topology: which sinks and subtrees merge with which, in what order. The embedding algorithm (DME) is optimal for a given topology, but the best topology is a hard combinatorial problem, so practical tools use geometric heuristics. The two classic ones are top-down partitioning by the method of means and medians, and bottom-up greedy matching. Both aim to pair nearby sinks of similar delay, which keeps merges short and avoids detour wire.

## Greedy nearest-pair merging

The simplest bottom-up rule merges the closest pair of current subtrees first, measured by Manhattan distance, replaces them by their merge node, and repeats. For sinks at $(60, 34)$, $(11, 41)$, $(11, 15)$ and $(14, 31)$, the closest pair is $(11, 41)$ and $(14, 31)$ at distance $3 + 10 = 13$ µm, so they merge first. With a nearest-neighbour structure each step is fast, and the result follows the placement closely. Its weakness is myopia: an early short merge can leave two far-apart leftovers that must be joined by a long, expensive edge.

## Means and medians

The method of means and medians (MMM) works top-down. It splits the sink set at the median x coordinate into two equal halves, then splits each half at its median y, alternating directions, until single sinks remain; the root connects the centres of mass of the two halves. For x coordinates 2, 6, 15, 56, 58 and 60, the median is $(15 + 56)/2 = 35.5$, so three sinks go to each side. Equal halves give a balanced topology whose depth is $\lceil \log_2 n \rceil$, which keeps latency and the number of merges per path uniform.

## Matching rounds

Recursive geometric matching builds one level at a time. In each round it pairs all current subtrees, for instance greedily by repeatedly taking the closest unmatched pair, then merges every pair at once. For sinks at $(18, 22)$, $(49, 37)$, $(6, 41)$ and $(15, 47)$, the closest pair is $(6, 41)$ and $(15, 47)$ at 15 µm, which forces the remaining pair at 46 µm, a total of 61 µm. A round halves the number of subtrees, so the topology is balanced. Minimum-cost perfect matching would do better than greedy per round, at higher cost. Either way the topology then goes to Tsay or DME for embedding.

## Going deeper

Jackson, Srinivasan and Kuh, "Clock routing for high-performance ICs", Design Automation Conference (1990); Kahng, Cong and Robins, "High-performance clock routing based on recursive geometric matching", Design Automation Conference (1991).
