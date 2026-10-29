# Functional-unit binding by clique partitioning

Functional-unit binding decides which unit executes each operation. Tseng and Siewiorek cast it as clique partitioning of a compatibility graph: operations that could share a unit are adjacent, and each clique becomes one unit. The formulation also covers register binding and interconnect-aware binding, so it is the general framework behind HLS sharing decisions.

## The compatibility graph

Two operations are compatible when one unit can execute both: they are of a type the unit supports and they never need it at the same time. For single-cycle operations of one type, that means they are scheduled in different steps. Draw a vertex per operation and an edge per compatible pair. With seven multiplications in steps 6, 5, 4, 2, 6, 4 and 5, only pairs in the same step are incompatible: (m1, m5), (m2, m7) and (m3, m6). Of the 21 pairs, 18 are compatible, so the graph has 18 edges. A compatibility graph is the complement of a conflict graph; working with compatibility suits sharing problems because a set of operations that can all share one unit is exactly a clique in it.

## Cliques as units

A partition of the vertices into cliques is a binding: every clique is a set of mutually compatible operations, executed on one unit. The fewest cliques is the fewest units. For single-cycle operations of one type, the graph is complete multipartite, with one part per step, and the minimum clique cover equals the size of the largest part, the busiest step: one operation from each step can always be grouped together. With multiplications in steps 1, 4, 5, 1, 2 and 2, the busiest steps hold two, so two multipliers suffice, for example $\{m_1, m_5, m_2, m_3\}$ and $\{m_4, m_6\}$. Clique partitioning is NP-hard on general graphs, and Tseng and Siewiorek's heuristic repeatedly merges the pair of vertices with the most common neighbours, which tends to keep large cliques possible.

## Multicycle operations and weighted edges

With non-pipelined multicycle units, each operation occupies its unit for an interval $[s, s + L - 1]$, and two operations are compatible when their intervals do not overlap. The conflict graph is then an interval graph, and its complement, the compatibility graph, has a minimum clique cover equal to the maximum number of overlapping busy intervals, by perfection of interval graphs. Multiplications of latency 2 starting in steps 4, 4, 4 and 3 overlap in step 4 four at a time, so four multipliers are needed. Real binders go further and weight the edges by how much sharing costs: two operations reading the same registers can share a unit with no new multiplexer inputs, so their edge gets a high weight. Clique partitioning that maximizes total weight then minimizes interconnect as well as unit count.

## Going deeper

Tseng and Siewiorek, "Automated synthesis of data paths in digital systems" (IEEE TCAD, 1986); De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter 6.
