# Graph coloring in EDA

Graph coloring assigns labels to vertices so that adjacent vertices differ. In EDA the labels are resources: registers for values in high-level synthesis, tracks for nets in channel routing, masks for features in multiple patterning, time slots for operations, and test sessions for cores. Coloring is NP-hard in general, so tools rely on greedy orderings, structure that makes special cases easy, and exact search only for small instances.

## Conflict graphs and greedy coloring

Build a conflict graph: one vertex per item, an edge between two items that cannot share a resource. In register allocation, an edge joins two values that are live at the same time. A proper coloring with $k$ colors is an assignment of $k$ registers with no conflict.

Greedy coloring visits vertices in a chosen order and gives each the smallest color not used by an already-colored neighbour. In the order $v_4, v_5, v_3, v_6, v_2, v_1, v_7$ with edges $v_1$-$v_2$, $v_1$-$v_4$, $v_2$-$v_4$, $v_2$-$v_6$, $v_4$-$v_6$, $v_3$-$v_4$, $v_3$-$v_7$, $v_5$-$v_7$: $v_4$ gets 1, $v_5$ gets 1, $v_3$ gets 2, $v_6$ gets 2, $v_2$ (adjacent to $v_4$ and $v_6$) gets 3.

## Counting colors and the effect of order

The number of colors greedy uses depends on the order. It never exceeds $\Delta + 1$, one more than the maximum degree, because a vertex has at most $\Delta$ colored neighbours. A bad order can use many more colors than necessary on some graphs; a good one, such as smallest-last ordering (repeatedly remove a minimum-degree vertex, then color in reverse), often comes close to optimal.

Chaitin's register allocator uses this idea: remove vertices of degree below $k$ (they can always be colored later), and when none remain, choose a value to spill to memory. The spill decision is where cost models enter: spilling a value used inside a hot loop is far more expensive than spilling one used once.

## Chromatic number and hardness

The **chromatic number** $\chi(G)$ is the fewest colors in any proper coloring. Deciding whether $\chi(G) \le 3$ is NP-complete, and $\chi$ is hard even to approximate in general graphs. For small conflict graphs exact search by backtracking is practical: try colors vertex by vertex and backtrack when a vertex has none left.

A greedy pass may use 4 colors on a graph whose chromatic number is 3, which exhaustive search confirms by finding a 3-coloring and proving no 2-coloring exists (an odd cycle or a triangle is enough proof). Two lower bounds are cheap: the size of the largest clique, and $n / \alpha(G)$, the vertex count over the largest independent set. EDA tools exploit structure where these bounds are tight, as in interval graphs.

## Going deeper

Chaitin, "Register allocation and spilling via graph coloring" (Proc. SIGPLAN Symposium on Compiler Construction, 1982); De Micheli, "Synthesis and Optimization of Digital Circuits" (1994), chapter on resource sharing.
