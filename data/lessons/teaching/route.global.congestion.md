# Routing congestion and overflow

Congestion is the language of global routing. Every gcell edge has a capacity (usable tracks) and a demand (wires routed across it), and overflow is the excess. Total overflow, maximum overflow and the count of overflowed edges summarise whether a routing can be realised, and the congestion map tells placement and detailed routing where the trouble is.

## Overflowed edges

For an edge $e$ with demand $d_e$ and capacity $c_e$, the overflow is $o_e = \max(0, d_e - c_e)$. An edge is overflowed when $d_e > c_e$: more wires want to cross the boundary than there are tracks. Every overflowed edge is a place where detailed routing will have to detour, or fail.

Example: edges with (demand, capacity) of (5, 7), (5, 4), (16, 14) and (7, 8) have overflows 0, 1, 2 and 0, so two edges are overflowed and the worst overflows by 2. A congestion map colours each edge or gcell by utilisation $d_e / c_e$ or by overflow, which makes hot spots visible to a designer and machine-readable for a placer.

## Total and maximum overflow

Total overflow (TOF) sums the overflow over all edges, $\sum_e \max(0, d_e - c_e)$, and maximum overflow (MOF) is the largest single value. The ISPD 2007 and 2008 global-routing contests ranked routers first by total overflow and then by wirelength, so TOF became the standard figure of merit.

Example: overflows of 2, 2, 0, 0, 0 and 1 give TOF 5 and MOF 2. Under-used edges do not offset over-used ones, because each term is clipped at zero: slack on one edge does not help a net that must cross another. A routing with TOF 0 respects every estimated capacity; one with a small TOF spread thinly may still detail-route, while the same TOF concentrated in one region usually will not.

## Moving demand

Rip-up and reroute changes overflow by moving demand. Taking $k$ wires off edge $a$ and putting them on edge $b$ changes TOF by the change in $\max(0, d_a - c_a)$ plus the change in $\max(0, d_b - c_b)$. It helps only if $b$ has slack: moving 3 wires off an edge overflowed by 3 onto an edge already at capacity simply moves the overflow.

Example: $e_1$ at demand 14, capacity 11 and $e_7$ at 4, capacity 4. Moving 3 wires gives $e_1$ 11 of 11 and $e_7$ 7 of 4, and TOF is unchanged. Global routers therefore price edges by a cost that rises steeply near capacity, so the search naturally prefers edges with slack, and they accept small detours to buy reductions in overflow.

## Going deeper

Kahng, Lienig, Markov and Hu, "VLSI Physical Design" (2nd ed.); the ISPD 2007 and 2008 global routing contest descriptions.
