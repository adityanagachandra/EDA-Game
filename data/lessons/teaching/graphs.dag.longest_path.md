# Longest paths in a DAG

Static timing analysis is, at its core, a longest-path computation on a DAG: the latest arrival time at a vertex is the longest path, by delay, from any timing start point. Longest paths are NP-hard in general graphs but linear on DAGs, which is the single most important fact behind block-based timing. The same recurrence gives logic depth, schedule length in high-level synthesis and positions in layout compaction.

## The max-plus recurrence

Give each arc $u \to v$ a delay $d(u,v)$ and each source an arrival time $a(s)$. For every other vertex,

$$a(v) = \max_{u \to v} \big(a(u) + d(u,v)\big).$$

The recurrence is well defined only because the graph is acyclic: process vertices in topological order and every $a(u)$ is final before it is used. Arcs $v_1 \to v_3$ (4), $v_3 \to v_4$ (3), $v_2 \to v_4$ (9), $v_3 \to v_5$ (5), $v_4 \to v_5$ (9) with $a(v_1) = a(v_2) = 0$ give $a(v_3) = 4$, $a(v_4) = \max(4+3, 0+9) = 9$, $a(v_5) = \max(4+5, 9+9) = 18$.

## Why acyclicity makes it easy

On a general graph, the longest simple path is NP-hard: it contains the Hamiltonian path problem. On a DAG every path is simple, and the recurrence above solves the problem in $O(V + E)$ by dynamic programming. Equivalently, negate every delay and run a shortest-path computation; on a DAG no negative cycle can exist.

Nonzero input arrivals simply initialize the sources. If inputs arrive at $v_1$ at 5 ps and $v_2$ at 0 ps, every path from $v_1$ gains 5 ps, and the maximum at each vertex is taken over these offset paths. Real timers carry separate rise and fall arrivals and propagate slews alongside, but the structure stays the same max-plus pass in topological order.

## Paths versus blocks

The recurrence keeps one number per vertex, the latest arrival, rather than every path. That is **block-based** analysis: linear time, but it loses which path produced the maximum unless a back-pointer is stored. **Path-based** analysis enumerates paths, which can be exponentially many, but can apply path-specific effects such as removing common clock pessimism or computing exact path delay with a better model.

Production timers do both: a block-based pass finds arrivals and slacks everywhere, then path enumeration extracts the worst few paths per endpoint for reporting and refinement. The max operator is also the reason timing is pessimistic under variation: the maximum of random arrivals is not Gaussian, which statistical timers must approximate.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design: From Graph Partitioning to Timing Closure" (2nd ed., 2022), chapter 8; Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms" (4th ed., 2022), section 22.2.
