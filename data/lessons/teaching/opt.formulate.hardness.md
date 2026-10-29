# Why exact optimization is hard

Most EDA problems worth solving are NP-hard: partitioning, floorplanning, placement, routing, covering, scheduling under resources. That is not a reason to give up but a statement about which promises a tool can make. Hardness tells the designer to expect heuristics, bounds and decomposition, to measure quality against baselines, and to look for structure in real instances that the worst case ignores.

## How fast enumeration runs out

Exhaustive search is the honest baseline, and arithmetic shows why it fails. Sizing 60 gates with 2 sizes each has $2^{60} \approx 10^{18.1}$ configurations. At $10^{9}$ evaluations per second, enumeration takes about $10^{9}$ seconds, roughly 37 years. With 100 gates it is about $10^{21}$ seconds, far beyond the age of the universe, about $4 \times 10^{17}$ seconds.

The computation is in logarithms: $\log_{10} k^n = n \log_{10} k$, and for orderings $\log_{10} n!$ (Stirling's formula, or a log-gamma function). Subtract the log of the evaluation rate to get the log of the time in seconds. Each extra gate multiplies the time by $k$, so faster hardware buys only a few more variables. The lesson is that any exact method must avoid visiting most of the space, which is what bounds, propagation and dynamic programming do.

## What NP-hardness does and does not say

NP-hardness is a worst-case statement: unless P equals NP, no algorithm solves every instance optimally in polynomial time. It does not say every instance is hard, and it does not forbid exact solvers on small or structured instances. Two-coloring a conflict graph is a linear-time bipartiteness check, while three-coloring, as in triple-patterning decomposition, is NP-complete; tools therefore split the graph into small components and solve those exactly or heuristically.

Integer programs are the common trap. Each LP relaxation in branch and bound may solve in seconds, yet the number of nodes can grow exponentially. Convexity is the property that usually makes continuous problems tractable, since every local optimum is then global. Nonconvex continuous problems can have many local optima, and descent finds only the one whose basin it starts in. Structure, such as planarity or small treewidth, can make restricted versions easy.

## Engineering under hardness

Production engines promise good solutions in bounded run time, usually without an optimality certificate, and are judged against QoR baselines. Where a certificate exists it is valuable: a MILP solver's gap between incumbent and best bound says how far the answer can be from optimal, so a flow can stop at 1% and know the price.

Run-time scaling is measured as carefully as quality. If a partitioner takes 1 s at $10^4$ cells and 100 s at $10^5$, its empirical complexity is about $n^2$, hopeless at $10^7$ cells. Multilevel methods answer by coarsening the netlist, solving the small clustered problem, and refining level by level. Relaxation gives bounds, decomposition gives smaller problems, and approximation algorithms give worst-case ratios: greedy set cover is within about $\ln n$ of optimal. These three moves, plus search, are the toolkit the next concept organizes.

## Going deeper

Garey and Johnson, "Computers and Intractability" (1979), for NP-completeness and the classic hard problems; Karypis and Kumar's papers on multilevel hypergraph partitioning (hMETIS) for the coarsen-solve-refine pattern.
