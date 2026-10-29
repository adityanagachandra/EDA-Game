# Minimum-area retiming

Once a clock period is fixed, many retimings achieve it, and they differ in how many registers they use. Minimum-area retiming picks the one with the fewest. It is a linear program in the lags whose constraints are all differences, which makes it the dual of a minimum-cost network flow problem and solvable exactly in polynomial time.

## The objective in terms of lags

The retimed register count is $\sum_e w_r(e) = \sum_e w(e) + \sum_v r(v)(|FI(v)| - |FO(v)|)$, a constant plus a linear function of the lags. Moving one register from each of a gate's three inputs to its single output ($r(v) = -1$, three fanins, one fanout) changes the count by $-1 \cdot (3 - 1) = -2$. Moving registers backward across a gate with one input and four registered fanouts replaces four registers by one. The objective therefore rewards pushing registers toward nodes with more fanin than fanout, subject to legality and the period constraints.

## A linear program dual to min-cost flow

Minimize $\sum_v r(v)(|FI(v)| - |FO(v)|)$ subject to $r(u) - r(v) \le w(e)$ for each edge and $r(u) - r(v) \le W(u, v) - 1$ for each pair with $D(u, v) > c$. Every constraint has one $+1$ and one -1, so the constraint matrix is totally unimodular and the LP optimum is integral without branch and bound. Its dual is a minimum-cost flow problem, which network-simplex or cost-scaling algorithms solve efficiently. Register sharing at multiple fanouts, which the plain edge sum overcounts, is modeled by adding a mirror vertex per multi-fanout node so that the shared count is the maximum over fanout edges.

## Practical limits

Real flows retime selectively. Moved registers need correct initial (reset) values: forward moves compute them by simulating the logic, but backward moves need a preimage that may not exist. Registers with enables, resets or scan connections constrain where they may move. Verification after retiming needs sequential equivalence checking, since register correspondence is lost. Finally the graph model ignores the library, so the mapped register count can differ from the LP optimum. Min-area retiming is exact for its model, and the model leaves out exactly the details that make production teams cautious.

## Going deeper

Leiserson and Saxe, "Retiming synchronous circuitry", Algorithmica (1991), on minimum-area retiming and mirror vertices; Touati and Brayton, "Computing the initial states of retimed circuits", IEEE TCAD (1993).
