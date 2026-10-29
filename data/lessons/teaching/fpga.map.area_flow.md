# Area recovery with area flow

After depth-optimal labeling, a LUT mapper recovers area by choosing, at every node with slack, the cut that leads to the fewest LUTs. Exact LUT counting is a global covering problem, so mappers use area flow, a heuristic that estimates each node's share of the LUTs in its cone. It is the main area heuristic in ABC's LUT mapper and in DAOmap-style mappers.

## Area flow of a node

With unit LUT area, the area flow of node $v$ under its chosen cut $C$ is $AF(v) = \left(1 + \sum_{u \in C} AF(u)\right) / |\mathrm{fanout}(v)|$, with $AF = 0$ at primary inputs. The numerator is the LUT at $v$ plus the area of the cones feeding it; dividing by the fanout spreads that area over the nodes that will reuse $v$'s output. If $v$ has fanout 2 and its cut's leaves have area flows 0.25, 2.75, 0.75 and 1.75, then $AF(v) = (1 + 5.5)/2 = 3.25$. Fanout here is an estimate (from the subject graph or the previous mapping pass), which makes area flow cheap to compute in one topological sweep.

## Choosing cuts by area flow

During area recovery the mapper visits nodes in topological order and, among the cuts whose arrival time still meets the node's required time, keeps the one with the smallest area flow. Because every cut of a node divides by the same fanout, the comparison reduces to the sum of leaf area flows. With cut 1 leaves at 1.5 and 1.5 and cut 2 leaves at 1, 2, 0.5 and 0.5, fanout 2: cut 1 gives $(1 + 3)/2 = 2$ and cut 2 gives $(1 + 4)/2 = 2.5$, so cut 1 wins. Required times come from the depth-optimal labels on critical paths and relax off them, which is how the mapper trades slack for area without increasing the overall depth.

## Why area flow is only a heuristic

Area flow propagates through the network, so it accounts for sharing that plain cone size misses: a node with high fanout looks cheap because its LUT will be used many times. But it can double count. If $n_1$ feeds both $n_2$ and $n_3$, and $n_3$'s cut includes both $n_2$ and $n_1$, then $n_1$'s area flow reaches $n_3$ twice: $AF(n_1) = 0.5$, $AF(n_2) = (1 + 0.5 + 2)/2 = 1.75$, $AF(n_3) = (1 + 1.75 + 0.5)/1 = 3.25$. Mappers therefore follow area flow with exact local area, which counts the LUTs a cut would actually add given the current mapping by referencing and dereferencing the cones. The pair, a global estimate then a local exact pass, is the standard recipe.

## Going deeper

Manohararajah, Brown and Vranesic, "Heuristics for Area Minimization in LUT-Based FPGA Technology Mapping", IEEE TCAD (2006); Chen and Cong, "DAOmap: A Depth-Optimal Area Optimization Mapping Algorithm for FPGA Designs", ICCAD (2004).
