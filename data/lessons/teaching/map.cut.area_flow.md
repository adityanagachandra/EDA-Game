# Area flow

Area-oriented mapping must decide which cut each node uses without knowing how many times each node will be shared in the final cover. Area flow is the standard estimate: it charges each cut its own area plus the area flows of its leaves, divided among the node's fanouts. It is cheap to compute in one forward pass and good at steering a mapper toward covers that share logic well.

## The definition

For node n with cut C implemented by a gate or LUT of area A(C), and an estimated fanout F(n) (number of references), the area flow is
$$AF(n) = \frac{A(C) + \sum_{l \in C} AF(l)}{\max(1, F(n))},$$
with AF = 0 at primary inputs. The numerator is the area needed to build n from scratch through this cut; the division spreads it over the consumers that will share it. Example: a LUT of area 1 whose leaves have area flows 1.4, 0 and 0, at a node of fanout 4, gets AF = (1 + 1.4) / 4 = 0.6.

## Choosing the cut with the least flow

In an area pass, each node keeps the cut with the smallest area flow, subject to its required time. Compute AF for every candidate and compare. With leaf flows p = 0.3, r = 2.4, s = 0.8 and area 1 at fanout 2, cut {p, s} gives (1 + 1.1) / 2 = 1.05, while {p, r, s} gives 2.25. Area flow rewards cuts whose leaves are cheap or widely shared, since a shared leaf has its flow divided by its fanout. The pass is one forward sweep, linear in the number of node-cut pairs.

## Propagation and its blind spot

Area flows propagate: AF(x) feeds AF(y) when x is a leaf of y's cut, so a chain x, y, n with fanouts 1 each simply accumulates areas, while a leaf with fanout 3 contributes a third of its flow to each user. The estimate uses fanouts of the subject graph or of the previous mapping, not of the final cover, so it can mis-charge nodes whose fanout changes when consumers absorb them. That is why mappers follow area flow with exact local area, which measures the real area added by a cut against the current cover's reference counts.

## Going deeper

Manohararajah, Brown and Vranesic, "Heuristics for area minimization in LUT-based FPGA technology mapping", IEEE TCAD 2006; Mishchenko, Chatterjee and Brayton, "Improvements to technology mapping for LUT-based FPGAs", IEEE TCAD 2007.
