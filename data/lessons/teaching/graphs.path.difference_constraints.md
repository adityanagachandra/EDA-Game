# Difference constraints as shortest paths

Many layout and timing problems reduce to constraints of the form $x_j - x_i \le c$: minimum spacing between cell edges in compaction, register moves in retiming, and offsets in clock scheduling. A system of such constraints is a graph, and it is solved by a single shortest-path computation. This reduction is one of the most reused tricks in EDA, and it explains why compactors and retimers are fast.

## From constraints to a graph

Write each constraint $x_j - x_i \le c$ as an arc $x_i \to x_j$ of weight $c$. Add a source $x_0$ with a weight-0 arc to every variable. If the graph has no negative cycle, setting $x_k = d(x_0, x_k)$, the shortest distance from the source, satisfies every constraint, because shortest distances obey the triangle inequality $d(j) \le d(i) + c$, which is exactly $x_j - x_i \le c$.

For $x_1 - x_2 \le 3$, $x_3 - x_1 \le -2$, $x_2 - x_3 \le 6$, the arcs are $x_2 \to x_1$ (3), $x_1 \to x_3$ (-2), $x_3 \to x_2$ (6). From the source every distance is at most 0, and $x_3 = -2$. Adding a constant to every variable gives another solution, since only differences are constrained.

## Compaction by longest paths

Layout compaction uses the mirror form. A minimum-spacing rule says cell $b$'s left edge is at least cell $a$'s left edge plus $a$'s width plus a gap: $x_b \ge x_a + w_a + g$. These are arcs $a \to b$ with weight $w_a + g$, and the smallest legal positions are **longest**-path distances from the left boundary. Because left-to-right constraints form a DAG, one topological pass solves it.

Cells of widths 1, 3, 5 with constraints cell 1 before cell 2 (gap 2), cell 1 before cell 3 (gap 1), cell 2 before cell 3 (gap 0) give $x_1 = 0$, $x_2 = 3$, $x_3 = \max(0 + 1 + 1, 3 + 3 + 0) = 6$, and the row width is $6 + 5 = 11$. The longest path through the constraint graph is the critical chain that sets the width.

## Bounds and feasibility

The shortest-path distance from $x_i$ to $x_j$ is the tightest upper bound the system places on $x_j - x_i$: chaining constraints along any path from $x_i$ to $x_j$ bounds the difference by the path's weight, and the shortest-path potentials attain it. So asking how far apart two variables can be is a single-source shortest-path query.

The system is infeasible exactly when the graph has a negative cycle, which Bellman-Ford detects. In compaction that means the spacing rules cannot all be met at once, typically because of a user constraint that fixes a cell too close to another. Retiming (Leiserson and Saxe) uses the same machinery: register counts on arcs must stay non-negative, which gives difference constraints on the retiming labels, and a clock-period target adds more.

## Going deeper

Cormen, Leiserson, Rivest and Stein, "Introduction to Algorithms" (4th ed., 2022), section 22.4; Leiserson and Saxe, "Retiming synchronous circuitry" (Algorithmica, 1991).
