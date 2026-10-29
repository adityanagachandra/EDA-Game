# Layer assignment by dynamic programming

Layer assignment must respect per-layer capacity: after 2D routing, some layers are already full in some gcells. Choosing a layer for each segment so that vias are minimal and no full layer is used is a shortest-path problem on a small layered graph, solved exactly by dynamic programming along each route or tree.

## The recurrence

Number a two-pin route's segments $S_1, \dots, S_n$ from pin to pin, and let $A_i$ be the layers still free for $S_i$ (the right direction, with capacity left in the gcells it crosses). Define $D(i, l)$ as the fewest vias for $S_1$ to $S_i$ with $S_i$ on layer $l$, counting the first pin's stack. Then $D(1, l) = l - 1$ and
$$D(i, l) = \min_{p \in A_{i-1}} D(i-1, p) + |l - p|.$$
The answer closes with the second pin: $\min_{l \in A_n} D(n, l) + l - 1$.

Example: if $S_2$ can only use M5 or M7, the stack from a horizontal $S_1$ to $S_2$ is forced to climb, and the recurrence finds where the climb is cheapest.

## Reading one table entry

The table has one entry per segment and free layer, so a five-segment route with three choices each has at most 15 entries, each computed from at most three predecessors. That is linear in route length, which is why assigners can afford to run it on every net and rerun it after capacity changes.

Example: if $S_1$ is free on M2 and M6 and $S_2$ on M5 and M7, then $D(1, 2) = 1$, $D(1, 6) = 5$, and $D(2, 5) = \min(1 + 3, 5 + 1) = 4$. Remembering which predecessor won at each entry recovers the assignment by walking back from the best final entry. On trees, the same recurrence runs from the leaves to the root, with a junction's cost depending on the layers of all its children.

## Ties and real objectives

Many assignments often achieve the minimum via count. Counting them is a small extension: carry, alongside each entry, the number of predecessors achieving it, and add counts on ties. Example: a route with free layers on its middle segments might have 3 distinct assignments at the minimum of 12 vias.

Ties are an opportunity. A production assigner breaks them, or trades a via or two, for other goals: balancing demand across layers so no layer is near capacity, putting timing-critical nets on low-resistance upper layers, and keeping via density per gcell low. These become extra terms in the same recurrence, or a sequential order in which critical nets are assigned first and claim their layers.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); T.-H. Lee and T.-C. Wang, "Congestion-Constrained Layer Assignment for Via Minimization in Global Routing" (2008).
