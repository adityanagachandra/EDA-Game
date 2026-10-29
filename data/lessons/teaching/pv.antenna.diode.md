# Antenna diodes

The standard antenna fix is a reverse-biased diode on the gate's net: a small diffusion region that leaks the plasma charge to the substrate before it can stress the oxide. Diodes are cheap in area but add capacitance and leakage to the net. Routers insert them from a library diode cell, and decks state how much protection a given diode area buys.

## Removing excess metal

Without protection, the metal on the net must satisfy $M \le R_0 G$, where $R_0$ is the limit and $G$ the gate area. A net with 74.6 square um of metal 3 on a 0.15 square um gate, with $R_0 = 250$, may keep 37.5 square um, so 37.1 square um must leave metal 3. In practice the metal is not deleted: the route moves to another layer or is broken by a jumper so that the excess is not connected when metal 3 is etched. The excess quantifies how hard the violation is and which fix can absorb it.

## Sizing the diode protection

Decks state the allowed ratio as a function of the diode (diffusion) area attached to the net, often as a piecewise-linear table. On a linear segment, the limit rises from $R_0$ by $k$ per square um of diode. With a ratio of 742.8 against $R_0 = 300$ and $k = 600$, the net needs 0.738 square um of diode; with diode cells of 0.25 square um, $\lceil 442.8 / 150 \rceil = 3$ cells. Real tables saturate, so very large violations may need a jumper instead. Each diode adds junction capacitance and leakage, a small timing and power cost that the flow must re-extract.

## Diodes on nets with several gates

One diode protects every gate on the net, because they share one node: the charge leaks through the diode regardless of which gate it would have reached. With several gates, the ratio uses their summed area, which already lowers it. Four gates totalling 0.905 square um with 443.5 square um of metal give a ratio of 490.1; with $R_0 = 400$ and 20 per diode cell, 5 cells are needed. Placement matters only electrically, not geometrically, so routers put the diode near a convenient gate or at the driver, and engines verify by recomputing the net's ratio with the diode area included.
